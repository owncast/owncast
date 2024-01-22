package cache

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

// CacheContainer is a container for all caches.
type CacheContainer struct {
	caches map[string]*CacheInstance
}

// CacheInstance is a single cache instance.
type CacheInstance struct {
	cache *ttlcache.Cache[string, []byte]
}

// This is the global singleton instance. (To be removed after refactor).
var _instance *CacheContainer

// NewCache creates a new cache instance.
func NewGlobalCache() *CacheContainer {
	_instance = &CacheContainer{
		caches: make(map[string]*CacheInstance),
	}

	return _instance
}

// GetCache returns the cache instance.
func GetGlobalCache() *CacheContainer {
	if _instance != nil {
		return _instance
	}
	return NewGlobalCache()
}

// GetOrCreateCache returns the cache instance or creates a new one.
func (c *CacheContainer) GetOrCreateCache(name string, expiration time.Duration) *CacheInstance {
	if _, ok := c.caches[name]; !ok {
		c.CreateCache(name, expiration)
	}
	return c.caches[name]
}

// CreateCache creates a new cache instance.
func (c *CacheContainer) CreateCache(name string, expiration time.Duration) *CacheInstance {
	cache := ttlcache.New[string, []byte](
		ttlcache.WithTTL[string, []byte](expiration),
		ttlcache.WithDisableTouchOnHit[string, []byte](),
	)
	ci := &CacheInstance{cache: cache}
	c.caches[name] = ci
	go cache.Start()
	return ci
}

// GetCache returns the cache instance.
func (c *CacheContainer) GetCache(name string) *CacheInstance {
	return c.caches[name]
}

// GetValueForKey returns the value for the given key.
func (ci *CacheInstance) GetValueForKey(key string) []byte {
	value := ci.cache.Get(key, ttlcache.WithDisableTouchOnHit[string, []byte]())
	if value == nil {
		return nil
	}

	if value.IsExpired() {
		return nil
	}

	val := value.Value()
	return val
}

// Set sets the value for the given key..
func (ci *CacheInstance) Set(key string, value []byte) {
	ci.cache.Set(key, value, 0)
}
