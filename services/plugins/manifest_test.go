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
	if m.DisplayName != "demo" {
		t.Errorf("name: got %q want %q", m.DisplayName, "demo")
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
		API: "1", DisplayName: "demo", Version: "1.0.0",
		Permissions: []string{"chat.send", "storage.kv"},
	}
	runtime := &Manifest{
		API: "1", DisplayName: "demo", Version: "1.0.0",
		Permissions: []string{"chat.send"}, // subset OK
	}
	if err := sidecar.AgreesWith(runtime); err != nil {
		t.Errorf("expected agreement, got error: %v", err)
	}
}

func TestAgreesWith_SlugMismatch(t *testing.T) {
	// AgreesWith identity-checks on Slug. When the runtime side ships no
	// explicit slug, the host derives one from the runtime DisplayName the
	// same way ParseManifest does on the sidecar side, so a clean string
	// like "demo" vs "other" still trips the mismatch path.
	sidecar := &Manifest{API: "1", DisplayName: "demo", Slug: "demo", Version: "1.0.0"}
	runtime := &Manifest{API: "1", DisplayName: "other", Version: "1.0.0"}
	err := sidecar.AgreesWith(runtime)
	if err == nil || !strings.Contains(err.Error(), "slug mismatch") {
		t.Errorf("expected slug mismatch error, got: %v", err)
	}
}

func TestAgreesWith_VersionMismatch(t *testing.T) {
	sidecar := &Manifest{API: "1", DisplayName: "demo", Version: "1.0.0"}
	runtime := &Manifest{API: "1", DisplayName: "demo", Version: "2.0.0"}
	err := sidecar.AgreesWith(runtime)
	if err == nil || !strings.Contains(err.Error(), "version mismatch") {
		t.Errorf("expected version mismatch error, got: %v", err)
	}
}

func TestAgreesWith_RuntimeExceedsDeclaredPermissions(t *testing.T) {
	sidecar := &Manifest{
		API: "1", DisplayName: "demo", Version: "1.0.0",
		Permissions: []string{"chat.send"},
	}
	runtime := &Manifest{
		API: "1", DisplayName: "demo", Version: "1.0.0",
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
	l := &Loaded{Manifest: &Manifest{DisplayName: "x"}}
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
	l := &Loaded{Manifest: &Manifest{DisplayName: "x"}}
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
		API: "1", DisplayName: "demo", Version: "1.0.0",
		Permissions: []string{"chat.send", "storage.kv", "network.fetch"},
	}
	runtime := &Manifest{
		API: "1", DisplayName: "demo", Version: "1.0.0",
		Permissions: []string{"chat.send"},
	}
	if err := sidecar.AgreesWith(runtime); err != nil {
		t.Errorf("sidecar should allow runtime to use subset of declared perms; got: %v", err)
	}
}

func TestParseManifest_Action_RelativeURLIsRewritten(t *testing.T) {
	m, err := ParseManifest([]byte(`{
		"api": "1", "name": "stats", "version": "1.0",
		"permissions": ["http.serve", "ui.modify"],
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
		"permissions": ["http.serve", "ui.modify"],
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
		"permissions": ["http.serve", "ui.modify"],
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
		"permissions": ["ui.modify"],
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
		"permissions": ["ui.modify"],
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
		"permissions": ["ui.modify"],
		"actions": [{"title": "Other", "url": "/plugins/other-plugin/page"}]
	}`))
	if err == nil {
		t.Fatal("expected error: action points at another plugin's namespace")
	}
	if !strings.Contains(err.Error(), "other plugin's namespace") {
		t.Errorf("error should call out cross-plugin URL, got: %v", err)
	}
}

