package scheduler

import (
	"context"
	"fmt"
	"log"
	"os"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	// TTLSecondsAfterFinished is how long to keep completed jobs
	TTLSecondsAfterFinished = 300

	// Default worker image - should be overridden by environment variable
	DefaultWorkerImage = "ko.local/ytarchive-worker:latest"

	// Job labels
	LabelApp       = "ytarchive"
	LabelComponent = "worker"
)

// K8sJobManager handles Kubernetes job operations
type K8sJobManager struct {
	client    *kubernetes.Clientset
	namespace string
}

// NewK8sJobManager creates a new K8sJobManager
func NewK8sJobManager(client *kubernetes.Clientset, namespace string) *K8sJobManager {
	return &K8sJobManager{
		client:    client,
		namespace: namespace,
	}
}

// CreateWorkerJob creates a Kubernetes Job for a worker
func (m *K8sJobManager) CreateWorkerJob(ctx context.Context, channelID string, workerNum int, jobID string) (string, error) {
	jobName := fmt.Sprintf("ytarchive-worker-%s-%d", channelID[:8], workerNum)

	// Truncate job name if too long (K8s limit is 63 characters)
	if len(jobName) > 63 {
		jobName = jobName[:63]
	}

	workerImage := getEnvWithDefault("WORKER_IMAGE", DefaultWorkerImage)
	redisURL := getEnvWithDefault("REDIS_URL", "redis:6379")
	storagePath := getEnvWithDefault("STORAGE_PATH", "/data/videos")

	ttlSeconds := int32(TTLSecondsAfterFinished)
	backoffLimit := int32(0) // No K8s-level retries; worker handles retries internally with exponential backoff
	parallelism := int32(1)
	completions := int32(1)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: m.namespace,
			Labels: map[string]string{
				"app":        LabelApp,
				"component":  LabelComponent,
				"channel-id": channelID,
				"job-id":     jobID,
				"worker-num": fmt.Sprintf("%d", workerNum),
			},
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: &ttlSeconds,
			BackoffLimit:            &backoffLimit,
			Parallelism:             &parallelism,
			Completions:             &completions,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":        LabelApp,
						"component":  LabelComponent,
						"channel-id": channelID,
						"job-id":     jobID,
					},
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{Name: getEnvWithDefault("IMAGE_PULL_SECRET", "ghcr-secret")},
					},
					NodeSelector: getNodeSelector(),
					DNSConfig: &corev1.PodDNSConfig{
						Options: []corev1.PodDNSConfigOption{
							{Name: "ndots", Value: stringPtr("1")},
							{Name: "single-request-reopen", Value: nil},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "worker",
							Image: workerImage,
							Env: []corev1.EnvVar{
								{
									Name:  "REDIS_URL",
									Value: redisURL,
								},
								{
									Name: "REDIS_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "ytarchive-secrets",
											},
											Key:      "redis-password",
											Optional: boolPtr(true),
										},
									},
								},
								{
									Name:  "STORAGE_PATH",
									Value: storagePath,
								},
								{
									Name:  "CHANNEL_ID",
									Value: channelID,
								},
								{
									Name:  "JOB_ID",
									Value: jobID,
								},
								{
									Name:  "WORKER_NUM",
									Value: fmt.Sprintf("%d", workerNum),
								},
								{
									Name:  "CONTROLLER_URL",
									Value: getEnvWithDefault("CONTROLLER_URL", "http://ytarchive-controller.ytarchive.svc.cluster.local"),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "video-storage",
									MountPath: storagePath,
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("100m"),
									corev1.ResourceMemory: resource.MustParse("256Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("500m"),
									corev1.ResourceMemory: resource.MustParse("512Mi"),
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "video-storage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: getEnvWithDefault("STORAGE_PVC_NAME", "ytarchive-storage"),
								},
							},
						},
					},
				},
			},
		},
	}

	createdJob, err := m.client.BatchV1().Jobs(m.namespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create K8s job: %w", err)
	}

	log.Printf("Created K8s job: %s in namespace %s", createdJob.Name, m.namespace)
	return createdJob.Name, nil
}

// GetJobStatus retrieves the status of a Kubernetes job as a simple string
// Returns: "Succeeded", "Failed", "Running", or "Unknown"
func (m *K8sJobManager) GetJobStatus(ctx context.Context, jobName string) (string, error) {
	job, err := m.client.BatchV1().Jobs(m.namespace).Get(ctx, jobName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get K8s job status: %w", err)
	}

	// Check completion conditions
	if job.Status.Succeeded > 0 {
		return "Succeeded", nil
	}
	if job.Status.Failed > 0 {
		return "Failed", nil
	}
	if job.Status.Active > 0 {
		return "Running", nil
	}

	// Check job conditions for more detail
	for _, cond := range job.Status.Conditions {
		if cond.Type == batchv1.JobComplete && cond.Status == corev1.ConditionTrue {
			return "Succeeded", nil
		}
		if cond.Type == batchv1.JobFailed && cond.Status == corev1.ConditionTrue {
			return "Failed", nil
		}
	}

	return "Running", nil
}

// GetJobFullStatus retrieves the full status struct of a Kubernetes job
func (m *K8sJobManager) GetJobFullStatus(ctx context.Context, jobName string) (*batchv1.JobStatus, error) {
	job, err := m.client.BatchV1().Jobs(m.namespace).Get(ctx, jobName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get K8s job status: %w", err)
	}
	return &job.Status, nil
}

// DeleteJob deletes a Kubernetes job
func (m *K8sJobManager) DeleteJob(ctx context.Context, jobName string) error {
	propagationPolicy := metav1.DeletePropagationBackground
	err := m.client.BatchV1().Jobs(m.namespace).Delete(ctx, jobName, metav1.DeleteOptions{
		PropagationPolicy: &propagationPolicy,
	})
	if err != nil {
		return fmt.Errorf("failed to delete K8s job: %w", err)
	}
	log.Printf("Deleted K8s job: %s", jobName)
	return nil
}

// ListJobsByChannel lists all jobs for a specific channel
func (m *K8sJobManager) ListJobsByChannel(ctx context.Context, channelID string) (*batchv1.JobList, error) {
	labelSelector := fmt.Sprintf("app=%s,component=%s,channel-id=%s", LabelApp, LabelComponent, channelID)

	jobs, err := m.client.BatchV1().Jobs(m.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list K8s jobs: %w", err)
	}
	return jobs, nil
}

// CleanupCompletedJobs removes completed jobs older than the TTL
func (m *K8sJobManager) CleanupCompletedJobs(ctx context.Context) error {
	labelSelector := fmt.Sprintf("app=%s,component=%s", LabelApp, LabelComponent)

	jobs, err := m.client.BatchV1().Jobs(m.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return fmt.Errorf("failed to list K8s jobs for cleanup: %w", err)
	}

	for _, job := range jobs.Items {
		// Check if job is completed or failed
		if job.Status.Succeeded > 0 || job.Status.Failed > 0 {
			// TTL controller should handle this, but we can force cleanup if needed
			if job.Spec.TTLSecondsAfterFinished == nil {
				if err := m.DeleteJob(ctx, job.Name); err != nil {
					log.Printf("Warning: failed to cleanup job %s: %v", job.Name, err)
				}
			}
		}
	}

	return nil
}

// getEnvWithDefault returns environment variable value or default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// boolPtr returns a pointer to a bool
func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}

func getNodeSelector() map[string]string {
	nodeSelector := os.Getenv("WORKER_NODE_SELECTOR")
	if nodeSelector != "" {
		return map[string]string{"kubernetes.io/hostname": nodeSelector}
	}
	return nil
}
