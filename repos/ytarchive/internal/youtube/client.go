package youtube

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/timholm/ytarchive/internal/metrics"
	"github.com/timholm/ytarchive/internal/ratelimit"
)

const (
	// Innertube API endpoints
	browseEndpoint  = "https://www.youtube.com/youtubei/v1/browse"
	playerEndpoint  = "https://www.youtube.com/youtubei/v1/player"
	resolveEndpoint = "https://www.youtube.com/youtubei/v1/navigation/resolve_url"

	// Default client configuration (WEB client)
	defaultClientName    = "WEB"
	defaultClientVersion = "2.20240101.00.00"
	defaultHL            = "en"
	defaultGL            = "US"

	// User agent for requests
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

// Default rate limits
const (
	defaultRateLimit = 2.0 // requests per second
	defaultBurstSize = 5   // max burst
)

// Client provides methods to fetch YouTube channel and video data using native HTTP requests
type Client struct {
	httpClient    *http.Client
	timeout       time.Duration
	clientName    string
	clientVersion string
	userAgent     string
	cookies       []*http.Cookie
	cookieHeader  string
	rateLimiter   *ratelimit.Limiter
	hlsParser     *HLSManifestParser
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithTimeout sets the timeout for HTTP requests
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithClientVersion sets a custom client version for innertube requests
func WithClientVersion(version string) ClientOption {
	return func(c *Client) {
		c.clientVersion = version
	}
}

// WithUserAgent sets a custom user agent
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithCookies sets cookies for authenticated requests
func WithCookies(cookies []*http.Cookie) ClientOption {
	return func(c *Client) {
		c.cookies = cookies
		c.cookieHeader = GetCookieHeader(cookies)
	}
}

// WithCookiesFromFile loads cookies from a Netscape format file
func WithCookiesFromFile(filepath string) ClientOption {
	return func(c *Client) {
		cookies, err := LoadCookiesFromFile(filepath)
		if err == nil && len(cookies) > 0 {
			c.cookies = cookies
			c.cookieHeader = GetCookieHeader(cookies)
		}
	}
}

// WithCookiesFromString loads cookies from a Netscape format string
func WithCookiesFromString(content string) ClientOption {
	return func(c *Client) {
		cookies, err := LoadCookiesFromString(content)
		if err == nil && len(cookies) > 0 {
			c.cookies = cookies
			c.cookieHeader = GetCookieHeader(cookies)
		}
	}
}

// WithRateLimit sets custom rate limiting
func WithRateLimit(requestsPerSecond float64, burst int) ClientOption {
	return func(c *Client) {
		c.rateLimiter = ratelimit.NewLimiter(requestsPerSecond, burst)
	}
}

// NewClient creates a new YouTube client using native HTTP requests
func NewClient(opts ...ClientOption) (*Client, error) {
	c := &Client{
		timeout:       5 * time.Minute,
		clientName:    defaultClientName,
		clientVersion: defaultClientVersion,
		userAgent:     defaultUserAgent,
		rateLimiter:   ratelimit.NewLimiter(defaultRateLimit, defaultBurstSize),
	}

	for _, opt := range opts {
		opt(c)
	}

	// Create HTTP client if not provided
	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Timeout: c.timeout,
		}
	}

	// Initialize HLS parser
	c.hlsParser = NewHLSManifestParser(c.httpClient, c.userAgent)

	return c, nil
}

// createContext creates the innertube context for API requests
func (c *Client) createContext() InnertubeContext {
	return InnertubeContext{
		Client: InnertubeClient{
			HL:            defaultHL,
			GL:            defaultGL,
			ClientName:    c.clientName,
			ClientVersion: c.clientVersion,
		},
	}
}

// Retry configuration
const (
	maxRetries     = 3
	baseRetryDelay = 1 * time.Second
	maxRetryDelay  = 30 * time.Second
)

// retryableStatusCodes are HTTP status codes that warrant a retry
var retryableStatusCodes = map[int]bool{
	http.StatusTooManyRequests:     true, // 429
	http.StatusInternalServerError: true, // 500
	http.StatusBadGateway:          true, // 502
	http.StatusServiceUnavailable:  true, // 503
	http.StatusGatewayTimeout:      true, // 504
}

