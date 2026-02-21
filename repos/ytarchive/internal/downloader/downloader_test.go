package downloader

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewDownloader(t *testing.T) {
	config := DefaultConfig("/data")
	reporter := NewProgressReporter("http://localhost:8080", "worker1")

	d := NewDownloader(config, reporter, "test-worker")

	if d == nil {
		t.Fatal("NewDownloader returned nil")
	}

	if d.config != config {
		t.Error("config was not set correctly")
	}

	if d.reporter != reporter {
		t.Error("reporter was not set correctly")
	}

	if d.workerID != "test-worker" {
		t.Errorf("workerID = %q, want %q", d.workerID, "test-worker")
	}
}

func TestGetRetryDelay(t *testing.T) {
	config := &Config{
		RetryDelays: []int{5, 15, 45},
	}
	d := NewDownloader(config, nil, "worker")

	tests := []struct {
		attempt  int
		expected int // seconds
	}{
		{1, 5},
		{2, 15},
		{3, 45},
		{4, 45}, // Beyond configured delays, use last
		{5, 45},
	}

	for _, tt := range tests {
		t.Run("attempt "+string(rune('0'+tt.attempt)), func(t *testing.T) {
			delay := d.getRetryDelay(tt.attempt)
			expectedDuration := int64(tt.expected) * int64(1e9) // Convert to nanoseconds
			if delay.Nanoseconds() != expectedDuration {
				t.Errorf("getRetryDelay(%d) = %v, want %ds", tt.attempt, delay, tt.expected)
			}
		})
	}
}

func TestGetRetryDelay_EmptyDelays(t *testing.T) {
	config := &Config{
		RetryDelays: []int{},
	}
	d := NewDownloader(config, nil, "worker")

	delay := d.getRetryDelay(1)
	expectedDuration := int64(30) * int64(1e9) // 30 seconds default

	if delay.Nanoseconds() != expectedDuration {
		t.Errorf("getRetryDelay with empty delays = %v, want 30s", delay)
	}
}

