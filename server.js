const express = require('express');
const { MongoClient, ObjectId } = require('mongodb');
const path = require('path');
const compression = require('compression');
const crypto = require('crypto');
const fs = require('fs');

const app = express();
app.set('etag', 'strong');
const PORT = process.env.PORT || 3000;
const MONGODB_URI = process.env.MONGODB_URI || 'mongodb://localhost:27017/holmvault';
const BASE_URL = process.env.BASE_URL || 'http://holm.local';
const DATA_CACHE_DIR = process.env.DATA_CACHE_DIR || path.join(__dirname, '.cache');

let db;
let mongoConnected = false;

// ══════════════════════════════════════════════════════════════
// LOCAL FILE CACHE — offline fallback for reads when MongoDB is down
// ══════════════════════════════════════════════════════════════
const localCache = {
  _ensureDir(subdir) {
    const dir = path.join(DATA_CACHE_DIR, subdir);
    if (!fs.existsSync(dir)) fs.mkdirSync(dir, { recursive: true });
    return dir;
  },

  /** Write a document or collection snapshot to disk */
  set(collection, key, data) {
    try {
      const dir = this._ensureDir(collection);
      const filePath = path.join(dir, `${encodeURIComponent(key)}.json`);
      fs.writeFileSync(filePath, JSON.stringify(data), 'utf8');
    } catch (err) {
      console.warn('[cache] Write failed:', err.message);
    }
  },

  /** Read a cached document from disk */
  get(collection, key) {
    try {
      const filePath = path.join(DATA_CACHE_DIR, collection, `${encodeURIComponent(key)}.json`);
      if (!fs.existsSync(filePath)) return null;
      return JSON.parse(fs.readFileSync(filePath, 'utf8'));
    } catch (err) {
      console.warn('[cache] Read failed:', err.message);
      return null;
    }
  },

  /** Check if cache directory exists and has data */
  stats() {
    try {
      if (!fs.existsSync(DATA_CACHE_DIR)) return { exists: false, collections: 0, totalFiles: 0, sizeBytes: 0 };
      const collections = fs.readdirSync(DATA_CACHE_DIR).filter(f =>
        fs.statSync(path.join(DATA_CACHE_DIR, f)).isDirectory()
      );
      let totalFiles = 0;
      let sizeBytes = 0;
      for (const col of collections) {
        const colPath = path.join(DATA_CACHE_DIR, col);
        const files = fs.readdirSync(colPath);
        totalFiles += files.length;
        for (const f of files) {
          sizeBytes += fs.statSync(path.join(colPath, f)).size;
        }
      }
      return { exists: true, collections: collections.length, totalFiles, sizeBytes };
    } catch {
      return { exists: false, collections: 0, totalFiles: 0, sizeBytes: 0 };
    }
  }
};

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
  client.on('connectionPoolClosed', () => {
    console.log('[mongo] Connection pool closed');
    mongoConnected = false;
  });
  client.on('connectionCheckOutFailed', (e) => console.warn('[mongo] Connection checkout failed:', e.reason));
  await client.connect();
  db = client.db();
  mongoConnected = true;
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

  // Activity stream (TTL 7 days)
  await db.collection('activity').createIndex({ createdAt: -1 });
  await db.collection('activity').createIndex({ createdAt: 1 }, { expireAfterSeconds: 604800 });

  // Tasks / roadmap
  try {
    await db.collection('tasks').createIndex({ taskId: 1 }, { unique: true });
  } catch (e) {
    console.warn('tasks taskId index:', e.codeName);
    await db.collection('tasks').createIndex({ taskId: 1 });
  }
  await db.collection('tasks').createIndex({ status: 1, priority: 1 });

  // Agent state singleton
  await db.collection('agent_state').updateOne(
    { _id: 'current' },
    { $setOnInsert: {
      autopilotRunning: false, currentIteration: 0, maxIterations: 0,
      ollamaStatus: 'idle', claudeStatus: 'idle', currentTask: '',
      lastUpdate: new Date(), startedAt: null
    }},
    { upsert: true }
  );

  // ── Orchestrator collections ──

  // Directives — user/system intentions
  try {
    await db.collection('directives').createIndex({ directiveId: 1 }, { unique: true });
  } catch (e) {
    console.warn('directives directiveId index:', e.codeName);
    await db.collection('directives').createIndex({ directiveId: 1 });
  }
  await db.collection('directives').createIndex({ status: 1, priority: 1 });

  // Workers — parallel Claude instances
  await db.collection('workers').createIndex({ status: 1 });
  await db.collection('workers').createIndex({ lastHeartbeat: 1 });

  // Orchestrator state singleton
  await db.collection('orchestrator_state').updateOne(
    { _id: 'orchestrator' },
    { $setOnInsert: {
      running: false, startedAt: null, maxWorkers: 5, activeWorkers: 0,
      totalIterations: 0, totalCost: 0,
      rateLimitBudget: { plan: 'max200', promptsUsed: 0, promptsLimit: 900, windowResetAt: null },
      lastUpdate: new Date()
    }},
    { upsert: true }
  );

  // Extended task indexes for orchestrator
  await db.collection('tasks').createIndex({ directiveId: 1 });
  await db.collection('tasks').createIndex({ assignedWorker: 1, status: 1 });

  console.log('Indexes created');
}

