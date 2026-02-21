package api

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/timholm/animus-dashboard/internal/ansible"
	"github.com/timholm/animus-dashboard/internal/k8s"
	"github.com/timholm/animus-dashboard/internal/loki"
	"github.com/timholm/animus-dashboard/internal/ssh"
)

type Handlers struct {
	k8sClient     *k8s.Client
	lokiClient    *loki.Client
	sshClient     *ssh.Client
	ansibleRunner *ansible.Runner
}

func NewHandlers(k8s *k8s.Client, loki *loki.Client, ssh *ssh.Client, ansible *ansible.Runner) *Handlers {
	return &Handlers{
		k8sClient:     k8s,
		lokiClient:    loki,
		sshClient:     ssh,
		ansibleRunner: ansible,
	}
}

// GetNodes returns all cluster nodes
func (h *Handlers) GetNodes(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	nodes, err := h.k8sClient.GetNodes(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(nodes)
}

// GetNode returns a specific node by ID
func (h *Handlers) GetNode(c *fiber.Ctx) error {
	nodeID := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	node, err := h.k8sClient.GetNode(ctx, nodeID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(node)
}

// GetNodePods returns pods running on a specific node
func (h *Handlers) GetNodePods(c *fiber.Ctx) error {
	nodeID := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pods, err := h.k8sClient.GetNodePods(ctx, nodeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(pods)
}

// GetClusterMetrics returns cluster-wide metrics
func (h *Handlers) GetClusterMetrics(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	metrics, err := h.k8sClient.GetClusterMetrics(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(metrics)
}

// GetLogs queries Loki for log entries
func (h *Handlers) GetLogs(c *fiber.Ctx) error {
	params := loki.QueryParams{
		Node:      c.Query("node"),
		Namespace: c.Query("namespace"),
		Query:     c.Query("query"),
		Limit:     c.QueryInt("limit", 100),
	}

	// Parse time range
	if start := c.Query("start"); start != "" {
		if t, err := time.Parse(time.RFC3339, start); err == nil {
			params.Start = t
		}
	}
	if end := c.Query("end"); end != "" {
		if t, err := time.Parse(time.RFC3339, end); err == nil {
			params.End = t
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	entries, err := h.lokiClient.Query(ctx, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(entries)
}

// GetScripts returns available Ansible scripts/playbooks
func (h *Handlers) GetScripts(c *fiber.Ctx) error {
	scripts, err := h.ansibleRunner.GetScripts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(scripts)
}

// RunScript executes an Ansible playbook
func (h *Handlers) RunScript(c *fiber.Ctx) error {
	var req struct {
		ScriptID    string   `json:"scriptId"`
		TargetNodes []string `json:"targetNodes"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.ScriptID == "" || len(req.TargetNodes) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "scriptId and targetNodes are required",
		})
	}

	ctx := context.Background()
	execution, err := h.ansibleRunner.RunScript(ctx, req.ScriptID, req.TargetNodes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(execution)
}

// GetScriptExecutions returns all script execution history
func (h *Handlers) GetScriptExecutions(c *fiber.Ctx) error {
	executions := h.ansibleRunner.GetExecutions()
	return c.JSON(executions)
}

// GetScriptExecution returns a specific execution by ID
func (h *Handlers) GetScriptExecution(c *fiber.Ctx) error {
	executionID := c.Params("id")

	execution, err := h.ansibleRunner.GetExecution(executionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(execution)
}

// CheckUpdates checks for available updates on nodes
func (h *Handlers) CheckUpdates(c *fiber.Ctx) error {
	if h.sshClient == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "SSH client not configured",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Get all nodes
	nodes, err := h.k8sClient.GetNodes(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	type NodeUpdate struct {
		Node    string `json:"node"`
		Updates []struct {
			Name      string `json:"name"`
			Current   string `json:"current"`
			Available string `json:"available"`
		} `json:"updates"`
	}

	results := []NodeUpdate{}
	for _, node := range nodes {
		updateCount, err := h.sshClient.GetPackageUpdates(ctx, node.IP)
		if err != nil {
			continue
		}

		nodeUpdate := NodeUpdate{
			Node: node.Name,
		}

		if updateCount > 0 {
			nodeUpdate.Updates = append(nodeUpdate.Updates, struct {
				Name      string `json:"name"`
				Current   string `json:"current"`
				Available string `json:"available"`
			}{
				Name:      "apt packages",
				Current:   "installed",
				Available: fmt.Sprintf("%d updates available", updateCount),
			})
		}

		results = append(results, nodeUpdate)
	}

	return c.JSON(results)
}
