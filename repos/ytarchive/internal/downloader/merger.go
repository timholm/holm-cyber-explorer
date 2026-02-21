package downloader

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/timholm/ytarchive/internal/logging"
)

// Merger handles merging separate audio and video streams using ffmpeg
type Merger struct {
	ffmpegPath string
	timeout    time.Duration
}

// MergerOption configures a Merger
type MergerOption func(*Merger)

// WithFFmpegPath sets a custom path to ffmpeg
func WithFFmpegPath(path string) MergerOption {
	return func(m *Merger) {
		m.ffmpegPath = path
	}
}

// WithMergeTimeout sets the merge operation timeout
func WithMergeTimeout(timeout time.Duration) MergerOption {
	return func(m *Merger) {
		m.timeout = timeout
	}
}

// NewMerger creates a new Merger
func NewMerger(opts ...MergerOption) *Merger {
	ffmpegPath := findFFmpeg()
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg" // fallback, will fail if not in PATH
	}

	m := &Merger{
		ffmpegPath: ffmpegPath,
		timeout:    30 * time.Minute,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// MergeResult contains the result of a merge operation
type MergeResult struct {
	OutputPath string
	Duration   time.Duration
	FileSize   int64
	Success    bool
	Error      error
}

// Merge combines video and audio files into a single output file
func (m *Merger) Merge(ctx context.Context, videoPath, audioPath, outputPath string) *MergeResult {
	startTime := time.Now()
	result := &MergeResult{
		OutputPath: outputPath,
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	// Build ffmpeg command
	args := m.buildMergeArgs(videoPath, audioPath, outputPath)

	logging.Debug("executing ffmpeg merge",
		"video_path", videoPath,
		"audio_path", audioPath,
		"output_path", outputPath,
		"args", strings.Join(args, " "),
	)

	cmd := exec.CommandContext(ctx, m.ffmpegPath, args...)

	// Capture output for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Error = fmt.Errorf("ffmpeg merge failed: %w - output: %s", err, string(output))
		result.Duration = time.Since(startTime)
		return result
	}

	// Get output file info
	info, err := os.Stat(outputPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to stat output file: %w", err)
		result.Duration = time.Since(startTime)
		return result
	}

	result.Success = true
	result.Duration = time.Since(startTime)
	result.FileSize = info.Size()

	logging.Info("merge completed",
		"output_path", outputPath,
		"duration", result.Duration.String(),
		"file_size", FormatFileSize(result.FileSize),
	)

	return result
}

// buildMergeArgs constructs ffmpeg command arguments for merging
func (m *Merger) buildMergeArgs(videoPath, audioPath, outputPath string) []string {
	return []string{
		"-y",            // Overwrite output file
		"-i", videoPath, // Input video
		"-i", audioPath, // Input audio
		"-c:v", "copy", // Copy video codec (no re-encoding)
		"-c:a", "aac", // Encode audio to AAC for compatibility
		"-b:a", "192k", // Audio bitrate
		"-movflags", "+faststart", // Enable fast start for streaming
		"-strict", "experimental",
		outputPath,
	}
}

// MergeWithCodecCopy merges without any re-encoding (fastest but may have compatibility issues)
func (m *Merger) MergeWithCodecCopy(ctx context.Context, videoPath, audioPath, outputPath string) *MergeResult {
	startTime := time.Now()
	result := &MergeResult{
		OutputPath: outputPath,
	}

	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	args := []string{
		"-y",
		"-i", videoPath,
		"-i", audioPath,
		"-c", "copy", // Copy both streams without re-encoding
		"-movflags", "+faststart",
		outputPath,
	}

	cmd := exec.CommandContext(ctx, m.ffmpegPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Error = fmt.Errorf("ffmpeg merge failed: %w - output: %s", err, string(output))
		result.Duration = time.Since(startTime)
		return result
	}

	info, err := os.Stat(outputPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to stat output file: %w", err)
		result.Duration = time.Since(startTime)
		return result
	}

	result.Success = true
	result.Duration = time.Since(startTime)
	result.FileSize = info.Size()

	return result
}

// ConcatSegments concatenates multiple segment files into one
func (m *Merger) ConcatSegments(ctx context.Context, segmentPaths []string, outputPath string) *MergeResult {
	startTime := time.Now()
	result := &MergeResult{
		OutputPath: outputPath,
	}

	if len(segmentPaths) == 0 {
		result.Error = fmt.Errorf("no segments to concatenate")
		return result
	}

	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	// Create a temporary file list for ffmpeg concat demuxer
	listFile := filepath.Join(filepath.Dir(outputPath), ".segments.txt")
	defer os.Remove(listFile)

	var listContent strings.Builder
	for _, seg := range segmentPaths {
		// Escape single quotes in path
		escapedPath := strings.ReplaceAll(seg, "'", "'\\''")
		listContent.WriteString(fmt.Sprintf("file '%s'\n", escapedPath))
	}

	if err := os.WriteFile(listFile, []byte(listContent.String()), 0644); err != nil {
		result.Error = fmt.Errorf("failed to create segment list file: %w", err)
		return result
	}

	args := []string{
		"-y",
		"-f", "concat",
		"-safe", "0",
		"-i", listFile,
		"-c", "copy",
		outputPath,
	}

	cmd := exec.CommandContext(ctx, m.ffmpegPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Error = fmt.Errorf("ffmpeg concat failed: %w - output: %s", err, string(output))
		result.Duration = time.Since(startTime)
		return result
	}

	info, err := os.Stat(outputPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to stat output file: %w", err)
		result.Duration = time.Since(startTime)
		return result
	}

	result.Success = true
	result.Duration = time.Since(startTime)
	result.FileSize = info.Size()

	return result
}

// ConvertToMP4 converts a video file to MP4 format
func (m *Merger) ConvertToMP4(ctx context.Context, inputPath, outputPath string) *MergeResult {
	startTime := time.Now()
	result := &MergeResult{
		OutputPath: outputPath,
	}

	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	args := []string{
		"-y",
		"-i", inputPath,
		"-c:v", "libx264",
		"-preset", "medium",
		"-crf", "23",
		"-c:a", "aac",
		"-b:a", "192k",
		"-movflags", "+faststart",
		outputPath,
	}

	cmd := exec.CommandContext(ctx, m.ffmpegPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Error = fmt.Errorf("ffmpeg convert failed: %w - output: %s", err, string(output))
		result.Duration = time.Since(startTime)
		return result
	}

	info, err := os.Stat(outputPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to stat output file: %w", err)
		result.Duration = time.Since(startTime)
		return result
	}

	result.Success = true
	result.Duration = time.Since(startTime)
	result.FileSize = info.Size()

	return result
}

// RemuxToMP4 remuxes a video to MP4 without re-encoding
func (m *Merger) RemuxToMP4(ctx context.Context, inputPath, outputPath string) *MergeResult {
	startTime := time.Now()
	result := &MergeResult{
		OutputPath: outputPath,
	}

	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	args := []string{
		"-y",
		"-i", inputPath,
		"-c", "copy",
		"-movflags", "+faststart",
		outputPath,
	}

	cmd := exec.CommandContext(ctx, m.ffmpegPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Error = fmt.Errorf("ffmpeg remux failed: %w - output: %s", err, string(output))
		result.Duration = time.Since(startTime)
		return result
	}

	info, err := os.Stat(outputPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to stat output file: %w", err)
		result.Duration = time.Since(startTime)
		return result
	}

	result.Success = true
	result.Duration = time.Since(startTime)
	result.FileSize = info.Size()

	return result
}

// ExtractAudio extracts audio from a video file
func (m *Merger) ExtractAudio(ctx context.Context, inputPath, outputPath string) *MergeResult {
	startTime := time.Now()
	result := &MergeResult{
		OutputPath: outputPath,
	}

	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	args := []string{
		"-y",
		"-i", inputPath,
		"-vn", // No video
		"-c:a", "aac",
		"-b:a", "192k",
		outputPath,
	}

	cmd := exec.CommandContext(ctx, m.ffmpegPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Error = fmt.Errorf("ffmpeg extract audio failed: %w - output: %s", err, string(output))
		result.Duration = time.Since(startTime)
		return result
	}

	info, err := os.Stat(outputPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to stat output file: %w", err)
		result.Duration = time.Since(startTime)
		return result
	}

	result.Success = true
	result.Duration = time.Since(startTime)
	result.FileSize = info.Size()

	return result
}

// VerifyFFmpeg checks if ffmpeg is installed and accessible
func VerifyFFmpeg() error {
	return VerifyFFmpegPath("ffmpeg")
}

// VerifyFFmpegPath checks if ffmpeg is available at the specified path
func VerifyFFmpegPath(path string) error {
	cmd := exec.Command(path, "-version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("ffmpeg not found at '%s': %w", path, err)
	}

	version := strings.Split(string(output), "\n")[0]
	logging.Info("found ffmpeg", "version", version)
	return nil
}

// GetVideoInfo gets information about a video file using ffprobe
func GetVideoInfo(videoPath string) (duration float64, width, height int, err error) {
	ffprobePath := findFFprobe()
	if ffprobePath == "" {
		return 0, 0, 0, fmt.Errorf("ffprobe not found")
	}
	cmd := exec.Command(ffprobePath,
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height,duration",
		"-of", "csv=p=0",
		videoPath,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, 0, 0, fmt.Errorf("ffprobe failed: %w", err)
	}

	// Parse output: width,height,duration
	parts := strings.Split(strings.TrimSpace(string(output)), ",")
	if len(parts) >= 2 {
		fmt.Sscanf(parts[0], "%d", &width)
		fmt.Sscanf(parts[1], "%d", &height)
	}
	if len(parts) >= 3 {
		fmt.Sscanf(parts[2], "%f", &duration)
	}

	return duration, width, height, nil
}

// GetMediaDuration returns the duration of a media file in seconds
func GetMediaDuration(path string) (float64, error) {
	ffprobePath := findFFprobe()
	if ffprobePath == "" {
		return 0, fmt.Errorf("ffprobe not found")
	}
	cmd := exec.Command(ffprobePath,
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "csv=p=0",
		path,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe failed: %w", err)
	}

	var duration float64
	fmt.Sscanf(strings.TrimSpace(string(output)), "%f", &duration)
	return duration, nil
}

// CleanupTempFiles removes temporary files created during merge operations
func CleanupTempFiles(dir string, patterns ...string) error {
	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			continue
		}
		for _, match := range matches {
			os.Remove(match)
		}
	}
	return nil
}

// MergerAvailable checks if ffmpeg is available for merging operations
func MergerAvailable() bool {
	return findFFmpeg() != ""
}

// findFFmpeg locates the ffmpeg binary, checking common paths
func findFFmpeg() string {
	// Check PATH first
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		return path
	}
	// Check static-ffmpeg image location
	if _, err := os.Stat("/ffmpeg"); err == nil {
		return "/ffmpeg"
	}
	// Check common install locations
	paths := []string{"/usr/bin/ffmpeg", "/usr/local/bin/ffmpeg"}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

// findFFprobe locates the ffprobe binary, checking common paths
func findFFprobe() string {
	// Check PATH first
	if path, err := exec.LookPath("ffprobe"); err == nil {
		return path
	}
	// Check static-ffmpeg image location
	if _, err := os.Stat("/ffprobe"); err == nil {
		return "/ffprobe"
	}
	// Check common install locations
	paths := []string{"/usr/bin/ffprobe", "/usr/local/bin/ffprobe"}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}
