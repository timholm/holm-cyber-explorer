const fs = require('fs');
const path = require('path');
const { marked } = require('marked');

// Hardcoded metadata for all repos — order determines display and PROJ-NNN numbering
const REPO_META = {
  'HolmOS': { order: 1, title: 'HolmOS', desc: 'Full web-based operating system on a 13-node Raspberry Pi Kubernetes cluster with 120+ microservices', category: 'core', relatedDomains: ['5', '4'] },
  'holm-hq-interface': { order: 2, title: 'HOLM HQ Interface', desc: 'Central headquarters web interface for the HOLM system', category: 'interface', relatedDomains: ['5', '10'] },
  'holm-mind': { order: 3, title: 'HOLM Mind', desc: 'Knowledge management and intelligence processing system', category: 'core', relatedDomains: ['7', '6'] },
  'animus-dashboard': { order: 4, title: 'Animus Dashboard', desc: 'Assassin\'s Creed themed Kubernetes cluster management dashboard with real-time monitoring', category: 'interface', relatedDomains: ['4', '8'] },
  'bookforge': { order: 5, title: 'BookForge', desc: 'AI-powered book generation system using local LLM with text-to-speech support', category: 'tools', relatedDomains: ['7', '6'] },
  'kube-sentinel': { order: 6, title: 'Kube Sentinel', desc: 'Kubernetes error prioritization and auto-remediation agent with web dashboard', category: 'infrastructure', relatedDomains: ['3', '4'] },
  'k8s-infrastructure': { order: 7, title: 'K8s Infrastructure', desc: 'GitOps repository for Kubernetes infrastructure managed by ArgoCD', category: 'infrastructure', relatedDomains: ['4'] },
  'pi-cluster-ansible': { order: 8, title: 'Pi Cluster Ansible', desc: 'Ansible playbooks for bootstrapping a 13-node Raspberry Pi Kubernetes cluster', category: 'infrastructure', relatedDomains: ['4', '12'] },
  'k8s-cluster-backup': { order: 9, title: 'K8s Cluster Backup', desc: 'Kubernetes cluster backup and disaster recovery utilities', category: 'infrastructure', relatedDomains: ['12', '4'] },
  'docs-framework': { order: 10, title: 'Docs Framework', desc: 'HTML document generation pipeline for the 20-domain knowledge base', category: 'tools', relatedDomains: ['20', 'META'] },
  'ytarchive': { order: 11, title: 'YouTube Archiver', desc: 'Kubernetes-native YouTube channel archiver with parallel workers and iSCSI storage', category: 'tools', relatedDomains: ['6', '18'] },
  'tiktok-archive': { order: 12, title: 'TikTok Archive', desc: 'TikTok content archiving and preservation service', category: 'tools', relatedDomains: ['6', '18'] },
  'claude-code-gitea-actions': { order: 13, title: 'Claude Code Gitea Actions', desc: 'CI/CD integration between Claude Code and Gitea for automated workflows', category: 'tools', relatedDomains: ['8', '19'] }
};

/**
 * Check if any file matching a simple glob pattern exists in a directory (immediate children only).
 * Supports patterns like '*.yaml', '*.yml', 'Dockerfile*'
 */
function hasFilePattern(dir, pattern) {
  try {
    const files = fs.readdirSync(dir);
    // Convert simple glob to regex: * becomes .*, rest is literal
    const regexStr = '^' + pattern.replace(/\./g, '\\.').replace(/\*/g, '.*') + '$';
    const regex = new RegExp(regexStr);
    return files.some(f => regex.test(f));
  } catch {
    return false;
  }
}

/**
 * Detect tech stack by inspecting well-known files in the repo directory.
 */
function detectTechStack(repoDir) {
  const stack = [];
  if (fs.existsSync(path.join(repoDir, 'package.json'))) {
    try {
      const pkg = JSON.parse(fs.readFileSync(path.join(repoDir, 'package.json'), 'utf8'));
      if (pkg.dependencies) {
        if (pkg.dependencies.next) stack.push('Next.js');
        else if (pkg.dependencies.react) stack.push('React');
        else if (pkg.dependencies.express) stack.push('Express');
        else stack.push('Node.js');
      } else {
        stack.push('Node.js');
      }
    } catch { stack.push('Node.js'); }
  }
  if (fs.existsSync(path.join(repoDir, 'go.mod'))) stack.push('Go');
  if (fs.existsSync(path.join(repoDir, 'requirements.txt'))) stack.push('Python');
  if (fs.existsSync(path.join(repoDir, 'Chart.yaml'))) stack.push('Helm');
  if (fs.existsSync(path.join(repoDir, 'Cargo.toml'))) stack.push('Rust');
  // Check for common files
  if (fs.existsSync(path.join(repoDir, 'Dockerfile')) || hasFilePattern(repoDir, 'Dockerfile*')) stack.push('Docker');
  if (hasFilePattern(repoDir, '*.yaml') || hasFilePattern(repoDir, '*.yml')) {
    if (fs.existsSync(path.join(repoDir, 'ansible.cfg')) || fs.existsSync(path.join(repoDir, 'playbooks'))) stack.push('Ansible');
  }
  return stack.length > 0 ? stack : ['Shell'];
}

/**
 * Generate fallback HTML content when no README exists.
 */
