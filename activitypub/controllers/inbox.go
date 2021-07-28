package controllers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/owncast/owncast/activitypub/inbox"
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

	inbox.Add(data)
	w.WriteHeader(http.StatusAccepted)
}

func returnInbox(w http.ResponseWriter, r *http.Request) {
}
