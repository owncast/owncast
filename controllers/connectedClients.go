package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/models"
)

// GetConnectedClients returns currently connected clients.
func GetConnectedClients(w http.ResponseWriter, r *http.Request) {
	clients := chat.GetClients()
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(clients); err != nil {
		InternalErrorHandler(w, err)
	}
}

// ExternalGetConnectedClients returns currently connected clients.
func ExternalGetConnectedClients(integration models.ExternalIntegration, w http.ResponseWriter, r *http.Request) {
	GetConnectedClients(w, r)
}
