package youtube

import (
	"testing"
)

func TestNormalizeChannelURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "handle with @ prefix",
			input:    "@aperturethinking",
			expected: "https://www.youtube.com/@aperturethinking",
		},
		{
			name:     "handle without @ prefix",
			input:    "aperturethinking",
			expected: "https://www.youtube.com/@aperturethinking",
		},
		{
			name:     "full URL with @ handle",
			input:    "https://youtube.com/@aperturethinking",
			expected: "https://www.youtube.com/@aperturethinking",
		},
		{
			name:     "full URL with www",
			input:    "https://www.youtube.com/@aperturethinking",
			expected: "https://www.youtube.com/@aperturethinking",
		},
		{
			name:     "URL without scheme",
			input:    "youtube.com/@aperturethinking",
			expected: "https://www.youtube.com/@aperturethinking",
		},
		{
			name:     "channel URL with channel ID",
			input:    "https://youtube.com/channel/UCxyz123abc",
			expected: "https://www.youtube.com/channel/UCxyz123abc",
		},
		{
			name:     "URL with leading/trailing whitespace",
			input:    "  @aperturethinking  ",
			expected: "https://www.youtube.com/@aperturethinking",
		},
		{
			name:     "http URL converted",
			input:    "http://youtube.com/@test",
			expected: "http://www.youtube.com/@test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeChannelURL(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeChannelURL(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseChannelJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		wantID      string
		wantName    string
		wantURL     string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid channel JSON with channel_id",
			input: []byte(`{
				"channel_id": "UCxyz123abc456def789ghi",
				"channel": "Test Channel",
				"channel_url": "https://www.youtube.com/channel/UCxyz123abc456def789ghi",
				"description": "A test channel",
				"playlist_count": 100,
				"thumbnails": [
					{"url": "https://yt3.ggpht.com/avatar.jpg", "id": "avatar_uncropped", "height": 100, "width": 100}
				]
			}`),
			wantID:   "UCxyz123abc456def789ghi",
			wantName: "Test Channel",
			wantURL:  "https://www.youtube.com/channel/UCxyz123abc456def789ghi",
			wantErr:  false,
		},
		{
			name: "fallback to uploader fields",
			input: []byte(`{
				"uploader_id": "UploaderID123",
				"uploader": "Uploader Name",
				"uploader_url": "https://www.youtube.com/@uploader",
				"description": "A test channel"
			}`),
			wantID:   "UploaderID123",
			wantName: "Uploader Name",
			wantURL:  "https://www.youtube.com/@uploader",
			wantErr:  false,
		},
		{
			name: "channel fields take precedence",
			input: []byte(`{
				"channel_id": "ChannelID",
				"channel": "Channel Name",
				"channel_url": "https://channel.url",
				"uploader_id": "UploaderID",
				"uploader": "Uploader Name",
				"uploader_url": "https://uploader.url"
			}`),
			wantID:   "ChannelID",
			wantName: "Channel Name",
			wantURL:  "https://channel.url",
			wantErr:  false,
		},
		{
			name:        "invalid JSON",
			input:       []byte(`{invalid json`),
			wantErr:     true,
			errContains: "failed to parse channel JSON",
		},
		{
			name:     "empty JSON object",
			input:    []byte(`{}`),
			wantID:   "",
			wantName: "",
			wantURL:  "",
			wantErr:  false,
		},
		{
			name: "extract avatar from thumbnails",
			input: []byte(`{
				"channel_id": "TestID",
				"channel": "Test",
				"thumbnails": [
					{"url": "https://yt3.ggpht.com/avatar.jpg", "height": 88, "width": 88},
					{"url": "https://banner.url/banner.jpg", "id": "banner_uncropped", "height": 1080, "width": 1920}
				]
			}`),
			wantID:   "TestID",
			wantName: "Test",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channel, err := ParseChannelJSON(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseChannelJSON() expected error, got nil")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("ParseChannelJSON() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseChannelJSON() unexpected error: %v", err)
				return
			}

			if channel.ID != tt.wantID {
				t.Errorf("ParseChannelJSON() ID = %q, want %q", channel.ID, tt.wantID)
			}
			if channel.Name != tt.wantName {
				t.Errorf("ParseChannelJSON() Name = %q, want %q", channel.Name, tt.wantName)
			}
			if channel.URL != tt.wantURL {
				t.Errorf("ParseChannelJSON() URL = %q, want %q", channel.URL, tt.wantURL)
			}
		})
	}
}

