package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAdminSessionStore_MintValidateExpire(t *testing.T) {
	s := newAdminSessionStore()

	token, err := s.mint()
	if err != nil {
		t.Fatalf("mint: %v", err)
	}
	if token == "" {
		t.Fatal("mint returned empty token")
	}
	if !s.valid(token) {
		t.Error("freshly minted token should validate")
	}

	// Force expiry by rewriting the recorded time, then confirm validate
	// returns false and the entry is pruned.
	s.mu.Lock()
	s.sessions[token] = time.Now().Add(-1 * time.Minute)
	s.mu.Unlock()

	if s.valid(token) {
		t.Error("expired token should not validate")
	}
	s.mu.Lock()
	_, present := s.sessions[token]
	s.mu.Unlock()
	if present {
		t.Error("expired token should have been pruned from store")
	}
}

func TestAdminSessionStore_ValidRejectsUnknownAndEmpty(t *testing.T) {
	s := newAdminSessionStore()
	if s.valid("") {
		t.Error("empty token should not validate")
	}
	if s.valid("not-a-real-token") {
		t.Error("unknown token should not validate")
	}
}

func TestEnsureAdminSessionCookie_SetsWhenAbsent(t *testing.T) {
	m := &Middleware{adminSessions: newAdminSessionStore()}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/admin/anything", nil)
	m.ensureAdminSessionCookie(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	var cookie *http.Cookie
	for _, c := range res.Cookies() {
		if c.Name == adminSessionCookieName {
			cookie = c
			break
		}
	}
	if cookie == nil {
		t.Fatalf("response did not set %s cookie", adminSessionCookieName)
	}
	if cookie.Value == "" {
		t.Error("cookie has empty value")
	}
	if !cookie.HttpOnly {
		t.Error("cookie should be HttpOnly")
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Errorf("cookie SameSite = %v want %v", cookie.SameSite, http.SameSiteLaxMode)
	}
	if cookie.Path != "/" {
		t.Errorf("cookie Path = %q want %q", cookie.Path, "/")
	}
	if !m.adminSessions.valid(cookie.Value) {
		t.Error("minted cookie value should validate against store")
	}
}

func TestEnsureAdminSessionCookie_NoOpWhenAlreadyValid(t *testing.T) {
	// A request that already carries a valid cookie should not mint a new
	// one — otherwise every authenticated API call would rotate the token
	// and create a growing list of live sessions per admin.
	m := &Middleware{adminSessions: newAdminSessionStore()}
	existing, _ := m.adminSessions.mint()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/admin/anything", nil)
	req.AddCookie(&http.Cookie{Name: adminSessionCookieName, Value: existing})

	m.ensureAdminSessionCookie(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if len(res.Cookies()) != 0 {
		t.Errorf("expected no Set-Cookie when cookie already valid, got %d", len(res.Cookies()))
	}
}

func TestHasValidAdminSessionCookie(t *testing.T) {
	m := &Middleware{adminSessions: newAdminSessionStore()}
	token, _ := m.adminSessions.mint()

	// Missing cookie → false.
	req := httptest.NewRequest(http.MethodGet, "/anything", nil)
	if m.hasValidAdminSessionCookie(req) {
		t.Error("no cookie should not validate")
	}

	// Unknown token → false.
	req.AddCookie(&http.Cookie{Name: adminSessionCookieName, Value: "bogus"})
	if m.hasValidAdminSessionCookie(req) {
		t.Error("unknown cookie value should not validate")
	}

	// Valid minted token → true.
	req = httptest.NewRequest(http.MethodGet, "/anything", nil)
	req.AddCookie(&http.Cookie{Name: adminSessionCookieName, Value: token})
	if !m.hasValidAdminSessionCookie(req) {
		t.Error("minted cookie should validate")
	}
}
