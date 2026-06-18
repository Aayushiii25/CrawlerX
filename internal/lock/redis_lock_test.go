package lock

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestRedisLock_AcquireRelease(t *testing.T) {
	client := testRedisClient(t)
	ctx := context.Background()

	lock := NewRedisLock(client, "test-node-1")

	client.Del(ctx, lockKeyPrefix+"test-url-1")

	ok, err := lock.Acquire(ctx, "test-url-1", 10*time.Second)
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}
	if !ok {
		t.Fatal("Expected lock to be acquired")
	}

	err = lock.Release(ctx, "test-url-1")
	if err != nil {
		t.Fatalf("Release failed: %v", err)
	}
}

func TestRedisLock_MutualExclusion(t *testing.T) {
	client := testRedisClient(t)
	ctx := context.Background()

	lock1 := NewRedisLock(client, "node-1")
	lock2 := NewRedisLock(client, "node-2")

	client.Del(ctx, lockKeyPrefix+"test-url-2")

	ok, err := lock1.Acquire(ctx, "test-url-2", 10*time.Second)
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}
	if !ok {
		t.Fatal("Node 1 should acquire lock")
	}

	ok, err = lock2.Acquire(ctx, "test-url-2", 10*time.Second)
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}
	if ok {
		t.Fatal("Node 2 should NOT acquire lock held by Node 1")
	}

	err = lock2.Release(ctx, "test-url-2")
	if err != ErrLockNotOwned {
		t.Fatalf("Expected ErrLockNotOwned, got: %v", err)
	}

	err = lock1.Release(ctx, "test-url-2")
	if err != nil {
		t.Fatalf("Release failed: %v", err)
	}

	ok, err = lock2.Acquire(ctx, "test-url-2", 10*time.Second)
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}
	if !ok {
		t.Fatal("Node 2 should acquire lock after Node 1 released")
	}

	lock2.Release(ctx, "test-url-2")
}

func TestRedisLock_TTLExpiry(t *testing.T) {
	client := testRedisClient(t)
	ctx := context.Background()

	lock1 := NewRedisLock(client, "node-1")
	lock2 := NewRedisLock(client, "node-2")

	client.Del(ctx, lockKeyPrefix+"test-url-3")

	ok, _ := lock1.Acquire(ctx, "test-url-3", 1*time.Second)
	if !ok {
		t.Fatal("Node 1 should acquire lock")
	}

	time.Sleep(1500 * time.Millisecond)

	ok, err := lock2.Acquire(ctx, "test-url-3", 10*time.Second)
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}
	if !ok {
		t.Fatal("Node 2 should acquire lock after TTL expiry")
	}

	lock2.Release(ctx, "test-url-3")
}

func TestRedisLock_Extend(t *testing.T) {
	client := testRedisClient(t)
	ctx := context.Background()

	lock := NewRedisLock(client, "test-node")

	client.Del(ctx, lockKeyPrefix+"test-url-4")

	ok, _ := lock.Acquire(ctx, "test-url-4", 2*time.Second)
	if !ok {
		t.Fatal("Should acquire lock")
	}

	err := lock.Extend(ctx, "test-url-4", 10*time.Second)
	if err != nil {
		t.Fatalf("Extend failed: %v", err)
	}

	time.Sleep(3 * time.Second)

	ttl := client.TTL(ctx, lockKeyPrefix+"test-url-4").Val()
	if ttl <= 0 {
		t.Error("Lock should still be alive after extension")
	}

	lock.Release(ctx, "test-url-4")
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
