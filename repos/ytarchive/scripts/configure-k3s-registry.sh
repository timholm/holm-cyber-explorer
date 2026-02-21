#!/bin/bash
# Configure k3s nodes to use insecure registry
# Run this script on EACH k3s node (via SSH)

REGISTRY_IP="${1:-192.168.8.197}"
REGISTRY_PORT="${2:-30500}"

echo "Configuring k3s to use insecure registry at $REGISTRY_IP:$REGISTRY_PORT"

# Create registries.yaml
sudo mkdir -p /etc/rancher/k3s/
sudo tee /etc/rancher/k3s/registries.yaml > /dev/null <<EOF
mirrors:
  "$REGISTRY_IP:$REGISTRY_PORT":
    endpoint:
      - "http://$REGISTRY_IP:$REGISTRY_PORT"
configs:
  "$REGISTRY_IP:$REGISTRY_PORT":
    tls:
      insecure_skip_verify: true
EOF

echo "Created /etc/rancher/k3s/registries.yaml"
cat /etc/rancher/k3s/registries.yaml

echo ""
echo "Restarting k3s service..."
sudo systemctl restart k3s || sudo systemctl restart k3s-agent

echo "Done! k3s should now be able to pull from $REGISTRY_IP:$REGISTRY_PORT"
