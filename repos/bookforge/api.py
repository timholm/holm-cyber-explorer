"""
REST API routes for BookForge book generation system.
Flask Blueprint with JSON responses and error handling.
"""

import os
import json
import uuid
import subprocess
import threading
from datetime import datetime
from functools import wraps
from flask import Blueprint, jsonify, request, send_file, current_app
from pathlib import Path

api = Blueprint('api', __name__, url_prefix='/api')

# In-memory storage for projects and jobs (in production, use a database)
projects = {}
jobs = {}

# Base directory for project data
DATA_DIR = Path('/home/tim/bookforge/data')
DATA_DIR.mkdir(exist_ok=True)


def load_projects():
    """Load projects from disk."""
    global projects
    projects_file = DATA_DIR / 'projects.json'
    if projects_file.exists():
        with open(projects_file, 'r') as f:
            projects = json.load(f)
    return projects


def save_projects():
    """Save projects to disk."""
    projects_file = DATA_DIR / 'projects.json'
    with open(projects_file, 'w') as f:
        json.dump(projects, f, indent=2, default=str)


def api_response(data=None, message=None, status=200):
    """Create a standardized API response."""
    response = {
        'success': 200 <= status < 300,
        'timestamp': datetime.utcnow().isoformat()
    }
    if data is not None:
        response['data'] = data
    if message:
        response['message'] = message
    return jsonify(response), status


def error_response(message, status=400, errors=None):
    """Create a standardized error response."""
    response = {
        'success': False,
        'error': message,
        'timestamp': datetime.utcnow().isoformat()
    }
    if errors:
        response['errors'] = errors
    return jsonify(response), status


def validate_project_exists(f):
    """Decorator to validate project exists."""
    @wraps(f)
    def decorated_function(project_id, *args, **kwargs):
        load_projects()
        if project_id not in projects:
            return error_response(f'Project {project_id} not found', 404)
        return f(project_id, *args, **kwargs)
    return decorated_function


def run_async_job(job_id, func, *args, **kwargs):
    """Run a function asynchronously and track job status."""
    def wrapper():
        try:
            jobs[job_id]['status'] = 'running'
            jobs[job_id]['started_at'] = datetime.utcnow().isoformat()
            result = func(*args, **kwargs)
            jobs[job_id]['status'] = 'completed'
            jobs[job_id]['result'] = result
            jobs[job_id]['completed_at'] = datetime.utcnow().isoformat()
        except Exception as e:
            jobs[job_id]['status'] = 'failed'
            jobs[job_id]['error'] = str(e)
            jobs[job_id]['completed_at'] = datetime.utcnow().isoformat()

    thread = threading.Thread(target=wrapper)
    thread.start()
    return thread


# ============================================================================
# Project Routes
# ============================================================================

@api.route('/projects', methods=['GET'])
def list_projects():
    """List all projects."""
    load_projects()
    project_list = []
    for pid, proj in projects.items():
        project_list.append({
            'id': pid,
            'title': proj.get('title', 'Untitled'),
            'genre': proj.get('genre', ''),
            'chapter_count': len(proj.get('chapters', [])),
            'created_at': proj.get('created_at'),
            'updated_at': proj.get('updated_at'),
            'status': proj.get('status', 'draft')
        })
    return api_response(data={'projects': project_list})


