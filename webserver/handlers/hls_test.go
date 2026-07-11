package handlers

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/owncast/owncast/config"
)

func TestHandleHLSRequestRejectsPathTraversal(t *testing.T) {
	originalStoragePath := config.HLSStoragePath
	config.HLSStoragePath = filepath.Join(t.TempDir(), "hls")
	t.Cleanup(func() {
		config.HLSStoragePath = originalStoragePath
	})

	r := httptest.NewRequest(http.MethodGet, "/hls/../outside.ts", nil)
	w := httptest.NewRecorder()

	(&Handlers{}).HandleHLSRequest(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
