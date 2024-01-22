package cache

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	expiration := 5 * time.Second
	globalCache := GetGlobalCache()
	assert.NotNil(t, globalCache, "NewGlobalCache should return a non-nil instance")
	assert.Equal(t, globalCache, GetGlobalCache(), "GetGlobalCache should return the created instance")

	cacheName := "testCache"
	globalCache.CreateCache(cacheName, expiration)

	createdCache := globalCache.GetCache(cacheName)
	assert.NotNil(t, createdCache, "GetCache should return a non-nil cache")

	key := "testKey"
	value := []byte("testValue")
	createdCache.Set(key, value)

	// Wait for cache to expire
	time.Sleep(expiration + 1*time.Second)

	// Verify that the cache has expired
	ci := globalCache.GetCache(cacheName)
	cachedValue := ci.GetValueForKey(key)
	assert.Nil(t, cachedValue, "Cache should not contain the value after expiration")
}

func TestConcurrentAccess(t *testing.T) {
	// Test concurrent access to the cache
	globalCache := NewGlobalCache()
	cacheName := "concurrentCache"
	expiration := 5 * time.Second
	globalCache.CreateCache(cacheName, expiration)

	// Start multiple goroutines to access the cache concurrently
	numGoroutines := 10
	keyPrefix := "key"
	valuePrefix := "value"

	done := make(chan struct{})
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer func() { done <- struct{}{} }()

			cache := globalCache.GetCache(cacheName)
			key := keyPrefix + strconv.Itoa(index)
			value := valuePrefix + strconv.Itoa(index)

			cache.Set(key, []byte(value))

			// Simulate some work
			time.Sleep(100 * time.Millisecond)

			ci := globalCache.GetCache(cacheName)
			cachedValue := string(ci.GetValueForKey(key))
			assert.Equal(t, value, cachedValue, "Cached value should match the set value")
		}(i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}
