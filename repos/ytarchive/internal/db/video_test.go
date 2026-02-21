package db

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

// setupTestDB creates an in-memory SQLite database with the schema applied.
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	_, err = db.Exec(Schema)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to apply schema: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestInsertVideo(t *testing.T) {
	tests := []struct {
		name    string
		video   *Video
		wantErr bool
	}{
		{
			name: "insert new video",
			video: &Video{
				ID:           "dQw4w9WgXcQ",
				Title:        "Test Video",
				Description:  "Test Description",
				Duration:     300,
				UploadDate:   "20240101",
				ThumbnailURL: "https://example.com/thumb.jpg",
				ViewCount:    1000000,
				Status:       StatusPending,
			},
			wantErr: false,
		},
		{
			name: "insert video with minimal fields",
			video: &Video{
				ID:    "minimalID",
				Title: "Minimal Video",
			},
			wantErr: false,
		},
		{
			name: "insert video with empty status defaults to pending",
			video: &Video{
				ID:     "emptyStatus",
				Title:  "Empty Status Video",
				Status: "",
			},
			wantErr: false,
		},
		{
			name: "upsert existing video",
			video: &Video{
				ID:        "dQw4w9WgXcQ",
				Title:     "Updated Title",
				ViewCount: 2000000,
			},
			wantErr: false,
		},
	}

	db := setupTestDB(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InsertVideo(db, tt.video)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertVideo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify video was inserted
				video, err := GetVideoByID(db, tt.video.ID)
				if err != nil {
					t.Errorf("GetVideoByID() error = %v", err)
					return
				}
				if video == nil {
					t.Error("GetVideoByID() returned nil for inserted video")
					return
				}
				if video.Title != tt.video.Title {
					t.Errorf("Video title = %q, want %q", video.Title, tt.video.Title)
				}
			}
		})
	}
}

func TestInsertVideo_NilInputs(t *testing.T) {
	db := setupTestDB(t)

	t.Run("nil database", func(t *testing.T) {
		err := InsertVideo(nil, &Video{ID: "test"})
		if err == nil {
			t.Error("InsertVideo() with nil db should return error")
		}
	})

	t.Run("nil video", func(t *testing.T) {
		err := InsertVideo(db, nil)
		if err == nil {
			t.Error("InsertVideo() with nil video should return error")
		}
	})
}

func TestUpdateVideoStatus(t *testing.T) {
	db := setupTestDB(t)

	// Insert a test video first
	video := &Video{
		ID:     "statusTest",
		Title:  "Status Test Video",
		Status: StatusPending,
	}
	if err := InsertVideo(db, video); err != nil {
		t.Fatalf("Failed to insert test video: %v", err)
	}

	tests := []struct {
		name      string
		videoID   string
		newStatus VideoStatus
		wantErr   bool
	}{
		{
			name:      "update to downloading",
			videoID:   "statusTest",
			newStatus: StatusDownloading,
			wantErr:   false,
		},
		{
			name:      "update to completed",
			videoID:   "statusTest",
			newStatus: StatusCompleted,
			wantErr:   false,
		},
		{
			name:      "update to failed",
			videoID:   "statusTest",
			newStatus: StatusFailed,
			wantErr:   false,
		},
		{
			name:      "update non-existent video",
			videoID:   "nonexistent",
			newStatus: StatusCompleted,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdateVideoStatus(db, tt.videoID, tt.newStatus)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateVideoStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify status was updated
				video, err := GetVideoByID(db, tt.videoID)
				if err != nil {
					t.Errorf("GetVideoByID() error = %v", err)
					return
				}
				if video.Status != tt.newStatus {
					t.Errorf("Video status = %q, want %q", video.Status, tt.newStatus)
				}
			}
		})
	}
}

func TestUpdateVideoStatus_NilDB(t *testing.T) {
	err := UpdateVideoStatus(nil, "test", StatusCompleted)
	if err == nil {
		t.Error("UpdateVideoStatus() with nil db should return error")
	}
}

