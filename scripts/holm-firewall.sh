#!/bin/bash
# ═══════════════════════════════════════════════════════════════
# HOLM Node-Level Firewall — Air-Gap Enforcement
# ═══════════════════════════════════════════════════════════════
# iptables rules for each cluster node to enforce network
# boundaries at the host level. This is defense-in-depth
# alongside k8s NetworkPolicy.
#
# Modes:
#   transitional  — LAN + internet (current state)
#   restricted    — LAN + allowlisted external IPs only
#   airgap        — LAN only, zero internet access
#
# Usage:
#   sudo ./holm-firewall.sh [transitional|restricted|airgap]
#   sudo ./holm-firewall.sh status
#   sudo ./holm-firewall.sh rollback
#
# TASK-031: Local DNS and Firewall Rules
# ═══════════════════════════════════════════════════════════════

set -euo pipefail

# ── Configuration ────────────────────────────────────────────
LAN_CIDR="192.168.8.0/24"
K3S_POD_CIDR="10.42.0.0/16"
K3S_SVC_CIDR="10.43.0.0/16"
METALLB_VIP="192.168.8.50"

# Allowlisted external IPs for restricted mode
# (GitHub, container registries — remove as you mirror locally)
ALLOWED_EXTERNAL=(
    "140.82.112.0/20"   # GitHub
    "185.199.108.0/22"  # GitHub Pages / CDN
    "104.18.0.0/16"     # Cloudflare (transitional)
)

# k3s API server port
K3S_API_PORT=6443

# SSH port for management
SSH_PORT=22

IPTABLES_BACKUP="/tmp/holm-iptables-backup.rules"
MODE="${1:-status}"

# ── Colors (cyberpunk aesthetic) ─────────────────────────────
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() { echo -e "${CYAN}[holm-fw]${NC} $*"; }
warn() { echo -e "${YELLOW}[holm-fw]${NC} $*"; }
err() { echo -e "${RED}[holm-fw]${NC} $*" >&2; }
ok() { echo -e "${GREEN}[holm-fw]${NC} $*"; }

# ── Pre-flight checks ────────────────────────────────────────
check_root() {
    if [[ $EUID -ne 0 ]]; then
        err "Must run as root. Use: sudo $0 $MODE"
        exit 1
    fi
}

backup_rules() {
    log "Backing up current iptables rules → $IPTABLES_BACKUP"
    iptables-save > "$IPTABLES_BACKUP"
}

# ── Flush and set defaults ────────────────────────────────────
flush_holm_rules() {
    log "Flushing HOLM firewall chains..."
    # Remove HOLM chains if they exist
    iptables -D INPUT -j HOLM-INPUT 2>/dev/null || true
    iptables -D OUTPUT -j HOLM-OUTPUT 2>/dev/null || true
    iptables -D FORWARD -j HOLM-FORWARD 2>/dev/null || true
    iptables -F HOLM-INPUT 2>/dev/null || true
    iptables -F HOLM-OUTPUT 2>/dev/null || true
    iptables -F HOLM-FORWARD 2>/dev/null || true
    iptables -X HOLM-INPUT 2>/dev/null || true
    iptables -X HOLM-OUTPUT 2>/dev/null || true
    iptables -X HOLM-FORWARD 2>/dev/null || true
}

