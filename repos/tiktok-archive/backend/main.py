from fastapi import FastAPI, Depends, HTTPException, BackgroundTasks, Query
from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import StaticFiles
from fastapi.responses import FileResponse
from sqlalchemy.orm import Session
from sqlalchemy import or_, func
from pydantic import BaseModel
from typing import List, Optional
from datetime import datetime
import os
import json

from database import init_db, get_db, Video, Collection, CollectionVideo, DownloadQueue
from downloader import download_video, parse_metadata, get_video_id, ARCHIVE_DIR, get_user_videos, save_cookies, has_cookies

app = FastAPI(title="TikTok Archive", version="1.0.0")

# CORS for frontend
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000", "http://localhost:5173"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Mount static files for archived content
os.makedirs(ARCHIVE_DIR, exist_ok=True)
app.mount("/archive", StaticFiles(directory=ARCHIVE_DIR), name="archive")


# Pydantic models
class DownloadRequest(BaseModel):
    url: str
    collection_id: Optional[int] = None


class BulkDownloadRequest(BaseModel):
    urls: List[str]
    collection_id: Optional[int] = None


class CollectionCreate(BaseModel):
    name: str
    description: Optional[str] = None


class UserDownloadRequest(BaseModel):
    username: str
    limit: Optional[int] = None
    collection_id: Optional[int] = None


class CookiesRequest(BaseModel):
    cookies: str


class VideoResponse(BaseModel):
    id: int
    tiktok_id: str
    url: str
    title: Optional[str]
    description: Optional[str]
    uploader: Optional[str]
    uploader_id: Optional[str]
    upload_date: Optional[datetime]
    duration: Optional[float]
    view_count: Optional[int]
    like_count: Optional[int]
    comment_count: Optional[int]
    video_path: Optional[str]
    thumbnail_path: Optional[str]
    audio_path: Optional[str]
    tags: Optional[List[str]]
    music_title: Optional[str]
    music_author: Optional[str]
    archived_at: datetime
    file_size: Optional[int]

    class Config:
        from_attributes = True


# Initialize database on startup
@app.on_event("startup")
async def startup():
    init_db()


# Background task for downloading
def download_task(url: str, queue_id: int, collection_id: Optional[int], db_session_maker):
    db = db_session_maker()
    try:
        # Update queue status
        queue_item = db.query(DownloadQueue).filter(DownloadQueue.id == queue_id).first()
        if queue_item:
            queue_item.status = "downloading"
            queue_item.started_at = datetime.utcnow()
            db.commit()

        # Download video
        video_id = get_video_id(url)
        result = download_video(url, video_id)

        if result["success"]:
            # Parse metadata
            metadata = parse_metadata(result["metadata"]) if result["metadata"] else {}

            # Check if video already exists
            existing = db.query(Video).filter(Video.tiktok_id == video_id).first()

            if existing:
                # Update existing record
                for key, value in metadata.items():
                    if value is not None:
                        setattr(existing, key, value)
                existing.video_path = result["video_path"]
                existing.thumbnail_path = result["thumbnail_path"]
                existing.audio_path = result["audio_path"]
                existing.updated_at = datetime.utcnow()
                video = existing
            else:
                # Create new record
                video = Video(
                    tiktok_id=video_id,
                    url=url,
                    video_path=result["video_path"],
                    thumbnail_path=result["thumbnail_path"],
                    audio_path=result["audio_path"],
                    **metadata
                )

                # Build search text
                search_parts = [
                    metadata.get("title", ""),
                    metadata.get("description", ""),
                    metadata.get("uploader", ""),
                    metadata.get("music_title", ""),
                ]
                if metadata.get("tags"):
                    try:
                        tags = json.loads(metadata["tags"])
                        search_parts.extend(tags)
                    except:
                        pass
                video.search_text = " ".join(filter(None, search_parts))

                # Get file size
                if result["video_path"] and os.path.exists(result["video_path"]):
                    video.file_size = os.path.getsize(result["video_path"])

                db.add(video)

            db.commit()
            db.refresh(video)

            # Add to collection if specified
            if collection_id:
                cv = CollectionVideo(collection_id=collection_id, video_id=video.id)
                db.add(cv)
                db.commit()

            # Update queue
            if queue_item:
                queue_item.status = "completed"
                queue_item.completed_at = datetime.utcnow()
                db.commit()
        else:
            if queue_item:
                queue_item.status = "failed"
                queue_item.error_message = result.get("error", "Unknown error")
                queue_item.completed_at = datetime.utcnow()
                db.commit()

    except Exception as e:
        if queue_item:
            queue_item.status = "failed"
            queue_item.error_message = str(e)
            queue_item.completed_at = datetime.utcnow()
            db.commit()
    finally:
        db.close()


