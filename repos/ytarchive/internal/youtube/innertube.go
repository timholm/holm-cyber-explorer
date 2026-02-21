package youtube

// Innertube API request/response types for YouTube's internal API
// These types mirror the JSON structures used by YouTube's web player

// InnertubeContext represents the client context required for innertube API requests
type InnertubeContext struct {
	Client InnertubeClient `json:"client"`
}

// InnertubeClient represents client information for innertube requests
type InnertubeClient struct {
	HL                string `json:"hl"`
	GL                string `json:"gl"`
	ClientName        string `json:"clientName"`
	ClientVersion     string `json:"clientVersion"`
	AndroidSDKVersion int    `json:"androidSdkVersion,omitempty"`
}

// BrowseRequest is the request body for the /browse endpoint
type BrowseRequest struct {
	Context      InnertubeContext `json:"context"`
	BrowseID     string           `json:"browseId,omitempty"`
	Params       string           `json:"params,omitempty"`
	Continuation string           `json:"continuation,omitempty"`
}

// PlayerRequest is the request body for the /player endpoint
type PlayerRequest struct {
	Context         InnertubeContext `json:"context"`
	VideoID         string           `json:"videoId"`
	PlaybackContext *PlaybackContext `json:"playbackContext,omitempty"`
}

// PlaybackContext provides additional context for player requests
type PlaybackContext struct {
	ContentPlaybackContext ContentPlaybackContext `json:"contentPlaybackContext"`
}

// ContentPlaybackContext contains playback settings
type ContentPlaybackContext struct {
	SignatureTimestamp int `json:"signatureTimestamp,omitempty"`
}

// BrowseResponse is the response from the /browse endpoint
type BrowseResponse struct {
	Header                    ChannelHeader              `json:"header,omitempty"`
	Metadata                  ChannelMetadata            `json:"metadata,omitempty"`
	Contents                  BrowseContents             `json:"contents,omitempty"`
	OnResponseReceivedActions []OnResponseReceivedAction `json:"onResponseReceivedActions,omitempty"`
}

// ChannelHeader contains channel header information
type ChannelHeader struct {
	C4TabbedHeaderRenderer *C4TabbedHeaderRenderer `json:"c4TabbedHeaderRenderer,omitempty"`
	PageHeaderRenderer     *PageHeaderRenderer     `json:"pageHeaderRenderer,omitempty"`
}

// C4TabbedHeaderRenderer is the channel header format
type C4TabbedHeaderRenderer struct {
	ChannelID       string        `json:"channelId"`
	Title           string        `json:"title"`
	Avatar          ThumbnailList `json:"avatar,omitempty"`
	Banner          ThumbnailList `json:"banner,omitempty"`
	SubscriberCount SimpleText    `json:"subscriberCountText,omitempty"`
	VideosCount     SimpleText    `json:"videosCountText,omitempty"`
}

// PageHeaderRenderer is an alternative header format for newer channels
type PageHeaderRenderer struct {
	PageTitle string            `json:"pageTitle"`
	Content   PageHeaderContent `json:"content,omitempty"`
}

// PageHeaderContent contains the content of a page header
type PageHeaderContent struct {
	PageHeaderViewModel *PageHeaderViewModel `json:"pageHeaderViewModel,omitempty"`
}

// PageHeaderViewModel is a newer channel header format
type PageHeaderViewModel struct {
	Title       TitleContainer       `json:"title,omitempty"`
	Image       ImageContainer       `json:"image,omitempty"`
	Banner      BannerContainer      `json:"banner,omitempty"`
	Description DescriptionContainer `json:"description,omitempty"`
}

// TitleContainer holds the title in a dynamic text view model
type TitleContainer struct {
	DynamicTextViewModel *DynamicTextViewModel `json:"dynamicTextViewModel,omitempty"`
}

// DynamicTextViewModel contains dynamic text
type DynamicTextViewModel struct {
	Text TextContent `json:"text,omitempty"`
}

// TextContent contains text and style
type TextContent struct {
	Content string `json:"content"`
}

// ImageContainer holds avatar images
type ImageContainer struct {
	DecoratedAvatarViewModel *DecoratedAvatarViewModel `json:"decoratedAvatarViewModel,omitempty"`
}

