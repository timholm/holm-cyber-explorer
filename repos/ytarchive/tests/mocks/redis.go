package mocks

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

// MockRedisClient provides an in-memory mock Redis implementation for testing
type MockRedisClient struct {
	mu          sync.RWMutex
	data        map[string]string
	sets        map[string]map[string]struct{}
	lists       map[string][]string
	expirations map[string]time.Time
	closed      bool
	shouldFail  bool
	failMessage string
}

// NewMockRedisClient creates a new mock Redis client
func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data:        make(map[string]string),
		sets:        make(map[string]map[string]struct{}),
		lists:       make(map[string][]string),
		expirations: make(map[string]time.Time),
	}
}

// Set stores a key-value pair
func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if m.shouldFail {
		return &MockRedisError{Message: m.failMessage}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	var strValue string
	switch v := value.(type) {
	case string:
		strValue = v
	case []byte:
		strValue = string(v)
	default:
		bytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		strValue = string(bytes)
	}

	m.data[key] = strValue

	if expiration > 0 {
		m.expirations[key] = time.Now().Add(expiration)
	}

	return nil
}

// Get retrieves a value by key
func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	if m.shouldFail {
		return "", &MockRedisError{Message: m.failMessage}
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	// Check expiration
	if exp, ok := m.expirations[key]; ok {
		if time.Now().After(exp) {
			delete(m.data, key)
			delete(m.expirations, key)
			return "", ErrNil
		}
	}

	value, ok := m.data[key]
	if !ok {
		return "", ErrNil
	}

	return value, nil
}

