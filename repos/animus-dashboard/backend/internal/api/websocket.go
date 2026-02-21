package api

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/timholm/animus-dashboard/internal/ansible"
	"github.com/timholm/animus-dashboard/internal/k8s"
	"github.com/timholm/animus-dashboard/internal/loki"
)

type WebSocketHandlers struct {
	k8sClient     *k8s.Client
	lokiClient    *loki.Client
	ansibleRunner *ansible.Runner
}

func NewWebSocketHandlers(k8s *k8s.Client, loki *loki.Client, ansible *ansible.Runner) *WebSocketHandlers {
	return &WebSocketHandlers{
		k8sClient:     k8s,
		lokiClient:    loki,
		ansibleRunner: ansible,
	}
}

// LogStream handles WebSocket connections for log streaming
func (h *WebSocketHandlers) LogStream(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return websocket.New(func(conn *websocket.Conn) {
			defer conn.Close()

			node := conn.Query("node")
			namespace := conn.Query("namespace")

			params := loki.QueryParams{
				Node:      node,
				Namespace: namespace,
				Limit:     50,
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Start tailing logs
			logCh, err := h.lokiClient.TailLogs(ctx, params)
			if err != nil {
				log.Printf("Failed to tail logs: %v", err)
				return
			}

			// Read messages from client (for ping/pong or close)
			go func() {
				for {
					_, _, err := conn.ReadMessage()
					if err != nil {
						cancel()
						return
					}
				}
			}()

			// Send logs to client
			for {
				select {
				case <-ctx.Done():
					return
				case entry, ok := <-logCh:
					if !ok {
						return
					}

					data := struct {
						Entries []loki.LogEntry `json:"entries"`
					}{
						Entries: []loki.LogEntry{entry},
					}

					msg, err := json.Marshal(data)
					if err != nil {
						continue
					}

					if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
						return
					}
				}
			}
		})(c)
	}

	return fiber.ErrUpgradeRequired
}

// MetricsStream handles WebSocket connections for real-time metrics
func (h *WebSocketHandlers) MetricsStream(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return websocket.New(func(conn *websocket.Conn) {
			defer conn.Close()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Read messages from client (for ping/pong or close)
			go func() {
				for {
					_, _, err := conn.ReadMessage()
					if err != nil {
						cancel()
						return
					}
				}
			}()

			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			// Send initial metrics
			h.sendMetrics(ctx, conn)

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					if err := h.sendMetrics(ctx, conn); err != nil {
						return
					}
				}
			}
		})(c)
	}

	return fiber.ErrUpgradeRequired
}

func (h *WebSocketHandlers) sendMetrics(ctx context.Context, conn *websocket.Conn) error {
	nodes, err := h.k8sClient.GetNodes(ctx)
	if err != nil {
		return err
	}

	clusterMetrics, _ := h.k8sClient.GetClusterMetrics(ctx)

	type NodeMetric struct {
		Name     string `json:"name"`
		CPU      int    `json:"cpu"`
		Memory   int    `json:"memory"`
		Disk     int    `json:"disk"`
		PodCount int    `json:"podCount"`
		Status   string `json:"status"`
	}

	type MetricsData struct {
		Nodes   []NodeMetric `json:"nodes"`
		Cluster struct {
			TotalCPU     int `json:"totalCpu"`
			TotalMemory  int `json:"totalMemory"`
			HealthyNodes int `json:"healthyNodes"`
			TotalNodes   int `json:"totalNodes"`
		} `json:"cluster"`
	}

	data := MetricsData{}
	for _, node := range nodes {
		data.Nodes = append(data.Nodes, NodeMetric{
			Name:     node.Name,
			CPU:      node.CPU,
			Memory:   node.Memory,
			Disk:     node.Disk,
			PodCount: node.PodCount,
			Status:   node.Status,
		})
	}

	if clusterMetrics != nil {
		data.Cluster.TotalCPU = clusterMetrics.CPUUsage
		data.Cluster.TotalMemory = clusterMetrics.MemoryUsage
		data.Cluster.HealthyNodes = clusterMetrics.HealthyNodes
		data.Cluster.TotalNodes = clusterMetrics.TotalNodes
	}

	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, msg)
}

// ScriptOutput handles WebSocket connections for script execution output
func (h *WebSocketHandlers) ScriptOutput(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return websocket.New(func(conn *websocket.Conn) {
			defer conn.Close()

			executionID := c.Params("id")
			if executionID == "" {
				return
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Get output channel for this execution
			outputCh, err := h.ansibleRunner.GetOutputChannel(executionID)
			if err != nil {
				log.Printf("Failed to get output channel: %v", err)
				return
			}

			// Read messages from client (for ping/pong or close)
			go func() {
				for {
					_, _, err := conn.ReadMessage()
					if err != nil {
						cancel()
						return
					}
				}
			}()

			// Send output to client
			for {
				select {
				case <-ctx.Done():
					return
				case line, ok := <-outputCh:
					if !ok {
						// Channel closed, send final status
						return
					}

					msg, err := json.Marshal(line)
					if err != nil {
						continue
					}

					if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
						return
					}
				}
			}
		})(c)
	}

	return fiber.ErrUpgradeRequired
}
