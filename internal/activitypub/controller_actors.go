package activitypub

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/internal/activitypub/follower"
)

// ActorHandler handles requests for a single actor.
func (s *Service) ActorHandler(w http.ResponseWriter, r *http.Request) {
	if !s.Persistence.Data.GetFederationEnabled() {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pathComponents := strings.Split(r.URL.Path, "/")
	accountName := pathComponents[3]

	if _, valid := s.Persistence.Data.GetFederatedInboxMap()[accountName]; !valid {
		// User is not valid
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// If this request is for an actor's inbox then pass
	// the request to the inbox controller.
	if len(pathComponents) == 5 && pathComponents[4] == "inbox" {
		s.InboxHandler(w, r)
		return
	} else if len(pathComponents) == 5 && pathComponents[4] == "outbox" {
		s.OutboxHandler(w, r)
		return
	} else if len(pathComponents) == 5 && pathComponents[4] == "followers" {
		// followers list
		s.FollowersHandler(w, r)
		return
	} else if len(pathComponents) == 5 && pathComponents[4] == "following" {
		// following list (none)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	actorIRI := s.Models.MakeLocalIRIForAccount(accountName)
	publicKey := s.Crypto.GetPublicKey(actorIRI)
	person := s.Models.MakeServiceForAccount(accountName)

	if err := follower.WriteStreamResponse(person, w, s.Crypto, publicKey); err != nil {
		log.Errorln("unable to write stream response for actor handler", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
