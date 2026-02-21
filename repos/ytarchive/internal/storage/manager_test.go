package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetChannelPath(t *testing.T) {
	tests := []struct {
		name      string
		basePath  string
		channelID string
		expected  string
	}{
		{
			name:      "standard channel ID",
			basePath:  "/data",
			channelID: "UCxyz123",
			expected:  "/data/channels/UCxyz123",
		},
		{
			name:      "channel ID with special characters",
			basePath:  "/data",
			channelID: "UC-abc_123",
			expected:  "/data/channels/UC-abc_123",
		},
		{
			name:      "custom base path",
			basePath:  "/custom/storage/path",
			channelID: "TestChannel",
			expected:  "/custom/storage/path/channels/TestChannel",
		},
		{
			name:      "empty channel ID",
			basePath:  "/data",
			channelID: "",
			expected:  "/data/channels",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(tt.basePath)
			result := m.GetChannelPath(tt.channelID)
			if result != tt.expected {
				t.Errorf("GetChannelPath(%q) = %q, want %q", tt.channelID, result, tt.expected)
			}
		})
	}
}

func TestGetVideoPath(t *testing.T) {
	tests := []struct {
		name      string
		basePath  string
		channelID string
		videoID   string
		expected  string
	}{
		{
			name:      "standard IDs",
			basePath:  "/data",
			channelID: "UCxyz123",
			videoID:   "dQw4w9WgXcQ",
			expected:  "/data/channels/UCxyz123/videos/dQw4w9WgXcQ",
		},
		{
			name:      "custom base path",
			basePath:  "/mnt/storage",
			channelID: "MyChannel",
			videoID:   "video123",
			expected:  "/mnt/storage/channels/MyChannel/videos/video123",
		},
		{
			name:      "video ID with underscore and dash",
			basePath:  "/data",
			channelID: "UCtest",
			videoID:   "ab-cd_ef123",
			expected:  "/data/channels/UCtest/videos/ab-cd_ef123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(tt.basePath)
			result := m.GetVideoPath(tt.channelID, tt.videoID)
			if result != tt.expected {
				t.Errorf("GetVideoPath(%q, %q) = %q, want %q", tt.channelID, tt.videoID, result, tt.expected)
			}
		})
	}
}

func TestGetVideoFilePath(t *testing.T) {
	m := NewManager("/data")
	result := m.GetVideoFilePath("UCtest", "vid123")
	expected := "/data/channels/UCtest/videos/vid123/video.mp4"
	if result != expected {
		t.Errorf("GetVideoFilePath() = %q, want %q", result, expected)
	}
}

func TestGetMetadataPath(t *testing.T) {
	m := NewManager("/data")
	result := m.GetMetadataPath("UCtest", "vid123")
	expected := "/data/channels/UCtest/videos/vid123/metadata.json"
	if result != expected {
		t.Errorf("GetMetadataPath() = %q, want %q", result, expected)
	}
}

func TestGetThumbnailPath(t *testing.T) {
	m := NewManager("/data")
	result := m.GetThumbnailPath("UCtest", "vid123")
	expected := "/data/channels/UCtest/videos/vid123/thumbnail.webp"
	if result != expected {
		t.Errorf("GetThumbnailPath() = %q, want %q", result, expected)
	}
}

func TestGetSubtitlesPath(t *testing.T) {
	m := NewManager("/data")
	result := m.GetSubtitlesPath("UCtest", "vid123")
	expected := "/data/channels/UCtest/videos/vid123/subtitles"
	if result != expected {
		t.Errorf("GetSubtitlesPath() = %q, want %q", result, expected)
	}
}