func TestParseManifest_Action_IconRelativeIsRewritten(t *testing.T) {
	m, err := ParseManifest([]byte(`{
		"api": "1", "name": "stats", "version": "1.0",
		"permissions": ["http.serve", "ui.modify"],
		"actions": [{"title": "Dashboard", "url": "/dashboard", "icon": "/star.png"}]
	}`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if m.Actions[0].Icon != "/plugins/stats/star.png" {
		t.Errorf("icon should auto-prefix to /plugins/stats/star.png, got %q", m.Actions[0].Icon)
	}
}

func TestParseManifest_Action_IconExternalPreserved(t *testing.T) {
	m, err := ParseManifest([]byte(`{
		"api": "1", "name": "stats", "version": "1.0",
		"permissions": ["ui.modify"],
		"actions": [{"title": "External", "url": "https://example.com", "icon": "https://cdn.example.com/star.png"}]
	}`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if m.Actions[0].Icon != "https://cdn.example.com/star.png" {
		t.Errorf("absolute icon URL should be preserved, got %q", m.Actions[0].Icon)
	}
}

func TestParseManifest_Action_IconMissingHttpServePerm(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "stats", "version": "1.0",
		"permissions": ["ui.modify"],
		"actions": [{"title": "External", "url": "https://example.com", "icon": "/star.png"}]
	}`))
	if err == nil {
		t.Fatal("expected error: icon targets self but http.serve not declared")
	}
	if !strings.Contains(err.Error(), "http.serve") {
		t.Errorf("error should mention http.serve, got: %v", err)
	}
}

func TestParseManifest_Action_IconPointsAtOtherPlugin(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "stats", "version": "1.0",
		"permissions": ["ui.modify"],
		"actions": [{"title": "External", "url": "https://example.com", "icon": "/plugins/other-plugin/star.png"}]
	}`))
	if err == nil {
		t.Fatal("expected error: icon points at another plugin's namespace")
	}
	if !strings.Contains(err.Error(), "other plugin's namespace") {
		t.Errorf("error should call out cross-plugin icon, got: %v", err)
	}
}

func TestParseManifest_Action_TitleRequired(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "stats", "version": "1.0",
		"permissions": ["ui.modify"],
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
				"permissions": ["ui.modify"],
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
		"permissions": ["ui.modify"],
		"actions": [{"title": "Info", "html": "<p>hello</p>"}]
	}`))
	if err != nil {
		t.Fatalf("html-only action should be valid (no http.serve needed): %v", err)
	}
	if m.Actions[0].Html != "<p>hello</p>" {
		t.Errorf("html lost: %q", m.Actions[0].Html)
	}
}

func TestParseManifest_UIModify_RequiredForActions(t *testing.T) {
	// A plugin that contributes a viewer action button must declare
	// "ui.modify" so the surface is visible in the permission list.
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "stats", "version": "1.0",
		"actions": [{"title": "x", "url": "https://example.com"}]
	}`))
	if err == nil {
		t.Fatal("expected error: actions without ui.modify permission")
	}
	if !strings.Contains(err.Error(), "ui.modify") {
		t.Errorf("error should mention ui.modify, got: %v", err)
	}
}

func TestParseManifest_AdminPages_NoUIModifyRequired(t *testing.T) {
	// Admin pages live in their own /plugins/<name>/admin namespace and
	// don't alter Owncast's existing UI, so they're baseline plugin
	// functionality, no ui.modify permission required.
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "stats", "version": "1.0",
		"permissions": ["http.serve"],
		"admin": {"pages": [{"title": "x", "path": "/admin"}]}
	}`))
	if err != nil {
		t.Errorf("admin.pages should not require ui.modify: %v", err)
	}
}

func TestParseManifest_UIModify_NotRequiredWhenNoUI(t *testing.T) {
	// A plugin that doesn't declare any UI surfaces should not be forced
	// to declare ui.modify — the permission is opt-in only for plugins
	// that contribute to the host's UI.
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "headless", "version": "1.0"
	}`))
	if err != nil {
		t.Errorf("headless plugin with no UI surfaces should be valid: %v", err)
	}
}

func TestParseManifest_Styles_RewritesRelativePaths(t *testing.T) {
	m, err := ParseManifest([]byte(`{
		"api": "1", "name": "styled", "version": "1.0",
		"permissions": ["ui.modify", "http.serve"],
		"styles": ["theme.css", "/extra.css", "/plugins/styled/another.css"]
	}`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	want := []string{"/plugins/styled/theme.css", "/plugins/styled/extra.css", "/plugins/styled/another.css"}
	if len(m.Styles) != len(want) {
		t.Fatalf("styles count: got %v want %v", m.Styles, want)
	}
	for i, w := range want {
		if m.Styles[i] != w {
			t.Errorf("styles[%d] = %q, want %q", i, m.Styles[i], w)
		}
	}
}

func TestParseManifest_Styles_RequiresUIModify(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "styled", "version": "1.0",
		"permissions": ["http.serve"],
		"styles": ["theme.css"]
	}`))
	if err == nil {
		t.Fatal("expected error when styles is set without ui.modify")
	}
	if !strings.Contains(err.Error(), "ui.modify") {
		t.Errorf("error should mention ui.modify, got: %v", err)
	}
}

func TestParseManifest_Styles_DoesNotRequireHttpServe(t *testing.T) {
	// styles are inlined from assets/ into /api/config — http.serve is not needed.
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "styled", "version": "1.0",
		"permissions": ["ui.modify"],
		"styles": ["theme.css"]
	}`))
	if err != nil {
		t.Fatalf("expected no error for styles without http.serve, got: %v", err)
	}
}

