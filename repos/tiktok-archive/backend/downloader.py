import yt_dlp
import subprocess
import os
import json
import hashlib
from datetime import datetime
from pathlib import Path
from typing import Optional, Dict, Any, List

ARCHIVE_DIR = os.getenv("ARCHIVE_DIR", "./archive")
VIDEOS_DIR = os.path.join(ARCHIVE_DIR, "videos")
THUMBNAILS_DIR = os.path.join(ARCHIVE_DIR, "thumbnails")
AUDIO_DIR = os.path.join(ARCHIVE_DIR, "audio")
COOKIES_FILE = os.path.join(ARCHIVE_DIR, "cookies.txt")

# Ensure directories exist
for d in [VIDEOS_DIR, THUMBNAILS_DIR, AUDIO_DIR]:
    os.makedirs(d, exist_ok=True)


def get_video_id(url: str) -> str:
    """Extract or generate a unique ID for the video."""
    if "tiktok.com" in url:
        parts = url.split("/")
        for i, part in enumerate(parts):
            if part == "video" and i + 1 < len(parts):
                video_id = parts[i + 1].split("?")[0]
                if video_id.isdigit():
                    return video_id
    return hashlib.md5(url.encode()).hexdigest()[:16]


def get_ydl_opts(output_template: str = None, extract_only: bool = False) -> dict:
    """Get yt-dlp options with cookie support."""
    opts = {
        "quiet": False,
        "no_warnings": False,
        "extract_flat": False,
        # Browser impersonation for better compatibility
        "http_headers": {
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
            "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
            "Accept-Language": "en-US,en;q=0.5",
        },
    }

    # Use cookies if available
    if os.path.exists(COOKIES_FILE):
        opts["cookiefile"] = COOKIES_FILE

    if output_template:
        opts["outtmpl"] = output_template

    if extract_only:
        opts["skip_download"] = True

    return opts


def extract_metadata(url: str) -> Optional[Dict[str, Any]]:
    """Extract metadata from a TikTok URL without downloading."""
    opts = get_ydl_opts(extract_only=True)

    try:
        with yt_dlp.YoutubeDL(opts) as ydl:
            info = ydl.extract_info(url, download=False)
            return info
    except Exception as e:
        print(f"Error extracting metadata: {e}")
        return None


def get_user_videos(username: str, limit: int = None) -> List[Dict[str, Any]]:
    """
    Get all video URLs from a TikTok user's profile.
    Username should be like '@username' or just 'username'.
    """
    if not username.startswith("@"):
        username = f"@{username}"

    user_url = f"https://www.tiktok.com/{username}"

    opts = get_ydl_opts()
    opts["extract_flat"] = "in_playlist"
    opts["playlistend"] = limit if limit else None

    videos = []
    try:
        with yt_dlp.YoutubeDL(opts) as ydl:
            result = ydl.extract_info(user_url, download=False)

            if result and "entries" in result:
                for entry in result["entries"]:
                    if entry:
                        videos.append({
                            "id": entry.get("id"),
                            "url": entry.get("url") or entry.get("webpage_url"),
                            "title": entry.get("title"),
                        })
    except Exception as e:
        print(f"Error fetching user videos: {e}")
        raise e

    return videos


