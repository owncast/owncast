package admin

// this is endpoint logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/utils"
)

// ExternalUpdateMessageVisibility updates an array of message IDs to have the same visiblity.
func (c *Controller) ExternalUpdateMessageVisibility(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	c.UpdateMessageVisibility(w, r)
}

// UpdateMessageVisibility updates an array of message IDs to have the same visiblity.
func (c *Controller) UpdateMessageVisibility(w http.ResponseWriter, r *http.Request) {
	type messageVisibilityUpdateRequest struct {
		IDArray []string `json:"idArray"`
		Visible bool     `json:"visible"`
	}

	if r.Method != controllers.POST {
		c.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request messageVisibilityUpdateRequest

	if err := decoder.Decode(&request); err != nil {
		log.Errorln(err)
		c.WriteSimpleResponse(w, false, "")
		return
	}

	if err := c.Core.Chat.SetMessagesVisibility(request.IDArray, request.Visible); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "changed")
}

// BanIPAddress will manually ban an IP address.
func (c *Controller) BanIPAddress(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		c.WriteSimpleResponse(w, false, "unable to ban IP address")
		return
	}

	if err := c.Data.BanIPAddress(configValue.Value.(string), "manually added"); err != nil {
		c.WriteSimpleResponse(w, false, "error saving IP address ban")
		return
	}

	c.WriteSimpleResponse(w, true, "IP address banned")
}

// UnBanIPAddress will remove an IP address ban.
func (c *Controller) UnBanIPAddress(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		c.WriteSimpleResponse(w, false, "unable to unban IP address")
		return
	}

	if err := c.Data.RemoveIPAddressBan(configValue.Value.(string)); err != nil {
		c.WriteSimpleResponse(w, false, "error removing IP address ban")
		return
	}

	c.WriteSimpleResponse(w, true, "IP address unbanned")
}

// GetIPAddressBans will return all the banned IP addresses.
func (c *Controller) GetIPAddressBans(w http.ResponseWriter, r *http.Request) {
	bans, err := c.Data.GetIPAddressBans()
	if err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteResponse(w, bans)
}

// UpdateUserEnabled enable or disable a single user by ID.
func (c *Controller) UpdateUserEnabled(w http.ResponseWriter, r *http.Request) {
	type blockUserRequest struct {
		UserID  string `json:"userId"`
		Enabled bool   `json:"enabled"`
	}

	if r.Method != controllers.POST {
		c.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request blockUserRequest

	if err := decoder.Decode(&request); err != nil {
		log.Errorln(err)
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	if request.UserID == "" {
		c.WriteSimpleResponse(w, false, "must provide userId")
		return
	}

	// Disable/enable the user
	if err := user.SetEnabled(request.UserID, request.Enabled, c.Core.Data.Store); err != nil {
		log.Errorln("error changing user enabled status", err)
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Hide/show the user's chat messages if disabling.
	// Leave hidden messages hidden to be safe.
	if !request.Enabled {
		if err := c.Core.Chat.SetMessageVisibilityForUserID(request.UserID, request.Enabled); err != nil {
			log.Errorln("error changing user messages visibility", err)
			c.WriteSimpleResponse(w, false, err.Error())
			return
		}
	}

	// Forcefully disconnect the user from the chat
	if !request.Enabled {
		clients, err := c.Core.Chat.GetClientsForUser(request.UserID)
		if len(clients) == 0 {
			// Nothing to do
			return
		}

		if err != nil {
			log.Errorln("error fetching clients for user: ", err)
			c.WriteSimpleResponse(w, false, err.Error())
			return
		}

		c.Core.Chat.DisconnectClients(clients)
		disconnectedUser := user.GetUserByID(request.UserID, c.Core.Data.Store)
		_ = c.Core.Chat.SendSystemAction(fmt.Sprintf("**%s** has been removed from chat.", disconnectedUser.DisplayName), true)

		// Ban this user's IP address.
		for _, client := range clients {
			ipAddress := client.IPAddress
			reason := fmt.Sprintf("Banning of %s", disconnectedUser.DisplayName)
			if err := c.Data.BanIPAddress(ipAddress, reason); err != nil {
				log.Errorln("error banning IP address: ", err)
			}
		}
	}

	c.WriteSimpleResponse(w, true, fmt.Sprintf("%s enabled: %t", request.UserID, request.Enabled))
}

// GetDisabledUsers will return all the disabled users.
func (c *Controller) GetDisabledUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := user.GetDisabledUsers(c.Data.Store)
	c.WriteResponse(w, users)
}

// UpdateUserModerator will set the moderator status for a user ID.
func (c *Controller) UpdateUserModerator(w http.ResponseWriter, r *http.Request) {
	type request struct {
		UserID      string `json:"userId"`
		IsModerator bool   `json:"isModerator"`
	}

	if r.Method != controllers.POST {
		c.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var req request

	if err := decoder.Decode(&req); err != nil {
		c.WriteSimpleResponse(w, false, "")
		return
	}

	// Update the user object with new moderation access.
	if err := user.SetModerator(req.UserID, req.IsModerator, c.Core.Data.Store); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update the clients for this user to know about the moderator access change.
	if err := c.Core.Chat.SendConnectedClientInfoToUser(req.UserID); err != nil {
		log.Debugln(err)
	}

	c.WriteSimpleResponse(w, true, fmt.Sprintf("%s is moderator: %t", req.UserID, req.IsModerator))
}

// GetModerators will return a list of moderator users.
func (c *Controller) GetModerators(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := user.GetModeratorUsers(c.Data.Store)
	c.WriteResponse(w, users)
}

// GetChatMessages returns all of the chat messages, unfiltered.
func (c *Controller) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	messages := chat.GetChatModerationHistory()
	c.WriteResponse(w, messages)
}

// SendSystemMessage will send an official "SYSTEM" message to chat on behalf of your server.
func (c *Controller) SendSystemMessage(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var message events.SystemMessageEvent
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		c.InternalErrorHandler(w, err)
		return
	}

	if err := c.Core.Chat.SendSystemMessage(message.Body, false); err != nil {
		c.BadRequestHandler(w, err)
	}

	c.WriteSimpleResponse(w, true, "sent")
}

// SendSystemMessageToConnectedClient will handle incoming requests to send a single message to a single connected client by ID.
func (c *Controller) SendSystemMessageToConnectedClient(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	clientIDText, err := utils.ReadRestURLParameter(r, "clientId")
	if err != nil {
		c.BadRequestHandler(w, err)
		return
	}

	clientIDNumeric, err := strconv.ParseUint(clientIDText, 10, 32)
	if err != nil {
		c.BadRequestHandler(w, err)
		return
	}

	var message events.SystemMessageEvent
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		c.InternalErrorHandler(w, err)
		return
	}

	c.Core.Chat.SendSystemMessageToClient(uint(clientIDNumeric), message.Body)
	c.WriteSimpleResponse(w, true, "sent")
}

