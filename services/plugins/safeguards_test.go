package plugins

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// findSafeguardStressWasm locates the safeguard-stress example. Tests skip
// (rather than fail) if it hasn't been built, mirroring the pattern used by
// findExampleWasm.
func findSafeguardStressWasm(t *testing.T) (wasmPath, manifestPath string) {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")
	wasmPath = filepath.Join(repoRoot, "examples", "js", "safeguard-stress", "safeguard-stress.wasm")
	manifestPath = filepath.Join(repoRoot, "examples", "js", "safeguard-stress", "plugin.manifest.json")
	if _, err := os.Stat(wasmPath); err != nil {
		t.Skipf("safeguard-stress wasm not built; run tools/build-plugin.sh examples/js/safeguard-stress")
	}
	return wasmPath, manifestPath
}

// loadStress is the common setup for safeguard tests — loads the
// safeguard-stress plugin with a no-op HostEnv. Returns the loaded plugin
// and a cleanup func.
func loadStress(t *testing.T) (*Loaded, func()) {
	t.Helper()
	wasmPath, manifestPath := findSafeguardStressWasm(t)
	ctx := context.Background()
	loaded, err := LoadPlugin(ctx, &HostEnv{}, wasmPath, manifestPath)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	return loaded, func() { loaded.Close(ctx) }
}

// --- Filter safeguards ---

func TestSafeguard_FilterHugeOutputRejected(t *testing.T) {
	loaded, cleanup := loadStress(t)
	defer cleanup()

	// QuickJS serializing a 1 MB JSON output legitimately doesn't fit in
	// the 50 ms production timeout — the timeout would fire first and we'd
	// never reach the size check. Production safety in practice comes from
	// the timeout; this test stretches the timeout to verify the *separate*
	// size-cap defense (which would matter for faster runtimes returning
	// a pre-allocated huge buffer instantly).
	orig := FilterTimeout
	FilterTimeout = 10 * time.Second
	defer func() { FilterTimeout = orig }()

	_, err := callOnFilter(context.Background(), loaded, "chat.message.received", map[string]any{
		"cmd": "huge-output",
	})
	if err == nil {
		t.Fatal("expected error for huge filter output, got nil")
	}
	if !strings.Contains(err.Error(), "too large") {
		t.Errorf("expected size error, got %v", err)
	}
}

func TestSafeguard_FilterSpinTriggersTimeout(t *testing.T) {
	loaded, cleanup := loadStress(t)
	defer cleanup()

	start := time.Now()
	_, err := callOnFilter(context.Background(), loaded, "chat.message.received", map[string]any{
		"cmd": "spin",
	})
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Errorf("expected timeout in error, got %v", err)
	}
	// 50 ms cap + generous slack for QuickJS cancellation latency. If this
	// takes seconds, the cancellation isn't working.
	if elapsed > 2*time.Second {
		t.Errorf("filter spin took %v — timeout enforcement is too slow", elapsed)
	}
}

func TestSafeguard_FilterFailuresAccumulateStrikes(t *testing.T) {
	loaded, cleanup := loadStress(t)
	defer cleanup()

	ctx := context.Background()
	for i := 0; i < FilterStrikeThreshold; i++ {
		if loaded.IsDisabled() {
			t.Fatalf("plugin disabled too early at iteration %d", i)
		}
		_, err := callOnFilter(ctx, loaded, "chat.message.received", map[string]any{"cmd": "huge-output"})
		if err == nil {
			t.Fatalf("iteration %d: expected error", i)
		}
		// Mirror dispatcher.Filter's strike bookkeeping.
		loaded.recordFilterFailure()
	}
	if !loaded.IsDisabled() {
		t.Fatal("plugin should be disabled after threshold reached")
	}
}

// --- Notification safeguards ---

func TestSafeguard_NotifyTimeoutEnforced(t *testing.T) {
	loaded, cleanup := loadStress(t)
	defer cleanup()

	envelope, _ := json.Marshal(Envelope{
		EventType: "chat.message.received",
		Payload:   map[string]any{"cmd": "spin"},
	})

	start := time.Now()
	err := callOnEvent(context.Background(), loaded, envelope)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Errorf("expected timeout in error, got %v", err)
	}
	if elapsed > 2*time.Second {
		t.Errorf("on_event spin took %v — timeout enforcement is too slow", elapsed)
	}
}

// --- HTTP handler safeguards ---

func TestSafeguard_HTTPHandlerTimeoutReturns504(t *testing.T) {
	loaded, cleanup := loadStress(t)
	defer cleanup()

	server := NewServer([]*Loaded{loaded})
	req := httptest.NewRequest("GET", "/plugins/safeguard-stress/spin", nil)
	rec := httptest.NewRecorder()

	start := time.Now()
	server.ServeHTTP(rec, req)
	elapsed := time.Since(start)

	if rec.Code != http.StatusGatewayTimeout {
		t.Errorf("status: got %d want %d (body: %q)", rec.Code, http.StatusGatewayTimeout, rec.Body.String())
	}
	// Allow some slack over the 5s cap for cancellation latency.
	if elapsed > 10*time.Second {
		t.Errorf("HTTP spin took %v — timeout enforcement is too slow", elapsed)
	}
}

func TestSafeguard_HTTPHandlerHugeOutputRejected(t *testing.T) {
	loaded, cleanup := loadStress(t)
	defer cleanup()

	server := NewServer([]*Loaded{loaded})
	req := httptest.NewRequest("GET", "/plugins/safeguard-stress/huge", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d want 500", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "too large") {
		t.Errorf("body should mention size: got %q", rec.Body.String())
	}
}

// --- Wasm memory cap ---

func TestSafeguard_WasmMemoryCapPreventsHugeAllocation(t *testing.T) {
	loaded, cleanup := loadStress(t)
	defer cleanup()

	_, err := callOnFilter(context.Background(), loaded, "chat.message.received", map[string]any{
		"cmd": "alloc",
	})
	if err == nil {
		t.Fatal("expected error when plugin tries to allocate past MaxWasmPages, got nil")
	}
	// The exact error wording depends on whether QuickJS surfaces this as
	// an OOM throw or wazero rejects the memory.grow. Either way, the call
	// must fail rather than the host swallowing a 64+ MB allocation.
}

// --- Manifest sandbox values ---

// Verifies the constants are actually wired into the extism manifest at load
// time. A regression that drops the Memory block would silently re-open the
// crash vector this file exists to prevent.
func TestSafeguard_ManifestSandboxConstantsSet(t *testing.T) {
	// Sanity-check the values themselves haven't been zeroed. The actual
	// wiring is verified by TestSafeguard_WasmMemoryCapPreventsHugeAllocation
	// (functional test); this guards against silent removal.
	if MaxWasmPages == 0 {
		t.Error("MaxWasmPages is 0 — wasm memory is uncapped")
	}
	if MaxExtismHTTPResponseBytes == 0 {
		t.Error("MaxExtismHTTPResponseBytes is 0 — outbound HTTP body is uncapped")
	}
	if MaxFilterOutputBytes == 0 {
		t.Error("MaxFilterOutputBytes is 0 — filter output is uncapped")
	}
	if MaxHTTPHandlerOutputBytes == 0 {
		t.Error("MaxHTTPHandlerOutputBytes is 0 — HTTP handler output is uncapped")
	}
	if NotifyTimeout == 0 {
		t.Error("NotifyTimeout is 0 — on_event has no per-call cap")
	}
	if HTTPHandlerTimeout == 0 {
		t.Error("HTTPHandlerTimeout is 0 — on_http_request has no per-call cap")
	}
}
