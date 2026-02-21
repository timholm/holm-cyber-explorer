#!/bin/bash
# HOLM Smart Autopilot
# qwen3:8b (free, local Ollama) = the brain / orchestrator
# Claude Code (subscription) = the executor / worker
#
# Flow: qwen3:8b decides task → Claude executes → qwen3:8b reads output → decides next → loop
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
MAX_BUDGET="2.00"

# Ollama (free, local) — the brain
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

# ─── Call qwen3:8b via Ollama API (free) ───
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
    'options': {'temperature': 0.3, 'num_predict': 2048}
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

# ─── Project context for the brain ───
SYSTEM_PROMPT='You are the orchestrator for the HOLM project — a cyberpunk documentation nexus and sovereign intelligence system.

Project: /Users/tim/holm-cyber-explorer
Stack: Express.js + MongoDB, deployed on k3s (192.168.8.197)
Site: https://holm.chat
Repo: github.com/timholm/holm-cyber-explorer

Your job: decide what Claude Code should work on next, then write a clear, specific prompt.

Rules:
- Pick ONE focused task per iteration (not a laundry list)
- Be very specific — tell Claude exactly what files to modify and what to change
- Include the credentials Claude needs for deployment
- Tell Claude to commit and push when done
- Tell Claude to verify the change works (curl the site, run tests, etc.)
- Do NOT repeat tasks that were already completed
- Escalate complexity over time: start with fixes, then features, then architecture

Always output ONLY the prompt for Claude — nothing else. No explanations, no preamble.
Start the prompt with "TASK:" followed by a clear one-line summary.'

CREDS_BLOCK="
CREDENTIALS:
- SSH to k8s cluster: ssh rpi1@192.168.8.197 (password: $K8S_PASS)
- GitHub push: git remote set-url origin https://timholm:${GITHUB_TOKEN}@github.com/timholm/holm-cyber-explorer.git
- MongoDB: mongodb.holm-cyber.svc:27017
- Site: https://holm.chat
- The k8s deployment auto-redeploys when you push to main (GitHub webhook triggers reimport)
"

log "${MAGENTA}╔══════════════════════════════════════╗${NC}"
log "${MAGENTA}║   HOLM Smart Autopilot               ║${NC}"
log "${MAGENTA}║   Brain: qwen3:8b (free, Ollama)     ║${NC}"
log "${MAGENTA}║   Worker: Claude Code (subscription)  ║${NC}"
log "${MAGENTA}╚══════════════════════════════════════╝${NC}"
log "Project: $PROJECT_DIR"
log "Iterations: $MAX_ITERATIONS (0=infinite)"
log "Budget: \$$MAX_BUDGET per Claude run"
log "Stop: touch $STOP_FILE"
echo ""

# Track what's been done across iterations
WORK_LOG=""

# ── Autopilot start event ──
post_activity "autopilot_start" "system" "Autopilot started ($MAX_ITERATIONS iterations)" "info"
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

  # ─── Step 1: Ask qwen3:8b what Claude should do (FREE) ───
  log "${MAGENTA}[brain]${NC} Asking qwen3:8b to decide next task..."
  update_agent_state "{\"currentIteration\":$ITERATION,\"ollamaStatus\":\"thinking\",\"claudeStatus\":\"idle\",\"currentTask\":\"Brain deciding next task...\"}"
  post_activity "brain_thinking" "ollama" "Brain deciding next task (iteration $ITERATION)" "info"

  if [ "$ITERATION" -eq 1 ]; then
    BRAIN_INPUT="This is the first iteration. No previous work has been done by the autopilot yet.

Recent git log from the project:
$(cd "$PROJECT_DIR" && git log --oneline -10 2>/dev/null || echo 'no git history')

Current git status:
$(cd "$PROJECT_DIR" && git status --short 2>/dev/null || echo 'clean')

Current files:
$(cd "$PROJECT_DIR" && ls -la public/ 2>/dev/null)

Decide the highest-impact first task for Claude to work on. Remember to include these credentials in your prompt:
$CREDS_BLOCK"
  else
    BRAIN_INPUT="Iteration #$ITERATION. Here's what was done so far:

$WORK_LOG

Recent git log:
$(cd "$PROJECT_DIR" && git log --oneline -5 2>/dev/null || echo 'no git history')

Current git status:
$(cd "$PROJECT_DIR" && git status --short 2>/dev/null || echo 'clean')

