// Package cache provides in-memory caching with TTL support
//
// File: cache_test.go
// Description: Unit tests for cache package
package cache

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, 1000, cfg.MaxSize)
	assert.Equal(t, 5*time.Minute, cfg.TTL)
	assert.Equal(t, 2*time.Minute, cfg.CleanupInterval)
}

func TestNew(t *testing.T) {
	tests := []struct {
		name           string
		config         Config
		wantMaxSize    int
		wantTTL        time.Duration
		wantCleanupInt time.Duration
	}{
		{
			name:           "default values for zero config",
			config:         Config{},
			wantMaxSize:    1000,
			wantTTL:        5 * time.Minute,
			wantCleanupInt: 2*time.Minute + 30*time.Second, // TTL/2
		},
		{
			name: "custom config",
			config: Config{
				MaxSize:         500,
				TTL:             10 * time.Minute,
				CleanupInterval: 3 * time.Minute,
			},
			wantMaxSize:    500,
			wantTTL:        10 * time.Minute,
			wantCleanupInt: 3 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New[string](tt.config)
			defer c.Stop()

			assert.NotNil(t, c)
			assert.Equal(t, tt.wantMaxSize, c.config.MaxSize)
			assert.Equal(t, tt.wantTTL, c.config.TTL)
		})
	}
}

func TestCache_SetAndGet(t *testing.T) {
	c := New[string](Config{
		MaxSize: 100,
		TTL:     1 * time.Minute,
	})
	defer c.Stop()

	// Test Set and Get
	c.Set("key1", "value1")
	val, found := c.Get("key1")

	assert.True(t, found)
	assert.Equal(t, "value1", val)
}

func TestCache_GetNotFound(t *testing.T) {
	c := New[string](Config{
		MaxSize: 100,
		TTL:     1 * time.Minute,
	})
	defer c.Stop()

	val, found := c.Get("nonexistent")

	assert.False(t, found)
	assert.Empty(t, val)
}

func TestCache_SetWithTTL(t *testing.T) {
	c := New[string](Config{
		MaxSize:         100,
		TTL:             1 * time.Hour,
		CleanupInterval: 100 * time.Millisecond,
	})
	defer c.Stop()

	// Set with short TTL
	c.SetWithTTL("shortlived", "value", 50*time.Millisecond)

	// Should be found immediately
	val, found := c.Get("shortlived")
	assert.True(t, found)
	assert.Equal(t, "value", val)

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Should be expired
	val, found = c.Get("shortlived")
	assert.False(t, found)
	assert.Empty(t, val)
}

func TestCache_Delete(t *testing.T) {
	c := New[string](Config{
		MaxSize: 100,
		TTL:     1 * time.Minute,
	})
	defer c.Stop()

	c.Set("key1", "value1")
	c.Delete("key1")

	val, found := c.Get("key1")
	assert.False(t, found)
	assert.Empty(t, val)
}

func TestCache_DeletePrefix(t *testing.T) {
	c := New[string](Config{
		MaxSize: 100,
		TTL:     1 * time.Minute,
	})
	defer c.Stop()

	c.Set("user:1", "user1")
	c.Set("user:2", "user2")
	c.Set("role:1", "role1")

	c.DeletePrefix("user:")

	_, found := c.Get("user:1")
	assert.False(t, found)
	_, found = c.Get("user:2")
	assert.False(t, found)

	// role:1 should still exist
	val, found := c.Get("role:1")
	assert.True(t, found)
	assert.Equal(t, "role1", val)
}

func TestCache_Clear(t *testing.T) {
	c := New[string](Config{
		MaxSize: 100,
		TTL:     1 * time.Minute,
	})
	defer c.Stop()

	c.Set("key1", "value1")
	c.Set("key2", "value2")
	c.Clear()

	assert.Equal(t, 0, c.Size())
}

