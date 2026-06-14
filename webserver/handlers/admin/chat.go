package admin

// this is endpoint logic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/chat/events"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/handlers/generated"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// ExternalUpdateMessageVisibility updates an array of message IDs to have the same visiblity.
func (a *Admin) ExternalUpdateMessageVisibility(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	a.UpdateMessageVisibility(w, r)
}

// UpdateMessageVisibility updates an array of message IDs to have the same visiblity.
func (a *Admin) UpdateMessageVisibility(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// nolint:goconst
		webutils.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request generated.MessageVisibilityUpdate

	if err := decoder.Decode(&request); err != nil {
		log.Errorln(err)
		webutils.WriteSimpleResponse(w, false, "")
		return
	}

	if request.IdArray == nil || request.Visible == nil {
		webutils.WriteSimpleResponse(w, false, "missing required fields: idArray and visible are required")
		return
	}

	if err := a.chat.SetMessagesVisibility(*request.IdArray, *request.Visible); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "changed")
}

// BanIPAddress will manually ban an IP address.
func (a *Admin) BanIPAddress(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		webutils.WriteSimpleResponse(w, false, "unable to ban IP address")
		return
	}

	if err := a.authRepository.BanIPAddress(configValue.Value.(string), "manually added"); err != nil {
		webutils.WriteSimpleResponse(w, false, "error saving IP address ban")
		return
	}

	webutils.WriteSimpleResponse(w, true, "IP address banned")
}

// UnBanIPAddress will remove an IP address ban.
func (a *Admin) UnBanIPAddress(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		webutils.WriteSimpleResponse(w, false, "unable to unban IP address")
		return
	}

	if err := a.authRepository.RemoveIPAddressBan(configValue.Value.(string)); err != nil {
		webutils.WriteSimpleResponse(w, false, "error removing IP address ban")
		return
	}

	webutils.WriteSimpleResponse(w, true, "IP address unbanned")
}

// GetIPAddressBans will return all the banned IP addresses.
func (a *Admin) GetIPAddressBans(w http.ResponseWriter, r *http.Request) {
	bans, err := a.authRepository.GetIPAddressBans()
	if err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteResponse(w, bans)
}

// UpdateUserEnabled enable or disable a single user by ID.
func (a *Admin) UpdateUserEnabled(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		webutils.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request generated.UpdateUserEnabledJSONBody

	if err := decoder.Decode(&request); err != nil {
		log.Errorln(err)
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	if request.UserId == nil || *request.UserId == "" || request.Enabled == nil {
		webutils.WriteSimpleResponse(w, false, "must provide userId and enabled state")
		return
	}

	if err := a.updateUserStatus(request); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	if !*request.Enabled {
		if err := a.handleUserDisabling(*request.UserId); err != nil {
			webutils.WriteSimpleResponse(w, false, err.Error())
			return
		}
	}

	webutils.WriteSimpleResponse(w, true, fmt.Sprintf("%s enabled: %t", *request.UserId, *request.Enabled))
}

func (a *Admin) updateUserStatus(request generated.UpdateUserEnabledJSONBody) error {
	if err := a.userRepository.SetEnabled(*request.UserId, *request.Enabled); err != nil {
		log.Errorln("error changing user enabled status", err)
		return err
	}

	messageIDs, err := a.chatMessageRepository.GetMessageIdsForUserID(*request.UserId)
	if err != nil {
		return errors.Wrap(err, "error fetching user messages")
	}

	if !*request.Enabled && len(messageIDs) > 0 {
		if err := a.chat.SetMessagesVisibility(messageIDs, *request.Enabled); err != nil {
			log.Errorln("error changing user messages visibility", err)
			return err
		}
	}
	return nil
}

func (a *Admin) handleUserDisabling(userID string) error {
	clients, err := a.chat.GetClientsForUser(userID)
	if len(clients) == 0 {
		return nil
	}

	if err != nil {
		log.Errorln("error fetching clients for user: ", err)
		return err
	}

	a.chat.DisconnectClients(clients)
	disconnectedUser := a.userRepository.GetUserByID(userID)
	_ = a.chat.SendSystemAction(fmt.Sprintf("**%s** has been removed from chat.", disconnectedUser.DisplayName), true)

	localIP4Address := "127.0.0.1"
	localIP6Address := "::1"

	for _, client := range clients {
		ipAddress := client.IPAddress
		if ipAddress != localIP4Address && ipAddress != localIP6Address {
			reason := fmt.Sprintf("Banning of %s", disconnectedUser.DisplayName)
			if err := a.authRepository.BanIPAddress(ipAddress, reason); err != nil {
				log.Errorln("error banning IP address: ", err)
			}
		}
	}
	return nil
}

// GetDisabledUsers will return all the disabled users.
func (a *Admin) GetDisabledUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := a.userRepository.GetDisabledUsers()
	webutils.WriteResponse(w, users)
}

