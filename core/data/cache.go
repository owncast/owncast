package data

import "errors"

func (ds *Datastore) GetCachedValue(key string) (ConfigEntry, error) {
	// Check for a cached value
	if val, ok := ds.cache[key]; ok {
		return val, nil
	}

	return ConfigEntry{}, errors.New(key + " not found in cache")
}

func (ds *Datastore) SetCachedValue(key string, e ConfigEntry) {
	ds.cache[key] = e
}
