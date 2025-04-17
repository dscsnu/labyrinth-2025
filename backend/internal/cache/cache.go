package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type Cache struct {
	c *cache.Cache
}

func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	return &Cache{
		c: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (c *Cache) Set(key string, value interface{}, expiration time.Duration) {
	c.c.Set(key, value, expiration)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	return c.c.Get(key)
}

func (c *Cache) Delete(key string) {
	c.c.Delete(key)
}

func (c *Cache) Flush() {
	c.c.Flush()
}

func (c *Cache) SetTTL(defaultExpiration, cleanupInterval time.Duration) {
	c.c = cache.New(defaultExpiration, cleanupInterval)
}

func DefaultExpiration() time.Duration {
	return cache.DefaultExpiration
}

func NoExpiration() time.Duration {
	return cache.NoExpiration
}