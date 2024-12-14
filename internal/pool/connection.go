package pool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrPoolExhausted = errors.New("connection pool exhausted")
	idCounter uint64
)

func generateID() uint64 {
	return atomic.AddUint64(&idCounter, 1)
}

type ConnectionPool struct {
	mu          sync.RWMutex
	connections chan *Connection
	maxSize     int
	timeout     time.Duration
}

type Connection struct {
	ID        string
	CreatedAt time.Time
	LastUsed  time.Time
}

func NewConnectionPool(maxSize int, timeout time.Duration) *ConnectionPool {
	return &ConnectionPool{
		connections: make(chan *Connection, maxSize),
		maxSize:     maxSize,
		timeout:     timeout,
	}
}

func (p *ConnectionPool) Acquire(ctx context.Context) (*Connection, error) {
	select {
	case conn := <-p.connections:
		conn.LastUsed = time.Now()
		return conn, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// Create new connection if pool not full
		p.mu.Lock()
		defer p.mu.Unlock()
		
		if len(p.connections) < p.maxSize {
			conn := &Connection{
				ID:        fmt.Sprintf("%d", generateID()),
				CreatedAt: time.Now(),
				LastUsed:  time.Now(),
			}
			return conn, nil
		}
		return nil, ErrPoolExhausted
	}
} 