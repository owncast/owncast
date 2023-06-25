package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/webserver/responses"
)

// GetConnectedChatClients returns currently connected clients.
func (h *Handlers) GetConnectedChatClients(w http.ResponseWriter, r *http.Request) {
	clients := chat.GetClients()
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(clients); err != nil {
		responses.InternalErrorHandler(w, err)
	}
}

// ExternalGetConnectedChatClients returns currently connected clients.
func (h *Handlers) ExternalGetConnectedChatClients(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	h.GetConnectedChatClients(w, r)
}
