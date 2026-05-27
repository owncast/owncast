package plugins

import (
	"strings"
	"testing"
)

func TestParseManifest_Valid(t *testing.T) {
	in := []byte(`{
		"api": "1",
		"name": "demo",
		"version": "1.2.3",
		"description": "a demo plugin",
		"permissions": ["chat.send", "storage.kv"]
	}`)
	m, err := ParseManifest(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Name != "demo" {
		t.Errorf("name: got %q want %q", m.Name, "demo")
	}
	if m.Version != "1.2.3" {
		t.Errorf("version: got %q want %q", m.Version, "1.2.3")
	}
	if len(m.Permissions) != 2 {
		t.Errorf("permissions: got %v want 2 entries", m.Permissions)
	}
}

func TestParseManifest_RejectsBadAPIVersion(t *testing.T) {
	in := []byte(`{"api": "99", "name": "x", "version": "1.0.0"}`)
	if _, err := ParseManifest(in); err == nil {
		t.Fatal("expected error for unsupported api version, got nil")
	} else if !strings.Contains(err.Error(), "unsupported api version") {
		t.Errorf("error mentions api version: got %v", err)
	}
}

func TestParseManifest_RequiresName(t *testing.T) {
	in := []byte(`{"api": "1", "version": "1.0.0"}`)
	if _, err := ParseManifest(in); err == nil {
		t.Fatal("expected error for missing name, got nil")
	} else if !strings.Contains(err.Error(), "name") {
		t.Errorf("error mentions name: got %v", err)
	}
}

func TestParseManifest_RequiresVersion(t *testing.T) {
	in := []byte(`{"api": "1", "name": "x"}`)
	if _, err := ParseManifest(in); err == nil {
		t.Fatal("expected error for missing version, got nil")
	} else if !strings.Contains(err.Error(), "version") {
		t.Errorf("error mentions version: got %v", err)
	}
}

func TestParseManifest_RejectsMalformedJSON(t *testing.T) {
	in := []byte(`{not valid json`)
	if _, err := ParseManifest(in); err == nil {
		t.Fatal("expected error for malformed json, got nil")
	}
}

func TestAgreesWith_HappyPath(t *testing.T) {
	sidecar := &Manifest{
		API: "1", Name: "demo", Version: "1.0.0",
		Permissions: []string{"chat.send", "storage.kv"},
	}
	runtime := &Manifest{
		API: "1", Name: "demo", Version: "1.0.0",
		Permissions: []string{"chat.send"}, // subset OK
	}
	if err := sidecar.AgreesWith(runtime); err != nil {
		t.Errorf("expected agreement, got error: %v", err)
	}
}

func TestAgreesWith_NameMismatch(t *testing.T) {
	sidecar := &Manifest{API: "1", Name: "demo", Version: "1.0.0"}
	runtime := &Manifest{API: "1", Name: "other", Version: "1.0.0"}
	err := sidecar.AgreesWith(runtime)
	if err == nil || !strings.Contains(err.Error(), "name mismatch") {
		t.Errorf("expected name mismatch error, got: %v", err)
	}
}

func TestAgreesWith_VersionMismatch(t *testing.T) {
	sidecar := &Manifest{API: "1", Name: "demo", Version: "1.0.0"}
	runtime := &Manifest{API: "1", Name: "demo", Version: "2.0.0"}
	err := sidecar.AgreesWith(runtime)
	if err == nil || !strings.Contains(err.Error(), "version mismatch") {
		t.Errorf("expected version mismatch error, got: %v", err)
	}
}

func TestAgreesWith_RuntimeExceedsDeclaredPermissions(t *testing.T) {
	sidecar := &Manifest{
		API: "1", Name: "demo", Version: "1.0.0",
		Permissions: []string{"chat.send"},
	}
	runtime := &Manifest{
		API: "1", Name: "demo", Version: "1.0.0",
		Permissions: []string{"chat.send", "events.emit"}, // events.emit not declared
	}
	err := sidecar.AgreesWith(runtime)
	if err == nil {
		t.Fatal("expected error for runtime claiming undeclared permission, got nil")
	}
	if !strings.Contains(err.Error(), "events.emit") {
		t.Errorf("error mentions the undeclared permission: got %v", err)
	}
}

func TestStrikeSystem_FilterFailuresAutoDisable(t *testing.T) {
	l := &Loaded{Manifest: &Manifest{Name: "x"}}
	for i := 0; i < FilterStrikeThreshold-1; i++ {
		if disabled := l.recordFilterFailure(); disabled {
			t.Fatalf("disabled too early at strike %d", i+1)
		}
		if l.IsDisabled() {
			t.Fatalf("IsDisabled() returned true at strike %d", i+1)
		}
	}
	if disabled := l.recordFilterFailure(); !disabled {
		t.Fatal("recordFilterFailure should have reported newly-disabled on the threshold strike")
	}
	if !l.IsDisabled() {
		t.Fatal("IsDisabled() should be true after threshold reached")
	}
	// Subsequent failures don't re-trigger the "newly disabled" signal.
	if disabled := l.recordFilterFailure(); disabled {
		t.Error("recordFilterFailure should not report newly-disabled twice")
	}
}

func TestStrikeSystem_SuccessResetsCounter(t *testing.T) {
	l := &Loaded{Manifest: &Manifest{Name: "x"}}
	// Rack up almost enough strikes to disable.
	for i := 0; i < FilterStrikeThreshold-1; i++ {
		l.recordFilterFailure()
	}
	// One success should reset the counter.
	l.recordFilterSuccess()
	if l.IsDisabled() {
		t.Fatal("success should have prevented auto-disable on the next failure batch")
	}
	// We can rack up failures again without tripping the threshold immediately.
	for i := 0; i < FilterStrikeThreshold-1; i++ {
		if disabled := l.recordFilterFailure(); disabled {
			t.Fatalf("disabled too early after reset, at strike %d", i+1)
		}
	}
}

func TestAgreesWith_SidecarMayDeclareMoreThanRuntimeUses(t *testing.T) {
	// The asymmetry: sidecar is the upper bound. Plugin author can declare
	// more than they end up using.
	sidecar := &Manifest{
		API: "1", Name: "demo", Version: "1.0.0",
		Permissions: []string{"chat.send", "storage.kv", "network.fetch"},
	}
	runtime := &Manifest{
		API: "1", Name: "demo", Version: "1.0.0",
		Permissions: []string{"chat.send"},
	}
	if err := sidecar.AgreesWith(runtime); err != nil {
		t.Errorf("sidecar should allow runtime to use subset of declared perms; got: %v", err)
	}
}

func TestParseManifest_Action_RelativeURLIsRewritten(t *testing.T) {
	m, err := ParseManifest([]byte(`{
		"api": "1", "name": "stats", "version": "1.0",
		"permissions": ["http.serve"],
		"actions": [{"title": "Dashboard", "url": "/dashboard"}]
	}`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(m.Actions) != 1 || m.Actions[0].Url != "/plugins/stats/dashboard" {
		t.Errorf("url should auto-prefix to /plugins/stats/dashboard, got %q", m.Actions[0].Url)
	}
}

func TestParseManifest_Action_BareSlashIsRewritten(t *testing.T) {
	m, err := ParseManifest([]byte(`{
		"api": "1", "name": "stats", "version": "1.0",
		"permissions": ["http.serve"],
		"actions": [{"title": "Dashboard", "url": "/"}]
	}`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if m.Actions[0].Url != "/plugins/stats/" {
		t.Errorf("url '/' should rewrite to /plugins/stats/, got %q", m.Actions[0].Url)
	}
}

func TestParseManifest_Action_ExplicitPluginPathPreserved(t *testing.T) {
	m, err := ParseManifest([]byte(`{
		"api": "1", "name": "stats", "version": "1.0",
		"permissions": ["http.serve"],
		"actions": [{"title": "Dashboard", "url": "/plugins/stats/foo"}]
	}`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if m.Actions[0].Url != "/plugins/stats/foo" {
		t.Errorf("explicit plugin path should be preserved, got %q", m.Actions[0].Url)
	}
}

func TestParseManifest_Action_ExternalURLPreserved(t *testing.T) {
	m, err := ParseManifest([]byte(`{
		"api": "1", "name": "stats", "version": "1.0",
		"actions": [{"title": "External", "url": "https://example.com/help"}]
	}`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if m.Actions[0].Url != "https://example.com/help" {
		t.Errorf("absolute URL should be preserved, got %q", m.Actions[0].Url)
	}
}

func TestParseManifest_Action_MissingHttpServePerm(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "stats", "version": "1.0",
		"actions": [{"title": "Dashboard", "url": "/dashboard"}]
	}`))
	if err == nil {
		t.Fatal("expected error: action targets self but http.serve not declared")
	}
	if !strings.Contains(err.Error(), "http.serve") {
		t.Errorf("error should mention http.serve, got: %v", err)
	}
}

func TestParseManifest_Action_PointsAtOtherPlugin(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "stats", "version": "1.0",
		"actions": [{"title": "Other", "url": "/plugins/other-plugin/page"}]
	}`))
	if err == nil {
		t.Fatal("expected error: action points at another plugin's namespace")
	}
	if !strings.Contains(err.Error(), "other plugin's namespace") {
		t.Errorf("error should call out cross-plugin URL, got: %v", err)
	}
}

func TestParseManifest_Action_TitleRequired(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "stats", "version": "1.0",
		"actions": [{"url": "https://example.com"}]
	}`))
	if err == nil {
		t.Fatal("expected error: title required")
	}
}

func TestParseManifest_Action_UrlXorHtml(t *testing.T) {
	cases := []struct {
		name, body string
	}{
		{"both", `{"title": "x", "url": "https://e.com", "html": "<p>x</p>"}`},
		{"neither", `{"title": "x"}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := ParseManifest([]byte(`{
				"api": "1", "name": "stats", "version": "1.0",
				"actions": [` + c.body + `]
			}`))
			if err == nil {
				t.Fatal("expected error: exactly one of url or html is required")
			}
		})
	}
}

func TestParseManifest_Action_HtmlOnly(t *testing.T) {
	m, err := ParseManifest([]byte(`{
		"api": "1", "name": "stats", "version": "1.0",
		"actions": [{"title": "Info", "html": "<p>hello</p>"}]
	}`))
	if err != nil {
		t.Fatalf("html-only action should be valid (no http.serve needed): %v", err)
	}
	if m.Actions[0].Html != "<p>hello</p>" {
		t.Errorf("html lost: %q", m.Actions[0].Html)
	}
}

func TestParseManifest_Network_AllowedHostsRequiredWithFetch(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "demo", "version": "1.0",
		"permissions": ["network.fetch"]
	}`))
	if err == nil {
		t.Fatal("expected error: network.fetch requires allowedHosts")
	}
	if !strings.Contains(err.Error(), "network.allowedHosts") {
		t.Errorf("error should reference network.allowedHosts: %v", err)
	}
}

func TestParseManifest_Network_AllowedHostsEmptyWithFetch(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "demo", "version": "1.0",
		"permissions": ["network.fetch"],
		"network": { "allowedHosts": [] }
	}`))
	if err == nil {
		t.Fatal("expected error: empty allowedHosts with network.fetch")
	}
}