async function connectWithRetry() {
  const MAX_ATTEMPTS = 10;
  const delays = [1000, 2000, 4000, 8000, 16000, 30000, 30000, 30000, 30000, 30000];

  for (let attempt = 1; attempt <= MAX_ATTEMPTS; attempt++) {
    try {
      await connectDB();
      console.log('[mongo] Populating local cache from MongoDB...');
      await populateCache();
      return;
    } catch (err) {
      if (attempt === MAX_ATTEMPTS) {
        console.error(`MongoDB connection failed after ${MAX_ATTEMPTS} attempts — running in OFFLINE mode`);
        console.warn('[offline] API reads will serve from local file cache if available');
        // Don't exit — keep running with cache fallback and retry in background
        scheduleReconnect();
        return;
      }
      const delaySec = delays[attempt - 1] / 1000;
      console.warn(
        `Retrying MongoDB connection (attempt ${attempt}/${MAX_ATTEMPTS}) in ${delaySec}s... (${err.message})`
      );
      await new Promise(resolve => setTimeout(resolve, delays[attempt - 1]));
    }
  }
}

/** Populate the local file cache with current MongoDB data */
async function populateCache() {
  if (!db) return;
  try {
    // Cache all documents (without content for the list, with content individually)
    const docs = await db.collection('documents')
      .find({}, { projection: { content: 0 } })
      .sort({ domain: 1, docId: 1 })
      .toArray();
    localCache.set('documents', '_index', docs);

    // Cache individual documents with content
    const fullDocs = await db.collection('documents').find({}).toArray();
    for (const doc of fullDocs) {
      localCache.set('documents', doc.docId, doc);
    }

    // Cache tasks
    const tasks = await db.collection('tasks').find({}).sort({ status: 1, priority: 1 }).toArray();
    localCache.set('tasks', '_index', tasks);

    // Cache graph data
    const graphDocs = await db.collection('documents')
      .find({}, { projection: { docId: 1, title: 1, domain: 1, domainName: 1, dependsOn: 1, dependedBy: 1 } })
      .toArray();
    localCache.set('graph', '_data', graphDocs);

    // Cache tags
    const tags = await db.collection('documents').aggregate([
      { $unwind: '$tags' },
      { $group: { _id: '$tags', count: { $sum: 1 } } },
      { $sort: { count: -1 } }
    ]).toArray();
    localCache.set('tags', '_index', tags.map(t => ({ tag: t._id, count: t.count })));

    const stats = localCache.stats();
    console.log(`[cache] Populated: ${stats.totalFiles} files, ${(stats.sizeBytes / 1024).toFixed(1)}KB`);
  } catch (err) {
    console.warn('[cache] Population failed:', err.message);
  }
}

