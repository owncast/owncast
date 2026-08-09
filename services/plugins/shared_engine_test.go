package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

// captureEnv returns a HostEnv whose OnChat records every send, plus the slice
// it records into. Enough to exercise the shared-engine dispatch path.
func captureEnv() (*HostEnv, *[]ChatSendRequest, *sync.Mutex) {
	var mu sync.Mutex
	var sends []ChatSendRequest
	env := &HostEnv{
		OnChat: func(r ChatSendRequest) {
			mu.Lock()
			sends = append(sends, r)
			mu.Unlock()
		},
	}
	return env, &sends, &mu
}

func chatPayload(name, body string) any {
	return map[string]any{
		"id":        "1",
		"user":      map[string]any{"id": "u1", "displayName": name},
		"clientId":  7,
		"body":      body,
		"timestamp": "2024-01-01T00:00:00Z",
	}
}

// loadShared loads an inline script as a shared-engine plugin through the real
// production load path (loadFromBytes → engine.Instance → register).
func loadShared(t *testing.T, ctx context.Context, env *HostEnv, runtime, slug, script string, perms []string) *Loaded {
	t.Helper()
	// No "type" in the manifest — the runtime is passed explicitly, mirroring
	// how the host infers it from the code artifact's filename at load.
	m := map[string]any{
		"api":         "1",
		"name":        slug,
		"slug":        slug,
		"version":     "0.1.0",
		"permissions": perms,
	}
	manifestBytes, _ := json.Marshal(m)
	loaded, err := loadFromBytes(ctx, env, manifestBytes, []byte(script), runtime, slug, nil)
	if err != nil {
		t.Fatalf("loadFromBytes(%s/%s): %v", runtime, slug, err)
	}
	return loaded
}

func TestSharedEngineJSRegisterAndDispatch(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, sends, mu := captureEnv()
	script := `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  onChatMessage(msg) { owncast.chat.send(msg.user.displayName + " said: " + msg.body); }
});`
	loaded := loadShared(t, ctx, env, RuntimeJavaScript, "echo-bot", script, []string{PermChatSend})
	defer loaded.Close(ctx)

	// register() must have derived the notify subscription for onChatMessage.
	if !subscribed(loaded.Manifest.Subscriptions.Notify, EventChatMessageReceived) {
		t.Fatalf("expected notify subscription for %s, got %+v", EventChatMessageReceived, loaded.Manifest.Subscriptions)
	}

	d := NewLiveDispatcher(func() []*Loaded { return []*Loaded{loaded} })
	d.Dispatch(ctx, EventChatMessageReceived, chatPayload("alice", "hi"))

	mu.Lock()
	defer mu.Unlock()
	if len(*sends) != 1 {
		t.Fatalf("expected 1 chat send, got %d: %+v", len(*sends), *sends)
	}
	if got := (*sends)[0].Text; got != "alice said: hi" {
		t.Errorf("chat send text = %q, want %q", got, "alice said: hi")
	}
	if got := (*sends)[0].PluginSlug; got != "echo-bot" {
		t.Errorf("chat send routed to slug %q, want echo-bot", got)
	}
}

func TestSharedEngineActionValidationErrorReachesJavaScript(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, sends, mu := captureEnv()
	script := `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  onChatMessage() {
    try {
      owncast.actions.add({ title: "broken" });
    } catch (err) {
      owncast.chat.send(err.message);
    }
  }
});`
	loaded := loadShared(t, ctx, env, RuntimeJavaScript, "action-errors", script,
		[]string{PermChatSend, PermUIModify})
	defer loaded.Close(ctx)

	d := NewLiveDispatcher(func() []*Loaded { return []*Loaded{loaded} })
	d.Dispatch(ctx, EventChatMessageReceived, chatPayload("alice", "hi"))

	mu.Lock()
	defer mu.Unlock()
	if len(*sends) != 1 {
		t.Fatalf("expected one error message, got %d: %+v", len(*sends), *sends)
	}
	if got := (*sends)[0].Text; got != "actions[0]: exactly one of url or html is required" {
		t.Fatalf("action error = %q", got)
	}
}

