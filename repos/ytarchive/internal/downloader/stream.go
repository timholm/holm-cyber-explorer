package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// StreamType represents the type of stream (video, audio, or combined)
type StreamType string

const (
	StreamTypeVideo    StreamType = "video"
	StreamTypeAudio    StreamType = "audio"
	StreamTypeCombined StreamType = "combined"
)

// Stream represents a downloadable video/audio stream
type Stream struct {
	FormatID      string     `json:"format_id"`
	URL           string     `json:"url"`
	Extension     string     `json:"ext"`
	Width         int        `json:"width"`
	Height        int        `json:"height"`
	Bitrate       int        `json:"bitrate"`
	FileSize      int64      `json:"filesize"`
	VCodec        string     `json:"vcodec"`
	ACodec        string     `json:"acodec"`
	Quality       string     `json:"quality"`
	QualityLabel  string     `json:"quality_label"`
	StreamType    StreamType `json:"stream_type"`
	MimeType      string     `json:"mime_type"`
	ContentLength int64      `json:"content_length"`
	FPS           int        `json:"fps"`
	// For DASH/HLS streams
	IsSegmented bool     `json:"is_segmented"`
	SegmentURLs []string `json:"segment_urls,omitempty"`
	InitURL     string   `json:"init_url,omitempty"`
}

// VideoInfo contains metadata about the video and available streams
type VideoInfo struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Duration     int      `json:"duration"`
	ChannelID    string   `json:"channel_id"`
	ChannelName  string   `json:"channel_name"`
	UploadDate   string   `json:"upload_date"`
	ThumbnailURL string   `json:"thumbnail_url"`
	ViewCount    int64    `json:"view_count"`
	Streams      []Stream `json:"streams"`
}

// StreamSelector handles stream selection based on quality preferences
type StreamSelector struct {
	maxHeight       int
	preferCombined  bool
	preferredCodecs []string
}

// NewStreamSelector creates a new StreamSelector with the given preferences
func NewStreamSelector(maxHeight int, preferCombined bool) *StreamSelector {
	return &StreamSelector{
		maxHeight:       maxHeight,
		preferCombined:  preferCombined,
		preferredCodecs: []string{"avc1", "h264", "mp4a", "aac"},
	}
}

// SelectBestStream selects the best stream based on quality preferences
// If preferCombined is true, it will ALWAYS use combined streams when available
// (this is required when ffmpeg is not available for merging video+audio).
// Returns: video stream, audio stream (nil if combined), error
func (s *StreamSelector) SelectBestStream(streams []Stream) (*Stream, *Stream, error) {
	if len(streams) == 0 {
		return nil, nil, fmt.Errorf("no streams available")
	}

	// Find best streams of each type
	bestVideo := s.findBestVideo(streams)
	bestCombined := s.findBestCombined(streams)

	// If we MUST use combined streams (no ffmpeg), use combined if available
	// This takes priority over resolution because video-only won't have audio
	if s.preferCombined && bestCombined != nil {
		return bestCombined, nil, nil
	}

	// Determine which has higher resolution
	videoHeight := 0
	combinedHeight := 0
	if bestVideo != nil {
		videoHeight = bestVideo.Height
	}
	if bestCombined != nil {
		combinedHeight = bestCombined.Height
	}

	// If separate video stream has higher resolution, use video+audio
	if bestVideo != nil && videoHeight > combinedHeight {
		audio := s.findBestAudio(streams)
		return bestVideo, audio, nil
	}

	// Use combined if available
	if bestCombined != nil {
		return bestCombined, nil, nil
	}

	// Fallback to video-only if no combined available
	if bestVideo != nil {
		audio := s.findBestAudio(streams)
		return bestVideo, audio, nil
	}

	return nil, nil, fmt.Errorf("no suitable video stream found")
}

