package youtube

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// parseChannelFromBrowseResponse extracts channel info from innertube browse response
func parseChannelFromBrowseResponse(resp *BrowseResponse, channelID string) *Channel {
	channel := &Channel{
		ID:        channelID,
		CreatedAt: time.Now(), // YouTube doesn't provide channel creation date
	}

	// Try to get info from metadata renderer (most reliable)
	if resp.Metadata.ChannelMetadataRenderer != nil {
		meta := resp.Metadata.ChannelMetadataRenderer
		channel.Name = meta.Title
		channel.Description = meta.Description
		channel.URL = meta.ChannelUrl
		if channel.URL == "" {
			channel.URL = meta.VanityChannelUrl
		}
		if meta.Avatar.GetBestThumbnail() != "" {
			channel.AvatarURL = meta.Avatar.GetBestThumbnail()
		}
	}

	// Try to get additional info from C4TabbedHeaderRenderer
	if resp.Header.C4TabbedHeaderRenderer != nil {
		header := resp.Header.C4TabbedHeaderRenderer

		if channel.Name == "" {
			channel.Name = header.Title
		}
		if channel.ID == "" {
			channel.ID = header.ChannelID
		}
		if channel.AvatarURL == "" {
			channel.AvatarURL = header.Avatar.GetBestThumbnail()
		}
		if channel.BannerURL == "" {
			channel.BannerURL = header.Banner.GetBestThumbnail()
		}

		// Parse video count from "X videos" text
		if videosText := header.VideosCount.GetText(); videosText != "" {
			channel.VideoCount = parseVideoCount(videosText)
		}
	}

	// Try PageHeaderRenderer for newer channel format
	if resp.Header.PageHeaderRenderer != nil {
		header := resp.Header.PageHeaderRenderer

		if channel.Name == "" {
			channel.Name = header.PageTitle
		}

		if header.Content.PageHeaderViewModel != nil {
			vm := header.Content.PageHeaderViewModel

			// Get title
			if channel.Name == "" && vm.Title.DynamicTextViewModel != nil {
				channel.Name = vm.Title.DynamicTextViewModel.Text.Content
			}

			// Get avatar
			if channel.AvatarURL == "" && vm.Image.DecoratedAvatarViewModel != nil {
				if vm.Image.DecoratedAvatarViewModel.Avatar.AvatarViewModel != nil {
					channel.AvatarURL = vm.Image.DecoratedAvatarViewModel.Avatar.AvatarViewModel.Image.GetBestThumbnail()
				}
			}

			// Get banner
			if channel.BannerURL == "" && vm.Banner.ImageBannerViewModel != nil {
				channel.BannerURL = vm.Banner.ImageBannerViewModel.Image.GetBestThumbnail()
			}

			// Get description
			if channel.Description == "" && vm.Description.DescriptionPreviewViewModel != nil {
				channel.Description = vm.Description.DescriptionPreviewViewModel.Description.Content
			}
		}
	}

	return channel
}

// parseVideoCount extracts numeric video count from text like "1,234 videos"
func parseVideoCount(text string) int {
	// Remove commas and extract number
	re := regexp.MustCompile(`[\d,]+`)
	match := re.FindString(text)
	if match == "" {
		return 0
	}

	match = strings.ReplaceAll(match, ",", "")
	count, _ := strconv.Atoi(match)
	return count
}

// parseVideosFromBrowseResponse extracts videos from initial browse response
func parseVideosFromBrowseResponse(resp *BrowseResponse) []Video {
	videos := make([]Video, 0)

	// Navigate to video grid content
	if resp.Contents.TwoColumnBrowseResultsRenderer == nil {
		return videos
	}

	for _, tab := range resp.Contents.TwoColumnBrowseResultsRenderer.Tabs {
		if tab.TabRenderer == nil || tab.TabRenderer.Content.RichGridRenderer == nil {
			continue
		}

		// Check if this is the videos tab
		grid := tab.TabRenderer.Content.RichGridRenderer
		for _, content := range grid.Contents {
			if video := parseRichGridContent(&content); video != nil {
				videos = append(videos, *video)
			}
		}
	}

	// Also check SectionListRenderer for alternative format
	for _, tab := range resp.Contents.TwoColumnBrowseResultsRenderer.Tabs {
		if tab.TabRenderer == nil || tab.TabRenderer.Content.SectionListRenderer == nil {
			continue
		}

		for _, section := range tab.TabRenderer.Content.SectionListRenderer.Contents {
			if section.ItemSectionRenderer == nil {
				continue
			}

			for _, item := range section.ItemSectionRenderer.Contents {
				if item.GridRenderer == nil {
					continue
				}

				for _, gridItem := range item.GridRenderer.Items {
					if video := parseGridVideoRenderer(gridItem.GridVideoRenderer); video != nil {
						videos = append(videos, *video)
					}
				}
			}
		}
	}

	return videos
}

