#!/usr/bin/env bash
#
# deploy.sh -- Deploy the holm.chat static documentation site to a Raspberry Pi
#              Kubernetes cluster and configure Cloudflare DNS.
#
# What this script does, step by step:
#
#   1. Prompts for connection details (Pi SSH target, auth method, Cloudflare
#      credentials, Pi public IP). Nothing is hardcoded.
#   2. Sets up SSH key authentication to the Pi. If a password is provided,
#      it uses ssh-copy-id to push a local key, then switches to key-based
#      auth for all subsequent commands.
#   3. Verifies that kubectl is reachable on the Pi (expects k3s or similar).
#   4. Copies the Kubernetes manifests (deploy/) to the Pi and applies them
#      in order: namespace, configmap, deployment (+ PVC), service, ingress.
#   5. Copies the built site (site/) into the PersistentVolume via a
#      temporary content-loader pod.
#   6. Creates or updates a Cloudflare DNS A record for holm.chat pointing
#      to the Pi's public IP.
#   7. Runs a basic reachability check against the deployed site.
#
# Safety:
#   - Idempotent: every step uses "create or update" semantics. Safe to
#     re-run at any point.
#   - Confirms before each destructive or mutating step.
#   - Prints rollback instructions if anything fails.
#   - No credentials are stored on disk or in the script.
#
# Usage:
#   ./deploy.sh                     # interactive, prompts for everything
#   ./deploy.sh --help              # show this header
#
# Prerequisites:
#   - ssh, ssh-keygen, ssh-copy-id, scp, rsync (standard on macOS/Linux)
#   - curl and jq (for Cloudflare API)
#   - The Pi must be running k3s (or another distro that provides kubectl)
#
# ============================================================================

set -euo pipefail

# ---------------------------------------------------------------------------
# Colors and helpers
# ---------------------------------------------------------------------------
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

info()    { printf "${CYAN}[INFO]${NC}  %s\n" "$*"; }
success() { printf "${GREEN}[OK]${NC}    %s\n" "$*"; }
warn()    { printf "${YELLOW}[WARN]${NC}  %s\n" "$*"; }
error()   { printf "${RED}[ERROR]${NC} %s\n" "$*" >&2; }
step()    { printf "\n${BOLD}==> Step %s${NC}\n" "$*"; }

confirm() {
    local prompt="${1:-Continue?}"
    printf "${YELLOW}%s [y/N]: ${NC}" "$prompt"
    read -r answer
    case "$answer" in
        [yY]|[yY][eE][sS]) return 0 ;;
        *) return 1 ;;
    esac
}

fail_with_rollback() {
    local msg="$1"
    local rollback="${2:-}"
    error "$msg"
    if [[ -n "$rollback" ]]; then
        printf "\n${RED}--- Rollback instructions ---${NC}\n"
        printf "%s\n" "$rollback"
        printf "${RED}-----------------------------${NC}\n"
    fi
    exit 1
}

# ---------------------------------------------------------------------------
# Show help
# ---------------------------------------------------------------------------
if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
    sed -n '2,/^# ====/p' "$0" | sed 's/^# \?//'
    exit 0
fi

# ---------------------------------------------------------------------------
# Pre-flight: verify local tools
# ---------------------------------------------------------------------------
for cmd in ssh scp curl jq; do
    if ! command -v "$cmd" &>/dev/null; then
        fail_with_rollback "Required command '$cmd' not found. Please install it first."
    fi
done

# ---------------------------------------------------------------------------
# Directory paths (relative to this script)
# ---------------------------------------------------------------------------
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="${SCRIPT_DIR}/deploy"
SITE_DIR="${SCRIPT_DIR}/site"

# ---------------------------------------------------------------------------
# 0. Gather inputs interactively
# ---------------------------------------------------------------------------
step "0: Gather deployment configuration"

printf "${BOLD}Pi SSH target${NC} (user@ip) [rpi1@192.168.8.197]: "
read -r PI_TARGET
PI_TARGET="${PI_TARGET:-rpi1@192.168.8.197}"
PI_USER="${PI_TARGET%%@*}"
PI_IP="${PI_TARGET##*@}"

