package config

import "time"

// internal/config/error_config.go
type ErrorConfig struct {
	RetryableCodes []int         `yaml:"retryable_codes"`
	MaxErrorAge    time.Duration `yaml:"max_error_age"`
	DetailedErrors bool          `yaml:"detailed_errors"`
}

func (ec *ErrorConfig) IsRetryableCode(code int) bool {
	for _, c := range ec.RetryableCodes {
		if c == code {
			return true
		}
	}
	return false
}
