package chat

import (
	"encoding/json"

	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

func (s *ChatServer) userNameChanged(eventData chatClientEvent) {
	var receivedEvent events.NameChangeEvent
	if err := json.Unmarshal(eventData.data, &receivedEvent); err != nil {
		log.Errorln("error unmarshalling to NameChangeEvent", err)
		return
	}

	proposedUsername := receivedEvent.NewName
	blocklist := data.GetUsernameBlocklist()

	if _, blocked := utils.FindInSlice(blocklist, proposedUsername); blocked {
		// Denied.
		log.Debugln(receivedEvent.User.DisplayName, "blocked from changing name to", proposedUsername)
		return
	}

	// Save the new name
	savedUser := user.GetUserByToken(eventData.client.accessToken)
	user.ChangeUsername(savedUser.Id, receivedEvent.NewName)

	// Send chat event letting everyone about about the name change
	oldName := savedUser.DisplayName
	savedUser.DisplayName = receivedEvent.NewName
	broadcastEvent := events.NameChangeBroadcast{
		Oldname: oldName,
	}
	broadcastEvent.User = savedUser
	broadcastEvent.SetDefaults()
	payload := broadcastEvent.GetBroadcastPayload()
	if err := s.Broadcast(payload); err != nil {
		log.Errorln("error broadcasting NameChangeEvent", err)
		return
	}

}

func (s *ChatServer) userMessageSent(eventData chatClientEvent) {
	var event events.UserMessageEvent
	if err := json.Unmarshal(eventData.data, &event); err != nil {
		log.Errorln("error unmarshalling to UserMessageEvent", err)
		return
	}

	event.SetDefaults()

	// Ignore empty messages
	if event.Empty() {
		return
	}

	event.User = user.GetUserByToken(eventData.client.accessToken)

	// Guard against nil users
	if event.User == nil {
		return
	}

	payload := event.GetBroadcastPayload()
	if err := s.Broadcast(payload); err != nil {
		log.Errorln("error broadcasting UserMessageEvent payload", err)
		return
	}

	addMessage(event)
}
