package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gabek/owncast/core"
	"github.com/gabek/owncast/models"
	"github.com/gabek/owncast/router/middleware"
)

//GetChatMessages gets all of the chat messages
func GetChatMessages(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)

	switch r.Method {
	case http.MethodGet:
		messages := core.GetAllChatMessages()

		json.NewEncoder(w).Encode(messages)
	case http.MethodPost:
		var message models.ChatMessage
		if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
			internalErrorHandler(w, err)
			return
		}

		if err := core.SendMessageToChat(message); err != nil {
			badRequestHandler(w, err)
			return
		}

		json.NewEncoder(w).Encode(j{"success": true})
	default:
		w.WriteHeader(http.StatusNotImplemented)
		json.NewEncoder(w).Encode(j{"error": "method not implemented (PRs are accepted)"})
	}
}