// doRequest performs an HTTP POST request with retry and exponential backoff
func (c *Client) doRequest(ctx context.Context, endpoint string, body interface{}) ([]byte, error) {
	// Apply rate limiting
	if c.rateLimiter != nil {
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limit wait cancelled: %w", err)
		}
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Extract endpoint name for metrics
	endpointName := extractEndpointName(endpoint)

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Calculate backoff delay with jitter
			delay := calculateBackoff(attempt)
			metrics.RecordAPIRetry(endpointName)

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		startTime := time.Now()
		respBody, statusCode, err := c.doSingleRequest(ctx, endpoint, jsonBody)
		latency := time.Since(startTime).Seconds()

		if err != nil {
			lastErr = err
			// Check if it's a context error (don't retry)
			if ctx.Err() != nil {
				metrics.RecordAPIRequest(endpointName, "context_cancelled", latency)
				return nil, ctx.Err()
			}
			// Network error - retry
			metrics.RecordAPIRequest(endpointName, "error", latency)
			continue
		}

		if statusCode == http.StatusOK {
			metrics.RecordAPIRequest(endpointName, "success", latency)
			return respBody, nil
		}

		// Check if status code is retryable
		if retryableStatusCodes[statusCode] {
			lastErr = fmt.Errorf("request failed with status %d", statusCode)
			metrics.RecordAPIRequest(endpointName, fmt.Sprintf("%d", statusCode), latency)
			continue
		}

		// Non-retryable error
		metrics.RecordAPIRequest(endpointName, fmt.Sprintf("%d", statusCode), latency)
		return nil, fmt.Errorf("request failed with status %d: %s", statusCode, string(respBody))
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", maxRetries+1, lastErr)
}

// doSingleRequest performs a single HTTP request without retry logic
func (c *Client) doSingleRequest(ctx context.Context, endpoint string, jsonBody []byte) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://www.youtube.com")
	req.Header.Set("Referer", "https://www.youtube.com/")

	// Add cookies if available for authenticated requests
	if c.cookieHeader != "" {
		req.Header.Set("Cookie", c.cookieHeader)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, resp.StatusCode, nil
}

// calculateBackoff returns the delay for the given retry attempt using exponential backoff
func calculateBackoff(attempt int) time.Duration {
	// Exponential backoff: baseDelay * 2^attempt
	delay := time.Duration(float64(baseRetryDelay) * math.Pow(2, float64(attempt-1)))
	if delay > maxRetryDelay {
		delay = maxRetryDelay
	}
	return delay
}

// extractEndpointName extracts a short name from the full endpoint URL
func extractEndpointName(endpoint string) string {
	switch {
	case strings.Contains(endpoint, "/browse"):
		return "browse"
	case strings.Contains(endpoint, "/player"):
		return "player"
	case strings.Contains(endpoint, "/resolve_url"):
		return "resolve_url"
	default:
		return "unknown"
	}
}

// normalizeChannelURL converts various YouTube channel URL formats to a standard format
func normalizeChannelURL(url string) string {
	url = strings.TrimSpace(url)

	// Handle @username format (e.g., @aperturethinking or https://youtube.com/@aperturethinking)
	if strings.HasPrefix(url, "@") {
		return fmt.Sprintf("https://www.youtube.com/%s", url)
	}

	// If it's just a username without @, assume it's a handle
	if !strings.Contains(url, "/") && !strings.Contains(url, ".") {
		return fmt.Sprintf("https://www.youtube.com/@%s", url)
	}

	// Ensure URL has scheme
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	// Normalize youtube.com to www.youtube.com
	url = strings.Replace(url, "://youtube.com", "://www.youtube.com", 1)

	return url
}

// extractChannelID attempts to extract channel ID from URL or returns the URL for further processing
func extractChannelID(url string) string {
	// Channel ID pattern: UC followed by 22 characters
	channelIDPattern := regexp.MustCompile(`UC[\w-]{22}`)

	if match := channelIDPattern.FindString(url); match != "" {
		return match
	}

	return url
}

// extractHandleFromURL extracts the @handle from a YouTube URL
func extractHandleFromURL(urlStr string) string {
	// Match @username in URL
	handlePattern := regexp.MustCompile(`@([\w.-]+)`)
	if match := handlePattern.FindStringSubmatch(urlStr); len(match) > 1 {
		return "@" + match[1]
	}
	return ""
}