// parseVideosFromContinuationResponse extracts videos from continuation response
func parseVideosFromContinuationResponse(resp *BrowseResponse) []Video {
	videos := make([]Video, 0)

	for _, action := range resp.OnResponseReceivedActions {
		if action.AppendContinuationItemsAction != nil {
			for _, content := range action.AppendContinuationItemsAction.ContinuationItems {
				if video := parseRichGridContent(&content); video != nil {
					videos = append(videos, *video)
				}
			}
		}

		if action.ReloadContinuationItemsCommand != nil {
			for _, content := range action.ReloadContinuationItemsCommand.ContinuationItems {
				if video := parseRichGridContent(&content); video != nil {
					videos = append(videos, *video)
				}
			}
		}
	}

	return videos
}

// parseRichGridContent extracts a video from RichGridContent
func parseRichGridContent(content *RichGridContent) *Video {
	if content.RichItemRenderer == nil {
		return nil
	}

	ric := content.RichItemRenderer.Content

	// Handle regular video
	if ric.VideoRenderer != nil {
		return parseVideoRenderer(ric.VideoRenderer)
	}

	// Handle short/reel
	if ric.ReelItemRenderer != nil {
		return parseReelItemRenderer(ric.ReelItemRenderer)
	}

	// Handle new shorts format
	if ric.ShortsLockupViewModel != nil {
		return parseShortsLockupViewModel(ric.ShortsLockupViewModel)
	}

	return nil
}

// parseVideoRenderer extracts video from VideoRenderer
func parseVideoRenderer(vr *VideoRenderer) *Video {
	if vr == nil || vr.VideoID == "" {
		return nil
	}

	video := &Video{
		ID:           vr.VideoID,
		Title:        extractTitle(vr.Title),
		Description:  vr.DescriptionSnippet.GetText(),
		ThumbnailURL: vr.Thumbnail.GetBestThumbnail(),
		Status:       StatusPending,
	}

	// Parse duration
	if durationText := vr.LengthText.GetText(); durationText != "" {
		video.Duration = parseDuration(durationText)
	}

	// Parse view count
	if viewText := vr.ViewCountText.GetText(); viewText != "" {
		video.ViewCount = parseViewCount(viewText)
	}

	// Parse upload date (relative date like "2 days ago")
	if publishedText := vr.PublishedTimeText.GetText(); publishedText != "" {
		video.UploadDate = parseRelativeDate(publishedText)
	}

	// Generate thumbnail if not provided
	if video.ThumbnailURL == "" && video.ID != "" {
		video.ThumbnailURL = fmt.Sprintf("https://i.ytimg.com/vi/%s/maxresdefault.jpg", video.ID)
	}

	return video
}

// parseReelItemRenderer extracts video from ReelItemRenderer (shorts)
func parseReelItemRenderer(rr *ReelItemRenderer) *Video {
	if rr == nil || rr.VideoID == "" {
		return nil
	}

	video := &Video{
		ID:           rr.VideoID,
		Title:        rr.Headline.GetText(),
		ThumbnailURL: rr.Thumbnail.GetBestThumbnail(),
		Status:       StatusPending,
	}

	// Generate thumbnail if not provided
	if video.ThumbnailURL == "" && video.ID != "" {
		video.ThumbnailURL = fmt.Sprintf("https://i.ytimg.com/vi/%s/maxresdefault.jpg", video.ID)
	}

	return video
}

