package pluginhost

// End-to-end tests for the gate middleware itself, as opposed to the pure
// decision in authgate_test.go. These exercise the parts that actually hold the
// credential: cookie verification, the 302-vs-401 split, identity propagation
// into the request context, and the onAuthCheck revocation path.
//
// serveGated takes the resolved gate state as arguments, so the whole policy —
// including the fail-closed branch, which needs a gate that is enabled but not
// running — is reachable without a wasm runtime.

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	plugins "github.com/owncast/owncast/services/plugins"
	"github.com/owncast/owncast/services/plugins/kv"
)

const testGateSlug = "gate"

var testAuthSecret = []byte("test-signing-secret-not-a-real-one")

// gatedHost builds a Host wired for the gate middleware, with the access policy
// pre-seeded so the settings read never touches storage.
func gatedHost(t *testing.T, settings AuthGateSettings) *Host {
	t.Helper()
	h := &Host{authSecret: testAuthSecret, kv: kv.NewMemory()}
	h.authGateSettings.Store(authGateSettingsSnapshot{slug: testGateSlug, settings: settings})
	return h
}

// nextRecorder is the downstream handler. It records whether the gate let the
// request through, and what identity it propagated.
type nextRecorder struct {
	served bool
	token  string
}

func (n *nextRecorder) handler() http.Handler {
	return http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		n.served = true
		n.token = plugins.SessionTokenFromContext(r.Context())
	})
}

// sessionCookie returns a valid gate cookie for token, expiring in ttl.
func sessionCookie(token string, ttl time.Duration) *http.Cookie {
	signed := plugins.SignSession(testAuthSecret, token, time.Now().Add(ttl).Unix())
	return plugins.NewSessionCookie(signed, ttl, false)
}

func TestServeGated_NoGateIsAPassthrough(t *testing.T) {
	h := gatedHost(t, defaultAuthGateSettings())
	next := &nextRecorder{}
	w := httptest.NewRecorder()

	h.serveGated(w, httptest.NewRequest(http.MethodGet, "/", nil), next.handler(), "", false)

	if !next.served {
		t.Fatal("no enabled gate must be a passthrough")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
}

func TestServeGated_AnonymousViewerIsRedirectedWithReturnTo(t *testing.T) {
	h := gatedHost(t, defaultAuthGateSettings())
	next := &nextRecorder{}
	w := httptest.NewRecorder()

	r := httptest.NewRequest(http.MethodGet, "/?utm=x", nil)
	h.serveGated(w, r, next.handler(), testGateSlug, true)

	if next.served {
		t.Fatal("anonymous viewer reached the downstream handler")
	}
	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302", w.Code)
	}
	want := "/plugins/" + testGateSlug + "/?return_to=%2F%3Futm%3Dx"
	if got := w.Header().Get("Location"); got != want {
		t.Fatalf("Location = %q, want %q", got, want)
	}
}

// A browser navigation gets an HTML redirect; anything else gets a status an
// XHR or API caller can act on.
func TestServeGated_NonNavigationGets401(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			h := gatedHost(t, defaultAuthGateSettings())
			next := &nextRecorder{}
			w := httptest.NewRecorder()

			h.serveGated(w, httptest.NewRequest(method, "/api/chat", nil), next.handler(), testGateSlug, true)

			if next.served {
				t.Fatal("unauthenticated request reached the downstream handler")
			}
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want 401", w.Code)
			}
		})
	}
}

func TestServeGated_ValidSessionPassesAndPropagatesIdentity(t *testing.T) {
	h := gatedHost(t, defaultAuthGateSettings())
	next := &nextRecorder{}
	w := httptest.NewRecorder()

	// Not the index: index navigation additionally runs onAuthCheck, which is
	// covered separately below.
	r := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	r.AddCookie(sessionCookie("access-token-123", time.Hour))
	h.serveGated(w, r, next.handler(), testGateSlug, true)

	if !next.served {
		t.Fatalf("valid session was rejected (status %d)", w.Code)
	}
	if next.token != "access-token-123" {
		t.Fatalf("propagated token = %q, want access-token-123", next.token)
	}
}

