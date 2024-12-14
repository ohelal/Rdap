package errors

import (
	"github.com/sony/gobreaker"
	"time"
)

type CircuitBreakerConfig struct {
	MaxRequests uint32
	Interval    time.Duration
	Timeout     time.Duration
}

type ErrorCircuitBreaker struct {
	breakers map[ErrorCategory]*gobreaker.CircuitBreaker
	config   CircuitBreakerConfig
}

func NewErrorCircuitBreaker(config CircuitBreakerConfig) *ErrorCircuitBreaker {
	return &ErrorCircuitBreaker{
		breakers: make(map[ErrorCategory]*gobreaker.CircuitBreaker),
		config:   config,
	}
}

func (ecb *ErrorCircuitBreaker) GetBreaker(category ErrorCategory) *gobreaker.CircuitBreaker {
	if cb, exists := ecb.breakers[category]; exists {
		return cb
	}

	settings := gobreaker.Settings{
		Name:        category.String(),
		MaxRequests: ecb.config.MaxRequests,
		Interval:    ecb.config.Interval,
		Timeout:     ecb.config.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	}

	ecb.breakers[category] = gobreaker.NewCircuitBreaker(settings)
	return ecb.breakers[category]
}
