package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// ProgressKeyPrefix is the Redis key prefix for download progress
	ProgressKeyPrefix = "download:progress:"
	// ActiveDownloadsKey is the Redis SET key for tracking active downloads
	ActiveDownloadsKey = "download:active"
	// ProgressTTL is the TTL for progress keys (5 minutes)
	ProgressTTL = 5 * time.Minute
	// MinProgressInterval is the minimum interval between progress updates
	MinProgressInterval = 1 * time.Second
)

// RedisProgressReporter reports download progress to Redis
type RedisProgressReporter struct {
	client     *redis.Client
	workerID   string
	lastUpdate map[string]time.Time // Track last update time per video
	ctx        context.Context
}

// NewRedisProgressReporter creates a new RedisProgressReporter
func NewRedisProgressReporter(client *redis.Client, workerID string) *RedisProgressReporter {
	return &RedisProgressReporter{
		client:     client,
		workerID:   workerID,
		lastUpdate: make(map[string]time.Time),
		ctx:        context.Background(),
	}
}

// progressKey returns the Redis key for a video's progress
func progressKey(videoID string) string {
	return ProgressKeyPrefix + videoID
}

// ReportProgress writes progress to Redis
// Progress is throttled to at most once per second per video (unless force is true)
func (r *RedisProgressReporter) ReportProgress(progress *DownloadProgress, force bool) error {
	if r.client == nil {
		return nil
	}

	// Throttle updates to at most once per second (unless force is true)
	if !force {
		if lastTime, ok := r.lastUpdate[progress.VideoID]; ok {
			if time.Since(lastTime) < MinProgressInterval {
				return nil
			}
		}
	}

	// Set worker ID and timestamp
	progress.WorkerID = r.workerID
	progress.UpdatedAt = time.Now().Unix()

	// Serialize progress to JSON
	data, err := json.Marshal(progress)
	if err != nil {
		return fmt.Errorf("failed to marshal progress: %w", err)
	}

	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	// Use a pipeline to execute multiple commands atomically
	pipe := r.client.Pipeline()

	// Set progress with TTL
	key := progressKey(progress.VideoID)
	pipe.Set(ctx, key, data, ProgressTTL)

	// Add to active downloads set (unless completed or error)
	if progress.Status == "downloading" || progress.Status == "processing" {
		pipe.SAdd(ctx, ActiveDownloadsKey, progress.VideoID)
		// Also set a TTL on the set entry tracking (we'll refresh it with each update)
		// Note: Redis SETs don't support per-member TTL, so we use the progress key TTL
		// to implicitly track active downloads
	} else {
		// Remove from active downloads when completed or errored
		pipe.SRem(ctx, ActiveDownloadsKey, progress.VideoID)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to write progress to Redis: %w", err)
	}

	r.lastUpdate[progress.VideoID] = time.Now()
	return nil
}

// StartDownload marks a download as starting and adds it to active downloads
func (r *RedisProgressReporter) StartDownload(videoID string) error {
	progress := &DownloadProgress{
		VideoID:   videoID,
		WorkerID:  r.workerID,
		Status:    "downloading",
		UpdatedAt: time.Now().Unix(),
	}
	return r.ReportProgress(progress, true)
}

// CompleteDownload marks a download as completed and removes it from active downloads
func (r *RedisProgressReporter) CompleteDownload(videoID string, totalBytes int64) error {
	progress := &DownloadProgress{
		VideoID:         videoID,
		WorkerID:        r.workerID,
		Status:          "completed",
		Percentage:      100,
		DownloadedBytes: totalBytes,
		TotalBytes:      totalBytes,
		UpdatedAt:       time.Now().Unix(),
	}
	return r.ReportProgress(progress, true)
}

// ErrorDownload marks a download as errored and removes it from active downloads
func (r *RedisProgressReporter) ErrorDownload(videoID string, errMsg string) error {
	progress := &DownloadProgress{
		VideoID:   videoID,
		WorkerID:  r.workerID,
		Status:    "error",
		UpdatedAt: time.Now().Unix(),
	}
	return r.ReportProgress(progress, true)
}

// GetProgress retrieves the current progress for a video
func (r *RedisProgressReporter) GetProgress(videoID string) (*DownloadProgress, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	key := progressKey(videoID)
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // No progress found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get progress: %w", err)
	}

	var progress DownloadProgress
	if err := json.Unmarshal([]byte(data), &progress); err != nil {
		return nil, fmt.Errorf("failed to unmarshal progress: %w", err)
	}

	return &progress, nil
}

// GetAllActiveDownloads returns all currently active video IDs
func (r *RedisProgressReporter) GetAllActiveDownloads() ([]string, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	videoIDs, err := r.client.SMembers(ctx, ActiveDownloadsKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get active downloads: %w", err)
	}

	return videoIDs, nil
}

// GetAllProgress returns progress for all active downloads
func (r *RedisProgressReporter) GetAllProgress() ([]*DownloadProgress, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	// Get all active video IDs
	videoIDs, err := r.client.SMembers(ctx, ActiveDownloadsKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get active downloads: %w", err)
	}

	if len(videoIDs) == 0 {
		return []*DownloadProgress{}, nil
	}

	// Build keys for all active videos
	keys := make([]string, len(videoIDs))
	for i, id := range videoIDs {
		keys[i] = progressKey(id)
	}

	// Fetch all progress data in one MGET call
	results, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get progress data: %w", err)
	}

	// Parse results
	progressList := make([]*DownloadProgress, 0, len(results))
	for i, result := range results {
		if result == nil {
			// Progress key expired but video is still in active set - clean it up
			r.client.SRem(ctx, ActiveDownloadsKey, videoIDs[i])
			continue
		}

		data, ok := result.(string)
		if !ok {
			continue
		}

		var progress DownloadProgress
		if err := json.Unmarshal([]byte(data), &progress); err != nil {
			continue
		}

		progressList = append(progressList, &progress)
	}

	return progressList, nil
}

// CleanupStaleEntries removes stale entries from the active downloads set
// (entries where the progress key has expired)
func (r *RedisProgressReporter) CleanupStaleEntries() error {
	ctx, cancel := context.WithTimeout(r.ctx, 10*time.Second)
	defer cancel()

	// Get all active video IDs
	videoIDs, err := r.client.SMembers(ctx, ActiveDownloadsKey).Result()
	if err != nil {
		return fmt.Errorf("failed to get active downloads: %w", err)
	}

	for _, videoID := range videoIDs {
		key := progressKey(videoID)
		exists, err := r.client.Exists(ctx, key).Result()
		if err != nil {
			continue
		}
		if exists == 0 {
			// Progress key expired, remove from active set
			r.client.SRem(ctx, ActiveDownloadsKey, videoID)
		}
	}

	return nil
}

// ClearProgress removes progress data for a video (used after completion/error)
func (r *RedisProgressReporter) ClearProgress(videoID string) error {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	key := progressKey(videoID)
	pipe := r.client.Pipeline()
	pipe.Del(ctx, key)
	pipe.SRem(ctx, ActiveDownloadsKey, videoID)
	_, err := pipe.Exec(ctx)
	return err
}

// Close cleans up resources (currently a no-op as we don't own the Redis client)
func (r *RedisProgressReporter) Close() error {
	return nil
}
