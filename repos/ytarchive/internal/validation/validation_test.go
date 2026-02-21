package validation

import (
	"testing"
)

func TestValidateChannelInput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid channel ID",
			input:   "UCddiUEpeqJcYeBxX1IVBKvQ",
			want:    "UCddiUEpeqJcYeBxX1IVBKvQ",
			wantErr: false,
		},
		{
			name:    "valid handle with @",
			input:   "@channelname",
			want:    "@channelname",
			wantErr: false,
		},
		{
			name:    "valid handle without @",
			input:   "channelname",
			want:    "@channelname",
			wantErr: false,
		},
		{
			name:    "valid YouTube URL with channel ID",
			input:   "https://www.youtube.com/channel/UCddiUEpeqJcYeBxX1IVBKvQ",
			want:    "UCddiUEpeqJcYeBxX1IVBKvQ",
			wantErr: false,
		},
		{
			name:    "valid YouTube URL with handle",
			input:   "https://www.youtube.com/@channelname",
			want:    "@channelname",
			wantErr: false,
		},
		{
			name:    "unsafe input with shell injection",
			input:   "channel; rm -rf /",
			wantErr: true,
		},
		{
			name:    "unsafe input with XSS",
			input:   "<script>alert('xss')</script>",
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid URL host",
			input:   "https://example.com/channel/test",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateChannelInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateChannelInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ValidateChannelInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateVideoID(t *testing.T) {
	tests := []struct {
		name    string
		videoID string
		wantErr bool
	}{
		{
			name:    "valid video ID",
			videoID: "dQw4w9WgXcQ",
			wantErr: false,
		},
		{
			name:    "valid video ID with underscore",
			videoID: "a1B2c3D4_-E",
			wantErr: false,
		},
		{
			name:    "too short",
			videoID: "abc123",
			wantErr: true,
		},
		{
			name:    "too long",
			videoID: "dQw4w9WgXcQabc",
			wantErr: true,
		},
		{
			name:    "empty",
			videoID: "",
			wantErr: true,
		},
		{
			name:    "unsafe characters",
			videoID: "abc<>123456",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVideoID(tt.videoID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVideoID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateVideoURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name:    "standard watch URL",
			url:     "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			want:    "dQw4w9WgXcQ",
			wantErr: false,
		},
		{
			name:    "short URL",
			url:     "https://youtu.be/dQw4w9WgXcQ",
			want:    "dQw4w9WgXcQ",
			wantErr: false,
		},
		{
			name:    "embed URL",
			url:     "https://www.youtube.com/embed/dQw4w9WgXcQ",
			want:    "dQw4w9WgXcQ",
			wantErr: false,
		},
		{
			name:    "URL without scheme",
			url:     "www.youtube.com/watch?v=dQw4w9WgXcQ",
			want:    "dQw4w9WgXcQ",
			wantErr: false,
		},
		{
			name:    "invalid host",
			url:     "https://example.com/watch?v=dQw4w9WgXcQ",
			wantErr: true,
		},
		{
			name:    "no video ID",
			url:     "https://www.youtube.com/",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateVideoURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVideoURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ValidateVideoURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "normal filename",
			filename: "my video.mp4",
			want:     "my video.mp4",
		},
		{
			name:     "filename with slashes",
			filename: "my/video\\file.mp4",
			want:     "my_video_file.mp4",
		},
		{
			name:     "filename with special chars",
			filename: "video:test*file?.mp4",
			want:     "video_test_file_.mp4",
		},
		{
			name:     "filename with leading dots",
			filename: "...hidden.mp4",
			want:     "hidden.mp4",
		},
		{
			name:     "empty filename",
			filename: "",
			want:     "unnamed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeFilename(tt.filename)
			if got != tt.want {
				t.Errorf("SanitizeFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}