func TestVideoURL(t *testing.T) {
	tests := []struct {
		videoID  string
		expected string
	}{
		{"dQw4w9WgXcQ", "https://www.youtube.com/watch?v=dQw4w9WgXcQ"},
		{"abc123", "https://www.youtube.com/watch?v=abc123"},
		{"", "https://www.youtube.com/watch?v="},
	}

	for _, tt := range tests {
		t.Run(tt.videoID, func(t *testing.T) {
			result := VideoURL(tt.videoID)
			if result != tt.expected {
				t.Errorf("VideoURL(%q) = %q, want %q", tt.videoID, result, tt.expected)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig("/data/output")

	if config.OutputPath != "/data/output" {
		t.Errorf("OutputPath = %q, want %q", config.OutputPath, "/data/output")
	}

	if config.MaxHeight != 4320 {
		t.Errorf("MaxHeight = %d, want 4320 (highest available)", config.MaxHeight)
	}

	if !config.WriteThumbnail {
		t.Error("WriteThumbnail should be true by default")
	}

	if !config.WriteInfoJSON {
		t.Error("WriteInfoJSON should be true by default")
	}

	if !config.WriteSubtitles {
		t.Error("WriteSubtitles should be true by default")
	}

	if config.MergeOutputFormat != "mp4" {
		t.Errorf("MergeOutputFormat = %q, want %q", config.MergeOutputFormat, "mp4")
	}

	if config.Retries != 3 {
		t.Errorf("Retries = %d, want 3", config.Retries)
	}

	if len(config.RetryDelays) != 3 {
		t.Errorf("RetryDelays length = %d, want 3", len(config.RetryDelays))
	}

	if config.PreferCombinedStream {
		t.Error("PreferCombinedStream should be false by default to prefer highest quality")
	}
}

func TestSubtitleLangsString(t *testing.T) {
	tests := []struct {
		name     string
		langs    []string
		expected string
	}{
		{
			name:     "single language",
			langs:    []string{"en"},
			expected: "en",
		},
		{
			name:     "multiple languages",
			langs:    []string{"en", "es", "fr"},
			expected: "en,es,fr",
		},
		{
			name:     "empty slice defaults to en",
			langs:    []string{},
			expected: "en",
		},
		{
			name:     "nil slice defaults to en",
			langs:    nil,
			expected: "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{SubtitleLangs: tt.langs}
			result := config.SubtitleLangsString()
			if result != tt.expected {
				t.Errorf("SubtitleLangsString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestOutputTemplate(t *testing.T) {
	config := &Config{OutputPath: "/data/videos"}

	template := config.OutputTemplate("testVideo123")

	if !strings.Contains(template, "/data/videos") {
		t.Errorf("OutputTemplate should contain output path, got: %s", template)
	}

	if !strings.Contains(template, "testVideo123") {
		t.Errorf("OutputTemplate should contain video ID, got: %s", template)
	}
}

func TestVideoDir(t *testing.T) {
	config := &Config{OutputPath: "/data"}

	// With channel ID
	dir := config.VideoDir("UC123", "video456")
	expected := "/data/channels/UC123/videos/video456"
	if dir != expected {
		t.Errorf("VideoDir with channel = %q, want %q", dir, expected)
	}

	// Without channel ID
	dir = config.VideoDir("", "video456")
	expected = "/data/video456"
	if dir != expected {
		t.Errorf("VideoDir without channel = %q, want %q", dir, expected)
	}
}

func TestStreamSelector_SelectBestStream(t *testing.T) {
	streams := []Stream{
		{FormatID: "18", Height: 360, StreamType: StreamTypeCombined, URL: "http://example.com/18", Extension: "mp4"},
		{FormatID: "22", Height: 720, StreamType: StreamTypeCombined, URL: "http://example.com/22", Extension: "mp4"},
		{FormatID: "137", Height: 1080, StreamType: StreamTypeVideo, URL: "http://example.com/137", Extension: "mp4"},
		{FormatID: "140", Height: 0, StreamType: StreamTypeAudio, URL: "http://example.com/140", Extension: "m4a", Bitrate: 128000},
	}

	t.Run("prefer combined stream always when preferCombined is true", func(t *testing.T) {
		// When preferCombined is true (no ffmpeg), ALWAYS use combined streams
		// even if separate video has higher resolution (video-only won't have audio)
		selector := NewStreamSelector(1080, true)
		video, audio, err := selector.SelectBestStream(streams)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if video == nil {
			t.Fatal("video stream should not be nil")
		}
		// Should use combined 22 (720p) because preferCombined is true
		if video.FormatID != "22" {
			t.Errorf("expected format 22 (combined), got %s", video.FormatID)
		}
		if audio != nil {
			t.Error("audio should be nil for combined stream")
		}
	})

	t.Run("prefer combined when same resolution", func(t *testing.T) {
		// Test with streams where combined has same resolution as video-only
		sameResStreams := []Stream{
			{FormatID: "22", Height: 720, StreamType: StreamTypeCombined, URL: "http://example.com/22", Extension: "mp4"},
			{FormatID: "136", Height: 720, StreamType: StreamTypeVideo, URL: "http://example.com/136", Extension: "mp4"},
			{FormatID: "140", Height: 0, StreamType: StreamTypeAudio, URL: "http://example.com/140", Extension: "m4a", Bitrate: 128000},
		}
		selector := NewStreamSelector(720, true)
		video, audio, err := selector.SelectBestStream(sameResStreams)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if video == nil {
			t.Fatal("video stream should not be nil")
		}
		// Should prefer combined 22 when heights are equal and preferCombined is true
		if video.FormatID != "22" {
			t.Errorf("expected format 22 (combined same height), got %s", video.FormatID)
		}
		if audio != nil {
			t.Error("audio should be nil for combined stream")
		}
	})

	t.Run("prefer separate streams", func(t *testing.T) {
		selector := NewStreamSelector(1080, false)
		video, audio, err := selector.SelectBestStream(streams)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if video == nil {
			t.Fatal("video stream should not be nil")
		}
		if video.FormatID != "137" {
			t.Errorf("expected format 137, got %s", video.FormatID)
		}
		if audio == nil {
			t.Fatal("audio stream should not be nil for separate streams")
		}
		if audio.FormatID != "140" {
			t.Errorf("expected audio format 140, got %s", audio.FormatID)
		}
	})

	t.Run("max height filter", func(t *testing.T) {
		selector := NewStreamSelector(720, true)
		video, _, err := selector.SelectBestStream(streams)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if video.Height > 720 {
			t.Errorf("expected height <= 720, got %d", video.Height)
		}
	})

	t.Run("no streams", func(t *testing.T) {
		selector := NewStreamSelector(1080, true)
		_, _, err := selector.SelectBestStream([]Stream{})
		if err == nil {
			t.Error("expected error for empty streams")
		}
	})
}

func TestExtractVideoID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{"plain ID", "dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"watch URL", "https://www.youtube.com/watch?v=dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"youtu.be URL", "https://youtu.be/dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"embed URL", "https://www.youtube.com/embed/dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"shorts URL", "https://www.youtube.com/shorts/dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"invalid URL", "not a url at all", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractVideoID(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("expected error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("ExtractVideoID(%q) = %q, want %q", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestQualityLabel(t *testing.T) {
	tests := []struct {
		height   int
		expected string
	}{
		{2160, "4K"},
		{1440, "1440p"},
		{1080, "1080p"},
		{720, "720p"},
		{480, "480p"},
		{360, "360p"},
		{240, "240p"},
		{144, "144p"},
		{100, "100p"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := QualityLabel(tt.height)
			if result != tt.expected {
				t.Errorf("QualityLabel(%d) = %q, want %q", tt.height, result, tt.expected)
			}
		})
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.0 KiB"},
		{1024 * 1024, "1.0 MiB"},
		{1024 * 1024 * 1024, "1.0 GiB"},
		{1536 * 1024, "1.5 MiB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := FormatFileSize(tt.bytes)
			if result != tt.expected {
				t.Errorf("FormatFileSize(%d) = %q, want %q", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestFormatSpeed(t *testing.T) {
	tests := []struct {
		speed    float64
		expected string
	}{
		{500, "500 B/s"},
		{1024, "1.0 KiB/s"},
		{1024 * 1024, "1.0 MiB/s"},
		{1024 * 1024 * 1024, "1.0 GiB/s"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := FormatSpeed(tt.speed)
			if result != tt.expected {
				t.Errorf("FormatSpeed(%f) = %q, want %q", tt.speed, result, tt.expected)
			}
		})
	}
}

func TestParseStreamType(t *testing.T) {
	tests := []struct {
		vcodec   string
		acodec   string
		expected StreamType
	}{
		{"avc1", "mp4a", StreamTypeCombined},
		{"avc1", "none", StreamTypeVideo},
		{"avc1", "", StreamTypeVideo},
		{"none", "mp4a", StreamTypeAudio},
		{"", "mp4a", StreamTypeAudio},
		{"", "", StreamTypeCombined}, // default
	}

	for _, tt := range tests {
		t.Run(string(tt.expected), func(t *testing.T) {
			result := ParseStreamType(tt.vcodec, tt.acodec)
			if result != tt.expected {
				t.Errorf("ParseStreamType(%q, %q) = %q, want %q", tt.vcodec, tt.acodec, result, tt.expected)
			}
		})
	}
}

func TestParseMimeType(t *testing.T) {
	tests := []struct {
		mimeType  string
		container string
		vcodec    string
		acodec    string
	}{
		{
			"video/mp4; codecs=\"avc1.42001E, mp4a.40.2\"",
			"mp4", "avc1.42001E", "mp4a.40.2",
		},
		{
			"video/webm; codecs=\"vp9\"",
			"webm", "vp9", "",
		},
		{
			"audio/mp4; codecs=\"mp4a.40.2\"",
			"mp4", "", "mp4a.40.2",
		},
		{
			"video/mp4",
			"mp4", "", "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.mimeType, func(t *testing.T) {
			container, vcodec, acodec := ParseMimeType(tt.mimeType)
			if container != tt.container {
				t.Errorf("container = %q, want %q", container, tt.container)
			}
			if vcodec != tt.vcodec {
				t.Errorf("vcodec = %q, want %q", vcodec, tt.vcodec)
			}
			if acodec != tt.acodec {
				t.Errorf("acodec = %q, want %q", acodec, tt.acodec)
			}
		})
	}
}

func TestProgressTracker(t *testing.T) {
	var received *DownloadProgress

	callback := func(progress *DownloadProgress) {
		received = progress
	}

	tracker := NewProgressTracker("video123", "worker1", 1000, callback)

	// Simulate progress updates with delay to trigger callback
	tracker.Update(100)
	time.Sleep(600 * time.Millisecond)
	tracker.Update(500)

	if received == nil {
		t.Fatal("callback was not called")
	}

	if received.VideoID != "video123" {
		t.Errorf("VideoID = %q, want %q", received.VideoID, "video123")
	}

	if received.Status != "downloading" {
		t.Errorf("Status = %q, want %q", received.Status, "downloading")
	}
}

func TestProgressTracker_Complete(t *testing.T) {
	var received *DownloadProgress

	callback := func(progress *DownloadProgress) {
		received = progress
	}

	tracker := NewProgressTracker("video123", "worker1", 1000, callback)
	tracker.Complete()

	if received == nil {
		t.Fatal("callback was not called")
	}

	if received.Status != "completed" {
		t.Errorf("Status = %q, want %q", received.Status, "completed")
	}

	if received.Percentage != 100 {
		t.Errorf("Percentage = %f, want 100", received.Percentage)
	}
}

func TestDownloader_downloadFile(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "13")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))
	defer server.Close()

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "downloader_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := DefaultConfig(tempDir)
	d := NewDownloader(config, nil, "test-worker")

	outputPath := filepath.Join(tempDir, "test.txt")
	err = d.downloadFile(context.Background(), server.URL, outputPath, nil)
	if err != nil {
		t.Fatalf("downloadFile failed: %v", err)
	}

	// Verify file contents
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if string(content) != "Hello, World!" {
		t.Errorf("file content = %q, want %q", string(content), "Hello, World!")
	}
}

func TestDownloader_downloadFileWithResume(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		rangeHeader := r.Header.Get("Range")
		if rangeHeader != "" && strings.HasPrefix(rangeHeader, "bytes=5-") {
			// Resume request - return the remaining content
			w.Header().Set("Content-Range", "bytes 5-12/13")
			w.Header().Set("Content-Length", "8")
			w.WriteHeader(http.StatusPartialContent)
			w.Write([]byte(", World!"))
		} else {
			// Initial request - return full content
			w.Header().Set("Accept-Ranges", "bytes")
			w.Header().Set("Content-Length", "13")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello, World!"))
		}
	}))
	defer server.Close()

	tempDir, err := os.MkdirTemp("", "downloader_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := DefaultConfig(tempDir)
	d := NewDownloader(config, nil, "test-worker")

	outputPath := filepath.Join(tempDir, "test.txt")

	// First download (full)
	err = d.downloadFileWithResume(context.Background(), server.URL, outputPath, 0, 13, nil)
	if err != nil {
		t.Fatalf("first download failed: %v", err)
	}

	// Verify first download content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(content) != "Hello, World!" {
		t.Errorf("first download content = %q, want %q", string(content), "Hello, World!")
	}

	// Create a new file to test resume
	outputPath2 := filepath.Join(tempDir, "test2.txt")

	// Simulate resume from byte 5 (write partial content first)
	if err := os.WriteFile(outputPath2, []byte("Hello"), 0644); err != nil {
		t.Fatalf("failed to write partial file: %v", err)
	}

	err = d.downloadFileWithResume(context.Background(), server.URL, outputPath2, 5, 13, nil)
	if err != nil {
		t.Fatalf("resume download failed: %v", err)
	}

	// Server should have been called twice
	if callCount != 2 {
		t.Errorf("server call count = %d, want 2", callCount)
	}

	// Verify final content (after resume)
	content2, err := os.ReadFile(outputPath2)
	if err != nil {
		t.Fatalf("failed to read resumed file: %v", err)
	}
	if string(content2) != "Hello, World!" {
		t.Errorf("resumed content = %q, want %q", string(content2), "Hello, World!")
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		hasError bool
	}{
		{
			name:     "valid config",
			config:   DefaultConfig("/data"),
			hasError: false,
		},
		{
			name: "empty output path",
			config: &Config{
				OutputPath: "",
				MaxHeight:  1080,
				Retries:    3,
			},
			hasError: true,
		},
		{
			name: "zero max height",
			config: &Config{
				OutputPath: "/data",
				MaxHeight:  0,
				Retries:    3,
			},
			hasError: true,
		},
		{
			name: "zero retries",
			config: &Config{
				OutputPath: "/data",
				MaxHeight:  1080,
				Retries:    0,
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.hasError && err == nil {
				t.Error("expected error")
			}
			if !tt.hasError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestNewProgressReporter(t *testing.T) {
	reporter := NewProgressReporter("http://localhost:8080/", "worker1")

	if reporter.controllerURL != "http://localhost:8080" {
		t.Errorf("controllerURL = %q, want trailing slash removed", reporter.controllerURL)
	}

	if reporter.workerID != "worker1" {
		t.Errorf("workerID = %q, want %q", reporter.workerID, "worker1")
	}
}

func TestKnownFormats(t *testing.T) {
	// Test that known formats are correctly defined
	combinedFormats := []int{18, 22}
	for _, itag := range combinedFormats {
		info, ok := KnownFormats[itag]
		if !ok {
			t.Errorf("format %d not found", itag)
			continue
		}
		if !info.HasAudio || !info.HasVideo {
			t.Errorf("format %d should be combined (has audio and video)", itag)
		}
	}

	videoOnlyFormats := []int{137, 136, 135}
	for _, itag := range videoOnlyFormats {
		info, ok := KnownFormats[itag]
		if !ok {
			t.Errorf("format %d not found", itag)
			continue
		}
		if info.HasAudio {
			t.Errorf("format %d should be video only", itag)
		}
		if !info.HasVideo {
			t.Errorf("format %d should have video", itag)
		}
	}

	audioOnlyFormats := []int{140, 139, 251}
	for _, itag := range audioOnlyFormats {
		info, ok := KnownFormats[itag]
		if !ok {
			t.Errorf("format %d not found", itag)
			continue
		}
		if info.HasVideo {
			t.Errorf("format %d should be audio only", itag)
		}
		if !info.HasAudio {
			t.Errorf("format %d should have audio", itag)
		}
	}
}

func TestParseStreamFromItag(t *testing.T) {
	tests := []struct {
		itag      int
		expected  StreamType
		hasStream bool
	}{
		{22, StreamTypeCombined, true},
		{18, StreamTypeCombined, true},
		{137, StreamTypeVideo, true},
		{140, StreamTypeAudio, true},
		{99999, "", false}, // Unknown format
	}

	for _, tt := range tests {
		t.Run(string(rune('0'+tt.itag)), func(t *testing.T) {
			stream := ParseStreamFromItag(tt.itag, "http://example.com/stream")
			if tt.hasStream {
				if stream == nil {
					t.Fatal("expected stream, got nil")
				}
				if stream.StreamType != tt.expected {
					t.Errorf("StreamType = %q, want %q", stream.StreamType, tt.expected)
				}
				if stream.URL != "http://example.com/stream" {
					t.Errorf("URL not set correctly")
				}
			} else {
				if stream != nil {
					t.Errorf("expected nil for unknown format, got %+v", stream)
				}
			}
		})
	}
}

func TestGetAvailableResolutions(t *testing.T) {
	streams := []Stream{
		{FormatID: "18", Height: 360, Width: 640, StreamType: StreamTypeCombined, Bitrate: 500000},
		{FormatID: "22", Height: 720, Width: 1280, StreamType: StreamTypeCombined, Bitrate: 2500000},
		{FormatID: "134", Height: 360, Width: 640, StreamType: StreamTypeVideo, Bitrate: 600000},
		{FormatID: "136", Height: 720, Width: 1280, StreamType: StreamTypeVideo, Bitrate: 3000000},
		{FormatID: "137", Height: 1080, Width: 1920, StreamType: StreamTypeVideo, Bitrate: 5000000, FPS: 30},
		{FormatID: "140", Height: 0, StreamType: StreamTypeAudio, Bitrate: 128000}, // Audio-only, should be excluded
	}

	resolutions := GetAvailableResolutions(streams)

	// Should have 3 unique heights (360, 720, 1080)
	if len(resolutions) != 3 {
		t.Errorf("expected 3 resolutions, got %d", len(resolutions))
	}

	// Should be sorted by height descending
	if resolutions[0].Height != 1080 {
		t.Errorf("first resolution should be 1080p, got %d", resolutions[0].Height)
	}
	if resolutions[1].Height != 720 {
		t.Errorf("second resolution should be 720p, got %d", resolutions[1].Height)
	}
	if resolutions[2].Height != 360 {
		t.Errorf("third resolution should be 360p, got %d", resolutions[2].Height)
	}

	// Check that labels are set correctly
	if resolutions[0].Label != "1080p" {
		t.Errorf("expected label '1080p', got %s", resolutions[0].Label)
	}

	// For 720p, should prefer the higher bitrate video stream (3000000) over combined (2500000)
	if resolutions[1].Bitrate != 3000000 {
		t.Errorf("expected 720p to have bitrate 3000000, got %d", resolutions[1].Bitrate)
	}
}
