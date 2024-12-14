package middleware

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"
)

type ResourceLimiter struct {
	limiter *rate.Limiter
}

func NewResourceLimiter(requestsPerSecond int, _ int) *ResourceLimiter {
	return &ResourceLimiter{
		limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond),
	}
}

func (rl *ResourceLimiter) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !rl.limiter.Allow() {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
			})
		}
		return c.Next()
	}
}
