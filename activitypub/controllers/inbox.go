package controllers

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/inbox"
	"github.com/owncast/owncast/core/data"

	log "github.com/sirupsen/logrus"
)

// InboxHandler handles inbound federated requests.
func InboxHandler(w http.ResponseWriter, r *http.Request) {
	if verified, err := inbox.Verify(r); err != nil || !verified {
		log.Warnln("Unable to verify remote request", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method == "POST" {
		acceptInboxRequest(w, r)
	} else if r.Method == "GET" { //nolint
		returnInbox(w)
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
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorln("Unable to read inbox request payload", err)
		return
	}

	log.Traceln("INBOX: ", string(data))

	inboxRequest := apmodels.InboxRequest{Request: r, ForLocalAccount: forLocalAccount, Body: data}

	inbox.Add(inboxRequest)
	w.WriteHeader(http.StatusAccepted)
}

func returnInbox(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}
