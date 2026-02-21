package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// UsageAnalyzer handles storage usage reporting
type UsageAnalyzer struct {
	manager *Manager
}

// NewUsageAnalyzer creates a new usage analyzer
func NewUsageAnalyzer(manager *Manager) *UsageAnalyzer {
	return &UsageAnalyzer{
		manager: manager,
	}
}

// FileInfo contains information about a file
type FileInfo struct {
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	SizeHuman string    `json:"size_human"`
	ModTime   time.Time `json:"mod_time"`
	ChannelID string    `json:"channel_id,omitempty"`
	VideoID   string    `json:"video_id,omitempty"`
}

// DuplicateSet represents a set of files that appear to be duplicates
type DuplicateSet struct {
	Checksum   string     `json:"checksum"`
	TotalSize  int64      `json:"total_size"`
	Files      []FileInfo `json:"files"`
	WastedSize int64      `json:"wasted_size"` // Size that could be reclaimed
}

// CleanupReport contains recommendations for cleaning up storage
// IMPORTANT: This report only provides recommendations - it NEVER auto-deletes files
type CleanupReport struct {
	GeneratedAt       time.Time `json:"generated_at"`
	TotalStorage      int64     `json:"total_storage"`
	TotalStorageHuman string    `json:"total_storage_human"`

	// Incomplete downloads that could be removed
	IncompleteDownloads []IncompleteDownload `json:"incomplete_downloads"`
	IncompleteSize      int64                `json:"incomplete_size"`

	// Duplicate files
	Duplicates     []DuplicateSet `json:"duplicates"`
	DuplicatesSize int64          `json:"duplicates_size"`

	// Orphaned files (no metadata)
	OrphanedFiles []FileInfo `json:"orphaned_files"`
	OrphanedSize  int64      `json:"orphaned_size"`

	// Temporary files
	TemporaryFiles []FileInfo `json:"temporary_files"`
	TemporarySize  int64      `json:"temporary_size"`

	// Old log files
	OldLogFiles      []FileInfo `json:"old_log_files"`
	OldLogSize       int64      `json:"old_log_size"`
	LogRetentionDays int        `json:"log_retention_days"`

	// Potential savings
	PotentialSavings      int64  `json:"potential_savings"`
	PotentialSavingsHuman string `json:"potential_savings_human"`

	// Warning: This is a report only
	Warning string `json:"warning"`
}

// ChannelUsage contains storage usage for a channel
type ChannelUsage struct {
	ChannelID     string    `json:"channel_id"`
	TotalSize     int64     `json:"total_size"`
	SizeHuman     string    `json:"size_human"`
	VideoCount    int       `json:"video_count"`
	VideoSize     int64     `json:"video_size"`
	MetadataSize  int64     `json:"metadata_size"`
	ThumbnailSize int64     `json:"thumbnail_size"`
	SubtitleSize  int64     `json:"subtitle_size"`
	LastModified  time.Time `json:"last_modified"`
}

// GetTotalUsage returns the total storage usage in bytes
func (ua *UsageAnalyzer) GetTotalUsage() (int64, error) {
	return ua.calculateDirUsage(ua.manager.GetBasePath())
}

// GetChannelUsage returns detailed storage usage for a specific channel
func (ua *UsageAnalyzer) GetChannelUsage(channelID string) (*ChannelUsage, error) {
	channelPath := ua.manager.GetChannelPath(channelID)

	if !ua.manager.FileExists(channelPath) {
		return nil, fmt.Errorf("channel directory does not exist: %s", channelID)
	}

	usage := &ChannelUsage{
		ChannelID: channelID,
	}

	var lastMod time.Time

	err := filepath.Walk(channelPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}

		if info.IsDir() {
			return nil
		}

		size := info.Size()
		usage.TotalSize += size

		if info.ModTime().After(lastMod) {
			lastMod = info.ModTime()
		}

		// Categorize by file type
		ext := filepath.Ext(path)
		switch ext {
		case ".mp4", ".webm", ".mkv", ".avi", ".flv":
			usage.VideoSize += size
			usage.VideoCount++
		case ".json", ".db":
			usage.MetadataSize += size
		case ".jpg", ".jpeg", ".png", ".webp":
			usage.ThumbnailSize += size
		case ".vtt", ".srt", ".ass":
			usage.SubtitleSize += size
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to calculate channel usage: %w", err)
	}

	usage.SizeHuman = formatSize(usage.TotalSize)
	usage.LastModified = lastMod

	return usage, nil
}

// GetAllChannelUsage returns usage for all channels
func (ua *UsageAnalyzer) GetAllChannelUsage() ([]ChannelUsage, error) {
	channels, err := ua.manager.ListChannels()
	if err != nil {
		return nil, err
	}

	var usages []ChannelUsage
	for _, channelID := range channels {
		usage, err := ua.GetChannelUsage(channelID)
		if err != nil {
			continue // Skip channels with errors
		}
		usages = append(usages, *usage)
	}

	// Sort by total size descending
	sort.Slice(usages, func(i, j int) bool {
		return usages[i].TotalSize > usages[j].TotalSize
	})

	return usages, nil
}

