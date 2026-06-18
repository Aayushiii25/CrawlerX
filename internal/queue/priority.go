package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	priorityQueueKey = "crawlerx:queue:priority"
)

type CrawlTask struct {
	URL        string    `json:"url"`
	Domain     string    `json:"domain"`
	Depth      int       `json:"depth"`
	Priority   float64   `json:"priority"`
	ParentURL  string    `json:"parent_url,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	RetryCount int       `json:"retry_count"`
}

type Queue interface {
	Enqueue(ctx context.Context, task CrawlTask) error

	Dequeue(ctx context.Context, count int) ([]CrawlTask, error)

	Len(ctx context.Context) (int64, error)

	Peek(ctx context.Context, count int) ([]CrawlTask, error)
}

type PriorityQueue struct {
	client *redis.Client
	key    string
}

func NewPriorityQueue(client *redis.Client) *PriorityQueue {
	return &PriorityQueue{
		client: client,
		key:    priorityQueueKey,
	}
}

func (pq *PriorityQueue) Enqueue(ctx context.Context, task CrawlTask) error {
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now().UTC()
	}

	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("queue: marshal task: %w", err)
	}

	score := (1.0 - task.Priority) + float64(task.CreatedAt.UnixNano())*1e-19

	return pq.client.ZAdd(ctx, pq.key, redis.Z{
		Score:  score,
		Member: string(data),
	}).Err()
}

func (pq *PriorityQueue) Dequeue(ctx context.Context, count int) ([]CrawlTask, error) {
	if count <= 0 {
		count = 1
	}

	results, err := pq.client.ZPopMin(ctx, pq.key, int64(count)).Result()
	if err != nil {
		return nil, fmt.Errorf("queue: dequeue: %w", err)
	}

	tasks := make([]CrawlTask, 0, len(results))
	for _, z := range results {
		var task CrawlTask
		if err := json.Unmarshal([]byte(z.Member.(string)), &task); err != nil {
			continue
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (pq *PriorityQueue) Len(ctx context.Context) (int64, error) {
	return pq.client.ZCard(ctx, pq.key).Result()
}

func (pq *PriorityQueue) Peek(ctx context.Context, count int) ([]CrawlTask, error) {
	results, err := pq.client.ZRangeWithScores(ctx, pq.key, 0, int64(count-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("queue: peek: %w", err)
	}

	tasks := make([]CrawlTask, 0, len(results))
	for _, z := range results {
		var task CrawlTask
		if err := json.Unmarshal([]byte(z.Member.(string)), &task); err != nil {
			continue
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
