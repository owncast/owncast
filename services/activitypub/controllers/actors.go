package controllers

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/services/activitypub/requests"
)

// ActorHandler handles requests for a single actor.
func (c *Controllers) ActorHandler(w http.ResponseWriter, r *http.Request) {
	if !c.configRepository.GetFederationEnabled() {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pathComponents := strings.Split(r.URL.Path, "/")
	accountName := pathComponents[3]

	if _, valid := c.configRepository.GetFederatedInboxMap()[accountName]; !valid {
		// User is not valid
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// If this request is for an actor's inbox then pass
	// the request to the inbox controller.
	if len(pathComponents) == 5 && pathComponents[4] == "inbox" {
		c.InboxHandler(w, r)
		return
	} else if len(pathComponents) == 5 && pathComponents[4] == "outbox" {
		c.OutboxHandler(w, r)
		return
	} else if len(pathComponents) == 5 && pathComponents[4] == "followers" {
		// followers list
		c.FollowersHandler(w, r)
		return
	} else if len(pathComponents) == 5 && pathComponents[4] == "following" {
		// following list (none)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	actorIRI := c.builder.MakeLocalIRIForAccount(accountName)
	publicKey := c.signer.GetPublicKey(actorIRI)
	person := c.builder.MakeServiceForAccount(accountName)

	if err := requests.WriteStreamResponse(person, w, publicKey, c.signer); err != nil {
		log.Errorln("unable to write stream response for actor handler", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
