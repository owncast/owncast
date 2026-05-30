package kv

import (
	"bytes"
	"sync"
	"testing"
)

// The in-memory implementation must satisfy the Store interface contract.
// Host-provided backends (e.g. Owncast's datastore-backed store) are tested
// in their own package against the same behaviors.
func TestStoreContract(t *testing.T) {
	t.Run("memory", func(t *testing.T) {
		runStoreSuite(t, func() (Store, func()) {
			return NewMemory(), func() {}
		})
	})
}

func runStoreSuite(t *testing.T, openStore func() (Store, func())) {
	t.Helper()

	t.Run("get missing returns nil", func(t *testing.T) {
		s, cleanup := openStore()
		defer cleanup()
		ns := s.Namespace("p1")
		got, err := ns.Get("missing")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != nil {
			t.Errorf("expected nil, got %q", got)
		}
	})

	t.Run("set then get round-trip", func(t *testing.T) {
		s, cleanup := openStore()
		defer cleanup()
		ns := s.Namespace("p1")
		if err := ns.Set("k", []byte("hello")); err != nil {
			t.Fatalf("set: %v", err)
		}
		got, err := ns.Get("k")
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if !bytes.Equal(got, []byte("hello")) {
			t.Errorf("get returned %q, want %q", got, "hello")
		}
	})

	t.Run("overwrite", func(t *testing.T) {
		s, cleanup := openStore()
		defer cleanup()
		ns := s.Namespace("p1")
		ns.Set("k", []byte("first"))
		ns.Set("k", []byte("second"))
		got, _ := ns.Get("k")
		if !bytes.Equal(got, []byte("second")) {
			t.Errorf("overwrite failed: got %q want %q", got, "second")
		}
	})

	t.Run("delete", func(t *testing.T) {
		s, cleanup := openStore()
		defer cleanup()
		ns := s.Namespace("p1")
		ns.Set("k", []byte("v"))
		if err := ns.Delete("k"); err != nil {
			t.Fatalf("delete: %v", err)
		}
		got, _ := ns.Get("k")
		if got != nil {
			t.Errorf("after delete, got %q want nil", got)
		}
	})

	t.Run("delete missing is no-op", func(t *testing.T) {
		s, cleanup := openStore()
		defer cleanup()
		ns := s.Namespace("p1")
		if err := ns.Delete("never-set"); err != nil {
			t.Errorf("delete on missing should be no-op, got: %v", err)
		}
	})

	t.Run("namespaces are isolated", func(t *testing.T) {
		s, cleanup := openStore()
		defer cleanup()
		a := s.Namespace("plugin-a")
		b := s.Namespace("plugin-b")
		a.Set("shared-key", []byte("a-value"))
		b.Set("shared-key", []byte("b-value"))

		gotA, _ := a.Get("shared-key")
		gotB, _ := b.Get("shared-key")
		if !bytes.Equal(gotA, []byte("a-value")) {
			t.Errorf("plugin-a got %q want %q", gotA, "a-value")
		}
		if !bytes.Equal(gotB, []byte("b-value")) {
			t.Errorf("plugin-b got %q want %q", gotB, "b-value")
		}

		// Delete in A must not affect B.
		a.Delete("shared-key")
		gotA2, _ := a.Get("shared-key")
		gotB2, _ := b.Get("shared-key")
		if gotA2 != nil {
			t.Errorf("plugin-a after delete: got %q want nil", gotA2)
		}
		if !bytes.Equal(gotB2, []byte("b-value")) {
			t.Errorf("plugin-b after a's delete: got %q want %q", gotB2, "b-value")
		}
	})

	t.Run("returned bytes are owned by caller", func(t *testing.T) {
		// Mutating the returned slice must not corrupt the store.
		s, cleanup := openStore()
		defer cleanup()
		ns := s.Namespace("p1")
		ns.Set("k", []byte("hello"))
		got, _ := ns.Get("k")
		for i := range got {
			got[i] = 'X'
		}
		again, _ := ns.Get("k")
		if !bytes.Equal(again, []byte("hello")) {
			t.Errorf("mutating returned slice corrupted store: got %q", again)
		}
	})
}

func TestMemory_ConcurrentAccess(t *testing.T) {
	// Smoke test for the mutex — race detector should catch issues.
	s := NewMemory()
	ns := s.Namespace("p1")
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ns.Set("k", []byte("v"))
			ns.Get("k")
		}(i)
	}
	wg.Wait()
}
