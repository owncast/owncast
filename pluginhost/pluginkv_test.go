package pluginhost

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/owncast/owncast/services/datastore"
)

func newTestDatastore(t *testing.T) *datastore.Datastore {
	t.Helper()
	ds, err := datastore.SetupPersistence(filepath.Join(t.TempDir(), "test.db"), t.TempDir())
	if err != nil {
		t.Fatalf("setup datastore: %v", err)
	}
	return ds
}

func TestDatastoreKVStore(t *testing.T) {
	store := newDatastoreKVStore(newTestDatastore(t))
	ns := store.Namespace("p1")

	t.Run("get missing returns nil", func(t *testing.T) {
		got, err := ns.Get("missing")
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if got != nil {
			t.Errorf("expected nil, got %q", got)
		}
	})

	t.Run("set then get round-trip", func(t *testing.T) {
		if err := ns.Set("k", []byte("hello")); err != nil {
			t.Fatalf("set: %v", err)
		}
		got, err := ns.Get("k")
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if !bytes.Equal(got, []byte("hello")) {
			t.Errorf("got %q want %q", got, "hello")
		}
	})

	t.Run("overwrite", func(t *testing.T) {
		_ = ns.Set("ow", []byte("first"))
		_ = ns.Set("ow", []byte("second"))
		got, _ := ns.Get("ow")
		if !bytes.Equal(got, []byte("second")) {
			t.Errorf("got %q want %q", got, "second")
		}
	})

	t.Run("binary value round-trip", func(t *testing.T) {
		blob := []byte{0x00, 0xff, 0x10, 'a', 0x00, 0x7f}
		if err := ns.Set("bin", blob); err != nil {
			t.Fatalf("set: %v", err)
		}
		got, _ := ns.Get("bin")
		if !bytes.Equal(got, blob) {
			t.Errorf("binary value corrupted: got %v want %v", got, blob)
		}
	})

	t.Run("delete then get returns nil", func(t *testing.T) {
		_ = ns.Set("del", []byte("v"))
		if err := ns.Delete("del"); err != nil {
			t.Fatalf("delete: %v", err)
		}
		got, _ := ns.Get("del")
		if got != nil {
			t.Errorf("after delete got %q want nil", got)
		}
	})

	t.Run("namespaces are isolated", func(t *testing.T) {
		a := store.Namespace("plugin-a")
		b := store.Namespace("plugin-b")
		_ = a.Set("shared", []byte("a-value"))
		_ = b.Set("shared", []byte("b-value"))
		gotA, _ := a.Get("shared")
		gotB, _ := b.Get("shared")
		if !bytes.Equal(gotA, []byte("a-value")) {
			t.Errorf("plugin-a got %q want %q", gotA, "a-value")
		}
		if !bytes.Equal(gotB, []byte("b-value")) {
			t.Errorf("plugin-b got %q want %q", gotB, "b-value")
		}
	})
}
