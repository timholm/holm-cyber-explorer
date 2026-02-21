package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/timholm/ytarchive/tests/mocks"
)

// TestChannel represents the channel structure for API tests
type TestChannel struct {
	ID          string    `json:"id"`
	YouTubeURL  string    `json:"youtube_url"`
	YouTubeID   string    `json:"youtube_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	VideoCount  int       `json:"video_count"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastSyncAt  time.Time `json:"last_sync_at,omitempty"`
}

// TestProgress represents the progress structure for API tests
type TestProgress struct {
	TotalChannels      int     `json:"total_channels"`
	TotalVideos        int     `json:"total_videos"`
	DownloadedVideos   int     `json:"downloaded_videos"`
	PendingVideos      int     `json:"pending_videos"`
	FailedVideos       int     `json:"failed_videos"`
	ActiveJobs         int     `json:"active_jobs"`
	DownloadPercentage float64 `json:"download_percentage"`
}

// MockAPIHandlers provides mock API handlers for testing
type MockAPIHandlers struct {
	redis    *mocks.MockRedisClient
	channels map[string]TestChannel
	mu       sync.RWMutex
}

// NewMockAPIHandlers creates a new mock API handlers instance
func NewMockAPIHandlers() *MockAPIHandlers {
	return &MockAPIHandlers{
		redis:    mocks.NewMockRedisClient(),
		channels: make(map[string]TestChannel),
	}
}

// SetupMockRouter creates a test router with mock handlers
func SetupMockRouter(handlers *MockAPIHandlers) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Channel endpoints
		api.POST("/channels", handlers.AddChannel)
		api.GET("/channels", handlers.ListChannels)
		api.GET("/channels/:id", handlers.GetChannel)
		api.POST("/channels/:id/sync", handlers.TriggerSync)
		api.DELETE("/channels/:id", handlers.DeleteChannel)

		// Progress endpoint
		api.GET("/progress", handlers.GetProgress)

		// Jobs endpoint
		api.GET("/jobs", handlers.ListJobs)
	}

	return router
}