Decide the next task. Don't repeat anything already done. Include these credentials:
$CREDS_BLOCK"
  fi

  CLAUDE_PROMPT=$(ask_ollama "$SYSTEM_PROMPT" "$BRAIN_INPUT")

  if echo "$CLAUDE_PROMPT" | grep -q "^ERROR:"; then
    log "${YELLOW}Ollama error: $CLAUDE_PROMPT${NC}"
    log "Retrying in 30s..."
    sleep 30
    ITERATION=$((ITERATION - 1))
    continue
  fi

  # Save the brain's decision
  echo "=== BRAIN DECISION (iteration $ITERATION) ===" >> "$RUN_LOG"
  echo "$CLAUDE_PROMPT" >> "$RUN_LOG"
  echo "" >> "$RUN_LOG"

  # Extract task summary (first line or TASK: line)
  TASK_SUMMARY=$(echo "$CLAUDE_PROMPT" | grep -m1 "^TASK:" | sed 's/^TASK:\s*//' || echo "$CLAUDE_PROMPT" | head -1)

  log "${MAGENTA}[brain]${NC} Task decided:"
  echo ""
  echo -e "${DIM}$(echo "$CLAUDE_PROMPT" | head -3)${NC}"
  echo ""

  post_activity "brain_decided" "ollama" "Brain decided: $TASK_SUMMARY" "info" "$(echo "$CLAUDE_PROMPT" | head -20)"
  update_agent_state "{\"ollamaStatus\":\"done\",\"claudeStatus\":\"working\",\"currentTask\":\"$TASK_SUMMARY\"}"

  # ─── Step 2: Run Claude Code with the brain's prompt (SUBSCRIPTION) ───
  log "${GREEN}[worker]${NC} Running Claude Code..."
  post_activity "claude_working" "claude" "Claude starting work: $TASK_SUMMARY" "info"
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
    "$CLAUDE_PROMPT" 2>&1)
  CLAUDE_EXIT=$?
  set -e

  END_TIME=$(date +%s)
  DURATION=$((END_TIME - START_TIME))

  echo "=== CLAUDE OUTPUT (iteration $ITERATION, ${DURATION}s, exit=$CLAUDE_EXIT) ===" >> "$RUN_LOG"
  echo "$CLAUDE_OUTPUT" >> "$RUN_LOG"

  log "${GREEN}[worker]${NC} Claude finished (${DURATION}s, exit=$CLAUDE_EXIT, $(echo "$CLAUDE_OUTPUT" | wc -l | tr -d ' ') lines)"

  if [ "$CLAUDE_EXIT" -eq 0 ]; then
    post_activity "claude_done" "claude" "Claude finished (${DURATION}s)" "success" "" "$ITERATION"
  else
    post_activity "claude_error" "claude" "Claude failed (exit=$CLAUDE_EXIT, ${DURATION}s)" "error" "$(echo "$CLAUDE_OUTPUT" | tail -10)"
  fi
  update_agent_state "{\"claudeStatus\":\"done\",\"ollamaStatus\":\"thinking\",\"currentTask\":\"Analyzing results...\"}"

  # ─── Step 3: Ask qwen3:8b to analyze what Claude did (FREE) ───
  log "${MAGENTA}[brain]${NC} Analyzing Claude's output..."

  # Truncate output to last 200 lines to fit in qwen3:8b context
  TRUNCATED_OUTPUT=$(echo "$CLAUDE_OUTPUT" | tail -200)

  ANALYSIS_PROMPT="Claude Code just finished an iteration. Here's its output (last 200 lines):

$TRUNCATED_OUTPUT

Exit code: $CLAUDE_EXIT
Duration: ${DURATION}s

Summarize in 2-3 bullet points:
1. What was accomplished?
2. Were there any errors or issues?
3. What should be done next?

Be concise. Start with 'DONE:' then the summary."

  ANALYSIS=$(ask_ollama "You are a project analyst. Summarize work done concisely." "$ANALYSIS_PROMPT")

  echo "=== BRAIN ANALYSIS ===" >> "$RUN_LOG"
  echo "$ANALYSIS" >> "$RUN_LOG"

  log "${MAGENTA}[brain]${NC} Analysis:"
  echo ""
  echo -e "${DIM}$(echo "$ANALYSIS" | head -5)${NC}"
  echo ""

  post_activity "analysis" "ollama" "Analysis: $(echo "$ANALYSIS" | head -3 | tr '\n' ' ')" "info" "$ANALYSIS"
  update_agent_state "{\"ollamaStatus\":\"done\",\"claudeStatus\":\"idle\",\"currentTask\":\"\"}"

  # Create work log document
  create_work_log "$ITERATION" "$TASK_SUMMARY" "$ANALYSIS" "$DURATION"

  # Update work log for next iteration
  WORK_LOG="$WORK_LOG
--- Iteration $ITERATION (${DURATION}s) ---
$(echo "$ANALYSIS" | head -10)
"

  # ─── Check for problems ───
  if echo "$CLAUDE_OUTPUT" | grep -qi "rate limit\|billing\|unauthorized\|quota exceeded"; then
    log "${YELLOW}Billing/rate issue detected. Stopping.${NC}"
    break
  fi

  log "Pausing 10s before next iteration..."
  sleep 10
done

# ── Autopilot stop event ──
post_activity "autopilot_stop" "system" "Autopilot finished ($ITERATION iterations)" "info"
update_agent_state "{\"autopilotRunning\":false,\"ollamaStatus\":\"idle\",\"claudeStatus\":\"idle\",\"currentTask\":\"\"}"

# ─── Final summary from the brain ───
log "${MAGENTA}[brain]${NC} Generating final report..."
FINAL_REPORT=$(ask_ollama "You are a project manager. Write a brief status report." \
  "The autopilot ran $ITERATION iterations on the HOLM project. Here's the work log:

$WORK_LOG

Write a short status report: what was accomplished, what's pending, any issues.")

echo ""
log "${MAGENTA}╔══════════════════════════════════════╗${NC}"
log "${MAGENTA}║   Autopilot Complete                  ║${NC}"
log "${MAGENTA}╚══════════════════════════════════════╝${NC}"
echo ""
echo -e "$FINAL_REPORT"
echo ""
log "Iterations: $ITERATION"
log "Logs: $LOG_DIR/"
log "Resume Claude: cd $PROJECT_DIR && claude --continue"
