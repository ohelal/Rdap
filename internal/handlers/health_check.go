package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func HealthCheckHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "healthy",
	})
} 