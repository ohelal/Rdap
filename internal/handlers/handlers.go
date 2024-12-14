// Copyright (C) 2024 Helal <mohamed@helal.me>
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package handlers provides HTTP request handlers for the RDAP service.
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ohelal/rdap/internal/kafka"
	"github.com/ohelal/rdap/internal/metrics"
	"github.com/ohelal/rdap/internal/service"
	"net"
	"strconv"
	"strings"
	"time"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	svc      *service.RDAPService
	metrics  *metrics.Metrics
	producer *kafka.Producer
}

// NewHandlers creates a new Handlers instance
func NewHandlers(svc *service.RDAPService, metrics *metrics.Metrics, producer *kafka.Producer) *Handlers {
	return &Handlers{svc: svc, metrics: metrics, producer: producer}
}

// LookupHandler handles RDAP lookup requests
func (h *Handlers) LookupHandler(c *fiber.Ctx) error {
	// Handle RDAP lookup requests
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing query parameter 'q'",
		})
	}

	// Validate query type
	queryType := getQueryType(query)
	if queryType == "invalid" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query format",
		})
	}

	var err error
	switch queryType {
	case "ip":
		err = h.IPLookupHandler(c)
	case "domain":
		err = h.DomainLookupHandler(c)
	case "asn":
		err = h.ASNLookupHandler(c)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unsupported query type",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return nil
}

// IPLookupHandler handles IP lookup requests
func (h *Handlers) IPLookupHandler(c *fiber.Ctx) error {
	// Handle IP lookup requests
	ip := c.Query("q")
	err := h.svc.HandleIPLookup(c)
	if err != nil {
		return err
	}

	h.sendToKafka("ip", ip, false)
	return nil
}

// DomainLookupHandler handles domain lookup requests
func (h *Handlers) DomainLookupHandler(c *fiber.Ctx) error {
	// Handle domain lookup requests
	domain := c.Query("q")
	err := h.svc.HandleDomainLookup(c)
	if err != nil {
		return err
	}

	h.sendToKafka("domain", domain, false)
	return nil
}

// ASNLookupHandler handles ASN lookup requests
func (h *Handlers) ASNLookupHandler(c *fiber.Ctx) error {
	// Handle ASN lookup requests
	asn := c.Query("q")
	err := h.svc.HandleASNLookup(c)
	if err != nil {
		return err
	}

	h.sendToKafka("asn", asn, false)
	return nil
}

// HealthHandler handles health check requests
func (h *Handlers) HealthHandler(c *fiber.Ctx) error {
	// Handle health check requests
	return c.JSON(fiber.Map{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *Handlers) sendToKafka(queryType, query string, cacheHit bool) {
	// Send message to Kafka
	msg := &kafka.Message{
		Type:      queryType,
		Query:     query,
		Timestamp: time.Now(),
		Source:    "api",
		CacheHit:  cacheHit,
	}

	if err := h.producer.SendMessage(msg); err != nil {
		// Just log the error, don't fail the request
		h.metrics.KafkaErrors.Inc()
	}
}

// getQueryType determines the type of RDAP query
func getQueryType(query string) string {
	// Check if query is an IP address
	if ip := net.ParseIP(query); ip != nil {
		return "ip"
	}

	// Check if query is an ASN
	if _, err := strconv.ParseInt(strings.TrimPrefix(query, "AS"), 10, 64); err == nil {
		return "asn"
	}

	// Check if query is a domain name
	if strings.Contains(query, ".") {
		return "domain"
	}

	return "invalid"
}
