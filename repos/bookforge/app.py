#!/usr/bin/env python3
import os
import re
import subprocess
import wave
import psycopg2
from psycopg2.extras import RealDictCursor
import uuid
from datetime import datetime
from flask import Flask, request, jsonify, send_from_directory
from flask_cors import CORS
import requests
import json
import threading
import time

app = Flask(__name__, static_folder='static')
CORS(app)

AUDIO_DIR = '/home/tim/bookforge/audio'
OLLAMA_URL = os.environ.get('OLLAMA_URL', 'http://localhost:11434')
PG_CONFIG = {'host': 'localhost', 'database': 'bookforge', 'user': 'bookforge', 'password': 'bookforge123'}

os.makedirs(AUDIO_DIR, exist_ok=True)

# ============ INDEPENDENT WRITER MANAGER ============

active_writers = {}
writer_manager_running = False

def background_write_loop(session_id):
    """Background thread that writes sections continuously for a session"""
    print(f"[Writer] Starting writer for session {session_id}")

    while session_id in active_writers and active_writers[session_id]:
        try:
            conn = get_db()
            cur = conn.cursor()

            cur.execute("SELECT * FROM generation_sessions WHERE id = %s", (session_id,))
            session = cur.fetchone()

            if not session or session["step"] != "writing":
                print(f"[Writer] Session {session_id} no longer in writing state")
                conn.close()
                break

            cur.execute("""
                SELECT cs.id as section_id, cs.title as section_title, cs.section_type,
                       c.id as chapter_id, c.chapter_number, c.title as chapter_title
                FROM chapter_sections cs
                JOIN chapters c ON cs.chapter_id = c.id
                JOIN projects p ON c.project_id = p.id
                WHERE p.title = %s AND cs.status = 'pending'
                ORDER BY c.chapter_number, cs.section_number
                LIMIT 1
            """, (session["selected_title"],))

            next_section = cur.fetchone()

            if not next_section:
                print(f"[Writer] Session {session_id} complete - all sections written")
                cur.execute("UPDATE generation_sessions SET step = 'complete', updated_at = %s WHERE id = %s",
                            (datetime.utcnow(), session_id))
                conn.commit()
                conn.close()
                break

            section_id = next_section["section_id"]
            chapter_id = next_section["chapter_id"]

            print(f"[Writer] Writing: {next_section['section_title']} (Ch {next_section['chapter_number']})")

            cur.execute("UPDATE chapter_sections SET status = 'generating', updated_at = %s WHERE id = %s",
                        (datetime.utcnow(), section_id))
            conn.commit()

            cur.execute("SELECT * FROM chapter_sections WHERE id = %s", (section_id,))
            section = cur.fetchone()
            cur.execute("SELECT * FROM chapters WHERE id = %s", (chapter_id,))
            chapter = cur.fetchone()
            conn.close()

            section_type = section["section_type"]
            book_title = session["selected_title"]
            chapter_title = chapter["title"]
            section_title = section["title"]
            book_concept = session["prompt"]
            genre = session["selected_genre"] or ""

            prompts = {
                'opening': f'Write an engaging opening scenario for {chapter_title} of "{book_title}". 2-3 paragraphs, 300-500 words. Book concept: {book_concept}',
                'objectives': f'Write 4-5 learning objectives for {chapter_title}. Format: "By the end of this chapter, you will be able to..." Book concept: {book_concept}',
                'summary': f'Write a chapter summary for {chapter_title}. 5-7 bullet points. Book concept: {book_concept}',
                'terms': f'Write 5-8 key terms for {chapter_title}. Format: **Term:** Definition. Book concept: {book_concept}',
                'exercises': f'Write 8 practice exercises for {chapter_title}. 3 Basic, 3 Intermediate, 2 Advanced. Book concept: {book_concept}'
            }

            default_prompt = f'Write "{section_title}" for {chapter_title} of "{book_title}". 1000-1500 words with examples and headers. Book concept: {book_concept}'

            try:
                response = requests.post(OLLAMA_URL + "/api/chat", json={
                    "model": session["model"] or "qwen2.5:7b-instruct",
                    "messages": [
                        {"role": "system", "content": "You are an expert author. Write engaging content using markdown."},
                        {"role": "user", "content": prompts.get(section_type, default_prompt)}
                    ],
                    "stream": False,
                    "options": {"num_predict": 4000}
                }, timeout=300)

                content_text = response.json().get("message", {}).get("content", "")
                word_count = len(content_text.split())

                print(f"[Writer] Generated {word_count} words for {section_title}")

                conn = get_db()
                cur = conn.cursor()

                cur.execute("UPDATE chapter_sections SET content = %s, word_count = %s, status = 'complete', updated_at = %s WHERE id = %s",
                            (content_text, word_count, datetime.utcnow(), section_id))

                cur.execute("UPDATE chapters SET word_count = (SELECT COALESCE(SUM(word_count), 0) FROM chapter_sections WHERE chapter_id = %s), updated_at = %s WHERE id = %s",
                            (chapter_id, datetime.utcnow(), chapter_id))

                cur.execute("SELECT COUNT(*) as pending FROM chapter_sections WHERE chapter_id = %s AND status != 'complete'", (chapter_id,))
                if cur.fetchone()["pending"] == 0:
                    cur.execute("SELECT content FROM chapter_sections WHERE chapter_id = %s ORDER BY section_number", (chapter_id,))
                    secs = cur.fetchall()
                    full = "\n\n---\n\n".join([s["content"] for s in secs if s["content"]])
                    cur.execute("UPDATE chapters SET content = %s, status = 'complete', updated_at = %s WHERE id = %s", (full, datetime.utcnow(), chapter_id))
                    print(f"[Writer] Chapter {next_section['chapter_number']} complete")

                cur.execute("""UPDATE generation_sessions SET
                    current_chapter = (SELECT COUNT(*) FROM chapters c JOIN projects p ON c.project_id = p.id WHERE p.title = %s AND c.status = 'complete'),
                    total_words = (SELECT COALESCE(SUM(c.word_count), 0) FROM chapters c JOIN projects p ON c.project_id = p.id WHERE p.title = %s),
                    updated_at = %s WHERE id = %s""", (book_title, book_title, datetime.utcnow(), session_id))

                conn.commit()
                conn.close()

            except Exception as e:
                print(f"[Writer] Write error: {e}")
                conn = get_db()
                cur = conn.cursor()
                cur.execute("UPDATE chapter_sections SET status = 'pending' WHERE id = %s", (section_id,))
                conn.commit()
                conn.close()
                time.sleep(5)

        except Exception as e:
            print(f"[Writer] Background writer error: {e}")
            time.sleep(5)

    if session_id in active_writers:
        del active_writers[session_id]
    print(f"[Writer] Writer for session {session_id} stopped")


