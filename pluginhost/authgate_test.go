package pluginhost

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecideGate(t *testing.T) {
	newReq := func(method, path, authHeader string) *http.Request {
		r := httptest.NewRequest(method, path, nil)
		if authHeader != "" {
			r.Header.Set("Authorization", authHeader)
		}
		return r
	}

	cases := []struct {
		name             string
		req              *http.Request
		slug             string
		loaded           bool
		secretConfigured bool
		sessionValid     bool
		wantOutcome      gateOutcome
	}{
		{"gate off allows everything", newReq("GET", "/", ""), "", false, true, false, gateAllow},
		{"admin always allowed", newReq("GET", "/admin/config", ""), "github-auth", true, true, false, gateAllow},
		{"gate plugin namespace allowed (login screen)", newReq("GET", "/plugins/github-auth/", ""), "github-auth", true, true, false, gateAllow},
		{"gate plugin callback allowed", newReq("GET", "/plugins/github-auth/callback?code=x", ""), "github-auth", true, true, false, gateAllow},
		{"armed but not loaded fails closed", newReq("GET", "/", ""), "github-auth", false, true, false, gateUnavailable},
		{"armed with no secret fails closed", newReq("GET", "/", ""), "github-auth", true, false, false, gateUnavailable},
		{"no session redirects to login", newReq("GET", "/", ""), "github-auth", true, true, false, gateLogin},
		{"valid session allowed", newReq("GET", "/hls/0.ts", ""), "github-auth", true, true, true, gateAllow},
		{"valid session survives a gate-plugin outage", newReq("GET", "/", ""), "github-auth", false, true, true, gateAllow},
		{"external-api route allowed", newReq("GET", "/api/integrations/status", "Bearer xyz"), "github-auth", true, true, false, gateAllow},
		// An Authorization header must NOT be a bypass on routes that never
		// check it — these are the gate-defeat vectors the exemption scoping fixes.
		{"auth header on viewer page is gated", newReq("GET", "/", "Bearer xyz"), "github-auth", true, true, false, gateLogin},
		{"auth header on HLS segment is gated", newReq("GET", "/hls/0.ts", "anything"), "github-auth", true, true, false, gateLogin},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, loginURL := decideGate(tc.req, tc.slug, tc.loaded, tc.secretConfigured, tc.sessionValid)
			if got != tc.wantOutcome {
				t.Fatalf("outcome: got %d want %d", got, tc.wantOutcome)
			}
			if got == gateLogin && !strings.HasPrefix(loginURL, "/plugins/github-auth/?return_to=") {
				t.Fatalf("login URL wrong: %q", loginURL)
			}
		})
	}
}

func TestGateExemptions(t *testing.T) {
	req := func(path, auth string) *http.Request {
		r := httptest.NewRequest("GET", path, nil)
		if auth != "" {
			r.Header.Set("Authorization", auth)
		}
		return r
	}
	cases := []struct {
		name     string
		req      *http.Request
		gateSlug string
		want     string // expected matching exemption, "" = gated
	}{
		{"admin root", req("/admin", ""), "g", "admin"},
		{"admin subpath", req("/admin/config", ""), "g", "admin"},
		{"admin API is exempt (cookie-authed, no Authorization header)", req("/api/admin/status", ""), "g", "admin"},
		{"admin API subpath", req("/api/admin/plugins/x/enable", ""), "g", "admin"},
		{"public api/config is still gated", req("/api/config", ""), "g", ""},
		{"active gate plugin root", req("/plugins/g", ""), "g", "active-gate-plugin"},
		{"active gate plugin subpath", req("/plugins/g/callback", ""), "g", "active-gate-plugin"},
		{"a different plugin is NOT exempt", req("/plugins/other/x", ""), "g", ""},
		{"external-api route exempt", req("/api/integrations/status", "Bearer t"), "g", "external-api"},
		{"external-api route exempt even without a header (route 401s itself)", req("/api/integrations/status", ""), "g", "external-api"},
		// An Authorization header is NOT a bypass off the external-API namespace.
		{"auth header on viewer page is NOT a bypass", req("/", "Bearer t"), "g", ""},
		{"auth header on HLS is NOT a bypass", req("/hls/0.ts", "Bearer t"), "g", ""},
		{"auth header on public api/config is NOT a bypass", req("/api/config", "Bearer t"), "g", ""},
		// Real files in the embedded web build are exempt (resolved by fs.Stat).
		{"admin stylesheet exempt", req("/styles/admin/chat.css", ""), "g", "static-assets"},
		{"service worker exempt", req("/sw.js", ""), "g", "static-assets"},
		{"web manifest exempt", req("/manifest.json", ""), "g", "static-assets"},
		// The SPA index is a page, not an asset → gated (would leak the viewer UI).
		{"index.html is a page, gated", req("/index.html", ""), "g", ""},
		// HLS media isn't in the web embed, so it's never matched here.
		{"HLS segment is NOT a web asset, stays gated", req("/hls/stream/0.ts", ""), "g", ""},
		{"plain viewer path is gated", req("/hls/0.ts", ""), "g", ""},
		{"plugin path with no active gate is gated", req("/plugins/g", ""), "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := gateExemptionFor(tc.req, tc.gateSlug); got != tc.want {
				t.Fatalf("gateExemptionFor: got %q want %q", got, tc.want)
			}
		})
	}
}
