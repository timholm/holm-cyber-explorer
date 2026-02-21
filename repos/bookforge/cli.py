"""
BookForge Command Line Interface

Provides CLI commands for managing the BookForge book generation system.
"""

import os
import sys
import click
from pathlib import Path
from datetime import datetime

# Add parent directory to path for imports
sys.path.insert(0, str(Path(__file__).parent.parent))


def get_db_session():
    """Get database session with proper error handling."""
    from bookforge.models import init_db
    db_path = os.environ.get('BOOKFORGE_DB', 'sqlite:///bookforge.db')
    engine, session = init_db(db_path)
    return session


def format_timestamp(dt):
    """Format datetime for display."""
    if dt is None:
        return 'N/A'
    return dt.strftime('%Y-%m-%d %H:%M:%S')


def format_status(status):
    """Format status enum for display with color."""
    status_colors = {
        'draft': 'yellow',
        'outlining': 'cyan',
        'writing': 'blue',
        'editing': 'magenta',
        'completed': 'green',
        'archived': 'white',
        'pending': 'yellow',
        'generating': 'cyan',
        'revised': 'blue',
        'final': 'green',
        'running': 'cyan',
        'failed': 'red',
        'cancelled': 'white',
    }
    value = status.value if hasattr(status, 'value') else str(status)
    color = status_colors.get(value, 'white')
    return click.style(value.upper(), fg=color)


@click.group()
@click.version_option(version='0.1.0', prog_name='BookForge')
@click.pass_context
def cli(ctx):
    """
    BookForge - AI-Powered Book Generation System

    Generate, manage, and convert books using AI with support for
    multiple export formats and text-to-speech audio generation.
    """
    ctx.ensure_object(dict)


@cli.command()
@click.option('--host', '-h', default='0.0.0.0', help='Host to bind to')
@click.option('--port', '-p', default=5000, type=int, help='Port to bind to')
@click.option('--debug', '-d', is_flag=True, help='Enable debug mode')
@click.option('--workers', '-w', default=4, type=int, help='Number of gunicorn workers')
def serve(host, port, debug, workers):
    """Start the BookForge web server."""
    click.echo(click.style('Starting BookForge Web Server...', fg='green', bold=True))
    click.echo(f'  Host: {host}')
    click.echo(f'  Port: {port}')
    click.echo(f'  Debug: {debug}')

    if debug:
        # Development mode with Flask's built-in server
        click.echo(click.style('Running in development mode', fg='yellow'))
        try:
            from bookforge.app import app
            app.run(host=host, port=port, debug=True)
        except ImportError as e:
            click.echo(click.style(f'Error: Could not import app module: {e}', fg='red'))
            click.echo('Make sure app.py exists in the bookforge directory.')
            sys.exit(1)
    else:
        # Production mode with gunicorn
        click.echo(f'  Workers: {workers}')
        click.echo(click.style('Running in production mode with Gunicorn', fg='cyan'))
        import subprocess
        cmd = [
            sys.executable.replace('python', 'gunicorn'),
            '--bind', f'{host}:{port}',
            '--workers', str(workers),
            '--timeout', '120',
            '--access-logfile', '-',
            'bookforge.app:app'
        ]
        # Try gunicorn from venv
        venv_gunicorn = Path(__file__).parent / 'venv' / 'bin' / 'gunicorn'
        if venv_gunicorn.exists():
            cmd[0] = str(venv_gunicorn)

        try:
            subprocess.run(cmd, cwd=str(Path(__file__).parent))
        except FileNotFoundError:
            click.echo(click.style('Error: gunicorn not found. Install with: pip install gunicorn', fg='red'))
            sys.exit(1)