// findBestCombined finds the best combined audio+video stream
// Prefers format 22 (720p) or 18 (360p) for maximum compatibility
func (s *StreamSelector) findBestCombined(streams []Stream) *Stream {
	var candidates []Stream

	for _, stream := range streams {
		if stream.StreamType != StreamTypeCombined {
			continue
		}
		if s.maxHeight > 0 && stream.Height > s.maxHeight {
			continue
		}
		// Only consider streams with valid URLs
		if stream.URL == "" {
			continue
		}
		candidates = append(candidates, stream)
	}

	if len(candidates) == 0 {
		return nil
	}

	// Sort by height (descending), prefer MP4 format, then by bitrate
	sort.Slice(candidates, func(i, j int) bool {
		// Prefer MP4
		iMP4 := candidates[i].Extension == "mp4"
		jMP4 := candidates[j].Extension == "mp4"
		if iMP4 != jMP4 {
			return iMP4
		}
		// Then by height
		if candidates[i].Height != candidates[j].Height {
			return candidates[i].Height > candidates[j].Height
		}
		// Then by bitrate
		return candidates[i].Bitrate > candidates[j].Bitrate
	})

	result := candidates[0]
	return &result
}

// findBestVideo finds the best video-only stream
func (s *StreamSelector) findBestVideo(streams []Stream) *Stream {
	var candidates []Stream

	for _, stream := range streams {
		if stream.StreamType != StreamTypeVideo {
			continue
		}
		if s.maxHeight > 0 && stream.Height > s.maxHeight {
			continue
		}
		if stream.URL == "" {
			continue
		}
		candidates = append(candidates, stream)
	}

	if len(candidates) == 0 {
		return nil
	}

	// Sort by height (descending), prefer MP4/H264, then by bitrate
	sort.Slice(candidates, func(i, j int) bool {
		// Prefer H264/AVC
		iH264 := strings.Contains(candidates[i].VCodec, "avc") || strings.Contains(candidates[i].VCodec, "h264")
		jH264 := strings.Contains(candidates[j].VCodec, "avc") || strings.Contains(candidates[j].VCodec, "h264")
		if iH264 != jH264 {
			return iH264
		}
		// Prefer MP4
		iMP4 := candidates[i].Extension == "mp4"
		jMP4 := candidates[j].Extension == "mp4"
		if iMP4 != jMP4 {
			return iMP4
		}
		// Then by height
		if candidates[i].Height != candidates[j].Height {
			return candidates[i].Height > candidates[j].Height
		}
		// Then by bitrate
		return candidates[i].Bitrate > candidates[j].Bitrate
	})

	result := candidates[0]
	return &result
}

// findBestAudio finds the best audio-only stream
func (s *StreamSelector) findBestAudio(streams []Stream) *Stream {
	var candidates []Stream

	for _, stream := range streams {
		if stream.StreamType != StreamTypeAudio {
			continue
		}
		if stream.URL == "" {
			continue
		}
		candidates = append(candidates, stream)
	}

	if len(candidates) == 0 {
		return nil
	}

	// Sort by bitrate (descending), prefer AAC/M4A
	sort.Slice(candidates, func(i, j int) bool {
		// Prefer AAC
		iAAC := strings.Contains(candidates[i].ACodec, "mp4a") || candidates[i].Extension == "m4a"
		jAAC := strings.Contains(candidates[j].ACodec, "mp4a") || candidates[j].Extension == "m4a"
		if iAAC != jAAC {
			return iAAC
		}
		// Then by bitrate
		return candidates[i].Bitrate > candidates[j].Bitrate
	})

	result := candidates[0]
	return &result
}

// ParseStreamType determines the stream type based on codec information
func ParseStreamType(vcodec, acodec string) StreamType {
	hasVideo := vcodec != "" && vcodec != "none"
	hasAudio := acodec != "" && acodec != "none"

	if hasVideo && hasAudio {
		return StreamTypeCombined
	}
	if hasVideo {
		return StreamTypeVideo
	}
	if hasAudio {
		return StreamTypeAudio
	}
	return StreamTypeCombined // Default assumption
}

// StreamFetcher handles fetching stream information and URLs
type StreamFetcher struct {
	httpClient *http.Client
	userAgent  string
}

// NewStreamFetcher creates a new StreamFetcher
func NewStreamFetcher() *StreamFetcher {
	return &StreamFetcher{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
		},
		userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
}

