package kafka

import (
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateHalfOpen
	StateOpen
)

type CircuitBreaker struct {
	mutex sync.RWMutex

	state                     State
	failureThreshold         int
	failureCount            int
	resetTimeout            time.Duration
	lastStateTransitionTime time.Time
	halfOpenMaxRequests     int
	halfOpenRequestCount    int
}

func NewCircuitBreaker(failureThreshold int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:             StateClosed,
		failureThreshold: failureThreshold,
		resetTimeout:     resetTimeout,
		halfOpenMaxRequests: 5,
	}
}

func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if now.Sub(cb.lastStateTransitionTime) >= cb.resetTimeout {
			cb.toHalfOpen(now)
			return true
		}
		return false
	case StateHalfOpen:
		if cb.halfOpenRequestCount < cb.halfOpenMaxRequests {
			cb.halfOpenRequestCount++
			return true
		}
		return false
	default:
		return false
	}
}

func (cb *CircuitBreaker) OnSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch cb.state {
	case StateHalfOpen:
		cb.toClosed(time.Now())
	case StateClosed:
		cb.failureCount = 0
	}
}

func (cb *CircuitBreaker) OnFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch cb.state {
	case StateClosed:
		cb.failureCount++
		if cb.failureCount >= cb.failureThreshold {
			cb.toOpen(time.Now())
		}
	case StateHalfOpen:
		cb.toOpen(time.Now())
	}
}

func (cb *CircuitBreaker) toClosed(now time.Time) {
	cb.state = StateClosed
	cb.failureCount = 0
	cb.lastStateTransitionTime = now
	cb.halfOpenRequestCount = 0
}

func (cb *CircuitBreaker) toOpen(now time.Time) {
	cb.state = StateOpen
	cb.lastStateTransitionTime = now
	cb.halfOpenRequestCount = 0
}

func (cb *CircuitBreaker) toHalfOpen(now time.Time) {
	cb.state = StateHalfOpen
	cb.lastStateTransitionTime = now
	cb.halfOpenRequestCount = 0
}
