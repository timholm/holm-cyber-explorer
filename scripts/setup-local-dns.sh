#!/bin/bash
# ═══════════════════════════════════════════════════════════════
# HOLM Local DNS Setup
# ═══════════════════════════════════════════════════════════════
# Configures local machines to resolve *.holm.chat to the
# MetalLB LoadBalancer IP on the local network.
#
# Run on each machine that needs to access HOLM services:
#   sudo ./setup-local-dns.sh [METALLB_IP]
#
# Default IP: 192.168.8.50 (first in MetalLB pool)
# ═══════════════════════════════════════════════════════════════

set -euo pipefail

VIP="${1:-192.168.8.50}"
HOSTS_FILE="/etc/hosts"
MARKER_START="# === HOLM LOCAL DNS START ==="
MARKER_END="# === HOLM LOCAL DNS END ==="

DOMAINS=(
  "holm.chat"
  "holm.local"
  "dev.holm.chat"
  "mind.holm.chat"
  "cmd.holm.chat"
  "hq.holm.chat"
  "uptime.holm.chat"
  "search.holm.chat"
  "animus.holm.chat"
  "bookforge.holm.chat"
  "vault.holm.chat"
  "pihole.holm.chat"
  "matrix.holm.chat"
  "auth.holm.chat"
  "odoo.holm.chat"
  "grafana.holm.chat"
)

echo "[holm] Setting up local DNS → $VIP"

# Remove existing HOLM entries if present
if grep -q "$MARKER_START" "$HOSTS_FILE" 2>/dev/null; then
  echo "[holm] Removing existing entries..."
  sed -i.bak "/$MARKER_START/,/$MARKER_END/d" "$HOSTS_FILE"
fi

# Add new entries
echo "" >> "$HOSTS_FILE"
echo "$MARKER_START" >> "$HOSTS_FILE"
for domain in "${DOMAINS[@]}"; do
  echo "$VIP $domain" >> "$HOSTS_FILE"
done
echo "$MARKER_END" >> "$HOSTS_FILE"

echo "[holm] Added ${#DOMAINS[@]} DNS entries pointing to $VIP"
echo "[holm] Verify: ping -c1 holm.chat"
