package controllers

import (
	"net/http"
	"strings"

	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/models"
	"github.com/owncast/owncast/activitypub/outbox"
	"github.com/owncast/owncast/activitypub/requests"
	log "github.com/sirupsen/logrus"
)

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