def writer_manager_loop():
    """Manager thread that monitors for sessions needing writing and starts writers automatically"""
    global writer_manager_running
    print("[WriterManager] Started - auto-starts writers for 'writing' sessions")

    while writer_manager_running:
        try:
            conn = get_db()
            cur = conn.cursor()
            cur.execute("SELECT id, selected_title FROM generation_sessions WHERE step = 'writing'")
            writing_sessions = cur.fetchall()
            conn.close()

            for session in writing_sessions:
                session_id = session["id"]
                if session_id not in active_writers or not active_writers[session_id]:
                    active_writers[session_id] = True
                    thread = threading.Thread(target=background_write_loop, args=(session_id,), daemon=True)
                    thread.start()
                    print(f"[WriterManager] Auto-started writer for: {session['selected_title']}")

        except Exception as e:
            print(f"[WriterManager] Error: {e}")

        time.sleep(5)

    print("[WriterManager] Stopped")


def start_writer_manager():
    global writer_manager_running
    if not writer_manager_running:
        writer_manager_running = True
        thread = threading.Thread(target=writer_manager_loop, daemon=True)
        thread.start()


# ============ INDEPENDENT AUDIO MANAGER ============

active_audio_generators = {}
audio_manager_running = False


def clean_text_for_speech(text):
    """Clean markdown text for natural TTS output"""
    text = re.sub(r'```.*?```', '', text, flags=re.DOTALL)
    text = re.sub(r'`[^`]+`', '', text)
    text = re.sub(r'<!--.*?-->', '', text, flags=re.DOTALL)
    text = re.sub(r'^---+$', '', text, flags=re.MULTILINE)
    text = re.sub(r'^#{1,6}\s*(.+)$', r'\1.', text, flags=re.MULTILINE)
    text = re.sub(r'\*\*([^*]+)\*\*', r'\1', text)
    text = re.sub(r'\*([^*]+)\*', r'\1', text)
    text = re.sub(r'__([^_]+)__', r'\1', text)
    text = re.sub(r'_([^_]+)_', r'\1', text)
    text = re.sub(r'^>\s*', '', text, flags=re.MULTILINE)
    text = re.sub(r'\[([^\]]+)\]\([^)]+\)', r'\1', text)
    text = re.sub(r'!\[.*?\]\(.*?\)', '', text)
    text = re.sub(r'^[\-\*]\s+', '', text, flags=re.MULTILINE)
    text = re.sub(r'^\d+\.\s+', '', text, flags=re.MULTILINE)
    text = re.sub(r'\|', ' ', text)
    text = re.sub(r'^[\-:]+$', '', text, flags=re.MULTILINE)
    text = re.sub(r'\[ \]', '', text)
    text = re.sub(r'\[x\]', '', text)
    text = re.sub(r'\n{3,}', '\n\n', text)
    text = re.sub(r'  +', ' ', text)
    return text.strip()


def get_audio_duration(wav_path):
    try:
        with wave.open(wav_path, 'rb') as wf:
            return wf.getnframes() / float(wf.getframerate())
    except:
        return 0


def generate_audio_xtts(text, output_path, settings):
    os.environ["COQUI_TOS_AGREED"] = "1"
    try:
        from TTS.api import TTS
        tts = TTS(model_name="tts_models/multilingual/multi-dataset/xtts_v2", gpu=False)

        if settings.get('voice_sample') and os.path.exists(settings['voice_sample']):
            tts.tts_to_file(
                text=text,
                file_path=output_path,
                speaker_wav=settings['voice_sample'],
                language=settings.get('language', 'en'),
                speed=settings.get('speed', 1.0)
            )
        else:
            tts.tts_to_file(
                text=text,
                file_path=output_path,
                language=settings.get('language', 'en'),
                speed=settings.get('speed', 1.0)
            )
        return True
    except ImportError:
        print("[Audio] XTTS not installed, using Piper")
        return generate_audio_piper(text, output_path, settings)
    except Exception as e:
        print(f"[Audio] XTTS error: {e}, using Piper")
        return generate_audio_piper(text, output_path, settings)


