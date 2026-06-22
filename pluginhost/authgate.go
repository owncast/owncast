package pluginhost

// Viewer-authentication gate middleware. When an admin has enabled a plugin
// holding the auth.gate permission, this blanket middleware (mounted ahead of
// every route) requires a valid signed session cookie before any request is
// served — the whole web server is gated, not just the page. A few things stay
// reachable so a visitor can actually log in (the gate plugin's own routes) or
// so other credentials still work (admin, external-API tokens).
//
// The per-request check is pure crypto (verify the cookie signature + expiry),
// no plugin call and no database lookup, so it is cheap enough for the
// per-segment HLS hot path.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	plugins "github.com/owncast/owncast/services/plugins"
	"github.com/owncast/owncast/static"
)

type gateOutcome int

const (
	gateAllow gateOutcome = iota
	gateLogin
	gateUnavailable
)

// authUnavailableHTML is shown (503) when a gate is armed but its plugin is not
// available — the fail-closed state. The admin can still reach /admin to fix it.
const authUnavailableHTML = `<!doctype html><html><head><meta charset="utf-8">` +
	`<title>Authentication unavailable</title></head><body>` +
	`<h1>Authentication temporarily unavailable</h1>` +
	`<p>This site requires sign-in, but the authentication plugin is not ` +
	`currently available. Please try again shortly.</p></body></html>`

// decideGate is the pure decision for one request: allow it, send it to the
// login screen, or serve the fail-closed page. gateSlug is the enabled
// auth.gate plugin's slug ("" when none is enabled); gateLoaded reports whether
// that plugin is actually running.
// The viewer-auth gate is opt-out: every request is gated unless a rule here
// exempts it. Keeping the exemptions as a small, named, individually-testable
// list — rather than string compares buried in the decision — makes the policy
// explicit: adding a bypass means adding a named rule, in the open, here.
type gateExemption struct {
	name  string
	allow func(r *http.Request, gateSlug string) bool
}

var gateExemptions = []gateExemption{
	{"admin", exemptAdminRoutes},
	{"static-assets", exemptStaticAssets},
	{"active-gate-plugin", exemptActiveGatePluginRoutes},
	{"external-api", exemptExternalAPIRoutes},
}

// externalAPIPrefix is the namespace for third-party API routes. Each handler
// there validates its own access token via RequireExternalAPIAccessToken.
const externalAPIPrefix = "/api/integrations/"

// webFS is the embedded web build (rooted at web/), resolved once. We test
// request paths against it to decide what's a static asset.
var webFS = static.GetWeb()

// exemptStaticAssets exempts any request that resolves to a real file in the
// embedded web build — the Next.js bundles, stylesheets, fonts, images, service
// worker, manifest, etc. that the (exempt) admin app needs to render. This is
// deliberately based on "does the build actually ship this file" rather than a
// hardcoded list of path prefixes or extensions, so it stays correct as webpack
// emits hash-named and relocated assets across builds.
//
// Two things are intentionally NOT matched, so they stay gated:
//   - the SPA HTML entry points (index.html, the embed shell): those are pages,
//     not assets — serving them unauthenticated would leak the viewer UI.
//   - HLS media (/hls/*.ts, *.m3u8): it lives in the HLS store, not the web
//     embed, so fs.Stat here never finds it.
func exemptStaticAssets(r *http.Request, _ string) bool {
	p := strings.TrimPrefix(r.URL.Path, "/")
	if p == "" || strings.HasSuffix(p, ".html") {
		return false
	}
	info, err := fs.Stat(webFS, p)
	return err == nil && !info.IsDir()
}

// exemptAdminRoutes: the admin surfaces carry their own credential (Basic auth
// or the admin session cookie) and have their own RequireAdminAuth gate, so
// they are never behind the viewer gate — an operator can always reach the
// dashboard and the admin API, even to disable a broken gate. Covers both the
// admin SPA (/admin/*) and the admin API (/api/admin/*); the latter often
// authenticates via the admin session cookie rather than an Authorization
// header, so the self-credentialed rule alone would miss it.
func exemptAdminRoutes(r *http.Request, _ string) bool {
	p := r.URL.Path
	return p == "/admin" || strings.HasPrefix(p, "/admin/") ||
		p == "/api/admin" || strings.HasPrefix(p, "/api/admin/")
}