@api.route('/projects', methods=['POST'])
def create_project():
    """Create a new project."""
    data = request.get_json()
    if not data:
        return error_response('Request body is required', 400)

    required_fields = ['title']
    missing = [f for f in required_fields if f not in data]
    if missing:
        return error_response(f'Missing required fields: {", ".join(missing)}', 400)

    project_id = str(uuid.uuid4())[:8]
    now = datetime.utcnow().isoformat()

    project = {
        'id': project_id,
        'title': data['title'],
        'genre': data.get('genre', 'fiction'),
        'description': data.get('description', ''),
        'target_audience': data.get('target_audience', 'general'),
        'num_chapters': data.get('num_chapters', 10),
        'words_per_chapter': data.get('words_per_chapter', 3000),
        'model': data.get('model', 'llama3.2'),
        'voice': data.get('voice', 'default'),
        'outline': None,
        'chapters': [],
        'status': 'draft',
        'created_at': now,
        'updated_at': now
    }

    load_projects()
    projects[project_id] = project
    save_projects()

    # Create project directory
    project_dir = DATA_DIR / project_id
    project_dir.mkdir(exist_ok=True)
    (project_dir / 'chapters').mkdir(exist_ok=True)
    (project_dir / 'audio').mkdir(exist_ok=True)

    return api_response(data={'project': project}, message='Project created successfully', status=201)


@api.route('/projects/<project_id>', methods=['GET'])
@validate_project_exists
def get_project(project_id):
    """Get a specific project."""
    return api_response(data={'project': projects[project_id]})


@api.route('/projects/<project_id>', methods=['DELETE'])
@validate_project_exists
def delete_project(project_id):
    """Delete a project."""
    import shutil

    # Remove project directory
    project_dir = DATA_DIR / project_id
    if project_dir.exists():
        shutil.rmtree(project_dir)

    del projects[project_id]
    save_projects()

    return api_response(message=f'Project {project_id} deleted successfully')


# ============================================================================
# Chapter Routes
# ============================================================================

@api.route('/projects/<project_id>/chapters', methods=['GET'])
@validate_project_exists
def list_chapters(project_id):
    """List all chapters for a project."""
    project = projects[project_id]
    chapters = project.get('chapters', [])

    chapter_list = []
    for i, chapter in enumerate(chapters):
        chapter_list.append({
            'number': i + 1,
            'title': chapter.get('title', f'Chapter {i + 1}'),
            'word_count': len(chapter.get('content', '').split()),
            'has_audio': chapter.get('has_audio', False),
            'generated_at': chapter.get('generated_at')
        })

    return api_response(data={
        'chapters': chapter_list,
        'total': len(chapter_list)
    })


@api.route('/projects/<project_id>/chapters/<int:chapter_num>', methods=['GET'])
@validate_project_exists
def get_chapter(project_id, chapter_num):
    """Get a specific chapter."""
    project = projects[project_id]
    chapters = project.get('chapters', [])

    if chapter_num < 1 or chapter_num > len(chapters):
        return error_response(f'Chapter {chapter_num} not found', 404)

    chapter = chapters[chapter_num - 1]
    return api_response(data={
        'chapter': {
            'number': chapter_num,
            'title': chapter.get('title', f'Chapter {chapter_num}'),
            'content': chapter.get('content', ''),
            'word_count': len(chapter.get('content', '').split()),
            'has_audio': chapter.get('has_audio', False),
            'audio_path': chapter.get('audio_path'),
            'generated_at': chapter.get('generated_at')
        }
    })


# ============================================================================
# Generation Routes
# ============================================================================

def generate_outline_task(project_id):
    """Background task to generate book outline using Ollama."""
    project = projects[project_id]

    prompt = f"""Create a detailed outline for a {project['genre']} book titled "{project['title']}".

Description: {project.get('description', 'No description provided')}
Target Audience: {project.get('target_audience', 'general')}
Number of Chapters: {project.get('num_chapters', 10)}

For each chapter, provide:
1. Chapter title
2. Brief summary (2-3 sentences)
3. Key events or topics covered

Format as JSON with structure:
{{
  "chapters": [
    {{
      "number": 1,
      "title": "Chapter Title",
      "summary": "Brief chapter summary",
      "key_points": ["point1", "point2"]
    }}
  ]
}}"""

    try:
        result = subprocess.run(
            ['ollama', 'run', project.get('model', 'llama3.2'), prompt],
            capture_output=True,
            text=True,
            timeout=300
        )

        if result.returncode == 0:
            # Try to parse JSON from response
            response_text = result.stdout.strip()
            try:
                # Find JSON in response
                json_start = response_text.find('{')
                json_end = response_text.rfind('}') + 1
                if json_start >= 0 and json_end > json_start:
                    outline = json.loads(response_text[json_start:json_end])
                else:
                    outline = {'raw_outline': response_text}
            except json.JSONDecodeError:
                outline = {'raw_outline': response_text}

            projects[project_id]['outline'] = outline
            projects[project_id]['status'] = 'outlined'
            projects[project_id]['updated_at'] = datetime.utcnow().isoformat()
            save_projects()
            return {'outline': outline}
        else:
            raise Exception(f'Ollama error: {result.stderr}')
    except subprocess.TimeoutExpired:
        raise Exception('Outline generation timed out')


