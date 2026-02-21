"""
BookForge - AI-Powered Book Generation System

A comprehensive toolkit for generating, managing, and converting books
using AI language models with support for multiple export formats
and text-to-speech audio generation.
"""

__version__ = "0.1.0"
__author__ = "BookForge Team"
__license__ = "MIT"

# Package-level imports for convenience
from .models import (
    Base,
    Project,
    Chapter,
    Outline,
    StyleGuide,
    GenerationJob,
    ProjectStatus,
    ChapterStatus,
    JobType,
    JobStatus,
    get_engine,
    create_tables,
    get_session,
    init_db,
)

__all__ = [
    # Version info
    "__version__",
    "__author__",
    "__license__",
    # Models
    "Base",
    "Project",
    "Chapter",
    "Outline",
    "StyleGuide",
    "GenerationJob",
    # Enums
    "ProjectStatus",
    "ChapterStatus",
    "JobType",
    "JobStatus",
    # Database utilities
    "get_engine",
    "create_tables",
    "get_session",
    "init_db",
]
