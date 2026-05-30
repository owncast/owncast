package pluginhost

import (
	"context"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/chat/events"
	"github.com/owncast/owncast/services/dispatcher"
	"github.com/owncast/owncast/services/plugins"
	"github.com/owncast/owncast/services/webhooks"
)

// Plugin-facing event payload shapes. These mirror the plugin SDK's
// TypeScript interfaces (ChatMessage, ChatUser, …) so a plugin's typed
// handlers receive exactly the documented JSON.

type pluginChatMessage struct {
	ID        string `json:"id"`
	User      string `json:"user"`
	Body      string `json:"body"`
	Timestamp string `json:"timestamp"`
}

type pluginChatUser struct {
	ID              string   `json:"id"`
	DisplayName     string   `json:"displayName"`
	IsBot           bool     `json:"isBot,omitempty"`
	IsAuthenticated bool     `json:"isAuthenticated,omitempty"`
	Scopes          []string `json:"scopes,omitempty"`
}

type pluginChatUserRename struct {
	User         pluginChatUser `json:"user"`
	PreviousName string         `json:"previousName"`
}

type pluginChatMessageModeration struct {
	MessageID string          `json:"messageId"`
	Visible   bool            `json:"visible"`
	Moderator *pluginChatUser `json:"moderator,omitempty"`
}

type pluginStreamLifecycleEvent struct {
	StartedAt string `json:"startedAt,omitempty"`
	StoppedAt string `json:"stoppedAt,omitempty"`
	Title     string `json:"title,omitempty"`
	Summary   string `json:"summary,omitempty"`
}

type pluginStreamTitleChange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// pluginEvent is one translated event ready to dispatch: a plugin event type
// and its SDK-shaped payload.
type pluginEvent struct {
	eventType string
	payload   any
}

// newPluginChatFilter returns a dispatcher.Filter that runs the plugin
// filterChatMessage chain on an inbound chat message, rewriting the message
// body in place and reporting whether it survived. Plugin errors are
// fail-open inside the plugin dispatcher, so a broken filter never blocks chat.
func newPluginChatFilter(pluginDispatcher *plugins.Dispatcher) dispatcher.Filter {
	return func(ctx context.Context, e dispatcher.Event) bool {
		msg, ok := e.Payload.(*events.UserMessageEvent)
		if !ok {
			return true // not a chat message we filter; let it pass
		}
		user := ""
		if msg.User != nil {
			user = msg.User.DisplayName
		}
		// Carry the timestamp through so filter plugins that gate on it
		// (slow-mode, rate-limit, etc.) can compare elapsed time. Nano-
		// precision matters: extism-js's Date.now() returns a frozen
		// WASI-default value, so plugins can't get wall-clock time on
		// their own — msg.timestamp is the only source of real-time
		// resolution they have. RFC3339 second-precision is too coarse
		// for sub-second rate limits.
		timestamp := ""
		if !msg.Timestamp.IsZero() {
			timestamp = msg.Timestamp.UTC().Format(time.RFC3339Nano)
		}
		final, allowed, _ := pluginDispatcher.Filter(ctx, plugins.EventChatMessageReceived, pluginChatMessage{ID: msg.ID, User: user, Body: msg.Body, Timestamp: timestamp})
		if !allowed {
			return false
		}
		// A plugin holding only chat.filter must not be able to
		// inject raw HTML into the broadcast. When the filter
		// returned a modified body, run it back through
		// RenderAndSanitize so any new <script>/<iframe>/etc the
		// plugin tried to introduce gets stripped by the same
		// bluemonday policy userMessageSent already applied to the
		// original message. Skip the round-trip when nothing
		// changed so a no-op filter doesn't re-render the body.
		//
		// Assign the result to msg.Body directly rather than via
		// RenderAndSanitizeMessageBody, because that helper resets
		// RawBody from Body first and would clobber whatever we
		// stored.
		next := filteredBody(final, msg.Body)
		if next != msg.Body {
			msg.Body = events.RenderAndSanitize(next)
		}
		return true
	}
}

