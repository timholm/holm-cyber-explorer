const http = require('http');
const fs = require('fs');
const crypto = require('crypto');

const PORT = process.env.PORT || 3000;
const GITEA_URL = process.env.GITEA_URL || 'http://gitea.gitea.svc.cluster.local:3000';
const GITEA_TOKEN = process.env.GITEA_TOKEN || '';
const CMD_API_URL = process.env.CMD_API_URL || 'http://cmd-api.default.svc.cluster.local:80';
const POLL_INTERVAL = parseInt(process.env.POLL_INTERVAL || '30', 10) * 1000;
const BOT_USER = process.env.BOT_USER || 'tim'; // user whose token we use

const events = [];
const MAX_EVENTS = 200;
const processedComments = new Set();
const processedIssues = new Set();
let pollCount = 0;
let lastPollTime = null;
let pollErrors = [];

// --- HTTP helpers ---

function httpRequest(baseUrl, method, path, body, timeout) {
  return new Promise((resolve, reject) => {
    const url = new URL(path, baseUrl);
    const isHttps = url.protocol === 'https:';
    const mod = isHttps ? require('https') : http;
    const opts = {
      hostname: url.hostname,
      port: url.port || (isHttps ? 443 : 80),
      path: url.pathname + url.search,
      method,
      headers: { 'Content-Type': 'application/json' },
    };
    if (baseUrl === GITEA_URL || path.startsWith('/api/v1')) {
      opts.headers['Authorization'] = `token ${GITEA_TOKEN}`;
    }
    const req = mod.request(opts, (res) => {
      let data = '';
      res.on('data', (c) => data += c);
      res.on('end', () => {
        try {
          resolve({ status: res.statusCode, data: JSON.parse(data) });
        } catch {
          resolve({ status: res.statusCode, data });
        }
      });
    });
    req.on('error', reject);
    req.setTimeout(timeout || 30000, () => { req.destroy(); reject(new Error('Request timeout')); });
    if (body) req.write(JSON.stringify(body));
    req.end();
  });
}

// --- Gitea API ---

function giteaGet(path) { return httpRequest(GITEA_URL, 'GET', `/api/v1${path}`); }
function giteaPost(path, body) { return httpRequest(GITEA_URL, 'POST', `/api/v1${path}`, body); }
function giteaPatch(path, body) { return httpRequest(GITEA_URL, 'PATCH', `/api/v1${path}`, body); }
function giteaPut(path, body) { return httpRequest(GITEA_URL, 'PUT', `/api/v1${path}`, body); }
function giteaDelete(path, body) { return httpRequest(GITEA_URL, 'DELETE', `/api/v1${path}`, body); }

async function getMyRepos() {
  const result = await giteaGet('/user/repos?limit=50');
  return Array.isArray(result.data) ? result.data : [];
}

async function getRepoIssues(owner, repo) {
  const result = await giteaGet(`/repos/${owner}/${repo}/issues?state=open&type=issues&sort=updated&limit=10`);
  return Array.isArray(result.data) ? result.data : [];
}

async function getIssueComments(owner, repo, number) {
  const result = await giteaGet(`/repos/${owner}/${repo}/issues/${number}/comments`);
  return Array.isArray(result.data) ? result.data : [];
}

async function getIssue(owner, repo, number) {
  return (await giteaGet(`/repos/${owner}/${repo}/issues/${number}`)).data;
}

async function createComment(owner, repo, number, body) {
  return (await giteaPost(`/repos/${owner}/${repo}/issues/${number}/comments`, { body })).data;
}

async function editComment(owner, repo, commentId, body) {
  return (await giteaPatch(`/repos/${owner}/${repo}/issues/comments/${commentId}`, { body }));
}

async function getPullRequest(owner, repo, number) {
  const pr = (await giteaGet(`/repos/${owner}/${repo}/pulls/${number}`)).data;
  let diff = '';
  try {
    const diffResult = await httpRequest(GITEA_URL, 'GET', `/api/v1/repos/${owner}/${repo}/pulls/${number}.diff`);
    diff = typeof diffResult.data === 'string' ? diffResult.data : '';
  } catch {}
  return { ...pr, diff };
}

