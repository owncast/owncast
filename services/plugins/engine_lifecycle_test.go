package plugins

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// engineCacheState reports the number of cached engines and the total live
// references across them, read under the cache lock.
func engineCacheState() (entries, refs int) {
	compiledEngines.mu.Lock()
	defer compiledEngines.mu.Unlock()
	for _, e := range compiledEngines.byKey {
		refs += e.refs
	}
	return len(compiledEngines.byKey), refs
}

// TestEngineRefcountLifecycle defends the shared-engine refcount invariant:
// every live plugin instance holds one reference on its language's compiled
// engine; the engine survives while any instance lives, is closed and evicted
// when the last one closes, and a later load recompiles it from scratch.
func TestEngineRefcountLifecycle(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, sends, mu := captureEnv()
	script := `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  onChatMessage(msg) { owncast.chat.send("survivor: " + msg.body); }
});`
	p1 := loadShared(t, ctx, env, RuntimeJavaScript, "ref-one", script, []string{PermChatSend})
	p2 := loadShared(t, ctx, env, RuntimeJavaScript, "ref-two", script, []string{PermChatSend})

	if n, refs := engineCacheState(); n != 1 || refs != 2 {
		t.Fatalf("after two JS loads: %d engines with %d total refs, want 1 engine with 2 refs", n, refs)
	}

	p1.Close(ctx)
	if n, refs := engineCacheState(); n != 1 || refs != 1 {
		t.Fatalf("after closing one of two: %d engines with %d refs, want the engine to survive with 1 ref", n, refs)
	}

	// The surviving instance must still dispatch on the shared engine.
	d := NewLiveDispatcher(func() []*Loaded { return []*Loaded{p2} })
	d.Dispatch(ctx, EventChatMessageReceived, chatPayload("alice", "still here"))
	mu.Lock()
	if len(*sends) != 1 || (*sends)[0].Text != "survivor: still here" {
		mu.Unlock()
		t.Fatalf("surviving plugin failed to dispatch after sibling close: sends=%+v", *sends)
	}
	mu.Unlock()

	p2.Close(ctx)
	if n, _ := engineCacheState(); n != 0 {
		t.Fatalf("after closing the last JS plugin the engine must be evicted; %d entries remain", n)
	}

	// A later load must recompile the engine and work again.
	p3 := loadShared(t, ctx, env, RuntimeJavaScript, "ref-three", script, []string{PermChatSend})
	defer p3.Close(ctx)
	if n, refs := engineCacheState(); n != 1 || refs != 1 {
		t.Fatalf("load after eviction: %d engines with %d refs, want a freshly compiled engine with 1 ref", n, refs)
	}
}

// TestLoadFailureReleasesEngine defends the load-error path: a script that
// fails after the engine reference was taken (register() evals the author
// source, which throws) must drop its reference so the engine is not leaked.
func TestLoadFailureReleasesEngine(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, _, _ := captureEnv()
	manifestBytes, _ := json.Marshal(map[string]any{
		"api":         "1",
		"name":        "boom",
		"slug":        "boom",
		"version":     "0.1.0",
		"permissions": []string{},
	})
	_, err := loadFromBytes(ctx, env, manifestBytes, []byte(`throw new Error("boom")`), RuntimeJavaScript, "boom", nil)
	if err == nil {
		t.Fatal("expected the throwing script to fail the load")
	}
	if n, refs := engineCacheState(); n != 0 {
		t.Fatalf("failed load leaked an engine reference: %d entries, %d refs", n, refs)
	}
}