// SendUserMessage will send a message to chat on behalf of a user. *Depreciated*.
func (c *Controller) SendUserMessage(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	c.BadRequestHandler(w, errors.New("no longer supported. see /api/integrations/chat/send"))
}

// SendIntegrationChatMessage will send a chat message on behalf of an external chat integration.
func (c *Controller) SendIntegrationChatMessage(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := integration.DisplayName

	if name == "" {
		c.BadRequestHandler(w, errors.New("unknown integration for provided access token"))
		return
	}

	var event events.UserMessageEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		c.InternalErrorHandler(w, err)
		return
	}
	event.SetDefaults()
	event.RenderBody()
	event.Type = "CHAT"

	if event.Empty() {
		c.BadRequestHandler(w, errors.New("invalid message"))
		return
	}

	event.User = &user.User{
		ID:           integration.ID,
		DisplayName:  name,
		DisplayColor: integration.DisplayColor,
		CreatedAt:    integration.CreatedAt,
		IsBot:        true,
	}

	if err := c.Core.Chat.Broadcast(&event); err != nil {
		c.BadRequestHandler(w, err)
		return
	}

	chat.SaveUserMessage(event)

	c.WriteSimpleResponse(w, true, "sent")
}

// SendChatAction will send a generic chat action.
func (c *Controller) SendChatAction(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var message events.SystemActionEvent
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		c.InternalErrorHandler(w, err)
		return
	}

	message.SetDefaults()
	message.RenderBody()

	if err := c.Core.Chat.SendSystemAction(message.Body, false); err != nil {
		c.BadRequestHandler(w, err)
		return
	}

	c.WriteSimpleResponse(w, true, "sent")
}

// SetEnableEstablishedChatUserMode sets the requirement for a chat user
// to be "established" for some time before taking part in chat.
func (c *Controller) SetEnableEstablishedChatUserMode(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		c.WriteSimpleResponse(w, false, "unable to update chat established user only mode")
		return
	}

	if err := c.Data.SetChatEstablishedUsersOnlyMode(configValue.Value.(bool)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "chat established users only mode updated")
}