// The credential is a signed cookie, so the gate must reject every way of
// presenting a bad one. Each of these previously had no coverage above the
// crypto primitive.
func TestServeGated_RejectsBadCredentials(t *testing.T) {
	valid := sessionCookie("access-token-123", time.Hour)

	tests := []struct {
		name   string
		cookie *http.Cookie
	}{
		{"no cookie", nil},
		{"expired session", sessionCookie("access-token-123", -time.Minute)},
		{"tampered payload", &http.Cookie{
			Name:  valid.Name,
			Value: strings.Replace(valid.Value, valid.Value[:4], "AAAA", 1),
		}},
		{"signature stripped", &http.Cookie{
			Name:  valid.Name,
			Value: strings.Split(valid.Value, ".")[0],
		}},
		{"garbage", &http.Cookie{Name: valid.Name, Value: "not-a-token"}},
		{"token signed with another secret", func() *http.Cookie {
			other := plugins.SignSession([]byte("a different secret"), "access-token-123", time.Now().Add(time.Hour).Unix())
			return plugins.NewSessionCookie(other, time.Hour, false)
		}()},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := gatedHost(t, defaultAuthGateSettings())
			next := &nextRecorder{}
			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.cookie != nil {
				r.AddCookie(tc.cookie)
			}
			h.serveGated(w, r, next.handler(), testGateSlug, true)

			if next.served {
				t.Fatal("gate accepted an invalid credential")
			}
			if w.Code != http.StatusFound {
				t.Fatalf("status = %d, want 302 to login", w.Code)
			}
		})
	}
}

// An armed gate whose plugin is not running must fail closed, and must not be
// talked out of it by a would-be viewer.
func TestServeGated_FailsClosedWhenGateUnavailable(t *testing.T) {
	h := gatedHost(t, defaultAuthGateSettings())
	next := &nextRecorder{}
	w := httptest.NewRecorder()

	h.serveGated(w, httptest.NewRequest(http.MethodGet, "/", nil), next.handler(), testGateSlug, false)

	if next.served {
		t.Fatal("request served while the gate was unavailable")
	}
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", w.Code)
	}
}

// Losing the signing secret means sessions cannot be verified at all, so the
// gate must close rather than wave everyone through.
func TestServeGated_FailsClosedWithoutASigningSecret(t *testing.T) {
	h := &Host{kv: kv.NewMemory()}
	h.authGateSettings.Store(authGateSettingsSnapshot{slug: testGateSlug, settings: defaultAuthGateSettings()})
	next := &nextRecorder{}
	w := httptest.NewRecorder()

	h.serveGated(w, httptest.NewRequest(http.MethodGet, "/", nil), next.handler(), testGateSlug, true)

	if next.served {
		t.Fatal("request served with no signing secret configured")
	}
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", w.Code)
	}
}

// An operator must always be able to reach the admin API to disable a gate that
// has broken the site, even while the gate is failing closed.
func TestServeGated_AdminStaysReachableWhileFailingClosed(t *testing.T) {
	h := gatedHost(t, AuthGateSettings{AccessMode: AuthGateAccessWebsiteStreamAndStatus})
	for _, path := range []string{"/admin/plugins", "/api/admin/plugins/gate/disable"} {
		t.Run(path, func(t *testing.T) {
			next := &nextRecorder{}
			w := httptest.NewRecorder()

			h.serveGated(w, httptest.NewRequest(http.MethodGet, path, nil), next.handler(), testGateSlug, false)

			if !next.served {
				t.Fatalf("admin path %s was blocked (status %d)", path, w.Code)
			}
		})
	}
}

// The access modes, exercised through the real middleware rather than the pure
// decision, including the status codes third-party clients actually see.
func TestServeGated_AccessPolicyThroughMiddleware(t *testing.T) {
	tests := []struct {
		name     string
		settings AuthGateSettings
		path     string
		wantPass bool
	}{
		{"website-only lets a player fetch the playlist", defaultAuthGateSettings(), "/hls/stream.m3u8", true},
		{"website-only lets a monitor read status", defaultAuthGateSettings(), "/api/status", true},
		{"website-only still gates the viewer page", defaultAuthGateSettings(), "/", false},
		{"website-and-stream blocks the playlist", AuthGateSettings{AccessMode: AuthGateAccessWebsiteAndStream}, "/hls/stream.m3u8", false},
		{"website-and-stream blocks segments", AuthGateSettings{AccessMode: AuthGateAccessWebsiteAndStream}, "/hls/0/stream0.ts", false},
		{"website-and-stream leaves status public", AuthGateSettings{AccessMode: AuthGateAccessWebsiteAndStream}, "/api/status", true},
		{"website-stream-and-status blocks status", AuthGateSettings{AccessMode: AuthGateAccessWebsiteStreamAndStatus}, "/api/status", false},
		{"website-stream-and-status also blocks HLS", AuthGateSettings{AccessMode: AuthGateAccessWebsiteStreamAndStatus}, "/hls/stream.m3u8", false},
		{"directory endpoint is never gate-challenged", AuthGateSettings{AccessMode: AuthGateAccessWebsiteStreamAndStatus}, "/api/yp", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := gatedHost(t, tc.settings)
			next := &nextRecorder{}
			w := httptest.NewRecorder()

			h.serveGated(w, httptest.NewRequest(http.MethodGet, tc.path, nil), next.handler(), testGateSlug, true)

			if next.served != tc.wantPass {
				t.Fatalf("%s: served = %v want %v (status %d)", tc.path, next.served, tc.wantPass, w.Code)
			}
		})
	}
}