def generate_audio_piper(text, output_path, settings):
    """Generate audio using Piper TTS"""
    try:
        piper_path = '/home/tim/audiobook/venv/bin/piper'
        voice_model = settings.get('piper_model', '/home/tim/audiobook/voices/en_US-lessac-medium.onnx')

        process = subprocess.Popen(
            [piper_path, '--model', voice_model, '--output_file', output_path],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE
        )
        stdout, stderr = process.communicate(input=text.encode('utf-8'), timeout=600)
        return process.returncode == 0
    except Exception as e:
        print(f"[Audio] Piper error: {e}")
        return False


def audio_generator_loop(project_id):
    """Background thread that generates audio for sections"""
    print(f"[Audio] Starting audio generator for project {project_id}")

    while project_id in active_audio_generators and active_audio_generators[project_id]:
        try:
            conn = get_db()
            cur = conn.cursor()

            cur.execute("SELECT * FROM audio_settings WHERE id = 'default'")
            settings = cur.fetchone()
            settings = dict(settings) if settings else {'tts_engine': 'piper', 'language': 'en', 'speed': 1.0}

            cur.execute("""
                SELECT cs.id as section_id, cs.title, cs.content, cs.section_number,
                       c.id as chapter_id, c.chapter_number, c.title as chapter_title,
                       p.id as project_id, p.title as project_title
                FROM chapter_sections cs
                JOIN chapters c ON cs.chapter_id = c.id
                JOIN projects p ON c.project_id = p.id
                WHERE p.id = %s
                    AND cs.status = 'complete'
                    AND cs.content IS NOT NULL
                    AND cs.content != ''
                    AND (cs.audio_status IS NULL OR cs.audio_status = 'pending')
                ORDER BY c.chapter_number, cs.section_number
                LIMIT 1
            """, (project_id,))

            next_section = cur.fetchone()

            if not next_section:
                cur.execute("""
                    SELECT COUNT(*) as pending FROM chapter_sections cs
                    JOIN chapters c ON cs.chapter_id = c.id
                    WHERE c.project_id = %s AND cs.status = 'complete'
                    AND (cs.audio_status IS NULL OR cs.audio_status = 'pending')
                """, (project_id,))
                pending = cur.fetchone()['pending']

                if pending == 0:
                    print(f"[Audio] All sections have audio for project {project_id}")
                    conn.close()
                    break

                conn.close()
                time.sleep(10)
                continue

            section_id = next_section['section_id']

            print(f"[Audio] Generating: Ch{next_section['chapter_number']} - {next_section['title']}")

            cur.execute("UPDATE chapter_sections SET audio_status = 'generating', updated_at = %s WHERE id = %s",
                        (datetime.utcnow(), section_id))
            conn.commit()
            conn.close()

            clean_text = clean_text_for_speech(next_section['content'])

            safe_title = "".join(c if c.isalnum() or c in ' -_' else '' for c in next_section['title'])[:50]
            output_filename = f"ch{next_section['chapter_number']:02d}_s{next_section['section_number']:02d}_{safe_title}.wav"
            output_path = os.path.join(AUDIO_DIR, project_id, output_filename)
            os.makedirs(os.path.dirname(output_path), exist_ok=True)

            success = False
            if settings.get('tts_engine') == 'xtts':
                success = generate_audio_xtts(clean_text, output_path, settings)
            else:
                success = generate_audio_piper(clean_text, output_path, settings)

            conn = get_db()
            cur = conn.cursor()

            if success and os.path.exists(output_path):
                duration = get_audio_duration(output_path)
                cur.execute("""
                    UPDATE chapter_sections
                    SET audio_path = %s, audio_status = 'complete', audio_duration = %s, updated_at = %s
                    WHERE id = %s
                """, (output_path, duration, datetime.utcnow(), section_id))
                print(f"[Audio] Generated {duration:.1f}s for {next_section['title']}")
            else:
                cur.execute("UPDATE chapter_sections SET audio_status = 'error', updated_at = %s WHERE id = %s",
                            (datetime.utcnow(), section_id))
                print(f"[Audio] Failed for {next_section['title']}")

            conn.commit()
            conn.close()

        except Exception as e:
            print(f"[Audio] Generator error: {e}")
            time.sleep(5)

    if project_id in active_audio_generators:
        del active_audio_generators[project_id]
    print(f"[Audio] Audio generator for project {project_id} stopped")


def start_audio_manager():
    global audio_manager_running
    if not audio_manager_running:
        audio_manager_running = True
        print("[AudioManager] Started")


# ============ DATABASE ============

SECTION_TYPES = [
    ('opening', 'Opening Scenario'),
    ('objectives', 'Learning Objectives'),
    ('section_1', 'Section 1'),
    ('section_2', 'Section 2'),
    ('section_3', 'Section 3'),
    ('section_4', 'Section 4'),
    ('pitfalls', 'Common Pitfalls'),
    ('summary', 'Chapter Summary'),
    ('terms', 'Key Terms'),
    ('exercises', 'Practice Exercises'),
]

