// Package db provides SQLite database operations for the YouTube Channel Archiver.
package db

import (
	"database/sql"
	"fmt"
	"strings"
)

// SearchResult represents a search result with relevance ranking.
type SearchResult struct {
	Video
	Rank float64 // BM25 relevance score (lower is more relevant)
}

// SearchVideos performs a full-text search on video titles and descriptions.
// Returns results ranked by relevance using BM25 algorithm.
func SearchVideos(db *sql.DB, query string, limit int) ([]SearchResult, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}
	if limit <= 0 {
		limit = 50
	}

	// Escape special FTS5 characters and prepare query
	searchQuery := prepareFTSQuery(query)

	sqlQuery := `
		SELECT v.id, v.title, v.description, v.duration, v.upload_date, v.thumbnail_url,
		       v.view_count, v.status, v.file_path, v.file_size, v.checksum,
		       v.download_started_at, v.download_completed_at, v.retry_count, v.last_error,
		       v.created_at, v.updated_at, bm25(videos_fts) as rank
		FROM videos_fts
		JOIN videos v ON videos_fts.rowid = v.rowid
		WHERE videos_fts MATCH ?
		ORDER BY rank
		LIMIT ?
	`

	rows, err := db.Query(sqlQuery, searchQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search videos: %w", err)
	}
	defer rows.Close()

	return scanSearchResults(rows)
}

// SearchVideosByTitle searches only in video titles.
func SearchVideosByTitle(db *sql.DB, query string, limit int) ([]SearchResult, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}
	if limit <= 0 {
		limit = 50
	}

	searchQuery := "title:" + prepareFTSQuery(query)

	sqlQuery := `
		SELECT v.id, v.title, v.description, v.duration, v.upload_date, v.thumbnail_url,
		       v.view_count, v.status, v.file_path, v.file_size, v.checksum,
		       v.download_started_at, v.download_completed_at, v.retry_count, v.last_error,
		       v.created_at, v.updated_at, bm25(videos_fts) as rank
		FROM videos_fts
		JOIN videos v ON videos_fts.rowid = v.rowid
		WHERE videos_fts MATCH ?
		ORDER BY rank
		LIMIT ?
	`

	rows, err := db.Query(sqlQuery, searchQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search videos by title: %w", err)
	}
	defer rows.Close()

	return scanSearchResults(rows)
}

// SearchVideosWithFilter searches videos with additional filters.
func SearchVideosWithFilter(db *sql.DB, query string, status VideoStatus, limit int) ([]SearchResult, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}
	if limit <= 0 {
		limit = 50
	}

	searchQuery := prepareFTSQuery(query)

	sqlQuery := `
		SELECT v.id, v.title, v.description, v.duration, v.upload_date, v.thumbnail_url,
		       v.view_count, v.status, v.file_path, v.file_size, v.checksum,
		       v.download_started_at, v.download_completed_at, v.retry_count, v.last_error,
		       v.created_at, v.updated_at, bm25(videos_fts) as rank
		FROM videos_fts
		JOIN videos v ON videos_fts.rowid = v.rowid
		WHERE videos_fts MATCH ? AND v.status = ?
		ORDER BY rank
		LIMIT ?
	`

	rows, err := db.Query(sqlQuery, searchQuery, status, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search videos with filter: %w", err)
	}
	defer rows.Close()

	return scanSearchResults(rows)
}

// RebuildFTSIndex rebuilds the FTS index from scratch.
// Use this if the index gets out of sync or corrupted.
func RebuildFTSIndex(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Delete all FTS entries
	_, err := db.Exec("DELETE FROM videos_fts")
	if err != nil {
		return fmt.Errorf("failed to clear FTS index: %w", err)
	}

	// Rebuild from videos table
	_, err = db.Exec(`
		INSERT INTO videos_fts(rowid, title, description)
		SELECT rowid, title, description FROM videos
	`)
	if err != nil {
		return fmt.Errorf("failed to rebuild FTS index: %w", err)
	}

	return nil
}

// prepareFTSQuery prepares a search query for FTS5.
// Escapes special characters and handles common search patterns.
func prepareFTSQuery(query string) string {
	// Trim whitespace
	query = strings.TrimSpace(query)

	// For simple queries, wrap each word in quotes to match exactly
	// This prevents FTS5 syntax errors from special characters
	words := strings.Fields(query)
	var escaped []string
	for _, word := range words {
		// Escape double quotes within the word
		word = strings.ReplaceAll(word, "\"", "\"\"")
		// Wrap in quotes for exact matching
		escaped = append(escaped, "\""+word+"\"")
	}

	// Join with implicit AND
	return strings.Join(escaped, " ")
}

// scanSearchResults scans search result rows into SearchResult slice.
func scanSearchResults(rows *sql.Rows) ([]SearchResult, error) {
	var results []SearchResult

	for rows.Next() {
		var result SearchResult
		var description, uploadDate, thumbnailURL, filePath, checksum, lastError sql.NullString
		var duration, viewCount, fileSize sql.NullInt64
		var downloadStartedAt, downloadCompletedAt sql.NullTime

		err := rows.Scan(
			&result.ID,
			&result.Title,
			&description,
			&duration,
			&uploadDate,
			&thumbnailURL,
			&viewCount,
			&result.Status,
			&filePath,
			&fileSize,
			&checksum,
			&downloadStartedAt,
			&downloadCompletedAt,
			&result.RetryCount,
			&lastError,
			&result.CreatedAt,
			&result.UpdatedAt,
			&result.Rank,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan search result: %w", err)
		}

		result.Description = description.String
		result.UploadDate = uploadDate.String
		result.ThumbnailURL = thumbnailURL.String
		result.FilePath = filePath.String
		result.Checksum = checksum.String
		result.LastError = lastError.String
		result.Duration = duration.Int64
		result.ViewCount = viewCount.Int64
		result.FileSize = fileSize.Int64

		if downloadStartedAt.Valid {
			result.DownloadStartedAt = &downloadStartedAt.Time
		}
		if downloadCompletedAt.Valid {
			result.DownloadCompletedAt = &downloadCompletedAt.Time
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating search results: %w", err)
	}

	return results, nil
}
