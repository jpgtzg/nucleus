package cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	CreatedAt time.Time
	ExpiresAt time.Time
}

type EventCache struct {
	Mu    sync.RWMutex          // Mutex to protect race conditions
	Store map[string]CacheEntry // Map to store event IDs and their corresponding cache entries
}

// cleanup removes expired entries from the cache
func (c *EventCache) Cleanup() {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	now := time.Now()
	for eventID, entry := range c.Store {
		if now.After(entry.ExpiresAt) {
			delete(c.Store, eventID)
		}
	}
}
