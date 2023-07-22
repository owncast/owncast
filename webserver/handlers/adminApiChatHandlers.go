package handlers

// this is endpoint logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/chat"
	"github.com/owncast/owncast/storage/chatrepository"
	"github.com/owncast/owncast/storage/userrepository"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/requests"
	"github.com/owncast/owncast/webserver/responses"

	log "github.com/sirupsen/logrus"
)

var (
	chatRepository = chatrepository.Get()
	userRepository = userrepository.Get()
)

// ExternalUpdateMessageVisibility updates an array of message IDs to have the same visiblity.
func (h *Handlers) ExternalUpdateMessageVisibility(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	h.UpdateMessageVisibility(w, r)
}

// UpdateMessageVisibility updates an array of message IDs to have the same visiblity.
func (h *Handlers) UpdateMessageVisibility(w http.ResponseWriter, r *http.Request) {
	type messageVisibilityUpdateRequest struct {
		IDArray []string `json:"idArray"`
		Visible bool     `json:"visible"`
	}

	if r.Method != http.MethodPost {
		responses.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request messageVisibilityUpdateRequest

	if err := decoder.Decode(&request); err != nil {
		log.Errorln(err)
		responses.WriteSimpleResponse(w, false, "")
		return
	}

	if err := h.chatService.SetMessagesVisibility(request.IDArray, request.Visible); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "changed")
}

// BanIPAddress will manually ban an IP address.
func (h *Handlers) BanIPAddress(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		responses.WriteSimpleResponse(w, false, "unable to ban IP address")
		return
	}

	if err := chatRepository.BanIPAddress(configValue.Value.(string), "manually added"); err != nil {
		responses.WriteSimpleResponse(w, false, "error saving IP address ban")
		return
	}

	responses.WriteSimpleResponse(w, true, "IP address banned")
}

// UnBanIPAddress will remove an IP address ban.
func (h *Handlers) UnBanIPAddress(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		responses.WriteSimpleResponse(w, false, "unable to unban IP address")
		return
	}

	if err := chatRepository.RemoveIPAddressBan(configValue.Value.(string)); err != nil {
		responses.WriteSimpleResponse(w, false, "error removing IP address ban")
		return
	}

	responses.WriteSimpleResponse(w, true, "IP address unbanned")
}

// GetIPAddressBans will return all the banned IP addresses.
func (h *Handlers) GetIPAddressBans(w http.ResponseWriter, r *http.Request) {
	bans, err := chatRepository.GetIPAddressBans()
	if err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteResponse(w, bans)
}