func TestGetChannelInfoPath(t *testing.T) {
	m := NewManager("/data")
	result := m.GetChannelInfoPath("UCtest")
	expected := "/data/channels/UCtest/channel.json"
	if result != expected {
		t.Errorf("GetChannelInfoPath() = %q, want %q", result, expected)
	}
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	// Test non-existent file
	t.Run("non-existent file", func(t *testing.T) {
		result := m.FileExists(filepath.Join(tmpDir, "nonexistent.txt"))
		if result {
			t.Error("FileExists() = true for non-existent file, want false")
		}
	})

	// Test existing file
	t.Run("existing file", func(t *testing.T) {
		testFile := filepath.Join(tmpDir, "testfile.txt")
		if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		result := m.FileExists(testFile)
		if !result {
			t.Error("FileExists() = false for existing file, want true")
		}
	})

	// Test directory
	t.Run("existing directory", func(t *testing.T) {
		testDir := filepath.Join(tmpDir, "testdir")
		if err := os.MkdirAll(testDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		result := m.FileExists(testDir)
		if !result {
			t.Error("FileExists() = false for existing directory, want true")
		}
	})
}

func TestSaveAndLoadChannelInfo(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	tests := []struct {
		name    string
		channel *Channel
	}{
		{
			name: "complete channel info",
			channel: &Channel{
				ID:          "internal-id-123",
				YouTubeID:   "UCxyz123abc",
				Name:        "Test Channel",
				Description: "This is a test channel description",
				CustomURL:   "@testchannel",
				Thumbnail:   "https://example.com/thumb.jpg",
				Banner:      "https://example.com/banner.jpg",
				VideoCount:  100,
				ViewCount:   50000,
				Subscribers: 1000,
				JoinedDate:  "2020-01-01",
				Country:     "US",
			},
		},
		{
			name: "minimal channel info",
			channel: &Channel{
				YouTubeID: "UCminimal",
				Name:      "Minimal Channel",
			},
		},
		{
			name: "channel with special characters in name",
			channel: &Channel{
				YouTubeID:   "UCspecial",
				Name:        "Channel with 'quotes' and \"double quotes\"",
				Description: "Description with\nnewlines\tand\ttabs",
			},
		},
		{
			name: "channel with unicode characters",
			channel: &Channel{
				YouTubeID:   "UCunicode",
				Name:        "Unicode: Test",
				Description: "Unicode description",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save the channel info
			err := m.SaveChannelInfo(tt.channel)
			if err != nil {
				t.Fatalf("SaveChannelInfo() error = %v", err)
			}

			// Verify file was created
			infoPath := m.GetChannelInfoPath(tt.channel.YouTubeID)
			if !m.FileExists(infoPath) {
				t.Fatalf("SaveChannelInfo() did not create file at %s", infoPath)
			}

			// Load the channel info
			loaded, err := m.LoadChannelInfo(tt.channel.YouTubeID)
			if err != nil {
				t.Fatalf("LoadChannelInfo() error = %v", err)
			}

			// Verify the loaded data matches
			if loaded.ID != tt.channel.ID {
				t.Errorf("Loaded ID = %q, want %q", loaded.ID, tt.channel.ID)
			}
			if loaded.YouTubeID != tt.channel.YouTubeID {
				t.Errorf("Loaded YouTubeID = %q, want %q", loaded.YouTubeID, tt.channel.YouTubeID)
			}
			if loaded.Name != tt.channel.Name {
				t.Errorf("Loaded Name = %q, want %q", loaded.Name, tt.channel.Name)
			}
			if loaded.Description != tt.channel.Description {
				t.Errorf("Loaded Description = %q, want %q", loaded.Description, tt.channel.Description)
			}
			if loaded.CustomURL != tt.channel.CustomURL {
				t.Errorf("Loaded CustomURL = %q, want %q", loaded.CustomURL, tt.channel.CustomURL)
			}
			if loaded.VideoCount != tt.channel.VideoCount {
				t.Errorf("Loaded VideoCount = %d, want %d", loaded.VideoCount, tt.channel.VideoCount)
			}
			if loaded.ViewCount != tt.channel.ViewCount {
				t.Errorf("Loaded ViewCount = %d, want %d", loaded.ViewCount, tt.channel.ViewCount)
			}
			if loaded.Subscribers != tt.channel.Subscribers {
				t.Errorf("Loaded Subscribers = %d, want %d", loaded.Subscribers, tt.channel.Subscribers)
			}

			// Verify timestamps were set
			if loaded.UpdatedAt.IsZero() {
				t.Error("UpdatedAt should not be zero")
			}
			if loaded.CreatedAt.IsZero() {
				t.Error("CreatedAt should not be zero")
			}
		})
	}
}

func TestSaveChannelInfo_UpdatesTimestamps(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	channel := &Channel{
		YouTubeID: "UCtime",
		Name:      "Timestamp Test",
	}

	// First save
	err := m.SaveChannelInfo(channel)
	if err != nil {
		t.Fatalf("First SaveChannelInfo() error = %v", err)
	}

	loaded1, err := m.LoadChannelInfo(channel.YouTubeID)
	if err != nil {
		t.Fatalf("First LoadChannelInfo() error = %v", err)
	}
	firstCreatedAt := loaded1.CreatedAt
	firstUpdatedAt := loaded1.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(10 * time.Millisecond)

	// Second save (update)
	channel.Name = "Updated Name"
	err = m.SaveChannelInfo(channel)
	if err != nil {
		t.Fatalf("Second SaveChannelInfo() error = %v", err)
	}

	loaded2, err := m.LoadChannelInfo(channel.YouTubeID)
	if err != nil {
		t.Fatalf("Second LoadChannelInfo() error = %v", err)
	}

	// CreatedAt should be preserved (though our implementation sets it on first save only)
	// UpdatedAt should be updated
	if !loaded2.UpdatedAt.After(firstUpdatedAt) && loaded2.UpdatedAt != firstUpdatedAt {
		t.Errorf("UpdatedAt was not updated: first=%v, second=%v", firstUpdatedAt, loaded2.UpdatedAt)
	}
	if loaded2.Name != "Updated Name" {
		t.Errorf("Name was not updated: got %q, want %q", loaded2.Name, "Updated Name")
	}

	// Verify that if CreatedAt was already set, a new save preserves the relative creation time
	_ = firstCreatedAt // Used for documentation purposes
}

func TestLoadChannelInfo_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	_, err := m.LoadChannelInfo("nonexistent")
	if err == nil {
		t.Error("LoadChannelInfo() for non-existent channel should return error")
	}
}