func TestParseManifest_Styles_RejectsCrossPluginPath(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "styled", "version": "1.0",
		"permissions": ["ui.modify", "http.serve"],
		"styles": ["/plugins/some-other-plugin/theme.css"]
	}`))
	if err == nil {
		t.Fatal("expected error: cross-plugin style path")
	}
	if !strings.Contains(err.Error(), "another plugin's namespace") {
		t.Errorf("error should call out cross-plugin path, got: %v", err)
	}
}

func TestParseManifest_Styles_RejectsAbsoluteURL(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "styled", "version": "1.0",
		"permissions": ["ui.modify", "http.serve"],
		"styles": ["https://fonts.example.com/theme.css"]
	}`))
	if err == nil {
		t.Fatal("expected error: absolute URLs rejected for styles")
	}
	if !strings.Contains(err.Error(), "absolute URL") {
		t.Errorf("error should call out absolute URL, got: %v", err)
	}
}

func TestParseManifest_Styles_RejectsNonCssExtension(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "styled", "version": "1.0",
		"permissions": ["ui.modify", "http.serve"],
		"styles": ["theme.scss"]
	}`))
	if err == nil {
		t.Fatal("expected error: only .css extension allowed")
	}
	if !strings.Contains(err.Error(), ".css") {
		t.Errorf("error should call out .css, got: %v", err)
	}
}

func TestParseManifest_Scripts_RewritesRelativePaths(t *testing.T) {
	m, err := ParseManifest([]byte(`{
		"api": "1", "name": "scripted", "version": "1.0",
		"permissions": ["ui.modify", "http.serve"],
		"scripts": ["client.js", "/extra.js", "/plugins/scripted/another.js"]
	}`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	want := []string{"/plugins/scripted/client.js", "/plugins/scripted/extra.js", "/plugins/scripted/another.js"}
	if len(m.Scripts) != len(want) {
		t.Fatalf("scripts count: got %v want %v", m.Scripts, want)
	}
	for i, w := range want {
		if m.Scripts[i] != w {
			t.Errorf("scripts[%d] = %q, want %q", i, m.Scripts[i], w)
		}
	}
}

func TestParseManifest_Scripts_RequiresUIModify(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "scripted", "version": "1.0",
		"permissions": ["http.serve"],
		"scripts": ["client.js"]
	}`))
	if err == nil {
		t.Fatal("expected error when scripts is set without ui.modify")
	}
	if !strings.Contains(err.Error(), "ui.modify") {
		t.Errorf("error should mention ui.modify, got: %v", err)
	}
}

func TestParseManifest_Scripts_DoesNotRequireHttpServe(t *testing.T) {
	// scripts are inlined from assets/ into /customjavascript — http.serve is not needed.
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "scripted", "version": "1.0",
		"permissions": ["ui.modify"],
		"scripts": ["client.js"]
	}`))
	if err != nil {
		t.Fatalf("expected no error for scripts without http.serve, got: %v", err)
	}
}

func TestParseManifest_Scripts_RejectsCrossPluginPath(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "scripted", "version": "1.0",
		"permissions": ["ui.modify", "http.serve"],
		"scripts": ["/plugins/some-other-plugin/client.js"]
	}`))
	if err == nil {
		t.Fatal("expected error: cross-plugin script path")
	}
	if !strings.Contains(err.Error(), "another plugin's namespace") {
		t.Errorf("error should call out cross-plugin path, got: %v", err)
	}
}

