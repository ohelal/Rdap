package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RequestTracing() fiber.Handler {
	return func(c *fiber.Ctx) error {
		traceID := c.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
			c.Set("X-Trace-ID", traceID)
		}
		c.Locals("trace_id", traceID)
		return c.Next()
	}
} 