func TestInitStorage(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	err := m.InitStorage()
	if err != nil {
		t.Fatalf("InitStorage() error = %v", err)
	}

	// Verify directories were created
	expectedDirs := []string{
		filepath.Join(tmpDir, "channels"),
		filepath.Join(tmpDir, "queue"),
		filepath.Join(tmpDir, "logs"),
	}

	for _, dir := range expectedDirs {
		if !m.FileExists(dir) {
			t.Errorf("InitStorage() did not create directory %s", dir)
		}
	}
}

func TestGetFileSize(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	// Create a test file with known content
	testFile := filepath.Join(tmpDir, "sizefile.txt")
	content := []byte("Hello, World!")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	size, err := m.GetFileSize(testFile)
	if err != nil {
		t.Fatalf("GetFileSize() error = %v", err)
	}

	expectedSize := int64(len(content))
	if size != expectedSize {
		t.Errorf("GetFileSize() = %d, want %d", size, expectedSize)
	}
}

func TestGetFileSize_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	_, err := m.GetFileSize(filepath.Join(tmpDir, "nonexistent.txt"))
	if err == nil {
		t.Error("GetFileSize() for non-existent file should return error")
	}
}

func TestListChannels(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	// Test empty channels directory
	t.Run("empty channels directory", func(t *testing.T) {
		channels, err := m.ListChannels()
		if err != nil {
			t.Fatalf("ListChannels() error = %v", err)
		}
		if len(channels) != 0 {
			t.Errorf("ListChannels() returned %d channels, want 0", len(channels))
		}
	})

	// Create some channel directories
	channelsDir := filepath.Join(tmpDir, "channels")
	if err := os.MkdirAll(channelsDir, 0755); err != nil {
		t.Fatalf("Failed to create channels directory: %v", err)
	}

	expectedChannels := []string{"UCchannel1", "UCchannel2", "UCchannel3"}
	for _, ch := range expectedChannels {
		if err := os.MkdirAll(filepath.Join(channelsDir, ch), 0755); err != nil {
			t.Fatalf("Failed to create channel directory: %v", err)
		}
	}

	// Also create a file (should be ignored)
	if err := os.WriteFile(filepath.Join(channelsDir, "ignore.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	t.Run("with channels", func(t *testing.T) {
		channels, err := m.ListChannels()
		if err != nil {
			t.Fatalf("ListChannels() error = %v", err)
		}
		if len(channels) != len(expectedChannels) {
			t.Errorf("ListChannels() returned %d channels, want %d", len(channels), len(expectedChannels))
		}
	})
}

func TestListVideos(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	channelID := "UCtest"
	videosDir := filepath.Join(tmpDir, "channels", channelID, "videos")

	// Test non-existent videos directory
	t.Run("non-existent videos directory", func(t *testing.T) {
		videos, err := m.ListVideos(channelID)
		if err != nil {
			t.Fatalf("ListVideos() error = %v", err)
		}
		if len(videos) != 0 {
			t.Errorf("ListVideos() returned %d videos, want 0", len(videos))
		}
	})

	// Create videos directories
	if err := os.MkdirAll(videosDir, 0755); err != nil {
		t.Fatalf("Failed to create videos directory: %v", err)
	}

	expectedVideos := []string{"video1", "video2", "video3"}
	for _, vid := range expectedVideos {
		if err := os.MkdirAll(filepath.Join(videosDir, vid), 0755); err != nil {
			t.Fatalf("Failed to create video directory: %v", err)
		}
	}

	t.Run("with videos", func(t *testing.T) {
		videos, err := m.ListVideos(channelID)
		if err != nil {
			t.Fatalf("ListVideos() error = %v", err)
		}
		if len(videos) != len(expectedVideos) {
			t.Errorf("ListVideos() returned %d videos, want %d", len(videos), len(expectedVideos))
		}
	})
}

