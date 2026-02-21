#!/bin/bash
# HOLM Multi-Agent Orchestrator
# Replaces autopilot.sh with parallel Claude Code workers
#
# Architecture:
#   Coordinator loop (this script) → spawns N Claude workers in parallel
#   Ollama (free, local) → decomposes directives into tasks, generates self-directives
#   Claude Code (subscription) → executes tasks autonomously
#
# Usage: ./orchestrator.sh [max_workers]
# Stop:  touch /tmp/orchestrator-stop  or  Ctrl+C

set -euo pipefail
unset CLAUDECODE 2>/dev/null || true

PROJECT_DIR="/Users/tim/holm-cyber-explorer"
LOG_DIR="$PROJECT_DIR/.orchestrator"
STOP_FILE="/tmp/orchestrator-stop"
MAX_WORKERS="${1:-5}"
MAX_BUDGET_PER_TASK="5.00"
HEARTBEAT_TIMEOUT=120
COORDINATOR_INTERVAL=10
IDLE_CYCLES_BEFORE_SELF_DIRECT=10

# Rate limit: ~900 prompts per 5hr window
RATE_LIMIT_MAX=900
RATE_LIMIT_SLOW_PCT=80
RATE_LIMIT_STOP_PCT=100

# Ollama (free, local)
OLLAMA_URL="http://192.168.8.230:11434"
OLLAMA_MODEL="qwen3:8b"

mkdir -p "$LOG_DIR/workers" "$LOG_DIR/streams"
rm -f "$STOP_FILE"

# Load creds
source "$PROJECT_DIR/.autopilot/creds.env"

HOLM_API_URL="https://holm.chat"
HOLM_API_KEY="${HOLM_API_KEY:-}"

# Colors
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
DIM='\033[2m'
NC='\033[0m'

log() { echo -e "${CYAN}[orch]${NC} $(date '+%H:%M:%S') $*"; }
log_worker() { echo -e "${GREEN}[worker-$1]${NC} $(date '+%H:%M:%S') ${@:2}"; }

# ─── Tracking ───
declare -A WORKER_PIDS
declare -A WORKER_TASKS
declare -A WORKER_SESSIONS
IDLE_CYCLES=0
TOTAL_ITERATIONS=0
TOTAL_COST=0
PROMPTS_USED=0
WINDOW_START=$(date +%s)

# ─── API helpers ───
api_get() {
  curl -s --max-time 10 "$HOLM_API_URL$1" 2>/dev/null || echo '[]'
}

api_post() {
  curl -s --max-time 10 -X POST "$HOLM_API_URL$1" \
    -H "Content-Type: application/json" \
    -H "X-Api-Key: $HOLM_API_KEY" \
    -d "$2" 2>/dev/null || echo '{}'
}

api_put() {
  curl -s --max-time 10 -X PUT "$HOLM_API_URL$1" \
    -H "Content-Type: application/json" \
    -H "X-Api-Key: $HOLM_API_KEY" \
    -d "$2" 2>/dev/null || echo '{}'
}

api_delete() {
  curl -s --max-time 10 -X DELETE "$HOLM_API_URL$1" \
    -H "X-Api-Key: $HOLM_API_KEY" 2>/dev/null || true
}