// GetStreamContentLength fetches the content length of a stream URL
func (f *StreamFetcher) GetStreamContentLength(streamURL string) (int64, error) {
	req, err := http.NewRequest(http.MethodHead, streamURL, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create HEAD request: %w", err)
	}

	// Set headers to mimic browser request
	req.Header.Set("User-Agent", f.userAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://www.youtube.com")
	req.Header.Set("Referer", "https://www.youtube.com/")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("HEAD request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HEAD request returned status %d", resp.StatusCode)
	}

	return resp.ContentLength, nil
}

// SupportsRangeRequests checks if a URL supports Range requests
func (f *StreamFetcher) SupportsRangeRequests(streamURL string) (bool, error) {
	req, err := http.NewRequest(http.MethodHead, streamURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create HEAD request: %w", err)
	}

	// Set headers to mimic browser request
	req.Header.Set("User-Agent", f.userAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://www.youtube.com")
	req.Header.Set("Referer", "https://www.youtube.com/")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("HEAD request failed: %w", err)
	}
	defer resp.Body.Close()

	acceptRanges := resp.Header.Get("Accept-Ranges")
	return acceptRanges == "bytes", nil
}

// VideoURLBuilder helps build YouTube video URLs
type VideoURLBuilder struct{}

// NewVideoURLBuilder creates a new VideoURLBuilder
func NewVideoURLBuilder() *VideoURLBuilder {
	return &VideoURLBuilder{}
}

// BuildWatchURL builds a YouTube watch URL from a video ID
func (b *VideoURLBuilder) BuildWatchURL(videoID string) string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
}

// BuildEmbedURL builds a YouTube embed URL from a video ID
func (b *VideoURLBuilder) BuildEmbedURL(videoID string) string {
	return fmt.Sprintf("https://www.youtube.com/embed/%s", videoID)
}

// BuildThumbnailURL builds a YouTube thumbnail URL from a video ID
func (b *VideoURLBuilder) BuildThumbnailURL(videoID string, quality string) string {
	// quality can be: default, mqdefault, hqdefault, sddefault, maxresdefault
	if quality == "" {
		quality = "maxresdefault"
	}
	return fmt.Sprintf("https://i.ytimg.com/vi/%s/%s.jpg", videoID, quality)
}

// ExtractVideoID extracts the video ID from various YouTube URL formats
func ExtractVideoID(urlStr string) (string, error) {
	// Handle plain video IDs
	if len(urlStr) == 11 && !strings.Contains(urlStr, "/") && !strings.Contains(urlStr, ".") {
		return urlStr, nil
	}

	// Parse the URL
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	// Handle youtu.be URLs
	if strings.Contains(u.Host, "youtu.be") {
		videoID := strings.TrimPrefix(u.Path, "/")
		if len(videoID) >= 11 {
			return videoID[:11], nil
		}
		return "", fmt.Errorf("invalid youtu.be URL")
	}

	// Handle youtube.com URLs
	if strings.Contains(u.Host, "youtube.com") {
		// Check query parameter
		if v := u.Query().Get("v"); v != "" {
			return v, nil
		}

		// Check embed URLs: /embed/VIDEO_ID
		if strings.HasPrefix(u.Path, "/embed/") {
			videoID := strings.TrimPrefix(u.Path, "/embed/")
			if len(videoID) >= 11 {
				return videoID[:11], nil
			}
		}

		// Check shorts URLs: /shorts/VIDEO_ID
		if strings.HasPrefix(u.Path, "/shorts/") {
			videoID := strings.TrimPrefix(u.Path, "/shorts/")
			if len(videoID) >= 11 {
				return videoID[:11], nil
			}
		}

		// Check live URLs: /live/VIDEO_ID
		if strings.HasPrefix(u.Path, "/live/") {
			videoID := strings.TrimPrefix(u.Path, "/live/")
			if len(videoID) >= 11 {
				return videoID[:11], nil
			}
		}
	}

	return "", fmt.Errorf("could not extract video ID from URL: %s", urlStr)
}

// QualityLabel returns a human-readable quality label
func QualityLabel(height int) string {
	switch {
	case height >= 2160:
		return "4K"
	case height >= 1440:
		return "1440p"
	case height >= 1080:
		return "1080p"
	case height >= 720:
		return "720p"
	case height >= 480:
		return "480p"
	case height >= 360:
		return "360p"
	case height >= 240:
		return "240p"
	case height >= 144:
		return "144p"
	default:
		return fmt.Sprintf("%dp", height)
	}
}

