package controllers

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/core/data"
)

func ActorHandler(w http.ResponseWriter, r *http.Request) {
	if !data.GetFederationEnabled() {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pathComponents := strings.Split(r.URL.Path, "/")
	accountName := pathComponents[3]

	if _, valid := data.GetFederatedInboxMap()[accountName]; !valid {
		// User is not valid
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// If this request is for an actor's inbox then pass
	// the request to the inbox controller.
	if len(pathComponents) == 5 && pathComponents[4] == "inbox" {
		InboxHandler(w, r)
		return
	} else if len(pathComponents) == 5 && pathComponents[4] == "outbox" {
		OutboxHandler(w, r)
		return
		// } else if len(pathComponents) == 5 {
		// 	ActorObjectHandler(w, r)
		// 	return
	} else if len(pathComponents) == 5 && pathComponents[4] == "followers" {
		// followers list
		FollowersHandler(w, r)
		return
	} else if len(pathComponents) == 5 && pathComponents[4] == "following" {
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
