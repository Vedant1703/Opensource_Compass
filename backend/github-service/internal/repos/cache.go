package repos

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Data      []RepoDTO
	Timestamp time.Time
}

type RepoCache struct {
	mu    sync.RWMutex
	cache map[string]CacheEntry
	ttl   time.Duration
}

func NewRepoCache(ttl time.Duration) *RepoCache {
	return &RepoCache{
		cache: make(map[string]CacheEntry),
		ttl:   ttl,
	}
}

func (c *RepoCache) Get(key string) ([]RepoDTO, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	// Check if cache entry has expired
	if time.Since(entry.Timestamp) > c.ttl {
		return nil, false
	}

	return entry.Data, true
}

func (c *RepoCache) Set(key string, data []RepoDTO) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = CacheEntry{
		Data:      data,
		Timestamp: time.Now(),
	}
}

// Clean removes expired entries from cache
func (c *RepoCache) Clean() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.cache {
		if now.Sub(entry.Timestamp) > c.ttl {
			delete(c.cache, key)
		}
	}
}

// Returns size of cache
func (c *RepoCache) Size() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return len(c.cache)
}