printf "\n${BOLD}Authentication method:${NC}\n"
printf "  1) SSH key (I already have key access)\n"
printf "  2) SSH key at a specific path\n"
printf "  3) Password (will set up key auth via ssh-copy-id)\n"
printf "Choose [1/2/3]: "
read -r AUTH_CHOICE
AUTH_CHOICE="${AUTH_CHOICE:-1}"

SSH_KEY_PATH=""
SSH_PASSWORD=""

case "$AUTH_CHOICE" in
    1)
        info "Using default SSH key (~/.ssh/id_* or agent)"
        SSH_OPTS="-o StrictHostKeyChecking=accept-new -o ConnectTimeout=10"
        ;;
    2)
        printf "Path to SSH private key: "
        read -r SSH_KEY_PATH
        if [[ ! -f "$SSH_KEY_PATH" ]]; then
            fail_with_rollback "Key file not found: $SSH_KEY_PATH"
        fi
        SSH_OPTS="-i ${SSH_KEY_PATH} -o StrictHostKeyChecking=accept-new -o ConnectTimeout=10"
        ;;
    3)
        printf "SSH password for ${PI_USER}@${PI_IP}: "
        read -rs SSH_PASSWORD
        printf "\n"
        if [[ -z "$SSH_PASSWORD" ]]; then
            fail_with_rollback "Password cannot be empty."
        fi
        # We will set up key auth first, then use the key for everything else
        SSH_OPTS="-o StrictHostKeyChecking=accept-new -o ConnectTimeout=10"
        ;;
    *)
        fail_with_rollback "Invalid choice: $AUTH_CHOICE"
        ;;
esac

printf "\n${BOLD}Cloudflare API token${NC} (Zone:DNS:Edit permission): "
read -rs CF_API_TOKEN
printf "\n"
if [[ -z "$CF_API_TOKEN" ]]; then
    fail_with_rollback "Cloudflare API token is required."
fi

printf "${BOLD}Cloudflare Zone ID${NC} for holm.chat: "
read -r CF_ZONE_ID
if [[ -z "$CF_ZONE_ID" ]]; then
    fail_with_rollback "Cloudflare Zone ID is required."
fi

printf "${BOLD}Pi's public IP${NC} (for the DNS A record): "
read -r PI_PUBLIC_IP
if [[ -z "$PI_PUBLIC_IP" ]]; then
    fail_with_rollback "Public IP is required for the DNS A record."
fi

# Summary
printf "\n${BOLD}--- Configuration Summary ---${NC}\n"
printf "  Pi target:    %s\n" "$PI_TARGET"
printf "  Auth method:  %s\n" "$(case $AUTH_CHOICE in 1) echo 'default key';; 2) echo "key: $SSH_KEY_PATH";; 3) echo 'password (will set up key)';; esac)"
printf "  Zone ID:      %s\n" "$CF_ZONE_ID"
printf "  Public IP:    %s\n" "$PI_PUBLIC_IP"
printf "  Manifests:    %s\n" "$DEPLOY_DIR"
printf "  Site content: %s\n" "$SITE_DIR"
printf "${BOLD}-----------------------------${NC}\n"

confirm "Proceed with deployment?" || { info "Aborted."; exit 0; }

# ---------------------------------------------------------------------------
# Helper: run a command on the Pi via SSH
# ---------------------------------------------------------------------------
remote() {
    # shellcheck disable=SC2086
    ssh $SSH_OPTS "${PI_TARGET}" "$@"
}

remote_sudo() {
    # shellcheck disable=SC2086
    ssh $SSH_OPTS "${PI_TARGET}" "sudo $*"
}

# ---------------------------------------------------------------------------
# 1. SSH key setup (if password auth was chosen)
# ---------------------------------------------------------------------------
step "1: SSH key authentication"

