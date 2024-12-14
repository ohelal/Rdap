package circuit

import (
    "context"
    "errors"
    "sync"
    "time"
)

var (
    ErrCircuitOpen = errors.New("circuit breaker is open")
)

type State int

const (
    StateClosed State = iota
    StateHalfOpen
    StateOpen
)

type CircuitBreaker struct {
    mu             sync.RWMutex
    state          State
    failures       int
    lastFailure    time.Time
    timeout        time.Duration
    maxFailures    int
    halfOpenQuota  int
}

func NewCircuitBreaker(timeout time.Duration, maxFailures, halfOpenQuota int) *CircuitBreaker {
    return &CircuitBreaker{
        state:         StateClosed,
        timeout:       timeout,
        maxFailures:   maxFailures,
        halfOpenQuota: halfOpenQuota,
    }
}

func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
    if !cb.AllowRequest() {
        return ErrCircuitOpen
    }

    err := fn()
    cb.RecordResult(err)
    return err
}

func (cb *CircuitBreaker) RecordResult(err error) {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()
        if cb.failures >= cb.maxFailures {
            cb.state = StateOpen
        }
    } else {
        cb.failures = 0
        cb.state = StateClosed
    }
}

func (cb *CircuitBreaker) AllowRequest() bool {
    cb.mu.RLock()
    defer cb.mu.RUnlock()

    switch cb.state {
    case StateClosed:
        return true
    case StateOpen:
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.mu.Lock()
            cb.state = StateHalfOpen
            cb.mu.Unlock()
            return true
        }
        return false
    case StateHalfOpen:
        return cb.failures < cb.halfOpenQuota
    default:
        return false
    }
}