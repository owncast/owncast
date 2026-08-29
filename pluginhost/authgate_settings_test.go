package pluginhost

// Tests for the admin-facing access-policy endpoint and the settings store
// behind it. The endpoint decides what an operator can change and for which
// plugin, and the store decides what the gate enforces on every request, so
// both need to be honest about failure.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/owncast/owncast/services/plugins"
	"github.com/owncast/owncast/services/plugins/kv"
)

// hostWithDiscoveredPlugins writes one .ocpkg per manifest into a temp plugins
// directory and starts a real Manager over it, so IsAuthGate is answered from
// genuine manifest discovery rather than a stub.
func hostWithDiscoveredPlugins(t *testing.T, manifests map[string][]byte) *Host {
	t.Helper()
	dir := t.TempDir()
	for slug, manifest := range manifests {
		pkg := buildPackageBytes(t, manifest, []byte("\x00asm\x01\x00\x00\x00"))
		if err := os.WriteFile(filepath.Join(dir, slug+".ocpkg"), pkg, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	mgr := plugins.NewManager(dir, &plugins.HostEnv{})
	if err := mgr.Start(t.Context()); err != nil {
		t.Fatalf("manager start: %v", err)
	}
	t.Cleanup(func() { mgr.Stop(t.Context()) })
	return &Host{manager: mgr, kv: kv.NewMemory()}
}

func gateManifest(slug string) []byte {
	return []byte(`{
		"api": "1",
		"name": "` + slug + `",
		"slug": "` + slug + `",
		"version": "0.1.0",
		"description": "auth gate fixture",
		"permissions": ["auth.gate"]
	}`)
}

func chatbotManifest(slug string) []byte {
	return []byte(`{
		"api": "1",
		"name": "` + slug + `",
		"slug": "` + slug + `",
		"version": "0.1.0",
		"description": "non-gate fixture",
		"permissions": ["chat.send"]
	}`)
}

func settingsRequest(t *testing.T, method, slug, body string) *http.Request {
	t.Helper()
	target := "/api/admin/plugins/" + slug + "/auth-settings"
	if body == "" {
		return httptest.NewRequest(method, target, nil)
	}
	return httptest.NewRequest(method, target, strings.NewReader(body))
}

// Settings must only be offered for plugins that can actually become the gate.
// Exposing them elsewhere would imply a plugin controls site access when it
// cannot.
func TestAuthGateSettings_OnlyForPluginsDeclaringAuthGate(t *testing.T) {
	h := hostWithDiscoveredPlugins(t, map[string][]byte{
		"gate":    gateManifest("gate"),
		"chatbot": chatbotManifest("chatbot"),
	})

	t.Run("gate plugin is served", func(t *testing.T) {
		w := httptest.NewRecorder()
		h.handleAuthGateSettings(w, settingsRequest(t, http.MethodGet, "gate", ""), "gate")
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
	})

	t.Run("non-gate plugin is not", func(t *testing.T) {
		w := httptest.NewRecorder()
		h.handleAuthGateSettings(w, settingsRequest(t, http.MethodGet, "chatbot", ""), "chatbot")
		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", w.Code)
		}
	})

	t.Run("unknown plugin is not", func(t *testing.T) {
		w := httptest.NewRecorder()
		h.handleAuthGateSettings(w, settingsRequest(t, http.MethodGet, "nope", ""), "nope")
		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", w.Code)
		}
	})
}

// The default is the compatibility contract: a freshly enabled gate protects
// the website and nothing else.
func TestAuthGateSettings_DefaultsAreWebsiteOnly(t *testing.T) {
	h := hostWithDiscoveredPlugins(t, map[string][]byte{"gate": gateManifest("gate")})
	w := httptest.NewRecorder()

	h.handleAuthGateSettings(w, settingsRequest(t, http.MethodGet, "gate", ""), "gate")

	var got AuthGateSettings
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.AccessMode != AuthGateAccessWebsiteOnly {
		t.Fatalf("default access mode = %q, want %q", got.AccessMode, AuthGateAccessWebsiteOnly)
	}
}