if [[ "$AUTH_CHOICE" == "3" ]]; then
    # Generate a dedicated key pair if one does not already exist
    PI_KEY_PATH="${HOME}/.ssh/id_ed25519_pi_cluster"
    if [[ ! -f "$PI_KEY_PATH" ]]; then
        info "Generating a new Ed25519 key pair at ${PI_KEY_PATH}"
        ssh-keygen -t ed25519 -C "deploy-to-pi" -f "$PI_KEY_PATH" -N "" -q
        success "Key pair created."
    else
        info "Key pair already exists at ${PI_KEY_PATH}. Reusing it."
    fi

    info "Copying public key to ${PI_TARGET} via ssh-copy-id..."
    if ! command -v sshpass &>/dev/null; then
        warn "sshpass is not installed. You will be prompted for the password by ssh-copy-id."
        ssh-copy-id -i "${PI_KEY_PATH}.pub" -o StrictHostKeyChecking=accept-new "$PI_TARGET"
    else
        sshpass -p "$SSH_PASSWORD" ssh-copy-id -i "${PI_KEY_PATH}.pub" -o StrictHostKeyChecking=accept-new "$PI_TARGET"
    fi

    # Switch to key-based auth for the rest of the script
    SSH_KEY_PATH="$PI_KEY_PATH"
    SSH_OPTS="-i ${SSH_KEY_PATH} -o StrictHostKeyChecking=accept-new -o ConnectTimeout=10"
    # Clear the password from memory
    SSH_PASSWORD=""

    success "SSH key authentication configured."
else
    info "Skipping key setup (already using key-based auth)."
fi

# Verify SSH connectivity
info "Testing SSH connection to ${PI_TARGET}..."
if ! remote "echo 'SSH connection OK'" 2>/dev/null; then
    fail_with_rollback \
        "Cannot SSH into ${PI_TARGET}." \
        "1. Verify the Pi is powered on and reachable at ${PI_IP}.
2. Check that your SSH key is authorized: ssh ${SSH_OPTS} ${PI_TARGET}
3. If using a password, ensure sshpass is installed or enter it when prompted."
fi
success "SSH connection verified."

# ---------------------------------------------------------------------------
# 2. Verify kubectl on the Pi
# ---------------------------------------------------------------------------
step "2: Verify kubectl on the Pi"

if ! remote "command -v kubectl" &>/dev/null; then
    # k3s installs kubectl at a non-standard path; try the k3s wrapper
    if remote "command -v k3s" &>/dev/null; then
        info "kubectl not in PATH, but k3s is available. Using 'k3s kubectl'."
        KUBECTL_CMD="sudo k3s kubectl"
    else
        fail_with_rollback \
            "Neither kubectl nor k3s found on the Pi." \
            "Install k3s on the Pi:
  curl -sfL https://get.k3s.io | sh -
Then re-run this script."
    fi
else
    KUBECTL_CMD="sudo kubectl"
fi

# Quick cluster health check
info "Checking cluster health..."
if ! remote "${KUBECTL_CMD} get nodes" 2>/dev/null; then
    fail_with_rollback \
        "kubectl is present but cannot reach the cluster." \
        "Verify k3s is running:
  sudo systemctl status k3s
  sudo k3s kubectl get nodes"
fi
success "Kubernetes cluster is reachable."

# ---------------------------------------------------------------------------
# 3. Verify site/ has content
# ---------------------------------------------------------------------------
step "3: Verify site content"

if [[ ! -d "$SITE_DIR" ]]; then
    warn "The site/ directory does not exist at ${SITE_DIR}."
    warn "Kubernetes manifests will still be applied, but no site content will be uploaded."
    SKIP_SITE_CONTENT=true
elif [[ -z "$(ls -A "$SITE_DIR" 2>/dev/null)" ]]; then
    warn "The site/ directory is empty."
    warn "Kubernetes manifests will still be applied, but no site content will be uploaded."
    SKIP_SITE_CONTENT=true
else
    FILE_COUNT=$(find "$SITE_DIR" -type f | wc -l | tr -d ' ')
    success "site/ contains ${FILE_COUNT} file(s)."
    SKIP_SITE_CONTENT=false
fi

# ---------------------------------------------------------------------------
# 4. Copy manifests to the Pi and apply them
# ---------------------------------------------------------------------------
step "4: Apply Kubernetes manifests"

if [[ ! -d "$DEPLOY_DIR" ]]; then
    fail_with_rollback \
        "Manifest directory not found: ${DEPLOY_DIR}" \
        "Ensure the deploy/ directory exists with namespace.yaml, configmap.yaml, etc."
fi