// resolveChannelID resolves a channel URL or handle to a channel ID (browseId)
func (c *Client) resolveChannelID(ctx context.Context, channelURL string) (string, error) {
	// If it's already a channel ID, return it
	if strings.HasPrefix(channelURL, "UC") && len(channelURL) == 24 {
		return channelURL, nil
	}

	normalizedURL := normalizeChannelURL(channelURL)

	// Try to resolve using the resolve_url endpoint
	req := ResolveURLRequest{
		Context: c.createContext(),
		URL:     normalizedURL,
	}

	data, err := c.doRequest(ctx, resolveEndpoint, req)
	if err != nil {
		return "", fmt.Errorf("failed to resolve URL: %w", err)
	}

	var resp ResolveURLResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", fmt.Errorf("failed to parse resolve response: %w", err)
	}

	if resp.Endpoint.BrowseEndpoint != nil && resp.Endpoint.BrowseEndpoint.BrowseID != "" {
		return resp.Endpoint.BrowseEndpoint.BrowseID, nil
	}

	return "", fmt.Errorf("could not resolve channel ID from URL: %s", channelURL)
}

// GetChannelInfo fetches channel metadata from a YouTube channel URL
func (c *Client) GetChannelInfo(url string) (*Channel, error) {
	return c.GetChannelInfoContext(context.Background(), url)
}

// GetChannelInfoContext fetches channel metadata with context support
func (c *Client) GetChannelInfoContext(ctx context.Context, channelURL string) (*Channel, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Resolve channel ID
	channelID, err := c.resolveChannelID(ctx, channelURL)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve channel: %w", err)
	}

	// Fetch channel info using browse endpoint
	req := BrowseRequest{
		Context:  c.createContext(),
		BrowseID: channelID,
	}

	data, err := c.doRequest(ctx, browseEndpoint, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch channel info: %w", err)
	}

	var resp BrowseResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse channel response: %w", err)
	}

	channel := parseChannelFromBrowseResponse(&resp, channelID)

	// Set the URL to the normalized version if not already set
	if channel.URL == "" {
		channel.URL = normalizeChannelURL(channelURL)
	}

	return channel, nil
}

// GetVideoList fetches all video IDs from a channel
func (c *Client) GetVideoList(channelID string) ([]Video, error) {
	return c.GetVideoListContext(context.Background(), channelID)
}

// GetVideoListContext fetches all video IDs from a channel with context support
func (c *Client) GetVideoListContext(ctx context.Context, channelID string) ([]Video, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Resolve to channel ID if needed
	resolvedID, err := c.resolveChannelID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve channel: %w", err)
	}

	// Fetch the videos tab using browse endpoint
	// The params "EgZ2aWRlb3PyBgQKAjoA" is base64 encoded protobuf for "videos" tab sorted by date
	req := BrowseRequest{
		Context:  c.createContext(),
		BrowseID: resolvedID,
		Params:   "EgZ2aWRlb3PyBgQKAjoA", // Videos tab, sorted by date
	}

	data, err := c.doRequest(ctx, browseEndpoint, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch video list: %w", err)
	}

	var resp BrowseResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse browse response: %w", err)
	}

	videos := parseVideosFromBrowseResponse(&resp)

	// Fetch continuation pages
	continuationToken := extractContinuationToken(&resp)
	for continuationToken != "" {
		moreVideos, nextToken, err := c.fetchVideoContinuation(ctx, continuationToken)
		if err != nil {
			// Log error but continue with what we have
			break
		}
		videos = append(videos, moreVideos...)
		continuationToken = nextToken
	}

	return videos, nil
}

// fetchVideoContinuation fetches the next page of videos using continuation token
func (c *Client) fetchVideoContinuation(ctx context.Context, continuationToken string) ([]Video, string, error) {
	req := BrowseRequest{
		Context:      c.createContext(),
		Continuation: continuationToken,
	}

	data, err := c.doRequest(ctx, browseEndpoint, req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch continuation: %w", err)
	}

	var resp BrowseResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, "", fmt.Errorf("failed to parse continuation response: %w", err)
	}

	videos := parseVideosFromContinuationResponse(&resp)
	nextToken := extractContinuationTokenFromActions(&resp)

	return videos, nextToken, nil
}

