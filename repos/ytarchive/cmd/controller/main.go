package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/timholm/ytarchive/internal/api"
	"github.com/timholm/ytarchive/internal/logging"
	"github.com/timholm/ytarchive/internal/scheduler"
)

func main() {
	logging.Info("starting YouTube Channel Archiver Controller")

	// Initialize Redis client
	redisClient, err := initRedis()
	if err != nil {
		logging.Error("failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()
	logging.Info("connected to Redis")

	// Initialize Kubernetes client
	k8sClient, err := initK8sClient()
	if err != nil {
		logging.Error("failed to create Kubernetes client", "error", err)
		os.Exit(1)
	}
	logging.Info("connected to Kubernetes")

	// Initialize scheduler
	jobScheduler := scheduler.NewScheduler(k8sClient, redisClient, getEnvWithDefault("K8S_NAMESPACE", "default"))

	// Initialize API handlers
	handlers := api.NewHandlers(redisClient, jobScheduler)

	// Set up Gin router
	router := api.SetupRoutes(handlers)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logging.Info("starting HTTP server", "addr", ":8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logging.Info("shutting down server")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logging.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logging.Info("server exited gracefully")
}

// initRedis creates and tests a Redis connection
func initRedis() (*redis.Client, error) {
	redisAddr := getEnvWithDefault("REDIS_URL", "localhost:6379")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0

	client := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		Password:     redisPassword,
		DB:           redisDB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

// initK8sClient creates a Kubernetes client
func initK8sClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	// Try in-cluster config first (when running inside K8s)
	config, err = rest.InClusterConfig()
	if err != nil {
		// Fall back to kubeconfig file (for local development)
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = os.Getenv("HOME") + "/.kube/config"
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		logging.Info("using kubeconfig from file", "path", kubeconfig)
	} else {
		logging.Info("using in-cluster Kubernetes config")
	}

	return kubernetes.NewForConfig(config)
}

// getEnvWithDefault returns environment variable value or default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