def get_db():
    return psycopg2.connect(**PG_CONFIG, cursor_factory=RealDictCursor)


# ============ ROUTES ============

@app.route('/')
def home():
    return send_from_directory('templates', 'index.html')


@app.route('/static/<path:filepath>')
def serve_static(filepath):
    return send_from_directory('static', filepath)


@app.route('/sessions', methods=['GET'])
def list_sessions():
    conn = get_db()
    cur = conn.cursor()
    cur.execute('SELECT * FROM generation_sessions ORDER BY created_at DESC LIMIT 20')
    sessions = cur.fetchall()
    conn.close()
    return jsonify({'success': True, 'sessions': [dict(s) for s in sessions]})


@app.route('/api/generate-titles', methods=['POST'])
def generate_titles():
    data = request.get_json() or {}
    prompt = data.get('prompt', '').strip()
    model = data.get('model', 'qwen2.5:7b-instruct')

    if not prompt:
        return jsonify({'success': False, 'error': 'Prompt is required'}), 400

    session_id = str(uuid.uuid4())

    conn = get_db()
    cur = conn.cursor()
    cur.execute('''
        INSERT INTO generation_sessions (id, prompt, model, step, created_at, updated_at)
        VALUES (%s, %s, %s, %s, %s, %s)
    ''', (session_id, prompt, model, 'titles_generating', datetime.utcnow(), datetime.utcnow()))
    conn.commit()
    conn.close()

    system_prompt = '''Generate 5 book title ideas as a JSON array. Output ONLY valid JSON, nothing else.
Format: [{"title":"Title","tagline":"Hook","genre":"Genre"}]'''

    try:
        response = requests.post(
            f'{OLLAMA_URL}/api/chat',
            json={
                'model': model,
                'messages': [
                    {'role': 'system', 'content': system_prompt},
                    {'role': 'user', 'content': f'Book concept: {prompt}'}
                ],
                'stream': False
            },
            timeout=120
        )
        response.raise_for_status()
        content = response.json().get('message', {}).get('content', '')

        conn = get_db()
        cur = conn.cursor()
        cur.execute('''
            UPDATE generation_sessions SET titles_json = %s, step = %s, updated_at = %s WHERE id = %s
        ''', (content, 'titles', datetime.utcnow(), session_id))
        conn.commit()
        conn.close()

        return jsonify({'success': True, 'session_id': session_id})
    except Exception as e:
        conn = get_db()
        cur = conn.cursor()
        cur.execute('DELETE FROM generation_sessions WHERE id = %s', (session_id,))
        conn.commit()
        conn.close()
        return jsonify({'success': False, 'error': str(e)}), 500


@app.route('/titles/<session_id>')
def view_titles(session_id):
    return send_from_directory('templates', 'titles.html')


@app.route('/api/titles/<session_id>')
def get_titles(session_id):
    conn = get_db()
    cur = conn.cursor()
    cur.execute('SELECT * FROM generation_sessions WHERE id = %s', (session_id,))
    session = cur.fetchone()
    conn.close()

    if not session:
        return jsonify({'success': False, 'error': 'Session not found'}), 404

    return jsonify({
        'success': True,
        'prompt': session['prompt'],
        'model': session['model'],
        'titles': session['titles_json'],
        'step': session['step'],
        'selected_title': session['selected_title']
    })


@app.route('/api/select-title', methods=['POST'])
def select_title():
    data = request.get_json() or {}
    session_id = data.get('session_id')
    title = data.get('title')
    tagline = data.get('tagline', '')
    genre = data.get('genre', '')

    if not session_id or not title:
        return jsonify({'success': False, 'error': 'session_id and title required'}), 400

    conn = get_db()
    cur = conn.cursor()
    cur.execute('''
        UPDATE generation_sessions
        SET selected_title = %s, selected_tagline = %s, selected_genre = %s, step = %s, updated_at = %s
        WHERE id = %s
    ''', (title, tagline, genre, 'title_selected', datetime.utcnow(), session_id))
    conn.commit()
    conn.close()

    return jsonify({'success': True})


@app.route('/projects', methods=['GET'])
def list_projects():
    conn = get_db()
    cur = conn.cursor()
    cur.execute('''SELECT * FROM projects ORDER BY updated_at DESC''')
    projects = cur.fetchall()
    conn.close()
    return jsonify({'success': True, 'count': len(projects), 'projects': [dict(p) for p in projects]})


@app.route('/projects/<project_id>', methods=['GET'])
def get_project(project_id):
    conn = get_db()
    cur = conn.cursor()
    cur.execute('SELECT * FROM projects WHERE id = %s', (project_id,))
    project = cur.fetchone()
    conn.close()
    if not project:
        return jsonify({'success': False, 'error': 'Not found'}), 404
    return jsonify({'success': True, 'project': dict(project)})


@app.route('/projects/<project_id>/chapters', methods=['GET'])
def list_chapters(project_id):
    conn = get_db()
    cur = conn.cursor()
    cur.execute('SELECT * FROM chapters WHERE project_id = %s ORDER BY chapter_number', (project_id,))
    chapters = cur.fetchall()
    conn.close()
    return jsonify({'success': True, 'chapters': [dict(c) for c in chapters]})