// GetVideoMetadata fetches detailed metadata for a specific video
func (c *Client) GetVideoMetadata(videoID string) (*VideoMetadata, error) {
	return c.GetVideoMetadataContext(context.Background(), videoID)
}

// GetVideoMetadataContext fetches detailed metadata for a specific video with context support
func (c *Client) GetVideoMetadataContext(ctx context.Context, videoID string) (*VideoMetadata, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Clean the video ID
	videoID = extractVideoID(videoID)
	if videoID == "" {
		return nil, fmt.Errorf("invalid video ID")
	}

	// Fetch video info using player endpoint
	req := PlayerRequest{
		Context: c.createContext(),
		VideoID: videoID,
	}

	data, err := c.doRequest(ctx, playerEndpoint, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch video metadata: %w", err)
	}

	var resp PlayerResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse player response: %w", err)
	}

	// Check playability status
	if resp.PlayabilityStatus != nil && resp.PlayabilityStatus.Status != "OK" {
		reason := resp.PlayabilityStatus.Reason
		if reason == "" {
			reason = resp.PlayabilityStatus.Status
		}
		return nil, fmt.Errorf("video not playable: %s", reason)
	}

	return parseVideoMetadataFromPlayerResponse(&resp), nil
}

// GetStreamURL gets the best stream URL for a video
func (c *Client) GetStreamURL(videoID string) (*StreamInfo, error) {
	return c.GetStreamURLContext(context.Background(), videoID)
}

// GetStreamURLContext gets the stream URL for a video with context support
func (c *Client) GetStreamURLContext(ctx context.Context, videoID string) (*StreamInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	videoID = extractVideoID(videoID)
	if videoID == "" {
		return nil, fmt.Errorf("invalid video ID")
	}

	// Try ANDROID client first - it returns direct URLs without signature encryption
	streamInfo, err := c.getStreamURLWithAndroidClient(ctx, videoID)
	if err == nil && len(streamInfo.Formats) > 0 {
		return streamInfo, nil
	}

	// Fall back to WEB client if ANDROID fails
	req := PlayerRequest{
		Context: c.createContext(),
		VideoID: videoID,
	}

	data, err := c.doRequest(ctx, playerEndpoint, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stream info: %w", err)
	}

	var resp PlayerResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse player response: %w", err)
	}

	if resp.PlayabilityStatus != nil && resp.PlayabilityStatus.Status != "OK" {
		reason := resp.PlayabilityStatus.Reason
		if reason == "" {
			reason = resp.PlayabilityStatus.Status
		}
		return nil, fmt.Errorf("video not playable: %s", reason)
	}

	streamInfo, err = parseStreamInfoFromPlayerResponse(&resp)
	if err != nil {
		return nil, err
	}

	// If no formats found, try HLS manifest as fallback
	if len(streamInfo.Formats) == 0 && streamInfo.HLSManifestURL != "" {
		hlsFormats, hlsErr := c.getFormatsFromHLS(ctx, streamInfo.HLSManifestURL)
		if hlsErr == nil && len(hlsFormats) > 0 {
			streamInfo.Formats = hlsFormats
		}
	}

	return streamInfo, nil
}

// getFormatsFromHLS fetches and parses HLS manifest to get downloadable formats
func (c *Client) getFormatsFromHLS(ctx context.Context, hlsURL string) ([]DownloadableFormat, error) {
	if c.hlsParser == nil {
		return nil, fmt.Errorf("HLS parser not initialized")
	}

	streams, err := c.hlsParser.ParseManifestURL(ctx, hlsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HLS manifest: %w", err)
	}

	return ConvertHLSToFormats(streams), nil
}

// clientConfig holds configuration for different YouTube client types
type clientConfig struct {
	ClientName    string
	ClientVersion string
	UserAgent     string
	APIKey        string
}

