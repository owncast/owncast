package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"sync"
	"time"
)

// adminSessionCookieName is the cookie carrying an opaque admin session
// token. It's sent on every request to the Owncast origin, so the host's
// admin auth check accepts it as an alternative to a Basic Auth header
// for clients that cannot inject one (notably plugin admin iframes).
const adminSessionCookieName = "owncast_admin_session"

// adminSessionTTL is how long a minted session is good for. The admin UI
// re-mints lazily, so an idle admin tab will silently get a fresh cookie
// on its next plugin-iframe navigation.
const adminSessionTTL = 8 * time.Hour

// adminSessionTokenBytes is the entropy of a minted session token before
// hex encoding. 32 bytes → 64 hex chars, comfortably resistant to guessing.
const adminSessionTokenBytes = 32

// adminSessionStore is an in-memory map of opaque session tokens to their
// expiry. Sessions are not persisted across restarts — the admin re-mints
// on the next request when the cookie no longer validates.
type adminSessionStore struct {
	mu       sync.Mutex
	sessions map[string]time.Time
}

func newAdminSessionStore() *adminSessionStore {
	return &adminSessionStore{sessions: make(map[string]time.Time)}
}

// mint generates a fresh token, records it with TTL adminSessionTTL, and
// returns the token. The caller is responsible for surfacing it to the
// client (typically as a Set-Cookie).
func (s *adminSessionStore) mint() (string, error) {
	buf := make([]byte, adminSessionTokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	token := hex.EncodeToString(buf)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[token] = time.Now().Add(adminSessionTTL)
	s.gcLocked(time.Now())
	return token, nil
}

// valid reports whether token is a currently-active session. An unknown
// or expired token returns false; expired tokens are pruned as a side
// effect so the map doesn't grow unboundedly.
func (s *adminSessionStore) valid(token string) bool {
	if token == "" {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	expiry, ok := s.sessions[token]
	if !ok {
		return false
	}
	if time.Now().After(expiry) {
		delete(s.sessions, token)
		return false
	}
	return true
}

// gcLocked drops any sessions whose expiry has passed. Called under mu
// from mint() so the cost is amortized over new sessions; the map size
// is bounded by concurrent admins, which is small.
func (s *adminSessionStore) gcLocked(now time.Time) {
	for token, expiry := range s.sessions {
		if now.After(expiry) {
			delete(s.sessions, token)
		}
	}
}

// hasValidAdminSessionCookie reports whether r carries an admin session
// cookie that is currently active. Used as an alternative auth path in
// IsAdminRequest for clients that cannot inject the Authorization header.
func (m *Middleware) hasValidAdminSessionCookie(r *http.Request) bool {
	c, err := r.Cookie(adminSessionCookieName)
	if err != nil {
		return false
	}
	return m.adminSessions.valid(c.Value)
}

// ensureAdminSessionCookie mints an admin session and sets the cookie
// on w when r doesn't already carry a valid one. Called from
// RequireAdminAuth after Basic Auth succeeds, so every authenticated
// admin API call keeps the cookie fresh as a side effect — there is no
// explicit "log in" endpoint the client must remember to call.
//
// A mint failure is silently swallowed; the admin API call itself still
// succeeds, and the next authenticated call will retry. The cookie is
// best-effort: if it's missing or expired, callers that can fall back
// to Basic Auth (the admin UI itself) keep working, and only contexts
// that can't (iframes) will get a one-time browser prompt.
func (m *Middleware) ensureAdminSessionCookie(w http.ResponseWriter, r *http.Request) {
	if m.hasValidAdminSessionCookie(r) {
		return
	}
	token, err := m.adminSessions.mint()
	if err != nil {
		return
	}
	// Secure mirrors whether the request itself came in over TLS:
	// always-on would silently break Owncast admin deployments running
	// over plain HTTP on a LAN, and always-off downgrades cookie
	// security on every HTTPS deployment. Reading r.TLS catches direct
	// HTTPS; checking X-Forwarded-Proto catches the common reverse-proxy
	// case where TLS terminates upstream.
	secure := r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
	// gosec G124: Secure is a constant `true`/`false` rather than a
	// literal, but it's request-scheme-aware (set when TLS is in play,
	// cleared when the admin is on a LAN HTTP deployment). The linter
	// can't model that and would force a single hard-coded value.
	http.SetCookie(w, &http.Cookie{ //nolint:gosec // G124: see comment above
		Name:     adminSessionCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(adminSessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   secure,
		// SameSite=Lax so the cookie reaches same-origin iframe loads and
		// top-level navigations triggered from the admin UI, without
		// being sent on third-party requests.
		SameSite: http.SameSiteLaxMode,
	})
}
