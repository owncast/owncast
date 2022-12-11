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
		log.Debugln("webfinger request rejected! Federation is not enabled")
		return
	}

	instanceHostURL := data.GetServerURL()
	if instanceHostURL == "" {
		w.WriteHeader(http.StatusNotFound)
		log.Warnln("webfinger request rejected! Federation is enabled but server URL is empty.")
		return
	}

	instanceHostString := utils.GetHostnameFromURLString(instanceHostURL)
	if instanceHostString == "" {
		w.WriteHeader(http.StatusNotFound)
		log.Warnln("webfinger request rejected! Federation is enabled but server URL is not set properly. data.GetServerURL(): " + data.GetServerURL())
		return
	}

	resource := r.URL.Query().Get("resource")
	preAcct, account, foundAcct := strings.Cut(resource, "acct:")

	if !foundAcct || preAcct != "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Debugln("webfinger request rejected! Malformed resource in query: " + resource)
		return
	}

	userComponents := strings.Split(account, "@")
	if len(userComponents) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		log.Debugln("webfinger request rejected! Malformed account in query: " + account)
		return
	}
	host := userComponents[1]
	user := userComponents[0]

	if _, valid := data.GetFederatedInboxMap()[user]; !valid {
		w.WriteHeader(http.StatusNotFound)
		log.Debugln("webfinger request rejected! Invalid user: " + user)
		return
	}

	// If the webfinger request doesn't match our server then it
	// should be rejected.
	if instanceHostString != host {
		w.WriteHeader(http.StatusNotImplemented)
		log.Debugln("webfinger request rejected! Invalid query host: " + host + " instanceHostString: " + instanceHostString)
		return
	}

	webfingerResponse := apmodels.MakeWebfingerResponse(user, user, host)

	w.Header().Set("Content-Type", "application/jrd+json")

	if err := json.NewEncoder(w).Encode(webfingerResponse); err != nil {
		log.Errorln("unable to write webfinger response", err)
	}
}