func TestGetPendingVideos(t *testing.T) {
	db := setupTestDB(t)

	// Insert videos with different statuses
	testVideos := []struct {
		id     string
		title  string
		status VideoStatus
	}{
		{"pending1", "Pending Video 1", StatusPending},
		{"pending2", "Pending Video 2", StatusPending},
		{"downloading1", "Downloading Video", StatusDownloading},
		{"completed1", "Completed Video", StatusCompleted},
		{"failed1", "Failed Video", StatusFailed},
		{"pending3", "Pending Video 3", StatusPending},
	}

	for _, tv := range testVideos {
		video := &Video{
			ID:     tv.id,
			Title:  tv.title,
			Status: tv.status,
		}
		if err := InsertVideo(db, video); err != nil {
			t.Fatalf("Failed to insert test video: %v", err)
		}
	}

	// Get pending videos
	pending, err := GetPendingVideos(db)
	if err != nil {
		t.Fatalf("GetPendingVideos() error = %v", err)
	}

	// Should have 3 pending videos
	expectedCount := 3
	if len(pending) != expectedCount {
		t.Errorf("GetPendingVideos() returned %d videos, want %d", len(pending), expectedCount)
	}

	// Verify all returned videos are pending
	for _, video := range pending {
		if video.Status != StatusPending {
			t.Errorf("GetPendingVideos() returned video with status %q, want %q", video.Status, StatusPending)
		}
	}
}

func TestGetPendingVideos_NilDB(t *testing.T) {
	_, err := GetPendingVideos(nil)
	if err == nil {
		t.Error("GetPendingVideos() with nil db should return error")
	}
}

func TestGetPendingVideos_EmptyDB(t *testing.T) {
	db := setupTestDB(t)

	pending, err := GetPendingVideos(db)
	if err != nil {
		t.Fatalf("GetPendingVideos() error = %v", err)
	}

	if len(pending) != 0 {
		t.Errorf("GetPendingVideos() on empty db returned %d videos, want 0", len(pending))
	}
}

func TestGetFailedVideos(t *testing.T) {
	db := setupTestDB(t)

	// Insert videos with different statuses
	testVideos := []struct {
		id     string
		status VideoStatus
	}{
		{"pending1", StatusPending},
		{"failed1", StatusFailed},
		{"failed2", StatusFailed},
		{"completed1", StatusCompleted},
	}

	for _, tv := range testVideos {
		video := &Video{
			ID:     tv.id,
			Title:  "Test",
			Status: tv.status,
		}
		if err := InsertVideo(db, video); err != nil {
			t.Fatalf("Failed to insert test video: %v", err)
		}
	}

	failed, err := GetFailedVideos(db)
	if err != nil {
		t.Fatalf("GetFailedVideos() error = %v", err)
	}

	if len(failed) != 2 {
		t.Errorf("GetFailedVideos() returned %d videos, want 2", len(failed))
	}
}

func TestGetVideoByID(t *testing.T) {
	db := setupTestDB(t)

	// Insert test video
	video := &Video{
		ID:           "testID123",
		Title:        "Test Video",
		Description:  "Test Description",
		Duration:     600,
		UploadDate:   "20240115",
		ThumbnailURL: "https://example.com/thumb.jpg",
		ViewCount:    5000,
		Status:       StatusPending,
	}
	if err := InsertVideo(db, video); err != nil {
		t.Fatalf("Failed to insert test video: %v", err)
	}

	tests := []struct {
		name      string
		videoID   string
		wantNil   bool
		wantTitle string
	}{
		{
			name:      "existing video",
			videoID:   "testID123",
			wantNil:   false,
			wantTitle: "Test Video",
		},
		{
			name:    "non-existent video",
			videoID: "nonexistent",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetVideoByID(db, tt.videoID)
			if err != nil {
				t.Fatalf("GetVideoByID() error = %v", err)
			}

			if tt.wantNil {
				if result != nil {
					t.Error("GetVideoByID() expected nil for non-existent video")
				}
			} else {
				if result == nil {
					t.Error("GetVideoByID() returned nil for existing video")
					return
				}
				if result.Title != tt.wantTitle {
					t.Errorf("GetVideoByID() title = %q, want %q", result.Title, tt.wantTitle)
				}
			}
		})
	}
}

func TestMarkDownloadStarted(t *testing.T) {
	db := setupTestDB(t)

	video := &Video{
		ID:     "downloadTest",
		Title:  "Download Test",
		Status: StatusPending,
	}
	if err := InsertVideo(db, video); err != nil {
		t.Fatalf("Failed to insert test video: %v", err)
	}

	err := MarkDownloadStarted(db, "downloadTest")
	if err != nil {
		t.Fatalf("MarkDownloadStarted() error = %v", err)
	}

	// Verify status and timestamp
	result, err := GetVideoByID(db, "downloadTest")
	if err != nil {
		t.Fatalf("GetVideoByID() error = %v", err)
	}

	if result.Status != StatusDownloading {
		t.Errorf("Status = %q, want %q", result.Status, StatusDownloading)
	}

	if result.DownloadStartedAt == nil {
		t.Error("DownloadStartedAt should not be nil")
	}
}

