.PHONY: build test run-coordinator run-crawler run-dashboard docker-up docker-down seed clean

# --- Build ---
build:
	@echo "Building coordinator..."
	go build -o bin/coordinator ./cmd/coordinator
	@echo "Building crawler..."
	go build -o bin/crawler ./cmd/crawler
	@echo "Build complete"

# --- Test ---
test:
	go test ./internal/hashring/... -v
	go test ./internal/scoring/... -v

test-integration:
	go test ./internal/dedup/... -v
	go test ./internal/lock/... -v
	go test ./internal/eventbus/... -v

test-all: test test-integration

# --- Run locally ---
run-coordinator:
	CRAWLERX_LOG_LEVEL=debug go run ./cmd/coordinator

run-crawler:
	CRAWLERX_LOG_LEVEL=debug CRAWLERX_NODE_ID=local-crawler-1 go run ./cmd/crawler

run-crawler-2:
	CRAWLERX_LOG_LEVEL=debug CRAWLERX_NODE_ID=local-crawler-2 go run ./cmd/crawler

run-dashboard:
	cd dashboard && npm run dev

# --- Docker ---
docker-up:
	cd deployments && docker compose up --build -d

docker-down:
	cd deployments && docker compose down

docker-logs:
	cd deployments && docker compose logs -f

# --- Seed URLs ---
seed:
	@echo "Seeding URLs..."
	curl -s -X POST http://localhost:8080/api/crawl \
		-H "Content-Type: application/json" \
		-d '{"urls":["https://example.com","https://httpbin.org","https://go.dev","https://news.ycombinator.com","https://github.com/trending"],"depth":0}' | jq .
	@echo "Done"

# --- Clean ---
clean:
	rm -rf bin/
	rm -rf data/
	@echo "Cleaned"

# --- Lint ---
lint:
	go vet ./...
	@echo "Lint passed"
