// Package db provides SQLite database operations for the YouTube Channel Archiver.
package db

// Schema defines the database schema for the channel metadata database.
const Schema = `
-- Videos table stores metadata and download status for each video
CREATE TABLE IF NOT EXISTS videos (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    duration INTEGER,
    upload_date TEXT,
    thumbnail_url TEXT,
    view_count INTEGER,
    status TEXT DEFAULT 'pending',
    file_path TEXT,
    file_size INTEGER,
    checksum TEXT,
    download_started_at DATETIME,
    download_completed_at DATETIME,
    retry_count INTEGER DEFAULT 0,
    last_error TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Index for efficient status queries
CREATE INDEX IF NOT EXISTS idx_videos_status ON videos(status);

-- Index for efficient upload date queries
CREATE INDEX IF NOT EXISTS idx_videos_upload_date ON videos(upload_date);

-- Sync history table tracks synchronization operations
CREATE TABLE IF NOT EXISTS sync_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    started_at DATETIME,
    completed_at DATETIME,
    videos_found INTEGER,
    videos_downloaded INTEGER,
    videos_failed INTEGER,
    status TEXT
);

-- Trigger to update the updated_at timestamp on video updates
CREATE TRIGGER IF NOT EXISTS update_videos_timestamp
AFTER UPDATE ON videos
FOR EACH ROW
BEGIN
    UPDATE videos SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

-- FTS5 virtual table for full-text search on title and description
CREATE VIRTUAL TABLE IF NOT EXISTS videos_fts USING fts5(
    title,
    description,
    content='videos',
    content_rowid='rowid'
);

-- Trigger to keep FTS index in sync on INSERT
CREATE TRIGGER IF NOT EXISTS videos_fts_insert AFTER INSERT ON videos
BEGIN
    INSERT INTO videos_fts(rowid, title, description)
    VALUES (new.rowid, new.title, new.description);
END;

-- Trigger to keep FTS index in sync on UPDATE
CREATE TRIGGER IF NOT EXISTS videos_fts_update AFTER UPDATE ON videos
BEGIN
    INSERT INTO videos_fts(videos_fts, rowid, title, description)
    VALUES ('delete', old.rowid, old.title, old.description);
    INSERT INTO videos_fts(rowid, title, description)
    VALUES (new.rowid, new.title, new.description);
END;

-- Trigger to keep FTS index in sync on DELETE
CREATE TRIGGER IF NOT EXISTS videos_fts_delete AFTER DELETE ON videos
BEGIN
    INSERT INTO videos_fts(videos_fts, rowid, title, description)
    VALUES ('delete', old.rowid, old.title, old.description);
END;
`

// VideoStatus represents the possible statuses for a video download.
type VideoStatus string

const (
	// StatusPending indicates the video is waiting to be downloaded
	StatusPending VideoStatus = "pending"
	// StatusDownloading indicates the video is currently being downloaded
	StatusDownloading VideoStatus = "downloading"
	// StatusCompleted indicates the video has been successfully downloaded
	StatusCompleted VideoStatus = "completed"
	// StatusFailed indicates the video download failed
	StatusFailed VideoStatus = "failed"
	// StatusSkipped indicates the video was skipped (e.g., already exists, too large)
	StatusSkipped VideoStatus = "skipped"
)

// SyncStatus represents the possible statuses for a sync operation.
type SyncStatus string

const (
	// SyncStatusRunning indicates the sync is in progress
	SyncStatusRunning SyncStatus = "running"
	// SyncStatusCompleted indicates the sync completed successfully
	SyncStatusCompleted SyncStatus = "completed"
	// SyncStatusFailed indicates the sync failed
	SyncStatusFailed SyncStatus = "failed"
)