func TestParseVideoList(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		wantCount   int
		wantIDs     []string
		wantTitles  []string
		wantErr     bool
		errContains string
	}{
		{
			name: "single video",
			input: []byte(`{"id": "abc123", "title": "Test Video", "duration": 120, "upload_date": "20240101", "view_count": 1000}
`),
			wantCount:  1,
			wantIDs:    []string{"abc123"},
			wantTitles: []string{"Test Video"},
			wantErr:    false,
		},
		{
			name: "multiple videos",
			input: []byte(`{"id": "video1", "title": "First Video", "duration": 100}
{"id": "video2", "title": "Second Video", "duration": 200}
{"id": "video3", "title": "Third Video", "duration": 300}
`),
			wantCount:  3,
			wantIDs:    []string{"video1", "video2", "video3"},
			wantTitles: []string{"First Video", "Second Video", "Third Video"},
			wantErr:    false,
		},
		{
			name: "with thumbnails",
			input: []byte(`{"id": "vid1", "title": "Video With Thumb", "thumbnails": [{"url": "https://thumb.jpg"}]}
`),
			wantCount:  1,
			wantIDs:    []string{"vid1"},
			wantTitles: []string{"Video With Thumb"},
			wantErr:    false,
		},
		{
			name:      "empty input",
			input:     []byte(``),
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "whitespace only returns error",
			input: []byte(`
   `),
			wantErr:     true,
			errContains: "failed to parse any videos",
		},
		{
			name: "skip malformed entries",
			input: []byte(`{"id": "good1", "title": "Good Video"}
{invalid json}
{"id": "good2", "title": "Another Good Video"}
`),
			wantCount:  2,
			wantIDs:    []string{"good1", "good2"},
			wantTitles: []string{"Good Video", "Another Good Video"},
			wantErr:    false,
		},
		{
			name:        "all invalid entries returns error",
			input:       []byte(`{invalid1}{invalid2}`),
			wantErr:     true,
			errContains: "failed to parse any videos",
		},
		{
			name: "video with generated thumbnail URL",
			input: []byte(`{"id": "novid123", "title": "No Thumbnail", "thumbnails": []}
`),
			wantCount:  1,
			wantIDs:    []string{"novid123"},
			wantTitles: []string{"No Thumbnail"},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			videos, err := ParseVideoList(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseVideoList() expected error, got nil")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("ParseVideoList() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseVideoList() unexpected error: %v", err)
				return
			}

			if len(videos) != tt.wantCount {
				t.Errorf("ParseVideoList() returned %d videos, want %d", len(videos), tt.wantCount)
				return
			}

			for i, wantID := range tt.wantIDs {
				if i >= len(videos) {
					break
				}
				if videos[i].ID != wantID {
					t.Errorf("ParseVideoList() video[%d].ID = %q, want %q", i, videos[i].ID, wantID)
				}
			}

			for i, wantTitle := range tt.wantTitles {
				if i >= len(videos) {
					break
				}
				if videos[i].Title != wantTitle {
					t.Errorf("ParseVideoList() video[%d].Title = %q, want %q", i, videos[i].Title, wantTitle)
				}
			}

			// Verify all videos have pending status
			for i, video := range videos {
				if video.Status != StatusPending {
					t.Errorf("ParseVideoList() video[%d].Status = %q, want %q", i, video.Status, StatusPending)
				}
			}
		})
	}
}