// DecoratedAvatarViewModel contains avatar info
type DecoratedAvatarViewModel struct {
	Avatar AvatarContainer `json:"avatar,omitempty"`
}

// AvatarContainer holds the avatar view model
type AvatarContainer struct {
	AvatarViewModel *AvatarViewModel `json:"avatarViewModel,omitempty"`
}

// AvatarViewModel contains avatar image
type AvatarViewModel struct {
	Image ThumbnailList `json:"image,omitempty"`
}

// BannerContainer holds banner images
type BannerContainer struct {
	ImageBannerViewModel *ImageBannerViewModel `json:"imageBannerViewModel,omitempty"`
}

// ImageBannerViewModel contains banner image
type ImageBannerViewModel struct {
	Image ThumbnailList `json:"image,omitempty"`
}

// DescriptionContainer holds channel description
type DescriptionContainer struct {
	DescriptionPreviewViewModel *DescriptionPreviewViewModel `json:"descriptionPreviewViewModel,omitempty"`
}

// DescriptionPreviewViewModel contains description
type DescriptionPreviewViewModel struct {
	Description TextRuns `json:"description,omitempty"`
}

// TextRuns contains text content with runs
type TextRuns struct {
	Content string    `json:"content,omitempty"`
	Runs    []TextRun `json:"runs,omitempty"`
}

// ChannelMetadata contains channel metadata
type ChannelMetadata struct {
	ChannelMetadataRenderer *ChannelMetadataRenderer `json:"channelMetadataRenderer,omitempty"`
}

// ChannelMetadataRenderer contains channel metadata details
type ChannelMetadataRenderer struct {
	Title            string        `json:"title"`
	Description      string        `json:"description"`
	ExternalID       string        `json:"externalId"`
	VanityChannelUrl string        `json:"vanityChannelUrl"`
	ChannelUrl       string        `json:"channelUrl"`
	Avatar           ThumbnailList `json:"avatar,omitempty"`
}

// BrowseContents contains the main content of a browse response
type BrowseContents struct {
	TwoColumnBrowseResultsRenderer *TwoColumnBrowseResultsRenderer `json:"twoColumnBrowseResultsRenderer,omitempty"`
}

// TwoColumnBrowseResultsRenderer is the main content renderer
type TwoColumnBrowseResultsRenderer struct {
	Tabs []Tab `json:"tabs,omitempty"`
}

// Tab represents a channel tab (Videos, Playlists, etc.)
type Tab struct {
	TabRenderer *TabRenderer `json:"tabRenderer,omitempty"`
}

// TabRenderer contains tab content
type TabRenderer struct {
	Title    string     `json:"title,omitempty"`
	Content  TabContent `json:"content,omitempty"`
	Endpoint Endpoint   `json:"endpoint,omitempty"`
}

// TabContent contains the content of a tab
type TabContent struct {
	RichGridRenderer    *RichGridRenderer    `json:"richGridRenderer,omitempty"`
	SectionListRenderer *SectionListRenderer `json:"sectionListRenderer,omitempty"`
}

// RichGridRenderer renders a grid of videos
type RichGridRenderer struct {
	Contents     []RichGridContent `json:"contents,omitempty"`
	Header       *RichGridHeader   `json:"header,omitempty"`
	Continuation []Continuation    `json:"continuation,omitempty"`
}

// RichGridHeader contains header info for the grid
type RichGridHeader struct {
	FeedFilterChipBarRenderer *FeedFilterChipBarRenderer `json:"feedFilterChipBarRenderer,omitempty"`
}

// FeedFilterChipBarRenderer renders filter chips
type FeedFilterChipBarRenderer struct {
	Contents []ChipContent `json:"contents,omitempty"`
}

// ChipContent contains filter chip data
type ChipContent struct {
	ChipCloudChipRenderer *ChipCloudChipRenderer `json:"chipCloudChipRenderer,omitempty"`
}

// ChipCloudChipRenderer renders a filter chip
type ChipCloudChipRenderer struct {
	Text       SimpleText `json:"text,omitempty"`
	IsSelected bool       `json:"isSelected"`
}