# API Routes

@app.get("/api/health")
async def health():
    return {"status": "ok"}


@app.post("/api/download")
async def download(
    request: DownloadRequest,
    background_tasks: BackgroundTasks,
    db: Session = Depends(get_db)
):
    """Queue a single video for download."""
    # Check if already downloaded
    video_id = get_video_id(request.url)
    existing = db.query(Video).filter(Video.tiktok_id == video_id).first()
    if existing:
        return {"status": "exists", "video_id": existing.id}

    # Add to queue
    queue_item = DownloadQueue(url=request.url)
    db.add(queue_item)
    db.commit()
    db.refresh(queue_item)

    # Start background download
    from database import SessionLocal
    background_tasks.add_task(
        download_task,
        request.url,
        queue_item.id,
        request.collection_id,
        SessionLocal
    )

    return {"status": "queued", "queue_id": queue_item.id}


@app.post("/api/download/bulk")
async def bulk_download(
    request: BulkDownloadRequest,
    background_tasks: BackgroundTasks,
    db: Session = Depends(get_db)
):
    """Queue multiple videos for download."""
    from database import SessionLocal

    queued = []
    for url in request.urls:
        video_id = get_video_id(url)
        existing = db.query(Video).filter(Video.tiktok_id == video_id).first()
        if existing:
            queued.append({"url": url, "status": "exists", "video_id": existing.id})
            continue

        queue_item = DownloadQueue(url=url)
        db.add(queue_item)
        db.commit()
        db.refresh(queue_item)

        background_tasks.add_task(
            download_task,
            url,
            queue_item.id,
            request.collection_id,
            SessionLocal
        )
        queued.append({"url": url, "status": "queued", "queue_id": queue_item.id})

    return {"queued": queued}


@app.get("/api/videos")
async def list_videos(
    search: Optional[str] = None,
    uploader: Optional[str] = None,
    collection_id: Optional[int] = None,
    sort: str = "archived_at",
    order: str = "desc",
    page: int = 1,
    limit: int = 20,
    db: Session = Depends(get_db)
):
    """List archived videos with filtering and pagination."""
    query = db.query(Video)

    # Filters
    if search:
        query = query.filter(
            or_(
                Video.search_text.ilike(f"%{search}%"),
                Video.title.ilike(f"%{search}%"),
                Video.description.ilike(f"%{search}%"),
            )
        )

    if uploader:
        query = query.filter(
            or_(
                Video.uploader.ilike(f"%{uploader}%"),
                Video.uploader_id.ilike(f"%{uploader}%"),
            )
        )

    if collection_id:
        video_ids = db.query(CollectionVideo.video_id).filter(
            CollectionVideo.collection_id == collection_id
        ).subquery()
        query = query.filter(Video.id.in_(video_ids))

    # Count total
    total = query.count()

    # Sorting
    sort_col = getattr(Video, sort, Video.archived_at)
    if order == "desc":
        query = query.order_by(sort_col.desc())
    else:
        query = query.order_by(sort_col.asc())

    # Pagination
    offset = (page - 1) * limit
    videos = query.offset(offset).limit(limit).all()

    # Format response
    result = []
    for v in videos:
        video_dict = {
            "id": v.id,
            "tiktok_id": v.tiktok_id,
            "url": v.url,
            "title": v.title,
            "description": v.description,
            "uploader": v.uploader,
            "uploader_id": v.uploader_id,
            "upload_date": v.upload_date,
            "duration": v.duration,
            "view_count": v.view_count,
            "like_count": v.like_count,
            "comment_count": v.comment_count,
            "video_url": f"/archive/videos/{os.path.basename(v.video_path)}" if v.video_path else None,
            "thumbnail_url": f"/archive/thumbnails/{os.path.basename(v.thumbnail_path)}" if v.thumbnail_path else None,
            "audio_url": f"/archive/audio/{os.path.basename(v.audio_path)}" if v.audio_path else None,
            "tags": json.loads(v.tags) if v.tags else [],
            "music_title": v.music_title,
            "music_author": v.music_author,
            "archived_at": v.archived_at,
            "file_size": v.file_size,
        }
        result.append(video_dict)

    return {
        "videos": result,
        "total": total,
        "page": page,
        "pages": (total + limit - 1) // limit,
    }


