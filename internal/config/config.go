// Copyright (C) 2024 Helal <mohamed@helal.me>
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package config provides configuration management for the RDAP service.
package config

import (
	"time"
)

// Config holds the service configuration
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Kafka     KafkaConfig     `mapstructure:"kafka"`
	Metrics   MetricsConfig   `mapstructure:"metrics"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	Security  SecurityConfig  `mapstructure:"security"`
	RDAP      RDAPConfig      `mapstructure:"rdap"`
	RateLimit RateLimitConfig `mapstructure:"rateLimit"`
	Error     ErrorConfig     `mapstructure:"error"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port               string        `mapstructure:"port" default:"8080"`
	ReadTimeout        time.Duration `mapstructure:"read_timeout" default:"10s"`
	WriteTimeout       time.Duration `mapstructure:"write_timeout" default:"10s"`
	MaxConcurrentReqs  int           `mapstructure:"max_concurrent_requests" default:"5000"`
	EnableCompression  bool          `mapstructure:"enable_compression" default:"true"`
	EnableRateLimit    bool          `mapstructure:"enable_rate_limit" default:"true"`
	RateLimitPerMinute int           `mapstructure:"rate_limit_per_minute" default:"100"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	URL      string        `mapstructure:"url" default:"redis:6379"`
	Password string        `mapstructure:"password"`
	DB       int           `mapstructure:"db" default:"0"`
	TTL      time.Duration `mapstructure:"ttl" default:"3600s"`
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Enabled bool     `mapstructure:"enabled" default:"false"`
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic" default:"rdap-events"`
	GroupID string   `mapstructure:"group_id" default:"rdap-service"`
}

// MetricsConfig holds metrics configuration
type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled" default:"true"`
	Port    string `mapstructure:"port" default:"9090"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level" default:"info"`
	Format string `mapstructure:"format" default:"json"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	TLSEnabled bool   `mapstructure:"tls_enabled" default:"false"`
	CertFile   string `mapstructure:"cert_file"`
	KeyFile    string `mapstructure:"key_file"`
}

// RDAPConfig holds RDAP configuration
type RDAPConfig struct {
	BaseURL    string        `mapstructure:"baseUrl"`
	Timeout    time.Duration `mapstructure:"timeout"`
	MaxRetries int           `mapstructure:"maxRetries"`
	RetryDelay time.Duration `mapstructure:"retryDelay"`
}

// RateLimitConfig holds rate limit configuration
type RateLimitConfig struct {
	RequestsPerSecond int `mapstructure:"requestsPerSecond"`
	Burst             int `mapstructure:"burst"`
}

// ErrorConfig holds error handling configuration
type ErrorConfig struct {
	RetryableCodes []int         `mapstructure:"retryable_codes" default:"[500,502,503,504]"`
	MaxErrorAge    time.Duration `mapstructure:"max_error_age" default:"24h"`
	DetailedErrors bool          `mapstructure:"detailed_errors" default:"false"`
}

func (ec *ErrorConfig) IsRetryableCode(code int) bool {
	for _, c := range ec.RetryableCodes {
		if c == code {
			return true
		}
	}
	return false
}

// LoadConfig loads the configuration from environment variables and config file
func LoadConfig() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port:               "8080",
			ReadTimeout:        10 * time.Second,
			WriteTimeout:       10 * time.Second,
			MaxConcurrentReqs:  5000,
			EnableCompression:  true,
			EnableRateLimit:    true,
			RateLimitPerMinute: 100,
		},
		Redis: RedisConfig{
			URL:      "redis:6379",
			Password: "",
			DB:       0,
			TTL:      3600 * time.Second,
		},
		Kafka: KafkaConfig{
			Enabled: false,
			Brokers: []string{},
			Topic:   "rdap-events",
			GroupID: "rdap-service",
		},
		Metrics: MetricsConfig{
			Enabled: true,
			Port:    "9090",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Security: SecurityConfig{
			TLSEnabled: false,
			CertFile:   "",
			KeyFile:    "",
		},
		RDAP: RDAPConfig{
			BaseURL:    "https://rdap.arin.net/registry",
			Timeout:    10 * time.Second,
			MaxRetries: 3,
			RetryDelay: time.Second,
		},
		RateLimit: RateLimitConfig{
			RequestsPerSecond: 2000,
			Burst:             200,
		},
		Error: ErrorConfig{
			RetryableCodes: []int{500, 502, 503, 504},
			MaxErrorAge:    24 * time.Hour,
			DetailedErrors: false,
		},
	}, nil
}
