package plugins

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// makePluginFiles drops a sidecar wasm + manifest pair into dir. The manifest
// name is what the manager keys discovered/enabled state by.
func makePluginFiles(t *testing.T, dir, name string, wasmBytes []byte) {
	t.Helper()
	// Version must match what the bundled example's register() returns —
	// the host enforces manifest/runtime agreement at load time.
	manifest := map[string]any{
		"api":         "1",
		"name":        name,
		"version":     "0.1.0",
		"description": name + " for tests",
		"permissions": []string{},
	}
	mb, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name+".manifest.json"), mb, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name+".wasm"), wasmBytes, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestManager_DiscoversWithoutLoading(t *testing.T) {
	wasmPath := findExampleWasm(t)
	wasmBytes, err := os.ReadFile(wasmPath)
	if err != nil {
		t.Fatalf("read example wasm: %v", err)
	}

	dir := t.TempDir()
	makePluginFiles(t, dir, "hello-world", wasmBytes)

	mgr := NewManager(dir, &HostEnv{})
	if err := mgr.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer mgr.Stop(context.Background())

	entries := mgr.List()
	if len(entries) != 1 {
		t.Fatalf("expected 1 discovered, got %d", len(entries))
	}
	if entries[0].Slug != "hello-world" {
		t.Errorf("slug: got %q want hello-world", entries[0].Slug)
	}
	if entries[0].Loaded {
		t.Error("plugin should not be loaded — admin never enabled it")
	}
	if entries[0].Enabled {
		t.Error("plugin should not be enabled — admin never enabled it")
	}
	if len(mgr.Snapshot()) != 0 {
		t.Errorf("snapshot should be empty for un-enabled plugin, got %d", len(mgr.Snapshot()))
	}
}

func TestManager_EnableLoadsAndPersists(t *testing.T) {
	wasmPath := findExampleWasm(t)
	wasmBytes, _ := os.ReadFile(wasmPath)

	dir := t.TempDir()
	makePluginFiles(t, dir, "hello-world", wasmBytes)

	ctx := context.Background()
	mgr := NewManager(dir, &HostEnv{})
	if err := mgr.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}

	if err := mgr.Enable(ctx, "hello-world"); err != nil {
		t.Fatalf("enable: %v", err)
	}
	if len(mgr.Snapshot()) != 1 {
		t.Errorf("snapshot count after enable: got %d want 1", len(mgr.Snapshot()))
	}

	// Persistence: stop, start a fresh manager, plugin should auto-load.
	mgr.Stop(ctx)

	mgr2 := NewManager(dir, &HostEnv{})
	if err := mgr2.Start(ctx); err != nil {
		t.Fatalf("restart: %v", err)
	}
	defer mgr2.Stop(ctx)

	if len(mgr2.Snapshot()) != 1 {
		t.Errorf("snapshot count after restart: got %d want 1 (enabled set should persist)",
			len(mgr2.Snapshot()))
	}
}

func TestManager_DisableUnloadsButKeepsDiscovered(t *testing.T) {
	wasmPath := findExampleWasm(t)
	wasmBytes, _ := os.ReadFile(wasmPath)

	dir := t.TempDir()
	makePluginFiles(t, dir, "hello-world", wasmBytes)

	ctx := context.Background()
	mgr := NewManager(dir, &HostEnv{})
	if err := mgr.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer mgr.Stop(ctx)

	if err := mgr.Enable(ctx, "hello-world"); err != nil {
		t.Fatalf("enable: %v", err)
	}
	if err := mgr.Disable(ctx, "hello-world"); err != nil {
		t.Fatalf("disable: %v", err)
	}
	if len(mgr.Snapshot()) != 0 {
		t.Errorf("snapshot after disable: got %d want 0", len(mgr.Snapshot()))
	}
	entries := mgr.List()
	if len(entries) != 1 {
		t.Errorf("discovered list after disable: got %d want 1", len(entries))
	}
	if entries[0].Enabled {
		t.Error("entry should not be marked enabled after disable")
	}
}

