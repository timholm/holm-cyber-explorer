const { MongoClient } = require('mongodb');
const fs = require('fs');
const path = require('path');

const MONGODB_URI = process.env.MONGODB_URI || 'mongodb://localhost:27017/holmvault';

// Domain number → name mapping
const DOMAIN_NAMES = {
  '1': 'Constitution & Philosophy',
  '2': 'Governance & Authority',
  '3': 'Security & Integrity',
  '4': 'Infrastructure & Power',
  '5': 'Platform & Core Systems',
  '6': 'Data & Archives',
  '7': 'Intelligence & Analysis',
  '8': 'Automation & Agents',
  '9': 'Education & Training',
  '10': 'User Operations',
  '11': 'Administration',
  '12': 'Disaster Recovery',
  '13': 'Evolution & Adaptation',
  '14': 'Research & Theory',
  '15': 'Ethics & Safeguards',
  '16': 'Federation & External Relations',
  '17': 'Federation Protocol',
  '18': 'Data Pipeline',
  '19': 'Advanced Automation',
  '20': 'Meta & Documentation',
  'META': 'Meta & Documentation',
  'FW': 'Framework',
  'HIC': 'HIC Architecture'
};

/**
 * Parse dependency references from an HTML dd element text.
 * Handles formats like:
 *   "ETH-001"
 *   "ETH-001, CON-001"
 *   "D18-002 (Validation Pipeline), D18-008 (Trust Scoring)"
 *   "All articles in all domains." (returns empty — it's prose, not IDs)
 */
function parseDependencies(text) {
  if (!text) return [];
  // Split by comma, then extract doc IDs (patterns like CON-001, GOV-001, D4-001, ETH-001, SEC-001, etc.)
  const idPattern = /\b([A-Z]{1,10}-\d{1,4}|D\d{1,2}-\d{1,4})\b/g;
  const ids = [];
  let match;
  while ((match = idPattern.exec(text)) !== null) {
    ids.push(match[1]);
  }
  return [...new Set(ids)];
}

/**
 * Extract metadata and content from an HTML file.
 * Parses the <aside class="metadata"> block and extracts the <main> content.
 */
function parseHTMLFile(htmlContent) {
  const result = {
    dependsOn: [],
    dependedBy: [],
    domainRaw: '',
    version: '',
    date: '',
    status: ''
  };

  // Extract metadata from <aside class="metadata">
  const metadataMatch = htmlContent.match(/<aside class="metadata">([\s\S]*?)<\/aside>/);
  if (metadataMatch) {
    const metaHTML = metadataMatch[1];

    // Parse dt/dd pairs
    const dtddPattern = /<dt>(.*?)<\/dt>\s*<dd>([\s\S]*?)<\/dd>/g;
    let m;
    while ((m = dtddPattern.exec(metaHTML)) !== null) {
      const key = m[1].trim();
      const val = m[2].replace(/<[^>]+>/g, '').trim();

      switch (key) {
        case 'Depends On':
          result.dependsOn = parseDependencies(val);
          break;
        case 'Depended Upon By':
          result.dependedBy = parseDependencies(val);
          break;
        case 'Domain':
          result.domainRaw = val;
          break;
        case 'Version':
          result.version = val;
          break;
        case 'Date':
          result.date = val;
          break;
        case 'Status':
          result.status = val;
          break;
      }
    }
  }

  // Extract body content from <main>...</main>
  const mainMatch = htmlContent.match(/<main>([\s\S]*?)<\/main>/);
  result.content = mainMatch ? mainMatch[1].trim() : htmlContent;

  return result;
}

/**
 * Generate tags for a document based on its domain and source.
 */
function generateTags(entry, domainName) {
  const tags = [];
  if (entry.domain && entry.domain !== 'FW') {
    tags.push(`domain-${entry.domain.toLowerCase()}`);
  }
  if (domainName) {
    // Add a simplified tag from domain name
    const simplified = domainName.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/-+$/, '');
    tags.push(simplified);
  }
  if (entry.source) {
    // Extract stage
    const stageMatch = entry.source.match(/stage(\d+)/);
    if (stageMatch) tags.push(`stage-${stageMatch[1]}`);
    tags.push(entry.source);
  }
  // Tag by ID prefix
  const prefixMatch = entry.id.match(/^([A-Z]+)/);
  if (prefixMatch) {
    tags.push(prefixMatch[1].toLowerCase());
  }
  return [...new Set(tags)];
}

async function main() {
  console.log('Starting HOLM Vault import...');
  console.log(`MongoDB: ${MONGODB_URI}`);

  const client = new MongoClient(MONGODB_URI);
  await client.connect();
  const db = client.db();
  const collection = db.collection('documents');

  // Check if already imported
  const count = await collection.countDocuments();
  if (count > 0) {
    console.log(`Database already has ${count} documents. Skipping import.`);
    await client.close();
    return;
  }

  // Read manifest
  const manifestPath = path.join(__dirname, 'manifest.json');
  if (!fs.existsSync(manifestPath)) {
    console.error('manifest.json not found at', manifestPath);
    await client.close();
    process.exit(1);
  }

  const manifest = JSON.parse(fs.readFileSync(manifestPath, 'utf8'));
  console.log(`Found ${manifest.length} entries in manifest`);

  const htmlDir = path.join(__dirname, 'html');
  const documents = [];
  let parsed = 0;
  let skipped = 0;

  for (const entry of manifest) {
    const htmlPath = path.join(htmlDir, entry.filename);

    if (!fs.existsSync(htmlPath)) {
      console.warn(`  SKIP: ${entry.id} — file not found: ${entry.filename}`);
      skipped++;
      continue;
    }

    const htmlContent = fs.readFileSync(htmlPath, 'utf8');
    const meta = parseHTMLFile(htmlContent);

    const domainName = DOMAIN_NAMES[entry.domain] || entry.domain || '';
    const tags = generateTags(entry, domainName);

    documents.push({
      docId: entry.id,
      title: entry.title,
      domain: entry.domain,
      domainName,
      content: meta.content,
      tags,
      dependsOn: meta.dependsOn,
      dependedBy: meta.dependedBy,
      source: entry.source,
      filename: entry.filename,
      version: meta.version,
      date: meta.date,
      status: meta.status,
      createdAt: new Date(),
      updatedAt: new Date()
    });
    parsed++;
  }

  console.log(`Parsed ${parsed} documents, skipped ${skipped}`);

  if (documents.length > 0) {
    await collection.insertMany(documents);
    console.log(`Inserted ${documents.length} documents into MongoDB`);
  }

  // Create indexes
  await collection.createIndex({ docId: 1 }, { unique: true });
  await collection.createIndex({ tags: 1 });
  await collection.createIndex({ domain: 1 });
  await collection.createIndex({ title: 'text', content: 'text' });
  console.log('Indexes created');

  await client.close();
  console.log('Import complete!');
}

main().catch(err => {
  console.error('Import failed:', err);
  process.exit(1);
});
