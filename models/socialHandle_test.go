package models

import (
	"encoding/xml"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestGetAllSocialHandlesIncludesRSS(t *testing.T) {
	handles := GetAllSocialHandles()

	handle, ok := handles["rss"]
	if !ok {
		t.Fatal("expected \"rss\" key to be present in GetAllSocialHandles()")
	}

	if handle.Platform != "RSS" {
		t.Errorf("expected Platform %q, got %q", "RSS", handle.Platform)
	}

	if handle.Icon != "/img/platformlogos/rss.svg" {
		t.Errorf("expected Icon %q, got %q", "/img/platformlogos/rss.svg", handle.Icon)
	}
}

func TestGetAllSocialHandlesIncludesTumblr(t *testing.T) {
	handles := GetAllSocialHandles()

	handle, ok := handles["tumblr"]
	if !ok {
		t.Fatal("expected \"tumblr\" key to be present in GetAllSocialHandles()")
	}

	if handle.Platform != "Tumblr" {
		t.Errorf("expected Platform %q, got %q", "Tumblr", handle.Platform)
	}

	if handle.Icon != "/img/platformlogos/tumblr.svg" {
		t.Errorf("expected Icon %q, got %q", "/img/platformlogos/tumblr.svg", handle.Icon)
	}
}

func TestSocialPlatformSVGsAreValidXML(t *testing.T) {
	files := []string{"rss.svg", "tumblr.svg"}

	for _, name := range files {
		name := name
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
