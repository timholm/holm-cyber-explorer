"""
SQLAlchemy database models for the BookForge book generation system.
"""

from datetime import datetime
from typing import Optional, List
from sqlalchemy import (
    create_engine,
    Column,
    Integer,
    String,
    Text,
    DateTime,
    ForeignKey,
    Index,
    JSON,
    Enum as SQLEnum
)
from sqlalchemy.orm import declarative_base, relationship, sessionmaker
import enum


Base = declarative_base()


class ProjectStatus(enum.Enum):
    """Status options for a project."""
    DRAFT = "draft"
    OUTLINING = "outlining"
    WRITING = "writing"
    EDITING = "editing"
    COMPLETED = "completed"
    ARCHIVED = "archived"


class ChapterStatus(enum.Enum):
    """Status options for a chapter."""
    PENDING = "pending"
    GENERATING = "generating"
    DRAFT = "draft"
    REVISED = "revised"
    FINAL = "final"


class JobType(enum.Enum):
    """Types of generation jobs."""
    OUTLINE = "outline"
    CHAPTER = "chapter"
    STYLE_GUIDE = "style_guide"
    AUDIO = "audio"
    REVISION = "revision"


class JobStatus(enum.Enum):
    """Status options for generation jobs."""
    PENDING = "pending"
    RUNNING = "running"
    COMPLETED = "completed"
    FAILED = "failed"
    CANCELLED = "cancelled"


class Project(Base):
    """
    Represents a book project with its metadata and settings.
    """
    __tablename__ = "projects"

    id = Column(Integer, primary_key=True, autoincrement=True)
    title = Column(String(255), nullable=False, index=True)
    description = Column(Text, nullable=True)
    target_audience = Column(String(255), nullable=True)
    genre = Column(String(100), nullable=True, index=True)
    status = Column(
        SQLEnum(ProjectStatus),
        default=ProjectStatus.DRAFT,
        nullable=False,
        index=True
    )
    created_at = Column(DateTime, default=datetime.utcnow, nullable=False)
    updated_at = Column(
        DateTime,
        default=datetime.utcnow,
        onupdate=datetime.utcnow,
        nullable=False
    )

    # Relationships
    chapters = relationship(
        "Chapter",
        back_populates="project",
        cascade="all, delete-orphan",
        order_by="Chapter.number"
    )
    outline = relationship(
        "Outline",
        back_populates="project",
        uselist=False,
        cascade="all, delete-orphan"
    )
    style_guide = relationship(
        "StyleGuide",
        back_populates="project",
        uselist=False,
        cascade="all, delete-orphan"
    )
    generation_jobs = relationship(
        "GenerationJob",
        back_populates="project",
        cascade="all, delete-orphan"
    )

    # Composite indexes
    __table_args__ = (
        Index("ix_projects_status_created", "status", "created_at"),
    )

    def __repr__(self) -> str:
        return f"<Project(id={self.id}, title='{self.title}', status={self.status.value})>"


class Chapter(Base):
    """
    Represents a chapter within a book project.
    """
    __tablename__ = "chapters"

    id = Column(Integer, primary_key=True, autoincrement=True)
    project_id = Column(
        Integer,
        ForeignKey("projects.id", ondelete="CASCADE"),
        nullable=False,
        index=True
    )
    number = Column(Integer, nullable=False)
    title = Column(String(255), nullable=True)
    content = Column(Text, nullable=True)
    word_count = Column(Integer, default=0, nullable=False)
    status = Column(
        SQLEnum(ChapterStatus),
        default=ChapterStatus.PENDING,
        nullable=False,
        index=True
    )
    audio_path = Column(String(512), nullable=True)
    created_at = Column(DateTime, default=datetime.utcnow, nullable=False)

    # Relationships
    project = relationship("Project", back_populates="chapters")
    generation_jobs = relationship(
        "GenerationJob",
        back_populates="chapter",
        cascade="all, delete-orphan"
    )

    # Composite indexes
    __table_args__ = (
        Index("ix_chapters_project_number", "project_id", "number", unique=True),
        Index("ix_chapters_project_status", "project_id", "status"),
    )

    def __repr__(self) -> str:
        return f"<Chapter(id={self.id}, project_id={self.project_id}, number={self.number}, title='{self.title}')>"


