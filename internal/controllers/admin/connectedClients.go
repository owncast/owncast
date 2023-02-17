package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/internal/controllers"
	"github.com/owncast/owncast/internal/core/chat"
	"github.com/owncast/owncast/internal/core/user"
)

// GetConnectedChatClients returns currently connected clients.
func GetConnectedChatClients(w http.ResponseWriter, r *http.Request) {
	clients := chat.GetClients()
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(clients); err != nil {
		controllers.InternalErrorHandler(w, err)
	}
}

// ExternalGetConnectedChatClients returns currently connected clients.
func ExternalGetConnectedChatClients(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	GetConnectedChatClients(w, r)
}