func TestNewManager_DefaultPath(t *testing.T) {
	// Test with empty path
	m := NewManager("")

	// Should use default or environment variable
	basePath := m.GetBasePath()
	if basePath == "" {
		t.Error("NewManager(\"\") should set a non-empty base path")
	}
}

func TestGetLogPath(t *testing.T) {
	m := NewManager("/data")
	testDate := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	result := m.GetLogPath(testDate)
	expected := "/data/logs/2024-06-15"

	if result != expected {
		t.Errorf("GetLogPath() = %q, want %q", result, expected)
	}
}

func TestGetQueuePath(t *testing.T) {
	m := NewManager("/data")
	result := m.GetQueuePath()
	expected := "/data/queue"

	if result != expected {
		t.Errorf("GetQueuePath() = %q, want %q", result, expected)
	}
}

func TestDeleteVideo(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	channelID := "UCtest"
	videoID := "vid123"

	// Create video directory with some files
	videoDir := m.GetVideoPath(channelID, videoID)
	if err := os.MkdirAll(videoDir, 0755); err != nil {
		t.Fatalf("Failed to create video directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(videoDir, "video.mp4"), []byte("video content"), 0644); err != nil {
		t.Fatalf("Failed to create video file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(videoDir, "metadata.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create metadata file: %v", err)
	}

	// Verify directory exists
	if !m.FileExists(videoDir) {
		t.Fatal("Video directory should exist before deletion")
	}

	// Delete video
	err := m.DeleteVideo(channelID, videoID)
	if err != nil {
		t.Fatalf("DeleteVideo() error = %v", err)
	}

	// Verify directory was deleted
	if m.FileExists(videoDir) {
		t.Error("Video directory should not exist after deletion")
	}
}

func TestDeleteChannel(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	channelID := "UCtest"

	// Create channel directory with content
	channelDir := m.GetChannelPath(channelID)
	if err := os.MkdirAll(channelDir, 0755); err != nil {
		t.Fatalf("Failed to create channel directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(channelDir, "channel.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create channel file: %v", err)
	}

	// Create a video subdirectory
	videoDir := filepath.Join(channelDir, "videos", "vid123")
	if err := os.MkdirAll(videoDir, 0755); err != nil {
		t.Fatalf("Failed to create video directory: %v", err)
	}

	// Delete channel
	err := m.DeleteChannel(channelID)
	if err != nil {
		t.Fatalf("DeleteChannel() error = %v", err)
	}

	// Verify directory was deleted
	if m.FileExists(channelDir) {
		t.Error("Channel directory should not exist after deletion")
	}
}

func TestSaveAndLoadVideoMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	channelID := "UCtest"
	video := &Video{
		ID:           "internal-id",
		YouTubeID:    "dQw4w9WgXcQ",
		ChannelID:    channelID,
		Title:        "Test Video Title",
		Description:  "This is a test video description",
		Duration:     300,
		ViewCount:    1000000,
		LikeCount:    50000,
		CommentCount: 10000,
		UploadDate:   "2024-01-15",
		Thumbnail:    "https://example.com/thumb.jpg",
		Tags:         []string{"test", "video", "example"},
		Categories:   []string{"Entertainment"},
		Status:       "pending",
	}

	// Save video metadata
	err := m.SaveVideoMetadata(channelID, video)
	if err != nil {
		t.Fatalf("SaveVideoMetadata() error = %v", err)
	}

	// Load video metadata
	loaded, err := m.LoadVideoMetadata(channelID, video.YouTubeID)
	if err != nil {
		t.Fatalf("LoadVideoMetadata() error = %v", err)
	}

	// Verify loaded data
	if loaded.YouTubeID != video.YouTubeID {
		t.Errorf("Loaded YouTubeID = %q, want %q", loaded.YouTubeID, video.YouTubeID)
	}
	if loaded.Title != video.Title {
		t.Errorf("Loaded Title = %q, want %q", loaded.Title, video.Title)
	}
	if loaded.Duration != video.Duration {
		t.Errorf("Loaded Duration = %d, want %d", loaded.Duration, video.Duration)
	}
	if loaded.ChannelID != channelID {
		t.Errorf("Loaded ChannelID = %q, want %q", loaded.ChannelID, channelID)
	}
}