// FormatFileSize formats a file size in bytes to a human-readable string
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// KnownFormats contains information about known YouTube format IDs (itags)
var KnownFormats = map[int]struct {
	Extension string
	Height    int
	HasAudio  bool
	HasVideo  bool
}{
	// Combined formats (preferred for no-merge downloads)
	18: {"mp4", 360, true, true},  // MP4 360p
	22: {"mp4", 720, true, true},  // MP4 720p
	37: {"mp4", 1080, true, true}, // MP4 1080p (rare)
	38: {"mp4", 3072, true, true}, // MP4 4K (rare)

	// Video-only formats (DASH)
	133: {"mp4", 240, false, true},  // MP4 240p
	134: {"mp4", 360, false, true},  // MP4 360p
	135: {"mp4", 480, false, true},  // MP4 480p
	136: {"mp4", 720, false, true},  // MP4 720p
	137: {"mp4", 1080, false, true}, // MP4 1080p
	138: {"mp4", 2160, false, true}, // MP4 4K
	160: {"mp4", 144, false, true},  // MP4 144p
	264: {"mp4", 1440, false, true}, // MP4 1440p
	266: {"mp4", 2160, false, true}, // MP4 4K
	298: {"mp4", 720, false, true},  // MP4 720p60
	299: {"mp4", 1080, false, true}, // MP4 1080p60

	// WebM video-only
	167: {"webm", 360, false, true},
	168: {"webm", 480, false, true},
	169: {"webm", 720, false, true},
	170: {"webm", 1080, false, true},
	218: {"webm", 480, false, true},
	219: {"webm", 480, false, true},
	242: {"webm", 240, false, true},
	243: {"webm", 360, false, true},
	244: {"webm", 480, false, true},
	245: {"webm", 480, false, true},
	246: {"webm", 480, false, true},
	247: {"webm", 720, false, true},
	248: {"webm", 1080, false, true},
	271: {"webm", 1440, false, true},
	272: {"webm", 2160, false, true},
	302: {"webm", 720, false, true},  // 60fps
	303: {"webm", 1080, false, true}, // 60fps
	308: {"webm", 1440, false, true}, // 60fps
	313: {"webm", 2160, false, true},
	315: {"webm", 2160, false, true}, // 60fps

	// Audio-only formats
	139: {"m4a", 0, true, false},  // M4A 48kbps
	140: {"m4a", 0, true, false},  // M4A 128kbps
	141: {"m4a", 0, true, false},  // M4A 256kbps
	171: {"webm", 0, true, false}, // WebM audio
	172: {"webm", 0, true, false}, // WebM audio
	249: {"webm", 0, true, false}, // Opus 50kbps
	250: {"webm", 0, true, false}, // Opus 70kbps
	251: {"webm", 0, true, false}, // Opus 160kbps
}

// ParseStreamFromItag creates a Stream from an itag (format ID)
func ParseStreamFromItag(itag int, streamURL string) *Stream {
	info, ok := KnownFormats[itag]
	if !ok {
		return nil
	}

	stream := &Stream{
		FormatID:  fmt.Sprintf("%d", itag),
		URL:       streamURL,
		Extension: info.Extension,
		Height:    info.Height,
	}

	if info.HasAudio && info.HasVideo {
		stream.StreamType = StreamTypeCombined
		stream.VCodec = "h264"
		stream.ACodec = "aac"
	} else if info.HasVideo {
		stream.StreamType = StreamTypeVideo
		stream.VCodec = "h264"
		stream.ACodec = "none"
	} else {
		stream.StreamType = StreamTypeAudio
		stream.VCodec = "none"
		stream.ACodec = "aac"
	}

	return stream
}

// ParseFormatFromID parses stream information from a known format ID string
func ParseFormatFromID(formatID string) *Stream {
	itag, err := strconv.Atoi(formatID)
	if err != nil {
		return nil
	}
	return ParseStreamFromItag(itag, "")
}

