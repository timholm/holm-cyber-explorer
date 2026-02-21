package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/timholm/ytarchive/internal/types"
)

// DefaultStoragePath is the default path for the iSCSI PVC mount
const DefaultStoragePath = "/data"

// Manager handles all storage operations for the YouTube archiver
type Manager struct {
	basePath string
	mu       sync.RWMutex
}

// Channel represents YouTube channel metadata for storage (extends types.Channel)
type Channel struct {
	ID          string    `json:"id"`
	YouTubeID   string    `json:"youtube_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CustomURL   string    `json:"custom_url,omitempty"`
	Thumbnail   string    `json:"thumbnail,omitempty"`
	Banner      string    `json:"banner,omitempty"`
	VideoCount  int       `json:"video_count"`
	ViewCount   int64     `json:"view_count"`
	Subscribers int64     `json:"subscribers"`
	JoinedDate  string    `json:"joined_date,omitempty"`
	Country     string    `json:"country,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Video represents YouTube video metadata for storage (extends types.Video)
type Video struct {
	ID           string    `json:"id"`
	YouTubeID    string    `json:"youtube_id"`
	ChannelID    string    `json:"channel_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	Duration     int       `json:"duration"` // in seconds
	ViewCount    int64     `json:"view_count"`
	LikeCount    int64     `json:"like_count"`
	CommentCount int64     `json:"comment_count"`
	UploadDate   string    `json:"upload_date"`
	Thumbnail    string    `json:"thumbnail,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
	Categories   []string  `json:"categories,omitempty"`
	FilePath     string    `json:"file_path,omitempty"`
	FileSize     int64     `json:"file_size,omitempty"`
	Checksum     string    `json:"checksum,omitempty"`
	Status       string    `json:"status"` // pending, downloading, downloaded, error
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Ensure storage types can be converted to canonical types
var (
	_ = convertChannelToCanonical
	_ = convertVideoToCanonical
)

// convertChannelToCanonical converts a storage Channel to the canonical types.Channel
func convertChannelToCanonical(c *Channel) *types.Channel {
	return &types.Channel{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		URL:         c.CustomURL,
		AvatarURL:   c.Thumbnail,
		BannerURL:   c.Banner,
		VideoCount:  c.VideoCount,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

// convertVideoToCanonical converts a storage Video to the canonical types.Video
func convertVideoToCanonical(v *Video) *types.Video {
	return &types.Video{
		ID:           v.ID,
		ChannelID:    v.ChannelID,
		Title:        v.Title,
		Description:  v.Description,
		Duration:     v.Duration,
		UploadDate:   v.UploadDate,
		ThumbnailURL: v.Thumbnail,
		ViewCount:    v.ViewCount,
		Status:       v.Status,
		FilePath:     v.FilePath,
		FileSize:     v.FileSize,
		Checksum:     v.Checksum,
		CreatedAt:    v.CreatedAt,
		UpdatedAt:    v.UpdatedAt,
	}
}

// ToCanonical converts a storage Channel to the canonical types.Channel
func (c *Channel) ToCanonical() *types.Channel {
	return convertChannelToCanonical(c)
}

// ToCanonical converts a storage Video to the canonical types.Video
func (v *Video) ToCanonical() *types.Video {
	return convertVideoToCanonical(v)
}

// NewManager creates a new storage manager
func NewManager(basePath string) *Manager {
	if basePath == "" {
		basePath = getStoragePath()
	}
	return &Manager{
		basePath: basePath,
	}
}

// getStoragePath returns the storage path from env or default
func getStoragePath() string {
	if path := os.Getenv("STORAGE_PATH"); path != "" {
		return path
	}
	return DefaultStoragePath
}

// InitStorage creates the base directory structure
func (m *Manager) InitStorage() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	directories := []string{
		filepath.Join(m.basePath, "channels"),
		filepath.Join(m.basePath, "queue"),
		filepath.Join(m.basePath, "logs"),
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// GetBasePath returns the base storage path
func (m *Manager) GetBasePath() string {
	return m.basePath
}

// GetChannelPath returns the path for a channel's directory
func (m *Manager) GetChannelPath(channelID string) string {
	return filepath.Join(m.basePath, "channels", channelID)
}

// GetVideoPath returns the path for a video's directory
func (m *Manager) GetVideoPath(channelID, videoID string) string {
	return filepath.Join(m.GetChannelPath(channelID), "videos", videoID)
}

// GetVideoFilePath returns the full path to a video file
func (m *Manager) GetVideoFilePath(channelID, videoID string) string {
	return filepath.Join(m.GetVideoPath(channelID, videoID), "video.mp4")
}

// GetMetadataPath returns the path for a video's metadata file
func (m *Manager) GetMetadataPath(channelID, videoID string) string {
	return filepath.Join(m.GetVideoPath(channelID, videoID), "metadata.json")
}

// GetThumbnailPath returns the path for a video's thumbnail
func (m *Manager) GetThumbnailPath(channelID, videoID string) string {
	return filepath.Join(m.GetVideoPath(channelID, videoID), "thumbnail.webp")
}

// GetSubtitlesPath returns the path for a video's subtitles directory
func (m *Manager) GetSubtitlesPath(channelID, videoID string) string {
	return filepath.Join(m.GetVideoPath(channelID, videoID), "subtitles")
}

// GetChannelInfoPath returns the path to channel.json
func (m *Manager) GetChannelInfoPath(channelID string) string {
	return filepath.Join(m.GetChannelPath(channelID), "channel.json")
}

// GetChannelAvatarPath returns the path to channel avatar
func (m *Manager) GetChannelAvatarPath(channelID string) string {
	return filepath.Join(m.GetChannelPath(channelID), "avatar.jpg")
}

// GetChannelDBPath returns the path to channel's metadata database
func (m *Manager) GetChannelDBPath(channelID string) string {
	return filepath.Join(m.GetChannelPath(channelID), "metadata.db")
}

// GetLogPath returns the path for a specific date's log directory
func (m *Manager) GetLogPath(date time.Time) string {
	return filepath.Join(m.basePath, "logs", date.Format("2006-01-02"))
}

// GetQueuePath returns the queue directory path
func (m *Manager) GetQueuePath() string {
	return filepath.Join(m.basePath, "queue")
}

// SaveChannelInfo saves channel metadata to channel.json
func (m *Manager) SaveChannelInfo(channel *Channel) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	channelPath := m.GetChannelPath(channel.YouTubeID)
	if err := os.MkdirAll(channelPath, 0755); err != nil {
		return fmt.Errorf("failed to create channel directory: %w", err)
	}

	channel.UpdatedAt = time.Now()
	if channel.CreatedAt.IsZero() {
		channel.CreatedAt = channel.UpdatedAt
	}

	data, err := json.MarshalIndent(channel, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal channel info: %w", err)
	}

	infoPath := m.GetChannelInfoPath(channel.YouTubeID)
	if err := os.WriteFile(infoPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write channel info: %w", err)
	}

	return nil
}

// LoadChannelInfo loads channel metadata from channel.json
func (m *Manager) LoadChannelInfo(channelID string) (*Channel, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	infoPath := m.GetChannelInfoPath(channelID)
	data, err := os.ReadFile(infoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read channel info: %w", err)
	}

	var channel Channel
	if err := json.Unmarshal(data, &channel); err != nil {
		return nil, fmt.Errorf("failed to unmarshal channel info: %w", err)
	}

	return &channel, nil
}

// SaveVideoMetadata saves video metadata to metadata.json
func (m *Manager) SaveVideoMetadata(channelID string, video *Video) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	videoPath := filepath.Join(m.GetChannelPath(channelID), "videos", video.YouTubeID)
	if err := os.MkdirAll(videoPath, 0755); err != nil {
		return fmt.Errorf("failed to create video directory: %w", err)
	}

	video.ChannelID = channelID
	video.UpdatedAt = time.Now()
	if video.CreatedAt.IsZero() {
		video.CreatedAt = video.UpdatedAt
	}

	data, err := json.MarshalIndent(video, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal video metadata: %w", err)
	}

	metadataPath := filepath.Join(videoPath, "metadata.json")
	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write video metadata: %w", err)
	}

	return nil
}

// LoadVideoMetadata loads video metadata from metadata.json
func (m *Manager) LoadVideoMetadata(channelID, videoID string) (*Video, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metadataPath := m.GetMetadataPath(channelID, videoID)
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read video metadata: %w", err)
	}

	var video Video
	if err := json.Unmarshal(data, &video); err != nil {
		return nil, fmt.Errorf("failed to unmarshal video metadata: %w", err)
	}

	return &video, nil
}

// FileExists checks if a file exists at the given path
func (m *Manager) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetFileSize returns the size of a file in bytes
func (m *Manager) GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("failed to stat file: %w", err)
	}
	return info.Size(), nil
}