// --- Gitea file operations ---

async function getRepoTree(owner, repo) {
  const result = await giteaGet(`/repos/${owner}/${repo}/git/trees/main?recursive=true`);
  if (!result.data?.tree) {
    // Try master branch
    const result2 = await giteaGet(`/repos/${owner}/${repo}/git/trees/master?recursive=true`);
    if (!result2.data?.tree) return [];
    return result2.data.tree.filter(f => f.type === 'blob').map(f => f.path);
  }
  return result.data.tree.filter(f => f.type === 'blob').map(f => f.path);
}

async function getFileContent(owner, repo, filepath) {
  const result = await giteaGet(`/repos/${owner}/${repo}/contents/${encodeURIComponent(filepath)}`);
  if (result.status !== 200 || !result.data?.content) return null;
  return { content: Buffer.from(result.data.content, 'base64').toString('utf-8'), sha: result.data.sha };
}

async function createOrUpdateFile(owner, repo, filepath, content, message, sha) {
  const payload = {
    content: Buffer.from(content).toString('base64'),
    message: message || `Update ${filepath}`,
  };
  if (sha) {
    payload.sha = sha;
    return await giteaPut(`/repos/${owner}/${repo}/contents/${encodeURIComponent(filepath)}`, payload);
  }
  return await giteaPost(`/repos/${owner}/${repo}/contents/${encodeURIComponent(filepath)}`, payload);
}

// --- Gather repo context ---

async function gatherRepoContext(owner, repo) {
  let context = '';
  try {
    const files = await getRepoTree(owner, repo);
    if (files.length > 0) {
      context += `**Repository files:**\n\`\`\`\n${files.join('\n')}\n\`\`\`\n\n`;

      // Read key files
      const keyFiles = files.filter(f =>
        /^(README|readme)/i.test(f) ||
        /^(package\.json|go\.mod|Cargo\.toml|pyproject\.toml|Makefile|Dockerfile)$/i.test(f)
      ).slice(0, 5);

      // If small repo, read all non-chart/template files
      if (files.length <= 20) {
        const allReadable = files.filter(f => !f.startsWith('.git') && !/\.(png|jpg|jpeg|gif|ico|woff|ttf)$/i.test(f));
        keyFiles.push(...allReadable.filter(f => !keyFiles.includes(f)));
      }

      for (const f of keyFiles.slice(0, 15)) {
        try {
          const file = await getFileContent(owner, repo, f);
          if (file && file.content.length < 8000) {
            context += `**File: ${f}**\n\`\`\`\n${file.content}\n\`\`\`\n\n`;
          }
        } catch {}
      }
    }
  } catch (err) {
    context += `(Could not read repo tree: ${err.message})\n\n`;
  }
  return context;
}

// --- Claude API ---

function callClaude(prompt) {
  return new Promise((resolve, reject) => {
    const url = new URL('/claude', CMD_API_URL);
    const payload = JSON.stringify({ prompt });
    const opts = {
      hostname: url.hostname,
      port: url.port || 80,
      path: url.pathname,
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Content-Length': Buffer.byteLength(payload),
      },
    };
    const req = http.request(opts, (res) => {
      let data = '';
      res.on('data', (c) => data += c);
      res.on('end', () => {
        try {
          const json = JSON.parse(data);
          resolve(json.response || json.error || 'No response');
        } catch {
          resolve(data || 'Empty response');
        }
      });
    });
    req.on('error', reject);
    req.setTimeout(300000, () => { req.destroy(); reject(new Error('Claude timeout (5min)')); });
    req.write(payload);
    req.end();
  });
}

// --- Parse and apply file changes ---

function parseFileChanges(response) {
  const changes = [];
  const pattern = /FILE:\s*(\S+)\s*\n```[\w]*\n([\s\S]*?)```/g;
  let match;
  while ((match = pattern.exec(response)) !== null) {
    changes.push({ path: match[1], content: match[2] });
  }
  const deletePattern = /DELETE:\s*(\S+)/g;
  while ((match = deletePattern.exec(response)) !== null) {
    changes.push({ path: match[1], delete: true });
  }
  return changes;
}

