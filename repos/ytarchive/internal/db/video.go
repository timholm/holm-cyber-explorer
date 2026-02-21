// Package db provides SQLite database operations for the YouTube Channel Archiver.
package db

import (
	"database/sql"
	"fmt"
	"time"
)

// Video represents a video record in the database.
type Video struct {
	ID                  string
	Title               string
	Description         string
	Duration            int64
	UploadDate          string
	ThumbnailURL        string
	ViewCount           int64
	Status              VideoStatus
	FilePath            string
	FileSize            int64
	Checksum            string
	DownloadStartedAt   *time.Time
	DownloadCompletedAt *time.Time
	RetryCount          int
	LastError           string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// SyncHistory represents a sync operation record in the database.
type SyncHistory struct {
	ID               int64
	StartedAt        time.Time
	CompletedAt      *time.Time
	VideosFound      int
	VideosDownloaded int
	VideosFailed     int
	Status           SyncStatus
}

// InsertVideo inserts a new video record into the database.
func InsertVideo(db *sql.DB, video *Video) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	if video == nil {
		return fmt.Errorf("video is nil")
	}

	query := `
		INSERT INTO videos (id, title, description, duration, upload_date, thumbnail_url, view_count, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title = excluded.title,
			description = excluded.description,
			duration = excluded.duration,
			upload_date = excluded.upload_date,
			thumbnail_url = excluded.thumbnail_url,
			view_count = excluded.view_count,
			updated_at = CURRENT_TIMESTAMP
	`

	status := video.Status
	if status == "" {
		status = StatusPending
	}

	_, err := db.Exec(query,
		video.ID,
		video.Title,
		video.Description,
		video.Duration,
		video.UploadDate,
		video.ThumbnailURL,
		video.ViewCount,
		status,
	)
	if err != nil {
		return fmt.Errorf("failed to insert video: %w", err)
	}

	return nil
}

// UpdateVideoStatus updates the status of a video.
func UpdateVideoStatus(db *sql.DB, videoID string, status VideoStatus) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := `UPDATE videos SET status = ? WHERE id = ?`
	result, err := db.Exec(query, status, videoID)
	if err != nil {
		return fmt.Errorf("failed to update video status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("video not found: %s", videoID)
	}

	return nil
}

// GetPendingVideos returns all videos with pending status.
func GetPendingVideos(db *sql.DB) ([]Video, error) {
	return getVideosByStatus(db, StatusPending)
}

// GetFailedVideos returns all videos with failed status.
func GetFailedVideos(db *sql.DB) ([]Video, error) {
	return getVideosByStatus(db, StatusFailed)
}

// getVideosByStatus returns all videos with the specified status.
func getVideosByStatus(db *sql.DB, status VideoStatus) ([]Video, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	query := `
		SELECT id, title, description, duration, upload_date, thumbnail_url, view_count,
		       status, file_path, file_size, checksum, download_started_at, download_completed_at,
		       retry_count, last_error, created_at, updated_at
		FROM videos
		WHERE status = ?
		ORDER BY upload_date DESC
	`

	rows, err := db.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query videos: %w", err)
	}
	defer rows.Close()

	return scanVideos(rows)
}

// GetVideoByID retrieves a video by its ID.
func GetVideoByID(db *sql.DB, videoID string) (*Video, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	query := `
		SELECT id, title, description, duration, upload_date, thumbnail_url, view_count,
		       status, file_path, file_size, checksum, download_started_at, download_completed_at,
		       retry_count, last_error, created_at, updated_at
		FROM videos
		WHERE id = ?
	`

	row := db.QueryRow(query, videoID)
	video, err := scanVideo(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get video: %w", err)
	}

	return video, nil
}

// MarkDownloadStarted marks a video as currently downloading.
func MarkDownloadStarted(db *sql.DB, videoID string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := `
		UPDATE videos
		SET status = ?, download_started_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	result, err := db.Exec(query, StatusDownloading, videoID)
	if err != nil {
		return fmt.Errorf("failed to mark download started: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("video not found: %s", videoID)
	}

	return nil
}

// MarkDownloadCompleted marks a video as successfully downloaded.
func MarkDownloadCompleted(db *sql.DB, videoID, filePath string, fileSize int64, checksum string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := `
		UPDATE videos
		SET status = ?,
		    file_path = ?,
		    file_size = ?,
		    checksum = ?,
		    download_completed_at = CURRENT_TIMESTAMP,
		    last_error = NULL
		WHERE id = ?
	`
	result, err := db.Exec(query, StatusCompleted, filePath, fileSize, checksum, videoID)
	if err != nil {
		return fmt.Errorf("failed to mark download completed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("video not found: %s", videoID)
	}

	return nil
}

// MarkDownloadFailed marks a video as failed with an error message.
func MarkDownloadFailed(db *sql.DB, videoID, errorMsg string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := `
		UPDATE videos
		SET status = ?, last_error = ?
		WHERE id = ?
	`
	result, err := db.Exec(query, StatusFailed, errorMsg, videoID)
	if err != nil {
		return fmt.Errorf("failed to mark download failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("video not found: %s", videoID)
	}

	return nil
}

// IncrementRetryCount increments the retry count for a video.
func IncrementRetryCount(db *sql.DB, videoID string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := `UPDATE videos SET retry_count = retry_count + 1 WHERE id = ?`
	result, err := db.Exec(query, videoID)
	if err != nil {
		return fmt.Errorf("failed to increment retry count: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("video not found: %s", videoID)
	}

	return nil
}

// GetAllVideos returns all videos from the database.
func GetAllVideos(db *sql.DB) ([]Video, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	query := `
		SELECT id, title, description, duration, upload_date, thumbnail_url, view_count,
		       status, file_path, file_size, checksum, download_started_at, download_completed_at,
		       retry_count, last_error, created_at, updated_at
		FROM videos
		ORDER BY upload_date DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query videos: %w", err)
	}
	defer rows.Close()

	return scanVideos(rows)
}

