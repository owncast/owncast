package moderation

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/userrepository"
	"github.com/owncast/owncast/webserver/utils"
	log "github.com/sirupsen/logrus"
)

// GetUserDetails returns the details of a chat user for moderators.
func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	type connectedClient struct {
		ConnectedAt  time.Time `json:"connectedAt"`
		UserAgent    string    `json:"userAgent"`
		Geo          string    `json:"geo,omitempty"`
		Id           uint      `json:"id"`
		MessageCount int       `json:"messageCount"`
	}

	type response struct {
		User             *models.User              `json:"user"`
		ConnectedClients []connectedClient         `json:"connectedClients"`
		Messages         []events.UserMessageEvent `json:"messages"`
	}

	pathComponents := strings.Split(r.URL.Path, "/")
	uid := pathComponents[len(pathComponents)-1]

	userRepository := userrepository.Get()

	u := userRepository.GetUserByID(uid)

	if u == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	c, _ := chat.GetClientsForUser(uid)
	clients := make([]connectedClient, len(c))
	for i, c := range c {
		client := connectedClient{
			Id:           c.Id,
			MessageCount: c.MessageCount,
			UserAgent:    c.UserAgent,
			ConnectedAt:  c.ConnectedAt,
		}
		if c.Geo != nil {
			client.Geo = c.Geo.CountryCode
		}

		clients[i] = client
	}

	messages, err := chat.GetMessagesFromUser(uid)
	if err != nil {
		log.Errorln(err)
	}

	res := response{
		User:             u,
		ConnectedClients: clients,
		Messages:         messages,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.InternalErrorHandler(w, err)
	}
}
