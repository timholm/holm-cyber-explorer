#!/bin/bash

# TikTok Archive - Local Development Runner

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}TikTok Archive${NC}"
echo "=================="

# Check dependencies
check_dep() {
    if ! command -v $1 &> /dev/null; then
        echo -e "${RED}Error: $1 is not installed${NC}"
        return 1
    fi
}

echo -e "\n${YELLOW}Checking dependencies...${NC}"

MISSING=0
check_dep python3 || MISSING=1
check_dep pip3 || MISSING=1
check_dep node || MISSING=1
check_dep npm || MISSING=1
check_dep ffmpeg || MISSING=1

if [ $MISSING -eq 1 ]; then
    echo -e "\n${RED}Please install missing dependencies:${NC}"
    echo "  - Python 3.9+: https://python.org"
    echo "  - Node.js 18+: https://nodejs.org"
    echo "  - FFmpeg: brew install ffmpeg (macOS) or apt install ffmpeg (Linux)"
    exit 1
fi

echo -e "${GREEN}All dependencies found!${NC}"

# Create directories
mkdir -p archive/videos archive/thumbnails archive/audio

# Backend setup
echo -e "\n${YELLOW}Setting up backend...${NC}"
cd backend

if [ ! -d "venv" ]; then
    python3 -m venv venv
fi

source venv/bin/activate
pip install -q -r requirements.txt
echo -e "${GREEN}Backend ready!${NC}"

# Start backend in background
echo -e "\n${YELLOW}Starting backend on http://localhost:8000${NC}"
ARCHIVE_DIR="../archive" uvicorn main:app --host 0.0.0.0 --port 8000 &
BACKEND_PID=$!

cd ..

# Frontend setup
echo -e "\n${YELLOW}Setting up frontend...${NC}"
cd frontend

if [ ! -d "node_modules" ]; then
    npm install
fi

echo -e "${GREEN}Frontend ready!${NC}"

# Start frontend
echo -e "\n${YELLOW}Starting frontend on http://localhost:3000${NC}"
npm run dev &
FRONTEND_PID=$!

cd ..

# Trap to cleanup on exit
cleanup() {
    echo -e "\n${YELLOW}Shutting down...${NC}"
    kill $BACKEND_PID 2>/dev/null || true
    kill $FRONTEND_PID 2>/dev/null || true
    exit 0
}

trap cleanup SIGINT SIGTERM

echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}TikTok Archive is running!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "Frontend: ${YELLOW}http://localhost:3000${NC}"
echo -e "API:      ${YELLOW}http://localhost:8000${NC}"
echo -e "API Docs: ${YELLOW}http://localhost:8000/docs${NC}"
echo ""
echo -e "Press Ctrl+C to stop"

# Wait for processes
wait
