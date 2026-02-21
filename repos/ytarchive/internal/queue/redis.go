// Package queue provides Redis queue operations for the YouTube Channel Archiver.
package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// QueueKeyPrefix is the prefix for queue keys in Redis
	QueueKeyPrefix = "ytarchive:queue:"
	// DefaultBatchSize is the default number of video IDs to claim in a batch
	DefaultBatchSize = 10
	// DefaultTimeout is the default timeout for Redis operations
	DefaultTimeout = 5 * time.Second
)

// Queue represents a Redis-backed queue for video downloads.
type Queue struct {
	client *redis.Client
	ctx    context.Context
}

// NewQueue creates a new Queue connected to the specified Redis URL.
// The redisURL should be in the format: redis://[:password@]host:port[/db]
func NewQueue(redisURL string) (*Queue, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)
	ctx := context.Background()

	// Test the connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Queue{
		client: client,
		ctx:    ctx,
	}, nil
}

// NewQueueWithClient creates a new Queue with an existing Redis client.
// This is useful for testing or when you need to share a client.
func NewQueueWithClient(client *redis.Client) *Queue {
	return &Queue{
		client: client,
		ctx:    context.Background(),
	}
}

// Close closes the Redis connection.
func (q *Queue) Close() error {
	return q.client.Close()
}

// getQueueKey returns the Redis key for a channel's queue.
func getQueueKey(channelID string) string {
	return QueueKeyPrefix + channelID
}

// PushVideoBatch pushes a batch of video IDs to the channel's queue.
// Videos are pushed to the left side of the list (LPUSH) so they can be
// popped from the right side (RPOP) in FIFO order.
func (q *Queue) PushVideoBatch(channelID string, videoIDs []string) error {
	if channelID == "" {
		return fmt.Errorf("channel ID cannot be empty")
	}
	if len(videoIDs) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(q.ctx, DefaultTimeout)
	defer cancel()

	key := getQueueKey(channelID)

	// Convert []string to []interface{} for Redis
	values := make([]interface{}, len(videoIDs))
	for i, id := range videoIDs {
		values[i] = id
	}

	if err := q.client.LPush(ctx, key, values...).Err(); err != nil {
		return fmt.Errorf("failed to push video batch: %w", err)
	}

	return nil
}

// ClaimBatch atomically claims a batch of video IDs from the channel's queue.
// Returns up to DefaultBatchSize video IDs using RPOP (FIFO order).
func (q *Queue) ClaimBatch(channelID string) ([]string, error) {
	return q.ClaimBatchN(channelID, DefaultBatchSize)
}

// ClaimBatchN atomically claims up to n video IDs from the channel's queue.
func (q *Queue) ClaimBatchN(channelID string, n int) ([]string, error) {
	if channelID == "" {
		return nil, fmt.Errorf("channel ID cannot be empty")
	}
	if n <= 0 {
		return nil, fmt.Errorf("batch size must be positive")
	}

	ctx, cancel := context.WithTimeout(q.ctx, DefaultTimeout)
	defer cancel()

	key := getQueueKey(channelID)

	// Use a Lua script to atomically pop multiple elements
	script := redis.NewScript(`
		local result = {}
		for i = 1, ARGV[1] do
			local item = redis.call('RPOP', KEYS[1])
			if not item then
				break
			end
			table.insert(result, item)
		end
		return result
	`)

	result, err := script.Run(ctx, q.client, []string{key}, n).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to claim batch: %w", err)
	}

	// Convert result to []string
	items, ok := result.([]interface{})
	if !ok {
		return nil, nil
	}

	videoIDs := make([]string, 0, len(items))
	for _, item := range items {
		if str, ok := item.(string); ok {
			videoIDs = append(videoIDs, str)
		}
	}

	return videoIDs, nil
}

// GetQueueLength returns the number of video IDs in the channel's queue.
func (q *Queue) GetQueueLength(channelID string) (int64, error) {
	if channelID == "" {
		return 0, fmt.Errorf("channel ID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(q.ctx, DefaultTimeout)
	defer cancel()

	key := getQueueKey(channelID)
	length, err := q.client.LLen(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get queue length: %w", err)
	}

	return length, nil
}

// ClearQueue removes all video IDs from the channel's queue.
func (q *Queue) ClearQueue(channelID string) error {
	if channelID == "" {
		return fmt.Errorf("channel ID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(q.ctx, DefaultTimeout)
	defer cancel()

	key := getQueueKey(channelID)
	if err := q.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to clear queue: %w", err)
	}

	return nil
}

// PeekQueue returns up to n video IDs from the queue without removing them.
// Useful for inspection and debugging.
func (q *Queue) PeekQueue(channelID string, n int64) ([]string, error) {
	if channelID == "" {
		return nil, fmt.Errorf("channel ID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(q.ctx, DefaultTimeout)
	defer cancel()

	key := getQueueKey(channelID)
	// Use negative indices to get from the right (FIFO order)
	result, err := q.client.LRange(ctx, key, -n, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to peek queue: %w", err)
	}

	// Reverse to match FIFO order
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}

// GetAllQueueKeys returns all queue keys in Redis.
// Useful for listing all active channel queues.
func (q *Queue) GetAllQueueKeys() ([]string, error) {
	ctx, cancel := context.WithTimeout(q.ctx, DefaultTimeout)
	defer cancel()

	pattern := QueueKeyPrefix + "*"
	keys, err := q.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get queue keys: %w", err)
	}

	return keys, nil
}

// Client returns the underlying Redis client for advanced operations.
func (q *Queue) Client() *redis.Client {
	return q.client
}

// Context returns the queue's context.
func (q *Queue) Context() context.Context {
	return q.ctx
}