func TestSharedEngineActionPermissionErrorReachesJavaScript(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, sends, mu := captureEnv()
	script := `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  onChatMessage() {
    try {
      owncast.actions.add({ title: "blocked" });
    } catch (err) {
      owncast.chat.send(err.message);
    }
  }
});`
	loaded := loadShared(t, ctx, env, RuntimeJavaScript, "action-permission-errors", script,
		[]string{PermChatSend})
	defer loaded.Close(ctx)

	d := NewLiveDispatcher(func() []*Loaded { return []*Loaded{loaded} })
	d.Dispatch(ctx, EventChatMessageReceived, chatPayload("alice", "hi"))

	mu.Lock()
	defer mu.Unlock()
	if len(*sends) != 1 {
		t.Fatalf("expected one permission error message, got %d: %+v", len(*sends), *sends)
	}
	if got := (*sends)[0].Text; got != "owncast.actions.add failed" {
		t.Fatalf("permission error = %q", got)
	}
}

// TestSharedEnginePermissionGate is the R1 check: a plugin that calls a host
// function for a permission it was not granted must be denied at call time
// (the shared engine links every host import, so the gate is the only
// enforcement). The plugin still loads; the call is simply a no-op.
func TestSharedEnginePermissionGate(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, sends, mu := captureEnv()
	script := `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  onChatMessage(msg) { owncast.chat.send("should be blocked"); }
});`
	// No chat.send permission granted.
	loaded := loadShared(t, ctx, env, RuntimeJavaScript, "sneaky", script, []string{})
	defer loaded.Close(ctx)

	d := NewLiveDispatcher(func() []*Loaded { return []*Loaded{loaded} })
	d.Dispatch(ctx, EventChatMessageReceived, chatPayload("alice", "hi"))

	mu.Lock()
	defer mu.Unlock()
	if len(*sends) != 0 {
		t.Fatalf("ungranted chat.send must be blocked; got %d sends: %+v", len(*sends), *sends)
	}
}

func TestSharedEngineLoggingIsAmbientAndPreservesIdentity(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	type logEntry struct {
		plugin  string
		level   PluginLogLevel
		message string
	}
	var entries []logEntry
	env := &HostEnv{
		Log: func(plugin string, level PluginLogLevel, message string) {
			entries = append(entries, logEntry{plugin: plugin, level: level, message: message})
		},
	}

	js := loadShared(t, ctx, env, RuntimeJavaScript, "js-logger", `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
owncast.log.info("ready");
owncast.log.warning("needs attention");
owncast.log.error("failed");
owncast.log.info("x".repeat(5000));
owncast.log.info("line one\nline two\u001b[31m");
module.exports = definePlugin({});`, nil)
	defer js.Close(ctx)

	py := loadShared(t, ctx, env, RuntimePython, "py-logger", `
owncast.log.info("ready")
owncast.log.warning("needs attention")
owncast.log.error("failed")
`, nil)
	defer py.Close(ctx)

	want := []logEntry{
		{plugin: "js-logger", level: PluginLogInfo, message: "ready"},
		{plugin: "js-logger", level: PluginLogWarning, message: "needs attention"},
		{plugin: "js-logger", level: PluginLogError, message: "failed"},
		{
			plugin:  "js-logger",
			level:   PluginLogInfo,
			message: strings.Repeat("x", MaxPluginLogMessageBytes-len(pluginLogTruncationSuffix)) + pluginLogTruncationSuffix,
		},
		{plugin: "js-logger", level: PluginLogInfo, message: "line one line two [31m"},
		{plugin: "py-logger", level: PluginLogInfo, message: "ready"},
		{plugin: "py-logger", level: PluginLogWarning, message: "needs attention"},
		{plugin: "py-logger", level: PluginLogError, message: "failed"},
	}
	if !reflect.DeepEqual(entries, want) {
		t.Fatalf("log entries = %+v, want %+v", entries, want)
	}
}

// TestDispatchCloseRace exercises the C1 lifecycle race: tearing a plugin down
// (Loaded.Close) while events are concurrently dispatched to it must neither
// data-race on the instance pointer nor nil-deref. Run under -race to catch the
// field race; before the snapshot-under-lock fix this also panicked (an
// unrecovered panic on the dispatch goroutine crashed the whole host).
func TestDispatchCloseRace(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, _, _ := captureEnv()
	script := `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  onChatMessage(msg) { owncast.chat.send("ok"); }
});`
	loaded := loadShared(t, ctx, env, RuntimeJavaScript, "racer", script, []string{PermChatSend})
	d := NewLiveDispatcher(func() []*Loaded { return []*Loaded{loaded} })

	var wg sync.WaitGroup
	// Many dispatchers fanning on_event out to the plugin...
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.Dispatch(ctx, EventChatMessageReceived, chatPayload("alice", "hi"))
		}()
	}
	// ...while it is torn down underneath them. The dispatches that land after
	// Close must safely no-op rather than crash.
	wg.Add(1)
	go func() {
		defer wg.Done()
		loaded.Close(ctx)
	}()
	wg.Wait()
}

