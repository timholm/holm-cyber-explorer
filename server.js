const express = require('express');
const { MongoClient, ObjectId } = require('mongodb');
const path = require('path');

const app = express();
const PORT = process.env.PORT || 3000;
const MONGODB_URI = process.env.MONGODB_URI || 'mongodb://localhost:27017/holmvault';

let db;

async function connectDB() {
  const client = new MongoClient(MONGODB_URI);
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
  console.log('Indexes created');
}

// Middleware
app.use(express.json({ limit: '5mb' }));
app.use(express.static(path.join(__dirname, 'public')));

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
app.put('/api/docs/:id', async (req, res) => {
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
app.post('/api/docs', async (req, res) => {
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
app.delete('/api/docs/:id', async (req, res) => {
  try {
    const result = await db.collection('documents').deleteOne({ docId: req.params.id.toUpperCase() });
    if (result.deletedCount === 0) return res.status(404).json({ error: 'Document not found' });
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Full-text search
app.get('/api/search', async (req, res) => {
  try {
    const { q } = req.query;
    if (!q) return res.status(400).json({ error: 'Query parameter q is required' });
    const docs = await db.collection('documents')
      .find(
        { $text: { $search: q } },
        { projection: { content: 0, score: { $meta: 'textScore' } } }
      )
      .sort({ score: { $meta: 'textScore' } })
      .limit(50)
      .toArray();
    res.json(docs);
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

// Trigger import
app.post('/api/import', async (req, res) => {
  try {
    const { execSync } = require('child_process');
    execSync('node import.js', { cwd: __dirname, stdio: 'pipe', timeout: 120000 });
    res.json({ success: true, message: 'Import completed' });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// SPA fallback — serve index.html for non-API, non-file routes
app.get('*', (req, res) => {
  res.sendFile(path.join(__dirname, 'public', 'index.html'));
});

// Start
connectDB().then(() => {
  app.listen(PORT, () => {
    console.log(`HOLM Vault API running on port ${PORT}`);
  });
}).catch(err => {
  console.error('Failed to connect to MongoDB:', err);
  process.exit(1);
});
