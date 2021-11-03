package chat

import (
	"errors"

	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/webhooks"
	log "github.com/sirupsen/logrus"
)

// SetMessagesVisibility will set the visibility of multiple messages by ID.
func SetMessagesVisibility(messageIDs []string, visibility bool) error {
	// Save new message visibility
	if err := saveMessageVisibility(messageIDs, visibility); err != nil {
		log.Errorln(err)
		return err
	}

	// Send an event letting the chat clients know to hide or show
	// the messages.
	event := events.SetMessageVisibilityEvent{
		MessageIDs: messageIDs,
		Visible:    visibility,
	}
	event.Event.SetDefaults()

	payload := event.GetBroadcastPayload()
	if err := _server.Broadcast(payload); err != nil {
		return errors.New("error broadcasting message visibility payload " + err.Error())
	}

	// Send webhook
	wh := webhooks.WebhookEvent{
		EventData: event,
		Type:      event.GetMessageType(),
	}

	webhooks.SendEventToWebhooks(wh)

	return nil
}
