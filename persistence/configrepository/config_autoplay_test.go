package configrepository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/owncast/owncast/services/datastore"
)

func newAutoplayTestRepo(t *testing.T) ConfigRepository {
	t.Helper()
	dir := t.TempDir()
	ds, err := datastore.SetupPersistence(filepath.Join(dir, "autoplay-test.db"), os.TempDir())
	if err != nil {
		t.Fatalf("failed to set up datastore: %v", err)
	}
	return New(ds)
}

func TestGetAutoplayDefaultsToOff(t *testing.T) {
	repo := newAutoplayTestRepo(t)
	if got := repo.GetAutoplay(); got != "off" {
		t.Errorf("expected default autoplay 'off' on an unset key, got %q", got)
	}
}

func TestSetAndGetAutoplayRoundTrip(t *testing.T) {
	repo := newAutoplayTestRepo(t)
	for _, value := range []string{"always", "sound-only", "off"} {
		if err := repo.SetAutoplay(value); err != nil {
			t.Fatalf("SetAutoplay(%q) returned error: %v", value, err)
		}
		if got := repo.GetAutoplay(); got != value {
			t.Errorf("after SetAutoplay(%q), GetAutoplay() = %q", value, got)
		}
	}
}