// UpdateUserModerator will set the moderator status for a user ID.
func (a *Admin) UpdateUserModerator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		webutils.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var req generated.UpdateUserModeratorJSONBody

	if err := decoder.Decode(&req); err != nil {
		webutils.WriteSimpleResponse(w, false, "")
		return
	}

	// Update the user object with new moderation access.
	if err := a.userRepository.SetModerator(*req.UserId, *req.IsModerator); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update the clients for this user to know about the moderator access change.
	if err := a.chat.SendConnectedClientInfoToUser(*req.UserId); err != nil {
		log.Debugln(err)
	}

	webutils.WriteSimpleResponse(w, true, fmt.Sprintf("%s is moderator: %t", *req.UserId, *req.IsModerator))
}

// GetModerators will return a list of moderator users.
func (a *Admin) GetModerators(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := a.userRepository.GetModeratorUsers()
	webutils.WriteResponse(w, users)
}

// GetChatMessages returns all of the chat messages, unfiltered.
func (a *Admin) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	messages := a.chatMessageRepository.GetChatModerationHistory()
	webutils.WriteResponse(w, messages)
}

// SendSystemMessage will send an official "SYSTEM" message to chat on behalf of your server.
func (a *Admin) SendSystemMessage(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var message events.SystemMessageEvent
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	if err := a.chat.SendSystemMessage(message.Body, false); err != nil {
		webutils.BadRequestHandler(w, err)
	}

	webutils.WriteSimpleResponse(w, true, "sent")
}

// SendSystemMessageToConnectedClient will handle incoming requests to send a single message to a single connected client by ID.
func (a *Admin) SendSystemMessageToConnectedClient(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	clientIDText, err := utils.GetURLParam(r, "clientId")
	if err != nil {
		webutils.BadRequestHandler(w, err)
		return
	}

	clientIDNumeric, err := strconv.ParseUint(clientIDText, 10, 32)
	if err != nil {
		webutils.BadRequestHandler(w, err)
		return
	}

	// var message events.SystemMessageEvent
	var message generated.SendSystemMessageJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	if message.Body == nil {
		webutils.WriteSimpleResponse(w, false, "no message body provided")
		return
	}

	a.chat.SendSystemMessageToClient(uint(clientIDNumeric), *message.Body)
	webutils.WriteSimpleResponse(w, true, "sent")
}

// SendUserMessage will send a message to chat on behalf of a user. *Depreciated*.
func (a *Admin) SendUserMessage(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	webutils.BadRequestHandler(w, errors.New("no longer supported. see /api/integrations/chat/send"))
}

// SendIntegrationChatMessage will send a chat message on behalf of an external chat integration.
func (a *Admin) SendIntegrationChatMessage(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := integration.DisplayName

	if name == "" {
		webutils.BadRequestHandler(w, errors.New("unknown integration for provided access token"))
		return
	}

	var event events.UserMessageEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}
	event.SetDefaults()
	event.RenderBody()
	event.Type = "CHAT"

	if event.Empty() {
		webutils.BadRequestHandler(w, errors.New("invalid message"))
		return
	}

	event.User = &models.User{
		ID:           integration.ID,
		DisplayName:  name,
		DisplayColor: integration.DisplayColor,
		CreatedAt:    integration.CreatedAt,
		IsBot:        true,
	}

	if err := a.chat.BroadcastEvent(&event); err != nil {
		webutils.BadRequestHandler(w, err)
		return
	}

	a.chatMessageRepository.SaveUserMessage(event)

	webutils.WriteSimpleResponse(w, true, "sent")
}

// SendChatAction will send a generic chat action.
func (a *Admin) SendChatAction(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var message events.SystemActionEvent
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	message.SetDefaults()
	message.RenderBody()

	if err := a.chat.SendSystemAction(message.Body, false); err != nil {
		webutils.BadRequestHandler(w, err)
		return
	}

	webutils.WriteSimpleResponse(w, true, "sent")
}

// SetEnableEstablishedChatUserMode sets the requirement for a chat user
// to be "established" for some time before taking part in chat.
func (a *Admin) SetEnableEstablishedChatUserMode(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		webutils.WriteSimpleResponse(w, false, "unable to update chat established user only mode")
		return
	}

	if err := a.configRepository.SetChatEstablishedUsersOnlyMode(configValue.Value.(bool)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "chat established users only mode updated")
}
