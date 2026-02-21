package ansible

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Runner struct {
	playbooksPath string
	inventoryPath string
	executions    map[string]*Execution
	mu            sync.RWMutex
}

type Script struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	TargetType  string `json:"targetType"` // "single" or "all"
	Filename    string `json:"-"`
}

type Execution struct {
	ID          string    `json:"id"`
	ScriptID    string    `json:"scriptId"`
	Status      string    `json:"status"` // "running", "completed", "failed"
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt time.Time `json:"completedAt,omitempty"`
	Output      []string  `json:"output"`
	TargetNodes []string  `json:"targetNodes"`
	outputCh    chan OutputLine
	cancel      context.CancelFunc
}

type OutputLine struct {
	Type      string    `json:"type"` // "stdout", "stderr", "status"
	Line      string    `json:"line"`
	Timestamp time.Time `json:"timestamp"`
}

func NewRunner(playbooksPath, inventoryPath string) *Runner {
	return &Runner{
		playbooksPath: playbooksPath,
		inventoryPath: inventoryPath,
		executions:    make(map[string]*Execution),
	}
}

func (r *Runner) GetScripts() ([]Script, error) {
	scripts := []Script{}

	// Read playbook files from directory
	entries, err := os.ReadDir(r.playbooksPath)
	if err != nil {
		// Return default scripts if directory doesn't exist
		return r.getDefaultScripts(), nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if filepath.Ext(name) != ".yml" && filepath.Ext(name) != ".yaml" {
			continue
		}

		script := r.parsePlaybookMeta(name)
		scripts = append(scripts, script)
	}

	if len(scripts) == 0 {
		return r.getDefaultScripts(), nil
	}

	return scripts, nil
}

func (r *Runner) getDefaultScripts() []Script {
	return []Script{
		{
			ID:          "update-packages",
			Name:        "Update Packages",
			Description: "Run apt update and upgrade on target nodes",
			Category:    "system",
			TargetType:  "single",
			Filename:    "update-packages.yml",
		},
		{
			ID:          "reboot-node",
			Name:        "Reboot Node",
			Description: "Safely reboot the target node",
			Category:    "system",
			TargetType:  "single",
			Filename:    "reboot-node.yml",
		},
		{
			ID:          "restart-k3s",
			Name:        "Restart K3s Agent",
			Description: "Restart the K3s agent service",
			Category:    "kubernetes",
			TargetType:  "single",
			Filename:    "restart-k3s.yml",
		},
		{
			ID:          "sync-time",
			Name:        "Sync Time (NTP)",
			Description: "Synchronize system time across nodes",
			Category:    "system",
			TargetType:  "all",
			Filename:    "sync-time.yml",
		},
		{
			ID:          "clear-logs",
			Name:        "Clear Old Logs",
			Description: "Remove old log files to free disk space",
			Category:    "monitoring",
			TargetType:  "single",
			Filename:    "clear-logs.yml",
		},
		{
			ID:          "drain-node",
			Name:        "Drain Node",
			Description: "Drain node for maintenance",
			Category:    "kubernetes",
			TargetType:  "single",
			Filename:    "drain-node.yml",
		},
		{
			ID:          "uncordon-node",
			Name:        "Uncordon Node",
			Description: "Mark node as schedulable again",
			Category:    "kubernetes",
			TargetType:  "single",
			Filename:    "uncordon-node.yml",
		},
	}
}

func (r *Runner) parsePlaybookMeta(filename string) Script {
	// Extract script info from filename
	baseName := filename[:len(filename)-len(filepath.Ext(filename))]

	return Script{
		ID:          baseName,
		Name:        formatName(baseName),
		Description: fmt.Sprintf("Run %s playbook", baseName),
		Category:    "custom",
		TargetType:  "single",
		Filename:    filename,
	}
}

func formatName(s string) string {
	// Convert kebab-case to Title Case
	result := ""
	nextUpper := true
	for _, c := range s {
		if c == '-' || c == '_' {
			result += " "
			nextUpper = true
		} else if nextUpper {
			if c >= 'a' && c <= 'z' {
				result += string(c - 32)
			} else {
				result += string(c)
			}
			nextUpper = false
		} else {
			result += string(c)
		}
	}
	return result
}

