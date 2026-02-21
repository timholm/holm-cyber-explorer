package storage

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// IntegrityChecker handles file verification operations
type IntegrityChecker struct {
	manager *Manager
}

// NewIntegrityChecker creates a new integrity checker
func NewIntegrityChecker(manager *Manager) *IntegrityChecker {
	return &IntegrityChecker{
		manager: manager,
	}
}

// VideoVerificationResult contains the result of video verification
type VideoVerificationResult struct {
	Path        string `json:"path"`
	Valid       bool   `json:"valid"`
	HasAudio    bool   `json:"has_audio"`
	HasVideo    bool   `json:"has_video"`
	Duration    string `json:"duration,omitempty"`
	Format      string `json:"format,omitempty"`
	ErrorMsg    string `json:"error_msg,omitempty"`
	FileSize    int64  `json:"file_size"`
	IsCorrupted bool   `json:"is_corrupted"`
}

// VerifyVideo checks if a file is a valid video using ffprobe
// Returns true if the video file appears to be valid and playable
func (ic *IntegrityChecker) VerifyVideo(path string) (*VideoVerificationResult, error) {
	result := &VideoVerificationResult{
		Path:  path,
		Valid: false,
	}

	// Check if file exists
	info, err := os.Stat(path)
	if err != nil {
		result.ErrorMsg = fmt.Sprintf("file not found: %v", err)
		return result, nil
	}
	result.FileSize = info.Size()

	// Check if file is too small to be a valid video
	if info.Size() < 1024 {
		result.ErrorMsg = "file too small to be a valid video"
		result.IsCorrupted = true
		return result, nil
	}

	// Try to verify with ffprobe if available
	if ffprobePath, err := exec.LookPath("ffprobe"); err == nil {
		return ic.verifyWithFFprobe(path, result, ffprobePath)
	}

	// Fallback: basic header check for common video formats
	return ic.verifyByHeader(path, result)
}

// verifyWithFFprobe uses ffprobe to thoroughly check video integrity
func (ic *IntegrityChecker) verifyWithFFprobe(path string, result *VideoVerificationResult, ffprobePath string) (*VideoVerificationResult, error) {
	// Run ffprobe to get stream information
	cmd := exec.Command(ffprobePath,
		"-v", "error",
		"-show_entries", "format=duration,format_name:stream=codec_type",
		"-of", "csv=p=0",
		path,
	)

	output, err := cmd.Output()
	if err != nil {
		// If ffprobe exits with error, the video is likely corrupted
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ErrorMsg = fmt.Sprintf("ffprobe error: %s", string(exitErr.Stderr))
			result.IsCorrupted = true
			return result, nil
		}
		result.ErrorMsg = fmt.Sprintf("failed to run ffprobe: %v", err)
		return result, nil
	}

	// Parse output
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) == 0 {
			continue
		}

		switch {
		case parts[0] == "video":
			result.HasVideo = true
		case parts[0] == "audio":
			result.HasAudio = true
		case len(parts) >= 2:
			// Format line: duration,format_name
			result.Duration = parts[0]
			result.Format = parts[1]
		}
	}

	// Additional integrity check - try to read some frames
	checkCmd := exec.Command(ffprobePath,
		"-v", "error",
		"-read_intervals", "%+5", // Read first 5 seconds
		"-show_entries", "frame=pkt_pts_time",
		"-of", "csv=p=0",
		path,
	)

	if err := checkCmd.Run(); err != nil {
		result.ErrorMsg = "video appears corrupted (unable to read frames)"
		result.IsCorrupted = true
		return result, nil
	}

	result.Valid = result.HasVideo || result.HasAudio
	return result, nil
}

// verifyByHeader performs basic header validation for video files
func (ic *IntegrityChecker) verifyByHeader(path string, result *VideoVerificationResult) (*VideoVerificationResult, error) {
	file, err := os.Open(path)
	if err != nil {
		result.ErrorMsg = fmt.Sprintf("unable to open file: %v", err)
		return result, nil
	}
	defer file.Close()

	// Read first 12 bytes for magic number detection
	header := make([]byte, 12)
	n, err := file.Read(header)
	if err != nil || n < 12 {
		result.ErrorMsg = "unable to read file header"
		result.IsCorrupted = true
		return result, nil
	}

	// Check for common video format signatures
	switch {
	case isMP4(header):
		result.Format = "mp4"
		result.Valid = true
	case isWebM(header):
		result.Format = "webm"
		result.Valid = true
	case isMKV(header):
		result.Format = "mkv"
		result.Valid = true
	case isAVI(header):
		result.Format = "avi"
		result.Valid = true
	case isFLV(header):
		result.Format = "flv"
		result.Valid = true
	default:
		result.ErrorMsg = "unknown or invalid video format"
	}

	return result, nil
}

// isMP4 checks for MP4/MOV file signature
func isMP4(header []byte) bool {
	// MP4 files have 'ftyp' at offset 4
	if len(header) >= 8 {
		return bytes.Equal(header[4:8], []byte("ftyp"))
	}
	return false
}

