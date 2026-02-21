package loki

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Source    string `json:"source"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Namespace string `json:"namespace,omitempty"`
}

type QueryParams struct {
	Node      string
	Namespace string
	Query     string
	Limit     int
	Start     time.Time
	End       time.Time
}

// Loki API response structures
type QueryResponse struct {
	Status string     `json:"status"`
	Data   QueryData  `json:"data"`
}

type QueryData struct {
	ResultType string   `json:"resultType"`
	Result     []Stream `json:"result"`
}

type Stream struct {
	Labels  map[string]string `json:"stream"`
	Values  [][]string        `json:"values"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Query(ctx context.Context, params QueryParams) ([]LogEntry, error) {
	// Build LogQL query
	query := c.buildQuery(params)

	// Set time range
	end := params.End
	if end.IsZero() {
		end = time.Now()
	}
	start := params.Start
	if start.IsZero() {
		start = end.Add(-1 * time.Hour)
	}

	limit := params.Limit
	if limit == 0 {
		limit = 100
	}

	// Build URL
	u, err := url.Parse(c.baseURL + "/loki/api/v1/query_range")
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("query", query)
	q.Set("start", fmt.Sprintf("%d", start.UnixNano()))
	q.Set("end", fmt.Sprintf("%d", end.UnixNano()))
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("direction", "backward")
	u.RawQuery = q.Encode()

	// Make request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query Loki: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Loki returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var queryResp QueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&queryResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to LogEntry slice
	entries := c.parseStreams(queryResp.Data.Result)
	return entries, nil
}

func (c *Client) buildQuery(params QueryParams) string {
	// Build label selectors
	selectors := []string{}

	if params.Node != "" {
		selectors = append(selectors, fmt.Sprintf(`node_name="%s"`, params.Node))
	}
	if params.Namespace != "" {
		selectors = append(selectors, fmt.Sprintf(`namespace="%s"`, params.Namespace))
	}

	// Default query if no selectors
	if len(selectors) == 0 {
		selectors = append(selectors, `job=~".+"`)
	}

	query := "{"
	for i, sel := range selectors {
		if i > 0 {
			query += ", "
		}
		query += sel
	}
	query += "}"

	// Add line filter if search query provided
	if params.Query != "" {
		query += fmt.Sprintf(` |~ "%s"`, params.Query)
	}

	return query
}

func (c *Client) parseStreams(streams []Stream) []LogEntry {
	entries := []LogEntry{}

	for _, stream := range streams {
		source := stream.Labels["app"]
		if source == "" {
			source = stream.Labels["container"]
		}
		if source == "" {
			source = stream.Labels["job"]
		}

		namespace := stream.Labels["namespace"]

		for _, value := range stream.Values {
			if len(value) < 2 {
				continue
			}

			// Parse timestamp (nanoseconds since epoch)
			ts := value[0]
			message := value[1]

			// Detect log level from message
			level := detectLogLevel(message)

			// Convert nanosecond timestamp to time
			var timestamp time.Time
			var nsec int64
			fmt.Sscanf(ts, "%d", &nsec)
			timestamp = time.Unix(0, nsec)

			entries = append(entries, LogEntry{
				Timestamp: timestamp.Format(time.RFC3339),
				Source:    source,
				Level:     level,
				Message:   message,
				Namespace: namespace,
			})
		}
	}

	return entries
}

func detectLogLevel(message string) string {
	// Simple level detection based on common patterns
	switch {
	case containsAny(message, []string{"ERROR", "error", "ERR", "FATAL", "fatal", "PANIC", "panic"}):
		return "error"
	case containsAny(message, []string{"WARN", "warn", "WARNING", "warning"}):
		return "warn"
	default:
		return "info"
	}
}

func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

// TailLogs returns a channel that streams log entries
func (c *Client) TailLogs(ctx context.Context, params QueryParams) (<-chan LogEntry, error) {
	ch := make(chan LogEntry, 100)

	go func() {
		defer close(ch)

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		lastTimestamp := time.Now()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				params.Start = lastTimestamp
				params.End = time.Now()
				params.Limit = 50

				entries, err := c.Query(ctx, params)
				if err != nil {
					continue
				}

				for _, entry := range entries {
					select {
					case ch <- entry:
						ts, _ := time.Parse(time.RFC3339, entry.Timestamp)
						if ts.After(lastTimestamp) {
							lastTimestamp = ts
						}
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return ch, nil
}
