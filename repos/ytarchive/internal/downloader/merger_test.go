package downloader

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestFindFFmpeg(t *testing.T) {
	path := findFFmpeg()

	// Test should pass if ffmpeg is found in any location, or return empty if not installed
	if path != "" {
		// Verify the path actually exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("findFFmpeg returned path %q but file does not exist", path)
		}
	}
}

func TestFindFFprobe(t *testing.T) {
	path := findFFprobe()

	// Test should pass if ffprobe is found in any location, or return empty if not installed
	if path != "" {
		// Verify the path actually exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("findFFprobe returned path %q but file does not exist", path)
		}
	}
}

func TestMergerAvailable(t *testing.T) {
	available := MergerAvailable()

	// Check if ffmpeg is actually available via exec.LookPath or known paths
	_, pathErr := exec.LookPath("ffmpeg")
	_, staticErr := os.Stat("/ffmpeg")

	expectedAvailable := pathErr == nil || staticErr == nil

	if available != expectedAvailable {
		t.Errorf("MergerAvailable() = %v, expected %v based on system state", available, expectedAvailable)
	}
}

func TestNewMerger(t *testing.T) {
	m := NewMerger()

	if m == nil {
		t.Fatal("NewMerger returned nil")
	}

	if m.timeout != 30*time.Minute {
		t.Errorf("default timeout = %v, want 30m", m.timeout)
	}

	// ffmpegPath should be set (either found or fallback)
	if m.ffmpegPath == "" {
		t.Error("ffmpegPath should not be empty")
	}
}

func TestNewMerger_WithOptions(t *testing.T) {
	customPath := "/custom/ffmpeg"
	customTimeout := 10 * time.Minute

	m := NewMerger(
		WithFFmpegPath(customPath),
		WithMergeTimeout(customTimeout),
	)

	if m.ffmpegPath != customPath {
		t.Errorf("ffmpegPath = %q, want %q", m.ffmpegPath, customPath)
	}

	if m.timeout != customTimeout {
		t.Errorf("timeout = %v, want %v", m.timeout, customTimeout)
	}
}

