package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/ohelal/rdap/internal/cache"
	"github.com/ohelal/rdap/internal/cdn"
	"github.com/ohelal/rdap/internal/coalescing"
	"time"
)

// CoalescedHandler handles coalesced requests
type CoalescedHandler struct {
	cache     *cache.CacheManager
	coalescer *coalescing.RequestCoalescer
	cdnConfig *cdn.CDNConfig
}

// NewCoalescedHandler creates a new coalesced handler
func NewCoalescedHandler(cacheManager *cache.CacheManager, timeout time.Duration) (*CoalescedHandler, error) {
	cdnConfig := cdn.NewCDNConfig("https://cdn.example.com", 30) // 30 seconds timeout
	if cdnConfig == nil {
		return nil, errors.New("failed to create cdn config")
	}

	return &CoalescedHandler{
		cache:     cacheManager,
		coalescer: coalescing.NewRequestCoalescer(timeout),
		cdnConfig: cdnConfig,
	}, nil
}

// HandleRequest handles a coalesced request
func (h *CoalescedHandler) HandleRequest(c *fiber.Ctx) error {
	key := c.Query("key")
	if key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing key parameter",
		})
	}

	result, err := h.coalescer.Execute(c.Context(), coalescing.RequestKey(key), func() (interface{}, error) {
		// Check cache first
		if cached, found := h.cache.Get(key); found {
			return cached, nil
		}

		// Make actual request and cache result
		result, err := h.makeRequest(key)
		if err != nil {
			return nil, err
		}

		h.cache.Set(key, result)
		return result, nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

// makeRequest makes the actual request
func (h *CoalescedHandler) makeRequest(key string) (interface{}, error) {
	// Implementation depends on your specific needs
	// This is just a placeholder
	return map[string]string{"key": key}, nil
}