func TestAuthGateSettings_SaveThenReadBack(t *testing.T) {
	h := hostWithDiscoveredPlugins(t, map[string][]byte{"gate": gateManifest("gate")})

	post := httptest.NewRecorder()
	body := `{"accessMode":"website-stream-and-status"}`
	h.handleAuthGateSettings(post, settingsRequest(t, http.MethodPost, "gate", body), "gate")
	if post.Code != http.StatusOK {
		t.Fatalf("save status = %d, want 200: %s", post.Code, post.Body)
	}

	// Read back through the endpoint.
	get := httptest.NewRecorder()
	h.handleAuthGateSettings(get, settingsRequest(t, http.MethodGet, "gate", ""), "gate")
	var got AuthGateSettings
	if err := json.NewDecoder(get.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.AccessMode != AuthGateAccessWebsiteStreamAndStatus {
		t.Fatalf("read back access mode = %q, want %q", got.AccessMode, AuthGateAccessWebsiteStreamAndStatus)
	}

	// And through the path the gate itself uses on every request.
	if enforced := h.authGateSettingsFor("gate"); enforced != got {
		t.Fatalf("enforced %+v != stored %+v", enforced, got)
	}
}

// A saved policy has to survive the in-memory cache being dropped, otherwise
// the gate would silently revert to defaults after a restart.
func TestAuthGateSettings_SurviveACacheReset(t *testing.T) {
	h := hostWithDiscoveredPlugins(t, map[string][]byte{"gate": gateManifest("gate")})
	if err := h.saveAuthGateSettings("gate", AuthGateSettings{AccessMode: AuthGateAccessWebsiteAndStream}); err != nil {
		t.Fatal(err)
	}

	h.authGateSettings.Store(authGateSettingsSnapshot{})
	if got := h.authGateSettingsFor("gate"); got.AccessMode != AuthGateAccessWebsiteAndStream {
		t.Fatalf("after cache reset got %q, want %q", got.AccessMode, AuthGateAccessWebsiteAndStream)
	}
}

func TestAuthGateSettings_RejectsBadInput(t *testing.T) {
	h := hostWithDiscoveredPlugins(t, map[string][]byte{"gate": gateManifest("gate")})

	tests := []struct {
		name string
		body string
	}{
		{"malformed JSON", `{"accessMode":`},
		{"wrong value type", `{"accessMode":true}`},
		{"unsupported mode", `{"accessMode":"status-only"}`},
		{"missing mode", `{}`},
		{"unknown field", `{"accessMode":"website-only","protectEverything":true}`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			h.handleAuthGateSettings(w, settingsRequest(t, http.MethodPost, "gate", tc.body), "gate")
			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400", w.Code)
			}
			if got := h.authGateSettingsFor("gate"); got.AccessMode != AuthGateAccessWebsiteOnly {
				t.Fatalf("rejected input changed the access mode: %q", got.AccessMode)
			}
		})
	}
}

func TestAuthGateSettings_RejectsUnsupportedMethods(t *testing.T) {
	h := hostWithDiscoveredPlugins(t, map[string][]byte{"gate": gateManifest("gate")})
	for _, method := range []string{http.MethodPut, http.MethodDelete, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			w := httptest.NewRecorder()
			h.handleAuthGateSettings(w, settingsRequest(t, method, "gate", `{}`), "gate")
			if w.Code != http.StatusMethodNotAllowed {
				t.Fatalf("status = %d, want 405", w.Code)
			}
		})
	}
}

// A corrupt or unreadable stored policy must fall back to the documented
// defaults rather than to whatever a partial decode produced, and must not be
// cached, so a later good read still wins.
func TestAuthGateSettingsFor_CorruptValueFallsBackToDefaults(t *testing.T) {
	h := &Host{kv: kv.NewMemory()}
	h.authGateSettings.Store(authGateSettingsSnapshot{})
	if err := h.kv.Namespace("gate").Set(authGateSettingsKey, []byte("{not json")); err != nil {
		t.Fatal(err)
	}

	if got := h.authGateSettingsFor("gate"); got != defaultAuthGateSettings() {
		t.Fatalf("corrupt settings produced %+v, want defaults", got)
	}

	// Not cached: writing a good value takes effect without a restart.
	if err := h.saveAuthGateSettings("gate", AuthGateSettings{AccessMode: AuthGateAccessWebsiteAndStream}); err != nil {
		t.Fatal(err)
	}
	if got := h.authGateSettingsFor("gate"); got.AccessMode != AuthGateAccessWebsiteAndStream {
		t.Fatalf("after repair got %q, want %q", got.AccessMode, AuthGateAccessWebsiteAndStream)
	}
}

// DirectoryAvailable is what YP consults before pinging the directory and
// before answering /api/yp. A plugin that merely declares auth.gate, without
// being enabled as the active gate, must not switch the directory off.
// (Enforcement of the "off" case lives in the YP service and is tested there.)
func TestDirectoryAvailable_UnenabledPluginPolicyIsIgnored(t *testing.T) {
	h := hostWithDiscoveredPlugins(t, map[string][]byte{"gate": gateManifest("gate")})
	h.authGateSettings.Store(authGateSettingsSnapshot{
		slug:     "gate",
		settings: AuthGateSettings{AccessMode: AuthGateAccessWebsiteStreamAndStatus},
	})

	if !h.DirectoryAvailable() {
		t.Fatal("a discovered-but-not-enabled plugin's policy must not disable the directory")
	}
}
