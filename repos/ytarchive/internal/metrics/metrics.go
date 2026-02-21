package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// YouTube API metrics
	YouTubeAPIRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ytarchive_youtube_api_requests_total",
			Help: "Total number of YouTube API requests",
		},
		[]string{"endpoint", "status"},
	)

	YouTubeAPILatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ytarchive_youtube_api_latency_seconds",
			Help:    "YouTube API request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)

	YouTubeAPIRetries = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ytarchive_youtube_api_retries_total",
			Help: "Total number of YouTube API request retries",
		},
		[]string{"endpoint"},
	)

	// Download metrics
	DownloadsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ytarchive_downloads_total",
			Help: "Total number of video downloads",
		},
		[]string{"status"},
	)

	DownloadDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ytarchive_download_duration_seconds",
			Help:    "Video download duration in seconds",
			Buckets: []float64{10, 30, 60, 120, 300, 600, 1800, 3600},
		},
		[]string{"channel_id"},
	)

	DownloadBytes = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ytarchive_download_bytes_total",
			Help: "Total bytes downloaded",
		},
		[]string{"channel_id"},
	)

	// Queue metrics
	QueueSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ytarchive_queue_size",
			Help: "Current size of the video download queue",
		},
		[]string{"channel_id"},
	)

	// Job metrics
	ActiveJobs = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "ytarchive_active_jobs",
			Help: "Number of currently active download jobs",
		},
	)

	JobsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ytarchive_jobs_total",
			Help: "Total number of jobs by status",
		},
		[]string{"status"},
	)

	// Channel metrics
	ChannelsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "ytarchive_channels_total",
			Help: "Total number of channels being archived",
		},
	)

	VideosDiscovered = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ytarchive_videos_discovered_total",
			Help: "Total number of videos discovered",
		},
		[]string{"channel_id"},
	)
)

// RecordDownloadSuccess records a successful download
func RecordDownloadSuccess(channelID string, duration float64, bytes int64) {
	DownloadsTotal.WithLabelValues("success").Inc()
	DownloadDuration.WithLabelValues(channelID).Observe(duration)
	DownloadBytes.WithLabelValues(channelID).Add(float64(bytes))
}

// RecordDownloadFailure records a failed download
func RecordDownloadFailure() {
	DownloadsTotal.WithLabelValues("failure").Inc()
}

// RecordAPIRequest records a YouTube API request
func RecordAPIRequest(endpoint, status string, latencySeconds float64) {
	YouTubeAPIRequests.WithLabelValues(endpoint, status).Inc()
	YouTubeAPILatency.WithLabelValues(endpoint).Observe(latencySeconds)
}

// RecordAPIRetry records a YouTube API retry
func RecordAPIRetry(endpoint string) {
	YouTubeAPIRetries.WithLabelValues(endpoint).Inc()
}
