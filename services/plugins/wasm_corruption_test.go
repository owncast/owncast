package plugins

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// TestIsWasmCorrupted validates the error-string predicate.
// Error strings were captured from live wazero runs; see
// TestFetchCorruptionDetectedAndRecovered for how they were collected.
func TestIsWasmCorrupted(t *testing.T) {
	cases := []struct {
		desc string
		err  error
		want bool
	}{
		// --- real runtime signals (captured from live runs) ---
		{
			// JS and Python engines, call 1: wazero recovers the http_request
			// panic; shadow stack permanently leaked. For QuickJS, one such
			// panic exhausts the guest stack entirely; Python tolerates repeated
			// panics without wedging.
			desc: "first-call recovered-by-wazero panic (JS and Python engines)",
			err: errors.New(
				`Get "https://bogus.invalid/": dial tcp: lookup bogus.invalid: no such host (recovered by wazero)
wasm stack trace:
	extism:host/env.http_request(i64,i64) i64`,
			),
			want: true,
		},
		{
			// JS engine, calls 2+: shadow stack exhausted; instance permanently
			// wedged. NOT matched here because (a) detection already fired on
			// call 1 via "(recovered by wazero)", (b) adding "unreachable"
			// would cause false positives for legitimate wasm guest panics
			// (Rust abort, etc.) that don't corrupt shadow-stack state.
			desc: "JS engine subsequent-call unreachable (not matched — call-1 detection is sufficient)",
			err:  errors.New("wasm error: unreachable\nwasm stack trace:\n\t.$1368(i32,i32,i32,i32,i32)"),
			want: false,
		},
		{
			// Self-contained .wasm modules (Rust/Go/Zig): smaller stack →
			// ~20 panics to exhaust → OOB trap on the final call.
			desc: "OOB trap (self-contained .wasm, extism/go-sdk#102 original report)",
			err:  errors.New("wasm error: out of bounds memory access"),
			want: true,
		},
		{
			desc: "OOB mid-message",
			err:  errors.New("some wrapper: out of bounds memory access in module"),
			want: true,
		},
		// --- false positives to reject ---
		{
			desc: "plain network error without recovery marker",
			err:  errors.New(`Get "https://bogus.invalid/": dial tcp: lookup bogus.invalid: no such host`),
			want: false,
		},
		{
			desc: "filter timeout",
			err:  errors.New("on_filter timed out after 50ms"),
			want: false,
		},
		{
			desc: "plain wasm unreachable without recovery marker",
			err:  errors.New("wasm error: unreachable"),
			want: false,
		},
	}
	for _, tc := range cases {
		got := isWasmCorrupted(tc.err)
		if got != tc.want {
			t.Errorf("[%s] isWasmCorrupted = %v, want %v\nerr: %q", tc.desc, got, tc.want, tc.err)
		}
	}
}

// TestFireCorruptedOnlyOnce confirms the CAS dedup: N concurrent callers
// invoke onCorrupted exactly once.
func TestFireCorruptedOnlyOnce(t *testing.T) {
	var mu sync.Mutex
	count := 0
	done := make(chan struct{})
	l := &Loaded{
		onCorrupted: func() {
			mu.Lock()
			count++
			mu.Unlock()
			close(done)
		},
	}

	const goroutines = 50
	var wg sync.WaitGroup
	for range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.fireCorrupted()
		}()
	}
	wg.Wait()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("onCorrupted goroutine did not complete")
	}

	mu.Lock()
	defer mu.Unlock()
	if count != 1 {
		t.Errorf("onCorrupted called %d times, want exactly 1", count)
	}
}

// TestFireCorruptedNilSafe ensures fireCorrupted with no callback set does not panic.
func TestFireCorruptedNilSafe(t *testing.T) {
	l := &Loaded{}
	l.fireCorrupted()
	l.fireCorrupted()
}

