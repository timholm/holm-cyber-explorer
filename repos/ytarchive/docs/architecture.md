# Architecture Deep-Dive

Comprehensive documentation of the YouTube Channel Archiver architecture.

## Table of Contents

1. [Overview](#overview)
2. [Component Interaction](#component-interaction)
3. [Data Flow](#data-flow)
4. [Worker Lifecycle](#worker-lifecycle)
5. [Storage Structure](#storage-structure)
6. [Queue Design](#queue-design)
7. [Scaling](#scaling)
8. [Fault Tolerance](#fault-tolerance)

---

## Overview

The YouTube Channel Archiver is designed as a cloud-native, Kubernetes-first application with the following design principles:

- **Horizontally scalable**: Workers scale based on queue depth
- **Fault tolerant**: Jobs can be retried, state is persisted
- **Observable**: Metrics, logs, and status endpoints
- **GitOps friendly**: Declarative configuration, ArgoCD integration

### High-Level Architecture

```
+----------------------------------------------------------+
|                      User Interface                       |
|  +------------+  +------------+  +-------------------+   |
|  |  Web UI    |  |  CLI       |  |  API Clients      |   |
|  +-----+------+  +-----+------+  +--------+----------+   |
|        |               |                  |              |
+--------|---------------|------------------|---------------+
         |               |                  |
         v               v                  v
+----------------------------------------------------------+
|                    API Controller                         |
|  +------------+  +------------+  +-------------------+   |
|  |  Gin HTTP  |  |  Handlers  |  |  Scheduler        |   |
|  |  Router    |  |            |  |                   |   |
|  +-----+------+  +-----+------+  +--------+----------+   |
|        |               |                  |              |
+--------|---------------|------------------|---------------+
         |               |                  |
    +----+----+     +----+----+       +-----+-----+
    |         |     |         |       |           |
    v         v     v         v       v           v
+-------+ +-------+ +-------+ +-------+ +------------------+
| Redis | |SQLite | | Redis | | K8s   | | Trident/Storage  |
| Cache | |  DB   | | Queue | | API   | |                  |
+-------+ +-------+ +-------+ +---+---+ +--------+---------+
                                  |              |
                                  v              v
                        +---------+--------+  +--+--+
                        |  Worker Pods     |  | PVC |
                        |  (K8s Jobs)      +->|     |
                        +------------------+  +-----+
```

---

## Component Interaction

### API Controller

The API Controller is the central orchestration component:

```
                        +------------------------+
                        |    API Controller      |
                        +------------------------+
                        |                        |
                        |  +------------------+  |
                        |  |   HTTP Server    |  |
                        |  |   (Gin/Port 8080)|  |
                        |  +--------+---------+  |
                        |           |            |
                        |  +--------v---------+  |
                        |  |    Handlers      |  |
                        |  |  - AddChannel    |  |
                        |  |  - ListChannels  |  |
                        |  |  - SyncChannel   |  |
                        |  |  - GetProgress   |  |
                        |  +--------+---------+  |
                        |           |            |
          +-------------+-----------+------------+--------------+
          |             |           |            |              |
          v             v           v            v              v
     +--------+    +--------+  +--------+  +---------+   +----------+
     | Redis  |    | SQLite |  | Redis  |  | K8s API |   | YouTube  |
     | Client |    |   DB   |  | Queue  |  | Client  |   |   API    |
     +--------+    +--------+  +--------+  +---------+   +----------+
```

### Scheduler

The Scheduler manages Kubernetes Job creation and worker coordination:

```go
type Scheduler struct {
    k8sClient     *kubernetes.Clientset
    redisClient   *redis.Client
    namespace     string
    workerImage   string
    maxWorkers    int
}

// StartSync initiates a channel sync operation
func (s *Scheduler) StartSync(ctx context.Context, channelID, youtubeID string) (string, error) {
    // 1. Fetch video list from YouTube
    // 2. Add videos to download queue
    // 3. Create Kubernetes Jobs for workers
    // 4. Return job ID
}
```

### Service Communication

```
+------------------+     REST/HTTP      +------------------+
|   Web Client     +<------------------>+  API Controller  |
+------------------+                    +--------+---------+
                                                 |
                    +----------------------------+
                    |              |             |
                    v              v             v
            +-------+------+ +-----+-----+ +----+------+
            |    Redis     | |  SQLite   | | K8s API   |
            |  TCP:6379    | |  File DB  | | In-Cluster|
            +--------------+ +-----------+ +-----------+
```

---

## Data Flow

### Channel Addition Flow

```
1. User Request                 2. Validation              3. Storage
+-------------+              +---------------+           +-----------+
| POST        |              | Extract       |           | Store in  |
| /api/       +------------->+ YouTube ID    +---------->+ Redis     |
| channels    |              | from URL      |           | and DB    |
+-------------+              +---------------+           +-----------+
                                    |
                                    v
                             +---------------+
                             | Check for     |
                             | Duplicates    |
                             +---------------+
```

### Sync Flow

```
+--------+     +----------+     +---------+     +--------+     +--------+
| Sync   |     | Fetch    |     | Parse   |     | Queue  |     | Create |
| Request|---->| YouTube  |---->| Video   |---->| Videos |---->| K8s    |
|        |     | API      |     | List    |     | Redis  |     | Jobs   |
+--------+     +----------+     +---------+     +--------+     +--------+
                                                     |
                                                     v
                                              +-------------+
                                              | Workers     |
                                              | Claim &     |
                                              | Download    |
                                              +-------------+
```

### Download Flow

```
+----------+     +----------+     +---------+     +----------+     +--------+
| Worker   |     | Claim    |     | Run     |     | Save     |     | Update |
| Starts   |---->| Video    |---->| yt-dlp  |---->| to       |---->| Status |
|          |     | from Q   |     |         |     | Storage  |     | Redis  |
+----------+     +----------+     +---------+     +----------+     +--------+
     |                                                                  |
     +------------------------------------------------------------------+
                            (Loop until queue empty)
```

---

## Worker Lifecycle

### Worker States

```
                    +-------------------+
                    |                   |
           +------->+     PENDING       |
           |        |                   |
           |        +--------+----------+
           |                 |
           |                 v
           |        +--------+----------+
           |        |                   |
           |        |     RUNNING       +--------+
           |        |                   |        |
           |        +--------+----------+        |
           |                 |                   |
           |    +------------+------------+      |
           |    |            |            |      |
           |    v            v            v      |
        +--+----+--+  +------+------+  +--+------+---+
        |          |  |             |  |             |
        | COMPLETED|  |   FAILED    |  |  TIMEOUT    |
        |          |  |             |  |             |
        +----------+  +------+------+  +-------------+
                             |
                             v
                      +------+------+
                      |   RETRY     |
                      | (if < max)  |
                      +-------------+
```

### Worker Pod Specification

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: ytarchive-download-{job-id}
  namespace: ytarchive
spec:
  backoffLimit: 3
  activeDeadlineSeconds: 3600
  ttlSecondsAfterFinished: 300
  template:
    spec:
      restartPolicy: OnFailure
      containers:
        - name: worker
          image: ytarchive-worker:latest
          env:
            - name: REDIS_ADDR
              value: "redis:6379"
            - name: JOB_ID
              value: "{job-id}"
            - name: CHANNEL_ID
              value: "{channel-id}"
          volumeMounts:
            - name: archive
              mountPath: /archive
          resources:
            requests:
              memory: "512Mi"
              cpu: "500m"
            limits:
              memory: "2Gi"
              cpu: "2"
      volumes:
        - name: archive
          persistentVolumeClaim:
            claimName: ytarchive-storage
```

### Worker Process Flow

```go
func WorkerMain() {
    // 1. Initialize connections
    redis := connectRedis()

    // 2. Process queue
    for {
        // Claim a video from queue
        videoID, err := redis.LPop(ctx, "download:queue")
        if err == redis.Nil {
            break // Queue empty
        }

        // Mark as processing
        redis.SAdd(ctx, "download:processing", videoID)

        // Download video
        err = downloadVideo(videoID)

        if err != nil {
            // Mark as failed
            redis.SAdd(ctx, "download:failed", videoID)
            redis.SRem(ctx, "download:processing", videoID)
            continue
        }

        // Mark as completed
        redis.SAdd(ctx, "download:completed", videoID)
        redis.SRem(ctx, "download:processing", videoID)
    }

    // 3. Exit gracefully
}
```

---

## Storage Structure

### Directory Layout

```
/archive/
├── {channel-id}/
│   ├── metadata.json           # Channel metadata
│   ├── videos/
│   │   ├── {video-id}/
│   │   │   ├── video.mp4       # Video file
│   │   │   ├── video.info.json # yt-dlp metadata
│   │   │   ├── thumbnail.jpg   # Video thumbnail
│   │   │   └── subtitles/      # Subtitle files
│   │   │       ├── en.vtt
│   │   │       └── es.vtt
│   │   └── {video-id}/
│   │       └── ...
│   └── playlists/
│       └── uploads.json        # Playlist data
└── {channel-id}/
    └── ...
```

### Storage Volume Architecture

```
+-------------------+
|  Kubernetes PVC   |
|  ytarchive-storage|
+--------+----------+
         |
         v
+--------+----------+
|  Trident CSI      |
|  (NetApp Driver)  |
+--------+----------+
         |
         v
+--------+----------+
|  iSCSI LUN        |
|  (NetApp ONTAP)   |
+-------------------+
```

### File Naming Convention

```
{channel-id}/{video-id}.{ext}

Examples:
- UC_aperturethinking/dQw4w9WgXcQ.mp4
- UC_aperturethinking/dQw4w9WgXcQ.info.json
- UC_aperturethinking/dQw4w9WgXcQ.en.vtt
```

---

## Queue Design

### Redis Data Structures

```
+------------------------------------------+
|              Redis Keys                   |
+------------------------------------------+

Channels:
  channels                  SET     Channel IDs
  channel:{id}             STRING  Channel JSON

Videos:
  video:{channel}:{video}  STRING  Video JSON

Queue:
  download:queue           LIST    Pending video IDs
  download:processing      SET     Currently downloading
  download:completed       SET     Successfully downloaded
  download:failed          SET     Failed downloads

Jobs:
  jobs:active              SET     Active job IDs
  job:{id}                 STRING  Job JSON

Progress:
  progress:{channel}       HASH    Download statistics
```

### Queue Operations

```
Push Video to Queue:
  RPUSH download:queue {video-id}

Claim Video (Atomic):
  video = LPOP download:queue
  SADD download:processing {video}

Complete Video:
  SREM download:processing {video}
  SADD download:completed {video}

Fail Video:
  SREM download:processing {video}
  SADD download:failed {video}
  RPUSH download:queue {video}  # Retry
```

### Queue Monitoring

```go
type QueueStats struct {
    Pending     int64  // LLEN download:queue
    Processing  int64  // SCARD download:processing
    Completed   int64  // SCARD download:completed
    Failed      int64  // SCARD download:failed
}

func GetQueueStats(ctx context.Context, redis *redis.Client) QueueStats {
    return QueueStats{
        Pending:    redis.LLen(ctx, "download:queue").Val(),
        Processing: redis.SCard(ctx, "download:processing").Val(),
        Completed:  redis.SCard(ctx, "download:completed").Val(),
        Failed:     redis.SCard(ctx, "download:failed").Val(),
    }
}
```

---

## Scaling

### Horizontal Pod Autoscaling

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: ytarchive-workers
  namespace: ytarchive
spec:
  scaleTargetRef:
    apiVersion: batch/v1
    kind: Job
    name: ytarchive-download
  minReplicas: 1
  maxReplicas: 20
  metrics:
    - type: External
      external:
        metric:
          name: redis_queue_length
          selector:
            matchLabels:
              queue: download
        target:
          type: AverageValue
          averageValue: "10"
```

### Worker Scaling Logic

```go
func (s *Scheduler) calculateWorkerCount(queueLength int) int {
    // Base: 1 worker per 50 videos
    workers := queueLength / 50

    // Minimum workers
    if workers < 1 {
        workers = 1
    }

    // Maximum workers
    if workers > s.maxWorkers {
        workers = s.maxWorkers
    }

    return workers
}
```

### Resource Recommendations

| Queue Size | Workers | Memory per Worker | CPU per Worker |
|------------|---------|------------------|----------------|
| 1-50       | 1       | 512Mi            | 500m           |
| 51-200     | 2-4     | 512Mi            | 500m           |
| 201-500    | 5-10    | 1Gi              | 1              |
| 500+       | 10-20   | 2Gi              | 2              |

---

## Fault Tolerance

### Retry Strategy

```go
const (
    MaxRetries     = 3
    RetryDelay     = 5 * time.Minute
    BackoffFactor  = 2.0
)

func (w *Worker) downloadWithRetry(videoID string) error {
    var lastErr error

    for attempt := 0; attempt < MaxRetries; attempt++ {
        err := w.download(videoID)
        if err == nil {
            return nil
        }

        lastErr = err
        delay := RetryDelay * time.Duration(math.Pow(BackoffFactor, float64(attempt)))
        time.Sleep(delay)
    }

    return fmt.Errorf("failed after %d attempts: %w", MaxRetries, lastErr)
}
```

### Failure Recovery

```
+------------------+
|  Video Download  |
|     Failure      |
+--------+---------+
         |
         v
+--------+---------+
|  Check Retry     |
|  Count           |
+--------+---------+
         |
    +----+----+
    |         |
    v         v
+---+---+ +---+----+
| Retry | | Dead   |
| Queue | | Letter |
+---+---+ +--------+
    |
    v
+---+--------+
| Exponential|
| Backoff    |
+------------+
```

### Health Checks

```yaml
# Liveness: Is the process alive?
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 30
  failureThreshold: 3

# Readiness: Can we handle requests?
readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
  failureThreshold: 3
```

### Data Durability

```
+------------------+     +------------------+     +------------------+
|     In-Memory    |     |   Persistent     |     |   Backed Up      |
|       Data       |     |     Storage      |     |     Storage      |
+------------------+     +------------------+     +------------------+
|  - Queue state   |     |  - Redis AOF     |     |  - Video files   |
|  - Progress      | --> |  - SQLite DB     | --> |  - ONTAP snaps   |
|  - Temp files    |     |  - Config files  |     |  - Off-site copy |
+------------------+     +------------------+     +------------------+

        |                        |                        |
        v                        v                        v
   Recoverable            Durable (Local)         Durable (Remote)
   from source            Trident volumes         Backup/DR
```

---

## Performance Considerations

### Network Optimization

- Use dedicated network namespace for downloads
- Configure rate limiting to avoid throttling
- Use connection pooling for Redis

### Storage Optimization

- Separate volumes for metadata and video files
- Use appropriate storage class for workload
- Consider compression for metadata

### Memory Optimization

- Stream downloads directly to disk
- Limit concurrent downloads per worker
- Use memory limits in pod specs

### Monitoring Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `ytarchive_queue_length` | Gauge | Number of pending downloads |
| `ytarchive_downloads_total` | Counter | Total downloads (by status) |
| `ytarchive_download_duration_seconds` | Histogram | Download time distribution |
| `ytarchive_storage_bytes` | Gauge | Total storage used |
| `ytarchive_active_workers` | Gauge | Number of active workers |
