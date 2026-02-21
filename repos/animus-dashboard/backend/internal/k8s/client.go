package k8s

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"
)

type Client struct {
	clientset     *kubernetes.Clientset
	metricsClient *metricsclient.Clientset
}

type Node struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Status     string   `json:"status"`
	CPU        int      `json:"cpu"`
	Memory     int      `json:"memory"`
	Disk       int      `json:"disk"`
	PodCount   int      `json:"podCount"`
	K3sVersion string   `json:"k3sVersion"`
	OSVersion  string   `json:"osVersion"`
	Uptime     string   `json:"uptime"`
	IP         string   `json:"ip"`
	Roles      []string `json:"roles"`
	HasUpdate  bool     `json:"hasUpdate"`
}

type Pod struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
	Restarts  int32  `json:"restarts"`
	Age       string `json:"age"`
	CPU       string `json:"cpu"`
	Memory    string `json:"memory"`
	NodeName  string `json:"nodeName"`
}

type ClusterMetrics struct {
	TotalNodes   int `json:"totalNodes"`
	HealthyNodes int `json:"healthyNodes"`
	TotalPods    int `json:"totalPods"`
	RunningPods  int `json:"runningPods"`
	CPUUsage     int `json:"cpuUsage"`
	MemoryUsage  int `json:"memoryUsage"`
}

func NewClient() (*Client, error) {
	var config *rest.Config
	var err error

	// Try in-cluster config first
	config, err = rest.InClusterConfig()
	if err != nil {
		// Fall back to kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create kubernetes config: %w", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	metricsClient, err := metricsclient.NewForConfig(config)
	if err != nil {
		// Metrics client is optional
		metricsClient = nil
	}

	return &Client{
		clientset:     clientset,
		metricsClient: metricsClient,
	}, nil
}

func (c *Client) GetNodes(ctx context.Context) ([]Node, error) {
	nodeList, err := c.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	// Get node metrics if available
	var nodeMetrics *metricsv1beta1.NodeMetricsList
	if c.metricsClient != nil {
		nodeMetrics, _ = c.metricsClient.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
	}

	// Get pods per node
	podList, err := c.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	podCountByNode := make(map[string]int)
	for _, pod := range podList.Items {
		podCountByNode[pod.Spec.NodeName]++
	}

	nodes := make([]Node, 0, len(nodeList.Items))
	for _, n := range nodeList.Items {
		node := c.convertNode(&n, nodeMetrics, podCountByNode[n.Name])
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (c *Client) GetNode(ctx context.Context, name string) (*Node, error) {
	n, err := c.clientset.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	// Get metrics for this node
	var nodeMetrics *metricsv1beta1.NodeMetricsList
	if c.metricsClient != nil {
		nodeMetrics, _ = c.metricsClient.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
	}

	// Count pods on this node
	podList, err := c.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", name),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	node := c.convertNode(n, nodeMetrics, len(podList.Items))
	return &node, nil
}

func (c *Client) GetNodePods(ctx context.Context, nodeName string) ([]Pod, error) {
	podList, err := c.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	pods := make([]Pod, 0, len(podList.Items))
	for _, p := range podList.Items {
		pod := c.convertPod(&p)
		pods = append(pods, pod)
	}

	return pods, nil
}

func (c *Client) GetClusterMetrics(ctx context.Context) (*ClusterMetrics, error) {
	nodeList, err := c.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	podList, err := c.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	healthyNodes := 0
	for _, node := range nodeList.Items {
		if isNodeReady(&node) {
			healthyNodes++
		}
	}

	runningPods := 0
	for _, pod := range podList.Items {
		if pod.Status.Phase == corev1.PodRunning {
			runningPods++
		}
	}

	// Calculate average CPU/Memory usage
	cpuUsage, memoryUsage := 0, 0
	if c.metricsClient != nil {
		nodeMetrics, err := c.metricsClient.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
		if err == nil && len(nodeMetrics.Items) > 0 {
			totalCPU, totalMem := int64(0), int64(0)
			for _, nm := range nodeMetrics.Items {
				totalCPU += nm.Usage.Cpu().MilliValue()
				totalMem += nm.Usage.Memory().Value()
			}
			// These are simplified calculations
			cpuUsage = int(totalCPU / int64(len(nodeMetrics.Items)) / 10) // Rough percentage
			memoryUsage = int(totalMem / int64(len(nodeMetrics.Items)) / (1024 * 1024 * 1024) * 100 / 8) // Assume 8GB per node
		}
	}

	return &ClusterMetrics{
		TotalNodes:   len(nodeList.Items),
		HealthyNodes: healthyNodes,
		TotalPods:    len(podList.Items),
		RunningPods:  runningPods,
		CPUUsage:     cpuUsage,
		MemoryUsage:  memoryUsage,
	}, nil
}

func (c *Client) convertNode(n *corev1.Node, metrics *metricsv1beta1.NodeMetricsList, podCount int) Node {
	// Get node status
	status := "Unknown"
	if isNodeReady(n) {
		status = "Ready"
	} else {
		status = "NotReady"
	}

	// Get internal IP
	ip := ""
	for _, addr := range n.Status.Addresses {
		if addr.Type == corev1.NodeInternalIP {
			ip = addr.Address
			break
		}
	}

	// Get roles
	roles := []string{}
	for label := range n.Labels {
		if label == "node-role.kubernetes.io/control-plane" || label == "node-role.kubernetes.io/master" {
			roles = append(roles, "control-plane")
		}
	}
	if len(roles) == 0 {
		roles = append(roles, "worker")
	}

	// Calculate uptime
	uptime := formatDuration(time.Since(n.CreationTimestamp.Time))

	// Get CPU/Memory usage from metrics
	cpu, memory := 0, 0
	if metrics != nil {
		for _, m := range metrics.Items {
			if m.Name == n.Name {
				// Convert to percentage (simplified)
				cpuQuantity := m.Usage.Cpu()
				memQuantity := m.Usage.Memory()
				cpu = int(cpuQuantity.MilliValue() / 10) // Rough percentage
				memory = int(memQuantity.Value() / (1024 * 1024 * 1024) * 100 / 8) // Assume 8GB
				if cpu > 100 {
					cpu = 100
				}
				if memory > 100 {
					memory = 100
				}
				break
			}
		}
	}

	return Node{
		ID:         n.Name,
		Name:       n.Name,
		Status:     status,
		CPU:        cpu,
		Memory:     memory,
		Disk:       30, // Would need to be fetched via other means
		PodCount:   podCount,
		K3sVersion: n.Status.NodeInfo.KubeletVersion,
		OSVersion:  n.Status.NodeInfo.OSImage,
		Uptime:     uptime,
		IP:         ip,
		Roles:      roles,
		HasUpdate:  false, // Would need external check
	}
}

func (c *Client) convertPod(p *corev1.Pod) Pod {
	// Calculate restarts
	restarts := int32(0)
	for _, cs := range p.Status.ContainerStatuses {
		restarts += cs.RestartCount
	}

	// Calculate age
	age := formatDuration(time.Since(p.CreationTimestamp.Time))

	return Pod{
		Name:      p.Name,
		Namespace: p.Namespace,
		Status:    string(p.Status.Phase),
		Restarts:  restarts,
		Age:       age,
		CPU:       "N/A",
		Memory:    "N/A",
		NodeName:  p.Spec.NodeName,
	}
}

func isNodeReady(node *corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
