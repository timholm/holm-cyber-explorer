const express = require('express');
const { MongoClient, ObjectId } = require('mongodb');
const path = require('path');
const compression = require('compression');
const crypto = require('crypto');

const app = express();
app.set('etag', 'strong');
const PORT = process.env.PORT || 3000;
const MONGODB_URI = process.env.MONGODB_URI || 'mongodb://localhost:27017/holmvault';

let db;

function validateGithubSignature(payload, signature, secret) {
  const hmac = crypto.createHmac('sha256', secret);
  const digest = 'sha256=' + hmac.update(payload).digest('hex');
  try {
    return crypto.timingSafeEqual(Buffer.from(digest), Buffer.from(signature));
  } catch {
    return false;
  }
}

async function connectDB() {
  const client = new MongoClient(MONGODB_URI, {
    maxPoolSize: 10,
    minPoolSize: 2,
    serverSelectionTimeoutMS: 5000,
    socketTimeoutMS: 30000
  });

  client.on('connectionPoolCreated', () => console.log('[mongo] Connection pool created'));
  client.on('connectionPoolClosed', () => console.log('[mongo] Connection pool closed'));
  client.on('connectionCheckOutFailed', (e) => console.warn('[mongo] Connection checkout failed:', e.reason));
  await client.connect();
  db = client.db();
  console.log('Connected to MongoDB');

  // Create indexes (skip unique if duplicates exist from import)
  const docs = db.collection('documents');
  try {
    await docs.createIndex({ docId: 1 }, { unique: true });
  } catch (e) {
    console.warn('docId unique index failed (duplicates exist), creating non-unique:', e.codeName);
    await docs.createIndex({ docId: 1 });
  }
  await docs.createIndex({ tags: 1 });
  await docs.createIndex({ domain: 1 });
  try {
    await docs.createIndex({ title: 'text', content: 'text' });
  } catch (e) {
    console.warn('Text index may already exist:', e.codeName);
  }
  await db.collection('comments').createIndex({ docId: 1 });
  await db.collection('comments').createIndex({ createdAt: -1 });
  console.log('Indexes created');
}

async function connectWithRetry() {
  const MAX_ATTEMPTS = 10;
  const delays = [1000, 2000, 4000, 8000, 16000, 30000, 30000, 30000, 30000, 30000];

  for (let attempt = 1; attempt <= MAX_ATTEMPTS; attempt++) {
    try {
      await connectDB();
      return;
    } catch (err) {
      if (attempt === MAX_ATTEMPTS) {
        console.error(`MongoDB connection failed after ${MAX_ATTEMPTS} attempts:`, err.message);
        process.exit(1);
      }
      const delaySec = delays[attempt - 1] / 1000;
      console.warn(
        `Retrying MongoDB connection (attempt ${attempt}/${MAX_ATTEMPTS}) in ${delaySec}s... (${err.message})`
      );
      await new Promise(resolve => setTimeout(resolve, delays[attempt - 1]));
    }
  }
}

// Security headers
app.use((req, res, next) => {
  res.setHeader('X-Content-Type-Options', 'nosniff');
  res.setHeader('X-Frame-Options', 'DENY');
  res.setHeader('X-XSS-Protection', '1; mode=block');
  res.setHeader('Referrer-Policy', 'strict-origin-when-cross-origin');
  res.setHeader('Content-Security-Policy', "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; font-src 'self'; img-src 'self' data:; connect-src 'self'; worker-src 'self'; frame-ancestors 'none'");
  next();
});

// Compression
app.use(compression());

// Middleware
app.use(express.json({ limit: '5mb' }));
app.use(express.static(path.join(__dirname, 'public'), {
  maxAge: '1d',
  setHeaders: (res, filePath) => {
    if (filePath.endsWith('.woff2') || filePath.endsWith('.svg')) {
      res.setHeader('Cache-Control', 'public, max-age=31536000, immutable');
    }
  }
}));