func TestParseManifest_Scripts_RejectsAbsoluteURL(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "scripted", "version": "1.0",
		"permissions": ["ui.modify", "http.serve"],
		"scripts": ["https://cdn.example.com/client.js"]
	}`))
	if err == nil {
		t.Fatal("expected error: absolute URLs rejected for scripts")
	}
	if !strings.Contains(err.Error(), "absolute URL") {
		t.Errorf("error should call out absolute URL, got: %v", err)
	}
}

func TestParseManifest_Scripts_RejectsNonJsExtension(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "scripted", "version": "1.0",
		"permissions": ["ui.modify", "http.serve"],
		"scripts": ["client.ts"]
	}`))
	if err == nil {
		t.Fatal("expected error: only .js extension allowed")
	}
	if !strings.Contains(err.Error(), ".js") {
		t.Errorf("error should call out .js, got: %v", err)
	}
}

func TestParseManifest_ExtraPageContent_RewritesRelativePath(t *testing.T) {
	m, err := ParseManifest([]byte(`{
		"api": "1", "name": "page", "version": "1.0",
		"permissions": ["ui.modify"],
		"extraPageContent": {"slug": "test-slot", "content": "content.html"}
	}`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if m.ExtraPageContent == nil {
		t.Fatal("extraPageContent should not be nil")
	}
	if m.ExtraPageContent.Content != "/plugins/page/content.html" {
		t.Errorf("extraPageContent.content = %q, want /plugins/page/content.html", m.ExtraPageContent.Content)
	}
}

func TestParseManifest_ExtraPageContent_DoesNotRequireHttpServe(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "page", "version": "1.0",
		"permissions": ["ui.modify"],
		"extraPageContent": {"slug": "test-slot", "content": "content.html"}
	}`))
	if err != nil {
		t.Errorf("extraPageContent should not require http.serve: %v", err)
	}
}

func TestParseManifest_ExtraPageContent_RequiresUIModify(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "page", "version": "1.0",
		"extraPageContent": {"slug": "test-slot", "content": "content.html"}
	}`))
	if err == nil {
		t.Fatal("expected error when extraPageContent is set without ui.modify")
	}
	if !strings.Contains(err.Error(), "ui.modify") {
		t.Errorf("error should mention ui.modify, got: %v", err)
	}
}

func TestParseManifest_ExtraPageContent_RejectsCrossPluginPath(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "page", "version": "1.0",
		"permissions": ["ui.modify"],
		"extraPageContent": {"slug": "test-slot", "content": "/plugins/some-other-plugin/content.html"}
	}`))
	if err == nil {
		t.Fatal("expected error: cross-plugin extraPageContent path")
	}
	if !strings.Contains(err.Error(), "another plugin's namespace") {
		t.Errorf("error should call out cross-plugin path, got: %v", err)
	}
}

func TestParseManifest_ExtraPageContent_RejectsAbsoluteURL(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "page", "version": "1.0",
		"permissions": ["ui.modify"],
		"extraPageContent": {"slug": "test-slot", "content": "https://cdn.example.com/content.html"}
	}`))
	if err == nil {
		t.Fatal("expected error: absolute URLs rejected for extraPageContent")
	}
	if !strings.Contains(err.Error(), "absolute URL") {
		t.Errorf("error should call out absolute URL, got: %v", err)
	}
}

func TestParseManifest_ExtraPageContent_RejectsNonHtmlExtension(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "page", "version": "1.0",
		"permissions": ["ui.modify"],
		"extraPageContent": {"slug": "test-slot", "content": "content.htm"}
	}`))
	if err == nil {
		t.Fatal("expected error: only .html extension allowed")
	}
	if !strings.Contains(err.Error(), ".html") {
		t.Errorf("error should call out .html, got: %v", err)
	}
}

func TestParseManifest_Tabs_RewritesContentPaths(t *testing.T) {
	m, err := ParseManifest([]byte(`{
		"api": "1", "name": "tabbed", "version": "1.0",
		"permissions": ["ui.modify"],
		"tabs": [
			{ "title": "Music", "slug": "music", "content": "music.html" },
			{ "title": "Photos", "slug": "photos", "content": "/photos.html" }
		]
	}`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	want := []Tab{
		{Title: "Music", Slug: "music", Content: "/plugins/tabbed/music.html"},
		{Title: "Photos", Slug: "photos", Content: "/plugins/tabbed/photos.html"},
	}
	if len(m.Tabs) != len(want) {
		t.Fatalf("tabs count: got %d want %d", len(m.Tabs), len(want))
	}
	for i, w := range want {
		if m.Tabs[i] != w {
			t.Errorf("tabs[%d] = %+v, want %+v", i, m.Tabs[i], w)
		}
	}
}

func TestParseManifest_Tabs_RequiresUIModify(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "tabbed", "version": "1.0",
		"tabs": [{ "title": "Music", "slug": "music", "content": "music.html" }]
	}`))
	if err == nil {
		t.Fatal("expected error when tabs is set without ui.modify")
	}
	if !strings.Contains(err.Error(), "ui.modify") {
		t.Errorf("error should mention ui.modify, got: %v", err)
	}
}