@cli.command()
@click.option('--poll-interval', '-i', default=5, type=int, help='Polling interval in seconds')
@click.option('--max-concurrent', '-c', default=2, type=int, help='Max concurrent jobs')
def worker(poll_interval, max_concurrent):
    """Start the BookForge background worker."""
    click.echo(click.style('Starting BookForge Background Worker...', fg='green', bold=True))
    click.echo(f'  Poll Interval: {poll_interval}s')
    click.echo(f'  Max Concurrent Jobs: {max_concurrent}')

    try:
        # Try to import and run the worker
        worker_path = Path(__file__).parent / 'worker.py'
        if worker_path.exists():
            import subprocess
            env = os.environ.copy()
            env['WORKER_POLL_INTERVAL'] = str(poll_interval)
            env['WORKER_MAX_CONCURRENT'] = str(max_concurrent)
            subprocess.run([sys.executable, str(worker_path)], env=env)
        else:
            # Inline worker implementation
            click.echo(click.style('Worker module not found, running inline worker...', fg='yellow'))
            run_inline_worker(poll_interval, max_concurrent)
    except KeyboardInterrupt:
        click.echo(click.style('\nWorker stopped.', fg='yellow'))


def run_inline_worker(poll_interval, max_concurrent):
    """Run a basic inline worker when worker.py is not available."""
    import time
    from bookforge.models import (
        init_db, GenerationJob, JobStatus, JobType
    )

    db_path = os.environ.get('BOOKFORGE_DB', 'sqlite:///bookforge.db')
    engine, session = init_db(db_path)

    click.echo('Worker running. Press Ctrl+C to stop.')

    while True:
        try:
            # Check for pending jobs
            pending_jobs = session.query(GenerationJob).filter(
                GenerationJob.status == JobStatus.PENDING
            ).limit(max_concurrent).all()

            if pending_jobs:
                for job in pending_jobs:
                    click.echo(f'Processing job {job.id} (type: {job.job_type.value})')
                    job.status = JobStatus.RUNNING
                    job.started_at = datetime.utcnow()
                    session.commit()

                    # Placeholder for actual job processing
                    click.echo(f'  Job {job.id} would be processed here')

            time.sleep(poll_interval)

        except KeyboardInterrupt:
            raise
        except Exception as e:
            click.echo(click.style(f'Worker error: {e}', fg='red'))
            time.sleep(poll_interval)


@cli.command()
@click.argument('title')
@click.option('--description', '-d', default='', help='Project description')
@click.option('--genre', '-g', default='general', help='Book genre')
@click.option('--audience', '-a', default='general readers', help='Target audience')
def create(title, description, genre, audience):
    """Create a new book project."""
    from bookforge.models import Project, ProjectStatus

    session = get_db_session()

    # Create the project
    project = Project(
        title=title,
        description=description,
        genre=genre,
        target_audience=audience,
        status=ProjectStatus.DRAFT
    )

    session.add(project)
    session.commit()

    click.echo(click.style('Project created successfully!', fg='green', bold=True))
    click.echo(f'  ID: {project.id}')
    click.echo(f'  Title: {project.title}')
    click.echo(f'  Genre: {project.genre}')
    click.echo(f'  Audience: {project.target_audience}')
    click.echo(f'  Status: {format_status(project.status)}')

    session.close()


@cli.command()
@click.argument('project_id', type=int)
@click.option('--chapters', '-c', default=None, type=str, help='Specific chapters (e.g., "1,3,5" or "1-5")')
@click.option('--force', '-f', is_flag=True, help='Force regeneration of existing chapters')
def generate(project_id, chapters, force):
    """Generate all chapters for a project."""
    from bookforge.models import (
        Project, Chapter, GenerationJob, JobType, JobStatus, ChapterStatus
    )

    session = get_db_session()

    # Find the project
    project = session.query(Project).filter(Project.id == project_id).first()
    if not project:
        click.echo(click.style(f'Error: Project {project_id} not found', fg='red'))
        sys.exit(1)

    click.echo(click.style(f'Generating chapters for: {project.title}', fg='green', bold=True))

    # Parse chapter selection
    chapter_nums = None
    if chapters:
        chapter_nums = parse_chapter_selection(chapters)
        click.echo(f'  Selected chapters: {chapter_nums}')

    # Get chapters to generate
    query = session.query(Chapter).filter(Chapter.project_id == project_id)
    if chapter_nums:
        query = query.filter(Chapter.number.in_(chapter_nums))
    if not force:
        query = query.filter(Chapter.status.in_([ChapterStatus.PENDING, ChapterStatus.DRAFT]))

    chapters_to_gen = query.order_by(Chapter.number).all()

    if not chapters_to_gen:
        # Create generation jobs for outline first if no chapters exist
        click.echo('  No chapters found. Creating outline generation job...')
        job = GenerationJob(
            project_id=project_id,
            job_type=JobType.OUTLINE,
            status=JobStatus.PENDING
        )
        session.add(job)
        session.commit()
        click.echo(f'  Created outline job: {job.id}')
    else:
        # Create chapter generation jobs
        jobs_created = 0
        for chapter in chapters_to_gen:
            job = GenerationJob(
                project_id=project_id,
                chapter_id=chapter.id,
                job_type=JobType.CHAPTER,
                status=JobStatus.PENDING
            )
            session.add(job)
            jobs_created += 1

        session.commit()
        click.echo(f'  Created {jobs_created} chapter generation job(s)')

    click.echo(click.style('Jobs queued. Start the worker to process them.', fg='cyan'))
    session.close()