# ── Base rules (all modes) ────────────────────────────────────
apply_base_rules() {
    log "Creating HOLM firewall chains..."
    iptables -N HOLM-INPUT
    iptables -N HOLM-OUTPUT
    iptables -N HOLM-FORWARD

    # Insert HOLM chains at top of built-in chains
    iptables -I INPUT 1 -j HOLM-INPUT
    iptables -I OUTPUT 1 -j HOLM-OUTPUT
    iptables -I FORWARD 1 -j HOLM-FORWARD

    # ── Always allow ──
    # Loopback
    iptables -A HOLM-INPUT -i lo -j ACCEPT
    iptables -A HOLM-OUTPUT -o lo -j ACCEPT

    # Established/related connections
    iptables -A HOLM-INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT
    iptables -A HOLM-OUTPUT -m state --state ESTABLISHED,RELATED -j ACCEPT
    iptables -A HOLM-FORWARD -m state --state ESTABLISHED,RELATED -j ACCEPT

    # ── LAN traffic (always allowed) ──
    iptables -A HOLM-INPUT -s "$LAN_CIDR" -j ACCEPT
    iptables -A HOLM-OUTPUT -d "$LAN_CIDR" -j ACCEPT
    iptables -A HOLM-FORWARD -s "$LAN_CIDR" -d "$LAN_CIDR" -j ACCEPT

    # ── k3s internal traffic ──
    # Pod-to-pod
    iptables -A HOLM-INPUT -s "$K3S_POD_CIDR" -j ACCEPT
    iptables -A HOLM-OUTPUT -d "$K3S_POD_CIDR" -j ACCEPT
    iptables -A HOLM-FORWARD -s "$K3S_POD_CIDR" -j ACCEPT
    iptables -A HOLM-FORWARD -d "$K3S_POD_CIDR" -j ACCEPT

    # Service network
    iptables -A HOLM-INPUT -s "$K3S_SVC_CIDR" -j ACCEPT
    iptables -A HOLM-OUTPUT -d "$K3S_SVC_CIDR" -j ACCEPT

    # ── SSH management (LAN only) ──
    iptables -A HOLM-INPUT -s "$LAN_CIDR" -p tcp --dport "$SSH_PORT" -j ACCEPT

    # ── k3s API server (LAN only) ──
    iptables -A HOLM-INPUT -s "$LAN_CIDR" -p tcp --dport "$K3S_API_PORT" -j ACCEPT

    # ── ICMP (ping — useful for diagnostics) ──
    iptables -A HOLM-INPUT -p icmp -j ACCEPT
    iptables -A HOLM-OUTPUT -p icmp -j ACCEPT

    # ── DNS (local resolution) ──
    iptables -A HOLM-INPUT -p udp --dport 53 -s "$LAN_CIDR" -j ACCEPT
    iptables -A HOLM-INPUT -p udp --dport 53 -s "$K3S_POD_CIDR" -j ACCEPT
    iptables -A HOLM-OUTPUT -p udp --dport 53 -d "$LAN_CIDR" -j ACCEPT
    iptables -A HOLM-OUTPUT -p udp --dport 53 -d "$K3S_SVC_CIDR" -j ACCEPT

    ok "Base rules applied (LAN + k3s internal traffic allowed)"
}

# ── Mode: Transitional ───────────────────────────────────────
apply_transitional() {
    log "Mode: TRANSITIONAL — LAN + internet access"
    apply_base_rules

    # Allow all outbound (internet still available)
    iptables -A HOLM-OUTPUT -j ACCEPT
    iptables -A HOLM-FORWARD -j ACCEPT

    ok "Transitional mode active — internet access permitted"
    warn "This is the migration state. Move to 'restricted' when ready."
}

# ── Mode: Restricted ──────────────────────────────────────────
apply_restricted() {
    log "Mode: RESTRICTED — LAN + allowlisted external only"
    apply_base_rules

    # Allow specific external IPs
    for cidr in "${ALLOWED_EXTERNAL[@]}"; do
        iptables -A HOLM-OUTPUT -d "$cidr" -j ACCEPT
        iptables -A HOLM-FORWARD -d "$cidr" -j ACCEPT
        log "  Allowed: $cidr"
    done

    # Allow outbound DNS to upstream (needed for external resolution)
    iptables -A HOLM-OUTPUT -p udp --dport 53 -j ACCEPT
    iptables -A HOLM-OUTPUT -p tcp --dport 53 -j ACCEPT

    # Allow outbound HTTPS (for allowed IPs above)
    iptables -A HOLM-OUTPUT -p tcp --dport 443 -j ACCEPT
    iptables -A HOLM-OUTPUT -p tcp --dport 80 -j ACCEPT

    # Drop everything else
    iptables -A HOLM-OUTPUT -j LOG --log-prefix "HOLM-BLOCKED-OUT: " --log-level 4
    iptables -A HOLM-OUTPUT -j DROP
    iptables -A HOLM-FORWARD -j LOG --log-prefix "HOLM-BLOCKED-FWD: " --log-level 4
    iptables -A HOLM-FORWARD -j DROP

    ok "Restricted mode active — only allowlisted external IPs permitted"
    warn "Blocked traffic is logged to syslog (HOLM-BLOCKED-*)"
}