func TestParseVideoList_ThumbnailURL(t *testing.T) {
	tests := []struct {
		name         string
		input        []byte
		wantThumbURL string
	}{
		{
			name:         "thumbnail from thumbnails array",
			input:        []byte(`{"id": "vid1", "title": "Test", "thumbnails": [{"url": "https://custom.thumb.jpg"}]}`),
			wantThumbURL: "https://custom.thumb.jpg",
		},
		{
			name:         "generated thumbnail when empty array",
			input:        []byte(`{"id": "vid2", "title": "Test", "thumbnails": []}`),
			wantThumbURL: "https://i.ytimg.com/vi/vid2/maxresdefault.jpg",
		},
		{
			name:         "generated thumbnail when no thumbnails field",
			input:        []byte(`{"id": "vid3", "title": "Test"}`),
			wantThumbURL: "https://i.ytimg.com/vi/vid3/maxresdefault.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			videos, err := ParseVideoList(tt.input)
			if err != nil {
				t.Fatalf("ParseVideoList() unexpected error: %v", err)
			}

			if len(videos) != 1 {
				t.Fatalf("ParseVideoList() returned %d videos, want 1", len(videos))
			}

			if videos[0].ThumbnailURL != tt.wantThumbURL {
				t.Errorf("ParseVideoList() ThumbnailURL = %q, want %q", videos[0].ThumbnailURL, tt.wantThumbURL)
			}
		})
	}
}

