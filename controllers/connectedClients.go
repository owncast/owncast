package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/core"
)

// GetConnectedClients returns currently connected clients
func GetConnectedClients(w http.ResponseWriter, r *http.Request) {
	clients := core.GetClients()
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(clients)
}