def parse_chapter_selection(selection):
    """Parse chapter selection string like '1,3,5' or '1-5'."""
    chapters = set()
    parts = selection.split(',')
    for part in parts:
        part = part.strip()
        if '-' in part:
            start, end = part.split('-')
            chapters.update(range(int(start), int(end) + 1))
        else:
            chapters.add(int(part))
    return sorted(chapters)


@cli.command()
@click.argument('project_id', type=int)
@click.option('--chapters', '-c', default=None, type=str, help='Specific chapters (e.g., "1,3,5" or "1-5")')
@click.option('--voice', '-v', default='en_US-lessac-medium', help='Voice model to use')
@click.option('--output-dir', '-o', default=None, help='Output directory for audio files')
def audio(project_id, chapters, voice, output_dir):
    """Convert project chapters to audio using TTS."""
    from bookforge.models import (
        Project, Chapter, GenerationJob, JobType, JobStatus
    )

    session = get_db_session()

    # Find the project
    project = session.query(Project).filter(Project.id == project_id).first()
    if not project:
        click.echo(click.style(f'Error: Project {project_id} not found', fg='red'))
        sys.exit(1)

    click.echo(click.style(f'Converting to audio: {project.title}', fg='green', bold=True))
    click.echo(f'  Voice: {voice}')

    # Parse chapter selection
    chapter_nums = None
    if chapters:
        chapter_nums = parse_chapter_selection(chapters)
        click.echo(f'  Selected chapters: {chapter_nums}')

    # Get chapters to convert
    query = session.query(Chapter).filter(Chapter.project_id == project_id)
    if chapter_nums:
        query = query.filter(Chapter.number.in_(chapter_nums))

    chapters_to_convert = query.filter(Chapter.content.isnot(None)).order_by(Chapter.number).all()

    if not chapters_to_convert:
        click.echo(click.style('No chapters with content found to convert.', fg='yellow'))
        session.close()
        return

    # Create audio generation jobs
    jobs_created = 0
    for chapter in chapters_to_convert:
        job = GenerationJob(
            project_id=project_id,
            chapter_id=chapter.id,
            job_type=JobType.AUDIO,
            status=JobStatus.PENDING
        )
        session.add(job)
        jobs_created += 1

    session.commit()
    click.echo(f'  Created {jobs_created} audio generation job(s)')
    click.echo(click.style('Jobs queued. Start the worker to process them.', fg='cyan'))

    session.close()


