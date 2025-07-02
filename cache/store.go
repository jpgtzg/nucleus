package cache

import (
	"nucleus/types/cache"
	"time"
)

type CacheEntry = cache.CacheEntry
type EventCache = cache.EventCache

// globalCache is the global cache for Stripe events, so we can avoid processing the same event multiple times
var globalCache *EventCache

func init() {
	globalCache = &EventCache{
		Store: make(map[string]CacheEntry),
	}

	// Checks for expired entries every 30 seconds and removes them
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			globalCache.Cleanup()
		}
	}()
}

// GetCacheEntry checks if an event ID exists in the cache
func GetCacheEntry(eventID string) (bool, error) {
	globalCache.Mu.RLock()
	defer globalCache.Mu.RUnlock()

	entry, exists := globalCache.Store[eventID]
	if !exists {
		return false, nil
	}

	// Check if the entry has expired
	if time.Now().After(entry.ExpiresAt) {
		return false, nil
	}

	return true, nil
}

// SetCacheEntry adds an event ID to the cache with a 30-second TTL
func SetCacheEntry(eventID string) (bool, error) {
	globalCache.Mu.Lock()
	defer globalCache.Mu.Unlock()

	// Check if the entry already exists
	if _, exists := globalCache.Store[eventID]; exists {
		return false, nil
	}

	now := time.Now()
	entry := CacheEntry{
		CreatedAt: now,
		ExpiresAt: now.Add(30 * time.Second),
	}

	globalCache.Store[eventID] = entry
	return true, nil
}

// DeleteCacheEntry removes an event ID from the cache
func DeleteCacheEntry(eventID string) {
	globalCache.Mu.Lock()
	defer globalCache.Mu.Unlock()
	delete(globalCache.Store, eventID)
}

// GetCacheStats returns basic statistics about the cache
func GetCacheStats() (int, error) {
	globalCache.Mu.RLock()
	defer globalCache.Mu.RUnlock()

	return len(globalCache.Store), nil
}
