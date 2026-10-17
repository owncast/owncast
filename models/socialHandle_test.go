package models

import (
	"os"
	"path/filepath"
	"testing"
)

// Every platform must point at an icon that is shipped in web/public, so a
// typo in a path is caught here instead of as a missing image on the page.
func TestSocialHandleIconsExist(t *testing.T) {
	for key, handle := range GetAllSocialHandles() {
		t.Run(key, func(t *testing.T) {
			path := filepath.Join("..", "web", "public", filepath.FromSlash(handle.Icon))
			if _, err := os.Stat(path); err != nil {
				t.Fatalf("icon for %q: %v", key, err)
			}
		})
	}
}
