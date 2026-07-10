package pluginhost

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	activityevents "github.com/owncast/owncast/services/activitypub/events"
	"github.com/owncast/owncast/services/dispatcher"
	"github.com/owncast/owncast/services/plugins"
	"github.com/owncast/owncast/services/webhooks"
)

func TestTranslateWebhookEvent_ChatMessageOnlyForUserMessages(t *testing.T) {
	ts := time.Date(2026, 5, 27, 12, 0, 0, 0, time.UTC)
	evt := webhooks.WebhookEvent{
		Type: models.MessageSent,
		EventData: &webhooks.WebhookChatMessage{
			ID: "m1",
			// Body is the HTML-rendered form the chat client sees; RawBody is
			// what plugins receive. Set both, to mirror the production payload.
			Body:      "<p>hello</p>",
			RawBody:   "hello",
			Timestamp: &ts,
			User:      &models.User{ID: "u1", DisplayName: "alice", Scopes: []string{"MODERATOR"}, Authenticated: true},
			ClientID:  42,
		},
	}
	out := translateWebhookEvent(evt)
	if len(out) != 1 {
		t.Fatalf("expected 1 event, got %d", len(out))
	}
	if out[0].eventType != plugins.EventChatMessageReceived {
		t.Errorf("eventType = %q want %q", out[0].eventType, plugins.EventChatMessageReceived)
	}
	msg, ok := out[0].payload.(pluginChatMessage)
	if !ok {
		t.Fatalf("payload type = %T want pluginChatMessage", out[0].payload)
	}
	if msg.ID != "m1" || msg.Body != "hello" || msg.ClientID != 42 {
		t.Errorf("unexpected message payload: %+v", msg)
	}
	// The full sender identity must ride along so plugins can do reliable
	// per-user state and scope-gated moderation without a display-name lookup.
	if msg.User == nil {
		t.Fatalf("message user identity was dropped")
	}
	if msg.User.ID != "u1" || msg.User.DisplayName != "alice" || !msg.User.IsAuthenticated || len(msg.User.Scopes) != 1 {
		t.Errorf("unexpected user identity: %+v", msg.User)
	}
	if msg.Timestamp != "2026-05-27T12:00:00Z" {
		t.Errorf("timestamp = %q", msg.Timestamp)
	}
}

