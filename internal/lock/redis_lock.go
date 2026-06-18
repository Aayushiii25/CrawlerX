package lock

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const lockKeyPrefix = "crawlerx:lock:"

var (
	ErrLockNotAcquired = errors.New("lock: not acquired")

	ErrLockNotOwned = errors.New("lock: not owned by this node")
)

type DistributedLock interface {
	Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error)

	Release(ctx context.Context, key string) error

	Extend(ctx context.Context, key string, ttl time.Duration) error
}

type RedisLock struct {
	client *redis.Client
	nodeID string
}

func NewRedisLock(client *redis.Client, nodeID string) *RedisLock {
	return &RedisLock{
		client: client,
		nodeID: nodeID,
	}
}

func (rl *RedisLock) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	lockKey := lockKeyPrefix + key
	ok, err := rl.client.SetNX(ctx, lockKey, rl.nodeID, ttl).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

var releaseScript = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	end
	return 0
`)

func (rl *RedisLock) Release(ctx context.Context, key string) error {
	lockKey := lockKeyPrefix + key
	result, err := releaseScript.Run(ctx, rl.client, []string{lockKey}, rl.nodeID).Int64()
	if err != nil {
		return err
	}
	if result == 0 {
		return ErrLockNotOwned
	}
	return nil
}

var extendScript = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("PEXPIRE", KEYS[1], ARGV[2])
	end
	return 0
`)

func (rl *RedisLock) Extend(ctx context.Context, key string, ttl time.Duration) error {
	lockKey := lockKeyPrefix + key
	result, err := extendScript.Run(ctx, rl.client, []string{lockKey}, rl.nodeID, ttl.Milliseconds()).Int64()
	if err != nil {
		return err
	}
	if result == 0 {
		return ErrLockNotOwned
	}
	return nil
}