async function applyFileChanges(owner, repo, changes, issueNumber) {
  const results = [];
  for (const change of changes) {
    try {
      if (change.delete) {
        results.push(`Would delete \`${change.path}\` (not implemented yet)`);
      } else {
        const existing = await getFileContent(owner, repo, change.path);
        const msg = `${existing ? 'Update' : 'Create'} ${change.path} (issue #${issueNumber})`;
        await createOrUpdateFile(owner, repo, change.path, change.content, msg, existing?.sha);
        results.push(`${existing ? 'Updated' : 'Created'} \`${change.path}\``);
      }
    } catch (err) {
      results.push(`Failed \`${change.path}\`: ${err.message}`);
    }
  }
  return results;
}

// --- Event store ---

function addEvent(event) {
  event.id = crypto.randomUUID();
  event.timestamp = new Date().toISOString();
  events.unshift(event);
  if (events.length > MAX_EVENTS) events.length = MAX_EVENTS;
  return event;
}

// --- Process an issue/comment with Claude ---

async function processWithClaude(owner, repo, issueNumber, isPR, userPrompt, triggerUser, triggerType) {
  const event = addEvent({
    type: triggerType,
    repo: `${owner}/${repo}`,
    issue: issueNumber,
    isPR,
    user: triggerUser,
    prompt: userPrompt?.substring(0, 500) || 'auto-triggered',
    status: 'working',
    response: null,
  });

  console.log(`[Process] ${owner}/${repo}#${issueNumber} - ${triggerType} by ${triggerUser}`);

  let thinkingComment;
  try {
    thinkingComment = await createComment(owner, repo, issueNumber,
      '> **Claude** is thinking...\n\n_Analyzing the repository and preparing a response..._');
  } catch (err) {
    console.error(`[Process] Failed to post thinking comment: ${err.message}`);
    event.status = 'error';
    event.response = `Failed to post comment: ${err.message}`;
    return;
  }

  const commentId = thinkingComment?.id;

  try {
    let context = '';

    if (isPR) {
      try {
        const pr = await getPullRequest(owner, repo, issueNumber);
        context += `## Pull Request: ${pr.title || ''}\n\n`;
        context += `**Description:**\n${pr.body || 'No description'}\n\n`;
        if (pr.diff) {
          const truncDiff = pr.diff.length > 15000 ? pr.diff.substring(0, 15000) + '\n...(truncated)' : pr.diff;
          context += `**Diff:**\n\`\`\`diff\n${truncDiff}\n\`\`\`\n\n`;
        }
      } catch {}
    } else {
      try {
        const issue = await getIssue(owner, repo, issueNumber);
        context += `## Issue #${issueNumber}: ${issue.title || ''}\n\n`;
        context += `**Description:**\n${issue.body || 'No description'}\n\n`;
      } catch {}
    }

    // Get recent comments
    try {
      const comments = await getIssueComments(owner, repo, issueNumber);
      const relevant = comments.filter(c => c.id !== commentId && !c.body?.includes('**Claude** is thinking'));
      if (relevant.length > 0) {
        context += `**Comments:**\n`;
        for (const c of relevant.slice(-10)) {
          const u = c.user?.login || c.user?.username || '?';
          context += `- @${u}: ${c.body?.substring(0, 500) || ''}\n`;
        }
        context += '\n';
      }
    } catch {}

    // Repo context
    const repoContext = await gatherRepoContext(owner, repo);
    context += repoContext;

    const fullPrompt = `You are an autonomous coding assistant working on Gitea repo: ${owner}/${repo}.

${context}

**Task:** ${userPrompt || 'Analyze this issue and implement any necessary code changes.'}

IMPORTANT: If you need to create or modify files in the repository, output them in this exact format:

FILE: path/to/file.ext
\`\`\`
full file content here
\`\`\`

You can output multiple FILE: blocks for multiple files. After the file blocks, write a summary of what you changed and why.

If no code changes are needed, just provide your analysis and recommendations in markdown.`;

    console.log(`[Claude] Calling for ${owner}/${repo}#${issueNumber}...`);
    const response = await callClaude(fullPrompt);
    console.log(`[Claude] Response: ${response.length} chars`);

    // Apply file changes
    const fileChanges = parseFileChanges(response);
    let changeSummary = '';
    if (fileChanges.length > 0) {
      console.log(`[Claude] Applying ${fileChanges.length} file changes`);
      const results = await applyFileChanges(owner, repo, fileChanges, issueNumber);
      changeSummary = `\n\n---\n**Changes applied to repository:**\n${results.map(r => `- ${r}`).join('\n')}`;
    }

    // Clean response for comment
    let cleanResponse = response;
    cleanResponse = cleanResponse.replace(/FILE:\s*\S+\s*\n```[\w]*\n[\s\S]*?```/g, '').trim();
    cleanResponse = cleanResponse.replace(/DELETE:\s*\S+/g, '').trim();
    cleanResponse = cleanResponse.replace(/\n{3,}/g, '\n\n');

    const formatted = `> **Claude** responding to @${triggerUser}:\n\n${cleanResponse || response}${changeSummary}`;

    if (commentId) {
      await editComment(owner, repo, commentId, formatted);
    }

    event.status = 'complete';
    event.response = response;
    console.log(`[Process] Complete: ${owner}/${repo}#${issueNumber}`);
  } catch (err) {
    console.error(`[Process] Error: ${err.message}`);
    if (commentId) {
      await editComment(owner, repo, commentId,
        `> **Claude** encountered an error:\n\n\`${err.message}\``).catch(() => {});
    }
    event.status = 'error';
    event.response = err.message;
  }
}