func TestParseManifest_Tabs_RequiresTitle(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "tabbed", "version": "1.0",
		"permissions": ["ui.modify"],
		"tabs": [{ "title": "", "slug": "music", "content": "music.html" }]
	}`))
	if err == nil {
		t.Fatal("expected error when tab.title is empty")
	}
	if !strings.Contains(err.Error(), "title is required") {
		t.Errorf("error should call out title, got: %v", err)
	}
}

func TestParseManifest_Tabs_RejectsDuplicateTitles(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "tabbed", "version": "1.0",
		"permissions": ["ui.modify"],
		"tabs": [
			{ "title": "Music", "slug": "music-a", "content": "a.html" },
			{ "title": "Music", "slug": "music-b", "content": "b.html" }
		]
	}`))
	if err == nil {
		t.Fatal("expected error: duplicate tab titles")
	}
	if !strings.Contains(err.Error(), "duplicate") {
		t.Errorf("error should call out duplicate, got: %v", err)
	}
}

func TestParseManifest_Tabs_RejectsCrossPluginPath(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "tabbed", "version": "1.0",
		"permissions": ["ui.modify"],
		"tabs": [{ "title": "Other", "slug": "other", "content": "/plugins/some-other-plugin/x.html" }]
	}`))
	if err == nil {
		t.Fatal("expected error: cross-plugin tab content path")
	}
	if !strings.Contains(err.Error(), "another plugin's namespace") {
		t.Errorf("error should call out cross-plugin path, got: %v", err)
	}
}

func TestParseManifest_Tabs_RejectsNonHtmlExtension(t *testing.T) {
	_, err := ParseManifest([]byte(`{
		"api": "1", "name": "tabbed", "version": "1.0",
		"permissions": ["ui.modify"],
		"tabs": [{ "title": "Bad", "slug": "bad", "content": "music.txt" }]
	}`))
	if err == nil {
		t.Fatal("expected error: only .html extension allowed for tabs")
	}
	if !strings.Contains(err.Error(), ".html") {
		t.Errorf("error should call out .html, got: %v", err)
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

func TestRequireChatFilterPermission_RejectsWhenMissing(t *testing.T) {
	// A plugin that subscribes to filterChatMessage at register-time must
	// declare the chat.filter permission in its manifest. The host refuses
	// to load otherwise so the admin can't be surprised by a plugin that
	// silently starts rewriting chat.
	m := &Manifest{DisplayName: "stealth", Permissions: nil}
	subs := Subscriptions{
		Filter: []Subscription{{Event: EventChatMessageReceived, Priority: 100}},
	}
	err := requireChatFilterPermission(m, subs)
	if err == nil {
		t.Fatal("expected error when filterChatMessage is subscribed without chat.filter")
	}
	if !strings.Contains(err.Error(), PermChatFilter) {
		t.Errorf("error should mention %q, got: %v", PermChatFilter, err)
	}
}

func TestRequireChatFilterPermission_AcceptsWhenDeclared(t *testing.T) {
	m := &Manifest{DisplayName: "honest", Permissions: []string{PermChatFilter}}
	subs := Subscriptions{
		Filter: []Subscription{{Event: EventChatMessageReceived, Priority: 100}},
	}
	if err := requireChatFilterPermission(m, subs); err != nil {
		t.Errorf("declared chat.filter should accept: %v", err)
	}
}

func TestRequireChatFilterPermission_NoOpWhenNoFilterSubscription(t *testing.T) {
	// A plugin that doesn't subscribe to filterChatMessage at all should
	// pass through, regardless of whether it declares chat.filter.
	m := &Manifest{DisplayName: "passive", Permissions: nil}
	subs := Subscriptions{
		Filter: []Subscription{{Event: "some-other-event", Priority: 100}},
	}
	if err := requireChatFilterPermission(m, subs); err != nil {
		t.Errorf("non-chat filter subscription should not require chat.filter: %v", err)
	}
}
