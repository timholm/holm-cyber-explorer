# Deploying holm.chat Documentation Site

Kubernetes deployment for a static documentation site running on a Raspberry Pi cluster (ARM64).

## Prerequisites

- Raspberry Pi 4 (or newer) running a lightweight K8s distribution (k3s recommended)
- `kubectl` installed on your local machine
- Domain `holm.chat` managed through Cloudflare
- SSH access to the Pi node

---

## 1. SSH Key Setup

Replace password authentication with SSH keys for the Pi node.

```bash
# Generate a key pair (if you don't have one)
ssh-keygen -t ed25519 -C "YOUR_EMAIL" -f ~/.ssh/id_pi_cluster

# Copy public key to the Pi
ssh-copy-id -i ~/.ssh/id_pi_cluster.pub YOUR_PI_USER@YOUR_PI_IP

# Disable password auth on the Pi (edit /etc/ssh/sshd_config)
#   PasswordAuthentication no
#   PubkeyAuthentication yes
# Then restart sshd:
#   sudo systemctl restart sshd
```

Add to `~/.ssh/config`:

```
Host pi-cluster
    HostName YOUR_PI_IP
    User YOUR_PI_USER
    IdentityFile ~/.ssh/id_pi_cluster
    Port 22
```

---

## 2. kubectl Context Configuration

Copy the kubeconfig from the Pi and set it as your active context.

```bash
# Copy kubeconfig from k3s
scp pi-cluster:/etc/rancher/k3s/k3s.yaml ~/.kube/pi-cluster.yaml

# Edit the file: replace 127.0.0.1 with YOUR_PI_IP
sed -i '' 's/127.0.0.1/YOUR_PI_IP/' ~/.kube/pi-cluster.yaml

# Merge into your kubeconfig (or use it standalone)
export KUBECONFIG=~/.kube/pi-cluster.yaml

# Verify connectivity
kubectl get nodes
```

---

## 3. Cloudflare DNS Setup

Configure DNS in the Cloudflare dashboard (do NOT store API keys in these manifests).

1. Log in to the Cloudflare dashboard
2. Select the `holm.chat` zone
3. Add a DNS record:
   - **Type:** A
   - **Name:** @ (root domain)
   - **Content:** YOUR_PI_PUBLIC_IP
   - **Proxy status:** Proxied (orange cloud)
4. SSL/TLS settings:
   - Set encryption mode to **Full**
   - This means Cloudflare handles the public HTTPS certificate
   - Traffic from Cloudflare to your cluster is over HTTP (port 80)

### Alternative: Let's Encrypt (no Cloudflare proxy)

If you prefer TLS termination at the cluster instead of Cloudflare proxy:

1. Install cert-manager on the cluster:
   ```bash
   kubectl apply -f https://github.com/cert-manager/cert-manager/releases/latest/download/cert-manager.yaml
   ```
2. Create a ClusterIssuer:
   ```yaml
   apiVersion: cert-manager.io/v1
   kind: ClusterIssuer
   metadata:
     name: letsencrypt-prod
   spec:
     acme:
       server: https://acme-v02.api.letsencrypt.org/directory
       email: YOUR_EMAIL
       privateKeySecretRef:
         name: letsencrypt-prod-key
       solvers:
         - http01:
             ingress:
               class: nginx
   ```
3. Uncomment the cert-manager annotations and `tls` block in `ingress.yaml`
4. Set Cloudflare DNS proxy status to **DNS only** (grey cloud) so the ACME challenge reaches your cluster

### Cloudflare API Token (for external-dns or cert-manager DNS01, if needed later)

If you ever automate DNS record management, you will need a scoped API token. **Never commit it to files.** Store it as a Kubernetes secret:

```bash
kubectl create secret generic cloudflare-api-token \
  --namespace docs \
  --from-literal=api-token=YOUR_CF_API_TOKEN
```

Required token permissions: `Zone:DNS:Edit` for the `holm.chat` zone.

---

## 4. Deploy the Manifests

Apply the manifests in order:

```bash
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f ingress.yaml
```

Or apply everything at once:

```bash
kubectl apply -f .
```

Verify the deployment:

```bash
kubectl -n docs get all
kubectl -n docs describe ingress docs-site
```

---

## 5. Content Publishing Workflow

Content is published manually by copying built HTML files into the PersistentVolume.

### Build your site locally

Use any static site generator (or plain HTML). The output should be a directory of HTML, CSS, and asset files with an `index.html` at the root.

### Copy content to the cluster

**Option A: kubectl cp (simplest)**

```bash
# Find the running pod name
POD=$(kubectl -n docs get pod -l app.kubernetes.io/name=docs-site -o jsonpath='{.items[0].metadata.name}')

# Copy built site content into the pod's volume
# Note: the volume is mounted read-only in the container, so we
# need to copy via a temporary pod or use Option B instead.

# Create a temporary content-loader pod with write access:
kubectl -n docs run content-loader \
  --image=busybox:latest \
  --restart=Never \
  --overrides='{
    "spec": {
      "containers": [{
        "name": "loader",
        "image": "busybox:latest",
        "command": ["sleep", "3600"],
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

# Wait for the loader pod to be ready
kubectl -n docs wait --for=condition=Ready pod/content-loader --timeout=60s

# Clear old content and copy new content
kubectl -n docs exec content-loader -- rm -rf /site/*
kubectl cp ./public/. docs/content-loader:/site/

# Verify the files landed
kubectl -n docs exec content-loader -- ls -la /site/

# Clean up the loader pod
kubectl -n docs delete pod content-loader

# Restart the docs pod to pick up fresh content
kubectl -n docs rollout restart deployment/docs-site
```

**Option B: rsync to the node, then mount as hostPath**

For a single-node cluster, you can rsync content directly to the Pi and use a `hostPath` volume instead of a PVC. Edit `deployment.yaml` to replace the PVC volume with:

```yaml
volumes:
  - name: site-content
    hostPath:
      path: /opt/docs-site/content
      type: DirectoryOrCreate
```

Then publish with:

```bash
rsync -avz --delete ./public/ pi-cluster:/opt/docs-site/content/
# No pod restart needed — nginx serves directly from the directory
```

---

## Architecture Overview

```
Internet
   |
Cloudflare (HTTPS termination, CDN, DDoS protection)
   |
YOUR_PI_PUBLIC_IP:80
   |
k3s Ingress Controller (nginx)
   |
Service: docs-site (ClusterIP :80)
   |
Pod: nginx:1.27-alpine
   |
Volume: /usr/share/nginx/html (PVC or hostPath)
```

---

## Resource Usage

These manifests are tuned for Raspberry Pi hardware:

| Resource | Request | Limit  |
|----------|---------|--------|
| CPU      | 50m     | 200m   |
| Memory   | 32Mi    | 128Mi  |
| Storage  | 1Gi     | -      |

---

## Troubleshooting

```bash
# Check pod status
kubectl -n docs get pods -o wide

# View nginx logs
kubectl -n docs logs -l app.kubernetes.io/name=docs-site --tail=50

# Check ingress routing
kubectl -n docs describe ingress docs-site

# Test the service directly (port-forward)
kubectl -n docs port-forward svc/docs-site 8080:80
# Then visit http://localhost:8080

# Verify DNS resolution
dig holm.chat
nslookup holm.chat

# Check if Cloudflare proxy is working
curl -I https://holm.chat
# Look for "cf-ray" header in the response
```