// GetLargestFiles returns the largest files in storage
func (ua *UsageAnalyzer) GetLargestFiles(limit int) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.Walk(ua.manager.GetBasePath(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}

		if info.IsDir() {
			return nil
		}

		fileInfo := FileInfo{
			Path:      path,
			Size:      info.Size(),
			SizeHuman: formatSize(info.Size()),
			ModTime:   info.ModTime(),
		}

		// Extract channel and video IDs from path
		relPath, _ := filepath.Rel(ua.manager.GetBasePath(), path)
		parts := filepath.SplitList(relPath)
		if len(parts) >= 2 {
			fileInfo.ChannelID = parts[1] // channels/{channel_id}
		}
		if len(parts) >= 4 {
			fileInfo.VideoID = parts[3] // channels/{channel_id}/videos/{video_id}
		}

		files = append(files, fileInfo)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan files: %w", err)
	}

	// Sort by size descending
	sort.Slice(files, func(i, j int) bool {
		return files[i].Size > files[j].Size
	})

	// Return only top N
	if limit > 0 && len(files) > limit {
		files = files[:limit]
	}

	return files, nil
}

// FindDuplicates finds potential duplicate files based on size and checksum
func (ua *UsageAnalyzer) FindDuplicates() ([]DuplicateSet, error) {
	// Group files by size first (quick filter)
	sizeMap := make(map[int64][]string)

	err := filepath.Walk(ua.manager.GetBasePath(), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// Only check video files
		ext := filepath.Ext(path)
		if ext != ".mp4" && ext != ".webm" && ext != ".mkv" {
			return nil
		}

		sizeMap[info.Size()] = append(sizeMap[info.Size()], path)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan for duplicates: %w", err)
	}

	// For files with same size, compute checksums
	checker := NewIntegrityChecker(ua.manager)
	checksumMap := make(map[string][]FileInfo)

	for size, paths := range sizeMap {
		if len(paths) < 2 {
			continue // No duplicates possible
		}

		for _, path := range paths {
			checksum, err := checker.CalculateChecksum(path)
			if err != nil {
				continue
			}

			info, err := os.Stat(path)
			if err != nil {
				continue
			}

			fileInfo := FileInfo{
				Path:      path,
				Size:      size,
				SizeHuman: formatSize(size),
				ModTime:   info.ModTime(),
			}

			checksumMap[checksum] = append(checksumMap[checksum], fileInfo)
		}
	}

	// Build duplicate sets
	var duplicates []DuplicateSet
	for checksum, files := range checksumMap {
		if len(files) < 2 {
			continue
		}

		totalSize := files[0].Size * int64(len(files))
		wastedSize := files[0].Size * int64(len(files)-1)

		duplicates = append(duplicates, DuplicateSet{
			Checksum:   checksum,
			TotalSize:  totalSize,
			WastedSize: wastedSize,
			Files:      files,
		})
	}

	// Sort by wasted size descending
	sort.Slice(duplicates, func(i, j int) bool {
		return duplicates[i].WastedSize > duplicates[j].WastedSize
	})

	return duplicates, nil
}

// GenerateCleanupReport generates a comprehensive cleanup report
// IMPORTANT: This function NEVER auto-deletes files - it only provides recommendations
func (ua *UsageAnalyzer) GenerateCleanupReport() (*CleanupReport, error) {
	report := &CleanupReport{
		GeneratedAt:      time.Now(),
		LogRetentionDays: 30, // Default retention period
		Warning:          "WARNING: This report is for informational purposes only. Files listed here are RECOMMENDATIONS for cleanup. NO files have been or will be automatically deleted. Manual review and deletion is required.",
	}

	// Get total storage
	totalSize, err := ua.GetTotalUsage()
	if err != nil {
		return nil, err
	}
	report.TotalStorage = totalSize
	report.TotalStorageHuman = formatSize(totalSize)

	// Find incomplete downloads
	incomplete, err := ua.manager.CleanupAllIncomplete()
	if err == nil {
		report.IncompleteDownloads = incomplete
		for _, inc := range incomplete {
			report.IncompleteSize += inc.Size
		}
	}

	// Find duplicates
	duplicates, err := ua.FindDuplicates()
	if err == nil {
		report.Duplicates = duplicates
		for _, dup := range duplicates {
			report.DuplicatesSize += dup.WastedSize
		}
	}

	// Find orphaned files (video files without metadata)
	orphaned, err := ua.findOrphanedFiles()
	if err == nil {
		report.OrphanedFiles = orphaned
		for _, f := range orphaned {
			report.OrphanedSize += f.Size
		}
	}

	// Find temporary files
	tempFiles, err := ua.findTemporaryFiles()
	if err == nil {
		report.TemporaryFiles = tempFiles
		for _, f := range tempFiles {
			report.TemporarySize += f.Size
		}
	}

	// Find old log files
	oldLogs, err := ua.findOldLogFiles(report.LogRetentionDays)
	if err == nil {
		report.OldLogFiles = oldLogs
		for _, f := range oldLogs {
			report.OldLogSize += f.Size
		}
	}

	// Calculate potential savings
	report.PotentialSavings = report.IncompleteSize + report.DuplicatesSize +
		report.OrphanedSize + report.TemporarySize + report.OldLogSize
	report.PotentialSavingsHuman = formatSize(report.PotentialSavings)

	return report, nil
}

