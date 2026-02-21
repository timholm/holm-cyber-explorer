package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"

	"github.com/timholm/ytarchive/internal/logging"
)

const (
	httpPort         = 8081
	maxUploadSize    = 10 * 1024 * 1024 * 1024 // 10GB max file size
	videoKeyPrefix   = "video:"
	channelKeyPrefix = "channel:"
)

// CollectorConfig holds the collector configuration
type CollectorConfig struct {
	StoragePath  string
	PostgresHost string
	PostgresPort string
	PostgresDB   string
	PostgresUser string
	PostgresPass string
	RedisURL     string
}

// Collector handles receiving and storing video files
type Collector struct {
	config *CollectorConfig
	db     *sql.DB
	redis  *redis.Client
	ready  atomic.Bool
}

// UploadRequest contains metadata for an uploaded video
type UploadRequest struct {
	VideoID       string `json:"video_id"`
	ChannelID     string `json:"channel_id"`
	ChannelName   string `json:"channel_name"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Duration      int    `json:"duration"`
	UploadDate    string `json:"upload_date"`
	EpisodeNumber int    `json:"episode_number"`
	Filename      string `json:"filename"`
	FileSize      int64  `json:"file_size"`
	Resolution    string `json:"resolution"`
	Format        string `json:"format"`
}

func main() {
	logging.Info("starting YouTube Channel Archiver Collector")

	config, err := loadConfig()
	if err != nil {
		logging.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	collector := &Collector{config: config}

	// Connect to PostgreSQL
	if err := collector.connectPostgres(); err != nil {
		logging.Error("failed to connect to PostgreSQL", "error", err)
		os.Exit(1)
	}
	defer collector.db.Close()
	logging.Info("connected to PostgreSQL")

	// Initialize database schema
	if err := collector.initSchema(); err != nil {
		logging.Error("failed to initialize database schema", "error", err)
		os.Exit(1)
	}

	// Connect to Redis
	if err := collector.connectRedis(); err != nil {
		logging.Warn("failed to connect to Redis, continuing without it", "error", err)
	} else {
		defer collector.redis.Close()
		logging.Info("connected to Redis")
	}

	// Ensure storage directory exists
	if err := os.MkdirAll(config.StoragePath, 0755); err != nil {
		logging.Error("failed to create storage directory", "error", err)
		os.Exit(1)
	}

	// Set up HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", collector.healthzHandler)
	mux.HandleFunc("/readyz", collector.readyzHandler)
	mux.HandleFunc("/upload", collector.uploadHandler)
	mux.HandleFunc("/api/videos", collector.listVideosHandler)
	mux.HandleFunc("/api/channels", collector.listChannelsHandler)
	mux.HandleFunc("/api/stats", collector.statsHandler)
	mux.HandleFunc("/stream/", collector.streamVideoHandler)
	mux.HandleFunc("/thumbnail/", collector.thumbnailHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", httpPort),
		Handler:      mux,
		ReadTimeout:  0, // No timeout for large uploads
		WriteTimeout: 0,
	}

	// Handle shutdown signals
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logging.Info("received shutdown signal", "signal", sig.String())
		cancel()
		server.Shutdown(context.Background())
	}()

	collector.ready.Store(true)
	logging.Info("collector ready, listening", "port", httpPort)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logging.Error("server error", "error", err)
	}

	logging.Info("collector shutdown complete")
}

func loadConfig() (*CollectorConfig, error) {
	config := &CollectorConfig{
		StoragePath:  os.Getenv("STORAGE_PATH"),
		PostgresHost: os.Getenv("POSTGRES_HOST"),
		PostgresPort: os.Getenv("POSTGRES_PORT"),
		PostgresDB:   os.Getenv("POSTGRES_DB"),
		PostgresUser: os.Getenv("POSTGRES_USER"),
		PostgresPass: os.Getenv("POSTGRES_PASSWORD"),
		RedisURL:     os.Getenv("REDIS_URL"),
	}

	if config.StoragePath == "" {
		config.StoragePath = "/data"
	}
	if config.PostgresHost == "" {
		config.PostgresHost = "postgres"
	}
	if config.PostgresPort == "" {
		config.PostgresPort = "5432"
	}
	if config.PostgresDB == "" {
		config.PostgresDB = "ytarchive"
	}
	if config.PostgresUser == "" {
		config.PostgresUser = "postgres"
	}
	if config.PostgresPass == "" {
		config.PostgresPass = "postgres"
	}

	return config, nil
}

func (c *Collector) connectPostgres() error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.config.PostgresHost,
		c.config.PostgresPort,
		c.config.PostgresUser,
		c.config.PostgresPass,
		c.config.PostgresDB,
	)

	var err error
	for i := 0; i < 30; i++ {
		c.db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = c.db.Ping()
			if err == nil {
				return nil
			}
		}
		logging.Info("waiting for PostgreSQL...", "attempt", i+1)
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("failed to connect to PostgreSQL after 30 attempts: %w", err)
}

func (c *Collector) connectRedis() error {
	if c.config.RedisURL == "" {
		return fmt.Errorf("Redis URL not configured")
	}

	opts := &redis.Options{
		Addr: c.config.RedisURL,
	}

	if len(c.config.RedisURL) > 8 && c.config.RedisURL[:8] == "redis://" {
		parsedOpts, err := redis.ParseURL(c.config.RedisURL)
		if err != nil {
			return err
		}
		opts = parsedOpts
	}

	c.redis = redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.redis.Ping(ctx).Err()
}

func (c *Collector) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS channels (
		id VARCHAR(64) PRIMARY KEY,
		youtube_id VARCHAR(64) NOT NULL,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		video_count INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS videos (
		id VARCHAR(64) PRIMARY KEY,
		channel_id VARCHAR(64) NOT NULL REFERENCES channels(id),
		title VARCHAR(500) NOT NULL,
		description TEXT,
		duration INTEGER,
		upload_date VARCHAR(20),
		episode_number INTEGER,
		filename VARCHAR(500),
		file_path VARCHAR(1000),
		file_size BIGINT,
		resolution VARCHAR(20),
		format VARCHAR(20),
		status VARCHAR(20) DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_videos_channel ON videos(channel_id);
	CREATE INDEX IF NOT EXISTS idx_videos_status ON videos(status);
	`

	_, err := c.db.Exec(schema)
	return err
}

