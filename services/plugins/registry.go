package plugins

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"sync"

	extism "github.com/extism/go-sdk"
	"github.com/owncast/owncast/services/plugins/kv"
)

// pluginIdentity is everything a host function needs to serve a call on behalf
// of a specific plugin. Host functions are built once and shared across all
// plugins (of every runtime), so the calling plugin's identity is resolved at
// call time from this registry, keyed by the slug the host stashed in the
// instance's Extism config under "__slug".
type pluginIdentity struct {
	slug        string
	granted     map[string]bool // permission set, gated at call time
	chatDisplay string          // manifest.ChatDisplayName(), for chat sends
	assetsFS    fs.FS           // per-plugin bundled assets/ (may be nil)
	manifest    *Manifest       // for config_get (declared config) + action buttons
	kvNamespace kv.Namespace    // env.KV.Namespace(slug), resolved once at load
}

// can reports whether the plugin was granted perm. Ambient capabilities (timer,
// config, asset read) don't call this.
func (id *pluginIdentity) can(perm string) bool {
	return id != nil && id.granted[perm]
}

// pluginRegistry maps slug -> identity for every currently-loaded shared-engine
// plugin. It has its own lock, independent of Manager.mu, because host
// functions read it on the hot call path with no Manager involved (the test
// runner and package preflight also load plugins without a Manager). Slug
// uniqueness is guaranteed upstream by the Manager's discovered/loaded maps.
type pluginRegistry struct {
	mu     sync.RWMutex
	byslug map[string]*pluginIdentity
}

// Reserved Extism config keys the host sets on a plugin instance at load. The
// host reads "__slug" back to resolve identity inside shared host functions.
// The guest reads "manifest" for every runtime and "script" for shared engines
// only, since a self-contained module already is its code. These live in
// Extism's config namespace, which is separate from the author's declared
// manifest.config that owncast_config_get serves. Manifest validation also
// rejects declared keys starting with "__".
const (
	configKeySlug     = "__slug"
	configKeyScript   = "script"
	configKeyManifest = "manifest"
)

var globalPluginRegistry = &pluginRegistry{byslug: map[string]*pluginIdentity{}}

func (r *pluginRegistry) put(id *pluginIdentity) {
	r.mu.Lock()
	r.byslug[id.slug] = id
	r.mu.Unlock()
}

func (r *pluginRegistry) get(slug string) (*pluginIdentity, bool) {
	r.mu.RLock()
	id, ok := r.byslug[slug]
	r.mu.RUnlock()
	return id, ok
}

func (r *pluginRegistry) remove(slug string) {
	r.mu.Lock()
	delete(r.byslug, slug)
	r.mu.Unlock()
}

// resolveCaller resolves the calling plugin's identity and enforces its
// permission grant. It is the single permission-gating boundary for host
// functions. Every plugin links every host import, so an ungranted capability
// can only be refused here. It used to simply not be linked when host functions
// were built per plugin. Pass perm == "" for ambient capabilities (timer,
// config, asset read) that need no grant. On a false return the host function
// must abort by returning zero or empty to the guest. Permission denials are
// logged.
func resolveCaller(ctx context.Context, fnName, perm string) (*pluginIdentity, bool) {
	id, ok := callerIdentity(ctx)
	if !ok {
		return nil, false
	}
	if perm != "" && !id.can(perm) {
		fmt.Fprintf(os.Stderr, "%s: plugin %q lacks permission %q\n", fnName, id.slug, perm)
		return nil, false
	}
	return id, true
}

// callerIdentity resolves the plugin making the current host-function call. The
// extism go-sdk injects the calling *extism.Plugin into the context at every
// Call (see PluginCtxKey), and the host sets inst.Config["__slug"] before the
// first call, so we read the slug back here and look up its identity.
func callerIdentity(ctx context.Context) (*pluginIdentity, bool) {
	p, ok := ctx.Value(extism.PluginCtxKey("plugin")).(*extism.Plugin)
	if !ok || p == nil {
		return nil, false
	}
	slug := p.Config[configKeySlug]
	if slug == "" {
		return nil, false
	}
	return globalPluginRegistry.get(slug)
}