// --- Polling loop ---

async function pollGitea() {
  pollCount++;
  lastPollTime = new Date().toISOString();

  try {
    const repos = await getMyRepos();
    console.log(`[Poll #${pollCount}] Checking ${repos.length} repos...`);

    for (const repo of repos) {
      const owner = repo.owner?.login || repo.owner?.username;
      const repoName = repo.name;
      if (!owner || !repoName) continue;

      try {
        const issues = await getRepoIssues(owner, repoName);

        for (const issue of issues) {
          const issueKey = `${owner}/${repoName}#${issue.number}`;
          const isPR = !!issue.pull_request;

          // Check if this is a new issue we haven't processed
          if (!processedIssues.has(issueKey)) {
            processedIssues.add(issueKey);
            // Skip issues older than 5 minutes (avoid processing old issues on startup)
            const issueAge = Date.now() - new Date(issue.created_at).getTime();
            if (issueAge < 5 * 60 * 1000) {
              console.log(`[Poll] New issue: ${issueKey} - ${issue.title}`);
              processWithClaude(owner, repoName, issue.number, isPR,
                `${issue.title}\n\n${issue.body || ''}`,
                issue.user?.login || 'unknown', 'issue_opened').catch(console.error);
              continue; // Don't also check comments for this issue this cycle
            }
          }

          // Check comments for @claude mentions
          try {
            const comments = await getIssueComments(owner, repoName, issue.number);
            for (const comment of comments) {
              const commentKey = `comment-${comment.id}`;
              if (processedComments.has(commentKey)) continue;
              processedComments.add(commentKey);

              // Skip our own comments
              if (comment.body?.includes('**Claude**')) continue;

              // Skip comments older than 5 minutes
              const commentAge = Date.now() - new Date(comment.created_at).getTime();
              if (commentAge > 5 * 60 * 1000) continue;

              // Check for @claude mention
              if (comment.body?.includes('@claude')) {
                const prompt = comment.body.replace(/@claude\s*/gi, '').trim();
                const user = comment.user?.login || comment.user?.username || 'unknown';
                console.log(`[Poll] @claude mention: ${issueKey} by ${user}`);
                processWithClaude(owner, repoName, issue.number, isPR,
                  prompt, user, 'comment_trigger').catch(console.error);
              }
            }
          } catch (err) {
            // Silently skip comment fetch errors
          }
        }
      } catch (err) {
        console.error(`[Poll] Error checking ${owner}/${repoName}: ${err.message}`);
      }
    }
  } catch (err) {
    console.error(`[Poll] Error fetching repos: ${err.message}`);
    pollErrors.push({ time: new Date().toISOString(), error: err.message });
    if (pollErrors.length > 20) pollErrors.shift();
  }
}