// TestEmitSelfSubscriptionDoesNotDeadlock is the C3 regression: a plugin that
// emits an event it is itself subscribed to (notify) must not deadlock. The
// synchronous owncast.emit re-enters the dispatcher while the plugin's on_event
// call still holds its mutex; delivering the event back into the same plugin
// would block on that mutex forever (and, from the tick, freeze onTick for
// every plugin). Notify skips self-delivery on the emit chain.
func TestEmitSelfSubscriptionDoesNotDeadlock(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, _, _ := captureEnv()
	// Emits its fully qualified hook on every chat message AND subscribes
	// (notify) to the local "ping" hook.
	script := `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  onChatMessage(msg) { owncast.events.emit("echoer.ping", {}); },
  on: { ping(payload) {} },
});`
	loaded := loadShared(t, ctx, env, RuntimeJavaScript, "echoer", script, []string{PermEmitEvent})
	defer loaded.Close(ctx)

	d := NewLiveDispatcher(func() []*Loaded { return []*Loaded{loaded} })
	env.Emit = d.Dispatch // route owncast.emit through the dispatcher, as in production

	done := make(chan struct{})
	go func() {
		d.Dispatch(ctx, EventChatMessageReceived, chatPayload("alice", "hi"))
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(NotifyTimeout / 2):
		t.Fatal("Dispatch deadlocked on a re-entrant self-subscribed emit")
	}
}

