package cache

import "time"

func SetupCaches(cm *CacheManager, defaultTTL time.Duration, cleanupInterval time.Duration) {
	entities := []string{}

	for _, entity := range entities {
		cm.InitCache(entity, defaultTTL, cleanupInterval)
	}
}