// ParseStreamsFromPlayerResponse parses streams from YouTube's player API response
func ParseStreamsFromPlayerResponse(data map[string]interface{}) ([]Stream, error) {
	var streams []Stream

	streamingData, ok := data["streamingData"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no streaming data in response")
	}

	// Parse combined formats
	if formats, ok := streamingData["formats"].([]interface{}); ok {
		for _, f := range formats {
			format, ok := f.(map[string]interface{})
			if !ok {
				continue
			}
			stream := parseStreamFormat(format, StreamTypeCombined)
			if stream != nil {
				streams = append(streams, *stream)
			}
		}
	}

	// Parse adaptive formats (separate video/audio)
	if adaptiveFormats, ok := streamingData["adaptiveFormats"].([]interface{}); ok {
		for _, f := range adaptiveFormats {
			format, ok := f.(map[string]interface{})
			if !ok {
				continue
			}
			stream := parseStreamFormat(format, "")
			if stream != nil {
				streams = append(streams, *stream)
			}
		}
	}

	return streams, nil
}

// parseStreamFormat parses a single stream format from the API response
func parseStreamFormat(format map[string]interface{}, defaultType StreamType) *Stream {
	stream := &Stream{}

	// Get itag
	if itag, ok := format["itag"].(float64); ok {
		stream.FormatID = fmt.Sprintf("%d", int(itag))
	}

	// Get URL
	if u, ok := format["url"].(string); ok {
		stream.URL = u
	} else if cipher, ok := format["signatureCipher"].(string); ok {
		// Handle signature cipher (requires decryption)
		parsed, _ := url.ParseQuery(cipher)
		if u := parsed.Get("url"); u != "" {
			stream.URL = u
			// Note: This URL may need signature decryption to work
		}
	}

	// Get MIME type
	if mimeType, ok := format["mimeType"].(string); ok {
		stream.MimeType = mimeType
		container, vcodec, acodec := ParseMimeType(mimeType)
		stream.Extension = container
		stream.VCodec = vcodec
		stream.ACodec = acodec

		if defaultType != "" {
			stream.StreamType = defaultType
		} else {
			stream.StreamType = ParseStreamType(vcodec, acodec)
		}
	}

	// Get dimensions
	if width, ok := format["width"].(float64); ok {
		stream.Width = int(width)
	}
	if height, ok := format["height"].(float64); ok {
		stream.Height = int(height)
	}

	// Get bitrate
	if bitrate, ok := format["bitrate"].(float64); ok {
		stream.Bitrate = int(bitrate)
	}

	// Get content length
	if cl, ok := format["contentLength"].(string); ok {
		stream.ContentLength, _ = strconv.ParseInt(cl, 10, 64)
	}

	// Get quality label
	if ql, ok := format["qualityLabel"].(string); ok {
		stream.QualityLabel = ql
	}
	if q, ok := format["quality"].(string); ok {
		stream.Quality = q
	}

	// Get FPS
	if fps, ok := format["fps"].(float64); ok {
		stream.FPS = int(fps)
	}

	return stream
}

// StreamFromJSON parses a Stream from JSON (useful for receiving stream data from external sources)
func StreamFromJSON(data []byte) (*Stream, error) {
	var stream Stream
	if err := json.Unmarshal(data, &stream); err != nil {
		return nil, fmt.Errorf("failed to parse stream JSON: %w", err)
	}
	return &stream, nil
}

// VideoInfoFromJSON parses VideoInfo from JSON
func VideoInfoFromJSON(data []byte) (*VideoInfo, error) {
	var info VideoInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("failed to parse video info JSON: %w", err)
	}
	return &info, nil
}

// IsValidStreamURL checks if a URL looks like a valid stream URL
func IsValidStreamURL(urlStr string) bool {
	if urlStr == "" {
		return false
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Must be http or https
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	// Must have a host
	if u.Host == "" {
		return false
	}

	return true
}

// ExtractSignatureCipher extracts the signature cipher from a URL if present
// This is used for streams that have encrypted signatures
func ExtractSignatureCipher(urlStr string) (string, string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", "", err
	}

	q := u.Query()

	// Check for signature cipher
	if sc := q.Get("sp"); sc != "" {
		sig := q.Get("sig")
		return sc, sig, nil
	}

	return "", "", nil
}

