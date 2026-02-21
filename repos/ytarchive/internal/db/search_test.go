package db

import (
	"testing"
)

func TestSearchVideos(t *testing.T) {
	db := setupTestDB(t)

	// Insert some test videos
	testVideos := []*Video{
		{ID: "vid1", Title: "Golang Tutorial for Beginners", Description: "Learn Go programming from scratch"},
		{ID: "vid2", Title: "Python Machine Learning", Description: "Build ML models with Python"},
		{ID: "vid3", Title: "Advanced Go Concurrency", Description: "Master goroutines and channels in Golang"},
		{ID: "vid4", Title: "JavaScript Basics", Description: "Web development fundamentals"},
	}

	for _, v := range testVideos {
		if err := InsertVideo(db, v); err != nil {
			t.Fatalf("Failed to insert video: %v", err)
		}
	}

	tests := []struct {
		name      string
		query     string
		wantCount int
	}{
		{"search for golang", "golang", 2},
		{"search for python", "python", 1},
		{"search for programming", "programming", 1},
		{"search for nonexistent", "nonexistent", 0},
		{"search for basics", "basics", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := SearchVideos(db, tt.query, 50)
			if err != nil {
				t.Fatalf("SearchVideos() error = %v", err)
			}
			if len(results) != tt.wantCount {
				t.Errorf("SearchVideos() got %d results, want %d", len(results), tt.wantCount)
			}
		})
	}
}

func TestSearchVideos_EmptyQuery(t *testing.T) {
	db := setupTestDB(t)

	_, err := SearchVideos(db, "", 50)
	if err == nil {
		t.Error("SearchVideos() should error on empty query")
	}
}

func TestSearchVideos_NilDB(t *testing.T) {
	_, err := SearchVideos(nil, "test", 50)
	if err == nil {
		t.Error("SearchVideos() should error on nil db")
	}
}

func TestSearchVideosByTitle(t *testing.T) {
	db := setupTestDB(t)

	// Insert test videos
	testVideos := []*Video{
		{ID: "vid1", Title: "Golang Programming", Description: "Learn Python too"},
		{ID: "vid2", Title: "Python Basics", Description: "With some Golang examples"},
	}

	for _, v := range testVideos {
		if err := InsertVideo(db, v); err != nil {
			t.Fatalf("Failed to insert video: %v", err)
		}
	}

	// Search by title only - "golang" appears in vid1 title only
	results, err := SearchVideosByTitle(db, "golang", 50)
	if err != nil {
		t.Fatalf("SearchVideosByTitle() error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("SearchVideosByTitle() got %d results, want 1", len(results))
	}
	if len(results) > 0 && results[0].ID != "vid1" {
		t.Errorf("SearchVideosByTitle() got wrong video, want vid1")
	}
}

func TestSearchVideosWithFilter(t *testing.T) {
	db := setupTestDB(t)

	// Insert test videos with different statuses
	testVideos := []*Video{
		{ID: "vid1", Title: "Golang Tutorial", Status: StatusPending},
		{ID: "vid2", Title: "Golang Advanced", Status: StatusCompleted},
		{ID: "vid3", Title: "Golang Basics", Status: StatusPending},
	}

	for _, v := range testVideos {
		if err := InsertVideo(db, v); err != nil {
			t.Fatalf("Failed to insert video: %v", err)
		}
	}

	// Search with status filter
	results, err := SearchVideosWithFilter(db, "golang", StatusPending, 50)
	if err != nil {
		t.Fatalf("SearchVideosWithFilter() error = %v", err)
	}
	if len(results) != 2 {
		t.Errorf("SearchVideosWithFilter() got %d results, want 2", len(results))
	}
}

func TestRebuildFTSIndex(t *testing.T) {
	db := setupTestDB(t)

	// Insert a video
	video := &Video{ID: "vid1", Title: "Test Video", Description: "Test Description"}
	if err := InsertVideo(db, video); err != nil {
		t.Fatalf("Failed to insert video: %v", err)
	}

	// Verify search works
	results, err := SearchVideos(db, "test", 50)
	if err != nil {
		t.Fatalf("SearchVideos() error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("SearchVideos() got %d results, want 1", len(results))
	}

	// Rebuild index
	if err := RebuildFTSIndex(db); err != nil {
		t.Fatalf("RebuildFTSIndex() error = %v", err)
	}

	// Verify search still works after rebuild
	results, err = SearchVideos(db, "test", 50)
	if err != nil {
		t.Fatalf("SearchVideos() after rebuild error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("SearchVideos() after rebuild got %d results, want 1", len(results))
	}
}

func TestPrepareFTSQuery(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "\"simple\""},
		{"two words", "\"two\" \"words\""},
		{"  spaces  ", "\"spaces\""},
		{"with\"quote", "\"with\"\"quote\""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := prepareFTSQuery(tt.input)
			if result != tt.expected {
				t.Errorf("prepareFTSQuery(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
