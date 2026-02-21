"""
BookForge Export System

Provides export functionality for completed books in multiple formats:
- Markdown (single file or per-chapter)
- HTML (styled, printable)
- EPUB (for e-readers)
- PDF (via markdown conversion)
- Audio ZIP (bundle all audio files)
"""

import os
import shutil
import zipfile
import uuid
import tempfile
import subprocess
from datetime import datetime
from pathlib import Path
from typing import Optional, List, Tuple, Dict, Any
from dataclasses import dataclass

import markdown
from sqlalchemy.orm import Session

from models import (
    Project, Chapter, Outline, StyleGuide,
    ProjectStatus, ChapterStatus,
    get_engine, get_session
)

# Try to import ebooklib for EPUB generation
try:
    from ebooklib import epub
    EBOOKLIB_AVAILABLE = True
except ImportError:
    EBOOKLIB_AVAILABLE = False


@dataclass
class ExportResult:
    """Result of an export operation."""
    success: bool
    output_path: Optional[str]
    format: str
    message: str
    file_size: Optional[int] = None


class BookExporter:
    """
    Handles exporting book projects to various formats.
    """

    # HTML template for styled, printable output
    HTML_TEMPLATE = '''<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{title}</title>
    <style>
        @page {{
            size: A4;
            margin: 2.5cm;
        }}

        body {{
            font-family: Georgia, 'Times New Roman', serif;
            font-size: 12pt;
            line-height: 1.6;
            color: #333;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #fff;
        }}

        .book-cover {{
            text-align: center;
            page-break-after: always;
            padding: 100px 20px;
        }}

        .book-title {{
            font-size: 2.5em;
            font-weight: bold;
            margin-bottom: 0.5em;
            color: #1a1a1a;
        }}

        .book-description {{
            font-size: 1.2em;
            color: #666;
            font-style: italic;
            margin-top: 2em;
        }}

        .book-meta {{
            margin-top: 3em;
            font-size: 0.9em;
            color: #888;
        }}

        .toc {{
            page-break-after: always;
        }}

        .toc h2 {{
            font-size: 1.5em;
            border-bottom: 2px solid #333;
            padding-bottom: 10px;
        }}

        .toc ul {{
            list-style-type: none;
            padding: 0;
        }}

        .toc li {{
            margin: 0.5em 0;
            padding: 0.3em 0;
        }}

        .toc a {{
            text-decoration: none;
            color: #333;
        }}

        .toc a:hover {{
            color: #0066cc;
        }}

        .toc .chapter-number {{
            font-weight: bold;
            margin-right: 1em;
        }}

        .chapter {{
            page-break-before: always;
        }}

        .chapter-header {{
            text-align: center;
            margin-bottom: 2em;
        }}

        .chapter-number-label {{
            font-size: 0.9em;
            text-transform: uppercase;
            letter-spacing: 0.2em;
            color: #888;
        }}

        .chapter-title {{
            font-size: 1.8em;
            font-weight: bold;
            margin-top: 0.3em;
            color: #1a1a1a;
        }}

        .chapter-content {{
            text-align: justify;
            hyphens: auto;
        }}

        .chapter-content p {{
            margin: 1em 0;
            text-indent: 1.5em;
        }}

        .chapter-content p:first-child {{
            text-indent: 0;
        }}

        .chapter-content p:first-child::first-letter {{
            font-size: 3em;
            float: left;
            line-height: 0.8;
            padding-right: 0.1em;
            font-weight: bold;
        }}

        h1, h2, h3, h4, h5, h6 {{
            color: #1a1a1a;
            margin-top: 1.5em;
            margin-bottom: 0.5em;
        }}

        blockquote {{
            border-left: 3px solid #ccc;
            margin-left: 0;
            padding-left: 1.5em;
            font-style: italic;
            color: #555;
        }}

        hr {{
            border: none;
            text-align: center;
            margin: 2em 0;
        }}

        hr::before {{
            content: "* * *";
            color: #888;
        }}

        @media print {{
            body {{
                max-width: none;
                padding: 0;
            }}

            .chapter {{
                page-break-before: always;
            }}

            .no-print {{
                display: none;
            }}
        }}
    </style>
</head>
<body>
{content}
</body>
</html>'''

    # EPUB CSS for e-reader formatting
    EPUB_CSS = '''
body {
    font-family: serif;
    line-height: 1.6;
    margin: 1em;
}

h1 {
    text-align: center;
    font-size: 2em;
    margin-top: 3em;
    margin-bottom: 1em;
}

h2 {
    text-align: center;
    font-size: 1.5em;
    margin-top: 2em;
    margin-bottom: 1em;
}

h3 {
    font-size: 1.2em;
    margin-top: 1.5em;
}

p {
    text-indent: 1.5em;
    margin: 0.5em 0;
    text-align: justify;
}

p.first {
    text-indent: 0;
}

blockquote {
    font-style: italic;
    margin: 1em 2em;
    padding-left: 1em;
    border-left: 2px solid #999;
}

.chapter-header {
    text-align: center;
    margin: 3em 0 2em 0;
}

.chapter-number {
    font-size: 0.9em;
    text-transform: uppercase;
    letter-spacing: 0.2em;
    color: #666;
}

.chapter-title {
    font-size: 1.8em;
    font-weight: bold;
    margin-top: 0.5em;
}

hr {
    border: none;
    text-align: center;
    margin: 2em 0;
}

hr::after {
    content: "* * *";
    color: #666;
}
'''

    def __init__(self, database_url: str = "sqlite:///bookforge.db"):
        """
        Initialize the exporter with database connection.

        Args:
            database_url: Database connection string.
        """
        self.engine = get_engine(database_url)
        self.md = markdown.Markdown(
            extensions=['extra', 'smarty', 'toc'],
            output_format='html5'
        )

    def _get_session(self) -> Session:
        """Create and return a new database session."""
        return get_session(self.engine)

    def _get_project_with_chapters(
        self,
        session: Session,
        project_id: int
    ) -> Tuple[Optional[Project], List[Chapter]]:
        """
        Fetch a project and its chapters from the database.

        Args:
            session: Database session.
            project_id: ID of the project to fetch.

        Returns:
            Tuple of (project, chapters) or (None, []) if not found.
        """
        project = session.query(Project).filter(Project.id == project_id).first()
        if not project:
            return None, []

        chapters = (
            session.query(Chapter)
            .filter(Chapter.project_id == project_id)
            .order_by(Chapter.number)
            .all()
        )
        return project, chapters

    def _generate_metadata(self, project: Project) -> Dict[str, Any]:
        """Generate metadata dictionary for a project."""
        return {
            'title': project.title,
            'description': project.description or '',
            'genre': project.genre or 'General',
            'target_audience': project.target_audience or '',
            'created_at': project.created_at.strftime('%Y-%m-%d'),
            'updated_at': project.updated_at.strftime('%Y-%m-%d'),
            'status': project.status.value,
        }

    def _ensure_output_dir(self, output_path: str) -> str:
        """Ensure the output directory exists and return the full path."""
        output_dir = os.path.dirname(output_path)
        if output_dir:
            os.makedirs(output_dir, exist_ok=True)
        return output_path

    def _get_file_size(self, path: str) -> Optional[int]:
        """Get file size in bytes, or None if file doesn't exist."""
        try:
            return os.path.getsize(path)
        except OSError:
            return None

    # =========================================================================
    # MARKDOWN EXPORT
    # =========================================================================

    def export_markdown(
        self,
        project_id: int,
        output_path: str,
        single_file: bool = True,
        include_toc: bool = True,
        include_metadata: bool = True
    ) -> ExportResult:
        """
        Export a project to Markdown format.

        Args:
            project_id: ID of the project to export.
            output_path: Path for the output file(s).
                        For single_file=True: path to the .md file
                        For single_file=False: path to the output directory
            single_file: If True, export as single file. If False, one file per chapter.
            include_toc: Include table of contents.
            include_metadata: Include project metadata at the beginning.

        Returns:
            ExportResult with success status and details.
        """
        session = self._get_session()
        try:
            project, chapters = self._get_project_with_chapters(session, project_id)

            if not project:
                return ExportResult(
                    success=False,
                    output_path=None,
                    format='markdown',
                    message=f"Project with ID {project_id} not found."
                )

            if not chapters:
                return ExportResult(
                    success=False,
                    output_path=None,
                    format='markdown',
                    message=f"Project '{project.title}' has no chapters to export."
                )

            if single_file:
                return self._export_markdown_single(
                    project, chapters, output_path, include_toc, include_metadata
                )
            else:
                return self._export_markdown_chapters(
                    project, chapters, output_path, include_toc, include_metadata
                )
        finally:
            session.close()

    def _export_markdown_single(
        self,
        project: Project,
        chapters: List[Chapter],
        output_path: str,
        include_toc: bool,
        include_metadata: bool
    ) -> ExportResult:
        """Export to a single Markdown file."""
        self._ensure_output_dir(output_path)

        lines = []

        # Title
        lines.append(f"# {project.title}")
        lines.append("")

        # Metadata
        if include_metadata:
            if project.description:
                lines.append(f"*{project.description}*")
                lines.append("")

            meta = self._generate_metadata(project)
            lines.append("---")
            lines.append("")
            lines.append(f"**Genre:** {meta['genre']}")
            if meta['target_audience']:
                lines.append(f"**Target Audience:** {meta['target_audience']}")
            lines.append(f"**Last Updated:** {meta['updated_at']}")
            lines.append("")
            lines.append("---")
            lines.append("")

        # Table of Contents
        if include_toc:
            lines.append("## Table of Contents")
            lines.append("")
            for chapter in chapters:
                chapter_title = chapter.title or f"Chapter {chapter.number}"
                anchor = f"chapter-{chapter.number}"
                lines.append(f"- [Chapter {chapter.number}: {chapter_title}](#{anchor})")
            lines.append("")
            lines.append("---")
            lines.append("")

        # Chapters
        for chapter in chapters:
            chapter_title = chapter.title or f"Chapter {chapter.number}"
            anchor = f"chapter-{chapter.number}"

            lines.append(f"## Chapter {chapter.number}: {chapter_title} {{#{anchor}}}")
            lines.append("")

            if chapter.content:
                lines.append(chapter.content)
            else:
                lines.append("*[Chapter content not yet available]*")

            lines.append("")
            lines.append("---")
            lines.append("")

        content = "\n".join(lines)

        with open(output_path, 'w', encoding='utf-8') as f:
            f.write(content)

        return ExportResult(
            success=True,
            output_path=output_path,
            format='markdown',
            message=f"Successfully exported '{project.title}' to Markdown.",
            file_size=self._get_file_size(output_path)
        )

    def _export_markdown_chapters(
        self,
        project: Project,
        chapters: List[Chapter],
        output_dir: str,
        include_toc: bool,
        include_metadata: bool
    ) -> ExportResult:
        """Export to separate Markdown files per chapter."""
        os.makedirs(output_dir, exist_ok=True)

        # Create index file
        index_lines = []
        index_lines.append(f"# {project.title}")
        index_lines.append("")

        if include_metadata:
            if project.description:
                index_lines.append(f"*{project.description}*")
                index_lines.append("")

            meta = self._generate_metadata(project)
            index_lines.append("---")
            index_lines.append("")
            index_lines.append(f"**Genre:** {meta['genre']}")
            if meta['target_audience']:
                index_lines.append(f"**Target Audience:** {meta['target_audience']}")
            index_lines.append(f"**Last Updated:** {meta['updated_at']}")
            index_lines.append("")
            index_lines.append("---")
            index_lines.append("")

        if include_toc:
            index_lines.append("## Chapters")
            index_lines.append("")

        # Create chapter files
        total_size = 0
        for chapter in chapters:
            chapter_title = chapter.title or f"Chapter {chapter.number}"
            filename = f"chapter_{chapter.number:03d}.md"
            filepath = os.path.join(output_dir, filename)

            chapter_lines = []
            chapter_lines.append(f"# Chapter {chapter.number}: {chapter_title}")
            chapter_lines.append("")

            if chapter.content:
                chapter_lines.append(chapter.content)
            else:
                chapter_lines.append("*[Chapter content not yet available]*")

            with open(filepath, 'w', encoding='utf-8') as f:
                f.write("\n".join(chapter_lines))

            total_size += self._get_file_size(filepath) or 0

            if include_toc:
                index_lines.append(f"- [Chapter {chapter.number}: {chapter_title}]({filename})")

        # Write index file
        index_path = os.path.join(output_dir, "index.md")
        with open(index_path, 'w', encoding='utf-8') as f:
            f.write("\n".join(index_lines))

        total_size += self._get_file_size(index_path) or 0

        return ExportResult(
            success=True,
            output_path=output_dir,
            format='markdown',
            message=f"Successfully exported '{project.title}' to {len(chapters) + 1} Markdown files.",
            file_size=total_size
        )

    # =========================================================================
    # HTML EXPORT
    # =========================================================================

    def export_html(
        self,
        project_id: int,
        output_path: str,
        include_cover: bool = True,
        include_toc: bool = True
    ) -> ExportResult:
        """
        Export a project to styled, printable HTML format.

        Args:
            project_id: ID of the project to export.
            output_path: Path for the output HTML file.
            include_cover: Include a cover page.
            include_toc: Include table of contents.

        Returns:
            ExportResult with success status and details.
        """
        session = self._get_session()
        try:
            project, chapters = self._get_project_with_chapters(session, project_id)

            if not project:
                return ExportResult(
                    success=False,
                    output_path=None,
                    format='html',
                    message=f"Project with ID {project_id} not found."
                )

            if not chapters:
                return ExportResult(
                    success=False,
                    output_path=None,
                    format='html',
                    message=f"Project '{project.title}' has no chapters to export."
                )

            self._ensure_output_dir(output_path)

            content_parts = []

            # Cover page
            if include_cover:
                cover_html = f'''
<div class="book-cover">
    <div class="book-title">{project.title}</div>
    {f'<div class="book-description">{project.description}</div>' if project.description else ''}
    <div class="book-meta">
        {f'<p>Genre: {project.genre}</p>' if project.genre else ''}
        {f'<p>For: {project.target_audience}</p>' if project.target_audience else ''}
        <p>Generated: {datetime.now().strftime('%B %d, %Y')}</p>
    </div>
</div>'''
                content_parts.append(cover_html)

            # Table of Contents
            if include_toc:
                toc_items = []
                for chapter in chapters:
                    chapter_title = chapter.title or f"Chapter {chapter.number}"
                    toc_items.append(
                        f'<li><span class="chapter-number">Chapter {chapter.number}</span>'
                        f'<a href="#chapter-{chapter.number}">{chapter_title}</a></li>'
                    )

                toc_html = f'''
<div class="toc">
    <h2>Table of Contents</h2>
    <ul>
        {"".join(toc_items)}
    </ul>
</div>'''
                content_parts.append(toc_html)

            # Chapters
            for chapter in chapters:
                chapter_title = chapter.title or f"Chapter {chapter.number}"

                # Convert chapter content from markdown to HTML
                self.md.reset()
                if chapter.content:
                    chapter_content_html = self.md.convert(chapter.content)
                else:
                    chapter_content_html = "<p><em>[Chapter content not yet available]</em></p>"

                chapter_html = f'''
<div class="chapter" id="chapter-{chapter.number}">
    <div class="chapter-header">
        <div class="chapter-number-label">Chapter {chapter.number}</div>
        <div class="chapter-title">{chapter_title}</div>
    </div>
    <div class="chapter-content">
        {chapter_content_html}
    </div>
</div>'''
                content_parts.append(chapter_html)

            # Combine all content
            full_content = "\n".join(content_parts)
            final_html = self.HTML_TEMPLATE.format(
                title=project.title,
                content=full_content
            )

            with open(output_path, 'w', encoding='utf-8') as f:
                f.write(final_html)

            return ExportResult(
                success=True,
                output_path=output_path,
                format='html',
                message=f"Successfully exported '{project.title}' to HTML.",
                file_size=self._get_file_size(output_path)
            )
        finally:
            session.close()

    # =========================================================================
    # EPUB EXPORT
    # =========================================================================

    def export_epub(
        self,
        project_id: int,
        output_path: str,
        author: str = "BookForge"
    ) -> ExportResult:
        """
        Export a project to EPUB format for e-readers.

        Uses ebooklib if available, otherwise creates EPUB manually.

        Args:
            project_id: ID of the project to export.
            output_path: Path for the output EPUB file.
            author: Author name for the EPUB metadata.

        Returns:
            ExportResult with success status and details.
        """
        session = self._get_session()
        try:
            project, chapters = self._get_project_with_chapters(session, project_id)

            if not project:
                return ExportResult(
                    success=False,
                    output_path=None,
                    format='epub',
                    message=f"Project with ID {project_id} not found."
                )

            if not chapters:
                return ExportResult(
                    success=False,
                    output_path=None,
                    format='epub',
                    message=f"Project '{project.title}' has no chapters to export."
                )

            self._ensure_output_dir(output_path)

            if EBOOKLIB_AVAILABLE:
                return self._export_epub_ebooklib(project, chapters, output_path, author)
            else:
                return self._export_epub_manual(project, chapters, output_path, author)
        finally:
            session.close()

    def _export_epub_ebooklib(
        self,
        project: Project,
        chapters: List[Chapter],
        output_path: str,
        author: str
    ) -> ExportResult:
        """Export EPUB using ebooklib."""
        book = epub.EpubBook()

        # Set metadata
        book.set_identifier(f"bookforge-{project.id}-{uuid.uuid4().hex[:8]}")
        book.set_title(project.title)
        book.set_language('en')
        book.add_author(author)

        if project.description:
            book.add_metadata('DC', 'description', project.description)
        if project.genre:
            book.add_metadata('DC', 'subject', project.genre)

        # Add CSS
        css = epub.EpubItem(
            uid="style",
            file_name="style/main.css",
            media_type="text/css",
            content=self.EPUB_CSS.encode('utf-8')
        )
        book.add_item(css)

        # Create title page
        title_content = f'''
<html>
<head>
    <title>{project.title}</title>
    <link href="style/main.css" rel="stylesheet" type="text/css"/>
</head>
<body>
    <h1>{project.title}</h1>
    {f'<p><em>{project.description}</em></p>' if project.description else ''}
    <p style="margin-top: 2em;">By {author}</p>
    {f'<p>Genre: {project.genre}</p>' if project.genre else ''}
</body>
</html>'''

        title_page = epub.EpubHtml(
            title='Title Page',
            file_name='title.xhtml',
            content=title_content,
            lang='en'
        )
        title_page.add_item(css)
        book.add_item(title_page)

        # Create chapter pages
        epub_chapters = []
        for chapter in chapters:
            chapter_title = chapter.title or f"Chapter {chapter.number}"

            self.md.reset()
            if chapter.content:
                chapter_content_html = self.md.convert(chapter.content)
            else:
                chapter_content_html = "<p><em>[Chapter content not yet available]</em></p>"

            chapter_html = f'''
<html>
<head>
    <title>Chapter {chapter.number}: {chapter_title}</title>
    <link href="style/main.css" rel="stylesheet" type="text/css"/>
</head>
<body>
    <div class="chapter-header">
        <p class="chapter-number">Chapter {chapter.number}</p>
        <h2 class="chapter-title">{chapter_title}</h2>
    </div>
    {chapter_content_html}
</body>
</html>'''

            epub_chapter = epub.EpubHtml(
                title=f"Chapter {chapter.number}: {chapter_title}",
                file_name=f"chapter_{chapter.number:03d}.xhtml",
                content=chapter_html,
                lang='en'
            )
            epub_chapter.add_item(css)
            book.add_item(epub_chapter)
            epub_chapters.append(epub_chapter)

        # Create table of contents
        book.toc = [title_page] + epub_chapters

        # Add navigation files
        book.add_item(epub.EpubNcx())
        book.add_item(epub.EpubNav())

        # Define spine
        book.spine = [title_page, 'nav'] + epub_chapters

        # Write EPUB file
        epub.write_epub(output_path, book)

        return ExportResult(
            success=True,
            output_path=output_path,
            format='epub',
            message=f"Successfully exported '{project.title}' to EPUB using ebooklib.",
            file_size=self._get_file_size(output_path)
        )

    def _export_epub_manual(
        self,
        project: Project,
        chapters: List[Chapter],
        output_path: str,
        author: str
    ) -> ExportResult:
        """Create EPUB manually without ebooklib."""
        epub_uuid = f"bookforge-{project.id}-{uuid.uuid4().hex[:8]}"

        with tempfile.TemporaryDirectory() as tmpdir:
            # Create EPUB structure
            meta_inf = os.path.join(tmpdir, "META-INF")
            oebps = os.path.join(tmpdir, "OEBPS")
            styles = os.path.join(oebps, "styles")

            os.makedirs(meta_inf)
            os.makedirs(oebps)
            os.makedirs(styles)

            # mimetype file (must be first, uncompressed)
            mimetype_path = os.path.join(tmpdir, "mimetype")
            with open(mimetype_path, 'w') as f:
                f.write("application/epub+zip")

            # container.xml
            container_xml = '''<?xml version="1.0" encoding="UTF-8"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
    <rootfiles>
        <rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml"/>
    </rootfiles>
</container>'''
            with open(os.path.join(meta_inf, "container.xml"), 'w') as f:
                f.write(container_xml)

            # CSS file
            with open(os.path.join(styles, "main.css"), 'w') as f:
                f.write(self.EPUB_CSS)

            # Create HTML files for chapters
            manifest_items = []
            spine_items = []
            toc_items = []

            # Title page
            title_content = f'''<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" lang="en">
<head>
    <title>{project.title}</title>
    <link href="styles/main.css" rel="stylesheet" type="text/css"/>
</head>
<body>
    <h1>{project.title}</h1>
    {f'<p><em>{project.description}</em></p>' if project.description else ''}
    <p style="margin-top: 2em;">By {author}</p>
</body>
</html>'''
            with open(os.path.join(oebps, "title.xhtml"), 'w', encoding='utf-8') as f:
                f.write(title_content)

            manifest_items.append('<item id="title" href="title.xhtml" media-type="application/xhtml+xml"/>')
            spine_items.append('<itemref idref="title"/>')
            toc_items.append(f'<navPoint id="title" playOrder="1"><navLabel><text>Title Page</text></navLabel><content src="title.xhtml"/></navPoint>')

            # Chapter files
            for i, chapter in enumerate(chapters, start=2):
                chapter_title = chapter.title or f"Chapter {chapter.number}"
                filename = f"chapter_{chapter.number:03d}.xhtml"

                self.md.reset()
                if chapter.content:
                    chapter_content_html = self.md.convert(chapter.content)
                else:
                    chapter_content_html = "<p><em>[Chapter content not yet available]</em></p>"

                chapter_html = f'''<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" lang="en">
<head>
    <title>Chapter {chapter.number}: {chapter_title}</title>
    <link href="styles/main.css" rel="stylesheet" type="text/css"/>
</head>
<body>
    <div class="chapter-header">
        <p class="chapter-number">Chapter {chapter.number}</p>
        <h2 class="chapter-title">{chapter_title}</h2>
    </div>
    {chapter_content_html}
</body>
</html>'''

                with open(os.path.join(oebps, filename), 'w', encoding='utf-8') as f:
                    f.write(chapter_html)

                item_id = f"chapter{chapter.number}"
                manifest_items.append(f'<item id="{item_id}" href="{filename}" media-type="application/xhtml+xml"/>')
                spine_items.append(f'<itemref idref="{item_id}"/>')
                toc_items.append(f'<navPoint id="{item_id}" playOrder="{i}"><navLabel><text>Chapter {chapter.number}: {chapter_title}</text></navLabel><content src="{filename}"/></navPoint>')

            # content.opf
            now = datetime.utcnow().strftime("%Y-%m-%dT%H:%M:%SZ")
            content_opf = f'''<?xml version="1.0" encoding="UTF-8"?>
<package xmlns="http://www.idpf.org/2007/opf" version="3.0" unique-identifier="uid">
    <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
        <dc:identifier id="uid">{epub_uuid}</dc:identifier>
        <dc:title>{project.title}</dc:title>
        <dc:creator>{author}</dc:creator>
        <dc:language>en</dc:language>
        {f'<dc:description>{project.description}</dc:description>' if project.description else ''}
        {f'<dc:subject>{project.genre}</dc:subject>' if project.genre else ''}
        <meta property="dcterms:modified">{now}</meta>
    </metadata>
    <manifest>
        <item id="css" href="styles/main.css" media-type="text/css"/>
        <item id="nav" href="nav.xhtml" media-type="application/xhtml+xml" properties="nav"/>
        <item id="ncx" href="toc.ncx" media-type="application/x-dtbncx+xml"/>
        {"".join(manifest_items)}
    </manifest>
    <spine toc="ncx">
        {"".join(spine_items)}
    </spine>
</package>'''
            with open(os.path.join(oebps, "content.opf"), 'w', encoding='utf-8') as f:
                f.write(content_opf)

            # toc.ncx (for EPUB 2 compatibility)
            toc_ncx = f'''<?xml version="1.0" encoding="UTF-8"?>
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1">
    <head>
        <meta name="dtb:uid" content="{epub_uuid}"/>
        <meta name="dtb:depth" content="1"/>
        <meta name="dtb:totalPageCount" content="0"/>
        <meta name="dtb:maxPageNumber" content="0"/>
    </head>
    <docTitle><text>{project.title}</text></docTitle>
    <navMap>
        {"".join(toc_items)}
    </navMap>
</ncx>'''
            with open(os.path.join(oebps, "toc.ncx"), 'w', encoding='utf-8') as f:
                f.write(toc_ncx)

            # nav.xhtml (EPUB 3 navigation)
            nav_items = [f'<li><a href="title.xhtml">Title Page</a></li>']
            for chapter in chapters:
                chapter_title = chapter.title or f"Chapter {chapter.number}"
                nav_items.append(f'<li><a href="chapter_{chapter.number:03d}.xhtml">Chapter {chapter.number}: {chapter_title}</a></li>')

            nav_xhtml = f'''<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en">
<head>
    <title>Table of Contents</title>
    <link href="styles/main.css" rel="stylesheet" type="text/css"/>
</head>
<body>
    <nav epub:type="toc" id="toc">
        <h1>Table of Contents</h1>
        <ol>
            {"".join(nav_items)}
        </ol>
    </nav>
</body>
</html>'''
            with open(os.path.join(oebps, "nav.xhtml"), 'w', encoding='utf-8') as f:
                f.write(nav_xhtml)

            # Create EPUB (ZIP) file
            with zipfile.ZipFile(output_path, 'w', zipfile.ZIP_DEFLATED) as zf:
                # mimetype must be first and uncompressed
                zf.write(mimetype_path, "mimetype", compress_type=zipfile.ZIP_STORED)

                # Add all other files
                for root, dirs, files in os.walk(tmpdir):
                    for file in files:
                        if file == "mimetype":
                            continue
                        filepath = os.path.join(root, file)
                        arcname = os.path.relpath(filepath, tmpdir)
                        zf.write(filepath, arcname)

        return ExportResult(
            success=True,
            output_path=output_path,
            format='epub',
            message=f"Successfully exported '{project.title}' to EPUB (manual creation).",
            file_size=self._get_file_size(output_path)
        )

    # =========================================================================
    # PDF EXPORT (via markdown conversion)
    # =========================================================================

    def export_pdf(
        self,
        project_id: int,
        output_path: str
    ) -> ExportResult:
        """
        Export a project to PDF format via markdown conversion.

        Attempts to use available PDF conversion tools in order:
        1. pandoc (if available)
        2. weasyprint (if available)
        3. Falls back to HTML export with print instructions

        Args:
            project_id: ID of the project to export.
            output_path: Path for the output PDF file.

        Returns:
            ExportResult with success status and details.
        """
        session = self._get_session()
        try:
            project, chapters = self._get_project_with_chapters(session, project_id)

            if not project:
                return ExportResult(
                    success=False,
                    output_path=None,
                    format='pdf',
                    message=f"Project with ID {project_id} not found."
                )

            if not chapters:
                return ExportResult(
                    success=False,
                    output_path=None,
                    format='pdf',
                    message=f"Project '{project.title}' has no chapters to export."
                )

            self._ensure_output_dir(output_path)

            # Try pandoc first
            if self._check_command_available('pandoc'):
                return self._export_pdf_pandoc(project, chapters, output_path)

            # Try weasyprint
            try:
                import weasyprint
                return self._export_pdf_weasyprint(project, chapters, output_path)
            except ImportError:
                pass

            # Fallback: Create HTML and provide instructions
            html_path = output_path.replace('.pdf', '.html')
            html_result = self.export_html(project_id, html_path)

            if html_result.success:
                return ExportResult(
                    success=True,
                    output_path=html_path,
                    format='pdf',
                    message=f"PDF tools not available. Created HTML at '{html_path}'. "
                            f"Open in browser and use Print -> Save as PDF to create PDF.",
                    file_size=html_result.file_size
                )
            else:
                return html_result
        finally:
            session.close()

    def _check_command_available(self, command: str) -> bool:
        """Check if a command is available in PATH."""
        try:
            subprocess.run(
                [command, '--version'],
                capture_output=True,
                check=True
            )
            return True
        except (subprocess.CalledProcessError, FileNotFoundError):
            return False

    def _export_pdf_pandoc(
        self,
        project: Project,
        chapters: List[Chapter],
        output_path: str
    ) -> ExportResult:
        """Export to PDF using pandoc."""
        # Create temporary markdown file
        with tempfile.NamedTemporaryFile(
            mode='w',
            suffix='.md',
            delete=False,
            encoding='utf-8'
        ) as tmp:
            # Write content
            tmp.write(f"---\n")
            tmp.write(f"title: {project.title}\n")
            if project.description:
                tmp.write(f"subtitle: {project.description}\n")
            tmp.write(f"date: {datetime.now().strftime('%Y-%m-%d')}\n")
            tmp.write(f"---\n\n")

            for chapter in chapters:
                chapter_title = chapter.title or f"Chapter {chapter.number}"
                tmp.write(f"# Chapter {chapter.number}: {chapter_title}\n\n")
                if chapter.content:
                    tmp.write(chapter.content)
                else:
                    tmp.write("*[Chapter content not yet available]*")
                tmp.write("\n\n")

            tmp_path = tmp.name

        try:
            # Run pandoc
            result = subprocess.run(
                [
                    'pandoc',
                    tmp_path,
                    '-o', output_path,
                    '--pdf-engine=pdflatex',
                    '-V', 'geometry:margin=1in',
                    '-V', 'fontsize=12pt',
                    '--toc',
                    '--toc-depth=2'
                ],
                capture_output=True,
                text=True
            )

            if result.returncode == 0:
                return ExportResult(
                    success=True,
                    output_path=output_path,
                    format='pdf',
                    message=f"Successfully exported '{project.title}' to PDF using pandoc.",
                    file_size=self._get_file_size(output_path)
                )
            else:
                # Try without pdflatex
                result = subprocess.run(
                    [
                        'pandoc',
                        tmp_path,
                        '-o', output_path,
                        '--toc'
                    ],
                    capture_output=True,
                    text=True
                )

                if result.returncode == 0:
                    return ExportResult(
                        success=True,
                        output_path=output_path,
                        format='pdf',
                        message=f"Successfully exported '{project.title}' to PDF.",
                        file_size=self._get_file_size(output_path)
                    )
                else:
                    return ExportResult(
                        success=False,
                        output_path=None,
                        format='pdf',
                        message=f"Pandoc PDF conversion failed: {result.stderr}"
                    )
        finally:
            os.unlink(tmp_path)

    def _export_pdf_weasyprint(
        self,
        project: Project,
        chapters: List[Chapter],
        output_path: str
    ) -> ExportResult:
        """Export to PDF using weasyprint."""
        import weasyprint

        # First create HTML content
        session = self._get_session()
        try:
            # Create HTML in memory
            html_result = self.export_html(
                project.id,
                output_path.replace('.pdf', '.html')
            )

            if not html_result.success:
                return ExportResult(
                    success=False,
                    output_path=None,
                    format='pdf',
                    message=f"Failed to create HTML for PDF conversion: {html_result.message}"
                )

            # Convert HTML to PDF
            html_path = html_result.output_path
            html_doc = weasyprint.HTML(filename=html_path)
            html_doc.write_pdf(output_path)

            # Clean up temporary HTML
            os.unlink(html_path)

            return ExportResult(
                success=True,
                output_path=output_path,
                format='pdf',
                message=f"Successfully exported '{project.title}' to PDF using weasyprint.",
                file_size=self._get_file_size(output_path)
            )
        finally:
            session.close()

    # =========================================================================
    # AUDIO BUNDLE EXPORT
    # =========================================================================

    def export_audio_bundle(
        self,
        project_id: int,
        output_path: str,
        include_metadata: bool = True
    ) -> ExportResult:
        """
        Bundle all audio files for a project into a ZIP archive.

        Args:
            project_id: ID of the project to export.
            output_path: Path for the output ZIP file.
            include_metadata: Include a metadata JSON file in the bundle.

        Returns:
            ExportResult with success status and details.
        """
        session = self._get_session()
        try:
            project, chapters = self._get_project_with_chapters(session, project_id)

            if not project:
                return ExportResult(
                    success=False,
                    output_path=None,
                    format='audio_zip',
                    message=f"Project with ID {project_id} not found."
                )

            # Filter chapters with audio
            chapters_with_audio = [
                ch for ch in chapters
                if ch.audio_path and os.path.exists(ch.audio_path)
            ]

            if not chapters_with_audio:
                return ExportResult(
                    success=False,
                    output_path=None,
                    format='audio_zip',
                    message=f"Project '{project.title}' has no audio files to bundle."
                )

            self._ensure_output_dir(output_path)

            with zipfile.ZipFile(output_path, 'w', zipfile.ZIP_DEFLATED) as zf:
                # Add audio files
                for chapter in chapters_with_audio:
                    audio_path = chapter.audio_path
                    ext = os.path.splitext(audio_path)[1]
                    chapter_title = chapter.title or f"Chapter {chapter.number}"

                    # Create clean filename
                    safe_title = "".join(
                        c for c in chapter_title
                        if c.isalnum() or c in (' ', '-', '_')
                    ).strip()
                    filename = f"{chapter.number:03d}_{safe_title}{ext}"

                    zf.write(audio_path, filename)

                # Add metadata file
                if include_metadata:
                    import json

                    metadata = {
                        'project': {
                            'id': project.id,
                            'title': project.title,
                            'description': project.description,
                            'genre': project.genre,
                            'target_audience': project.target_audience,
                            'created_at': project.created_at.isoformat(),
                        },
                        'chapters': [
                            {
                                'number': ch.number,
                                'title': ch.title,
                                'audio_file': f"{ch.number:03d}_{(''.join(c for c in (ch.title or f'Chapter {ch.number}') if c.isalnum() or c in (' ', '-', '_'))).strip()}{os.path.splitext(ch.audio_path)[1]}",
                            }
                            for ch in chapters_with_audio
                        ],
                        'total_chapters': len(chapters_with_audio),
                        'exported_at': datetime.utcnow().isoformat(),
                    }

                    zf.writestr('metadata.json', json.dumps(metadata, indent=2))

            return ExportResult(
                success=True,
                output_path=output_path,
                format='audio_zip',
                message=f"Successfully bundled {len(chapters_with_audio)} audio files for '{project.title}'.",
                file_size=self._get_file_size(output_path)
            )
        finally:
            session.close()


