package admin

// this is endpoint logic

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

// UpdateMessageVisibility updates an array of message IDs to have the same visiblity.
func UpdateMessageVisibility(w http.ResponseWriter, r *http.Request) {
	if r.Method != controllers.POST {
		controllers.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request messageVisibilityUpdateRequest // creates an empty struc

	err := decoder.Decode(&request) // decode the json into `request`
	if err != nil {
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

type messageVisibilityUpdateRequest struct {
	IDArray []string `json:"idArray"`
	Visible bool     `json:"visible"`
}

// GetChatMessages returns all of the chat messages, unfiltered.
func GetChatMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	messages := core.GetAllChatMessages(false)

	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Errorln(err)
	}
}

// SendSystemMessage will send an official "SYSTEM" message
// to chat on behalf of your server.
func SendSystemMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var message models.ChatEvent
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		controllers.InternalErrorHandler(w, err)
		return
	}

	message.MessageType = models.SystemMessageSent
	message.Author = config.Config.InstanceDetails.Name
	message.ClientID = "owncast-server"
	message.ID = shortid.MustGenerate()
	message.Visible = true

	message.SetDefaults()
	message.RenderAndSanitizeMessageBody()

	if err := core.SendMessageToChat(message); err != nil {
		controllers.BadRequestHandler(w, err)
		return
	}

	controllers.WriteSimpleResponse(w, true, "sent")
}

// SendUserMessage will send a message to chat on behalf of a user.
func SendUserMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var message models.ChatEvent
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		controllers.InternalErrorHandler(w, err)
		return
	}

	message.MessageType = models.MessageSent
	message.ClientID = "external-request"
	message.ID = shortid.MustGenerate()
	message.Visible = true

	message.SetDefaults()
	message.RenderAndSanitizeMessageBody()

	if err := core.SendMessageToChat(message); err != nil {
		controllers.BadRequestHandler(w, err)
		return
	}

	controllers.WriteSimpleResponse(w, true, "sent")
}