func TestTranslatePluginEvent_FediverseEvents(t *testing.T) {
	actor := activityevents.FediverseActor{
		Name:   "Alice Example",
		Handle: "@@alice@social.example",
		URL:    "https://social.example/users/alice",
		Image:  "https://social.example/users/alice/avatar.png",
	}
	engagement := func(target string) *activityevents.FediverseEngagementEvent {
		return &activityevents.FediverseEngagementEvent{
			Actor:  actor,
			Target: &activityevents.FediverseTarget{URL: target},
		}
	}
	post := &activityevents.FediverseInboundPostEvent{
		Actor:       actor,
		Content:     "<p>Hello <a href=\"https://owncast.example\">@streamer</a></p>",
		ContentText: "Hello @streamer",
		URL:         "https://social.example/posts/123",
		PostedAt:    "2026-07-10T12:00:00Z",
		InReplyTo:   "https://owncast.example/federation/post/parent",
		Attachments: []activityevents.FediverseAttachment{{
			URL:       "https://social.example/media/photo.jpg",
			MediaType: "image/jpeg",
			Alt:       "A photo",
		}},
		Language: "en",
	}

	tests := []struct {
		name      string
		event     dispatcher.Event
		eventType string
		json      string
	}{
		{
			name: "follow",
			event: dispatcher.Event{Payload: webhooks.WebhookEvent{
				Type: models.FediverseEngagementFollow,
				EventData: &webhooks.WebhookFediverseEngagementFollowEventData{
					Name:      actor.Name,
					Username:  "alice@social.example",
					Image:     actor.Image,
					ServerURL: "https://owncast.example",
				},
			}},
			eventType: "fediverse.follow",
			json:      `{"actor":{"name":"Alice Example","handle":"@alice@social.example","image":"https://social.example/users/alice/avatar.png"}}`,
		},
		{
			name:      "like",
			event:     dispatcher.Event{Type: models.FediverseEngagementLike, Payload: engagement("https://owncast.example/federation/post/liked")},
			eventType: "fediverse.like",
			json:      `{"actor":{"name":"Alice Example","handle":"@alice@social.example","url":"https://social.example/users/alice","image":"https://social.example/users/alice/avatar.png"},"target":{"url":"https://owncast.example/federation/post/liked"}}`,
		},
		{
			name:      "repost",
			event:     dispatcher.Event{Type: models.FediverseEngagementRepost, Payload: engagement("https://owncast.example/federation/post/reposted")},
			eventType: "fediverse.repost",
			json:      `{"actor":{"name":"Alice Example","handle":"@alice@social.example","url":"https://social.example/users/alice","image":"https://social.example/users/alice/avatar.png"},"target":{"url":"https://owncast.example/federation/post/reposted"}}`,
		},
		{
			name:      "mention",
			event:     dispatcher.Event{Type: models.FediverseMention, Payload: post},
			eventType: "fediverse.mention",
			json:      `{"actor":{"name":"Alice Example","handle":"@alice@social.example","url":"https://social.example/users/alice","image":"https://social.example/users/alice/avatar.png"},"content":"\u003cp\u003eHello \u003ca href=\"https://owncast.example\"\u003e@streamer\u003c/a\u003e\u003c/p\u003e","contentText":"Hello @streamer","url":"https://social.example/posts/123","postedAt":"2026-07-10T12:00:00Z","inReplyTo":"https://owncast.example/federation/post/parent","attachments":[{"url":"https://social.example/media/photo.jpg","mediaType":"image/jpeg","alt":"A photo"}],"language":"en"}`,
		},
		{
			name:      "reply",
			event:     dispatcher.Event{Type: models.FediverseReply, Payload: post},
			eventType: "fediverse.reply",
			json:      `{"actor":{"name":"Alice Example","handle":"@alice@social.example","url":"https://social.example/users/alice","image":"https://social.example/users/alice/avatar.png"},"content":"\u003cp\u003eHello \u003ca href=\"https://owncast.example\"\u003e@streamer\u003c/a\u003e\u003c/p\u003e","contentText":"Hello @streamer","url":"https://social.example/posts/123","postedAt":"2026-07-10T12:00:00Z","inReplyTo":"https://owncast.example/federation/post/parent","attachments":[{"url":"https://social.example/media/photo.jpg","mediaType":"image/jpeg","alt":"A photo"}],"language":"en"}`,
		},
		{
			name:      "quote",
			event:     dispatcher.Event{Type: models.FediverseEngagementQuote, Payload: engagement("https://owncast.example/federation/post/quoted")},
			eventType: "fediverse.quote",
			json:      `{"actor":{"name":"Alice Example","handle":"@alice@social.example","url":"https://social.example/users/alice","image":"https://social.example/users/alice/avatar.png"},"target":{"url":"https://owncast.example/federation/post/quoted"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := translatePluginEvent(tt.event)
			if len(out) != 1 {
				t.Fatalf("expected 1 event, got %d", len(out))
			}
			if out[0].eventType != tt.eventType {
				t.Fatalf("eventType = %q want %q", out[0].eventType, tt.eventType)
			}
			payload, err := json.Marshal(out[0].payload)
			if err != nil {
				t.Fatalf("marshal payload: %v", err)
			}
			if string(payload) != tt.json {
				t.Errorf("payload JSON = %s want %s", payload, tt.json)
			}
		})
	}
}

func TestTranslatePluginEvent_FediverseActivityPassesThroughRawJSON(t *testing.T) {
	raw := json.RawMessage(`{"type":"Like","actor":"https://social.example/users/alice"}`)
	out := translatePluginEvent(dispatcher.Event{Type: models.FediverseActivity, Payload: raw})
	if len(out) != 1 {
		t.Fatalf("expected 1 event, got %d", len(out))
	}
	if out[0].eventType != "fediverse.activity" {
		t.Fatalf("eventType = %q want %q", out[0].eventType, "fediverse.activity")
	}
	payload, err := json.Marshal(out[0].payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	if string(payload) != string(raw) {
		t.Errorf("payload JSON = %s want %s", payload, raw)
	}
}

func TestTranslatePluginEvent_FediverseWrongPayloadTypes(t *testing.T) {
	tests := []dispatcher.Event{
		{Payload: webhooks.WebhookEvent{Type: models.FediverseEngagementFollow, EventData: struct{}{}}},
		{Type: models.FediverseEngagementLike, Payload: struct{}{}},
		{Type: models.FediverseEngagementRepost, Payload: struct{}{}},
		{Type: models.FediverseMention, Payload: struct{}{}},
		{Type: models.FediverseReply, Payload: struct{}{}},
		{Type: models.FediverseEngagementQuote, Payload: struct{}{}},
		{Type: models.FediverseActivity, Payload: []byte(`{}`)},
	}
	for _, event := range tests {
		if out := translatePluginEvent(event); len(out) != 0 {
			t.Errorf("%s with payload %T produced %d events", event.Type, event.Payload, len(out))
		}
	}
}

func TestTranslateWebhookEvent_SkipsBotAuthoredMessages(t *testing.T) {
	// A message a plugin posted (under its bot identity) must not be delivered
	// back to plugins, or chat-reacting plugins would echo-loop forever.
	evt := webhooks.WebhookEvent{
		Type: models.MessageSent,
		EventData: &webhooks.WebhookChatMessage{
			ID:   "b1",
			Body: "echo",
			User: &models.User{DisplayName: "echo-bot", IsBot: true},
		},
	}
	if out := translateWebhookEvent(evt); len(out) != 0 {
		t.Errorf("bot-authored message should produce no plugin events, got %d", len(out))
	}
}

func TestTranslateWebhookEvent_SystemMessageProducesNothing(t *testing.T) {
	// A plugin's own chat.send posts a system message; it must not become a
	// chat.message.received event (no feedback loop).
	evt := webhooks.WebhookEvent{
		Type:      models.SystemMessageSent,
		EventData: &webhooks.WebhookChatMessage{ID: "s1", Body: "system"},
	}
	if out := translateWebhookEvent(evt); len(out) != 0 {
		t.Errorf("system message should produce no plugin events, got %d", len(out))
	}
}

func TestTranslateWebhookEvent_UserJoined(t *testing.T) {
	evt := webhooks.WebhookEvent{
		Type: models.UserJoined,
		EventData: &webhooks.WebhookUserJoinedEventData{
			User: &models.User{ID: "u1", DisplayName: "bob", DisplayColor: 5, IsBot: true, Scopes: []string{"MODERATOR"}},
		},
	}
	out := translateWebhookEvent(evt)
	if len(out) != 1 || out[0].eventType != plugins.EventChatUserJoined {
		t.Fatalf("unexpected events: %+v", out)
	}
	user, ok := out[0].payload.(plugins.HostUser)
	if !ok {
		t.Fatalf("payload type = %T want plugins.HostUser", out[0].payload)
	}
	if user.ID != "u1" || user.DisplayName != "bob" || user.DisplayColor != 5 || !user.IsBot || len(user.Scopes) != 1 {
		t.Errorf("unexpected user payload: %+v", user)
	}
}

func TestTranslateWebhookEvent_RenameUsesPreviousName(t *testing.T) {
	evt := webhooks.WebhookEvent{
		Type: models.UserNameChanged,
		EventData: &webhooks.WebhookNameChangeEventData{
			NewName: "newname",
			User:    &models.User{ID: "u1", DisplayName: "oldname", PreviousNames: []string{"first", "oldname"}},
		},
	}
	out := translateWebhookEvent(evt)
	if len(out) != 1 || out[0].eventType != plugins.EventChatUserRenamed {
		t.Fatalf("unexpected events: %+v", out)
	}
	rename, ok := out[0].payload.(pluginChatUserRename)
	if !ok {
		t.Fatalf("payload type = %T want pluginChatUserRename", out[0].payload)
	}
	if rename.User.DisplayName != "newname" {
		t.Errorf("renamed user display name = %q want newname", rename.User.DisplayName)
	}
	if rename.PreviousName != "oldname" {
		t.Errorf("previousName = %q want oldname", rename.PreviousName)
	}
}

func TestTranslateWebhookEvent_VisibilityFansOutPerMessage(t *testing.T) {
	evt := webhooks.WebhookEvent{
		Type: models.VisibiltyToggled,
		EventData: &webhooks.WebhookVisibilityToggleEventData{
			Visible:    false,
			MessageIDs: []string{"a", "b", "c"},
		},
	}
	out := translateWebhookEvent(evt)
	if len(out) != 3 {
		t.Fatalf("expected 3 moderation events, got %d", len(out))
	}
	for i, id := range []string{"a", "b", "c"} {
		mod, ok := out[i].payload.(pluginChatMessageModeration)
		if !ok {
			t.Fatalf("payload[%d] type = %T", i, out[i].payload)
		}
		if mod.MessageID != id || mod.Visible {
			t.Errorf("event %d = %+v want id %q visible false", i, mod, id)
		}
	}
}

func TestTranslateWebhookEvent_StreamLifecycle(t *testing.T) {
	ts := time.Date(2026, 5, 27, 9, 30, 0, 0, time.UTC)
	started := webhooks.WebhookEvent{
		Type: models.StreamStarted,
		EventData: map[string]interface{}{
			"streamTitle": "Live",
			"summary":     "a show",
			"timestamp":   ts,
		},
	}
	out := translateWebhookEvent(started)
	if len(out) != 1 || out[0].eventType != plugins.EventStreamStarted {
		t.Fatalf("unexpected events: %+v", out)
	}
	life, ok := out[0].payload.(pluginStreamLifecycleEvent)
	if !ok {
		t.Fatalf("payload type = %T", out[0].payload)
	}
	if life.StartedAt != "2026-05-27T09:30:00Z" || life.Title != "Live" || life.Summary != "a show" {
		t.Errorf("unexpected lifecycle payload: %+v", life)
	}
	if life.StoppedAt != "" {
		t.Errorf("started event should not set stoppedAt: %+v", life)
	}
}

func TestFilteredBody(t *testing.T) {
	cases := []struct {
		name     string
		final    any
		fallback string
		want     string
	}{
		{"unchanged passes through the struct", pluginChatMessage{Body: "hello"}, "hello", "hello"},
		{"modified comes back as a decoded map", map[string]interface{}{"body": "HELLO"}, "hello", "HELLO"},
		{"map without body uses fallback", map[string]interface{}{"id": "x"}, "hello", "hello"},
		{"unexpected shape uses fallback", 42, "hello", "hello"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := filteredBody(tc.final, tc.fallback); got != tc.want {
				t.Errorf("filteredBody = %q want %q", got, tc.want)
			}
		})
	}
}

func TestTranslateWebhookEvent_TitleChange(t *testing.T) {
	evt := webhooks.WebhookEvent{
		Type:      models.StreamTitleUpdated,
		EventData: map[string]interface{}{"streamTitle": "New Title"},
	}
	out := translateWebhookEvent(evt)
	if len(out) != 1 || out[0].eventType != plugins.EventStreamTitleChanged {
		t.Fatalf("unexpected events: %+v", out)
	}
	change, ok := out[0].payload.(pluginStreamTitleChange)
	if !ok {
		t.Fatalf("payload type = %T", out[0].payload)
	}
	if change.To != "New Title" {
		t.Errorf("title change To = %q want New Title", change.To)
	}
}