if confirm "Apply Kubernetes manifests to the Pi cluster?"; then
    REMOTE_MANIFEST_DIR="/tmp/docs-k8s-manifests"

    info "Copying manifests to ${PI_TARGET}:${REMOTE_MANIFEST_DIR}..."
    remote "mkdir -p ${REMOTE_MANIFEST_DIR}"
    # shellcheck disable=SC2086
    scp $SSH_OPTS \
        "${DEPLOY_DIR}/namespace.yaml" \
        "${DEPLOY_DIR}/configmap.yaml" \
        "${DEPLOY_DIR}/deployment.yaml" \
        "${DEPLOY_DIR}/service.yaml" \
        "${DEPLOY_DIR}/ingress.yaml" \
        "${PI_TARGET}:${REMOTE_MANIFEST_DIR}/"
    success "Manifests copied."

    # Apply in dependency order
    MANIFEST_ORDER="namespace.yaml configmap.yaml deployment.yaml service.yaml ingress.yaml"
    for manifest in $MANIFEST_ORDER; do
        info "Applying ${manifest}..."
        if ! remote "${KUBECTL_CMD} apply -f ${REMOTE_MANIFEST_DIR}/${manifest}"; then
            fail_with_rollback \
                "Failed to apply ${manifest}." \
                "1. Check the manifest for syntax errors:
     cat ${DEPLOY_DIR}/${manifest}
2. Check cluster state:
     ssh ${PI_TARGET} '${KUBECTL_CMD} get all -n docs'
3. To undo everything:
     ssh ${PI_TARGET} '${KUBECTL_CMD} delete namespace docs'"
        fi
        success "${manifest} applied."
    done

    # Wait for the deployment to be ready
    info "Waiting for docs-site deployment to become ready (timeout 120s)..."
    if ! remote "${KUBECTL_CMD} -n docs rollout status deployment/docs-site --timeout=120s" 2>/dev/null; then
        warn "Deployment did not become ready within 120s. Continuing anyway."
        warn "Check pod status: ssh ${PI_TARGET} '${KUBECTL_CMD} -n docs get pods'"
    else
        success "Deployment is ready."
    fi

    # Clean up remote manifests
    remote "rm -rf ${REMOTE_MANIFEST_DIR}" 2>/dev/null || true
else
    info "Skipping manifest application."
fi

# ---------------------------------------------------------------------------
# 5. Copy site content into the PersistentVolume
# ---------------------------------------------------------------------------
step "5: Upload site content to the cluster"

if [[ "$SKIP_SITE_CONTENT" == "true" ]]; then
    warn "Skipping site content upload (no content in site/)."
else
    if confirm "Upload site/ contents to the PersistentVolume?"; then
        info "Creating temporary content-loader pod..."

        # Delete any leftover content-loader pod from a previous run (idempotent)
        remote "${KUBECTL_CMD} -n docs delete pod content-loader --ignore-not-found=true" 2>/dev/null || true

        # Create the loader pod with write access to the PVC
        LOADER_OVERRIDES='{
            "spec": {
                "containers": [{
                    "name": "loader",
                    "image": "busybox:latest",
                    "command": ["sleep", "600"],
                    "volumeMounts": [{
                        "name": "content",
                        "mountPath": "/site"
                    }]
                }],
                "volumes": [{
                    "name": "content",
                    "persistentVolumeClaim": {
                        "claimName": "docs-content"
                    }
                }]
            }
        }'

        remote "${KUBECTL_CMD} -n docs run content-loader \
            --image=busybox:latest \
            --restart=Never \
            --overrides='${LOADER_OVERRIDES}'"

        info "Waiting for content-loader pod to be ready..."
        if ! remote "${KUBECTL_CMD} -n docs wait --for=condition=Ready pod/content-loader --timeout=90s"; then
            fail_with_rollback \
                "content-loader pod did not become ready." \
                "1. Check pod status:
     ssh ${PI_TARGET} '${KUBECTL_CMD} -n docs describe pod content-loader'