post_activity() {
  local ev_type="$1" ev_agent="$2" ev_message="$3" ev_status="${4:-info}" ev_detail="${5:-}"
  [ -z "$HOLM_API_KEY" ] && return 0
  local payload
  payload=$(python3 -c "
import json, sys
d = {'type': sys.argv[1], 'agent': sys.argv[2], 'message': sys.argv[3][:500], 'status': sys.argv[4], 'detail': sys.argv[5][:2000], 'iteration': 0}
print(json.dumps(d))
" "$ev_type" "$ev_agent" "$ev_message" "$ev_status" "$ev_detail")
  curl -s --max-time 5 -X POST "$HOLM_API_URL/api/activity" \
    -H "Content-Type: application/json" \
    -H "X-Api-Key: $HOLM_API_KEY" \
    -d "$payload" >/dev/null 2>&1 &
}

update_orchestrator_state() {
  local active_count=0
  for wid in "${!WORKER_PIDS[@]}"; do
    if kill -0 "${WORKER_PIDS[$wid]}" 2>/dev/null; then
      active_count=$((active_count + 1))
    fi
  done

  # Reset rate window if 5 hours elapsed
  local now=$(date +%s)
  local elapsed=$((now - WINDOW_START))
  if [ $elapsed -ge 18000 ]; then
    PROMPTS_USED=0
    WINDOW_START=$now
  fi

  local reset_at=$((WINDOW_START + 18000))
  api_put "/api/orchestrator" "{
    \"running\": true,
    \"activeWorkers\": $active_count,
    \"maxWorkers\": $MAX_WORKERS,
    \"totalIterations\": $TOTAL_ITERATIONS,
    \"totalCost\": $TOTAL_COST,
    \"rateLimitBudget\": {
      \"plan\": \"max200\",
      \"promptsUsed\": $PROMPTS_USED,
      \"promptsLimit\": $RATE_LIMIT_MAX,
      \"windowResetAt\": \"$(date -u -r $reset_at +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || date -u -d @$reset_at +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo '')\"
    },
    \"startedAt\": \"$(date -u -r $WINDOW_START +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || date -u +%Y-%m-%dT%H:%M:%SZ)\"
  }" >/dev/null &
}

