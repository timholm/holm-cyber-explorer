#!/bin/bash
# HOLM Autopilot v2
# Claude Code (subscription) = the brain AND the worker — decides tasks, executes, commits
# qwen3:8b (free, Ollama) = the monitor — checks Claude's work, flags issues, keeps it honest
#
# Flow: Claude decides + executes → Ollama reviews output → if drift/error, Ollama writes correction → loop
#
# Usage: ./autopilot.sh [max_iterations]
# Stop:  touch /tmp/autopilot-stop  or  Ctrl+C

set -euo pipefail
unset CLAUDECODE 2>/dev/null || true

PROJECT_DIR="/Users/tim/holm-cyber-explorer"
LOG_DIR="$PROJECT_DIR/.autopilot"
STOP_FILE="/tmp/autopilot-stop"
MAX_ITERATIONS="${1:-5}"
ITERATION=0
MAX_BUDGET="5.00"

# Ollama (free, local) — the monitor
OLLAMA_URL="http://192.168.8.230:11434"
OLLAMA_MODEL="qwen3:8b"

mkdir -p "$LOG_DIR"
rm -f "$STOP_FILE"

# Load creds
source "$LOG_DIR/creds.env"

HOLM_API_URL="https://holm.chat"
HOLM_API_KEY="${HOLM_API_KEY:-}"

CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
DIM='\033[2m'
NC='\033[0m'

log() { echo -e "${CYAN}[autopilot]${NC} $(date '+%H:%M:%S') $*"; }

