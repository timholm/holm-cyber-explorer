#!/bin/bash
# Test YouTube Channel Archiver API on Kubernetes
# Prerequisites: kubectl configured, ytarchive namespace exists
# Usage: ./scripts/test-k8s.sh

set -e

NAMESPACE="${NAMESPACE:-ytarchive}"
TEST_CHANNEL="https://youtube.com/@aperturethinking"

echo "=== YouTube Channel Archiver - K8s API Test ==="
echo "Namespace: $NAMESPACE"
echo ""

# Start port-forward in background
echo "Starting port-forward..."
kubectl -n "$NAMESPACE" port-forward svc/ytarchive-controller 8080:80 &
PF_PID=$!
sleep 3

# Cleanup on exit
cleanup() {
    echo "Stopping port-forward..."
    kill $PF_PID 2>/dev/null || true
}
trap cleanup EXIT

API_URL="http://localhost:8080"

# Health check
echo "1. Health Check..."
curl -s "$API_URL/healthz" && echo ""
echo ""

# Ready check
echo "2. Ready Check..."
curl -s "$API_URL/readyz" && echo ""
echo ""

# Add channel
echo "3. Adding test channel: $TEST_CHANNEL"
RESPONSE=$(curl -s -X POST "$API_URL/api/channels" \
  -H "Content-Type: application/json" \
  -d "{\"youtube_url\": \"$TEST_CHANNEL\"}")
echo "$RESPONSE" | jq . 2>/dev/null || echo "$RESPONSE"
CHANNEL_ID=$(echo "$RESPONSE" | jq -r '.id' 2>/dev/null || echo "")
echo "Channel ID: $CHANNEL_ID"
echo ""

if [ -z "$CHANNEL_ID" ] || [ "$CHANNEL_ID" = "null" ]; then
    echo "Failed to create channel. Checking existing channels..."
    curl -s "$API_URL/api/channels" | jq .
    exit 1
fi

# List channels
echo "4. List all channels..."
curl -s "$API_URL/api/channels" | jq .
echo ""

# Trigger sync
echo "5. Trigger channel sync..."
SYNC_RESPONSE=$(curl -s -X POST "$API_URL/api/channels/$CHANNEL_ID/sync")
echo "$SYNC_RESPONSE" | jq . 2>/dev/null || echo "$SYNC_RESPONSE"
JOB_ID=$(echo "$SYNC_RESPONSE" | jq -r '.job_id' 2>/dev/null || echo "")
echo "Job ID: $JOB_ID"
echo ""

# Monitor progress
echo "6. Monitoring sync progress (Ctrl+C to stop)..."
echo "   Workers will be spawned as K8s Jobs automatically."
echo ""

for i in {1..30}; do
    echo "--- Progress check $i ---"

    # Check progress
    echo "Progress:"
    curl -s "$API_URL/api/progress" | jq .

    # Check active jobs
    echo "Active Jobs:"
    curl -s "$API_URL/api/jobs" | jq '.count'

    # Check K8s jobs
    echo "K8s Worker Jobs:"
    kubectl -n "$NAMESPACE" get jobs -l app.kubernetes.io/part-of=ytarchive 2>/dev/null || echo "No jobs yet"

    # Check if sync is complete
    CHANNEL_STATUS=$(curl -s "$API_URL/api/channels/$CHANNEL_ID" | jq -r '.channel.status' 2>/dev/null)
    echo "Channel Status: $CHANNEL_STATUS"

    if [ "$CHANNEL_STATUS" = "synced" ] || [ "$CHANNEL_STATUS" = "completed" ]; then
        echo ""
        echo "=== Sync Complete! ==="
        break
    fi

    echo ""
    sleep 10
done

# Final stats
echo ""
echo "=== Final Stats ==="
curl -s "$API_URL/api/stats" | jq .

echo ""
echo "=== Test Complete ==="
echo ""
echo "To continue monitoring:"
echo "  watch -n 5 'kubectl -n $NAMESPACE port-forward svc/ytarchive-controller 8080:80 & sleep 2 && curl -s http://localhost:8080/api/progress | jq . && kill %1'"
