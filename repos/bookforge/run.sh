#!/bin/bash
#
# BookForge Development Runner
# Starts all services manually for development/testing
#

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VENV_DIR="${SCRIPT_DIR}/venv"
LOG_DIR="${SCRIPT_DIR}/logs"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# PID files for cleanup
WORKER_PID_FILE="${SCRIPT_DIR}/.worker.pid"
WEB_PID_FILE="${SCRIPT_DIR}/.web.pid"

echo_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

echo_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

echo_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

cleanup() {
    echo ""
    echo_info "Shutting down BookForge..."
    
    # Kill worker process
    if [[ -f "${WORKER_PID_FILE}" ]]; then
        WORKER_PID=$(cat "${WORKER_PID_FILE}")
        if kill -0 "${WORKER_PID}" 2>/dev/null; then
            echo_info "Stopping worker (PID: ${WORKER_PID})..."
            kill "${WORKER_PID}" 2>/dev/null || true
        fi
        rm -f "${WORKER_PID_FILE}"
    fi
    
    # Kill web process
    if [[ -f "${WEB_PID_FILE}" ]]; then
        WEB_PID=$(cat "${WEB_PID_FILE}")
        if kill -0 "${WEB_PID}" 2>/dev/null; then
            echo_info "Stopping web server (PID: ${WEB_PID})..."
            kill "${WEB_PID}" 2>/dev/null || true
        fi
        rm -f "${WEB_PID_FILE}"
    fi
    
    echo_info "BookForge stopped."
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM

# Check if virtual environment exists
if [[ ! -d "${VENV_DIR}" ]]; then
    echo_error "Virtual environment not found at ${VENV_DIR}"
    echo_info "Create it with: python3 -m venv ${VENV_DIR}"
    exit 1
fi

# Activate virtual environment
source "${VENV_DIR}/bin/activate"

# Create local log directory for development
mkdir -p "${LOG_DIR}"

# Change to project directory
cd "${SCRIPT_DIR}"

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}    BookForge Development Server${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Start background worker
echo_info "Starting background worker..."
python worker.py > "${LOG_DIR}/worker.log" 2>&1 &
WORKER_PID=$!
echo "${WORKER_PID}" > "${WORKER_PID_FILE}"
echo_info "Worker started (PID: ${WORKER_PID})"

# Wait a moment for worker to initialize
sleep 1

# Check if worker is still running
if ! kill -0 "${WORKER_PID}" 2>/dev/null; then
    echo_error "Worker failed to start. Check ${LOG_DIR}/worker.log for details."
    rm -f "${WORKER_PID_FILE}"
    exit 1
fi

# Start Flask development server
echo_info "Starting Flask development server..."
echo_info "Web interface will be available at: http://localhost:5000"
echo ""
echo_warn "Press Ctrl+C to stop all services"
echo ""

# Export Flask environment variables
export FLASK_APP=app.py
export FLASK_ENV=development
export FLASK_DEBUG=1

# Run Flask in foreground (allows Ctrl+C to work)
python -m flask run --host=0.0.0.0 --port=5000 &
WEB_PID=$!
echo "${WEB_PID}" > "${WEB_PID_FILE}"

# Wait for either process to exit
wait ${WEB_PID} 2>/dev/null || true

# If we get here, web server stopped - cleanup
cleanup
