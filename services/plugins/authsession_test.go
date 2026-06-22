package plugins

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSignVerifySession_Roundtrip(t *testing.T) {
	secret := []byte("test-secret-0123456789")
	now := time.Now().Unix()
	exp := now + 3600
	tok := SignSession(secret, "user-abc", exp)

	got, ok := VerifySession(secret, tok, now)
	if !ok {
		t.Fatal("expected valid token")
	}
	if got != "user-abc" {
		t.Fatalf("userID: got %q want %q", got, "user-abc")
	}
}

func TestVerifySession_Expired(t *testing.T) {
	secret := []byte("test-secret")
	exp := int64(1000)
	tok := SignSession(secret, "u1", exp)
	if _, ok := VerifySession(secret, tok, exp+1); ok {
		t.Fatal("expected expired token to be rejected")
	}
	if _, ok := VerifySession(secret, tok, exp); !ok {
		t.Fatal("token should still be valid at exactly exp")
	}
}

func TestVerifySession_Tampered(t *testing.T) {
	secret := []byte("test-secret")
	now := time.Now().Unix()
	tok := SignSession(secret, "u1", now+3600)

	// Wrong secret.
	if _, ok := VerifySession([]byte("other-secret"), tok, now); ok {
		t.Fatal("expected rejection under a different secret")
	}
	// Tampered signature: flip the FIRST signature character. (Not the last —
	// for a 32-byte HMAC the final base64 character carries only 4 significant
	// bits, so several characters decode to the same signature; flipping it is
	// sometimes a no-op, which made this test flaky. The first character always
	// changes a decoded byte.)
	dot := strings.IndexByte(tok, '.')
	repl := byte('A')
	if tok[dot+1] == 'A' {
		repl = 'B'
	}
	bad := tok[:dot+1] + string(repl) + tok[dot+2:]
	if _, ok := VerifySession(secret, bad, now); ok {
		t.Fatal("expected rejection of a tampered signature")
	}
	// Garbage.
	for _, junk := range []string{"", "no-dot", "a.b", ".", "x."} {
		if _, ok := VerifySession(secret, junk, now); ok {
			t.Fatalf("expected rejection of malformed token %q", junk)
		}
	}
}

func TestSessionFromRequest(t *testing.T) {
	secret := []byte("test-secret")
	now := time.Now().Unix()
	tok := SignSession(secret, "viewer-1", now+3600)

	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: SessionCookieName, Value: tok})
	if got, ok := SessionFromRequest(secret, r, now); !ok || got != "viewer-1" {
		t.Fatalf("SessionFromRequest: got (%q,%v) want (viewer-1,true)", got, ok)
	}

	// No cookie.
	if _, ok := SessionFromRequest(secret, httptest.NewRequest("GET", "/", nil), now); ok {
		t.Fatal("expected no session when cookie absent")
	}
}

func TestClampSessionTTL(t *testing.T) {
	if got := ClampSessionTTL(0); got != DefaultSessionTTL {
		t.Fatalf("ttl 0: got %v want default %v", got, DefaultSessionTTL)
	}
	if got := ClampSessionTTL(60); got != time.Minute {
		t.Fatalf("ttl 60: got %v want 1m", got)
	}
	if got := ClampSessionTTL(int64(MaxSessionTTL.Seconds()) + 10_000); got != MaxSessionTTL {
		t.Fatalf("oversized ttl: got %v want max %v", got, MaxSessionTTL)
	}
}

func TestAuthSink_CookiesGrantAndClear(t *testing.T) {
	s := &authSink{}
	if len(s.cookies(false)) != 0 {
		t.Fatal("empty sink should yield no cookies")
	}
	s.grant("tok123", time.Hour)
	cookies := s.cookies(false)
	if len(cookies) != 1 || cookies[0].Name != SessionCookieName || cookies[0].Value != "tok123" {
		t.Fatalf("grant cookie wrong: %+v", cookies)
	}
	if cookies[0].MaxAge != int(time.Hour.Seconds()) || !cookies[0].HttpOnly || cookies[0].SameSite != http.SameSiteLaxMode {
		t.Fatalf("grant cookie attributes wrong: %+v", cookies[0])
	}

	s2 := &authSink{}
	s2.end()
	cl := s2.cookies(false)
	if len(cl) != 1 || cl[0].MaxAge != -1 || cl[0].Value != "" {
		t.Fatalf("clear cookie wrong: %+v", cl)
	}
}

func TestSessionCookie_SecureFollowsRequest(t *testing.T) {
	// Secure off on plain HTTP, on over HTTPS.
	s := &authSink{}
	s.grant("tok", time.Hour)
	if s.cookies(false)[0].Secure {
		t.Fatal("cookie should not be Secure on plain HTTP")
	}
	if !s.cookies(true)[0].Secure {
		t.Fatal("cookie should be Secure over HTTPS")
	}

	// RequestIsSecure honors r.TLS and X-Forwarded-Proto.
	plain := httptest.NewRequest("GET", "http://x/", nil)
	if RequestIsSecure(plain) {
		t.Fatal("plain http request should not be secure")
	}
	fwd := httptest.NewRequest("GET", "http://x/", nil)
	fwd.Header.Set("X-Forwarded-Proto", "https")
	if !RequestIsSecure(fwd) {
		t.Fatal("X-Forwarded-Proto=https should be treated as secure")
	}
}
