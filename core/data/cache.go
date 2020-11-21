package data

import "errors"

func (ds *Datastore) GetCachedValue(key string) ([]byte, error) {
	// Check for a cached value
	if val, ok := ds.cache[key]; ok {
		return val, nil
	}

	return nil, errors.New(key + " not found in cache")
}

func (ds *Datastore) SetCachedValue(key string, b []byte) {
	ds.cache[key] = b
}
