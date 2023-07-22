package webserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/owncast/owncast/storage/data"
)

var srv *webServer

func TestMain(m *testing.M) {
	_, err := data.NewStore(":memory:")
	if err != nil {
		panic(err)
	}
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