func TestExtractChannelID(t *testing.T) {
	// YouTube channel IDs are "UC" followed by exactly 22 characters
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "URL with channel ID",
			input:    "https://www.youtube.com/channel/UCxyz123abc456def789gh12",
			expected: "UCxyz123abc456def789gh12",
		},
		{
			name:     "just channel ID",
			input:    "UCxyz123abc456def789gh12",
			expected: "UCxyz123abc456def789gh12",
		},
		{
			name:     "handle URL returns original",
			input:    "https://www.youtube.com/@testchannel",
			expected: "https://www.youtube.com/@testchannel",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "channel ID too short returns original",
			input:    "UCshort",
			expected: "UCshort",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractChannelID(tt.input)
			if result != tt.expected {
				t.Errorf("extractChannelID(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractVideoID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "11 character video ID",
			input:    "dQw4w9WgXcQ",
			expected: "dQw4w9WgXcQ",
		},
		{
			name:     "standard watch URL",
			input:    "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			expected: "dQw4w9WgXcQ",
		},
		{
			name:     "short URL",
			input:    "https://youtu.be/dQw4w9WgXcQ",
			expected: "dQw4w9WgXcQ",
		},
		{
			name:     "embed URL",
			input:    "https://www.youtube.com/embed/dQw4w9WgXcQ",
			expected: "dQw4w9WgXcQ",
		},
		{
			name:     "shorts URL",
			input:    "https://www.youtube.com/shorts/dQw4w9WgXcQ",
			expected: "dQw4w9WgXcQ",
		},
		{
			name:     "URL with extra parameters",
			input:    "https://www.youtube.com/watch?v=dQw4w9WgXcQ&list=PLrAXtmErZgOeiKm4sgNOknGvNjby9efdf",
			expected: "dQw4w9WgXcQ",
		},
		{
			name:     "invalid - too short",
			input:    "abc123",
			expected: "",
		},
		{
			name:     "invalid - too long",
			input:    "dQw4w9WgXcQextra",
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "whitespace around video ID",
			input:    "  dQw4w9WgXcQ  ",
			expected: "dQw4w9WgXcQ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractVideoID(tt.input)
			if result != tt.expected {
				t.Errorf("extractVideoID(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "seconds only",
			input:    "45",
			expected: 45,
		},
		{
			name:     "minutes and seconds",
			input:    "3:45",
			expected: 225,
		},
		{
			name:     "hours, minutes, seconds",
			input:    "1:30:45",
			expected: 5445,
		},
		{
			name:     "zero duration",
			input:    "0:00",
			expected: 0,
		},
		{
			name:     "single digit minutes",
			input:    "5:30",
			expected: 330,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDuration(tt.input)
			if result != tt.expected {
				t.Errorf("parseDuration(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseViewCount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name:     "numeric views",
			input:    "1234 views",
			expected: 1234,
		},
		{
			name:     "views with commas",
			input:    "1,234,567 views",
			expected: 1234567,
		},
		{
			name:     "K suffix",
			input:    "1.5K views",
			expected: 1500,
		},
		{
			name:     "M suffix",
			input:    "2.3M views",
			expected: 2300000,
		},
		{
			name:     "B suffix",
			input:    "1.2B views",
			expected: 1200000000,
		},
		{
			name:     "no views",
			input:    "No views",
			expected: 0,
		},
		{
			name:     "single view",
			input:    "1 view",
			expected: 1,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseViewCount(tt.input)
			if result != tt.expected {
				t.Errorf("parseViewCount(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	// Test that client can be created without error
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client == nil {
		t.Error("NewClient() returned nil client")
	}

	// Test with custom options
	client, err = NewClient(
		WithTimeout(10000000000), // 10 seconds
		WithUserAgent("TestAgent/1.0"),
		WithClientVersion("2.0"),
	)
	if err != nil {
		t.Fatalf("NewClient() with options error = %v", err)
	}
	if client == nil {
		t.Error("NewClient() with options returned nil client")
	}
}

func TestStreamInfoGetBestVideoFormat(t *testing.T) {
	info := &StreamInfo{
		Formats: []DownloadableFormat{
			{ITag: 18, Width: 640, Height: 360, Bitrate: 500000},
			{ITag: 22, Width: 1280, Height: 720, Bitrate: 2000000},
			{ITag: 137, Width: 1920, Height: 1080, Bitrate: 4000000},
			{ITag: 140, AudioQuality: "AUDIO_QUALITY_MEDIUM", Bitrate: 128000}, // Audio only
		},
	}

	best := info.GetBestVideoFormat()
	if best == nil {
		t.Fatal("GetBestVideoFormat() returned nil")
	}
	if best.ITag != 137 {
		t.Errorf("GetBestVideoFormat() ITag = %d, want 137", best.ITag)
	}
	if best.Height != 1080 {
		t.Errorf("GetBestVideoFormat() Height = %d, want 1080", best.Height)
	}
}

func TestStreamInfoGetBestAudioFormat(t *testing.T) {
	info := &StreamInfo{
		Formats: []DownloadableFormat{
			{ITag: 18, Width: 640, Height: 360, Bitrate: 500000},
			{ITag: 140, AudioQuality: "AUDIO_QUALITY_MEDIUM", Bitrate: 128000},
			{ITag: 251, AudioQuality: "AUDIO_QUALITY_HIGH", Bitrate: 160000},
		},
	}

	best := info.GetBestAudioFormat()
	if best == nil {
		t.Fatal("GetBestAudioFormat() returned nil")
	}
	if best.ITag != 251 {
		t.Errorf("GetBestAudioFormat() ITag = %d, want 251", best.ITag)
	}
}

func TestStreamInfoGetFormatByITag(t *testing.T) {
	info := &StreamInfo{
		Formats: []DownloadableFormat{
			{ITag: 18, Width: 640, Height: 360},
			{ITag: 22, Width: 1280, Height: 720},
		},
	}

	// Test finding existing format
	format := info.GetFormatByITag(22)
	if format == nil {
		t.Fatal("GetFormatByITag(22) returned nil")
	}
	if format.Height != 720 {
		t.Errorf("GetFormatByITag(22) Height = %d, want 720", format.Height)
	}

	// Test non-existent format
	format = info.GetFormatByITag(999)
	if format != nil {
		t.Error("GetFormatByITag(999) should return nil for non-existent format")
	}
}

func TestStreamInfoGetCombinedFormats(t *testing.T) {
	info := &StreamInfo{
		Formats: []DownloadableFormat{
			{ITag: 18, Width: 640, Height: 360, AudioChannels: 2},    // Combined
			{ITag: 137, Width: 1920, Height: 1080, AudioChannels: 0}, // Video only
			{ITag: 140, Width: 0, Height: 0, AudioChannels: 2},       // Audio only
			{ITag: 22, Width: 1280, Height: 720, AudioChannels: 2},   // Combined
		},
	}

	combined := info.GetCombinedFormats()
	if len(combined) != 2 {
		t.Errorf("GetCombinedFormats() returned %d formats, want 2", len(combined))
	}

	// Check that the combined formats have the expected ITags
	itags := make(map[int]bool)
	for _, f := range combined {
		itags[f.ITag] = true
	}
	if !itags[18] {
		t.Error("GetCombinedFormats() should include ITag 18")
	}
	if !itags[22] {
		t.Error("GetCombinedFormats() should include ITag 22")
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && searchString(s, substr)))
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
