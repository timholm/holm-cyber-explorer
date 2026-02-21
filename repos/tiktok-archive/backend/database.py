from sqlalchemy import create_engine, Column, Integer, String, DateTime, Text, Boolean, Float
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from datetime import datetime
import os

DATABASE_URL = os.getenv("DATABASE_URL", "sqlite:///./tiktok_archive.db")

engine = create_engine(DATABASE_URL, connect_args={"check_same_thread": False})
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
Base = declarative_base()


class Video(Base):
    __tablename__ = "videos"

    id = Column(Integer, primary_key=True, index=True)
    tiktok_id = Column(String(64), unique=True, index=True)
    url = Column(String(512), nullable=False)

    # Metadata from yt-dlp
    title = Column(String(512))
    description = Column(Text)
    uploader = Column(String(256), index=True)
    uploader_id = Column(String(256), index=True)
    upload_date = Column(DateTime)
    duration = Column(Float)
    view_count = Column(Integer)
    like_count = Column(Integer)
    comment_count = Column(Integer)
    share_count = Column(Integer)

    # Local storage
    video_path = Column(String(512))
    thumbnail_path = Column(String(512))
    audio_path = Column(String(512))

    # Tags and categories
    tags = Column(Text)  # JSON array stored as text
    music_title = Column(String(512))
    music_author = Column(String(256))

    # Archive metadata
    archived_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)
    is_available = Column(Boolean, default=True)  # Track if original is still up
    file_size = Column(Integer)  # bytes

    # Search optimization
    search_text = Column(Text, index=True)  # Combined searchable text


class Collection(Base):
    __tablename__ = "collections"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String(256), nullable=False)
    description = Column(Text)
    created_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)


class CollectionVideo(Base):
    __tablename__ = "collection_videos"

    id = Column(Integer, primary_key=True, index=True)
    collection_id = Column(Integer, index=True)
    video_id = Column(Integer, index=True)
    added_at = Column(DateTime, default=datetime.utcnow)


class DownloadQueue(Base):
    __tablename__ = "download_queue"

    id = Column(Integer, primary_key=True, index=True)
    url = Column(String(512), nullable=False)
    status = Column(String(32), default="pending")  # pending, downloading, completed, failed
    error_message = Column(Text)
    created_at = Column(DateTime, default=datetime.utcnow)
    started_at = Column(DateTime)
    completed_at = Column(DateTime)


def init_db():
    Base.metadata.create_all(bind=engine)


def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
