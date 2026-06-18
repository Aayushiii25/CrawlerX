#!/bin/bash
# Seed URLs for testing CrawlerX
# Usage: ./scripts/seed.sh

set -e

COORDINATOR=${CRAWLERX_COORDINATOR_ADDR:-http://localhost:8080}

echo "Seeding URLs to $COORDINATOR..."

curl -s -X POST "$COORDINATOR/api/crawl" \
  -H "Content-Type: application/json" \
  -d '{
    "urls": [
      "https://example.com",
      "https://httpbin.org",
      "https://go.dev",
      "https://go.dev/doc",
      "https://go.dev/blog",
      "https://news.ycombinator.com",
      "https://github.com/trending",
      "https://en.wikipedia.org/wiki/Web_crawler",
      "https://developer.mozilla.org/en-US/docs/Web",
      "https://www.rust-lang.org"
    ],
    "depth": 0
  }' | python3 -m json.tool 2>/dev/null || echo "(install python3 for pretty output)"

echo ""
echo "Seeding complete. Check dashboard at http://localhost:3000"