// filteredBody extracts the body from a filter chain's result. The result is
// the original pluginChatMessage when no plugin modified it, or a decoded JSON
// object when one did; fallback covers any unexpected shape.
func filteredBody(final any, fallback string) string {
	switch v := final.(type) {
	case pluginChatMessage:
		return v.Body
	case map[string]interface{}:
		if body, ok := v["body"].(string); ok {
			return body
		}
	}
	return fallback
}

// newPluginEventListener returns a dispatcher.Listener that translates each
// Owncast event (carried as a webhooks.WebhookEvent) into the plugin SDK's
// payload shape and dispatches it to subscribed plugins. Dispatch runs on its
// own goroutine so a slow plugin never blocks the event source (the chat hot
// path).
func newPluginEventListener(pluginDispatcher *plugins.Dispatcher) dispatcher.Listener {
	return func(ctx context.Context, e dispatcher.Event) {
		webhookEvent, ok := e.Payload.(webhooks.WebhookEvent)
		if !ok {
			return
		}
		// Notifications are fire-and-forget: detach from the publisher's
		// context so the dispatch isn't cancelled when the publishing call
		// returns. Per-plugin timeouts still apply inside Dispatch.
		ctx = context.WithoutCancel(ctx)
		for _, pe := range translateWebhookEvent(webhookEvent) {
			// Notify-only on purpose: the inbound-chat path already ran
			// the filter chain through newPluginChatFilter before the
			// event was broadcast. Calling Dispatch (filter + notify)
			// here would re-run every plugin's on_filter on a message
			// that already survived once, doubling work and double-
			// triggering rate-limit logic (slow-mode, etc.).
			go pluginDispatcher.Notify(ctx, pe.eventType, pe.payload)
		}
	}
}

// translateWebhookEvent maps an Owncast webhook event onto the plugin events
// it should produce (zero, one, or — for a multi-message moderation toggle —
// several). It's pure so the mapping can be tested without a live dispatcher.
//
// Only genuine user chat messages (models.MessageSent) become
// chat.message.received — system messages and actions (including a plugin's
// own chat.send output) are intentionally excluded, so plugins don't react to
// their own posts.
func translateWebhookEvent(evt webhooks.WebhookEvent) []pluginEvent {
	switch evt.Type {
	case models.MessageSent, models.UserJoined, models.UserParted, models.UserNameChanged, models.VisibiltyToggled:
		return translateChatEvent(evt)
	case models.StreamStarted, models.StreamStopped, models.StreamTitleUpdated:
		return translateStreamEvent(evt)
	}
	return nil
}

func translateChatEvent(evt webhooks.WebhookEvent) []pluginEvent {
	switch evt.Type {
	case models.MessageSent:
		return chatMessageEvent(evt)

	case models.UserJoined:
		data, ok := evt.EventData.(*webhooks.WebhookUserJoinedEventData)
		if !ok {
			return nil
		}
		return []pluginEvent{{plugins.EventChatUserJoined, toPluginChatUser(data.User)}}

	case models.UserParted:
		data, ok := evt.EventData.(*webhooks.WebhookUserPartEventData)
		if !ok {
			return nil
		}
		return []pluginEvent{{plugins.EventChatUserParted, toPluginChatUser(data.User)}}

	case models.UserNameChanged:
		data, ok := evt.EventData.(*webhooks.WebhookNameChangeEventData)
		if !ok {
			return nil
		}
		user := toPluginChatUser(data.User)
		user.DisplayName = data.NewName
		return []pluginEvent{{plugins.EventChatUserRenamed, pluginChatUserRename{
			User:         user,
			PreviousName: previousName(data.User),
		}}}

	case models.VisibiltyToggled:
		data, ok := evt.EventData.(*webhooks.WebhookVisibilityToggleEventData)
		if !ok {
			return nil
		}
		var moderator *pluginChatUser
		if data.User != nil {
			m := toPluginChatUser(data.User)
			moderator = &m
		}
		// Owncast toggles a set of IDs at once; the SDK payload is
		// per-message, so fan one event out per affected message.
		out := make([]pluginEvent, 0, len(data.MessageIDs))
		for _, id := range data.MessageIDs {
			out = append(out, pluginEvent{plugins.EventChatMessageModerated, pluginChatMessageModeration{
				MessageID: id,
				Visible:   data.Visible,
				Moderator: moderator,
			}})
		}
		return out
	}
	return nil
}

