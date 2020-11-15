package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/router/middleware"
	log "github.com/sirupsen/logrus"
)

// GetChatMessages gets all of the chat messages.
func GetChatMessages(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)

	switch r.Method {
	case http.MethodGet:
		messages := core.GetAllChatMessages()

		err := json.NewEncoder(w).Encode(messages)
		if err != nil {
			log.Errorln(err)
		}
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

		if err := json.NewEncoder(w).Encode(j{"success": true}); err != nil {
			internalErrorHandler(w, err)
			return
		}
	default:
		w.WriteHeader(http.StatusNotImplemented)
		if err := json.NewEncoder(w).Encode(j{"error": "method not implemented (PRs are accepted)"}); err != nil {
			internalErrorHandler(w, err)
		}
	}
}
