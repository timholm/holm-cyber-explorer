// Package queue provides Redis queue operations for the YouTube Channel Archiver.
package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// JobKeyPrefix is the prefix for job status keys in Redis
	JobKeyPrefix = "ytarchive:job:"
	// JobsSetKey is the key for the set of all active job keys
	JobsSetKey = "ytarchive:active_jobs"
	// ChannelJobsKeyPrefix is the prefix for channel-specific job sets
	ChannelJobsKeyPrefix = "ytarchive:channel_jobs:"
	// JobTTL is the time-to-live for job status entries
	JobTTL = 24 * time.Hour
)

// JobStatus represents the status of a download job.
type JobStatus struct {
	ChannelID  string    `json:"channel_id"`
	WorkerID   string    `json:"worker_id"`
	VideoID    string    `json:"video_id"`
	Progress   float64   `json:"progress"`
	BytesTotal int64     `json:"bytes_total"`
	BytesDone  int64     `json:"bytes_done"`
	Status     string    `json:"status"`
	StartedAt  time.Time `json:"started_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// JobStatusType represents the possible status values for a job.
type JobStatusType string

const (
	// JobStatusPending indicates the job is waiting to start
	JobStatusPending JobStatusType = "pending"
	// JobStatusDownloading indicates the job is actively downloading
	JobStatusDownloading JobStatusType = "downloading"
	// JobStatusProcessing indicates the job is processing the download
	JobStatusProcessing JobStatusType = "processing"
	// JobStatusCompleted indicates the job completed successfully
	JobStatusCompleted JobStatusType = "completed"
	// JobStatusFailed indicates the job failed
	JobStatusFailed JobStatusType = "failed"
)

// getJobKey returns the Redis key for a job.
func getJobKey(channelID, videoID string) string {
	return JobKeyPrefix + channelID + ":" + videoID
}

// getChannelJobsKey returns the Redis key for a channel's job set.
func getChannelJobsKey(channelID string) string {
	return ChannelJobsKeyPrefix + channelID
}

// UpdateJobProgress updates the progress of a job in Redis.
func (q *Queue) UpdateJobProgress(status *JobStatus) error {
	if status == nil {
		return fmt.Errorf("status is nil")
	}
	if status.ChannelID == "" {
		return fmt.Errorf("channel ID cannot be empty")
	}
	if status.VideoID == "" {
		return fmt.Errorf("video ID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(q.ctx, DefaultTimeout)
	defer cancel()

	// Update the timestamp
	status.UpdatedAt = time.Now()

	// Serialize the status
	data, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal job status: %w", err)
	}

	jobKey := getJobKey(status.ChannelID, status.VideoID)
	channelJobsKey := getChannelJobsKey(status.ChannelID)

	// Use a pipeline for atomic operations
	pipe := q.client.Pipeline()

	// Set the job status with TTL
	pipe.Set(ctx, jobKey, data, JobTTL)

	// Add to the global active jobs set
	pipe.SAdd(ctx, JobsSetKey, jobKey)

	// Add to the channel-specific jobs set
	pipe.SAdd(ctx, channelJobsKey, jobKey)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update job progress: %w", err)
	}

	return nil
}

// GetJobStatus retrieves the status of a specific job.
func (q *Queue) GetJobStatus(channelID, videoID string) (*JobStatus, error) {
	if channelID == "" {
		return nil, fmt.Errorf("channel ID cannot be empty")
	}
	if videoID == "" {
		return nil, fmt.Errorf("video ID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(q.ctx, DefaultTimeout)
	defer cancel()

	jobKey := getJobKey(channelID, videoID)
	data, err := q.client.Get(ctx, jobKey).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get job status: %w", err)
	}

	var status JobStatus
	if err := json.Unmarshal([]byte(data), &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job status: %w", err)
	}

	return &status, nil
}

// GetActiveJobs returns all active jobs for a specific channel.
func (q *Queue) GetActiveJobs(channelID string) ([]JobStatus, error) {
	if channelID == "" {
		return nil, fmt.Errorf("channel ID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(q.ctx, DefaultTimeout)
	defer cancel()

	channelJobsKey := getChannelJobsKey(channelID)

	// Get all job keys for this channel
	jobKeys, err := q.client.SMembers(ctx, channelJobsKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get channel job keys: %w", err)
	}

	if len(jobKeys) == 0 {
		return []JobStatus{}, nil
	}

	return q.getJobsFromKeys(ctx, jobKeys, channelJobsKey)
}

// GetAllActiveJobs returns all active jobs across all channels.
func (q *Queue) GetAllActiveJobs() ([]JobStatus, error) {
	ctx, cancel := context.WithTimeout(q.ctx, DefaultTimeout)
	defer cancel()

	// Get all job keys
	jobKeys, err := q.client.SMembers(ctx, JobsSetKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get all job keys: %w", err)
	}

	if len(jobKeys) == 0 {
		return []JobStatus{}, nil
	}

	return q.getJobsFromKeys(ctx, jobKeys, JobsSetKey)
}

// getJobsFromKeys retrieves job statuses from a list of keys.
func (q *Queue) getJobsFromKeys(ctx context.Context, jobKeys []string, setKey string) ([]JobStatus, error) {
	// Get all job data using MGET
	data, err := q.client.MGet(ctx, jobKeys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get job data: %w", err)
	}

	jobs := make([]JobStatus, 0, len(data))
	expiredKeys := make([]string, 0)

	for i, item := range data {
		if item == nil {
			// Job has expired, mark for cleanup
			expiredKeys = append(expiredKeys, jobKeys[i])
			continue
		}

		str, ok := item.(string)
		if !ok {
			continue
		}

		var status JobStatus
		if err := json.Unmarshal([]byte(str), &status); err != nil {
			continue
		}

		// Filter out completed or failed jobs that are stale
		if (status.Status == string(JobStatusCompleted) || status.Status == string(JobStatusFailed)) &&
			time.Since(status.UpdatedAt) > 5*time.Minute {
			expiredKeys = append(expiredKeys, jobKeys[i])
			continue
		}

		jobs = append(jobs, status)
	}

	// Clean up expired keys
	if len(expiredKeys) > 0 {
		go q.cleanupExpiredJobs(expiredKeys, setKey)
	}

	return jobs, nil
}

// cleanupExpiredJobs removes expired job keys from the tracking sets.
func (q *Queue) cleanupExpiredJobs(keys []string, setKey string) {
	ctx, cancel := context.WithTimeout(q.ctx, DefaultTimeout)
	defer cancel()

	pipe := q.client.Pipeline()

	for _, key := range keys {
		pipe.SRem(ctx, setKey, key)
		pipe.SRem(ctx, JobsSetKey, key)
	}

	pipe.Exec(ctx) // Ignore errors for cleanup
}

// RemoveJob removes a job from tracking.
func (q *Queue) RemoveJob(channelID, videoID string) error {
	if channelID == "" {
		return fmt.Errorf("channel ID cannot be empty")
	}
	if videoID == "" {
		return fmt.Errorf("video ID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(q.ctx, DefaultTimeout)
	defer cancel()

	jobKey := getJobKey(channelID, videoID)
	channelJobsKey := getChannelJobsKey(channelID)

	pipe := q.client.Pipeline()
	pipe.Del(ctx, jobKey)
	pipe.SRem(ctx, JobsSetKey, jobKey)
	pipe.SRem(ctx, channelJobsKey, jobKey)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to remove job: %w", err)
	}

	return nil
}

// ClearChannelJobs removes all jobs for a channel.
func (q *Queue) ClearChannelJobs(channelID string) error {
	if channelID == "" {
		return fmt.Errorf("channel ID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(q.ctx, DefaultTimeout)
	defer cancel()

	channelJobsKey := getChannelJobsKey(channelID)

	// Get all job keys for this channel
	jobKeys, err := q.client.SMembers(ctx, channelJobsKey).Result()
	if err != nil {
		return fmt.Errorf("failed to get channel job keys: %w", err)
	}

	if len(jobKeys) == 0 {
		return nil
	}

	pipe := q.client.Pipeline()

	// Delete all job data
	for _, key := range jobKeys {
		pipe.Del(ctx, key)
		pipe.SRem(ctx, JobsSetKey, key)
	}

	// Delete the channel jobs set
	pipe.Del(ctx, channelJobsKey)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to clear channel jobs: %w", err)
	}

	return nil
}

// GetActiveJobCount returns the number of active jobs for a channel.
func (q *Queue) GetActiveJobCount(channelID string) (int64, error) {
	if channelID == "" {
		return 0, fmt.Errorf("channel ID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(q.ctx, DefaultTimeout)
	defer cancel()

	channelJobsKey := getChannelJobsKey(channelID)
	count, err := q.client.SCard(ctx, channelJobsKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get active job count: %w", err)
	}

	return count, nil
}

// GetTotalActiveJobCount returns the total number of active jobs across all channels.
func (q *Queue) GetTotalActiveJobCount() (int64, error) {
	ctx, cancel := context.WithTimeout(q.ctx, DefaultTimeout)
	defer cancel()

	count, err := q.client.SCard(ctx, JobsSetKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get total active job count: %w", err)
	}

	return count, nil
}
