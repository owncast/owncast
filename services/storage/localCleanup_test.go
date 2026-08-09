package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/owncast/owncast/config"
)

// fakeProtector reports a fixed set of protected segment filenames.
type fakeProtector struct {
	protected map[string]bool
	err       error
}

func (f fakeProtector) ProtectedSegmentFilenames() (map[string]bool, error) {
	return f.protected, f.err
}

// writeSegments creates numbered segments in a variant directory, oldest
// first, so cleanup's newest-first ordering is deterministic.
func writeSegments(t *testing.T, variantDir string, names []string) {
	t.Helper()

	if err := os.MkdirAll(variantDir, 0o750); err != nil {
		t.Fatalf("unable to create variant directory: %v", err)
	}

	modTime := time.Now().Add(-time.Duration(len(names)) * time.Minute)
	for _, name := range names {
		path := filepath.Join(variantDir, name)
		if err := os.WriteFile(path, []byte("segment"), 0o600); err != nil {
			t.Fatalf("unable to write segment: %v", err)
		}
		if err := os.Chtimes(path, modTime, modTime); err != nil {
			t.Fatalf("unable to set segment time: %v", err)
		}
		modTime = modTime.Add(time.Minute)
	}
}

// withHLSStorage points config.HLSStoragePath at a temporary directory for the
// duration of a test.
func withHLSStorage(t *testing.T) string {
	t.Helper()

	original := config.HLSStoragePath
	dir := t.TempDir()
	config.HLSStoragePath = dir
	t.Cleanup(func() { config.HLSStoragePath = original })

	return dir
}

func remaining(t *testing.T, variantDir string) []string {
	t.Helper()

	entries, err := os.ReadDir(variantDir)
	if err != nil {
		t.Fatalf("unable to read variant directory: %v", err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names
}

// TestLocalCleanupKeepsProtectedSegments is the core of clip retention: a
// segment a clip or replay references survives cleanup even though it is old
// enough to be pruned.
func TestLocalCleanupKeepsProtectedSegments(t *testing.T) {
	dir := withHLSStorage(t)
	variantDir := filepath.Join(dir, "0")

	writeSegments(t, variantDir, []string{
		"stream-1.ts", "stream-2.ts", "stream-3.ts", "stream-4.ts", "stream-5.ts",
	})

	// Keep only the 2 newest, but protect one of the older ones.
	protector := fakeProtector{protected: map[string]bool{"stream-2.ts": true}}
	if err := localCleanup(2, protector); err != nil {
		t.Fatalf("cleanup: %v", err)
	}

	left := remaining(t, variantDir)
	want := map[string]bool{"stream-4.ts": true, "stream-5.ts": true, "stream-2.ts": true}

	if len(left) != len(want) {
		t.Fatalf("expected %d segments to remain, got %v", len(want), left)
	}
	for _, name := range left {
		if !want[name] {
			t.Errorf("unexpected segment %q survived cleanup", name)
		}
	}
}

// TestLocalCleanupWithoutProtectorPrunes verifies the ordinary live-window
// behavior is unchanged when nothing is protected.
func TestLocalCleanupWithoutProtectorPrunes(t *testing.T) {
	dir := withHLSStorage(t)
	variantDir := filepath.Join(dir, "0")

	writeSegments(t, variantDir, []string{"a.ts", "b.ts", "c.ts", "d.ts"})

	if err := localCleanup(2, nil); err != nil {
		t.Fatalf("cleanup: %v", err)
	}

	if left := remaining(t, variantDir); len(left) != 2 {
		t.Errorf("expected 2 segments to remain, got %v", left)
	}
}

// TestLocalCleanupProtectorFailureIsNotFatal verifies a ledger that cannot
// answer degrades to ordinary cleanup rather than blocking it.
func TestLocalCleanupProtectorFailureIsNotFatal(t *testing.T) {
	dir := withHLSStorage(t)
	variantDir := filepath.Join(dir, "0")

	writeSegments(t, variantDir, []string{"a.ts", "b.ts", "c.ts", "d.ts"})

	protector := fakeProtector{err: os.ErrClosed}
	if err := localCleanup(2, protector); err != nil {
		t.Fatalf("cleanup: %v", err)
	}

	if left := remaining(t, variantDir); len(left) != 2 {
		t.Errorf("expected 2 segments to remain, got %v", left)
	}
}
