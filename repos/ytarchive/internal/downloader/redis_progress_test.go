package downloader

import (
	"testing"
	"time"
)

func TestProgressKey(t *testing.T) {
	key := progressKey("abc123")
	expected := "download:progress:abc123"
	if key != expected {
		t.Errorf("progressKey() = %s, want %s", key, expected)
	}
}

func TestRedisProgressReporterNilClient(t *testing.T) {
	// Test that the reporter doesn't panic with a nil client
	reporter := NewRedisProgressReporter(nil, "worker-1")

	progress := &DownloadProgress{
		VideoID:         "test123",
		WorkerID:        "worker-1",
		Status:          "downloading",
		Percentage:      50.0,
		DownloadedBytes: 1024,
		TotalBytes:      2048,
		Speed:           "1 MiB/s",
		ETA:             "00:05",
		UpdatedAt:       time.Now().Unix(),
	}

	// Should not panic with nil client
	err := reporter.ReportProgress(progress, false)
	if err != nil {
		t.Errorf("ReportProgress with nil client should not return error, got: %v", err)
	}
}

func TestProgressConstants(t *testing.T) {
	// Verify constants are set correctly
	if ProgressKeyPrefix != "download:progress:" {
		t.Errorf("ProgressKeyPrefix = %s, want download:progress:", ProgressKeyPrefix)
	}

	if ActiveDownloadsKey != "download:active" {
		t.Errorf("ActiveDownloadsKey = %s, want download:active", ActiveDownloadsKey)
	}

	if ProgressTTL != 5*time.Minute {
		t.Errorf("ProgressTTL = %v, want 5 minutes", ProgressTTL)
	}

	if MinProgressInterval != 1*time.Second {
		t.Errorf("MinProgressInterval = %v, want 1 second", MinProgressInterval)
	}
}

func TestDownloadProgressJSONFields(t *testing.T) {
	// Test that DownloadProgress has all required JSON fields
	progress := &DownloadProgress{
		VideoID:         "video123",
		WorkerID:        "worker456",
		Status:          "downloading",
		Percentage:      75.5,
		DownloadedBytes: 1000000,
		TotalBytes:      1500000,
		Speed:           "2.5 MiB/s",
		ETA:             "02:30",
		Fragment:        "5/10",
		UpdatedAt:       1234567890,
	}

	// Verify all fields are accessible
	if progress.VideoID != "video123" {
		t.Error("VideoID not set correctly")
	}
	if progress.WorkerID != "worker456" {
		t.Error("WorkerID not set correctly")
	}
	if progress.Status != "downloading" {
		t.Error("Status not set correctly")
	}
	if progress.Percentage != 75.5 {
		t.Error("Percentage not set correctly")
	}
	if progress.DownloadedBytes != 1000000 {
		t.Error("DownloadedBytes not set correctly")
	}
	if progress.TotalBytes != 1500000 {
		t.Error("TotalBytes not set correctly")
	}
	if progress.Speed != "2.5 MiB/s" {
		t.Error("Speed not set correctly")
	}
	if progress.ETA != "02:30" {
		t.Error("ETA not set correctly")
	}
	if progress.Fragment != "5/10" {
		t.Error("Fragment not set correctly")
	}
	if progress.UpdatedAt != 1234567890 {
		t.Error("UpdatedAt not set correctly")
	}
}
