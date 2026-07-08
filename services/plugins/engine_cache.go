package plugins

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"

	extism "github.com/extism/go-sdk"
	"github.com/owncast/owncast/services/plugins/engines"
	"github.com/tetratelabs/wazero"
)

// compiledEngines holds the per-language interpreter engines, each compiled
// once (extism.NewCompiledPlugin) from the embedded engine wasm and then
// instantiated per plugin. This is what collapses per-plugin memory: instead of
// every plugin compiling its own copy of QuickJS/CPython, one compiled engine
// is shared across all plugins of that language.
//
// Engines are reference-counted: every live plugin instance holds one
// reference (taken by acquire, dropped by the release func it returns). When
// the last plugin of a language unloads, the engine's runtime is closed and
// its compiled native code is freed, so a server with all plugins disabled
// pays no engine memory at all. The next load of that language recompiles the
// engine from the embedded wasm.
// ponytail: no linger on the last release — disabling then re-enabling the
// sole plugin of a language recompiles the engine (sub-second for JS, a
// couple of seconds for Python). Add a delayed close if that churn ever
// matters.
//
// It is package-level (not Manager-owned) for the same reason the plugin
// registry is: the load path (loadFromBytes, LoadPlugin) runs without a
// Manager in the test runner and package preflight.
// engineKey identifies a compiled engine. It includes the *HostEnv because the
// engine's host functions close over it: in production there is exactly one
// HostEnv (one Manager), so all plugins of a language share one engine and the
// per-plugin memory win holds. Tests construct a fresh HostEnv per case, so
// keying by env keeps each case's host calls routed to its own env rather than
// whichever env happened to compile the engine first.
type engineKey struct {
	lang string
	env  *HostEnv
}

// engineEntry pairs a compiled engine with the number of live plugin
// instances built from it.
type engineEntry struct {
	cp   *extism.CompiledPlugin
	refs int
}

type engineCache struct {
	mu    sync.Mutex
	byKey map[engineKey]*engineEntry
}

var compiledEngines = &engineCache{byKey: map[engineKey]*engineEntry{}}

// acquire returns the compiled engine for a (language, env) pair, compiling
// and memoizing it on first use, and takes a reference on it. The returned
// release func drops that reference; the caller must invoke it exactly once
// when the plugin instance built from the engine is closed. Host functions
// are built once per engine and shared across every instance it produces;
// they resolve the calling plugin's identity at call time (see registry.go).
func (c *engineCache) acquire(ctx context.Context, env *HostEnv, lang string) (*extism.CompiledPlugin, func(), error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := engineKey{lang: lang, env: env}
	entry, ok := c.byKey[key]
	if !ok {
		cp, err := compileEngine(ctx, env, lang)
		if err != nil {
			return nil, nil, err
		}
		entry = &engineEntry{cp: cp}
		c.byKey[key] = entry
	}
	entry.refs++
	return entry.cp, func() { c.release(key, entry) }, nil
}

// release drops one reference on an engine. When the last reference goes, the
// engine's runtime (and with it the compiled native code) is closed and freed,
// and the freed pages are handed back to the OS so the drop is visible to an
// operator watching process memory.
func (c *engineCache) release(key engineKey, entry *engineEntry) {
	c.mu.Lock()
	entry.refs--
	closeIt := entry.refs == 0 && c.byKey[key] == entry
	if closeIt {
		delete(c.byKey, key)
	}
	c.mu.Unlock()
	if closeIt {
		_ = entry.cp.Close(context.Background())
		debug.FreeOSMemory()
	}
}

// compileEngine compiles the embedded engine wasm for a language. Called with
// the cache lock held; compilation is rare (first plugin of a language after
// none were loaded) so holding the lock through it is fine.
func compileEngine(ctx context.Context, env *HostEnv, lang string) (*extism.CompiledPlugin, error) {
	wasmBytes, ok := engines.Bytes(lang)
	if !ok {
		return nil, fmt.Errorf("no embedded engine for runtime %q", lang)
	}

	em := extism.Manifest{
		Wasm:    []extism.Wasm{extism.WasmData{Data: wasmBytes, Name: "engine-" + lang}},
		Timeout: 10_000,
		Memory: &extism.ManifestMemory{
			MaxPages:             MaxWasmPages,
			MaxHttpResponseBytes: MaxExtismHTTPResponseBytes,
			MaxVarBytes:          MaxExtismVarBytes,
		},
		// AllowedHosts is intentionally omitted: network scope is per plugin,
		// not per engine, so it's set on each instance (inst.AllowedHosts)
		// from the plugin's manifest at load time.
	}
	// No external compilation cache: the engine's compiled code must die with
	// its runtime when the last plugin of the language unloads. Reuse while
	// plugins are loaded comes from the engineCache map itself.
	// ModuleConfig is set per-instance (Instance()), not here.
	pc := extism.PluginConfig{
		EnableWasi:    true,
		RuntimeConfig: wazero.NewRuntimeConfig(),
	}
	cp, err := extism.NewCompiledPlugin(ctx, em, pc, BuildHostFunctions(env))
	if err != nil {
		return nil, fmt.Errorf("compile %s engine: %w", lang, err)
	}
	// Compilation churns through large transient buffers (parsing, IR) that
	// the Go runtime would otherwise sit on for minutes. Hand them back now so
	// the steady-state cost of an enabled plugin is the compiled code and its
	// instances, not the compiler's scratch space (measured: ~90 MiB retained
	// across the JS+Python engine compiles without this).
	debug.FreeOSMemory()
	return cp, nil
}