// ListChannels returns a list of all channel IDs in storage
func (m *Manager) ListChannels() ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	channelsPath := filepath.Join(m.basePath, "channels")
	entries, err := os.ReadDir(channelsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read channels directory: %w", err)
	}

	var channels []string
	for _, entry := range entries {
		if entry.IsDir() {
			channels = append(channels, entry.Name())
		}
	}

	return channels, nil
}

// ListVideos returns a list of all video IDs for a channel
func (m *Manager) ListVideos(channelID string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	videosPath := filepath.Join(m.GetChannelPath(channelID), "videos")
	entries, err := os.ReadDir(videosPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read videos directory: %w", err)
	}

	var videos []string
	for _, entry := range entries {
		if entry.IsDir() {
			videos = append(videos, entry.Name())
		}
	}

	return videos, nil
}

// DeleteVideo removes a video directory and all its contents
func (m *Manager) DeleteVideo(channelID, videoID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	videoPath := m.GetVideoPath(channelID, videoID)
	if err := os.RemoveAll(videoPath); err != nil {
		return fmt.Errorf("failed to delete video directory: %w", err)
	}

	return nil
}

// DeleteChannel removes a channel directory and all its contents
func (m *Manager) DeleteChannel(channelID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	channelPath := m.GetChannelPath(channelID)
	if err := os.RemoveAll(channelPath); err != nil {
		return fmt.Errorf("failed to delete channel directory: %w", err)
	}

	return nil
}
