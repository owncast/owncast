package chat

import (
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

	// Send an update event to all clients for each message.
	// Note: Our client expects a single message at a time, so we can't just
	// send an array of messages in a single update.
	for _, id := range messageIDs {
		message, err := getMessageByID(id)
		if err != nil {
			log.Errorln(err)
			continue
		}
		payload := message.GetBroadcastPayload()
		payload["type"] = events.VisibiltyToggled
		if err := _server.Broadcast(payload); err != nil {
			log.Debugln(err)
		}

		go webhooks.SendChatEvent(message)
	}

	return nil
}
