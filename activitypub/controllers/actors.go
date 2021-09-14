package controllers

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/utils"
)

// ActorHandler handles requests for a single actor.
func ActorHandler(w http.ResponseWriter, r *http.Request) {
	if !data.GetFederationEnabled() {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	accountName, err := utils.ReadRestURLParameter(r, "user")
	if err != nil || accountName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resource, _ := utils.ReadRestURLParameter(r, "resource")

	if _, valid := data.GetFederatedInboxMap()[accountName]; !valid {
		// User is not valid
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// If this request is for an actor's inbox then pass
	// the request to the inbox controller.
	if resource == "inbox" {
		InboxHandler(w, r)
		return
	} else if resource == "outbox" {
		OutboxHandler(w, r)
		return
	} else if resource == "followers" {
		// followers list
		FollowersHandler(w, r)
		return
	} else if resource == "following" {
		// following list (none)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	actorIRI := apmodels.MakeLocalIRIForAccount(accountName)
	publicKey := crypto.GetPublicKey(actorIRI)
	person := apmodels.MakeActor(accountName)

	if err := requests.WriteStreamResponse(person, w, publicKey); err != nil {
		log.Errorln("unable to write stream response for actor handler", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