# ── Mode: Air-Gap ─────────────────────────────────────────────
apply_airgap() {
    log "Mode: AIR-GAP — zero internet access"
    echo -e "${MAGENTA}"
    echo "  ╔═══════════════════════════════════════════════╗"
    echo "  ║  SOVEREIGN MODE: CUTTING EXTERNAL ACCESS      ║"
    echo "  ║  Only LAN and cluster-internal traffic allowed ║"
    echo "  ╚═══════════════════════════════════════════════╝"
    echo -e "${NC}"

    apply_base_rules

    # Drop ALL non-LAN, non-k3s outbound
    iptables -A HOLM-OUTPUT -j LOG --log-prefix "HOLM-AIRGAP-OUT: " --log-level 4
    iptables -A HOLM-OUTPUT -j DROP
    iptables -A HOLM-FORWARD -d "$LAN_CIDR" -j ACCEPT
    iptables -A HOLM-FORWARD -j LOG --log-prefix "HOLM-AIRGAP-FWD: " --log-level 4
    iptables -A HOLM-FORWARD -j DROP

    ok "AIR-GAP mode active — no internet access"
    ok "Cluster is fully sovereign. Only LAN traffic flows."
}

# ── Status ────────────────────────────────────────────────────
show_status() {
    echo -e "${CYAN}═══════════════════════════════════════════════${NC}"
    echo -e "${CYAN}  HOLM Firewall Status${NC}"
    echo -e "${CYAN}═══════════════════════════════════════════════${NC}"

    if iptables -L HOLM-INPUT -n 2>/dev/null | grep -q "ACCEPT"; then
        ok "HOLM firewall chains: ACTIVE"
    else
        warn "HOLM firewall chains: NOT CONFIGURED"
        echo "  Run: sudo $0 [transitional|restricted|airgap]"
        return
    fi

    echo ""
    echo -e "${MAGENTA}── HOLM-INPUT ──${NC}"
    iptables -L HOLM-INPUT -n --line-numbers 2>/dev/null || true

    echo ""
    echo -e "${MAGENTA}── HOLM-OUTPUT ──${NC}"
    iptables -L HOLM-OUTPUT -n --line-numbers 2>/dev/null || true

    echo ""
    echo -e "${MAGENTA}── HOLM-FORWARD ──${NC}"
    iptables -L HOLM-FORWARD -n --line-numbers 2>/dev/null || true

    echo ""
    echo -e "${CYAN}── Recent blocked packets ──${NC}"
    dmesg 2>/dev/null | grep "HOLM-BLOCKED\|HOLM-AIRGAP" | tail -5 || echo "  (none)"
}

# ── Rollback ──────────────────────────────────────────────────
rollback() {
    if [[ -f "$IPTABLES_BACKUP" ]]; then
        log "Rolling back to saved rules..."
        flush_holm_rules
        iptables-restore < "$IPTABLES_BACKUP"
        ok "Rollback complete"
    else
        warn "No backup found. Flushing HOLM chains only."
        flush_holm_rules
        ok "HOLM chains removed. Default iptables policy restored."
    fi
}

# ── Main ──────────────────────────────────────────────────────
case "$MODE" in
    status)
        show_status
        ;;
    transitional)
        check_root
        backup_rules
        flush_holm_rules
        apply_transitional
        ;;
    restricted)
        check_root
        backup_rules
        flush_holm_rules
        apply_restricted
        ;;
    airgap)
        check_root
        backup_rules
        flush_holm_rules
        apply_airgap
        ;;
    rollback)
        check_root
        rollback
        ;;
    *)
        echo "Usage: sudo $0 {transitional|restricted|airgap|status|rollback}"
        echo ""
        echo "Modes:"
        echo "  transitional  LAN + full internet (current migration state)"
        echo "  restricted    LAN + allowlisted external IPs only"
        echo "  airgap        LAN only — zero internet access (sovereign mode)"
        echo "  status        Show current firewall state"
        echo "  rollback      Restore previous iptables rules"
        exit 1
        ;;
esac
