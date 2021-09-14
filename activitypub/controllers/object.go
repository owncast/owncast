package controllers

import (
	"net/http"

	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

// ObjectHandler handles requests for a single federated ActivityPub object.
func ObjectHandler(w http.ResponseWriter, r *http.Request) {
	if !data.GetFederationEnabled() {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	objectID, err := utils.ReadRestURLParameter(r, "object")
	if err != nil || objectID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	object, err := persistence.GetObjectByID(objectID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	accountName := data.GetDefaultFederationUsername()
	actorIRI := apmodels.MakeLocalIRIForAccount(accountName)
	publicKey := crypto.GetPublicKey(actorIRI)

	if err := requests.WriteResponse([]byte(object), w, publicKey); err != nil {
		log.Errorln(err)
	}
}
