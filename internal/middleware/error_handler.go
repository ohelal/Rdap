package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ohelal/rdap/internal/errors"
	"github.com/ohelal/rdap/internal/logger"
)

func ErrorHandler(log *logger.ContextualLogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			if e, ok := err.(*errors.Error); ok {
				log.WithContext(c.Context()).Err(err).
					Str("trace_id", e.TraceID).
					Msg(e.Message)
				
				return c.Status(e.Code).JSON(fiber.Map{
					"error": e.Message,
					"trace_id": e.TraceID,
				})
			}
			log.WithContext(c.Context()).Err(err).Msg("Unhandled error")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		}
		return nil
	}
}