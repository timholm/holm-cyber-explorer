package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/timholm/ytarchive/internal/downloader"
	"github.com/timholm/ytarchive/internal/logging"
	"github.com/timholm/ytarchive/internal/youtube"
)

const (
	// Unified queue key for all channels - KEDA monitors this
	unifiedQueueKey = "ytarchive:download:queue"

	// How long to wait when queue is empty before checking again
	emptyQueueWaitTime = 5 * time.Second

	// Health check server port
	healthCheckPort = 8080

	// Default delay between video downloads to avoid YouTube rate limiting
	defaultDownloadDelaySeconds = 2
)

// WorkerConfig holds the worker configuration from environment variables
type WorkerConfig struct {
	RedisURL      string
	StoragePath   string
	WorkerID      string
	ControllerURL string
	CollectorURL  string
}

// healthStatus tracks worker health for readiness probes
type healthStatus struct {
	ready   atomic.Bool
	healthy atomic.Bool
}

func main() {
	logging.Info("starting YouTube Channel Archiver Worker (continuous mode)")

	// Load configuration from environment
	config, err := loadConfig()
	if err != nil {
		logging.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	logging.Info("worker configuration loaded",
		"worker_id", config.WorkerID,
		"storage_path", config.StoragePath,
	)

	// Check if ffmpeg is available for merging (optional)
	if downloader.MergerAvailable() {
		logging.Info("ffmpeg available for stream merging")
	} else {
		logging.Warn("ffmpeg not found - will use combined streams only (limited quality)")
	}

	// Initialize health status
	health := &healthStatus{}
	health.healthy.Store(true)

	// Start health check server
	go startHealthServer(health)

	// Connect to Redis
	redisClient, err := connectRedis(config.RedisURL)
	if err != nil {
		logging.Error("failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()
	logging.Info("connected to Redis")

	// Set up context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		logging.Info("received shutdown signal, finishing current video...", "signal", sig.String())
		cancel()
	}()

	// Create HTTP progress reporter (for controller API)
	var reporter *downloader.ProgressReporter
	if config.ControllerURL != "" {
		reporter = downloader.NewProgressReporter(config.ControllerURL, config.WorkerID)
	}

	// Create Redis progress reporter (for UI polling)
	redisProgressReporter := downloader.NewRedisProgressReporter(redisClient, config.WorkerID)
	logging.Info("Redis progress reporter initialized")

	// Load cookies for authenticated requests (bypass bot detection)
	var ytClientOpts []youtube.ClientOption
	var cookieHeader string

	// Check for cookies in Redis first (set via web UI)
	redisCookies, err := redisClient.Get(ctx, "config:cookies").Result()
	if err == nil && redisCookies != "" {
		cookies, err := youtube.LoadCookiesFromString(redisCookies)
		if err != nil {
			logging.Warn("failed to parse cookies from Redis", "error", err)
		} else if len(cookies) > 0 {
			ytClientOpts = append(ytClientOpts, youtube.WithCookies(cookies))
			cookieHeader = youtube.GetCookieHeader(cookies)
			logging.Info("loaded cookies from Redis", "cookie_count", len(cookies))
		}
	}

	// Fallback: Check for cookies file path
	if cookieHeader == "" {
		if cookiesFile := os.Getenv("YOUTUBE_COOKIES_FILE"); cookiesFile != "" {
			cookies, err := youtube.LoadCookiesFromFile(cookiesFile)
			if err != nil {
				logging.Warn("failed to load cookies from file", "path", cookiesFile, "error", err)
			} else if len(cookies) > 0 {
				ytClientOpts = append(ytClientOpts, youtube.WithCookies(cookies))
				cookieHeader = youtube.GetCookieHeader(cookies)
				logging.Info("loaded cookies from file", "path", cookiesFile, "cookie_count", len(cookies))
			}
		}
	}

	// Fallback: Check for cookies in environment variable (Netscape format)
	if cookiesEnv := os.Getenv("YOUTUBE_COOKIES"); cookiesEnv != "" && cookieHeader == "" {
		cookies, err := youtube.LoadCookiesFromString(cookiesEnv)
		if err != nil {
			logging.Warn("failed to load cookies from environment", "error", err)
		} else if len(cookies) > 0 {
			ytClientOpts = append(ytClientOpts, youtube.WithCookies(cookies))
			cookieHeader = youtube.GetCookieHeader(cookies)
			logging.Info("loaded cookies from environment", "cookie_count", len(cookies))
		}
	}

	// Create YouTube client for fetching stream URLs
	ytClient, err := youtube.NewClient(ytClientOpts...)
	if err != nil {
		logging.Error("failed to create YouTube client", "error", err)
		os.Exit(1)
	}
	logging.Info("YouTube client initialized", "has_cookies", cookieHeader != "")

	// Create downloader with default config
	dlConfig := downloader.DefaultConfig(config.StoragePath)

	// If ffmpeg is not available, we must use combined streams (audio+video in one file)
	// Otherwise downloads will result in video-only files without audio
	if !downloader.MergerAvailable() {
		dlConfig.PreferCombinedStream = true
		dlConfig.MaxHeight = 720 // Combined streams max out at 720p (format 22)
		logging.Info("ffmpeg not available, limiting to combined streams (max 720p)")
	}

	dl := downloader.NewDownloader(dlConfig, reporter, config.WorkerID)

	// Set Redis progress reporter for UI polling
	dl.SetRedisReporter(redisProgressReporter)

	// Pass cookies to downloader
	if cookieHeader != "" {
		dl.SetCookies(cookieHeader)
	}

	// Mark worker as ready
	health.ready.Store(true)
	logging.Info("worker ready, starting continuous processing loop")

	// Statistics
	successCount := 0
	failCount := 0

	// Get download delay from environment (helps avoid YouTube rate limiting)
	downloadDelay := time.Duration(getEnvInt("DOWNLOAD_DELAY_SECONDS", defaultDownloadDelaySeconds)) * time.Second
	logging.Info("download delay configured", "delay", downloadDelay)

	// Main processing loop - runs continuously until shutdown
	for {
		// Check if context is cancelled (shutdown requested)
		if ctx.Err() != nil {
			logging.Info("shutdown requested, exiting processing loop")
			break
		}

		// Claim ONE video from the unified queue
		channelID, videoID, err := claimOneVideo(ctx, redisClient)
		if err != nil {
			if ctx.Err() != nil {
				break // Context cancelled during claim
			}
			logging.Error("failed to claim video from queue", "error", err)
			time.Sleep(emptyQueueWaitTime)
			continue
		}

		if videoID == "" {
			// Queue is empty, wait before checking again
			logging.Debug("queue empty, waiting...", "wait_time", emptyQueueWaitTime)
			time.Sleep(emptyQueueWaitTime)
			continue
		}

		logging.Info("processing video",
			"worker_id", config.WorkerID,
			"channel_id", channelID,
			"video_id", videoID,
		)

		// Process the video
		success := processVideo(ctx, config, redisClient, ytClient, dl, reporter, channelID, videoID)
		if success {
			successCount++
		} else {
			failCount++
		}

		logging.Info("video processing complete",
			"worker_id", config.WorkerID,
			"channel_id", channelID,
			"video_id", videoID,
			"success", success,
			"total_success", successCount,
			"total_failed", failCount,
		)

		// Add delay between downloads to avoid YouTube rate limiting
		if downloadDelay > 0 && ctx.Err() == nil {
			logging.Debug("waiting before next download", "delay", downloadDelay)
			time.Sleep(downloadDelay)
		}
	}

	// Report final statistics
	logging.Info("worker shutting down",
		"worker_id", config.WorkerID,
		"success_count", successCount,
		"failed_count", failCount,
	)
}

// VideoInfo holds video metadata from Redis
type VideoInfo struct {
	ID            string `json:"id"`
	YouTubeID     string `json:"youtube_id"`
	ChannelID     string `json:"channel_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Duration      int    `json:"duration"`
	UploadDate    string `json:"upload_date"`
	ThumbnailURL  string `json:"thumbnail_url"`
	ViewCount     int64  `json:"view_count"`
	EpisodeNumber int    `json:"episode_number"`
	ChannelName   string `json:"channel_name"`
	Status        string `json:"status"`
}

// processVideo downloads a single video and returns success/failure
func processVideo(ctx context.Context, config *WorkerConfig, redisClient *redis.Client, ytClient *youtube.Client, dl *downloader.Downloader, reporter *downloader.ProgressReporter, channelID, videoID string) bool {
	// Report download starting
	if reporter != nil {
		status := &downloader.VideoStatus{
			VideoID: videoID,
			Status:  "downloading",
		}
		if err := reporter.ReportVideoStatus(status); err != nil {
			logging.Warn("failed to report video status",
				"worker_id", config.WorkerID,
				"video_id", videoID,
				"error", err,
			)
		}
	}

	// Fetch video info from Redis (for episode number and channel name)
	var videoInfo *VideoInfo
	if redisClient != nil {
		info, err := fetchVideoInfo(ctx, redisClient, channelID, videoID)
		if err != nil {
			logging.Warn("failed to fetch video info from Redis, will use defaults",
				"worker_id", config.WorkerID,
				"video_id", videoID,
				"error", err,
			)
		} else {
			videoInfo = info
		}
	}

	// Fetch stream URLs from YouTube
	streamInfo, err := ytClient.GetStreamURLContext(ctx, videoID)
	if err != nil {
		logging.Error("failed to fetch stream info",
			"worker_id", config.WorkerID,
			"video_id", videoID,
			"error", err,
		)
		if reporter != nil {
			status := &downloader.VideoStatus{
				VideoID: videoID,
				Status:  "error",
				Error:   fmt.Sprintf("failed to fetch stream info: %v", err),
			}
			reporter.ReportVideoStatus(status)
		}
		return false
	}

	// Convert YouTube formats to downloader streams
	streams := convertFormatsToStreams(streamInfo.Formats)
	if len(streams) == 0 {
		logging.Error("no downloadable streams found",
			"worker_id", config.WorkerID,
			"video_id", videoID,
		)
		if reporter != nil {
			status := &downloader.VideoStatus{
				VideoID: videoID,
				Status:  "error",
				Error:   "no downloadable streams found",
			}
			reporter.ReportVideoStatus(status)
		}
		return false
	}

	logging.Info("found streams",
		"worker_id", config.WorkerID,
		"video_id", videoID,
		"stream_count", len(streams),
	)

	// Create download request with stream URLs and video info
	req := &downloader.DownloadRequest{
		VideoID:   videoID,
		ChannelID: channelID,
		Title:     streamInfo.Title,
		Streams:   streams,
	}

	// Populate additional info from Redis if available
	if videoInfo != nil {
		if req.Title == "" {
			req.Title = videoInfo.Title
		}
		req.ChannelName = videoInfo.ChannelName
		req.EpisodeNumber = videoInfo.EpisodeNumber
		req.Description = videoInfo.Description
		req.Duration = videoInfo.Duration
		req.UploadDate = videoInfo.UploadDate
		req.ThumbnailURL = videoInfo.ThumbnailURL
		req.ViewCount = videoInfo.ViewCount
	}

	// Download the video
	result := dl.Download(ctx, req)

	if result.Success {
		logging.Info("successfully downloaded video",
			"worker_id", config.WorkerID,
			"channel_id", channelID,
			"video_id", videoID,
			"file_path", result.FilePath,
			"file_size", result.FileSize,
			"duration", result.Duration.String(),
		)

		// Upload to collector for permanent storage
		if config.CollectorURL != "" {
			uploadMeta := &UploadMetadata{
				VideoID:       videoID,
				ChannelID:     channelID,
				ChannelName:   req.ChannelName,
				Title:         req.Title,
				Description:   req.Description,
				Duration:      req.Duration,
				UploadDate:    req.UploadDate,
				EpisodeNumber: req.EpisodeNumber,
				Filename:      filepath.Base(result.FilePath),
				FileSize:      result.FileSize,
			}

			// Upload with retry logic
			var uploadErr error
			maxUploadRetries := 3
			for attempt := 1; attempt <= maxUploadRetries; attempt++ {
				uploadErr = uploadToCollector(ctx, config.CollectorURL, result.FilePath, uploadMeta)
				if uploadErr == nil {
					break
				}
				logging.Warn("upload attempt failed, retrying",
					"worker_id", config.WorkerID,
					"video_id", videoID,
					"attempt", attempt,
					"max_attempts", maxUploadRetries,
					"error", uploadErr,
				)
				if attempt < maxUploadRetries {
					// Exponential backoff: 5s, 15s, 45s
					backoff := time.Duration(5*attempt*attempt) * time.Second
					select {
					case <-ctx.Done():
						uploadErr = ctx.Err()
						break
					case <-time.After(backoff):
					}
				}
			}
			if uploadErr != nil {
				logging.Error("failed to upload to collector after retries",
					"worker_id", config.WorkerID,
					"video_id", videoID,
					"error", uploadErr,
				)
				// Continue - report as error since upload failed
				if reporter != nil {
					status := &downloader.VideoStatus{
						VideoID: videoID,
						Status:  "error",
						Error:   fmt.Sprintf("upload to collector failed: %v", uploadErr),
					}
					reporter.ReportVideoStatus(status)
				}
				return false
			}

			logging.Info("uploaded video to collector",
				"worker_id", config.WorkerID,
				"video_id", videoID,
			)

			// Clean up local file after successful upload
			os.RemoveAll(filepath.Dir(result.FilePath))
		}

		// Report success
		if reporter != nil {
			status := &downloader.VideoStatus{
				VideoID:  videoID,
				Status:   "downloaded",
				FilePath: result.FilePath,
				FileSize: result.FileSize,
			}
			if err := reporter.ReportVideoStatus(status); err != nil {
				logging.Warn("failed to report video status",
					"worker_id", config.WorkerID,
					"video_id", videoID,
					"error", err,
				)
			}
		}
		return true
	}

	logging.Error("failed to download video",
		"worker_id", config.WorkerID,
		"channel_id", channelID,
		"video_id", videoID,
		"attempts", result.Attempts,
		"error", result.Error,
	)

	// Report failure
	if reporter != nil {
		status := &downloader.VideoStatus{
			VideoID: videoID,
			Status:  "error",
			Error:   result.Error.Error(),
		}
		if err := reporter.ReportVideoStatus(status); err != nil {
			logging.Warn("failed to report video status",
				"worker_id", config.WorkerID,
				"video_id", videoID,
				"error", err,
			)
		}
	}

	// Clean up partial downloads
	if err := dl.Cleanup(channelID, videoID); err != nil {
		logging.Warn("failed to cleanup partial download",
			"worker_id", config.WorkerID,
			"video_id", videoID,
			"error", err,
		)
	}

	return false
}

// startHealthServer starts an HTTP server for health checks
func startHealthServer(health *healthStatus) {
	mux := http.NewServeMux()

	// Liveness probe - always returns 200 if the server is running
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if health.healthy.Load() {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("unhealthy"))
		}
	})

	// Readiness probe - returns 200 only when worker is ready to process
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if health.ready.Load() {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ready"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("not ready"))
		}
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", healthCheckPort),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	logging.Info("starting health check server", "port", healthCheckPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logging.Error("health check server error", "error", err)
	}
}

// loadConfig loads worker configuration from environment variables
func loadConfig() (*WorkerConfig, error) {
	config := &WorkerConfig{
		RedisURL:      os.Getenv("REDIS_URL"),
		StoragePath:   os.Getenv("STORAGE_PATH"),
		WorkerID:      os.Getenv("WORKER_ID"),
		ControllerURL: os.Getenv("CONTROLLER_URL"),
		CollectorURL:  os.Getenv("COLLECTOR_URL"),
	}

	// Validate required fields
	if config.RedisURL == "" {
		// Default to localhost if not specified
		config.RedisURL = "localhost:6379"
	}

	if config.StoragePath == "" {
		config.StoragePath = "/tmp/downloads"
	}

	if config.WorkerID == "" {
		// Generate a default worker ID from hostname and timestamp
		hostname, _ := os.Hostname()
		if hostname == "" {
			hostname = "worker"
		}
		config.WorkerID = fmt.Sprintf("%s-%d", hostname, time.Now().UnixNano()%10000)
	}

	if config.CollectorURL == "" {
		config.CollectorURL = "http://collector.ytarchive.svc.cluster.local:8081"
	}

	return config, nil
}

// connectRedis creates and tests a Redis connection
func connectRedis(redisURL string) (*redis.Client, error) {
	// Parse Redis URL (supports redis://host:port or just host:port)
	opts := &redis.Options{
		Addr:         redisURL,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     5,
	}

	// If it's a full URL, parse it
	if len(redisURL) > 8 && redisURL[:8] == "redis://" {
		parsedOpts, err := redis.ParseURL(redisURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
		}
		opts = parsedOpts
		opts.DialTimeout = 5 * time.Second
		opts.ReadTimeout = 3 * time.Second
		opts.WriteTimeout = 3 * time.Second
		opts.PoolSize = 5
	}

	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return client, nil
}

// claimOneVideo atomically claims a single video from the unified queue using RPOP.
// Queue items are formatted as "channelID:videoID".
// Returns the channelID and videoID parsed from the queue item.
func claimOneVideo(ctx context.Context, client *redis.Client) (channelID, videoID string, err error) {
	queueItem, err := client.RPop(ctx, unifiedQueueKey).Result()
	if err == redis.Nil {
		// Queue is empty
		return "", "", nil
	}
	if err != nil {
		return "", "", fmt.Errorf("failed to pop from queue: %w", err)
	}

	// Parse queue item format: "channelID:videoID"
	parts := strings.SplitN(queueItem, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid queue item format: %s (expected channelID:videoID)", queueItem)
	}

	return parts[0], parts[1], nil
}

// fetchVideoInfo retrieves video metadata from Redis
func fetchVideoInfo(ctx context.Context, client *redis.Client, channelID, videoID string) (*VideoInfo, error) {
	videoKey := "video:" + channelID + ":" + videoID
	data, err := client.Get(ctx, videoKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch video info: %w", err)
	}

	var info VideoInfo
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		return nil, fmt.Errorf("failed to parse video info: %w", err)
	}

	return &info, nil
}

// UploadMetadata contains metadata for uploading a video to the collector
type UploadMetadata struct {
	VideoID       string `json:"video_id"`
	ChannelID     string `json:"channel_id"`
	ChannelName   string `json:"channel_name"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Duration      int    `json:"duration"`
	UploadDate    string `json:"upload_date"`
	EpisodeNumber int    `json:"episode_number"`
	Filename      string `json:"filename"`
	FileSize      int64  `json:"file_size"`
	Resolution    string `json:"resolution"`
	Format        string `json:"format"`
}

// uploadToCollector uploads a video file to the collector service using streaming
func uploadToCollector(ctx context.Context, collectorURL, filePath string, metadata *UploadMetadata) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	// Create a pipe to stream the multipart form data
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	// Error channel for the goroutine
	errChan := make(chan error, 1)

	// Write the multipart form in a goroutine to enable streaming
	go func() {
		defer pw.Close()
		defer writer.Close()

		// Add metadata field first
		metadataJSON, err := json.Marshal(metadata)
		if err != nil {
			errChan <- fmt.Errorf("failed to marshal metadata: %w", err)
			return
		}
		if err := writer.WriteField("metadata", string(metadataJSON)); err != nil {
			errChan <- fmt.Errorf("failed to write metadata field: %w", err)
			return
		}

		// Add file part - streams directly from file
		part, err := writer.CreateFormFile("video", filepath.Base(filePath))
		if err != nil {
			errChan <- fmt.Errorf("failed to create form file: %w", err)
			return
		}

		// Stream file content to the multipart writer
		if _, err := io.Copy(part, file); err != nil {
			errChan <- fmt.Errorf("failed to stream file: %w", err)
			return
		}

		errChan <- nil
	}()

	// Create request with the pipe reader as body (streams data)
	uploadURL := strings.TrimSuffix(collectorURL, "/") + "/upload"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadURL, pr)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// Don't set Content-Length - let it use chunked transfer encoding

	logging.Info("uploading to collector (streaming)",
		"url", uploadURL,
		"file_size", fileInfo.Size(),
		"filename", metadata.Filename,
	)

	// Send request with long timeout for large files
	client := &http.Client{
		Timeout: 60 * time.Minute, // 60 minutes for large uploads
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("upload request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors from the writer goroutine
	if writeErr := <-errChan; writeErr != nil {
		return writeErr
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// convertFormatsToStreams converts YouTube DownloadableFormats to downloader Streams
func convertFormatsToStreams(formats []youtube.DownloadableFormat) []downloader.Stream {
	streams := make([]downloader.Stream, 0, len(formats))

	for _, f := range formats {
		stream := downloader.Stream{
			FormatID:      fmt.Sprintf("%d", f.ITag),
			URL:           f.URL,
			Width:         f.Width,
			Height:        f.Height,
			Bitrate:       f.Bitrate,
			FileSize:      f.ContentLength,
			Quality:       f.Quality,
			QualityLabel:  f.QualityLabel,
			MimeType:      f.MimeType,
			ContentLength: f.ContentLength,
			FPS:           f.FPS,
		}

		// Parse mime type to determine codecs and extension
		if strings.Contains(f.MimeType, "video/mp4") {
			stream.Extension = "mp4"
			stream.VCodec = "h264"
			if strings.Contains(f.MimeType, "avc1") {
				stream.VCodec = "h264"
			}
		} else if strings.Contains(f.MimeType, "video/webm") {
			stream.Extension = "webm"
			stream.VCodec = "vp9"
		} else if strings.Contains(f.MimeType, "audio/mp4") {
			stream.Extension = "m4a"
			stream.ACodec = "aac"
		} else if strings.Contains(f.MimeType, "audio/webm") {
			stream.Extension = "webm"
			stream.ACodec = "opus"
		}

		// Determine stream type based on codecs presence
		hasVideo := f.Width > 0 || f.Height > 0 || strings.Contains(f.MimeType, "video")
		hasAudio := f.AudioQuality != "" || strings.Contains(f.MimeType, "audio")

		if hasVideo && hasAudio {
			stream.StreamType = downloader.StreamTypeCombined
		} else if hasVideo {
			stream.StreamType = downloader.StreamTypeVideo
			stream.ACodec = ""
		} else if hasAudio {
			stream.StreamType = downloader.StreamTypeAudio
			stream.VCodec = ""
		}

		streams = append(streams, stream)
	}

	return streams
}

// getEnvInt gets an integer environment variable with a default value
func getEnvInt(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultValue
}
