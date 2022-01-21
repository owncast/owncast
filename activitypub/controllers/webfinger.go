package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

// WebfingerHandler will handle webfinger lookup requests.
func WebfingerHandler(w http.ResponseWriter, r *http.Request) {
	if !data.GetFederationEnabled() {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	resource := r.URL.Query().Get("resource")
	resourceComponents := strings.Split(resource, ":")

	var account string
	if len(resourceComponents) == 2 {
		account = resourceComponents[1]
	} else {
		account = resourceComponents[0]
	}

	userComponents := strings.Split(account, "@")
	if len(userComponents) < 2 {
		return
	}
	host := userComponents[1]
	user := userComponents[0]

	if _, valid := data.GetFederatedInboxMap()[user]; !valid {
		// User is not valid
		w.WriteHeader(http.StatusNotFound)
		log.Debugln("webfinger request rejected for user:", user)
		return
	}

	// If the webfinger request doesn't match our server then it
	// should be rejected.
	instanceHostString := data.GetServerURL()
	if instanceHostString == "" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	instanceHostString = utils.GetHostnameFromURLString(instanceHostString)
	if instanceHostString == "" || instanceHostString != host {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	webfingerResponse := apmodels.MakeWebfingerResponse(user, user, host)

	w.Header().Set("Content-Type", "application/jrd+json")

	if err := json.NewEncoder(w).Encode(webfingerResponse); err != nil {
		log.Errorln("unable to write webfinger response", err)
	}
}
