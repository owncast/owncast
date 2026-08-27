package plugins

import (
	"context"
	"strings"
	"testing"
)

func TestCustomEventHooksBelongToDeclaringPlugin(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, sends, mu := captureEnv()
	emitter := loadShared(t, ctx, env, RuntimeJavaScript, "emitter", `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  onChatMessage() { owncast.events.emit("plugin-one.foo", {}); }
});`, []string{PermEmitEvent})
	defer emitter.Close(ctx)

	pluginOne := loadShared(t, ctx, env, RuntimeJavaScript, "plugin-one", `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  on: { foo() { owncast.chat.send("plugin-one"); } }
});`, []string{PermChatSend})
	defer pluginOne.Close(ctx)

	pluginTwo := loadShared(t, ctx, env, RuntimeJavaScript, "plugin-two", `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  on: { "plugin-one.foo"() { owncast.chat.send("plugin-two"); } }
});`, []string{PermChatSend})
	defer pluginTwo.Close(ctx)

	if !subscribed(pluginOne.Manifest.Subscriptions.Notify, "plugin-one.foo") {
		t.Fatalf("plugin-one subscriptions = %+v", pluginOne.Manifest.Subscriptions.Notify)
	}
	if subscribed(pluginTwo.Manifest.Subscriptions.Notify, "plugin-one.foo") ||
		!subscribed(pluginTwo.Manifest.Subscriptions.Notify, "plugin-two.plugin-one.foo") {
		t.Fatalf("plugin-two subscriptions = %+v", pluginTwo.Manifest.Subscriptions.Notify)
	}

	d := NewLiveDispatcher(func() []*Loaded { return []*Loaded{emitter, pluginOne, pluginTwo} })
	env.Emit = d.Dispatch
	d.Dispatch(ctx, EventChatMessageReceived, chatPayload("alice", "run"))

	mu.Lock()
	if len(*sends) != 1 || (*sends)[0].Text != "plugin-one" {
		t.Fatalf("plugin-one.foo deliveries = %+v", *sends)
	}
	*sends = nil
	mu.Unlock()

	// The attempted squat is still a valid hook, but only under plugin-two's
	// own namespace.
	d.Dispatch(ctx, "plugin-two.plugin-one.foo", map[string]any{})
	mu.Lock()
	defer mu.Unlock()
	if len(*sends) != 1 || (*sends)[0].Text != "plugin-two" {
		t.Fatalf("plugin-two.plugin-one.foo deliveries = %+v", *sends)
	}
}

func TestCustomEventHookCannotCollideWithBuiltIn(t *testing.T) {
	subs := Subscriptions{Notify: []Subscription{{Event: "message.received"}}}
	if err := namespaceCustomSubscriptions("chat", &subs); err == nil ||
		!strings.Contains(err.Error(), EventChatMessageReceived) {
		t.Fatalf("expected %q collision, got %v", EventChatMessageReceived, err)
	}
}

func TestBuiltInSubscriptionKeepsCanonicalName(t *testing.T) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	env, sends, mu := captureEnv()
	loaded := loadShared(t, ctx, env, RuntimeJavaScript, "chat", `
const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  onChatMessage() { owncast.chat.send("received"); }
});`, []string{PermChatSend})
	defer loaded.Close(ctx)

	d := NewLiveDispatcher(func() []*Loaded { return []*Loaded{loaded} })
	d.Dispatch(ctx, EventChatMessageReceived, chatPayload("alice", "run"))

	mu.Lock()
	defer mu.Unlock()
	if len(*sends) != 1 || (*sends)[0].Text != "received" {
		t.Fatalf("%s deliveries for chat slug = %+v", EventChatMessageReceived, *sends)
	}
}