// parseShortsLockupViewModel extracts video from ShortsLockupViewModel
func parseShortsLockupViewModel(slvm *ShortsLockupViewModel) *Video {
	if slvm == nil {
		return nil
	}

	videoID := ""

	// Try to extract video ID from onTap action
	if slvm.OnTap.InnertubeCommand.ReelWatchEndpoint != nil {
		videoID = slvm.OnTap.InnertubeCommand.ReelWatchEndpoint.VideoID
	}

	// Try to extract from entity ID (format: "shorts-shelf-item-VIDEO_ID")
	if videoID == "" && slvm.EntityID != "" {
		parts := strings.Split(slvm.EntityID, "-")
		if len(parts) > 0 {
			videoID = parts[len(parts)-1]
		}
	}

	if videoID == "" {
		return nil
	}

	video := &Video{
		ID:           videoID,
		ThumbnailURL: fmt.Sprintf("https://i.ytimg.com/vi/%s/maxresdefault.jpg", videoID),
		Status:       StatusPending,
	}

	return video
}

// parseGridVideoRenderer extracts video from GridVideoRenderer
func parseGridVideoRenderer(gvr *GridVideoRenderer) *Video {
	if gvr == nil || gvr.VideoID == "" {
		return nil
	}

	video := &Video{
		ID:           gvr.VideoID,
		Title:        gvr.Title.GetText(),
		ThumbnailURL: gvr.Thumbnail.GetBestThumbnail(),
		Status:       StatusPending,
	}

	// Parse view count
	if viewText := gvr.ViewCountText.GetText(); viewText != "" {
		video.ViewCount = parseViewCount(viewText)
	}

	// Parse upload date
	if publishedText := gvr.PublishedTimeText.GetText(); publishedText != "" {
		video.UploadDate = parseRelativeDate(publishedText)
	}

	// Generate thumbnail if not provided
	if video.ThumbnailURL == "" && video.ID != "" {
		video.ThumbnailURL = fmt.Sprintf("https://i.ytimg.com/vi/%s/maxresdefault.jpg", video.ID)
	}

	return video
}

// extractTitle extracts title from TextRuns (handling both formats)
func extractTitle(tr TextRuns) string {
	return tr.GetText()
}

// parseDuration parses duration string like "10:30" or "1:02:30" to seconds
func parseDuration(durationStr string) int {
	parts := strings.Split(durationStr, ":")
	if len(parts) == 0 {
		return 0
	}

	var seconds int
	multipliers := []int{1, 60, 3600} // seconds, minutes, hours

	for i := len(parts) - 1; i >= 0; i-- {
		idx := len(parts) - 1 - i
		if idx >= len(multipliers) {
			break
		}
		val, _ := strconv.Atoi(strings.TrimSpace(parts[i]))
		seconds += val * multipliers[idx]
	}

	return seconds
}

// parseViewCount parses view count from text like "1,234 views" or "1.2M views"
func parseViewCount(viewStr string) int64 {
	viewStr = strings.ToLower(viewStr)
	viewStr = strings.ReplaceAll(viewStr, ",", "")
	viewStr = strings.ReplaceAll(viewStr, " views", "")
	viewStr = strings.ReplaceAll(viewStr, " view", "")
	viewStr = strings.TrimSpace(viewStr)

	// Handle "No views" case
	if viewStr == "no" || viewStr == "" {
		return 0
	}

	// Handle multiplier suffixes
	multiplier := int64(1)
	if strings.HasSuffix(viewStr, "k") {
		multiplier = 1000
		viewStr = strings.TrimSuffix(viewStr, "k")
	} else if strings.HasSuffix(viewStr, "m") {
		multiplier = 1000000
		viewStr = strings.TrimSuffix(viewStr, "m")
	} else if strings.HasSuffix(viewStr, "b") {
		multiplier = 1000000000
		viewStr = strings.TrimSuffix(viewStr, "b")
	}

	// Parse the number
	val, err := strconv.ParseFloat(viewStr, 64)
	if err != nil {
		return 0
	}

	return int64(val * float64(multiplier))
}

