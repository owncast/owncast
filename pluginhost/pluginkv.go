package pluginhost

import (
	"database/sql"
	"errors"

	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/services/plugins/kv"
)

// pluginKVKeyPrefix namespaces plugin KV entries within Owncast's shared
// config datastore so they don't collide with server config keys.
const pluginKVKeyPrefix = "plugins.kv."

// datastoreKVStore implements the plugin runtime's kv.Store on top of
// Owncast's native key/value datastore. Each plugin's keys live under
// "plugins.kv.<plugin>." so plugins can neither read nor overwrite each
// other's data.
type datastoreKVStore struct {
	datastore *datastore.Datastore
}

func newDatastoreKVStore(ds *datastore.Datastore) *datastoreKVStore {
	return &datastoreKVStore{datastore: ds}
}

func (s *datastoreKVStore) Namespace(plugin string) kv.Namespace {
	return &datastoreKVNamespace{datastore: s.datastore, prefix: pluginKVKeyPrefix + plugin + "."}
}

// Close is a no-op: the datastore's lifecycle is owned by Owncast, not the
// plugin host.
func (s *datastoreKVStore) Close() error { return nil }

type datastoreKVNamespace struct {
	datastore *datastore.Datastore
	prefix    string
}

func (n *datastoreKVNamespace) Get(key string) ([]byte, error) {
	value, err := n.datastore.GetString(n.prefix + key)
	if err != nil {
		// Only a missing key (or unset on a fresh install) reads as
		// nil, matching the kv.Store contract. Anything else is a
		// real persistence error (decode/type/database) that the
		// plugin needs to see; silently translating to nil would
		// look like data loss to the plugin.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return []byte(value), nil
}

func (n *datastoreKVNamespace) Set(key string, value []byte) error {
	return n.datastore.SetString(n.prefix+key, string(value))
}

func (n *datastoreKVNamespace) Delete(key string) error {
	return n.datastore.Delete(n.prefix + key)
}
