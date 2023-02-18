package activitypub

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/internal/activitypub/follower"
	"github.com/owncast/owncast/internal/activitypub/persistence"
)

// ObjectHandler handles requests for a single federated ActivityPub object.
func (s *Service) ObjectHandler(w http.ResponseWriter, r *http.Request) {
	if !s.Persistence.Data.GetFederationEnabled() {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// If private federation mode is enabled do not allow access to objects.
	if s.Persistence.Data.GetFederationIsPrivate() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	p, err := persistence.New(s.Persistence.Data)
	if err != nil {
		log.Errorf("getting persistent data store: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	iri := strings.Join([]string{strings.TrimSuffix(s.Persistence.Data.GetServerURL(), "/"), r.URL.Path}, "")
	object, _, _, err := p.GetObjectByIRI(iri)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	accountName := s.Persistence.Data.GetDefaultFederationUsername()
	actorIRI := s.Models.MakeLocalIRIForAccount(accountName)
	publicKey := s.Crypto.GetPublicKey(actorIRI)

	if err := follower.WriteResponse([]byte(object), w, s.Crypto, publicKey); err != nil {
		log.Errorln(err)
	}
}
