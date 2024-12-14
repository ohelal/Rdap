package errors

import (
	"sync"
	"time"
)

// ErrorTracker handles error statistics and aggregation
type ErrorTracker struct {
	mu    sync.RWMutex
	stats map[ErrorCategory]*ErrorStats
}

func NewErrorTracker() *ErrorTracker {
	return &ErrorTracker{
		stats: make(map[ErrorCategory]*ErrorStats),
	}
}

func (t *ErrorTracker) TrackError(err *Error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if stat, exists := t.stats[err.Category]; exists {
		stat.Count++
		stat.LastSeen = time.Now()
		if err.Source != "" {
			stat.Sources[err.Source]++
		}
	} else {
		t.stats[err.Category] = &ErrorStats{
			Count:     1,
			FirstSeen: time.Now(),
			LastSeen:  time.Now(),
			Sources:   make(map[string]int64),
		}
		if err.Source != "" {
			t.stats[err.Category].Sources[err.Source] = 1
		}
	}
} 