// parseRelativeDate converts relative date like "2 days ago" to YYYYMMDD format
func parseRelativeDate(relativeStr string) string {
	relativeStr = strings.ToLower(relativeStr)
	now := time.Now()

	// Patterns for different time units
	patterns := map[string]time.Duration{
		"second": time.Second,
		"minute": time.Minute,
		"hour":   time.Hour,
		"day":    24 * time.Hour,
		"week":   7 * 24 * time.Hour,
		"month":  30 * 24 * time.Hour,
		"year":   365 * 24 * time.Hour,
	}

	// Extract number and unit
	re := regexp.MustCompile(`(\d+)\s*(second|minute|hour|day|week|month|year)s?\s*ago`)
	match := re.FindStringSubmatch(relativeStr)

	if len(match) >= 3 {
		num, _ := strconv.Atoi(match[1])
		unit := match[2]

		if duration, ok := patterns[unit]; ok {
			uploadTime := now.Add(-time.Duration(num) * duration)
			return uploadTime.Format("20060102")
		}
	}

	// Handle special cases
	if strings.Contains(relativeStr, "just now") || strings.Contains(relativeStr, "moments ago") {
		return now.Format("20060102")
	}

	if strings.Contains(relativeStr, "yesterday") {
		return now.AddDate(0, 0, -1).Format("20060102")
	}

	// If we can't parse, return today's date
	return now.Format("20060102")
}

// extractContinuationToken extracts continuation token from browse response
func extractContinuationToken(resp *BrowseResponse) string {
	if resp.Contents.TwoColumnBrowseResultsRenderer == nil {
		return ""
	}

	for _, tab := range resp.Contents.TwoColumnBrowseResultsRenderer.Tabs {
		if tab.TabRenderer == nil {
			continue
		}

		// Check RichGridRenderer
		if tab.TabRenderer.Content.RichGridRenderer != nil {
			grid := tab.TabRenderer.Content.RichGridRenderer

			// Check continuation in grid contents
			for _, content := range grid.Contents {
				if content.ContinuationItemRenderer != nil {
					return content.ContinuationItemRenderer.ContinuationEndpoint.ContinuationCommand.Token
				}
			}

			// Check continuation array
			for _, cont := range grid.Continuation {
				if cont.NextContinuationData != nil {
					return cont.NextContinuationData.Continuation
				}
				if cont.ReloadContinuationData != nil {
					return cont.ReloadContinuationData.Continuation
				}
			}
		}

		// Check SectionListRenderer
		if tab.TabRenderer.Content.SectionListRenderer != nil {
			for _, cont := range tab.TabRenderer.Content.SectionListRenderer.Continuation {
				if cont.NextContinuationData != nil {
					return cont.NextContinuationData.Continuation
				}
			}
		}
	}

	return ""
}

// extractContinuationTokenFromActions extracts continuation token from response actions
func extractContinuationTokenFromActions(resp *BrowseResponse) string {
	for _, action := range resp.OnResponseReceivedActions {
		if action.AppendContinuationItemsAction != nil {
			for _, content := range action.AppendContinuationItemsAction.ContinuationItems {
				if content.ContinuationItemRenderer != nil {
					return content.ContinuationItemRenderer.ContinuationEndpoint.ContinuationCommand.Token
				}
			}
		}

		if action.ReloadContinuationItemsCommand != nil {
			for _, content := range action.ReloadContinuationItemsCommand.ContinuationItems {
				if content.ContinuationItemRenderer != nil {
					return content.ContinuationItemRenderer.ContinuationEndpoint.ContinuationCommand.Token
				}
			}
		}
	}

	return ""
}

