package admin

// this is endpoint logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

// ExternalUpdateMessageVisibility updates an array of message IDs to have the same visiblity.
func ExternalUpdateMessageVisibility(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	UpdateMessageVisibility(w, r)
}

// UpdateMessageVisibility updates an array of message IDs to have the same visiblity.
func UpdateMessageVisibility(w http.ResponseWriter, r *http.Request) {
	type messageVisibilityUpdateRequest struct {
		IDArray []string `json:"idArray"`
		Visible bool     `json:"visible"`
	}

	if r.Method != controllers.POST {
		controllers.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request messageVisibilityUpdateRequest

	if err := decoder.Decode(&request); err != nil {
		log.Errorln(err)
		controllers.WriteSimpleResponse(w, false, "")
		return
	}

	if err := chat.SetMessagesVisibility(request.IDArray, request.Visible); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "changed")
}

// UpdateUserEnabled enable or disable a single user by ID.
func UpdateUserEnabled(w http.ResponseWriter, r *http.Request) {
	type blockUserRequest struct {
		UserID  string `json:"userId"`
		Enabled bool   `json:"enabled"`
	}

	if r.Method != controllers.POST {
		controllers.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request blockUserRequest

	if err := decoder.Decode(&request); err != nil {
		log.Errorln(err)
		controllers.WriteSimpleResponse(w, false, "")
		return
	}

	// Disable/enable the user
	if err := user.SetEnabled(request.UserID, request.Enabled); err != nil {
		log.Errorln("error changing user enabled status", err)
	}

	// Hide/show the user's chat messages if disabling.
	// Leave hidden messages hidden to be safe.
	if !request.Enabled {
		if err := chat.SetMessageVisibilityForUserID(request.UserID, request.Enabled); err != nil {
			log.Errorln("error changing user messages visibility", err)
		}
	}

	// Forcefully disconnect the user from the chat
	if !request.Enabled {
		chat.DisconnectUser(request.UserID)
		disconnectedUser := user.GetUserByID(request.UserID)
		_ = chat.SendSystemAction(fmt.Sprintf("**%s** has been removed from chat.", disconnectedUser.DisplayName), true)
	}

	controllers.WriteSimpleResponse(w, true, fmt.Sprintf("%s enabled: %t", request.UserID, request.Enabled))
}

// GetDisabledUsers will return all the disabled users.
func GetDisabledUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := user.GetDisabledUsers()
	controllers.WriteResponse(w, users)
}

// GetChatMessages returns all of the chat messages, unfiltered.
func GetChatMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	messages := chat.GetChatModerationHistory()
	controllers.WriteResponse(w, messages)
}

// SendSystemMessage will send an official "SYSTEM" message to chat on behalf of your server.
func SendSystemMessage(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var message events.SystemMessageEvent
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		controllers.InternalErrorHandler(w, err)
		return
	}

	if err := chat.SendSystemMessage(message.Body, false); err != nil {
		controllers.BadRequestHandler(w, err)
	}

	controllers.WriteSimpleResponse(w, true, "sent")
}

func SendSystemMessageToConnectedClient(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	clientIdText, err := utils.ReadRestUrlParameter(r, "clientId")
	if err != nil {
		controllers.BadRequestHandler(w, err)
		return
	}

	clientIdNumeric, err := strconv.ParseUint(clientIdText, 10, 32)
	if err != nil {
		controllers.BadRequestHandler(w, err)
		return
	}

	var message events.SystemMessageEvent
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		controllers.InternalErrorHandler(w, err)
		return
	}

	chat.SendSystemMessageToClient(uint(clientIdNumeric), message.Body)
	controllers.WriteSimpleResponse(w, true, "sent")
}

// SendUserMessage will send a message to chat on behalf of a user. *Depreciated*.
func SendUserMessage(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	controllers.BadRequestHandler(w, errors.New("no longer supported. see /api/integrations/chat/send"))
}

// SendIntegrationChatMessage will send a chat message on behalf of an external chat integration.
func SendIntegrationChatMessage(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := integration.DisplayName

	if name == "" {
		controllers.BadRequestHandler(w, errors.New("unknown integration for provided access token"))
		return
	}

	var event events.UserMessageEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		controllers.InternalErrorHandler(w, err)
		return
	}
	event.SetDefaults()
	event.RenderBody()
	event.Type = "CHAT"

	if event.Empty() {
		controllers.BadRequestHandler(w, errors.New("invalid message"))
		return
	}

	event.User = &user.User{
		ID:           integration.ID,
		DisplayName:  name,
		DisplayColor: integration.DisplayColor,
		CreatedAt:    integration.CreatedAt,
	}

	if err := chat.Broadcast(&event); err != nil {
		controllers.BadRequestHandler(w, err)
		return
	}

	chat.SaveUserMessage(event)

	controllers.WriteSimpleResponse(w, true, "sent")
}

// SendChatAction will send a generic chat action.
func SendChatAction(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var message events.SystemActionEvent
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		controllers.InternalErrorHandler(w, err)
		return
	}

	message.SetDefaults()
	message.RenderBody()

	if err := chat.SendSystemAction(message.Body, false); err != nil {
		controllers.BadRequestHandler(w, err)
		return
	}

	controllers.WriteSimpleResponse(w, true, "sent")
}
