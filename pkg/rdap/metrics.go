package rdap

import (
    "sync/atomic"
    "time"
)

// Metrics holds client usage statistics
type Metrics struct {
    TotalRequests      uint64
    SuccessfulRequests uint64
    FailedRequests     uint64
    AverageLatencyNs   uint64
    lastReset         time.Time
}

// NewMetrics creates a new metrics instance
func NewMetrics() *Metrics {
    return &Metrics{
        lastReset: time.Now(),
    }
}

// RecordRequest records metrics for a request
func (m *Metrics) RecordRequest(start time.Time, err error) {
    atomic.AddUint64(&m.TotalRequests, 1)
    if err != nil {
        atomic.AddUint64(&m.FailedRequests, 1)
    } else {
        atomic.AddUint64(&m.SuccessfulRequests, 1)
    }
    
    duration := time.Since(start)
    current := atomic.LoadUint64(&m.AverageLatencyNs)
    atomic.StoreUint64(&m.AverageLatencyNs, (current+uint64(duration.Nanoseconds()))/2)
}

// Reset resets all metrics
func (m *Metrics) Reset() {
    atomic.StoreUint64(&m.TotalRequests, 0)
    atomic.StoreUint64(&m.SuccessfulRequests, 0)
    atomic.StoreUint64(&m.FailedRequests, 0)
    atomic.StoreUint64(&m.AverageLatencyNs, 0)
    m.lastReset = time.Now()
}

// GetAverageLatency returns the average latency as a Duration
func (m *Metrics) GetAverageLatency() time.Duration {
    return time.Duration(atomic.LoadUint64(&m.AverageLatencyNs))
}
