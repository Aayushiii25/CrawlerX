package eventbus

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestRedisEventBus_PublishSubscribe(t *testing.T) {
	client := testRedisClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bus := NewRedisEventBus(client)
	bus.channel = "crawlerx:events:test"
	defer bus.Close()

	events, err := bus.Subscribe(ctx)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	testEvent := Event{
		Type:   EventPageCrawled,
		NodeID: "test-node",
		Payload: Payload{
			URL:        "https://example.com",
			StatusCode: 200,
			Duration:   150,
		},
	}

	if err := bus.Publish(ctx, testEvent); err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	select {
	case event := <-events:
		if event.Type != EventPageCrawled {
			t.Errorf("Expected type %s, got %s", EventPageCrawled, event.Type)
		}
		if event.Payload.URL != "https://example.com" {
			t.Errorf("Expected URL https://example.com, got %s", event.Payload.URL)
		}
		if event.ID == "" {
			t.Error("Expected auto-generated event ID")
		}
		t.Logf("Received event: %+v", event)
	case <-time.After(5 * time.Second):
		t.Fatal("Timed out waiting for event")
	}
}

func TestRedisEventBus_MultipleSubscribers(t *testing.T) {
	client := testRedisClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bus := NewRedisEventBus(client)
	bus.channel = "crawlerx:events:test"
	defer bus.Close()

	events1, err := bus.Subscribe(ctx)
	if err != nil {
		t.Fatalf("Subscribe 1 failed: %v", err)
	}

	events2, err := bus.Subscribe(ctx)
	if err != nil {
		t.Fatalf("Subscribe 2 failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	bus.Publish(ctx, Event{
		Type:   EventNodeJoined,
		NodeID: "node-1",
	})

	for i, ch := range []<-chan Event{events1, events2} {
		select {
		case event := <-ch:
			if event.Type != EventNodeJoined {
				t.Errorf("Subscriber %d: wrong event type: %s", i, event.Type)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("Subscriber %d: timed out", i)
		}
	}
}

func testRedisClient(t *testing.T) *redis.Client {
	t.Helper()

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("Skipping test: Redis not available: %v", err)
	}

	t.Cleanup(func() {
		client.Close()
	})

	return client
}
