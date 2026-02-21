package youtube

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// HLSStream represents a stream parsed from an HLS manifest
type HLSStream struct {
	URL        string
	Bandwidth  int
	Resolution string
	Width      int
	Height     int
	Codecs     string
	FrameRate  float64
	IsAudio    bool
}

// HLSManifestParser parses HLS (m3u8) manifests
type HLSManifestParser struct {
	httpClient *http.Client
	userAgent  string
}

// NewHLSManifestParser creates a new HLS manifest parser
func NewHLSManifestParser(httpClient *http.Client, userAgent string) *HLSManifestParser {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if userAgent == "" {
		userAgent = defaultUserAgent
	}
	return &HLSManifestParser{
		httpClient: httpClient,
		userAgent:  userAgent,
	}
}

// ParseManifestURL fetches and parses an HLS manifest from a URL
func (p *HLSManifestParser) ParseManifestURL(ctx context.Context, manifestURL string) ([]HLSStream, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", manifestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", p.userAgent)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("manifest request failed with status %d", resp.StatusCode)
	}

	return p.ParseManifest(resp.Body, manifestURL)
}

// ParseManifest parses an HLS manifest from a reader
func (p *HLSManifestParser) ParseManifest(r io.Reader, baseURL string) ([]HLSStream, error) {
	var streams []HLSStream
	var currentStream *HLSStream

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || line == "#EXTM3U" {
			continue
		}

		// Parse EXT-X-STREAM-INF (video variants)
		if strings.HasPrefix(line, "#EXT-X-STREAM-INF:") {
			attrs := parseAttributes(strings.TrimPrefix(line, "#EXT-X-STREAM-INF:"))
			currentStream = &HLSStream{
				Bandwidth:  parseInt(attrs["BANDWIDTH"]),
				Resolution: attrs["RESOLUTION"],
				Codecs:     attrs["CODECS"],
			}

			// Parse resolution
			if res := attrs["RESOLUTION"]; res != "" {
				parts := strings.Split(res, "x")
				if len(parts) == 2 {
					currentStream.Width, _ = strconv.Atoi(parts[0])
					currentStream.Height, _ = strconv.Atoi(parts[1])
				}
			}

			// Parse frame rate
			if fr := attrs["FRAME-RATE"]; fr != "" {
				currentStream.FrameRate, _ = strconv.ParseFloat(fr, 64)
			}

			continue
		}

		// Parse EXT-X-MEDIA (audio tracks)
		if strings.HasPrefix(line, "#EXT-X-MEDIA:") {
			attrs := parseAttributes(strings.TrimPrefix(line, "#EXT-X-MEDIA:"))
			if attrs["TYPE"] == "AUDIO" && attrs["URI"] != "" {
				audioStream := HLSStream{
					URL:       resolveURL(baseURL, attrs["URI"]),
					Bandwidth: parseInt(attrs["BANDWIDTH"]),
					Codecs:    attrs["CODECS"],
					IsAudio:   true,
				}
				streams = append(streams, audioStream)
			}
			continue
		}

		// URL line (follows EXT-X-STREAM-INF)
		if currentStream != nil && !strings.HasPrefix(line, "#") {
			currentStream.URL = resolveURL(baseURL, line)
			streams = append(streams, *currentStream)
			currentStream = nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading manifest: %w", err)
	}

	return streams, nil
}

// GetSegmentURLs parses a media playlist and returns segment URLs
func (p *HLSManifestParser) GetSegmentURLs(ctx context.Context, playlistURL string) ([]string, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", playlistURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", p.userAgent)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch playlist: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("playlist request failed with status %d", resp.StatusCode)
	}

	var segments []string
	var initSegment string

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		// Parse initialization segment
		if strings.HasPrefix(line, "#EXT-X-MAP:") {
			attrs := parseAttributes(strings.TrimPrefix(line, "#EXT-X-MAP:"))
			if uri := attrs["URI"]; uri != "" {
				initSegment = resolveURL(playlistURL, uri)
			}
			continue
		}

		// Skip tags
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Segment URL
		segments = append(segments, resolveURL(playlistURL, line))
	}

	if err := scanner.Err(); err != nil {
		return nil, "", fmt.Errorf("error reading playlist: %w", err)
	}

	return segments, initSegment, nil
}

// parseAttributes parses HLS tag attributes
func parseAttributes(s string) map[string]string {
	attrs := make(map[string]string)

	// Regex to match key=value or key="value"
	re := regexp.MustCompile(`([A-Z0-9-]+)=("([^"]+)"|([^,]+))`)
	matches := re.FindAllStringSubmatch(s, -1)

	for _, match := range matches {
		key := match[1]
		value := match[3]
		if value == "" {
			value = match[4]
		}
		attrs[key] = value
	}

	return attrs
}

// parseInt safely parses an integer
func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// resolveURL resolves a potentially relative URL against a base URL
func resolveURL(baseURL, relativeURL string) string {
	if strings.HasPrefix(relativeURL, "http://") || strings.HasPrefix(relativeURL, "https://") {
		return relativeURL
	}

	// Find the base path
	lastSlash := strings.LastIndex(baseURL, "/")
	if lastSlash == -1 {
		return relativeURL
	}

	return baseURL[:lastSlash+1] + relativeURL
}

// ConvertHLSToFormats converts HLS streams to DownloadableFormat
func ConvertHLSToFormats(streams []HLSStream) []DownloadableFormat {
	var formats []DownloadableFormat

	for i, stream := range streams {
		format := DownloadableFormat{
			ITag:     9000 + i, // Use high ITag numbers for HLS streams
			URL:      stream.URL,
			Bitrate:  stream.Bandwidth,
			Width:    stream.Width,
			Height:   stream.Height,
			Quality:  qualityFromHeight(stream.Height),
			MimeType: "video/mp4", // HLS typically uses fMP4 or TS
		}

		if stream.IsAudio {
			format.MimeType = "audio/mp4"
			format.AudioQuality = "AUDIO_QUALITY_MEDIUM"
		} else if stream.Height > 0 {
			format.QualityLabel = fmt.Sprintf("%dp", stream.Height)
			if stream.FrameRate >= 50 {
				format.QualityLabel = fmt.Sprintf("%dp%d", stream.Height, int(stream.FrameRate))
			}
			format.FPS = int(stream.FrameRate)
		}

		formats = append(formats, format)
	}

	return formats
}

// qualityFromHeight returns a quality label from video height
func qualityFromHeight(height int) string {
	switch {
	case height >= 2160:
		return "hd2160"
	case height >= 1440:
		return "hd1440"
	case height >= 1080:
		return "hd1080"
	case height >= 720:
		return "hd720"
	case height >= 480:
		return "large"
	case height >= 360:
		return "medium"
	case height >= 240:
		return "small"
	default:
		return "tiny"
	}
}
