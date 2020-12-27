package chat

import "github.com/owncast/owncast/models"

func SetMessagesVisability(messageIDs []string, visability bool) error {
	// Send an update event to all clients for each message.
	// Note: Our client expects a single message at a time, so we can't just
	// send an array of messages in a single update.
	for _, id := range messageIDs {
		event := models.ChatEvent{
			ID:          id,
			Visible:     visability,
			MessageType: UPDATE,
		}
		_server.sendAll(event)
	}

	// Save new message visability to chat database
	return saveMessageVisibility(messageIDs, visability)
}
