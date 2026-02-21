# Deployment Guide

Step-by-step guide for deploying the YouTube Channel Archiver to Kubernetes.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Storage Setup](#storage-setup)
3. [Redis Deployment](#redis-deployment)
4. [Controller Deployment](#controller-deployment)
5. [ArgoCD Deployment](#argocd-deployment)
6. [Verification](#verification)
7. [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Kubernetes Cluster

- Kubernetes v1.28 or later
- kubectl configured with cluster access
- Cluster admin permissions

```bash
# Verify cluster access
kubectl cluster-info
kubectl get nodes
```

### NetApp Trident

Trident provides dynamic storage provisioning for iSCSI volumes.

```bash
# Check if Trident is installed
kubectl get pods -n trident

# Verify Trident backend
kubectl get tbc -n trident
```

If Trident is not installed, follow the [Trident installation guide](https://docs.netapp.com/us-en/trident/trident-get-started/kubernetes-deploy.html).

### Argo Workflows (Optional)

For workflow-based download pipelines:

```bash
# Install Argo Workflows
kubectl create namespace argo
kubectl apply -n argo -f https://github.com/argoproj/argo-workflows/releases/download/v3.5.0/install.yaml

# Verify installation
kubectl get pods -n argo
```

### ArgoCD (Optional)

For GitOps-based deployments:

```bash
# Install ArgoCD
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Get initial admin password
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```

### Container Registry Access

Ensure your cluster can pull images from your container registry:

```bash
# Create image pull secret if needed
kubectl create secret docker-registry regcred \
  --docker-server=<your-registry> \
  --docker-username=<username> \
  --docker-password=<password> \
  --namespace=ytarchive
```

---

## Storage Setup

### Create Namespace

```bash
kubectl create namespace ytarchive
```

### Trident Storage Class

Create a storage class for video archives:

```yaml
# deploy/kubernetes/storage-class.yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: ytarchive-storage
provisioner: csi.trident.netapp.io
parameters:
  backendType: "ontap-san"
  storagePools: "san-pool"
  fsType: "ext4"
reclaimPolicy: Retain
allowVolumeExpansion: true
volumeBindingMode: Immediate
```

```bash
kubectl apply -f deploy/kubernetes/storage-class.yaml
```

### Persistent Volume Claim

Create a PVC for video storage:

```yaml
# deploy/kubernetes/pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: ytarchive-storage
  namespace: ytarchive
spec:
  accessModes:
    - ReadWriteMany
  storageClassName: ytarchive-storage
  resources:
    requests:
      storage: 1Ti
```

```bash
kubectl apply -f deploy/kubernetes/pvc.yaml

# Verify PVC is bound
kubectl get pvc -n ytarchive
```

---

## Redis Deployment

### Redis StatefulSet

```yaml
# deploy/kubernetes/redis.yaml
apiVersion: v1
kind: Service
metadata:
  name: redis
  namespace: ytarchive
spec:
  ports:
    - port: 6379
      targetPort: 6379
  selector:
    app: redis
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: redis
  namespace: ytarchive
spec:
  serviceName: redis
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
        - name: redis
          image: redis:7-alpine
          ports:
            - containerPort: 6379
          resources:
            requests:
              memory: "256Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"
          volumeMounts:
            - name: redis-data
              mountPath: /data
          command:
            - redis-server
            - --appendonly
            - "yes"
            - --maxmemory
            - "256mb"
            - --maxmemory-policy
            - "allkeys-lru"
  volumeClaimTemplates:
    - metadata:
        name: redis-data
      spec:
        accessModes: ["ReadWriteOnce"]
        storageClassName: ytarchive-storage
        resources:
          requests:
            storage: 10Gi
```

```bash
kubectl apply -f deploy/kubernetes/redis.yaml

# Verify Redis is running
kubectl get pods -n ytarchive -l app=redis
kubectl logs -n ytarchive redis-0
```

---

## Controller Deployment

### ConfigMap

```yaml
# deploy/kubernetes/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: ytarchive-config
  namespace: ytarchive
data:
  REDIS_ADDR: "redis:6379"
  K8S_NAMESPACE: "ytarchive"
  STORAGE_PATH: "/archive"
  MAX_WORKERS: "5"
  DOWNLOAD_FORMAT: "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best"
  RATE_LIMIT: "5M"
  LOG_LEVEL: "info"
```

### Service Account and RBAC

```yaml
# deploy/kubernetes/rbac.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: ytarchive
  namespace: ytarchive
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: ytarchive
  namespace: ytarchive
rules:
  - apiGroups: [""]
    resources: ["pods", "pods/log"]
    verbs: ["get", "list", "watch"]
  - apiGroups: ["batch"]
    resources: ["jobs"]
    verbs: ["get", "list", "watch", "create", "delete"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: ytarchive
  namespace: ytarchive
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: ytarchive
subjects:
  - kind: ServiceAccount
    name: ytarchive
    namespace: ytarchive
```

### Deployment

```yaml
# deploy/kubernetes/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ytarchive-controller
  namespace: ytarchive
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ytarchive-controller
  template:
    metadata:
      labels:
        app: ytarchive-controller
    spec:
      serviceAccountName: ytarchive
      containers:
        - name: controller
          image: ghcr.io/timholm/ytarchive:latest
          ports:
            - containerPort: 8080
          envFrom:
            - configMapRef:
                name: ytarchive-config
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 10
            periodSeconds: 30
          readinessProbe:
            httpGet:
              path: /ready
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10
          volumeMounts:
            - name: archive-storage
              mountPath: /archive
      volumes:
        - name: archive-storage
          persistentVolumeClaim:
            claimName: ytarchive-storage
---
apiVersion: v1
kind: Service
metadata:
  name: ytarchive
  namespace: ytarchive
spec:
  type: ClusterIP
  ports:
    - port: 8080
      targetPort: 8080
  selector:
    app: ytarchive-controller
```

### Deploy All Components

```bash
# Apply all manifests
kubectl apply -f deploy/kubernetes/configmap.yaml
kubectl apply -f deploy/kubernetes/rbac.yaml
kubectl apply -f deploy/kubernetes/deployment.yaml

# Verify deployment
kubectl get pods -n ytarchive
kubectl get svc -n ytarchive
```

---

## ArgoCD Deployment

For GitOps-based deployment using ArgoCD:

### ArgoCD Application

```yaml
# deploy/argocd/application.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: ytarchive
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/timholm/ytarchive.git
    targetRevision: main
    path: deploy/kubernetes
  destination:
    server: https://kubernetes.default.svc
    namespace: ytarchive
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
```

```bash
# Apply ArgoCD application
kubectl apply -f deploy/argocd/application.yaml

# Check sync status
argocd app get ytarchive

# Manual sync if needed
argocd app sync ytarchive
```

### ArgoCD Project (Optional)

```yaml
# deploy/argocd/project.yaml
apiVersion: argoproj.io/v1alpha1
kind: AppProject
metadata:
  name: ytarchive
  namespace: argocd
spec:
  description: YouTube Channel Archiver
  sourceRepos:
    - https://github.com/timholm/ytarchive.git
  destinations:
    - namespace: ytarchive
      server: https://kubernetes.default.svc
  clusterResourceWhitelist:
    - group: ""
      kind: Namespace
  namespaceResourceWhitelist:
    - group: ""
      kind: "*"
    - group: "apps"
      kind: "*"
    - group: "batch"
      kind: "*"
```

---

## Verification

### Check All Components

```bash
# Check all pods
kubectl get pods -n ytarchive

# Expected output:
# NAME                                    READY   STATUS    RESTARTS   AGE
# redis-0                                 1/1     Running   0          5m
# ytarchive-controller-xxx                1/1     Running   0          2m

# Check services
kubectl get svc -n ytarchive

# Check PVC
kubectl get pvc -n ytarchive

# Check logs
kubectl logs -n ytarchive deployment/ytarchive-controller
```

### Test API

```bash
# Port forward to access locally
kubectl port-forward -n ytarchive svc/ytarchive 8080:8080

# In another terminal, test endpoints
curl http://localhost:8080/health
curl http://localhost:8080/ready

# Add a test channel
curl -X POST http://localhost:8080/api/channels \
  -H "Content-Type: application/json" \
  -d '{"youtube_url": "https://www.youtube.com/@aperturethinking"}'

# Check channels
curl http://localhost:8080/api/channels

# Check progress
curl http://localhost:8080/api/progress
```

### Ingress (Optional)

```yaml
# deploy/kubernetes/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ytarchive
  namespace: ytarchive
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  tls:
    - hosts:
        - ytarchive.example.com
      secretName: ytarchive-tls
  rules:
    - host: ytarchive.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: ytarchive
                port:
                  number: 8080
```

---

## Troubleshooting

### Common Issues

#### Pod Not Starting

```bash
# Check pod events
kubectl describe pod -n ytarchive <pod-name>

# Check logs
kubectl logs -n ytarchive <pod-name> --previous
```

#### Redis Connection Failed

```bash
# Verify Redis is running
kubectl exec -n ytarchive -it redis-0 -- redis-cli ping

# Check Redis logs
kubectl logs -n ytarchive redis-0
```

#### Storage Issues

```bash
# Check PVC status
kubectl describe pvc -n ytarchive ytarchive-storage

# Check Trident logs
kubectl logs -n trident -l app=trident-csi
```

#### Download Jobs Not Running

```bash
# Check job status
kubectl get jobs -n ytarchive

# Describe failed job
kubectl describe job -n ytarchive <job-name>

# Check worker pod logs
kubectl logs -n ytarchive <worker-pod-name>
```

### Useful Commands

```bash
# Watch pod status
watch kubectl get pods -n ytarchive

# Stream controller logs
kubectl logs -n ytarchive -f deployment/ytarchive-controller

# Check Redis queue length
kubectl exec -n ytarchive redis-0 -- redis-cli llen download:queue

# Check storage usage
kubectl exec -n ytarchive deployment/ytarchive-controller -- df -h /archive
```

### Cleanup

```bash
# Delete all resources
kubectl delete namespace ytarchive

# Or delete specific resources
kubectl delete -f deploy/kubernetes/
```
