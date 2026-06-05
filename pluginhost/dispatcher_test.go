package pluginhost

import (
	"context"
	"strings"
	"testing"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/chat/events"
	"github.com/owncast/owncast/services/dispatcher"
	"github.com/owncast/owncast/services/plugins"
	"github.com/owncast/owncast/services/webhooks"
)

// The host's chat filter, registered on a real dispatcher, passes messages
// through unchanged when no plugins are loaded — and the dispatcher delivers
// the inbound message to it. Exercises the real chat -> dispatcher -> plugin
// filter seam (full rewriting needs a wasm plugin, covered separately).
func TestPluginChatFilter_PassesThroughWithNoPlugins(t *testing.T) {
	d := dispatcher.New()
	d.AddFilter(newPluginChatFilter(plugins.NewDispatcher(nil)))

	msg := &events.UserMessageEvent{}
	msg.Body = "hello"

	allow := d.ApplyFilters(context.Background(), dispatcher.Event{Type: models.MessageSent, Payload: msg})
	if !allow {
		t.Fatal("message should pass with no plugins loaded")
	}
	if msg.Body != "hello" {
		t.Errorf("body changed unexpectedly: %q", msg.Body)
	}
}

// The chat filter only acts on chat-message payloads; anything else passes.
func TestPluginChatFilter_IgnoresNonChatPayload(t *testing.T) {
	d := dispatcher.New()
	d.AddFilter(newPluginChatFilter(plugins.NewDispatcher(nil)))

	if !d.ApplyFilters(context.Background(), dispatcher.Event{Type: "other", Payload: "not a chat message"}) {
		t.Fatal("non-chat payloads must pass through unchanged")
	}
}

// The host's event listener, registered on a real dispatcher, handles both a
// translatable webhook event and an unrelated payload without panicking. (The
// plugin notify side fans out to loaded plugins, which needs a wasm plugin to
// observe; here we verify the dispatcher -> listener wiring and payload guard.)
func TestPluginEventListener_HandlesPublishedEvents(t *testing.T) {
	d := dispatcher.New()
	d.AddListener(newPluginEventListener(plugins.NewDispatcher(nil)))

	// Unrelated payload type: ignored.
	d.Publish(context.Background(), dispatcher.Event{Type: "x", Payload: 42})

	// A real chat-message webhook event: translated and dispatched to (no) plugins.
	d.Publish(context.Background(), dispatcher.Event{
		Type: models.MessageSent,
		Payload: webhooks.WebhookEvent{
			Type:      models.MessageSent,
			EventData: &webhooks.WebhookChatMessage{ID: "m1", Body: "hi", User: &models.User{DisplayName: "alice"}},
		},
	})
}

// The dispatcher applies a filter chain to an inbound chat message exactly as
// chat relies on: filters mutate the message body in place and can drop it.
// This stands in for a plugin's filterChatMessage (which needs wasm to run)
// and proves the mechanism the chat hot path depends on.
func TestDispatcherChatFilterContract(t *testing.T) {
	d := dispatcher.New()
	// Rewrite: uppercase the body.
	d.AddFilter(func(_ context.Context, e dispatcher.Event) bool {
		if m, ok := e.Payload.(*events.UserMessageEvent); ok {
			m.Body = strings.ToUpper(m.Body)
		}
		return true
	})
	// Drop: reject anything containing "SPAM" (post-uppercase).
	d.AddFilter(func(_ context.Context, e dispatcher.Event) bool {
		m, ok := e.Payload.(*events.UserMessageEvent)
		return !ok || !strings.Contains(m.Body, "SPAM")
	})

	t.Run("rewrites body in place", func(t *testing.T) {
		msg := &events.UserMessageEvent{}
		msg.Body = "hello"
		if !d.ApplyFilters(context.Background(), dispatcher.Event{Payload: msg}) {
			t.Fatal("expected message to pass")
		}
		if msg.Body != "HELLO" {
			t.Errorf("body = %q, want HELLO", msg.Body)
		}
	})

	t.Run("drops a rejected message", func(t *testing.T) {
		msg := &events.UserMessageEvent{}
		msg.Body = "spam"
		if d.ApplyFilters(context.Background(), dispatcher.Event{Payload: msg}) {
			t.Fatal("expected message to be dropped")
		}
	})
}
