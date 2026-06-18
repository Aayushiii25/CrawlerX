package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type EventBus interface {
	Publish(ctx context.Context, event Event) error

	Subscribe(ctx context.Context) (<-chan Event, error)

	Close() error
}

type RedisEventBus struct {
	client      *redis.Client
	pubsub      *redis.PubSub
	channel     string
	mu          sync.Mutex
	closed      bool
	subscribers []chan Event
}

func NewRedisEventBus(client *redis.Client) *RedisEventBus {
	return &RedisEventBus{
		client:  client,
		channel: ChannelEvents,
	}
}

func (eb *RedisEventBus) Publish(ctx context.Context, event Event) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("eventbus: marshal event: %w", err)
	}

	return eb.client.Publish(ctx, eb.channel, string(data)).Err()
}

func (eb *RedisEventBus) Subscribe(ctx context.Context) (<-chan Event, error) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.closed {
		return nil, fmt.Errorf("eventbus: bus is closed")
	}

	if eb.pubsub == nil {
		eb.pubsub = eb.client.Subscribe(ctx, eb.channel)

		_, err := eb.pubsub.Receive(ctx)
		if err != nil {
			return nil, fmt.Errorf("eventbus: subscribe: %w", err)
		}

		go eb.dispatch(eb.pubsub.Channel())
	}

	ch := make(chan Event, 1000)
	eb.subscribers = append(eb.subscribers, ch)

	go func() {
		<-ctx.Done()
		eb.removeSubscriber(ch)
	}()

	return ch, nil
}

func (eb *RedisEventBus) dispatch(msgs <-chan *redis.Message) {
	for msg := range msgs {
		var event Event
		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			continue
		}

		eb.mu.Lock()
		for _, ch := range eb.subscribers {
			select {
			case ch <- event:
			default:

			}
		}
		eb.mu.Unlock()
	}
}

func (eb *RedisEventBus) removeSubscriber(ch chan Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	for i, sub := range eb.subscribers {
		if sub == ch {
			eb.subscribers = append(eb.subscribers[:i], eb.subscribers[i+1:]...)
			close(ch)
			return
		}
	}
}

func (eb *RedisEventBus) Close() error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.closed {
		return nil
	}
	eb.closed = true

	for _, ch := range eb.subscribers {
		close(ch)
	}
	eb.subscribers = nil

	if eb.pubsub != nil {
		return eb.pubsub.Close()
	}
	return nil
}