@app.route('/book/<project_id>')
def view_book(project_id):
    return send_from_directory('templates', 'reader.html')


@app.route('/health')
def health():
    return jsonify({'status': 'healthy', 'service': 'BookForge', 'version': '2.2.0'})


@app.route("/api/models")
def list_models():
    try:
        res = requests.get(f"{OLLAMA_URL}/api/tags", timeout=5)
        models = res.json().get("models", [])
        return jsonify({"success": True, "models": [m["name"] for m in models]})
    except:
        return jsonify({"success": True, "models": ["qwen2.5:7b-instruct", "llama3.1:8b"]})


@app.route("/api/generate-outline", methods=["POST"])
def generate_outline():
    data = request.get_json() or {}
    session_id = data.get("session_id")
    if not session_id:
        return jsonify({"success": False, "error": "session_id required"}), 400

    conn = get_db()
    cur = conn.cursor()
    cur.execute("SELECT * FROM generation_sessions WHERE id = %s", (session_id,))
    session = cur.fetchone()

    if not session or not session["selected_title"]:
        conn.close()
        return jsonify({"success": False, "error": "Session or title not found"}), 404

    cur.execute("UPDATE generation_sessions SET step = %s, updated_at = %s WHERE id = %s",
                ("outline_generating", datetime.utcnow(), session_id))
    conn.commit()

    system_prompt = '''Generate a detailed book outline as JSON. Output ONLY valid JSON:
{
  "chapters": [
    {"number": 1, "title": "Chapter Title", "summary": "Brief description",
     "sections": [{"type": "opening", "title": "Opening Scenario", "description": "..."}]}
  ]
}
Generate 10-15 chapters with sections.'''

    user_prompt = f"Book: {session['selected_title']}\nGenre: {session['selected_genre']}\nConcept: {session['prompt']}"

    try:
        response = requests.post(OLLAMA_URL + "/api/chat", json={
            "model": session["model"] or "qwen2.5:7b-instruct",
            "messages": [{"role": "system", "content": system_prompt}, {"role": "user", "content": user_prompt}],
            "stream": False
        }, timeout=180)
        content = response.json().get("message", {}).get("content", "")

        total_chapters = 0
        try:
            match = content.find('{')
            if match >= 0:
                json_str = content[match:]
                end = json_str.rfind('}') + 1
                parsed = json.loads(json_str[:end])
                total_chapters = len(parsed.get('chapters', []))
        except:
            pass

        cur.execute("""UPDATE generation_sessions
                      SET outline = %s, step = %s, total_chapters = %s, updated_at = %s
                      WHERE id = %s""",
                    (content, "outline_complete", total_chapters, datetime.utcnow(), session_id))
        conn.commit()
        conn.close()
        return jsonify({"success": True})
    except Exception as e:
        cur.execute("UPDATE generation_sessions SET step = %s, updated_at = %s WHERE id = %s",
                    ("title_selected", datetime.utcnow(), session_id))
        conn.commit()
        conn.close()
        return jsonify({"success": False, "error": str(e)}), 500


@app.route("/outline/<session_id>")
def view_outline(session_id):
    return send_from_directory("templates", "outline.html")


@app.route("/api/outline/<session_id>")
def get_outline(session_id):
    conn = get_db()
    cur = conn.cursor()
    cur.execute("SELECT * FROM generation_sessions WHERE id = %s", (session_id,))
    session = cur.fetchone()
    conn.close()
    if not session:
        return jsonify({"success": False, "error": "Not found"}), 404
    return jsonify({
        "success": True,
        "title": session["selected_title"],
        "tagline": session["selected_tagline"],
        "genre": session["selected_genre"],
        "prompt": session["prompt"],
        "outline": session["outline"],
        "step": session["step"],
        "total_chapters": session.get("total_chapters", 0),
        "current_chapter": session.get("current_chapter", 0),
        "total_words": session.get("total_words", 0)
    })


@app.route("/api/session/<session_id>/status")
def get_session_status(session_id):
    conn = get_db()
    cur = conn.cursor()
    cur.execute("SELECT id, step, selected_title, outline, current_chapter, current_section, total_chapters, total_words FROM generation_sessions WHERE id = %s", (session_id,))
    session = cur.fetchone()
    conn.close()
    if not session:
        return jsonify({"success": False, "error": "Not found"}), 404
    return jsonify({
        "success": True,
        "step": session["step"],
        "has_outline": session["outline"] is not None,
        "current_chapter": session.get("current_chapter", 0),
        "current_section": session.get("current_section", 0),
        "total_chapters": session.get("total_chapters", 0),
        "total_words": session.get("total_words", 0)
    })


# ============ WRITING ENDPOINTS ============

@app.route("/writing/<session_id>")
def view_writing(session_id):
    return send_from_directory("templates", "writing.html")


