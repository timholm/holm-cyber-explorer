#!/bin/bash
# Deploy YouTube Channel Archiver to Kubernetes for testing
# Usage: ./scripts/deploy-test.sh [REGISTRY]
# Example: ./scripts/deploy-test.sh ghcr.io/timholm

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
REGISTRY="${1:-ko.local}"

echo "=== YouTube Channel Archiver - K8s Test Deployment ==="
echo "Registry: $REGISTRY"
echo "Project: $PROJECT_DIR"
echo ""

# Check prerequisites
command -v kubectl >/dev/null 2>&1 || { echo "kubectl required but not installed"; exit 1; }
command -v ko >/dev/null 2>&1 || { echo "ko required but not installed. Install: go install github.com/google/ko@latest"; exit 1; }

# Set Ko registry
export KO_DOCKER_REPO="$REGISTRY"

echo "1. Creating namespace..."
kubectl apply -f "$PROJECT_DIR/deploy/kubernetes/namespace.yaml"

echo "2. Deploying ConfigMap and Secrets..."
kubectl apply -f "$PROJECT_DIR/deploy/kubernetes/configmap.yaml"
kubectl apply -f "$PROJECT_DIR/deploy/kubernetes/secret.yaml"

echo "3. Setting up RBAC..."
kubectl apply -f "$PROJECT_DIR/deploy/kubernetes/rbac.yaml"

echo "4. Deploying Redis..."
kubectl apply -f "$PROJECT_DIR/deploy/kubernetes/redis.yaml"

echo "5. Waiting for Redis to be ready..."
kubectl -n ytarchive rollout status statefulset/redis --timeout=120s || {
    echo "Warning: Redis not ready yet, continuing..."
}

echo "6. Creating PVC (if not using iSCSI, will create hostPath)..."
kubectl apply -f "$PROJECT_DIR/deploy/kubernetes/pvc-iscsi.yaml" 2>/dev/null || {
    echo "iSCSI PVC failed, creating test PVC with hostPath..."
    cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: ytarchive-data
  namespace: ytarchive
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
EOF
}

echo "7. Building and deploying Controller with Ko..."
cd "$PROJECT_DIR"
ko apply -f deploy/kubernetes/controller.yaml

echo "8. Deploying Service..."
kubectl apply -f "$PROJECT_DIR/deploy/kubernetes/controller-service.yaml"

echo "9. Deploying HPA for auto-scaling..."
kubectl apply -f "$PROJECT_DIR/deploy/kubernetes/hpa.yaml"

echo "10. Waiting for Controller to be ready..."
kubectl -n ytarchive rollout status deployment/ytarchive-controller --timeout=180s

echo ""
echo "=== Deployment Complete ==="
echo ""
echo "Controller is running. To access the API:"
echo ""
echo "  # Port forward to localhost:8080"
echo "  kubectl -n ytarchive port-forward svc/ytarchive-controller 8080:80"
echo ""
echo "  # Then test with:"
echo "  curl http://localhost:8080/healthz"
echo "  curl http://localhost:8080/api/channels"
echo ""
echo "  # Add test channel:"
echo "  curl -X POST http://localhost:8080/api/channels \\"
echo "    -H 'Content-Type: application/json' \\"
echo "    -d '{\"youtube_url\": \"https://youtube.com/@aperturethinking\"}'"
echo ""
echo "To view logs:"
echo "  kubectl -n ytarchive logs -f deployment/ytarchive-controller"
echo ""
echo "To check HPA status:"
echo "  kubectl -n ytarchive get hpa"