def download_video(url: str, video_id: str = None) -> Dict[str, Any]:
    """
    Download a TikTok video with metadata.
    Returns dict with paths and metadata.
    """
    if not video_id:
        video_id = get_video_id(url)

    video_path = os.path.join(VIDEOS_DIR, f"{video_id}.mp4")
    thumbnail_path = os.path.join(THUMBNAILS_DIR, f"{video_id}.jpg")
    audio_path = os.path.join(AUDIO_DIR, f"{video_id}.mp3")

    result = {
        "video_id": video_id,
        "success": False,
        "video_path": None,
        "thumbnail_path": None,
        "audio_path": None,
        "metadata": None,
        "error": None,
    }

    opts = get_ydl_opts(output_template=os.path.join(VIDEOS_DIR, f"{video_id}.%(ext)s"))
    opts["format"] = "best"
    opts["writeinfojson"] = True
    opts["writethumbnail"] = True

    try:
        with yt_dlp.YoutubeDL(opts) as ydl:
            info = ydl.extract_info(url, download=True)

            # Find the actual downloaded video file
            downloaded_file = None
            for ext in ["mp4", "webm", "mkv"]:
                potential_path = os.path.join(VIDEOS_DIR, f"{video_id}.{ext}")
                if os.path.exists(potential_path):
                    downloaded_file = potential_path
                    break

            if downloaded_file and downloaded_file != video_path:
                if not downloaded_file.endswith(".mp4"):
                    subprocess.run([
                        "ffmpeg", "-i", downloaded_file,
                        "-c:v", "copy", "-c:a", "copy",
                        video_path, "-y"
                    ], capture_output=True)
                    os.remove(downloaded_file)
                else:
                    video_path = downloaded_file

            result["video_path"] = video_path
            result["metadata"] = info

            # Handle thumbnail
            thumb_candidates = [
                os.path.join(VIDEOS_DIR, f"{video_id}.jpg"),
                os.path.join(VIDEOS_DIR, f"{video_id}.webp"),
                os.path.join(VIDEOS_DIR, f"{video_id}.png"),
            ]

            for thumb in thumb_candidates:
                if os.path.exists(thumb):
                    if thumb != thumbnail_path:
                        subprocess.run([
                            "ffmpeg", "-i", thumb,
                            thumbnail_path, "-y"
                        ], capture_output=True)
                        os.remove(thumb)
                    result["thumbnail_path"] = thumbnail_path
                    break

            # Generate thumbnail from video if not available
            if not result["thumbnail_path"] and os.path.exists(video_path):
                subprocess.run([
                    "ffmpeg", "-i", video_path,
                    "-ss", "00:00:01",
                    "-vframes", "1",
                    "-vf", "scale=480:-1",
                    thumbnail_path, "-y"
                ], capture_output=True)
                if os.path.exists(thumbnail_path):
                    result["thumbnail_path"] = thumbnail_path

            # Extract audio
            if os.path.exists(video_path):
                subprocess.run([
                    "ffmpeg", "-i", video_path,
                    "-vn", "-acodec", "libmp3lame",
                    "-q:a", "2",
                    audio_path, "-y"
                ], capture_output=True)
                if os.path.exists(audio_path):
                    result["audio_path"] = audio_path

            result["success"] = True

    except Exception as e:
        result["error"] = str(e)

    return result


def parse_metadata(info: Dict[str, Any]) -> Dict[str, Any]:
    """Parse yt-dlp metadata into our database schema."""
    upload_date = None
    if info.get("upload_date"):
        try:
            upload_date = datetime.strptime(info["upload_date"], "%Y%m%d")
        except:
            pass

    tags = info.get("tags", [])
    if isinstance(tags, list):
        tags = json.dumps(tags)

    return {
        "tiktok_id": info.get("id"),
        "title": info.get("title", "")[:512],
        "description": info.get("description", ""),
        "uploader": info.get("uploader", info.get("creator", ""))[:256],
        "uploader_id": info.get("uploader_id", info.get("channel_id", ""))[:256],
        "upload_date": upload_date,
        "duration": info.get("duration"),
        "view_count": info.get("view_count"),
        "like_count": info.get("like_count"),
        "comment_count": info.get("comment_count"),
        "share_count": info.get("repost_count"),
        "tags": tags,
        "music_title": info.get("track", "")[:512] if info.get("track") else None,
        "music_author": info.get("artist", "")[:256] if info.get("artist") else None,
    }


def save_cookies(cookies_content: str) -> bool:
    """Save cookies content to file."""
    try:
        with open(COOKIES_FILE, "w") as f:
            f.write(cookies_content)
        return True
    except Exception as e:
        print(f"Error saving cookies: {e}")
        return False


def has_cookies() -> bool:
    """Check if cookies file exists."""
    return os.path.exists(COOKIES_FILE)


def check_video_available(url: str) -> bool:
    """Check if a TikTok video is still available."""
    try:
        metadata = extract_metadata(url)
        return metadata is not None
    except:
        return False
