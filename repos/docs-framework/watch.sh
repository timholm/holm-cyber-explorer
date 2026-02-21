#!/usr/bin/env bash
#
# watch.sh -- Watch for new/changed .md files, build HTML, commit and push.
# Polls every 30 seconds. Run in background: nohup ./watch.sh &
#
set -uo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$DIR"

LOG="$DIR/.watch.log"
SNAP_FILE="$DIR/.md-snapshot"

log() { printf "[%s] %s\n" "$(date +%H:%M:%S)" "$*" | tee -a "$LOG"; }

# Snapshot: list of all .md files with their sizes and mod times
snapshot() {
    find "$DIR" -maxdepth 1 -name '*.md' -exec stat -f '%N %z %m' {} + 2>/dev/null | sort
}

# Initialize
touch "$SNAP_FILE"
log "Watcher started. Polling every 30s for .md changes in $DIR"

while true; do
    CURRENT="$(snapshot)"
    PREVIOUS="$(cat "$SNAP_FILE" 2>/dev/null || echo "")"

    if [[ "$CURRENT" != "$PREVIOUS" ]]; then
        log "Changes detected in .md files"

        # Wait for writes to finish
        sleep 8

        # Re-snapshot after wait (in case still writing)
        CURRENT="$(snapshot)"

        # Rebuild site
        log "Running build.py..."
        if python3 "$DIR/build.py" >> "$LOG" 2>&1; then
            log "Build succeeded"
        else
            log "Build failed, skipping commit"
            echo "$CURRENT" > "$SNAP_FILE"
            sleep 30
            continue
        fi

        # Pull first to avoid conflicts from GitHub Actions bot commits
        git pull --rebase >> "$LOG" 2>&1 || true

        # Stage everything
        git add *.md site/*.html site/*.css build.py 2>/dev/null

        if git diff --cached --quiet; then
            log "No changes to commit"
        else
            CHANGED=$(git diff --cached --name-only | wc -l | tr -d ' ')
            MSG="Auto-commit: build and deploy ${CHANGED} updated files $(date +%Y-%m-%d\ %H:%M)"
            git commit -m "$MSG" >> "$LOG" 2>&1
            if git push origin main >> "$LOG" 2>&1; then
                log "Pushed $CHANGED files to GitHub"
            else
                # Retry once after pull
                git pull --rebase >> "$LOG" 2>&1 || true
                git push origin main >> "$LOG" 2>&1 && log "Pushed (retry)" || log "Push failed"
            fi
        fi

        echo "$CURRENT" > "$SNAP_FILE"
    fi

    sleep 30
done
