package api

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/timholm/ytarchive/web"
)

// SetupRoutes configures all API routes
func SetupRoutes(handlers *Handlers) *gin.Engine {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(requestLogger())
	router.Use(corsMiddleware())

	// Health check handler
	healthHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	}

	// Readiness check handler
	readyHandler := func(c *gin.Context) {
		// Check Redis connectivity
		if err := handlers.redis.Ping(c.Request.Context()).Err(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"error":  "Redis connection failed",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	}

	// Serve embedded frontend
	frontendFS, err := web.GetFS()
	if err == nil {
		// Read index.html once at startup
		indexHTML, readErr := fs.ReadFile(frontendFS, "index.html")
		if readErr != nil {
			// Fallback to inline HTML if index.html read fails
			router.GET("/", handlers.ServeFrontend)
		} else {
			// Serve static assets from /assets/
			assetsFS, _ := fs.Sub(frontendFS, "assets")
			router.StaticFS("/assets", http.FS(assetsFS))

			// Serve index.html for root
			router.GET("/", func(c *gin.Context) {
				c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
			})

			// SPA fallback - serve index.html for any non-API routes
			router.NoRoute(func(c *gin.Context) {
				// Don't serve index.html for API routes
				if len(c.Request.URL.Path) > 4 && c.Request.URL.Path[:4] == "/api" {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
			})
		}
	} else {
		// Fallback to inline HTML if embed fails
		router.GET("/", handlers.ServeFrontend)
	}

	// Health check endpoints (both standard and K8s-style)
	router.GET("/health", healthHandler)
	router.GET("/healthz", healthHandler)

	// Readiness check endpoints (both standard and K8s-style)
	router.GET("/ready", readyHandler)
	router.GET("/readyz", readyHandler)

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	api := router.Group("/api")
	{
		// Dashboard endpoints
		api.GET("/stats", handlers.GetStats)
		api.GET("/activity", handlers.GetActivity)

		// Channel endpoints
		channels := api.Group("/channels")
		{
			channels.POST("", handlers.AddChannel)
			channels.GET("", handlers.ListChannels)
			channels.GET("/:id", handlers.GetChannel)
			channels.POST("/:id/sync", handlers.SyncChannel)
			channels.POST("/:id/index", handlers.IndexChannelVideos)
			channels.DELETE("/:id", handlers.DeleteChannel)
			channels.GET("/:id/videos", handlers.GetChannelVideos)
		}

		// Index endpoint - rebuild FTS index for all channels
		api.POST("/index", handlers.IndexAllChannels)

		// Search endpoint
		api.GET("/search", handlers.SearchVideos)

		// Video endpoints
		videos := api.Group("/videos")
		{
			videos.GET("", handlers.ListVideos)
			videos.GET("/:id", handlers.GetVideo)
			videos.GET("/:id/stream", handlers.GetVideoStream)
			videos.GET("/:id/thumbnail", handlers.GetVideoThumbnail)
			videos.GET("/:id/metadata", handlers.GetVideoMetadata)
			videos.GET("/:id/audio", handlers.GetVideoAudio)
			videos.GET("/:id/subtitles", handlers.GetVideoSubtitles)
			videos.GET("/:id/subtitles/list", handlers.ListVideoSubtitles)
			videos.GET("/:id/files", handlers.GetVideoFiles)
			videos.POST("/:id/download", handlers.TriggerDownload)
			videos.POST("/:id/status", handlers.UpdateVideoStatus)
		}

		// Job endpoints
		jobs := api.Group("/jobs")
		{
			jobs.GET("", handlers.ListJobs)
			jobs.GET("/progress", handlers.GetJobsProgress)
			jobs.GET("/:id", handlers.GetJob)
			jobs.POST("/:id/cancel", handlers.CancelJob)
		}

		// Progress endpoint (legacy, kept for compatibility)
		api.GET("/progress", handlers.GetProgress)

		// Worker progress reporting endpoint
		api.POST("/progress/:id", handlers.UpdateProgress)

		// Downloads progress endpoint - returns all active download progress
		downloads := api.Group("/downloads")
		{
			downloads.GET("/progress", handlers.GetDownloadsProgress)
		}

		// Cookies endpoints
		api.GET("/cookies", handlers.GetCookies)
		api.POST("/cookies", handlers.SaveCookies)
		api.DELETE("/cookies", handlers.DeleteCookies)
	}

	return router
}

// requestLogger is a middleware that logs HTTP requests
func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request details
		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		if query != "" {
			path = path + "?" + query
		}

		gin.DefaultWriter.Write([]byte(fmt.Sprintf(
			"%s | %s | %s | %s | %d %s | %s\n",
			time.Now().Format(time.RFC3339),
			method,
			path,
			clientIP,
			status,
			http.StatusText(status),
			latency.String(),
		)))
	}
}

// corsMiddleware adds CORS headers to responses
func corsMiddleware() gin.HandlerFunc {
	// Get CORS origin from environment variable, default to "*" if not set
	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "*"
	}

	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", corsOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
