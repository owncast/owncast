package webserver

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/owncast/owncast/core/data"
)

var srv *webServer

func TestMain(m *testing.M) {
	dbFile, err := os.CreateTemp(os.TempDir(), "owncast-test-db.db")
	if err != nil {
		panic(err)
	}

	data.SetupPersistence(dbFile.Name())
	srv = New()

	m.Run()
}

// TestPrometheusPath tests that the /debug/vars endpoint that
// Prometheus automatically enables is not exposed.
func TestPrometheusDebugPath(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/debug/vars", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, r)
	if w.Result().StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Result().StatusCode)
	}
}

func TestTestingEndpoint(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, r)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Result().StatusCode)
	}
}
