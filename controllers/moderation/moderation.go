package moderation

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/user"
)

func New(s *controllers.Service) (*Controller, error) {
	c := &Controller{
		Service: s,
	}

	return c, nil
}

type Controller struct {
	*controllers.Service
}

// GetUserDetails returns the details of a chat user for moderators.
func (c *Controller) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	type connectedClient struct {
		Id           uint      `json:"id"`
		MessageCount int       `json:"messageCount"`
		UserAgent    string    `json:"userAgent"`
		ConnectedAt  time.Time `json:"connectedAt"`
		Geo          string    `json:"geo,omitempty"`
	}

	type response struct {
		User             *user.User                `json:"user"`
		ConnectedClients []connectedClient         `json:"connectedClients"`
		Messages         []events.UserMessageEvent `json:"messages"`
	}

	pathComponents := strings.Split(r.URL.Path, "/")
	uid := pathComponents[len(pathComponents)-1]

	u := user.GetUserByID(uid, c.Data.Store)

	if u == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	clients, _ := c.Core.Chat.GetClientsForUser(uid)
	clientsInfo := make([]connectedClient, len(clients))

	for i, client := range clients {
		info := connectedClient{
			Id:           client.Id,
			MessageCount: client.MessageCount,
			UserAgent:    client.UserAgent,
			ConnectedAt:  client.ConnectedAt,
		}
		if client.Geo != nil {
			info.Geo = client.Geo.CountryCode
		}

		clientsInfo[i] = info
	}

	messages, err := chat.GetMessagesFromUser(uid)
	if err != nil {
		log.Errorln(err)
	}

	res := response{
		User:             u,
		ConnectedClients: clientsInfo,
		Messages:         messages,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		c.Service.InternalErrorHandler(w, err)
	}
}
