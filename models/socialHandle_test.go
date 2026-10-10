package models

import (
	"encoding/xml"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestGetAllSocialHandlesIncludesNewPlatforms(t *testing.T) {
	handles := GetAllSocialHandles()

	tests := []struct {
		key      string
		platform string
		icon     string
	}{
		{key: "rss", platform: "RSS", icon: "/img/platformlogos/rss.svg"},
		{key: "tumblr", platform: "Tumblr", icon: "/img/platformlogos/tumblr.svg"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			handle, ok := handles[tt.key]
			if !ok {
				t.Fatalf("expected %q key to be present in GetAllSocialHandles()", tt.key)
			}

			if handle.Platform != tt.platform {
				t.Errorf("expected Platform %q, got %q", tt.platform, handle.Platform)
			}

			if handle.Icon != tt.icon {
				t.Errorf("expected Icon %q, got %q", tt.icon, handle.Icon)
			}
		})
	}
}

func TestSocialPlatformSVGsAreValidXML(t *testing.T) {
	files := []string{"rss.svg", "tumblr.svg"}

	for _, name := range files {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join("..", "web", "public", "img", "platformlogos", name)

			f, err := os.Open(path)
			if err != nil {
				t.Fatalf("expected %s to exist: %v", path, err)
			}
			defer f.Close()

			decoder := xml.NewDecoder(f)
			for {
				_, err := decoder.Token()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					t.Fatalf("expected %s to be well-formed XML: %v", path, err)
				}
			}
		})
	}
}