// AddChannel handles POST /api/channels
func (h *MockAPIHandlers) AddChannel(c *gin.Context) {
	var req struct {
		YouTubeURL string `json:"youtube_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Extract YouTube ID from URL (simplified for testing)
	youtubeID := extractTestChannelID(req.YouTubeURL)
	if youtubeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid YouTube URL"})
		return
	}

	// Check for duplicates
	h.mu.RLock()
	for _, ch := range h.channels {
		if ch.YouTubeID == youtubeID {
			h.mu.RUnlock()
			c.JSON(http.StatusConflict, gin.H{"error": "Channel already exists"})
			return
		}
	}
	h.mu.RUnlock()

	// Create channel
	now := time.Now()
	channel := TestChannel{
		ID:         generateTestID(),
		YouTubeURL: req.YouTubeURL,
		YouTubeID:  youtubeID,
		Name:       "Test Channel " + youtubeID,
		Status:     "pending",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Store in mock
	h.mu.Lock()
	h.channels[channel.ID] = channel
	h.mu.Unlock()
	channelJSON, _ := json.Marshal(channel)
	h.redis.Set(c.Request.Context(), "channel:"+channel.ID, string(channelJSON), 0)
	h.redis.SAdd(c.Request.Context(), "channels", channel.ID)

	c.JSON(http.StatusCreated, channel)
}

// ListChannels handles GET /api/channels
func (h *MockAPIHandlers) ListChannels(c *gin.Context) {
	h.mu.RLock()
	channels := make([]TestChannel, 0, len(h.channels))
	for _, ch := range h.channels {
		channels = append(channels, ch)
	}
	h.mu.RUnlock()
	c.JSON(http.StatusOK, gin.H{"channels": channels, "count": len(channels)})
}

// GetChannel handles GET /api/channels/:id
func (h *MockAPIHandlers) GetChannel(c *gin.Context) {
	channelID := c.Param("id")
	h.mu.RLock()
	channel, ok := h.channels[channelID]
	h.mu.RUnlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"channel": channel, "videos": []interface{}{}})
}

// TriggerSync handles POST /api/channels/:id/sync
func (h *MockAPIHandlers) TriggerSync(c *gin.Context) {
	channelID := c.Param("id")
	h.mu.Lock()
	channel, ok := h.channels[channelID]
	if !ok {
		h.mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	if channel.Status == "syncing" {
		h.mu.Unlock()
		c.JSON(http.StatusConflict, gin.H{"error": "Channel is already syncing"})
		return
	}

	// Update status
	channel.Status = "syncing"
	channel.UpdatedAt = time.Now()
	h.channels[channelID] = channel
	h.mu.Unlock()

	jobID := generateTestID()
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Sync started",
		"job_id":  jobID,
		"channel": channel,
	})
}

// DeleteChannel handles DELETE /api/channels/:id
func (h *MockAPIHandlers) DeleteChannel(c *gin.Context) {
	channelID := c.Param("id")
	h.mu.Lock()
	if _, ok := h.channels[channelID]; !ok {
		h.mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	delete(h.channels, channelID)
	h.mu.Unlock()
	h.redis.Del(c.Request.Context(), "channel:"+channelID)
	h.redis.SRem(c.Request.Context(), "channels", channelID)

	c.JSON(http.StatusOK, gin.H{"message": "Channel deleted successfully"})
}

// GetProgress handles GET /api/progress
func (h *MockAPIHandlers) GetProgress(c *gin.Context) {
	progress := TestProgress{
		TotalChannels:      len(h.channels),
		TotalVideos:        150,
		DownloadedVideos:   75,
		PendingVideos:      70,
		FailedVideos:       5,
		ActiveJobs:         2,
		DownloadPercentage: 50.0,
	}
	c.JSON(http.StatusOK, progress)
}

// ListJobs handles GET /api/jobs
func (h *MockAPIHandlers) ListJobs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"jobs": []interface{}{}, "count": 0})
}

// Helper functions
func extractTestChannelID(url string) string {
	// Simplified extraction for testing
	if len(url) > 20 {
		// Extract @username or channel ID
		for i := len(url) - 1; i >= 0; i-- {
			if url[i] == '/' || url[i] == '@' {
				return url[i+1:]
			}
		}
	}
	return ""
}

var testIDCounter int64

func generateTestID() string {
	id := atomic.AddInt64(&testIDCounter, 1)
	return "test-id-" + string(rune('0'+id%10)) + string(rune('0'+id/10))
}

// TestAddChannel tests the POST /api/channels endpoint
func TestAddChannel(t *testing.T) {
	handlers := NewMockAPIHandlers()
	router := SetupMockRouter(handlers)

	tests := []struct {
		name           string
		payload        map[string]string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "ValidChannel",
			payload:        map[string]string{"youtube_url": "https://www.youtube.com/@aperturethinking"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "MissingURL",
			payload:        map[string]string{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name:           "InvalidURL",
			payload:        map[string]string{"youtube_url": "not-a-url"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid YouTube URL",
		},
		{
			name:           "ValidChannelWithChannelID",
			payload:        map[string]string{"youtube_url": "https://www.youtube.com/channel/UC_test_channel"},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/api/channels", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				var resp map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &resp)
				if errMsg, ok := resp["error"].(string); ok {
					if errMsg != tt.expectedError && !containsString(errMsg, tt.expectedError) {
						t.Errorf("Expected error containing '%s', got '%s'", tt.expectedError, errMsg)
					}
				}
			}
		})
	}

	// Test duplicate channel
	t.Run("DuplicateChannel", func(t *testing.T) {
		// Add channel first
		payload := map[string]string{"youtube_url": "https://www.youtube.com/@duplicatechannel"}
		payloadBytes, _ := json.Marshal(payload)

		req1, _ := http.NewRequest("POST", "/api/channels", bytes.NewBuffer(payloadBytes))
		req1.Header.Set("Content-Type", "application/json")
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)

		if w1.Code != http.StatusCreated {
			t.Fatalf("First request should succeed, got %d", w1.Code)
		}

		// Try to add same channel again
		req2, _ := http.NewRequest("POST", "/api/channels", bytes.NewBuffer(payloadBytes))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		if w2.Code != http.StatusConflict {
			t.Errorf("Expected status 409 for duplicate, got %d", w2.Code)
		}
	})
}

// TestListChannels tests the GET /api/channels endpoint
func TestListChannels(t *testing.T) {
	handlers := NewMockAPIHandlers()
	router := SetupMockRouter(handlers)

	// Add some channels first
	channelURLs := []string{
		"https://www.youtube.com/@channel1",
		"https://www.youtube.com/@channel2",
		"https://www.youtube.com/@channel3",
	}

	for _, url := range channelURLs {
		payload := map[string]string{"youtube_url": url}
		payloadBytes, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/channels", bytes.NewBuffer(payloadBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// Test list channels
	t.Run("ListAllChannels", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/channels", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp struct {
			Channels []TestChannel `json:"channels"`
			Count    int           `json:"count"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)

		if resp.Count != 3 {
			t.Errorf("Expected count 3, got %d", resp.Count)
		}

		if len(resp.Channels) != 3 {
			t.Errorf("Expected 3 channels, got %d", len(resp.Channels))
		}
	})

	// Test empty list
	t.Run("EmptyList", func(t *testing.T) {
		emptyHandlers := NewMockAPIHandlers()
		emptyRouter := SetupMockRouter(emptyHandlers)

		req, _ := http.NewRequest("GET", "/api/channels", nil)
		w := httptest.NewRecorder()
		emptyRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp struct {
			Channels []TestChannel `json:"channels"`
			Count    int           `json:"count"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)

		if resp.Count != 0 {
			t.Errorf("Expected count 0, got %d", resp.Count)
		}
	})
}

// TestTriggerSync tests the POST /api/channels/:id/sync endpoint
func TestTriggerSync(t *testing.T) {
	handlers := NewMockAPIHandlers()
	router := SetupMockRouter(handlers)

	// Add a channel first
	payload := map[string]string{"youtube_url": "https://www.youtube.com/@syncchannel"}
	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/channels", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var createResp TestChannel
	json.Unmarshal(w.Body.Bytes(), &createResp)
	channelID := createResp.ID

	// Test sync
	t.Run("TriggerSyncSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/channels/"+channelID+"/sync", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusAccepted {
			t.Errorf("Expected status 202, got %d", w.Code)
		}

		var resp struct {
			Message string      `json:"message"`
			JobID   string      `json:"job_id"`
			Channel TestChannel `json:"channel"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)

		if resp.Message != "Sync started" {
			t.Errorf("Expected message 'Sync started', got '%s'", resp.Message)
		}

		if resp.JobID == "" {
			t.Error("Expected non-empty job ID")
		}

		if resp.Channel.Status != "syncing" {
			t.Errorf("Expected status 'syncing', got '%s'", resp.Channel.Status)
		}
	})

	// Test sync while already syncing
	t.Run("SyncWhileAlreadySyncing", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/channels/"+channelID+"/sync", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("Expected status 409, got %d", w.Code)
		}
	})

	// Test sync non-existent channel
	t.Run("SyncNonExistentChannel", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/channels/non-existent-id/sync", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}

// TestGetProgress tests the GET /api/progress endpoint
func TestGetProgress(t *testing.T) {
	handlers := NewMockAPIHandlers()
	router := SetupMockRouter(handlers)

	// Add some channels
	for i := 0; i < 3; i++ {
		payload := map[string]string{"youtube_url": "https://www.youtube.com/@progresschannel" + string(rune('0'+i))}
		payloadBytes, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/channels", bytes.NewBuffer(payloadBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	t.Run("GetProgressSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/progress", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var progress TestProgress
		json.Unmarshal(w.Body.Bytes(), &progress)

		if progress.TotalChannels != 3 {
			t.Errorf("Expected 3 channels, got %d", progress.TotalChannels)
		}

		if progress.TotalVideos <= 0 {
			t.Error("Expected positive total videos")
		}

		if progress.DownloadPercentage < 0 || progress.DownloadPercentage > 100 {
			t.Errorf("Download percentage should be between 0 and 100, got %f", progress.DownloadPercentage)
		}

		// Verify video counts add up
		expectedTotal := progress.DownloadedVideos + progress.PendingVideos + progress.FailedVideos
		if expectedTotal != progress.TotalVideos {
			t.Errorf("Video counts don't add up: %d + %d + %d != %d",
				progress.DownloadedVideos, progress.PendingVideos, progress.FailedVideos, progress.TotalVideos)
		}
	})
}

