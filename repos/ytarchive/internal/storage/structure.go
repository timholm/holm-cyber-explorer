package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DirectoryStructure defines the expected directory layout
/*
/data
├── channels/
│   └── {channel_id}/
│       ├── channel.json
│       ├── avatar.jpg
│       ├── metadata.db
│       └── videos/
│           └── {video_id}/
│               ├── video.mp4
│               ├── metadata.json
│               ├── thumbnail.webp
│               └── subtitles/
│                   └── en.vtt
├── queue/
└── logs/
    └── {date}/
*/

// IncompleteDownload represents a partial or failed download
type IncompleteDownload struct {
	VideoID     string    `json:"video_id"`
	ChannelID   string    `json:"channel_id"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	ModTime     time.Time `json:"mod_time"`
	Reason      string    `json:"reason"`
	HasMetadata bool      `json:"has_metadata"`
	HasVideo    bool      `json:"has_video"`
}

// CreateChannelDirectory creates the full directory structure for a channel
func (m *Manager) CreateChannelDirectory(channelID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	channelPath := m.GetChannelPath(channelID)

	// Create channel directory and subdirectories
	directories := []string{
		channelPath,
		filepath.Join(channelPath, "videos"),
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// CreateVideoDirectory creates the directory structure for a video
func (m *Manager) CreateVideoDirectory(channelID, videoID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	videoPath := filepath.Join(m.GetChannelPath(channelID), "videos", videoID)

	// Create video directory and subtitles subdirectory
	directories := []string{
		videoPath,
		filepath.Join(videoPath, "subtitles"),
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// CreateLogDirectory creates a log directory for a specific date
func (m *Manager) CreateLogDirectory(date time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	logPath := m.GetLogPath(date)
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	return nil
}

// CleanupIncomplete finds partial or incomplete downloads for a channel
// Returns a list of incomplete downloads for review - NEVER auto-deletes
func (m *Manager) CleanupIncomplete(channelID string) ([]IncompleteDownload, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var incomplete []IncompleteDownload

	videosPath := filepath.Join(m.GetChannelPath(channelID), "videos")
	entries, err := os.ReadDir(videosPath)
	if err != nil {
		if os.IsNotExist(err) {
			return incomplete, nil
		}
		return nil, fmt.Errorf("failed to read videos directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		videoID := entry.Name()
		videoPath := filepath.Join(videosPath, videoID)

		download := IncompleteDownload{
			VideoID:   videoID,
			ChannelID: channelID,
			Path:      videoPath,
		}

		// Check for metadata.json
		metadataPath := filepath.Join(videoPath, "metadata.json")
		download.HasMetadata = fileExists(metadataPath)

		// Check for video file
		videoFilePath := filepath.Join(videoPath, "video.mp4")
		download.HasVideo = fileExists(videoFilePath)

		// Check for partial downloads (.part files)
		partFiles, _ := filepath.Glob(filepath.Join(videoPath, "*.part"))
		hasPartFiles := len(partFiles) > 0

		// Check for temporary files
		tmpFiles, _ := filepath.Glob(filepath.Join(videoPath, "*.tmp"))
		hasTmpFiles := len(tmpFiles) > 0

		// Determine if this is an incomplete download
		isIncomplete := false
		reasons := []string{}

		if hasPartFiles {
			isIncomplete = true
			reasons = append(reasons, "has .part files")
		}

		if hasTmpFiles {
			isIncomplete = true
			reasons = append(reasons, "has .tmp files")
		}

		if download.HasMetadata && !download.HasVideo {
			isIncomplete = true
			reasons = append(reasons, "metadata exists but no video file")
		}

		if !download.HasMetadata && !download.HasVideo {
			isIncomplete = true
			reasons = append(reasons, "empty video directory")
		}

		// Check if video file is suspiciously small (might be corrupted)
		if download.HasVideo {
			if info, err := os.Stat(videoFilePath); err == nil {
				download.Size = info.Size()
				download.ModTime = info.ModTime()
				// Flag videos smaller than 1KB as potentially corrupt
				if info.Size() < 1024 {
					isIncomplete = true
					reasons = append(reasons, "video file too small (potentially corrupt)")
				}
			}
		}

		if isIncomplete {
			download.Reason = strings.Join(reasons, "; ")

			// Calculate directory size
			if size, err := calculateDirSize(videoPath); err == nil {
				download.Size = size
			}

			// Get mod time from directory if not set
			if download.ModTime.IsZero() {
				if info, err := entry.Info(); err == nil {
					download.ModTime = info.ModTime()
				}
			}

			incomplete = append(incomplete, download)
		}
	}

	return incomplete, nil
}

// CleanupAllIncomplete finds incomplete downloads across all channels
func (m *Manager) CleanupAllIncomplete() ([]IncompleteDownload, error) {
	channels, err := m.ListChannels()
	if err != nil {
		return nil, fmt.Errorf("failed to list channels: %w", err)
	}

	var allIncomplete []IncompleteDownload
	for _, channelID := range channels {
		incomplete, err := m.CleanupIncomplete(channelID)
		if err != nil {
			continue // Skip channels with errors
		}
		allIncomplete = append(allIncomplete, incomplete...)
	}

	return allIncomplete, nil
}

// ValidateChannelStructure checks if a channel directory has the correct structure
func (m *Manager) ValidateChannelStructure(channelID string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var issues []string

	channelPath := m.GetChannelPath(channelID)

	// Check if channel directory exists
	if !fileExists(channelPath) {
		return []string{"channel directory does not exist"}, nil
	}

	// Check for required files
	requiredFiles := map[string]string{
		"channel.json": m.GetChannelInfoPath(channelID),
	}

	for name, path := range requiredFiles {
		if !fileExists(path) {
			issues = append(issues, fmt.Sprintf("missing required file: %s", name))
		}
	}

	// Check for videos directory
	videosPath := filepath.Join(channelPath, "videos")
	if !fileExists(videosPath) {
		issues = append(issues, "missing videos directory")
	}

	return issues, nil
}

// ValidateVideoStructure checks if a video directory has the correct structure
func (m *Manager) ValidateVideoStructure(channelID, videoID string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var issues []string

	videoPath := m.GetVideoPath(channelID, videoID)

	// Check if video directory exists
	if !fileExists(videoPath) {
		return []string{"video directory does not exist"}, nil
	}

	// Check for metadata
	metadataPath := m.GetMetadataPath(channelID, videoID)
	if !fileExists(metadataPath) {
		issues = append(issues, "missing metadata.json")
	}

	// Check for video file
	videoFilePath := filepath.Join(videoPath, "video.mp4")
	if !fileExists(videoFilePath) {
		issues = append(issues, "missing video.mp4")
	}

	// Check for subtitles directory
	subtitlesPath := m.GetSubtitlesPath(channelID, videoID)
	if !fileExists(subtitlesPath) {
		issues = append(issues, "missing subtitles directory")
	}

	return issues, nil
}

// GetDirectoryTree returns a string representation of the directory structure
func (m *Manager) GetDirectoryTree(maxDepth int) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var builder strings.Builder
	builder.WriteString(m.basePath + "\n")

	err := walkDir(m.basePath, "", 0, maxDepth, &builder)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

// walkDir recursively builds a directory tree string
func walkDir(path, prefix string, depth, maxDepth int, builder *strings.Builder) error {
	if maxDepth > 0 && depth >= maxDepth {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for i, entry := range entries {
		isLast := i == len(entries)-1
		connector := "├── "
		if isLast {
			connector = "└── "
		}

		builder.WriteString(prefix + connector + entry.Name())
		if entry.IsDir() {
			builder.WriteString("/")
		}
		builder.WriteString("\n")

		if entry.IsDir() {
			newPrefix := prefix + "│   "
			if isLast {
				newPrefix = prefix + "    "
			}
			err := walkDir(filepath.Join(path, entry.Name()), newPrefix, depth+1, maxDepth, builder)
			if err != nil {
				continue // Skip directories we can't read
			}
		}
	}

	return nil
}

// EnsureDirectoryStructure ensures all base directories exist
func (m *Manager) EnsureDirectoryStructure() error {
	return m.InitStorage()
}

// fileExists is a helper function to check if a file/directory exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// calculateDirSize calculates the total size of a directory
func calculateDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}
