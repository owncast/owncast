package plugins

import (
	"context"
	"fmt"
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
// It is package-level (not Manager-owned) for the same reason the plugin
// registry and the compilation cache are: the load path (loadFromBytes,
// LoadPlugin) runs without a Manager in the test runner and package preflight.
// In production there is exactly one *HostEnv; the engine's host functions are
// built once against the env captured on first compile.
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

type engineCache struct {
	mu    sync.Mutex
	byKey map[engineKey]*extism.CompiledPlugin
}

var compiledEngines = &engineCache{byKey: map[engineKey]*extism.CompiledPlugin{}}

// get returns the compiled engine for a (language, env) pair, compiling and
// memoizing it on first use. Host functions are built once per engine and
// shared across every instance it produces; they resolve the calling plugin's
// identity at call time (see registry.go).
func (c *engineCache) get(ctx context.Context, env *HostEnv, lang string) (*extism.CompiledPlugin, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := engineKey{lang: lang, env: env}
	if cp, ok := c.byKey[key]; ok {
		return cp, nil
	}
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
	// Share the compilation cache and build the full host-function set once.
	// ModuleConfig is set per-instance (Instance()), not here.
	pc := extism.PluginConfig{
		EnableWasi:    true,
		RuntimeConfig: wazero.NewRuntimeConfig().WithCompilationCache(sharedCompilationCache()),
	}
	cp, err := extism.NewCompiledPlugin(ctx, em, pc, BuildHostFunctions(env))
	if err != nil {
		return nil, fmt.Errorf("compile %s engine: %w", lang, err)
	}
	c.byKey[key] = cp
	return cp, nil
}