// getClientConfigs returns client configurations to try in order
func getClientConfigs() []clientConfig {
	return []clientConfig{
		// ANDROID client - most reliable for getting direct URLs
		{
			ClientName:    "ANDROID",
			ClientVersion: "19.09.37",
			UserAgent:     "com.google.android.youtube/19.09.37 (Linux; U; Android 13) gzip",
			APIKey:        "AIzaSyA8eiZmM1FaDVjRy-df2KTyQ_vz_yYM39w",
		},
		// ANDROID_CREATOR - backup option
		{
			ClientName:    "ANDROID_CREATOR",
			ClientVersion: "24.06.100",
			UserAgent:     "com.google.android.apps.youtube.creator/24.06.100 (Linux; U; Android 13) gzip",
			APIKey:        "AIzaSyD_ICzD_MV9O5kyyrv0s3j3G1_823Wp0JA",
		},
		// TV embedded client - sometimes bypasses restrictions
		{
			ClientName:    "TVHTML5_SIMPLY_EMBEDDED_PLAYER",
			ClientVersion: "2.0",
			UserAgent:     "Mozilla/5.0 (PlayStation; PlayStation 4/11.00) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4 Safari/605.1.15",
			APIKey:        "",
		},
	}
}

// getStreamURLWithAndroidClient tries multiple client types to get working stream URLs
func (c *Client) getStreamURLWithAndroidClient(ctx context.Context, videoID string) (*StreamInfo, error) {
	var lastErr error

	for _, cfg := range getClientConfigs() {
		clientCtx := InnertubeContext{
			Client: InnertubeClient{
				HL:            defaultHL,
				GL:            defaultGL,
				ClientName:    cfg.ClientName,
				ClientVersion: cfg.ClientVersion,
			},
		}

		req := PlayerRequest{
			Context: clientCtx,
			VideoID: videoID,
		}

		jsonBody, err := json.Marshal(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to marshal request body: %w", err)
			continue
		}

		endpoint := playerEndpoint
		if cfg.APIKey != "" {
			endpoint = fmt.Sprintf("%s?key=%s", playerEndpoint, cfg.APIKey)
		}

		httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonBody))
		if err != nil {
			lastErr = fmt.Errorf("failed to create request: %w", err)
			continue
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("User-Agent", cfg.UserAgent)
		httpReq.Header.Set("Accept", "*/*")
		httpReq.Header.Set("X-Goog-Api-Format-Version", "2")

		// Add cookies if available
		if c.cookieHeader != "" {
			httpReq.Header.Set("Cookie", c.cookieHeader)
		}

		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			lastErr = fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
			continue
		}

		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to read response: %w", err)
			continue
		}

		var playerResp PlayerResponse
		if err := json.Unmarshal(data, &playerResp); err != nil {
			lastErr = fmt.Errorf("failed to parse player response: %w", err)
			continue
		}

		if playerResp.PlayabilityStatus != nil && playerResp.PlayabilityStatus.Status != "OK" {
			reason := playerResp.PlayabilityStatus.Reason
			if reason == "" {
				reason = playerResp.PlayabilityStatus.Status
			}
			lastErr = fmt.Errorf("video not playable with %s: %s", cfg.ClientName, reason)
			continue
		}

		streamInfo, err := parseStreamInfoFromPlayerResponse(&playerResp)
		if err != nil {
			lastErr = err
			continue
		}

		// Successfully got streams with this client
		if len(streamInfo.Formats) > 0 {
			return streamInfo, nil
		}
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("no working client found for video %s", videoID)
}

// extractVideoID extracts video ID from various formats
func extractVideoID(input string) string {
	input = strings.TrimSpace(input)

	// Check if it's already a video ID (11 characters)
	videoIDPattern := regexp.MustCompile(`^[\w-]{11}$`)
	if videoIDPattern.MatchString(input) {
		return input
	}

	// Try to extract from URL
	patterns := []string{
		`youtube\.com/watch\?v=([\w-]{11})`,
		`youtube\.com/embed/([\w-]{11})`,
		`youtube\.com/v/([\w-]{11})`,
		`youtu\.be/([\w-]{11})`,
		`youtube\.com/shorts/([\w-]{11})`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if match := re.FindStringSubmatch(input); len(match) > 1 {
			return match[1]
		}
	}

	return ""
}

// GetVideoListPaginated fetches videos with pagination support
func (c *Client) GetVideoListPaginated(channelID string, pageSize int, pageToken string) ([]Video, string, error) {
	return c.GetVideoListPaginatedContext(context.Background(), channelID, pageSize, pageToken)
}