@app.route("/api/start-writing", methods=["POST"])
def start_writing():
    data = request.get_json() or {}
    session_id = data.get("session_id")
    write_model = data.get("model")

    if not session_id:
        return jsonify({"success": False, "error": "session_id required"}), 400

    conn = get_db()
    cur = conn.cursor()
    cur.execute("SELECT * FROM generation_sessions WHERE id = %s", (session_id,))
    session = cur.fetchone()

    if not session or not session["outline"]:
        conn.close()
        return jsonify({"success": False, "error": "Session or outline not found"}), 404

    outline = session["outline"]
    chapters = []
    try:
        match = outline.find('{')
        if match >= 0:
            json_str = outline[match:]
            end = json_str.rfind('}') + 1
            parsed = json.loads(json_str[:end])
            chapters = parsed.get('chapters', [])
    except Exception as e:
        conn.close()
        return jsonify({"success": False, "error": f"Could not parse outline: {e}"}), 400

    if not chapters:
        conn.close()
        return jsonify({"success": False, "error": "No chapters found in outline"}), 400

    project_id = str(uuid.uuid4())
    cur.execute("""
        INSERT INTO projects (id, title, description, status, created_at, updated_at)
        VALUES (%s, %s, %s, %s, %s, %s)
    """, (project_id, session["selected_title"], session["prompt"], 'writing', datetime.utcnow(), datetime.utcnow()))

    for ch in chapters:
        chapter_id = str(uuid.uuid4())
        ch_num = ch.get('number', 1)
        ch_title = ch.get('title', f'Chapter {ch_num}')
        ch_summary = ch.get('summary', '')

        cur.execute("""
            INSERT INTO chapters (id, project_id, chapter_number, title, summary, status, created_at, updated_at)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
        """, (chapter_id, project_id, ch_num, ch_title, ch_summary, 'pending', datetime.utcnow(), datetime.utcnow()))

        sections = ch.get('sections', [])
        for i, sec in enumerate(sections):
            section_id = str(uuid.uuid4())
            sec_type = sec.get('type', f'section_{i+1}')
            sec_title = sec.get('title', f'Section {i+1}')

            cur.execute("""
                INSERT INTO chapter_sections (id, chapter_id, section_number, section_type, title, status, created_at, updated_at)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
            """, (section_id, chapter_id, i+1, sec_type, sec_title, 'pending', datetime.utcnow(), datetime.utcnow()))

        existing_types = [s.get('type', '') for s in sections]
        standard_sections = [('summary', 'Chapter Summary'), ('terms', 'Key Terms'), ('exercises', 'Practice Exercises')]
        section_num = len(sections) + 1
        for sec_type, sec_title in standard_sections:
            if sec_type not in existing_types:
                section_id = str(uuid.uuid4())
                cur.execute("""
                    INSERT INTO chapter_sections (id, chapter_id, section_number, section_type, title, status)
                    VALUES (%s, %s, %s, %s, %s, %s)
                """, (section_id, chapter_id, section_num, sec_type, sec_title, 'pending'))
                section_num += 1

    final_model = write_model or session.get("model") or "qwen2.5:7b-instruct"
    cur.execute("""
        UPDATE generation_sessions
        SET step = %s, current_chapter = %s, total_chapters = %s, model = %s, updated_at = %s
        WHERE id = %s
    """, ('writing', 0, len(chapters), final_model, datetime.utcnow(), session_id))

    conn.commit()
    conn.close()

    print(f"[API] Started writing for session {session_id}")

    return jsonify({
        "success": True,
        "project_id": project_id,
        "total_chapters": len(chapters),
        "message": "Writing started automatically"
    })


@app.route("/api/writing/<session_id>/progress")
def get_writing_progress(session_id):
    conn = get_db()
    cur = conn.cursor()

    cur.execute("SELECT * FROM generation_sessions WHERE id = %s", (session_id,))
    session = cur.fetchone()

    if not session:
        conn.close()
        return jsonify({"success": False, "error": "Session not found"}), 404

    cur.execute("""
        SELECT p.id, p.title FROM projects p
        WHERE p.title = %s
        ORDER BY p.created_at DESC LIMIT 1
    """, (session["selected_title"],))
    project = cur.fetchone()

    if not project:
        conn.close()
        return jsonify({"success": True, "status": "not_started", "session": dict(session)})

    cur.execute("""
        SELECT c.*,
            (SELECT COUNT(*) FROM chapter_sections cs WHERE cs.chapter_id = c.id) as total_sections,
            (SELECT COUNT(*) FROM chapter_sections cs WHERE cs.chapter_id = c.id AND cs.status = 'complete') as completed_sections
        FROM chapters c
        WHERE c.project_id = %s
        ORDER BY c.chapter_number
    """, (project["id"],))
    chapters = cur.fetchall()

    chapters_with_sections = []
    for ch in chapters:
        cur.execute("""
            SELECT id, section_number, section_type, title, word_count, status
            FROM chapter_sections WHERE chapter_id = %s ORDER BY section_number
        """, (ch["id"],))
        sections = cur.fetchall()
        ch_dict = dict(ch)
        ch_dict["sections"] = [dict(s) for s in sections]
        chapters_with_sections.append(ch_dict)
    chapters = chapters_with_sections

    cur.execute("""
        SELECT cs.*, c.chapter_number, c.title as chapter_title
        FROM chapter_sections cs
        JOIN chapters c ON cs.chapter_id = c.id
        WHERE c.project_id = %s AND cs.status = 'generating'
        LIMIT 1
    """, (project["id"],))
    current = cur.fetchone()

    conn.close()

    total_words = sum(ch.get('word_count', 0) or 0 for ch in chapters)
    completed_chapters = sum(1 for ch in chapters if ch['status'] == 'complete')
    is_writing = session_id in active_writers and active_writers[session_id]

    return jsonify({
        "success": True,
        "status": "writing" if (current or is_writing) else ("complete" if completed_chapters == len(chapters) else "paused"),
        "project_id": project["id"],
        "project_title": project["title"],
        "total_chapters": len(chapters),
        "completed_chapters": completed_chapters,
        "total_words": total_words,
        "model": session.get("model"),
        "chapters": chapters,
        "current_section": dict(current) if current else None,
        "auto_writing": is_writing
    })


