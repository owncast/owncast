package admin

// this is endpoint logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/userrepository"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/handlers/generated"
	webutils "github.com/owncast/owncast/webserver/utils"
	log "github.com/sirupsen/logrus"
)

// ExternalUpdateMessageVisibility updates an array of message IDs to have the same visiblity.
func ExternalUpdateMessageVisibility(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	UpdateMessageVisibility(w, r)
}

// UpdateMessageVisibility updates an array of message IDs to have the same visiblity.
func UpdateMessageVisibility(w http.ResponseWriter, r *http.Request) {
	// type messageVisibilityUpdateRequest struct {
	// 	IDArray []string `json:"idArray"`
	// 	Visible bool     `json:"visible"`
	// }

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

	if err := chat.SetMessagesVisibility(*request.IdArray, *request.Visible); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "changed")
}

// BanIPAddress will manually ban an IP address.
func BanIPAddress(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		webutils.WriteSimpleResponse(w, false, "unable to ban IP address")
		return
	}

	if err := data.BanIPAddress(configValue.Value.(string), "manually added"); err != nil {
		webutils.WriteSimpleResponse(w, false, "error saving IP address ban")
		return
	}

	webutils.WriteSimpleResponse(w, true, "IP address banned")
}

// UnBanIPAddress will remove an IP address ban.
func UnBanIPAddress(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		webutils.WriteSimpleResponse(w, false, "unable to unban IP address")
		return
	}

	if err := data.RemoveIPAddressBan(configValue.Value.(string)); err != nil {
		webutils.WriteSimpleResponse(w, false, "error removing IP address ban")
		return
	}

	webutils.WriteSimpleResponse(w, true, "IP address unbanned")
}

// GetIPAddressBans will return all the banned IP addresses.
func GetIPAddressBans(w http.ResponseWriter, r *http.Request) {
	bans, err := data.GetIPAddressBans()
	if err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteResponse(w, bans)
}