/** Background reconnection — tries every 60s when MongoDB is down */
function scheduleReconnect() {
  const RETRY_INTERVAL = 60000;
  const timer = setInterval(async () => {
    if (mongoConnected) {
      clearInterval(timer);
      return;
    }
    console.log('[reconnect] Attempting MongoDB reconnection...');
    try {
      await connectDB();
      console.log('[reconnect] MongoDB reconnected — repopulating cache...');
      await populateCache();
      clearInterval(timer);
    } catch (err) {
      console.warn(`[reconnect] Still offline: ${err.message}`);
    }
  }, RETRY_INTERVAL);
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
    if (filePath.endsWith('.html') || filePath.endsWith('/')) {
      res.setHeader('Cache-Control', 'no-cache');
    } else if (filePath.endsWith('.woff2') || filePath.endsWith('.svg')) {
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

// For /api routes: if MongoDB is down, allow GET routes to try cache fallback
app.use('/api', (req, res, next) => {
  if (!db) {
    // Write operations require MongoDB — no cache fallback
    if (req.method !== 'GET') {
      return res.status(503).json({ error: 'Database unavailable — write operations disabled in offline mode' });
    }
    // Mark request as offline so route handlers can try cache
    req.offlineMode = true;
  }
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
    if (req.offlineMode) {
      const cached = localCache.get('documents', '_index');
      if (cached) {
        res.setHeader('X-Served-From', 'cache');
        const { domain, tag } = req.query;
        let filtered = cached;
        if (domain) filtered = filtered.filter(d => d.domain === domain);
        if (tag) filtered = filtered.filter(d => (d.tags || []).includes(tag));
        return res.json(filtered);
      }
      return res.status(503).json({ error: 'Database offline and no cache available' });
    }
    const { domain, tag } = req.query;
    const filter = {};
    if (domain) filter.domain = domain;
    if (tag) filter.tags = tag;
    const docs = await db.collection('documents')
      .find(filter, { projection: { content: 0 } })
      .sort({ domain: 1, docId: 1 })
      .toArray();
    // Populate cache on successful read (unfiltered)
    if (!domain && !tag) localCache.set('documents', '_index', docs);
    res.setHeader('Cache-Control', 'no-cache');
    res.json(docs);
  } catch (err) {
    // Try cache fallback on MongoDB error
    const cached = localCache.get('documents', '_index');
    if (cached) {
      res.setHeader('X-Served-From', 'cache');
      return res.json(cached);
    }
    res.status(500).json({ error: err.message });
  }
});

// Get single doc by docId
app.get('/api/docs/:id', async (req, res) => {
  const docId = req.params.id.toUpperCase();
  try {
    if (req.offlineMode) {
      const cached = localCache.get('documents', docId);
      if (cached) {
        res.setHeader('X-Served-From', 'cache');
        return res.json(cached);
      }
      return res.status(503).json({ error: 'Database offline and document not cached' });
    }
    const doc = await db.collection('documents').findOne({ docId });
    if (!doc) return res.status(404).json({ error: 'Document not found' });
    // Cache individual document
    localCache.set('documents', docId, doc);
    res.json(doc);
  } catch (err) {
    const cached = localCache.get('documents', docId);
    if (cached) {
      res.setHeader('X-Served-From', 'cache');
      return res.json(cached);
    }
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
  function buildGraph(docs) {
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
    return { nodes, edges };
  }

  try {
    if (req.offlineMode) {
      const cached = localCache.get('graph', '_data');
      if (cached) {
        res.setHeader('X-Served-From', 'cache');
        return res.json(buildGraph(cached));
      }
      return res.status(503).json({ error: 'Database offline and graph not cached' });
    }
    const docs = await db.collection('documents')
      .find({}, { projection: { docId: 1, title: 1, domain: 1, domainName: 1, dependsOn: 1, dependedBy: 1 } })
      .toArray();
    localCache.set('graph', '_data', docs);
    res.json(buildGraph(docs));
  } catch (err) {
    const cached = localCache.get('graph', '_data');
    if (cached) {
      res.setHeader('X-Served-From', 'cache');
      return res.json(buildGraph(cached));
    }
    res.status(500).json({ error: err.message });
  }
});

// All tags with counts
app.get('/api/tags', async (req, res) => {
  try {
    if (req.offlineMode) {
      const cached = localCache.get('tags', '_index');
      if (cached) { res.setHeader('X-Served-From', 'cache'); return res.json(cached); }
      return res.status(503).json({ error: 'Database offline and tags not cached' });
    }
    const tags = await db.collection('documents').aggregate([
      { $unwind: '$tags' },
      { $group: { _id: '$tags', count: { $sum: 1 } } },
      { $sort: { count: -1 } }
    ]).toArray();
    const result = tags.map(t => ({ tag: t._id, count: t.count }));
    localCache.set('tags', '_index', result);
    res.json(result);
  } catch (err) {
    const cached = localCache.get('tags', '_index');
    if (cached) { res.setHeader('X-Served-From', 'cache'); return res.json(cached); }
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
    const detail = req.query.detail === 'true';
    const docs = await db.collection('documents')
      .find({}, { projection: { docId: 1, domain: 1, domainName: 1, dependsOn: 1, dependedBy: 1, status: 1, tags: 1 } })
      .toArray();
    const docIds = new Set(docs.map(d => d.docId));
    const domains = {};
    const brokenList = [];
    for (const doc of docs) {
      const key = doc.domain || 'unknown';
      if (!domains[key]) domains[key] = { name: doc.domainName || key, total: 0, issues: 0, orphaned: 0, broken: 0 };
      domains[key].total++;
      const brokenDeps = (doc.dependsOn || []).filter(d => !docIds.has(d));
      if (brokenDeps.length > 0) {
        domains[key].broken++;
        if (detail) brokenList.push({ docId: doc.docId, domain: key, brokenRefs: brokenDeps });
      }
      if ((doc.dependsOn || []).length === 0 && (doc.dependedBy || []).length === 0) domains[key].orphaned++;
    }
    const result = { totalDocs: docs.length, domains };
    if (detail) result.brokenDetails = brokenList;
    res.json(result);
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

// Repair broken dependencies in-place
// POST /api/repair-deps — removes broken dependsOn refs, rebuilds dependedBy from dependsOn
app.post('/api/repair-deps', authMiddleware, async (req, res) => {
  try {
    const docs = await db.collection('documents')
      .find({}, { projection: { docId: 1, dependsOn: 1, dependedBy: 1 } })
      .toArray();
    const validIds = new Set(docs.map(d => d.docId));

    // 1) Strip broken dependsOn
    let brokenRemoved = 0;
    const brokenDetails = [];
    const ops = [];
    for (const doc of docs) {
      const broken = (doc.dependsOn || []).filter(id => !validIds.has(id));
      if (broken.length > 0) {
        brokenRemoved += broken.length;
        brokenDetails.push({ docId: doc.docId, broken });
      }
    }

    // 2) Build dependedBy inverse map from clean dependsOn
    const dependedByMap = {};
    for (const doc of docs) {
      const cleanDeps = (doc.dependsOn || []).filter(id => validIds.has(id));
      for (const depId of cleanDeps) {
        if (!dependedByMap[depId]) dependedByMap[depId] = new Set();
        dependedByMap[depId].add(doc.docId);
      }
    }

    // 3) Build bulk update ops
    let updated = 0;
    for (const doc of docs) {
      const cleanDependsOn = (doc.dependsOn || []).filter(id => validIds.has(id));
      const newDependedBy = dependedByMap[doc.docId] ? [...dependedByMap[doc.docId]].sort() : [];
      const oldDependsOn = doc.dependsOn || [];
      const oldDependedBy = doc.dependedBy || [];

      // Only update if something changed
      const depsChanged = JSON.stringify(cleanDependsOn) !== JSON.stringify(oldDependsOn);
      const refsChanged = JSON.stringify(newDependedBy) !== JSON.stringify(oldDependedBy);
      if (depsChanged || refsChanged) {
        ops.push({
          updateOne: {
            filter: { docId: doc.docId },
            update: { $set: { dependsOn: cleanDependsOn, dependedBy: newDependedBy, updatedAt: new Date() } }
          }
        });
        updated++;
      }
    }

    if (ops.length > 0) {
      await db.collection('documents').bulkWrite(ops);
    }

    res.json({
      totalDocs: docs.length,
      brokenRefsRemoved: brokenRemoved,
      docsUpdated: updated,
      brokenDetails
    });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Projects data (cached 5 min)
let projectsCache = null;
let projectsCacheTime = 0;
app.get('/api/projects', (req, res) => {
  const now = Date.now();
  if (projectsCache && now - projectsCacheTime < 300000) {
    return res.json(projectsCache);
  }
  const projectsPath = path.join(__dirname, 'projects.json');
  if (!fs.existsSync(projectsPath)) {
    return res.json([]);
  }
  try {
    projectsCache = JSON.parse(fs.readFileSync(projectsPath, 'utf8'));
    projectsCacheTime = now;
    res.json(projectsCache);
  } catch (err) {
    res.status(500).json({ error: 'Failed to read projects data' });
  }
});

// Random document
app.get('/api/random', async (req, res) => {
  try {
    const docs = await db.collection('documents').aggregate([
      { $sample: { size: 1 } },
      { $project: { content: 0 } }
    ]).toArray();
    if (docs.length === 0) return res.status(404).json({ error: 'No documents' });
    res.json(docs[0]);
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

      console.log('[webhook] Indexing project repos...');
      await execAsync('node index-repos.js', { cwd: __dirname, timeout: 60000 });

      console.log('[webhook] Running upsert reimport...');
      await execAsync('node import.js --force', { cwd: __dirname, timeout: 120000 });

      console.log('[webhook] Reimport completed successfully.');

      // Check if server.js changed — if so, exit for k8s restart
      const { stdout: changed } = await execAsync('git diff HEAD~1 --name-only');
      if (changed && changed.includes('server.js')) {
        console.log('[WEBHOOK] server.js changed — exiting for container restart');
        process.exit(0);
      }
    } catch (err) {
      console.error('[webhook] Reimport failed:', err.message);
    }
  })();
});

// Health check — used by Kubernetes liveness/readiness probes
// Returns 200 even in offline mode (app is alive, just degraded)
app.get('/healthz', async (req, res) => {
  const cache = localCache.stats();
  if (db && mongoConnected) {
    try {
      await db.admin().ping();
      return res.status(200).json({ status: 'ok', mongo: 'connected', cache, mode: 'online' });
    } catch (err) {
      mongoConnected = false;
      return res.status(200).json({ status: 'degraded', mongo: 'disconnected', cache, mode: 'offline', error: err.message });
    }
  }
  // Offline but alive — respond 200 so k8s doesn't kill the pod
  res.status(200).json({ status: 'degraded', mongo: 'disconnected', cache, mode: 'offline' });
});

// Storage status — detailed view of persistent storage and cache health
app.get('/api/storage/status', async (req, res) => {
  const cache = localCache.stats();
  let mongo = { connected: false, dbName: null, collections: [] };
  if (db && mongoConnected) {
    try {
      await db.admin().ping();
      const collections = await db.listCollections().toArray();
      const collStats = [];
      for (const col of collections) {
        const count = await db.collection(col.name).estimatedDocumentCount();
        collStats.push({ name: col.name, documents: count });
      }
      mongo = { connected: true, dbName: db.databaseName, collections: collStats };
    } catch (err) {
      mongo.error = err.message;
    }
  }
  res.json({
    mode: mongoConnected ? 'online' : 'offline',
    mongodb: mongo,
    localCache: {
      directory: DATA_CACHE_DIR,
      ...cache
    },
    persistentVolumes: {
      appContent: { mountPath: '/app', description: 'Cached app code — survives pod restarts without GitHub' },
      localCache: { mountPath: DATA_CACHE_DIR, description: 'Document read cache — offline fallback' },
      mongodbData: { mountPath: '/data/db', description: 'MongoDB data directory (Longhorn PVC)' },
      mongodbBackup: { description: 'Nightly mongodump archives (CronJob)' }
    }
  });
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
<meta property="og:url" content="${BASE_URL}/doc/${doc.docId}">
<meta property="og:site_name" content="holm.chat">
<meta name="twitter:card" content="summary">
<meta http-equiv="refresh" content="0;url=/#${doc.docId}">
</head><body><p>Redirecting to <a href="/#${doc.docId}">${title.replace(/</g, '&lt;')}</a>...</p></body></html>`);
  } catch (err) {
    res.redirect('/#' + req.params.id.toUpperCase());
  }
});

// ── Worker Log Buffers (in-memory, ephemeral) ──
const workerLogs = new Map();
const WORKER_LOG_MAX = 500;

// ── SSE Infrastructure ──
const sseClients = [];

function broadcastSSE(eventName, data) {
  const msg = `event: ${eventName}\ndata: ${JSON.stringify(data)}\n\n`;
  for (let i = sseClients.length - 1; i >= 0; i--) {
    try {
      sseClients[i].write(msg);
    } catch {
      sseClients.splice(i, 1);
    }
  }
}

// SSE stream — long-lived EventSource connection
app.get('/api/activity/stream', (req, res) => {
  res.writeHead(200, {
    'Content-Type': 'text/event-stream',
    'Cache-Control': 'no-cache',
    'Connection': 'keep-alive',
    'X-Accel-Buffering': 'no'
  });
  res.write('event: connected\ndata: {"status":"ok"}\n\n');
  sseClients.push(res);
  req.on('close', () => {
    const idx = sseClients.indexOf(res);
    if (idx !== -1) sseClients.splice(idx, 1);
  });
});

// Post activity event
app.post('/api/activity', authMiddleware, async (req, res) => {
  try {
    const { type, iteration, message, detail, agent, status } = req.body;
    if (!type || !message) return res.status(400).json({ error: 'type and message required' });
    const event = {
      type,
      iteration: iteration || 0,
      message,
      detail: (detail || '').slice(0, 2000),
      agent: agent || 'system',
      status: status || 'info',
      duration: req.body.duration || null,
      createdAt: new Date()
    };
    await db.collection('activity').insertOne(event);
    broadcastSSE('activity', event);
    res.status(201).json(event);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Recent activity
app.get('/api/activity', async (req, res) => {
  try {
    const events = await db.collection('activity')
      .find({})
      .sort({ createdAt: -1 })
      .limit(50)
      .toArray();
    res.json(events);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Get agent state — synthesizes from orchestrator_state if running
app.get('/api/agent-state', async (req, res) => {
  try {
    const orchState = await db.collection('orchestrator_state').findOne({ _id: 'orchestrator' });
    if (orchState && orchState.running) {
      // Synthesize old format from orchestrator state
      return res.json({
        autopilotRunning: true,
        currentIteration: orchState.totalIterations || 0,
        maxIterations: 0,
        ollamaStatus: 'monitoring',
        claudeStatus: orchState.activeWorkers > 0 ? 'working' : 'idle',
        currentTask: `Orchestrator: ${orchState.activeWorkers || 0} workers active`,
        lastUpdate: orchState.lastUpdate,
        startedAt: orchState.startedAt
      });
    }
    const state = await db.collection('agent_state').findOne({ _id: 'current' });
    res.json(state || { autopilotRunning: false, ollamaStatus: 'idle', claudeStatus: 'idle' });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Update agent state
app.put('/api/agent-state', authMiddleware, async (req, res) => {
  try {
    const update = { ...req.body, lastUpdate: new Date() };
    delete update._id;
    await db.collection('agent_state').updateOne(
      { _id: 'current' },
      { $set: update },
      { upsert: true }
    );
    const state = await db.collection('agent_state').findOne({ _id: 'current' });
    broadcastSSE('state', state);
    res.json(state);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// List tasks (roadmap) — supports orchestrator query filters
app.get('/api/tasks', async (req, res) => {
  try {
    if (req.offlineMode) {
      const cached = localCache.get('tasks', '_index');
      if (cached) { res.setHeader('X-Served-From', 'cache'); return res.json(cached); }
      return res.status(503).json({ error: 'Database offline and tasks not cached' });
    }
    const filter = {};
    if (req.query.status) filter.status = req.query.status;
    if (req.query.directiveId) filter.directiveId = req.query.directiveId;
    if (req.query.unassigned === 'true') {
      filter.$or = [{ assignedWorker: null }, { assignedWorker: { $exists: false } }];
    }
    const tasks = await db.collection('tasks')
      .find(filter)
      .sort({ status: 1, priority: 1, createdAt: -1 })
      .toArray();
    // Cache unfiltered task list
    if (!req.query.status && !req.query.directiveId && !req.query.unassigned) {
      localCache.set('tasks', '_index', tasks);
    }
    res.json(tasks);
  } catch (err) {
    const cached = localCache.get('tasks', '_index');
    if (cached) { res.setHeader('X-Served-From', 'cache'); return res.json(cached); }
    res.status(500).json({ error: err.message });
  }
});

// Create task
app.post('/api/tasks', authMiddleware, async (req, res) => {
  try {
    const { title, description, priority, source, tags } = req.body;
    if (!title) return res.status(400).json({ error: 'title required' });

    // Auto-generate TASK-NNN
    const last = await db.collection('tasks')
      .find({})
      .sort({ createdAt: -1 })
      .limit(1)
      .toArray();
    let nextNum = 1;
    if (last.length > 0 && last[0].taskId) {
      const m = last[0].taskId.match(/TASK-(\d+)/);
      if (m) nextNum = parseInt(m[1]) + 1;
    }
    const taskId = 'TASK-' + String(nextNum).padStart(3, '0');

    const task = {
      taskId,
      title,
      description: description || '',
      status: 'planned',
      priority: Math.min(5, Math.max(1, parseInt(priority) || 3)),
      iteration: req.body.iteration || null,
      source: source || 'manual',
      tags: tags || [],
      createdAt: new Date(),
      updatedAt: new Date(),
      completedAt: null
    };
    await db.collection('tasks').insertOne(task);
    broadcastSSE('task', task);
    res.status(201).json(task);
  } catch (err) {
    if (err.code === 11000) return res.status(409).json({ error: 'Task already exists' });
    res.status(500).json({ error: err.message });
  }
});

// Update task — extended for orchestrator fields
app.put('/api/tasks/:id', authMiddleware, async (req, res) => {
  try {
    const update = { $set: { updatedAt: new Date() } };
    const allowed = ['title', 'description', 'status', 'priority', 'tags', 'iteration',
      'directiveId', 'assignedWorker', 'sessionId', 'dependencies', 'estimatedCost',
      'actualCost', 'attempt', 'maxAttempts', 'output', 'startedAt', 'failureReason'];
    for (const key of allowed) {
      if (req.body[key] !== undefined) update.$set[key] = req.body[key];
    }
    if (req.body.status === 'completed') update.$set.completedAt = new Date();

    const result = await db.collection('tasks').updateOne(
      { taskId: req.params.id.toUpperCase() },
      update
    );
    if (result.matchedCount === 0) return res.status(404).json({ error: 'Task not found' });
    const task = await db.collection('tasks').findOne({ taskId: req.params.id.toUpperCase() });
    broadcastSSE('task', task);
    res.json(task);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// ── Directive Routes ──

// List all directives
app.get('/api/directives', async (req, res) => {
  try {
    const filter = {};
    if (req.query.status) filter.status = req.query.status;
    const directives = await db.collection('directives')
      .find(filter)
      .sort({ createdAt: -1 })
      .toArray();
    res.json(directives);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Create directive (no auth — single-operator personal system)
app.post('/api/directives', async (req, res) => {
  try {
    const { intent, priority, source, context, maxWorkers } = req.body;
    if (!intent) return res.status(400).json({ error: 'intent required' });

    // Auto-generate DIR-NNN
    const last = await db.collection('directives')
      .find({}).sort({ createdAt: -1 }).limit(1).toArray();
    let nextNum = 1;
    if (last.length > 0 && last[0].directiveId) {
      const m = last[0].directiveId.match(/DIR-(\d+)/);
      if (m) nextNum = parseInt(m[1]) + 1;
    }
    const directiveId = 'DIR-' + String(nextNum).padStart(3, '0');

    const directive = {
      directiveId,
      intent,
      priority: Math.min(5, Math.max(1, parseInt(priority) || 3)),
      status: 'pending',
      source: source || 'user',
      decomposition: [],
      context: context || '',
      maxWorkers: maxWorkers || null,
      createdAt: new Date(),
      updatedAt: new Date(),
      completedAt: null
    };
    await db.collection('directives').insertOne(directive);
    broadcastSSE('directive', directive);
    res.status(201).json(directive);
  } catch (err) {
    if (err.code === 11000) return res.status(409).json({ error: 'Directive already exists' });
    res.status(500).json({ error: err.message });
  }
});

// Update directive status
app.put('/api/directives/:id', authMiddleware, async (req, res) => {
  try {
    const update = { $set: { updatedAt: new Date() } };
    const allowed = ['status', 'priority', 'intent', 'context', 'maxWorkers'];
    for (const key of allowed) {
      if (req.body[key] !== undefined) update.$set[key] = req.body[key];
    }
    if (req.body.status === 'completed') update.$set.completedAt = new Date();

    const result = await db.collection('directives').updateOne(
      { directiveId: req.params.id.toUpperCase() },
      update
    );
    if (result.matchedCount === 0) return res.status(404).json({ error: 'Directive not found' });
    const directive = await db.collection('directives').findOne({ directiveId: req.params.id.toUpperCase() });
    broadcastSSE('directive', directive);
    res.json(directive);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Decompose directive — accept task array, create linked tasks
app.post('/api/directives/:id/decompose', authMiddleware, async (req, res) => {
  try {
    const { tasks: taskDefs } = req.body;
    if (!Array.isArray(taskDefs) || taskDefs.length === 0) {
      return res.status(400).json({ error: 'tasks array required' });
    }
    const directiveId = req.params.id.toUpperCase();
    const directive = await db.collection('directives').findOne({ directiveId });
    if (!directive) return res.status(404).json({ error: 'Directive not found' });

    // Get next task number
    const lastTask = await db.collection('tasks')
      .find({}).sort({ createdAt: -1 }).limit(1).toArray();
    let nextNum = 1;
    if (lastTask.length > 0 && lastTask[0].taskId) {
      const m = lastTask[0].taskId.match(/TASK-(\d+)/);
      if (m) nextNum = parseInt(m[1]) + 1;
    }

    const createdTasks = [];
    for (const td of taskDefs) {
      const taskId = 'TASK-' + String(nextNum++).padStart(3, '0');
      const task = {
        taskId,
        title: td.title || 'Untitled task',
        description: td.description || '',
        status: 'queued',
        priority: Math.min(5, Math.max(1, parseInt(td.priority) || directive.priority)),
        directiveId,
        dependencies: td.dependencies || [],
        estimatedCost: td.estimatedCost || null,
        assignedWorker: null,
        sessionId: null,
        attempt: 0,
        maxAttempts: 3,
        output: null,
        source: 'orchestrator',
        tags: td.tags || [],
        createdAt: new Date(),
        updatedAt: new Date(),
        startedAt: null,
        completedAt: null,
        failureReason: null
      };
      await db.collection('tasks').insertOne(task);
      broadcastSSE('task', task);
      createdTasks.push(task);
    }

    // Link tasks to directive
    const taskIds = createdTasks.map(t => t.taskId);
    await db.collection('directives').updateOne(
      { directiveId },
      { $set: { decomposition: taskIds, status: 'active', updatedAt: new Date() } }
    );
    const updated = await db.collection('directives').findOne({ directiveId });
    broadcastSSE('directive', updated);

    res.status(201).json({ directive: updated, tasks: createdTasks });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// ── Worker Routes ──

// List all workers
app.get('/api/workers', async (req, res) => {
  try {
    const workers = await db.collection('workers')
      .find({})
      .sort({ _id: 1 })
      .toArray();
    res.json(workers);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Upsert worker state
app.put('/api/workers/:id', authMiddleware, async (req, res) => {
  try {
    const workerId = req.params.id;
    const update = { ...req.body, lastHeartbeat: new Date() };
    delete update._id;
    await db.collection('workers').updateOne(
      { _id: workerId },
      { $set: update },
      { upsert: true }
    );
    const worker = await db.collection('workers').findOne({ _id: workerId });
    broadcastSSE('worker', worker);
    res.json(worker);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Remove stopped worker
app.delete('/api/workers/:id', authMiddleware, async (req, res) => {
  try {
    const result = await db.collection('workers').deleteOne({ _id: req.params.id });
    if (result.deletedCount === 0) return res.status(404).json({ error: 'Worker not found' });
    broadcastSSE('worker', { _id: req.params.id, status: 'removed' });
    workerLogs.delete(req.params.id);
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Append log lines to worker buffer and broadcast via SSE
app.post('/api/workers/:id/log', authMiddleware, (req, res) => {
  const wid = req.params.id;
  const { lines } = req.body;
  if (!lines || !Array.isArray(lines)) return res.status(400).json({ error: 'lines array required' });

  if (!workerLogs.has(wid)) workerLogs.set(wid, []);
  const buf = workerLogs.get(wid);

  for (const line of lines) {
    const entry = { text: (line.text || '').slice(0, 500), type: line.type || 'text', tool: line.tool || null, ts: line.ts || new Date().toISOString() };
    buf.push(entry);
    broadcastSSE('worker_log', { workerId: wid, ...entry });
  }

  while (buf.length > WORKER_LOG_MAX) buf.shift();
  res.json({ buffered: buf.length });
});

// Get full log buffer for a worker (initial load)
app.get('/api/workers/:id/logs', (req, res) => {
  res.json(workerLogs.get(req.params.id) || []);
});

// ── Orchestrator Routes ──

// Get orchestrator state
app.get('/api/orchestrator', async (req, res) => {
  try {
    const state = await db.collection('orchestrator_state').findOne({ _id: 'orchestrator' });
    res.json(state || { running: false, activeWorkers: 0, totalCost: 0 });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Update orchestrator state
app.put('/api/orchestrator', authMiddleware, async (req, res) => {
  try {
    const update = { ...req.body, lastUpdate: new Date() };
    delete update._id;
    await db.collection('orchestrator_state').updateOne(
      { _id: 'orchestrator' },
      { $set: update },
      { upsert: true }
    );
    const state = await db.collection('orchestrator_state').findOne({ _id: 'orchestrator' });
    broadcastSSE('orchestrator', state);
    res.json(state);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// ── Backward-compatible agent-state from orchestrator ──
// Existing GET /api/agent-state already works. Extend it to also read orchestrator_state
// for richer data when orchestrator is running.

// Create work log document from autopilot
app.post('/api/activity/doc', authMiddleware, async (req, res) => {
  try {
    const { iteration, summary, analysis, duration, tasksCompleted } = req.body;
    if (!summary) return res.status(400).json({ error: 'summary required' });

    // Auto-generate LOG-NNN
    const lastLog = await db.collection('documents')
      .find({ domain: 'LOG' })
      .sort({ createdAt: -1 })
      .limit(1)
      .toArray();
    let nextNum = 1;
    if (lastLog.length > 0 && lastLog[0].docId) {
      const m = lastLog[0].docId.match(/LOG-(\d+)/);
      if (m) nextNum = parseInt(m[1]) + 1;
    }
    const docId = 'LOG-' + String(nextNum).padStart(3, '0');
    const now = new Date();
    const dateStr = now.toISOString().split('T')[0];

    const content = `<h1>Autopilot Work Log — Iteration ${iteration || '?'}</h1>
<p><strong>Date:</strong> ${dateStr} &nbsp; <strong>Duration:</strong> ${duration || '?'}s</p>
<h2>Summary</h2>
<p>${(summary || '').replace(/</g, '&lt;').replace(/\n/g, '<br>')}</p>
${analysis ? `<h2>Analysis</h2><p>${analysis.replace(/</g, '&lt;').replace(/\n/g, '<br>')}</p>` : ''}
${tasksCompleted ? `<h2>Tasks Completed</h2><p>${tasksCompleted.replace(/</g, '&lt;').replace(/\n/g, '<br>')}</p>` : ''}`;

    const doc = {
      docId,
      title: `Work Log ${docId} — ${dateStr}`,
      domain: 'LOG',
      domainName: 'Autopilot Work Log',
      content,
      tags: ['autopilot', 'work-log', `iteration-${iteration || 0}`],
      dependsOn: [],
      dependedBy: [],
      source: 'autopilot',
      filename: `${docId.toLowerCase()}.html`,
      createdAt: now,
      updatedAt: now
    };
    await db.collection('documents').insertOne(doc);
    broadcastSSE('doc_created', { docId, title: doc.title });
    res.status(201).json(doc);
  } catch (err) {
    if (err.code === 11000) return res.status(409).json({ error: 'Log already exists' });
    res.status(500).json({ error: err.message });
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
