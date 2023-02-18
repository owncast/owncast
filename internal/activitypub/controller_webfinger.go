package activitypub

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/utils"
)

// WebfingerHandler will handle webfinger lookup requests.
func (s *Service) WebfingerHandler(w http.ResponseWriter, r *http.Request) {
	if !s.Persistence.Data.GetFederationEnabled() {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Debugln("webfinger request rejected! Federation is not enabled")
		return
	}

	instanceHostURL := s.Persistence.Data.GetServerURL()
	if instanceHostURL == "" {
		w.WriteHeader(http.StatusNotFound)
		log.Warnln("webfinger request rejected! Federation is enabled but server URL is empty.")
		return
	}

	instanceHostString := utils.GetHostnameFromURLString(instanceHostURL)
	if instanceHostString == "" {
		w.WriteHeader(http.StatusNotFound)
		log.Warnln("webfinger request rejected! Federation is enabled but server URL is not set properly. s.Persistence.Data.GetServerURL(): " + s.Persistence.Data.GetServerURL())
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

	if _, valid := s.Persistence.Data.GetFederatedInboxMap()[user]; !valid {
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

	webfingerResponse := s.Models.MakeWebfingerResponse(user, user, host)

	w.Header().Set("Content-Type", "application/jrd+json")

	if err := json.NewEncoder(w).Encode(webfingerResponse); err != nil {
		log.Errorln("unable to write webfinger response", err)
	}
}
