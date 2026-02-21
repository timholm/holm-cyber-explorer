package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/timholm/ytarchive/internal/logging"
	"github.com/timholm/ytarchive/internal/metrics"
)

// Downloader handles downloading videos using native HTTP
type Downloader struct {
	config        *Config
	reporter      *ProgressReporter
	redisReporter *RedisProgressReporter
	workerID      string
	httpClient    *http.Client
	userAgent     string
	cookieHeader  string
	useAndroid    bool
}

// NewDownloader creates a new Downloader with the given configuration
func NewDownloader(config *Config, reporter *ProgressReporter, workerID string) *Downloader {
	return &Downloader{
		config:   config,
		reporter: reporter,
		workerID: workerID,
		httpClient: &http.Client{
			Timeout: 0, // No timeout for downloads, we handle it ourselves
			Transport: &http.Transport{
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 5,
				IdleConnTimeout:     90 * time.Second,
				DisableCompression:  true, // Don't compress video streams
			},
		},
		// Use Android user agent to match the player request
		userAgent:  "com.google.android.youtube/19.09.37 (Linux; U; Android 13) gzip",
		useAndroid: true,
	}
}

// SetRedisReporter sets the Redis progress reporter for direct Redis progress updates
func (d *Downloader) SetRedisReporter(reporter *RedisProgressReporter) {
	d.redisReporter = reporter
}

// SetCookies sets the cookie header for authenticated downloads
func (d *Downloader) SetCookies(cookieHeader string) {
	d.cookieHeader = cookieHeader
}

// DownloadResult contains the result of a download operation
type DownloadResult struct {
	VideoID   string
	ChannelID string
	Success   bool
	FilePath  string
	FileSize  int64
	Error     error
	Attempts  int
	Duration  time.Duration
}

// DownloadRequest contains all information needed to download a video
type DownloadRequest struct {
	VideoID              string
	ChannelID            string
	Title                string
	Description          string
	Duration             int
	UploadDate           string
	ThumbnailURL         string
	ViewCount            int64
	ChannelName          string
	EpisodeNumber        int // Episode number based on upload date order (oldest = 1)
	Streams              []Stream
	AvailableResolutions []ResolutionOption
}

// Download downloads a video with retry logic and progress reporting
func (d *Downloader) Download(ctx context.Context, req *DownloadRequest) *DownloadResult {
	startTime := time.Now()
	result := &DownloadResult{
		VideoID:   req.VideoID,
		ChannelID: req.ChannelID,
	}

	var lastErr error
	for attempt := 1; attempt <= d.config.Retries; attempt++ {
		result.Attempts = attempt

		logging.Info("starting download attempt",
			"worker_id", d.workerID,
			"video_id", req.VideoID,
			"channel_id", req.ChannelID,
			"attempt", attempt,
			"max_attempts", d.config.Retries,
		)

		err := d.downloadVideo(ctx, req)
		if err == nil {
			result.Success = true
			result.Duration = time.Since(startTime)

			// Find the downloaded file and get its size
			filePath, fileSize := d.findDownloadedFile(req.ChannelID, req.VideoID)
			result.FilePath = filePath
			result.FileSize = fileSize

			logging.Info("successfully downloaded video",
				"worker_id", d.workerID,
				"video_id", req.VideoID,
				"duration", result.Duration.String(),
				"file_path", filePath,
				"file_size", fileSize,
			)

			// Record success metrics
			metrics.RecordDownloadSuccess(req.ChannelID, result.Duration.Seconds(), fileSize)

			return result
		}

		lastErr = err
		logging.Warn("download attempt failed",
			"worker_id", d.workerID,
			"video_id", req.VideoID,
			"attempt", attempt,
			"error", err,
		)

		// Check if context is cancelled
		if ctx.Err() != nil {
			result.Error = ctx.Err()
			return result
		}

		// Wait before retrying (exponential backoff)
		if attempt < d.config.Retries {
			delay := d.getRetryDelay(attempt)
			logging.Debug("waiting before retry",
				"worker_id", d.workerID,
				"video_id", req.VideoID,
				"delay", delay.String(),
			)

			select {
			case <-ctx.Done():
				result.Error = ctx.Err()
				return result
			case <-time.After(delay):
				// Continue to next attempt
			}
		}
	}

	result.Error = fmt.Errorf("all %d download attempts failed: %w", d.config.Retries, lastErr)
	result.Duration = time.Since(startTime)

	// Record failure metrics
	metrics.RecordDownloadFailure()

	return result
}

