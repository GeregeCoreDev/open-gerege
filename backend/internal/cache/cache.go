// Package cache provides in-memory caching with TTL support
//
// File: cache.go
// Description: Generic cache implementation for frequently accessed data
//
// This package provides a thread-safe, generic cache with:
//   - TTL (Time-To-Live) support
//   - LRU eviction when max size reached
//   - Atomic operations
//   - Metrics tracking (hits, misses)
//
// Usage:
//
//	userCache := cache.New[domain.User](cache.Config{
//	    MaxSize: 1000,
//	    TTL:     5 * time.Minute,
//	})
//	userCache.Set("user:1", user)
//	user, found := userCache.Get("user:1")
package cache

import (
	"sync"
	"sync/atomic"
	"time"
)

// Config holds cache configuration
type Config struct {
	// MaxSize is the maximum number of items in cache
	MaxSize int
	// TTL is the time-to-live for cache entries
	TTL time.Duration
	// CleanupInterval is how often to run cleanup (default: TTL/2)
	CleanupInterval time.Duration
}

// DefaultConfig returns sensible defaults
func DefaultConfig() Config {
	return Config{
		MaxSize:         1000,
		TTL:             5 * time.Minute,
		CleanupInterval: 2 * time.Minute,
	}
}

// item represents a cached item with expiration
type item[T any] struct {
	value      T
	expiration int64
}

// isExpired checks if the item has expired
func (i item[T]) isExpired() bool {
	return time.Now().UnixNano() > i.expiration
}

// Cache is a generic in-memory cache with TTL support
type Cache[T any] struct {
	items   map[string]item[T]
	mu      sync.RWMutex
	config  Config
	hits    atomic.Int64
	misses  atomic.Int64
	stopCh  chan struct{}
	stopped atomic.Bool
}

// New creates a new cache with the given configuration
func New[T any](cfg Config) *Cache[T] {
	if cfg.MaxSize <= 0 {
		cfg.MaxSize = 1000
	}
	if cfg.TTL <= 0 {
		cfg.TTL = 5 * time.Minute
	}
	if cfg.CleanupInterval <= 0 {
		cfg.CleanupInterval = cfg.TTL / 2
	}

	c := &Cache[T]{
		items:  make(map[string]item[T]),
		config: cfg,
		stopCh: make(chan struct{}),
	}

	// Start cleanup goroutine
	go c.cleanup()

	return c
}

// Get retrieves an item from the cache
// Returns the item and true if found and not expired, zero value and false otherwise
func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if !found {
		c.misses.Add(1)
		var zero T
		return zero, false
	}

	if item.isExpired() {
		c.misses.Add(1)
		c.Delete(key)
		var zero T
		return zero, false
	}

	c.hits.Add(1)
	return item.value, true
}

// Set adds or updates an item in the cache
func (c *Cache[T]) Set(key string, value T) {
	c.SetWithTTL(key, value, c.config.TTL)
}

// SetWithTTL adds or updates an item with a custom TTL
func (c *Cache[T]) SetWithTTL(key string, value T, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Evict if at max size
	if len(c.items) >= c.config.MaxSize {
		c.evictOldest()
	}

	c.items[key] = item[T]{
		value:      value,
		expiration: time.Now().Add(ttl).UnixNano(),
	}
}

// Delete removes an item from the cache
func (c *Cache[T]) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

// DeletePrefix removes all items with keys starting with prefix
func (c *Cache[T]) DeletePrefix(prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key := range c.items {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(c.items, key)
		}
	}
}

// Clear removes all items from the cache
func (c *Cache[T]) Clear() {
	c.mu.Lock()
	c.items = make(map[string]item[T])
	c.mu.Unlock()
}

// Size returns the number of items in the cache
func (c *Cache[T]) Size() int {
	c.mu.RLock()
	size := len(c.items)
	c.mu.RUnlock()
	return size
}

// Stats returns cache statistics
type Stats struct {
	Size   int
	Hits   int64
	Misses int64
	Ratio  float64 // Hit ratio (0.0 to 1.0)
}

// Stats returns cache statistics
func (c *Cache[T]) Stats() Stats {
	hits := c.hits.Load()
	misses := c.misses.Load()
	total := hits + misses

	var ratio float64
	if total > 0 {
		ratio = float64(hits) / float64(total)
	}

	return Stats{
		Size:   c.Size(),
		Hits:   hits,
		Misses: misses,
		Ratio:  ratio,
	}
}

// Stop stops the cleanup goroutine
func (c *Cache[T]) Stop() {
	if c.stopped.CompareAndSwap(false, true) {
		close(c.stopCh)
	}
}

// cleanup periodically removes expired items
func (c *Cache[T]) cleanup() {
	ticker := time.NewTicker(c.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-c.stopCh:
			return
		}
	}
}

// deleteExpired removes all expired items
func (c *Cache[T]) deleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixNano()
	for key, item := range c.items {
		if now > item.expiration {
			delete(c.items, key)
		}
	}
}

// evictOldest removes the oldest item (simple LRU approximation)
func (c *Cache[T]) evictOldest() {
	var oldestKey string
	var oldestExp int64 = time.Now().Add(time.Hour).UnixNano()

	for key, item := range c.items {
		if item.expiration < oldestExp {
			oldestExp = item.expiration
			oldestKey = key
		}
	}

	if oldestKey != "" {
		delete(c.items, oldestKey)
	}
}

// GetOrSet gets an item from cache, or sets it using the provided function if not found
func (c *Cache[T]) GetOrSet(key string, fn func() (T, error)) (T, error) {
	// Try to get from cache first
	if value, found := c.Get(key); found {
		return value, nil
	}

	// Not in cache, call function to get value
	value, err := fn()
	if err != nil {
		var zero T
		return zero, err
	}

	// Store in cache
	c.Set(key, value)
	return value, nil
}
