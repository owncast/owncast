package chat

import (
	"errors"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/webhooks"
	"github.com/owncast/owncast/storage/chatrepository"
	log "github.com/sirupsen/logrus"
)

// SetMessagesVisibility will set the visibility of multiple messages by ID.
func (c *Chat) SetMessagesVisibility(messageIDs []string, visibility bool) error {
	cr := chatrepository.Get()

	// Save new message visibility
	if err := cr.SaveMessageVisibility(messageIDs, visibility); err != nil {
		log.Errorln(err)
		return err
	}

	// Send an event letting the chat clients know to hide or show
	// the messages.
	event := models.SetMessageVisibilityEvent{
		MessageIDs: messageIDs,
		Visible:    visibility,
	}
	event.Event.SetDefaults()

	payload := event.GetBroadcastPayload()
	if err := c.server.Broadcast(payload); err != nil {
		return errors.New("error broadcasting message visibility payload " + err.Error())
	}

	webhookManager := webhooks.Get()
	webhookManager.SendChatEventSetMessageVisibility(event)

	return nil
}
