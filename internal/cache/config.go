package cache

import (
	"github.com/go-redis/redis/v8"
	"time"
)

// CacheConfig holds cache configuration
type CacheConfig struct {
	LocalTTL      time.Duration
	RedisTTL      time.Duration
	MaxLocalSize  int64
	EnableRedis   bool
	RedisURL      string
	RedisConfig   *redis.Options
}

// ErrCacheMiss is returned when a key is not found in the cache
var ErrCacheMiss = redis.Nil
