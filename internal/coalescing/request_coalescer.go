package coalescing

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrTimeout = errors.New("request timed out")
)

type RequestKey string

type InFlightRequest struct {
	Done     chan struct{}
	Response interface{}
	Error    error
}

type RequestCoalescer struct {
	mu          sync.RWMutex
	inFlight    map[RequestKey]*InFlightRequest
	maxWaitTime time.Duration
}

func NewRequestCoalescer(maxWaitTime time.Duration) *RequestCoalescer {
	return &RequestCoalescer{
		inFlight:    make(map[RequestKey]*InFlightRequest),
		maxWaitTime: maxWaitTime,
	}
}

func (rc *RequestCoalescer) Execute(ctx context.Context, key RequestKey, fn func() (interface{}, error)) (interface{}, error) {
	// Check if request is already in flight
	if inFlight := rc.getInFlight(key); inFlight != nil {
		select {
		case <-inFlight.Done:
			return inFlight.Response, inFlight.Error
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(rc.maxWaitTime):
			return nil, ErrTimeout
		}
	}

	// Create new in-flight request
	inFlight := &InFlightRequest{
		Done: make(chan struct{}),
	}

	rc.mu.Lock()
	rc.inFlight[key] = inFlight
	rc.mu.Unlock()

	// Execute the actual request
	inFlight.Response, inFlight.Error = fn()
	close(inFlight.Done)

	// Cleanup
	rc.mu.Lock()
	delete(rc.inFlight, key)
	rc.mu.Unlock()

	return inFlight.Response, inFlight.Error
}

func (rc *RequestCoalescer) getInFlight(key RequestKey) *InFlightRequest {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.inFlight[key]
}