// DownloadFromURL downloads a video directly from a stream URL (legacy compatibility)
func (d *Downloader) DownloadFromURL(ctx context.Context, videoID string, streamURL string) *DownloadResult {
	req := &DownloadRequest{
		VideoID: videoID,
		Streams: []Stream{
			{
				URL:        streamURL,
				StreamType: StreamTypeCombined,
				Extension:  "mp4",
			},
		},
	}
	return d.Download(ctx, req)
}

// downloadVideo performs the actual download operation
func (d *Downloader) downloadVideo(ctx context.Context, req *DownloadRequest) error {
	// Create video directory
	videoDir := d.getVideoDir(req.ChannelID, req.VideoID)
	if err := os.MkdirAll(videoDir, 0755); err != nil {
		return fmt.Errorf("failed to create video directory: %w", err)
	}

	// Select best stream
	selector := NewStreamSelector(d.config.MaxHeight, d.config.PreferCombinedStream)
	videoStream, audioStream, err := selector.SelectBestStream(req.Streams)
	if err != nil {
		return fmt.Errorf("failed to select stream: %w", err)
	}

	logging.Info("selected stream",
		"video_id", req.VideoID,
		"format_id", videoStream.FormatID,
		"quality", QualityLabel(videoStream.Height),
		"type", videoStream.StreamType,
		"has_audio_stream", audioStream != nil,
	)

	// Generate the final filename using channel name, episode number, and title
	// Format: {channel}-ep{number}-{title}.{ext}
	var finalFilename string
	if req.ChannelName != "" && req.EpisodeNumber > 0 && req.Title != "" {
		finalFilename = GenerateVideoFilename(req.ChannelName, req.EpisodeNumber, req.Title, "mp4")
	} else {
		// Fallback to video.mp4 if we don't have complete info
		finalFilename = "video.mp4"
	}

	// Report progress: starting
	d.reportProgress(req.VideoID, "downloading", 0, 0, 0, "", "")

	// Download based on stream type
	var videoPath string
	if videoStream.IsSegmented {
		videoPath, err = d.downloadSegmentedStream(ctx, req.VideoID, videoStream, videoDir)
	} else {
		videoPath, err = d.downloadSingleStream(ctx, req.VideoID, videoStream, videoDir)
	}
	if err != nil {
		return fmt.Errorf("failed to download video stream: %w", err)
	}

	// If we have a separate audio stream, download and merge it
	if audioStream != nil && videoStream.StreamType == StreamTypeVideo {
		audioPath := filepath.Join(videoDir, "audio.m4a")
		if audioStream.IsSegmented {
			audioPath, err = d.downloadSegmentedStream(ctx, req.VideoID, audioStream, videoDir)
		} else {
			_, err = d.downloadSingleStream(ctx, req.VideoID, audioStream, videoDir)
		}
		if err != nil {
			logging.Warn("failed to download audio stream, continuing with video only",
				"video_id", req.VideoID,
				"error", err,
			)
		} else {
			// Merge audio and video if ffmpeg is available
			if MergerAvailable() {
				d.reportProgress(req.VideoID, "processing", 95, 0, 0, "", "")
				merger := NewMerger()
				outputPath := filepath.Join(videoDir, finalFilename)

				// If output would be the same as input video, use a temp file
				// FFmpeg cannot read and write to the same file
				useTempOutput := outputPath == videoPath
				if useTempOutput {
					outputPath = filepath.Join(videoDir, "merged_temp.mp4")
				}

				result := merger.MergeWithCodecCopy(ctx, videoPath, audioPath, outputPath)
				if result.Error != nil {
					logging.Warn("failed to merge streams, keeping video only",
						"video_id", req.VideoID,
						"error", result.Error,
					)
					// Clean up temp file if we created one
					if useTempOutput {
						os.Remove(outputPath)
					}
				} else {
					// Clean up original video and audio files
					os.Remove(videoPath)
					os.Remove(audioPath)

					// If we used a temp file, rename it to the final name
					if useTempOutput {
						finalPath := filepath.Join(videoDir, finalFilename)
						if err := os.Rename(outputPath, finalPath); err != nil {
							logging.Warn("failed to rename merged file",
								"video_id", req.VideoID,
								"error", err,
							)
							videoPath = outputPath
						} else {
							videoPath = finalPath
						}
					} else {
						videoPath = outputPath
					}
				}
			}
		}
	} else if finalFilename != "video.mp4" && !videoStream.IsSegmented {
		// If no merging needed, rename the video file to the final filename
		newPath := filepath.Join(videoDir, finalFilename)
		if err := os.Rename(videoPath, newPath); err == nil {
			videoPath = newPath
		}
	}

	// Download thumbnail
	if d.config.WriteThumbnail && req.ThumbnailURL != "" {
		thumbPath := filepath.Join(videoDir, "thumbnail.jpg")
		if err := d.downloadFile(ctx, req.ThumbnailURL, thumbPath, nil); err != nil {
			logging.Warn("failed to download thumbnail",
				"video_id", req.VideoID,
				"error", err,
			)
		}
	}

	// Write metadata JSON
	if d.config.WriteInfoJSON {
		metadataPath := filepath.Join(videoDir, "metadata.json")
		if err := d.writeMetadata(req, metadataPath, videoPath, videoStream); err != nil {
			logging.Warn("failed to write metadata",
				"video_id", req.VideoID,
				"error", err,
			)
		}
	}

	// Report progress: completed
	d.reportProgress(req.VideoID, "completed", 100, 0, 0, "", "")

	return nil
}

