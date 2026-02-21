package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"github.com/timholm/ytarchive/internal/db"
	"github.com/timholm/ytarchive/internal/scheduler"
	"github.com/timholm/ytarchive/internal/types"
	"github.com/timholm/ytarchive/internal/validation"
)

const (
	// Redis key prefixes
	channelKeyPrefix  = "channel:"
	channelListKey    = "channels"
	videoKeyPrefix    = "video:"
	jobKeyPrefix      = "job:"
	activeJobsKey     = "jobs:active"
	progressKeyPrefix = "progress:"
	cookiesKey        = "config:cookies"
)

// Channel represents a YouTube channel being tracked (API-specific extension of types.Channel)
type Channel struct {
	ID          string    `json:"id"`
	YouTubeURL  string    `json:"youtube_url"`
	YouTubeID   string    `json:"youtube_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	VideoCount  int       `json:"video_count"`
	Status      string    `json:"status"` // pending, syncing, synced, error
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastSyncAt  time.Time `json:"last_sync_at,omitempty"`
}

// Video represents a video from a channel (API-specific extension of types.Video)
type Video struct {
	ID           string    `json:"id"`
	YouTubeID    string    `json:"youtube_id"`
	ChannelID    string    `json:"channel_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	Duration     int       `json:"duration"` // in seconds
	UploadDate   string    `json:"upload_date,omitempty"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	ViewCount    int64     `json:"view_count,omitempty"`
	Status       string    `json:"status"` // pending, downloading, downloaded, error
	FilePath     string    `json:"file_path,omitempty"`
	FileSize     int64     `json:"file_size,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Job represents a download job (API-specific extension of types.Job)
type Job struct {
	ID          string    `json:"id"`
	ChannelID   string    `json:"channel_id"`
	WorkerNum   int       `json:"worker_num"`
	Status      string    `json:"status"` // pending, running, completed, failed
	VideoCount  int       `json:"video_count"`
	Downloaded  int       `json:"downloaded"`
	Failed      int       `json:"failed"`
	K8sJobName  string    `json:"k8s_job_name,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

// Ensure API types can be converted to canonical types
var (
	_ = convertChannelToCanonical
	_ = convertVideoToCanonical
	_ = convertJobToCanonical
)

// convertChannelToCanonical converts an API Channel to the canonical types.Channel
func convertChannelToCanonical(c *Channel) *types.Channel {
	return &types.Channel{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		URL:         c.YouTubeURL,
		VideoCount:  c.VideoCount,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

// convertVideoToCanonical converts an API Video to the canonical types.Video
func convertVideoToCanonical(v *Video) *types.Video {
	return &types.Video{
		ID:          v.ID,
		ChannelID:   v.ChannelID,
		Title:       v.Title,
		Description: v.Description,
		Duration:    v.Duration,
		Status:      v.Status,
		FilePath:    v.FilePath,
		FileSize:    v.FileSize,
		CreatedAt:   v.CreatedAt,
		UpdatedAt:   v.UpdatedAt,
	}
}

// convertJobToCanonical converts an API Job to the canonical types.Job
func convertJobToCanonical(j *Job) *types.Job {
	return &types.Job{
		ID:         j.ID,
		ChannelID:  j.ChannelID,
		Status:     j.Status,
		VideoCount: j.VideoCount,
		Completed:  j.Downloaded,
		Failed:     j.Failed,
		StartedAt:  j.CreatedAt,
		UpdatedAt:  j.UpdatedAt,
	}
}

// Progress represents overall download progress
type Progress struct {
	TotalChannels      int     `json:"total_channels"`
	TotalVideos        int     `json:"total_videos"`
	DownloadedVideos   int     `json:"downloaded_videos"`
	PendingVideos      int     `json:"pending_videos"`
	FailedVideos       int     `json:"failed_videos"`
	ActiveJobs         int     `json:"active_jobs"`
	DownloadPercentage float64 `json:"download_percentage"`
}

// Stats represents dashboard statistics
type Stats struct {
	TotalChannels    int    `json:"total_channels"`
	TotalVideos      int    `json:"total_videos"`
	StorageUsed      int64  `json:"storage_used"`
	StorageUsedHuman string `json:"storage_used_human"`
	ActiveDownloads  int    `json:"active_downloads"`
	DownloadedVideos int    `json:"downloaded_videos"`
	PendingVideos    int    `json:"pending_videos"`
	FailedVideos     int    `json:"failed_videos"`
}

// Activity represents an activity event
type Activity struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"` // channel_added, video_downloaded, job_started, job_completed, error
	Message   string    `json:"message"`
	ChannelID string    `json:"channel_id,omitempty"`
	VideoID   string    `json:"video_id,omitempty"`
	JobID     string    `json:"job_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// JobProgress represents aggregated progress for jobs
type JobProgress struct {
	TotalJobs       int     `json:"total_jobs"`
	RunningJobs     int     `json:"running_jobs"`
	CompletedJobs   int     `json:"completed_jobs"`
	FailedJobs      int     `json:"failed_jobs"`
	TotalVideos     int     `json:"total_videos"`
	Downloaded      int     `json:"downloaded"`
	Failed          int     `json:"failed"`
	ProgressPercent float64 `json:"progress_percent"`
}

// Redis key for activity feed
const activityKey = "activity:feed"

// AddChannelRequest is the request body for adding a channel
type AddChannelRequest struct {
	YouTubeURL string `json:"youtube_url" binding:"required"`
}

// Handlers contains all API handlers
type Handlers struct {
	redis     *redis.Client
	scheduler *scheduler.Scheduler
}

// NewHandlers creates a new Handlers instance
func NewHandlers(redisClient *redis.Client, sched *scheduler.Scheduler) *Handlers {
	return &Handlers{
		redis:     redisClient,
		scheduler: sched,
	}
}

// AddChannel handles POST /api/channels
func (h *Handlers) AddChannel(c *gin.Context) {
	var req AddChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Validate and extract YouTube channel ID from URL
	youtubeID, err := validation.ValidateChannelInput(req.YouTubeURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid YouTube URL: " + err.Error()})
		return
	}

	// Additional extraction for legacy URLs (keep backward compatibility)
	if extractedID, extractErr := extractChannelID(req.YouTubeURL); extractErr == nil && extractedID != "" {
		youtubeID = extractedID
	}

	// Check if channel already exists
	ctx := c.Request.Context()
	existingChannels, err := h.redis.SMembers(ctx, channelListKey).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error checking existing channels: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing channels"})
		return
	}

	for _, chID := range existingChannels {
		channelData, err := h.redis.Get(ctx, channelKeyPrefix+chID).Result()
		if err != nil {
			continue
		}
		var existingChannel Channel
		if err := json.Unmarshal([]byte(channelData), &existingChannel); err != nil {
			continue
		}
		if existingChannel.YouTubeID == youtubeID {
			c.JSON(http.StatusConflict, gin.H{"error": "Channel already exists", "channel": existingChannel})
			return
		}
	}

	// Create new channel
	now := time.Now()
	channel := Channel{
		ID:         uuid.New().String(),
		YouTubeURL: req.YouTubeURL,
		YouTubeID:  youtubeID,
		Status:     "pending",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Save to Redis
	channelJSON, err := json.Marshal(channel)
	if err != nil {
		log.Printf("Error marshaling channel: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create channel"})
		return
	}

	pipe := h.redis.Pipeline()
	pipe.Set(ctx, channelKeyPrefix+channel.ID, channelJSON, 0)
	pipe.SAdd(ctx, channelListKey, channel.ID)
	if _, err := pipe.Exec(ctx); err != nil {
		log.Printf("Error saving channel to Redis: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save channel"})
		return
	}

	log.Printf("Channel added: %s (YouTube ID: %s)", channel.ID, channel.YouTubeID)
	c.JSON(http.StatusCreated, channel)
}

// ListChannels handles GET /api/channels
func (h *Handlers) ListChannels(c *gin.Context) {
	ctx := c.Request.Context()

	channelIDs, err := h.redis.SMembers(ctx, channelListKey).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching channel list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch channels"})
		return
	}

	channels := make([]Channel, 0, len(channelIDs))
	for _, id := range channelIDs {
		channelData, err := h.redis.Get(ctx, channelKeyPrefix+id).Result()
		if err != nil {
			log.Printf("Error fetching channel %s: %v", id, err)
			continue
		}

		var channel Channel
		if err := json.Unmarshal([]byte(channelData), &channel); err != nil {
			log.Printf("Error unmarshaling channel %s: %v", id, err)
			continue
		}
		channels = append(channels, channel)
	}

	c.JSON(http.StatusOK, gin.H{"channels": channels, "count": len(channels)})
}

// GetChannel handles GET /api/channels/:id
func (h *Handlers) GetChannel(c *gin.Context) {
	channelID := c.Param("id")
	ctx := c.Request.Context()

	// Get channel
	channelData, err := h.redis.Get(ctx, channelKeyPrefix+channelID).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}
	if err != nil {
		log.Printf("Error fetching channel %s: %v", channelID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch channel"})
		return
	}

	var channel Channel
	if err := json.Unmarshal([]byte(channelData), &channel); err != nil {
		log.Printf("Error unmarshaling channel %s: %v", channelID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse channel data"})
		return
	}

	// Get videos for this channel
	videos, err := h.getChannelVideos(ctx, channelID)
	if err != nil {
		log.Printf("Error fetching videos for channel %s: %v", channelID, err)
	}

	c.JSON(http.StatusOK, gin.H{"channel": channel, "videos": videos})
}

// SyncChannel handles POST /api/channels/:id/sync
func (h *Handlers) SyncChannel(c *gin.Context) {
	channelID := c.Param("id")
	ctx := c.Request.Context()

	// Verify channel exists
	channelData, err := h.redis.Get(ctx, channelKeyPrefix+channelID).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}
	if err != nil {
		log.Printf("Error fetching channel %s: %v", channelID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch channel"})
		return
	}

	var channel Channel
	if err := json.Unmarshal([]byte(channelData), &channel); err != nil {
		log.Printf("Error unmarshaling channel %s: %v", channelID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse channel data"})
		return
	}

	// Check if already syncing
	if channel.Status == "syncing" {
		c.JSON(http.StatusConflict, gin.H{"error": "Channel is already syncing"})
		return
	}

	// Start sync using scheduler
	jobID, err := h.scheduler.StartSync(ctx, channel.ID, channel.YouTubeID)
	if err != nil {
		log.Printf("Error starting sync for channel %s: %v", channelID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start sync: " + err.Error()})
		return
	}

	// Update channel status
	channel.Status = "syncing"
	channel.UpdatedAt = time.Now()
	channelJSON, _ := json.Marshal(channel)
	if err := h.redis.Set(ctx, channelKeyPrefix+channelID, channelJSON, 0).Err(); err != nil {
		log.Printf("Error updating channel status: %v", err)
	}

	log.Printf("Sync started for channel %s, job ID: %s", channelID, jobID)
	c.JSON(http.StatusAccepted, gin.H{"message": "Sync started", "job_id": jobID, "channel": channel})
}

// DeleteChannel handles DELETE /api/channels/:id
func (h *Handlers) DeleteChannel(c *gin.Context) {
	channelID := c.Param("id")
	ctx := c.Request.Context()

	// Check if channel exists
	exists, err := h.redis.Exists(ctx, channelKeyPrefix+channelID).Result()
	if err != nil {
		log.Printf("Error checking channel existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check channel"})
		return
	}
	if exists == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	// Delete channel and remove from list
	pipe := h.redis.Pipeline()
	pipe.Del(ctx, channelKeyPrefix+channelID)
	pipe.SRem(ctx, channelListKey, channelID)

	// Also delete associated videos
	videoKeys, _ := h.redis.Keys(ctx, videoKeyPrefix+channelID+":*").Result()
	for _, key := range videoKeys {
		pipe.Del(ctx, key)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		log.Printf("Error deleting channel %s: %v", channelID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete channel"})
		return
	}

	log.Printf("Channel deleted: %s", channelID)
	c.JSON(http.StatusOK, gin.H{"message": "Channel deleted successfully"})
}

// ListJobs handles GET /api/jobs
func (h *Handlers) ListJobs(c *gin.Context) {
	ctx := c.Request.Context()

	jobIDs, err := h.redis.SMembers(ctx, activeJobsKey).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching job list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		return
	}

	jobs := make([]Job, 0, len(jobIDs))
	for _, id := range jobIDs {
		jobData, err := h.redis.Get(ctx, jobKeyPrefix+id).Result()
		if err != nil {
			log.Printf("Error fetching job %s: %v", id, err)
			continue
		}

		var job Job
		if err := json.Unmarshal([]byte(jobData), &job); err != nil {
			log.Printf("Error unmarshaling job %s: %v", id, err)
			continue
		}
		jobs = append(jobs, job)
	}

	c.JSON(http.StatusOK, gin.H{"jobs": jobs, "count": len(jobs)})
}

// GetProgress handles GET /api/progress
func (h *Handlers) GetProgress(c *gin.Context) {
	ctx := c.Request.Context()

	// Get all channels
	channelIDs, err := h.redis.SMembers(ctx, channelListKey).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching channel list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch progress"})
		return
	}

	// Get active jobs
	activeJobCount, err := h.redis.SCard(ctx, activeJobsKey).Result()
	if err != nil && err != redis.Nil {
		activeJobCount = 0
	}

	// Calculate video statistics
	var totalVideos, downloadedVideos, pendingVideos, failedVideos int
	for _, channelID := range channelIDs {
		videos, err := h.getChannelVideos(ctx, channelID)
		if err != nil {
			continue
		}
		for _, video := range videos {
			totalVideos++
			switch video.Status {
			case "downloaded":
				downloadedVideos++
			case "pending", "downloading":
				pendingVideos++
			case "error":
				failedVideos++
			}
		}
	}

	var downloadPercentage float64
	if totalVideos > 0 {
		downloadPercentage = float64(downloadedVideos) / float64(totalVideos) * 100
	}

	progress := Progress{
		TotalChannels:      len(channelIDs),
		TotalVideos:        totalVideos,
		DownloadedVideos:   downloadedVideos,
		PendingVideos:      pendingVideos,
		FailedVideos:       failedVideos,
		ActiveJobs:         int(activeJobCount),
		DownloadPercentage: downloadPercentage,
	}

	c.JSON(http.StatusOK, progress)
}

// GetStats handles GET /api/stats - Returns dashboard statistics
func (h *Handlers) GetStats(c *gin.Context) {
	ctx := c.Request.Context()

	// Get all channels
	channelIDs, err := h.redis.SMembers(ctx, channelListKey).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching channel list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}

	// Get active jobs count (these are active downloads)
	activeJobCount, err := h.redis.SCard(ctx, activeJobsKey).Result()
	if err != nil && err != redis.Nil {
		activeJobCount = 0
	}

	// Calculate video statistics and storage
	var totalVideos, downloadedVideos, pendingVideos, failedVideos int
	var storageUsed int64

	for _, channelID := range channelIDs {
		videos, err := h.getChannelVideos(ctx, channelID)
		if err != nil {
			continue
		}
		for _, video := range videos {
			totalVideos++
			storageUsed += video.FileSize
			switch video.Status {
			case "downloaded":
				downloadedVideos++
			case "pending", "downloading":
				pendingVideos++
			case "error":
				failedVideos++
			}
		}
	}

	stats := Stats{
		TotalChannels:    len(channelIDs),
		TotalVideos:      totalVideos,
		StorageUsed:      storageUsed,
		StorageUsedHuman: formatBytes(storageUsed),
		ActiveDownloads:  int(activeJobCount),
		DownloadedVideos: downloadedVideos,
		PendingVideos:    pendingVideos,
		FailedVideos:     failedVideos,
	}

	c.JSON(http.StatusOK, stats)
}

// GetActivity handles GET /api/activity - Returns recent activity feed
func (h *Handlers) GetActivity(c *gin.Context) {
	ctx := c.Request.Context()

	// Get the last 20 activity events from Redis list
	activities, err := h.redis.LRange(ctx, activityKey, 0, 19).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching activity feed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activity"})
		return
	}

	result := make([]Activity, 0, len(activities))
	for _, activityData := range activities {
		var activity Activity
		if err := json.Unmarshal([]byte(activityData), &activity); err != nil {
			log.Printf("Error unmarshaling activity: %v", err)
			continue
		}
		result = append(result, activity)
	}

	c.JSON(http.StatusOK, gin.H{"activities": result, "count": len(result)})
}

// SearchVideos handles GET /api/search - Full-text search across videos using SQLite FTS5
func (h *Handlers) SearchVideos(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query 'q' is required"})
		return
	}

	channelID := c.Query("channel")
	status := c.Query("status")
	limitStr := c.DefaultQuery("limit", "50")

	limit := 50
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 200 {
		limit = l
	}

	// If no channel specified, search across all channels
	ctx := c.Request.Context()
	channelIDs := []string{channelID}
	if channelID == "" {
		var err error
		channelIDs, err = h.redis.SMembers(ctx, channelListKey).Result()
		if err != nil && err != redis.Nil {
			log.Printf("Error fetching channel list: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch channels"})
			return
		}
	}

	var allResults []gin.H
	for _, chID := range channelIDs {
		if chID == "" {
			continue
		}

		// Open the channel's SQLite database
		channelDB, err := db.OpenChannelDB(chID)
		if err != nil {
			log.Printf("Error opening channel DB for %s: %v", chID, err)
			continue
		}

		var results []db.SearchResult
		if status != "" {
			results, err = db.SearchVideosWithFilter(channelDB, query, db.VideoStatus(status), limit)
		} else {
			results, err = db.SearchVideos(channelDB, query, limit)
		}
		channelDB.Close()

		if err != nil {
			log.Printf("Error searching channel %s: %v", chID, err)
			continue
		}

		for _, r := range results {
			allResults = append(allResults, gin.H{
				"id":          r.ID,
				"channel_id":  chID,
				"title":       r.Title,
				"description": r.Description,
				"duration":    r.Duration,
				"upload_date": r.UploadDate,
				"status":      r.Status,
				"file_path":   r.FilePath,
				"file_size":   r.FileSize,
				"rank":        r.Rank,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"query":   query,
		"results": allResults,
		"count":   len(allResults),
	})
}

// ListVideos handles GET /api/videos - Lists all videos with optional filters
func (h *Handlers) ListVideos(c *gin.Context) {
	ctx := c.Request.Context()

	// Get optional query parameters
	search := c.Query("search")
	status := c.Query("status")
	channelID := c.Query("channel")

	// Get all channels
	channelIDs, err := h.redis.SMembers(ctx, channelListKey).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching channel list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch videos"})
		return
	}

	// Filter by channel if specified
	if channelID != "" {
		channelIDs = []string{channelID}
	}

	var allVideos []Video
	for _, chID := range channelIDs {
		videos, err := h.getChannelVideos(ctx, chID)
		if err != nil {
			continue
		}
		allVideos = append(allVideos, videos...)
	}

	// Apply filters
	filteredVideos := make([]Video, 0)
	for _, video := range allVideos {
		// Filter by status
		if status != "" && video.Status != status {
			continue
		}
		// Filter by search term (title contains search string)
		if search != "" && !containsIgnoreCase(video.Title, search) {
			continue
		}
		filteredVideos = append(filteredVideos, video)
	}

	c.JSON(http.StatusOK, gin.H{"videos": filteredVideos, "count": len(filteredVideos)})
}

// GetVideo handles GET /api/videos/:id - Get single video details
func (h *Handlers) GetVideo(c *gin.Context) {
	videoID := c.Param("id")
	ctx := c.Request.Context()

	// Search for the video across all channels
	channelIDs, err := h.redis.SMembers(ctx, channelListKey).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching channel list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch video"})
		return
	}

	for _, channelID := range channelIDs {
		videoData, err := h.redis.Get(ctx, videoKeyPrefix+channelID+":"+videoID).Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			log.Printf("Error fetching video %s: %v", videoID, err)
			continue
		}

		var video Video
		if err := json.Unmarshal([]byte(videoData), &video); err != nil {
			log.Printf("Error unmarshaling video %s: %v", videoID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse video data"})
			return
		}

		c.JSON(http.StatusOK, video)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
}

// TriggerDownload handles POST /api/videos/:id/download - Trigger download for specific video
func (h *Handlers) TriggerDownload(c *gin.Context) {
	videoID := c.Param("id")
	ctx := c.Request.Context()

	// Search for the video across all channels
	channelIDs, err := h.redis.SMembers(ctx, channelListKey).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching channel list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to trigger download"})
		return
	}

	for _, channelID := range channelIDs {
		videoKey := videoKeyPrefix + channelID + ":" + videoID
		videoData, err := h.redis.Get(ctx, videoKey).Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			log.Printf("Error fetching video %s: %v", videoID, err)
			continue
		}

		var video Video
		if err := json.Unmarshal([]byte(videoData), &video); err != nil {
			log.Printf("Error unmarshaling video %s: %v", videoID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse video data"})
			return
		}

		// Check if already downloading or downloaded
		if video.Status == "downloading" {
			c.JSON(http.StatusConflict, gin.H{"error": "Video is already downloading"})
			return
		}
		if video.Status == "downloaded" {
			c.JSON(http.StatusConflict, gin.H{"error": "Video is already downloaded"})
			return
		}

		// Update video status to pending for download
		video.Status = "pending"
		video.UpdatedAt = time.Now()
		videoJSON, _ := json.Marshal(video)
		if err := h.redis.Set(ctx, videoKey, videoJSON, 0).Err(); err != nil {
			log.Printf("Error updating video status: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue video for download"})
			return
		}

		// Add to download queue (unified queue for all channels)
		// Format: "channelID:videoID"
		queueKey := "ytarchive:download:queue"
		queueItem := channelID + ":" + videoID
		if err := h.redis.RPush(ctx, queueKey, queueItem).Err(); err != nil {
			log.Printf("Error adding video to download queue: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue video for download"})
			return
		}

		log.Printf("Download triggered for video %s", videoID)
		c.JSON(http.StatusAccepted, gin.H{"message": "Download queued", "video": video})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
}

// GetChannelVideos handles GET /api/channels/:id/videos - Get videos for a specific channel
func (h *Handlers) GetChannelVideos(c *gin.Context) {
	channelID := c.Param("id")
	ctx := c.Request.Context()

	// Verify channel exists
	exists, err := h.redis.Exists(ctx, channelKeyPrefix+channelID).Result()
	if err != nil {
		log.Printf("Error checking channel existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch channel videos"})
		return
	}
	if exists == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	videos, err := h.getChannelVideos(ctx, channelID)
	if err != nil {
		log.Printf("Error fetching videos for channel %s: %v", channelID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch videos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"videos": videos, "count": len(videos), "channel_id": channelID})
}

// GetJob handles GET /api/jobs/:id - Get specific job details
func (h *Handlers) GetJob(c *gin.Context) {
	jobID := c.Param("id")
	ctx := c.Request.Context()

	jobData, err := h.redis.Get(ctx, jobKeyPrefix+jobID).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	if err != nil {
		log.Printf("Error fetching job %s: %v", jobID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch job"})
		return
	}

	var job Job
	if err := json.Unmarshal([]byte(jobData), &job); err != nil {
		log.Printf("Error unmarshaling job %s: %v", jobID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse job data"})
		return
	}

	c.JSON(http.StatusOK, job)
}

// CancelJob handles POST /api/jobs/:id/cancel - Cancel a running job
func (h *Handlers) CancelJob(c *gin.Context) {
	jobID := c.Param("id")
	ctx := c.Request.Context()

	jobData, err := h.redis.Get(ctx, jobKeyPrefix+jobID).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	if err != nil {
		log.Printf("Error fetching job %s: %v", jobID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch job"})
		return
	}

	var job Job
	if err := json.Unmarshal([]byte(jobData), &job); err != nil {
		log.Printf("Error unmarshaling job %s: %v", jobID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse job data"})
		return
	}

	// Check if job can be cancelled
	if job.Status == "completed" {
		c.JSON(http.StatusConflict, gin.H{"error": "Job is already completed"})
		return
	}
	if job.Status == "failed" {
		c.JSON(http.StatusConflict, gin.H{"error": "Job has already failed"})
		return
	}
	if job.Status == "cancelled" {
		c.JSON(http.StatusConflict, gin.H{"error": "Job is already cancelled"})
		return
	}

	// Update job status to cancelled
	job.Status = "cancelled"
	job.UpdatedAt = time.Now()
	job.CompletedAt = time.Now()
	jobJSON, _ := json.Marshal(job)

	pipe := h.redis.Pipeline()
	pipe.Set(ctx, jobKeyPrefix+jobID, jobJSON, 0)
	pipe.SRem(ctx, activeJobsKey, jobID)
	if _, err := pipe.Exec(ctx); err != nil {
		log.Printf("Error cancelling job %s: %v", jobID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel job"})
		return
	}

	log.Printf("Job cancelled: %s", jobID)
	c.JSON(http.StatusOK, gin.H{"message": "Job cancelled successfully", "job": job})
}

// GetJobsProgress handles GET /api/jobs/progress - Get aggregated progress for all jobs
func (h *Handlers) GetJobsProgress(c *gin.Context) {
	ctx := c.Request.Context()

	jobIDs, err := h.redis.SMembers(ctx, activeJobsKey).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching job list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs progress"})
		return
	}

	progress := JobProgress{}
	for _, id := range jobIDs {
		jobData, err := h.redis.Get(ctx, jobKeyPrefix+id).Result()
		if err != nil {
			continue
		}

		var job Job
		if err := json.Unmarshal([]byte(jobData), &job); err != nil {
			continue
		}

		progress.TotalJobs++
		progress.TotalVideos += job.VideoCount
		progress.Downloaded += job.Downloaded
		progress.Failed += job.Failed

		switch job.Status {
		case "running", "pending":
			progress.RunningJobs++
		case "completed":
			progress.CompletedJobs++
		case "failed":
			progress.FailedJobs++
		}
	}

	if progress.TotalVideos > 0 {
		progress.ProgressPercent = float64(progress.Downloaded) / float64(progress.TotalVideos) * 100
	}

	c.JSON(http.StatusOK, progress)
}

// IndexChannelVideos handles POST /api/channels/:id/index - Rebuild FTS index for a channel
func (h *Handlers) IndexChannelVideos(c *gin.Context) {
	channelID := c.Param("id")
	ctx := c.Request.Context()

	// Get all videos from Redis
	videos, err := h.getChannelVideos(ctx, channelID)
	if err != nil {
		log.Printf("Error fetching videos for channel %s: %v", channelID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch videos"})
		return
	}

	// Open/create SQLite database for this channel
	channelDB, err := db.OpenChannelDB(channelID)
	if err != nil {
		log.Printf("Error opening channel DB for %s: %v", channelID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open channel database"})
		return
	}
	defer channelDB.Close()

	// Insert each video into SQLite
	indexed := 0
	for _, video := range videos {
		dbVideo := &db.Video{
			ID:           video.YouTubeID,
			Title:        video.Title,
			Description:  video.Description,
			Duration:     int64(video.Duration),
			UploadDate:   video.UploadDate,
			ThumbnailURL: video.ThumbnailURL,
			ViewCount:    video.ViewCount,
			Status:       db.VideoStatus(video.Status),
			FilePath:     video.FilePath,
			FileSize:     video.FileSize,
		}
		if err := db.InsertVideo(channelDB, dbVideo); err != nil {
			log.Printf("Error inserting video %s: %v", video.YouTubeID, err)
			continue
		}
		indexed++
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Index rebuilt successfully",
		"channel_id":    channelID,
		"videos_found":  len(videos),
		"videos_indexed": indexed,
	})
}

// IndexAllChannels handles POST /api/index - Rebuild FTS index for all channels
func (h *Handlers) IndexAllChannels(c *gin.Context) {
	ctx := c.Request.Context()

	// Get all channels
	channelIDs, err := h.redis.SMembers(ctx, channelListKey).Result()
	if err != nil {
		log.Printf("Error fetching channel list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch channels"})
		return
	}

	results := make([]gin.H, 0, len(channelIDs))
	totalIndexed := 0

	for _, channelID := range channelIDs {
		// Get videos from Redis
		videos, err := h.getChannelVideos(ctx, channelID)
		if err != nil {
			log.Printf("Error fetching videos for channel %s: %v", channelID, err)
			results = append(results, gin.H{"channel_id": channelID, "error": err.Error()})
			continue
		}

		// Open/create SQLite database
		channelDB, err := db.OpenChannelDB(channelID)
		if err != nil {
			log.Printf("Error opening channel DB for %s: %v", channelID, err)
			results = append(results, gin.H{"channel_id": channelID, "error": err.Error()})
			continue
		}

		// Insert videos
		indexed := 0
		for _, video := range videos {
			dbVideo := &db.Video{
				ID:           video.YouTubeID,
				Title:        video.Title,
				Description:  video.Description,
				Duration:     int64(video.Duration),
				UploadDate:   video.UploadDate,
				ThumbnailURL: video.ThumbnailURL,
				ViewCount:    video.ViewCount,
				Status:       db.VideoStatus(video.Status),
				FilePath:     video.FilePath,
				FileSize:     video.FileSize,
			}
			if err := db.InsertVideo(channelDB, dbVideo); err == nil {
				indexed++
			}
		}
		channelDB.Close()

		totalIndexed += indexed
		results = append(results, gin.H{
			"channel_id":     channelID,
			"videos_found":   len(videos),
			"videos_indexed": indexed,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Index rebuilt for all channels",
		"channels":       len(channelIDs),
		"total_indexed":  totalIndexed,
		"results":        results,
	})
}

// getChannelVideos retrieves all videos for a channel
func (h *Handlers) getChannelVideos(ctx context.Context, channelID string) ([]Video, error) {
	videoKeys, err := h.redis.Keys(ctx, videoKeyPrefix+channelID+":*").Result()
	if err != nil {
		return nil, err
	}

	videos := make([]Video, 0, len(videoKeys))
	for _, key := range videoKeys {
		videoData, err := h.redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var video Video
		if err := json.Unmarshal([]byte(videoData), &video); err != nil {
			continue
		}
		videos = append(videos, video)
	}

	return videos, nil
}

// extractChannelID extracts the YouTube channel ID from various URL formats
func extractChannelID(url string) (string, error) {
	patterns := []struct {
		regex *regexp.Regexp
		group int
	}{
		// https://www.youtube.com/channel/UC...
		{regexp.MustCompile(`youtube\.com/channel/([a-zA-Z0-9_-]+)`), 1},
		// https://www.youtube.com/@username
		{regexp.MustCompile(`youtube\.com/@([a-zA-Z0-9_-]+)`), 1},
		// https://www.youtube.com/c/channelname
		{regexp.MustCompile(`youtube\.com/c/([a-zA-Z0-9_-]+)`), 1},
		// https://www.youtube.com/user/username
		{regexp.MustCompile(`youtube\.com/user/([a-zA-Z0-9_-]+)`), 1},
	}

	for _, p := range patterns {
		matches := p.regex.FindStringSubmatch(url)
		if len(matches) > p.group {
			return matches[p.group], nil
		}
	}

	return "", &InvalidURLError{URL: url}
}

// InvalidURLError represents an invalid YouTube URL error
type InvalidURLError struct {
	URL string
}

func (e *InvalidURLError) Error() string {
	return "could not extract channel ID from URL: " + e.URL
}

// formatBytes converts bytes to human-readable format
func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case bytes >= TB:
		return formatFloat64(float64(bytes)/TB) + " TB"
	case bytes >= GB:
		return formatFloat64(float64(bytes)/GB) + " GB"
	case bytes >= MB:
		return formatFloat64(float64(bytes)/MB) + " MB"
	case bytes >= KB:
		return formatFloat64(float64(bytes)/KB) + " KB"
	default:
		return formatInt64(bytes) + " B"
	}
}

// formatFloat64 formats a float64 with 2 decimal places
func formatFloat64(f float64) string {
	// Simple formatting without fmt to avoid import
	intPart := int64(f)
	decPart := int64((f - float64(intPart)) * 100)
	if decPart < 0 {
		decPart = -decPart
	}
	decStr := formatInt64(decPart)
	if len(decStr) < 2 {
		decStr = "0" + decStr
	}
	return formatInt64(intPart) + "." + decStr
}

// formatInt64 formats an int64 as string
func formatInt64(i int64) string {
	if i == 0 {
		return "0"
	}
	s := ""
	negative := i < 0
	if negative {
		i = -i
	}
	for i > 0 {
		s = string(rune('0'+i%10)) + s
		i /= 10
	}
	if negative {
		s = "-" + s
	}
	return s
}

// containsIgnoreCase checks if s contains substr (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	return regexp.MustCompile("(?i)" + regexp.QuoteMeta(substr)).MatchString(s)
}

// UpdateProgressRequest is the request body for updating video progress
type UpdateProgressRequest struct {
	VideoID    string  `json:"video_id"`
	WorkerID   string  `json:"worker_id"`
	Percent    float64 `json:"percent"`
	Downloaded int64   `json:"downloaded"`
	Total      int64   `json:"total"`
	Speed      string  `json:"speed"`
	ETA        string  `json:"eta"`
}

// DownloadProgressInfo represents detailed progress for a single video download
type DownloadProgressInfo struct {
	VideoID         string  `json:"video_id"`
	WorkerID        string  `json:"worker_id"`
	Status          string  `json:"status"`
	Percentage      float64 `json:"percentage"`
	DownloadedBytes int64   `json:"downloaded_bytes"`
	TotalBytes      int64   `json:"total_bytes"`
	Speed           string  `json:"speed"`
	ETA             string  `json:"eta"`
	Fragment        string  `json:"fragment,omitempty"`
	UpdatedAt       int64   `json:"updated_at"`
}

// UpdateVideoStatusRequest is the request body for updating video status
type UpdateVideoStatusRequest struct {
	VideoID   string `json:"video_id"`
	ChannelID string `json:"channel_id"`
	Status    string `json:"status"`
	FilePath  string `json:"file_path,omitempty"`
	FileSize  int64  `json:"file_size,omitempty"`
	Error     string `json:"error,omitempty"`
}

// GetDownloadsProgress handles GET /api/downloads/progress - Get all active download progress
func (h *Handlers) GetDownloadsProgress(c *gin.Context) {
	ctx := c.Request.Context()

	// Find all progress keys in Redis
	progressKeys, err := h.redis.Keys(ctx, progressKeyPrefix+"*").Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching progress keys: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch download progress"})
		return
	}

	downloads := make([]DownloadProgressInfo, 0, len(progressKeys))
	now := time.Now().Unix()

	for _, key := range progressKeys {
		progressData, err := h.redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var progress map[string]interface{}
		if err := json.Unmarshal([]byte(progressData), &progress); err != nil {
			continue
		}

		// Check if this progress entry is recent (within last 60 seconds)
		updatedAt, _ := progress["updated_at"].(float64)
		if now-int64(updatedAt) > 60 {
			// Stale progress entry, skip it
			continue
		}

		info := DownloadProgressInfo{
			VideoID:   getStringValue(progress, "video_id"),
			WorkerID:  getStringValue(progress, "worker_id"),
			Status:    getStringValue(progress, "status"),
			Speed:     getStringValue(progress, "speed"),
			ETA:       getStringValue(progress, "eta"),
			Fragment:  getStringValue(progress, "fragment"),
			UpdatedAt: int64(updatedAt),
		}

		// Handle percentage - can be "percent" or "percentage"
		if pct, ok := progress["percent"].(float64); ok {
			info.Percentage = pct
		} else if pct, ok := progress["percentage"].(float64); ok {
			info.Percentage = pct
		}

		// Handle downloaded bytes
		if dl, ok := progress["downloaded"].(float64); ok {
			info.DownloadedBytes = int64(dl)
		} else if dl, ok := progress["downloaded_bytes"].(float64); ok {
			info.DownloadedBytes = int64(dl)
		}

		// Handle total bytes
		if total, ok := progress["total"].(float64); ok {
			info.TotalBytes = int64(total)
		} else if total, ok := progress["total_bytes"].(float64); ok {
			info.TotalBytes = int64(total)
		}

		// Default status if not set
		if info.Status == "" {
			info.Status = "downloading"
		}

		downloads = append(downloads, info)
	}

	c.JSON(http.StatusOK, gin.H{
		"downloads": downloads,
		"count":     len(downloads),
	})
}

// getStringValue safely extracts a string value from a map
func getStringValue(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

// UpdateProgress handles POST /api/progress/:id - Update download progress for a video
func (h *Handlers) UpdateProgress(c *gin.Context) {
	videoID := c.Param("id")
	ctx := c.Request.Context()

	var req UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Use videoID from path if not in body
	if req.VideoID == "" {
		req.VideoID = videoID
	}

	// Store progress in Redis with key progress:{videoID}
	progressKey := progressKeyPrefix + videoID
	progressData := map[string]interface{}{
		"video_id":         req.VideoID,
		"worker_id":        req.WorkerID,
		"status":           "downloading",
		"percentage":       req.Percent,
		"downloaded_bytes": req.Downloaded,
		"total_bytes":      req.Total,
		"speed":            req.Speed,
		"eta":              req.ETA,
		"updated_at":       time.Now().Unix(),
	}

	progressJSON, err := json.Marshal(progressData)
	if err != nil {
		log.Printf("Error marshaling progress: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save progress"})
		return
	}

	// Set with 1 hour expiration (progress data is ephemeral)
	if err := h.redis.Set(ctx, progressKey, progressJSON, time.Hour).Err(); err != nil {
		log.Printf("Error saving progress to Redis: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save progress"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Progress updated", "video_id": videoID})
}

// UpdateVideoStatus handles POST /api/videos/:id/status - Update video status
func (h *Handlers) UpdateVideoStatus(c *gin.Context) {
	videoID := c.Param("id")
	ctx := c.Request.Context()

	var req UpdateVideoStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Use videoID from path if not in body
	if req.VideoID == "" {
		req.VideoID = videoID
	}

	// If channelID is provided, update the specific video
	if req.ChannelID != "" {
		videoKey := videoKeyPrefix + req.ChannelID + ":" + videoID
		if err := h.updateVideoInRedis(ctx, videoKey, req); err != nil {
			log.Printf("Error updating video %s: %v", videoID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update video status"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Video status updated", "video_id": videoID})
		return
	}

	// If no channelID, search for video across all channels
	channelIDs, err := h.redis.SMembers(ctx, channelListKey).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching channel list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find video"})
		return
	}

	for _, channelID := range channelIDs {
		videoKey := videoKeyPrefix + channelID + ":" + videoID
		exists, err := h.redis.Exists(ctx, videoKey).Result()
		if err != nil {
			continue
		}
		if exists > 0 {
			if err := h.updateVideoInRedis(ctx, videoKey, req); err != nil {
				log.Printf("Error updating video %s: %v", videoID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update video status"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Video status updated", "video_id": videoID})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
}

// updateVideoInRedis updates a video's status in Redis
func (h *Handlers) updateVideoInRedis(ctx context.Context, videoKey string, req UpdateVideoStatusRequest) error {
	videoData, err := h.redis.Get(ctx, videoKey).Result()
	if err != nil {
		return err
	}

	var video map[string]interface{}
	if err := json.Unmarshal([]byte(videoData), &video); err != nil {
		return err
	}

	// Update fields
	video["status"] = req.Status
	video["updated_at"] = time.Now()

	if req.FilePath != "" {
		video["file_path"] = req.FilePath
	}
	if req.FileSize > 0 {
		video["file_size"] = req.FileSize
	}
	if req.Error != "" {
		video["error"] = req.Error
	}

	videoJSON, err := json.Marshal(video)
	if err != nil {
		return err
	}

	return h.redis.Set(ctx, videoKey, videoJSON, 0).Err()
}

// CookiesRequest is the request body for saving cookies
type CookiesRequest struct {
	Cookies string `json:"cookies" binding:"required"`
}

// CookiesResponse is the response for cookie-related endpoints
type CookiesResponse struct {
	Configured bool   `json:"configured"`
	Message    string `json:"message,omitempty"`
}

// GetCookies handles GET /api/cookies - Check if cookies are configured
func (h *Handlers) GetCookies(c *gin.Context) {
	ctx := c.Request.Context()
	exists, err := h.redis.Exists(ctx, cookiesKey).Result()
	if err != nil {
		log.Printf("Error checking cookies: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check cookies"})
		return
	}
	c.JSON(http.StatusOK, CookiesResponse{Configured: exists > 0})
}

// SaveCookies handles POST /api/cookies - Save cookies configuration
func (h *Handlers) SaveCookies(c *gin.Context) {
	var req CookiesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}
	if !isValidNetscapeCookies(req.Cookies) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cookies format. Please paste cookies in Netscape/cookies.txt format."})
		return
	}
	ctx := c.Request.Context()
	if err := h.redis.Set(ctx, cookiesKey, req.Cookies, 0).Err(); err != nil {
		log.Printf("Error saving cookies: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cookies"})
		return
	}
	log.Printf("Cookies saved successfully")
	c.JSON(http.StatusOK, CookiesResponse{Configured: true, Message: "Cookies saved successfully"})
}

// DeleteCookies handles DELETE /api/cookies - Delete cookies configuration
func (h *Handlers) DeleteCookies(c *gin.Context) {
	ctx := c.Request.Context()
	if err := h.redis.Del(ctx, cookiesKey).Err(); err != nil {
		log.Printf("Error deleting cookies: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cookies"})
		return
	}
	log.Printf("Cookies deleted")
	c.JSON(http.StatusOK, CookiesResponse{Configured: false, Message: "Cookies deleted successfully"})
}

// GetVideoStream handles GET /api/videos/:id/stream - Stream video file
func (h *Handlers) GetVideoStream(c *gin.Context) {
	videoID := c.Param("id")
	ctx := c.Request.Context()

	// Find video to get channel ID
	channelIDs, err := h.redis.SMembers(ctx, channelListKey).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching channel list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch video"})
		return
	}

	var channelID string
	for _, cID := range channelIDs {
		videoData, err := h.redis.Get(ctx, videoKeyPrefix+cID+":"+videoID).Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			continue
		}

		var video Video
		if err := json.Unmarshal([]byte(videoData), &video); err != nil {
			continue
		}

		channelID = cID
		break
	}

	if channelID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Proxy to collector service
	collectorURL := getCollectorURL()
	streamURL := fmt.Sprintf("%s/stream/%s/%s", collectorURL, channelID, videoID)

	// Create request to collector
	req, err := http.NewRequestWithContext(ctx, "GET", streamURL, nil)
	if err != nil {
		log.Printf("Error creating collector request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stream video"})
		return
	}

	// Forward range headers for seeking
	if rangeHeader := c.GetHeader("Range"); rangeHeader != "" {
		req.Header.Set("Range", rangeHeader)
	}

	// Make request to collector
	client := &http.Client{Timeout: 0} // No timeout for streaming
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error streaming from collector: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to stream video from storage"})
		return
	}
	defer resp.Body.Close()

	// Forward response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Set status and stream body
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

// GetVideoThumbnail handles GET /api/videos/:id/thumbnail - Serve thumbnail
func (h *Handlers) GetVideoThumbnail(c *gin.Context) {
	videoID := c.Param("id")
	ctx := c.Request.Context()

	// Find video to get channel ID
	channelIDs, err := h.redis.SMembers(ctx, channelListKey).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching channel list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch thumbnail"})
		return
	}

	for _, channelID := range channelIDs {
		exists, err := h.redis.Exists(ctx, videoKeyPrefix+channelID+":"+videoID).Result()
		if err != nil || exists == 0 {
			continue
		}

		storagePath := getStoragePath()
		videoDir := storagePath + "/channels/" + channelID + "/videos/" + videoID

		// Try different thumbnail formats
		patterns := []string{"thumbnail.jpg", "thumbnail.webp", "thumbnail.png"}
		for _, pattern := range patterns {
			path := videoDir + "/" + pattern
			if _, err := os.Stat(path); err == nil {
				c.File(path)
				return
			}
		}

		// If no local thumbnail, try to get from YouTube
		c.Redirect(http.StatusTemporaryRedirect, "https://img.youtube.com/vi/"+videoID+"/maxresdefault.jpg")
		return
	}

	// Fallback to YouTube thumbnail
	c.Redirect(http.StatusTemporaryRedirect, "https://img.youtube.com/vi/"+videoID+"/maxresdefault.jpg")
}

// GetVideoMetadata handles GET /api/videos/:id/metadata - Get video metadata file
func (h *Handlers) GetVideoMetadata(c *gin.Context) {
	videoID := c.Param("id")
	ctx := c.Request.Context()

	channelIDs, err := h.redis.SMembers(ctx, channelListKey).Result()
	if err != nil && err != redis.Nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch metadata"})
		return
	}

	for _, channelID := range channelIDs {
		exists, err := h.redis.Exists(ctx, videoKeyPrefix+channelID+":"+videoID).Result()
		if err != nil || exists == 0 {
			continue
		}

		storagePath := getStoragePath()
		metadataPath := storagePath + "/channels/" + channelID + "/videos/" + videoID + "/metadata.json"

		if _, err := os.Stat(metadataPath); err == nil {
			c.File(metadataPath)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Metadata file not found"})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
}

// GetVideoAudio handles GET /api/videos/:id/audio - Serve audio file
func (h *Handlers) GetVideoAudio(c *gin.Context) {
	videoID := c.Param("id")
	ctx := c.Request.Context()

	channelIDs, err := h.redis.SMembers(ctx, channelListKey).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching channel list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audio"})
		return
	}

	for _, channelID := range channelIDs {
		exists, err := h.redis.Exists(ctx, videoKeyPrefix+channelID+":"+videoID).Result()
		if err != nil || exists == 0 {
			continue
		}

		storagePath := getStoragePath()
		videoDir := storagePath + "/channels/" + channelID + "/videos/" + videoID

		// Look for audio file with various extensions
		patterns := []string{"audio.m4a", "audio.mp3", "audio.webm", "audio.opus"}
		var audioPath string
		for _, pattern := range patterns {
			path := videoDir + "/" + pattern
			if _, err := os.Stat(path); err == nil {
				audioPath = path
				break
			}
		}

		if audioPath == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Audio file not found"})
			return
		}

		// Serve the file with proper content type
		c.Header("Accept-Ranges", "bytes")
		c.File(audioPath)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
}

// GetVideoSubtitles handles GET /api/videos/:id/subtitles - Serve subtitles file
func (h *Handlers) GetVideoSubtitles(c *gin.Context) {
	videoID := c.Param("id")
	lang := c.DefaultQuery("lang", "en")
	ctx := c.Request.Context()

	channelIDs, err := h.redis.SMembers(ctx, channelListKey).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching channel list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subtitles"})
		return
	}

	for _, channelID := range channelIDs {
		exists, err := h.redis.Exists(ctx, videoKeyPrefix+channelID+":"+videoID).Result()
		if err != nil || exists == 0 {
			continue
		}

		storagePath := getStoragePath()
		videoDir := storagePath + "/channels/" + channelID + "/videos/" + videoID

		// Look for subtitle file - try specific language first, then default
		var subtitlePath string
		patterns := []string{
			"subtitles." + lang + ".vtt",
			"subtitles.vtt",
			"subtitles." + lang + ".srt",
			"subtitles.srt",
		}

		for _, pattern := range patterns {
			path := videoDir + "/" + pattern
			if _, err := os.Stat(path); err == nil {
				subtitlePath = path
				break
			}
		}

		if subtitlePath == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subtitle file not found"})
			return
		}

		// Set content type for WebVTT
		c.Header("Content-Type", "text/vtt; charset=utf-8")
		c.File(subtitlePath)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
}

// ListVideoSubtitles handles GET /api/videos/:id/subtitles/list - List available subtitle languages
func (h *Handlers) ListVideoSubtitles(c *gin.Context) {
	videoID := c.Param("id")
	ctx := c.Request.Context()

	channelIDs, err := h.redis.SMembers(ctx, channelListKey).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching channel list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list subtitles"})
		return
	}

	for _, channelID := range channelIDs {
		exists, err := h.redis.Exists(ctx, videoKeyPrefix+channelID+":"+videoID).Result()
		if err != nil || exists == 0 {
			continue
		}

		storagePath := getStoragePath()
		videoDir := storagePath + "/channels/" + channelID + "/videos/" + videoID

		// Find all subtitle files
		subtitles := []map[string]string{}

		// Check for default subtitle file
		if _, err := os.Stat(videoDir + "/subtitles.vtt"); err == nil {
			subtitles = append(subtitles, map[string]string{
				"lang": "en",
				"file": "subtitles.vtt",
				"url":  "/api/videos/" + videoID + "/subtitles",
			})
		}

		// Look for language-specific files
		entries, err := os.ReadDir(videoDir)
		if err == nil {
			for _, entry := range entries {
				name := entry.Name()
				// Match pattern: subtitles.{lang}.vtt
				if len(name) > 14 && name[:10] == "subtitles." && name[len(name)-4:] == ".vtt" {
					lang := name[10 : len(name)-4]
					if lang != "" {
						subtitles = append(subtitles, map[string]string{
							"lang": lang,
							"file": name,
							"url":  "/api/videos/" + videoID + "/subtitles?lang=" + lang,
						})
					}
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"video_id":  videoID,
			"subtitles": subtitles,
			"count":     len(subtitles),
		})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
}

// VideoFilesResponse contains information about available files for a video
type VideoFilesResponse struct {
	VideoID  string `json:"video_id"`
	HasVideo bool   `json:"has_video"`
	HasAudio bool   `json:"has_audio"`
	HasSubs  bool   `json:"has_subtitles"`
	HasThumb bool   `json:"has_thumbnail"`
	HasMeta  bool   `json:"has_metadata"`
	VideoURL string `json:"video_url,omitempty"`
	AudioURL string `json:"audio_url,omitempty"`
	SubsURL  string `json:"subtitles_url,omitempty"`
	ThumbURL string `json:"thumbnail_url,omitempty"`
	MetaURL  string `json:"metadata_url,omitempty"`
}

// GetVideoFiles handles GET /api/videos/:id/files - List available files for a video
func (h *Handlers) GetVideoFiles(c *gin.Context) {
	videoID := c.Param("id")
	ctx := c.Request.Context()

	channelIDs, err := h.redis.SMembers(ctx, channelListKey).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching channel list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch video files"})
		return
	}

	for _, channelID := range channelIDs {
		exists, err := h.redis.Exists(ctx, videoKeyPrefix+channelID+":"+videoID).Result()
		if err != nil || exists == 0 {
			continue
		}

		storagePath := getStoragePath()
		videoDir := storagePath + "/channels/" + channelID + "/videos/" + videoID

		response := VideoFilesResponse{
			VideoID: videoID,
		}

		// Check for video file
		for _, ext := range []string{"mp4", "mkv", "webm"} {
			if _, err := os.Stat(videoDir + "/video." + ext); err == nil {
				response.HasVideo = true
				response.VideoURL = "/api/videos/" + videoID + "/stream"
				break
			}
		}

		// Check for audio file
		for _, ext := range []string{"m4a", "mp3", "webm", "opus"} {
			if _, err := os.Stat(videoDir + "/audio." + ext); err == nil {
				response.HasAudio = true
				response.AudioURL = "/api/videos/" + videoID + "/audio"
				break
			}
		}

		// Check for subtitles
		if _, err := os.Stat(videoDir + "/subtitles.vtt"); err == nil {
			response.HasSubs = true
			response.SubsURL = "/api/videos/" + videoID + "/subtitles"
		}

		// Check for thumbnail
		for _, ext := range []string{"jpg", "webp", "png"} {
			if _, err := os.Stat(videoDir + "/thumbnail." + ext); err == nil {
				response.HasThumb = true
				response.ThumbURL = "/api/videos/" + videoID + "/thumbnail"
				break
			}
		}

		// Check for metadata
		if _, err := os.Stat(videoDir + "/metadata.json"); err == nil {
			response.HasMeta = true
			response.MetaURL = "/api/videos/" + videoID + "/metadata"
		}

		c.JSON(http.StatusOK, response)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
}

// getStoragePath returns the configured storage path
func getStoragePath() string {
	path := os.Getenv("STORAGE_PATH")
	if path == "" {
		path = "/data"
	}
	return path
}

func getCollectorURL() string {
	url := os.Getenv("COLLECTOR_URL")
	if url == "" {
		url = "http://collector.ytarchive.svc.cluster.local:8081"
	}
	return url
}

// isValidNetscapeCookies checks if the string looks like Netscape format cookies
func isValidNetscapeCookies(cookies string) bool {
	if cookies == "" {
		return false
	}
	lines := regexp.MustCompile(`\r?\n`).Split(cookies, -1)
	validLineCount := 0
	for _, line := range lines {
		trimmed := regexp.MustCompile(`^\s+`).ReplaceAllString(line, "")
		if trimmed == "" || (len(trimmed) > 0 && trimmed[0] == '#') {
			continue
		}
		fields := regexp.MustCompile(`\t`).Split(line, -1)
		if len(fields) >= 7 {
			validLineCount++
		}
	}
	return validLineCount > 0
}

// ServeFrontend serves the main HTML frontend
func (h *Handlers) ServeFrontend(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, frontendHTML)
}

const frontendHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>YT Archive</title>
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #1a1a2e; color: #eee; min-height: 100vh; }
        .container { max-width: 1400px; margin: 0 auto; padding: 2rem; }
        h1 { font-size: 2rem; margin-bottom: 1rem; color: #fff; }
        h2 { font-size: 1.5rem; margin-bottom: 1rem; color: #ccc; }
        .card { background: #16213e; border-radius: 12px; padding: 2rem; margin-bottom: 1.5rem; }
        .setup-card { max-width: 650px; margin: 2rem auto; }
        textarea { width: 100%; height: 180px; padding: 1rem; border: 1px solid #333; border-radius: 8px; background: #0f0f23; color: #eee; font-family: monospace; font-size: 0.85rem; }
        input[type="text"] { width: 100%; padding: 0.75rem 1rem; border: 1px solid #333; border-radius: 8px; background: #0f0f23; color: #eee; }
        .btn { display: inline-block; padding: 0.75rem 1.5rem; border: none; border-radius: 8px; cursor: pointer; font-size: 1rem; font-weight: 500; transition: opacity 0.2s, transform 0.1s; }
        .btn:hover { opacity: 0.9; }
        .btn:active { transform: scale(0.98); }
        .btn:disabled { opacity: 0.5; cursor: not-allowed; }
        .btn-primary { background: #e94560; color: #fff; }
        .btn-secondary { background: #333; color: #fff; }
        .btn-success { background: #27ae60; color: #fff; }
        .btn-danger { background: #c0392b; color: #fff; }
        .btn-download { background: #9b59b6; color: #fff; }
        .btn-sm { padding: 0.4rem 0.8rem; font-size: 0.8rem; }
        .btn-large { padding: 1rem 2rem; font-size: 1.1rem; }
        .error { background: #c0392b; color: #fff; padding: 1rem; border-radius: 8px; margin-bottom: 1rem; }
        .success { background: #27ae60; color: #fff; padding: 1rem; border-radius: 8px; margin-bottom: 1rem; }
        .hidden { display: none !important; }
        .flex { display: flex; gap: 1rem; align-items: center; }
        .flex-wrap { flex-wrap: wrap; }
        .mt-1 { margin-top: 1rem; }
        .mt-2 { margin-top: 1.5rem; }
        .mb-1 { margin-bottom: 1rem; }
        .mb-2 { margin-bottom: 1.5rem; }
        .text-center { text-align: center; }
        .stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(120px, 1fr)); gap: 1rem; }
        .stat-card { background: #0f0f23; padding: 1rem; border-radius: 8px; text-align: center; }
        .stat-value { font-size: 1.75rem; font-weight: bold; color: #e94560; }
        .stat-label { font-size: 0.85rem; color: #888; }
        .channel-list { list-style: none; }
        .channel-item { background: #0f0f23; padding: 1rem; border-radius: 8px; margin-bottom: 0.75rem; display: flex; justify-content: space-between; align-items: center; cursor: pointer; transition: background 0.2s; }
        .channel-item:hover { background: #1a1a3e; }
        .channel-status, .video-status { font-size: 0.75rem; padding: 0.2rem 0.5rem; border-radius: 4px; font-weight: 500; text-transform: uppercase; letter-spacing: 0.5px; }
        .status-pending { background: #f39c12; color: #000; }
        .status-syncing, .status-downloading { background: #3498db; color: #fff; }
        .status-synced, .status-downloaded { background: #27ae60; color: #fff; }
        .status-error { background: #c0392b; color: #fff; }
        .step { background: #0f0f23; border-radius: 12px; padding: 1.5rem; margin-bottom: 1rem; border-left: 4px solid #333; }
        .step.active { border-left-color: #e94560; }
        .step.done { border-left-color: #27ae60; }
        .step-header { display: flex; align-items: center; gap: 1rem; margin-bottom: 0.75rem; }
        .step-number { width: 32px; height: 32px; border-radius: 50%; background: #333; display: flex; align-items: center; justify-content: center; font-weight: bold; }
        .step.active .step-number { background: #e94560; }
        .step.done .step-number { background: #27ae60; }
        .step-title { font-size: 1.1rem; font-weight: 500; }
        .step-content { margin-left: 48px; color: #aaa; }
        .step-actions { margin-top: 1rem; }
        .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; gap: 1rem; flex-wrap: wrap; }
        .header-left { display: flex; align-items: center; gap: 1rem; }
        .settings-btn { background: none; border: none; color: #888; cursor: pointer; font-size: 1.5rem; }
        .back-btn { background: none; border: none; color: #888; cursor: pointer; font-size: 1.5rem; padding: 0.5rem; }
        .back-btn:hover { color: #fff; }

        /* Video List Styles */
        .video-list { margin-top: 1rem; }
        .video-item { background: #0f0f23; border-radius: 8px; margin-bottom: 0.75rem; display: flex; align-items: stretch; overflow: hidden; transition: background 0.2s, transform 0.2s; }
        .video-item:hover { background: #1a1a3e; transform: translateX(4px); }
        .video-thumb-container { position: relative; flex-shrink: 0; width: 180px; }
        .video-thumb { width: 100%; height: 100%; object-fit: cover; background: #333; aspect-ratio: 16/9; }
        .video-duration-badge { position: absolute; bottom: 4px; right: 4px; background: rgba(0,0,0,0.8); color: #fff; padding: 2px 6px; border-radius: 4px; font-size: 0.75rem; font-weight: 500; }
        .video-content { flex: 1; padding: 1rem; display: flex; flex-direction: column; justify-content: space-between; min-width: 0; }
        .video-title { font-size: 1rem; font-weight: 500; margin-bottom: 0.5rem; color: #fff; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; line-height: 1.4; }
        .video-meta-row { display: flex; align-items: center; gap: 0.75rem; flex-wrap: wrap; margin-bottom: 0.5rem; }
        .video-meta-item { font-size: 0.8rem; color: #888; display: flex; align-items: center; gap: 0.25rem; }
        .video-meta-item svg { width: 14px; height: 14px; }
        .video-description { font-size: 0.85rem; color: #666; display: -webkit-box; -webkit-line-clamp: 1; -webkit-box-orient: vertical; overflow: hidden; margin-bottom: 0.5rem; }
        .video-actions { flex-shrink: 0; display: flex; align-items: center; padding: 1rem; gap: 0.5rem; }

        /* Action Bar */
        .action-bar { display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 1rem; margin-bottom: 1rem; padding: 1rem; background: #0f0f23; border-radius: 8px; }
        .action-bar-left { display: flex; align-items: center; gap: 1rem; flex-wrap: wrap; }
        .action-bar-right { display: flex; align-items: center; gap: 0.5rem; }
        .filter-tabs { display: flex; gap: 0.25rem; background: #16213e; padding: 0.25rem; border-radius: 8px; }
        .filter-tab { padding: 0.5rem 1rem; background: none; border: none; color: #888; cursor: pointer; border-radius: 6px; font-size: 0.9rem; transition: all 0.2s; }
        .filter-tab:hover { color: #fff; }
        .filter-tab.active { background: #e94560; color: #fff; }
        .search-input { padding: 0.5rem 1rem; border: 1px solid #333; border-radius: 8px; background: #16213e; color: #eee; width: 200px; }
        .search-input:focus { outline: none; border-color: #e94560; }

        /* Modal Styles */
        .modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.85); z-index: 100; display: flex; align-items: center; justify-content: center; padding: 1rem; }
        .modal-content { background: #16213e; border-radius: 12px; padding: 0; max-width: 1000px; width: 100%; max-height: 90vh; overflow: auto; }
        .modal-header { padding: 1rem 1.5rem; border-bottom: 1px solid #333; display: flex; justify-content: space-between; align-items: flex-start; gap: 1rem; }
        .modal-title { font-size: 1.25rem; font-weight: 500; line-height: 1.4; }
        .modal-body { padding: 0; }
        .modal-close { background: none; border: none; color: #888; font-size: 1.5rem; cursor: pointer; flex-shrink: 0; }
        .modal-close:hover { color: #fff; }
        .video-player { width: 100%; aspect-ratio: 16/9; background: #000; }
        .video-details { padding: 1.5rem; }
        .video-details-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 1rem; margin-bottom: 1rem; }
        .detail-item { background: #0f0f23; padding: 0.75rem; border-radius: 8px; }
        .detail-label { font-size: 0.75rem; color: #888; text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 0.25rem; }
        .detail-value { font-size: 1rem; color: #fff; font-weight: 500; }

        /* Progress indicator */
        .progress-bar { height: 4px; background: #333; border-radius: 2px; overflow: hidden; margin-top: 0.5rem; }
        .progress-fill { height: 100%; background: #3498db; transition: width 0.3s; }

        /* Notification toast */
        .toast { position: fixed; bottom: 2rem; right: 2rem; background: #27ae60; color: #fff; padding: 1rem 1.5rem; border-radius: 8px; z-index: 200; animation: slideIn 0.3s ease; }
        .toast.error { background: #c0392b; }
        @keyframes slideIn { from { transform: translateY(100px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }

        /* Responsive */
        @media (max-width: 768px) {
            .video-item { flex-direction: column; }
            .video-thumb-container { width: 100%; }
            .video-actions { justify-content: flex-end; padding: 0.75rem 1rem; }
            .action-bar { flex-direction: column; align-items: stretch; }
            .action-bar-left, .action-bar-right { justify-content: center; }
            .search-input { width: 100%; }
        }
    </style>
</head>
<body>
    <div id="app" class="container"><div class="loading" style="text-align:center;padding:2rem;color:#888">Loading...</div></div>
    <script>
        const API = '/api';
        const EXT_URL = 'https://chromewebstore.google.com/detail/get-cookiestxt-locally/cclelndahbckbenkjhflpdbgdldlbecc';
        let step = 1;
        let currentView = 'main';
        let currentChannel = null;
        let currentFilter = 'all';
        let searchQuery = '';

        const api = {
            cookies: () => fetch(API+'/cookies').then(r=>r.json()),
            saveCookies: c => fetch(API+'/cookies',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({cookies:c})}).then(r=>{if(!r.ok)throw r;return r.json()}),
            delCookies: () => fetch(API+'/cookies',{method:'DELETE'}),
            stats: () => fetch(API+'/stats').then(r=>r.json()),
            channels: () => fetch(API+'/channels').then(r=>r.json()),
            channel: id => fetch(API+'/channels/'+id).then(r=>r.json()),
            channelVideos: id => fetch(API+'/channels/'+id+'/videos').then(r=>r.json()),
            addChannel: u => fetch(API+'/channels',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({youtube_url:u})}).then(r=>{if(!r.ok)throw r;return r.json()}),
            syncChannel: id => fetch(API+'/channels/'+id+'/sync',{method:'POST'}).then(r=>{if(!r.ok)throw r;return r.json()}),
            delChannel: id => fetch(API+'/channels/'+id,{method:'DELETE'}),
            downloadVideo: id => fetch(API+'/videos/'+id+'/download',{method:'POST'}).then(r=>{if(!r.ok)throw r;return r.json()})
        };

        function goStep(n) { step=n; document.querySelectorAll('.step').forEach((e,i)=>{e.classList.remove('active','done');if(i+1<n)e.classList.add('done');if(i+1===n)e.classList.add('active');}); }

        function formatDuration(s) {
            if(!s) return '--:--';
            const h = Math.floor(s/3600);
            const m = Math.floor((s%3600)/60);
            const sec = s%60;
            if(h > 0) return h+':'+(m<10?'0':'')+m+':'+(sec<10?'0':'')+sec;
            return m+':'+(sec<10?'0':'')+sec;
        }
        function formatSize(b) { if(!b) return '-'; if(b>1e9) return (b/1e9).toFixed(2)+' GB'; if(b>1e6) return (b/1e6).toFixed(1)+' MB'; return (b/1e3).toFixed(1)+' KB'; }
        function escapeHtml(str) { if(!str) return ''; return str.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;'); }

        function showToast(msg, isError=false) {
            const existing = document.querySelector('.toast');
            if(existing) existing.remove();
            const toast = document.createElement('div');
            toast.className = 'toast' + (isError ? ' error' : '');
            toast.textContent = msg;
            document.body.appendChild(toast);
            setTimeout(() => toast.remove(), 3000);
        }

        function setupPage() {
            return '<div class="setup-card card"><h1 class="text-center mb-2">YT Archive Setup</h1><p class="text-center mb-2" style="color:#888">Follow these steps to configure YouTube authentication</p><div id="err" class="error hidden"></div><div id="ok" class="success hidden"></div><div class="step active" data-step="1"><div class="step-header"><div class="step-number">1</div><div class="step-title">Install Browser Extension</div></div><div class="step-content"><p>Install the "Get cookies.txt LOCALLY" extension for Chrome/Edge.</p><div class="step-actions"><a href="'+EXT_URL+'" target="_blank" class="btn btn-primary">Install Extension</a> <button class="btn btn-secondary" id="s1done">I\'ve Installed It</button></div></div></div><div class="step" data-step="2"><div class="step-header"><div class="step-number">2</div><div class="step-title">Export YouTube Cookies</div></div><div class="step-content"><p>Open YouTube, sign in if needed, then click the extension icon and export cookies.</p><div class="step-actions"><a href="https://www.youtube.com" target="_blank" class="btn btn-primary">Open YouTube</a> <button class="btn btn-secondary" id="s2done">I Have the Cookies</button></div></div></div><div class="step" data-step="3"><div class="step-header"><div class="step-number">3</div><div class="step-title">Paste Cookies</div></div><div class="step-content"><p>Open the exported cookies.txt file and paste its contents below:</p><div class="step-actions"><textarea id="cookies" placeholder="# Netscape HTTP Cookie File..."></textarea><div class="mt-1"><button id="save" class="btn btn-success btn-large">Save & Continue</button></div></div></div></div></div>';
        }

        function mainPage(stats, channels) {
            const ch = channels.channels || [];
            return '<div class="header"><h1>YT Archive</h1><button id="settings" class="settings-btn">&#9881;</button></div><div class="card"><div class="stats-grid"><div class="stat-card"><div class="stat-value">'+(stats.total_channels||0)+'</div><div class="stat-label">Channels</div></div><div class="stat-card"><div class="stat-value">'+(stats.total_videos||0)+'</div><div class="stat-label">Videos</div></div><div class="stat-card"><div class="stat-value">'+(stats.downloaded_videos||0)+'</div><div class="stat-label">Downloaded</div></div><div class="stat-card"><div class="stat-value">'+(stats.pending_videos||0)+'</div><div class="stat-label">Pending</div></div><div class="stat-card"><div class="stat-value">'+(stats.storage_used_human||'0 B')+'</div><div class="stat-label">Storage</div></div></div></div><div class="card"><h2>Add Channel</h2><div id="aerr" class="error hidden"></div><div class="flex mt-1"><input type="text" id="url" placeholder="https://youtube.com/@channelname"><button id="add" class="btn btn-primary">Add</button></div></div><div class="card"><h2>Channels</h2>'+(ch.length?'':'<p style="color:#888">No channels added yet</p>')+'<ul class="channel-list mt-1">'+ch.map(c=>'<li class="channel-item" data-id="'+c.id+'"><div class="channel-info"><span class="channel-name" style="font-weight:500">'+(c.name||c.youtube_id)+'</span> <span class="channel-status status-'+c.status+'">'+c.status+'</span><span style="color:#666;margin-left:0.5rem;font-size:0.85rem">'+c.video_count+' videos</span></div><div class="flex" onclick="event.stopPropagation()"><button class="btn btn-sm btn-secondary sync" data-id="'+c.id+'" '+(c.status==='syncing'?'disabled':'')+'>Sync</button><button class="btn btn-sm btn-danger del" data-id="'+c.id+'">Delete</button></div></li>').join('')+'</ul></div><div id="modal" class="card hidden" style="position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:100;width:400px;"><h2>Settings</h2><p class="mt-1" style="color:#888">Cookies are currently configured.</p><div class="mt-2 flex"><button id="delcookies" class="btn btn-danger">Delete Cookies</button><button id="closemodal" class="btn btn-secondary">Close</button></div></div>';
        }

        function channelPage(channel, videos) {
            const vids = videos.videos || [];
            const downloaded = vids.filter(v=>v.status==='downloaded');
            const downloading = vids.filter(v=>v.status==='downloading');
            const pending = vids.filter(v=>v.status==='pending');
            const errors = vids.filter(v=>v.status==='error');
            const notDownloaded = vids.filter(v=>v.status!=='downloaded');

            let filteredVids = vids;
            if(currentFilter === 'downloaded') filteredVids = downloaded;
            else if(currentFilter === 'pending') filteredVids = notDownloaded;
            else if(currentFilter === 'downloading') filteredVids = downloading;
            else if(currentFilter === 'error') filteredVids = errors;

            if(searchQuery) {
                const q = searchQuery.toLowerCase();
                filteredVids = filteredVids.filter(v => (v.title||'').toLowerCase().includes(q) || (v.description||'').toLowerCase().includes(q));
            }

            return '<div class="header"><div class="header-left"><button id="back" class="back-btn">&larr;</button><div><h1 style="margin-bottom:0.25rem">'+(channel.name||channel.youtube_id)+'</h1><div class="flex" style="gap:0.5rem"><span class="channel-status status-'+channel.status+'">'+channel.status+'</span><span style="color:#888;font-size:0.9rem">'+vids.length+' videos</span></div></div></div></div>'+
            '<div class="card"><div class="stats-grid"><div class="stat-card"><div class="stat-value">'+vids.length+'</div><div class="stat-label">Total</div></div><div class="stat-card"><div class="stat-value" style="color:#27ae60">'+downloaded.length+'</div><div class="stat-label">Downloaded</div></div><div class="stat-card"><div class="stat-value" style="color:#3498db">'+downloading.length+'</div><div class="stat-label">Downloading</div></div><div class="stat-card"><div class="stat-value" style="color:#f39c12">'+pending.length+'</div><div class="stat-label">Pending</div></div>'+(errors.length?'<div class="stat-card"><div class="stat-value" style="color:#c0392b">'+errors.length+'</div><div class="stat-label">Errors</div></div>':'')+'</div></div>'+
            '<div class="card"><div class="action-bar"><div class="action-bar-left"><div class="filter-tabs"><button class="filter-tab'+(currentFilter==='all'?' active':'')+'" data-filter="all">All ('+vids.length+')</button><button class="filter-tab'+(currentFilter==='downloaded'?' active':'')+'" data-filter="downloaded">Downloaded ('+downloaded.length+')</button><button class="filter-tab'+(currentFilter==='pending'?' active':'')+'" data-filter="pending">Not Downloaded ('+notDownloaded.length+')</button></div><input type="text" class="search-input" id="searchInput" placeholder="Search videos..." value="'+escapeHtml(searchQuery)+'"></div><div class="action-bar-right">'+(notDownloaded.length > 0 ? '<button id="downloadAll" class="btn btn-download"><svg style="width:16px;height:16px;margin-right:0.5rem;vertical-align:middle" viewBox="0 0 24 24" fill="currentColor"><path d="M19 9h-4V3H9v6H5l7 7 7-7zM5 18v2h14v-2H5z"/></svg>Download All ('+notDownloaded.length+')</button>' : '')+'<button id="syncChannel" class="btn btn-secondary" '+(channel.status==='syncing'?'disabled':'')+'>Sync Channel</button></div></div>'+
            '<div class="video-list" id="videoList">'+filteredVids.map(v => {
                const isDownloaded = v.status === 'downloaded';
                const isDownloading = v.status === 'downloading';
                const isPending = v.status === 'pending';
                const isError = v.status === 'error';
                return '<div class="video-item" data-id="'+v.id+'" data-status="'+v.status+'">'+
                    '<div class="video-thumb-container">'+
                        '<img class="video-thumb" src="/api/videos/'+v.id+'/thumbnail" onerror="this.src=\'https://img.youtube.com/vi/'+v.youtube_id+'/mqdefault.jpg\'" loading="lazy">'+
                        '<span class="video-duration-badge">'+formatDuration(v.duration)+'</span>'+
                    '</div>'+
                    '<div class="video-content">'+
                        '<div class="video-title">'+escapeHtml(v.title||v.youtube_id||v.id)+'</div>'+
                        '<div class="video-meta-row">'+
                            '<span class="video-status status-'+v.status+'">'+v.status+'</span>'+
                            (v.file_size ? '<span class="video-meta-item">'+formatSize(v.file_size)+'</span>' : '')+
                            '<span class="video-meta-item"><a href="https://youtube.com/watch?v='+v.youtube_id+'" target="_blank" style="color:#888;text-decoration:none" onclick="event.stopPropagation()">YouTube</a></span>'+
                        '</div>'+
                        (v.description ? '<div class="video-description">'+escapeHtml(v.description)+'</div>' : '')+
                    '</div>'+
                    '<div class="video-actions" onclick="event.stopPropagation()">'+
                        (isDownloaded ? '<a href="/api/videos/'+v.id+'/stream" class="btn btn-sm btn-success" download>Save</a><button class="btn btn-sm btn-primary play-btn" data-id="'+v.id+'">Play</button>' : '')+
                        (isPending || isError ? '<button class="btn btn-sm btn-download download-btn" data-id="'+v.id+'">Download</button>' : '')+
                        (isDownloading ? '<span style="color:#3498db;font-size:0.85rem">Downloading...</span>' : '')+
                    '</div>'+
                '</div>';
            }).join('')+'</div>'+
            (filteredVids.length === 0 ? '<p style="text-align:center;color:#888;padding:2rem">No videos match your filter</p>' : '')+
            '</div>';
        }

        function videoModal(video) {
            const isDownloaded = video.status === 'downloaded';
            const isPending = video.status === 'pending';
            const isDownloading = video.status === 'downloading';
            const isError = video.status === 'error';
            return '<div class="modal-overlay" id="videoModal"><div class="modal-content"><div class="modal-header"><div class="modal-title">'+escapeHtml(video.title||video.id)+'</div><button class="modal-close" id="closeVideo">&times;</button></div><div class="modal-body">'+
                (isDownloaded ? '<video class="video-player" controls autoplay><source src="/api/videos/'+video.id+'/stream" type="video/mp4">Your browser does not support video.</video>' : '<div style="aspect-ratio:16/9;background:#000;display:flex;align-items:center;justify-content:center;flex-direction:column;gap:1rem"><img src="https://img.youtube.com/vi/'+video.youtube_id+'/maxresdefault.jpg" style="max-width:100%;max-height:100%;object-fit:contain" onerror="this.style.display=\'none\'"><p style="color:#888">'+(isDownloading?'Downloading...':'Video not downloaded yet')+'</p></div>')+
                '<div class="video-details">'+
                    '<div class="video-details-grid">'+
                        '<div class="detail-item"><div class="detail-label">Status</div><div class="detail-value"><span class="video-status status-'+video.status+'">'+video.status+'</span></div></div>'+
                        '<div class="detail-item"><div class="detail-label">Duration</div><div class="detail-value">'+formatDuration(video.duration)+'</div></div>'+
                        (video.file_size ? '<div class="detail-item"><div class="detail-label">File Size</div><div class="detail-value">'+formatSize(video.file_size)+'</div></div>' : '')+
                        '<div class="detail-item"><div class="detail-label">YouTube ID</div><div class="detail-value"><a href="https://youtube.com/watch?v='+video.youtube_id+'" target="_blank" style="color:#e94560">'+video.youtube_id+'</a></div></div>'+
                    '</div>'+
                    (video.description ? '<div style="margin-top:1rem"><div class="detail-label" style="margin-bottom:0.5rem">Description</div><p style="color:#aaa;font-size:0.9rem;line-height:1.6;white-space:pre-wrap">'+escapeHtml(video.description)+'</p></div>' : '')+
                    '<div class="mt-2 flex">'+
                        (isDownloaded ? '<a href="/api/videos/'+video.id+'/stream" download class="btn btn-success">Download File</a>' : '')+
                        (isPending || isError ? '<button id="modalDownload" class="btn btn-download" data-id="'+video.id+'">Download Video</button>' : '')+
                        (isDownloading ? '<span style="color:#3498db">Download in progress...</span>' : '')+
                    '</div>'+
                '</div>'+
            '</div></div></div>';
        }

        function bindSetup() {
            document.getElementById('s1done').onclick = () => goStep(2);
            document.getElementById('s2done').onclick = () => { goStep(3); document.getElementById('cookies').focus(); };
            document.getElementById('cookies').onfocus = () => { if(step<3) goStep(3); };
            document.getElementById('save').onclick = async () => {
                const btn = document.getElementById('save'), txt = document.getElementById('cookies').value.trim();
                if(!txt) { document.getElementById('err').textContent='Please paste cookies'; document.getElementById('err').classList.remove('hidden'); return; }
                btn.disabled=true; btn.textContent='Saving...';
                try { await api.saveCookies(txt); document.getElementById('ok').textContent='Saved! Loading...'; document.getElementById('ok').classList.remove('hidden'); setTimeout(()=>location.reload(),1000); }
                catch(e) { const d=await e.json().catch(()=>({})); document.getElementById('err').textContent=d.error||'Failed'; document.getElementById('err').classList.remove('hidden'); btn.disabled=false; btn.textContent='Save & Continue'; }
            };
        }

        function bindMain() {
            document.getElementById('add').onclick = async () => {
                const url = document.getElementById('url').value.trim(); if(!url) return;
                try { await api.addChannel(url); location.reload(); } catch(e) { document.getElementById('aerr').textContent=(await e.json().catch(()=>({}))).error||'Failed'; document.getElementById('aerr').classList.remove('hidden'); }
            };
            document.querySelectorAll('.sync').forEach(b => b.onclick = async (e) => { e.stopPropagation(); b.disabled=true; try { await api.syncChannel(b.dataset.id); location.reload(); } catch(e) { alert('Failed'); b.disabled=false; } });
            document.querySelectorAll('.del').forEach(b => b.onclick = async (e) => { e.stopPropagation(); if(!confirm('Delete?')) return; await api.delChannel(b.dataset.id); location.reload(); });
            document.querySelectorAll('.channel-item').forEach(item => item.onclick = () => showChannel(item.dataset.id));
            document.getElementById('settings').onclick = () => document.getElementById('modal').classList.remove('hidden');
            document.getElementById('closemodal').onclick = () => document.getElementById('modal').classList.add('hidden');
            document.getElementById('delcookies').onclick = async () => { if(!confirm('Delete cookies?')) return; await api.delCookies(); location.reload(); };
        }

        function bindChannel(channelData, videosData) {
            document.getElementById('back').onclick = () => { currentView='main'; currentChannel=null; currentFilter='all'; searchQuery=''; render(); };

            // Filter tabs
            document.querySelectorAll('.filter-tab').forEach(t => t.onclick = () => {
                currentFilter = t.dataset.filter;
                showChannel(currentChannel);
            });

            // Search
            const searchInput = document.getElementById('searchInput');
            if(searchInput) {
                let debounceTimer;
                searchInput.oninput = () => {
                    clearTimeout(debounceTimer);
                    debounceTimer = setTimeout(() => {
                        searchQuery = searchInput.value;
                        showChannel(currentChannel);
                    }, 300);
                };
            }

            // Sync channel button
            const syncBtn = document.getElementById('syncChannel');
            if(syncBtn) {
                syncBtn.onclick = async () => {
                    syncBtn.disabled = true;
                    syncBtn.textContent = 'Syncing...';
                    try {
                        await api.syncChannel(currentChannel);
                        showToast('Sync started');
                        showChannel(currentChannel);
                    } catch(e) {
                        showToast('Failed to start sync', true);
                        syncBtn.disabled = false;
                        syncBtn.textContent = 'Sync Channel';
                    }
                };
            }

            // Download All button
            const downloadAllBtn = document.getElementById('downloadAll');
            if(downloadAllBtn) {
                downloadAllBtn.onclick = async () => {
                    const notDownloaded = (videosData.videos||[]).filter(v => v.status !== 'downloaded' && v.status !== 'downloading');
                    if(notDownloaded.length === 0) return;

                    downloadAllBtn.disabled = true;
                    downloadAllBtn.textContent = 'Queueing...';

                    let queued = 0;
                    for(const video of notDownloaded) {
                        try {
                            await api.downloadVideo(video.id);
                            queued++;
                        } catch(e) {
                            console.error('Failed to queue', video.id, e);
                        }
                    }

                    showToast(queued + ' videos queued for download');
                    showChannel(currentChannel);
                };
            }

            // Individual download buttons
            document.querySelectorAll('.download-btn').forEach(btn => {
                btn.onclick = async (e) => {
                    e.stopPropagation();
                    btn.disabled = true;
                    btn.textContent = 'Queueing...';
                    try {
                        await api.downloadVideo(btn.dataset.id);
                        showToast('Video queued for download');
                        showChannel(currentChannel);
                    } catch(e) {
                        const msg = (await e.json().catch(()=>({}))).error || 'Failed to queue download';
                        showToast(msg, true);
                        btn.disabled = false;
                        btn.textContent = 'Download';
                    }
                };
            });

            // Play buttons
            document.querySelectorAll('.play-btn').forEach(btn => {
                btn.onclick = (e) => {
                    e.stopPropagation();
                    const video = (videosData.videos||[]).find(x=>x.id===btn.dataset.id);
                    if(video) showVideoModal(video);
                };
            });

            // Video item click
            document.querySelectorAll('.video-item').forEach(v => v.onclick = async () => {
                const video = (videosData.videos||[]).find(x=>x.id===v.dataset.id);
                if(video) showVideoModal(video);
            });
        }

        function showVideoModal(video) {
            const modal = document.createElement('div');
            modal.innerHTML = videoModal(video);
            document.body.appendChild(modal.firstChild);

            document.getElementById('closeVideo').onclick = () => document.getElementById('videoModal').remove();
            document.getElementById('videoModal').onclick = (e) => { if(e.target.id==='videoModal') document.getElementById('videoModal').remove(); };

            const modalDownloadBtn = document.getElementById('modalDownload');
            if(modalDownloadBtn) {
                modalDownloadBtn.onclick = async () => {
                    modalDownloadBtn.disabled = true;
                    modalDownloadBtn.textContent = 'Queueing...';
                    try {
                        await api.downloadVideo(video.id);
                        showToast('Video queued for download');
                        document.getElementById('videoModal').remove();
                        showChannel(currentChannel);
                    } catch(e) {
                        const msg = (await e.json().catch(()=>({}))).error || 'Failed to queue download';
                        showToast(msg, true);
                        modalDownloadBtn.disabled = false;
                        modalDownloadBtn.textContent = 'Download Video';
                    }
                };
            }
        }

        async function showChannel(id) {
            currentView='channel'; currentChannel=id;
            const app = document.getElementById('app');
            if(!document.querySelector('.video-list')) {
                app.innerHTML = '<div style="text-align:center;padding:2rem;color:#888">Loading...</div>';
            }
            try {
                const [ch, vids] = await Promise.all([api.channel(id), api.channelVideos(id)]);
                app.innerHTML = channelPage(ch.channel, vids);
                bindChannel(ch, vids);
            } catch(e) {
                console.error(e);
                app.innerHTML = '<p style="color:#c0392b;text-align:center;padding:2rem">Error loading channel</p>';
            }
        }

        async function render() {
            const app = document.getElementById('app');
            const {configured} = await api.cookies().catch(()=>({configured:false}));
            if(!configured) { app.innerHTML = setupPage(); bindSetup(); return; }
            if(currentView==='channel' && currentChannel) { showChannel(currentChannel); return; }
            const [s,c] = await Promise.all([api.stats(),api.channels()]);
            app.innerHTML = mainPage(s,c);
            bindMain();
        }

        render();
    </script>
</body>
</html>
`
