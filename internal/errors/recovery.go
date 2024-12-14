package errors

import (
	"context"
	"time"
)

type RecoveryStrategy interface {
	Recover(ctx context.Context, err *Error) error
}

type RetryStrategy struct {
	maxRetries  int
	backoffBase time.Duration
}

func NewRetryStrategy(maxRetries int, backoffBase time.Duration) *RetryStrategy {
	return &RetryStrategy{
		maxRetries:  maxRetries,
		backoffBase: backoffBase,
	}
}

func (rs *RetryStrategy) Recover(ctx context.Context, err *Error) error {
	if !err.Retryable {
		return err
	}

	for attempt := 0; attempt < rs.maxRetries; attempt++ {
		backoff := rs.backoffBase * time.Duration(1<<attempt)
		timer := time.NewTimer(backoff)

		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
			// Attempt recovery logic here
			retryCounter.WithLabelValues(err.Category.String()).Inc()
		}
	}

	return err
} 