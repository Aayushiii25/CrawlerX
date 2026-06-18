package dedup

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestBloomFilter_BasicOperations(t *testing.T) {
	client := testRedisClient(t)
	ctx := context.Background()

	bf := NewRedisBloomFilter(client, 1000, 0.01)
	defer bf.Reset(ctx)

	if err := bf.Reset(ctx); err != nil {
		t.Fatalf("Reset failed: %v", err)
	}

	found, err := bf.MightContain(ctx, "https://example.com/not-added")
	if err != nil {
		t.Fatalf("MightContain failed: %v", err)
	}
	if found {
		t.Error("Expected non-added URL to not be found")
	}

	if err := bf.Add(ctx, "https://example.com/page1"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	found, err = bf.MightContain(ctx, "https://example.com/page1")
	if err != nil {
		t.Fatalf("MightContain failed: %v", err)
	}
	if !found {
		t.Error("Expected added URL to be found (no false negatives)")
	}

	found, err = bf.MightContain(ctx, "https://example.com/page2")
	if err != nil {
		t.Fatalf("MightContain failed: %v", err)
	}
	if found {
		t.Log("Note: false positive detected — acceptable but unusual for small set")
	}
}

func TestBloomFilter_FalsePositiveRate(t *testing.T) {
	client := testRedisClient(t)
	ctx := context.Background()

	n := uint(10000)
	targetFP := 0.01
	bf := NewRedisBloomFilter(client, n, targetFP)
	defer bf.Reset(ctx)

	if err := bf.Reset(ctx); err != nil {
		t.Fatalf("Reset failed: %v", err)
	}

	for i := uint(0); i < n; i++ {
		url := fmt.Sprintf("https://example.com/added/%d", i)
		if err := bf.Add(ctx, url); err != nil {
			t.Fatalf("Add failed at %d: %v", i, err)
		}
	}

	falsePositives := 0
	testCount := int(n) * 10
	for i := 0; i < testCount; i++ {
		url := fmt.Sprintf("https://other-domain.com/not-added/%d", i)
		found, err := bf.MightContain(ctx, url)
		if err != nil {
			t.Fatalf("MightContain failed at %d: %v", i, err)
		}
		if found {
			falsePositives++
		}
	}

	actualFP := float64(falsePositives) / float64(testCount)
	t.Logf("False positive rate: %.4f%% (%d/%d)", actualFP*100, falsePositives, testCount)

	if actualFP > targetFP*2 {
		t.Errorf("False positive rate %.4f exceeds 2x target %.4f", actualFP, targetFP)
	}
}

func TestBloomFilter_Stats(t *testing.T) {
	client := testRedisClient(t)
	ctx := context.Background()

	bf := NewRedisBloomFilter(client, 1000, 0.01)
	defer bf.Reset(ctx)

	if err := bf.Reset(ctx); err != nil {
		t.Fatalf("Reset failed: %v", err)
	}

	for i := 0; i < 100; i++ {
		bf.Add(ctx, fmt.Sprintf("https://example.com/%d", i))
	}

	stats, err := bf.Stats(ctx)
	if err != nil {
		t.Fatalf("Stats failed: %v", err)
	}

	if stats.ItemsAdded != 100 {
		t.Errorf("Expected 100 items added, got %d", stats.ItemsAdded)
	}
	if stats.NumBits == 0 {
		t.Error("Expected non-zero bit count")
	}
	if stats.NumHashFuncs == 0 {
		t.Error("Expected non-zero hash function count")
	}
	t.Logf("Stats: %+v", stats)
}

func testRedisClient(t *testing.T) *redis.Client {
	t.Helper()

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("Skipping test: Redis not available at localhost:6379: %v", err)
	}

	t.Cleanup(func() {
		client.Close()
	})

	return client
}