// ParseMimeType extracts codec information from a MIME type string
// Example: "video/mp4; codecs=\"avc1.42001E, mp4a.40.2\""
func ParseMimeType(mimeType string) (container string, vcodec string, acodec string) {
	// Extract the container type
	parts := strings.Split(mimeType, ";")
	if len(parts) == 0 {
		return "", "", ""
	}

	containerParts := strings.Split(strings.TrimSpace(parts[0]), "/")
	if len(containerParts) == 2 {
		container = containerParts[1]
		// Normalize container names
		if container == "webm" || container == "mp4" {
			// Keep as is
		} else if container == "3gpp" {
			container = "3gp"
		}
	}

	// Extract codecs if present
	if len(parts) > 1 {
		codecsMatch := regexp.MustCompile(`codecs="([^"]+)"`).FindStringSubmatch(mimeType)
		if len(codecsMatch) > 1 {
			codecs := strings.Split(codecsMatch[1], ",")
			for _, codec := range codecs {
				codec = strings.TrimSpace(codec)
				if strings.HasPrefix(codec, "avc") || strings.HasPrefix(codec, "vp") ||
					strings.HasPrefix(codec, "av01") || strings.HasPrefix(codec, "hev") {
					vcodec = codec
				} else if strings.HasPrefix(codec, "mp4a") || strings.HasPrefix(codec, "opus") ||
					strings.HasPrefix(codec, "vorbis") || strings.HasPrefix(codec, "ac-3") {
					acodec = codec
				}
			}
		}
	}

	return container, vcodec, acodec
}

// ParseDASHManifest parses a DASH manifest and extracts segment URLs
func ParseDASHManifest(manifestData []byte) ([]Stream, error) {
	// Basic DASH manifest parsing
	// In a production system, you'd want a full MPD parser

	var streams []Stream

	// This is a simplified parser - real DASH manifests are more complex
	// For now, we'll rely on pre-parsed stream data

	return streams, nil
}

// ParseHLSManifest parses an HLS manifest and extracts segment URLs
func ParseHLSManifest(manifestData []byte) ([]string, error) {
	var segmentURLs []string

	lines := strings.Split(string(manifestData), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Lines that aren't comments and aren't empty are segment URLs
		segmentURLs = append(segmentURLs, line)
	}

	return segmentURLs, nil
}

// FetchHLSSegmentURLs fetches and parses an HLS manifest to get segment URLs
func (f *StreamFetcher) FetchHLSSegmentURLs(manifestURL string) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, manifestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", f.userAgent)

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("manifest request returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	return ParseHLSManifest(data)
}

// CalculateTotalSize calculates the total size for segmented streams
func (f *StreamFetcher) CalculateTotalSize(segmentURLs []string) (int64, error) {
	var totalSize int64

	for _, segURL := range segmentURLs {
		size, err := f.GetStreamContentLength(segURL)
		if err != nil {
			// If we can't get size for one segment, estimate based on others
			continue
		}
		totalSize += size
	}

	return totalSize, nil
}

// StreamInfoString returns a string representation of stream info
func (s *Stream) InfoString() string {
	var parts []string

	if s.FormatID != "" {
		parts = append(parts, fmt.Sprintf("Format: %s", s.FormatID))
	}

	if s.Height > 0 {
		parts = append(parts, QualityLabel(s.Height))
	}

	if s.Extension != "" {
		parts = append(parts, strings.ToUpper(s.Extension))
	}

	if s.Bitrate > 0 {
		parts = append(parts, fmt.Sprintf("%d kbps", s.Bitrate/1000))
	}

	if s.ContentLength > 0 {
		parts = append(parts, FormatFileSize(s.ContentLength))
	}

	return strings.Join(parts, " | ")
}

// BitrateFromQuality estimates bitrate from video height
func BitrateFromQuality(height int) int {
	switch {
	case height >= 2160:
		return 20000000 // 20 Mbps
	case height >= 1440:
		return 10000000 // 10 Mbps
	case height >= 1080:
		return 5000000 // 5 Mbps
	case height >= 720:
		return 2500000 // 2.5 Mbps
	case height >= 480:
		return 1000000 // 1 Mbps
	case height >= 360:
		return 500000 // 500 kbps
	default:
		return 250000 // 250 kbps
	}
}

