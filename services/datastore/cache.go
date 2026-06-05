package datastore

import (
	"errors"
)

// GetCachedValue will return a value for key from the cache.
func (ds *Datastore) GetCachedValue(key string) ([]byte, error) {
	ds.cacheLock.Lock()
	defer ds.cacheLock.Unlock()

	// Check for a cached value
	if val, ok := ds.cache[key]; ok {
		return val, nil
	}

	return nil, errors.New(key + " not found in cache")
}

// SetCachedValue will set a value for key in the cache.
func (ds *Datastore) SetCachedValue(key string, b []byte) {
	ds.cacheLock.Lock()
	defer ds.cacheLock.Unlock()

	ds.cache[key] = b
}

// DeleteCachedValue removes key from the cache.
func (ds *Datastore) DeleteCachedValue(key string) {
	ds.cacheLock.Lock()
	defer ds.cacheLock.Unlock()

	delete(ds.cache, key)
}
