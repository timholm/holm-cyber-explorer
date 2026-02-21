package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/timholm/animus-dashboard/internal/api"
	"github.com/timholm/animus-dashboard/internal/ansible"
	"github.com/timholm/animus-dashboard/internal/auth"
	"github.com/timholm/animus-dashboard/internal/k8s"
	"github.com/timholm/animus-dashboard/internal/loki"
	"github.com/timholm/animus-dashboard/internal/ssh"
)

func main() {
	// Configuration
	config := loadConfig()

	// Initialize Kubernetes client
	k8sClient, err := k8s.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}
	log.Println("Kubernetes client initialized")

	// Initialize Loki client
	lokiClient := loki.NewClient(config.LokiURL)
	log.Printf("Loki client initialized: %s", config.LokiURL)

	// Initialize SSH client
	sshClient, err := ssh.NewClient(config.SSHKeyPath, config.SSHUser)
	if err != nil {
		log.Printf("Warning: SSH client initialization failed: %v", err)
	} else {
		log.Println("SSH client initialized")
	}

	// Initialize Ansible runner
	ansibleRunner := ansible.NewRunner(config.AnsiblePath, config.InventoryPath)
	log.Println("Ansible runner initialized")

	// Initialize Keycloak auth (optional)
	var authMiddleware fiber.Handler
	if config.KeycloakURL != "" {
		keycloak, err := auth.NewKeycloakAuth(config.KeycloakURL, config.KeycloakRealm, config.KeycloakClientID)
		if err != nil {
			log.Printf("Warning: Keycloak auth initialization failed: %v", err)
		} else {
			authMiddleware = keycloak.Middleware()
			log.Println("Keycloak authentication enabled")
		}
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Animus Dashboard API",
		ServerHeader: "Animus",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.CORSOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Create API handlers
	handlers := api.NewHandlers(k8sClient, lokiClient, sshClient, ansibleRunner)

	// API routes
	apiGroup := app.Group("/api")
	if authMiddleware != nil {
		apiGroup.Use(authMiddleware)
	}

	// Node routes
	apiGroup.Get("/nodes", handlers.GetNodes)
	apiGroup.Get("/nodes/:id", handlers.GetNode)
	apiGroup.Get("/nodes/:id/pods", handlers.GetNodePods)

	// Cluster routes
	apiGroup.Get("/cluster/metrics", handlers.GetClusterMetrics)

	// Log routes
	apiGroup.Get("/logs", handlers.GetLogs)

	// Script routes
	apiGroup.Get("/scripts", handlers.GetScripts)
	apiGroup.Post("/scripts/run", handlers.RunScript)
	apiGroup.Get("/scripts/executions", handlers.GetScriptExecutions)
	apiGroup.Get("/scripts/executions/:id", handlers.GetScriptExecution)

	// Version routes
	apiGroup.Get("/versions/check", handlers.CheckUpdates)

	// WebSocket routes
	wsHandlers := api.NewWebSocketHandlers(k8sClient, lokiClient, ansibleRunner)
	app.Get("/ws/logs", wsHandlers.LogStream)
	app.Get("/ws/metrics", wsHandlers.MetricsStream)
	app.Get("/ws/scripts/:id", wsHandlers.ScriptOutput)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "animus-dashboard",
		})
	})

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()

	// Start server
	addr := ":" + config.Port
	log.Printf("Animus Dashboard API starting on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

type Config struct {
	Port             string
	LokiURL          string
	SSHKeyPath       string
	SSHUser          string
	AnsiblePath      string
	InventoryPath    string
	KeycloakURL      string
	KeycloakRealm    string
	KeycloakClientID string
	CORSOrigins      string
}

func loadConfig() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		LokiURL:          getEnv("LOKI_URL", "http://loki.monitoring:3100"),
		SSHKeyPath:       getEnv("SSH_KEY_PATH", "/etc/animus/ssh/id_ed25519"),
		SSHUser:          getEnv("SSH_USER", "tim"),
		AnsiblePath:      getEnv("ANSIBLE_PATH", "/app/ansible/playbooks"),
		InventoryPath:    getEnv("INVENTORY_PATH", "/app/ansible/inventory.ini"),
		KeycloakURL:      getEnv("KEYCLOAK_URL", ""),
		KeycloakRealm:    getEnv("KEYCLOAK_REALM", "master"),
		KeycloakClientID: getEnv("KEYCLOAK_CLIENT_ID", "animus-dashboard"),
		CORSOrigins:      getEnv("CORS_ORIGINS", "*"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
