package chat

import (
	"github.com/owncast/owncast/core/webhooks"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

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
		message, err := getMessageById(id)
		if err != nil {
			log.Errorln(err)
			continue
		}
		message.MessageType = models.VisibiltyToggled
		_server.sendAll(message)

		go webhooks.SendChatEvent(message)
	}

	return nil
}
