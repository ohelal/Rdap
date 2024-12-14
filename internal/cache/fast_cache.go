package cache

import (
	"github.com/VictoriaMetrics/fastcache"
	"sync/atomic"
)

type FastCache struct {
	cache    *fastcache.Cache
	hits     uint64
	misses   uint64
	maxBytes int
}

func NewFastCache() *FastCache {
	return &FastCache{
		cache:    fastcache.New(1024 * 1024 * 1024), // 1GB cache
		maxBytes: 1024 * 1024 * 1024,
	}
}

func (c *FastCache) Get(key []byte) ([]byte, bool) {
	value := c.cache.Get(nil, key)
	if value == nil {
		atomic.AddUint64(&c.misses, 1)
		return nil, false
	}
	atomic.AddUint64(&c.hits, 1)
	return value, true
}

func (c *FastCache) Set(key, value []byte) {
	c.cache.Set(key, value)
}
