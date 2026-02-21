package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestUploadToCollector_Success(t *testing.T) {
	// Create a test server that accepts uploads
	var receivedMetadata *UploadMetadata
	var receivedFileContent []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Errorf("expected multipart/form-data content type")
		}

		// Parse multipart form
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Errorf("failed to parse multipart form: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Get metadata
		metadataStr := r.FormValue("metadata")
		if metadataStr != "" {
			receivedMetadata = &UploadMetadata{}
			json.Unmarshal([]byte(metadataStr), receivedMetadata)
		}

		// Get file
		file, _, err := r.FormFile("video")
		if err != nil {
			t.Errorf("failed to get file: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()

		receivedFileContent, _ = io.ReadAll(file)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a temp file to upload
	tempDir, err := os.MkdirTemp("", "upload_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFilePath := filepath.Join(tempDir, "test_video.mp4")
	testContent := []byte("fake video content for testing")
	if err := os.WriteFile(testFilePath, testContent, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	metadata := &UploadMetadata{
		VideoID:     "test123",
		ChannelID:   "channel456",
		ChannelName: "Test Channel",
		Title:       "Test Video",
		Filename:    "test_video.mp4",
		FileSize:    int64(len(testContent)),
	}

	ctx := context.Background()
	err = uploadToCollector(ctx, server.URL, testFilePath, metadata)
	if err != nil {
		t.Fatalf("uploadToCollector failed: %v", err)
	}

	// Verify metadata was received
	if receivedMetadata == nil {
		t.Fatal("metadata was not received")
	}
	if receivedMetadata.VideoID != "test123" {
		t.Errorf("VideoID = %q, want %q", receivedMetadata.VideoID, "test123")
	}
	if receivedMetadata.ChannelID != "channel456" {
		t.Errorf("ChannelID = %q, want %q", receivedMetadata.ChannelID, "channel456")
	}

	// Verify file content was received
	if string(receivedFileContent) != string(testContent) {
		t.Errorf("file content mismatch")
	}
}

func TestUploadToCollector_FileNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	metadata := &UploadMetadata{
		VideoID:  "test123",
		Filename: "nonexistent.mp4",
	}

	ctx := context.Background()
	err := uploadToCollector(ctx, server.URL, "/nonexistent/path/video.mp4", metadata)
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
	if !strings.Contains(err.Error(), "failed to open file") {
		t.Errorf("expected 'failed to open file' error, got: %v", err)
	}
}

func TestUploadToCollector_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	// Create a temp file
	tempDir, err := os.MkdirTemp("", "upload_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFilePath := filepath.Join(tempDir, "test_video.mp4")
	if err := os.WriteFile(testFilePath, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	metadata := &UploadMetadata{
		VideoID:  "test123",
		Filename: "test_video.mp4",
	}

	ctx := context.Background()
	err = uploadToCollector(ctx, server.URL, testFilePath, metadata)
	if err == nil {
		t.Error("expected error for server 500 response")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error to contain status 500, got: %v", err)
	}
}

func TestUploadToCollector_ContextCancellation(t *testing.T) {
	// Server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a temp file
	tempDir, err := os.MkdirTemp("", "upload_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFilePath := filepath.Join(tempDir, "test_video.mp4")
	if err := os.WriteFile(testFilePath, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	metadata := &UploadMetadata{
		VideoID:  "test123",
		Filename: "test_video.mp4",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err = uploadToCollector(ctx, server.URL, testFilePath, metadata)
	if err == nil {
		t.Error("expected error for context timeout")
	}
}

func TestUploadToCollector_InvalidURL(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "upload_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFilePath := filepath.Join(tempDir, "test_video.mp4")
	if err := os.WriteFile(testFilePath, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	metadata := &UploadMetadata{
		VideoID:  "test123",
		Filename: "test_video.mp4",
	}

	ctx := context.Background()
	err = uploadToCollector(ctx, "http://nonexistent.invalid:99999", testFilePath, metadata)
	if err == nil {
		t.Error("expected error for invalid URL")
	}
}

func TestUploadMetadata_JSON(t *testing.T) {
	metadata := &UploadMetadata{
		VideoID:       "vid123",
		ChannelID:     "ch456",
		ChannelName:   "Test Channel",
		Title:         "Test Title",
		Description:   "Test Description",
		Duration:      300,
		UploadDate:    "2024-01-01",
		EpisodeNumber: 42,
		Filename:      "video.mp4",
		FileSize:      1024000,
		Resolution:    "1080p",
		Format:        "mp4",
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		t.Fatalf("failed to marshal metadata: %v", err)
	}

	var decoded UploadMetadata
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal metadata: %v", err)
	}

	if decoded.VideoID != metadata.VideoID {
		t.Errorf("VideoID mismatch")
	}
	if decoded.ChannelID != metadata.ChannelID {
		t.Errorf("ChannelID mismatch")
	}
	if decoded.EpisodeNumber != metadata.EpisodeNumber {
		t.Errorf("EpisodeNumber mismatch")
	}
	if decoded.FileSize != metadata.FileSize {
		t.Errorf("FileSize mismatch")
	}
}

func TestUploadToCollector_LargeFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping large file test in short mode")
	}

	var receivedSize int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart and count bytes
		reader, err := r.MultipartReader()
		if err != nil {
			t.Errorf("failed to get multipart reader: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Errorf("error reading part: %v", err)
				break
			}

			if part.FormName() == "video" {
				n, _ := io.Copy(io.Discard, part)
				receivedSize = n
			}
			part.Close()
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a larger temp file (1MB)
	tempDir, err := os.MkdirTemp("", "upload_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFilePath := filepath.Join(tempDir, "large_video.mp4")
	largeContent := make([]byte, 1024*1024) // 1MB
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}
	if err := os.WriteFile(testFilePath, largeContent, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	metadata := &UploadMetadata{
		VideoID:  "large123",
		Filename: "large_video.mp4",
		FileSize: int64(len(largeContent)),
	}

	ctx := context.Background()
	err = uploadToCollector(ctx, server.URL, testFilePath, metadata)
	if err != nil {
		t.Fatalf("uploadToCollector failed for large file: %v", err)
	}

	if receivedSize != int64(len(largeContent)) {
		t.Errorf("received size = %d, want %d", receivedSize, len(largeContent))
	}
}

func TestVideoInfo_Fields(t *testing.T) {
	info := VideoInfo{
		ID:            "id123",
		YouTubeID:     "yt123",
		ChannelID:     "ch456",
		Title:         "Test Title",
		Description:   "Test Description",
		Duration:      600,
		UploadDate:    "2024-01-15",
		ThumbnailURL:  "https://example.com/thumb.jpg",
		ViewCount:     1000000,
		EpisodeNumber: 10,
		ChannelName:   "Test Channel",
		Status:        "pending",
	}

	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("failed to marshal VideoInfo: %v", err)
	}

	var decoded VideoInfo
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal VideoInfo: %v", err)
	}

	if decoded.ID != info.ID {
		t.Errorf("ID = %q, want %q", decoded.ID, info.ID)
	}
	if decoded.YouTubeID != info.YouTubeID {
		t.Errorf("YouTubeID = %q, want %q", decoded.YouTubeID, info.YouTubeID)
	}
	if decoded.Duration != info.Duration {
		t.Errorf("Duration = %d, want %d", decoded.Duration, info.Duration)
	}
	if decoded.ViewCount != info.ViewCount {
		t.Errorf("ViewCount = %d, want %d", decoded.ViewCount, info.ViewCount)
	}
}

func TestGetEnvInt(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue int
		expected     int
	}{
		{
			name:         "env not set, use default",
			envKey:       "TEST_ENV_NOT_SET_12345",
			envValue:     "",
			defaultValue: 42,
			expected:     42,
		},
		{
			name:         "env set with valid int",
			envKey:       "TEST_ENV_INT_VALID",
			envValue:     "100",
			defaultValue: 42,
			expected:     100,
		},
		{
			name:         "env set with invalid int",
			envKey:       "TEST_ENV_INT_INVALID",
			envValue:     "not-a-number",
			defaultValue: 42,
			expected:     42,
		},
		{
			name:         "env set with empty string",
			envKey:       "TEST_ENV_INT_EMPTY",
			envValue:     "",
			defaultValue: 42,
			expected:     42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}

			result := getEnvInt(tt.envKey, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getEnvInt(%q, %d) = %d, want %d", tt.envKey, tt.defaultValue, result, tt.expected)
			}
		})
	}
}
