#!/bin/bash
set -e

# Clone docs-framework for the build context
echo "Fetching docs-framework..."
rm -rf /tmp/docs-framework
git clone --depth 1 https://github.com/timholm/docs-framework.git /tmp/docs-framework
cp /tmp/docs-framework/html/manifest.json ./manifest.json
cp -r /tmp/docs-framework/html ./html

# Build and push
echo "Building Docker image..."
docker buildx build --platform linux/arm64 -t gitea.holm.chat/tim/holm-cyber-explorer:latest --push .

# Clean up
rm -rf ./html ./manifest.json /tmp/docs-framework

echo "Done. Image pushed to gitea.holm.chat/tim/holm-cyber-explorer:latest"