func TestMerger_MergeWithCodecCopy_InvalidInputs(t *testing.T) {
	m := NewMerger()
	ctx := context.Background()

	t.Run("nonexistent video file", func(t *testing.T) {
		result := m.MergeWithCodecCopy(ctx, "/nonexistent/video.mp4", "/nonexistent/audio.m4a", "/tmp/output.mp4")
		if result.Error == nil {
			t.Error("expected error for nonexistent video file")
		}
		if result.Success {
			t.Error("Success should be false for failed merge")
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		result := m.MergeWithCodecCopy(ctx, "/tmp/video.mp4", "/tmp/audio.m4a", "/tmp/output.mp4")
		if result.Error == nil {
			t.Error("expected error for cancelled context")
		}
	})
}

func TestMerger_ConcatSegments_EmptyList(t *testing.T) {
	m := NewMerger()
	ctx := context.Background()

	result := m.ConcatSegments(ctx, []string{}, "/tmp/output.mp4")
	if result.Error == nil {
		t.Error("expected error for empty segment list")
	}
}

func TestMerger_ConcatSegments_NonexistentFiles(t *testing.T) {
	m := NewMerger()
	ctx := context.Background()

	result := m.ConcatSegments(ctx, []string{"/nonexistent/seg1.ts", "/nonexistent/seg2.ts"}, "/tmp/output.mp4")
	if result.Error == nil {
		t.Error("expected error for nonexistent segment files")
	}
}

func TestGetVideoInfo_NoFFprobe(t *testing.T) {
	// This test checks error handling when ffprobe is not available
	// Skip if ffprobe is actually available
	if findFFprobe() != "" {
		t.Skip("ffprobe is available, skipping error test")
	}

	_, _, _, err := GetVideoInfo("/tmp/test.mp4")
	if err == nil {
		t.Error("expected error when ffprobe is not available")
	}
}

func TestGetVideoInfo_NonexistentFile(t *testing.T) {
	// Skip if ffprobe is not available
	if findFFprobe() == "" {
		t.Skip("ffprobe not available")
	}

	_, _, _, err := GetVideoInfo("/nonexistent/video.mp4")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestGetMediaDuration_NoFFprobe(t *testing.T) {
	// This test checks error handling when ffprobe is not available
	// Skip if ffprobe is actually available
	if findFFprobe() != "" {
		t.Skip("ffprobe is available, skipping error test")
	}

	_, err := GetMediaDuration("/tmp/test.mp4")
	if err == nil {
		t.Error("expected error when ffprobe is not available")
	}
}

func TestGetMediaDuration_NonexistentFile(t *testing.T) {
	// Skip if ffprobe is not available
	if findFFprobe() == "" {
		t.Skip("ffprobe not available")
	}

	_, err := GetMediaDuration("/nonexistent/video.mp4")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestCleanupTempFiles(t *testing.T) {
	// Create a temp directory with some files
	tempDir, err := os.MkdirTemp("", "cleanup_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := []string{"test1.tmp", "test2.tmp", "keep.txt"}
	for _, f := range testFiles {
		if err := os.WriteFile(filepath.Join(tempDir, f), []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	// Clean up .tmp files
	err = CleanupTempFiles(tempDir, "*.tmp")
	if err != nil {
		t.Errorf("CleanupTempFiles returned error: %v", err)
	}

	// Verify .tmp files are removed
	for _, f := range []string{"test1.tmp", "test2.tmp"} {
		if _, err := os.Stat(filepath.Join(tempDir, f)); !os.IsNotExist(err) {
			t.Errorf("file %s should have been removed", f)
		}
	}

	// Verify keep.txt still exists
	if _, err := os.Stat(filepath.Join(tempDir, "keep.txt")); os.IsNotExist(err) {
		t.Error("keep.txt should not have been removed")
	}
}

func TestCleanupTempFiles_NonexistentDir(t *testing.T) {
	// Should not return error for nonexistent directory
	err := CleanupTempFiles("/nonexistent/directory", "*.tmp")
	if err != nil {
		t.Errorf("CleanupTempFiles should not error on nonexistent dir: %v", err)
	}
}

func TestVerifyFFmpegPath_Invalid(t *testing.T) {
	err := VerifyFFmpegPath("/nonexistent/ffmpeg")
	if err == nil {
		t.Error("expected error for invalid ffmpeg path")
	}
}

func TestMergeResult_Fields(t *testing.T) {
	result := &MergeResult{
		OutputPath: "/tmp/output.mp4",
		Duration:   5 * time.Second,
		FileSize:   1024,
		Success:    true,
		Error:      nil,
	}

	if result.OutputPath != "/tmp/output.mp4" {
		t.Errorf("OutputPath = %q, want /tmp/output.mp4", result.OutputPath)
	}

	if result.Duration != 5*time.Second {
		t.Errorf("Duration = %v, want 5s", result.Duration)
	}

	if result.FileSize != 1024 {
		t.Errorf("FileSize = %d, want 1024", result.FileSize)
	}

	if !result.Success {
		t.Error("Success should be true")
	}
}

func TestMerger_ExtractAudio_InvalidInput(t *testing.T) {
	m := NewMerger()
	ctx := context.Background()

	result := m.ExtractAudio(ctx, "/nonexistent/video.mp4", "/tmp/audio.m4a")
	if result.Error == nil {
		t.Error("expected error for nonexistent input file")
	}
	if result.Success {
		t.Error("Success should be false")
	}
}

func TestMerger_ConvertToMP4_InvalidInput(t *testing.T) {
	m := NewMerger()
	ctx := context.Background()

	result := m.ConvertToMP4(ctx, "/nonexistent/input.webm", "/tmp/output.mp4")
	if result.Error == nil {
		t.Error("expected error for nonexistent input file")
	}
	if result.Success {
		t.Error("Success should be false")
	}
}

func TestMerger_RemuxToMP4_InvalidInput(t *testing.T) {
	m := NewMerger()
	ctx := context.Background()

	result := m.RemuxToMP4(ctx, "/nonexistent/input.mkv", "/tmp/output.mp4")
	if result.Error == nil {
		t.Error("expected error for nonexistent input file")
	}
	if result.Success {
		t.Error("Success should be false")
	}
}

func TestMerger_Merge_Timeout(t *testing.T) {
	// Skip if ffmpeg is not available
	if !MergerAvailable() {
		t.Skip("ffmpeg not available")
	}

	// Create merger with very short timeout
	m := NewMerger(WithMergeTimeout(1 * time.Nanosecond))
	ctx := context.Background()

	// This should timeout immediately
	result := m.Merge(ctx, "/nonexistent/video.mp4", "/nonexistent/audio.m4a", "/tmp/output.mp4")
	// The error might be from timeout or from nonexistent files - either is acceptable
	if result.Success {
		t.Error("Success should be false for failed merge")
	}
}

func TestBuildMergeArgs(t *testing.T) {
	m := NewMerger()

	args := m.buildMergeArgs("/input/video.mp4", "/input/audio.m4a", "/output/merged.mp4")

	// Verify required arguments are present
	hasY := false
	hasInputVideo := false
	hasInputAudio := false
	hasOutput := false
	hasFaststart := false

	for i, arg := range args {
		if arg == "-y" {
			hasY = true
		}
		if arg == "-i" && i+1 < len(args) {
			if args[i+1] == "/input/video.mp4" {
				hasInputVideo = true
			}
			if args[i+1] == "/input/audio.m4a" {
				hasInputAudio = true
			}
		}
		if arg == "/output/merged.mp4" {
			hasOutput = true
		}
		if arg == "+faststart" {
			hasFaststart = true
		}
	}

	if !hasY {
		t.Error("missing -y flag for overwrite")
	}
	if !hasInputVideo {
		t.Error("missing video input")
	}
	if !hasInputAudio {
		t.Error("missing audio input")
	}
	if !hasOutput {
		t.Error("missing output path")
	}
	if !hasFaststart {
		t.Error("missing faststart flag")
	}
}