// UpdateUserEnabled enable or disable a single user by ID.
func (h *Handlers) UpdateUserEnabled(w http.ResponseWriter, r *http.Request) {
	type blockUserRequest struct {
		UserID  string `json:"userId"`
		Enabled bool   `json:"enabled"`
	}

	if r.Method != http.MethodPost {
		responses.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request blockUserRequest

	if err := decoder.Decode(&request); err != nil {
		log.Errorln(err)
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	if request.UserID == "" {
		responses.WriteSimpleResponse(w, false, "must provide userId")
		return
	}

	// Disable/enable the user
	if err := userRepository.SetEnabled(request.UserID, request.Enabled); err != nil {
		log.Errorln("error changing user enabled status", err)
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Hide/show the user's chat messages if disabling.
	// Leave hidden messages hidden to be safe.
	if !request.Enabled {
		messageIds, err := h.chatRepository.SetMessageVisibilityForUserID(request.UserID, request.Enabled)
		if err != nil {
			log.Errorln("error changing user messages visibility", err)
			responses.WriteSimpleResponse(w, false, err.Error())
			return
		}
		h.chatService.SetMessagesVisibility(messageIds, request.Enabled)
	}

	// Forcefully disconnect the user from the chat
	if !request.Enabled {
		clients, err := h.chatService.GetClientsForUser(request.UserID)
		if len(clients) == 0 {
			// Nothing to do
			return
		}

		if err != nil {
			log.Errorln("error fetching clients for user: ", err)
			responses.WriteSimpleResponse(w, false, err.Error())
			return
		}

		h.chatService.DisconnectClients(clients)
		disconnectedUser := userRepository.GetUserByID(request.UserID)
		_ = h.chatService.SendSystemAction(fmt.Sprintf("**%s** has been removed from chat.", disconnectedUser.DisplayName), true)

		localIP4Address := "127.0.0.1"
		localIP6Address := "::1"

		// Ban this user's IP address.
		for _, client := range clients {
			ipAddress := client.IPAddress
			if ipAddress != localIP4Address && ipAddress != localIP6Address {
				reason := fmt.Sprintf("Banning of %s", disconnectedUser.DisplayName)
				if err := h.chatRepository.BanIPAddress(ipAddress, reason); err != nil {
					log.Errorln("error banning IP address: ", err)
				}
			}
		}
	}

	responses.WriteSimpleResponse(w, true, fmt.Sprintf("%s enabled: %t", request.UserID, request.Enabled))
}

// GetDisabledUsers will return all the disabled users.
func (h *Handlers) GetDisabledUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := userRepository.GetDisabledUsers()
	responses.WriteResponse(w, users)
}

// UpdateUserModerator will set the moderator status for a user ID.
func (h *Handlers) UpdateUserModerator(w http.ResponseWriter, r *http.Request) {
	type request struct {
		UserID      string `json:"userId"`
		IsModerator bool   `json:"isModerator"`
	}

	if r.Method != http.MethodPost {
		responses.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var req request

	if err := decoder.Decode(&req); err != nil {
		responses.WriteSimpleResponse(w, false, "")
		return
	}

	// Update the user object with new moderation access.
	if err := userRepository.SetModerator(req.UserID, req.IsModerator); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update the clients for this user to know about the moderator access change.
	if err := h.chatService.SendConnectedClientInfoToUser(req.UserID); err != nil {
		log.Debugln(err)
	}

	responses.WriteSimpleResponse(w, true, fmt.Sprintf("%s is moderator: %t", req.UserID, req.IsModerator))
}

// GetModerators will return a list of moderator users.
func (h *Handlers) GetModerators(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := userRepository.GetModeratorUsers()
	responses.WriteResponse(w, users)
}

// GetChatMessages returns all of the chat messages, unfiltered.
func (h *Handlers) GetAdminChatMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	messages := h.chatRepository.GetChatModerationHistory()
	responses.WriteResponse(w, messages)
}

// SendSystemMessage will send an official "SYSTEM" message to chat on behalf of your server.
func (h *Handlers) SendSystemMessage(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var message chat.SystemMessageEvent
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		responses.InternalErrorHandler(w, err)
		return
	}

	if err := h.chatService.SendSystemMessage(message.Body, false); err != nil {
		responses.BadRequestHandler(w, err)
	}

	responses.WriteSimpleResponse(w, true, "sent")
}

// SendSystemMessageToConnectedClient will handle incoming requests to send a single message to a single connected client by ID.
func (h *Handlers) SendSystemMessageToConnectedClient(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	clientIDText, err := utils.ReadRestURLParameter(r, "clientId")
	if err != nil {
		responses.BadRequestHandler(w, err)
		return
	}

	clientIDNumeric, err := strconv.ParseUint(clientIDText, 10, 32)
	if err != nil {
		responses.BadRequestHandler(w, err)
		return
	}

	var message chat.SystemMessageEvent
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		responses.InternalErrorHandler(w, err)
		return
	}

	h.chatService.SendSystemMessageToClient(uint(clientIDNumeric), message.Body)
	responses.WriteSimpleResponse(w, true, "sent")
}

// SendUserMessage will send a message to chat on behalf of a user. *Depreciated*.
func (h *Handlers) SendUserMessage(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	responses.BadRequestHandler(w, errors.New("no longer supported. see /api/integrations/chat/send"))
}

// SendIntegrationChatMessage will send a chat message on behalf of an external chat integration.
func (h *Handlers) SendIntegrationChatMessage(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := integration.DisplayName

	if name == "" {
		responses.BadRequestHandler(w, errors.New("unknown integration for provided access token"))
		return
	}

	var event chat.UserMessageEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		responses.InternalErrorHandler(w, err)
		return
	}
	event.SetDefaults()
	event.RenderBody()
	event.Type = "CHAT"

	if event.Empty() {
		responses.BadRequestHandler(w, errors.New("invalid message"))
		return
	}

	event.User = &models.User{
		ID:           integration.ID,
		DisplayName:  name,
		DisplayColor: integration.DisplayColor,
		CreatedAt:    integration.CreatedAt,
		IsBot:        true,
	}

	if err := h.chatService.Broadcast(&event); err != nil {
		responses.BadRequestHandler(w, err)
		return
	}

	h.chatRepository.SaveUserMessage(event)

	responses.WriteSimpleResponse(w, true, "sent")
}

// SendChatAction will send a generic chat action.
func (h *Handlers) SendChatAction(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var message chat.SystemActionEvent
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		responses.InternalErrorHandler(w, err)
		return
	}

	message.SetDefaults()
	message.RenderBody()

	if err := h.chatService.SendSystemAction(message.Body, false); err != nil {
		responses.BadRequestHandler(w, err)
		return
	}

	responses.WriteSimpleResponse(w, true, "sent")
}

// SetEnableEstablishedChatUserMode sets the requirement for a chat user
// to be "established" for some time before taking part in chat.
func (h *Handlers) SetEnableEstablishedChatUserMode(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		responses.WriteSimpleResponse(w, false, "unable to update chat established user only mode")
		return
	}

	if err := configRepository.SetChatEstablishedUsersOnlyMode(configValue.Value.(bool)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "chat established users only mode updated")
}
