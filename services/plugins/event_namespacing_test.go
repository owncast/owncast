package plugins

import (
	"context"
	"sort"
	"strings"
	"sync"
	"testing"
)

// Custom plugin events are namespaced by the host, not by the plugin. The
// emitting plugin passes a suffix and hostEmitEvent prefixes it with the slug
// that resolveCaller derived from the wasm call, which the guest cannot
// influence. These tests pin the properties that follow from that:
//
//  1. a plugin can only publish under "<its-own-slug>.",
//  2. it therefore cannot impersonate another plugin's custom event, and
//  3. it cannot forge a core host event: the prefix keeps almost every
//     composed name distinct from the built-ins, and hostEmitEvent rejects
//     the residual collisions (a plugin slugged "chat" composing
//     "chat.message.received") against reservedEventTypes.
//
// The tests deliberately assert on the FULL set of deliveries rather than
// "the good one arrived": a forged delivery is an extra entry, so an
// exact-set comparison is what actually catches a regression.

// drainSends returns every chat send recorded so far and clears the buffer,
// so one loaded plugin pair can serve several scenarios.
func drainSends(sends *[]ChatSendRequest, mu *sync.Mutex) []string {
	mu.Lock()
	defer mu.Unlock()
	out := make([]string, 0, len(*sends))
	for _, s := range *sends {
		out = append(out, s.Text)
	}
	*sends = (*sends)[:0]
	sort.Strings(out)
	return out
}