# =============================================================================
# CONVENIENCE FUNCTIONS
# =============================================================================

def export_markdown(
    project_id: int,
    output_path: str,
    database_url: str = "sqlite:///bookforge.db",
    single_file: bool = True,
    **kwargs
) -> ExportResult:
    """
    Export a project to Markdown format.

    Args:
        project_id: ID of the project to export.
        output_path: Path for the output file(s).
        database_url: Database connection string.
        single_file: If True, export as single file.
        **kwargs: Additional arguments passed to BookExporter.export_markdown()

    Returns:
        ExportResult with success status and details.
    """
    exporter = BookExporter(database_url)
    return exporter.export_markdown(project_id, output_path, single_file=single_file, **kwargs)


def export_html(
    project_id: int,
    output_path: str,
    database_url: str = "sqlite:///bookforge.db",
    **kwargs
) -> ExportResult:
    """
    Export a project to styled HTML format.

    Args:
        project_id: ID of the project to export.
        output_path: Path for the output HTML file.
        database_url: Database connection string.
        **kwargs: Additional arguments passed to BookExporter.export_html()

    Returns:
        ExportResult with success status and details.
    """
    exporter = BookExporter(database_url)
    return exporter.export_html(project_id, output_path, **kwargs)


def export_epub(
    project_id: int,
    output_path: str,
    database_url: str = "sqlite:///bookforge.db",
    **kwargs
) -> ExportResult:
    """
    Export a project to EPUB format.

    Args:
        project_id: ID of the project to export.
        output_path: Path for the output EPUB file.
        database_url: Database connection string.
        **kwargs: Additional arguments passed to BookExporter.export_epub()

    Returns:
        ExportResult with success status and details.
    """
    exporter = BookExporter(database_url)
    return exporter.export_epub(project_id, output_path, **kwargs)


