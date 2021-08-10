package controllers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/inbox"
	"github.com/owncast/owncast/activitypub/models"
	"github.com/owncast/owncast/activitypub/outbox"
	"github.com/owncast/owncast/activitypub/requests"
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

func OutboxHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	orderedCollection := outbox.Get()
	pathComponents := strings.Split(r.URL.Path, "/")
	accountName := pathComponents[3]
	actorIRI := models.MakeLocalIRIForAccount(accountName)
	publicKey := crypto.GetPublicKey(actorIRI)

	if err := requests.WriteStreamResponse(orderedCollection, w, publicKey); err != nil {
		log.Errorln(err)
	}
}

func acceptInboxRequest(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	fmt.Println("INBOX: ", string(data))

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	// TODO: Verify inbound requests.
	// if verified, err := utils.Verify(r); !verified || err != nil {
	// 	fmt.Println("Inbox payload not verified", err)
	// 	return
	// }

	// https://gek-ap-test.ngrok.io/federation/user/live/inbox
	urlPathComponents := strings.Split(r.URL.Path, "/")
	var forLocalAccount string
	if len(urlPathComponents) == 5 {
		forLocalAccount = urlPathComponents[3]
	} else {
		log.Errorln("Unable to determine username from url path", r.URL.Path)
		return
	}

	inbox.Add(data, forLocalAccount)
	w.WriteHeader(http.StatusAccepted)
}

func returnInbox(w http.ResponseWriter, r *http.Request) {
}
