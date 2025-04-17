package cache

import (
	"sync"
	"time"
)

type CacheManager struct {
	caches map[string]*Cache
	mu     sync.RWMutex
}

func NewManager() *CacheManager {
	return &CacheManager{
		caches: make(map[string]*Cache),
	}
}

func (cm *CacheManager) InitCache(entity string, defaultTTL, cleanupInterval time.Duration) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if _, exists := cm.caches[entity]; !exists {
		cm.caches[entity] = New(defaultTTL, cleanupInterval)
	}
}

func (cm *CacheManager) Set(entity string, key string, value interface{}, expiration time.Duration) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if c, exists := cm.caches[entity]; exists {
		c.Set(key, value, expiration)
	}
}

func (cm *CacheManager) Get(entity string, key string) (interface{}, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if c, exists := cm.caches[entity]; exists {
		return c.Get(key)
	}
	return nil, false
}

func (cm *CacheManager) Delete(entity string, key string) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if c, exists := cm.caches[entity]; exists {
		c.Delete(key)
	}
}

func (cm *CacheManager) Flush(entity string) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if c, exists := cm.caches[entity]; exists {
		c.Flush()
	}
}
