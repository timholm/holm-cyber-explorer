package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

// YouTubeChannelResponse represents a mock YouTube channel API response
type YouTubeChannelResponse struct {
	Kind     string `json:"kind"`
	Etag     string `json:"etag"`
	PageInfo struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []YouTubeChannelItem `json:"items"`
}

// YouTubeChannelItem represents a single channel in the response
type YouTubeChannelItem struct {
	Kind    string `json:"kind"`
	Etag    string `json:"etag"`
	ID      string `json:"id"`
	Snippet struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		CustomURL   string `json:"customUrl"`
		PublishedAt string `json:"publishedAt"`
		Thumbnails  struct {
			Default struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"default"`
			Medium struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"medium"`
			High struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"high"`
		} `json:"thumbnails"`
		Country string `json:"country"`
	} `json:"snippet"`
	Statistics struct {
		ViewCount             string `json:"viewCount"`
		SubscriberCount       string `json:"subscriberCount"`
		HiddenSubscriberCount bool   `json:"hiddenSubscriberCount"`
		VideoCount            string `json:"videoCount"`
	} `json:"statistics"`
}

// YouTubeVideoListResponse represents a mock YouTube video list API response
type YouTubeVideoListResponse struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken,omitempty"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []YouTubeVideoItem `json:"items"`
}

// YouTubeVideoItem represents a single video in the response
type YouTubeVideoItem struct {
	Kind    string `json:"kind"`
	Etag    string `json:"etag"`
	ID      string `json:"id"`
	Snippet struct {
		PublishedAt string `json:"publishedAt"`
		ChannelID   string `json:"channelId"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Thumbnails  struct {
			Default struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"default"`
			Medium struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"medium"`
			High struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"high"`
		} `json:"thumbnails"`
		ChannelTitle         string   `json:"channelTitle"`
		Tags                 []string `json:"tags,omitempty"`
		CategoryID           string   `json:"categoryId"`
		LiveBroadcastContent string   `json:"liveBroadcastContent"`
	} `json:"snippet"`
	ContentDetails struct {
		Duration        string `json:"duration"`
		Dimension       string `json:"dimension"`
		Definition      string `json:"definition"`
		Caption         string `json:"caption"`
		LicensedContent bool   `json:"licensedContent"`
		Projection      string `json:"projection"`
	} `json:"contentDetails"`
	Statistics struct {
		ViewCount     string `json:"viewCount"`
		LikeCount     string `json:"likeCount"`
		FavoriteCount string `json:"favoriteCount"`
		CommentCount  string `json:"commentCount"`
	} `json:"statistics"`
}

// MockYouTubeClient provides mock YouTube API responses for testing
type MockYouTubeClient struct {
	Server       *httptest.Server
	Channels     map[string]YouTubeChannelItem
	Videos       map[string][]YouTubeVideoItem
	ShouldFail   bool
	FailMessage  string
	RequestCount int
}

// NewMockYouTubeClient creates a new mock YouTube client
func NewMockYouTubeClient() *MockYouTubeClient {
	mock := &MockYouTubeClient{
		Channels: make(map[string]YouTubeChannelItem),
		Videos:   make(map[string][]YouTubeVideoItem),
	}

	// Add default test data
	mock.AddDefaultTestData()

	// Create test server
	mock.Server = httptest.NewServer(http.HandlerFunc(mock.handleRequest))

	return mock
}

// AddDefaultTestData adds sample data for testing
func (m *MockYouTubeClient) AddDefaultTestData() {
	// Add aperturethinking channel
	m.Channels["aperturethinking"] = YouTubeChannelItem{
		Kind: "youtube#channel",
		Etag: "test-etag-123",
		ID:   "UC_aperturethinking_ID",
		Snippet: struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			CustomURL   string `json:"customUrl"`
			PublishedAt string `json:"publishedAt"`
			Thumbnails  struct {
				Default struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"default"`
				Medium struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"medium"`
				High struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"high"`
			} `json:"thumbnails"`
			Country string `json:"country"`
		}{
			Title:       "Aperture Thinking",
			Description: "A channel about photography, filmmaking, and creative visual thinking.",
			CustomURL:   "@aperturethinking",
			PublishedAt: "2015-03-15T00:00:00Z",
			Country:     "US",
		},
		Statistics: struct {
			ViewCount             string `json:"viewCount"`
			SubscriberCount       string `json:"subscriberCount"`
			HiddenSubscriberCount bool   `json:"hiddenSubscriberCount"`
			VideoCount            string `json:"videoCount"`
		}{
			ViewCount:             "1500000",
			SubscriberCount:       "50000",
			HiddenSubscriberCount: false,
			VideoCount:            "150",
		},
	}

	// Add sample videos for the channel
	m.Videos["UC_aperturethinking_ID"] = []YouTubeVideoItem{
		{
			Kind: "youtube#video",
			Etag: "video-etag-1",
			ID:   "video_001",
			Snippet: struct {
				PublishedAt string `json:"publishedAt"`
				ChannelID   string `json:"channelId"`
				Title       string `json:"title"`
				Description string `json:"description"`
				Thumbnails  struct {
					Default struct {
						URL    string `json:"url"`
						Width  int    `json:"width"`
						Height int    `json:"height"`
					} `json:"default"`
					Medium struct {
						URL    string `json:"url"`
						Width  int    `json:"width"`
						Height int    `json:"height"`
					} `json:"medium"`
					High struct {
						URL    string `json:"url"`
						Width  int    `json:"width"`
						Height int    `json:"height"`
					} `json:"high"`
				} `json:"thumbnails"`
				ChannelTitle         string   `json:"channelTitle"`
				Tags                 []string `json:"tags,omitempty"`
				CategoryID           string   `json:"categoryId"`
				LiveBroadcastContent string   `json:"liveBroadcastContent"`
			}{
				PublishedAt:          "2024-01-15T10:00:00Z",
				ChannelID:            "UC_aperturethinking_ID",
				Title:                "Understanding Aperture in Photography",
				Description:          "Learn the basics of aperture and how it affects your photos.",
				ChannelTitle:         "Aperture Thinking",
				Tags:                 []string{"photography", "aperture", "tutorial"},
				CategoryID:           "22",
				LiveBroadcastContent: "none",
			},
			ContentDetails: struct {
				Duration        string `json:"duration"`
				Dimension       string `json:"dimension"`
				Definition      string `json:"definition"`
				Caption         string `json:"caption"`
				LicensedContent bool   `json:"licensedContent"`
				Projection      string `json:"projection"`
			}{
				Duration:        "PT15M30S",
				Dimension:       "2d",
				Definition:      "hd",
				Caption:         "true",
				LicensedContent: true,
				Projection:      "rectangular",
			},
			Statistics: struct {
				ViewCount     string `json:"viewCount"`
				LikeCount     string `json:"likeCount"`
				FavoriteCount string `json:"favoriteCount"`
				CommentCount  string `json:"commentCount"`
			}{
				ViewCount:     "25000",
				LikeCount:     "1500",
				FavoriteCount: "0",
				CommentCount:  "120",
			},
		},
		{
			Kind: "youtube#video",
			Etag: "video-etag-2",
			ID:   "video_002",
			Snippet: struct {
				PublishedAt string `json:"publishedAt"`
				ChannelID   string `json:"channelId"`
				Title       string `json:"title"`
				Description string `json:"description"`
				Thumbnails  struct {
					Default struct {
						URL    string `json:"url"`
						Width  int    `json:"width"`
						Height int    `json:"height"`
					} `json:"default"`
					Medium struct {
						URL    string `json:"url"`
						Width  int    `json:"width"`
						Height int    `json:"height"`
					} `json:"medium"`
					High struct {
						URL    string `json:"url"`
						Width  int    `json:"width"`
						Height int    `json:"height"`
					} `json:"high"`
				} `json:"thumbnails"`
				ChannelTitle         string   `json:"channelTitle"`
				Tags                 []string `json:"tags,omitempty"`
				CategoryID           string   `json:"categoryId"`
				LiveBroadcastContent string   `json:"liveBroadcastContent"`
			}{
				PublishedAt:          "2024-01-20T14:00:00Z",
				ChannelID:            "UC_aperturethinking_ID",
				Title:                "Low Light Photography Tips",
				Description:          "Master low light photography with these essential tips and techniques.",
				ChannelTitle:         "Aperture Thinking",
				Tags:                 []string{"photography", "low light", "tips"},
				CategoryID:           "22",
				LiveBroadcastContent: "none",
			},
			ContentDetails: struct {
				Duration        string `json:"duration"`
				Dimension       string `json:"dimension"`
				Definition      string `json:"definition"`
				Caption         string `json:"caption"`
				LicensedContent bool   `json:"licensedContent"`
				Projection      string `json:"projection"`
			}{
				Duration:        "PT22M45S",
				Dimension:       "2d",
				Definition:      "hd",
				Caption:         "true",
				LicensedContent: true,
				Projection:      "rectangular",
			},
			Statistics: struct {
				ViewCount     string `json:"viewCount"`
				LikeCount     string `json:"likeCount"`
				FavoriteCount string `json:"favoriteCount"`
				CommentCount  string `json:"commentCount"`
			}{
				ViewCount:     "18000",
				LikeCount:     "1200",
				FavoriteCount: "0",
				CommentCount:  "95",
			},
		},
	}
}

