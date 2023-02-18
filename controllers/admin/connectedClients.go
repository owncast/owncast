package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/core/user"
)

// GetConnectedChatClients returns currently connected clients.
func (c *Controller) GetConnectedChatClients(w http.ResponseWriter, r *http.Request) {
	clients := c.Core.Chat.GetClients()
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(clients); err != nil {
		c.InternalErrorHandler(w, err)
	}
}

// ExternalGetConnectedChatClients returns currently connected clients.
func (c *Controller) ExternalGetConnectedChatClients(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	c.GetConnectedChatClients(w, r)
}