@app.get("/api/videos/{video_id}")
async def get_video(video_id: int, db: Session = Depends(get_db)):
    """Get a single video by ID."""
    video = db.query(Video).filter(Video.id == video_id).first()
    if not video:
        raise HTTPException(status_code=404, detail="Video not found")

    return {
        "id": video.id,
        "tiktok_id": video.tiktok_id,
        "url": video.url,
        "title": video.title,
        "description": video.description,
        "uploader": video.uploader,
        "uploader_id": video.uploader_id,
        "upload_date": video.upload_date,
        "duration": video.duration,
        "view_count": video.view_count,
        "like_count": video.like_count,
        "comment_count": video.comment_count,
        "share_count": video.share_count,
        "video_url": f"/archive/videos/{os.path.basename(video.video_path)}" if video.video_path else None,
        "thumbnail_url": f"/archive/thumbnails/{os.path.basename(video.thumbnail_path)}" if video.thumbnail_path else None,
        "audio_url": f"/archive/audio/{os.path.basename(video.audio_path)}" if video.audio_path else None,
        "tags": json.loads(video.tags) if video.tags else [],
        "music_title": video.music_title,
        "music_author": video.music_author,
        "archived_at": video.archived_at,
        "file_size": video.file_size,
        "is_available": video.is_available,
    }


@app.delete("/api/videos/{video_id}")
async def delete_video(video_id: int, delete_files: bool = True, db: Session = Depends(get_db)):
    """Delete a video from the archive."""
    video = db.query(Video).filter(Video.id == video_id).first()
    if not video:
        raise HTTPException(status_code=404, detail="Video not found")

    # Delete files if requested
    if delete_files:
        for path in [video.video_path, video.thumbnail_path, video.audio_path]:
            if path and os.path.exists(path):
                os.remove(path)

    # Remove from collections
    db.query(CollectionVideo).filter(CollectionVideo.video_id == video_id).delete()

    # Delete video record
    db.delete(video)
    db.commit()

    return {"status": "deleted"}


@app.get("/api/collections")
async def list_collections(db: Session = Depends(get_db)):
    """List all collections."""
    collections = db.query(Collection).order_by(Collection.created_at.desc()).all()

    result = []
    for c in collections:
        video_count = db.query(CollectionVideo).filter(
            CollectionVideo.collection_id == c.id
        ).count()
        result.append({
            "id": c.id,
            "name": c.name,
            "description": c.description,
            "video_count": video_count,
            "created_at": c.created_at,
        })

    return {"collections": result}


@app.post("/api/collections")
async def create_collection(request: CollectionCreate, db: Session = Depends(get_db)):
    """Create a new collection."""
    collection = Collection(name=request.name, description=request.description)
    db.add(collection)
    db.commit()
    db.refresh(collection)
    return {"id": collection.id, "name": collection.name}


@app.post("/api/collections/{collection_id}/videos/{video_id}")
async def add_to_collection(collection_id: int, video_id: int, db: Session = Depends(get_db)):
    """Add a video to a collection."""
    # Check existence
    collection = db.query(Collection).filter(Collection.id == collection_id).first()
    if not collection:
        raise HTTPException(status_code=404, detail="Collection not found")

    video = db.query(Video).filter(Video.id == video_id).first()
    if not video:
        raise HTTPException(status_code=404, detail="Video not found")

    # Check if already in collection
    existing = db.query(CollectionVideo).filter(
        CollectionVideo.collection_id == collection_id,
        CollectionVideo.video_id == video_id
    ).first()
    if existing:
        return {"status": "already_exists"}

    cv = CollectionVideo(collection_id=collection_id, video_id=video_id)
    db.add(cv)
    db.commit()
    return {"status": "added"}


