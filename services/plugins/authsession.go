package plugins

// Viewer-authentication sessions. A viewer-auth plugin (one holding the
// auth.gate permission) calls owncast.auth.grantSession after it verifies a
// visitor; the host mints a signed, stateless session cookie carrying the
// viewer's Owncast access token and an expiry, and attaches it to the in-flight
// on_http_request response. The gate (added separately) verifies the cookie
// signature + expiry on every request — pure crypto, no database lookup, so it
// is cheap enough for the per-segment HLS hot path.
//
// The plugin never sees the cookie or the signing secret: it only names the
// user (via grantSession) and the host owns the credential end to end.

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// SessionCookieName is the gate session cookie. Core owns it; plugins
	// cannot set it (the response writer strips a plugin-supplied value).
	SessionCookieName = "owncast_session"

	// DefaultSessionTTL is used when a plugin requests ttl 0. It is intentionally
	// short-ish because, with page-load-granularity revocation, the TTL is the
	// hard backstop on how long a revoked-but-idle viewer can linger.
	DefaultSessionTTL = 24 * time.Hour
	// MaxSessionTTL caps a plugin-requested ttl.
	MaxSessionTTL = 30 * 24 * time.Hour
)

// ClampSessionTTL resolves a plugin-requested ttl (in seconds) to a duration,
// applying the default for 0 and the ceiling for anything too large. Both the
// cookie MaxAge and the signed token's expiry derive from it, so they stay in
// sync.
func ClampSessionTTL(ttlSeconds int64) time.Duration {
	if ttlSeconds <= 0 {
		return DefaultSessionTTL
	}
	ttl := time.Duration(ttlSeconds) * time.Second
	if ttl > MaxSessionTTL {
		return MaxSessionTTL
	}
	return ttl
}

// SignSession produces a stateless session token: base64url(payload) "."
// base64url(HMAC-SHA256(payload)). payload is "userID|expUnix".
func SignSession(secret []byte, userID string, expUnix int64) string {
	payload := userID + "|" + strconv.FormatInt(expUnix, 10)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." +
		base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// VerifySession checks a token's signature against secret and its expiry
// against nowUnix, returning the user ID it carries. ok is false for a
// malformed, tampered, or expired token.
func VerifySession(secret []byte, token string, nowUnix int64) (userID string, ok bool) {
	dot := strings.IndexByte(token, '.')
	if dot < 0 {
		return "", false
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(token[:dot])
	if err != nil {
		return "", false
	}
	sig, err := base64.RawURLEncoding.DecodeString(token[dot+1:])
	if err != nil {
		return "", false
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write(payloadBytes)
	if !hmac.Equal(sig, mac.Sum(nil)) {
		return "", false
	}
	payload := string(payloadBytes)
	bar := strings.LastIndexByte(payload, '|')
	if bar < 0 {
		return "", false
	}
	exp, err := strconv.ParseInt(payload[bar+1:], 10, 64)
	if err != nil || nowUnix > exp {
		return "", false
	}
	return payload[:bar], true
}

// newSessionCookie builds the gate session cookie. SameSite=Lax is required:
// the auth provider's OAuth callback is a cross-site top-level redirect, and
// Strict would drop the cookie on return. Secure is set whenever the request
// arrived over HTTPS — this is an auth credential, so it should never travel in
// cleartext when the connection is encrypted. It is intentionally NOT forced on
// plain-HTTP requests, since many deployments terminate TLS at a proxy and
// serve the app over HTTP internally (the proxy hop is the secure boundary).
func NewSessionCookie(token string, ttl time.Duration, secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(ttl.Seconds()),
	}
}

// clearSessionCookie expires the gate session cookie (logout).
func ClearSessionCookie(secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}
}

// SessionFromRequest reads and verifies the gate session cookie on a request,
// returning the user ID it carries. ok is false when the cookie is absent,
// malformed, tampered, or expired. The gate middleware and the chat /ws bridge
// use this to resolve identity without a database hit.
func SessionFromRequest(secret []byte, r *http.Request, nowUnix int64) (userID string, ok bool) {
	c, err := r.Cookie(SessionCookieName)
	if err != nil || c.Value == "" {
		return "", false
	}
	return VerifySession(secret, c.Value, nowUnix)
}

// --- request-scoped grant sink ---------------------------------------------
//
// grantSession/endSession are called from inside a plugin's on_http_request
// wasm execution, but the cookie must be attached to the response the host
// writes AFTER the wasm returns. We bridge the two by stashing a sink in the
// call context (which extism propagates to host functions). serveDynamic seeds
// it before the call and drains it after.

type sessionAction struct {
	clear bool
	token string
	ttl   time.Duration
}

type authSink struct {
	mu      sync.Mutex
	actions []sessionAction
}

func (s *authSink) grant(token string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.actions = append(s.actions, sessionAction{token: token, ttl: ttl})
}

func (s *authSink) end() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.actions = append(s.actions, sessionAction{clear: true})
}

// cookies renders the recorded actions into Set-Cookie values, in order (a
// plugin should issue at most one per request, but the browser applies the
// last regardless). secure marks the cookies Secure (set when the request
// arrived over HTTPS).
func (s *authSink) cookies(secure bool) []*http.Cookie {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]*http.Cookie, 0, len(s.actions))
	for _, a := range s.actions {
		if a.clear {
			out = append(out, ClearSessionCookie(secure))
		} else {
			out = append(out, NewSessionCookie(a.token, a.ttl, secure))
		}
	}
	return out
}

// RequestIsSecure reports whether a request arrived over HTTPS, either directly
// (r.TLS) or via a TLS-terminating proxy that set X-Forwarded-Proto. Used to
// decide whether the session cookie should be marked Secure.
func RequestIsSecure(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	return strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}

type authSinkKeyType struct{}

var authSinkKey authSinkKeyType

// --- resolved-session propagation -------------------------------------------
//
// The gate middleware verifies the session cookie once per request; it stashes
// the access token it carries in the request context so downstream handlers
// (notably chat's /ws) can resolve the viewer's identity without re-verifying
// or re-reading the cookie. This is how a gate login flows through to chat.

type sessionTokenKeyType struct{}

var sessionTokenKey sessionTokenKeyType

// WithSessionToken returns a context carrying the verified gate-session access
// token.
func WithSessionToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, sessionTokenKey, token)
}

// SessionTokenFromContext returns the gate-session access token the middleware
// resolved for this request, or "" if the request had no valid gate session.
func SessionTokenFromContext(ctx context.Context) string {
	t, _ := ctx.Value(sessionTokenKey).(string)
	return t
}

// withAuthSink returns a context carrying a fresh sink and the sink itself.
func withAuthSink(ctx context.Context) (context.Context, *authSink) {
	sink := &authSink{}
	return context.WithValue(ctx, authSinkKey, sink), sink
}

// authSinkFrom returns the sink in ctx, or nil when none was seeded (e.g. the
// test harness, which exercises host functions without an HTTP response).
func authSinkFrom(ctx context.Context) *authSink {
	sink, _ := ctx.Value(authSinkKey).(*authSink)
	return sink
}
