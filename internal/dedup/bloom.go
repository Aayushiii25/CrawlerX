package dedup

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash"
	"hash/fnv"
	"math"

	"github.com/redis/go-redis/v9"
)

const (
	bloomKeyPrefix = "crawlerx:bloom:"
	bloomBitsKey   = bloomKeyPrefix + "bits"
	bloomCountKey  = bloomKeyPrefix + "count"
)

type BloomFilter interface {
	Add(ctx context.Context, url string) error

	MightContain(ctx context.Context, url string) (bool, error)

	Reset(ctx context.Context) error

	Stats(ctx context.Context) (BloomStats, error)
}

type BloomStats struct {
	NumBits           uint    `json:"num_bits"`
	NumHashFuncs      uint    `json:"num_hash_funcs"`
	ItemsAdded        int64   `json:"items_added"`
	ExpectedItems     uint    `json:"expected_items"`
	FalsePositiveRate float64 `json:"false_positive_rate"`
	EstimatedFPRate   float64 `json:"estimated_fp_rate"`
}

type RedisBloomFilter struct {
	client   *redis.Client
	numBits  uint
	numHash  uint
	expected uint
	fpRate   float64
}

func NewRedisBloomFilter(client *redis.Client, expectedItems uint, fpRate float64) *RedisBloomFilter {

	m := uint(math.Ceil(-float64(expectedItems) * math.Log(fpRate) / (math.Log(2) * math.Log(2))))

	k := uint(math.Ceil(float64(m) / float64(expectedItems) * math.Log(2)))

	if k < 1 {
		k = 1
	}

	return &RedisBloomFilter{
		client:   client,
		numBits:  m,
		numHash:  k,
		expected: expectedItems,
		fpRate:   fpRate,
	}
}

func (bf *RedisBloomFilter) Add(ctx context.Context, url string) error {
	positions := bf.hashPositions(url)

	pipe := bf.client.Pipeline()
	for _, pos := range positions {
		pipe.SetBit(ctx, bloomBitsKey, int64(pos), 1)
	}

	pipe.Incr(ctx, bloomCountKey)

	_, err := pipe.Exec(ctx)
	return err
}

func (bf *RedisBloomFilter) MightContain(ctx context.Context, url string) (bool, error) {
	positions := bf.hashPositions(url)

	pipe := bf.client.Pipeline()
	cmds := make([]*redis.IntCmd, len(positions))
	for i, pos := range positions {
		cmds[i] = pipe.GetBit(ctx, bloomBitsKey, int64(pos))
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	for _, cmd := range cmds {
		if cmd.Val() == 0 {
			return false, nil
		}
	}
	return true, nil
}

func (bf *RedisBloomFilter) Reset(ctx context.Context) error {
	pipe := bf.client.Pipeline()
	pipe.Del(ctx, bloomBitsKey)
	pipe.Del(ctx, bloomCountKey)
	_, err := pipe.Exec(ctx)
	return err
}

func (bf *RedisBloomFilter) Stats(ctx context.Context) (BloomStats, error) {
	count, err := bf.client.Get(ctx, bloomCountKey).Int64()
	if err == redis.Nil {
		count = 0
	} else if err != nil {
		return BloomStats{}, err
	}

	estimatedFP := math.Pow(1-math.Exp(-float64(bf.numHash)*float64(count)/float64(bf.numBits)), float64(bf.numHash))

	return BloomStats{
		NumBits:           bf.numBits,
		NumHashFuncs:      bf.numHash,
		ItemsAdded:        count,
		ExpectedItems:     bf.expected,
		FalsePositiveRate: bf.fpRate,
		EstimatedFPRate:   estimatedFP,
	}, nil
}

func (bf *RedisBloomFilter) hashPositions(url string) []uint {
	h1, h2 := doubleHash(url)
	positions := make([]uint, bf.numHash)
	for i := uint(0); i < bf.numHash; i++ {
		positions[i] = uint((h1 + uint64(i)*h2) % uint64(bf.numBits))
	}
	return positions
}

func doubleHash(data string) (uint64, uint64) {
	var h hash.Hash64 = fnv.New64a()
	h.Write([]byte(data))
	sum := h.Sum64()

	h1 := uint64(sum >> 32)
	h2 := uint64(sum & 0xFFFFFFFF)

	if h2%2 == 0 {
		h2++
	}

	lengthHash := fnv.New64a()
	lengthHash.Write([]byte(fmt.Sprintf("%s:%d", data, len(data))))
	mixed := lengthHash.Sum64()

	h1 = h1 ^ (mixed >> 32)

	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, mixed)
	h2Mixer := fnv.New64a()
	h2Mixer.Write(buf)
	h2 = h2 ^ (h2Mixer.Sum64() & 0xFFFFFFFF)
	if h2%2 == 0 {
		h2++
	}

	return h1, h2
}