// UpdateUserEnabled enable or disable a single user by ID.
func UpdateUserEnabled(w http.ResponseWriter, r *http.Request) {
	type blockUserRequest struct {
		UserID  string `json:"userId"`
		Enabled bool   `json:"enabled"`
	}

	if r.Method != http.MethodPost {
		webutils.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request blockUserRequest

	if err := decoder.Decode(&request); err != nil {
		log.Errorln(err)
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	if request.UserID == "" {
		webutils.WriteSimpleResponse(w, false, "must provide userId")
		return
	}

	userRepository := userrepository.Get()

	// Disable/enable the user
	if err := userRepository.SetEnabled(request.UserID, request.Enabled); err != nil {
		log.Errorln("error changing user enabled status", err)
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Hide/show the user's chat messages if disabling.
	// Leave hidden messages hidden to be safe.
	if !request.Enabled {
		if err := chat.SetMessageVisibilityForUserID(request.UserID, request.Enabled); err != nil {
			log.Errorln("error changing user messages visibility", err)
			webutils.WriteSimpleResponse(w, false, err.Error())
			return
		}
	}

	// Forcefully disconnect the user from the chat
	if !request.Enabled {
		clients, err := chat.GetClientsForUser(request.UserID)
		if len(clients) == 0 {
			// Nothing to do
			return
		}

		if err != nil {
			log.Errorln("error fetching clients for user: ", err)
			webutils.WriteSimpleResponse(w, false, err.Error())
			return
		}

		chat.DisconnectClients(clients)
		disconnectedUser := userRepository.GetUserByID(request.UserID)
		_ = chat.SendSystemAction(fmt.Sprintf("**%s** has been removed from chat.", disconnectedUser.DisplayName), true)

		localIP4Address := "127.0.0.1"
		localIP6Address := "::1"

		// Ban this user's IP address.
		for _, client := range clients {
			ipAddress := client.IPAddress
			if ipAddress != localIP4Address && ipAddress != localIP6Address {
				reason := fmt.Sprintf("Banning of %s", disconnectedUser.DisplayName)
				if err := data.BanIPAddress(ipAddress, reason); err != nil {
					log.Errorln("error banning IP address: ", err)
				}
			}
		}
	}

	webutils.WriteSimpleResponse(w, true, fmt.Sprintf("%s enabled: %t", request.UserID, request.Enabled))
}

// GetDisabledUsers will return all the disabled users.
func GetDisabledUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userRepository := userrepository.Get()

	users := userRepository.GetDisabledUsers()
	webutils.WriteResponse(w, users)
}

// UpdateUserModerator will set the moderator status for a user ID.
func UpdateUserModerator(w http.ResponseWriter, r *http.Request) {
	type request struct {
		UserID      string `json:"userId"`
		IsModerator bool   `json:"isModerator"`
	}

	if r.Method != http.MethodPost {
		webutils.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var req request

	if err := decoder.Decode(&req); err != nil {
		webutils.WriteSimpleResponse(w, false, "")
		return
	}

	userRepository := userrepository.Get()

	// Update the user object with new moderation access.
	if err := userRepository.SetModerator(req.UserID, req.IsModerator); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update the clients for this user to know about the moderator access change.
	if err := chat.SendConnectedClientInfoToUser(req.UserID); err != nil {
		log.Debugln(err)
	}

	webutils.WriteSimpleResponse(w, true, fmt.Sprintf("%s is moderator: %t", req.UserID, req.IsModerator))
}

// GetModerators will return a list of moderator users.
func GetModerators(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userRepository := userrepository.Get()

	users := userRepository.GetModeratorUsers()
	webutils.WriteResponse(w, users)
}

// GetChatMessages returns all of the chat messages, unfiltered.
func GetChatMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	messages := chat.GetChatModerationHistory()
	webutils.WriteResponse(w, messages)
}

// SendSystemMessage will send an official "SYSTEM" message to chat on behalf of your server.
func SendSystemMessage(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var message events.SystemMessageEvent
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	if err := chat.SendSystemMessage(message.Body, false); err != nil {
		webutils.BadRequestHandler(w, err)
	}

	webutils.WriteSimpleResponse(w, true, "sent")
}

// SendSystemMessageToConnectedClient will handle incoming requests to send a single message to a single connected client by ID.
func SendSystemMessageToConnectedClient(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
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

	var message events.SystemMessageEvent
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	chat.SendSystemMessageToClient(uint(clientIDNumeric), message.Body)
	webutils.WriteSimpleResponse(w, true, "sent")
}

// SendUserMessage will send a message to chat on behalf of a user. *Depreciated*.
func SendUserMessage(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	webutils.BadRequestHandler(w, errors.New("no longer supported. see /api/integrations/chat/send"))
}

// SendIntegrationChatMessage will send a chat message on behalf of an external chat integration.
func SendIntegrationChatMessage(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
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

	if err := chat.Broadcast(&event); err != nil {
		webutils.BadRequestHandler(w, err)
		return
	}

	chat.SaveUserMessage(event)

	webutils.WriteSimpleResponse(w, true, "sent")
}

// SendChatAction will send a generic chat action.
func SendChatAction(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var message events.SystemActionEvent
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	message.SetDefaults()
	message.RenderBody()

	if err := chat.SendSystemAction(message.Body, false); err != nil {
		webutils.BadRequestHandler(w, err)
		return
	}

	webutils.WriteSimpleResponse(w, true, "sent")
}

// SetEnableEstablishedChatUserMode sets the requirement for a chat user
// to be "established" for some time before taking part in chat.
func SetEnableEstablishedChatUserMode(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		webutils.WriteSimpleResponse(w, false, "unable to update chat established user only mode")
		return
	}

	if err := data.SetChatEstablishedUsersOnlyMode(configValue.Value.(bool)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "chat established users only mode updated")
}
