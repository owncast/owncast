package chat

import (
	"context"
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/chatmessagerepository"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/webhookrepository"
	"github.com/owncast/owncast/services/chat/events"
	"github.com/owncast/owncast/services/dispatcher"
	"github.com/owncast/owncast/services/webhooks"
)

type moderationTestChatRepository struct{}

func (moderationTestChatRepository) SaveUserMessage(events.UserMessageEvent)             {}
func (moderationTestChatRepository) SaveFederatedAction(events.FediverseEngagementEvent) {}
func (moderationTestChatRepository) SaveEvent(string, *string, string, string, *time.Time, time.Time, *string, *string, *string, *string) {
}
func (moderationTestChatRepository) GetChatModerationHistory() []interface{} { return nil }
func (moderationTestChatRepository) GetChatHistory() []interface{}           { return nil }
func (moderationTestChatRepository) GetMessagesFromUser(string) ([]events.UserMessageEvent, error) {
	return nil, nil
}

func (moderationTestChatRepository) GetMessageIdsForUserID(string) ([]string, error) {
	return nil, nil
}

func (moderationTestChatRepository) SetMessageVisibilityForMessageIDs([]string, bool) error {
	return nil
}
func (moderationTestChatRepository) GetMessagesCount() int64 { return 0 }

var _ chatmessagerepository.ChatMessageRepository = moderationTestChatRepository{}

type moderationTestConfig struct {
	configrepository.ConfigRepository
}

func (moderationTestConfig) GetServerURL() string { return "http://owncast.test" }

type moderationTestWebhookRepository struct {
	webhookrepository.WebhookRepository
}

func (moderationTestWebhookRepository) GetWebhooksForEvent(models.EventType) []models.Webhook {
	return nil
}

func TestSetMessagesVisibilityIncludesModeratorInWebhookEvent(t *testing.T) {
	eventDispatcher := dispatcher.New()
	var received webhooks.WebhookEvent
	eventDispatcher.AddListener(func(_ context.Context, event dispatcher.Event) {
		received = event.Payload.(webhooks.WebhookEvent)
	})

	webhookService := webhooks.New(webhooks.Deps{
		GetStatus:         func() models.Status { return models.Status{Online: true} },
		ConfigRepository:  moderationTestConfig{},
		WebhookRepository: moderationTestWebhookRepository{},
		Events:            eventDispatcher,
	})
	chatService := New(Deps{
		ChatMessageRepository: moderationTestChatRepository{},
		Webhooks:              webhookService,
	})
	moderator := &models.User{
		ID:          "moderator-id",
		DisplayName: "Moderator",
	}

	if err := chatService.SetMessagesVisibility([]string{"message-id"}, false, moderator); err != nil {
		t.Fatal(err)
	}

	payload, ok := received.EventData.(*webhooks.WebhookVisibilityToggleEventData)
	if !ok {
		t.Fatalf("event data type = %T, want *WebhookVisibilityToggleEventData", received.EventData)
	}
	if payload.User == nil || payload.User.ID != moderator.ID {
		t.Fatalf("moderator = %#v, want ID %q", payload.User, moderator.ID)
	}
}