func TestMarkDownloadCompleted(t *testing.T) {
	db := setupTestDB(t)

	video := &Video{
		ID:     "completeTest",
		Title:  "Complete Test",
		Status: StatusDownloading,
	}
	if err := InsertVideo(db, video); err != nil {
		t.Fatalf("Failed to insert test video: %v", err)
	}

	filePath := "/data/videos/completeTest/video.mp4"
	fileSize := int64(1024 * 1024 * 100) // 100MB
	checksum := "abc123def456"

	err := MarkDownloadCompleted(db, "completeTest", filePath, fileSize, checksum)
	if err != nil {
		t.Fatalf("MarkDownloadCompleted() error = %v", err)
	}

	result, err := GetVideoByID(db, "completeTest")
	if err != nil {
		t.Fatalf("GetVideoByID() error = %v", err)
	}

	if result.Status != StatusCompleted {
		t.Errorf("Status = %q, want %q", result.Status, StatusCompleted)
	}
	if result.FilePath != filePath {
		t.Errorf("FilePath = %q, want %q", result.FilePath, filePath)
	}
	if result.FileSize != fileSize {
		t.Errorf("FileSize = %d, want %d", result.FileSize, fileSize)
	}
	if result.Checksum != checksum {
		t.Errorf("Checksum = %q, want %q", result.Checksum, checksum)
	}
	if result.DownloadCompletedAt == nil {
		t.Error("DownloadCompletedAt should not be nil")
	}
}

func TestMarkDownloadFailed(t *testing.T) {
	db := setupTestDB(t)

	video := &Video{
		ID:     "failTest",
		Title:  "Fail Test",
		Status: StatusDownloading,
	}
	if err := InsertVideo(db, video); err != nil {
		t.Fatalf("Failed to insert test video: %v", err)
	}

	errorMsg := "Network error: connection timeout"
	err := MarkDownloadFailed(db, "failTest", errorMsg)
	if err != nil {
		t.Fatalf("MarkDownloadFailed() error = %v", err)
	}

	result, err := GetVideoByID(db, "failTest")
	if err != nil {
		t.Fatalf("GetVideoByID() error = %v", err)
	}

	if result.Status != StatusFailed {
		t.Errorf("Status = %q, want %q", result.Status, StatusFailed)
	}
	if result.LastError != errorMsg {
		t.Errorf("LastError = %q, want %q", result.LastError, errorMsg)
	}
}

func TestIncrementRetryCount(t *testing.T) {
	db := setupTestDB(t)

	video := &Video{
		ID:     "retryTest",
		Title:  "Retry Test",
		Status: StatusFailed,
	}
	if err := InsertVideo(db, video); err != nil {
		t.Fatalf("Failed to insert test video: %v", err)
	}

	// Increment retry count multiple times
	for i := 1; i <= 3; i++ {
		err := IncrementRetryCount(db, "retryTest")
		if err != nil {
			t.Fatalf("IncrementRetryCount() attempt %d error = %v", i, err)
		}

		result, err := GetVideoByID(db, "retryTest")
		if err != nil {
			t.Fatalf("GetVideoByID() error = %v", err)
		}

		if result.RetryCount != i {
			t.Errorf("After attempt %d: RetryCount = %d, want %d", i, result.RetryCount, i)
		}
	}
}

func TestGetAllVideos(t *testing.T) {
	db := setupTestDB(t)

	// Insert multiple videos
	for i := 1; i <= 5; i++ {
		video := &Video{
			ID:     "video" + string(rune('0'+i)),
			Title:  "Video " + string(rune('0'+i)),
			Status: StatusPending,
		}
		if err := InsertVideo(db, video); err != nil {
			t.Fatalf("Failed to insert test video: %v", err)
		}
	}

	videos, err := GetAllVideos(db)
	if err != nil {
		t.Fatalf("GetAllVideos() error = %v", err)
	}

	if len(videos) != 5 {
		t.Errorf("GetAllVideos() returned %d videos, want 5", len(videos))
	}
}