// exemptActiveGatePluginRoutes: the active gate plugin's own namespace must stay
// reachable while unauthenticated — it serves the login screen, the auth
// callback, and its assets. Dynamic: only the currently-enabled gate plugin is
// exempt; every other plugin's routes stay gated.
func exemptActiveGatePluginRoutes(r *http.Request, gateSlug string) bool {
	if gateSlug == "" {
		return false
	}
	base := "/plugins/" + gateSlug
	return r.URL.Path == base || strings.HasPrefix(r.URL.Path, base+"/")
}

// exemptExternalAPIRoutes: the external-API (third-party token) routes under
// /api/integrations/ carry their own Bearer token, validated per-route by
// RequireExternalAPIAccessToken. A token client carries no session cookie, so
// gating these would bounce a valid API call into the HTML login screen; a
// token-less call still hits the route's own 401. (Admin Basic-auth /
// admin-session routes are exempted separately by exemptAdminRoutes.)
//
// Scoped to that namespace on purpose. An earlier version exempted ANY request
// carrying an Authorization header, on any path. But the viewer page (/) and
// the HLS handlers never inspect Authorization, so "Authorization: anything"
// let an anonymous visitor walk straight through the gate to the page and the
// live video. Gating is opt-out: a credential a route never checks must never
// be a bypass on that route.
func exemptExternalAPIRoutes(r *http.Request, _ string) bool {
	return strings.HasPrefix(r.URL.Path, externalAPIPrefix)
}

// gateExemptionFor returns the name of the first exemption that applies, or ""
// if the request is subject to the gate.
func gateExemptionFor(r *http.Request, gateSlug string) string {
	for _, e := range gateExemptions {
		if e.allow(r, gateSlug) {
			return e.name
		}
	}
	return ""
}

// decideGate is the pure gate decision. secretConfigured reports whether a
// signing secret exists (without one an armed gate can't verify sessions, so it
// fails closed); sessionValid reports whether the request carries a valid
// session cookie (the middleware verifies it once, up front).
func decideGate(r *http.Request, gateSlug string, gateLoaded bool, secretConfigured bool, sessionValid bool) (gateOutcome, string) {
	// Gate off: no enabled auth.gate plugin.
	if gateSlug == "" {
		return gateAllow, ""
	}

	// A valid session is always allowed — checked first so it short-circuits the
	// hot path (no exemption work, no per-request file stat), and so existing
	// sessions keep working even while the gate plugin is unavailable (the cookie
	// is verified with crypto alone, no plugin needed).
	if sessionValid {
		return gateAllow, ""
	}

	// Explicit, opt-out policy: anything on the named exemption list bypasses
	// the gate (this also keeps admin + the gate plugin's own routes reachable
	// while the gate is unavailable, so a broken gate can be recovered from).
	if gateExemptionFor(r, gateSlug) != "" {
		return gateAllow, ""
	}

	// Armed but we can't verify sessions (no signing secret) or the gate plugin
	// isn't running (crashed, failed to load, auto-disabled) → fail closed.
	// Never silently open.
	if !secretConfigured || !gateLoaded {
		return gateUnavailable, ""
	}

	// Otherwise the visitor is unauthenticated: send them to the plugin's login
	// screen with a sanitized, same-origin return_to (it is built from this
	// server's own request URI, so it can't be an open redirect).
	loginURL := "/plugins/" + gateSlug + "/?return_to=" + url.QueryEscape(r.URL.RequestURI())
	return gateLogin, loginURL
}