@app.route("/api/write-next", methods=["POST"])
def write_next_section():
    data = request.get_json() or {}
    session_id = data.get("session_id")

    if not session_id:
        return jsonify({"success": False, "error": "session_id required"}), 400

    conn = get_db()
    cur = conn.cursor()

    cur.execute("SELECT * FROM generation_sessions WHERE id = %s", (session_id,))
    session = cur.fetchone()

    if not session:
        conn.close()
        return jsonify({"success": False, "error": "Session not found"}), 404

    cur.execute("""
        SELECT cs.id as section_id, c.id as chapter_id
        FROM chapter_sections cs
        JOIN chapters c ON cs.chapter_id = c.id
        JOIN projects p ON c.project_id = p.id
        WHERE p.title = %s AND cs.status = 'pending'
        ORDER BY c.chapter_number, cs.section_number
        LIMIT 1
    """, (session["selected_title"],))

    next_section = cur.fetchone()

    if not next_section:
        cur.execute("UPDATE generation_sessions SET step = 'complete', updated_at = %s WHERE id = %s",
                    (datetime.utcnow(), session_id))
        conn.commit()
        conn.close()
        return jsonify({"success": True, "complete": True, "message": "All sections written!"})

    conn.close()
    return jsonify({"success": True, "message": "Section queued"})


@app.route("/api/auto-write-status/<session_id>")
def auto_write_status(session_id):
    is_active = session_id in active_writers and active_writers[session_id]
    return jsonify({"success": True, "active": is_active})


@app.route("/api/stop-writing", methods=["POST"])
def stop_writing():
    data = request.get_json() or {}
    session_id = data.get("session_id")
    if session_id in active_writers:
        active_writers[session_id] = False
        print(f"[API] Stopped writing for session {session_id}")
    return jsonify({"success": True, "message": "Writing stopped"})


@app.route("/api/resume-writing", methods=["POST"])
def resume_writing():
    data = request.get_json() or {}
    session_id = data.get("session_id")

    if not session_id:
        return jsonify({"success": False, "error": "session_id required"}), 400

    conn = get_db()
    cur = conn.cursor()
    cur.execute("UPDATE generation_sessions SET step = 'writing', updated_at = %s WHERE id = %s",
                (datetime.utcnow(), session_id))
    conn.commit()
    conn.close()

    print(f"[API] Resumed writing for session {session_id}")
    return jsonify({"success": True, "message": "Writing resumed"})


@app.route("/api/section/<section_id>")
def get_section(section_id):
    conn = get_db()
    cur = conn.cursor()
    cur.execute("SELECT * FROM chapter_sections WHERE id = %s", (section_id,))
    section = cur.fetchone()
    conn.close()
    if not section:
        return jsonify({"success": False, "error": "Section not found"}), 404
    return jsonify({"success": True, "section": dict(section)})


# Legacy endpoints
@app.route("/api/start-auto-write", methods=["POST"])
def start_auto_write():
    return jsonify({"success": True, "message": "Writing is automatic"})


@app.route("/api/stop-auto-write", methods=["POST"])
def stop_auto_write():
    return stop_writing()


# ============ AUDIO ENDPOINTS ============

@app.route("/audio")
def view_audio():
    return send_from_directory("templates", "audio.html")


@app.route("/api/audio-settings", methods=["GET"])
def get_audio_settings():
    conn = get_db()
    cur = conn.cursor()
    cur.execute("SELECT * FROM audio_settings WHERE id = 'default'")
    settings = cur.fetchone()

    voice_samples = []
    voice_dir = '/home/tim/bookforge/voices'
    if os.path.exists(voice_dir):
        for f in os.listdir(voice_dir):
            if f.endswith(('.wav', '.mp3', '.ogg')):
                voice_samples.append({'name': os.path.splitext(f)[0], 'path': os.path.join(voice_dir, f)})

    conn.close()

    return jsonify({
        "success": True,
        "settings": dict(settings) if settings else {},
        "voice_samples": voice_samples,
        "engines": [
            {"id": "xtts", "name": "XTTS-v2 (Coqui)", "supports_cloning": True},
            {"id": "piper", "name": "Piper TTS", "supports_cloning": False}
        ],
        "languages": [
            {"code": "en", "name": "English"}, {"code": "es", "name": "Spanish"},
            {"code": "fr", "name": "French"}, {"code": "de", "name": "German"},
            {"code": "it", "name": "Italian"}, {"code": "pt", "name": "Portuguese"},
            {"code": "pl", "name": "Polish"}, {"code": "ru", "name": "Russian"},
            {"code": "nl", "name": "Dutch"}, {"code": "zh-cn", "name": "Chinese"},
            {"code": "ja", "name": "Japanese"}, {"code": "ko", "name": "Korean"}
        ]
    })