func TestGetVideoCount(t *testing.T) {
	db := setupTestDB(t)

	// Empty database
	count, err := GetVideoCount(db)
	if err != nil {
		t.Fatalf("GetVideoCount() error = %v", err)
	}
	if count != 0 {
		t.Errorf("GetVideoCount() on empty db = %d, want 0", count)
	}

	// Insert videos
	for i := 1; i <= 3; i++ {
		video := &Video{
			ID:     "count" + string(rune('0'+i)),
			Title:  "Count Test",
			Status: StatusPending,
		}
		if err := InsertVideo(db, video); err != nil {
			t.Fatalf("Failed to insert test video: %v", err)
		}
	}

	count, err = GetVideoCount(db)
	if err != nil {
		t.Fatalf("GetVideoCount() error = %v", err)
	}
	if count != 3 {
		t.Errorf("GetVideoCount() = %d, want 3", count)
	}
}

func TestGetVideoCountByStatus(t *testing.T) {
	db := setupTestDB(t)

	testVideos := []struct {
		id     string
		status VideoStatus
	}{
		{"v1", StatusPending},
		{"v2", StatusPending},
		{"v3", StatusDownloading},
		{"v4", StatusCompleted},
		{"v5", StatusCompleted},
		{"v6", StatusCompleted},
		{"v7", StatusFailed},
	}

	for _, tv := range testVideos {
		video := &Video{
			ID:     tv.id,
			Title:  "Test",
			Status: tv.status,
		}
		if err := InsertVideo(db, video); err != nil {
			t.Fatalf("Failed to insert test video: %v", err)
		}
	}

	tests := []struct {
		status    VideoStatus
		wantCount int
	}{
		{StatusPending, 2},
		{StatusDownloading, 1},
		{StatusCompleted, 3},
		{StatusFailed, 1},
		{StatusSkipped, 0},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			count, err := GetVideoCountByStatus(db, tt.status)
			if err != nil {
				t.Fatalf("GetVideoCountByStatus() error = %v", err)
			}
			if count != tt.wantCount {
				t.Errorf("GetVideoCountByStatus(%q) = %d, want %d", tt.status, count, tt.wantCount)
			}
		})
	}
}

func TestSyncHistory(t *testing.T) {
	db := setupTestDB(t)

	// Test insert
	sync := &SyncHistory{
		VideosFound:      100,
		VideosDownloaded: 50,
		VideosFailed:     5,
		Status:           SyncStatusRunning,
	}

	id, err := InsertSyncHistory(db, sync)
	if err != nil {
		t.Fatalf("InsertSyncHistory() error = %v", err)
	}
	if id <= 0 {
		t.Errorf("InsertSyncHistory() returned invalid id = %d", id)
	}

	sync.ID = id

	// Test get latest
	latest, err := GetLatestSyncHistory(db)
	if err != nil {
		t.Fatalf("GetLatestSyncHistory() error = %v", err)
	}
	if latest == nil {
		t.Fatal("GetLatestSyncHistory() returned nil")
	}
	if latest.ID != id {
		t.Errorf("GetLatestSyncHistory() ID = %d, want %d", latest.ID, id)
	}
	if latest.VideosFound != 100 {
		t.Errorf("GetLatestSyncHistory() VideosFound = %d, want 100", latest.VideosFound)
	}

	// Test update
	sync.VideosDownloaded = 75
	sync.Status = SyncStatusCompleted
	err = UpdateSyncHistory(db, sync)
	if err != nil {
		t.Fatalf("UpdateSyncHistory() error = %v", err)
	}

	updated, err := GetLatestSyncHistory(db)
	if err != nil {
		t.Fatalf("GetLatestSyncHistory() after update error = %v", err)
	}
	if updated.VideosDownloaded != 75 {
		t.Errorf("After update: VideosDownloaded = %d, want 75", updated.VideosDownloaded)
	}
	if updated.Status != SyncStatusCompleted {
		t.Errorf("After update: Status = %q, want %q", updated.Status, SyncStatusCompleted)
	}
}

func TestGetLatestSyncHistory_Empty(t *testing.T) {
	db := setupTestDB(t)

	latest, err := GetLatestSyncHistory(db)
	if err != nil {
		t.Fatalf("GetLatestSyncHistory() error = %v", err)
	}
	if latest != nil {
		t.Errorf("GetLatestSyncHistory() on empty db should return nil, got %+v", latest)
	}
}
