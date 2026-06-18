package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	retryQueueKey = "crawlerx:queue:retry"
)

type RetryQueue struct {
	client     *redis.Client
	key        string
	maxRetries int
	baseDelay  time.Duration
	dlq        *DeadLetterQueue
}

func NewRetryQueue(client *redis.Client, maxRetries int, baseDelay time.Duration, dlq *DeadLetterQueue) *RetryQueue {
	return &RetryQueue{
		client:     client,
		key:        retryQueueKey,
		maxRetries: maxRetries,
		baseDelay:  baseDelay,
		dlq:        dlq,
	}
}

func (rq *RetryQueue) Retry(ctx context.Context, task CrawlTask, errMsg string) (bool, error) {
	task.RetryCount++

	if task.RetryCount > rq.maxRetries {

		return false, rq.dlq.Add(ctx, task, errMsg)
	}

	delay := time.Duration(float64(rq.baseDelay) * math.Pow(2, float64(task.RetryCount-1)))

	retryAt := time.Now().Add(delay).Unix()

	data, err := json.Marshal(task)
	if err != nil {
		return false, fmt.Errorf("retry: marshal task: %w", err)
	}

	err = rq.client.ZAdd(ctx, rq.key, redis.Z{
		Score:  float64(retryAt),
		Member: string(data),
	}).Err()
	if err != nil {
		return false, fmt.Errorf("retry: enqueue: %w", err)
	}

	return true, nil
}

func (rq *RetryQueue) GetReady(ctx context.Context, count int) ([]CrawlTask, error) {
	now := float64(time.Now().Unix())

	script := redis.NewScript(`
		local results = redis.call('ZRANGEBYSCORE', KEYS[1], '-inf', ARGV[1], 'LIMIT', 0, ARGV[2])
		if #results > 0 then
			redis.call('ZREM', KEYS[1], unpack(results))
		end
		return results
	`)

	results, err := script.Run(ctx, rq.client, []string{rq.key}, now, count).StringSlice()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("retry: get ready: %w", err)
	}

	tasks := make([]CrawlTask, 0, len(results))
	for _, r := range results {
		var task CrawlTask
		if err := json.Unmarshal([]byte(r), &task); err != nil {
			continue
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (rq *RetryQueue) Len(ctx context.Context) (int64, error) {
	return rq.client.ZCard(ctx, rq.key).Result()
}