def generate_chapter_task(project_id, chapter_num):
    """Background task to generate a chapter using Ollama."""
    project = projects[project_id]
    outline = project.get('outline', {})

    chapter_outline = None
    if 'chapters' in outline and len(outline['chapters']) >= chapter_num:
        chapter_outline = outline['chapters'][chapter_num - 1]

    prompt = f"""Write Chapter {chapter_num} of a {project['genre']} book titled "{project['title']}".

{f'Chapter Title: {chapter_outline["title"]}' if chapter_outline else ''}
{f'Chapter Summary: {chapter_outline.get("summary", "")}' if chapter_outline else ''}
{f'Key Points: {", ".join(chapter_outline.get("key_points", []))}' if chapter_outline else ''}

Target word count: {project.get('words_per_chapter', 3000)} words
Target audience: {project.get('target_audience', 'general')}

Write engaging, well-paced prose. Include dialogue, description, and character development as appropriate."""

    try:
        result = subprocess.run(
            ['ollama', 'run', project.get('model', 'llama3.2'), prompt],
            capture_output=True,
            text=True,
            timeout=600
        )

        if result.returncode == 0:
            content = result.stdout.strip()

            # Ensure chapters list is long enough
            while len(projects[project_id]['chapters']) < chapter_num:
                projects[project_id]['chapters'].append({})

            projects[project_id]['chapters'][chapter_num - 1] = {
                'title': chapter_outline['title'] if chapter_outline else f'Chapter {chapter_num}',
                'content': content,
                'has_audio': False,
                'generated_at': datetime.utcnow().isoformat()
            }

            # Save chapter to file
            chapter_file = DATA_DIR / project_id / 'chapters' / f'chapter_{chapter_num:02d}.txt'
            with open(chapter_file, 'w') as f:
                f.write(content)

            projects[project_id]['updated_at'] = datetime.utcnow().isoformat()
            save_projects()

            return {'chapter_num': chapter_num, 'word_count': len(content.split())}
        else:
            raise Exception(f'Ollama error: {result.stderr}')
    except subprocess.TimeoutExpired:
        raise Exception(f'Chapter {chapter_num} generation timed out')


def generate_all_chapters_task(project_id):
    """Background task to generate all chapters."""
    project = projects[project_id]
    num_chapters = project.get('num_chapters', 10)
    results = []

    for i in range(1, num_chapters + 1):
        result = generate_chapter_task(project_id, i)
        results.append(result)
        # Update job progress
        job_id = f'gen_all_{project_id}'
        if job_id in jobs:
            jobs[job_id]['progress'] = i / num_chapters

    projects[project_id]['status'] = 'completed'
    save_projects()
    return {'chapters_generated': len(results), 'results': results}


