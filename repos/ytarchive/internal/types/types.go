package types

import "time"

type Channel struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	AvatarURL   string    `json:"avatar_url"`
	BannerURL   string    `json:"banner_url"`
	VideoCount  int       `json:"video_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Video struct {
	ID                  string     `json:"id"`
	ChannelID           string     `json:"channel_id"`
	Title               string     `json:"title"`
	Description         string     `json:"description"`
	Duration            int        `json:"duration"`
	UploadDate          string     `json:"upload_date"`
	ThumbnailURL        string     `json:"thumbnail_url"`
	ViewCount           int64      `json:"view_count"`
	Status              string     `json:"status"`
	FilePath            string     `json:"file_path,omitempty"`
	FileSize            int64      `json:"file_size,omitempty"`
	Checksum            string     `json:"checksum,omitempty"`
	RetryCount          int        `json:"retry_count"`
	LastError           string     `json:"last_error,omitempty"`
	DownloadStartedAt   *time.Time `json:"download_started_at,omitempty"`
	DownloadCompletedAt *time.Time `json:"download_completed_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type Job struct {
	ID         string    `json:"id"`
	ChannelID  string    `json:"channel_id"`
	WorkerID   string    `json:"worker_id"`
	Status     string    `json:"status"`
	Progress   float64   `json:"progress"`
	VideoCount int       `json:"video_count"`
	Completed  int       `json:"completed"`
	Failed     int       `json:"failed"`
	StartedAt  time.Time `json:"started_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Status constants
const (
	StatusPending     = "pending"
	StatusDownloading = "downloading"
	StatusCompleted   = "completed"
	StatusFailed      = "failed"
	StatusCancelled   = "cancelled"
)
