package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/core/data"
)

const (
	followersPageSize = 50
)

// FollowersHandler will return the list of remote followers on the Fediverse.
func FollowersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var response interface{}
	var err error
	if r.URL.Query().Get("page") != "" {
		response, err = getFollowersPage(r.URL.Query().Get("page"), r)
	} else {
		response, err = getInitialFollowersRequest(r)
	}

	if response == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pathComponents := strings.Split(r.URL.Path, "/")
	accountName := pathComponents[3]
	actorIRI := apmodels.MakeLocalIRIForAccount(accountName)
	publicKey := crypto.GetPublicKey(actorIRI)

	if err := requests.WriteStreamResponse(response.(vocab.Type), w, publicKey); err != nil {
		log.Errorln("unable to write stream response for followers handler", err)
	}
}

func getInitialFollowersRequest(r *http.Request) (vocab.ActivityStreamsOrderedCollection, error) {
	followerCount, _ := persistence.GetFollowerCount()
	collection := streams.NewActivityStreamsOrderedCollection()
	idProperty := streams.NewJSONLDIdProperty()
	id, err := createPageURL(r, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create followers page property")
	}
	idProperty.SetIRI(id)
	collection.SetJSONLDId(idProperty)

	totalItemsProperty := streams.NewActivityStreamsTotalItemsProperty()
	totalItemsProperty.Set(int(followerCount))
	collection.SetActivityStreamsTotalItems(totalItemsProperty)

	first := streams.NewActivityStreamsFirstProperty()
	page := "1"
	firstIRI, err := createPageURL(r, &page)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create first page property")
	}

	first.SetIRI(firstIRI)
	collection.SetActivityStreamsFirst(first)

	return collection, nil
}

func getFollowersPage(page string, r *http.Request) (vocab.ActivityStreamsOrderedCollectionPage, error) {
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse page number")
	}

	followerCount, err := persistence.GetFollowerCount()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get follower count")
	}

	followers, err := persistence.GetFederationFollowers(followersPageSize, (pageInt-1)*followersPageSize)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get federation followers")
	}

	collectionPage := streams.NewActivityStreamsOrderedCollectionPage()
	idProperty := streams.NewJSONLDIdProperty()
	id, err := createPageURL(r, &page)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create followers page ID")
	}
	idProperty.SetIRI(id)
	collectionPage.SetJSONLDId(idProperty)

	orderedItems := streams.NewActivityStreamsOrderedItemsProperty()

	for _, follower := range followers {
		u, _ := url.Parse(follower.Link)
		orderedItems.AppendIRI(u)
	}
	collectionPage.SetActivityStreamsOrderedItems(orderedItems)

	partOf := streams.NewActivityStreamsPartOfProperty()
	partOfIRI, err := createPageURL(r, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create partOf property for followers page")
	}

	partOf.SetIRI(partOfIRI)
	collectionPage.SetActivityStreamsPartOf(partOf)

	if pageInt*followersPageSize < int(followerCount) {
		next := streams.NewActivityStreamsNextProperty()
		nextPage := fmt.Sprintf("%d", pageInt+1)
		nextIRI, err := createPageURL(r, &nextPage)
		if err != nil {
			return nil, errors.Wrap(err, "unable to create next page property")
		}

		next.SetIRI(nextIRI)
		collectionPage.SetActivityStreamsNext(next)
	}

	return collectionPage, nil
}

func createPageURL(r *http.Request, page *string) (*url.URL, error) {
	domain := data.GetServerURL()
	if domain == "" {
		return nil, errors.New("unable to get server URL")
	}

	pageURL, err := url.Parse(domain)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse server URL")
	}

	if page != nil {
		query := pageURL.Query()
		query.Add("page", *page)
		pageURL.RawQuery = query.Encode()
	}
	pageURL.Path = r.URL.Path

	return pageURL, nil
}