// buildFetchCrasherOcpkg returns an in-memory .ocpkg archive containing a
// plugin.js that calls owncast.http.fetch("https://bogus.invalid/") when it
// receives a chat event with body "fetch", and owncast.chat.send("ok:"+body)
// for all other bodies. Using a conditional fetch lets the test prove that
// after a reload, non-fetching events succeed.
func buildFetchCrasherOcpkg(t *testing.T) []byte {
	t.Helper()
	const script = `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  onChatMessage(msg) {
    if (msg.body === "fetch") {
      // The JS try/catch is cosmetic: the panic fires inside the extism
      // http_request HOST function and is recovered by wazero before control
      // returns to the guest, so the error surfaces from CallWithContext.
      try { owncast.http.fetch("https://bogus.invalid/"); } catch (e) {}
    } else {
      owncast.chat.send("ok:" + msg.body);
    }
  }
});`

	manifest := map[string]any{
		"api":         "1",
		"name":        "fetch-crasher",
		"slug":        "fetch-crasher",
		"version":     "0.1.0",
		"permissions": []string{PermNetworkFetch, PermChatSend},
		"network":     map[string]any{"allowedHosts": []string{"bogus.invalid"}},
	}
	mb, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("marshal manifest: %v", err)
	}

	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	for name, data := range map[string][]byte{
		"plugin.manifest.json": mb,
		"plugin.js":            []byte(script),
	} {
		f, err := w.Create(name)
		if err != nil {
			t.Fatalf("zip create %s: %v", name, err)
		}
		if _, err := f.Write(data); err != nil {
			t.Fatalf("zip write %s: %v", name, err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatalf("zip close: %v", err)
	}
	return buf.Bytes()
}

// TestFetchCorruptionDetectedAndRecovered is the end-to-end proof of
// extism/go-sdk#102 and the fix.
//
// What it proves:
//  1. Reproduced: a single owncast.http.fetch call to an unresolvable domain
//     causes the extism http_request host function to panic. wazero recovers
//     the panic but permanently corrupts the QuickJS engine instance's shadow
//     stack. The first-call error is captured independently of isWasmCorrupted
//     so a wrong detection string causes FAIL, not a false Skip.
//  2. String match: isWasmCorrupted matches the ACTUAL wazero error string.
//  3. Reload executes: callPlugin fires the real Manager.loadInternal-wired
//     onCorrupted callback, which calls Manager.Reload. The test waits for a
//     distinct new Loaded instance to appear in Manager.loaded.
//  4. Healthy after reload: a non-fetching event dispatched to the reloaded
//     instance triggers owncast.chat.send, proving the plugin is functional.
func TestFetchCorruptionDetectedAndRecovered(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	// Write an .ocpkg to a temp dir and use a real Manager so the production
	// loadInternal→onCorrupted→Reload chain executes (not a stub).
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "fetch-crasher.ocpkg")
	if err := os.WriteFile(pkgPath, buildFetchCrasherOcpkg(t), 0o644); err != nil {
		t.Fatalf("write ocpkg: %v", err)
	}

	env, sends, sendsMu := captureEnv()
	mgr := NewManager(dir, env)
	if err := mgr.Start(ctx); err != nil {
		t.Fatalf("manager start: %v", err)
	}
	defer mgr.Stop(ctx)
	if err := mgr.Enable(ctx, "fetch-crasher"); err != nil {
		t.Fatalf("enable: %v", err)
	}

	// --- Phase 1: reproduce the bug ---
	// Break on the FIRST non-nil error regardless of isWasmCorrupted so a
	// wrong detection string causes FAIL (not a false Skip). We call the plugin
	// directly (bypassing the dispatcher) for a synchronous, deterministic error.
	mgr.mu.RLock()
	original := mgr.loaded["fetch-crasher"]
	mgr.mu.RUnlock()
	if original == nil {
		t.Fatal("plugin not loaded after Enable")
	}

	fetchInput, _ := json.Marshal(Envelope{
		EventType: EventChatMessageReceived,
		Payload:   map[string]any{"body": "fetch"},
	})

	var firstErr error
	for range 5 {
		callCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		original.mu.Lock()
		_, _, callErr := original.callPlugin(callCtx, "on_event", fetchInput)
		original.mu.Unlock()
		cancel()
		if callErr != nil {
			firstErr = callErr
			break // stop at the first error regardless of kind
		}
	}
	if firstErr == nil {
		t.Skip("no error from fetch on bogus.invalid — " +
			"upstream extism/go-sdk fix may be in this build; nothing to recover from")
		return
	}
	t.Logf("first error: %v", firstErr)

	// --- Phase 2: verify isWasmCorrupted matches the REAL error ---
	// If this FAILS (not skips), update the substring in isWasmCorrupted to
	// the string logged above.
	if !isWasmCorrupted(firstErr) {
		t.Fatalf(
			"isWasmCorrupted did not match the actual wazero error — "+
				"update the detection substring.\nActual error: %q",
			firstErr,
		)
	}

	// --- Phase 3: wait for Manager.Reload to replace the instance ---
	// onCorrupted (wired by loadInternal) fires the real Manager.Reload in a
	// goroutine. Reload deletes the old entry, runs loadInternal, and replaces
	// mgr.loaded["fetch-crasher"] with a new *Loaded pointer.
	var reloaded *Loaded
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		mgr.mu.RLock()
		reloaded = mgr.loaded["fetch-crasher"]
		mgr.mu.RUnlock()
		if reloaded != nil && reloaded != original {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if reloaded == nil || reloaded == original {
		t.Fatal("Manager.Reload did not replace the instance within 10s")
	}
	t.Logf("instance replaced: original=%p reloaded=%p", original, reloaded)

	// --- Phase 4: prove the reloaded instance is healthy ---
	// Dispatch a non-fetch event through the live-snapshot dispatcher; the
	// plugin should call owncast.chat.send("ok:ping"). This exercises the
	// full callPlugin→on_event path on the new instance.
	sendsMu.Lock()
	*sends = nil
	sendsMu.Unlock()

	d := NewLiveDispatcher(mgr.Snapshot)
	d.Dispatch(ctx, EventChatMessageReceived, map[string]any{"body": "ping"})

	// Notify is async; give the goroutine time to complete.
	var got []ChatSendRequest
	deadline = time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		sendsMu.Lock()
		got = make([]ChatSendRequest, len(*sends))
		copy(got, *sends)
		sendsMu.Unlock()
		if len(got) > 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	found := false
	for _, s := range got {
		if s.Text == "ok:ping" {
			found = true
		}
	}
	if !found {
		t.Errorf("reloaded instance did not send chat on non-fetch event; "+
			"sends=%v — plugin may still be broken after reload", got)
	} else {
		t.Logf("reloaded instance healthy: chat.send produced %v", got)
	}
}