@app.route("/api/audio-settings", methods=["POST"])
def update_audio_settings():
    data = request.get_json() or {}

    conn = get_db()
    cur = conn.cursor()

    cur.execute("""
        UPDATE audio_settings SET
            tts_engine = COALESCE(%s, tts_engine),
            voice_sample = COALESCE(%s, voice_sample),
            language = COALESCE(%s, language),
            speed = COALESCE(%s, speed),
            temperature = COALESCE(%s, temperature),
            top_k = COALESCE(%s, top_k),
            top_p = COALESCE(%s, top_p),
            updated_at = %s
        WHERE id = 'default'
    """, (
        data.get('tts_engine'),
        data.get('voice_sample'),
        data.get('language'),
        data.get('speed'),
        data.get('temperature'),
        data.get('top_k'),
        data.get('top_p'),
        datetime.utcnow()
    ))

    conn.commit()
    conn.close()

    return jsonify({"success": True, "message": "Settings updated"})


@app.route("/api/upload-voice-sample", methods=["POST"])
def upload_voice_sample():
    if 'file' not in request.files:
        return jsonify({"success": False, "error": "No file provided"}), 400

    file = request.files['file']
    if file.filename == '':
        return jsonify({"success": False, "error": "No file selected"}), 400

    voice_dir = '/home/tim/bookforge/voices'
    os.makedirs(voice_dir, exist_ok=True)

    filename = file.filename.replace(' ', '_')
    filepath = os.path.join(voice_dir, filename)
    file.save(filepath)

    return jsonify({"success": True, "path": filepath, "name": os.path.splitext(filename)[0]})


@app.route("/api/start-audio-generation", methods=["POST"])
def start_audio_generation():
    data = request.get_json() or {}
    project_id = data.get("project_id")

    if not project_id:
        return jsonify({"success": False, "error": "project_id required"}), 400

    if project_id in active_audio_generators and active_audio_generators[project_id]:
        return jsonify({"success": True, "message": "Already generating"})

    active_audio_generators[project_id] = True
    thread = threading.Thread(target=audio_generator_loop, args=(project_id,), daemon=True)
    thread.start()

    return jsonify({"success": True, "message": "Audio generation started"})


@app.route("/api/stop-audio-generation", methods=["POST"])
def stop_audio_generation():
    data = request.get_json() or {}
    project_id = data.get("project_id")

    if project_id in active_audio_generators:
        active_audio_generators[project_id] = False

    return jsonify({"success": True, "message": "Stopping..."})


@app.route("/api/audio-status/<project_id>")
def get_audio_status(project_id):
    conn = get_db()
    cur = conn.cursor()

    cur.execute("""
        SELECT
            COUNT(*) as total,
            SUM(CASE WHEN audio_status = 'complete' THEN 1 ELSE 0 END) as completed,
            SUM(CASE WHEN audio_status = 'generating' THEN 1 ELSE 0 END) as generating,
            SUM(CASE WHEN audio_status = 'error' THEN 1 ELSE 0 END) as errors,
            SUM(COALESCE(cs.audio_duration, 0)) as total_duration
        FROM chapter_sections cs
        JOIN chapters c ON cs.chapter_id = c.id
        WHERE c.project_id = %s AND cs.status = 'complete'
    """, (project_id,))
    stats = cur.fetchone()

    cur.execute("""
        SELECT cs.id, cs.title, cs.section_number, cs.audio_status, cs.audio_path, cs.audio_duration,
               c.chapter_number, c.title as chapter_title
        FROM chapter_sections cs
        JOIN chapters c ON cs.chapter_id = c.id
        WHERE c.project_id = %s AND cs.status = 'complete'
        ORDER BY c.chapter_number, cs.section_number
    """, (project_id,))
    sections = cur.fetchall()

    conn.close()

    is_active = project_id in active_audio_generators and active_audio_generators[project_id]

    return jsonify({
        "success": True,
        "active": is_active,
        "total_sections": stats['total'] or 0,
        "completed": stats['completed'] or 0,
        "generating": stats['generating'] or 0,
        "errors": stats['errors'] or 0,
        "total_duration": stats['total_duration'] or 0,
        "sections": [dict(s) for s in sections]
    })


@app.route("/api/regenerate-audio/<section_id>", methods=["POST"])
def regenerate_section_audio(section_id):
    conn = get_db()
    cur = conn.cursor()
    cur.execute("UPDATE chapter_sections SET audio_status = 'pending', audio_path = NULL, audio_duration = NULL WHERE id = %s",
                (section_id,))
    conn.commit()
    conn.close()
    return jsonify({"success": True, "message": "Section queued for regeneration"})


@app.route("/audio/<path:filepath>")
def serve_generated_audio(filepath):
    return send_from_directory(AUDIO_DIR, filepath)


@app.errorhandler(404)
def not_found(e):
    return jsonify({"success": False, "error": "Not found"}), 404


if __name__ == "__main__":
    print("BookForge v2.2 - Starting on http://0.0.0.0:8080")
    print("Independent writer manager enabled")
    print("Audio generation with XTTS-v2 support enabled")

    start_writer_manager()
    start_audio_manager()

    app.run(host="0.0.0.0", port=8080, debug=True, use_reloader=False)
