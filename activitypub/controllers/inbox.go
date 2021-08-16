package controllers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/owncast/owncast/activitypub/inbox"
	"github.com/owncast/owncast/activitypub/models"

	log "github.com/sirupsen/logrus"
)

func InboxHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		acceptInboxRequest(w, r)
	} else if r.Method == "GET" {
		returnInbox(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func acceptInboxRequest(w http.ResponseWriter, r *http.Request) {
	// https://gek-ap-test.ngrok.io/federation/user/live/inbox
	urlPathComponents := strings.Split(r.URL.Path, "/")
	var forLocalAccount string
	if len(urlPathComponents) == 5 {
		forLocalAccount = urlPathComponents[3]
	} else {
		log.Errorln("Unable to determine username from url path", r.URL.Path)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorln("Unable to read inbox request payload", err)
		return
	}

	fmt.Println("INBOX: ", string(data))

	inboxRequest := models.InboxRequest{Request: r, ForLocalAccount: forLocalAccount, Body: data}

	inbox.Add(inboxRequest)
	w.WriteHeader(http.StatusAccepted)
}

func returnInbox(w http.ResponseWriter, r *http.Request) {
}
