#!/bin/bash
# Run CrawlerX locally for development
# Prerequisites: Redis running on localhost:6379
# Usage: ./scripts/run-local.sh

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}CrawlerX Local Development Runner${NC}"
echo "=================================="

# Check Redis
if ! redis-cli ping > /dev/null 2>&1; then
  echo -e "${YELLOW}Redis not running. Starting Redis...${NC}"
  redis-server --daemonize yes
  sleep 1
fi
echo "✓ Redis is running"

# Build
echo "Building..."
go build -o bin/coordinator ./cmd/coordinator
go build -o bin/crawler ./cmd/crawler
echo "✓ Build complete"

# Create data directory
mkdir -p data

# Start coordinator
echo "Starting coordinator on :8080..."
CRAWLERX_CORS_ORIGINS="http://localhost:3000,http://localhost:3001,http://localhost:3002" CRAWLERX_LOG_LEVEL=info ./bin/coordinator &
COORD_PID=$!
sleep 2

# Start crawler nodes
echo "Starting crawler node 1..."
CRAWLERX_NODE_ID=local-node-1 CRAWLERX_WORKER_COUNT=5 CRAWLERX_LOG_LEVEL=info ./bin/crawler &
CRAWLER1_PID=$!
sleep 1

echo "Starting crawler node 2..."
CRAWLERX_NODE_ID=local-node-2 CRAWLERX_WORKER_COUNT=5 CRAWLERX_LOG_LEVEL=info ./bin/crawler &
CRAWLER2_PID=$!

echo ""
echo -e "${GREEN}CrawlerX is running!${NC}"
echo "  Coordinator: http://localhost:8080"
echo "  Dashboard:   cd dashboard && npm run dev  (run separately)"
echo ""
echo "  Seed URLs:   make seed"
echo ""
echo "Press Ctrl+C to stop all services"

# Cleanup on exit
cleanup() {
  echo ""
  echo "Shutting down..."
  kill $COORD_PID $CRAWLER1_PID $CRAWLER2_PID 2>/dev/null
  wait 2>/dev/null
  echo "All services stopped"
}

trap cleanup EXIT
wait
