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
		// check it. The HLS case below is public by explicit policy.
		{"auth header on viewer page is gated", newReq("GET", "/", "Bearer xyz"), "github-auth", true, true, false, gateLogin},
		{"auth header on HLS segment is allowed in website-only mode", newReq("GET", "/hls/0.ts", "anything"), "github-auth", true, true, false, gateAllow},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, loginURL := decideGate(tc.req, tc.slug, tc.loaded, tc.secretConfigured, tc.sessionValid, defaultAuthGateSettings())
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
		{"auth header on HLS is allowed in website-only mode", req("/hls/0.ts", "Bearer t"), "g", "third-party-player"},
		{"auth header on public api/config is NOT a bypass", req("/api/config", "Bearer t"), "g", ""},
		// Real files in the embedded web build are exempt (resolved by fs.Stat).
		{"admin stylesheet exempt", req("/styles/admin/chat.css", ""), "g", "static-assets"},
		{"service worker exempt", req("/sw.js", ""), "g", "static-assets"},
		{"web manifest exempt", req("/manifest.json", ""), "g", "static-assets"},
		// The SPA index is a page, not an asset → gated (would leak the viewer UI).
		{"index.html is a page, gated", req("/index.html", ""), "g", ""},
		// HLS media isn't in the web embed, so it's never matched here.
		{"HLS segment is an external-player exemption", req("/hls/stream/0.ts", ""), "g", "third-party-player"},
		{"plain viewer path is an external-player exemption", req("/hls/0.ts", ""), "g", "third-party-player"},
		{"plugin path with no active gate is gated", req("/plugins/g", ""), "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := gateExemptionFor(tc.req, tc.gateSlug, defaultAuthGateSettings()); got != tc.want {
				t.Fatalf("gateExemptionFor: got %q want %q", got, tc.want)
			}
		})
	}
}

// TestAuthGateAccessPolicyMatrix walks every access mode against the paths
// whose treatment the mode is supposed to change, plus the paths it must never
// change. Table-complete on purpose: this is the whole operator-facing
// contract for what a gate does and does not cover.
func TestAuthGateAccessPolicyMatrix(t *testing.T) {
	websiteOnly := defaultAuthGateSettings()
	websiteAndStream := AuthGateSettings{AccessMode: AuthGateAccessWebsiteAndStream}
	websiteStreamAndStatus := AuthGateSettings{AccessMode: AuthGateAccessWebsiteStreamAndStatus}

	tests := []struct {
		mode     string
		settings AuthGateSettings
		path     string
		want     gateOutcome
	}{
		// Website-only (the default): the viewer UI is the only thing gated.
		{"website-only", websiteOnly, "/", gateLogin},
		{"website-only", websiteOnly, "/embed/video", gateLogin},
		{"website-only", websiteOnly, "/api/config", gateLogin},
		{"website-only", websiteOnly, "/hls/stream.m3u8", gateAllow},
		{"website-only", websiteOnly, "/hls/0/stream0.ts", gateAllow},
		{"website-only", websiteOnly, "/api/status", gateAllow},
		{"website-only", websiteOnly, "/api/yp", gateAllow},

		// The second mode adds HLS while leaving status public.
		{"website-and-stream", websiteAndStream, "/hls/stream.m3u8", gateLogin},
		{"website-and-stream", websiteAndStream, "/hls/0/stream0.ts", gateLogin},
		{"website-and-stream", websiteAndStream, "/api/status", gateAllow},
		{"website-and-stream", websiteAndStream, "/api/yp", gateAllow},
		{"website-and-stream", websiteAndStream, "/", gateLogin},

		// The final mode adds stream status. The directory endpoint stays past
		// the gate and is switched off by the YP handler instead.
		{"website-stream-and-status", websiteStreamAndStatus, "/hls/0/stream0.ts", gateLogin},
		{"website-stream-and-status", websiteStreamAndStatus, "/api/status", gateLogin},
		{"website-stream-and-status", websiteStreamAndStatus, "/api/yp", gateAllow},

		// No mode may ever re-gate the recovery or self-credentialed surfaces.
		{"website-stream-and-status", websiteStreamAndStatus, "/admin", gateAllow},
		{"website-stream-and-status", websiteStreamAndStatus, "/admin/plugins", gateAllow},
		{"website-stream-and-status", websiteStreamAndStatus, "/api/admin/plugins", gateAllow},
		{"website-stream-and-status", websiteStreamAndStatus, "/plugins/gate/", gateAllow},
		{"website-stream-and-status", websiteStreamAndStatus, "/api/integrations/streamtitle", gateAllow},
	}

	for _, tc := range tests {
		t.Run(tc.mode+" "+tc.path, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tc.path, nil)
			got, _ := decideGate(r, "gate", true, true, false, tc.settings)
			if got != tc.want {
				t.Fatalf("%s %s: outcome got %d want %d", tc.mode, tc.path, got, tc.want)
			}
		})
	}
}

// A valid session must outrank every access mode: a signed-in viewer never gets
// bounced off a path the mode closed to anonymous callers.
func TestAuthGateSessionOutranksPolicy(t *testing.T) {
	settings := AuthGateSettings{AccessMode: AuthGateAccessWebsiteStreamAndStatus}
	for _, path := range []string{"/", "/hls/0/stream0.ts", "/api/status", "/api/config"} {
		r := httptest.NewRequest(http.MethodGet, path, nil)
		if got, _ := decideGate(r, "gate", true, true, true, settings); got != gateAllow {
			t.Fatalf("%s with a valid session: got %d want gateAllow", path, got)
		}
	}
}

// An access mode must never turn a fail-closed gate into an open one, and a
// broken gate must not start blocking paths the operator declared public.
func TestAuthGateBrokenGateInteractsCorrectlyWithPolicy(t *testing.T) {
	tests := []struct {
		name             string
		path             string
		settings         AuthGateSettings
		loaded           bool
		secretConfigured bool
		want             gateOutcome
	}{
		{"gated path, plugin not loaded", "/", defaultAuthGateSettings(), false, true, gateUnavailable},
		{"gated path, no signing secret", "/", defaultAuthGateSettings(), true, false, gateUnavailable},
		{"protected stream, plugin not loaded", "/hls/0.ts", AuthGateSettings{AccessMode: AuthGateAccessWebsiteAndStream}, false, true, gateUnavailable},
		{"blocked status, no signing secret", "/api/status", AuthGateSettings{AccessMode: AuthGateAccessWebsiteStreamAndStatus}, true, false, gateUnavailable},
		// Paths the mode leaves public stay public: a broken gate must not take
		// down HLS or status for operators who never gated them.
		{"public stream survives a broken gate", "/hls/0.ts", defaultAuthGateSettings(), false, true, gateAllow},
		{"public status survives a broken gate", "/api/status", defaultAuthGateSettings(), false, true, gateAllow},
		// And the recovery surfaces are always reachable.
		{"admin survives a broken gate", "/admin/plugins", AuthGateSettings{AccessMode: AuthGateAccessWebsiteStreamAndStatus}, false, false, gateAllow},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tc.path, nil)
			got, _ := decideGate(r, "gate", tc.loaded, tc.secretConfigured, false, tc.settings)
			if got != tc.want {
				t.Fatalf("got %d want %d", got, tc.want)
			}
		})
	}
}