# ─── Call qwen3:8b via Ollama API (free) — monitor role ───
ask_ollama() {
  local system_prompt="$1"
  local user_prompt="$2"
  local retries=3

  local payload
  payload=$(python3 -c "
import json, sys
print(json.dumps({
    'model': sys.argv[3],
    'messages': [
        {'role': 'system', 'content': sys.argv[1]},
        {'role': 'user', 'content': sys.argv[2]}
    ],
    'stream': False,
    'options': {'temperature': 0.3, 'num_predict': 1024}
}))
" "$system_prompt" "$user_prompt" "$OLLAMA_MODEL")

  local attempt=0
  while [ $attempt -lt $retries ]; do
    attempt=$((attempt + 1))

    local response
    response=$(curl -s --max-time 300 "$OLLAMA_URL/api/chat" \
      -H "Content-Type: application/json" \
      -d "$payload" 2>/dev/null) || true

    if [ -z "$response" ]; then
      log "${YELLOW}Ollama returned empty response (attempt $attempt/$retries)${NC}"
      sleep 5
      continue
    fi

    local content
    content=$(echo "$response" | python3 -c "
import json, sys, re
try:
    data = json.load(sys.stdin)
    text = data.get('message', {}).get('content', '')
    # Strip qwen3 think tokens
    text = re.sub(r'<think>.*?</think>', '', text, flags=re.DOTALL).strip()
    if text:
        print(text)
    else:
        print('ERROR: Empty content after parsing')
except Exception as e:
    print(f'ERROR: {e}')
" 2>/dev/null)

    if [ -z "$content" ] || echo "$content" | grep -q "^ERROR:"; then
      log "${YELLOW}Ollama parse issue (attempt $attempt/$retries): $content${NC}"
      sleep 5
      continue
    fi

    echo "$content"
    return 0
  done

  echo "ERROR: Ollama failed after $retries attempts"
  return 1
}

# ─── Dashboard helpers (fire-and-forget) ───
post_activity() {
  local type="$1" agent="$2" message="$3" status="${4:-info}" detail="${5:-}"
  [ -z "$HOLM_API_KEY" ] && return 0
  curl -s --max-time 5 -X POST "$HOLM_API_URL/api/activity" \
    -H "Content-Type: application/json" \
    -H "X-Api-Key: $HOLM_API_KEY" \
    -d "$(python3 -c "
import json, sys
print(json.dumps({
    'type': sys.argv[1], 'agent': sys.argv[2], 'message': sys.argv[3],
    'status': sys.argv[4], 'detail': sys.argv[5][:2000] if len(sys.argv) > 5 else '',
    'iteration': int(sys.argv[6]) if len(sys.argv) > 6 else 0
}))
" "$type" "$agent" "$message" "$status" "$detail" "${ITERATION:-0}")" >/dev/null 2>&1 &
}

update_agent_state() {
  local json_data="$1"
  [ -z "$HOLM_API_KEY" ] && return 0
  curl -s --max-time 5 -X PUT "$HOLM_API_URL/api/agent-state" \
    -H "Content-Type: application/json" \
    -H "X-Api-Key: $HOLM_API_KEY" \
    -d "$json_data" >/dev/null 2>&1 &
}

create_work_log() {
  local iteration="$1" summary="$2" analysis="${3:-}" duration="${4:-0}"
  [ -z "$HOLM_API_KEY" ] && return 0
  curl -s --max-time 10 -X POST "$HOLM_API_URL/api/activity/doc" \
    -H "Content-Type: application/json" \
    -H "X-Api-Key: $HOLM_API_KEY" \
    -d "$(python3 -c "
import json, sys
print(json.dumps({
    'iteration': int(sys.argv[1]),
    'summary': sys.argv[2],
    'analysis': sys.argv[3] if len(sys.argv) > 3 else '',
    'duration': int(sys.argv[4]) if len(sys.argv) > 4 else 0
}))
" "$iteration" "$summary" "$analysis" "$duration")" >/dev/null 2>&1 &
}

# ─── Credentials block injected into Claude's prompt ───
CREDS_BLOCK="
CREDENTIALS:
- SSH to k8s cluster: ssh rpi1@192.168.8.197 (password: $K8S_PASS)
- GitHub push: git remote set-url origin https://timholm:${GITHUB_TOKEN}@github.com/timholm/holm-cyber-explorer.git
- MongoDB: mongodb.holm-cyber.svc:27017
- Site: https://holm.chat
- holm.chat API key: $HOLM_API_KEY
- The k8s deployment auto-redeploys when you push to main (GitHub webhook triggers reimport)
"

# ─── Claude's master prompt — it decides AND executes ───
CLAUDE_SYSTEM="You are the Chief Systems Architect and sole developer for HOLM — a cyberpunk documentation nexus and sovereign intelligence system.

Project: /Users/tim/holm-cyber-explorer
Stack: Express.js + MongoDB, deployed on k3s (Raspberry Pi cluster at 192.168.8.197)
Site: https://holm.chat
Repo: github.com/timholm/holm-cyber-explorer

You have full autonomy. Each iteration, YOU decide the highest-impact task and execute it completely:
1. Assess the current state (git log, git status, CLAUDE.md, roadmap tasks)
2. Pick ONE focused, high-impact task — prioritize from the roadmap at $HOLM_API_URL/api/tasks
3. Implement it fully — edit files, test, commit, push
4. Verify the change works (curl the site, check endpoints, etc.)
5. Update the roadmap task status via the API when done

Rules:
- Do NOT repeat work already done (check git log)
- Always commit and push when done
- Always verify your changes work on the live site
- Be ambitious — you have a \$${MAX_BUDGET} budget per iteration, use it
- Fix bugs before adding features
- Keep the cyberpunk aesthetic (dark theme, cyan/magenta neon, monospace)
- This is a REAL production site — test before pushing

$CREDS_BLOCK"

# ─── Ollama monitor prompt — watches Claude, flags drift ───
MONITOR_PROMPT='You are a quality monitor for the HOLM autopilot system. Your ONLY job is to review Claude'\''s work output and determine:

1. VERDICT: Did Claude do useful, real work? (YES/NO/PARTIAL)
2. SUMMARY: 2-3 bullet points of what was accomplished
3. ISSUES: Any errors, regressions, or problems detected
4. CORRECTION: If Claude drifted or failed, write a one-line correction directive. If work was good, write "NONE"

Rules:
- Be harsh. If Claude talked about what it WOULD do but didn'\''t actually write code or make commits, verdict is NO.
- If Claude made commits and pushed, verdict is YES.
- If there were errors but some progress, verdict is PARTIAL.
- Keep summaries to 2-3 lines max.
- Output in this exact format:
  VERDICT: YES|NO|PARTIAL
  SUMMARY: ...
  ISSUES: ...
  CORRECTION: ...'

log "${MAGENTA}╔══════════════════════════════════════╗${NC}"
log "${MAGENTA}║   HOLM Autopilot v2                   ║${NC}"
log "${MAGENTA}║   Claude: brain + worker (subscription)║${NC}"
log "${MAGENTA}║   Ollama: monitor (free)               ║${NC}"
log "${MAGENTA}╚══════════════════════════════════════╝${NC}"
log "Project: $PROJECT_DIR"
log "Iterations: $MAX_ITERATIONS (0=infinite)"
log "Budget: \$$MAX_BUDGET per Claude run"
log "Stop: touch $STOP_FILE"
echo ""

# Track what's been done across iterations
WORK_LOG=""
CORRECTION=""

# ── Autopilot start event ──
post_activity "autopilot_start" "system" "Autopilot v2 started — Claude decides + works, Ollama monitors ($MAX_ITERATIONS iterations)" "info"
update_agent_state "{\"autopilotRunning\":true,\"currentIteration\":0,\"maxIterations\":$MAX_ITERATIONS,\"ollamaStatus\":\"idle\",\"claudeStatus\":\"idle\",\"currentTask\":\"\",\"startedAt\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"}"

while true; do
  if [ -f "$STOP_FILE" ]; then
    log "${YELLOW}Stop file detected. Shutting down.${NC}"
    rm -f "$STOP_FILE"
    break
  fi

  if [ "$MAX_ITERATIONS" -gt 0 ] && [ "$ITERATION" -ge "$MAX_ITERATIONS" ]; then
    log "${YELLOW}Max iterations reached ($MAX_ITERATIONS). Done.${NC}"
    break
  fi

  ITERATION=$((ITERATION + 1))
  RUN_LOG="$LOG_DIR/run-$(printf '%03d' $ITERATION)-$(date '+%Y%m%d-%H%M%S').log"

  log "${GREEN}══════ Iteration $ITERATION/$MAX_ITERATIONS ══════${NC}"

  # ─── Step 1: Build Claude's prompt with full context ───
  update_agent_state "{\"currentIteration\":$ITERATION,\"ollamaStatus\":\"idle\",\"claudeStatus\":\"working\",\"currentTask\":\"Claude deciding and executing...\"}"
  post_activity "claude_working" "claude" "Claude starting iteration $ITERATION (deciding + executing)" "info"

  CLAUDE_INPUT="$CLAUDE_SYSTEM

--- CURRENT STATE ---
Recent git log:
$(cd "$PROJECT_DIR" && git log --oneline -10 2>/dev/null || echo 'no git history')

Git status:
$(cd "$PROJECT_DIR" && git status --short 2>/dev/null || echo 'clean')

Previous work this session:
${WORK_LOG:-None yet — this is the first iteration.}
"

  # If the monitor flagged a correction, inject it
  if [ -n "$CORRECTION" ] && [ "$CORRECTION" != "NONE" ]; then
    CLAUDE_INPUT="$CLAUDE_INPUT
--- MONITOR CORRECTION FROM LAST ITERATION ---
$CORRECTION
--- END CORRECTION ---
Address this correction before moving on to new work.
"
    log "${YELLOW}[monitor]${NC} Injecting correction: $CORRECTION"
    post_activity "monitor_correction" "ollama" "Monitor correction: $CORRECTION" "warning"
  fi

  CLAUDE_INPUT="$CLAUDE_INPUT
Now: assess the project state, pick the highest-impact task, implement it, commit, push, and verify. Go."

  log "${GREEN}[claude]${NC} Claude is thinking and working..."

  echo "=== CLAUDE PROMPT (iteration $ITERATION) ===" >> "$RUN_LOG"
  echo "$CLAUDE_INPUT" >> "$RUN_LOG"
  echo "" >> "$RUN_LOG"

  # ─── Step 2: Run Claude — it decides AND executes ───
  START_TIME=$(date +%s)

  set +e
  if [ "$ITERATION" -eq 1 ]; then
    RESUME_FLAG=""
  else
    RESUME_FLAG="--continue"
  fi

  CLAUDE_OUTPUT=$(cd "$PROJECT_DIR" && claude -p \
    $RESUME_FLAG \
    --dangerously-skip-permissions \
    --max-budget-usd "$MAX_BUDGET" \
    "$CLAUDE_INPUT" 2>&1)
  CLAUDE_EXIT=$?
  set -e

  END_TIME=$(date +%s)
  DURATION=$((END_TIME - START_TIME))

  echo "=== CLAUDE OUTPUT (iteration $ITERATION, ${DURATION}s, exit=$CLAUDE_EXIT) ===" >> "$RUN_LOG"
  echo "$CLAUDE_OUTPUT" >> "$RUN_LOG"

  LINES=$(echo "$CLAUDE_OUTPUT" | wc -l | tr -d ' ')
  log "${GREEN}[claude]${NC} Claude finished (${DURATION}s, exit=$CLAUDE_EXIT, $LINES lines)"

  if [ "$CLAUDE_EXIT" -eq 0 ]; then
    post_activity "claude_done" "claude" "Claude finished iteration $ITERATION (${DURATION}s, $LINES lines)" "success"
  else
    post_activity "claude_error" "claude" "Claude failed iteration $ITERATION (exit=$CLAUDE_EXIT, ${DURATION}s)" "error" "$(echo "$CLAUDE_OUTPUT" | tail -10)"
  fi

  # Extract what Claude said it did (first few and last few lines)
  TASK_SUMMARY=$(echo "$CLAUDE_OUTPUT" | head -5 | tr '\n' ' ' | cut -c1-200)

  # ─── Step 3: Ollama monitors Claude's output (FREE) ───
  log "${MAGENTA}[monitor]${NC} Ollama reviewing Claude's work..."
  update_agent_state "{\"claudeStatus\":\"done\",\"ollamaStatus\":\"monitoring\",\"currentTask\":\"Monitor reviewing Claude output...\"}"
  post_activity "monitor_reviewing" "ollama" "Monitor reviewing Claude's work" "info"

  TRUNCATED_OUTPUT=$(echo "$CLAUDE_OUTPUT" | tail -200)

  REVIEW_INPUT="Claude Code just ran iteration $ITERATION on the HOLM project. Review its output:

--- CLAUDE OUTPUT (last 200 lines) ---
$TRUNCATED_OUTPUT
--- END OUTPUT ---

Exit code: $CLAUDE_EXIT
Duration: ${DURATION}s

Check:
- Did Claude actually write code and make commits? (look for git commit/push output)
- Did Claude verify changes work? (look for curl/test output)
- Any errors or regressions?
- Is Claude staying focused on high-impact work?

Respond in the exact format specified."

  REVIEW=$(ask_ollama "$MONITOR_PROMPT" "$REVIEW_INPUT")

  echo "=== MONITOR REVIEW ===" >> "$RUN_LOG"
  echo "$REVIEW" >> "$RUN_LOG"

  # Parse the review
  VERDICT=$(echo "$REVIEW" | grep -i "^VERDICT:" | head -1 | sed 's/^VERDICT:\s*//' | tr '[:lower:]' '[:upper:]' | tr -d ' ')
  CORRECTION=$(echo "$REVIEW" | grep -i "^CORRECTION:" | head -1 | sed 's/^CORRECTION:\s*//')
  SUMMARY=$(echo "$REVIEW" | grep -i "^SUMMARY:" | head -1 | sed 's/^SUMMARY:\s*//')

  log "${MAGENTA}[monitor]${NC} Verdict: $VERDICT"
  echo ""
  echo -e "${DIM}$(echo "$REVIEW" | head -6)${NC}"
  echo ""

  if [ "$VERDICT" = "NO" ]; then
    post_activity "monitor_verdict" "ollama" "VERDICT: NO — Claude did not produce useful work. Correction: $CORRECTION" "warning" "$REVIEW"
  elif [ "$VERDICT" = "PARTIAL" ]; then
    post_activity "monitor_verdict" "ollama" "VERDICT: PARTIAL — $SUMMARY" "info" "$REVIEW"
  else
    post_activity "monitor_verdict" "ollama" "VERDICT: YES — $SUMMARY" "success" "$REVIEW"
    CORRECTION=""  # Clear correction if work was good
  fi

  update_agent_state "{\"ollamaStatus\":\"done\",\"claudeStatus\":\"idle\",\"currentTask\":\"\"}"

  # Create work log document
  create_work_log "$ITERATION" "${TASK_SUMMARY:-Claude iteration $ITERATION}" "$REVIEW" "$DURATION"

  # Update work log for next iteration
  WORK_LOG="$WORK_LOG
--- Iteration $ITERATION (${DURATION}s, verdict: $VERDICT) ---
$(echo "$REVIEW" | head -6)
"

  # ─── Check for billing problems ───
  if echo "$CLAUDE_OUTPUT" | grep -qi "rate limit\|billing\|unauthorized\|quota exceeded"; then
    log "${YELLOW}Billing/rate issue detected. Stopping.${NC}"
    break
  fi

  # Short pause — keep Claude running hot
  log "Next iteration in 5s..."
  sleep 5
done

# ── Autopilot stop event ──
post_activity "autopilot_stop" "system" "Autopilot v2 finished ($ITERATION iterations)" "info"
update_agent_state "{\"autopilotRunning\":false,\"ollamaStatus\":\"idle\",\"claudeStatus\":\"idle\",\"currentTask\":\"\"}"

# ─── Final summary from the monitor ───
log "${MAGENTA}[monitor]${NC} Generating final report..."
FINAL_REPORT=$(ask_ollama "You are a project monitor. Write a brief status report." \
  "The HOLM autopilot v2 ran $ITERATION iterations. Claude decided tasks and executed them. Here's the session log:

$WORK_LOG

Write a concise status report: what Claude accomplished, any issues the monitor flagged, what's pending.")

echo ""
log "${MAGENTA}╔══════════════════════════════════════╗${NC}"
log "${MAGENTA}║   Autopilot v2 Complete               ║${NC}"
log "${MAGENTA}╚══════════════════════════════════════╝${NC}"
echo ""
echo -e "$FINAL_REPORT"
echo ""
log "Iterations: $ITERATION"
log "Logs: $LOG_DIR/"
log "Resume Claude: cd $PROJECT_DIR && claude --continue"
