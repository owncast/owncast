package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChatAccessTokenFromRequest(t *testing.T) {
	tests := []struct {
		name   string
		cookie string // cookie value; "" means no cookie set
		query  string // accessToken query value; "" means absent
		want   string
	}{
		{name: "cookie preferred over query", cookie: "ck", query: "qp", want: "ck"},
		{name: "cookie only", cookie: "ck", query: "", want: "ck"},
		{name: "query fallback", cookie: "", query: "qp", want: "qp"},
		{name: "empty cookie falls back to query", cookie: "", query: "qp", want: "qp"},
		{name: "neither", cookie: "", query: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := "/api/thing"
			if tt.query != "" {
				target += "?accessToken=" + tt.query
			}
			r := httptest.NewRequest(http.MethodGet, target, nil)
			if tt.cookie != "" {
				r.AddCookie(&http.Cookie{Name: ChatAccessTokenCookieName, Value: tt.cookie})
			}
			if got := ChatAccessTokenFromRequest(r); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSetChatAccessTokenCookie_Attributes(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/chat/register", nil)
	r.Header.Set("X-Forwarded-Proto", "https") // simulate TLS-terminating proxy
	rec := httptest.NewRecorder()

	SetChatAccessTokenCookie(rec, r, "tok-123")

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	c := cookies[0]
	if c.Name != ChatAccessTokenCookieName {
		t.Errorf("name: got %q, want %q", c.Name, ChatAccessTokenCookieName)
	}
	if c.Value != "tok-123" {
		t.Errorf("value: got %q", c.Value)
	}
	if !c.HttpOnly {
		t.Error("cookie should be HttpOnly so page JS can't read the token")
	}
	if c.SameSite != http.SameSiteLaxMode {
		t.Errorf("SameSite: got %v, want Lax", c.SameSite)
	}
	if c.Path != "/" {
		t.Errorf("path: got %q, want /", c.Path)
	}
	if !c.Secure {
		t.Error("Secure should be set when the request arrived over https")
	}
	if c.MaxAge <= 0 {
		t.Errorf("expected a positive MaxAge, got %d", c.MaxAge)
	}
}

func TestSetChatAccessTokenCookie_InsecureOnPlainHTTP(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/chat/register", nil) // no TLS, no X-Forwarded-Proto
	rec := httptest.NewRecorder()

	SetChatAccessTokenCookie(rec, r, "tok-123")

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	if cookies[0].Secure {
		t.Error("Secure must not be set on a plain-HTTP request, or LAN deployments break")
	}
}

func TestSetChatAccessTokenCookie_EmptyTokenNoOp(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/chat/register", nil)
	rec := httptest.NewRecorder()

	SetChatAccessTokenCookie(rec, r, "")

	if len(rec.Result().Cookies()) != 0 {
		t.Error("no cookie should be set for an empty token")
	}
}

func TestAddChatAccessTokenCookieHeader(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/ws?accessToken=tok", nil)
	h := http.Header{}

	AddChatAccessTokenCookieHeader(h, r, "tok")

	setCookie := h.Get("Set-Cookie")
	if setCookie == "" {
		t.Fatal("expected a Set-Cookie header for the handshake response")
	}

	// Parse it back the way a browser would, via a response.
	resp := http.Response{Header: http.Header{"Set-Cookie": h["Set-Cookie"]}}
	cookies := resp.Cookies()
	if len(cookies) != 1 || cookies[0].Name != ChatAccessTokenCookieName || cookies[0].Value != "tok" {
		t.Fatalf("unexpected parsed cookie: %+v", cookies)
	}

	// Empty token must not add a header.
	h2 := http.Header{}
	AddChatAccessTokenCookieHeader(h2, r, "")
	if h2.Get("Set-Cookie") != "" {
		t.Error("no Set-Cookie should be added for an empty token")
	}
}
