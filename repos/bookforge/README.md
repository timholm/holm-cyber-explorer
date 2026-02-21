# BookForge

AI-powered book generation system using Ollama LLM with text-to-speech support.

## Features

- Generate book outlines, chapters, and full books using AI
- Multiple title suggestions with genre detection
- Section-by-section writing with progress tracking
- Text-to-speech audio generation (Piper/XTTS)
- Export to EPUB, PDF, and audiobook formats
- Web UI for project management
- PostgreSQL for data persistence
- Kubernetes-ready with Helm chart

## Quick Start

### Local Development

```bash
# Create virtual environment
python3 -m venv venv
source venv/bin/activate

# Install dependencies
pip install -r requirements.txt

# Set environment variables
export OLLAMA_URL=http://localhost:11434
export POSTGRES_HOST=localhost
export POSTGRES_DB=bookforge
export POSTGRES_USER=bookforge
export POSTGRES_PASSWORD=bookforge123

# Run the app
python app.py
```

### Docker

```bash
docker build -t bookforge .
docker run -p 8080:8080 \
  -e OLLAMA_URL=http://host.docker.internal:11434 \
  bookforge
```

### Kubernetes with Helm

```bash
helm install bookforge ./chart -n default
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Web UI home |
| `/health` | GET | Health check |
| `/sessions` | GET | List generation sessions |
| `/projects` | GET | List all projects |
| `/api/generate-titles` | POST | Generate book title ideas |
| `/api/select-title` | POST | Select a title and start outline |
| `/api/generate-outline` | POST | Generate chapter outline |
| `/api/start-writing` | POST | Start background writing |
| `/book/<id>` | GET | View/read a book |

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Web UI    │────▶│  Flask API  │────▶│   Ollama    │
└─────────────┘     └─────────────┘     └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │ PostgreSQL  │
                    └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │  TTS Engine │
                    │ (Piper/XTTS)│
                    └─────────────┘
```

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `OLLAMA_URL` | http://localhost:11434 | Ollama API URL |
| `POSTGRES_HOST` | localhost | PostgreSQL host |
| `POSTGRES_DB` | bookforge | Database name |
| `POSTGRES_USER` | bookforge | Database user |
| `POSTGRES_PASSWORD` | - | Database password |
| `PORT` | 8080 | Server port |

## Helm Values

See `chart/values.yaml` for all configurable options:

```yaml
image:
  repository: ghcr.io/timholm/bookforge
  tag: latest

ollama:
  url: http://ollama-external:11434

postgres:
  host: postgres
  database: bookforge
```

## License

MIT
