package config

import (
	"fmt"
	"strconv"
)

func (cfg *Config) Validate() error {
	port, err := strconv.Atoi(cfg.Server.Port)
	if err != nil || port <= 0 || port > 65535 {
		return fmt.Errorf("invalid server port: %s", cfg.Server.Port)
	}
	if cfg.RDAP.BaseURL == "" {
		return fmt.Errorf("RDAP base URL cannot be empty")
	}
	if cfg.RateLimit.RequestsPerSecond <= 0 {
		return fmt.Errorf("rate limit must be greater than zero")
	}
	return nil
}