// RichGridContent is a single item in a rich grid
type RichGridContent struct {
	RichItemRenderer         *RichItemRenderer         `json:"richItemRenderer,omitempty"`
	ContinuationItemRenderer *ContinuationItemRenderer `json:"continuationItemRenderer,omitempty"`
}

// RichItemRenderer renders a single video item
type RichItemRenderer struct {
	Content RichItemContent `json:"content,omitempty"`
}

// RichItemContent contains the video renderer
type RichItemContent struct {
	VideoRenderer         *VideoRenderer         `json:"videoRenderer,omitempty"`
	ReelItemRenderer      *ReelItemRenderer      `json:"reelItemRenderer,omitempty"`
	ShortsLockupViewModel *ShortsLockupViewModel `json:"shortsLockupViewModel,omitempty"`
}

// ShortsLockupViewModel represents a YouTube Short
type ShortsLockupViewModel struct {
	EntityID string      `json:"entityId"`
	OnTap    OnTapAction `json:"onTap,omitempty"`
}

// OnTapAction contains the action when tapping on a short
type OnTapAction struct {
	InnertubeCommand InnertubeCommand `json:"innertubeCommand,omitempty"`
}

// InnertubeCommand contains command details
type InnertubeCommand struct {
	ReelWatchEndpoint *ReelWatchEndpoint `json:"reelWatchEndpoint,omitempty"`
}

// ReelWatchEndpoint contains the video ID for shorts
type ReelWatchEndpoint struct {
	VideoID string `json:"videoId"`
}

// ReelItemRenderer renders a short video (reel)
type ReelItemRenderer struct {
	VideoID   string        `json:"videoId"`
	Headline  SimpleText    `json:"headline,omitempty"`
	Thumbnail ThumbnailList `json:"thumbnail,omitempty"`
}

// VideoRenderer contains video information
type VideoRenderer struct {
	VideoID            string        `json:"videoId"`
	Title              TextRuns      `json:"title,omitempty"`
	DescriptionSnippet TextRuns      `json:"descriptionSnippet,omitempty"`
	LengthText         SimpleText    `json:"lengthText,omitempty"`
	ViewCountText      SimpleText    `json:"viewCountText,omitempty"`
	PublishedTimeText  SimpleText    `json:"publishedTimeText,omitempty"`
	Thumbnail          ThumbnailList `json:"thumbnail,omitempty"`
}

// SectionListRenderer renders a list of sections
type SectionListRenderer struct {
	Contents     []SectionContent `json:"contents,omitempty"`
	Continuation []Continuation   `json:"continuations,omitempty"`
}

// SectionContent is content within a section
type SectionContent struct {
	ItemSectionRenderer *ItemSectionRenderer `json:"itemSectionRenderer,omitempty"`
}

// ItemSectionRenderer renders items in a section
type ItemSectionRenderer struct {
	Contents []ItemContent `json:"contents,omitempty"`
}

// ItemContent is a single item within a section
type ItemContent struct {
	GridRenderer *GridRenderer `json:"gridRenderer,omitempty"`
}

// GridRenderer renders a grid of items
type GridRenderer struct {
	Items        []GridItem     `json:"items,omitempty"`
	Continuation []Continuation `json:"continuations,omitempty"`
}

// GridItem is a single item in a grid
type GridItem struct {
	GridVideoRenderer *GridVideoRenderer `json:"gridVideoRenderer,omitempty"`
}

// GridVideoRenderer contains grid video information
type GridVideoRenderer struct {
	VideoID           string        `json:"videoId"`
	Title             SimpleText    `json:"title,omitempty"`
	ViewCountText     SimpleText    `json:"viewCountText,omitempty"`
	PublishedTimeText SimpleText    `json:"publishedTimeText,omitempty"`
	Thumbnail         ThumbnailList `json:"thumbnail,omitempty"`
}

// ContinuationItemRenderer contains continuation data for pagination
type ContinuationItemRenderer struct {
	ContinuationEndpoint ContinuationEndpoint `json:"continuationEndpoint,omitempty"`
}

// ContinuationEndpoint contains the continuation token
type ContinuationEndpoint struct {
	ContinuationCommand ContinuationCommand `json:"continuationCommand,omitempty"`
}

