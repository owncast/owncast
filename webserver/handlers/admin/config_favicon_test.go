package admin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/datastore"
)

// testAdmin is the *Admin used by tests in this package. Only the
// configRepository handle is needed by favicon tests; expand if other
// handler tests land that need additional deps.
var (
	testAdmin     *Admin
	testDatastore *datastore.Datastore
)

func TestMain(m *testing.M) {
	dbFile, err := os.CreateTemp(os.TempDir(), "owncast-test-db.db")
	if err != nil {
		panic(err)
	}
	dbFile.Close()

	ds, err := datastore.SetupPersistence(dbFile.Name(), os.TempDir())
	if err != nil {
		panic(err)
	}
	testDatastore = ds

	// Ensure data directory exists for file operations.
	if err := os.MkdirAll("data", 0o755); err != nil {
		panic(err)
	}

	testAdmin = &Admin{configRepository: configrepository.New(testDatastore)}

	code := m.Run()
	os.Remove(dbFile.Name())
	os.Exit(code)
}

// Minimal valid 1x1 PNG.
var minimalPNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49,
	0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x02,
	0x00, 0x00, 0x00, 0x90, 0x77, 0x53, 0xde, 0x00, 0x00, 0x00, 0x0c, 0x49, 0x44,
	0x41, 0x54, 0x08, 0xd7, 0x63, 0xf8, 0xff, 0xff, 0x3f, 0x00, 0x05, 0xfe, 0x02,
	0xfe, 0xdc, 0xcc, 0x59, 0xe7, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44,
	0xae, 0x42, 0x60, 0x82,
}

// Minimal valid ICO bytes.
var minimalICO = []byte{
	0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x18,
	0x00, 0x30, 0x00, 0x00, 0x00, 0x16, 0x00, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00,
	0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x18, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff,
	0x00, 0x00, 0x00, 0x00, 0x00,
}

func makeBase64DataURL(contentType string, imgData []byte) string {
	return fmt.Sprintf("data:%s;base64,%s", contentType, base64.StdEncoding.EncodeToString(imgData))
}

func makeConfigValueBody(value string) string {
	return fmt.Sprintf(`{"value": %q}`, value)
}

func parseResponse(t *testing.T, w *httptest.ResponseRecorder) models.BaseAPIResponse {
	t.Helper()
	var resp models.BaseAPIResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return resp
}

func cleanupFaviconFiles(t *testing.T) {
	t.Helper()
	for _, ext := range []string{".png", ".ico"} {
		os.Remove(filepath.Join("data", "favicon"+ext))
	}
	configRepository := configrepository.New(testDatastore)
	configRepository.SetFaviconPath("")
}