// Start polling
console.log(`Starting poll loop (every ${POLL_INTERVAL / 1000}s)...`);
// Initial poll after 5s to let things settle
setTimeout(() => {
  pollGitea();
  setInterval(pollGitea, POLL_INTERVAL);
}, 5000);

// --- Manual trigger handler ---

async function handleManualTrigger(body) {
  let payload;
  try { payload = JSON.parse(body); } catch { return { status: 400, data: { error: 'Invalid JSON' } }; }

  const { owner, repo, issue, prompt } = payload;
  if (!prompt) return { status: 400, data: { error: 'Missing prompt' } };

  if (owner && repo && issue) {
    processWithClaude(owner, repo, parseInt(issue), false, prompt, 'dashboard', 'manual').catch(console.error);
  } else {
    const event = addEvent({
      type: 'manual',
      repo: `${owner || '?'}/${repo || '?'}`,
      issue: null,
      isPR: false,
      user: 'dashboard',
      prompt,
      status: 'working',
      response: null,
    });

    callClaude(prompt).then(response => {
      event.status = 'complete';
      event.response = response;
    }).catch(err => {
      event.status = 'error';
      event.response = err.message;
    });
  }

  return { status: 202, data: { message: 'Processing' } };
}

// --- HTTP Server ---

let indexHtml = '';
try { indexHtml = fs.readFileSync('/app/public/index.html', 'utf-8'); } catch { indexHtml = '<h1>Dashboard not found</h1>'; }

const server = http.createServer(async (req, res) => {
  const url = new URL(req.url, `http://${req.headers.host}`);
  const path = url.pathname;

  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type');
  if (req.method === 'OPTIONS') { res.writeHead(204); return res.end(); }

  if (path === '/' && req.method === 'GET') {
    res.writeHead(200, { 'Content-Type': 'text/html' });
    return res.end(indexHtml);
  }

  if (path === '/health' && req.method === 'GET') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    return res.end(JSON.stringify({
      status: 'ok',
      pollCount,
      lastPoll: lastPollTime,
      events: events.length,
      trackedIssues: processedIssues.size,
      trackedComments: processedComments.size,
      errors: pollErrors.length,
    }));
  }

  if (path === '/api/events' && req.method === 'GET') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    return res.end(JSON.stringify(events));
  }

  if (path === '/api/status' && req.method === 'GET') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    return res.end(JSON.stringify({
      pollCount,
      lastPoll: lastPollTime,
      interval: POLL_INTERVAL / 1000,
      trackedIssues: processedIssues.size,
      trackedComments: processedComments.size,
      recentErrors: pollErrors.slice(-5),
    }));
  }

  // Keep webhook endpoint as backup
  if (path === '/webhook' && req.method === 'POST') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    return res.end(JSON.stringify({ message: 'Webhook received (polling mode active)' }));
  }

  if (path === '/api/trigger' && req.method === 'POST') {
    let body = '';
    req.on('data', (c) => body += c);
    req.on('end', async () => {
      try {
        const result = await handleManualTrigger(body);
        res.writeHead(result.status, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify(result.data));
      } catch (err) {
        res.writeHead(500, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: err.message }));
      }
    });
    return;
  }

  res.writeHead(404, { 'Content-Type': 'application/json' });
  res.end(JSON.stringify({ error: 'Not found' }));
});

server.listen(PORT, () => {
  console.log(`Claude Code Gitea Actions server running on port ${PORT}`);
  console.log(`Mode: Polling (every ${POLL_INTERVAL / 1000}s)`);
  console.log(`Gitea URL: ${GITEA_URL}`);
  console.log(`CMD API URL: ${CMD_API_URL}`);
});