func (r *Runner) RunScript(ctx context.Context, scriptID string, targetNodes []string) (*Execution, error) {
	// Find the script
	scripts, _ := r.GetScripts()
	var script *Script
	for _, s := range scripts {
		if s.ID == scriptID {
			script = &s
			break
		}
	}
	if script == nil {
		return nil, fmt.Errorf("script not found: %s", scriptID)
	}

	// Create execution record
	execCtx, cancel := context.WithCancel(ctx)
	execution := &Execution{
		ID:          uuid.New().String(),
		ScriptID:    scriptID,
		Status:      "running",
		StartedAt:   time.Now(),
		TargetNodes: targetNodes,
		Output:      []string{},
		outputCh:    make(chan OutputLine, 100),
		cancel:      cancel,
	}

	r.mu.Lock()
	r.executions[execution.ID] = execution
	r.mu.Unlock()

	// Run in background
	go r.executePlaybook(execCtx, execution, script, targetNodes)

	return execution, nil
}

func (r *Runner) executePlaybook(ctx context.Context, execution *Execution, script *Script, targets []string) {
	defer close(execution.outputCh)

	playbookPath := filepath.Join(r.playbooksPath, script.Filename)

	// Build ansible-playbook command
	args := []string{
		"-i", r.inventoryPath,
	}

	// Add target hosts
	if len(targets) == 1 && targets[0] != "all" {
		args = append(args, "-l", targets[0])
	}

	args = append(args, playbookPath)

	cmd := exec.CommandContext(ctx, "ansible-playbook", args...)

	// Set up output pipes
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		execution.Status = "failed"
		execution.Output = append(execution.Output, fmt.Sprintf("Error: %v", err))
		execution.CompletedAt = time.Now()
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		execution.Status = "failed"
		execution.Output = append(execution.Output, fmt.Sprintf("Error: %v", err))
		execution.CompletedAt = time.Now()
		return
	}

	// Start command
	if err := cmd.Start(); err != nil {
		execution.Status = "failed"
		execution.Output = append(execution.Output, fmt.Sprintf("Error: %v", err))
		execution.CompletedAt = time.Now()
		return
	}

	// Read output
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			execution.Output = append(execution.Output, line)
			execution.outputCh <- OutputLine{
				Type:      "stdout",
				Line:      line,
				Timestamp: time.Now(),
			}
		}
	}()

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			execution.Output = append(execution.Output, line)
			execution.outputCh <- OutputLine{
				Type:      "stderr",
				Line:      line,
				Timestamp: time.Now(),
			}
		}
	}()

	wg.Wait()

	// Wait for command to complete
	err = cmd.Wait()
	execution.CompletedAt = time.Now()

	if err != nil {
		execution.Status = "failed"
		execution.outputCh <- OutputLine{
			Type:      "status",
			Line:      fmt.Sprintf("Execution failed: %v", err),
			Timestamp: time.Now(),
		}
	} else {
		execution.Status = "completed"
		execution.outputCh <- OutputLine{
			Type:      "status",
			Line:      "Execution completed successfully",
			Timestamp: time.Now(),
		}
	}
}

func (r *Runner) GetExecution(id string) (*Execution, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	execution, ok := r.executions[id]
	if !ok {
		return nil, fmt.Errorf("execution not found: %s", id)
	}

	return execution, nil
}

func (r *Runner) GetExecutions() []*Execution {
	r.mu.RLock()
	defer r.mu.RUnlock()

	executions := make([]*Execution, 0, len(r.executions))
	for _, e := range r.executions {
		executions = append(executions, e)
	}

	return executions
}

func (r *Runner) GetOutputChannel(id string) (<-chan OutputLine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	execution, ok := r.executions[id]
	if !ok {
		return nil, fmt.Errorf("execution not found: %s", id)
	}

	return execution.outputCh, nil
}

func (r *Runner) StopExecution(id string) error {
	r.mu.RLock()
	execution, ok := r.executions[id]
	r.mu.RUnlock()

	if !ok {
		return fmt.Errorf("execution not found: %s", id)
	}

	if execution.cancel != nil {
		execution.cancel()
	}

	return nil
}
