# API Reference

Complete API documentation for the YouTube Channel Archiver service.

## Base URL

```
http://localhost:8080
```

In Kubernetes, the service is typically exposed at:
```
http://ytarchive.ytarchive.svc.cluster.local:8080
```

## Authentication

Currently, the API does not require authentication. In production, it's recommended to use a service mesh or API gateway for authentication.

## Response Format

All responses are in JSON format. Successful responses return the requested data directly, while errors return an error object:

```json
{
  "error": "Error message describing what went wrong"
}
```

## Endpoints

### Health Check

#### GET /health

Check if the service is running.

**Response**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Status Codes**
- `200 OK` - Service is healthy

---

#### GET /ready

Check if the service is ready to accept requests (includes Redis connectivity check).

**Response**
```json
{
  "status": "ready"
}
```

**Error Response**
```json
{
  "status": "not ready",
  "error": "Redis connection failed"
}
```

**Status Codes**
- `200 OK` - Service is ready
- `503 Service Unavailable` - Service is not ready

---

### Channels

#### POST /api/channels

Add a new YouTube channel to track and archive.

**Request Body**
```json
{
  "youtube_url": "https://www.youtube.com/@channelname"
}
```

**Supported URL Formats**
- `https://www.youtube.com/@username`
- `https://www.youtube.com/channel/UC...`
- `https://www.youtube.com/c/channelname`
- `https://www.youtube.com/user/username`

**Response**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "youtube_url": "https://www.youtube.com/@channelname",
  "youtube_id": "channelname",
  "name": "",
  "description": "",
  "video_count": 0,
  "status": "pending",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Status Codes**
- `201 Created` - Channel added successfully
- `400 Bad Request` - Invalid request body or YouTube URL
- `409 Conflict` - Channel already exists

**Example**
```bash
curl -X POST http://localhost:8080/api/channels \
  -H "Content-Type: application/json" \
  -d '{"youtube_url": "https://www.youtube.com/@aperturethinking"}'
```

---

#### GET /api/channels

List all tracked channels.

**Response**
```json
{
  "channels": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "youtube_url": "https://www.youtube.com/@channelname",
      "youtube_id": "channelname",
      "name": "Channel Name",
      "description": "Channel description",
      "video_count": 150,
      "status": "synced",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T12:00:00Z",
      "last_sync_at": "2024-01-15T12:00:00Z"
    }
  ],
  "count": 1
}
```

**Channel Status Values**
- `pending` - Channel added, not yet synced
- `syncing` - Currently fetching video list
- `synced` - Video list fetched, ready for downloads
- `error` - An error occurred during sync

**Status Codes**
- `200 OK` - Success

**Example**
```bash
curl http://localhost:8080/api/channels
```

---

#### GET /api/channels/:id

Get details for a specific channel, including its videos.

**Parameters**
- `id` (path) - Channel UUID

**Response**
```json
{
  "channel": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "youtube_url": "https://www.youtube.com/@channelname",
    "youtube_id": "channelname",
    "name": "Channel Name",
    "description": "Channel description",
    "video_count": 150,
    "status": "synced",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T12:00:00Z",
    "last_sync_at": "2024-01-15T12:00:00Z"
  },
  "videos": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "youtube_id": "dQw4w9WgXcQ",
      "channel_id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Video Title",
      "description": "Video description",
      "duration": 212,
      "status": "downloaded",
      "file_path": "/archive/channelname/dQw4w9WgXcQ.mp4",
      "file_size": 52428800,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T11:00:00Z"
    }
  ]
}
```

**Video Status Values**
- `pending` - Video discovered, not yet downloaded
- `downloading` - Currently being downloaded
- `downloaded` - Successfully downloaded
- `error` - Download failed

**Status Codes**
- `200 OK` - Success
- `404 Not Found` - Channel not found

**Example**
```bash
curl http://localhost:8080/api/channels/550e8400-e29b-41d4-a716-446655440000
```

---

#### POST /api/channels/:id/sync

Trigger a sync operation for a channel. This will:
1. Fetch the latest video list from YouTube
2. Add new videos to the download queue
3. Create Kubernetes Jobs for download workers

**Parameters**
- `id` (path) - Channel UUID

**Response**
```json
{
  "message": "Sync started",
  "job_id": "770e8400-e29b-41d4-a716-446655440002",
  "channel": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "youtube_url": "https://www.youtube.com/@channelname",
    "youtube_id": "channelname",
    "name": "Channel Name",
    "status": "syncing",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T14:00:00Z"
  }
}
```

**Status Codes**
- `202 Accepted` - Sync started successfully
- `404 Not Found` - Channel not found
- `409 Conflict` - Channel is already syncing
- `500 Internal Server Error` - Failed to start sync

**Example**
```bash
curl -X POST http://localhost:8080/api/channels/550e8400-e29b-41d4-a716-446655440000/sync
```