// findOrphanedFiles finds video files without corresponding metadata
func (ua *UsageAnalyzer) findOrphanedFiles() ([]FileInfo, error) {
	var orphaned []FileInfo

	channelsPath := filepath.Join(ua.manager.GetBasePath(), "channels")

	err := filepath.Walk(channelsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// Check video files
		ext := filepath.Ext(path)
		if ext != ".mp4" && ext != ".webm" && ext != ".mkv" {
			return nil
		}

		// Check if metadata.json exists in the same directory
		metadataPath := filepath.Join(filepath.Dir(path), "metadata.json")
		if !ua.manager.FileExists(metadataPath) {
			orphaned = append(orphaned, FileInfo{
				Path:      path,
				Size:      info.Size(),
				SizeHuman: formatSize(info.Size()),
				ModTime:   info.ModTime(),
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return orphaned, nil
}

// findTemporaryFiles finds .tmp and .part files
func (ua *UsageAnalyzer) findTemporaryFiles() ([]FileInfo, error) {
	var tempFiles []FileInfo

	err := filepath.Walk(ua.manager.GetBasePath(), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext == ".tmp" || ext == ".part" || ext == ".temp" {
			tempFiles = append(tempFiles, FileInfo{
				Path:      path,
				Size:      info.Size(),
				SizeHuman: formatSize(info.Size()),
				ModTime:   info.ModTime(),
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tempFiles, nil
}

// findOldLogFiles finds log files older than retention period
func (ua *UsageAnalyzer) findOldLogFiles(retentionDays int) ([]FileInfo, error) {
	var oldLogs []FileInfo

	logsPath := filepath.Join(ua.manager.GetBasePath(), "logs")
	if !ua.manager.FileExists(logsPath) {
		return oldLogs, nil
	}

	cutoff := time.Now().AddDate(0, 0, -retentionDays)

	err := filepath.Walk(logsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		if info.ModTime().Before(cutoff) {
			oldLogs = append(oldLogs, FileInfo{
				Path:      path,
				Size:      info.Size(),
				SizeHuman: formatSize(info.Size()),
				ModTime:   info.ModTime(),
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return oldLogs, nil
}

// calculateDirUsage calculates the total size of a directory
func (ua *UsageAnalyzer) calculateDirUsage(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to calculate directory size: %w", err)
	}

	return size, nil
}

// GetUsageSummary returns a high-level summary of storage usage
type UsageSummary struct {
	TotalSize      int64  `json:"total_size"`
	TotalSizeHuman string `json:"total_size_human"`
	ChannelCount   int    `json:"channel_count"`
	VideoCount     int    `json:"video_count"`
	VideoSize      int64  `json:"video_size"`
	MetadataSize   int64  `json:"metadata_size"`
	LogSize        int64  `json:"log_size"`
	QueueSize      int64  `json:"queue_size"`
}

// GetUsageSummary returns a summary of storage usage
func (ua *UsageAnalyzer) GetUsageSummary() (*UsageSummary, error) {
	summary := &UsageSummary{}

	// Get channel usages
	channelUsages, err := ua.GetAllChannelUsage()
	if err != nil {
		return nil, err
	}

	summary.ChannelCount = len(channelUsages)
	for _, cu := range channelUsages {
		summary.VideoCount += cu.VideoCount
		summary.VideoSize += cu.VideoSize
		summary.MetadataSize += cu.MetadataSize
	}

	// Get log size
	logsPath := filepath.Join(ua.manager.GetBasePath(), "logs")
	if ua.manager.FileExists(logsPath) {
		logSize, _ := ua.calculateDirUsage(logsPath)
		summary.LogSize = logSize
	}

	// Get queue size
	queuePath := ua.manager.GetQueuePath()
	if ua.manager.FileExists(queuePath) {
		queueSize, _ := ua.calculateDirUsage(queuePath)
		summary.QueueSize = queueSize
	}

	// Calculate total
	summary.TotalSize, _ = ua.GetTotalUsage()
	summary.TotalSizeHuman = formatSize(summary.TotalSize)

	return summary, nil
}

// formatSize converts bytes to human-readable format
func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
