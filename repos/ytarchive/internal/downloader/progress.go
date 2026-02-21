package downloader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// DownloadProgress represents the current download progress
type DownloadProgress struct {
	VideoID         string  `json:"video_id"`
	WorkerID        string  `json:"worker_id"`
	Status          string  `json:"status"` // downloading, processing, completed, error
	Percentage      float64 `json:"percentage"`
	DownloadedBytes int64   `json:"downloaded_bytes"`
	TotalBytes      int64   `json:"total_bytes"`
	Speed           string  `json:"speed"`
	ETA             string  `json:"eta"`
	Fragment        string  `json:"fragment,omitempty"` // e.g., "5/10" for segmented downloads
	UpdatedAt       int64   `json:"updated_at"`
}

// ProgressCallback is called when download progress is updated
type ProgressCallback func(progress *DownloadProgress)

// ProgressReporter handles reporting progress to the controller API
type ProgressReporter struct {
	controllerURL string
	workerID      string
	httpClient    *http.Client
}

// NewProgressReporter creates a new ProgressReporter
func NewProgressReporter(controllerURL, workerID string) *ProgressReporter {
	return &ProgressReporter{
		controllerURL: strings.TrimSuffix(controllerURL, "/"),
		workerID:      workerID,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ReportProgress sends progress to the controller API
func (r *ProgressReporter) ReportProgress(progress *DownloadProgress) error {
	progress.WorkerID = r.workerID
	progress.UpdatedAt = time.Now().Unix()

	// Construct the API request payload with all required fields
	payload := map[string]interface{}{
		"video_id":   progress.VideoID,
		"worker_id":  progress.WorkerID,
		"status":     progress.Status,
		"percent":    progress.Percentage,
		"downloaded": progress.DownloadedBytes,
		"total":      progress.TotalBytes,
		"speed":      progress.Speed,
		"eta":        progress.ETA,
		"fragment":   progress.Fragment,
		"updated_at": progress.UpdatedAt,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal progress: %w", err)
	}

	url := fmt.Sprintf("%s/api/progress/%s", r.controllerURL, progress.VideoID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send progress: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("progress report failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// VideoStatus represents the status update for a video
type VideoStatus struct {
	VideoID   string `json:"video_id"`
	WorkerID  string `json:"worker_id"`
	Status    string `json:"status"` // pending, downloading, downloaded, error
	Error     string `json:"error,omitempty"`
	FilePath  string `json:"file_path,omitempty"`
	FileSize  int64  `json:"file_size,omitempty"`
	UpdatedAt int64  `json:"updated_at"`
}

// ReportVideoStatus sends video status update to the controller API
func (r *ProgressReporter) ReportVideoStatus(status *VideoStatus) error {
	status.WorkerID = r.workerID
	status.UpdatedAt = time.Now().Unix()

	data, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal status: %w", err)
	}

	url := fmt.Sprintf("%s/api/videos/%s/status", r.controllerURL, status.VideoID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status report failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ProgressTracker tracks download progress with speed calculation
type ProgressTracker struct {
	videoID      string
	workerID     string
	totalBytes   int64
	startTime    time.Time
	lastUpdate   time.Time
	lastBytes    int64
	downloaded   int64
	callback     ProgressCallback
	mu           sync.Mutex
	speedSamples []float64 // Rolling window of speed samples
	maxSamples   int
}

// NewProgressTracker creates a new ProgressTracker
func NewProgressTracker(videoID, workerID string, totalBytes int64, callback ProgressCallback) *ProgressTracker {
	return &ProgressTracker{
		videoID:      videoID,
		workerID:     workerID,
		totalBytes:   totalBytes,
		startTime:    time.Now(),
		lastUpdate:   time.Now(),
		callback:     callback,
		speedSamples: make([]float64, 0, 10),
		maxSamples:   10,
	}
}

// Update updates the progress with new downloaded bytes
func (pt *ProgressTracker) Update(downloadedBytes int64) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	pt.downloaded = downloadedBytes

	now := time.Now()
	elapsed := now.Sub(pt.lastUpdate).Seconds()

	// Only update if at least 500ms has passed
	if elapsed < 0.5 {
		return
	}

	// Calculate speed
	bytesDelta := downloadedBytes - pt.lastBytes
	speed := float64(bytesDelta) / elapsed

	// Add to rolling window
	pt.speedSamples = append(pt.speedSamples, speed)
	if len(pt.speedSamples) > pt.maxSamples {
		pt.speedSamples = pt.speedSamples[1:]
	}

	// Calculate average speed
	avgSpeed := pt.averageSpeed()

	pt.lastUpdate = now
	pt.lastBytes = downloadedBytes

	if pt.callback != nil {
		progress := &DownloadProgress{
			VideoID:         pt.videoID,
			WorkerID:        pt.workerID,
			Status:          "downloading",
			DownloadedBytes: downloadedBytes,
			TotalBytes:      pt.totalBytes,
			Speed:           FormatSpeed(avgSpeed),
			UpdatedAt:       now.Unix(),
		}

		if pt.totalBytes > 0 {
			progress.Percentage = float64(downloadedBytes) / float64(pt.totalBytes) * 100
			progress.ETA = pt.calculateETA(avgSpeed)
		}

		pt.callback(progress)
	}
}

// SetTotal sets the total bytes (useful when total is discovered during download)
func (pt *ProgressTracker) SetTotal(totalBytes int64) {
	pt.mu.Lock()
	pt.totalBytes = totalBytes
	pt.mu.Unlock()
}

// averageSpeed returns the average speed from the rolling window
func (pt *ProgressTracker) averageSpeed() float64 {
	if len(pt.speedSamples) == 0 {
		return 0
	}
	var sum float64
	for _, s := range pt.speedSamples {
		sum += s
	}
	return sum / float64(len(pt.speedSamples))
}

// calculateETA calculates estimated time of arrival
func (pt *ProgressTracker) calculateETA(speed float64) string {
	if speed <= 0 || pt.totalBytes <= 0 {
		return "calculating..."
	}

	remaining := pt.totalBytes - pt.downloaded
	seconds := float64(remaining) / speed

	if seconds < 0 {
		return "00:00"
	}

	hours := int(seconds / 3600)
	minutes := int((seconds - float64(hours)*3600) / 60)
	secs := int(seconds) % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
	}
	return fmt.Sprintf("%02d:%02d", minutes, secs)
}

// Complete marks the progress as complete
func (pt *ProgressTracker) Complete() {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	if pt.callback != nil {
		pt.callback(&DownloadProgress{
			VideoID:         pt.videoID,
			WorkerID:        pt.workerID,
			Status:          "completed",
			Percentage:      100,
			DownloadedBytes: pt.totalBytes,
			TotalBytes:      pt.totalBytes,
			UpdatedAt:       time.Now().Unix(),
		})
	}
}

// Error marks the progress as errored
func (pt *ProgressTracker) Error(err error) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	if pt.callback != nil {
		pt.callback(&DownloadProgress{
			VideoID:         pt.videoID,
			WorkerID:        pt.workerID,
			Status:          "error",
			DownloadedBytes: pt.downloaded,
			TotalBytes:      pt.totalBytes,
			UpdatedAt:       time.Now().Unix(),
		})
	}
}

// FormatSpeed formats bytes per second to a human-readable string
func FormatSpeed(bytesPerSecond float64) string {
	if bytesPerSecond < 1024 {
		return fmt.Sprintf("%.0f B/s", bytesPerSecond)
	}
	if bytesPerSecond < 1024*1024 {
		return fmt.Sprintf("%.1f KiB/s", bytesPerSecond/1024)
	}
	if bytesPerSecond < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MiB/s", bytesPerSecond/(1024*1024))
	}
	return fmt.Sprintf("%.1f GiB/s", bytesPerSecond/(1024*1024*1024))
}

// SegmentProgressTracker tracks progress for segmented downloads
type SegmentProgressTracker struct {
	videoID       string
	workerID      string
	totalSegments int
	completed     int
	callback      ProgressCallback
	mu            sync.Mutex
}

// NewSegmentProgressTracker creates a new SegmentProgressTracker
func NewSegmentProgressTracker(videoID, workerID string, totalSegments int, callback ProgressCallback) *SegmentProgressTracker {
	return &SegmentProgressTracker{
		videoID:       videoID,
		workerID:      workerID,
		totalSegments: totalSegments,
		callback:      callback,
	}
}

// SegmentComplete marks a segment as complete
func (spt *SegmentProgressTracker) SegmentComplete() {
	spt.mu.Lock()
	defer spt.mu.Unlock()

	spt.completed++

	if spt.callback != nil {
		percentage := float64(spt.completed) / float64(spt.totalSegments) * 100
		spt.callback(&DownloadProgress{
			VideoID:    spt.videoID,
			WorkerID:   spt.workerID,
			Status:     "downloading",
			Percentage: percentage,
			Fragment:   fmt.Sprintf("%d/%d", spt.completed, spt.totalSegments),
			UpdatedAt:  time.Now().Unix(),
		})
	}
}

// Complete marks all segments as complete
func (spt *SegmentProgressTracker) Complete() {
	spt.mu.Lock()
	defer spt.mu.Unlock()

	if spt.callback != nil {
		spt.callback(&DownloadProgress{
			VideoID:    spt.videoID,
			WorkerID:   spt.workerID,
			Status:     "completed",
			Percentage: 100,
			Fragment:   fmt.Sprintf("%d/%d", spt.totalSegments, spt.totalSegments),
			UpdatedAt:  time.Now().Unix(),
		})
	}
}
