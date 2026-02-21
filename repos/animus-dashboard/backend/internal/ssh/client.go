package ssh

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

type Client struct {
	privateKey []byte
	user       string
}

type CommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func NewClient(keyPath, user string) (*Client, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH key: %w", err)
	}

	return &Client{
		privateKey: key,
		user:       user,
	}, nil
}

func (c *Client) connect(host string) (*ssh.Client, error) {
	signer, err := ssh.ParsePrivateKey(c.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	config := &ssh.ClientConfig{
		User: c.user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // In production, use known hosts
		Timeout:         10 * time.Second,
	}

	// Try common SSH ports
	client, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", host, err)
	}

	return client, nil
}

func (c *Client) RunCommand(ctx context.Context, host, command string) (*CommandResult, error) {
	client, err := c.connect(host)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	// Run command with context timeout
	done := make(chan error, 1)
	go func() {
		done <- session.Run(command)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-done:
		result := &CommandResult{
			Stdout:   stdout.String(),
			Stderr:   stderr.String(),
			ExitCode: 0,
		}
		if err != nil {
			if exitErr, ok := err.(*ssh.ExitError); ok {
				result.ExitCode = exitErr.ExitStatus()
			} else {
				return nil, err
			}
		}
		return result, nil
	}
}

// GetK3sVersion retrieves the K3s version from a node
func (c *Client) GetK3sVersion(ctx context.Context, host string) (string, error) {
	result, err := c.RunCommand(ctx, host, "k3s --version 2>/dev/null | head -1")
	if err != nil {
		return "", err
	}
	return result.Stdout, nil
}

// GetPackageUpdates checks for available package updates
func (c *Client) GetPackageUpdates(ctx context.Context, host string) (int, error) {
	result, err := c.RunCommand(ctx, host, "apt list --upgradable 2>/dev/null | tail -n +2 | wc -l")
	if err != nil {
		return 0, err
	}

	var count int
	fmt.Sscanf(result.Stdout, "%d", &count)
	return count, nil
}

// GetDiskUsage retrieves disk usage percentage for root partition
func (c *Client) GetDiskUsage(ctx context.Context, host string) (int, error) {
	result, err := c.RunCommand(ctx, host, "df / | tail -1 | awk '{print $5}' | tr -d '%'")
	if err != nil {
		return 0, err
	}

	var usage int
	fmt.Sscanf(result.Stdout, "%d", &usage)
	return usage, nil
}

// GetUptime retrieves system uptime
func (c *Client) GetUptime(ctx context.Context, host string) (string, error) {
	result, err := c.RunCommand(ctx, host, "uptime -p")
	if err != nil {
		return "", err
	}
	return result.Stdout, nil
}
