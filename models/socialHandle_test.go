package models

import (
	"encoding/xml"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// Every platform must point at an icon that is shipped in web/public, and
// SVG icons must be well-formed, so a typo in a path or a broken file is
// caught here instead of as a missing image on the page.
func TestSocialHandleIconsExist(t *testing.T) {
	for key, handle := range GetAllSocialHandles() {
		t.Run(key, func(t *testing.T) {
			path := filepath.Join("..", "web", "public", filepath.FromSlash(handle.Icon))

			f, err := os.Open(path)
			if err != nil {
				t.Fatalf("icon for %q: %v", key, err)
			}
			defer f.Close()

			if filepath.Ext(path) != ".svg" {
				return
			}
			decoder := xml.NewDecoder(f)
			for {
				if _, err := decoder.Token(); err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					t.Fatalf("icon for %q is not well-formed XML: %v", key, err)
				}
			}
		})
	}
}