def convert_to_audio_task(project_id):
    """Background task to convert all chapters to audio."""
    project = projects[project_id]
    chapters = project.get('chapters', [])
    voice = project.get('voice', 'default')
    results = []

    for i, chapter in enumerate(chapters):
        if not chapter.get('content'):
            continue

        chapter_num = i + 1
        audio_file = DATA_DIR / project_id / 'audio' / f'chapter_{chapter_num:02d}.mp3'

        try:
            # Using piper TTS (or espeak as fallback)
            text = chapter['content']

            # Try piper first
            piper_result = subprocess.run(
                ['which', 'piper'],
                capture_output=True
            )

            if piper_result.returncode == 0:
                # Use piper TTS
                process = subprocess.Popen(
                    ['piper', '--model', voice, '--output_file', str(audio_file)],
                    stdin=subprocess.PIPE,
                    stdout=subprocess.PIPE,
                    stderr=subprocess.PIPE
                )
                process.communicate(input=text.encode())
            else:
                # Fallback to espeak
                subprocess.run(
                    ['espeak', '-w', str(audio_file), text[:5000]],  # Limit text length for espeak
                    capture_output=True,
                    timeout=300
                )

            projects[project_id]['chapters'][i]['has_audio'] = True
            projects[project_id]['chapters'][i]['audio_path'] = str(audio_file)
            results.append({'chapter': chapter_num, 'audio_file': str(audio_file)})

        except Exception as e:
            results.append({'chapter': chapter_num, 'error': str(e)})

    save_projects()
    return {'converted': len([r for r in results if 'audio_file' in r]), 'results': results}


@api.route('/projects/<project_id>/generate-outline', methods=['POST'])
@validate_project_exists
def generate_outline(project_id):
    """Start outline generation for a project."""
    job_id = f'outline_{project_id}'

    if job_id in jobs and jobs[job_id]['status'] == 'running':
        return error_response('Outline generation already in progress', 409)

    jobs[job_id] = {
        'id': job_id,
        'type': 'outline_generation',
        'project_id': project_id,
        'status': 'pending',
        'created_at': datetime.utcnow().isoformat()
    }

    run_async_job(job_id, generate_outline_task, project_id)

    return api_response(
        data={'job_id': job_id},
        message='Outline generation started',
        status=202
    )


@api.route('/projects/<project_id>/generate-chapter/<int:chapter_num>', methods=['POST'])
@validate_project_exists
def generate_chapter(project_id, chapter_num):
    """Start chapter generation."""
    project = projects[project_id]

    if chapter_num < 1 or chapter_num > project.get('num_chapters', 10):
        return error_response(f'Invalid chapter number: {chapter_num}', 400)

    job_id = f'chapter_{project_id}_{chapter_num}'

    if job_id in jobs and jobs[job_id]['status'] == 'running':
        return error_response(f'Chapter {chapter_num} generation already in progress', 409)

    jobs[job_id] = {
        'id': job_id,
        'type': 'chapter_generation',
        'project_id': project_id,
        'chapter_num': chapter_num,
        'status': 'pending',
        'created_at': datetime.utcnow().isoformat()
    }

    run_async_job(job_id, generate_chapter_task, project_id, chapter_num)

    return api_response(
        data={'job_id': job_id},
        message=f'Chapter {chapter_num} generation started',
        status=202
    )


@api.route('/projects/<project_id>/generate-all', methods=['POST'])
@validate_project_exists
def generate_all_chapters(project_id):
    """Start generation of all chapters."""
    job_id = f'gen_all_{project_id}'

    if job_id in jobs and jobs[job_id]['status'] == 'running':
        return error_response('Full book generation already in progress', 409)

    jobs[job_id] = {
        'id': job_id,
        'type': 'full_generation',
        'project_id': project_id,
        'status': 'pending',
        'progress': 0,
        'created_at': datetime.utcnow().isoformat()
    }

    run_async_job(job_id, generate_all_chapters_task, project_id)

    return api_response(
        data={'job_id': job_id},
        message='Full book generation started',
        status=202
    )


