package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
)

var (
	asContext = []string{"https://www.w3.org/ns/activitystreams", "https://w3id.org/security/v1"}
)

func GetActor(w http.ResponseWriter, r *http.Request) {
	actor := models.Actor{
		Context:           asContext,
		Id:                urlFor("/actor"),
		Type:              "Person",
		PreferredUsername: config.Config.Federation.Username,
		Name:              config.Config.InstanceDetails.Name,
		Summary:           config.Config.InstanceDetails.Summary,
		URL:               urlFor("/"),
		Inbox:             urlFor("/inbox"),
		Outbox:            urlFor("/outbox"),
		Following:         urlFor("/following"),
		Followers:         urlFor("/followers"),
	}

	if err := json.NewEncoder(w).Encode(actor); err != nil {
		internalErrorHandler(w, err)
	}
}
