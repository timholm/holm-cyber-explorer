package youtube

import "time"

// Channel represents a YouTube channel with its metadata
type Channel struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AvatarURL   string    `json:"avatar_url"`
	BannerURL   string    `json:"banner_url"`
	VideoCount  int       `json:"video_count"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
}

// Video represents a YouTube video with basic metadata
type Video struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Duration     int    `json:"duration"`
	UploadDate   string `json:"upload_date"`
	ThumbnailURL string `json:"thumbnail_url"`
	ViewCount    int64  `json:"view_count"`
	Status       string `json:"status"` // pending, downloading, completed, failed
}

// Format represents a video format option available for download
type Format struct {
	FormatID   string `json:"format_id"`
	Extension  string `json:"ext"`
	Resolution string `json:"resolution"`
	FileSize   int64  `json:"filesize"`
	VCodec     string `json:"vcodec"`
	ACodec     string `json:"acodec"`
	Quality    string `json:"quality"`
}

// VideoMetadata extends Video with additional detailed information
type VideoMetadata struct {
	Video
	Formats    []Format `json:"formats"`
	Subtitles  []string `json:"subtitles"`
	Tags       []string `json:"tags"`
	Categories []string `json:"categories"`
}

// StreamInfo contains stream URLs and format information for downloading
type StreamInfo struct {
	VideoID         string               `json:"video_id"`
	Title           string               `json:"title"`
	ExpiresIn       time.Duration        `json:"expires_in"`
	Formats         []DownloadableFormat `json:"formats"`
	HLSManifestURL  string               `json:"hls_manifest_url,omitempty"`
	DASHManifestURL string               `json:"dash_manifest_url,omitempty"`
}

// DownloadableFormat represents a single downloadable stream format with URL
type DownloadableFormat struct {
	ITag            int    `json:"itag"`
	URL             string `json:"url"`
	MimeType        string `json:"mime_type"`
	Bitrate         int    `json:"bitrate"`
	Width           int    `json:"width,omitempty"`
	Height          int    `json:"height,omitempty"`
	ContentLength   int64  `json:"content_length,omitempty"`
	Quality         string `json:"quality"`
	QualityLabel    string `json:"quality_label,omitempty"`
	AudioQuality    string `json:"audio_quality,omitempty"`
	AudioSampleRate string `json:"audio_sample_rate,omitempty"`
	AudioChannels   int    `json:"audio_channels,omitempty"`
	FPS             int    `json:"fps,omitempty"`
}

// GetBestVideoFormat returns the best video format (highest resolution)
func (s *StreamInfo) GetBestVideoFormat() *DownloadableFormat {
	var best *DownloadableFormat
	for i := range s.Formats {
		f := &s.Formats[i]
		// Skip audio-only formats
		if f.Width == 0 && f.Height == 0 {
			continue
		}
		if best == nil || (f.Height > best.Height) || (f.Height == best.Height && f.Bitrate > best.Bitrate) {
			best = f
		}
	}
	return best
}

// GetBestAudioFormat returns the best audio format (highest bitrate)
func (s *StreamInfo) GetBestAudioFormat() *DownloadableFormat {
	var best *DownloadableFormat
	for i := range s.Formats {
		f := &s.Formats[i]
		// Only audio formats (no width/height)
		if f.Width > 0 || f.Height > 0 {
			continue
		}
		if f.AudioQuality == "" {
			continue
		}
		if best == nil || f.Bitrate > best.Bitrate {
			best = f
		}
	}
	return best
}

// GetFormatByITag returns a format by its itag
func (s *StreamInfo) GetFormatByITag(itag int) *DownloadableFormat {
	for i := range s.Formats {
		if s.Formats[i].ITag == itag {
			return &s.Formats[i]
		}
	}
	return nil
}

// GetCombinedFormats returns formats that have both video and audio
func (s *StreamInfo) GetCombinedFormats() []DownloadableFormat {
	var combined []DownloadableFormat
	for _, f := range s.Formats {
		// Combined formats have both video dimensions and audio quality
		if f.Width > 0 && f.Height > 0 && f.AudioChannels > 0 {
			combined = append(combined, f)
		}
	}
	return combined
}

// VideoStatus constants
const (
	StatusPending     = "pending"
	StatusDownloading = "downloading"
	StatusCompleted   = "completed"
	StatusFailed      = "failed"
)