def export_pdf(
    project_id: int,
    output_path: str,
    database_url: str = "sqlite:///bookforge.db"
) -> ExportResult:
    """
    Export a project to PDF format.

    Args:
        project_id: ID of the project to export.
        output_path: Path for the output PDF file.
        database_url: Database connection string.

    Returns:
        ExportResult with success status and details.
    """
    exporter = BookExporter(database_url)
    return exporter.export_pdf(project_id, output_path)


def export_audio_bundle(
    project_id: int,
    output_path: str,
    database_url: str = "sqlite:///bookforge.db",
    **kwargs
) -> ExportResult:
    """
    Bundle all audio files for a project into a ZIP archive.

    Args:
        project_id: ID of the project to export.
        output_path: Path for the output ZIP file.
        database_url: Database connection string.
        **kwargs: Additional arguments passed to BookExporter.export_audio_bundle()

    Returns:
        ExportResult with success status and details.
    """
    exporter = BookExporter(database_url)
    return exporter.export_audio_bundle(project_id, output_path, **kwargs)


if __name__ == "__main__":
    # Example usage and testing
    import argparse

    parser = argparse.ArgumentParser(description="Export BookForge projects")
    parser.add_argument("project_id", type=int, help="Project ID to export")
    parser.add_argument("format", choices=['markdown', 'html', 'epub', 'pdf', 'audio'],
                       help="Export format")
    parser.add_argument("-o", "--output", required=True, help="Output path")
    parser.add_argument("--db", default="sqlite:///bookforge.db",
                       help="Database URL")
    parser.add_argument("--multi-file", action="store_true",
                       help="For markdown: export as multiple files")

    args = parser.parse_args()

    exporter = BookExporter(args.db)

    if args.format == 'markdown':
        result = exporter.export_markdown(
            args.project_id,
            args.output,
            single_file=not args.multi_file
        )
    elif args.format == 'html':
        result = exporter.export_html(args.project_id, args.output)
    elif args.format == 'epub':
        result = exporter.export_epub(args.project_id, args.output)
    elif args.format == 'pdf':
        result = exporter.export_pdf(args.project_id, args.output)
    elif args.format == 'audio':
        result = exporter.export_audio_bundle(args.project_id, args.output)

    print(f"Success: {result.success}")
    print(f"Message: {result.message}")
    if result.output_path:
        print(f"Output: {result.output_path}")
    if result.file_size:
        print(f"Size: {result.file_size:,} bytes")
