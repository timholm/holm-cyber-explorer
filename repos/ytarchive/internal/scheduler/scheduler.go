package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"k8s.io/client-go/kubernetes"

	"github.com/timholm/ytarchive/internal/logging"
	"github.com/timholm/ytarchive/internal/youtube"
)

const (
	// Redis key prefixes
	channelKeyPrefix    = "channel:"
	videoKeyPrefix      = "video:"
	syncJobKeyPrefix    = "syncjob:"
	videoQueueKeyPrefix = "queue:videos:"
	// Unified queue for all channels - KEDA monitors this
	unifiedQueueKey = "ytarchive:download:queue"
)

// SyncJob represents a full channel sync operation
type SyncJob struct {
	ID         string    `json:"id"`
	ChannelID  string    `json:"channel_id"`
	YouTubeID  string    `json:"youtube_id"`
	Status     string    `json:"status"` // pending, discovering, running, completed, failed
	VideoCount int       `json:"video_count"`
	Downloaded int       `json:"downloaded"`
	Failed     int       `json:"failed"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Scheduler manages job scheduling and worker coordination
// With KEDA-based scaling, the scheduler only pushes videos to Redis queues.
// KEDA monitors queue length and scales workers automatically.
type Scheduler struct {
	k8sClient     *kubernetes.Clientset
	redis         *redis.Client
	namespace     string
	k8sManager    *K8sJobManager // Kept for cleanup operations
	youtubeClient *youtube.Client
	mu            sync.Mutex
}

// NewScheduler creates a new Scheduler instance
func NewScheduler(k8sClient *kubernetes.Clientset, redisClient *redis.Client, namespace string) *Scheduler {
	// Create YouTube client for video discovery
	ytClient, err := youtube.NewClient()
	if err != nil {
		logging.Warn("failed to create YouTube client, video discovery will be limited", "error", err)
	}

	s := &Scheduler{
		k8sClient:     k8sClient,
		redis:         redisClient,
		namespace:     namespace,
		k8sManager:    NewK8sJobManager(k8sClient, namespace), // Kept for cleanup operations
		youtubeClient: ytClient,
	}

	// Recover any channels stuck in "syncing" state from previous controller instance
	go s.recoverStuckChannels()

	return s
}

// recoverStuckChannels checks for channels stuck in "syncing" state and fixes them
func (s *Scheduler) recoverStuckChannels() {
	ctx := context.Background()

	// Give the system a moment to initialize
	time.Sleep(5 * time.Second)

	logging.Info("checking for stuck channels to recover")

	// Get all channel IDs
	channelIDs, err := s.redis.SMembers(ctx, "channels").Result()
	if err != nil {
		logging.Warn("failed to get channel list for recovery", "error", err)
		return
	}

	for _, channelID := range channelIDs {
		channelData, err := s.redis.Get(ctx, channelKeyPrefix+channelID).Result()
		if err != nil {
			continue
		}

		var channelMap map[string]interface{}
		if err := json.Unmarshal([]byte(channelData), &channelMap); err != nil {
			continue
		}

		status, _ := channelMap["status"].(string)
		if status != "syncing" {
			continue
		}

		logging.Info("found stuck channel, checking if sync is complete", "channel_id", channelID)

		// Check video statuses to see if all are processed
		// With unified queue, we check video status rather than per-channel queue length
		allProcessed := s.checkAllVideosProcessed(ctx, channelID)

		if allProcessed {
			logging.Info("recovering stuck channel - marking as synced", "channel_id", channelID)
			s.updateChannelStatus(ctx, channelID, "synced")
		} else {
			logging.Info("channel still has pending work",
				"channel_id", channelID,
				"all_processed", allProcessed,
			)
		}
	}

	logging.Info("finished checking for stuck channels")
}

// StartSync initiates a full sync for a channel.
// With KEDA-based scaling:
// 1. Discovers videos from YouTube
// 2. Pushes video IDs to Redis queue (queue:videos:{channelID})
// 3. KEDA monitors the queue and scales workers automatically
// 4. Workers pull videos from queue and download them
// 5. When queue is empty, KEDA scales workers to 0
func (s *Scheduler) StartSync(ctx context.Context, channelID, youtubeID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create sync job record
	syncJobID := uuid.New().String()
	syncJob := SyncJob{
		ID:        syncJobID,
		ChannelID: channelID,
		YouTubeID: youtubeID,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save sync job to Redis
	syncJobJSON, err := json.Marshal(syncJob)
	if err != nil {
		return "", fmt.Errorf("failed to marshal sync job: %w", err)
	}
	if err := s.redis.Set(ctx, syncJobKeyPrefix+syncJobID, syncJobJSON, 24*time.Hour).Err(); err != nil {
		return "", fmt.Errorf("failed to save sync job: %w", err)
	}

	// Start async sync process
	go s.executeSyncJob(syncJobID, channelID, youtubeID)

	logging.Info("sync job created",
		"job_id", syncJobID,
		"channel_id", channelID,
	)
	return syncJobID, nil
}

// executeSyncJob performs the actual sync operation asynchronously.
// With KEDA-based scaling and streaming discovery, videos are pushed to the queue
// incrementally as they're discovered. KEDA monitors the queue and scales workers automatically.
func (s *Scheduler) executeSyncJob(syncJobID, channelID, youtubeID string) {
	ctx := context.Background()

	// Update channel status to syncing
	s.updateChannelStatus(ctx, channelID, "syncing")

	// Update status to discovering
	s.updateSyncJobStatus(ctx, syncJobID, "discovering")

	// Discover videos from YouTube using streaming pagination.
	// Videos are pushed to the queue incrementally as they're discovered,
	// allowing workers to start downloading before discovery is complete.
	videoIDs, err := s.discoverAndSaveVideos(ctx, channelID, youtubeID)
	if err != nil {
		logging.Error("error discovering videos",
			"job_id", syncJobID,
			"channel_id", channelID,
			"error", err,
		)
		s.updateSyncJobStatus(ctx, syncJobID, "failed")
		s.updateChannelStatus(ctx, channelID, "error")
		return
	}

	videoCount := len(videoIDs)
	if videoCount == 0 {
		logging.Info("no new videos to download for channel",
			"job_id", syncJobID,
			"channel_id", channelID,
		)
		s.updateSyncJobStatus(ctx, syncJobID, "completed")
		s.updateChannelStatus(ctx, channelID, "synced")
		return
	}

	// Videos have already been pushed to queue during streaming discovery
	logging.Info("discovery complete - videos already queued, KEDA scaling workers",
		"job_id", syncJobID,
		"channel_id", channelID,
		"video_count", videoCount,
	)

	// Update sync job with video count
	s.updateSyncJobVideoCount(ctx, syncJobID, videoCount)
	s.updateSyncJobStatus(ctx, syncJobID, "running")

	// Start monitoring queue completion
	// Workers are scaled by KEDA, we just monitor progress
	go s.monitorSyncProgress(ctx, syncJobID, channelID, videoCount)
}

// discoverAndSaveVideos fetches videos from YouTube using streaming pagination.
// Instead of loading all videos into memory at once, it:
// 1. First pass: counts new videos (lightweight, needed for episode numbering)
// 2. Second pass: streams through pages, saves to Redis, and pushes to queue after each batch
// This allows workers to start downloading while discovery is still in progress.
func (s *Scheduler) discoverAndSaveVideos(ctx context.Context, channelID, youtubeID string) ([]string, error) {
	if s.youtubeClient == nil {
		logging.Warn("YouTube client not available, falling back to Redis-only video discovery",
			"channel_id", channelID,
		)
		return s.getChannelVideoIDs(ctx, channelID)
	}

	logging.Info("discovering videos from YouTube (streaming mode)",
		"channel_id", channelID,
		"youtube_id", youtubeID,
	)

	// Phase 1: Count new videos and collect their IDs for ordering
	// This is a lightweight pass - we only check existence, not save data
	newVideoIDs, requeueVideoIDs, err := s.countNewVideosStreaming(ctx, channelID, youtubeID)
	if err != nil {
		logging.Warn("error counting videos from YouTube, falling back to Redis",
			"channel_id", channelID,
			"youtube_id", youtubeID,
			"error", err,
		)
		return s.getChannelVideoIDs(ctx, channelID)
	}

	totalNew := len(newVideoIDs)
	totalRequeue := len(requeueVideoIDs)

	logging.Info("video discovery count complete",
		"channel_id", channelID,
		"new_videos", totalNew,
		"requeue_videos", totalRequeue,
	)

	if totalNew == 0 && totalRequeue == 0 {
		return []string{}, nil
	}

	// Phase 2: Stream through and save new videos with correct episode numbers
	// Videos from YouTube come newest-first, but newVideoIDs is in that order
	// We want oldest = episode 1, so we reverse the order for numbering
	// Use unified queue so KEDA can monitor a single queue for all channels
	channelName := s.getChannelName(ctx, channelID)
	existingEpisodeCount := s.getExistingEpisodeCount(ctx, channelID)
	now := time.Now()

	// Create a map from video ID to episode number (oldest new video = existingCount + 1)
	episodeMap := make(map[string]int, totalNew)
	for i, videoID := range newVideoIDs {
		// newVideoIDs is newest-first, so reverse the numbering
		episodeMap[videoID] = existingEpisodeCount + totalNew - i
	}

	// Now stream through pages again, this time saving and pushing to queue
	var pageToken string
	var processedNew, processedRequeue int
	var allVideoIDs []string

	for {
		videos, nextToken, err := s.youtubeClient.GetVideoListPaginatedContext(ctx, youtubeID, 0, pageToken)
		if err != nil {
			logging.Error("error fetching video page during save phase",
				"channel_id", channelID,
				"error", err,
			)
			break
		}

		var batchIDs []string

		for _, video := range videos {
			videoKey := videoKeyPrefix + channelID + ":" + video.ID

			// Check if this is a new video we need to save
			if episodeNum, isNew := episodeMap[video.ID]; isNew {
				// Save new video with episode number
				videoData := map[string]interface{}{
					"id":             video.ID,
					"youtube_id":     video.ID,
					"channel_id":     channelID,
					"title":          video.Title,
					"description":    video.Description,
					"duration":       video.Duration,
					"upload_date":    video.UploadDate,
					"thumbnail_url":  video.ThumbnailURL,
					"view_count":     video.ViewCount,
					"episode_number": episodeNum,
					"channel_name":   channelName,
					"status":         "pending",
					"created_at":     now,
					"updated_at":     now,
				}
				videoJSON, err := json.Marshal(videoData)
				if err != nil {
					logging.Warn("error marshaling video",
						"channel_id", channelID,
						"video_id", video.ID,
						"error", err,
					)
					continue
				}
				if err := s.redis.Set(ctx, videoKey, videoJSON, 0).Err(); err != nil {
					logging.Warn("error saving video to Redis",
						"channel_id", channelID,
						"video_id", video.ID,
						"error", err,
					)
					continue
				}
				batchIDs = append(batchIDs, video.ID)
				processedNew++
				delete(episodeMap, video.ID) // Mark as processed
			} else {
				// Check if this is an existing video that needs requeue
				for _, requeueID := range requeueVideoIDs {
					if requeueID == video.ID {
						batchIDs = append(batchIDs, video.ID)
						processedRequeue++
						break
					}
				}
			}
		}

		// Push this batch to unified queue immediately - workers can start downloading now
		// Format: channelID:videoID so workers know which channel the video belongs to
		if len(batchIDs) > 0 {
			for _, id := range batchIDs {
				queueItem := channelID + ":" + id
				if err := s.redis.LPush(ctx, unifiedQueueKey, queueItem).Err(); err != nil {
					logging.Error("error pushing video to queue",
						"channel_id", channelID,
						"video_id", id,
						"error", err,
					)
				}
			}
			allVideoIDs = append(allVideoIDs, batchIDs...)

			logging.Info("pushed video batch to queue",
				"channel_id", channelID,
				"batch_size", len(batchIDs),
				"total_queued", len(allVideoIDs),
				"queue_key", unifiedQueueKey,
			)
		}

		if nextToken == "" {
			break
		}
		pageToken = nextToken
	}

	logging.Info("streaming discovery complete",
		"channel_id", channelID,
		"new_saved", processedNew,
		"requeued", processedRequeue,
		"total_queued", len(allVideoIDs),
	)

	return allVideoIDs, nil
}

// countNewVideosStreaming does a lightweight streaming pass to count new videos.
// Returns two slices: new video IDs (in discovery order, newest first) and requeue video IDs.
func (s *Scheduler) countNewVideosStreaming(ctx context.Context, channelID, youtubeID string) (newVideoIDs, requeueVideoIDs []string, err error) {
	var pageToken string
	pageCount := 0

	for {
		videos, nextToken, err := s.youtubeClient.GetVideoListPaginatedContext(ctx, youtubeID, 0, pageToken)
		if err != nil {
			return nil, nil, fmt.Errorf("error fetching video page %d: %w", pageCount, err)
		}

		pageCount++

		for _, video := range videos {
			videoKey := videoKeyPrefix + channelID + ":" + video.ID

			exists, err := s.redis.Exists(ctx, videoKey).Result()
			if err != nil {
				logging.Warn("error checking video existence",
					"channel_id", channelID,
					"video_id", video.ID,
					"error", err,
				)
				continue
			}

			if exists == 0 {
				// New video
				newVideoIDs = append(newVideoIDs, video.ID)
			} else {
				// Check if it needs requeue (pending or error status)
				videoData, err := s.redis.Get(ctx, videoKey).Result()
				if err != nil {
					continue
				}
				var existingVideo map[string]interface{}
				if err := json.Unmarshal([]byte(videoData), &existingVideo); err != nil {
					continue
				}
				status, _ := existingVideo["status"].(string)
				if status == "pending" || status == "error" {
					requeueVideoIDs = append(requeueVideoIDs, video.ID)
				}
			}
		}

		logging.Info("counting videos - page processed",
			"channel_id", channelID,
			"page", pageCount,
			"videos_in_page", len(videos),
			"new_so_far", len(newVideoIDs),
			"requeue_so_far", len(requeueVideoIDs),
		)

		if nextToken == "" {
			break
		}
		pageToken = nextToken
	}

	return newVideoIDs, requeueVideoIDs, nil
}

// monitorSyncProgress monitors queue length and video statuses for a sync operation.
// With KEDA-based scaling, workers are managed automatically.
// We just need to check when the queue is empty and all videos are processed.
func (s *Scheduler) monitorSyncProgress(ctx context.Context, syncJobID, channelID string, totalVideos int) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			queueLen, downloaded, failed := s.checkSyncProgress(ctx, channelID)

			logging.Info("sync progress",
				"job_id", syncJobID,
				"channel_id", channelID,
				"queue_length", queueLen,
				"downloaded", downloaded,
				"failed", failed,
				"total", totalVideos,
			)

			// Update sync job progress
			s.updateSyncJobProgress(ctx, syncJobID, downloaded, failed)

			// Check if sync is complete: queue is empty AND all videos have been processed
			if queueLen == 0 {
				// Verify all videos are actually processed (not just queue empty)
				if s.checkAllVideosProcessed(ctx, channelID) {
					if failed > 0 && downloaded == 0 {
						s.updateSyncJobStatus(ctx, syncJobID, "failed")
						s.updateChannelStatus(ctx, channelID, "error")
					} else {
						s.updateSyncJobStatus(ctx, syncJobID, "completed")
						s.updateChannelStatus(ctx, channelID, "synced")
					}

					logging.Info("sync completed",
						"job_id", syncJobID,
						"channel_id", channelID,
						"downloaded", downloaded,
						"failed", failed,
					)
					return
				}
			}
		}
	}
}

// checkSyncProgress checks the unified queue length and video statuses for a channel
func (s *Scheduler) checkSyncProgress(ctx context.Context, channelID string) (queueLen int64, downloaded, failed int) {
	// Check unified queue length (shared across all channels)
	queueLen, err := s.redis.LLen(ctx, unifiedQueueKey).Result()
	if err != nil {
		logging.Warn("failed to check queue length", "error", err)
	}

	// Count video statuses for this specific channel
	downloaded, failed = s.countVideoStatuses(ctx, channelID)

	return queueLen, downloaded, failed
}

// countVideoStatuses counts downloaded and failed videos for a channel
func (s *Scheduler) countVideoStatuses(ctx context.Context, channelID string) (downloaded, failed int) {
	// Scan for all videos belonging to this channel
	videoPattern := videoKeyPrefix + channelID + ":*"
	var cursor uint64

	for {
		keys, nextCursor, err := s.redis.Scan(ctx, cursor, videoPattern, 100).Result()
		if err != nil {
			logging.Warn("failed to scan video keys", "channel_id", channelID, "error", err)
			break
		}

		for _, key := range keys {
			videoData, err := s.redis.Get(ctx, key).Result()
			if err != nil {
				continue
			}

			var video map[string]interface{}
			if err := json.Unmarshal([]byte(videoData), &video); err != nil {
				continue
			}

			status, _ := video["status"].(string)
			switch status {
			case "downloaded", "completed":
				downloaded++
			case "failed", "error":
				failed++
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return downloaded, failed
}

// checkAllVideosProcessed verifies that all videos for a channel have been processed
func (s *Scheduler) checkAllVideosProcessed(ctx context.Context, channelID string) bool {
	videoPattern := videoKeyPrefix + channelID + ":*"
	var cursor uint64

	for {
		keys, nextCursor, err := s.redis.Scan(ctx, cursor, videoPattern, 100).Result()
		if err != nil {
			return false
		}

		for _, key := range keys {
			videoData, err := s.redis.Get(ctx, key).Result()
			if err != nil {
				continue
			}

			var video map[string]interface{}
			if err := json.Unmarshal([]byte(videoData), &video); err != nil {
				continue
			}

			status, _ := video["status"].(string)
			// If any video is still pending or downloading, not all are processed
			if status == "pending" || status == "downloading" || status == "queued" {
				return false
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return true
}

// getChannelVideoIDs retrieves video IDs for a channel from Redis
func (s *Scheduler) getChannelVideoIDs(ctx context.Context, channelID string) ([]string, error) {
	videoKeys, err := s.redis.Keys(ctx, videoKeyPrefix+channelID+":*").Result()
	if err != nil {
		return nil, err
	}

	videoIDs := make([]string, 0, len(videoKeys))
	for _, key := range videoKeys {
		// Extract video ID from key (format: video:channelID:videoID)
		videoID := key[len(videoKeyPrefix+channelID+":"):]
		videoIDs = append(videoIDs, videoID)
	}

	return videoIDs, nil
}

// updateSyncJobStatus updates the status of a sync job
func (s *Scheduler) updateSyncJobStatus(ctx context.Context, syncJobID, status string) {
	syncJobData, err := s.redis.Get(ctx, syncJobKeyPrefix+syncJobID).Result()
	if err != nil {
		logging.Warn("error fetching sync job",
			"job_id", syncJobID,
			"error", err,
		)
		return
	}

	var syncJob SyncJob
	if err := json.Unmarshal([]byte(syncJobData), &syncJob); err != nil {
		logging.Warn("error unmarshaling sync job",
			"job_id", syncJobID,
			"error", err,
		)
		return
	}

	syncJob.Status = status
	syncJob.UpdatedAt = time.Now()

	syncJobJSON, _ := json.Marshal(syncJob)
	s.redis.Set(ctx, syncJobKeyPrefix+syncJobID, syncJobJSON, 24*time.Hour)
}

// updateSyncJobVideoCount updates the video count of a sync job
func (s *Scheduler) updateSyncJobVideoCount(ctx context.Context, syncJobID string, videoCount int) {
	syncJobData, err := s.redis.Get(ctx, syncJobKeyPrefix+syncJobID).Result()
	if err != nil {
		return
	}

	var syncJob SyncJob
	if err := json.Unmarshal([]byte(syncJobData), &syncJob); err != nil {
		return
	}

	syncJob.VideoCount = videoCount
	syncJob.UpdatedAt = time.Now()

	syncJobJSON, _ := json.Marshal(syncJob)
	s.redis.Set(ctx, syncJobKeyPrefix+syncJobID, syncJobJSON, 24*time.Hour)
}

// updateSyncJobProgress updates the progress counters of a sync job
func (s *Scheduler) updateSyncJobProgress(ctx context.Context, syncJobID string, downloaded, failed int) {
	syncJobData, err := s.redis.Get(ctx, syncJobKeyPrefix+syncJobID).Result()
	if err != nil {
		return
	}

	var syncJob SyncJob
	if err := json.Unmarshal([]byte(syncJobData), &syncJob); err != nil {
		return
	}

	syncJob.Downloaded = downloaded
	syncJob.Failed = failed
	syncJob.UpdatedAt = time.Now()

	syncJobJSON, _ := json.Marshal(syncJob)
	s.redis.Set(ctx, syncJobKeyPrefix+syncJobID, syncJobJSON, 24*time.Hour)
}

// updateChannelStatus updates the status of a channel
func (s *Scheduler) updateChannelStatus(ctx context.Context, channelID, status string) {
	channelData, err := s.redis.Get(ctx, channelKeyPrefix+channelID).Result()
	if err != nil {
		logging.Warn("error fetching channel",
			"channel_id", channelID,
			"error", err,
		)
		return
	}

	var channelMap map[string]interface{}
	if err := json.Unmarshal([]byte(channelData), &channelMap); err != nil {
		logging.Warn("error unmarshaling channel",
			"channel_id", channelID,
			"error", err,
		)
		return
	}

	channelMap["status"] = status
	channelMap["updated_at"] = time.Now()
	if status == "synced" {
		channelMap["last_sync_at"] = time.Now()
	}

	channelJSON, _ := json.Marshal(channelMap)
	s.redis.Set(ctx, channelKeyPrefix+channelID, channelJSON, 0)
}

// GetProgress returns the current progress for a channel sync.
// With KEDA-based scaling and unified queue, we count video statuses directly.
func (s *Scheduler) GetProgress(ctx context.Context, channelID string) (downloaded, failed, total int, err error) {
	// Count all video statuses for this channel
	downloaded, failed, pending := s.countAllVideoStatuses(ctx, channelID)

	// Total is all videos
	total = pending + downloaded + failed

	return downloaded, failed, total, nil
}

// countAllVideoStatuses counts videos by status for a channel
func (s *Scheduler) countAllVideoStatuses(ctx context.Context, channelID string) (downloaded, failed, pending int) {
	videoPattern := videoKeyPrefix + channelID + ":*"
	var cursor uint64

	for {
		keys, nextCursor, err := s.redis.Scan(ctx, cursor, videoPattern, 100).Result()
		if err != nil {
			logging.Warn("failed to scan video keys", "channel_id", channelID, "error", err)
			break
		}

		for _, key := range keys {
			videoData, err := s.redis.Get(ctx, key).Result()
			if err != nil {
				continue
			}

			var video map[string]interface{}
			if err := json.Unmarshal([]byte(videoData), &video); err != nil {
				continue
			}

			status, _ := video["status"].(string)
			switch status {
			case "completed":
				downloaded++
			case "failed", "error":
				failed++
			case "pending", "downloading":
				pending++
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return downloaded, failed, pending
}

// CleanupChannel cleans up any stale K8s jobs for a channel.
// This can be used for manual cleanup if needed.
func (s *Scheduler) CleanupChannel(ctx context.Context, channelID string) error {
	if s.k8sManager == nil {
		return nil
	}

	jobs, err := s.k8sManager.ListJobsByChannel(ctx, channelID)
	if err != nil {
		return fmt.Errorf("failed to list jobs for channel: %w", err)
	}

	for _, job := range jobs.Items {
		if err := s.k8sManager.DeleteJob(ctx, job.Name); err != nil {
			logging.Warn("failed to delete stale job",
				"channel_id", channelID,
				"job_name", job.Name,
				"error", err,
			)
		}
	}

	return nil
}

// getChannelName retrieves the channel name from Redis
func (s *Scheduler) getChannelName(ctx context.Context, channelID string) string {
	channelKey := channelKeyPrefix + channelID
	channelData, err := s.redis.Get(ctx, channelKey).Result()
	if err != nil {
		return ""
	}

	var channelMap map[string]interface{}
	if err := json.Unmarshal([]byte(channelData), &channelMap); err != nil {
		return ""
	}

	if name, ok := channelMap["name"].(string); ok {
		return name
	}
	return ""
}

// getExistingEpisodeCount counts existing videos for a channel to determine starting episode number
func (s *Scheduler) getExistingEpisodeCount(ctx context.Context, channelID string) int {
	videoPattern := videoKeyPrefix + channelID + ":*"
	var count int
	var cursor uint64

	for {
		keys, nextCursor, err := s.redis.Scan(ctx, cursor, videoPattern, 100).Result()
		if err != nil {
			break
		}
		count += len(keys)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return count
}
