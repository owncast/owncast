package controllers

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/services/activitypub/requests"
)

// ObjectHandler handles requests for a single federated ActivityPub object.
func (c *Controllers) ObjectHandler(w http.ResponseWriter, r *http.Request) {
	if !c.configRepository.GetFederationEnabled() {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// If private federation mode is enabled do not allow access to objects.
	if c.configRepository.GetFederationIsPrivate() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	serverURL, err := c.builder.GetCanonicalServerURL()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	serverURL.Path = r.URL.Path
	iri := serverURL.String()
	object, _, _, err := c.persistence.GetObjectByIRI(iri)
	if err != nil {
		legacyIRI := strings.Join([]string{strings.TrimSuffix(c.configRepository.GetServerURL(), "/"), r.URL.Path}, "")
		if legacyIRI == iri {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		object, _, _, err = c.persistence.GetObjectByIRI(legacyIRI)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	accountName := c.configRepository.GetDefaultFederationUsername()
	actorIRI := c.builder.MakeLocalIRIForAccount(accountName)
	publicKey := c.signer.GetPublicKey(actorIRI)

	if err := requests.WriteResponse([]byte(object), w, publicKey, c.signer); err != nil {
		log.Errorln(err)
	}
}
