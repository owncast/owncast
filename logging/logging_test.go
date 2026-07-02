package logging

import (
	"os"
	"path/filepath"
	"testing"
)

// A previous run that dies mid-rotation can leave rotatelogs _lock and
// _symlink staging files behind, which block every future rotation of that
// log. Verify startup cleanup removes them without touching real log files.
func TestRemoveStaleRotationFiles(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "transcoder.log")

	staleFiles := []string{
		"transcoder.log.202607020000_lock",
		"transcoder.log.202607020000.1_lock",
		"transcoder.log.202607020000_symlink",
	}
	keepFiles := []string{
		"transcoder.log.202607020000",
		"transcoder.log.202607020000.1",
		"owncast.log.202607020000_lock",
	}

	for _, name := range append(append([]string{}, staleFiles...), keepFiles...) {
		if err := os.WriteFile(filepath.Join(dir, name), nil, 0o600); err != nil {
			t.Fatalf("unable to create test file %s: %v", name, err)
		}
	}

	removeStaleRotationFiles(logPath)

	for _, name := range staleFiles {
		if _, err := os.Stat(filepath.Join(dir, name)); !os.IsNotExist(err) {
			t.Errorf("stale rotation file %s should have been removed", name)
		}
	}

	for _, name := range keepFiles {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Errorf("file %s should not have been removed: %v", name, err)
		}
	}
}
