package eventbus

import "time"

type EventType string

const (
	EventNodeJoined        EventType = "node.joined"
	EventNodeFailed        EventType = "node.failed"
	EventNodeLeft          EventType = "node.left"
	EventPageCrawled       EventType = "page.crawled"
	EventDuplicateDetected EventType = "duplicate.detected"
	EventQueueOverflow     EventType = "queue.overflow"
	EventCrawlError        EventType = "crawl.error"
	EventTaskAssigned      EventType = "task.assigned"
	EventTaskCompleted     EventType = "task.completed"
)

type Event struct {
	ID        string    `json:"id"`
	Type      EventType `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	NodeID    string    `json:"node_id"`
	Payload   Payload   `json:"payload"`
}

type Payload struct {
	URL        string `json:"url,omitempty"`
	Domain     string `json:"domain,omitempty"`
	StatusCode int    `json:"status_code,omitempty"`
	Depth      int    `json:"depth,omitempty"`

	WorkerCount int    `json:"worker_count,omitempty"`
	NodeAddr    string `json:"node_addr,omitempty"`

	QueueSize int64  `json:"queue_size,omitempty"`
	QueueName string `json:"queue_name,omitempty"`

	Error      string `json:"error,omitempty"`
	RetryCount int    `json:"retry_count,omitempty"`

	Duration   int64 `json:"duration_ms,omitempty"`
	ContentLen int64 `json:"content_length,omitempty"`
	LinksFound int   `json:"links_found,omitempty"`

	Message string `json:"message,omitempty"`
}

const (
	ChannelEvents = "crawlerx:events"
)
