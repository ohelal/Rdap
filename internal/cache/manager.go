package cache

import (
	"sync"

	"github.com/VictoriaMetrics/fastcache"
)

// CacheManager manages both local and distributed caches
type CacheManager struct {
	local       *fastcache.Cache
	distributed *DistributedCache
	mu          sync.RWMutex
	config      CacheConfig
}

// NewCacheManager creates a new cache manager
func NewCacheManager(config *CacheConfig) (*CacheManager, error) {
	local := fastcache.New(int(config.MaxLocalSize))

	var distributed *DistributedCache
	var err error
	if config.EnableRedis {
		distributed, err = NewDistributedCache(config)
		if err != nil {
			return nil, err
		}
	}

	return &CacheManager{
		local:       local,
		distributed: distributed,
		config:      *config,
	}, nil
}

// Get retrieves a value from the cache
func (cm *CacheManager) Get(key string) (interface{}, bool) {
	// Try local cache first
	if val := cm.local.Get(nil, []byte(key)); val != nil {
		return val, true
	}

	// Try distributed cache if available
	if cm.distributed != nil {
		if val, found := cm.distributed.Get(key); found {
			// Update local cache
			cm.local.Set([]byte(key), []byte(val.(string)))
			return val, true
		}
	}

	return nil, false
}

// Set stores a value in both caches
func (cm *CacheManager) Set(key string, value interface{}) error {
	// Update local cache
	cm.local.Set([]byte(key), []byte(value.(string)))

	// Update distributed cache if available
	if cm.distributed != nil {
		return cm.distributed.Set(key, value)
	}

	return nil
}

// Delete removes a value from both caches
func (cm *CacheManager) Delete(key string) error {
	// Remove from local cache
	cm.local.Del([]byte(key))

	// Remove from distributed cache if available
	if cm.distributed != nil {
		return cm.distributed.Delete(key)
	}

	return nil
}

// Close closes all caches
func (cm *CacheManager) Close() error {
	if cm.distributed != nil {
		return cm.distributed.Close()
	}
	return nil
}
