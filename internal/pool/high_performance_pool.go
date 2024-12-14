package pool

import (
    "sync"
    "sync/atomic"
)

type HighPerformancePool struct {
    maxWorkers uint32
    active     uint32
    pool       chan struct{}
    mu         sync.RWMutex
}

func NewHighPerformancePool(maxWorkers uint32) *HighPerformancePool {
    return &HighPerformancePool{
        maxWorkers: maxWorkers,
        pool:       make(chan struct{}, maxWorkers),
    }
}

func (p *HighPerformancePool) Acquire() bool {
    if atomic.LoadUint32(&p.active) >= p.maxWorkers {
        return false
    }
    
    select {
    case p.pool <- struct{}{}:
        atomic.AddUint32(&p.active, 1)
        return true
    default:
        return false
    }
}

func (p *HighPerformancePool) Release() {
    <-p.pool
    atomic.AddUint32(&p.active, ^uint32(0))
} 