// GetVideoListPaginatedContext fetches videos with pagination support and context
func (c *Client) GetVideoListPaginatedContext(ctx context.Context, channelID string, pageSize int, pageToken string) ([]Video, string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	var videos []Video
	var nextToken string

	if pageToken == "" {
		// First page - resolve channel and fetch videos tab
		resolvedID, err := c.resolveChannelID(ctx, channelID)
		if err != nil {
			return nil, "", fmt.Errorf("failed to resolve channel: %w", err)
		}

		req := BrowseRequest{
			Context:  c.createContext(),
			BrowseID: resolvedID,
			Params:   "EgZ2aWRlb3PyBgQKAjoA",
		}

		data, err := c.doRequest(ctx, browseEndpoint, req)
		if err != nil {
			return nil, "", fmt.Errorf("failed to fetch video list: %w", err)
		}

		var resp BrowseResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, "", fmt.Errorf("failed to parse browse response: %w", err)
		}

		videos = parseVideosFromBrowseResponse(&resp)
		nextToken = extractContinuationToken(&resp)
	} else {
		// Continuation page
		req := BrowseRequest{
			Context:      c.createContext(),
			Continuation: pageToken,
		}

		data, err := c.doRequest(ctx, browseEndpoint, req)
		if err != nil {
			return nil, "", fmt.Errorf("failed to fetch continuation: %w", err)
		}

		var resp BrowseResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, "", fmt.Errorf("failed to parse continuation response: %w", err)
		}

		videos = parseVideosFromContinuationResponse(&resp)
		nextToken = extractContinuationTokenFromActions(&resp)
	}

	// Limit to pageSize if needed
	if pageSize > 0 && len(videos) > pageSize {
		videos = videos[:pageSize]
	}

	return videos, nextToken, nil
}

// decodeSignatureCipher decodes a signature cipher URL
func decodeSignatureCipher(cipher string) (string, string, string, error) {
	params, err := url.ParseQuery(cipher)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse signature cipher: %w", err)
	}

	streamURL := params.Get("url")
	sig := params.Get("s")
	sp := params.Get("sp")
	if sp == "" {
		sp = "signature"
	}

	return streamURL, sig, sp, nil
}

// parseStreamInfoFromPlayerResponse extracts stream information from player response
func parseStreamInfoFromPlayerResponse(resp *PlayerResponse) (*StreamInfo, error) {
	if resp.StreamingData == nil {
		return nil, fmt.Errorf("no streaming data available")
	}

	info := &StreamInfo{
		VideoID:   resp.VideoDetails.VideoID,
		Title:     resp.VideoDetails.Title,
		ExpiresIn: 0,
		Formats:   make([]DownloadableFormat, 0),
	}

	if resp.StreamingData.ExpiresInSeconds != "" {
		expires, _ := strconv.Atoi(resp.StreamingData.ExpiresInSeconds)
		info.ExpiresIn = time.Duration(expires) * time.Second
	}

	// Collect all formats
	allFormats := append(resp.StreamingData.Formats, resp.StreamingData.AdaptiveFormats...)

	for _, f := range allFormats {
		streamURL := f.URL

		// Handle signature cipher if present
		if streamURL == "" && f.SignatureCipher != "" {
			baseURL, sig, sp, err := decodeSignatureCipher(f.SignatureCipher)
			if err != nil {
				continue
			}
			// Note: Signature decoding would require fetching and parsing the player.js
			// For now, we'll skip ciphered streams as they require complex JS interpretation
			// Most videos have at least some streams without cipher
			_ = sig
			_ = sp
			streamURL = baseURL
		}

		if streamURL == "" {
			continue
		}

		format := DownloadableFormat{
			ITag:            f.ITag,
			URL:             streamURL,
			MimeType:        f.MimeType,
			Bitrate:         f.Bitrate,
			Width:           f.Width,
			Height:          f.Height,
			Quality:         f.Quality,
			QualityLabel:    f.QualityLabel,
			AudioQuality:    f.AudioQuality,
			AudioSampleRate: f.AudioSampleRate,
			AudioChannels:   f.AudioChannels,
			FPS:             f.FPS,
		}

		if f.ContentLength != "" {
			format.ContentLength, _ = strconv.ParseInt(f.ContentLength, 10, 64)
		}

		info.Formats = append(info.Formats, format)
	}

	// Set HLS and DASH manifest URLs if available
	info.HLSManifestURL = resp.StreamingData.HLSManifestURL
	info.DASHManifestURL = resp.StreamingData.DashManifestURL

	if len(info.Formats) == 0 && info.HLSManifestURL == "" && info.DASHManifestURL == "" {
		return nil, fmt.Errorf("no playable streams found")
	}

	return info, nil
}

