package controllers

import (
	"net/http"
	"strings"

	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/persistence/configrepository"
	log "github.com/sirupsen/logrus"
)

// ObjectHandler handles requests for a single federated ActivityPub object.
func ObjectHandler(w http.ResponseWriter, r *http.Request) {
	configRepository := configrepository.Get()

	if !configRepository.GetFederationEnabled() {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// If private federation mode is enabled do not allow access to objects.
	if configRepository.GetFederationIsPrivate() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	iri := strings.Join([]string{strings.TrimSuffix(configRepository.GetServerURL(), "/"), r.URL.Path}, "")
	object, _, _, err := persistence.GetObjectByIRI(iri)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	accountName := configRepository.GetDefaultFederationUsername()
	actorIRI := apmodels.MakeLocalIRIForAccount(accountName)
	publicKey := crypto.GetPublicKey(actorIRI)

	if err := requests.WriteResponse([]byte(object), w, publicKey); err != nil {
		log.Errorln(err)
	}
}
