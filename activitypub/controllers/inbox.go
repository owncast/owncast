package controllers

import (
	"io"
	"net/http"
	"strings"

	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/inbox"
	"github.com/owncast/owncast/core/data"

	log "github.com/sirupsen/logrus"
)

// InboxHandler handles inbound federated requests.
func InboxHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		acceptInboxRequest(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func acceptInboxRequest(w http.ResponseWriter, r *http.Request) {
	if !data.GetFederationEnabled() {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	urlPathComponents := strings.Split(r.URL.Path, "/")
	var forLocalAccount string
	if len(urlPathComponents) == 5 {
		forLocalAccount = urlPathComponents[3]
	} else {
		log.Errorln("Unable to determine username from url path", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// The account this request is for must match the account name we have set
	// for federation.
	if forLocalAccount != data.GetFederationUsername() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorln("Unable to read inbox request payload", err)
		return
	}

	inboxRequest := apmodels.InboxRequest{Request: r, ForLocalAccount: forLocalAccount, Body: data}
	inbox.AddToQueue(inboxRequest)
	w.WriteHeader(http.StatusAccepted)
}
