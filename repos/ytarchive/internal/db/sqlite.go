// Package db provides SQLite database operations for the YouTube Channel Archiver.
// Each channel gets its own SQLite database at: {STORAGE_PATH}/channels/{channel_id}/metadata.db
package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const (
	// DefaultDataDir is the default base directory for all channel data
	DefaultDataDir = "/data/channels"
	// DBFileName is the name of the SQLite database file for each channel
	DBFileName = "metadata.db"
)

// getDataDir returns the data directory from STORAGE_PATH env var or the default
func getDataDir() string {
	// Try STORAGE_PATH first
	if storagePath := os.Getenv("STORAGE_PATH"); storagePath != "" {
		channelDir := filepath.Join(storagePath, "channels")
		// Check if we can create the directory
		if err := os.MkdirAll(channelDir, 0755); err == nil {
			return channelDir
		}
	}

	// Try DB_PATH for dedicated database storage
	if dbPath := os.Getenv("DB_PATH"); dbPath != "" {
		if err := os.MkdirAll(dbPath, 0755); err == nil {
			return dbPath
		}
	}

	// Fallback to temp directory for read-only root filesystem
	tmpDir := filepath.Join(os.TempDir(), "ytarchive", "channels")
	os.MkdirAll(tmpDir, 0755)
	return tmpDir
}

// DataDir returns the base directory for all channel data
func DataDir() string {
	return getDataDir()
}

// OpenChannelDB opens or creates a SQLite database for the specified channel.
// The database is stored at /data/channels/{channel_id}/metadata.db
func OpenChannelDB(channelID string) (*sql.DB, error) {
	if channelID == "" {
		return nil, fmt.Errorf("channel ID cannot be empty")
	}

	// Construct the path to the channel's database
	channelDir := filepath.Join(DataDir(), channelID)
	dbPath := filepath.Join(channelDir, DBFileName)

	// Ensure the channel directory exists
	if err := os.MkdirAll(channelDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create channel directory: %w", err)
	}

	// Open the SQLite database with foreign keys enabled
	db, err := sql.Open("sqlite", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize the schema
	if err := InitSchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return db, nil
}

// InitSchema creates the required database tables if they don't exist.
func InitSchema(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Execute the schema creation SQL
	_, err := db.Exec(Schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

// Close closes the database connection.
func Close(db *sql.DB) error {
	if db == nil {
		return nil
	}
	return db.Close()
}

// GetDBPath returns the path to the database file for a channel.
func GetDBPath(channelID string) string {
	return filepath.Join(DataDir(), channelID, DBFileName)
}

// ChannelDBExists checks if a database exists for the specified channel.
func ChannelDBExists(channelID string) bool {
	dbPath := GetDBPath(channelID)
	_, err := os.Stat(dbPath)
	return err == nil
}
