package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/router/middleware"
	log "github.com/sirupsen/logrus"
)

// GetChatMessages gets all of the chat messages.
func GetChatMessages(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		messages := core.GetAllChatMessages()

		err := json.NewEncoder(w).Encode(messages)
		if err != nil {
			log.Errorln(err)
		}
	default:
		w.WriteHeader(http.StatusNotImplemented)
		if err := json.NewEncoder(w).Encode(j{"error": "method not implemented (PRs are accepted)"}); err != nil {
			InternalErrorHandler(w, err)
		}
	}
}
