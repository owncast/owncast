package webhooks

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/webhookrepository"
	"github.com/owncast/owncast/services/chat/events"
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
		testSvc.SendChatEvent(&events.UserMessageEvent{
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
		"serverURL": "http://localhost:8080",
		"status": {
			"lastConnectTime": null,
			"lastDisconnectTime": null,
			"online": true,
			"overallMaxViewerCount": 420,
			"sessionMaxViewerCount": 69,
			"streamTitle": "my stream",
			"versionNumber": "1.2.3",
			"viewerCount": 5
		},
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
		testSvc.SendChatEventUsernameChanged(events.NameChangeEvent{
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
		"id": "id",
		"newName": "new name",
		"serverURL": "http://localhost:8080",
		"status": {
			"lastConnectTime": null,
			"lastDisconnectTime": null,
			"online": true,
			"overallMaxViewerCount": 420,
			"sessionMaxViewerCount": 69,
			"streamTitle": "my stream",
			"versionNumber": "1.2.3",
			"viewerCount": 5
		},
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
		testSvc.SendChatEventUserJoined(events.UserJoinedEvent{
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
		"id": "id",
		"serverURL": "http://localhost:8080",
		"status": {
			"lastConnectTime": null,
			"lastDisconnectTime": null,
			"online": true,
			"overallMaxViewerCount": 420,
			"sessionMaxViewerCount": 69,
			"streamTitle": "my stream",
			"versionNumber": "1.2.3",
			"viewerCount": 5
		},
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
		testSvc.SendChatEventSetMessageVisibility(events.SetMessageVisibilityEvent{
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
		"id": "id",
		"ids": [
			"message1",
			"message2"
		],
		"serverURL": "http://localhost:8080",
		"status": {
			"lastConnectTime": null,
			"lastDisconnectTime": null,
			"online": true,
			"overallMaxViewerCount": 420,
			"sessionMaxViewerCount": 69,
			"streamTitle": "my stream",
			"versionNumber": "1.2.3",
			"viewerCount": 5
		},
		"timestamp": "1970-01-01T00:01:12.000000006Z",
		"user": null,
		"visible": false
	}`)
}

// TestWebhookHasServerStatus verifies that all webhook events include server status
func TestWebhookHasServerStatus(t *testing.T) {
	// Set up server configuration
	configRepository := configrepository.New(testDatastore)
	configRepository.SetServerURL("http://localhost:8080")

	eventChannel := make(chan WebhookEvent)

	// Set up a server.
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := WebhookEvent{}
		json.NewDecoder(r.Body).Decode(&data)
		eventChannel <- data
	}))
	defer svr.Close()

	webhooksRepo := webhookrepository.New(testDatastore)

	// Subscribe to the webhook.
	hook, err := webhooksRepo.InsertWebhook(svr.URL, []models.EventType{models.UserJoined})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := webhooksRepo.DeleteWebhook(hook); err != nil {
			t.Error(err)
		}
	}()

	// Send a chat event
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

	testSvc.SendChatEventUserJoined(events.UserJoinedEvent{
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

	// Capture the event
	event := <-eventChannel

	// Verify the webhook event has a status field in eventData
	eventData, ok := event.EventData.(map[string]interface{})
	if !ok {
		t.Error("Expected EventData to be a map")
	}

	status, ok := eventData["status"].(map[string]interface{})
	if !ok {
		t.Error("Expected eventData to contain status field")
	}

	versionNumber, ok := status["versionNumber"].(string)
	if !ok || versionNumber == "" {
		t.Error("Expected eventData.status to have versionNumber, but it was empty")
	}

	serverURL, ok := eventData["serverURL"].(string)
	if !ok || serverURL == "" {
		t.Error("Expected eventData to have serverURL, but it was empty")
	}

	if event.Type != models.UserJoined {
		t.Errorf("Expected event type %v but got %v", models.UserJoined, event.Type)
	}
}
