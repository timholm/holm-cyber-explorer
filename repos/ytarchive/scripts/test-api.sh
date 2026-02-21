#!/bin/bash
# Test script for YouTube Channel Archiver API
# Prerequisites: Redis running on localhost:6379, yt-dlp installed

set -e

API_URL="${API_URL:-http://localhost:8080}"
TEST_CHANNEL="https://youtube.com/@aperturethinking"

echo "=== YouTube Channel Archiver API Test ==="
echo "API URL: $API_URL"
echo ""

# Health check
echo "1. Health Check..."
curl -s "$API_URL/healthz" | jq .
echo ""

# Add channel
echo "2. Adding test channel: $TEST_CHANNEL"
CHANNEL_RESPONSE=$(curl -s -X POST "$API_URL/api/channels" \
  -H "Content-Type: application/json" \
  -d "{\"youtube_url\": \"$TEST_CHANNEL\"}")
echo "$CHANNEL_RESPONSE" | jq .
CHANNEL_ID=$(echo "$CHANNEL_RESPONSE" | jq -r '.id')
echo "Channel ID: $CHANNEL_ID"
echo ""

# List channels
echo "3. List all channels..."
curl -s "$API_URL/api/channels" | jq .
echo ""

# Get channel details
echo "4. Get channel details..."
curl -s "$API_URL/api/channels/$CHANNEL_ID" | jq .
echo ""

# Trigger sync
echo "5. Trigger channel sync..."
curl -s -X POST "$API_URL/api/channels/$CHANNEL_ID/sync" | jq .
echo ""

# Check progress
echo "6. Check progress (wait 5 seconds)..."
sleep 5
curl -s "$API_URL/api/progress" | jq .
echo ""

# List jobs
echo "7. List active jobs..."
curl -s "$API_URL/api/jobs" | jq .
echo ""

# Get stats
echo "8. Get stats..."
curl -s "$API_URL/api/stats" | jq .
echo ""

echo "=== Test Complete ==="
echo ""
echo "To monitor sync progress, run:"
echo "  watch -n 2 'curl -s $API_URL/api/progress | jq .'"
