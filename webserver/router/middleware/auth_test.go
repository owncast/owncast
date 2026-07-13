package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	logtest "github.com/sirupsen/logrus/hooks/test"
)

func TestInvalidAccessTokenLogIncludesRequestContext(t *testing.T) {
	hook := logtest.NewGlobal()
	defer hook.Reset()

	req := httptest.NewRequest("POST", "http://owncast.test/api/integrations/chat?ignored=true", nil)
	req.RemoteAddr = "203.0.113.10:1234"
	req.Header.Set("Authorization", "Bearer secret-token")

	logInvalidAccessToken(req, "CAN_SEND_MESSAGES")

	entry := hook.LastEntry()
	if entry == nil {
		t.Fatal("expected an invalid access token log entry")
	}
	for _, value := range []string{"203.0.113.10", "POST", "/api/integrations/chat", "CAN_SEND_MESSAGES"} {
		if !strings.Contains(entry.Message, value) {
			t.Errorf("log entry %q does not include %q", entry.Message, value)
		}
	}
	if strings.Contains(entry.Message, "secret-token") || strings.Contains(entry.Message, "ignored=true") {
		t.Errorf("log entry includes sensitive request data: %q", entry.Message)
	}
}

func TestRequireAdminAuthRejectsCrossOriginRequests(t *testing.T) {
	m := &Middleware{adminSessions: newAdminSessionStore()}
	handler := m.RequireAdminAuth(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodOptions} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "http://owncast.example/api/admin/config", nil)
			req.Header.Set("Origin", "https://attacker.example")
			req.Header.Set(adminCSRFHeader, "1")
			rec := httptest.NewRecorder()

			handler(rec, req)

			if rec.Code != http.StatusForbidden {
				t.Fatalf("expected %d, got %d", http.StatusForbidden, rec.Code)
			}
		})
	}
}

func TestRequireAdminAuthRequiresCSRFHeaderWithoutOrigin(t *testing.T) {
	m := &Middleware{adminSessions: newAdminSessionStore()}
	handler := m.RequireAdminAuth(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodPost, "http://owncast.example/api/admin/config", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected missing CSRF header to return %d, got %d", http.StatusForbidden, rec.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "http://owncast.example/api/admin/config", nil)
	req.Header.Set(adminCSRFHeader, "1")
	rec = httptest.NewRecorder()
	handler(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected CSRF-protected request to reach authentication, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "http://owncast.example/api/admin/config", nil)
	req.Header.Set("Authorization", "Bearer test")
	rec = httptest.NewRecorder()
	handler(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected originless Authorization client to reach authentication, got %d", rec.Code)
	}
}

func TestRequireAdminAuthAllowsSameOriginMutation(t *testing.T) {
	m := &Middleware{adminSessions: newAdminSessionStore()}
	handler := m.RequireAdminAuth(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodPost, "https://owncast.example/api/admin/config", nil)
	req.Host = "owncast.example"
	req.Header.Set("Origin", "https://owncast.example")
	req.Header.Set("X-Forwarded-Proto", "https")
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected same-origin request to reach authentication, got %d", rec.Code)
	}
}

func TestRequireAdminAuthAllowsDevelopmentPreflight(t *testing.T) {
	m := &Middleware{adminSessions: newAdminSessionStore()}
	handler := m.RequireAdminAuth(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	req := httptest.NewRequest(http.MethodOptions, "http://localhost:8080/api/admin/config", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected development preflight to return %d, got %d", http.StatusNoContent, rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Fatalf("unexpected allowed origin %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Headers"); !containsHeader(got, adminCSRFHeader) {
		t.Fatalf("expected %s in allowed headers, got %q", adminCSRFHeader, got)
	}
}

func containsHeader(value, header string) bool {
	for _, candidate := range strings.Split(value, ",") {
		if strings.EqualFold(strings.TrimSpace(candidate), header) {
			return true
		}
	}
	return false
}
