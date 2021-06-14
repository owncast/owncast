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
		panic(err)
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
		panic(err)
	}

}

func (s *ChatServer) userMessageSent(eventData chatClientEvent) {
	// fmt.Println("server:userMessageSent", event)
	var event events.UserMessageEvent
	if err := json.Unmarshal(eventData.data, &event); err != nil {
		panic(err)
	}
	event.SetDefaults()
	event.User = user.GetUserByToken(eventData.client.accessToken)
	payload := event.GetBroadcastPayload()
	if err := s.Broadcast(payload); err != nil {
		panic(err)
	}

	addMessage(event)
}