// downloadSingleStream downloads a non-segmented stream
func (d *Downloader) downloadSingleStream(ctx context.Context, videoID string, stream *Stream, videoDir string) (string, error) {
	filename := fmt.Sprintf("video.%s", stream.Extension)
	if stream.StreamType == StreamTypeAudio {
		filename = fmt.Sprintf("audio.%s", stream.Extension)
	}
	outputPath := filepath.Join(videoDir, filename)
	partPath := outputPath + ".part"

	// Check if we have a partial download
	var startByte int64 = 0
	if info, err := os.Stat(partPath); err == nil {
		startByte = info.Size()
		logging.Info("resuming download",
			"video_id", videoID,
			"start_byte", startByte,
		)
	}

	// Get total size for progress reporting
	totalSize := stream.ContentLength
	if totalSize == 0 {
		fetcher := NewStreamFetcher()
		if size, err := fetcher.GetStreamContentLength(stream.URL); err == nil {
			totalSize = size
		}
	}

	// Progress callback with speed calculation
	var lastReportTime time.Time
	var lastBytes int64
	var speedSamples []float64
	const maxSpeedSamples = 5

	progressCallback := func(downloaded, total int64) {
		now := time.Now()
		elapsed := now.Sub(lastReportTime).Seconds()

		// Throttle progress reports to once per second
		if elapsed < 1.0 {
			return
		}

		// Calculate speed
		var speed float64
		if elapsed > 0 && lastBytes > 0 {
			speed = float64(downloaded-lastBytes) / elapsed
			speedSamples = append(speedSamples, speed)
			if len(speedSamples) > maxSpeedSamples {
				speedSamples = speedSamples[1:]
			}
		}

		// Calculate average speed for smoother display
		avgSpeed := float64(0)
		if len(speedSamples) > 0 {
			for _, s := range speedSamples {
				avgSpeed += s
			}
			avgSpeed /= float64(len(speedSamples))
		}

		// Calculate ETA
		eta := ""
		if avgSpeed > 0 && total > 0 {
			remaining := total - downloaded
			etaSeconds := float64(remaining) / avgSpeed
			if etaSeconds > 0 {
				hours := int(etaSeconds / 3600)
				minutes := int((etaSeconds - float64(hours)*3600) / 60)
				secs := int(etaSeconds) % 60
				if hours > 0 {
					eta = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
				} else {
					eta = fmt.Sprintf("%02d:%02d", minutes, secs)
				}
			}
		}

		lastReportTime = now
		lastBytes = downloaded

		percentage := float64(0)
		if total > 0 {
			percentage = float64(downloaded) / float64(total) * 100
		}
		d.reportProgress(videoID, "downloading", percentage, downloaded, total, FormatSpeed(avgSpeed), eta)
	}

	// Download with resume support
	err := d.downloadFileWithResume(ctx, stream.URL, partPath, startByte, totalSize, progressCallback)
	if err != nil {
		return "", err
	}

	// Rename .part file to final name
	if err := os.Rename(partPath, outputPath); err != nil {
		return "", fmt.Errorf("failed to rename part file: %w", err)
	}

	return outputPath, nil
}

