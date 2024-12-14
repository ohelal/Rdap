package config

import (
	"time"
)

// Config holds the service configuration
type Config struct {
	Server    ServerConfig    `json:"server"`
	RDAP      RDAPConfig      `json:"rdap"`
	RateLimit RateLimitConfig `json:"rateLimit"`
	Logging   LoggingConfig   `json:"logging"`
}

type ServerConfig struct {
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"readTimeout"`
	WriteTimeout time.Duration `json:"writeTimeout"`
	IdleTimeout  time.Duration `json:"idleTimeout"`
}

type RDAPConfig struct {
	BaseURL    string        `json:"baseUrl"`
	Timeout    time.Duration `json:"timeout"`
	MaxRetries int           `json:"maxRetries"`
	RetryDelay time.Duration `json:"retryDelay"`
}

type RateLimitConfig struct {
	RequestsPerSecond int `json:"requestsPerSecond"`
	Burst             int `json:"burst"`
}

type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

// LoadConfig loads all configuration files
func LoadConfig() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  10 * time.Second,
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
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
	}, nil
}
