# YouTube Channel Archiver

A Kubernetes-native service for archiving YouTube channels with parallel workers.

## Features

- **Parallel download workers** - Auto-scaling workers process downloads concurrently
- **iSCSI storage via Trident** - NetApp Trident integration for persistent block storage
- **Web UI for browsing archives** - Browse and search downloaded content
- **Argo Workflows integration** - Define complex download pipelines as workflows
- **Ko-based container builds** - Fast, reproducible container builds without Dockerfiles
- **Redis queue management** - Reliable job queuing with Redis
- **SQLite metadata storage** - Lightweight local metadata persistence
- **REST API** - Full-featured API for channel management and monitoring

## Quick Start

### Prerequisites

- Go 1.22+
- Kubernetes cluster (v1.28+)
- Redis 7.0+
- NetApp Trident (for iSCSI storage)
- yt-dlp (for video downloads)
- ko (for container builds)

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/timholm/ytarchive.git
   cd ytarchive
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up Redis**
   ```bash
   # Using Docker
   docker run -d --name redis -p 6379:6379 redis:7-alpine

   # Or using Homebrew on macOS
   brew install redis
   brew services start redis
   ```

4. **Set environment variables**
   ```bash
   export REDIS_ADDR=localhost:6379
   export K8S_NAMESPACE=default
   ```

5. **Run the controller**
   ```bash
   go run ./cmd/controller
   ```

6. **Access the API**
   ```bash
   curl http://localhost:8080/health
   ```

### Docker Deployment

```bash
# Build with ko
ko build ./cmd/controller

# Or use Docker directly
docker build -t ytarchive:latest .
docker run -p 8080:8080 -e REDIS_ADDR=redis:6379 ytarchive:latest
```

### Kubernetes Deployment

See [docs/deployment.md](docs/deployment.md) for detailed Kubernetes deployment instructions.

```bash
# Quick deployment with ArgoCD
kubectl apply -f deploy/argocd/application.yaml

# Or manual deployment
kubectl apply -f deploy/kubernetes/
```

## Architecture

```
                                    +------------------+
                                    |   Web Browser    |
                                    +--------+---------+
                                             |
                                             v
+----------------+              +------------+------------+
|   ArgoCD       |              |     API Controller      |
|   (GitOps)     +------------->|     (Port 8080)         |
+----------------+              +------------+------------+
                                             |
                    +------------------------+------------------------+
                    |                        |                        |
                    v                        v                        v
           +-------+-------+        +-------+-------+        +-------+-------+
           |    Redis      |        |    SQLite     |        |  Kubernetes   |
           |    Queue      |        |   Metadata    |        |     API       |
           +-------+-------+        +---------------+        +-------+-------+
                   |                                                 |
                   |                 +--------------------------------+
                   |                 |
                   v                 v
           +-------+-----------------+-------+
           |      Download Workers           |
           |   (Kubernetes Jobs/Pods)        |
           +----------------+----------------+
                            |
                            v
                   +--------+--------+
                   |   iSCSI/Trident |
                   |   Storage       |
                   +-----------------+
```

### Components

- **API Controller**: REST API server handling channel management, job scheduling, and progress tracking
- **Redis Queue**: Manages download job queue, worker coordination, and progress tracking
- **SQLite**: Stores channel and video metadata locally
- **Download Workers**: Kubernetes Jobs that process videos from the queue using yt-dlp
- **Trident Storage**: Provides persistent iSCSI volumes for video storage

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `REDIS_ADDR` | Redis server address | `localhost:6379` |
| `REDIS_PASSWORD` | Redis password | (empty) |
| `K8S_NAMESPACE` | Kubernetes namespace for jobs | `default` |
| `STORAGE_PATH` | Path to video storage | `/archive` |
| `WORKER_IMAGE` | Docker image for workers | `ytarchive-worker:latest` |
| `MAX_WORKERS` | Maximum concurrent workers | `5` |
| `LOG_LEVEL` | Logging level | `info` |

### ConfigMap Options

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: ytarchive-config
data:
  REDIS_ADDR: "redis:6379"
  K8S_NAMESPACE: "ytarchive"
  STORAGE_PATH: "/archive"
  MAX_WORKERS: "10"
  DOWNLOAD_FORMAT: "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best"
  RATE_LIMIT: "5M"
```

## API Reference

### Health Check

```bash
GET /health
```

Returns service health status.

### Channels

```bash
# Add a channel
POST /api/channels
Content-Type: application/json
{"youtube_url": "https://www.youtube.com/@channelname"}

# List all channels
GET /api/channels

# Get channel details
GET /api/channels/:id

# Trigger sync for a channel
POST /api/channels/:id/sync

# Delete a channel
DELETE /api/channels/:id
```

### Jobs

```bash
# List active jobs
GET /api/jobs
```

### Progress

```bash
# Get overall download progress
GET /api/progress
```

See [docs/api.md](docs/api.md) for complete API documentation with examples.

## Deployment

### Prerequisites

1. **Kubernetes Cluster** - v1.28 or later
2. **NetApp Trident** - For iSCSI storage provisioning
3. **Argo Workflows** - For workflow management (optional)
4. **ArgoCD** - For GitOps deployment (optional)

### Quick Deploy

```bash
# 1. Create namespace
kubectl create namespace ytarchive

# 2. Deploy Redis
kubectl apply -f deploy/kubernetes/redis.yaml

# 3. Deploy storage class (Trident)
kubectl apply -f deploy/kubernetes/storage-class.yaml

# 4. Deploy the controller
kubectl apply -f deploy/kubernetes/deployment.yaml

# 5. Verify deployment
kubectl get pods -n ytarchive
```

See [docs/deployment.md](docs/deployment.md) for step-by-step deployment guide.

## Development

### Project Structure

```
ytarchive/
├── cmd/
│   └── controller/         # Main application entry point
├── internal/
│   ├── api/               # HTTP handlers and routes
│   ├── db/                # SQLite database operations
│   ├── downloader/        # yt-dlp wrapper
│   ├── queue/             # Redis queue management
│   ├── scheduler/         # Kubernetes job scheduler
│   ├── storage/           # Storage management
│   └── youtube/           # YouTube API client
├── deploy/
│   ├── argocd/           # ArgoCD application manifests
│   ├── kubernetes/       # Kubernetes manifests
│   └── workflows/        # Argo Workflow templates
├── tests/
│   ├── mocks/            # Mock implementations
│   └── testdata/         # Test fixtures
├── web/                   # Web UI files
└── docs/                  # Documentation
```

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run integration tests only (requires build tag)
go test -v ./tests/... -tags=integration

# Test with real YouTube (requires network)
INTEGRATION_NETWORK=true go test -v ./tests/... -tags=integration -run TestFetch

# Run with coverage
go test -v -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Building

```bash
# Build locally
go build -o bin/controller ./cmd/controller

# Build container with ko
ko build ./cmd/controller

# Build and push to registry
KO_DOCKER_REPO=myregistry.io/ytarchive ko build ./cmd/controller
```

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [yt-dlp](https://github.com/yt-dlp/yt-dlp) - Video downloading
- [ko](https://ko.build/) - Container building
- [Gin](https://gin-gonic.com/) - HTTP web framework
- [go-redis](https://redis.uptrace.dev/) - Redis client
- [client-go](https://github.com/kubernetes/client-go) - Kubernetes client