func TestSetFaviconPNG(t *testing.T) {
	defer cleanupFaviconFiles(t)

	body := makeConfigValueBody(makeBase64DataURL("image/png", minimalPNG))
	req := httptest.NewRequest(http.MethodPost, "/api/admin/config/favicon", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	testAdmin.SetFavicon(w, req)

	resp := parseResponse(t, w)
	if !resp.Success {
		t.Fatalf("expected success, got error: %s", resp.Message)
	}
	if resp.Message != "favicon updated" {
		t.Errorf("expected message 'favicon updated', got %q", resp.Message)
	}

	// Verify file was written.
	if _, err := os.Stat(filepath.Join("data", "favicon.png")); os.IsNotExist(err) {
		t.Error("expected favicon.png to be written to disk")
	}

	// Verify config was updated.
	configRepository := configrepository.New(testDatastore)
	if path := configRepository.GetFaviconPath(); path != "favicon.png" {
		t.Errorf("expected favicon path 'favicon.png', got %q", path)
	}
}

func TestSetFaviconICO(t *testing.T) {
	defer cleanupFaviconFiles(t)

	body := makeConfigValueBody(makeBase64DataURL("image/x-icon", minimalICO))
	req := httptest.NewRequest(http.MethodPost, "/api/admin/config/favicon", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	testAdmin.SetFavicon(w, req)

	resp := parseResponse(t, w)
	if !resp.Success {
		t.Fatalf("expected success, got error: %s", resp.Message)
	}

	// Verify file was written.
	if _, err := os.Stat(filepath.Join("data", "favicon.ico")); os.IsNotExist(err) {
		t.Error("expected favicon.ico to be written to disk")
	}

	// Verify config was updated.
	configRepository := configrepository.New(testDatastore)
	if path := configRepository.GetFaviconPath(); path != "favicon.ico" {
		t.Errorf("expected favicon path 'favicon.ico', got %q", path)
	}
}

func TestSetFaviconRejectsJPEG(t *testing.T) {
	defer cleanupFaviconFiles(t)

	body := makeConfigValueBody(makeBase64DataURL("image/jpeg", minimalPNG))
	req := httptest.NewRequest(http.MethodPost, "/api/admin/config/favicon", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	testAdmin.SetFavicon(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	resp := parseResponse(t, w)
	if resp.Success {
		t.Error("expected failure for JPEG upload")
	}
	if resp.Message != "favicon must be PNG or ICO format" {
		t.Errorf("unexpected error message: %s", resp.Message)
	}
}

func TestSetFaviconRejectsOversized(t *testing.T) {
	defer cleanupFaviconFiles(t)

	// Create a PNG data URL with >200KB of decoded data but under the
	// MaxBytesReader limit so the size check in the handler is exercised.
	oversized := make([]byte, 210*1024)
	body := makeConfigValueBody(makeBase64DataURL("image/png", oversized))
	req := httptest.NewRequest(http.MethodPost, "/api/admin/config/favicon", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	testAdmin.SetFavicon(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	resp := parseResponse(t, w)
	if resp.Success {
		t.Error("expected failure for oversized file")
	}
	if resp.Message != "file too large, max 200KB" {
		t.Errorf("unexpected error message: %s", resp.Message)
	}
}

func TestSetFaviconRejectsEmptyBody(t *testing.T) {
	defer cleanupFaviconFiles(t)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/config/favicon", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	testAdmin.SetFavicon(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestSetFaviconRejectsGETMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/admin/config/favicon", nil)
	w := httptest.NewRecorder()

	testAdmin.SetFavicon(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for GET request, got %d", w.Code)
	}
}

func TestResetFavicon(t *testing.T) {
	// First set a favicon so there's something to reset.
	if err := os.MkdirAll("data", 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join("data", "favicon.png"), minimalPNG, 0o600); err != nil {
		t.Fatal(err)
	}

	configRepository := configrepository.New(testDatastore)
	if err := configRepository.SetFaviconPath("favicon.png"); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/admin/config/favicon", nil)
	w := httptest.NewRecorder()

	testAdmin.ResetFavicon(w, req)

	resp := parseResponse(t, w)
	if !resp.Success {
		t.Fatalf("expected success, got error: %s", resp.Message)
	}
	if resp.Message != "favicon reset to default" {
		t.Errorf("expected message 'favicon reset to default', got %q", resp.Message)
	}

	// Verify config was cleared.
	if path := configRepository.GetFaviconPath(); path != "" {
		t.Errorf("expected empty favicon path after reset, got %q", path)
	}

	// Verify file was removed.
	if _, err := os.Stat(filepath.Join("data", "favicon.png")); !os.IsNotExist(err) {
		t.Error("expected favicon.png to be deleted after reset")
	}
}

func TestResetFaviconWhenNoneSet(t *testing.T) {
	// Reset when no custom favicon is set should still succeed.
	configRepository := configrepository.New(testDatastore)
	configRepository.SetFaviconPath("")

	req := httptest.NewRequest(http.MethodDelete, "/api/admin/config/favicon", nil)
	w := httptest.NewRecorder()

	testAdmin.ResetFavicon(w, req)

	resp := parseResponse(t, w)
	if !resp.Success {
		t.Fatalf("expected success, got error: %s", resp.Message)
	}
}
