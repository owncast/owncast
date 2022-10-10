package webhooks

import (
	"testing"
	"time"

	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/models"
)

func TestSendChatEvent(t *testing.T) {
	timestamp := time.Unix(72, 6).UTC()
	user := user.User{
		ID:              "user id",
		DisplayName:     "display name",
		DisplayColor:    4,
		CreatedAt:       time.Unix(3, 26).UTC(),
		DisabledAt:      nil,
		PreviousNames:   []string{"somebody"},
		NameChangedAt:   nil,
		Scopes:          []string{},
		IsBot:           false,
		AuthenticatedAt: nil,
		Authenticated:   false,
	}

	checkPayload(t, models.MessageSent, func() {
		SendChatEvent(&events.UserMessageEvent{
			Event: events.Event{
				Type:      events.MessageSent,
				ID:        "id",
				Timestamp: timestamp,
			},
			UserEvent: events.UserEvent{
				User:     &user,
				ClientID: 51,
				HiddenAt: nil,
			},
			MessageEvent: events.MessageEvent{
				OutboundEvent: nil,
				Body:          "body",
				RawBody:       "raw body",
			},
		})
	}, `{
		"body": "body",
		"clientId": 51,
		"id": "id",
		"rawBody": "raw body",
		"timestamp": "1970-01-01T00:01:12.000000006Z",
		"user": {
			"authenticated": false,
			"createdAt": "1970-01-01T00:00:03.000000026Z",
			"displayColor": 4,
			"displayName": "display name",
			"id": "user id",
			"isBot": false,
			"previousNames": ["somebody"]
		},
		"visible": true
	}`)
}

func TestSendChatEventUsernameChanged(t *testing.T) {
	timestamp := time.Unix(72, 6).UTC()
	user := user.User{
		ID:              "user id",
		DisplayName:     "display name",
		DisplayColor:    4,
		CreatedAt:       time.Unix(3, 26).UTC(),
		DisabledAt:      nil,
		PreviousNames:   []string{"somebody"},
		NameChangedAt:   nil,
		Scopes:          []string{},
		IsBot:           false,
		AuthenticatedAt: nil,
		Authenticated:   false,
	}

	checkPayload(t, models.UserNameChanged, func() {
		SendChatEventUsernameChanged(events.NameChangeEvent{
			Event: events.Event{
				Type:      events.UserNameChanged,
				ID:        "id",
				Timestamp: timestamp,
			},
			UserEvent: events.UserEvent{
				User:     &user,
				ClientID: 51,
				HiddenAt: nil,
			},
			NewName: "new name",
		})
	}, `{
		"clientId": 51,
		"id": "id",
		"newName": "new name",
		"timestamp": "1970-01-01T00:01:12.000000006Z",
		"type": "NAME_CHANGE",
		"user": {
			"authenticated": false,
			"createdAt": "1970-01-01T00:00:03.000000026Z",
			"displayColor": 4,
			"displayName": "display name",
			"id": "user id",
			"isBot": false,
			"previousNames": ["somebody"]
		}
	}`)
}

func TestSendChatEventUserJoined(t *testing.T) {
	timestamp := time.Unix(72, 6).UTC()
	user := user.User{
		ID:              "user id",
		DisplayName:     "display name",
		DisplayColor:    4,
		CreatedAt:       time.Unix(3, 26).UTC(),
		DisabledAt:      nil,
		PreviousNames:   []string{"somebody"},
		NameChangedAt:   nil,
		Scopes:          []string{},
		IsBot:           false,
		AuthenticatedAt: nil,
		Authenticated:   false,
	}

	checkPayload(t, models.UserJoined, func() {
		SendChatEventUserJoined(events.UserJoinedEvent{
			Event: events.Event{
				Type:      events.UserJoined,
				ID:        "id",
				Timestamp: timestamp,
			},
			UserEvent: events.UserEvent{
				User:     &user,
				ClientID: 51,
				HiddenAt: nil,
			},
		})
	}, `{
		"clientId": 51,
		"id": "id",
		"type": "USER_JOINED",
		"timestamp": "1970-01-01T00:01:12.000000006Z",
		"user": {
			"authenticated": false,
			"createdAt": "1970-01-01T00:00:03.000000026Z",
			"displayColor": 4,
			"displayName": "display name",
			"id": "user id",
			"isBot": false,
			"previousNames": ["somebody"]
		}
	}`)
}

func TestSendChatEventSetMessageVisibility(t *testing.T) {
	timestamp := time.Unix(72, 6).UTC()

	checkPayload(t, models.VisibiltyToggled, func() {
		SendChatEventSetMessageVisibility(events.SetMessageVisibilityEvent{
			Event: events.Event{
				Type:      events.VisibiltyUpdate,
				ID:        "id",
				Timestamp: timestamp,
			},
			UserMessageEvent: events.UserMessageEvent{},
			MessageIDs:       []string{"message1", "message2"},
			Visible:          false,
		})
	}, `{
		"MessageIDs": [
			"message1",
			"message2"
		],
		"Visible": false,
		"body": "",
		"id": "id",
		"timestamp": "1970-01-01T00:01:12.000000006Z",
		"type": "VISIBILITY-UPDATE",
		"user": null
	}`)
}