func TestManager_ScanRemovesDeletedFiles(t *testing.T) {
	wasmPath := findExampleWasm(t)
	wasmBytes, _ := os.ReadFile(wasmPath)

	dir := t.TempDir()
	makePluginFiles(t, dir, "hello-world", wasmBytes)

	ctx := context.Background()
	mgr := NewManager(dir, &HostEnv{})
	mgr.scanInterval = 20 * time.Millisecond
	if err := mgr.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer mgr.Stop(ctx)

	if len(mgr.List()) != 1 {
		t.Fatalf("setup: expected 1 discovered, got %d", len(mgr.List()))
	}

	// Delete both files; the next scan should drop the entry.
	if err := os.Remove(filepath.Join(dir, "hello-world.wasm")); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(dir, "hello-world.manifest.json")); err != nil {
		t.Fatal(err)
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if len(mgr.List()) == 0 {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Errorf("expected scan to drop deleted plugin within 2s, still have %d", len(mgr.List()))
}

func TestManager_EnableUnknownPluginErrors(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManager(dir, &HostEnv{})
	if err := mgr.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer mgr.Stop(context.Background())

	err := mgr.Enable(context.Background(), "does-not-exist")
	if err == nil {
		t.Error("expected error enabling unknown plugin")
	}
}

// rewriteManifestPerms updates the sidecar manifest's permissions field
// without touching name or version, simulating an author dropping a new
// .ocpkg whose manifest declares additional permissions.
func rewriteManifestPerms(t *testing.T, dir, name string, perms []string) {
	t.Helper()
	manifest := map[string]any{
		"api":         "1",
		"name":        name,
		"version":     "0.1.0",
		"description": name + " for tests",
		"permissions": perms,
	}
	mb, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name+".manifest.json"), mb, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestPendingPermissions(t *testing.T) {
	cases := []struct {
		name     string
		manifest []string
		approved []string
		want     []string
	}{
		{"empty manifest", nil, []string{"a"}, nil},
		{"all approved", []string{"a", "b"}, []string{"a", "b"}, nil},
		{"new perm pending", []string{"a", "b"}, []string{"a"}, []string{"b"}},
		{"approved superset", []string{"a"}, []string{"a", "b"}, nil},
		{"no approval baseline", []string{"a", "b"}, nil, []string{"a", "b"}},
		{"unsorted manifest order", []string{"c", "a"}, []string{"a"}, []string{"c"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := pendingPermissions(tc.manifest, tc.approved)
			if !stringSliceEqual(got, tc.want) {
				t.Errorf("pendingPermissions(%v, %v) = %v, want %v",
					tc.manifest, tc.approved, got, tc.want)
			}
		})
	}
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestManager_EnableCapturesApprovalBaseline(t *testing.T) {
	wasmPath := findExampleWasm(t)
	wasmBytes, _ := os.ReadFile(wasmPath)
	dir := t.TempDir()
	makePluginFiles(t, dir, "hello-world", wasmBytes)

	ctx := context.Background()
	mgr := NewManager(dir, &HostEnv{})
	if err := mgr.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer mgr.Stop(ctx)

	if err := mgr.Enable(ctx, "hello-world"); err != nil {
		t.Fatalf("enable: %v", err)
	}

	// The on-disk .enabled.json should now contain an empty approved-perm
	// list for hello-world (the bundled manifest has no permissions).
	raw, err := os.ReadFile(filepath.Join(dir, ".enabled.json"))
	if err != nil {
		t.Fatalf("read enabled file: %v", err)
	}
	var f enabledFileContents
	if err := json.Unmarshal(raw, &f); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, ok := f.ApprovedPermissions["hello-world"]; !ok {
		t.Errorf("expected approved perm entry for hello-world after enable; got %v", f.ApprovedPermissions)
	}
}

func TestManager_PermExpansionAutoDisables(t *testing.T) {
	wasmPath := findExampleWasm(t)
	wasmBytes, _ := os.ReadFile(wasmPath)
	dir := t.TempDir()
	makePluginFiles(t, dir, "hello-world", wasmBytes)

	ctx := context.Background()
	mgr := NewManager(dir, &HostEnv{})
	if err := mgr.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := mgr.Enable(ctx, "hello-world"); err != nil {
		t.Fatalf("enable: %v", err)
	}
	mgr.Stop(ctx)

	// Author updates the manifest on disk to declare a permission the
	// admin never approved. The wasm itself stays put (so AgreesWith
	// still passes: empty runtime perms is a subset of any manifest set).
	rewriteManifestPerms(t, dir, "hello-world", []string{"chat.send"})

	mgr2 := NewManager(dir, &HostEnv{})
	if err := mgr2.Start(ctx); err != nil {
		t.Fatalf("restart: %v", err)
	}
	defer mgr2.Stop(ctx)

	entries := mgr2.List()
	if len(entries) != 1 {
		t.Fatalf("expected 1 discovered, got %d", len(entries))
	}
	entry := entries[0]
	if !entry.Enabled {
		t.Error("admin intent should survive: Enabled=true")
	}
	if entry.Loaded {
		t.Error("plugin must not load while permissions are pending approval")
	}
	if !stringSliceEqual(entry.PendingPermissions, []string{"chat.send"}) {
		t.Errorf("PendingPermissions = %v, want [chat.send]", entry.PendingPermissions)
	}
	if len(mgr2.Snapshot()) != 0 {
		t.Errorf("snapshot must be empty: plugin is pending approval, got %d loaded", len(mgr2.Snapshot()))
	}
}

func TestManager_ReEnableApprovesAndLoads(t *testing.T) {
	wasmPath := findExampleWasm(t)
	wasmBytes, _ := os.ReadFile(wasmPath)
	dir := t.TempDir()
	makePluginFiles(t, dir, "hello-world", wasmBytes)

	ctx := context.Background()
	mgr := NewManager(dir, &HostEnv{})
	if err := mgr.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := mgr.Enable(ctx, "hello-world"); err != nil {
		t.Fatalf("enable: %v", err)
	}
	mgr.Stop(ctx)

	rewriteManifestPerms(t, dir, "hello-world", []string{"chat.send"})

	mgr2 := NewManager(dir, &HostEnv{})
	if err := mgr2.Start(ctx); err != nil {
		t.Fatalf("restart: %v", err)
	}
	defer mgr2.Stop(ctx)

	// Sanity: still pending.
	if entries := mgr2.List(); entries[0].Loaded {
		t.Fatal("plugin should not be loaded before re-enable")
	}

	// Admin re-enables; this captures the expanded perm set as the new
	// approved baseline, clears PendingPermissions, and loads the plugin.
	if err := mgr2.Enable(ctx, "hello-world"); err != nil {
		t.Fatalf("re-enable: %v", err)
	}
	entries := mgr2.List()
	if len(entries[0].PendingPermissions) != 0 {
		t.Errorf("PendingPermissions should be cleared after re-enable, got %v",
			entries[0].PendingPermissions)
	}
	if !entries[0].Loaded {
		t.Error("plugin should be loaded after re-enable")
	}
}

func TestManager_ExistingInstallGetsApprovalBaseline(t *testing.T) {
	// An installation that predates the approved-permissions field: the
	// enabled file lists the plugin but has no approved perms map. On
	// next Start, the manager should silently capture the current
	// manifest perms as the baseline rather than treat them as pending.
	wasmPath := findExampleWasm(t)
	wasmBytes, _ := os.ReadFile(wasmPath)
	dir := t.TempDir()
	rewriteManifestPerms(t, dir, "hello-world", []string{"chat.send"})
	if err := os.WriteFile(filepath.Join(dir, "hello-world.wasm"), wasmBytes, 0o644); err != nil {
		t.Fatal(err)
	}
	// Legacy enabled.json: enabled set only, no approvals.
	legacy := []byte(`{"enabled":["hello-world"]}`)
	if err := os.WriteFile(filepath.Join(dir, ".enabled.json"), legacy, 0o644); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	mgr := NewManager(dir, &HostEnv{})
	if err := mgr.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer mgr.Stop(ctx)

	entries := mgr.List()
	if len(entries) != 1 {
		t.Fatalf("expected 1 discovered, got %d", len(entries))
	}
	if len(entries[0].PendingPermissions) != 0 {
		t.Errorf("existing install should not flag any pending perms, got %v",
			entries[0].PendingPermissions)
	}
	if !entries[0].Loaded {
		t.Error("existing install should load on next Start without re-approval")
	}
}

// buildPackageBytes builds an in-memory .ocpkg with the given manifest and
// wasm bytes (and optional assets) using the same on-disk layout the
// scan path expects.
func buildPackageBytes(t *testing.T, manifest []byte, wasm []byte, assets map[string][]byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	if manifest != nil {
		w, err := zw.Create(pkgManifestFilename)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write(manifest); err != nil {
			t.Fatal(err)
		}
	}
	if wasm != nil {
		w, err := zw.Create(pkgWasmFilename)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write(wasm); err != nil {
			t.Fatal(err)
		}
	}
	for name, data := range assets {
		w, err := zw.Create(pkgAssetsPrefix + name)
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

func TestManager_Install_WritesPackageAndDiscovers(t *testing.T) {
	wasmPath := findExampleWasm(t)
	wasmBytes, _ := os.ReadFile(wasmPath)

	dir := t.TempDir()
	mgr := NewManager(dir, &HostEnv{})
	if err := mgr.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer mgr.Stop(context.Background())

	pkg := buildPackageBytes(t, validManifestBytes(), wasmBytes, nil)
	entry, err := mgr.Install(context.Background(), pkg)
	if err != nil {
		t.Fatalf("install: %v", err)
	}
	if entry.Slug != "hello-world" {
		t.Errorf("entry slug: got %q want hello-world", entry.Slug)
	}
	expectedPath := filepath.Join(dir, "hello-world.ocpkg")
	if _, err := os.Stat(expectedPath); err != nil {
		t.Errorf("expected file %s to exist after install: %v", expectedPath, err)
	}
	if len(mgr.List()) != 1 {
		t.Errorf("install should leave one discovered plugin, got %d", len(mgr.List()))
	}
}

func TestManager_Install_RejectsEmptyUpload(t *testing.T) {
	mgr := NewManager(t.TempDir(), &HostEnv{})
	_, err := mgr.Install(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for empty upload")
	}
}

func TestManager_Install_RejectsOversizedUpload(t *testing.T) {
	mgr := NewManager(t.TempDir(), &HostEnv{})
	oversized := make([]byte, MaxUploadBytes+1)
	_, err := mgr.Install(context.Background(), oversized)
	if err == nil {
		t.Fatal("expected error for oversized upload")
	}
	if !strings.Contains(err.Error(), "cap") {
		t.Errorf("error should mention the cap: got %v", err)
	}
}

func TestManager_Install_RejectsNonOcpkg(t *testing.T) {
	mgr := NewManager(t.TempDir(), &HostEnv{})
	_, err := mgr.Install(context.Background(), []byte("not a zip"))
	if err == nil {
		t.Fatal("expected error for garbage upload")
	}
}

func TestManager_Install_RejectsMissingManifest(t *testing.T) {
	wasmPath := findExampleWasm(t)
	wasmBytes, _ := os.ReadFile(wasmPath)

	mgr := NewManager(t.TempDir(), &HostEnv{})
	pkg := buildPackageBytes(t, nil, wasmBytes, nil)
	_, err := mgr.Install(context.Background(), pkg)
	if err == nil {
		t.Fatal("expected error for package missing manifest")
	}
	if !strings.Contains(err.Error(), "manifest") {
		t.Errorf("error should mention manifest: got %v", err)
	}
}

func TestManager_Uninstall_DeletesFileAndClearsState(t *testing.T) {
	wasmPath := findExampleWasm(t)
	wasmBytes, _ := os.ReadFile(wasmPath)
	dir := t.TempDir()

	ctx := context.Background()
	mgr := NewManager(dir, &HostEnv{})
	if err := mgr.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer mgr.Stop(ctx)

	pkg := buildPackageBytes(t, validManifestBytes(), wasmBytes, nil)
	if _, err := mgr.Install(ctx, pkg); err != nil {
		t.Fatalf("install: %v", err)
	}
	if err := mgr.Enable(ctx, "hello-world"); err != nil {
		t.Fatalf("enable: %v", err)
	}
	if len(mgr.Snapshot()) != 1 {
		t.Fatalf("expected 1 loaded before uninstall, got %d", len(mgr.Snapshot()))
	}

	if err := mgr.Uninstall(ctx, "hello-world"); err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	if len(mgr.List()) != 0 {
		t.Errorf("expected 0 discovered after uninstall, got %d", len(mgr.List()))
	}
	if len(mgr.Snapshot()) != 0 {
		t.Errorf("expected 0 loaded after uninstall, got %d", len(mgr.Snapshot()))
	}
	if _, err := os.Stat(filepath.Join(dir, "hello-world.ocpkg")); !os.IsNotExist(err) {
		t.Errorf("expected .ocpkg file to be removed, stat err = %v", err)
	}

	// Persistence: the cleared state survives a restart. A fresh manager
	// should see no enabled set and no approval entry.
	mgr.Stop(ctx)
	mgr2 := NewManager(dir, &HostEnv{})
	if err := mgr2.Start(ctx); err != nil {
		t.Fatalf("restart: %v", err)
	}
	defer mgr2.Stop(ctx)
	raw, err := os.ReadFile(filepath.Join(dir, ".enabled.json"))
	if err == nil {
		var f enabledFileContents
		if err := json.Unmarshal(raw, &f); err != nil {
			t.Fatalf("decode persisted state: %v", err)
		}
		if len(f.Enabled) != 0 {
			t.Errorf("enabled set should be empty after uninstall, got %v", f.Enabled)
		}
		if _, ok := f.ApprovedPermissions["hello-world"]; ok {
			t.Error("approval snapshot should be cleared after uninstall")
		}
	}
}

func TestManager_Uninstall_UnknownPluginErrors(t *testing.T) {
	mgr := NewManager(t.TempDir(), &HostEnv{})
	if err := mgr.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer mgr.Stop(context.Background())
	err := mgr.Uninstall(context.Background(), "nothing-here")
	if err == nil {
		t.Fatal("expected error uninstalling unknown plugin")
	}
}

func TestManager_Install_UpdatesExistingPlugin(t *testing.T) {
	wasmPath := findExampleWasm(t)
	wasmBytes, _ := os.ReadFile(wasmPath)

	dir := t.TempDir()
	mgr := NewManager(dir, &HostEnv{})
	if err := mgr.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer mgr.Stop(context.Background())

	// First install: hello-world 0.1.0
	pkg1 := buildPackageBytes(t, validManifestBytes(), wasmBytes, nil)
	if _, err := mgr.Install(context.Background(), pkg1); err != nil {
		t.Fatalf("first install: %v", err)
	}

	// Second install of the same plugin name with a bumped description.
	updated := []byte(`{
		"api": "1",
		"name": "hello-world",
		"version": "0.1.0",
		"description": "updated description",
		"permissions": []
	}`)
	pkg2 := buildPackageBytes(t, updated, wasmBytes, nil)
	entry, err := mgr.Install(context.Background(), pkg2)
	if err != nil {
		t.Fatalf("second install: %v", err)
	}
	if entry.Description != "updated description" {
		t.Errorf("description after update: got %q want %q", entry.Description, "updated description")
	}
	if len(mgr.List()) != 1 {
		t.Errorf("update should not duplicate the entry, got %d", len(mgr.List()))
	}
}