// handleRequest processes incoming mock API requests
func (m *MockYouTubeClient) handleRequest(w http.ResponseWriter, r *http.Request) {
	m.RequestCount++

	if m.ShouldFail {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{
				"code":    500,
				"message": m.FailMessage,
			},
		})
		return
	}

	path := r.URL.Path
	w.Header().Set("Content-Type", "application/json")

	switch {
	case path == "/youtube/v3/channels":
		m.handleChannelsRequest(w, r)
	case path == "/youtube/v3/search":
		m.handleSearchRequest(w, r)
	case path == "/youtube/v3/videos":
		m.handleVideosRequest(w, r)
	case path == "/youtube/v3/playlistItems":
		m.handlePlaylistItemsRequest(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Unknown endpoint: " + path,
		})
	}
}

// handleChannelsRequest handles /youtube/v3/channels requests
func (m *MockYouTubeClient) handleChannelsRequest(w http.ResponseWriter, r *http.Request) {
	forUsername := r.URL.Query().Get("forUsername")
	forHandle := r.URL.Query().Get("forHandle")
	id := r.URL.Query().Get("id")

	var channelKey string
	if forUsername != "" {
		channelKey = forUsername
	} else if forHandle != "" {
		channelKey = forHandle
	} else if id != "" {
		channelKey = id
	}

	response := YouTubeChannelResponse{
		Kind: "youtube#channelListResponse",
		Etag: "response-etag",
	}

	if channel, ok := m.Channels[channelKey]; ok {
		response.Items = []YouTubeChannelItem{channel}
		response.PageInfo.TotalResults = 1
		response.PageInfo.ResultsPerPage = 1
	} else {
		// Check if any channel ID matches
		for _, ch := range m.Channels {
			if ch.ID == channelKey {
				response.Items = []YouTubeChannelItem{ch}
				response.PageInfo.TotalResults = 1
				response.PageInfo.ResultsPerPage = 1
				break
			}
		}
	}

	json.NewEncoder(w).Encode(response)
}