// AuthGateMiddleware returns the blanket viewer-auth gate as chi middleware.
// It is a no-op while no auth.gate plugin is enabled.
func (h *Host) AuthGateMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now().Unix()
			secretConfigured := len(h.authSecret) > 0

			// Verify the session once, up front, and propagate the identity it
			// carries into the request context — for ALL requests with a valid
			// cookie, including exempt paths (so the gate plugin's own login
			// screen can tell the viewer is already signed in) and chat (/ws).
			sessionValid := false
			sessionToken := ""
			if secretConfigured {
				if token, ok := plugins.SessionFromRequest(h.authSecret, r, now); ok {
					sessionValid = true
					sessionToken = token
					r = r.WithContext(plugins.WithSessionToken(r.Context(), token))
				}
			}

			slug, loaded := h.manager.ActiveAuthGate()
			outcome, loginURL := decideGate(r, slug, loaded, secretConfigured, sessionValid)
			switch outcome {
			case gateAllow:
				// On the viewer's page load, let the active gate plugin
				// re-validate (and optionally revoke or refresh) the session via
				// onAuthCheck. Only when the allow is due to a real session on a
				// non-exempt index navigation — not for exempt paths or API/asset
				// requests.
				if sessionValid && slug != "" && isIndexNavigation(r) && gateExemptionFor(r, slug) == "" {
					h.runAuthCheckOnIndex(w, r, next, slug, sessionToken)
					return
				}
				next.ServeHTTP(w, r)
			case gateUnavailable:
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = io.WriteString(w, authUnavailableHTML)
			case gateLogin:
				// Redirect navigations; 401 anything else so XHR/API callers
				// get a status they can act on instead of an HTML body.
				if r.Method == http.MethodGet || r.Method == http.MethodHead {
					//nolint:gosec // G710: loginURL is a same-origin /plugins/<slug>/ path with a query-escaped return_to, not attacker-controlled.
					http.Redirect(w, r, loginURL, http.StatusFound)
				} else {
					http.Error(w, "authentication required", http.StatusUnauthorized)
				}
			}
		})
	}
}

// authVerdict is the JSON a gate plugin's onAuthCheck returns.
type authVerdict struct {
	Action string `json:"action"` // "ok" | "refresh" | "deny"
	TTL    int64  `json:"ttl"`    // refresh: new lifetime in seconds (0 = host default)
	Reason string `json:"reason"`
}

// isIndexNavigation reports whether a request is a viewer landing on the main
// page — the point at which the gate re-validates the session via onAuthCheck.
func isIndexNavigation(r *http.Request) bool {
	return r.URL.Path == "/" && (r.Method == http.MethodGet || r.Method == http.MethodHead)
}

// runAuthCheckOnIndex calls the gate plugin's onAuthCheck for the current
// viewer and acts on the verdict: ok → continue, refresh → re-issue the cookie
// with a fresh expiry and continue, deny (or any error) → clear the session and
// bounce to the login screen. Errors fail closed.
func (h *Host) runAuthCheckOnIndex(w http.ResponseWriter, r *http.Request, next http.Handler, slug, token string) {
	secure := plugins.RequestIsSecure(r)
	verdict, err := h.callAuthCheck(r.Context(), slug, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "plugin auth: onAuthCheck failed for %s, denying session: %v\n", slug, err)
		h.denySession(w, r, slug, secure)
		return
	}
	switch verdict.Action {
	case "deny":
		h.denySession(w, r, slug, secure)
	case "refresh":
		ttl := plugins.ClampSessionTTL(verdict.TTL)
		refreshed := plugins.SignSession(h.authSecret, token, time.Now().Add(ttl).Unix())
		http.SetCookie(w, plugins.NewSessionCookie(refreshed, ttl, secure))
		next.ServeHTTP(w, r)
	default: // "ok" (and anything unrecognized) → let them through
		next.ServeHTTP(w, r)
	}
}

// callAuthCheck builds the identity envelope for token and invokes the gate
// plugin's onAuthCheck, returning its verdict.
func (h *Host) callAuthCheck(ctx context.Context, slug, token string) (authVerdict, error) {
	if h.userByToken == nil {
		return authVerdict{}, fmt.Errorf("no user resolver configured")
	}
	user := h.userByToken(token)
	if user == nil {
		return authVerdict{}, fmt.Errorf("session token resolves to no user")
	}
	input, err := json.Marshal(map[string]any{"user": user})
	if err != nil {
		return authVerdict{}, err
	}
	out, err := h.server.CallAuthCheck(ctx, slug, input)
	if err != nil {
		return authVerdict{}, err
	}
	var v authVerdict
	if err := json.Unmarshal(out, &v); err != nil {
		return authVerdict{}, err
	}
	return v, nil
}

// denySession clears the gate cookie and redirects the viewer to the login
// screen (called on a deny verdict or an onAuthCheck failure).
func (h *Host) denySession(w http.ResponseWriter, r *http.Request, slug string, secure bool) {
	http.SetCookie(w, plugins.ClearSessionCookie(secure))
	loginURL := "/plugins/" + slug + "/?return_to=" + url.QueryEscape(r.URL.RequestURI())
	http.Redirect(w, r, loginURL, http.StatusFound)
}
