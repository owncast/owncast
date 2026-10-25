package controllers

import (
	"errors"
	"io"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/services/activitypub/apmodels"
)

// InboxHandler handles inbound federated requests.
func (c *Controllers) InboxHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c.acceptInboxRequest(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (c *Controllers) acceptInboxRequest(w http.ResponseWriter, r *http.Request) {
	if !c.configRepository.GetFederationEnabled() {
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
	if forLocalAccount != c.configRepository.GetFederationUsername() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MiB — consistent with indieauth server handler
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorln("Unable to read inbox request payload", err)
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	inboxRequest := apmodels.InboxRequest{Request: r, ForLocalAccount: forLocalAccount, Body: data}
	c.inbox.AddToQueue(inboxRequest)
	w.WriteHeader(http.StatusAccepted)
}