// handleSearchRequest handles /youtube/v3/search requests
func (m *MockYouTubeClient) handleSearchRequest(w http.ResponseWriter, r *http.Request) {
	channelID := r.URL.Query().Get("channelId")
	searchType := r.URL.Query().Get("type")

	if searchType != "video" {
		json.NewEncoder(w).Encode(YouTubeVideoListResponse{
			Kind: "youtube#searchListResponse",
		})
		return
	}

	response := struct {
		Kind          string `json:"kind"`
		Etag          string `json:"etag"`
		NextPageToken string `json:"nextPageToken,omitempty"`
		PageInfo      struct {
			TotalResults   int `json:"totalResults"`
			ResultsPerPage int `json:"resultsPerPage"`
		} `json:"pageInfo"`
		Items []struct {
			Kind string `json:"kind"`
			Etag string `json:"etag"`
			ID   struct {
				Kind    string `json:"kind"`
				VideoID string `json:"videoId"`
			} `json:"id"`
		} `json:"items"`
	}{
		Kind: "youtube#searchListResponse",
		Etag: "search-etag",
	}

	if videos, ok := m.Videos[channelID]; ok {
		for _, v := range videos {
			response.Items = append(response.Items, struct {
				Kind string `json:"kind"`
				Etag string `json:"etag"`
				ID   struct {
					Kind    string `json:"kind"`
					VideoID string `json:"videoId"`
				} `json:"id"`
			}{
				Kind: "youtube#searchResult",
				Etag: v.Etag,
				ID: struct {
					Kind    string `json:"kind"`
					VideoID string `json:"videoId"`
				}{
					Kind:    "youtube#video",
					VideoID: v.ID,
				},
			})
		}
		response.PageInfo.TotalResults = len(videos)
		response.PageInfo.ResultsPerPage = len(videos)
	}

	json.NewEncoder(w).Encode(response)
}