class Outline(Base):
    """
    Stores the structural outline for a book project as JSON.
    """
    __tablename__ = "outlines"

    id = Column(Integer, primary_key=True, autoincrement=True)
    project_id = Column(
        Integer,
        ForeignKey("projects.id", ondelete="CASCADE"),
        nullable=False,
        unique=True,
        index=True
    )
    structure = Column(JSON, nullable=True)
    created_at = Column(DateTime, default=datetime.utcnow, nullable=False)

    # Relationships
    project = relationship("Project", back_populates="outline")

    def __repr__(self) -> str:
        return f"<Outline(id={self.id}, project_id={self.project_id})>"


class StyleGuide(Base):
    """
    Stores the style guide content for a book project.
    """
    __tablename__ = "style_guides"

    id = Column(Integer, primary_key=True, autoincrement=True)
    project_id = Column(
        Integer,
        ForeignKey("projects.id", ondelete="CASCADE"),
        nullable=False,
        unique=True,
        index=True
    )
    content = Column(Text, nullable=True)
    created_at = Column(DateTime, default=datetime.utcnow, nullable=False)

    # Relationships
    project = relationship("Project", back_populates="style_guide")

    def __repr__(self) -> str:
        return f"<StyleGuide(id={self.id}, project_id={self.project_id})>"


class GenerationJob(Base):
    """
    Tracks generation jobs for various tasks like outline, chapter, or audio generation.
    """
    __tablename__ = "generation_jobs"

    id = Column(Integer, primary_key=True, autoincrement=True)
    project_id = Column(
        Integer,
        ForeignKey("projects.id", ondelete="CASCADE"),
        nullable=False,
        index=True
    )
    chapter_id = Column(
        Integer,
        ForeignKey("chapters.id", ondelete="SET NULL"),
        nullable=True,
        index=True
    )
    job_type = Column(SQLEnum(JobType), nullable=False, index=True)
    status = Column(
        SQLEnum(JobStatus),
        default=JobStatus.PENDING,
        nullable=False,
        index=True
    )
    started_at = Column(DateTime, nullable=True)
    completed_at = Column(DateTime, nullable=True)
    error = Column(Text, nullable=True)

    # Relationships
    project = relationship("Project", back_populates="generation_jobs")
    chapter = relationship("Chapter", back_populates="generation_jobs")

    # Composite indexes
    __table_args__ = (
        Index("ix_jobs_project_status", "project_id", "status"),
        Index("ix_jobs_type_status", "job_type", "status"),
    )

    def __repr__(self) -> str:
        return f"<GenerationJob(id={self.id}, project_id={self.project_id}, type={self.job_type.value}, status={self.status.value})>"


# Database initialization utilities
def get_engine(database_url: str = "sqlite:///bookforge.db"):
    """
    Create and return a SQLAlchemy engine.
    
    Args:
        database_url: Database connection string. Defaults to SQLite.
    
    Returns:
        SQLAlchemy Engine instance.
    """
    return create_engine(database_url, echo=False)


def create_tables(engine):
    """
    Create all tables in the database.
    
    Args:
        engine: SQLAlchemy Engine instance.
    """
    Base.metadata.create_all(engine)


def get_session(engine):
    """
    Create and return a new database session.
    
    Args:
        engine: SQLAlchemy Engine instance.
    
    Returns:
        SQLAlchemy Session instance.
    """
    Session = sessionmaker(bind=engine)
    return Session()


def init_db(database_url: str = "sqlite:///bookforge.db"):
    """
    Initialize the database and return engine and session.
    
    Args:
        database_url: Database connection string. Defaults to SQLite.
    
    Returns:
        Tuple of (engine, session).
    """
    engine = get_engine(database_url)
    create_tables(engine)
    session = get_session(engine)
    return engine, session


if __name__ == "__main__":
    # Example usage: Initialize database and create tables
    engine, session = init_db()
    print("Database initialized successfully!")
    print(f"Tables created: {list(Base.metadata.tables.keys())}")