@cli.command('export')
@click.argument('project_id', type=int)
@click.argument('format', type=click.Choice(['pdf', 'epub', 'mobi', 'html', 'markdown', 'docx']))
@click.option('--output', '-o', default=None, help='Output file path')
@click.option('--include-toc', '-t', is_flag=True, default=True, help='Include table of contents')
def export_project(project_id, format, output, include_toc):
    """Export a project to the specified format."""
    from bookforge.models import Project, Chapter

    session = get_db_session()

    # Find the project
    project = session.query(Project).filter(Project.id == project_id).first()
    if not project:
        click.echo(click.style(f'Error: Project {project_id} not found', fg='red'))
        sys.exit(1)

    click.echo(click.style(f'Exporting: {project.title}', fg='green', bold=True))
    click.echo(f'  Format: {format.upper()}')

    # Get chapters
    chapters = session.query(Chapter).filter(
        Chapter.project_id == project_id,
        Chapter.content.isnot(None)
    ).order_by(Chapter.number).all()

    if not chapters:
        click.echo(click.style('No chapters with content found to export.', fg='yellow'))
        session.close()
        return

    click.echo(f'  Chapters: {len(chapters)}')

    # Determine output path
    if output is None:
        safe_title = ''.join(c if c.isalnum() or c in ' -_' else '' for c in project.title)
        safe_title = safe_title.replace(' ', '_').lower()
        output = f'{safe_title}.{format}'

    click.echo(f'  Output: {output}')

    # Export based on format
    try:
        if format == 'markdown':
            export_markdown(project, chapters, output, include_toc)
        elif format == 'html':
            export_html(project, chapters, output, include_toc)
        elif format == 'epub':
            export_epub(project, chapters, output, include_toc)
        else:
            click.echo(click.style(f'Export to {format} not yet implemented.', fg='yellow'))
            click.echo('Supported formats: markdown, html, epub')
            session.close()
            return

        click.echo(click.style(f'Successfully exported to: {output}', fg='green'))
    except Exception as e:
        click.echo(click.style(f'Export failed: {e}', fg='red'))

    session.close()


def export_markdown(project, chapters, output_path, include_toc):
    """Export project to Markdown format."""
    with open(output_path, 'w', encoding='utf-8') as f:
        # Title
        f.write(f'# {project.title}\n\n')

        if project.description:
            f.write(f'{project.description}\n\n')

        f.write('---\n\n')

        # Table of contents
        if include_toc:
            f.write('## Table of Contents\n\n')
            for chapter in chapters:
                title = chapter.title or f'Chapter {chapter.number}'
                anchor = title.lower().replace(' ', '-')
                f.write(f'- [{title}](#{anchor})\n')
            f.write('\n---\n\n')

        # Chapters
        for chapter in chapters:
            title = chapter.title or f'Chapter {chapter.number}'
            f.write(f'## {title}\n\n')
            f.write(f'{chapter.content}\n\n')
            f.write('---\n\n')


def export_html(project, chapters, output_path, include_toc):
    """Export project to HTML format."""
    import markdown

    md = markdown.Markdown(extensions=['extra', 'toc'])

    with open(output_path, 'w', encoding='utf-8') as f:
        f.write('<!DOCTYPE html>\n<html>\n<head>\n')
        f.write(f'<title>{project.title}</title>\n')
        f.write('<meta charset="utf-8">\n')
        f.write('<style>\n')
        f.write('body { max-width: 800px; margin: 0 auto; padding: 20px; font-family: Georgia, serif; }\n')
        f.write('h1 { text-align: center; }\n')
        f.write('h2 { margin-top: 2em; border-bottom: 1px solid #ccc; }\n')
        f.write('.toc { background: #f5f5f5; padding: 1em; border-radius: 5px; }\n')
        f.write('</style>\n')
        f.write('</head>\n<body>\n')

        f.write(f'<h1>{project.title}</h1>\n')

        if project.description:
            f.write(f'<p><em>{project.description}</em></p>\n')

        # Table of contents
        if include_toc:
            f.write('<div class="toc">\n<h3>Table of Contents</h3>\n<ul>\n')
            for chapter in chapters:
                title = chapter.title or f'Chapter {chapter.number}'
                anchor = f'chapter-{chapter.number}'
                f.write(f'<li><a href="#{anchor}">{title}</a></li>\n')
            f.write('</ul>\n</div>\n')

        # Chapters
        for chapter in chapters:
            title = chapter.title or f'Chapter {chapter.number}'
            anchor = f'chapter-{chapter.number}'
            f.write(f'<h2 id="{anchor}">{title}</h2>\n')
            content_html = md.convert(chapter.content)
            f.write(content_html)
            f.write('\n')
            md.reset()

        f.write('</body>\n</html>')


