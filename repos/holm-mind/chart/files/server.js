const http = require('http');
const fs = require('fs');
const path = require('path');
const crypto = require('crypto');

const PORT = parseInt(process.env.PORT || '3000');
const AUTH_PASSWORD = process.env.AUTH_PASSWORD || 'changeme';
const TOKEN_SECRET = crypto.randomBytes(32).toString('hex');
const K8S_HOST = process.env.KUBERNETES_SERVICE_HOST || 'kubernetes.default.svc';
const K8S_PORT = process.env.KUBERNETES_SERVICE_PORT || '443';
const CACHE_TTL = 30000;

// Read ServiceAccount token and CA
let saToken = '';
let caCert = null;
try {
  saToken = fs.readFileSync('/var/run/secrets/kubernetes.io/serviceaccount/token', 'utf8').trim();
  caCert = fs.readFileSync('/var/run/secrets/kubernetes.io/serviceaccount/ca.crt');
} catch (e) {
  console.log('No ServiceAccount token found, K8s API calls will fail');
}

// Simple token signing
function signToken(payload) {
  const data = Buffer.from(JSON.stringify(payload)).toString('base64url');
  const sig = crypto.createHmac('sha256', TOKEN_SECRET).update(data).digest('base64url');
  return `${data}.${sig}`;
}

function verifyToken(token) {
  if (!token) return null;
  const [data, sig] = token.split('.');
  if (!data || !sig) return null;
  const expected = crypto.createHmac('sha256', TOKEN_SECRET).update(data).digest('base64url');
  if (sig !== expected) return null;
  try { return JSON.parse(Buffer.from(data, 'base64url').toString()); } catch { return null; }
}

// Kubernetes API client
function k8sRequest(apiPath) {
  return new Promise((resolve, reject) => {
    const options = {
      hostname: K8S_HOST,
      port: K8S_PORT,
      path: apiPath,
      method: 'GET',
      headers: { Authorization: `Bearer ${saToken}` },
      ca: caCert,
      rejectUnauthorized: !!caCert,
    };
    const req = (require('https')).request(options, (res) => {
      let body = '';
      res.on('data', (chunk) => body += chunk);
      res.on('end', () => {
        try { resolve(JSON.parse(body)); } catch { reject(new Error('Invalid JSON from K8s API')); }
      });
    });
    req.on('error', reject);
    req.setTimeout(10000, () => { req.destroy(); reject(new Error('K8s API timeout')); });
    req.end();
  });
}

// Cache
let clusterCache = null;
let cacheTime = 0;

async function getClusterData() {
  if (clusterCache && Date.now() - cacheTime < CACHE_TTL) return clusterCache;

  const [namespaces, deployments, pods] = await Promise.all([
    k8sRequest('/api/v1/namespaces'),
    k8sRequest('/apis/apps/v1/deployments'),
    k8sRequest('/api/v1/pods'),
  ]);

  const nsMap = {};
  for (const ns of (namespaces.items || [])) {
    const name = ns.metadata.name;
    if (name.startsWith('kube-') || name === 'envoy-gateway-system') continue;
    nsMap[name] = { name, deployments: [], podCount: 0, readyCount: 0 };
  }

  for (const dep of (deployments.items || [])) {
    const ns = dep.metadata.namespace;
    if (!nsMap[ns]) continue;
    const ready = dep.status?.readyReplicas || 0;
    const desired = dep.spec?.replicas || 1;
    nsMap[ns].deployments.push({
      name: dep.metadata.name,
      replicas: desired,
      ready,
      image: dep.spec?.template?.spec?.containers?.[0]?.image || 'unknown',
      status: ready >= desired ? 'running' : ready > 0 ? 'degraded' : 'down',
      age: dep.metadata.creationTimestamp,
    });
  }

  for (const pod of (pods.items || [])) {
    const ns = pod.metadata.namespace;
    if (!nsMap[ns]) continue;
    nsMap[ns].podCount++;
    const phase = pod.status?.phase;
    if (phase === 'Running' || phase === 'Succeeded') nsMap[ns].readyCount++;
  }

  // Compute namespace health
  for (const ns of Object.values(nsMap)) {
    const total = ns.deployments.length;
    const healthy = ns.deployments.filter(d => d.status === 'running').length;
    ns.health = total === 0 ? 'green' : healthy === total ? 'green' : healthy > 0 ? 'yellow' : 'red';
    ns.serviceCount = total;
  }

  clusterCache = Object.values(nsMap).sort((a, b) => b.serviceCount - a.serviceCount);
  cacheTime = Date.now();
  return clusterCache;
}

