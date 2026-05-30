package kv

import "sync"

// NewMemory returns an in-memory Store. Used by the test runner so each
// scenario gets fresh state without touching disk.
func NewMemory() Store {
	return &memoryStore{namespaces: map[string]*memoryNamespace{}}
}

type memoryStore struct {
	mu         sync.Mutex
	namespaces map[string]*memoryNamespace
}

func (s *memoryStore) Close() error { return nil }

func (s *memoryStore) Namespace(plugin string) Namespace {
	s.mu.Lock()
	defer s.mu.Unlock()
	ns, ok := s.namespaces[plugin]
	if !ok {
		ns = &memoryNamespace{data: map[string][]byte{}}
		s.namespaces[plugin] = ns
	}
	return ns
}

type memoryNamespace struct {
	mu   sync.Mutex
	data map[string][]byte
}

func (n *memoryNamespace) Get(key string) ([]byte, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	v, ok := n.data[key]
	if !ok {
		return nil, nil
	}
	out := make([]byte, len(v))
	copy(out, v)
	return out, nil
}

func (n *memoryNamespace) Set(key string, value []byte) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	cp := make([]byte, len(value))
	copy(cp, value)
	n.data[key] = cp
	return nil
}

func (n *memoryNamespace) Delete(key string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	delete(n.data, key)
	return nil
}

// Snapshot returns a copy of all data in the namespace. Useful for tests
// that want to assert on final KV state.
func (n *memoryNamespace) Snapshot() map[string][]byte {
	n.mu.Lock()
	defer n.mu.Unlock()
	out := make(map[string][]byte, len(n.data))
	for k, v := range n.data {
		cp := make([]byte, len(v))
		copy(cp, v)
		out[k] = cp
	}
	return out
}