// parseVideoMetadataFromPlayerResponse extracts video metadata from player response
func parseVideoMetadataFromPlayerResponse(resp *PlayerResponse) *VideoMetadata {
	metadata := &VideoMetadata{
		Video: Video{
			Status: StatusPending,
		},
		Tags:       []string{},
		Categories: []string{},
		Subtitles:  []string{},
		Formats:    []Format{},
	}

	// Extract from VideoDetails
	if resp.VideoDetails != nil {
		vd := resp.VideoDetails
		metadata.ID = vd.VideoID
		metadata.Title = vd.Title
		metadata.Description = vd.ShortDescription
		metadata.ThumbnailURL = vd.Thumbnail.GetBestThumbnail()
		metadata.Tags = vd.Keywords

		// Parse duration
		if vd.LengthSeconds != "" {
			duration, _ := strconv.Atoi(vd.LengthSeconds)
			metadata.Duration = duration
		}

		// Parse view count
		if vd.ViewCount != "" {
			viewCount, _ := strconv.ParseInt(vd.ViewCount, 10, 64)
			metadata.ViewCount = viewCount
		}
	}

	// Extract from Microformat for more details
	if resp.Microformat != nil && resp.Microformat.PlayerMicroformatRenderer != nil {
		mf := resp.Microformat.PlayerMicroformatRenderer

		if metadata.Title == "" {
			metadata.Title = mf.Title.GetText()
		}
		if metadata.Description == "" {
			metadata.Description = mf.Description.GetText()
		}

		if mf.Category != "" {
			metadata.Categories = []string{mf.Category}
		}

		// Use upload date from microformat
		if mf.UploadDate != "" {
			metadata.UploadDate = strings.ReplaceAll(mf.UploadDate, "-", "")
		} else if mf.PublishDate != "" {
			metadata.UploadDate = strings.ReplaceAll(mf.PublishDate, "-", "")
		}

		if metadata.ThumbnailURL == "" {
			metadata.ThumbnailURL = mf.Thumbnail.GetBestThumbnail()
		}
	}

	// Extract formats from streaming data
	if resp.StreamingData != nil {
		allFormats := append(resp.StreamingData.Formats, resp.StreamingData.AdaptiveFormats...)

		for _, f := range allFormats {
			format := Format{
				FormatID:   strconv.Itoa(f.ITag),
				Extension:  extractExtension(f.MimeType),
				Resolution: "",
				VCodec:     extractCodec(f.MimeType, "video"),
				ACodec:     extractCodec(f.MimeType, "audio"),
				Quality:    f.QualityLabel,
			}

			if f.Width > 0 && f.Height > 0 {
				format.Resolution = fmt.Sprintf("%dx%d", f.Width, f.Height)
			}

			if f.ContentLength != "" {
				format.FileSize, _ = strconv.ParseInt(f.ContentLength, 10, 64)
			}

			if format.Quality == "" {
				format.Quality = f.Quality
			}

			metadata.Formats = append(metadata.Formats, format)
		}
	}

	// Extract subtitles
	if resp.Captions != nil && resp.Captions.PlayerCaptionsTracklistRenderer != nil {
		for _, track := range resp.Captions.PlayerCaptionsTracklistRenderer.CaptionTracks {
			lang := track.LanguageCode
			if track.Kind == "asr" {
				lang += " (auto)"
			}
			metadata.Subtitles = append(metadata.Subtitles, lang)
		}
	}

	// Generate thumbnail if not provided
	if metadata.ThumbnailURL == "" && metadata.ID != "" {
		metadata.ThumbnailURL = fmt.Sprintf("https://i.ytimg.com/vi/%s/maxresdefault.jpg", metadata.ID)
	}

	// Initialize nil slices
	if metadata.Tags == nil {
		metadata.Tags = []string{}
	}
	if metadata.Categories == nil {
		metadata.Categories = []string{}
	}
	if metadata.Subtitles == nil {
		metadata.Subtitles = []string{}
	}

	return metadata
}

// extractExtension extracts file extension from mime type
func extractExtension(mimeType string) string {
	// mimeType format: "video/mp4; codecs=\"avc1.4d401f\""
	if strings.HasPrefix(mimeType, "video/mp4") {
		return "mp4"
	}
	if strings.HasPrefix(mimeType, "video/webm") {
		return "webm"
	}
	if strings.HasPrefix(mimeType, "audio/mp4") {
		return "m4a"
	}
	if strings.HasPrefix(mimeType, "audio/webm") {
		return "webm"
	}
	return ""
}