function generateFallback(meta, repoDir, repoName) {
  let html = `<h2>About</h2>\n<p>${meta.desc}</p>\n`;

  // List top-level directory contents
  try {
    const entries = fs.readdirSync(repoDir);
    if (entries.length > 0) {
      html += `<h2>Repository Contents</h2>\n<ul>\n`;
      for (const entry of entries.sort()) {
        const fullPath = path.join(repoDir, entry);
        const isDir = fs.statSync(fullPath).isDirectory();
        html += `<li>${isDir ? '<strong>' + entry + '/</strong>' : entry}</li>\n`;
      }
      html += `</ul>\n`;
    }
  } catch {
    // If we can't read the directory, just skip the listing
  }

  html += `<p><em>This project does not have a README. Content was auto-generated from repository metadata.</em></p>`;
  return html;
}

function main() {
  const rootDir = __dirname;
  const reposDir = path.join(rootDir, 'repos');
  const htmlDir = path.join(rootDir, 'html');
  const manifestPath = path.join(rootDir, 'manifest.json');
  const projectsPath = path.join(rootDir, 'projects.json');

  // Ensure html directory exists
  if (!fs.existsSync(htmlDir)) {
    fs.mkdirSync(htmlDir, { recursive: true });
  }

  // Read repo directories
  let repoDirs;
  try {
    repoDirs = fs.readdirSync(reposDir).filter(d => {
      return fs.statSync(path.join(reposDir, d)).isDirectory();
    });
  } catch (err) {
    console.error('Failed to read repos/ directory:', err.message);
    process.exit(1);
  }

  const projects = [];
  let htmlCount = 0;

  // Process each known repo in order
  const sortedRepos = Object.entries(REPO_META)
    .sort((a, b) => a[1].order - b[1].order);

  for (const [repoName, meta] of sortedRepos) {
    const repoDir = path.join(reposDir, repoName);

    if (!fs.existsSync(repoDir)) {
      console.warn(`  SKIP: ${repoName} — directory not found under repos/`);
      continue;
    }

    const orderNum = String(meta.order).padStart(3, '0');
    const docId = `PROJ-${orderNum}`;
    const filename = `proj-${orderNum}.html`;

    // Detect tech stack
    const techStack = detectTechStack(repoDir);

    // Try to read README.md
    let readmeHTML = '';
    const readmePath = path.join(repoDir, 'README.md');
    if (fs.existsSync(readmePath)) {
      try {
        const readmeMd = fs.readFileSync(readmePath, 'utf8');
        readmeHTML = marked(readmeMd);
      } catch (err) {
        console.warn(`  WARN: Failed to parse README for ${repoName}: ${err.message}`);
        readmeHTML = generateFallback(meta, repoDir, repoName);
      }
    } else {
      readmeHTML = generateFallback(meta, repoDir, repoName);
    }

    // Generate HTML file matching the existing import format
    const htmlContent = `<aside class="metadata">
<dl>
<dt>Domain</dt><dd>Projects &amp; Repositories</dd>
<dt>Status</dt><dd>Active</dd>
<dt>Version</dt><dd>1.0</dd>
</dl>
</aside>
<main>
<h1>${meta.title}</h1>
<p><em>${meta.desc}</em></p>
<h2>Tech Stack</h2>
<p>${techStack.join(', ')}</p>
<h2>Repository</h2>
<p><code>${repoName}</code> — Category: ${meta.category}</p>
<hr>
${readmeHTML}
</main>`;

    const htmlPath = path.join(htmlDir, filename);
    fs.writeFileSync(htmlPath, htmlContent, 'utf8');
    htmlCount++;

    // Build project entry
    projects.push({
      id: docId,
      name: meta.title,
      repo: repoName,
      description: meta.desc,
      techStack,
      category: meta.category,
      status: 'active',
      docId
    });

    console.log(`  ${docId}: ${meta.title} (${repoName}) — ${techStack.join(', ')}`);
  }

  // Read existing manifest, remove any existing PROJ entries, append new ones
  let manifest = [];
  try {
    manifest = JSON.parse(fs.readFileSync(manifestPath, 'utf8'));
  } catch (err) {
    console.warn('Could not read manifest.json, starting fresh:', err.message);
  }

  const beforeCount = manifest.length;
  manifest = manifest.filter(entry => entry.domain !== 'PROJ');
  const removedCount = beforeCount - manifest.length;
  if (removedCount > 0) {
    console.log(`Removed ${removedCount} existing PROJ entries from manifest`);
  }

  // Append new PROJ entries
  const newEntries = projects.map(p => ({
    id: p.id,
    title: p.name,
    domain: 'PROJ',
    filename: `proj-${String(REPO_META[p.repo].order).padStart(3, '0')}.html`,
    source: 'index-repos'
  }));

  manifest.push(...newEntries);
  fs.writeFileSync(manifestPath, JSON.stringify(manifest, null, 2) + '\n', 'utf8');

  // Write projects.json
  fs.writeFileSync(projectsPath, JSON.stringify(projects, null, 2) + '\n', 'utf8');

  // Print stats
  console.log('\n--- Index Complete ---');
  console.log(`Repos indexed:         ${projects.length}`);
  console.log(`HTML files generated:  ${htmlCount}`);
  console.log(`Manifest entries added: ${newEntries.length}`);
  console.log(`Total manifest size:   ${manifest.length}`);
}

main();
