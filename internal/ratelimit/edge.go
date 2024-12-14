package ratelimit

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type EdgeRateLimiter struct {
	redis      *redis.Client
	window     time.Duration
	maxRequest int64
}

func NewEdgeRateLimiter(redisURL string, window time.Duration, maxRequest int64) (*EdgeRateLimiter, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	return &EdgeRateLimiter{
		redis:      client,
		window:     window,
		maxRequest: maxRequest,
	}, nil
}

func (rl *EdgeRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	pipe := rl.redis.Pipeline()
	now := time.Now().UnixNano()
	windowStart := now - rl.window.Nanoseconds()

	// Remove old requests
	pipe.ZRemRangeByScore(ctx, key, "0", string(windowStart))

	// Add current request
	pipe.ZAdd(ctx, key, &redis.Z{
		Score:  float64(now),
		Member: now,
	})

	// Get count of requests in window
	pipe.ZCard(ctx, key)

	// Execute pipeline
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	// Get the count from the last command
	count := cmds[2].(*redis.IntCmd).Val()

	return count <= rl.maxRequest, nil
}
