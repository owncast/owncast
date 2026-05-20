package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/utils"
)

// WebfingerHandler will handle webfinger lookup requests.
func (c *Controllers) WebfingerHandler(w http.ResponseWriter, r *http.Request) {
	if !c.configRepository.GetFederationEnabled() {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Debugln("webfinger request rejected! Federation is not enabled")
		return
	}

	instanceURL, err := c.builder.GetCanonicalServerURL()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Warnln("webfinger request rejected! Federation is enabled but server URL host cannot be canonicalized: " + c.configRepository.GetServerURL())
		return
	}
	instanceHostString := instanceURL.Host

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
	host, err = utils.CanonicalizeHost(host)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Debugln("webfinger request rejected! Invalid host: " + userComponents[1])
		return
	}

	if _, valid := c.configRepository.GetFederatedInboxMap()[user]; !valid {
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

	webfingerResponse := c.builder.MakeWebfingerResponse(user, user, host)

	w.Header().Set("Content-Type", "application/jrd+json")

	if err := json.NewEncoder(w).Encode(webfingerResponse); err != nil {
		log.Errorln("unable to write webfinger response", err)
	}
}
