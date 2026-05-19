package chat

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/services/chat/events"
)

// SetMessagesVisibility sets the visibility of multiple messages by ID.
func (s *Service) SetMessagesVisibility(messageIDs []string, visibility bool) error {
	// Save new message visibility
	if err := s.chatMessageRepository.SetMessageVisibilityForMessageIDs(messageIDs, visibility); err != nil {
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
	if err := s.Broadcast(payload); err != nil {
		return errors.New("error broadcasting message visibility payload " + err.Error())
	}

	s.webhooks.SendChatEventSetMessageVisibility(event)

	return nil
}