// extractCodec extracts codec from mime type
func extractCodec(mimeType string, codecType string) string {
	// mimeType format: "video/mp4; codecs=\"avc1.4d401f, mp4a.40.2\""
	codecsMatch := regexp.MustCompile(`codecs="([^"]+)"`).FindStringSubmatch(mimeType)
	if len(codecsMatch) < 2 {
		return ""
	}

	codecs := strings.Split(codecsMatch[1], ",")
	for _, codec := range codecs {
		codec = strings.TrimSpace(codec)
		if codecType == "video" && (strings.HasPrefix(codec, "avc") || strings.HasPrefix(codec, "vp") || strings.HasPrefix(codec, "av0")) {
			return codec
		}
		if codecType == "audio" && (strings.HasPrefix(codec, "mp4a") || strings.HasPrefix(codec, "opus") || strings.HasPrefix(codec, "vorbis")) {
			return codec
		}
	}

	return ""
}

// Legacy parsing functions for backward compatibility (used by tests)

// ytdlpChannelOutput represents the JSON structure returned by yt-dlp for channel info (legacy)
type ytdlpChannelOutput struct {
	ID            string `json:"id"`
	ChannelID     string `json:"channel_id"`
	Channel       string `json:"channel"`
	Uploader      string `json:"uploader"`
	UploaderID    string `json:"uploader_id"`
	UploaderURL   string `json:"uploader_url"`
	ChannelURL    string `json:"channel_url"`
	Description   string `json:"description"`
	PlaylistCount int    `json:"playlist_count"`
	Thumbnails    []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
		ID     string `json:"id"`
	} `json:"thumbnails"`
}

// ytdlpVideoEntry represents a video entry in flat-playlist output (legacy)
type ytdlpVideoEntry struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Duration    float64 `json:"duration"`
	UploadDate  string  `json:"upload_date"`
	ViewCount   int64   `json:"view_count"`
	URL         string  `json:"url"`
	Thumbnails  []struct {
		URL string `json:"url"`
	} `json:"thumbnails"`
}

// ytdlpVideoMetadata represents the full JSON structure for a single video (legacy)
type ytdlpVideoMetadata struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Duration    float64 `json:"duration"`
	UploadDate  string  `json:"upload_date"`
	ViewCount   int64   `json:"view_count"`
	Thumbnail   string  `json:"thumbnail"`
	Thumbnails  []struct {
		URL string `json:"url"`
	} `json:"thumbnails"`
	Tags       []string `json:"tags"`
	Categories []string `json:"categories"`
	Subtitles  map[string][]struct {
		Ext string `json:"ext"`
		URL string `json:"url"`
	} `json:"subtitles"`
	AutomaticCaptions map[string][]struct {
		Ext string `json:"ext"`
		URL string `json:"url"`
	} `json:"automatic_captions"`
	Formats []struct {
		FormatID       string  `json:"format_id"`
		Ext            string  `json:"ext"`
		Resolution     string  `json:"resolution"`
		FileSize       int64   `json:"filesize"`
		FileSizeApprox int64   `json:"filesize_approx"`
		VCodec         string  `json:"vcodec"`
		ACodec         string  `json:"acodec"`
		Quality        float64 `json:"quality"`
		Height         int     `json:"height"`
		Width          int     `json:"width"`
	} `json:"formats"`
}

// ParseChannelJSON parses yt-dlp JSON output for channel information (legacy compatibility)
func ParseChannelJSON(data []byte) (*Channel, error) {
	var raw ytdlpChannelOutput
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse channel JSON: %w", err)
	}

	channel := &Channel{
		ID:          raw.ChannelID,
		Name:        raw.Channel,
		Description: raw.Description,
		VideoCount:  raw.PlaylistCount,
		URL:         raw.ChannelURL,
		CreatedAt:   time.Now(),
	}

	// Use uploader info as fallback
	if channel.ID == "" {
		channel.ID = raw.UploaderID
	}
	if channel.Name == "" {
		channel.Name = raw.Uploader
	}
	if channel.URL == "" {
		channel.URL = raw.UploaderURL
	}

	// Extract avatar and banner from thumbnails
	for _, thumb := range raw.Thumbnails {
		if thumb.ID == "avatar_uncropped" || strings.Contains(thumb.URL, "yt3.ggpht.com") {
			if channel.AvatarURL == "" {
				channel.AvatarURL = thumb.URL
			}
		}
		if thumb.ID == "banner_uncropped" || strings.Contains(thumb.URL, "banner") {
			if channel.BannerURL == "" {
				channel.BannerURL = thumb.URL
			}
		}
	}

	// If no specific avatar found, use the first small thumbnail
	if channel.AvatarURL == "" && len(raw.Thumbnails) > 0 {
		for _, thumb := range raw.Thumbnails {
			if thumb.Height <= 200 && thumb.Height > 0 {
				channel.AvatarURL = thumb.URL
				break
			}
		}
	}

	return channel, nil
}

