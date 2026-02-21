# Animus Dashboard

**Assassin's Creed Themed Kubernetes Cluster Management System**

A modern web dashboard for managing your Raspberry Pi Kubernetes cluster with an "Animus" theme inspired by Assassin's Creed.

## Features

- **Node-by-node management** - View and manage each node (Memory Core) individually
- **Real-time pod monitoring** - See running sequences (pods) with live status updates
- **Loki log streaming** - Stream logs from your cluster with filtering support
- **Ansible script execution** - Run predefined protocols (playbooks) on your nodes
- **Keycloak authentication** - Secure access with SSO integration

## Architecture

```
┌────────────────────────────────────────────────────────┐
│                    ANIMUS DASHBOARD                     │
│                   (Next.js 14 Frontend)                 │
└────────────────────────────────────────────────────────┘
                           │
                           ▼
┌────────────────────────────────────────────────────────┐
│                   ANIMUS API (Go Backend)               │
│                        Port 8080                        │
└────────────────────────────────────────────────────────┘
        │              │              │              │
        ▼              ▼              ▼              ▼
   ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐
   │ K8s API │   │  Loki   │   │ Ansible │   │   SSH   │
   └─────────┘   └─────────┘   └─────────┘   └─────────┘
```

## Quick Start

### Prerequisites

- Node.js 20+
- Go 1.22+
- Access to a Kubernetes cluster
- Loki for log aggregation (optional)
- Ansible installed on the backend host

### Development

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

**Backend:**
```bash
cd backend
go mod download
go run ./cmd/server
```

### Production Deployment

1. Generate SSH keys for Ansible:
```bash
ssh-keygen -t ed25519 -f animus-key -N ""
```

2. Create Kubernetes secrets:
```bash
kubectl create namespace animus
kubectl create secret generic animus-ssh-key \
  --from-file=id_ed25519=animus-key \
  --from-file=id_ed25519.pub=animus-key.pub \
  -n animus
```

3. Deploy to Kubernetes:
```bash
kubectl apply -f deploy/kubernetes/
```

## Configuration

### Environment Variables

**Backend:**
| Variable | Description | Default |
|----------|-------------|---------|
| PORT | API server port | 8080 |
| LOKI_URL | Loki server URL | http://loki.monitoring:3100 |
| SSH_KEY_PATH | Path to SSH private key | /etc/animus/ssh/id_ed25519 |
| SSH_USER | SSH username | tim |
| KEYCLOAK_URL | Keycloak server URL | (disabled if empty) |

**Frontend:**
| Variable | Description | Default |
|----------|-------------|---------|
| BACKEND_URL | Backend API URL | http://localhost:8080 |
| NEXT_PUBLIC_API_URL | Public API URL | (empty for same-origin) |

## Theme

The dashboard uses an Assassin's Creed "Animus" aesthetic:

| Color | Hex | Usage |
|-------|-----|-------|
| Deep Black | #0a0a0f | Primary background |
| Dark Navy | #12121a | Secondary background |
| Animus Gold | #c9a227 | Primary accent |
| Memory Cyan | #00d4ff | Secondary accent |
| Synchronized | #00ff88 | Success states |
| Desynchronized | #ff3366 | Error states |

## Terminology

| Dashboard Term | Kubernetes Equivalent |
|----------------|----------------------|
| Memory Cores | Nodes |
| Sequences | Pods |
| Memory Stream | Logs |
| Protocols | Ansible Playbooks |
| Synchronization | Health Status |

## License

MIT
