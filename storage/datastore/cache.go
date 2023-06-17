package datastore

import (
	"errors"
	"sync"

	log "github.com/sirupsen/logrus"
)

var _cacheLock = sync.Mutex{}

// GetCachedValue will return a value for key from the cache.
func (ds *Datastore) GetCachedValue(key string) ([]byte, error) {
	_cacheLock.Lock()
	defer _cacheLock.Unlock()

	// Check for a cached value
	if val, ok := ds.cache[key]; ok {
		return val, nil
	}

	return nil, errors.New(key + " not found in cache")
}

// SetCachedValue will set a value for key in the cache.
func (ds *Datastore) SetCachedValue(key string, b []byte) {
	_cacheLock.Lock()
	defer _cacheLock.Unlock()

	ds.cache[key] = b
}

func (ds *Datastore) warmCache() {
	log.Traceln("Warming config value cache")

	res, err := ds.DB.Query("SELECT key, value FROM datastore")
	if err != nil || res.Err() != nil {
		log.Errorln("error warming config cache", err, res.Err())
	}
	defer res.Close()

	for res.Next() {
		var rowKey string
		var rowValue []byte
		if err := res.Scan(&rowKey, &rowValue); err != nil {
			log.Errorln("error pre-caching config row", err)
		}
		ds.cache[rowKey] = rowValue
	}
}