// TestHealthEndpoint tests the GET /health endpoint
func TestHealthEndpoint(t *testing.T) {
	handlers := NewMockAPIHandlers()
	router := SetupMockRouter(handlers)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", resp["status"])
	}

	if resp["timestamp"] == nil {
		t.Error("Expected timestamp in response")
	}
}

// TestDeleteChannel tests the DELETE /api/channels/:id endpoint
func TestDeleteChannel(t *testing.T) {
	handlers := NewMockAPIHandlers()
	router := SetupMockRouter(handlers)

	// Add a channel first
	payload := map[string]string{"youtube_url": "https://www.youtube.com/@deletechannel"}
	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/channels", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var createResp TestChannel
	json.Unmarshal(w.Body.Bytes(), &createResp)
	channelID := createResp.ID

	// Test delete
	t.Run("DeleteSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/channels/"+channelID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Verify channel is deleted
		getReq, _ := http.NewRequest("GET", "/api/channels/"+channelID, nil)
		getW := httptest.NewRecorder()
		router.ServeHTTP(getW, getReq)

		if getW.Code != http.StatusNotFound {
			t.Errorf("Expected status 404 for deleted channel, got %d", getW.Code)
		}
	})

	// Test delete non-existent channel
	t.Run("DeleteNonExistent", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/channels/non-existent", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}

// TestGetChannel tests the GET /api/channels/:id endpoint
func TestGetChannel(t *testing.T) {
	handlers := NewMockAPIHandlers()
	router := SetupMockRouter(handlers)

	// Add a channel first
	payload := map[string]string{"youtube_url": "https://www.youtube.com/@getchanneltest"}
	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/channels", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var createResp TestChannel
	json.Unmarshal(w.Body.Bytes(), &createResp)
	channelID := createResp.ID

	t.Run("GetExistingChannel", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/channels/"+channelID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp struct {
			Channel TestChannel   `json:"channel"`
			Videos  []interface{} `json:"videos"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)

		if resp.Channel.ID != channelID {
			t.Errorf("Expected channel ID '%s', got '%s'", channelID, resp.Channel.ID)
		}
	})

	t.Run("GetNonExistentChannel", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/channels/non-existent", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}

// TestConcurrentRequests tests handling of concurrent API requests
func TestConcurrentRequests(t *testing.T) {
	handlers := NewMockAPIHandlers()
	router := SetupMockRouter(handlers)

	// Run concurrent channel additions
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			payload := map[string]string{"youtube_url": "https://www.youtube.com/@concurrent" + string(rune('a'+idx))}
			payloadBytes, _ := json.Marshal(payload)
			req, _ := http.NewRequest("POST", "/api/channels", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			done <- w.Code == http.StatusCreated
		}(i)
	}

	// Wait for all requests to complete
	successCount := 0
	for i := 0; i < 10; i++ {
		if <-done {
			successCount++
		}
	}

	if successCount != 10 {
		t.Errorf("Expected all 10 concurrent requests to succeed, got %d", successCount)
	}

	// Verify all channels were created
	req, _ := http.NewRequest("GET", "/api/channels", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp struct {
		Count int `json:"count"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Count != 10 {
		t.Errorf("Expected 10 channels, got %d", resp.Count)
	}
}

// TestAPIWithMockRedis tests API handlers with mock Redis backend
func TestAPIWithMockRedis(t *testing.T) {
	ctx := context.Background()
	redis := mocks.NewMockRedisClient()

	// Test Redis operations through API handlers
	t.Run("RedisChannelStorage", func(t *testing.T) {
		channel := TestChannel{
			ID:        "test-channel-1",
			YouTubeID: "UC_test_1",
			Name:      "Test Channel",
			Status:    "pending",
		}

		channelJSON, _ := json.Marshal(channel)

		// Store channel
		err := redis.Set(ctx, "channel:"+channel.ID, string(channelJSON), 0)
		if err != nil {
			t.Fatalf("Failed to store channel: %v", err)
		}

		// Add to channel list
		_, err = redis.SAdd(ctx, "channels", channel.ID)
		if err != nil {
			t.Fatalf("Failed to add to channel list: %v", err)
		}

		// Retrieve channel
		stored, err := redis.Get(ctx, "channel:"+channel.ID)
		if err != nil {
			t.Fatalf("Failed to get channel: %v", err)
		}

		var retrieved TestChannel
		json.Unmarshal([]byte(stored), &retrieved)

		if retrieved.Name != channel.Name {
			t.Errorf("Expected name '%s', got '%s'", channel.Name, retrieved.Name)
		}

		// List channels
		channelIDs, err := redis.SMembers(ctx, "channels")
		if err != nil {
			t.Fatalf("Failed to get channel list: %v", err)
		}

		if len(channelIDs) != 1 {
			t.Errorf("Expected 1 channel in list, got %d", len(channelIDs))
		}
	})

	// Test Redis failure handling
	t.Run("RedisFailure", func(t *testing.T) {
		redis.SetFailure(true, "connection refused")

		_, err := redis.Get(ctx, "any-key")
		if err == nil {
			t.Error("Expected error when Redis fails")
		}

		redis.SetFailure(false, "")
	})
}

// Helper function to check if string contains substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
