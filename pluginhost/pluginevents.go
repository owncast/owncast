package pluginhost

import (
	"context"
	"time"

	"github.com/owncast/owncast/models"
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

// newPluginChatFilter returns a chat message-filter hook that runs the
// plugin filterChatMessage chain synchronously and returns the (possibly
// rewritten) body plus whether the message survived. Plugin errors are
// fail-open inside the dispatcher, so a broken filter never blocks chat.
func newPluginChatFilter(dispatcher *plugins.Dispatcher) func(messageID, user, body string) (string, bool) {
	return func(messageID, user, body string) (string, bool) {
		msg := pluginChatMessage{ID: messageID, User: user, Body: body}
		final, allowed, _ := dispatcher.Filter(context.Background(), plugins.EventChatMessageReceived, msg)
		if !allowed {
			return "", false
		}
		return filteredBody(final, body), true
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

// newPluginEventListener returns a webhooks event listener that translates
// each Owncast event into the plugin SDK's payload shape and dispatches it to
// subscribed plugins. Dispatch runs on its own goroutine so a slow plugin
// never blocks the event source (the chat hot path).
func newPluginEventListener(dispatcher *plugins.Dispatcher) func(webhooks.WebhookEvent) {
	return func(evt webhooks.WebhookEvent) {
		for _, e := range translateWebhookEvent(evt) {
			go dispatcher.Dispatch(context.Background(), e.eventType, e.payload)
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
		data, ok := evt.EventData.(*webhooks.WebhookChatMessage)
		if !ok {
			return nil
		}
		msg := pluginChatMessage{ID: data.ID, Body: data.Body, Timestamp: formatTimePtr(data.Timestamp)}
		if data.User != nil {
			msg.User = data.User.DisplayName
		}
		return []pluginEvent{{plugins.EventChatMessageReceived, msg}}

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
	return t.UTC().Format(time.RFC3339)
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
		timestamp = t.UTC().Format(time.RFC3339)
	}
	if started {
		out.StartedAt = timestamp
	} else {
		out.StoppedAt = timestamp
	}
	return out
}