2. Clean up:
     ssh ${PI_TARGET} '${KUBECTL_CMD} -n docs delete pod content-loader'"
        fi
        success "content-loader pod is ready."

        # Clear old content
        info "Clearing old site content from the volume..."
        remote "${KUBECTL_CMD} -n docs exec content-loader -- sh -c 'rm -rf /site/*'"

        # Tar the site locally, pipe through SSH, extract into the pod
        # This is more reliable than kubectl cp for large directory trees.
        info "Uploading site content (this may take a moment)..."

        # First, tar the site content to a temp file on the Pi
        REMOTE_TAR="/tmp/docs-site-content.tar.gz"
        # shellcheck disable=SC2086
        tar -czf - -C "$SITE_DIR" . | ssh $SSH_OPTS "$PI_TARGET" "cat > ${REMOTE_TAR}"

        # Then copy it into the pod and extract
        remote "${KUBECTL_CMD} -n docs cp ${REMOTE_TAR} content-loader:/tmp/site-content.tar.gz"
        remote "${KUBECTL_CMD} -n docs exec content-loader -- tar -xzf /tmp/site-content.tar.gz -C /site"
        remote "${KUBECTL_CMD} -n docs exec content-loader -- rm /tmp/site-content.tar.gz"

        # Clean up remote tar
        remote "rm -f ${REMOTE_TAR}" 2>/dev/null || true

        # Verify
        info "Verifying uploaded content..."
        remote "${KUBECTL_CMD} -n docs exec content-loader -- ls -la /site/"

        # Clean up the loader pod
        info "Removing content-loader pod..."
        remote "${KUBECTL_CMD} -n docs delete pod content-loader --wait=false" 2>/dev/null || true

        # Restart the deployment so nginx picks up the new files
        info "Restarting docs-site deployment to serve fresh content..."
        remote "${KUBECTL_CMD} -n docs rollout restart deployment/docs-site"
        remote "${KUBECTL_CMD} -n docs rollout status deployment/docs-site --timeout=120s" 2>/dev/null || true

        success "Site content deployed."
    else
        info "Skipping content upload."
    fi
fi

# ---------------------------------------------------------------------------
# 6. Configure Cloudflare DNS
# ---------------------------------------------------------------------------
step "6: Configure Cloudflare DNS A record for holm.chat"

if confirm "Create/update Cloudflare DNS A record (holm.chat -> ${PI_PUBLIC_IP})?"; then
    CF_API="https://api.cloudflare.com/client/v4"

    # Verify the API token works
    info "Verifying Cloudflare API token..."
    CF_VERIFY=$(curl -s -X GET "${CF_API}/user/tokens/verify" \
        -H "Authorization: Bearer ${CF_API_TOKEN}" \
        -H "Content-Type: application/json")

    if [[ "$(echo "$CF_VERIFY" | jq -r '.success')" != "true" ]]; then
        fail_with_rollback \
            "Cloudflare API token verification failed." \
            "Ensure your token has Zone:DNS:Edit permissions.
Response: $(echo "$CF_VERIFY" | jq -r '.errors')"
    fi
    success "Cloudflare API token is valid."

    # Check if an A record for holm.chat already exists
    info "Checking for existing A record..."
    EXISTING_RECORDS=$(curl -s -X GET \
        "${CF_API}/zones/${CF_ZONE_ID}/dns_records?type=A&name=holm.chat" \
        -H "Authorization: Bearer ${CF_API_TOKEN}" \
        -H "Content-Type: application/json")

    RECORD_COUNT=$(echo "$EXISTING_RECORDS" | jq -r '.result | length')

    if [[ "$RECORD_COUNT" -gt 0 ]]; then
        EXISTING_ID=$(echo "$EXISTING_RECORDS" | jq -r '.result[0].id')
        EXISTING_IP=$(echo "$EXISTING_RECORDS" | jq -r '.result[0].content')
        info "Existing A record found (ID: ${EXISTING_ID}, IP: ${EXISTING_IP})."

        if [[ "$EXISTING_IP" == "$PI_PUBLIC_IP" ]]; then
            success "A record already points to ${PI_PUBLIC_IP}. No update needed."
        else
            info "Updating A record from ${EXISTING_IP} to ${PI_PUBLIC_IP}..."
            UPDATE_RESULT=$(curl -s -X PUT \
                "${CF_API}/zones/${CF_ZONE_ID}/dns_records/${EXISTING_ID}" \
                -H "Authorization: Bearer ${CF_API_TOKEN}" \
                -H "Content-Type: application/json" \
                --data "{
                    \"type\": \"A\",
                    \"name\": \"holm.chat\",
                    \"content\": \"${PI_PUBLIC_IP}\",
                    \"ttl\": 1,
                    \"proxied\": true
                }")

            if [[ "$(echo "$UPDATE_RESULT" | jq -r '.success')" != "true" ]]; then
                fail_with_rollback \
                    "Failed to update DNS record." \
                    "Cloudflare error: $(echo "$UPDATE_RESULT" | jq -r '.errors')
