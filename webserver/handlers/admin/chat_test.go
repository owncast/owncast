package admin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBanIPAddressRejectsInvalidValue(t *testing.T) {
	handler := &Admin{}
	for _, body := range []string{
		`{"value": "not-an-ip"}`,
		`{"value": true}`,
		`{"value": 123}`,
		`{"value": ["127.0.0.1"]}`,
	} {
		t.Run(body, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
			rr := httptest.NewRecorder()

			handler.BanIPAddress(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("body %s: got status %d, want %d", body, rr.Code, http.StatusBadRequest)
			}
		})
	}
}