func sameStrings(a, b []string) bool {
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

// receiverScript builds a plugin that subscribes to every name in names and
// reports each delivery as "got:<subscribed name>:<payload.n>". Subscribing to
// both the honest and the forged spelling of an event is the point: only the
// honest one may ever fire.
func receiverScript(names []string) string {
	return `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
const names = ["` + strings.Join(names, `","`) + `"];
const on = {};
for (const n of names) {
  on[n] = function (payload) {
    owncast.chat.send("got:" + n + ":" + (payload && payload.n));
  };
}
module.exports = definePlugin({ on });`
}

// TestEmitIsNamespacedByCallerSlug drives one emitter through several event
// names. The emitter emits whatever the chat body says, so a single load
// covers the honest case, the impersonation attempt, and the core-event
// forgery attempt.
func TestEmitIsNamespacedByCallerSlug(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, sends, mu := captureEnv()

	// Emits the chat body verbatim as the event-name suffix.
	emitter := loadShared(t, ctx, env, RuntimeJavaScript, "emitter", `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  onChatMessage(msg) { owncast.events.emit(msg.body, { n: 7 }); },
});`, []string{PermEmitEvent})
	defer emitter.Close(ctx)

	// Every honest name the emitter can produce, paired with the spelling an
	// attacker would need the host to accept.
	receiver := loadShared(t, ctx, env, RuntimeJavaScript, "receiver",
		receiverScript([]string{
			"emitter.ping", "ping",
			"emitter.victim.alert", "victim.alert",
			"emitter.chat.message.received", "chat.message.received",
			"emitter.orders.created", "orders.created",
		}), []string{PermChatSend})
	defer receiver.Close(ctx)

	d := NewLiveDispatcher(func() []*Loaded { return []*Loaded{emitter, receiver} })
	env.Emit = d.Dispatch

	// The receiver also subscribes to the real chat.message.received, so every
	// scenario sees exactly one host-originated delivery on top of the custom
	// one. Its absence would mean core dispatch broke; a second copy would mean
	// a plugin managed to forge it.
	const coreDelivery = "got:chat.message.received:undefined"

	for _, tc := range []struct {
		name string
		emit string
		want []string
	}{
		{
			name: "plain suffix is prefixed with the emitter slug",
			emit: "ping",
			want: []string{coreDelivery, "got:emitter.ping:7"},
		},
		{
			name: "dotted suffix keeps its hierarchy under the slug",
			emit: "orders.created",
			want: []string{coreDelivery, "got:emitter.orders.created:7"},
		},
		{
			name: "cannot impersonate another plugin's namespace",
			emit: "victim.alert",
			want: []string{coreDelivery, "got:emitter.victim.alert:7"},
		},
		{
			name: "cannot forge a core host event",
			emit: "chat.message.received",
			want: []string{coreDelivery, "got:emitter.chat.message.received:7"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d.Dispatch(ctx, EventChatMessageReceived, chatPayload("alice", tc.emit))
			got := drainSends(sends, mu)
			sort.Strings(tc.want)
			if !sameStrings(got, tc.want) {
				t.Fatalf("emit(%q)\n got:  %q\n want: %q", tc.emit, got, tc.want)
			}
		})
	}
}

// TestEmitFromPythonIsNamespaced covers the same guarantee through the Python
// engine. The prefix is applied host-side, so both SDKs must land on it
// without either one cooperating.
func TestEmitFromPythonIsNamespaced(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, sends, mu := captureEnv()

	// Tries to publish directly into another plugin's namespace.
	emitter := loadShared(t, ctx, env, RuntimePython, "py-emitter", `
@plugin.on_chat_message
def handle(msg):
    owncast.events.emit("victim.alert", {"n": 7})
`, []string{PermEmitEvent})
	defer emitter.Close(ctx)

	receiver := loadShared(t, ctx, env, RuntimeJavaScript, "receiver",
		receiverScript([]string{"py-emitter.victim.alert", "victim.alert"}),
		[]string{PermChatSend})
	defer receiver.Close(ctx)

	d := NewLiveDispatcher(func() []*Loaded { return []*Loaded{emitter, receiver} })
	env.Emit = d.Dispatch
	d.Dispatch(ctx, EventChatMessageReceived, chatPayload("alice", "hi"))

	want := []string{"got:py-emitter.victim.alert:7"}
	if got := drainSends(sends, mu); !sameStrings(got, want) {
		t.Fatalf("python emit\n got:  %q\n want: %q", got, want)
	}
}

// TestEmitWithoutPermissionDispatchesNothing pins the permission gate ahead of
// the prefixing: resolveCaller rejects the call, so no event reaches the
// dispatcher at all, prefixed or not.
func TestEmitWithoutPermissionDispatchesNothing(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, sends, mu := captureEnv()

	// No events.emit in the manifest. The call raises inside the guest, which
	// is swallowed here so the scenario reports on deliveries, not on the
	// guest-side error shape.
	emitter := loadShared(t, ctx, env, RuntimeJavaScript, "unprivileged", `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  onChatMessage() {
    try { owncast.events.emit("ping", { n: 7 }); } catch (e) {}
  },
});`, nil)
	defer emitter.Close(ctx)

	receiver := loadShared(t, ctx, env, RuntimeJavaScript, "receiver",
		receiverScript([]string{"unprivileged.ping", "ping"}),
		[]string{PermChatSend})
	defer receiver.Close(ctx)

	d := NewLiveDispatcher(func() []*Loaded { return []*Loaded{emitter, receiver} })
	env.Emit = d.Dispatch
	d.Dispatch(ctx, EventChatMessageReceived, chatPayload("alice", "hi"))

	if got := drainSends(sends, mu); len(got) != 0 {
		t.Fatalf("expected no delivery without events.emit, got %q", got)
	}
}

// TestSubscriptionsAreNotRewritten is the counterpart to the namespacing: only
// the emit side is touched. Subscriptions stay literal, otherwise a subscriber
// could never name the plugin it wants to hear from and cross-plugin events
// would stop being deliverable at all.
func TestSubscriptionsAreNotRewritten(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, _, _ := captureEnv()
	receiver := loadShared(t, ctx, env, RuntimeJavaScript, "receiver",
		receiverScript([]string{"emitter.ping"}), []string{PermChatSend})
	defer receiver.Close(ctx)

	if !subscribed(receiver.Manifest.Subscriptions.Notify, "emitter.ping") {
		t.Fatalf("subscription was rewritten, notify = %+v", receiver.Manifest.Subscriptions.Notify)
	}
	if subscribed(receiver.Manifest.Subscriptions.Notify, "receiver.emitter.ping") {
		t.Fatal("subscription was prefixed with the subscriber slug; cross-plugin events would never match")
	}
}

// TestEmitCannotComposeCoreEventName closes the gap the slug prefix leaves
// open: a plugin whose slug matches a core name's first segment could
// otherwise compose an exact built-in name and have it fan out as if the
// host originated it. The emitter here is slugged "chat" and emits
// "message.received"; the composed "chat.message.received" must be dropped,
// so the receiver sees only the genuine host dispatch.
func TestEmitCannotComposeCoreEventName(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, sends, mu := captureEnv()

	emitter := loadShared(t, ctx, env, RuntimeJavaScript, "chat", `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  onChatMessage() { owncast.events.emit("message.received", { n: 7 }); },
});`, []string{PermEmitEvent})
	defer emitter.Close(ctx)

	receiver := loadShared(t, ctx, env, RuntimeJavaScript, "receiver",
		receiverScript([]string{"chat.message.received"}), []string{PermChatSend})
	defer receiver.Close(ctx)

	d := NewLiveDispatcher(func() []*Loaded { return []*Loaded{emitter, receiver} })
	env.Emit = d.Dispatch
	d.Dispatch(ctx, EventChatMessageReceived, chatPayload("alice", "hi"))

	want := []string{"got:chat.message.received:undefined"}
	if got := drainSends(sends, mu); !sameStrings(got, want) {
		t.Fatalf("forged core event reached subscribers\n got:  %q\n want: %q", got, want)
	}
}

// TestEmitCannotSpoofAnotherPluginsEvents is issue #5093's scenario run
// against this design. The attacker wants events under the victim's name to
// fire with attacker-controlled payloads. Under host-side emit prefixing the
// attempt dispatches as "attacker.victim.foo", so "victim.foo" only ever
// fires when the victim itself emits - and registration squatting gains the
// attacker nothing, because subscribing to "victim.foo" only ever yields
// events the victim really sent.
func TestEmitCannotSpoofAnotherPluginsEvents(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, sends, mu := captureEnv()

	// Emits under its own name when told to.
	victim := loadShared(t, ctx, env, RuntimeJavaScript, "victim", `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  onChatMessage(msg) { if (msg.body === "victim") owncast.events.emit("foo", { n: 1 }); },
});`, []string{PermEmitEvent})
	defer victim.Close(ctx)

	// Tries to emit the victim's fully qualified name.
	attacker := loadShared(t, ctx, env, RuntimeJavaScript, "attacker", `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  onChatMessage(msg) { if (msg.body === "attacker") owncast.events.emit("victim.foo", { n: 2 }); },
});`, []string{PermEmitEvent})
	defer attacker.Close(ctx)

	listener := loadShared(t, ctx, env, RuntimeJavaScript, "listener",
		receiverScript([]string{"victim.foo", "attacker.victim.foo"}), []string{PermChatSend})
	defer listener.Close(ctx)

	d := NewLiveDispatcher(func() []*Loaded { return []*Loaded{victim, attacker, listener} })
	env.Emit = d.Dispatch

	d.Dispatch(ctx, EventChatMessageReceived, chatPayload("alice", "victim"))
	if got, want := drainSends(sends, mu), []string{"got:victim.foo:1"}; !sameStrings(got, want) {
		t.Fatalf("genuine victim emit\n got:  %q\n want: %q", got, want)
	}

	d.Dispatch(ctx, EventChatMessageReceived, chatPayload("alice", "attacker"))
	if got, want := drainSends(sends, mu), []string{"got:attacker.victim.foo:2"}; !sameStrings(got, want) {
		t.Fatalf("spoof attempt must surface under the attacker's own name only\n got:  %q\n want: %q", got, want)
	}
}
