package activitypub

import (
	"io"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/internal/activitypub/apmodels"
)

// InboxHandler handles inbound federated requests.
func (s *Service) InboxHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.acceptInboxRequest(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Service) acceptInboxRequest(w http.ResponseWriter, r *http.Request) {
	if !s.Persistence.Data.GetFederationEnabled() {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	urlPathComponents := strings.Split(r.URL.Path, "/")
	var forLocalAccount string
	if len(urlPathComponents) == 5 {
		forLocalAccount = urlPathComponents[3]
	} else {
		log.Errorln("Unable to determine username from url path")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// The account this request is for must match the account name we have set
	// for federation.
	if forLocalAccount != s.Persistence.Data.GetFederationUsername() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorln("Unable to read inbox request payload", err)
		return
	}

	inboxRequest := apmodels.InboxRequest{Request: r, ForLocalAccount: forLocalAccount, Body: data}
	s.Inbox.AddToQueue(inboxRequest)
	w.WriteHeader(http.StatusAccepted)
}