// CaptionInfo represents information about a caption track
type CaptionInfo struct {
	LanguageCode string `json:"language_code"`
	LanguageName string `json:"language_name"`
	BaseURL      string `json:"base_url"`
	VssID        string `json:"vss_id"`
	IsAutomatic  bool   `json:"is_automatic"`
}

// GetCaptions fetches available caption tracks for a video
func (c *Client) GetCaptions(videoID string) ([]CaptionInfo, error) {
	return c.GetCaptionsContext(context.Background(), videoID)
}

// GetCaptionsContext fetches available caption tracks for a video with context support
func (c *Client) GetCaptionsContext(ctx context.Context, videoID string) ([]CaptionInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	videoID = extractVideoID(videoID)
	if videoID == "" {
		return nil, fmt.Errorf("invalid video ID")
	}

	// Fetch video info using player endpoint
	req := PlayerRequest{
		Context: c.createContext(),
		VideoID: videoID,
	}

	data, err := c.doRequest(ctx, playerEndpoint, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch captions: %w", err)
	}

	var resp PlayerResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse player response: %w", err)
	}

	if resp.PlayabilityStatus != nil && resp.PlayabilityStatus.Status != "OK" {
		reason := resp.PlayabilityStatus.Reason
		if reason == "" {
			reason = resp.PlayabilityStatus.Status
		}
		return nil, fmt.Errorf("video not playable: %s", reason)
	}

	return parseCaptionsFromPlayerResponse(&resp), nil
}

// parseCaptionsFromPlayerResponse extracts caption information from player response
func parseCaptionsFromPlayerResponse(resp *PlayerResponse) []CaptionInfo {
	var captions []CaptionInfo

	if resp.Captions == nil || resp.Captions.PlayerCaptionsTracklistRenderer == nil {
		return captions
	}

	for _, track := range resp.Captions.PlayerCaptionsTracklistRenderer.CaptionTracks {
		caption := CaptionInfo{
			LanguageCode: track.LanguageCode,
			LanguageName: track.Name.GetText(),
			BaseURL:      track.BaseURL,
			VssID:        track.VssID,
			IsAutomatic:  track.Kind == "asr", // "asr" = Automatic Speech Recognition
		}
		captions = append(captions, caption)
	}

	return captions
}

// DownloadCaption downloads a caption track and returns the content
func (c *Client) DownloadCaption(captionURL string) ([]byte, error) {
	return c.DownloadCaptionContext(context.Background(), captionURL)
}

// DownloadCaptionContext downloads a caption track with context support
func (c *Client) DownloadCaptionContext(ctx context.Context, captionURL string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Add format parameter to get WebVTT format
	if !strings.Contains(captionURL, "fmt=") {
		if strings.Contains(captionURL, "?") {
			captionURL += "&fmt=vtt"
		} else {
			captionURL += "?fmt=vtt"
		}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", captionURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "*/*")

	if c.cookieHeader != "" {
		req.Header.Set("Cookie", c.cookieHeader)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch caption: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("caption request failed with status %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read caption content: %w", err)
	}

	return content, nil
}

// GetCaptionByLanguage finds a caption track by language code
func (c *Client) GetCaptionByLanguage(captions []CaptionInfo, langCode string, preferAutomatic bool) *CaptionInfo {
	var manual, automatic *CaptionInfo

	for i := range captions {
		if captions[i].LanguageCode == langCode {
			if captions[i].IsAutomatic {
				automatic = &captions[i]
			} else {
				manual = &captions[i]
			}
		}
	}

	// Return based on preference
	if preferAutomatic {
		if automatic != nil {
			return automatic
		}
		return manual
	}

	// Default: prefer manual over automatic
	if manual != nil {
		return manual
	}
	return automatic
}