func (c *Collector) healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (c *Collector) readyzHandler(w http.ResponseWriter, r *http.Request) {
	if c.ready.Load() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("not ready"))
	}
}

func (c *Collector) uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(100 * 1024 * 1024); err != nil { // 100MB memory limit
		logging.Error("failed to parse multipart form", "error", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	// Get metadata
	metadataStr := r.FormValue("metadata")
	if metadataStr == "" {
		http.Error(w, "Missing metadata", http.StatusBadRequest)
		return
	}

	var metadata UploadRequest
	if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
		logging.Error("failed to parse metadata", "error", err)
		http.Error(w, "Invalid metadata", http.StatusBadRequest)
		return
	}

	// Get file
	file, header, err := r.FormFile("video")
	if err != nil {
		logging.Error("failed to get video file", "error", err)
		http.Error(w, "Missing video file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	logging.Info("receiving video upload",
		"video_id", metadata.VideoID,
		"channel_id", metadata.ChannelID,
		"filename", header.Filename,
		"size", header.Size,
	)

	// Ensure channel exists in database
	if err := c.ensureChannel(metadata.ChannelID, metadata.ChannelName); err != nil {
		logging.Error("failed to ensure channel exists", "error", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Create destination directory
	destDir := filepath.Join(c.config.StoragePath, "channels", metadata.ChannelID, "videos", metadata.VideoID)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		logging.Error("failed to create destination directory", "error", err)
		http.Error(w, "Storage error", http.StatusInternalServerError)
		return
	}

	// Determine filename
	filename := metadata.Filename
	if filename == "" {
		filename = header.Filename
	}
	destPath := filepath.Join(destDir, filename)

	// Write file
	destFile, err := os.Create(destPath)
	if err != nil {
		logging.Error("failed to create destination file", "error", err)
		http.Error(w, "Storage error", http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	written, err := io.Copy(destFile, file)
	if err != nil {
		logging.Error("failed to write file", "error", err)
		os.Remove(destPath)
		http.Error(w, "Storage error", http.StatusInternalServerError)
		return
	}

	// Store metadata in PostgreSQL
	if err := c.storeVideoMetadata(&metadata, destPath, written); err != nil {
		logging.Error("failed to store video metadata", "error", err)
		// Don't delete the file, just log the error
	}

	// Update Redis if available
	if c.redis != nil {
		c.updateRedisVideoStatus(metadata.ChannelID, metadata.VideoID, "downloaded", destPath, written)
	}

	logging.Info("video upload complete",
		"video_id", metadata.VideoID,
		"channel_id", metadata.ChannelID,
		"file_path", destPath,
		"file_size", written,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"video_id":  metadata.VideoID,
		"file_path": destPath,
		"file_size": written,
	})
}

func (c *Collector) ensureChannel(channelID, channelName string) error {
	query := `
		INSERT INTO channels (id, youtube_id, name)
		VALUES ($1, $1, $2)
		ON CONFLICT (id) DO UPDATE SET
			name = COALESCE(NULLIF($2, ''), channels.name),
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := c.db.Exec(query, channelID, channelName)
	return err
}

func (c *Collector) storeVideoMetadata(metadata *UploadRequest, filePath string, fileSize int64) error {
	query := `
		INSERT INTO videos (id, channel_id, title, description, duration, upload_date,
		                    episode_number, filename, file_path, file_size, resolution, format, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, 'downloaded')
		ON CONFLICT (id) DO UPDATE SET
			title = $3,
			description = $4,
			file_path = $9,
			file_size = $10,
			resolution = $11,
			format = $12,
			status = 'downloaded',
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := c.db.Exec(query,
		metadata.VideoID,
		metadata.ChannelID,
		metadata.Title,
		metadata.Description,
		metadata.Duration,
		metadata.UploadDate,
		metadata.EpisodeNumber,
		metadata.Filename,
		filePath,
		fileSize,
		metadata.Resolution,
		metadata.Format,
	)
	return err
}

func (c *Collector) updateRedisVideoStatus(channelID, videoID, status, filePath string, fileSize int64) {
	ctx := context.Background()
	videoKey := videoKeyPrefix + channelID + ":" + videoID

	// Get existing video data
	videoData, err := c.redis.Get(ctx, videoKey).Result()
	if err != nil {
		return
	}

	var video map[string]interface{}
	if err := json.Unmarshal([]byte(videoData), &video); err != nil {
		return
	}

	video["status"] = status
	video["file_path"] = filePath
	video["file_size"] = fileSize
	video["updated_at"] = time.Now()

	updatedData, _ := json.Marshal(video)
	c.redis.Set(ctx, videoKey, updatedData, 0)
}

func (c *Collector) listVideosHandler(w http.ResponseWriter, r *http.Request) {
	channelID := r.URL.Query().Get("channel_id")

	var rows *sql.Rows
	var err error

	if channelID != "" {
		rows, err = c.db.Query(`
			SELECT id, channel_id, title, duration, upload_date, episode_number,
			       filename, file_path, file_size, resolution, status
			FROM videos
			WHERE channel_id = $1
			ORDER BY episode_number DESC
		`, channelID)
	} else {
		rows, err = c.db.Query(`
			SELECT id, channel_id, title, duration, upload_date, episode_number,
			       filename, file_path, file_size, resolution, status
			FROM videos
			ORDER BY created_at DESC
			LIMIT 100
		`)
	}

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var videos []map[string]interface{}
	for rows.Next() {
		var id, channelID, title, status string
		var uploadDate, filename, filePath, resolution sql.NullString
		var duration, episodeNumber sql.NullInt64
		var fileSize sql.NullInt64

		err := rows.Scan(&id, &channelID, &title, &duration, &uploadDate,
			&episodeNumber, &filename, &filePath, &fileSize, &resolution, &status)
		if err != nil {
			continue
		}

		videos = append(videos, map[string]interface{}{
			"id":             id,
			"channel_id":     channelID,
			"title":          title,
			"duration":       duration.Int64,
			"upload_date":    uploadDate.String,
			"episode_number": episodeNumber.Int64,
			"filename":       filename.String,
			"file_path":      filePath.String,
			"file_size":      fileSize.Int64,
			"resolution":     resolution.String,
			"status":         status,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"videos": videos,
		"count":  len(videos),
	})
}

func (c *Collector) listChannelsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := c.db.Query(`
		SELECT c.id, c.name, c.video_count,
		       (SELECT COUNT(*) FROM videos WHERE channel_id = c.id AND status = 'downloaded') as downloaded
		FROM channels c
		ORDER BY c.name
	`)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var channels []map[string]interface{}
	for rows.Next() {
		var id, name string
		var videoCount, downloaded int

		if err := rows.Scan(&id, &name, &videoCount, &downloaded); err != nil {
			continue
		}

		channels = append(channels, map[string]interface{}{
			"id":          id,
			"name":        name,
			"video_count": videoCount,
			"downloaded":  downloaded,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"channels": channels,
		"count":    len(channels),
	})
}

func (c *Collector) statsHandler(w http.ResponseWriter, r *http.Request) {
	var totalChannels, totalVideos, downloadedVideos int
	var totalSize int64

	c.db.QueryRow("SELECT COUNT(*) FROM channels").Scan(&totalChannels)
	c.db.QueryRow("SELECT COUNT(*) FROM videos").Scan(&totalVideos)
	c.db.QueryRow("SELECT COUNT(*) FROM videos WHERE status = 'downloaded'").Scan(&downloadedVideos)
	c.db.QueryRow("SELECT COALESCE(SUM(file_size), 0) FROM videos").Scan(&totalSize)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_channels":    totalChannels,
		"total_videos":      totalVideos,
		"downloaded_videos": downloadedVideos,
		"total_size_bytes":  totalSize,
		"total_size_human":  formatSize(totalSize),
	})
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// streamVideoHandler serves video files for streaming
// URL format: /stream/{channel_id}/{video_id}
func (c *Collector) streamVideoHandler(w http.ResponseWriter, r *http.Request) {
	// Parse path: /stream/{channel_id}/{video_id}
	path := r.URL.Path[len("/stream/"):]
	parts := filepath.SplitList(path)
	if len(parts) == 1 {
		// Try splitting by /
		parts = splitPath(path)
	}

	if len(parts) < 2 {
		http.Error(w, "Invalid path. Expected /stream/{channel_id}/{video_id}", http.StatusBadRequest)
		return
	}

	channelID := parts[0]
	videoID := parts[1]

	// Find video file
	videoDir := filepath.Join(c.config.StoragePath, "channels", channelID, "videos", videoID)

	// Look for video file with various extensions
	patterns := []string{"video.mp4", "video.mkv", "video.webm", "*.mp4", "*.mkv", "*.webm"}
	var videoPath string
	for _, pattern := range patterns {
		matches, _ := filepath.Glob(filepath.Join(videoDir, pattern))
		if len(matches) > 0 {
			videoPath = matches[0]
			break
		}
	}

	if videoPath == "" {
		logging.Warn("video file not found", "channel_id", channelID, "video_id", videoID, "dir", videoDir)
		http.Error(w, "Video file not found", http.StatusNotFound)
		return
	}

	// Get file info for Content-Length
	fileInfo, err := os.Stat(videoPath)
	if err != nil {
		http.Error(w, "Failed to read video file", http.StatusInternalServerError)
		return
	}

	// Set headers for video streaming
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Type", getContentType(videoPath))

	// Serve the file (http.ServeFile handles range requests)
	http.ServeFile(w, r, videoPath)

	logging.Info("streaming video",
		"channel_id", channelID,
		"video_id", videoID,
		"file_size", fileInfo.Size(),
	)
}

// thumbnailHandler serves video thumbnails
// URL format: /thumbnail/{channel_id}/{video_id}
func (c *Collector) thumbnailHandler(w http.ResponseWriter, r *http.Request) {
	// Parse path: /thumbnail/{channel_id}/{video_id}
	path := r.URL.Path[len("/thumbnail/"):]
	parts := splitPath(path)

	if len(parts) < 2 {
		http.Error(w, "Invalid path. Expected /thumbnail/{channel_id}/{video_id}", http.StatusBadRequest)
		return
	}

	channelID := parts[0]
	videoID := parts[1]

	// Find thumbnail file
	videoDir := filepath.Join(c.config.StoragePath, "channels", channelID, "videos", videoID)

	// Look for thumbnail
	patterns := []string{"thumbnail.jpg", "thumbnail.png", "thumbnail.webp", "*.jpg", "*.png", "*.webp"}
	var thumbPath string
	for _, pattern := range patterns {
		matches, _ := filepath.Glob(filepath.Join(videoDir, pattern))
		if len(matches) > 0 {
			thumbPath = matches[0]
			break
		}
	}

	if thumbPath == "" {
		http.Error(w, "Thumbnail not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, thumbPath)
}

func splitPath(path string) []string {
	var parts []string
	for _, p := range filepath.SplitList(path) {
		parts = append(parts, p)
	}
	// If SplitList didn't work (Unix), split by /
	if len(parts) == 1 && len(path) > 0 {
		parts = []string{}
		for _, p := range filepath.Clean(path) {
			if p == filepath.Separator || p == '/' {
				continue
			}
		}
		// Manual split
		current := ""
		for _, c := range path {
			if c == '/' {
				if current != "" {
					parts = append(parts, current)
					current = ""
				}
			} else {
				current += string(c)
			}
		}
		if current != "" {
			parts = append(parts, current)
		}
	}
	return parts
}

func getContentType(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".mkv":
		return "video/x-matroska"
	case ".webm":
		return "video/webm"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
