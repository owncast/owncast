package admin

// this is endpoint logic

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/core/chat"
	log "github.com/sirupsen/logrus"
)

// UpdateMessageVisibility updates an array of message IDs to have the same visiblity.
func UpdateMessageVisibility(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
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

	// // make sql update call here.
	// // := means create a new var
	// _db := data.GetDatabase()
	// updateMessageVisibility(_db, request)

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
	// middleware.EnableCors(&w)
	w.Header().Set("Content-Type", "application/json")

	messages := core.GetAllChatMessages(false)

	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Errorln(err)
	}
}
