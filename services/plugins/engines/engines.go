// Package engines holds the shared language interpreter engines (QuickJS for
// JS, CPython for Python) compiled to WebAssembly. One engine per language is
// embedded into the Owncast binary and compiled once at runtime, then
// instantiated per plugin with the plugin's script injected via Extism config.
// This replaces the old model where every plugin embedded its own copy of the
// interpreter (which made memory scale linearly per plugin).
//
// The .wasm files are build artifacts produced in the owncast-plugin-sdk repo
// (engines/build.mjs, engines/build_py.py) from a fixed bootstrap entry, and
// committed here so Owncast's Go build stays toolchain-free. Rebuild and
// re-commit them when the SDK runtime or the bootstrap changes.
package engines

import _ "embed"

//go:embed javascript/engine.wasm
var jsEngine []byte

//go:embed python/engine.wasm
var pyEngine []byte

// Bytes returns the embedded engine wasm for a manifest Type ("javascript" or
// "python"), and whether that language has an engine.
func Bytes(lang string) ([]byte, bool) {
	switch lang {
	case "javascript":
		return jsEngine, true
	case "python":
		return pyEngine, true
	default:
		return nil, false
	}
}
