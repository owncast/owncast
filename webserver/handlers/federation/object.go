package federation

import (
	"net/http"
	"strings"

	"github.com/owncast/owncast/services/apfederation/apmodels"
	"github.com/owncast/owncast/services/apfederation/crypto"
	"github.com/owncast/owncast/services/apfederation/requests"
	"github.com/owncast/owncast/storage/configrepository"
	"github.com/owncast/owncast/storage/federationrepository"
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

	federationRespository := federationrepository.Get()
	iri := strings.Join([]string{strings.TrimSuffix(configRepository.GetServerURL(), "/"), r.URL.Path}, "")
	object, _, _, err := federationRespository.GetObjectByIRI(iri)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	accountName := configRepository.GetDefaultFederationUsername()
	actorIRI := apmodels.MakeLocalIRIForAccount(accountName)
	publicKey := crypto.GetPublicKey(actorIRI)

	req := requests.Get()
	if err := req.WriteResponse([]byte(object), w, publicKey); err != nil {
		log.Errorln(err)
	}
}
