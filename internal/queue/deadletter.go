package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	dlqKey = "crawlerx:queue:dlq"
)

type DeadLetterEntry struct {
	Task     CrawlTask `json:"task"`
	Error    string    `json:"error"`
	FailedAt time.Time `json:"failed_at"`
	Retries  int       `json:"retries"`
}

type DeadLetterQueue struct {
	client  *redis.Client
	key     string
	maxSize int64
}

func NewDeadLetterQueue(client *redis.Client, maxSize int64) *DeadLetterQueue {
	return &DeadLetterQueue{
		client:  client,
		key:     dlqKey,
		maxSize: maxSize,
	}
}

func (dlq *DeadLetterQueue) Add(ctx context.Context, task CrawlTask, errMsg string) error {
	entry := DeadLetterEntry{
		Task:     task,
		Error:    errMsg,
		FailedAt: time.Now().UTC(),
		Retries:  task.RetryCount,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("dlq: marshal entry: %w", err)
	}

	pipe := dlq.client.Pipeline()
	pipe.RPush(ctx, dlq.key, string(data))

	pipe.LTrim(ctx, dlq.key, -dlq.maxSize, -1)
	_, err = pipe.Exec(ctx)

	return err
}

func (dlq *DeadLetterQueue) List(ctx context.Context, start, stop int64) ([]DeadLetterEntry, error) {
	results, err := dlq.client.LRange(ctx, dlq.key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("dlq: list: %w", err)
	}

	entries := make([]DeadLetterEntry, 0, len(results))
	for _, r := range results {
		var entry DeadLetterEntry
		if err := json.Unmarshal([]byte(r), &entry); err != nil {
			continue
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (dlq *DeadLetterQueue) Len(ctx context.Context) (int64, error) {
	return dlq.client.LLen(ctx, dlq.key).Result()
}

func (dlq *DeadLetterQueue) Flush(ctx context.Context) error {
	return dlq.client.Del(ctx, dlq.key).Err()
}

func (dlq *DeadLetterQueue) Requeue(ctx context.Context, index int64, pq *PriorityQueue) error {

	results, err := dlq.client.LRange(ctx, dlq.key, index, index).Result()
	if err != nil {
		return fmt.Errorf("dlq: requeue: %w", err)
	}
	if len(results) == 0 {
		return fmt.Errorf("dlq: entry not found at index %d", index)
	}

	var entry DeadLetterEntry
	if err := json.Unmarshal([]byte(results[0]), &entry); err != nil {
		return fmt.Errorf("dlq: requeue unmarshal: %w", err)
	}

	entry.Task.RetryCount = 0
	if err := pq.Enqueue(ctx, entry.Task); err != nil {
		return fmt.Errorf("dlq: requeue enqueue: %w", err)
	}

	return nil
}
