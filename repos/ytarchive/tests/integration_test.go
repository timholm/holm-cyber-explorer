//go:build integration
// +build integration

package tests

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/timholm/ytarchive/tests/mocks"
)

// TestFetchChannelInfo tests fetching channel information from YouTube
// Uses @aperturethinking as the test channel
func TestFetchChannelInfo(t *testing.T) {
	// Skip if not running integration tests with network
	if os.Getenv("INTEGRATION_NETWORK") != "true" {
		t.Skip("Skipping network integration test. Set INTEGRATION_NETWORK=true to run.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with mock first
	mockClient := mocks.NewMockYouTubeClient()
	defer mockClient.Close()

	// Verify mock has aperturethinking channel
	channel, ok := mockClient.Channels["aperturethinking"]
	if !ok {
		t.Fatal("Expected aperturethinking channel in mock data")
	}

	if channel.Snippet.Title != "Aperture Thinking" {
		t.Errorf("Expected channel title 'Aperture Thinking', got '%s'", channel.Snippet.Title)
	}

	if channel.Statistics.VideoCount != "150" {
		t.Errorf("Expected video count '150', got '%s'", channel.Statistics.VideoCount)
	}

	// Test mock server endpoint
	resp, err := http.Get(mockClient.URL() + "/youtube/v3/channels?forHandle=aperturethinking")
	if err != nil {
		t.Fatalf("Failed to query mock server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var response mocks.YouTubeChannelResponse
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response.Items) == 0 {
		t.Fatal("Expected at least one channel in response")
	}

	if response.Items[0].Snippet.Title != "Aperture Thinking" {
		t.Errorf("Expected channel title 'Aperture Thinking', got '%s'", response.Items[0].Snippet.Title)
	}

	_ = ctx // Used for timeout
}

// TestFetchVideoList tests extracting video list from a channel
func TestFetchVideoList(t *testing.T) {
	mockClient := mocks.NewMockYouTubeClient()
	defer mockClient.Close()

	// Get channel first to get the channel ID
	channel := mockClient.Channels["aperturethinking"]

	// Query for videos via search endpoint
	resp, err := http.Get(mockClient.URL() + "/youtube/v3/search?channelId=" + channel.ID + "&type=video")
	if err != nil {
		t.Fatalf("Failed to query mock server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var searchResponse struct {
		Kind  string `json:"kind"`
		Items []struct {
			ID struct {
				VideoID string `json:"videoId"`
			} `json:"id"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &searchResponse); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify we got videos
	if len(searchResponse.Items) == 0 {
		t.Error("Expected videos in response, got none")
	}

	// Check video IDs match expected mock data
	expectedVideos := map[string]bool{
		"video_001": false,
		"video_002": false,
	}

	for _, item := range searchResponse.Items {
		if _, ok := expectedVideos[item.ID.VideoID]; ok {
			expectedVideos[item.ID.VideoID] = true
		}
	}

	for videoID, found := range expectedVideos {
		if !found {
			t.Errorf("Expected video %s not found in response", videoID)
		}
	}
}

// TestDownloadSingleVideo tests the download pipeline for a single video
func TestDownloadSingleVideo(t *testing.T) {
	downloader := mocks.NewMockYouTubeDownloader()

	testVideoID := "video_001"
	testOutputPath := t.TempDir()

	// Test successful download
	result, err := downloader.Download(testVideoID, testOutputPath)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	if result.VideoID != testVideoID {
		t.Errorf("Expected video ID '%s', got '%s'", testVideoID, result.VideoID)
	}

	expectedPath := filepath.Join(testOutputPath, testVideoID+".mp4")
	if result.FilePath != expectedPath {
		t.Errorf("Expected file path '%s', got '%s'", expectedPath, result.FilePath)
	}

	if result.FileSize <= 0 {
		t.Error("Expected positive file size")
	}

	// Verify video is tracked in downloaded list
	downloaded := downloader.GetDownloadedVideos()
	if _, ok := downloaded[testVideoID]; !ok {
		t.Error("Expected video to be in downloaded videos map")
	}

	// Test download failure
	downloader.SetFailure(true, "network error")
	_, err = downloader.Download("video_002", testOutputPath)
	if err == nil {
		t.Error("Expected download to fail")
	}
}

// TestRedisQueue tests batch push and claim operations for the download queue
func TestRedisQueue(t *testing.T) {
	ctx := context.Background()
	redisClient := mocks.NewMockRedisClient()
	queue := mocks.NewMockRedisQueue(redisClient)

	queueKey := "download:queue"
	processingKey := "download:processing"
	completedKey := "download:completed"

	// Test batch push
	videoIDs := []string{"video_001", "video_002", "video_003", "video_004", "video_005"}
	err := queue.PushVideos(ctx, queueKey, videoIDs)
	if err != nil {
		t.Fatalf("Failed to push videos: %v", err)
	}

	// Verify queue length
	length, err := queue.QueueLength(ctx, queueKey)
	if err != nil {
		t.Fatalf("Failed to get queue length: %v", err)
	}
	if length != int64(len(videoIDs)) {
		t.Errorf("Expected queue length %d, got %d", len(videoIDs), length)
	}

	// Test claiming videos
	claimedVideos := make([]string, 0)
	for i := 0; i < 3; i++ {
		videoID, err := queue.ClaimVideo(ctx, queueKey, processingKey)
		if err != nil {
			t.Fatalf("Failed to claim video: %v", err)
		}
		claimedVideos = append(claimedVideos, videoID)
	}

	// Verify claimed videos are in processing set
	for _, videoID := range claimedVideos {
		isMember, err := redisClient.SIsMember(ctx, processingKey, videoID)
		if err != nil {
			t.Fatalf("Failed to check set membership: %v", err)
		}
		if !isMember {
			t.Errorf("Expected video %s to be in processing set", videoID)
		}
	}

	// Verify queue length decreased
	length, err = queue.QueueLength(ctx, queueKey)
	if err != nil {
		t.Fatalf("Failed to get queue length: %v", err)
	}
	if length != 2 {
		t.Errorf("Expected queue length 2, got %d", length)
	}

	// Test completing a video
	err = queue.CompleteVideo(ctx, processingKey, completedKey, claimedVideos[0])
	if err != nil {
		t.Fatalf("Failed to complete video: %v", err)
	}

	// Verify video moved from processing to completed
	isMember, _ := redisClient.SIsMember(ctx, processingKey, claimedVideos[0])
	if isMember {
		t.Error("Expected video to be removed from processing set")
	}

	isMember, _ = redisClient.SIsMember(ctx, completedKey, claimedVideos[0])
	if !isMember {
		t.Error("Expected video to be in completed set")
	}

	// Test concurrent claims (simulate multiple workers)
	t.Run("ConcurrentClaims", func(t *testing.T) {
		// Reset the queue
		redisClient.FlushAll(ctx)
		queue.PushVideos(ctx, queueKey, []string{"v1", "v2", "v3", "v4", "v5", "v6", "v7", "v8", "v9", "v10"})

		results := make(chan string, 10)
		errors := make(chan error, 10)

		// Simulate 5 concurrent workers claiming videos
		for i := 0; i < 5; i++ {
			go func() {
				for {
					videoID, err := queue.ClaimVideo(ctx, queueKey, processingKey)
					if err != nil {
						if mocks.IsNilError(err) {
							return // Queue is empty
						}
						errors <- err
						return
					}
					results <- videoID
				}
			}()
		}

		// Wait for all claims
		time.Sleep(100 * time.Millisecond)

		// Check that all videos were claimed without duplicates
		claimed := make(map[string]bool)
		close(results)
		for v := range results {
			if claimed[v] {
				t.Errorf("Video %s was claimed more than once", v)
			}
			claimed[v] = true
		}

		// Verify no errors
		close(errors)
		for err := range errors {
			t.Errorf("Unexpected error during concurrent claims: %v", err)
		}
	})
}

// TestSQLiteOperations tests CRUD operations for the SQLite database
func TestSQLiteOperations(t *testing.T) {
	// Create temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS channels (
			id TEXT PRIMARY KEY,
			youtube_id TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			video_count INTEGER DEFAULT 0,
			status TEXT DEFAULT 'pending',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			last_sync_at DATETIME
		);

		CREATE TABLE IF NOT EXISTS videos (
			id TEXT PRIMARY KEY,
			youtube_id TEXT UNIQUE NOT NULL,
			channel_id TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT,
			duration INTEGER DEFAULT 0,
			status TEXT DEFAULT 'pending',
			file_path TEXT,
			file_size INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (channel_id) REFERENCES channels(id)
		);

		CREATE INDEX IF NOT EXISTS idx_videos_channel_id ON videos(channel_id);
		CREATE INDEX IF NOT EXISTS idx_videos_status ON videos(status);
	`)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	// Test CREATE - Insert channel
	t.Run("CreateChannel", func(t *testing.T) {
		_, err := db.Exec(`
			INSERT INTO channels (id, youtube_id, name, description, video_count, status)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "ch-001", "UC_aperturethinking_ID", "Aperture Thinking", "Photography channel", 150, "pending")

		if err != nil {
			t.Fatalf("Failed to insert channel: %v", err)
		}
	})

	// Test READ - Query channel
	t.Run("ReadChannel", func(t *testing.T) {
		var id, youtubeID, name, status string
		var videoCount int

		err := db.QueryRow(`
			SELECT id, youtube_id, name, video_count, status FROM channels WHERE id = ?
		`, "ch-001").Scan(&id, &youtubeID, &name, &videoCount, &status)

		if err != nil {
			t.Fatalf("Failed to query channel: %v", err)
		}

		if name != "Aperture Thinking" {
			t.Errorf("Expected name 'Aperture Thinking', got '%s'", name)
		}

		if videoCount != 150 {
			t.Errorf("Expected video count 150, got %d", videoCount)
		}
	})

	// Test CREATE - Insert videos
	t.Run("CreateVideos", func(t *testing.T) {
		videos := []struct {
			id        string
			youtubeID string
			title     string
			duration  int
		}{
			{"v-001", "video_001", "Understanding Aperture", 930},
			{"v-002", "video_002", "Low Light Photography", 1365},
			{"v-003", "video_003", "Composition Rules", 1200},
		}

		for _, v := range videos {
			_, err := db.Exec(`
				INSERT INTO videos (id, youtube_id, channel_id, title, duration, status)
				VALUES (?, ?, ?, ?, ?, ?)
			`, v.id, v.youtubeID, "ch-001", v.title, v.duration, "pending")

			if err != nil {
				t.Fatalf("Failed to insert video %s: %v", v.id, err)
			}
		}
	})

	// Test READ - Query videos by channel
	t.Run("ReadVideosByChannel", func(t *testing.T) {
		rows, err := db.Query(`
			SELECT id, title, duration, status FROM videos WHERE channel_id = ?
		`, "ch-001")

		if err != nil {
			t.Fatalf("Failed to query videos: %v", err)
		}
		defer rows.Close()

		var count int
		for rows.Next() {
			var id, title, status string
			var duration int
			if err := rows.Scan(&id, &title, &duration, &status); err != nil {
				t.Fatalf("Failed to scan video row: %v", err)
			}
			count++
		}

		if count != 3 {
			t.Errorf("Expected 3 videos, got %d", count)
		}
	})

	// Test UPDATE - Update video status
	t.Run("UpdateVideoStatus", func(t *testing.T) {
		result, err := db.Exec(`
			UPDATE videos SET status = ?, file_path = ?, file_size = ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ?
		`, "downloaded", "/archive/video_001.mp4", 104857600, "v-001")

		if err != nil {
			t.Fatalf("Failed to update video: %v", err)
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected != 1 {
			t.Errorf("Expected 1 row affected, got %d", rowsAffected)
		}

		// Verify update
		var status, filePath string
		var fileSize int64
		err = db.QueryRow(`SELECT status, file_path, file_size FROM videos WHERE id = ?`, "v-001").
			Scan(&status, &filePath, &fileSize)

		if err != nil {
			t.Fatalf("Failed to verify update: %v", err)
		}

		if status != "downloaded" {
			t.Errorf("Expected status 'downloaded', got '%s'", status)
		}

		if fileSize != 104857600 {
			t.Errorf("Expected file size 104857600, got %d", fileSize)
		}
	})

	// Test DELETE - Delete a video
	t.Run("DeleteVideo", func(t *testing.T) {
		result, err := db.Exec(`DELETE FROM videos WHERE id = ?`, "v-003")
		if err != nil {
			t.Fatalf("Failed to delete video: %v", err)
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected != 1 {
			t.Errorf("Expected 1 row affected, got %d", rowsAffected)
		}

		// Verify deletion
		var count int
		db.QueryRow(`SELECT COUNT(*) FROM videos WHERE channel_id = ?`, "ch-001").Scan(&count)
		if count != 2 {
			t.Errorf("Expected 2 videos after deletion, got %d", count)
		}
	})

	// Test aggregate query - Download progress
	t.Run("DownloadProgress", func(t *testing.T) {
		var total, downloaded, pending int

		err := db.QueryRow(`
			SELECT
				COUNT(*) as total,
				SUM(CASE WHEN status = 'downloaded' THEN 1 ELSE 0 END) as downloaded,
				SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) as pending
			FROM videos WHERE channel_id = ?
		`, "ch-001").Scan(&total, &downloaded, &pending)

		if err != nil {
			t.Fatalf("Failed to get download progress: %v", err)
		}

		if total != 2 {
			t.Errorf("Expected total 2, got %d", total)
		}

		if downloaded != 1 {
			t.Errorf("Expected downloaded 1, got %d", downloaded)
		}

		if pending != 1 {
			t.Errorf("Expected pending 1, got %d", pending)
		}
	})

	// Test transaction
	t.Run("Transaction", func(t *testing.T) {
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		// Insert multiple videos in transaction
		for i := 4; i <= 6; i++ {
			_, err := tx.Exec(`
				INSERT INTO videos (id, youtube_id, channel_id, title, duration, status)
				VALUES (?, ?, ?, ?, ?, ?)
			`, "v-00"+string(rune('0'+i)), "video_00"+string(rune('0'+i)), "ch-001", "Video "+string(rune('0'+i)), 600, "pending")

			if err != nil {
				tx.Rollback()
				t.Fatalf("Failed to insert video in transaction: %v", err)
			}
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("Failed to commit transaction: %v", err)
		}

		// Verify all videos were inserted
		var count int
		db.QueryRow(`SELECT COUNT(*) FROM videos WHERE channel_id = ?`, "ch-001").Scan(&count)
		if count != 5 {
			t.Errorf("Expected 5 videos after transaction, got %d", count)
		}
	})

	// Test rollback
	t.Run("Rollback", func(t *testing.T) {
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		// Insert a video
		_, err = tx.Exec(`
			INSERT INTO videos (id, youtube_id, channel_id, title, duration, status)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "v-rollback", "video_rollback", "ch-001", "Rollback Video", 600, "pending")

		if err != nil {
			t.Fatalf("Failed to insert video: %v", err)
		}

		// Rollback instead of commit
		tx.Rollback()

		// Verify video was not inserted
		var count int
		db.QueryRow(`SELECT COUNT(*) FROM videos WHERE id = ?`, "v-rollback").Scan(&count)
		if count != 0 {
			t.Error("Expected rolled back video to not exist")
		}
	})
}

// TestYouTubeAPIFailures tests handling of YouTube API failures
func TestYouTubeAPIFailures(t *testing.T) {
	mockClient := mocks.NewMockYouTubeClient()
	defer mockClient.Close()

	// Test API failure mode
	mockClient.SetFailure(true, "quota exceeded")

	resp, err := http.Get(mockClient.URL() + "/youtube/v3/channels?forHandle=aperturethinking")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var errorResp map[string]interface{}
	json.Unmarshal(body, &errorResp)

	if errorResp["error"] == nil {
		t.Error("Expected error in response")
	}

	// Reset and verify recovery
	mockClient.SetFailure(false, "")
	resp2, _ := http.Get(mockClient.URL() + "/youtube/v3/channels?forHandle=aperturethinking")
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 after recovery, got %d", resp2.StatusCode)
	}
}

// TestLoadTestData tests loading test data from JSON files
func TestLoadTestData(t *testing.T) {
	// Load channel.json
	channelFile := filepath.Join("testdata", "channel.json")
	channelData, err := os.ReadFile(channelFile)
	if err != nil {
		t.Fatalf("Failed to read channel.json: %v", err)
	}

	var channelResp mocks.YouTubeChannelResponse
	if err := json.Unmarshal(channelData, &channelResp); err != nil {
		t.Fatalf("Failed to parse channel.json: %v", err)
	}

	if len(channelResp.Items) == 0 {
		t.Error("Expected channel items in test data")
	}

	if channelResp.Items[0].Snippet.Title != "Aperture Thinking" {
		t.Errorf("Expected channel title 'Aperture Thinking', got '%s'", channelResp.Items[0].Snippet.Title)
	}

	// Load video_list.json
	videoFile := filepath.Join("testdata", "video_list.json")
	videoData, err := os.ReadFile(videoFile)
	if err != nil {
		t.Fatalf("Failed to read video_list.json: %v", err)
	}

	var playlistResp struct {
		Kind          string `json:"kind"`
		NextPageToken string `json:"nextPageToken"`
		Items         []struct {
			Snippet struct {
				Title      string `json:"title"`
				ResourceID struct {
					VideoID string `json:"videoId"`
				} `json:"resourceId"`
			} `json:"snippet"`
		} `json:"items"`
	}

	if err := json.Unmarshal(videoData, &playlistResp); err != nil {
		t.Fatalf("Failed to parse video_list.json: %v", err)
	}

	if len(playlistResp.Items) == 0 {
		t.Error("Expected video items in test data")
	}

	// Verify expected videos
	expectedTitles := []string{
		"Understanding Aperture in Photography",
		"Low Light Photography Tips",
		"Composition Rules Every Photographer Should Know",
	}

	for i, expected := range expectedTitles {
		if i >= len(playlistResp.Items) {
			break
		}
		if playlistResp.Items[i].Snippet.Title != expected {
			t.Errorf("Expected video title '%s', got '%s'", expected, playlistResp.Items[i].Snippet.Title)
		}
	}
}
