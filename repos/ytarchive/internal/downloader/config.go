package downloader

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Config holds the configuration for video downloads
type Config struct {
	// OutputPath is the base directory for downloaded videos
	OutputPath string

	// MaxHeight is the maximum video height (default 1080)
	MaxHeight int

	// SubtitleLangs specifies which subtitle languages to download
	SubtitleLangs []string

	// WriteThumbnail enables thumbnail download
	WriteThumbnail bool

	// WriteInfoJSON enables metadata JSON dump
	WriteInfoJSON bool

	// WriteSubtitles enables subtitle download
	WriteSubtitles bool

	// MergeOutputFormat specifies the output format after merging
	MergeOutputFormat string

	// Retries is the number of download retries
	Retries int

	// RetryDelays contains the delay in seconds for each retry attempt
	RetryDelays []int

	// PreferCombinedStream when true, prefers combined audio+video streams (format 18/22)
	// to avoid needing ffmpeg for merging
	PreferCombinedStream bool

	// ChunkSize is the size of each download chunk for progress tracking (default 1MB)
	ChunkSize int64

	// ConnectionTimeout is the timeout for establishing connections (seconds)
	ConnectionTimeout int

	// ReadTimeout is the timeout for read operations (seconds)
	ReadTimeout int
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig(outputPath string) *Config {
	return &Config{
		OutputPath:           outputPath,
		MaxHeight:            getEnvInt("MAX_VIDEO_HEIGHT", 4320), // Default to highest available (8K)
		SubtitleLangs:        getEnvSlice("SUBTITLE_LANGS", []string{"en"}),
		WriteThumbnail:       getEnvBool("WRITE_THUMBNAIL", true),
		WriteInfoJSON:        getEnvBool("WRITE_INFO_JSON", true),
		WriteSubtitles:       getEnvBool("WRITE_SUBTITLES", true),
		MergeOutputFormat:    getEnvString("MERGE_OUTPUT_FORMAT", "mp4"),
		Retries:              getEnvInt("DOWNLOAD_RETRIES", 3),
		RetryDelays:          []int{5, 15, 45},                            // Exponential backoff: 5s, 15s, 45s
		PreferCombinedStream: getEnvBool("PREFER_COMBINED_STREAM", false), // False to prefer highest quality
		ChunkSize:            int64(getEnvInt("CHUNK_SIZE", 1024*1024)),   // 1MB default
		ConnectionTimeout:    getEnvInt("CONNECTION_TIMEOUT", 30),
		ReadTimeout:          getEnvInt("READ_TIMEOUT", 60),
	}
}

// ConfigOption is a function that configures a Config
type ConfigOption func(*Config)

// WithMaxHeight sets the maximum video height
func WithMaxHeight(height int) ConfigOption {
	return func(c *Config) {
		c.MaxHeight = height
	}
}

// WithRetries sets the number of download retries
func WithRetries(retries int) ConfigOption {
	return func(c *Config) {
		c.Retries = retries
	}
}

// WithPreferCombinedStream sets whether to prefer combined streams
func WithPreferCombinedStream(prefer bool) ConfigOption {
	return func(c *Config) {
		c.PreferCombinedStream = prefer
	}
}

// WithWriteThumbnail sets whether to download thumbnails
func WithWriteThumbnail(write bool) ConfigOption {
	return func(c *Config) {
		c.WriteThumbnail = write
	}
}

// WithWriteInfoJSON sets whether to write metadata JSON
func WithWriteInfoJSON(write bool) ConfigOption {
	return func(c *Config) {
		c.WriteInfoJSON = write
	}
}

// NewConfig creates a new Config with the given options
func NewConfig(outputPath string, opts ...ConfigOption) *Config {
	cfg := DefaultConfig(outputPath)
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// getEnvInt returns an integer from an environment variable or a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getEnvBool returns a boolean from an environment variable or a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		value = strings.ToLower(value)
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}

// getEnvString returns a string from an environment variable or a default value
func getEnvString(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvSlice returns a slice from an environment variable (comma-separated) or a default value
func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

// FormatString returns a format selection string for quality preferences
func (c *Config) FormatString() string {
	if c.PreferCombinedStream {
		// Prefer combined streams (no merging needed), but still prefer higher quality first
		return "22/18/bestvideo+bestaudio/best"
	}
	// Prefer best quality, may need merging
	return "bestvideo+bestaudio/best"
}

// OutputTemplate returns the output template for a given video ID
// Output format: {outputPath}/{videoID}/video.{ext}
func (c *Config) OutputTemplate(videoID string) string {
	return filepath.Join(c.OutputPath, videoID, "video.%(ext)s")
}

// VideoDir returns the video directory path for a given channel and video ID
func (c *Config) VideoDir(channelID, videoID string) string {
	if channelID != "" {
		return filepath.Join(c.OutputPath, "channels", channelID, "videos", videoID)
	}
	return filepath.Join(c.OutputPath, videoID)
}

// SubtitleLangsString returns the subtitle languages as a comma-separated string
func (c *Config) SubtitleLangsString() string {
	if len(c.SubtitleLangs) == 0 {
		return "en"
	}
	return strings.Join(c.SubtitleLangs, ",")
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.OutputPath == "" {
		return ErrInvalidConfig("output path is required")
	}
	if c.MaxHeight <= 0 {
		return ErrInvalidConfig("max height must be positive")
	}
	if c.Retries <= 0 {
		return ErrInvalidConfig("retries must be positive")
	}
	return nil
}

// ErrInvalidConfig represents a configuration error
type ErrInvalidConfig string

func (e ErrInvalidConfig) Error() string {
	return "invalid config: " + string(e)
}

// SanitizeFilename removes or replaces characters that are invalid in filenames
func SanitizeFilename(name string) string {
	// Remove invalid filename characters
	invalidChars := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)
	name = invalidChars.ReplaceAllString(name, "")

	// Replace spaces and multiple consecutive special chars with single dash
	name = strings.TrimSpace(name)
	name = regexp.MustCompile(`[\s_]+`).ReplaceAllString(name, "-")

	// Convert to lowercase for consistency
	name = strings.ToLower(name)

	// Remove any leading/trailing dashes
	name = strings.Trim(name, "-")

	// Truncate to reasonable length (max 100 chars for title portion)
	if len(name) > 100 {
		name = name[:100]
		// Don't cut in the middle of a word if possible
		if lastDash := strings.LastIndex(name, "-"); lastDash > 80 {
			name = name[:lastDash]
		}
	}

	return name
}

// SanitizeChannelName converts a channel name to a URL-friendly slug
func SanitizeChannelName(name string) string {
	// Remove special characters, keep alphanumeric and spaces
	result := strings.Builder{}
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(unicode.ToLower(r))
		} else if unicode.IsSpace(r) || r == '-' || r == '_' {
			result.WriteRune('-')
		}
	}

	name = result.String()
	// Remove consecutive dashes
	name = regexp.MustCompile(`-+`).ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")

	return name
}

// GenerateVideoFilename creates a filename in the format: {channel}-ep{number}-{title}.{ext}
func GenerateVideoFilename(channelName string, episodeNumber int, title string, ext string) string {
	sanitizedChannel := SanitizeChannelName(channelName)
	sanitizedTitle := SanitizeFilename(title)

	// Format: channelname-ep00001-video-title.mp4
	return filepath.Clean(strings.Join([]string{
		sanitizedChannel,
		"-ep",
		padEpisodeNumber(episodeNumber),
		"-",
		sanitizedTitle,
		".",
		ext,
	}, ""))
}

// padEpisodeNumber pads the episode number with leading zeros (5 digits)
func padEpisodeNumber(num int) string {
	s := strconv.Itoa(num)
	for len(s) < 5 {
		s = "0" + s
	}
	return s
}