// ParseVideoList parses yt-dlp flat-playlist JSON output for video list (legacy compatibility)
func ParseVideoList(data []byte) ([]Video, error) {
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	videos := make([]Video, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var entry ytdlpVideoEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		video := Video{
			ID:          entry.ID,
			Title:       entry.Title,
			Description: entry.Description,
			Duration:    int(entry.Duration),
			UploadDate:  entry.UploadDate,
			ViewCount:   entry.ViewCount,
			Status:      StatusPending,
		}

		if len(entry.Thumbnails) > 0 {
			video.ThumbnailURL = entry.Thumbnails[0].URL
		}

		if video.ThumbnailURL == "" && video.ID != "" {
			video.ThumbnailURL = fmt.Sprintf("https://i.ytimg.com/vi/%s/maxresdefault.jpg", video.ID)
		}

		videos = append(videos, video)
	}

	if len(videos) == 0 && len(data) > 0 {
		return nil, fmt.Errorf("failed to parse any videos from data")
	}

	return videos, nil
}

// ParseVideoMetadata parses yt-dlp JSON output for detailed video metadata (legacy compatibility)
func ParseVideoMetadata(data []byte) (*VideoMetadata, error) {
	var raw ytdlpVideoMetadata
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse video metadata JSON: %w", err)
	}

	metadata := &VideoMetadata{
		Video: Video{
			ID:          raw.ID,
			Title:       raw.Title,
			Description: raw.Description,
			Duration:    int(raw.Duration),
			UploadDate:  raw.UploadDate,
			ViewCount:   raw.ViewCount,
			Status:      StatusPending,
		},
		Tags:       raw.Tags,
		Categories: raw.Categories,
	}

	if raw.Thumbnail != "" {
		metadata.ThumbnailURL = raw.Thumbnail
	} else if len(raw.Thumbnails) > 0 {
		metadata.ThumbnailURL = raw.Thumbnails[0].URL
	} else if raw.ID != "" {
		metadata.ThumbnailURL = fmt.Sprintf("https://i.ytimg.com/vi/%s/maxresdefault.jpg", raw.ID)
	}

	metadata.Formats = make([]Format, 0, len(raw.Formats))
	for _, f := range raw.Formats {
		format := Format{
			FormatID:   f.FormatID,
			Extension:  f.Ext,
			Resolution: f.Resolution,
			FileSize:   f.FileSize,
			VCodec:     f.VCodec,
			ACodec:     f.ACodec,
			Quality:    fmt.Sprintf("%.0f", f.Quality),
		}

		if format.FileSize == 0 {
			format.FileSize = f.FileSizeApprox
		}

		if format.Resolution == "" && f.Width > 0 && f.Height > 0 {
			format.Resolution = fmt.Sprintf("%dx%d", f.Width, f.Height)
		}

		metadata.Formats = append(metadata.Formats, format)
	}

	subtitleLangs := make(map[string]bool)
	for lang := range raw.Subtitles {
		subtitleLangs[lang] = true
	}
	for lang := range raw.AutomaticCaptions {
		subtitleLangs[lang+" (auto)"] = true
	}

	metadata.Subtitles = make([]string, 0, len(subtitleLangs))
	for lang := range subtitleLangs {
		metadata.Subtitles = append(metadata.Subtitles, lang)
	}

	if metadata.Tags == nil {
		metadata.Tags = []string{}
	}
	if metadata.Categories == nil {
		metadata.Categories = []string{}
	}

	return metadata, nil
}