// denySession is the revocation primitive: a gate that says "this viewer is no
// longer entitled" must actually strip the credential, not merely redirect. A
// redirect alone would leave the cookie in place and the viewer signed in.
func TestDenySession_ClearsTheCookieAndBounces(t *testing.T) {
	h := gatedHost(t, defaultAuthGateSettings())
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(sessionCookie("access-token-123", time.Hour))

	h.denySession(w, r, testGateSlug, false)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302", w.Code)
	}
	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("got %d cookies, want exactly the cleared session cookie", len(cookies))
	}
	cleared := cookies[0]
	if cleared.Value != "" || cleared.MaxAge >= 0 {
		t.Fatalf("session cookie was not cleared: value=%q maxAge=%d", cleared.Value, cleared.MaxAge)
	}

	// Prove it: a request replaying the cleared cookie no longer authenticates.
	replay := httptest.NewRequest(http.MethodGet, "/", nil)
	replay.AddCookie(cleared)
	if _, ok := plugins.SessionFromRequest(testAuthSecret, replay, time.Now().Unix()); ok {
		t.Fatal("cleared cookie still verifies as a valid session")
	}
}

// The refresh verdict re-issues the cookie for sliding expiry. The new cookie
// must carry the same identity and still verify.
func TestRefreshedSessionCookieStaysValid(t *testing.T) {
	ttl := plugins.ClampSessionTTL(0)
	refreshed := plugins.SignSession(testAuthSecret, "access-token-123", time.Now().Add(ttl).Unix())
	cookie := plugins.NewSessionCookie(refreshed, ttl, false)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(cookie)

	token, ok := plugins.SessionFromRequest(testAuthSecret, r, time.Now().Unix())
	if !ok {
		t.Fatal("refreshed cookie does not verify")
	}
	if token != "access-token-123" {
		t.Fatalf("refreshed token = %q, want access-token-123", token)
	}
}

// onAuthCheck runs only where the gate intends it to: a viewer landing on the
// page. It must not fire on assets, API calls, or exempt paths, or a revocation
// check would run on the HLS hot path.
func TestIsIndexNavigation(t *testing.T) {
	tests := []struct {
		method string
		path   string
		want   bool
	}{
		{http.MethodGet, "/", true},
		{http.MethodHead, "/", true},
		{http.MethodPost, "/", false},
		{http.MethodGet, "/index.html", false},
		{http.MethodGet, "/api/config", false},
		{http.MethodGet, "/hls/0/stream0.ts", false},
	}
	for _, tc := range tests {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			if got := isIndexNavigation(httptest.NewRequest(tc.method, tc.path, nil)); got != tc.want {
				t.Fatalf("isIndexNavigation = %v want %v", got, tc.want)
			}
		})
	}
}

// callAuthCheck must report an error (which the caller turns into a deny) when
// it cannot build the identity to hand the plugin, rather than returning a
// zero verdict that reads as "ok".
func TestCallAuthCheck_FailsClosedWithoutAnIdentity(t *testing.T) {
	t.Run("no user resolver wired", func(t *testing.T) {
		h := &Host{}
		if _, err := h.callAuthCheck(t.Context(), testGateSlug, "tok"); err == nil {
			t.Fatal("expected an error when no user resolver is configured")
		}
	})

	t.Run("token resolves to no user", func(t *testing.T) {
		h := &Host{userByToken: func(string) *plugins.HostUser { return nil }}
		if _, err := h.callAuthCheck(t.Context(), testGateSlug, "tok"); err == nil {
			t.Fatal("expected an error when the token resolves to no user")
		}
	})
}