// TestManifestParse covers the R4 reserved-config-key guard and confirms the
// manifest no longer requires (or validates) a "type" field — the runtime is
// inferred from the code artifact at load time.
func TestManifestParse(t *testing.T) {
	base := `{"api":"1","name":"X","slug":"x","version":"0.1.0"%s}`
	cases := []struct {
		name    string
		extra   string
		wantErr bool
	}{
		{"no type field is fine", "", false},
		{"stray type ignored, not validated", `,"type":"anything"`, false},
		{"reserved config key rejected", `,"config":{"__slug":{"type":"string"}}`, true},
		{"normal config key ok", `,"config":{"greeting":{"type":"string"}}`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := ParseManifest([]byte(fmt.Sprintf(base, tc.extra)))
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil (manifest %+v)", m)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// TestSharedEngineNetworkIsolation is the R3 check: a shared engine carries no
// host allowlist, so each instance must get its own network scope from its own
// manifest. Two plugins on the same engine must not inherit each other's
// allowed hosts.
func TestSharedEngineNetworkIsolation(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, _, _ := captureEnv()
	script := `const { definePlugin } = require("@owncast/plugin-sdk"); module.exports = definePlugin({});`

	load := func(slug, host string) *Loaded {
		m := map[string]any{
			"api": "1", "name": slug, "slug": slug, "version": "0.1.0",
			"permissions": []string{PermNetworkFetch},
			"network":     map[string]any{"allowedHosts": []string{host}},
		}
		mb, _ := json.Marshal(m)
		l, err := loadFromBytes(ctx, env, mb, []byte(script), RuntimeJavaScript, slug, nil)
		if err != nil {
			t.Fatalf("load %s: %v", slug, err)
		}
		return l
	}
	a := load("net-a", "a.example.com")
	defer a.Close(ctx)
	b := load("net-b", "b.example.com")
	defer b.Close(ctx)

	if got := a.plugin.AllowedHosts; len(got) != 1 || got[0] != "a.example.com" {
		t.Errorf("plugin A AllowedHosts = %v, want [a.example.com]", got)
	}
	if got := b.plugin.AllowedHosts; len(got) != 1 || got[0] != "b.example.com" {
		t.Errorf("plugin B AllowedHosts = %v, want [b.example.com]", got)
	}
}

// TestSharedEngineConcurrency is the R2 check: many instances of one shared
// compiled engine, dispatched concurrently, must route correctly and not race.
// Run the package with -race to exercise the shared wazero runtime.
func TestSharedEngineConcurrency(t *testing.T) {
	if runtime.GOOS == "openbsd" {
		// wazero has no compiler backend on OpenBSD, so plugins run in
		// interpreter mode there; the 500ms on_event budget is routinely
		// exceeded under concurrency and the test flakes on timing, not
		// on the routing/race contract it defends.
		t.Skip("wazero interpreter mode on openbsd is too slow for the on_event time budget")
	}
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, sends, mu := captureEnv()
	script := `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  onChatMessage(msg) { owncast.chat.send(msg.body); }
});`
	var loadeds []*Loaded
	for _, slug := range []string{"c1", "c2", "c3"} {
		l := loadShared(t, ctx, env, RuntimeJavaScript, slug, script, []string{PermChatSend})
		defer l.Close(ctx)
		loadeds = append(loadeds, l)
	}
	d := NewLiveDispatcher(func() []*Loaded { return loadeds })

	const perPlugin = 20
	var wg sync.WaitGroup
	for i := 0; i < perPlugin; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			d.Dispatch(ctx, EventChatMessageReceived, chatPayload("u", "m"))
		}(i)
	}
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	// Each of the perPlugin dispatches notifies all 3 plugins.
	if want := perPlugin * len(loadeds); len(*sends) != want {
		t.Fatalf("expected %d sends, got %d", want, len(*sends))
	}
}

// TestSharedEngineReportsAndDispatchesCommands verifies that JS and Python
// declarations reach core registration, avoid broad chat subscriptions, and
// receive the targeted command event.
func TestSharedEngineReportsAndDispatchesCommands(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, sends, mu := captureEnv()

	jsScript := `
const { definePlugin } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  commands: { uptime: { description: "Stream uptime", aliases: ["up"], run: (ctx) => ctx.reply("js:" + ctx.invokedAs + ":" + ctx.args.join(",")) } },
});`
	js := loadShared(t, ctx, env, RuntimeJavaScript, "js-cmds", jsScript, []string{PermChatSend})
	defer js.Close(ctx)
	assertHasCommand(t, js, "uptime", "Stream uptime")

	pyScript := `
plugin.commands({"uptime": {"description": "Stream uptime", "aliases": ["up"], "run": lambda ctx: ctx.reply("py:" + ctx.invoked_as + ":" + ",".join(ctx.args))}})
`
	py := loadShared(t, ctx, env, RuntimePython, "py-cmds", pyScript, []string{PermChatSend})
	defer py.Close(ctx)
	assertHasCommand(t, py, "uptime", "Stream uptime")

	for _, loaded := range []*Loaded{js, py} {
		if subscribed(loaded.Manifest.Subscriptions.Notify, EventChatMessageReceived) {
			t.Errorf("%s: command declaration must not subscribe to ordinary chat", loaded.Manifest.Slug)
		}
	}

	d := NewDispatcher([]*Loaded{js, py})
	d.DispatchCommands(ctx, HostChatMessage{
		ID:        "m1",
		User:      &HostUser{ID: "u1", DisplayName: "alice"},
		ClientID:  7,
		Body:      "!up one two",
		Timestamp: "2026-01-01T00:00:00Z",
	})

	mu.Lock()
	got := make(map[string]bool, len(*sends))
	for _, send := range *sends {
		got[send.Text] = true
	}
	mu.Unlock()
	for _, want := range []string{"js:up:one,two", "py:up:one,two"} {
		if !got[want] {
			t.Errorf("duplicate command dispatch missing %q; got %+v", want, got)
		}
	}

	// The host still aggregates the same registrations for its built-in help.
	help := BuildHelpMessage([]*Loaded{js, py}, false)
	if !strings.Contains(help, "`!uptime`") || !strings.Contains(help, "Stream uptime") {
		t.Errorf("help should list the reported command:\n%s", help)
	}
}

func assertHasCommand(t *testing.T, l *Loaded, name, desc string) {
	t.Helper()
	for _, c := range l.Manifest.Commands {
		if c.Name == name {
			if c.Description != desc {
				t.Errorf("%s: command %q description = %q, want %q", l.Manifest.Slug, name, c.Description, desc)
			}
			if c.Prefix != "!" {
				t.Errorf("%s: command %q prefix = %q, want \"!\"", l.Manifest.Slug, name, c.Prefix)
			}
			return
		}
	}
	t.Errorf("%s: register() did not report command %q; got %+v", l.Manifest.Slug, name, l.Manifest.Commands)
}

func TestSharedEngineActionValidationErrorReachesPython(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, sends, mu := captureEnv()
	script := `
@plugin.on_chat_message
def handle(msg):
    try:
        owncast.actions.add({"title": "broken"})
    except Exception as err:
        owncast.chat.send(str(err))
`
	loaded := loadShared(t, ctx, env, RuntimePython, "python-action-errors", script,
		[]string{PermChatSend, PermUIModify})
	defer loaded.Close(ctx)

	d := NewLiveDispatcher(func() []*Loaded { return []*Loaded{loaded} })
	d.Dispatch(ctx, EventChatMessageReceived, chatPayload("alice", "hi"))

	mu.Lock()
	defer mu.Unlock()
	if len(*sends) != 1 {
		t.Fatalf("expected one error message, got %d: %+v", len(*sends), *sends)
	}
	if got := (*sends)[0].Text; got != "actions[0]: exactly one of url or html is required" {
		t.Fatalf("action error = %q", got)
	}
}

func TestSharedEngineActionPermissionErrorReachesPython(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, sends, mu := captureEnv()
	script := `
@plugin.on_chat_message
def handle(msg):
    try:
        owncast.actions.add({"title": "blocked"})
    except Exception as err:
        owncast.chat.send(str(err))
`
	loaded := loadShared(t, ctx, env, RuntimePython, "python-action-permission-errors", script,
		[]string{PermChatSend})
	defer loaded.Close(ctx)

	d := NewLiveDispatcher(func() []*Loaded { return []*Loaded{loaded} })
	d.Dispatch(ctx, EventChatMessageReceived, chatPayload("alice", "hi"))

	mu.Lock()
	defer mu.Unlock()
	if len(*sends) != 1 {
		t.Fatalf("expected one permission error message, got %d: %+v", len(*sends), *sends)
	}
	if got := (*sends)[0].Text; got != "owncast.actions.add failed" {
		t.Fatalf("permission error = %q", got)
	}
}

func TestSharedEnginePythonRegister(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, sends, mu := captureEnv()
	// The SDK build strips the `from owncast_plugin import ...` line; the SDK
	// names (plugin, owncast) are already globals in the engine. Injecting the
	// import would make the frozen CPython try to load a nonexistent module.
	script := `
@plugin.on_chat_message
def greet(msg):
    owncast.chat.send(msg.user.display_name + " said: " + msg.body)
`
	loaded := loadShared(t, ctx, env, RuntimePython, "py-echo", script, []string{PermChatSend})
	defer loaded.Close(ctx)

	if !subscribed(loaded.Manifest.Subscriptions.Notify, EventChatMessageReceived) {
		t.Fatalf("python: expected notify subscription for %s, got %+v", EventChatMessageReceived, loaded.Manifest.Subscriptions)
	}

	d := NewLiveDispatcher(func() []*Loaded { return []*Loaded{loaded} })
	d.Dispatch(ctx, EventChatMessageReceived, chatPayload("bob", "yo"))

	mu.Lock()
	defer mu.Unlock()
	if len(*sends) != 1 {
		t.Fatalf("python: expected 1 chat send, got %d: %+v", len(*sends), *sends)
	}
	if got := (*sends)[0].Text; !strings.Contains(got, "bob said: yo") {
		t.Errorf("python chat send text = %q, want to contain %q", got, "bob said: yo")
	}
}

// resetForTest closes and clears the compiled engines and the plugin registry.
// Tests construct fresh HostEnvs per case; without this an engine compiled
// against an earlier test's env would leak into the next. It lives in a test
// file so it isn't flagged as unused production code.
func (c *engineCache) resetForTest(ctx context.Context) {
	c.mu.Lock()
	for k, entry := range c.byKey {
		_ = entry.cp.Close(ctx)
		delete(c.byKey, k)
	}
	c.mu.Unlock()
	globalPluginRegistry.mu.Lock()
	globalPluginRegistry.byslug = map[string]*pluginIdentity{}
	globalPluginRegistry.mu.Unlock()
}