def export_epub(project, chapters, output_path, include_toc):
    """Export project to EPUB format."""
    try:
        from ebooklib import epub
    except ImportError:
        raise click.ClickException('ebooklib not installed. Run: pip install ebooklib')

    import markdown
    md = markdown.Markdown(extensions=['extra'])

    book = epub.EpubBook()
    book.set_identifier(f'bookforge-{project.id}')
    book.set_title(project.title)
    book.set_language('en')

    if project.description:
        book.add_metadata('DC', 'description', project.description)

    # Add chapters
    epub_chapters = []
    for chapter in chapters:
        title = chapter.title or f'Chapter {chapter.number}'
        c = epub.EpubHtml(title=title, file_name=f'ch{chapter.number:02d}.xhtml')
        content_html = md.convert(chapter.content)
        c.content = f'<h1>{title}</h1>{content_html}'
        book.add_item(c)
        epub_chapters.append(c)
        md.reset()

    # Add navigation
    book.toc = epub_chapters
    book.add_item(epub.EpubNcx())
    book.add_item(epub.EpubNav())

    # Spine
    book.spine = ['nav'] + epub_chapters

    # Write the book
    epub.write_epub(output_path, book)


@cli.command('list')
@click.option('--status', '-s', default=None, help='Filter by status')
@click.option('--limit', '-l', default=20, type=int, help='Maximum number of projects to show')
def list_projects(status, limit):
    """List all projects."""
    from bookforge.models import Project, ProjectStatus

    session = get_db_session()

    query = session.query(Project)

    if status:
        try:
            status_enum = ProjectStatus(status.lower())
            query = query.filter(Project.status == status_enum)
        except ValueError:
            click.echo(click.style(f'Invalid status: {status}', fg='red'))
            click.echo(f'Valid statuses: {[s.value for s in ProjectStatus]}')
            session.close()
            return

    projects = query.order_by(Project.updated_at.desc()).limit(limit).all()

    if not projects:
        click.echo(click.style('No projects found.', fg='yellow'))
        session.close()
        return

    click.echo(click.style(f'Projects ({len(projects)} found):', fg='green', bold=True))
    click.echo()

    # Header
    click.echo(f'{"ID":>4}  {"Status":<12}  {"Chapters":>8}  {"Updated":<20}  Title')
    click.echo('-' * 80)

    for project in projects:
        chapter_count = len(project.chapters)
        updated = format_timestamp(project.updated_at)
        status_str = format_status(project.status)
        click.echo(f'{project.id:>4}  {status_str:<20}  {chapter_count:>8}  {updated:<20}  {project.title[:30]}')

    session.close()


@cli.command()
@click.option('--verbose', '-v', is_flag=True, help='Show detailed status')
def status(verbose):
    """Show system status and statistics."""
    from bookforge.models import (
        Project, Chapter, GenerationJob,
        ProjectStatus, ChapterStatus, JobStatus
    )

    session = get_db_session()

    click.echo(click.style('BookForge System Status', fg='green', bold=True))
    click.echo('=' * 50)

    # Project statistics
    total_projects = session.query(Project).count()
    projects_by_status = {}
    for s in ProjectStatus:
        count = session.query(Project).filter(Project.status == s).count()
        if count > 0:
            projects_by_status[s.value] = count

    click.echo(f'\nProjects: {total_projects}')
    if verbose and projects_by_status:
        for s, count in projects_by_status.items():
            click.echo(f'  - {s}: {count}')

    # Chapter statistics
    total_chapters = session.query(Chapter).count()
    total_words = session.query(Chapter).with_entities(
        Chapter.word_count
    ).all()
    total_word_count = sum(w[0] or 0 for w in total_words)

    click.echo(f'\nChapters: {total_chapters}')
    click.echo(f'Total Words: {total_word_count:,}')

    if verbose:
        chapters_by_status = {}
        for s in ChapterStatus:
            count = session.query(Chapter).filter(Chapter.status == s).count()
            if count > 0:
                chapters_by_status[s.value] = count
        if chapters_by_status:
            for s, count in chapters_by_status.items():
                click.echo(f'  - {s}: {count}')

    # Job statistics
    pending_jobs = session.query(GenerationJob).filter(
        GenerationJob.status == JobStatus.PENDING
    ).count()
    running_jobs = session.query(GenerationJob).filter(
        GenerationJob.status == JobStatus.RUNNING
    ).count()
    failed_jobs = session.query(GenerationJob).filter(
        GenerationJob.status == JobStatus.FAILED
    ).count()

    click.echo(f'\nGeneration Jobs:')
    click.echo(f'  - Pending: {pending_jobs}')
    click.echo(f'  - Running: {running_jobs}')
    if failed_jobs > 0:
        click.echo(click.style(f'  - Failed: {failed_jobs}', fg='red'))

    # Database info
    if verbose:
        db_path = os.environ.get('BOOKFORGE_DB', 'sqlite:///bookforge.db')
        click.echo(f'\nDatabase: {db_path}')

    click.echo()
    session.close()