// Del deletes one or more keys
func (m *MockRedisClient) Del(ctx context.Context, keys ...string) (int64, error) {
	if m.shouldFail {
		return 0, &MockRedisError{Message: m.failMessage}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	var deleted int64
	for _, key := range keys {
		if _, ok := m.data[key]; ok {
			delete(m.data, key)
			delete(m.expirations, key)
			deleted++
		}
		if _, ok := m.sets[key]; ok {
			delete(m.sets, key)
			deleted++
		}
		if _, ok := m.lists[key]; ok {
			delete(m.lists, key)
			deleted++
		}
	}

	return deleted, nil
}

// Exists checks if keys exist
func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	if m.shouldFail {
		return 0, &MockRedisError{Message: m.failMessage}
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var count int64
	for _, key := range keys {
		// Check expiration
		if exp, ok := m.expirations[key]; ok {
			if time.Now().After(exp) {
				continue
			}
		}

		if _, ok := m.data[key]; ok {
			count++
		} else if _, ok := m.sets[key]; ok {
			count++
		} else if _, ok := m.lists[key]; ok {
			count++
		}
	}

	return count, nil
}

// Keys returns all keys matching a pattern (simplified glob matching)
func (m *MockRedisClient) Keys(ctx context.Context, pattern string) ([]string, error) {
	if m.shouldFail {
		return nil, &MockRedisError{Message: m.failMessage}
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var keys []string

	// Collect keys from all storage types
	allKeys := make(map[string]struct{})
	for k := range m.data {
		allKeys[k] = struct{}{}
	}
	for k := range m.sets {
		allKeys[k] = struct{}{}
	}
	for k := range m.lists {
		allKeys[k] = struct{}{}
	}

	for key := range allKeys {
		if matchPattern(pattern, key) {
			// Check expiration
			if exp, ok := m.expirations[key]; ok {
				if time.Now().After(exp) {
					continue
				}
			}
			keys = append(keys, key)
		}
	}

	return keys, nil
}

// SAdd adds members to a set
func (m *MockRedisClient) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	if m.shouldFail {
		return 0, &MockRedisError{Message: m.failMessage}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.sets[key] == nil {
		m.sets[key] = make(map[string]struct{})
	}

	var added int64
	for _, member := range members {
		str := toString(member)
		if _, exists := m.sets[key][str]; !exists {
			m.sets[key][str] = struct{}{}
			added++
		}
	}

	return added, nil
}

// SMembers returns all members of a set
func (m *MockRedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	if m.shouldFail {
		return nil, &MockRedisError{Message: m.failMessage}
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	set, ok := m.sets[key]
	if !ok {
		return []string{}, nil
	}

	members := make([]string, 0, len(set))
	for member := range set {
		members = append(members, member)
	}

	return members, nil
}

// SRem removes members from a set
func (m *MockRedisClient) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	if m.shouldFail {
		return 0, &MockRedisError{Message: m.failMessage}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	set, ok := m.sets[key]
	if !ok {
		return 0, nil
	}

	var removed int64
	for _, member := range members {
		str := toString(member)
		if _, exists := set[str]; exists {
			delete(set, str)
			removed++
		}
	}

	return removed, nil
}

// SCard returns the number of members in a set
func (m *MockRedisClient) SCard(ctx context.Context, key string) (int64, error) {
	if m.shouldFail {
		return 0, &MockRedisError{Message: m.failMessage}
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	set, ok := m.sets[key]
	if !ok {
		return 0, nil
	}

	return int64(len(set)), nil
}

// SIsMember checks if a member is in a set
func (m *MockRedisClient) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	if m.shouldFail {
		return false, &MockRedisError{Message: m.failMessage}
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	set, ok := m.sets[key]
	if !ok {
		return false, nil
	}

	_, exists := set[toString(member)]
	return exists, nil
}

// LPush pushes values to the left of a list
func (m *MockRedisClient) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	if m.shouldFail {
		return 0, &MockRedisError{Message: m.failMessage}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, v := range values {
		m.lists[key] = append([]string{toString(v)}, m.lists[key]...)
	}

	return int64(len(m.lists[key])), nil
}

// RPush pushes values to the right of a list
func (m *MockRedisClient) RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	if m.shouldFail {
		return 0, &MockRedisError{Message: m.failMessage}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, v := range values {
		m.lists[key] = append(m.lists[key], toString(v))
	}

	return int64(len(m.lists[key])), nil
}

// LPop removes and returns the first element from a list
func (m *MockRedisClient) LPop(ctx context.Context, key string) (string, error) {
	if m.shouldFail {
		return "", &MockRedisError{Message: m.failMessage}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	list, ok := m.lists[key]
	if !ok || len(list) == 0 {
		return "", ErrNil
	}

	value := list[0]
	m.lists[key] = list[1:]

	return value, nil
}

// RPop removes and returns the last element from a list
func (m *MockRedisClient) RPop(ctx context.Context, key string) (string, error) {
	if m.shouldFail {
		return "", &MockRedisError{Message: m.failMessage}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	list, ok := m.lists[key]
	if !ok || len(list) == 0 {
		return "", ErrNil
	}

	value := list[len(list)-1]
	m.lists[key] = list[:len(list)-1]

	return value, nil
}

// LLen returns the length of a list
func (m *MockRedisClient) LLen(ctx context.Context, key string) (int64, error) {
	if m.shouldFail {
		return 0, &MockRedisError{Message: m.failMessage}
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	list, ok := m.lists[key]
	if !ok {
		return 0, nil
	}

	return int64(len(list)), nil
}

// LRange returns elements from a list
func (m *MockRedisClient) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if m.shouldFail {
		return nil, &MockRedisError{Message: m.failMessage}
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	list, ok := m.lists[key]
	if !ok {
		return []string{}, nil
	}

	length := int64(len(list))
	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop
	}
	if start < 0 {
		start = 0
	}
	if stop >= length {
		stop = length - 1
	}
	if start > stop || start >= length {
		return []string{}, nil
	}

	return list[start : stop+1], nil
}

// Ping checks connection (always succeeds for mock unless shouldFail is set)
func (m *MockRedisClient) Ping(ctx context.Context) error {
	if m.shouldFail {
		return &MockRedisError{Message: m.failMessage}
	}
	if m.closed {
		return &MockRedisError{Message: "connection closed"}
	}
	return nil
}

// Close closes the mock connection
func (m *MockRedisClient) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	return nil
}

// FlushAll clears all data
func (m *MockRedisClient) FlushAll(ctx context.Context) error {
	if m.shouldFail {
		return &MockRedisError{Message: m.failMessage}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[string]string)
	m.sets = make(map[string]map[string]struct{})
	m.lists = make(map[string][]string)
	m.expirations = make(map[string]time.Time)

	return nil
}

// SetFailure configures the mock to return errors
func (m *MockRedisClient) SetFailure(shouldFail bool, message string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail = shouldFail
	m.failMessage = message
}

// Pipeline returns a mock pipeline
func (m *MockRedisClient) Pipeline() *MockRedisPipeline {
	return &MockRedisPipeline{
		client: m,
		ops:    make([]pipelineOp, 0),
	}
}

// MockRedisPipeline represents a mock Redis pipeline
type MockRedisPipeline struct {
	client *MockRedisClient
	ops    []pipelineOp
}

type pipelineOp struct {
	opType string
	key    string
	value  interface{}
	values []interface{}
	exp    time.Duration
}

// Set queues a SET operation
func (p *MockRedisPipeline) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	p.ops = append(p.ops, pipelineOp{
		opType: "set",
		key:    key,
		value:  value,
		exp:    expiration,
	})
}

// Get queues a GET operation
func (p *MockRedisPipeline) Get(ctx context.Context, key string) {
	p.ops = append(p.ops, pipelineOp{
		opType: "get",
		key:    key,
	})
}

// Del queues a DEL operation
func (p *MockRedisPipeline) Del(ctx context.Context, keys ...string) {
	for _, key := range keys {
		p.ops = append(p.ops, pipelineOp{
			opType: "del",
			key:    key,
		})
	}
}

// SAdd queues an SADD operation
func (p *MockRedisPipeline) SAdd(ctx context.Context, key string, members ...interface{}) {
	p.ops = append(p.ops, pipelineOp{
		opType: "sadd",
		key:    key,
		values: members,
	})
}

// SRem queues an SREM operation
func (p *MockRedisPipeline) SRem(ctx context.Context, key string, members ...interface{}) {
	p.ops = append(p.ops, pipelineOp{
		opType: "srem",
		key:    key,
		values: members,
	})
}

// LPush queues an LPUSH operation
func (p *MockRedisPipeline) LPush(ctx context.Context, key string, values ...interface{}) {
	p.ops = append(p.ops, pipelineOp{
		opType: "lpush",
		key:    key,
		values: values,
	})
}

// RPush queues an RPUSH operation
func (p *MockRedisPipeline) RPush(ctx context.Context, key string, values ...interface{}) {
	p.ops = append(p.ops, pipelineOp{
		opType: "rpush",
		key:    key,
		values: values,
	})
}

// Exec executes all queued operations
func (p *MockRedisPipeline) Exec(ctx context.Context) ([]interface{}, error) {
	results := make([]interface{}, 0, len(p.ops))

	for _, op := range p.ops {
		var result interface{}
		var err error

		switch op.opType {
		case "set":
			err = p.client.Set(ctx, op.key, op.value, op.exp)
			result = "OK"
		case "get":
			result, err = p.client.Get(ctx, op.key)
		case "del":
			result, err = p.client.Del(ctx, op.key)
		case "sadd":
			result, err = p.client.SAdd(ctx, op.key, op.values...)
		case "srem":
			result, err = p.client.SRem(ctx, op.key, op.values...)
		case "lpush":
			result, err = p.client.LPush(ctx, op.key, op.values...)
		case "rpush":
			result, err = p.client.RPush(ctx, op.key, op.values...)
		}

		if err != nil {
			return results, err
		}
		results = append(results, result)
	}

	return results, nil
}

// MockRedisError represents a mock Redis error
type MockRedisError struct {
	Message string
}

func (e *MockRedisError) Error() string {
	return e.Message
}

// ErrNil is returned when a key doesn't exist
var ErrNil = &MockRedisError{Message: "redis: nil"}

// IsNilError checks if the error is a nil error (key not found)
func IsNilError(err error) bool {
	if err == nil {
		return false
	}
	return err == ErrNil || err.Error() == "redis: nil"
}

// Helper function to convert interface{} to string
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	default:
		bytes, _ := json.Marshal(val)
		return string(bytes)
	}
}

// matchPattern provides simple glob pattern matching for Keys()
// Supports * as wildcard for any characters
func matchPattern(pattern, str string) bool {
	if pattern == "*" {
		return true
	}

	// Simple implementation for common patterns like "prefix:*"
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(str) >= len(prefix) && str[:len(prefix)] == prefix
	}

	// Patterns like "*:suffix"
	if len(pattern) > 0 && pattern[0] == '*' {
		suffix := pattern[1:]
		return len(str) >= len(suffix) && str[len(str)-len(suffix):] == suffix
	}

	// Patterns like "prefix:*:suffix"
	for i := 0; i < len(pattern); i++ {
		if pattern[i] == '*' {
			prefix := pattern[:i]
			suffix := pattern[i+1:]
			if len(str) < len(prefix)+len(suffix) {
				return false
			}
			return str[:len(prefix)] == prefix && str[len(str)-len(suffix):] == suffix
		}
	}

	return pattern == str
}

// MockRedisQueue provides queue-specific operations for the download queue
type MockRedisQueue struct {
	client *MockRedisClient
}

// NewMockRedisQueue creates a new mock Redis queue
func NewMockRedisQueue(client *MockRedisClient) *MockRedisQueue {
	return &MockRedisQueue{client: client}
}

// PushVideos pushes multiple video IDs to the download queue
func (q *MockRedisQueue) PushVideos(ctx context.Context, queueKey string, videoIDs []string) error {
	for _, id := range videoIDs {
		_, err := q.client.RPush(ctx, queueKey, id)
		if err != nil {
			return err
		}
	}
	return nil
}

// ClaimVideo claims a video from the queue for processing
func (q *MockRedisQueue) ClaimVideo(ctx context.Context, queueKey, processingKey string) (string, error) {
	videoID, err := q.client.LPop(ctx, queueKey)
	if err != nil {
		return "", err
	}

	// Add to processing set
	_, err = q.client.SAdd(ctx, processingKey, videoID)
	if err != nil {
		// Try to put it back
		q.client.LPush(ctx, queueKey, videoID)
		return "", err
	}

	return videoID, nil
}

// CompleteVideo marks a video as completed
func (q *MockRedisQueue) CompleteVideo(ctx context.Context, processingKey, completedKey, videoID string) error {
	// Remove from processing
	_, err := q.client.SRem(ctx, processingKey, videoID)
	if err != nil {
		return err
	}

	// Add to completed
	_, err = q.client.SAdd(ctx, completedKey, videoID)
	return err
}

// QueueLength returns the number of items in the queue
func (q *MockRedisQueue) QueueLength(ctx context.Context, queueKey string) (int64, error) {
	return q.client.LLen(ctx, queueKey)
}
