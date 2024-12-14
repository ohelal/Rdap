package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

// RateLimiterConfig holds the configuration for rate limiting
type RateLimiterConfig struct {
	RedisClient  *redis.Client
	MaxRequests  map[string]int // Requests per window per endpoint
	WindowSize   time.Duration
	DefaultMax   int // Default max requests for unspecified endpoints
}

// RateLimiter creates a new rate limiting middleware with Redis backend
func RateLimiter(config RateLimiterConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		endpoint := c.Path()
		clientIP := c.IP()
		
		// Determine max requests for this endpoint
		maxRequests := config.DefaultMax
		if limit, exists := config.MaxRequests[endpoint]; exists {
			maxRequests = limit
		}

		// Create a unique key for this client and endpoint
		key := fmt.Sprintf("ratelimit:%s:%s", clientIP, endpoint)
		
		// Use Redis to implement a sliding window rate limit
		ctx := context.Background()
		now := time.Now().UnixNano()
		windowStart := now - config.WindowSize.Nanoseconds()

		pipe := config.RedisClient.Pipeline()
		
		// Remove old requests outside the window
		pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))
		
		// Add current request
		pipe.ZAdd(ctx, key, &redis.Z{
			Score:  float64(now),
			Member: now,
		})
		
		// Count requests in current window
		pipe.ZCard(ctx, key)
		
		// Set key expiration
		pipe.Expire(ctx, key, config.WindowSize)
		
		results, err := pipe.Exec(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Rate limiting error",
				"details": "Failed to process request rate limiting",
			})
		}

		// Get the current request count
		requestCount := results[2].(*redis.IntCmd).Val()
		
		if requestCount > int64(maxRequests) {
			retryAfter := time.Duration(windowStart+config.WindowSize.Nanoseconds()-now) * time.Nanosecond
			c.Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
			
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
				"details": fmt.Sprintf("Maximum of %d requests per %v allowed", maxRequests, config.WindowSize),
				"retry_after_seconds": int(retryAfter.Seconds()),
			})
		}

		return c.Next()
	}
}

// NewDefaultRateLimiter creates a rate limiter with default configuration
func NewDefaultRateLimiter(redisClient *redis.Client) fiber.Handler {
	config := RateLimiterConfig{
		RedisClient: redisClient,
		MaxRequests: map[string]int{
			"/ip":       100,  // IP endpoint
			"/domain":   200,  // Domain endpoint
			"/autnum":   150,  // AS number endpoint
			"/nameserver": 100, // Nameserver endpoint
		},
		WindowSize:  time.Minute,
		DefaultMax: 50,  // Default limit for unspecified endpoints
	}
	
	return RateLimiter(config)
}