package chat

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/core/webhooks"
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

	for _, blockedName := range blocklist {
		normalizedName := strings.TrimSpace(blockedName)
		normalizedName = strings.ToLower(normalizedName)
		if strings.Contains(normalizedName, proposedUsername) {
			// Denied.
			log.Debugln(eventData.client.User.DisplayName, "blocked from changing name to", proposedUsername, "due to blocked name", normalizedName)
			message := fmt.Sprintf("You cannot change your name to **%s**.", proposedUsername)
			s.sendActionToClient(eventData.client, message)
			return
		}
	}

	savedUser := user.GetUserByToken(eventData.client.accessToken)
	oldName := savedUser.DisplayName

	// Save the new name
	user.ChangeUsername(eventData.client.User.Id, receivedEvent.NewName)

	// Send chat event letting everyone about about the name change
	savedUser.DisplayName = receivedEvent.NewName

	// Update the connected clients associated user with the new name
	eventData.client.User = savedUser

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

	// Send chat user name changed webhook
	receivedEvent.User = savedUser
	webhooks.SendChatEventUsernameChanged(receivedEvent)
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

	// Send chat message sent webhook
	webhooks.SendChatEvent(&event)

	SaveUserMessage(event)

	eventData.client.MessageCount = eventData.client.MessageCount + 1
}