async function getNamespaceDetail(nsName) {
  const [deployments, pods] = await Promise.all([
    k8sRequest(`/apis/apps/v1/namespaces/${nsName}/deployments`),
    k8sRequest(`/api/v1/namespaces/${nsName}/pods`),
  ]);

  const services = [];
  for (const dep of (deployments.items || [])) {
    const ready = dep.status?.readyReplicas || 0;
    const desired = dep.spec?.replicas || 1;
    const depPods = (pods.items || []).filter(p =>
      dep.spec?.selector?.matchLabels &&
      Object.entries(dep.spec.selector.matchLabels).every(([k, v]) => p.metadata?.labels?.[k] === v)
    );
    services.push({
      name: dep.metadata.name,
      replicas: desired,
      ready,
      image: dep.spec?.template?.spec?.containers?.[0]?.image || 'unknown',
      status: ready >= desired ? 'running' : ready > 0 ? 'degraded' : 'down',
      age: dep.metadata.creationTimestamp,
      pods: depPods.map(p => ({
        name: p.metadata.name,
        phase: p.status?.phase || 'Unknown',
        restarts: p.status?.containerStatuses?.[0]?.restartCount || 0,
        node: p.spec?.nodeName || 'unscheduled',
      })),
    });
  }
  return { namespace: nsName, services };
}

// Read index.html
let indexHtml = '';
try {
  indexHtml = fs.readFileSync(path.join(__dirname, 'public', 'index.html'), 'utf8');
} catch {
  try { indexHtml = fs.readFileSync(path.join(__dirname, 'index.html'), 'utf8'); } catch { indexHtml = '<h1>index.html not found</h1>'; }
}

const MANIFEST = JSON.stringify({
  name: 'holm mind',
  short_name: 'mind',
  start_url: '/',
  display: 'standalone',
  background_color: '#0a0a1a',
  theme_color: '#0a0a1a',
  description: 'Kubernetes cluster mind map',
  icons: [{ src: 'data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text y=".9em" font-size="90">🧠</text></svg>', sizes: 'any', type: 'image/svg+xml' }],
});

function parseBody(req) {
  return new Promise((resolve) => {
    let body = '';
    req.on('data', (c) => { body += c; if (body.length > 1e5) req.destroy(); });
    req.on('end', () => { try { resolve(JSON.parse(body)); } catch { resolve({}); } });
  });
}

function send(res, status, data, contentType = 'application/json') {
  res.writeHead(status, {
    'Content-Type': contentType,
    'Cache-Control': 'no-store',
    'Access-Control-Allow-Origin': '*',
  });
  res.end(typeof data === 'string' ? data : JSON.stringify(data));
}

function authCheck(req) {
  const auth = req.headers.authorization;
  if (!auth?.startsWith('Bearer ')) return null;
  return verifyToken(auth.slice(7));
}

const server = http.createServer(async (req, res) => {
  const url = new URL(req.url, `http://${req.headers.host}`);
  const pathname = url.pathname;

  // CORS preflight
  if (req.method === 'OPTIONS') {
    res.writeHead(204, {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'GET, POST, OPTIONS',
      'Access-Control-Allow-Headers': 'Content-Type, Authorization',
    });
    return res.end();
  }

  try {
    // Static routes
    if (req.method === 'GET' && (pathname === '/' || pathname === '/index.html')) {
      return send(res, 200, indexHtml, 'text/html');
    }
    if (req.method === 'GET' && pathname === '/manifest.json') {
      return send(res, 200, MANIFEST, 'application/manifest+json');
    }
    if (req.method === 'GET' && pathname === '/health') {
      return send(res, 200, { status: 'ok', uptime: process.uptime() });
    }

    // Auth
    if (req.method === 'POST' && pathname === '/api/login') {
      const body = await parseBody(req);
      if (body.password === AUTH_PASSWORD) {
        const token = signToken({ iat: Date.now(), sub: 'user' });
        return send(res, 200, { token });
      }
      return send(res, 401, { error: 'Invalid password' });
    }

    // Protected routes
    if (pathname.startsWith('/api/')) {
      if (!authCheck(req)) return send(res, 401, { error: 'Unauthorized' });

      if (req.method === 'GET' && pathname === '/api/cluster') {
        const data = await getClusterData();
        return send(res, 200, data);
      }

      const nsMatch = pathname.match(/^\/api\/namespace\/([a-z0-9-]+)$/);
      if (req.method === 'GET' && nsMatch) {
        const data = await getNamespaceDetail(nsMatch[1]);
        return send(res, 200, data);
      }

      return send(res, 404, { error: 'Not found' });
    }

    send(res, 404, { error: 'Not found' });
  } catch (err) {
    console.error('Request error:', err.message);
    send(res, 500, { error: 'Internal server error' });
  }
});

server.listen(PORT, () => console.log(`holm-mind listening on :${PORT}`));
