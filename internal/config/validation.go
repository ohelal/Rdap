package config

import (
	"fmt"
)

func (cfg *Config) Validate() error {
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}
	if cfg.RDAP.BaseURL == "" {
		return fmt.Errorf("RDAP base URL cannot be empty")
	}
	if cfg.RateLimit.RequestsPerSecond <= 0 {
		return fmt.Errorf("rate limit must be greater than zero")
	}
	return nil
}
