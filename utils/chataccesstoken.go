package utils

import (
	"net/http"
	"strings"
	"time"
)

// ChatAccessTokenCookieName is the cookie carrying a chat user's access
// token. It is set when a user registers or connects to chat so that all
// subsequent same-origin web requests (notably to plugin HTTP handlers under
// /plugins/) carry the user's identity without the frontend having to append
// the token to every URL. The host resolves the token to a user server-side;
// the raw token is never handed to plugins.
const ChatAccessTokenCookieName = "owncast_chat_token" //nolint:gosec // cookie name, not a credential

// chatAccessTokenCookieTTL is how long the identity cookie persists. Access
// tokens themselves don't expire, so this is generous. The cookie is
// refreshed on every chat connection anyway.
const chatAccessTokenCookieTTL = 365 * 24 * time.Hour

// requestIsSecure reports whether the request reached us over TLS, either
// directly (r.TLS) or via a TLS-terminating reverse proxy (X-Forwarded-Proto).
// Used to set the Secure cookie attribute only when it won't break a plain
// HTTP LAN deployment. Mirrors the admin session cookie logic.
func requestIsSecure(r *http.Request) bool {
	return r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}

// chatAccessTokenCookie builds the identity cookie for the given token.
// HttpOnly keeps it out of reach of page JavaScript (the frontend already
// holds the token in localStorage, so nothing needs to read it from here),
// and SameSite=Lax sends it on same-origin requests and top-level navigations
// without leaking it to third-party sites.
func chatAccessTokenCookie(r *http.Request, token string) *http.Cookie {
	return &http.Cookie{ //nolint:gosec // Secure is set conditionally via requestIsSecure to support plain-HTTP LAN deployments; HttpOnly and SameSite are always set.
		Name:     ChatAccessTokenCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(chatAccessTokenCookieTTL.Seconds()),
		HttpOnly: true,
		Secure:   requestIsSecure(r),
		SameSite: http.SameSiteLaxMode,
	}
}

// SetChatAccessTokenCookie writes the identity cookie to an HTTP response.
func SetChatAccessTokenCookie(w http.ResponseWriter, r *http.Request, token string) {
	if token == "" {
		return
	}
	http.SetCookie(w, chatAccessTokenCookie(r, token))
}

// AddChatAccessTokenCookieHeader adds the identity cookie to a header map.
// This is the path used for the WebSocket handshake response, where the
// cookie must be supplied as part of the upgrade response headers rather than
// written to a ResponseWriter after the connection is hijacked.
func AddChatAccessTokenCookieHeader(h http.Header, r *http.Request, token string) {
	if token == "" {
		return
	}
	h.Add("Set-Cookie", chatAccessTokenCookie(r, token).String())
}

// ChatAccessTokenFromRequest extracts a chat user's access token from a
// request, preferring the identity cookie and falling back to the legacy
// accessToken query parameter so existing callers keep working.
func ChatAccessTokenFromRequest(r *http.Request) string {
	if c, err := r.Cookie(ChatAccessTokenCookieName); err == nil && c.Value != "" {
		return c.Value
	}
	return r.URL.Query().Get("accessToken")
}