func chatMessageEvent(evt webhooks.WebhookEvent) []pluginEvent {
	data, ok := evt.EventData.(*webhooks.WebhookChatMessage)
	if !ok {
		return nil
	}
	// Don't deliver messages authored by a bot (e.g. another plugin's
	// chat.send, posted under its bot identity) back to plugins. This prevents
	// echo loops: a plugin that replies to chat would otherwise be re-triggered
	// by its own reply, forever.
	if data.User != nil && data.User.IsBot {
		return nil
	}
	// Use RawBody (what the user actually typed), not Body (the HTML-rendered
	// form like `<p>!broadcaster</p>`). Plugins doing command matching or
	// content analysis want the raw text; the chat client handles rendering.
	msg := pluginChatMessage{ID: data.ID, Body: data.RawBody, Timestamp: formatTimePtr(data.Timestamp)}
	if data.User != nil {
		msg.User = data.User.DisplayName
	}
	return []pluginEvent{{plugins.EventChatMessageReceived, msg}}
}

func translateStreamEvent(evt webhooks.WebhookEvent) []pluginEvent {
	switch evt.Type {
	case models.StreamStarted:
		return []pluginEvent{{plugins.EventStreamStarted, streamLifecycleEvent(evt.EventData, true)}}

	case models.StreamStopped:
		return []pluginEvent{{plugins.EventStreamStopped, streamLifecycleEvent(evt.EventData, false)}}

	case models.StreamTitleUpdated:
		m, ok := evt.EventData.(map[string]interface{})
		if !ok {
			return nil
		}
		to, _ := m["streamTitle"].(string)
		// Owncast's title-changed event carries only the new title.
		return []pluginEvent{{plugins.EventStreamTitleChanged, pluginStreamTitleChange{To: to}}}
	}
	return nil
}

func toPluginChatUser(u *models.User) pluginChatUser {
	if u == nil {
		return pluginChatUser{}
	}
	return pluginChatUser{
		ID:              u.ID,
		DisplayName:     u.DisplayName,
		IsBot:           u.IsBot,
		IsAuthenticated: u.Authenticated,
		Scopes:          u.Scopes,
	}
}

// previousName returns the user's most recent prior display name, the closest
// available value for a rename's "from". Empty when there's no history.
func previousName(u *models.User) string {
	if u == nil || len(u.PreviousNames) == 0 {
		return ""
	}
	return u.PreviousNames[len(u.PreviousNames)-1]
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.UTC().Format(time.RFC3339Nano)
}

// streamLifecycleEvent builds a stream.started/stopped payload from the
// map-shaped EventData the webhooks stream builder emits.
func streamLifecycleEvent(data interface{}, started bool) pluginStreamLifecycleEvent {
	out := pluginStreamLifecycleEvent{}
	m, ok := data.(map[string]interface{})
	if !ok {
		return out
	}
	out.Title, _ = m["streamTitle"].(string)
	out.Summary, _ = m["summary"].(string)
	timestamp := ""
	if t, ok := m["timestamp"].(time.Time); ok {
		timestamp = t.UTC().Format(time.RFC3339Nano)
	}
	if started {
		out.StartedAt = timestamp
	} else {
		out.StoppedAt = timestamp
	}
	return out
}