@api.route('/projects/<project_id>/convert-audio', methods=['POST'])
@validate_project_exists
def convert_to_audio(project_id):
    """Start audio conversion for all chapters."""
    job_id = f'audio_{project_id}'

    if job_id in jobs and jobs[job_id]['status'] == 'running':
        return error_response('Audio conversion already in progress', 409)

    jobs[job_id] = {
        'id': job_id,
        'type': 'audio_conversion',
        'project_id': project_id,
        'status': 'pending',
        'created_at': datetime.utcnow().isoformat()
    }

    run_async_job(job_id, convert_to_audio_task, project_id)

    return api_response(
        data={'job_id': job_id},
        message='Audio conversion started',
        status=202
    )


# ============================================================================
# Export Routes
# ============================================================================

@api.route('/projects/<project_id>/export/<format>', methods=['GET'])
@validate_project_exists
def export_project(project_id, format):
    """Export project in specified format."""
    project = projects[project_id]
    valid_formats = ['txt', 'md', 'html', 'epub', 'json']

    if format not in valid_formats:
        return error_response(f'Invalid format. Supported: {", ".join(valid_formats)}', 400)

    export_dir = DATA_DIR / project_id / 'exports'
    export_dir.mkdir(exist_ok=True)

    title = project['title']
    chapters = project.get('chapters', [])

    if format == 'txt':
        content = f"{title}\n{'=' * len(title)}\n\n"
        for i, ch in enumerate(chapters):
            content += f"\nChapter {i + 1}: {ch.get('title', '')}\n"
            content += '-' * 40 + '\n'
            content += ch.get('content', '') + '\n\n'

        export_file = export_dir / f'{project_id}.txt'
        with open(export_file, 'w') as f:
            f.write(content)

        return send_file(export_file, as_attachment=True, download_name=f'{title}.txt')

    elif format == 'md':
        content = f"# {title}\n\n"
        if project.get('description'):
            content += f"*{project['description']}*\n\n"

        for i, ch in enumerate(chapters):
            content += f"\n## Chapter {i + 1}: {ch.get('title', '')}\n\n"
            content += ch.get('content', '') + '\n\n'

        export_file = export_dir / f'{project_id}.md'
        with open(export_file, 'w') as f:
            f.write(content)

        return send_file(export_file, as_attachment=True, download_name=f'{title}.md')

    elif format == 'html':
        content = f"""<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>{title}</title>
    <style>
        body {{ font-family: Georgia, serif; max-width: 800px; margin: 0 auto; padding: 20px; }}
        h1 {{ text-align: center; }}
        h2 {{ margin-top: 2em; border-bottom: 1px solid #ccc; }}
        p {{ line-height: 1.6; text-indent: 2em; }}
    </style>
</head>
<body>
    <h1>{title}</h1>
"""
        if project.get('description'):
            content += f"    <p><em>{project['description']}</em></p>\n"

        for i, ch in enumerate(chapters):
            content += f"    <h2>Chapter {i + 1}: {ch.get('title', '')}</h2>\n"
            paragraphs = ch.get('content', '').split('\n\n')
            for p in paragraphs:
                if p.strip():
                    content += f"    <p>{p.strip()}</p>\n"

        content += "</body>\n</html>"

        export_file = export_dir / f'{project_id}.html'
        with open(export_file, 'w') as f:
            f.write(content)

        return send_file(export_file, as_attachment=True, download_name=f'{title}.html')

    elif format == 'json':
        export_file = export_dir / f'{project_id}.json'
        with open(export_file, 'w') as f:
            json.dump(project, f, indent=2, default=str)

        return send_file(export_file, as_attachment=True, download_name=f'{title}.json')

    elif format == 'epub':
        # Basic EPUB generation (requires ebooklib for full support)
        try:
            from ebooklib import epub

            book = epub.EpubBook()
            book.set_identifier(project_id)
            book.set_title(title)
            book.set_language('en')

            chapters_epub = []
            for i, ch in enumerate(chapters):
                c = epub.EpubHtml(title=ch.get('title', f'Chapter {i+1}'),
                                  file_name=f'chapter_{i+1}.xhtml',
                                  lang='en')
                content_html = ''.join(f'<p>{p}</p>' for p in ch.get('content', '').split('\n\n') if p.strip())
                c.content = f'<h1>{ch.get("title", f"Chapter {i+1}")}</h1>{content_html}'
                book.add_item(c)
                chapters_epub.append(c)

            book.toc = chapters_epub
            book.add_item(epub.EpubNcx())
            book.add_item(epub.EpubNav())
            book.spine = ['nav'] + chapters_epub

            export_file = export_dir / f'{project_id}.epub'
            epub.write_epub(str(export_file), book)

            return send_file(export_file, as_attachment=True, download_name=f'{title}.epub')
        except ImportError:
            return error_response('EPUB export requires ebooklib package', 501)


