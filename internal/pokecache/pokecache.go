package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries  map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

func NewCache(interval time.Duration) *Cache {
	newCache := &Cache{
		entries:  map[string]cacheEntry{},
		interval: interval}

	go newCache.reapLoop()

	return newCache
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.entries[key]
	if exists {
		return entry.val, true
	}

	return nil, false

}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)

	for range ticker.C {
		c.mu.Lock()

		for key, entry := range c.entries {
			timeSinceEntry := time.Since(entry.createdAt)
			if timeSinceEntry > c.interval {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}

}