// downloadSegmentedStream downloads a segmented (DASH/HLS) stream
func (d *Downloader) downloadSegmentedStream(ctx context.Context, videoID string, stream *Stream, videoDir string) (string, error) {
	if len(stream.SegmentURLs) == 0 {
		return "", fmt.Errorf("no segment URLs provided")
	}

	filename := fmt.Sprintf("video.%s", stream.Extension)
	if stream.StreamType == StreamTypeAudio {
		filename = fmt.Sprintf("audio.%s", stream.Extension)
	}
	outputPath := filepath.Join(videoDir, filename)

	// Create a temporary directory for segments
	segmentDir := filepath.Join(videoDir, ".segments")
	if err := os.MkdirAll(segmentDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create segment directory: %w", err)
	}
	defer os.RemoveAll(segmentDir)

	// Download init segment if present
	var segmentPaths []string
	if stream.InitURL != "" {
		initPath := filepath.Join(segmentDir, "init.mp4")
		if err := d.downloadFile(ctx, stream.InitURL, initPath, nil); err != nil {
			return "", fmt.Errorf("failed to download init segment: %w", err)
		}
		segmentPaths = append(segmentPaths, initPath)
	}

	// Download segments with progress
	totalSegments := len(stream.SegmentURLs)
	for i, segURL := range stream.SegmentURLs {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		segPath := filepath.Join(segmentDir, fmt.Sprintf("segment_%05d.ts", i))
		if err := d.downloadFile(ctx, segURL, segPath, nil); err != nil {
			return "", fmt.Errorf("failed to download segment %d: %w", i, err)
		}
		segmentPaths = append(segmentPaths, segPath)

		// Report progress with fragment info
		percentage := float64(i+1) / float64(totalSegments) * 100
		d.reportProgressWithFragment(videoID, "downloading", percentage, int64(i+1), int64(totalSegments),
			"", "", fmt.Sprintf("%d/%d", i+1, totalSegments))
	}

	// Concatenate segments
	if MergerAvailable() {
		merger := NewMerger()
		result := merger.ConcatSegments(ctx, segmentPaths, outputPath)
		if result.Error != nil {
			return "", fmt.Errorf("failed to concatenate segments: %w", result.Error)
		}
	} else {
		// Fallback: simple binary concatenation (may not work for all formats)
		if err := d.concatenateFiles(segmentPaths, outputPath); err != nil {
			return "", fmt.Errorf("failed to concatenate segments: %w", err)
		}
	}

	return outputPath, nil
}