# ─── Ollama helper ───
ask_ollama() {
  local system_prompt="$1"
  local user_prompt="$2"
  local retries=3
  local attempt=0

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

  while [ $attempt -lt $retries ]; do
    attempt=$((attempt + 1))
    local response
    response=$(curl -s --max-time 300 "$OLLAMA_URL/api/chat" \
      -H "Content-Type: application/json" \
      -d "$payload" 2>/dev/null) || true

    if [ -z "$response" ]; then
      sleep 5; continue
    fi

    local content
    content=$(echo "$response" | python3 -c "
import json, sys, re
try:
    data = json.load(sys.stdin)
    text = data.get('message', {}).get('content', '')
    text = re.sub(r'<think>.*?</think>', '', text, flags=re.DOTALL).strip()
    if text: print(text)
    else: print('ERROR: Empty')
except Exception as e: print(f'ERROR: {e}')
" 2>/dev/null)

    if [ -z "$content" ] || echo "$content" | grep -q "^ERROR:"; then
      sleep 5; continue
    fi
    echo "$content"
    return 0
  done
  echo "ERROR: Ollama failed after $retries attempts"
  return 1
}

# ─── Directive Decomposition ───
decompose_directive() {
  local directive_id="$1"
  local intent="$2"

  log "${MAGENTA}[ollama]${NC} Decomposing directive $directive_id..."
  post_activity "decompose_start" "ollama" "Decomposing: $intent" "info"

  # Get current tasks to avoid dupes
  local existing_tasks
  existing_tasks=$(api_get "/api/tasks")

  local decomp_prompt="You are a task decomposer for the HOLM system — a cyberpunk documentation nexus (Express.js + MongoDB on k3s).
Project dir: /Users/tim/holm-cyber-explorer

Given a high-level directive, break it into 2-5 concrete, parallelizable tasks.
Each task should be independently executable by a Claude Code worker.

Rules:
- Tasks that touch the same files should have dependencies on each other
- Each task title should be specific and actionable
- Include a brief description of what needs to be done
- Return ONLY valid JSON, no markdown, no explanation

Output format:
[
  {\"title\": \"...\", \"description\": \"...\", \"priority\": 1-5, \"tags\": [\"...\"], \"dependencies\": []},
  ...
]"

  local decomp_result
  decomp_result=$(ask_ollama "$decomp_prompt" "Directive: $intent

Existing tasks (avoid duplicates):
$(echo "$existing_tasks" | python3 -c "
import json, sys
try:
    tasks = json.load(sys.stdin)
    for t in tasks[-20:]:
        print(f\"- {t.get('taskId','?')}: {t.get('title','?')} [{t.get('status','?')}]\")
except: print('(none)')
")")

  # Extract JSON array from response
  local json_tasks
  json_tasks=$(echo "$decomp_result" | python3 -c "
import json, sys, re
text = sys.stdin.read()
# Find JSON array in text
match = re.search(r'\[[\s\S]*\]', text)
if match:
    try:
        tasks = json.loads(match.group())
        print(json.dumps(tasks))
    except:
        print('[]')
else:
    print('[]')
")

  local task_count
  task_count=$(echo "$json_tasks" | python3 -c "import json,sys; print(len(json.load(sys.stdin)))")

  if [ "$task_count" -eq 0 ]; then
    log "${YELLOW}[ollama]${NC} Decomposition produced 0 tasks"
    post_activity "decompose_fail" "ollama" "Decomposition failed for $directive_id" "warning"
    return 1
  fi

  log "${GREEN}[ollama]${NC} Created $task_count tasks for $directive_id"

  # POST to decompose endpoint
  api_post "/api/directives/$directive_id/decompose" "{\"tasks\": $json_tasks}" >/dev/null
  post_activity "decompose_done" "ollama" "Decomposed $directive_id into $task_count tasks" "success"
}

# ─── Worker Stream Parser ───
# Reads NDJSON from Claude's stdout and updates the dashboard
parse_worker_stream() {
  local worker_id="$1"
  local task_id="$2"
  local last_update=0
  local line_count=0

  while IFS= read -r line; do
    line_count=$((line_count + 1))

    # Throttle dashboard updates to 1/sec
    local now=$(date +%s)
    if [ $((now - last_update)) -ge 1 ]; then
      last_update=$now

      # Extract text from the JSON line
      local text
      text=$(echo "$line" | python3 -c "
import json, sys
try:
    d = json.load(sys.stdin)
    t = d.get('type', '')
    if t == 'assistant':
        content = d.get('message', {}).get('content', [])
        for c in content:
            if c.get('type') == 'text':
                text = c.get('text', '')[-200:]
                print(text)
                break
    elif t == 'tool_use':
        print('[tool] ' + d.get('name', '?'))
    elif t == 'result':
        cost = d.get('cost_usd', 0)
        if cost: print(f'COST:{cost}')
except: pass
" 2>/dev/null)

      if [ -n "$text" ]; then
        # Check for cost info
        if echo "$text" | grep -q "^COST:"; then
          local cost_val
          cost_val=$(echo "$text" | sed 's/COST://')
          # Update total cost (approximate)
          TOTAL_COST=$(python3 -c "print(round($TOTAL_COST + ${cost_val:-0}, 2))")
        fi

        # Update worker output
        api_put "/api/workers/$worker_id" "{\"currentOutput\": $(python3 -c "import json; print(json.dumps('$text'[:200])")}" >/dev/null 2>&1 &

        # Broadcast stream event for live output
        # (SSE broadcast happens server-side via worker update)
      fi
    fi

    # Heartbeat — update worker
    if [ $((line_count % 10)) -eq 0 ]; then
      api_put "/api/workers/$worker_id" "{\"status\": \"working\", \"lastHeartbeat\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"}" >/dev/null 2>&1 &
    fi
  done
}

# ─── Spawn Worker ───
spawn_worker() {
  local worker_id="$1"
  local task_id="$2"
  local task_title="$3"
  local task_desc="$4"
  local session_id

  session_id=$(uuidgen 2>/dev/null || python3 -c "import uuid; print(uuid.uuid4())")

  log_worker "$worker_id" "Spawning for $task_id: $task_title"
  post_activity "worker_spawn" "worker-$worker_id" "Starting: $task_title" "info"

  # Register worker
  api_put "/api/workers/$worker_id" "{
    \"sessionId\": \"$session_id\",
    \"status\": \"working\",
    \"currentTaskId\": \"$task_id\",
    \"currentDirectiveId\": null,
    \"role\": \"executor\",
    \"lastHeartbeat\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"startedAt\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"totalTasks\": 0,
    \"totalCost\": 0,
    \"currentOutput\": \"Starting...\"
  }" >/dev/null

  # Mark task as in_progress
  api_put "/api/tasks/$task_id" "{\"status\": \"in_progress\", \"assignedWorker\": \"$worker_id\", \"sessionId\": \"$session_id\", \"startedAt\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"}" >/dev/null

  # Load creds for Claude prompt
  source "$PROJECT_DIR/.autopilot/creds.env"

  local task_prompt="You are a worker for the HOLM multi-agent orchestrator. Execute this task completely:

TASK: $task_title
DETAILS: $task_desc
TASK ID: $task_id

Project: /Users/tim/holm-cyber-explorer
Stack: Express.js + MongoDB, deployed on k3s (Raspberry Pi cluster)
Site: https://holm.chat
Repo: github.com/timholm/holm-cyber-explorer

CREDENTIALS:
- SSH to k8s cluster: ssh rpi1@192.168.8.197 (password: $K8S_PASS)
- GitHub push: git remote set-url origin https://timholm:${GITHUB_TOKEN}@github.com/timholm/holm-cyber-explorer.git
- MongoDB: mongodb.holm-cyber.svc:27017
- Site: https://holm.chat
- holm.chat API key: $HOLM_API_KEY

Rules:
- Implement the task fully — edit files, test, commit, push
- Pull before making changes (git pull origin main)
- Commit with descriptive messages
- Verify changes work
- Keep the cyberpunk aesthetic
- When done, update the task via: curl -X PUT '$HOLM_API_URL/api/tasks/$task_id' -H 'Content-Type: application/json' -H 'X-Api-Key: $HOLM_API_KEY' -d '{\"status\":\"completed\",\"output\":\"<summary of what you did>\"}'
"

  # Spawn Claude in background
  local worker_log="$LOG_DIR/workers/$worker_id-$task_id-$(date '+%Y%m%d-%H%M%S').log"

  (
    cd "$PROJECT_DIR"
    claude -p \
      --dangerously-skip-permissions \
      --output-format stream-json \
      --max-budget-usd "$MAX_BUDGET_PER_TASK" \
      "$task_prompt" 2>&1 | tee "$worker_log" | parse_worker_stream "$worker_id" "$task_id"

    local exit_code=${PIPESTATUS[0]}
    PROMPTS_USED=$((PROMPTS_USED + 1))
    TOTAL_ITERATIONS=$((TOTAL_ITERATIONS + 1))

    if [ $exit_code -eq 0 ]; then
      log_worker "$worker_id" "${GREEN}Task $task_id completed${NC}"
      post_activity "worker_done" "worker-$worker_id" "Completed: $task_title" "success"
      api_put "/api/tasks/$task_id" "{\"status\": \"completed\", \"completedAt\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"}" >/dev/null
    else
      log_worker "$worker_id" "${RED}Task $task_id failed (exit=$exit_code)${NC}"
      post_activity "worker_error" "worker-$worker_id" "Failed: $task_title (exit=$exit_code)" "error"

      # Increment attempt counter
      local current_attempt
      current_attempt=$(api_get "/api/tasks" | python3 -c "
import json, sys
tasks = json.load(sys.stdin)
for t in tasks:
    if t.get('taskId') == '$task_id':
        print(t.get('attempt', 0))
        break
" 2>/dev/null || echo "0")
      local next_attempt=$((current_attempt + 1))

      if [ $next_attempt -ge 3 ]; then
        api_put "/api/tasks/$task_id" "{\"status\": \"failed\", \"attempt\": $next_attempt, \"failureReason\": \"Failed after $next_attempt attempts\"}" >/dev/null
      else
        api_put "/api/tasks/$task_id" "{\"status\": \"queued\", \"assignedWorker\": null, \"attempt\": $next_attempt}" >/dev/null
      fi
    fi

    # Mark worker idle
    api_put "/api/workers/$worker_id" "{\"status\": \"idle\", \"currentTaskId\": null, \"currentOutput\": \"\"}" >/dev/null
  ) &

  local pid=$!
  WORKER_PIDS[$worker_id]=$pid
  WORKER_TASKS[$worker_id]=$task_id
  WORKER_SESSIONS[$worker_id]=$session_id

  log_worker "$worker_id" "PID=$pid, session=$session_id"
}

# ─── Self-Direction via Ollama ───
generate_self_directive() {
  log "${MAGENTA}[ollama]${NC} Generating self-directive from CLAUDE.md vision..."
  post_activity "self_direct" "ollama" "Generating autonomous directive" "info"

  local claude_md=""
  if [ -f "$PROJECT_DIR/CLAUDE.md" ]; then
    claude_md=$(head -100 "$PROJECT_DIR/CLAUDE.md")
  fi

  local recent_tasks
  recent_tasks=$(api_get "/api/tasks" | python3 -c "
import json, sys
try:
    tasks = json.load(sys.stdin)
    for t in tasks[-15:]:
        print(f\"- {t.get('taskId','?')}: {t.get('title','?')} [{t.get('status','?')}]\")
except: print('(none)')
")

  local intent
  intent=$(ask_ollama "You are the strategic planner for HOLM — a personal sovereign intelligence system.
Generate ONE high-level directive for the next autonomous work cycle.
Focus on: improving the system, fixing issues, enhancing documentation, or building toward the air-gapped vision.
Return ONLY the directive intent as a single sentence. No JSON, no formatting." \
"System vision (from CLAUDE.md):
$claude_md

Recent task history:
$recent_tasks

Generate the next most impactful directive. One sentence only.")

  if [ -z "$intent" ] || echo "$intent" | grep -q "^ERROR:"; then
    log "${YELLOW}[ollama]${NC} Self-directive generation failed"
    return 1
  fi

  log "${GREEN}[ollama]${NC} Self-directive: $intent"

  # Create directive via API
  api_post "/api/directives" "{\"intent\": $(python3 -c "import json; print(json.dumps('$intent'[:500]))"), \"priority\": 3, \"source\": \"system\"}" >/dev/null
  post_activity "self_direct_created" "ollama" "Self-directive: $intent" "success"
}

# ─── Check Dependencies Met ───
deps_met() {
  local task_deps="$1"
  local all_tasks="$2"

  if [ -z "$task_deps" ] || [ "$task_deps" = "[]" ] || [ "$task_deps" = "null" ]; then
    return 0
  fi

  python3 -c "
import json, sys
deps = json.loads(sys.argv[1])
tasks = json.loads(sys.argv[2])
task_map = {t['taskId']: t['status'] for t in tasks}
for d in deps:
    if task_map.get(d) != 'completed':
        sys.exit(1)
sys.exit(0)
" "$task_deps" "$all_tasks" 2>/dev/null
}

# ─── Count Active Workers ───
count_active_workers() {
  local count=0
  for wid in "${!WORKER_PIDS[@]}"; do
    if kill -0 "${WORKER_PIDS[$wid]}" 2>/dev/null; then
      count=$((count + 1))
    fi
  done
  echo $count
}

# ─── Cleanup Dead Workers ───
cleanup_workers() {
  for wid in "${!WORKER_PIDS[@]}"; do
    if ! kill -0 "${WORKER_PIDS[$wid]}" 2>/dev/null; then
      unset "WORKER_PIDS[$wid]"
      unset "WORKER_TASKS[$wid]"
      unset "WORKER_SESSIONS[$wid]"
    fi
  done
}

# ─── Find Next Available Worker ID ───
next_worker_id() {
  for i in $(seq -w 1 99); do
    local wid="worker-$i"
    if [ -z "${WORKER_PIDS[$wid]+x}" ]; then
      echo "$wid"
      return
    fi
  done
  echo ""
}

# ─── Rate Limit Check ───
check_rate_limit() {
  local pct=$((PROMPTS_USED * 100 / RATE_LIMIT_MAX))
  if [ $pct -ge $RATE_LIMIT_STOP_PCT ]; then
    log "${RED}Rate limit reached ($PROMPTS_USED/$RATE_LIMIT_MAX). Stopping spawns.${NC}"
    return 2  # Hard stop
  elif [ $pct -ge $RATE_LIMIT_SLOW_PCT ]; then
    log "${YELLOW}Rate limit warning ($PROMPTS_USED/$RATE_LIMIT_MAX). Slowing down.${NC}"
    return 1  # Soft slow
  fi
  return 0
}

# ─── MAIN: Startup ───
log "${MAGENTA}╔══════════════════════════════════════════════╗${NC}"
log "${MAGENTA}║   HOLM Multi-Agent Orchestrator              ║${NC}"
log "${MAGENTA}║   Workers: up to $MAX_WORKERS Claude instances       ║${NC}"
log "${MAGENTA}║   Coordinator: Ollama $OLLAMA_MODEL              ║${NC}"
log "${MAGENTA}╚══════════════════════════════════════════════╝${NC}"
log "Project: $PROJECT_DIR"
log "Max workers: $MAX_WORKERS"
log "Budget: \$$MAX_BUDGET_PER_TASK per task"
log "Rate limit: $RATE_LIMIT_MAX prompts/5hr"
log "Stop: touch $STOP_FILE"
echo ""

post_activity "orchestrator_start" "system" "Orchestrator started (max $MAX_WORKERS workers)" "info"
update_orchestrator_state

# ─── Graceful Shutdown ───
cleanup() {
  log "${YELLOW}Shutting down orchestrator...${NC}"
  for wid in "${!WORKER_PIDS[@]}"; do
    if kill -0 "${WORKER_PIDS[$wid]}" 2>/dev/null; then
      log "Stopping $wid (PID ${WORKER_PIDS[$wid]})"
      kill "${WORKER_PIDS[$wid]}" 2>/dev/null || true
      api_put "/api/workers/$wid" '{"status": "stopped", "currentOutput": ""}' >/dev/null 2>&1
    fi
  done
  api_put "/api/orchestrator" '{"running": false, "activeWorkers": 0}' >/dev/null 2>&1
  post_activity "orchestrator_stop" "system" "Orchestrator stopped" "info"
  log "${GREEN}Orchestrator shutdown complete.${NC}"
  exit 0
}
trap cleanup SIGINT SIGTERM

# ─── MAIN: Coordinator Loop ───
while true; do
  # Check stop file
  if [ -f "$STOP_FILE" ]; then
    log "${YELLOW}Stop file detected.${NC}"
    rm -f "$STOP_FILE"
    cleanup
  fi

  cleanup_workers
  local_active=$(count_active_workers)

  # ── 1. Check for pending directives → decompose ──
  local pending_directives
  pending_directives=$(api_get "/api/directives?status=pending")
  local pending_count
  pending_count=$(echo "$pending_directives" | python3 -c "import json,sys; print(len(json.load(sys.stdin)))" 2>/dev/null || echo "0")

  if [ "$pending_count" -gt 0 ]; then
    IDLE_CYCLES=0
    # Decompose first pending directive
    local dir_id dir_intent
    dir_id=$(echo "$pending_directives" | python3 -c "import json,sys; d=json.load(sys.stdin); print(d[0]['directiveId'])" 2>/dev/null)
    dir_intent=$(echo "$pending_directives" | python3 -c "import json,sys; d=json.load(sys.stdin); print(d[0]['intent'])" 2>/dev/null)

    if [ -n "$dir_id" ] && [ -n "$dir_intent" ]; then
      decompose_directive "$dir_id" "$dir_intent"
    fi
  fi

  # ── 2. Find queued tasks with met dependencies ──
  local all_tasks
  all_tasks=$(api_get "/api/tasks")
  local queued_tasks
  queued_tasks=$(echo "$all_tasks" | python3 -c "
import json, sys
tasks = json.load(sys.stdin)
task_map = {t['taskId']: t.get('status') for t in tasks}
queued = []
for t in tasks:
    if t.get('status') != 'queued': continue
    if t.get('assignedWorker'): continue
    deps = t.get('dependencies', [])
    if all(task_map.get(d) == 'completed' for d in deps):
        queued.append(t)
print(json.dumps(queued))
" 2>/dev/null || echo '[]')

  local queued_count
  queued_count=$(echo "$queued_tasks" | python3 -c "import json,sys; print(len(json.load(sys.stdin)))" 2>/dev/null || echo "0")

  # ── 3. Assign tasks to idle workers / spawn new workers ──
  if [ "$queued_count" -gt 0 ]; then
    IDLE_CYCLES=0

    # Check rate limit before spawning
    set +e
    check_rate_limit
    local rate_status=$?
    set -e

    if [ $rate_status -lt 2 ]; then
      # Get task details for first queued task
      local task_info
      task_info=$(echo "$queued_tasks" | python3 -c "
import json, sys
tasks = json.load(sys.stdin)
if tasks:
    t = tasks[0]
    print(json.dumps({'id': t['taskId'], 'title': t.get('title',''), 'desc': t.get('description','')}))
else:
    print('{}')
" 2>/dev/null)

      local t_id t_title t_desc
      t_id=$(echo "$task_info" | python3 -c "import json,sys; print(json.load(sys.stdin).get('id',''))" 2>/dev/null)
      t_title=$(echo "$task_info" | python3 -c "import json,sys; print(json.load(sys.stdin).get('title',''))" 2>/dev/null)
      t_desc=$(echo "$task_info" | python3 -c "import json,sys; print(json.load(sys.stdin).get('desc',''))" 2>/dev/null)

      if [ -n "$t_id" ] && [ "$local_active" -lt "$MAX_WORKERS" ]; then
        local wid
        wid=$(next_worker_id)
        if [ -n "$wid" ]; then
          # Add delay if rate limit is soft warning
          if [ $rate_status -eq 1 ]; then
            log "${YELLOW}Rate limit soft warning — adding 30s delay${NC}"
            sleep 30
          fi
          spawn_worker "$wid" "$t_id" "$t_title" "$t_desc"
        fi
      fi
    fi
  else
    # No queued tasks
    IDLE_CYCLES=$((IDLE_CYCLES + 1))
  fi

  # ── 4. Monitor heartbeats — kill stale workers ──
  local now_epoch
  now_epoch=$(date +%s)
  for wid in "${!WORKER_PIDS[@]}"; do
    if kill -0 "${WORKER_PIDS[$wid]}" 2>/dev/null; then
      # Check if worker process has been running too long without output
      # (heartbeat timeout check via API would be more accurate,
      #  but simple PID check suffices for now)
      :
    fi
  done

  # ── 5. Check directive completion ──
  local active_directives
  active_directives=$(api_get "/api/directives?status=active")
  echo "$active_directives" | python3 -c "
import json, sys
directives = json.load(sys.stdin)
all_tasks_raw = '''$(echo "$all_tasks" | tr "'" "'")'''
try:
    all_tasks = json.loads(all_tasks_raw)
except:
    all_tasks = []
task_map = {t['taskId']: t.get('status') for t in all_tasks}
for d in directives:
    decomp = d.get('decomposition', [])
    if not decomp: continue
    statuses = [task_map.get(tid, 'unknown') for tid in decomp]
    if all(s in ('completed', 'failed', 'skipped') for s in statuses):
        print(d['directiveId'])
" 2>/dev/null | while read -r completed_dir; do
    if [ -n "$completed_dir" ]; then
      log "${GREEN}Directive $completed_dir completed!${NC}"
      api_put "/api/directives/$completed_dir" '{"status": "completed"}' >/dev/null
      post_activity "directive_complete" "system" "Directive $completed_dir completed" "success"
    fi
  done

  # ── 6. Self-direction — if idle too long, generate directive ──
  if [ "$IDLE_CYCLES" -ge "$IDLE_CYCLES_BEFORE_SELF_DIRECT" ] && [ "$local_active" -eq 0 ]; then
    log "${MAGENTA}Idle for $IDLE_CYCLES cycles — self-directing...${NC}"
    set +e
    generate_self_directive
    set -e
    IDLE_CYCLES=0
  fi

  # ── 7. Update orchestrator state ──
  update_orchestrator_state

  # ── Sleep until next cycle ──
  sleep "$COORDINATOR_INTERVAL"
done