// ContinuationCommand contains the actual continuation token
type ContinuationCommand struct {
	Token string `json:"token"`
}

// Continuation holds continuation data
type Continuation struct {
	NextContinuationData   *NextContinuationData `json:"nextContinuationData,omitempty"`
	ReloadContinuationData *NextContinuationData `json:"reloadContinuationData,omitempty"`
}

// NextContinuationData contains the next page token
type NextContinuationData struct {
	Continuation string `json:"continuation"`
}

// OnResponseReceivedAction contains actions after a continuation request
type OnResponseReceivedAction struct {
	AppendContinuationItemsAction  *AppendContinuationItemsAction  `json:"appendContinuationItemsAction,omitempty"`
	ReloadContinuationItemsCommand *ReloadContinuationItemsCommand `json:"reloadContinuationItemsCommand,omitempty"`
}

// AppendContinuationItemsAction contains items to append
type AppendContinuationItemsAction struct {
	ContinuationItems []RichGridContent `json:"continuationItems,omitempty"`
}

// ReloadContinuationItemsCommand contains items for reload
type ReloadContinuationItemsCommand struct {
	ContinuationItems []RichGridContent `json:"continuationItems,omitempty"`
}

// Endpoint contains navigation endpoint data
type Endpoint struct {
	BrowseEndpoint *BrowseEndpoint `json:"browseEndpoint,omitempty"`
}

// BrowseEndpoint contains browse navigation data
type BrowseEndpoint struct {
	BrowseID string `json:"browseId"`
	Params   string `json:"params,omitempty"`
}

// SimpleText is a text container with simpleText field
type SimpleText struct {
	SimpleText string    `json:"simpleText,omitempty"`
	Runs       []TextRun `json:"runs,omitempty"`
}

// TextRun is a single text run
type TextRun struct {
	Text string `json:"text"`
}

// ThumbnailList contains a list of thumbnails
type ThumbnailList struct {
	Thumbnails []Thumbnail `json:"thumbnails,omitempty"`
	Sources    []Thumbnail `json:"sources,omitempty"`
}

// Thumbnail represents a single thumbnail image
type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// PlayerResponse is the response from the /player endpoint
type PlayerResponse struct {
	VideoDetails      *VideoDetails      `json:"videoDetails,omitempty"`
	StreamingData     *StreamingData     `json:"streamingData,omitempty"`
	PlayabilityStatus *PlayabilityStatus `json:"playabilityStatus,omitempty"`
	Captions          *Captions          `json:"captions,omitempty"`
	Microformat       *Microformat       `json:"microformat,omitempty"`
}

// VideoDetails contains detailed video information
type VideoDetails struct {
	VideoID          string        `json:"videoId"`
	Title            string        `json:"title"`
	LengthSeconds    string        `json:"lengthSeconds"`
	ChannelID        string        `json:"channelId"`
	ShortDescription string        `json:"shortDescription"`
	Thumbnail        ThumbnailList `json:"thumbnail,omitempty"`
	ViewCount        string        `json:"viewCount"`
	Author           string        `json:"author"`
	Keywords         []string      `json:"keywords,omitempty"`
	IsLiveContent    bool          `json:"isLiveContent"`
	IsPrivate        bool          `json:"isPrivate"`
}

// StreamingData contains stream URLs and formats
type StreamingData struct {
	ExpiresInSeconds string         `json:"expiresInSeconds"`
	Formats          []StreamFormat `json:"formats,omitempty"`
	AdaptiveFormats  []StreamFormat `json:"adaptiveFormats,omitempty"`
	HLSManifestURL   string         `json:"hlsManifestUrl,omitempty"`
	DashManifestURL  string         `json:"dashManifestUrl,omitempty"`
}

