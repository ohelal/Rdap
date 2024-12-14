package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"sync/atomic"
	"time"
)

// DistributedCache represents a distributed cache using Redis
type DistributedCache struct {
	client  *redis.Client
	ttl     time.Duration
	pool    *redis.Client
	poolMu  sync.RWMutex
	metrics struct {
		hits    uint64
		misses  uint64
		errors  uint64
		latency time.Duration
	}
}

// CacheMetrics tracks cache performance
type CacheMetrics struct {
	hits    uint64
	misses  uint64
	errors  uint64
	latency time.Duration
}

// NewDistributedCache creates a new distributed cache instance
func NewDistributedCache(config *CacheConfig) (*DistributedCache, error) {
	// Create connection pool
	poolOptions := &redis.Options{
		Addr:               config.RedisURL,
		DB:                 0,
		MaxRetries:         3,
		MinRetryBackoff:    8 * time.Millisecond,
		MaxRetryBackoff:    512 * time.Millisecond,
		DialTimeout:        5 * time.Second,
		ReadTimeout:        3 * time.Second,
		WriteTimeout:       3 * time.Second,
		PoolSize:           50,
		MinIdleConns:       10,
		MaxConnAge:         30 * time.Minute,
		PoolTimeout:        4 * time.Second,
		IdleTimeout:        5 * time.Minute,
		IdleCheckFrequency: 1 * time.Minute,
	}

	client := redis.NewClient(poolOptions)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &DistributedCache{
		client: client,
		ttl:    config.RedisTTL,
	}, nil
}

// Get retrieves a value from the cache with metrics
func (c *DistributedCache) Get(key string) (interface{}, bool) {
	start := time.Now()
	ctx := context.Background()

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			atomic.AddUint64(&c.metrics.misses, 1)
			return nil, false
		}
		atomic.AddUint64(&c.metrics.errors, 1)
		return nil, false
	}

	atomic.AddUint64(&c.metrics.hits, 1)
	c.metrics.latency = time.Since(start)

	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		atomic.AddUint64(&c.metrics.errors, 1)
		return nil, false
	}

	return result, true
}

// Set stores a value in the cache with retry logic
func (c *DistributedCache) Set(key string, value interface{}) error {
	ctx := context.Background()
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	// Retry logic for Set operations
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err = c.client.Set(ctx, key, data, c.ttl).Err()
		if err == nil {
			return nil
		}

		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
		}
	}

	return fmt.Errorf("failed to set cache after %d retries: %v", maxRetries, err)
}

// Delete removes a value from the cache
func (c *DistributedCache) Delete(key string) error {
	ctx := context.Background()
	return c.client.Del(ctx, key).Err()
}

// Close closes the Redis connection
func (c *DistributedCache) Close() error {
	return c.client.Close()
}

// GetMetrics returns cache performance metrics
func (c *DistributedCache) GetMetrics() *CacheMetrics {
	return &CacheMetrics{
		hits:    atomic.LoadUint64(&c.metrics.hits),
		misses:  atomic.LoadUint64(&c.metrics.misses),
		errors:  atomic.LoadUint64(&c.metrics.errors),
		latency: c.metrics.latency,
	}
}
