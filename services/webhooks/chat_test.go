package webhooks

import (
	"testing"
	"time"

	"github.com/owncast/owncast/models"
)

func TestSendChatEvent(t *testing.T) {
	timestamp := time.Unix(72, 6).UTC()
	user := models.User{
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
		manager.SendChatEvent(&models.UserMessageEvent{
			Event: models.Event{
				Type:      models.MessageSent,
				ID:        "id",
				Timestamp: timestamp,
			},
			UserEvent: models.UserEvent{
				User:     &user,
				ClientID: 51,
				HiddenAt: nil,
			},
			MessageEvent: models.MessageEvent{
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
	user := models.User{
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
		manager.SendChatEventUsernameChanged(models.NameChangeEvent{
			Event: models.Event{
				Type:      models.UserNameChanged,
				ID:        "id",
				Timestamp: timestamp,
			},
			UserEvent: models.UserEvent{
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
	user := models.User{
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
		manager.SendChatEventUserJoined(models.UserJoinedEvent{
			Event: models.Event{
				Type:      models.UserJoined,
				ID:        "id",
				Timestamp: timestamp,
			},
			UserEvent: models.UserEvent{
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
		manager.SendChatEventSetMessageVisibility(models.SetMessageVisibilityEvent{
			Event: models.Event{
				Type:      models.VisibiltyUpdate,
				ID:        "id",
				Timestamp: timestamp,
			},
			UserMessageEvent: models.UserMessageEvent{},
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
