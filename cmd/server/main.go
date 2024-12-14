// Copyright (C) 2024 Helal <mohamed@helal.me>
// SPDX-License-Identifier: AGPL-3.0-or-later

// Command server runs the RDAP service with Redis caching and Kafka integration.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/ohelal/rdap/internal/cache"
	"github.com/ohelal/rdap/internal/config"
	"github.com/ohelal/rdap/internal/errors"
	"github.com/ohelal/rdap/internal/handlers"
	"github.com/ohelal/rdap/internal/kafka"
	"github.com/ohelal/rdap/internal/metrics"
	"github.com/ohelal/rdap/internal/middleware"
	"github.com/ohelal/rdap/internal/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Default configuration values
const (
	defaultRedisURL = "redis:6379"
	defaultPort    = "8080"
	defaultMetricsPort = "9090"
)

func main() {
	// Create a context that we'll use to cancel goroutines on shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize metrics collector
	metricsCollector := metrics.NewMetrics()

	// Initialize Redis client with fallback to default
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = defaultRedisURL
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:         redisURL,
		DB:           0,
		PoolSize:     50,
		MinIdleConns: 10,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
	defer redisClient.Close()

	// Initialize Kafka producer
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	if len(kafkaBrokers) == 0 || kafkaBrokers[0] == "" {
		kafkaBrokers = []string{"localhost:9092"}
	}

	producer, err := kafka.NewProducer(kafkaBrokers, "rdap-queries")
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Initialize Kafka consumer
	consumer, err := kafka.NewConsumer(
		kafkaBrokers,
		"rdap-consumer-group",
		[]string{"rdap-queries"},
		kafka.NewMessageHandler(),
	)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	// Start the Kafka consumer in a goroutine
	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Printf("Error starting Kafka consumer: %v", err)
		}
	}()

	// Initialize cache
	cacheConfig := &cache.CacheConfig{
		MaxLocalSize: 1024 * 1024 * 1024, // 1GB
		LocalTTL:     time.Hour,
		EnableRedis:  true,
		RedisURL:     redisURL,
		RedisTTL:     time.Hour,
	}

	cacheManager, err := cache.NewCacheManager(cacheConfig)
	if err != nil {
		log.Fatalf("Failed to initialize cache: %v", err)
	}
	defer cacheManager.Close()

	// Load bootstrap configurations
	configDir := os.Getenv("CONFIG_DIR")
	if configDir == "" {
		configDir = "/app/config"
	}

	dnsConfig, ipConfig, asnConfig, err := service.LoadAllBootstrapConfigs(configDir)
	if err != nil {
		log.Fatalf("Failed to load bootstrap configs: %v", err)
	}

	// Initialize RDAP service
	rdapService, err := service.NewRDAPService(dnsConfig, ipConfig, asnConfig, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize RDAP service: %v", err)
	}

	// Initialize handlers
	handlers := handlers.NewHandlers(rdapService, metricsCollector, producer)

	// Initialize Fiber app with logging
	log.Println("Starting RDAP service...")

	app := fiber.New(fiber.Config{
		ErrorHandler:          errors.HandleError,
		DisableStartupMessage: false,
	})

	// Add logging middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path}\n",
	}))

	// Middleware
	app.Use(compress.New())
	app.Use(recover.New())
	app.Use(middleware.NewDefaultRateLimiter(redisClient))

	// Routes
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	app.Get("/health", handlers.HealthHandler)
	app.Get("/ip/:ip", handlers.IPLookupHandler)
	app.Get("/domain/:domain", handlers.DomainLookupHandler)
	app.Get("/autnum/:asn", handlers.ASNLookupHandler)

	// Add graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start server
	log.Printf("Server starting on port %d...", cfg.Server.Port)
	if err := app.Listen(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

	// Add graceful shutdown handler
	go func() {
		<-quit
		log.Println("Shutting down server...")
		if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}()
}