// downloadFileWithResume downloads a file with resume support using Range headers
func (d *Downloader) downloadFileWithResume(ctx context.Context, url, outputPath string, startByte, totalSize int64, progressCallback func(downloaded, total int64)) error {
	// Open file for writing (append if resuming)
	flags := os.O_CREATE | os.O_WRONLY
	if startByte > 0 {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}
	file, err := os.OpenFile(outputPath, flags, 0644)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer file.Close()

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to match the YouTube player request (Android client)
	req.Header.Set("User-Agent", d.userAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// For Android client, use different headers
	if d.useAndroid {
		req.Header.Set("X-Goog-Api-Format-Version", "2")
	} else {
		req.Header.Set("Origin", "https://www.youtube.com")
		req.Header.Set("Referer", "https://www.youtube.com/")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Site", "cross-site")
	}

	// Add cookies for authenticated downloads
	if d.cookieHeader != "" {
		req.Header.Set("Cookie", d.cookieHeader)
	}

	// Add Range header for resume
	if startByte > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", startByte))
	}

	// Execute request
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Update total size if we didn't have it
	if totalSize == 0 {
		if resp.StatusCode == http.StatusPartialContent {
			// Parse Content-Range header
			totalSize = resp.ContentLength + startByte
		} else {
			totalSize = resp.ContentLength
		}
	}

	// Create progress reader
	downloaded := startByte
	reader := &progressReader{
		reader: resp.Body,
		onProgress: func(n int64) {
			downloaded += n
			if progressCallback != nil {
				progressCallback(downloaded, totalSize)
			}
		},
	}

	// Copy data to file
	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// downloadFile downloads a file without resume support
func (d *Downloader) downloadFile(ctx context.Context, url, outputPath string, progressCallback func(downloaded, total int64)) error {
	return d.downloadFileWithResume(ctx, url, outputPath, 0, 0, progressCallback)
}

// concatenateFiles concatenates multiple files into one (simple binary concat)
func (d *Downloader) concatenateFiles(inputPaths []string, outputPath string) error {
	output, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer output.Close()

	for _, inputPath := range inputPaths {
		input, err := os.Open(inputPath)
		if err != nil {
			return err
		}
		_, err = io.Copy(output, input)
		input.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// writeMetadata writes video metadata to a JSON file
func (d *Downloader) writeMetadata(req *DownloadRequest, path string, videoPath string, selectedStream *Stream) error {
	var fileSize int64
	if info, err := os.Stat(videoPath); err == nil {
		fileSize = info.Size()
	}

	// Get available resolutions from streams if not already populated
	availableResolutions := req.AvailableResolutions
	if len(availableResolutions) == 0 && len(req.Streams) > 0 {
		availableResolutions = GetAvailableResolutions(req.Streams)
	}

	// Build selected resolution info
	var selectedResolution map[string]interface{}
	if selectedStream != nil {
		selectedResolution = map[string]interface{}{
			"height":      selectedStream.Height,
			"width":       selectedStream.Width,
			"label":       QualityLabel(selectedStream.Height),
			"format_id":   selectedStream.FormatID,
			"stream_type": string(selectedStream.StreamType),
			"codec":       selectedStream.VCodec,
			"bitrate":     selectedStream.Bitrate,
		}
	}

	metadata := map[string]interface{}{
		"id":                    req.VideoID,
		"title":                 req.Title,
		"description":           req.Description,
		"duration":              req.Duration,
		"upload_date":           req.UploadDate,
		"channel_id":            req.ChannelID,
		"channel_name":          req.ChannelName,
		"view_count":            req.ViewCount,
		"thumbnail_url":         req.ThumbnailURL,
		"file_size":             fileSize,
		"downloaded_at":         time.Now().Format(time.RFC3339),
		"available_resolutions": availableResolutions,
		"selected_resolution":   selectedResolution,
	}

	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// getVideoDir returns the directory path for a video
func (d *Downloader) getVideoDir(channelID, videoID string) string {
	if channelID != "" {
		return filepath.Join(d.config.OutputPath, "channels", channelID, "videos", videoID)
	}
	return filepath.Join(d.config.OutputPath, videoID)
}

// getRetryDelay returns the delay before the next retry attempt
func (d *Downloader) getRetryDelay(attempt int) time.Duration {
	if attempt-1 < len(d.config.RetryDelays) {
		return time.Duration(d.config.RetryDelays[attempt-1]) * time.Second
	}
	// Default to last delay if we've exhausted the configured delays
	if len(d.config.RetryDelays) > 0 {
		return time.Duration(d.config.RetryDelays[len(d.config.RetryDelays)-1]) * time.Second
	}
	return 30 * time.Second
}

// findDownloadedFile locates the downloaded video file and returns its path and size
func (d *Downloader) findDownloadedFile(channelID, videoID string) (string, int64) {
	videoDir := d.getVideoDir(channelID, videoID)

	// Look for the main video file - check for new naming format first (*.mp4)
	// then fall back to legacy naming (video.mp4, video.mkv, etc.)
	patterns := []string{"*-ep*-*.mp4", "*-ep*-*.mkv", "*-ep*-*.webm", "video.mp4", "video.mkv", "video.webm", "*.mp4", "*.mkv", "*.webm"}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(videoDir, pattern))
		if err != nil {
			continue
		}

		for _, match := range matches {
			// Skip partial files
			if filepath.Ext(match) == ".part" {
				continue
			}

			// Skip audio files
			base := filepath.Base(match)
			if len(base) >= 5 && base[:5] == "audio" {
				continue
			}

			info, err := os.Stat(match)
			if err != nil {
				continue
			}

			return match, info.Size()
		}
	}

	return "", 0
}

// reportProgress sends a progress update to both HTTP and Redis reporters
func (d *Downloader) reportProgress(videoID, status string, percentage float64, downloaded, total int64, speed, eta string) {
	d.reportProgressWithFragment(videoID, status, percentage, downloaded, total, speed, eta, "")
}

// reportProgressWithFragment sends a progress update with fragment info (for segmented downloads)
func (d *Downloader) reportProgressWithFragment(videoID, status string, percentage float64, downloaded, total int64, speed, eta, fragment string) {
	progress := &DownloadProgress{
		VideoID:         videoID,
		WorkerID:        d.workerID,
		Status:          status,
		Percentage:      percentage,
		DownloadedBytes: downloaded,
		TotalBytes:      total,
		Speed:           speed,
		ETA:             eta,
		Fragment:        fragment,
		UpdatedAt:       time.Now().Unix(),
	}

	// Report to HTTP endpoint (controller API)
	if d.reporter != nil {
		if err := d.reporter.ReportProgress(progress); err != nil {
			logging.Warn("failed to report progress to HTTP",
				"worker_id", d.workerID,
				"video_id", videoID,
				"error", err,
			)
		}
	}

	// Report to Redis for UI polling
	if d.redisReporter != nil {
		// Force update for status changes (completed, error, processing)
		force := status == "completed" || status == "error" || status == "processing"
		if err := d.redisReporter.ReportProgress(progress, force); err != nil {
			logging.Warn("failed to report progress to Redis",
				"worker_id", d.workerID,
				"video_id", videoID,
				"error", err,
			)
		}
	}
}

// Cleanup removes partial downloads for a video
func (d *Downloader) Cleanup(channelID, videoID string) error {
	videoDir := d.getVideoDir(channelID, videoID)

	// Check if directory exists
	if _, err := os.Stat(videoDir); os.IsNotExist(err) {
		return nil
	}

	// Look for partial files
	partFiles, err := filepath.Glob(filepath.Join(videoDir, "*.part"))
	if err != nil {
		return fmt.Errorf("failed to find partial files: %w", err)
	}

	for _, partFile := range partFiles {
		if err := os.Remove(partFile); err != nil {
			logging.Warn("failed to remove partial file",
				"worker_id", d.workerID,
				"video_id", videoID,
				"file", partFile,
				"error", err,
			)
		}
	}

	// Remove segment directory if it exists
	segmentDir := filepath.Join(videoDir, ".segments")
	os.RemoveAll(segmentDir)

	return nil
}

// progressReader wraps an io.Reader to track read progress
type progressReader struct {
	reader     io.Reader
	onProgress func(n int64)
	mu         sync.Mutex
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	if n > 0 {
		pr.mu.Lock()
		pr.onProgress(int64(n))
		pr.mu.Unlock()
	}
	return n, err
}

// VideoURL returns the YouTube video URL for a given video ID
func VideoURL(videoID string) string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
}

// AudioDownloadResult contains the result of an audio-only download
type AudioDownloadResult struct {
	VideoID   string
	ChannelID string
	Success   bool
	FilePath  string
	FileSize  int64
	Error     error
}

// DownloadAudioOnly downloads just the audio stream from a video
func (d *Downloader) DownloadAudioOnly(ctx context.Context, req *DownloadRequest) *AudioDownloadResult {
	result := &AudioDownloadResult{
		VideoID:   req.VideoID,
		ChannelID: req.ChannelID,
	}

	// Create video directory
	videoDir := d.getVideoDir(req.ChannelID, req.VideoID)
	if err := os.MkdirAll(videoDir, 0755); err != nil {
		result.Error = fmt.Errorf("failed to create video directory: %w", err)
		return result
	}

	// Find best audio stream
	selector := NewStreamSelector(0, false)
	audioStream := selector.findBestAudio(req.Streams)
	if audioStream == nil {
		result.Error = fmt.Errorf("no audio stream found")
		return result
	}

	logging.Info("downloading audio-only",
		"video_id", req.VideoID,
		"format_id", audioStream.FormatID,
		"bitrate", audioStream.Bitrate,
	)

	// Download audio stream
	var audioPath string
	var err error
	if audioStream.IsSegmented {
		audioPath, err = d.downloadSegmentedStream(ctx, req.VideoID, audioStream, videoDir)
	} else {
		audioPath, err = d.downloadAudioStream(ctx, req.VideoID, audioStream, videoDir)
	}

	if err != nil {
		result.Error = fmt.Errorf("failed to download audio: %w", err)
		return result
	}

	// Get file info
	if info, err := os.Stat(audioPath); err == nil {
		result.FileSize = info.Size()
	}

	result.Success = true
	result.FilePath = audioPath
	return result
}

// downloadAudioStream downloads an audio-only stream to audio.m4a
func (d *Downloader) downloadAudioStream(ctx context.Context, videoID string, stream *Stream, videoDir string) (string, error) {
	ext := stream.Extension
	if ext == "" {
		ext = "m4a"
	}
	outputPath := filepath.Join(videoDir, fmt.Sprintf("audio.%s", ext))
	partPath := outputPath + ".part"

	// Check if we have a partial download
	var startByte int64 = 0
	if info, err := os.Stat(partPath); err == nil {
		startByte = info.Size()
	}

	// Get total size for progress reporting
	totalSize := stream.ContentLength
	if totalSize == 0 {
		fetcher := NewStreamFetcher()
		if size, err := fetcher.GetStreamContentLength(stream.URL); err == nil {
			totalSize = size
		}
	}

	// Download with resume support
	err := d.downloadFileWithResume(ctx, stream.URL, partPath, startByte, totalSize, nil)
	if err != nil {
		return "", err
	}

	// Rename .part file to final name
	if err := os.Rename(partPath, outputPath); err != nil {
		return "", fmt.Errorf("failed to rename part file: %w", err)
	}

	return outputPath, nil
}

// DownloadSubtitles downloads subtitles for a video and saves them as VTT files
func (d *Downloader) DownloadSubtitles(ctx context.Context, channelID, videoID string, captions map[string]string) error {
	videoDir := d.getVideoDir(channelID, videoID)
	if err := os.MkdirAll(videoDir, 0755); err != nil {
		return fmt.Errorf("failed to create video directory: %w", err)
	}

	for langCode, captionURL := range captions {
		// Add format parameter to get WebVTT format
		url := captionURL
		if !stringContains(url, "fmt=") {
			if stringContains(url, "?") {
				url += "&fmt=vtt"
			} else {
				url += "?fmt=vtt"
			}
		}

		// Determine filename - use "subtitles.vtt" for primary language (usually en)
		// and "subtitles.{lang}.vtt" for additional languages
		filename := "subtitles.vtt"
		if langCode != "en" && langCode != "" {
			filename = fmt.Sprintf("subtitles.%s.vtt", langCode)
		}
		outputPath := filepath.Join(videoDir, filename)

		logging.Info("downloading subtitles",
			"video_id", videoID,
			"language", langCode,
			"output", filename,
		)

		err := d.downloadFile(ctx, url, outputPath, nil)
		if err != nil {
			logging.Warn("failed to download subtitles",
				"video_id", videoID,
				"language", langCode,
				"error", err,
			)
			continue
		}
	}

	return nil
}

// ExtractAudioFromVideo extracts audio from an existing video file using ffmpeg
func (d *Downloader) ExtractAudioFromVideo(ctx context.Context, channelID, videoID string) (*AudioDownloadResult, error) {
	result := &AudioDownloadResult{
		VideoID:   videoID,
		ChannelID: channelID,
	}

	videoDir := d.getVideoDir(channelID, videoID)

	// Find the video file
	var videoPath string
	patterns := []string{"video.mp4", "video.mkv", "video.webm"}
	for _, pattern := range patterns {
		path := filepath.Join(videoDir, pattern)
		if _, err := os.Stat(path); err == nil {
			videoPath = path
			break
		}
	}

	if videoPath == "" {
		result.Error = fmt.Errorf("video file not found")
		return result, result.Error
	}

	// Check if ffmpeg is available
	if !MergerAvailable() {
		result.Error = fmt.Errorf("ffmpeg not available for audio extraction")
		return result, result.Error
	}

	audioPath := filepath.Join(videoDir, "audio.m4a")

	// Use ffmpeg to extract audio
	merger := NewMerger()
	mergeResult := merger.ExtractAudio(ctx, videoPath, audioPath)
	if mergeResult.Error != nil {
		result.Error = mergeResult.Error
		return result, result.Error
	}

	// Get file info
	if info, err := os.Stat(audioPath); err == nil {
		result.FileSize = info.Size()
	}

	result.Success = true
	result.FilePath = audioPath
	return result, nil
}

// stringContains checks if a string contains a substring
func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