// isWebM checks for WebM file signature
func isWebM(header []byte) bool {
	// WebM starts with EBML header
	if len(header) >= 4 {
		return bytes.Equal(header[0:4], []byte{0x1A, 0x45, 0xDF, 0xA3})
	}
	return false
}

// isMKV checks for MKV file signature (same as WebM)
func isMKV(header []byte) bool {
	return isWebM(header)
}

// isAVI checks for AVI file signature
func isAVI(header []byte) bool {
	if len(header) >= 12 {
		return bytes.Equal(header[0:4], []byte("RIFF")) && bytes.Equal(header[8:12], []byte("AVI "))
	}
	return false
}

// isFLV checks for FLV file signature
func isFLV(header []byte) bool {
	if len(header) >= 3 {
		return bytes.Equal(header[0:3], []byte("FLV"))
	}
	return false
}

// CalculateChecksum computes SHA256 hash of a file
func (ic *IntegrityChecker) CalculateChecksum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate checksum: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// CompareChecksums compares two checksum strings (case-insensitive)
func (ic *IntegrityChecker) CompareChecksums(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

// VerifyChecksum verifies a file against an expected checksum
func (ic *IntegrityChecker) VerifyChecksum(path, expectedChecksum string) (bool, error) {
	actualChecksum, err := ic.CalculateChecksum(path)
	if err != nil {
		return false, err
	}
	return ic.CompareChecksums(actualChecksum, expectedChecksum), nil
}

// VerifyVideoWithMetadata verifies a video file against its stored metadata
func (ic *IntegrityChecker) VerifyVideoWithMetadata(channelID, videoID string) (*VideoVerificationResult, error) {
	// Load metadata
	video, err := ic.manager.LoadVideoMetadata(channelID, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to load metadata: %w", err)
	}

	// Verify the video file
	videoPath := ic.manager.GetVideoFilePath(channelID, videoID)
	result, err := ic.VerifyVideo(videoPath)
	if err != nil {
		return nil, err
	}

	// If metadata has checksum, verify it
	if video.Checksum != "" && result.Valid {
		match, err := ic.VerifyChecksum(videoPath, video.Checksum)
		if err != nil {
			result.ErrorMsg = fmt.Sprintf("checksum verification failed: %v", err)
			result.IsCorrupted = true
			result.Valid = false
		} else if !match {
			result.ErrorMsg = "checksum mismatch - file may have been modified"
			result.IsCorrupted = true
			result.Valid = false
		}
	}

	// Verify file size if available
	if video.FileSize > 0 && result.FileSize != video.FileSize {
		result.ErrorMsg = fmt.Sprintf("file size mismatch: expected %d, got %d", video.FileSize, result.FileSize)
		result.IsCorrupted = true
		result.Valid = false
	}

	return result, nil
}

// VerifyAllVideos verifies all videos for a channel
func (ic *IntegrityChecker) VerifyAllVideos(channelID string) (map[string]*VideoVerificationResult, error) {
	videos, err := ic.manager.ListVideos(channelID)
	if err != nil {
		return nil, err
	}

	results := make(map[string]*VideoVerificationResult)
	for _, videoID := range videos {
		videoPath := ic.manager.GetVideoFilePath(channelID, videoID)
		if !ic.manager.FileExists(videoPath) {
			results[videoID] = &VideoVerificationResult{
				Path:     videoPath,
				Valid:    false,
				ErrorMsg: "video file does not exist",
			}
			continue
		}

		result, err := ic.VerifyVideo(videoPath)
		if err != nil {
			results[videoID] = &VideoVerificationResult{
				Path:     videoPath,
				Valid:    false,
				ErrorMsg: fmt.Sprintf("verification error: %v", err),
			}
			continue
		}
		results[videoID] = result
	}

	return results, nil
}

// BatchVerificationResult contains results for batch verification
type BatchVerificationResult struct {
	TotalFiles     int                                 `json:"total_files"`
	ValidFiles     int                                 `json:"valid_files"`
	CorruptedFiles int                                 `json:"corrupted_files"`
	MissingFiles   int                                 `json:"missing_files"`
	Results        map[string]*VideoVerificationResult `json:"results,omitempty"`
}

// VerifyDirectory verifies all video files in a directory
func (ic *IntegrityChecker) VerifyDirectory(dirPath string) (*BatchVerificationResult, error) {
	result := &BatchVerificationResult{
		Results: make(map[string]*VideoVerificationResult),
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}

		if info.IsDir() {
			return nil
		}

		// Check if it's a video file
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".mp4" && ext != ".webm" && ext != ".mkv" && ext != ".avi" && ext != ".flv" {
			return nil
		}

		result.TotalFiles++

		verifyResult, err := ic.VerifyVideo(path)
		if err != nil {
			result.CorruptedFiles++
			return nil
		}

		result.Results[path] = verifyResult
		if verifyResult.Valid {
			result.ValidFiles++
		} else if verifyResult.IsCorrupted {
			result.CorruptedFiles++
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return result, nil
}