@app.delete("/api/collections/{collection_id}/videos/{video_id}")
async def remove_from_collection(collection_id: int, video_id: int, db: Session = Depends(get_db)):
    """Remove a video from a collection."""
    deleted = db.query(CollectionVideo).filter(
        CollectionVideo.collection_id == collection_id,
        CollectionVideo.video_id == video_id
    ).delete()

    if not deleted:
        raise HTTPException(status_code=404, detail="Video not in collection")

    db.commit()
    return {"status": "removed"}


@app.get("/api/queue")
async def get_queue(db: Session = Depends(get_db)):
    """Get download queue status."""
    pending = db.query(DownloadQueue).filter(
        DownloadQueue.status.in_(["pending", "downloading"])
    ).order_by(DownloadQueue.created_at.asc()).all()

    recent = db.query(DownloadQueue).filter(
        DownloadQueue.status.in_(["completed", "failed"])
    ).order_by(DownloadQueue.completed_at.desc()).limit(20).all()

    return {
        "pending": [{"id": q.id, "url": q.url, "status": q.status} for q in pending],
        "recent": [
            {
                "id": q.id,
                "url": q.url,
                "status": q.status,
                "error": q.error_message,
                "completed_at": q.completed_at
            }
            for q in recent
        ]
    }


@app.get("/api/cookies/status")
async def cookies_status():
    """Check if cookies are configured."""
    return {"has_cookies": has_cookies()}


@app.post("/api/cookies")
async def upload_cookies(request: CookiesRequest):
    """
    Upload TikTok cookies in Netscape format.
    Export cookies from browser using a cookie export extension.
    """
    if save_cookies(request.cookies):
        return {"status": "saved"}
    raise HTTPException(status_code=500, detail="Failed to save cookies")


@app.post("/api/download/user")
async def download_user_profile(
    request: UserDownloadRequest,
    background_tasks: BackgroundTasks,
    db: Session = Depends(get_db)
):
    """
    Download all videos from a TikTok user's profile.
    Requires cookies to be configured for most accounts.
    """
    from database import SessionLocal

    username = request.username.lstrip("@")

    try:
        # Fetch user's video list
        videos = get_user_videos(username, limit=request.limit)
    except Exception as e:
        raise HTTPException(
            status_code=400,
            detail=f"Failed to fetch user videos: {str(e)}. Make sure cookies are configured."
        )

    if not videos:
        return {"status": "no_videos", "message": f"No videos found for @{username}"}

    # Create a collection for this user if not specified
    collection_id = request.collection_id
    if not collection_id:
        collection = Collection(
            name=f"@{username}",
            description=f"Videos from TikTok user @{username}"
        )
        db.add(collection)
        db.commit()
        db.refresh(collection)
        collection_id = collection.id

    # Queue all videos for download
    queued = []
    skipped = []
    for video in videos:
        url = video.get("url")
        if not url:
            continue

        video_id = video.get("id") or get_video_id(url)
        existing = db.query(Video).filter(Video.tiktok_id == video_id).first()
        if existing:
            skipped.append({"id": video_id, "url": url})
            continue

        queue_item = DownloadQueue(url=url)
        db.add(queue_item)
        db.commit()
        db.refresh(queue_item)

        background_tasks.add_task(
            download_task,
            url,
            queue_item.id,
            collection_id,
            SessionLocal
        )
        queued.append({"id": video_id, "url": url, "queue_id": queue_item.id})

    return {
        "status": "queued",
        "username": username,
        "collection_id": collection_id,
        "total_found": len(videos),
        "queued": len(queued),
        "skipped": len(skipped),
    }


@app.get("/api/stats")
async def get_stats(db: Session = Depends(get_db)):
    """Get archive statistics."""
    total_videos = db.query(Video).count()
    total_collections = db.query(Collection).count()

    total_size = db.query(func.sum(Video.file_size)).scalar() or 0

    uploaders = db.query(Video.uploader, func.count(Video.id)).group_by(
        Video.uploader
    ).order_by(func.count(Video.id).desc()).limit(10).all()

    pending_downloads = db.query(DownloadQueue).filter(
        DownloadQueue.status.in_(["pending", "downloading"])
    ).count()

    return {
        "total_videos": total_videos,
        "total_collections": total_collections,
        "total_size_bytes": total_size,
        "total_size_gb": round(total_size / (1024 ** 3), 2),
        "pending_downloads": pending_downloads,
        "top_uploaders": [{"name": u[0] or "Unknown", "count": u[1]} for u in uploaders],
        "has_cookies": has_cookies(),
    }


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
