// Package cache is a registry of named TTL-keyed caches. A Container is
// constructed once in main() and injected into any service that needs to
// cache byte payloads with an expiration.
package cache

import (
	"sync"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

// Container is a registry of named caches. Use New() to construct one in
// main() and pass it to consumers as a constructor dependency.
type Container struct {
	mu     sync.Mutex
	caches map[string]*Instance
}

// Instance is a single TTL-keyed cache.
type Instance struct {
	cache *ttlcache.Cache[string, []byte]
}

// New returns an empty Container.
func New() *Container {
	return &Container{caches: make(map[string]*Instance)}
}

// GetOrCreate returns the named cache, creating it with the given TTL on
// first use. Concurrent calls are safe.
func (c *Container) GetOrCreate(name string, ttl time.Duration) *Instance {
	c.mu.Lock()
	defer c.mu.Unlock()
	if existing, ok := c.caches[name]; ok {
		return existing
	}
	cache := ttlcache.New[string, []byte](
		ttlcache.WithTTL[string, []byte](ttl),
		ttlcache.WithDisableTouchOnHit[string, []byte](),
	)
	inst := &Instance{cache: cache}
	c.caches[name] = inst
	go cache.Start()
	return inst
}

// Stop halts every cache's TTL-eviction goroutine. Safe to call once at
// shutdown.
func (c *Container) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, inst := range c.caches {
		inst.cache.Stop()
	}
}

// Get fetches the value for the given key. Returns nil if the key is
// absent or has expired.
func (i *Instance) Get(key string) []byte {
	v := i.cache.Get(key, ttlcache.WithDisableTouchOnHit[string, []byte]())
	if v == nil || v.IsExpired() {
		return nil
	}
	return v.Value()
}

// Set stores the value for the given key under the cache's configured TTL.
func (i *Instance) Set(key string, value []byte) {
	i.cache.Set(key, value, 0)
}
