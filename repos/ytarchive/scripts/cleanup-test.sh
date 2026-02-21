#!/bin/bash
# Cleanup YouTube Channel Archiver test deployment
# Usage: ./scripts/cleanup-test.sh

set -e

NAMESPACE="${NAMESPACE:-ytarchive}"

echo "=== Cleaning up ytarchive test deployment ==="
echo "Namespace: $NAMESPACE"
echo ""

read -p "This will delete all ytarchive resources. Continue? [y/N] " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 1
fi

echo "Deleting all jobs..."
kubectl -n "$NAMESPACE" delete jobs --all 2>/dev/null || true

echo "Deleting deployments..."
kubectl -n "$NAMESPACE" delete deployment ytarchive-controller 2>/dev/null || true

echo "Deleting services..."
kubectl -n "$NAMESPACE" delete svc ytarchive-controller 2>/dev/null || true

echo "Deleting Redis..."
kubectl -n "$NAMESPACE" delete statefulset redis 2>/dev/null || true
kubectl -n "$NAMESPACE" delete svc redis-service 2>/dev/null || true

echo "Deleting HPA..."
kubectl -n "$NAMESPACE" delete hpa ytarchive-controller 2>/dev/null || true

echo "Deleting ConfigMaps and Secrets..."
kubectl -n "$NAMESPACE" delete configmap ytarchive-config 2>/dev/null || true
kubectl -n "$NAMESPACE" delete secret ytarchive-secrets 2>/dev/null || true

echo "Deleting RBAC..."
kubectl -n "$NAMESPACE" delete serviceaccount ytarchive-controller 2>/dev/null || true
kubectl -n "$NAMESPACE" delete role ytarchive-controller 2>/dev/null || true
kubectl -n "$NAMESPACE" delete rolebinding ytarchive-controller 2>/dev/null || true

echo ""
read -p "Delete PVC (will lose all downloaded videos)? [y/N] " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Deleting PVC..."
    kubectl -n "$NAMESPACE" delete pvc ytarchive-data 2>/dev/null || true
    kubectl -n "$NAMESPACE" delete pvc redis-data-redis-0 2>/dev/null || true
fi

echo ""
read -p "Delete namespace? [y/N] " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Deleting namespace..."
    kubectl delete namespace "$NAMESPACE" 2>/dev/null || true
fi

echo ""
echo "=== Cleanup Complete ==="