---

#### DELETE /api/channels/:id

Delete a channel and its associated videos from tracking.

**Note**: This does not delete downloaded video files from storage.

**Parameters**
- `id` (path) - Channel UUID

**Response**
```json
{
  "message": "Channel deleted successfully"
}
```

**Status Codes**
- `200 OK` - Channel deleted
- `404 Not Found` - Channel not found

**Example**
```bash
curl -X DELETE http://localhost:8080/api/channels/550e8400-e29b-41d4-a716-446655440000
```

---

### Jobs

#### GET /api/jobs

List all active download jobs.

**Response**
```json
{
  "jobs": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440002",
      "channel_id": "550e8400-e29b-41d4-a716-446655440000",
      "worker_num": 3,
      "status": "running",
      "video_count": 50,
      "downloaded": 25,
      "failed": 2,
      "k8s_job_name": "ytarchive-download-abc123",
      "created_at": "2024-01-15T14:00:00Z",
      "updated_at": "2024-01-15T14:30:00Z"
    }
  ],
  "count": 1
}
```

**Job Status Values**
- `pending` - Job created, waiting for workers
- `running` - Workers actively downloading
- `completed` - All videos processed
- `failed` - Job failed

**Status Codes**
- `200 OK` - Success

**Example**
```bash
curl http://localhost:8080/api/jobs
```

---

### Progress

#### GET /api/progress

Get overall download progress across all channels.

**Response**
```json
{
  "total_channels": 5,
  "total_videos": 750,
  "downloaded_videos": 500,
  "pending_videos": 240,
  "failed_videos": 10,
  "active_jobs": 2,
  "download_percentage": 66.67
}
```

**Fields**
- `total_channels` - Number of tracked channels
- `total_videos` - Total videos discovered across all channels
- `downloaded_videos` - Successfully downloaded videos
- `pending_videos` - Videos waiting to be downloaded
- `failed_videos` - Videos that failed to download
- `active_jobs` - Currently running download jobs
- `download_percentage` - Percentage of videos downloaded (0-100)

**Status Codes**
- `200 OK` - Success

**Example**
```bash
curl http://localhost:8080/api/progress
```

---

## Error Handling

### Common Error Responses

#### 400 Bad Request
```json
{
  "error": "Invalid request body: Key: 'AddChannelRequest.YouTubeURL' Error:Field validation for 'YouTubeURL' failed on the 'required' tag"
}
```

#### 404 Not Found
```json
{
  "error": "Channel not found"
}
```

#### 409 Conflict
```json
{
  "error": "Channel already exists",
  "channel": {
    "id": "existing-channel-id",
    "youtube_id": "channelname"
  }
}
```

#### 500 Internal Server Error
```json
{
  "error": "Failed to save channel"
}
```

---

## Rate Limiting

The API does not implement rate limiting by default. For production deployments, consider using:

- Kubernetes Ingress rate limiting
- Service mesh (Istio, Linkerd) rate limiting
- API gateway rate limiting

---

## WebSocket Events (Future)

A WebSocket endpoint for real-time progress updates is planned:

```
ws://localhost:8080/ws/progress
```

This will stream events such as:
- Download started
- Download completed
- Download failed
- Sync completed

---

## SDK Examples

### Go Client

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

func main() {
    // Add a channel
    payload := map[string]string{
        "youtube_url": "https://www.youtube.com/@channelname",
    }
    body, _ := json.Marshal(payload)

    resp, err := http.Post(
        "http://localhost:8080/api/channels",
        "application/json",
        bytes.NewBuffer(body),
    )
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Printf("Status: %d\n", resp.StatusCode)
}
```

### Python Client

```python
import requests

# Add a channel
response = requests.post(
    "http://localhost:8080/api/channels",
    json={"youtube_url": "https://www.youtube.com/@channelname"}
)
print(response.json())

# Get progress
progress = requests.get("http://localhost:8080/api/progress")
print(f"Downloaded: {progress.json()['downloaded_videos']} videos")
```

### curl Examples

```bash
# Add multiple channels
for channel in aperturethinking techlinked mkbhd; do
  curl -X POST http://localhost:8080/api/channels \
    -H "Content-Type: application/json" \
    -d "{\"youtube_url\": \"https://www.youtube.com/@$channel\"}"
  echo
done

# Trigger sync for all channels
curl -s http://localhost:8080/api/channels | \
  jq -r '.channels[].id' | \
  xargs -I {} curl -X POST http://localhost:8080/api/channels/{}/sync

# Watch progress
watch -n 5 'curl -s http://localhost:8080/api/progress | jq .'
```

---

## OpenAPI Specification

An OpenAPI 3.0 specification is available at:

```
GET /api/openapi.yaml
```

You can use this with tools like Swagger UI or generate client libraries.
