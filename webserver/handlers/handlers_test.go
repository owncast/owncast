package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/persistence/configrepository"
)

func TestMain(m *testing.M) {
	dbFile, err := os.CreateTemp(os.TempDir(), "owncast-test-db.db")
	if err != nil {
		panic(err)
	}
	dbFile.Close()

	if err := data.SetupPersistence(dbFile.Name()); err != nil {
		panic(err)
	}

	code := m.Run()
	os.Remove(dbFile.Name())
	os.Exit(code)
}

// Minimal 1x1 PNG.
var minimalPNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49,
	0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x02,
	0x00, 0x00, 0x00, 0x90, 0x77, 0x53, 0xde, 0x00, 0x00, 0x00, 0x0c, 0x49, 0x44,
	0x41, 0x54, 0x08, 0xd7, 0x63, 0xf8, 0xff, 0xff, 0x3f, 0x00, 0x05, 0xfe, 0x02,
	0xfe, 0xdc, 0xcc, 0x59, 0xe7, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44,
	0xae, 0x42, 0x60, 0x82,
}

func makeBase64DataURL(contentType string, data []byte) string {
	return fmt.Sprintf("data:%s;base64,%s", contentType, base64.StdEncoding.EncodeToString(data))
}

func TestGetFaviconDefault(t *testing.T) {
	// With no custom favicon set, should return the default.
	configRepository := configrepository.Get()
	_ = configRepository.SetFaviconPath("")

	req := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
	w := httptest.NewRecorder()

	GetFavicon(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "image/png" {
		t.Errorf("expected Content-Type image/png for default favicon, got %s", ct)
	}
}

func TestGetFaviconCustomPNG(t *testing.T) {
	// Set up a custom PNG favicon on disk.
	if err := os.MkdirAll("data", 0o755); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filepath.Join("data", "favicon.png"))

	if err := os.WriteFile(filepath.Join("data", "favicon.png"), minimalPNG, 0o600); err != nil {
		t.Fatal(err)
	}

	configRepository := configrepository.Get()
	if err := configRepository.SetFaviconPath("favicon.png"); err != nil {
		t.Fatal(err)
	}
	defer configRepository.SetFaviconPath("")

	req := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
	w := httptest.NewRecorder()

	GetFavicon(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "image/png" {
		t.Errorf("expected Content-Type image/png, got %s", ct)
	}
}

func TestGetFaviconCustomICO(t *testing.T) {
	// Set up a custom ICO favicon on disk.
	if err := os.MkdirAll("data", 0o755); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filepath.Join("data", "favicon.ico"))

	icoData := []byte{0x00, 0x00, 0x01, 0x00} // minimal ICO header bytes
	if err := os.WriteFile(filepath.Join("data", "favicon.ico"), icoData, 0o600); err != nil {
		t.Fatal(err)
	}

	configRepository := configrepository.Get()
	if err := configRepository.SetFaviconPath("favicon.ico"); err != nil {
		t.Fatal(err)
	}
	defer configRepository.SetFaviconPath("")

	req := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
	w := httptest.NewRecorder()

	GetFavicon(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "image/x-icon" {
		t.Errorf("expected Content-Type image/x-icon, got %s", ct)
	}
}

func TestGetFaviconMissingFileFallsBackToDefault(t *testing.T) {
	// Point config at a file that doesn't exist; should fall back to default.
	configRepository := configrepository.Get()
	if err := configRepository.SetFaviconPath("favicon-nonexistent.png"); err != nil {
		t.Fatal(err)
	}
	defer configRepository.SetFaviconPath("")

	req := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
	w := httptest.NewRecorder()

	GetFavicon(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "image/png" {
		t.Errorf("expected fallback Content-Type image/png, got %s", ct)
	}
}