// ParseBitrateString parses a bitrate string like "128kbps" or "5Mbps"
func ParseBitrateString(s string) int {
	s = strings.ToLower(strings.TrimSpace(s))

	multiplier := 1
	if strings.HasSuffix(s, "kbps") {
		s = strings.TrimSuffix(s, "kbps")
		multiplier = 1000
	} else if strings.HasSuffix(s, "mbps") {
		s = strings.TrimSuffix(s, "mbps")
		multiplier = 1000000
	} else if strings.HasSuffix(s, "k") {
		s = strings.TrimSuffix(s, "k")
		multiplier = 1000
	} else if strings.HasSuffix(s, "m") {
		s = strings.TrimSuffix(s, "m")
		multiplier = 1000000
	}

	val, _ := strconv.ParseFloat(s, 64)
	return int(val * float64(multiplier))
}

// SortStreamsByQuality sorts streams by quality (highest first)
func SortStreamsByQuality(streams []Stream) {
	sort.Slice(streams, func(i, j int) bool {
		// Combined streams first
		if streams[i].StreamType == StreamTypeCombined && streams[j].StreamType != StreamTypeCombined {
			return true
		}
		if streams[i].StreamType != StreamTypeCombined && streams[j].StreamType == StreamTypeCombined {
			return false
		}
		// Then by height
		if streams[i].Height != streams[j].Height {
			return streams[i].Height > streams[j].Height
		}
		// Then by bitrate
		return streams[i].Bitrate > streams[j].Bitrate
	})
}

// FilterStreamsByType filters streams by type
func FilterStreamsByType(streams []Stream, streamType StreamType) []Stream {
	var result []Stream
	for _, s := range streams {
		if s.StreamType == streamType {
			result = append(result, s)
		}
	}
	return result
}

// GetCombinedStreams returns only combined (audio+video) streams
func GetCombinedStreams(streams []Stream) []Stream {
	return FilterStreamsByType(streams, StreamTypeCombined)
}

// GetVideoOnlyStreams returns only video-only streams
func GetVideoOnlyStreams(streams []Stream) []Stream {
	return FilterStreamsByType(streams, StreamTypeVideo)
}

// GetAudioOnlyStreams returns only audio-only streams
func GetAudioOnlyStreams(streams []Stream) []Stream {
	return FilterStreamsByType(streams, StreamTypeAudio)
}

// ResolutionOption represents an available resolution with its details
type ResolutionOption struct {
	Height     int    `json:"height"`
	Width      int    `json:"width"`
	Label      string `json:"label"`
	FormatID   string `json:"format_id"`
	StreamType string `json:"stream_type"`
	Codec      string `json:"codec,omitempty"`
	Bitrate    int    `json:"bitrate,omitempty"`
	FPS        int    `json:"fps,omitempty"`
}

// GetAvailableResolutions returns all unique resolution options from the streams
func GetAvailableResolutions(streams []Stream) []ResolutionOption {
	// Use a map to deduplicate by height (keeping best quality for each height)
	resolutionMap := make(map[int]ResolutionOption)

	for _, stream := range streams {
		if stream.Height <= 0 {
			continue // Skip audio-only streams
		}

		existing, exists := resolutionMap[stream.Height]
		// Keep the one with higher bitrate or prefer combined streams
		if !exists ||
			stream.Bitrate > existing.Bitrate ||
			(stream.StreamType == StreamTypeCombined && existing.StreamType != string(StreamTypeCombined)) {
			resolutionMap[stream.Height] = ResolutionOption{
				Height:     stream.Height,
				Width:      stream.Width,
				Label:      QualityLabel(stream.Height),
				FormatID:   stream.FormatID,
				StreamType: string(stream.StreamType),
				Codec:      stream.VCodec,
				Bitrate:    stream.Bitrate,
				FPS:        stream.FPS,
			}
		}
	}

	// Convert map to slice and sort by height descending
	var options []ResolutionOption
	for _, opt := range resolutionMap {
		options = append(options, opt)
	}

	sort.Slice(options, func(i, j int) bool {
		return options[i].Height > options[j].Height
	})

	return options
}