// handleVideosRequest handles /youtube/v3/videos requests
func (m *MockYouTubeClient) handleVideosRequest(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	response := YouTubeVideoListResponse{
		Kind: "youtube#videoListResponse",
		Etag: "videos-etag",
	}

	// Find video by ID across all channels
	for _, videos := range m.Videos {
		for _, v := range videos {
			if v.ID == id {
				response.Items = []YouTubeVideoItem{v}
				response.PageInfo.TotalResults = 1
				response.PageInfo.ResultsPerPage = 1
				break
			}
		}
	}

	json.NewEncoder(w).Encode(response)
}

// handlePlaylistItemsRequest handles /youtube/v3/playlistItems requests
func (m *MockYouTubeClient) handlePlaylistItemsRequest(w http.ResponseWriter, r *http.Request) {
	playlistID := r.URL.Query().Get("playlistId")

	// Uploads playlist ID typically starts with UU instead of UC for channel
	channelID := "UC" + playlistID[2:]

	response := struct {
		Kind          string `json:"kind"`
		Etag          string `json:"etag"`
		NextPageToken string `json:"nextPageToken,omitempty"`
		PageInfo      struct {
			TotalResults   int `json:"totalResults"`
			ResultsPerPage int `json:"resultsPerPage"`
		} `json:"pageInfo"`
		Items []struct {
			Kind    string `json:"kind"`
			Etag    string `json:"etag"`
			ID      string `json:"id"`
			Snippet struct {
				PublishedAt string `json:"publishedAt"`
				ChannelID   string `json:"channelId"`
				Title       string `json:"title"`
				Description string `json:"description"`
				ResourceID  struct {
					Kind    string `json:"kind"`
					VideoID string `json:"videoId"`
				} `json:"resourceId"`
			} `json:"snippet"`
		} `json:"items"`
	}{
		Kind: "youtube#playlistItemListResponse",
		Etag: "playlist-etag",
	}

	if videos, ok := m.Videos[channelID]; ok {
		for _, v := range videos {
			response.Items = append(response.Items, struct {
				Kind    string `json:"kind"`
				Etag    string `json:"etag"`
				ID      string `json:"id"`
				Snippet struct {
					PublishedAt string `json:"publishedAt"`
					ChannelID   string `json:"channelId"`
					Title       string `json:"title"`
					Description string `json:"description"`
					ResourceID  struct {
						Kind    string `json:"kind"`
						VideoID string `json:"videoId"`
					} `json:"resourceId"`
				} `json:"snippet"`
			}{
				Kind: "youtube#playlistItem",
				Etag: v.Etag,
				ID:   fmt.Sprintf("PLI_%s", v.ID),
				Snippet: struct {
					PublishedAt string `json:"publishedAt"`
					ChannelID   string `json:"channelId"`
					Title       string `json:"title"`
					Description string `json:"description"`
					ResourceID  struct {
						Kind    string `json:"kind"`
						VideoID string `json:"videoId"`
					} `json:"resourceId"`
				}{
					PublishedAt: v.Snippet.PublishedAt,
					ChannelID:   v.Snippet.ChannelID,
					Title:       v.Snippet.Title,
					Description: v.Snippet.Description,
					ResourceID: struct {
						Kind    string `json:"kind"`
						VideoID string `json:"videoId"`
					}{
						Kind:    "youtube#video",
						VideoID: v.ID,
					},
				},
			})
		}
		response.PageInfo.TotalResults = len(videos)
		response.PageInfo.ResultsPerPage = len(videos)
	}

	json.NewEncoder(w).Encode(response)
}