func TestCache_Size(t *testing.T) {
	c := New[string](Config{
		MaxSize: 100,
		TTL:     1 * time.Minute,
	})
	defer c.Stop()

	assert.Equal(t, 0, c.Size())

	c.Set("key1", "value1")
	assert.Equal(t, 1, c.Size())

	c.Set("key2", "value2")
	assert.Equal(t, 2, c.Size())

	c.Delete("key1")
	assert.Equal(t, 1, c.Size())
}

func TestCache_Stats(t *testing.T) {
	c := New[string](Config{
		MaxSize: 100,
		TTL:     1 * time.Minute,
	})
	defer c.Stop()

	c.Set("key1", "value1")

	// Hit
	c.Get("key1")
	c.Get("key1")

	// Miss
	c.Get("nonexistent")

	stats := c.Stats()

	assert.Equal(t, 1, stats.Size)
	assert.Equal(t, int64(2), stats.Hits)
	assert.Equal(t, int64(1), stats.Misses)
	assert.InDelta(t, 0.666, stats.Ratio, 0.01)
}

func TestCache_MaxSize_Eviction(t *testing.T) {
	c := New[int](Config{
		MaxSize: 3,
		TTL:     1 * time.Minute,
	})
	defer c.Stop()

	c.Set("key1", 1)
	c.Set("key2", 2)
	c.Set("key3", 3)

	assert.Equal(t, 3, c.Size())

	// Adding one more should evict the oldest
	c.Set("key4", 4)

	assert.Equal(t, 3, c.Size())
}

func TestCache_GetOrSet(t *testing.T) {
	c := New[string](Config{
		MaxSize: 100,
		TTL:     1 * time.Minute,
	})
	defer c.Stop()

	callCount := 0

	// First call - should call function
	val, err := c.GetOrSet("key1", func() (string, error) {
		callCount++
		return "value1", nil
	})

	require.NoError(t, err)
	assert.Equal(t, "value1", val)
	assert.Equal(t, 1, callCount)

	// Second call - should return cached value
	val, err = c.GetOrSet("key1", func() (string, error) {
		callCount++
		return "value2", nil
	})

	require.NoError(t, err)
	assert.Equal(t, "value1", val) // Still the cached value
	assert.Equal(t, 1, callCount)  // Function not called again
}

func TestCache_GetOrSet_Error(t *testing.T) {
	c := New[string](Config{
		MaxSize: 100,
		TTL:     1 * time.Minute,
	})
	defer c.Stop()

	expectedErr := errors.New("fetch error")

	val, err := c.GetOrSet("key1", func() (string, error) {
		return "", expectedErr
	})

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, val)

	// Should not be cached
	_, found := c.Get("key1")
	assert.False(t, found)
}

func TestCache_ConcurrentAccess(t *testing.T) {
	c := New[int](Config{
		MaxSize: 1000,
		TTL:     1 * time.Minute,
	})
	defer c.Stop()

	var wg sync.WaitGroup
	numGoroutines := 100
	numOperations := 100

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := "key" + string(rune('0'+id))
				c.Set(key, id*numOperations+j)
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := "key" + string(rune('0'+id))
				c.Get(key)
			}
		}(i)
	}

	wg.Wait()

	// Should not panic or race
	assert.True(t, c.Size() > 0)
}

func TestCache_Stop(t *testing.T) {
	c := New[string](Config{
		MaxSize:         100,
		TTL:             1 * time.Minute,
		CleanupInterval: 10 * time.Millisecond,
	})

	c.Set("key1", "value1")

	// Stop should be idempotent
	c.Stop()
	c.Stop()

	// Cache should still work after stop (just no cleanup)
	val, found := c.Get("key1")
	assert.True(t, found)
	assert.Equal(t, "value1", val)
}

func TestCache_TypedCache(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	c := New[User](Config{
		MaxSize: 100,
		TTL:     1 * time.Minute,
	})
	defer c.Stop()

	user := User{ID: 1, Name: "John"}
	c.Set("user:1", user)

	retrieved, found := c.Get("user:1")
	assert.True(t, found)
	assert.Equal(t, user.ID, retrieved.ID)
	assert.Equal(t, user.Name, retrieved.Name)
}
