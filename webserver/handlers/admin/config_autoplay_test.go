package admin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/owncast/owncast/persistence/configrepository"
)

func TestSetAutoplayValidValues(t *testing.T) {
	configRepository := configrepository.New(testDatastore)

	for _, value := range []string{"off", "always", "sound-only"} {
		body := makeConfigValueBody(value)
		req := httptest.NewRequest(http.MethodPost, "/api/admin/config/autoplay", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		testAdmin.SetAutoplay(w, req)

		resp := parseResponse(t, w)
		if !resp.Success {
			t.Fatalf("value %q: expected success, got error: %s", value, resp.Message)
		}
		if got := configRepository.GetAutoplay(); got != value {
			t.Errorf("value %q: expected persisted autoplay %q, got %q", value, value, got)
		}
	}
}

func TestSetAutoplayRejectsInvalidValue(t *testing.T) {
	configRepository := configrepository.New(testDatastore)

	// Set a known-good value so we can confirm a rejected set leaves it unchanged.
	good := httptest.NewRequest(http.MethodPost, "/api/admin/config/autoplay", strings.NewReader(makeConfigValueBody("always")))
	good.Header.Set("Content-Type", "application/json")
	testAdmin.SetAutoplay(httptest.NewRecorder(), good)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/config/autoplay", strings.NewReader(makeConfigValueBody("bogus")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	testAdmin.SetAutoplay(w, req)

	resp := parseResponse(t, w)
	if resp.Success {
		t.Fatal("expected failure for an invalid autoplay value")
	}
	if resp.Message != "invalid autoplay value" {
		t.Errorf("expected message 'invalid autoplay value', got %q", resp.Message)
	}
	if got := configRepository.GetAutoplay(); got != "always" {
		t.Errorf("expected autoplay to stay 'always' after a rejected set, got %q", got)
	}
}

func TestSetAutoplayRejectsNonStringValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/admin/config/autoplay", strings.NewReader(`{"value": true}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	testAdmin.SetAutoplay(w, req)

	resp := parseResponse(t, w)
	if resp.Success {
		t.Fatal("expected failure for a non-string autoplay value")
	}
}

func TestSetAutoplayRejectsGETMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/admin/config/autoplay", nil)
	w := httptest.NewRecorder()

	testAdmin.SetAutoplay(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for a GET request, got %d", w.Code)
	}
}
