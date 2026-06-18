package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	CoordinatorAddr    string
	CoordinatorPort    int
	HeartbeatInterval  time.Duration
	HeartbeatTimeout   time.Duration
	MetricsCacheTTL    time.Duration
	CORSAllowedOrigins []string

	NodeID         string
	WorkerCount    int
	MaxCrawlDepth  int
	MaxRetries     int
	RequestTimeout time.Duration
	UserAgent      string
	MaxRedirects   int
	RobotsCacheTTL time.Duration

	DefaultRateLimit float64
	RateLimitBurst   int

	BloomExpectedItems uint
	BloomFPRate        float64

	VirtualNodes int

	LockTTL time.Duration

	QueueBatchSize int
	RetryBaseDelay time.Duration
	DLQMaxSize     int64

	DataDir string

	DashboardPort int

	LLMProvider string
	LLMAPIKey   string
	LLMModel    string

	EventChannelBuffer int
}

func Load() *Config {
	cfg := &Config{

		RedisAddr:     envOrDefault("CRAWLERX_REDIS_ADDR", "localhost:6379"),
		RedisPassword: envOrDefault("CRAWLERX_REDIS_PASSWORD", ""),
		RedisDB:       envOrDefaultInt("CRAWLERX_REDIS_DB", 0),

		CoordinatorAddr:    envOrDefault("CRAWLERX_COORDINATOR_ADDR", "http://localhost:8080"),
		CoordinatorPort:    envOrDefaultInt("CRAWLERX_COORDINATOR_PORT", 8080),
		HeartbeatInterval:  envOrDefaultDuration("CRAWLERX_HEARTBEAT_INTERVAL", 10*time.Second),
		HeartbeatTimeout:   envOrDefaultDuration("CRAWLERX_HEARTBEAT_TIMEOUT", 30*time.Second),
		MetricsCacheTTL:    envOrDefaultDuration("CRAWLERX_METRICS_CACHE_TTL", 1*time.Second),
		CORSAllowedOrigins: strings.Split(envOrDefault("CRAWLERX_CORS_ORIGINS", "http://localhost:3000"), ","),

		NodeID:         envOrDefault("CRAWLERX_NODE_ID", generateNodeID()),
		WorkerCount:    envOrDefaultInt("CRAWLERX_WORKER_COUNT", 10),
		MaxCrawlDepth:  envOrDefaultInt("CRAWLERX_MAX_CRAWL_DEPTH", 5),
		MaxRetries:     envOrDefaultInt("CRAWLERX_MAX_RETRIES", 3),
		RequestTimeout: envOrDefaultDuration("CRAWLERX_REQUEST_TIMEOUT", 30*time.Second),
		UserAgent:      envOrDefault("CRAWLERX_USER_AGENT", "CrawlerX/1.0 (+https://github.com/crawlerx/crawlerx)"),
		MaxRedirects:   envOrDefaultInt("CRAWLERX_MAX_REDIRECTS", 5),
		RobotsCacheTTL: envOrDefaultDuration("CRAWLERX_ROBOTS_CACHE_TTL", 24*time.Hour),

		DefaultRateLimit: envOrDefaultFloat("CRAWLERX_DEFAULT_RATE_LIMIT", 1.0),
		RateLimitBurst:   envOrDefaultInt("CRAWLERX_RATE_LIMIT_BURST", 3),

		BloomExpectedItems: uint(envOrDefaultInt("CRAWLERX_BLOOM_EXPECTED_ITEMS", 10_000_000)),
		BloomFPRate:        envOrDefaultFloat("CRAWLERX_BLOOM_FP_RATE", 0.01),

		VirtualNodes: envOrDefaultInt("CRAWLERX_VIRTUAL_NODES", 150),

		LockTTL: envOrDefaultDuration("CRAWLERX_LOCK_TTL", 60*time.Second),

		QueueBatchSize: envOrDefaultInt("CRAWLERX_QUEUE_BATCH_SIZE", 10),
		RetryBaseDelay: envOrDefaultDuration("CRAWLERX_RETRY_BASE_DELAY", 5*time.Second),
		DLQMaxSize:     int64(envOrDefaultInt("CRAWLERX_DLQ_MAX_SIZE", 100_000)),

		DataDir: envOrDefault("CRAWLERX_DATA_DIR", "./data"),

		DashboardPort: envOrDefaultInt("CRAWLERX_DASHBOARD_PORT", 3000),

		LLMProvider: envOrDefault("CRAWLERX_LLM_PROVIDER", ""),
		LLMAPIKey:   envOrDefault("CRAWLERX_LLM_API_KEY", ""),
		LLMModel:    envOrDefault("CRAWLERX_LLM_MODEL", ""),

		EventChannelBuffer: envOrDefaultInt("CRAWLERX_EVENT_CHANNEL_BUFFER", 1000),
	}

	return cfg
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envOrDefaultInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func envOrDefaultFloat(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return fallback
}

func envOrDefaultDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

func generateNodeID() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return hostname + "-" + strconv.Itoa(os.Getpid())
}
