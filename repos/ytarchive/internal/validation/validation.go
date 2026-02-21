package validation

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// Validation errors
var (
	ErrInvalidURL       = fmt.Errorf("invalid URL format")
	ErrInvalidChannelID = fmt.Errorf("invalid channel ID format")
	ErrInvalidVideoID   = fmt.Errorf("invalid video ID format")
	ErrUnsafeInput      = fmt.Errorf("input contains potentially unsafe characters")
)

// Regular expressions for validation
var (
	// YouTube channel ID: UC followed by 22 alphanumeric characters
	channelIDRegex = regexp.MustCompile(`^UC[a-zA-Z0-9_-]{22}$`)

	// YouTube video ID: 11 alphanumeric characters (including - and _)
	videoIDRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{11}$`)

	// YouTube handle: @ followed by alphanumeric, underscores, dots, hyphens (3-30 chars)
	handleRegex = regexp.MustCompile(`^@[a-zA-Z0-9._-]{3,30}$`)

	// Safe URL characters (no shell injection, XSS, etc.)
	unsafeCharsRegex = regexp.MustCompile(`[<>'";&|$\x60\\]`)

	// Valid YouTube URL hosts
	validYouTubeHosts = map[string]bool{
		"youtube.com":     true,
		"www.youtube.com": true,
		"youtu.be":        true,
		"m.youtube.com":   true,
	}
)

// ValidateChannelInput validates and sanitizes channel input (URL, ID, or handle)
func ValidateChannelInput(input string) (string, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return "", fmt.Errorf("channel input cannot be empty")
	}

	// Check for unsafe characters
	if unsafeCharsRegex.MatchString(input) {
		return "", ErrUnsafeInput
	}

	// If it looks like a URL, validate it
	if strings.Contains(input, "://") || strings.HasPrefix(input, "www.") {
		return validateYouTubeURL(input)
	}

	// If it's a channel ID
	if channelIDRegex.MatchString(input) {
		return input, nil
	}

	// If it's a handle (with @)
	if handleRegex.MatchString(input) {
		return input, nil
	}

	// If it's a handle without @, add it
	if len(input) >= 3 && len(input) <= 30 && !strings.HasPrefix(input, "@") {
		handle := "@" + input
		if handleRegex.MatchString(handle) {
			return handle, nil
		}
	}

	return "", fmt.Errorf("invalid channel identifier: %s", input)
}

// validateYouTubeURL validates a YouTube URL and extracts the channel identifier
func validateYouTubeURL(rawURL string) (string, error) {
	// Add scheme if missing
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", ErrInvalidURL
	}

	// Validate host
	if !validYouTubeHosts[parsed.Host] {
		return "", fmt.Errorf("not a valid YouTube URL: %s", parsed.Host)
	}

	// Extract channel identifier from path
	path := strings.TrimPrefix(parsed.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		return "", fmt.Errorf("no channel identifier in URL")
	}

	switch parts[0] {
	case "channel":
		if len(parts) >= 2 && channelIDRegex.MatchString(parts[1]) {
			return parts[1], nil
		}
		return "", ErrInvalidChannelID

	case "c", "user":
		if len(parts) >= 2 {
			// Custom URL or username - return as-is, API will resolve
			return parts[1], nil
		}
		return "", fmt.Errorf("missing channel name in URL")

	default:
		// Could be a handle (@username)
		if strings.HasPrefix(parts[0], "@") {
			if handleRegex.MatchString(parts[0]) {
				return parts[0], nil
			}
			return "", fmt.Errorf("invalid handle format")
		}
		// Could be a custom URL
		return parts[0], nil
	}
}

// ValidateVideoID validates a YouTube video ID
func ValidateVideoID(videoID string) error {
	videoID = strings.TrimSpace(videoID)

	if videoID == "" {
		return fmt.Errorf("video ID cannot be empty")
	}

	if unsafeCharsRegex.MatchString(videoID) {
		return ErrUnsafeInput
	}

	if !videoIDRegex.MatchString(videoID) {
		return ErrInvalidVideoID
	}

	return nil
}

// ValidateVideoURL validates a YouTube video URL and extracts the video ID
func ValidateVideoURL(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)

	if unsafeCharsRegex.MatchString(rawURL) {
		return "", ErrUnsafeInput
	}

	// Add scheme if missing
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", ErrInvalidURL
	}

	// Validate host
	if !validYouTubeHosts[parsed.Host] {
		return "", fmt.Errorf("not a valid YouTube URL")
	}

	// Extract video ID
	var videoID string

	if parsed.Host == "youtu.be" {
		// Short URL format: youtu.be/VIDEO_ID
		videoID = strings.TrimPrefix(parsed.Path, "/")
	} else {
		// Standard URL format: youtube.com/watch?v=VIDEO_ID
		videoID = parsed.Query().Get("v")
		if videoID == "" {
			// Try embed format: youtube.com/embed/VIDEO_ID
			path := strings.TrimPrefix(parsed.Path, "/")
			if strings.HasPrefix(path, "embed/") {
				videoID = strings.TrimPrefix(path, "embed/")
			} else if strings.HasPrefix(path, "v/") {
				videoID = strings.TrimPrefix(path, "v/")
			}
		}
	}

	if videoID == "" {
		return "", fmt.Errorf("no video ID found in URL")
	}

	// Remove any trailing path or query components
	if idx := strings.IndexAny(videoID, "/?&"); idx != -1 {
		videoID = videoID[:idx]
	}

	if err := ValidateVideoID(videoID); err != nil {
		return "", err
	}

	return videoID, nil
}

// SanitizeFilename sanitizes a filename for safe filesystem use
func SanitizeFilename(filename string) string {
	// Replace unsafe characters
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		"\x00", "",
	)

	sanitized := replacer.Replace(filename)

	// Remove leading/trailing spaces and dots
	sanitized = strings.Trim(sanitized, " .")

	// Limit length
	if len(sanitized) > 200 {
		sanitized = sanitized[:200]
	}

	// Ensure non-empty
	if sanitized == "" {
		sanitized = "unnamed"
	}

	return sanitized
}