// GetVideoCount returns the total number of videos in the database.
func GetVideoCount(db *sql.DB) (int, error) {
	if db == nil {
		return 0, fmt.Errorf("database connection is nil")
	}

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM videos").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count videos: %w", err)
	}

	return count, nil
}

// GetVideoCountByStatus returns the number of videos with a specific status.
func GetVideoCountByStatus(db *sql.DB, status VideoStatus) (int, error) {
	if db == nil {
		return 0, fmt.Errorf("database connection is nil")
	}

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM videos WHERE status = ?", status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count videos: %w", err)
	}

	return count, nil
}

// InsertSyncHistory inserts a new sync history record.
func InsertSyncHistory(db *sql.DB, sync *SyncHistory) (int64, error) {
	if db == nil {
		return 0, fmt.Errorf("database connection is nil")
	}

	query := `
		INSERT INTO sync_history (started_at, completed_at, videos_found, videos_downloaded, videos_failed, status)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := db.Exec(query,
		sync.StartedAt,
		sync.CompletedAt,
		sync.VideosFound,
		sync.VideosDownloaded,
		sync.VideosFailed,
		sync.Status,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert sync history: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return id, nil
}

// UpdateSyncHistory updates an existing sync history record.
func UpdateSyncHistory(db *sql.DB, sync *SyncHistory) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := `
		UPDATE sync_history
		SET completed_at = ?, videos_found = ?, videos_downloaded = ?, videos_failed = ?, status = ?
		WHERE id = ?
	`

	_, err := db.Exec(query,
		sync.CompletedAt,
		sync.VideosFound,
		sync.VideosDownloaded,
		sync.VideosFailed,
		sync.Status,
		sync.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update sync history: %w", err)
	}

	return nil
}

// GetLatestSyncHistory returns the most recent sync history record.
func GetLatestSyncHistory(db *sql.DB) (*SyncHistory, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	query := `
		SELECT id, started_at, completed_at, videos_found, videos_downloaded, videos_failed, status
		FROM sync_history
		ORDER BY id DESC
		LIMIT 1
	`

	var sync SyncHistory
	var completedAt sql.NullTime

	err := db.QueryRow(query).Scan(
		&sync.ID,
		&sync.StartedAt,
		&completedAt,
		&sync.VideosFound,
		&sync.VideosDownloaded,
		&sync.VideosFailed,
		&sync.Status,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest sync history: %w", err)
	}

	if completedAt.Valid {
		sync.CompletedAt = &completedAt.Time
	}

	return &sync, nil
}

// scanVideo scans a single video row.
func scanVideo(row *sql.Row) (*Video, error) {
	var video Video
	var description, uploadDate, thumbnailURL, filePath, checksum, lastError sql.NullString
	var duration, viewCount, fileSize sql.NullInt64
	var downloadStartedAt, downloadCompletedAt sql.NullTime

	err := row.Scan(
		&video.ID,
		&video.Title,
		&description,
		&duration,
		&uploadDate,
		&thumbnailURL,
		&viewCount,
		&video.Status,
		&filePath,
		&fileSize,
		&checksum,
		&downloadStartedAt,
		&downloadCompletedAt,
		&video.RetryCount,
		&lastError,
		&video.CreatedAt,
		&video.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	video.Description = description.String
	video.UploadDate = uploadDate.String
	video.ThumbnailURL = thumbnailURL.String
	video.FilePath = filePath.String
	video.Checksum = checksum.String
	video.LastError = lastError.String
	video.Duration = duration.Int64
	video.ViewCount = viewCount.Int64
	video.FileSize = fileSize.Int64

	if downloadStartedAt.Valid {
		video.DownloadStartedAt = &downloadStartedAt.Time
	}
	if downloadCompletedAt.Valid {
		video.DownloadCompletedAt = &downloadCompletedAt.Time
	}

	return &video, nil
}

// scanVideos scans multiple video rows.
func scanVideos(rows *sql.Rows) ([]Video, error) {
	var videos []Video

	for rows.Next() {
		var video Video
		var description, uploadDate, thumbnailURL, filePath, checksum, lastError sql.NullString
		var duration, viewCount, fileSize sql.NullInt64
		var downloadStartedAt, downloadCompletedAt sql.NullTime

		err := rows.Scan(
			&video.ID,
			&video.Title,
			&description,
			&duration,
			&uploadDate,
			&thumbnailURL,
			&viewCount,
			&video.Status,
			&filePath,
			&fileSize,
			&checksum,
			&downloadStartedAt,
			&downloadCompletedAt,
			&video.RetryCount,
			&lastError,
			&video.CreatedAt,
			&video.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan video: %w", err)
		}

		video.Description = description.String
		video.UploadDate = uploadDate.String
		video.ThumbnailURL = thumbnailURL.String
		video.FilePath = filePath.String
		video.Checksum = checksum.String
		video.LastError = lastError.String
		video.Duration = duration.Int64
		video.ViewCount = viewCount.Int64
		video.FileSize = fileSize.Int64

		if downloadStartedAt.Valid {
			video.DownloadStartedAt = &downloadStartedAt.Time
		}
		if downloadCompletedAt.Valid {
			video.DownloadCompletedAt = &downloadCompletedAt.Time
		}

		videos = append(videos, video)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return videos, nil
}