// Close shuts down the mock server
func (m *MockYouTubeClient) Close() {
	if m.Server != nil {
		m.Server.Close()
	}
}

// URL returns the mock server URL
func (m *MockYouTubeClient) URL() string {
	return m.Server.URL
}

// SetFailure configures the mock to return errors
func (m *MockYouTubeClient) SetFailure(shouldFail bool, message string) {
	m.ShouldFail = shouldFail
	m.FailMessage = message
}

// ResetRequestCount resets the request counter
func (m *MockYouTubeClient) ResetRequestCount() {
	m.RequestCount = 0
}

// AddChannel adds a channel to the mock data
func (m *MockYouTubeClient) AddChannel(handle string, channel YouTubeChannelItem) {
	m.Channels[handle] = channel
}

// AddVideo adds a video to a channel's video list
func (m *MockYouTubeClient) AddVideo(channelID string, video YouTubeVideoItem) {
	m.Videos[channelID] = append(m.Videos[channelID], video)
}

// MockYouTubeDownloader simulates yt-dlp download behavior
type MockYouTubeDownloader struct {
	DownloadedVideos map[string]DownloadResult
	ShouldFail       bool
	FailMessage      string
	DownloadDelay    time.Duration
}

// DownloadResult represents the result of a download operation
type DownloadResult struct {
	VideoID   string
	FilePath  string
	FileSize  int64
	Duration  int
	Format    string
	Error     error
	Timestamp time.Time
}

// NewMockYouTubeDownloader creates a new mock downloader
func NewMockYouTubeDownloader() *MockYouTubeDownloader {
	return &MockYouTubeDownloader{
		DownloadedVideos: make(map[string]DownloadResult),
		DownloadDelay:    100 * time.Millisecond,
	}
}

// Download simulates downloading a video
func (m *MockYouTubeDownloader) Download(videoID, outputPath string) (DownloadResult, error) {
	time.Sleep(m.DownloadDelay)

	if m.ShouldFail {
		result := DownloadResult{
			VideoID:   videoID,
			Error:     fmt.Errorf(m.FailMessage),
			Timestamp: time.Now(),
		}
		return result, result.Error
	}

	result := DownloadResult{
		VideoID:   videoID,
		FilePath:  fmt.Sprintf("%s/%s.mp4", outputPath, videoID),
		FileSize:  1024 * 1024 * 100, // 100MB mock file
		Duration:  930,               // 15:30 mock duration
		Format:    "mp4",
		Timestamp: time.Now(),
	}

	m.DownloadedVideos[videoID] = result
	return result, nil
}

// GetDownloadedVideos returns all downloaded videos
func (m *MockYouTubeDownloader) GetDownloadedVideos() map[string]DownloadResult {
	return m.DownloadedVideos
}

// SetFailure configures the mock to return errors
func (m *MockYouTubeDownloader) SetFailure(shouldFail bool, message string) {
	m.ShouldFail = shouldFail
	m.FailMessage = message
}

// Reset clears the downloaded videos map
func (m *MockYouTubeDownloader) Reset() {
	m.DownloadedVideos = make(map[string]DownloadResult)
}
