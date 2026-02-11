package rabbitmq

import (
	"sync"
	"time"
)

// DedupCache is an in-memory cache that tracks processed event_ids
// to prevent duplicate event handling (idempotency).
// Entries expire after the configured TTL to avoid unbounded memory growth.
//
// For production at scale, consider replacing this with a Redis-based
// implementation (e.g. SET event_id EX 3600 NX) for cross-instance dedup.
type DedupCache struct {
	mu      sync.Mutex
	entries map[string]time.Time
	ttl     time.Duration
}

// NewDedupCache creates a new DedupCache with the given TTL for entries.
// It starts a background goroutine that periodically evicts expired entries.
func NewDedupCache(ttl time.Duration) *DedupCache {
	c := &DedupCache{
		entries: make(map[string]time.Time),
		ttl:     ttl,
	}
	go c.evictLoop()
	return c
}

// IsDuplicate returns true if the eventID has already been seen (and is not expired).
// If it has not been seen, it marks it as seen and returns false.
func (c *DedupCache) IsDuplicate(eventID string) bool {
	if eventID == "" {
		// No event_id means we cannot deduplicate; treat as unique.
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if ts, ok := c.entries[eventID]; ok {
		if time.Since(ts) < c.ttl {
			return true
		}
		// Expired entry, allow re-processing
	}
	c.entries[eventID] = time.Now()
	return false
}

// evictLoop removes expired entries every ttl/2 interval.
func (c *DedupCache) evictLoop() {
	ticker := time.NewTicker(c.ttl / 2)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for id, ts := range c.entries {
			if now.Sub(ts) >= c.ttl {
				delete(c.entries, id)
			}
		}
		c.mu.Unlock()
	}
}