// lifecycleManifest builds manifest bytes for the Install lifecycle tests:
// fixed slug, caller-chosen version and permissions.
func lifecycleManifest(t *testing.T, version string, perms []string) []byte {
	t.Helper()
	b, err := json.Marshal(map[string]any{
		"api":         "1",
		"name":        "lifecycle-bot",
		"slug":        "lifecycle-bot",
		"version":     version,
		"permissions": perms,
	})
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// buildJSPackageBytes builds an in-memory .ocpkg carrying a plugin.js code
// entry (the JS twin of buildPackageBytes, which hardcodes plugin.wasm).
func buildJSPackageBytes(t *testing.T, manifest, script []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, data := range map[string][]byte{
		pkgManifestFilename: manifest,
		pkgJSFilename:       script,
	} {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write(data); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// TestManager_Install_ReloadsRunningInstanceOnUpdate defends the Install
// contract for updates within the approved permission set: the returned entry
// reports Enabled/Loaded like List() does, and installing v2 over a running v1
// swaps the live instance so the new code answers dispatches immediately.
func TestManager_Install_ReloadsRunningInstanceOnUpdate(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, sends, mu := captureEnv()
	mgr := NewManager(t.TempDir(), env)
	if err := mgr.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer mgr.Stop(ctx)

	script := func(tag string) []byte {
		return fmt.Appendf(nil, `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  onChatMessage(msg) { owncast.chat.send(%q + msg.body); }
});`, tag+": ")
	}
	perms := []string{PermChatSend}

	entry, err := mgr.Install(ctx, buildJSPackageBytes(t, lifecycleManifest(t, "0.1.0", perms), script("v1")))
	if err != nil {
		t.Fatalf("install v1: %v", err)
	}
	if entry.Enabled || entry.Loaded {
		t.Fatalf("fresh install must be neither enabled nor loaded: Enabled=%v Loaded=%v", entry.Enabled, entry.Loaded)
	}

	if err := mgr.Enable(ctx, "lifecycle-bot"); err != nil {
		t.Fatalf("enable: %v", err)
	}

	entry, err = mgr.Install(ctx, buildJSPackageBytes(t, lifecycleManifest(t, "0.2.0", perms), script("v2")))
	if err != nil {
		t.Fatalf("install v2: %v", err)
	}
	if !entry.Enabled || !entry.Loaded {
		t.Fatalf("update of an enabled+loaded plugin must stay enabled+loaded: Enabled=%v Loaded=%v LastError=%q",
			entry.Enabled, entry.Loaded, entry.LastError)
	}
	if entry.Version != "0.2.0" {
		t.Errorf("entry version after update: got %q, want 0.2.0", entry.Version)
	}
	if len(entry.PendingPermissions) != 0 {
		t.Errorf("same-permission update must not report pending permissions, got %v", entry.PendingPermissions)
	}

	// The running instance must be the NEW code: v2 answers the dispatch.
	d := NewLiveDispatcher(mgr.Snapshot)
	d.Dispatch(ctx, EventChatMessageReceived, chatPayload("alice", "ping"))
	mu.Lock()
	defer mu.Unlock()
	if len(*sends) != 1 || (*sends)[0].Text != "v2: ping" {
		t.Fatalf("dispatch after update: sends=%+v, want exactly one %q", *sends, "v2: ping")
	}
}

// TestManager_Install_PermissionExpansionDefersReload defends the other half
// of the Install update contract: when the new package declares a permission
// the admin never approved, the running v1 instance is left alone (it holds
// only approved permissions) and the entry reports the pending set for the
// re-approval flow.
func TestManager_Install_PermissionExpansionDefersReload(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, _, _ := captureEnv()
	mgr := NewManager(t.TempDir(), env)
	if err := mgr.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer mgr.Stop(ctx)

	script := []byte(`
const { definePlugin } = require("@owncast/plugin-sdk");
module.exports = definePlugin({ onChatMessage(msg) {} });`)

	if _, err := mgr.Install(ctx, buildJSPackageBytes(t, lifecycleManifest(t, "0.1.0", []string{}), script)); err != nil {
		t.Fatalf("install v1: %v", err)
	}
	if err := mgr.Enable(ctx, "lifecycle-bot"); err != nil {
		t.Fatalf("enable: %v", err)
	}

	entry, err := mgr.Install(ctx, buildJSPackageBytes(t, lifecycleManifest(t, "0.2.0", []string{PermChatSend}), script))
	if err != nil {
		t.Fatalf("install v2: %v", err)
	}
	if !entry.Enabled {
		t.Error("permission-expanding update must leave the plugin enabled")
	}
	if !entry.Loaded {
		t.Error("permission-expanding update must leave the old instance loaded")
	}
	if len(entry.PendingPermissions) != 1 || entry.PendingPermissions[0] != PermChatSend {
		t.Errorf("pending permissions: got %v, want [%s]", entry.PendingPermissions, PermChatSend)
	}

	// The old instance is still the one running: its manifest is v1.
	snap := mgr.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected exactly one loaded plugin, got %d", len(snap))
	}
	if got := snap[0].Manifest.Version; got != "0.1.0" {
		t.Errorf("running instance version after deferred update: got %q, want the old 0.1.0", got)
	}
}

// TestManager_Enable_ReApprovesPendingPermissionsAndSwapsInstance defends the
// re-approval contract: Enable on an already-enabled, already-loaded plugin
// with pending permissions re-captures the approval baseline, clears the
// pending set, and swaps the running old-version instance for one built from
// the updated package — releasing the old instance's engine reference. Enable
// with nothing pending stays a no-op and must not restart the instance.
func TestManager_Enable_ReApprovesPendingPermissionsAndSwapsInstance(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, _, _ := captureEnv()
	mgr := NewManager(t.TempDir(), env)
	if err := mgr.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer mgr.Stop(ctx)

	script := []byte(`
const { definePlugin } = require("@owncast/plugin-sdk");
module.exports = definePlugin({ onChatMessage(msg) {} });`)

	// v1 with no permissions, enabled; v2 expands to chat.send → deferred.
	if _, err := mgr.Install(ctx, buildJSPackageBytes(t, lifecycleManifest(t, "0.1.0", []string{}), script)); err != nil {
		t.Fatalf("install v1: %v", err)
	}
	if err := mgr.Enable(ctx, "lifecycle-bot"); err != nil {
		t.Fatalf("enable v1: %v", err)
	}
	if _, err := mgr.Install(ctx, buildJSPackageBytes(t, lifecycleManifest(t, "0.2.0", []string{PermChatSend}), script)); err != nil {
		t.Fatalf("install v2: %v", err)
	}

	// Re-approval: the admin consents to the expanded permission set.
	if err := mgr.Enable(ctx, "lifecycle-bot"); err != nil {
		t.Fatalf("re-approving Enable: %v", err)
	}

	var entry *DiscoveredEntry
	for _, e := range mgr.List() {
		if e.Slug == "lifecycle-bot" {
			entry = &e
			break
		}
	}
	if entry == nil {
		t.Fatal("lifecycle-bot missing from List() after re-approval")
	}
	if len(entry.PendingPermissions) != 0 {
		t.Errorf("re-approval must clear pending permissions, got %v", entry.PendingPermissions)
	}
	if !entry.Enabled || !entry.Loaded {
		t.Errorf("re-approved plugin must be enabled+loaded: Enabled=%v Loaded=%v LastError=%q",
			entry.Enabled, entry.Loaded, entry.LastError)
	}

	// The running instance is now the updated package.
	snap := mgr.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected exactly one loaded plugin after re-approval, got %d", len(snap))
	}
	if got := snap[0].Manifest.Version; got != "0.2.0" {
		t.Errorf("running instance version after re-approval: got %q, want 0.2.0", got)
	}

	// One live instance → one engine holding exactly one reference; the
	// swapped-out old instance must have released its ref.
	if n, refs := engineCacheState(); n != 1 || refs != 1 {
		t.Fatalf("engine cache after instance swap: %d engines with %d refs, want 1 engine with 1 ref (old ref leaked?)", n, refs)
	}

	// Enable with nothing pending stays a no-op: the same instance keeps
	// running (a restart would produce a new *Loaded).
	cur := snap[0]
	if err := mgr.Enable(ctx, "lifecycle-bot"); err != nil {
		t.Fatalf("no-op Enable: %v", err)
	}
	snap = mgr.Snapshot()
	if len(snap) != 1 || snap[0] != cur {
		t.Fatal("Enable on an enabled+loaded plugin with nothing pending must not restart the instance")
	}
	if n, refs := engineCacheState(); n != 1 || refs != 1 {
		t.Fatalf("engine cache after no-op Enable: %d engines with %d refs, want unchanged 1/1", n, refs)
	}
}

// TestManager_Install_LoadsEnabledButUnloadedPluginOnUpdate defends the last
// corner of the Install update contract: a plugin can be enabled but not
// running (an on-disk update expanded its permissions, so the host held it
// unloaded across a restart). Installing a version whose permissions are
// covered by the approved baseline must bring it back up — the admin already
// consented to everything this package asks for.
func TestManager_Install_LoadsEnabledButUnloadedPluginOnUpdate(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, _, _ := captureEnv()
	dir := t.TempDir()
	mgr := NewManager(dir, env)
	if err := mgr.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}

	script := []byte(`
const { definePlugin } = require("@owncast/plugin-sdk");
module.exports = definePlugin({ onChatMessage(msg) {} });`)
	perms := []string{PermChatSend}

	if _, err := mgr.Install(ctx, buildJSPackageBytes(t, lifecycleManifest(t, "0.1.0", perms), script)); err != nil {
		t.Fatalf("install v1: %v", err)
	}
	if err := mgr.Enable(ctx, "lifecycle-bot"); err != nil {
		t.Fatalf("enable: %v", err)
	}
	mgr.Stop(ctx)

	// The author drops a permission-expanding v2 on disk while the host is
	// down; the restarted manager discovers it, sees an unapproved
	// permission, and keeps the plugin unloaded.
	v2 := buildJSPackageBytes(t, lifecycleManifest(t, "0.2.0", []string{PermChatSend, PermStorageKV}), script)
	if err := os.WriteFile(filepath.Join(dir, "lifecycle-bot.ocpkg"), v2, 0o600); err != nil {
		t.Fatalf("drop v2 on disk: %v", err)
	}
	mgr = NewManager(dir, env)
	if err := mgr.Start(ctx); err != nil {
		t.Fatalf("restart: %v", err)
	}
	defer mgr.Stop(ctx)
	if len(mgr.Snapshot()) != 0 {
		t.Fatal("precondition: the expanded-permission v2 must be held unloaded")
	}

	// The admin installs v3, back inside the approved baseline. The plugin
	// is enabled and everything v3 asks for is approved, so the update must
	// leave it running — not stranded at "enabled, not loaded" until the
	// next restart.
	entry, err := mgr.Install(ctx, buildJSPackageBytes(t, lifecycleManifest(t, "0.3.0", perms), script))
	if err != nil {
		t.Fatalf("install v3: %v", err)
	}
	if !entry.Enabled || !entry.Loaded {
		t.Fatalf("baseline-conforming update of an enabled plugin must end up loaded: Enabled=%v Loaded=%v LastError=%q",
			entry.Enabled, entry.Loaded, entry.LastError)
	}
	if len(entry.PendingPermissions) != 0 {
		t.Errorf("pending permissions after conforming update: got %v, want none", entry.PendingPermissions)
	}
	snap := mgr.Snapshot()
	if len(snap) != 1 || snap[0].Manifest.Version != "0.3.0" {
		t.Fatalf("running instance after conforming update: got %d loaded (want 1) with version %q (want 0.3.0)",
			len(snap), func() string {
				if len(snap) > 0 {
					return snap[0].Manifest.Version
				}
				return ""
			}())
	}
}