You can update the record manually in the Cloudflare dashboard:
  Zone: holm.chat -> DNS -> A record -> change content to ${PI_PUBLIC_IP}"
            fi
            success "DNS A record updated to ${PI_PUBLIC_IP}."
        fi
    else
        info "No existing A record found. Creating one..."
        CREATE_RESULT=$(curl -s -X POST \
            "${CF_API}/zones/${CF_ZONE_ID}/dns_records" \
            -H "Authorization: Bearer ${CF_API_TOKEN}" \
            -H "Content-Type: application/json" \
            --data "{
                \"type\": \"A\",
                \"name\": \"holm.chat\",
                \"content\": \"${PI_PUBLIC_IP}\",
                \"ttl\": 1,
                \"proxied\": true
            }")

        if [[ "$(echo "$CREATE_RESULT" | jq -r '.success')" != "true" ]]; then
            fail_with_rollback \
                "Failed to create DNS record." \
                "Cloudflare error: $(echo "$CREATE_RESULT" | jq -r '.errors')
You can create the record manually in the Cloudflare dashboard:
  Zone: holm.chat -> DNS -> Add record -> Type: A, Name: @, Content: ${PI_PUBLIC_IP}, Proxy: On"
        fi
        success "DNS A record created: holm.chat -> ${PI_PUBLIC_IP} (proxied)."
    fi

    info "Note: Cloudflare proxy is enabled (orange cloud). SSL/TLS mode should"
    info "be set to 'Full' in the Cloudflare dashboard for holm.chat."
else
    info "Skipping DNS configuration."
fi

# ---------------------------------------------------------------------------
# 7. Verify the site is reachable
# ---------------------------------------------------------------------------
step "7: Verify site reachability"

info "Waiting a few seconds for DNS propagation and pod readiness..."
sleep 5

# Check the health endpoint via the Pi's internal service first
info "Checking internal health endpoint on the Pi..."
if remote "${KUBECTL_CMD} -n docs exec deploy/docs-site -- curl -sf http://localhost/healthz" 2>/dev/null; then
    success "Internal health check passed (/healthz returns OK)."
else
    warn "Internal health check did not respond. The pod may still be starting."
fi

# Check via the public URL (may take time for DNS to propagate)
info "Checking https://holm.chat ..."
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" --max-time 10 "https://holm.chat" 2>/dev/null || echo "000")

if [[ "$HTTP_STATUS" == "200" ]]; then
    success "https://holm.chat is returning HTTP 200. Site is live!"
elif [[ "$HTTP_STATUS" == "000" ]]; then
    warn "Could not connect to https://holm.chat (DNS may not have propagated yet)."
    info "This is normal if you just created the DNS record."
    info "Try again in a few minutes: curl -I https://holm.chat"
else
    warn "https://holm.chat returned HTTP ${HTTP_STATUS}."
    info "The site may still be starting. Check:"
    info "  curl -I https://holm.chat"
    info "  ssh ${PI_TARGET} '${KUBECTL_CMD} -n docs get pods'"
fi

# ---------------------------------------------------------------------------
# Done
# ---------------------------------------------------------------------------
printf "\n${GREEN}${BOLD}========================================${NC}\n"
printf "${GREEN}${BOLD}  Deployment complete!${NC}\n"
printf "${GREEN}${BOLD}========================================${NC}\n\n"

printf "${BOLD}Useful commands:${NC}\n"
printf "  Check pods:       ssh %s '%s -n docs get pods'\n" "$PI_TARGET" "$KUBECTL_CMD"
printf "  View logs:        ssh %s '%s -n docs logs -l app.kubernetes.io/name=docs-site --tail=50'\n" "$PI_TARGET" "$KUBECTL_CMD"
printf "  Describe ingress: ssh %s '%s -n docs describe ingress docs-site'\n" "$PI_TARGET" "$KUBECTL_CMD"
printf "  Port-forward:     ssh -L 8080:localhost:8080 %s '%s -n docs port-forward svc/docs-site 8080:80'\n" "$PI_TARGET" "$KUBECTL_CMD"
printf "  DNS check:        dig holm.chat\n"
printf "  Full rollback:    ssh %s '%s delete namespace docs'\n\n" "$PI_TARGET" "$KUBECTL_CMD"
