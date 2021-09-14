package controllers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/go-fed/activity/streams"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/outbox"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/requests"
	log "github.com/sirupsen/logrus"
)

// OutboxHandler will handle requests for the local ActivityPub outbox.
func OutboxHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	orderedCollection := outbox.Get()
	pathComponents := strings.Split(r.URL.Path, "/")
	accountName := pathComponents[3]
	actorIRI := apmodels.MakeLocalIRIForAccount(accountName)
	publicKey := crypto.GetPublicKey(actorIRI)

	if err := requests.WriteStreamResponse(orderedCollection, w, publicKey); err != nil {
		log.Errorln("unable to write stream response for outbox handler", err)
	}
}

// FollowersHandler will return the list of remote followers on the Fediverse.
func FollowersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	followers, err := persistence.GetFederationFollowers()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	followerCount, _ := persistence.GetFollowerCount()

	collection := streams.NewActivityStreamsOrderedCollectionPage()
	totalItemsProperty := streams.NewActivityStreamsTotalItemsProperty()
	totalItemsProperty.Set(int(followerCount))
	collection.SetActivityStreamsTotalItems(totalItemsProperty)
	orderedItems := streams.NewActivityStreamsOrderedItemsProperty()

	for _, follower := range followers {
		u, _ := url.Parse(follower.Link)
		orderedItems.AppendIRI(u)
	}
	collection.SetActivityStreamsOrderedItems(orderedItems)

	pathComponents := strings.Split(r.URL.Path, "/")
	accountName := pathComponents[3]
	actorIRI := apmodels.MakeLocalIRIForAccount(accountName)
	publicKey := crypto.GetPublicKey(actorIRI)

	if err := requests.WriteStreamResponse(collection, w, publicKey); err != nil {
		log.Errorln("unable to write stream response for followers handler", err)
	}
}

// ActorObjectHandler will handle the request for our local ActivityPub actor.
func ActorObjectHandler(w http.ResponseWriter, r *http.Request) {
	object, err := persistence.GetObjectByIRI(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
		// controllers.WriteSimpleResponse(w, false, err.Error())
	}

	if _, err := w.Write([]byte(object)); err != nil {
		log.Errorln(err)
	}
}