func TestParseManifest_Network_PassesThroughHosts(t *testing.T) {
	m, err := ParseManifest([]byte(`{
		"api": "1", "name": "demo", "version": "1.0",
		"permissions": ["network.fetch"],
		"network": { "allowedHosts": ["api.discord.com", "*.weather.com"] }
	}`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(m.Network.AllowedHosts) != 2 ||
		m.Network.AllowedHosts[0] != "api.discord.com" ||
		m.Network.AllowedHosts[1] != "*.weather.com" {
		t.Errorf("allowedHosts not preserved: %v", m.Network.AllowedHosts)
	}
}

func TestParseManifest_Network_ExplicitWildcardOK(t *testing.T) {
	m, err := ParseManifest([]byte(`{
		"api": "1", "name": "demo", "version": "1.0",
		"permissions": ["network.fetch"],
		"network": { "allowedHosts": ["*"] }
	}`))
	if err != nil {
		t.Fatalf("explicit \"*\" should be valid: %v", err)
	}
	if len(m.Network.AllowedHosts) != 1 || m.Network.AllowedHosts[0] != "*" {
		t.Errorf("wildcard not preserved: %v", m.Network.AllowedHosts)
	}
}

func TestParseManifest_Network_BlankHostRejected(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "demo", "version": "1.0",
		"permissions": ["network.fetch"],
		"network": { "allowedHosts": ["api.example.com", "   "] }
	}`))
	if err == nil {
		t.Fatal("expected error: blank host in allowedHosts")
	}
}

func TestParseManifest_Network_NoFetchPermissionAllowsAnyShape(t *testing.T) {
	// A manifest with allowedHosts but without network.fetch is valid;
	// the field is just inert.
	m, err := ParseManifest([]byte(`{
		"api": "1", "name": "demo", "version": "1.0",
		"network": { "allowedHosts": ["api.discord.com"] }
	}`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(m.Network.AllowedHosts) != 1 {
		t.Errorf("allowedHosts should be preserved even without fetch perm: %v", m.Network.AllowedHosts)
	}
}