// Request logging
app.use((req, res, next) => {
  const start = Date.now();
  res.on('finish', () => {
    const ms = Date.now() - start;
    if (req.path.startsWith('/api/')) {
      console.log(`${req.method} ${req.path} ${res.statusCode} ${ms}ms`);
    }
  });
  next();
});

// Return 503 for /api routes while MongoDB is still connecting
app.use('/api', (req, res, next) => {
  if (!db) return res.status(503).json({ error: 'Database connecting...' });
  next();
});

// API key authentication
const API_KEY = process.env.API_KEY || (() => {
  const generated = crypto.randomBytes(32).toString('hex');
  console.warn(`No API_KEY set — generated for this session: ${generated}`);
  return generated;
})();

function authMiddleware(req, res, next) {
  const authHeader = req.headers['authorization'];
  const xApiKey = req.headers['x-api-key'];
  let providedKey = null;
  if (authHeader && authHeader.startsWith('Bearer ')) {
    providedKey = authHeader.slice(7);
  } else if (xApiKey) {
    providedKey = xApiKey;
  }
  if (!providedKey || providedKey !== API_KEY) {
    return res.status(401).json({ error: 'Unauthorized' });
  }
  next();
}

// List all docs (lightweight — no content field)
app.get('/api/docs', async (req, res) => {
  try {
    const { domain, tag } = req.query;
    const filter = {};
    if (domain) filter.domain = domain;
    if (tag) filter.tags = tag;
    const docs = await db.collection('documents')
      .find(filter, { projection: { content: 0 } })
      .sort({ domain: 1, docId: 1 })
      .toArray();
    // Set ETag-friendly cache headers
    res.setHeader('Cache-Control', 'no-cache');
    res.json(docs);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Get single doc by docId
app.get('/api/docs/:id', async (req, res) => {
  try {
    const doc = await db.collection('documents').findOne({ docId: req.params.id.toUpperCase() });
    if (!doc) return res.status(404).json({ error: 'Document not found' });
    res.json(doc);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Update doc
app.put('/api/docs/:id', authMiddleware, async (req, res) => {
  try {
    const { content, tags, title } = req.body;
    const update = { $set: { updatedAt: new Date() } };
    if (content !== undefined) update.$set.content = content;
    if (tags !== undefined) update.$set.tags = tags;
    if (title !== undefined) update.$set.title = title;
    const result = await db.collection('documents').updateOne(
      { docId: req.params.id.toUpperCase() },
      update
    );
    if (result.matchedCount === 0) return res.status(404).json({ error: 'Document not found' });
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Create doc
app.post('/api/docs', authMiddleware, async (req, res) => {
  try {
    const { docId, title, domain, domainName, content, tags, dependsOn, dependedBy, source } = req.body;
    if (!docId || !title) return res.status(400).json({ error: 'docId and title are required' });
    const doc = {
      docId: docId.toUpperCase(),
      title,
      domain: domain || '',
      domainName: domainName || '',
      content: content || '',
      tags: tags || [],
      dependsOn: dependsOn || [],
      dependedBy: dependedBy || [],
      source: source || 'manual',
      filename: `${docId.toLowerCase()}.html`,
      createdAt: new Date(),
      updatedAt: new Date()
    };
    await db.collection('documents').insertOne(doc);
    res.status(201).json(doc);
  } catch (err) {
    if (err.code === 11000) return res.status(409).json({ error: 'Document already exists' });
    res.status(500).json({ error: err.message });
  }
});

// Delete doc
app.delete('/api/docs/:id', authMiddleware, async (req, res) => {
  try {
    const result = await db.collection('documents').deleteOne({ docId: req.params.id.toUpperCase() });
    if (result.deletedCount === 0) return res.status(404).json({ error: 'Document not found' });
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Full-text search with snippets
app.get('/api/search', async (req, res) => {
  try {
    const { q } = req.query;
    if (!q) return res.status(400).json({ error: 'Query parameter q is required' });
    const docs = await db.collection('documents')
      .find(
        { $text: { $search: q } },
        { projection: { score: { $meta: 'textScore' } } }
      )
      .sort({ score: { $meta: 'textScore' } })
      .limit(50)
      .toArray();

    const term = q.trim();
    const termRe = new RegExp(term.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'i');

    const results = docs.map(doc => {
      const plain = (doc.content || '').replace(/<[^>]*>/g, '');
      const matchIdx = plain.search(termRe);
      let snippet = '';
      if (matchIdx !== -1) {
        const start = Math.max(0, matchIdx - 60);
        const end = Math.min(plain.length, matchIdx + term.length + 90);
        const raw = plain.slice(start, end).trim();
        snippet = raw.replace(termRe, m => `<mark>${m}</mark>`);
      } else {
        snippet = plain.slice(0, 150).trim();
      }
      const { content, ...rest } = doc;
      return { ...rest, snippet };
    });

    res.json(results);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Graph data — all nodes and edges
app.get('/api/graph', async (req, res) => {
  try {
    const docs = await db.collection('documents')
      .find({}, { projection: { docId: 1, title: 1, domain: 1, domainName: 1, dependsOn: 1, dependedBy: 1 } })
      .toArray();

    const nodes = docs.map(d => ({
      id: d.docId,
      title: d.title,
      domain: d.domain,
      domainName: d.domainName
    }));

    const edges = [];
    const edgeSet = new Set();
    for (const d of docs) {
      for (const dep of (d.dependsOn || [])) {
        const key = `${d.docId}->${dep}`;
        if (!edgeSet.has(key)) {
          edgeSet.add(key);
          edges.push({ source: d.docId, target: dep });
        }
      }
    }
    res.json({ nodes, edges });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// All tags with counts
app.get('/api/tags', async (req, res) => {
  try {
    const tags = await db.collection('documents').aggregate([
      { $unwind: '$tags' },
      { $group: { _id: '$tags', count: { $sum: 1 } } },
      { $sort: { count: -1 } }
    ]).toArray();
    res.json(tags.map(t => ({ tag: t._id, count: t.count })));
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// List comments for a document
app.get('/api/docs/:id/comments', async (req, res) => {
  try {
    const comments = await db.collection('comments')
      .find({ docId: req.params.id.toUpperCase() })
      .sort({ createdAt: -1 })
      .toArray();
    res.json(comments);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Add a comment (rate-limited, server-side sanitized)
const commentRateLimit = {};
// Clean up expired rate limit entries every 5 minutes
setInterval(() => {
  const now = Date.now();
  for (const ip in commentRateLimit) {
    commentRateLimit[ip] = commentRateLimit[ip].filter(t => now - t < 60000);
    if (commentRateLimit[ip].length === 0) delete commentRateLimit[ip];
  }
}, 300000);
app.post('/api/docs/:id/comments', async (req, res) => {
  try {
    // Simple IP-based rate limiting (5 comments per minute)
    const ip = req.ip || req.connection.remoteAddress;
    const now = Date.now();
    if (!commentRateLimit[ip]) commentRateLimit[ip] = [];
    commentRateLimit[ip] = commentRateLimit[ip].filter(t => now - t < 60000);
    if (commentRateLimit[ip].length >= 5) {
      return res.status(429).json({ error: 'Rate limit exceeded. Try again in a minute.' });
    }
    commentRateLimit[ip].push(now);

    const { author, text } = req.body;
    if (!text || !text.trim()) return res.status(400).json({ error: 'Comment text is required' });
    if (!author || !author.trim()) return res.status(400).json({ error: 'Author name is required' });

    // Server-side sanitization: strip all HTML tags
    const stripHtml = s => s.replace(/<[^>]*>/g, '');
    const comment = {
      docId: req.params.id.toUpperCase(),
      author: stripHtml(author.trim()).slice(0, 50),
      text: stripHtml(text.trim()).slice(0, 2000),
      createdAt: new Date()
    };
    const result = await db.collection('comments').insertOne(comment);
    comment._id = result.insertedId;
    res.status(201).json(comment);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Delete a comment (auth required)
app.delete('/api/comments/:id', authMiddleware, async (req, res) => {
  try {
    const result = await db.collection('comments').deleteOne({ _id: new ObjectId(req.params.id) });
    if (result.deletedCount === 0) return res.status(404).json({ error: 'Comment not found' });
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Domain health summary
app.get('/api/health', async (req, res) => {
  try {
    const docs = await db.collection('documents')
      .find({}, { projection: { docId: 1, domain: 1, domainName: 1, dependsOn: 1, dependedBy: 1, status: 1, tags: 1 } })
      .toArray();
    const docIds = new Set(docs.map(d => d.docId));
    const domains = {};
    for (const doc of docs) {
      const key = doc.domain || 'unknown';
      if (!domains[key]) domains[key] = { name: doc.domainName || key, total: 0, issues: 0, orphaned: 0, broken: 0 };
      domains[key].total++;
      const brokenDeps = (doc.dependsOn || []).filter(d => !docIds.has(d));
      if (brokenDeps.length > 0) domains[key].broken++;
      if ((doc.dependsOn || []).length === 0 && (doc.dependedBy || []).length === 0) domains[key].orphaned++;
    }
    res.json({ totalDocs: docs.length, domains });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});


// Aggregate statistics
app.get('/api/stats', async (req, res) => {
  try {
    const docs = await db.collection('documents')
      .find({}, { projection: { docId: 1, domain: 1, dependsOn: 1, dependedBy: 1, tags: 1, status: 1, content: 1 } })
      .toArray();
    
    const docIds = new Set(docs.map(d => d.docId));
    let totalWords = 0;
    let brokenDeps = 0;
    let orphaned = 0;
    let totalDeps = 0;
    const tagCounts = {};
    const statusCounts = {};
    const domainCounts = {};
    
    for (const doc of docs) {
      // Word count
      const plain = (doc.content || '').replace(/<[^>]*>/g, '');
      totalWords += plain.split(/\s+/).filter(w => w.length > 0).length;
      
      // Dependencies
      const deps = doc.dependsOn || [];
      const refs = doc.dependedBy || [];
      totalDeps += deps.length;
      deps.forEach(d => { if (!docIds.has(d)) brokenDeps++; });
      if (deps.length === 0 && refs.length === 0) orphaned++;
      
      // Tags
      (doc.tags || []).forEach(t => { tagCounts[t] = (tagCounts[t] || 0) + 1; });
      
      // Status
      const st = doc.status || 'unknown';
      statusCounts[st] = (statusCounts[st] || 0) + 1;
      
      // Domain
      const dk = doc.domain || 'unknown';
      domainCounts[dk] = (domainCounts[dk] || 0) + 1;
    }
    
    const commentCount = await db.collection('comments').countDocuments();
    
    res.json({
      documents: docs.length,
      domains: Object.keys(domainCounts).length,
      totalWords,
      totalDependencies: totalDeps,
      brokenDependencies: brokenDeps,
      orphanedDocuments: orphaned,
      comments: commentCount,
      topTags: Object.entries(tagCounts).sort((a, b) => b[1] - a[1]).slice(0, 20).map(([tag, count]) => ({ tag, count })),
      statusBreakdown: statusCounts,
      domainBreakdown: domainCounts
    });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Trigger import (async, non-blocking)
app.post('/api/import', authMiddleware, async (req, res) => {
  try {
    const { exec } = require('child_process');
    const { promisify } = require('util');
    const execAsync = promisify(exec);
    await execAsync('node import.js --force', { cwd: __dirname, timeout: 120000 });
    res.json({ success: true, message: 'Import completed' });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// GitHub webhook — drop collection and reimport on push to main
app.post('/api/webhook', express.raw({ type: 'application/json' }), async (req, res) => {
  const secret = process.env.WEBHOOK_SECRET;
  const signature = req.headers['x-hub-signature-256'];

  if (secret) {
    if (!signature || !validateGithubSignature(req.body, signature, secret)) {
      return res.status(401).json({ error: 'Invalid webhook signature' });
    }
  }

  let payload;
  try {
    payload = JSON.parse(req.body.toString());
  } catch {
    return res.status(400).json({ error: 'Invalid JSON payload' });
  }

  const pushedBranch = (payload.ref || '').replace('refs/heads/', '');
  if (pushedBranch !== 'main') {
    return res.status(200).json({ success: true, message: 'Push to non-main branch ignored' });
  }

  res.json({ success: true, message: 'Reimport triggered' });

  (async () => {
    try {
      const { exec } = require('child_process');
      const { promisify } = require('util');
      const execAsync = promisify(exec);

      console.log('[webhook] Pulling latest docs...');
      try {
        await execAsync('git -C /app pull');
      } catch (pullErr) {
        console.warn('[webhook] git pull failed, re-cloning:', pullErr.message);
        await execAsync('rm -rf /tmp/docs-update && git clone --depth 1 https://github.com/timholm/docs-framework.git /tmp/docs-update && cp /tmp/docs-update/html/manifest.json /app/manifest.json && cp -r /tmp/docs-update/html /app/html');
      }

      console.log('[webhook] Running upsert reimport...');
      await execAsync('node import.js --force', { cwd: __dirname, timeout: 120000 });

      console.log('[webhook] Reimport completed successfully.');
    } catch (err) {
      console.error('[webhook] Reimport failed:', err.message);
    }
  })();
});

// Health check — used by Kubernetes liveness/readiness probes
app.get('/healthz', async (req, res) => {
  try {
    await db.admin().ping();
    res.status(200).json({ status: 'ok', mongo: 'connected' });
  } catch (err) {
    res.status(503).json({ status: 'error', mongo: 'disconnected', error: err.message });
  }
});

// Document permalink with Open Graph tags for link previews
app.get('/doc/:id', async (req, res) => {
  try {
    const doc = await db.collection('documents').findOne(
      { docId: req.params.id.toUpperCase() },
      { projection: { docId: 1, title: 1, domainName: 1, content: 1 } }
    );
    if (!doc) {
      return res.redirect('/#' + req.params.id.toUpperCase());
    }
    const plain = (doc.content || '').replace(/<[^>]*>/g, '').slice(0, 200).trim();
    const title = doc.title + ' — holm.chat';
    const desc = plain || doc.domainName || 'Document ' + doc.docId;

    res.send(`<!DOCTYPE html>
<html><head>
<meta charset="UTF-8">
<title>${title.replace(/"/g, '&quot;')}</title>
<meta property="og:title" content="${title.replace(/"/g, '&quot;')}">
<meta property="og:description" content="${desc.replace(/"/g, '&quot;')}">
<meta property="og:type" content="article">
<meta property="og:url" content="https://holm.chat/doc/${doc.docId}">
<meta property="og:site_name" content="holm.chat">
<meta name="twitter:card" content="summary">
<meta http-equiv="refresh" content="0;url=/#${doc.docId}">
</head><body><p>Redirecting to <a href="/#${doc.docId}">${title.replace(/</g, '&lt;')}</a>...</p></body></html>`);
  } catch (err) {
    res.redirect('/#' + req.params.id.toUpperCase());
  }
});

// SPA fallback — serve index.html for navigation, 404 for everything else
app.get('*', (req, res) => {
  // Don't serve index.html for API routes that weren't matched
  if (req.path.startsWith('/api/')) {
    return res.status(404).json({ error: 'Not found' });
  }
  // Don't serve index.html for direct file requests that don't exist
  if (req.path.match(/\.\w{2,4}$/) && !req.path.endsWith('.html')) {
    return res.status(404).send('Not found');
  }
  res.sendFile(path.join(__dirname, 'public', 'index.html'));
});

// Start the HTTP server immediately so k8s probes can respond while MongoDB connects.
// API routes return 503 until db is populated by connectWithRetry().
const server = app.listen(PORT, () => {
  console.log(`holm.chat running on port ${PORT}`);
});

// Graceful shutdown
process.on('SIGTERM', () => {
  console.log('SIGTERM received, shutting down gracefully...');
  server.close(() => {
    console.log('HTTP server closed');
    process.exit(0);
  });
  setTimeout(() => { console.error('Forced shutdown'); process.exit(1); }, 10000);
});

connectWithRetry();