// StreamFormat represents a single stream format
type StreamFormat struct {
	ITag             int    `json:"itag"`
	URL              string `json:"url,omitempty"`
	SignatureCipher  string `json:"signatureCipher,omitempty"`
	MimeType         string `json:"mimeType"`
	Bitrate          int    `json:"bitrate"`
	Width            int    `json:"width,omitempty"`
	Height           int    `json:"height,omitempty"`
	ContentLength    string `json:"contentLength,omitempty"`
	Quality          string `json:"quality"`
	QualityLabel     string `json:"qualityLabel,omitempty"`
	AudioQuality     string `json:"audioQuality,omitempty"`
	AudioSampleRate  string `json:"audioSampleRate,omitempty"`
	AudioChannels    int    `json:"audioChannels,omitempty"`
	AverageBitrate   int    `json:"averageBitrate,omitempty"`
	FPS              int    `json:"fps,omitempty"`
	ProjectionType   string `json:"projectionType,omitempty"`
	ApproxDurationMs string `json:"approxDurationMs,omitempty"`
}

// PlayabilityStatus contains video availability information
type PlayabilityStatus struct {
	Status          string `json:"status"`
	Reason          string `json:"reason,omitempty"`
	PlayableInEmbed bool   `json:"playableInEmbed"`
}

// Captions contains caption/subtitle information
type Captions struct {
	PlayerCaptionsTracklistRenderer *PlayerCaptionsTracklistRenderer `json:"playerCaptionsTracklistRenderer,omitempty"`
}

// PlayerCaptionsTracklistRenderer contains caption tracks
type PlayerCaptionsTracklistRenderer struct {
	CaptionTracks []CaptionTrack `json:"captionTracks,omitempty"`
}

// CaptionTrack represents a single caption track
type CaptionTrack struct {
	BaseURL      string     `json:"baseUrl"`
	Name         SimpleText `json:"name,omitempty"`
	VssID        string     `json:"vssId"`
	LanguageCode string     `json:"languageCode"`
	Kind         string     `json:"kind,omitempty"`
}

// Microformat contains additional video metadata
type Microformat struct {
	PlayerMicroformatRenderer *PlayerMicroformatRenderer `json:"playerMicroformatRenderer,omitempty"`
}

// PlayerMicroformatRenderer contains microformat metadata
type PlayerMicroformatRenderer struct {
	Title             SimpleText    `json:"title,omitempty"`
	Description       SimpleText    `json:"description,omitempty"`
	LengthSeconds     string        `json:"lengthSeconds"`
	OwnerChannelName  string        `json:"ownerChannelName"`
	ExternalChannelID string        `json:"externalChannelId"`
	ViewCount         string        `json:"viewCount"`
	Category          string        `json:"category"`
	PublishDate       string        `json:"publishDate"`
	UploadDate        string        `json:"uploadDate"`
	Thumbnail         ThumbnailList `json:"thumbnail,omitempty"`
}

// ResolveURLRequest is the request body for resolving URLs
type ResolveURLRequest struct {
	Context InnertubeContext `json:"context"`
	URL     string           `json:"url"`
}

// ResolveURLResponse is the response from URL resolution
type ResolveURLResponse struct {
	Endpoint ResolvedEndpoint `json:"endpoint,omitempty"`
}

// ResolvedEndpoint contains the resolved endpoint
type ResolvedEndpoint struct {
	BrowseEndpoint *BrowseEndpoint `json:"browseEndpoint,omitempty"`
}

// GetText extracts text from SimpleText, handling both simpleText and runs formats
func (st SimpleText) GetText() string {
	if st.SimpleText != "" {
		return st.SimpleText
	}
	if len(st.Runs) > 0 {
		text := ""
		for _, run := range st.Runs {
			text += run.Text
		}
		return text
	}
	return ""
}

// GetText extracts text from TextRuns, handling both content and runs formats
func (tr TextRuns) GetText() string {
	if tr.Content != "" {
		return tr.Content
	}
	if len(tr.Runs) > 0 {
		text := ""
		for _, run := range tr.Runs {
			text += run.Text
		}
		return text
	}
	return ""
}

// GetBestThumbnail returns the highest quality thumbnail URL
func (tl ThumbnailList) GetBestThumbnail() string {
	thumbnails := tl.Thumbnails
	if len(thumbnails) == 0 {
		thumbnails = tl.Sources
	}
	if len(thumbnails) == 0 {
		return ""
	}

	// Find the largest thumbnail
	best := thumbnails[0]
	for _, thumb := range thumbnails {
		if thumb.Width > best.Width || thumb.Height > best.Height {
			best = thumb
		}
	}
	return best.URL
}