@cli.command()
@click.option('--force', '-f', is_flag=True, help='Force initialization even if tables exist')
def init(force):
    """Initialize the database."""
    from bookforge.models import init_db, Base, get_engine

    db_path = os.environ.get('BOOKFORGE_DB', 'sqlite:///bookforge.db')

    click.echo(click.style('Initializing BookForge database...', fg='green', bold=True))
    click.echo(f'  Database: {db_path}')

    try:
        engine, session = init_db(db_path)
        tables = list(Base.metadata.tables.keys())
        click.echo(f'  Tables created: {len(tables)}')
        if tables:
            for table in tables:
                click.echo(f'    - {table}')
        click.echo(click.style('Database initialized successfully!', fg='green'))
        session.close()
    except Exception as e:
        click.echo(click.style(f'Error initializing database: {e}', fg='red'))
        sys.exit(1)


@cli.command()
@click.argument('project_id', type=int)
@click.option('--verbose', '-v', is_flag=True, help='Show full chapter content')
def show(project_id, verbose):
    """Show details of a specific project."""
    from bookforge.models import Project, Chapter

    session = get_db_session()

    project = session.query(Project).filter(Project.id == project_id).first()
    if not project:
        click.echo(click.style(f'Error: Project {project_id} not found', fg='red'))
        session.close()
        sys.exit(1)

    click.echo(click.style(f'Project: {project.title}', fg='green', bold=True))
    click.echo('=' * 50)
    click.echo(f'ID: {project.id}')
    click.echo(f'Status: {format_status(project.status)}')
    click.echo(f'Genre: {project.genre or "N/A"}')
    click.echo(f'Audience: {project.target_audience or "N/A"}')
    click.echo(f'Created: {format_timestamp(project.created_at)}')
    click.echo(f'Updated: {format_timestamp(project.updated_at)}')

    if project.description:
        click.echo(f'\nDescription:\n{project.description}')

    # Show chapters
    chapters = session.query(Chapter).filter(
        Chapter.project_id == project_id
    ).order_by(Chapter.number).all()

    if chapters:
        click.echo(f'\nChapters ({len(chapters)}):')
        click.echo('-' * 50)
        for chapter in chapters:
            title = chapter.title or f'Chapter {chapter.number}'
            status_str = format_status(chapter.status)
            words = f'{chapter.word_count:,} words' if chapter.word_count else 'no content'
            audio_marker = ' [audio]' if chapter.audio_path else ''
            click.echo(f'  {chapter.number:2}. {title:<30} {status_str} ({words}){audio_marker}')

            if verbose and chapter.content:
                click.echo(f'\n{chapter.content[:500]}...')
                click.echo()
    else:
        click.echo(click.style('\nNo chapters yet.', fg='yellow'))

    session.close()


def main():
    """Main entry point for the CLI."""
    cli()


if __name__ == '__main__':
    main()