# ============================================================================
# Job Routes
# ============================================================================

@api.route('/jobs', methods=['GET'])
def list_jobs():
    """List all jobs."""
    status_filter = request.args.get('status')

    job_list = []
    for job_id, job in jobs.items():
        if status_filter and job['status'] != status_filter:
            continue
        job_list.append(job)

    # Sort by creation time, newest first
    job_list.sort(key=lambda x: x.get('created_at', ''), reverse=True)

    return api_response(data={
        'jobs': job_list,
        'total': len(job_list),
        'running': len([j for j in jobs.values() if j['status'] == 'running'])
    })


@api.route('/jobs/<job_id>', methods=['GET'])
def get_job(job_id):
    """Get job status."""
    if job_id not in jobs:
        return error_response(f'Job {job_id} not found', 404)

    return api_response(data={'job': jobs[job_id]})


# ============================================================================
# Model and Voice Routes
# ============================================================================

@api.route('/models', methods=['GET'])
def list_models():
    """List available Ollama models."""
    try:
        result = subprocess.run(
            ['ollama', 'list'],
            capture_output=True,
            text=True,
            timeout=30
        )

        if result.returncode == 0:
            lines = result.stdout.strip().split('\n')
            models = []

            # Skip header line
            for line in lines[1:]:
                parts = line.split()
                if parts:
                    model_name = parts[0]
                    size = parts[1] if len(parts) > 1 else 'unknown'
                    models.append({
                        'name': model_name,
                        'size': size
                    })

            return api_response(data={'models': models})
        else:
            return error_response(f'Failed to list models: {result.stderr}', 500)
    except FileNotFoundError:
        return error_response('Ollama is not installed', 503)
    except subprocess.TimeoutExpired:
        return error_response('Timeout listing models', 504)


@api.route('/voices', methods=['GET'])
def list_voices():
    """List available TTS voices."""
    # Check for piper voices
    voices = []

    # Default system voices
    voices.append({'id': 'default', 'name': 'System Default', 'type': 'system'})

    # Check for piper
    try:
        piper_models_dir = Path('/home/tim/.local/share/piper/voices')
        if piper_models_dir.exists():
            for voice_file in piper_models_dir.glob('*.onnx'):
                voices.append({
                    'id': str(voice_file),
                    'name': voice_file.stem,
                    'type': 'piper'
                })
    except Exception:
        pass

    # Check for espeak voices
    try:
        result = subprocess.run(
            ['espeak', '--voices'],
            capture_output=True,
            text=True,
            timeout=10
        )
        if result.returncode == 0:
            lines = result.stdout.strip().split('\n')
            for line in lines[1:6]:  # Limit to first 5 espeak voices
                parts = line.split()
                if len(parts) >= 4:
                    voices.append({
                        'id': parts[4],
                        'name': parts[3],
                        'type': 'espeak',
                        'language': parts[1]
                    })
    except Exception:
        pass

    return api_response(data={'voices': voices})


# ============================================================================
# Error Handlers
# ============================================================================

@api.errorhandler(400)
def bad_request(e):
    return error_response('Bad request', 400)


@api.errorhandler(404)
def not_found(e):
    return error_response('Resource not found', 404)


@api.errorhandler(500)
def internal_error(e):
    return error_response('Internal server error', 500)


# Initialize projects on module load
load_projects()
