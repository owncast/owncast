package activitypub

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

	"github.com/owncast/owncast/internal/activitypub/follower"
)

const (
	followersPageSize = 50
)

// FollowersHandler will return the list of remote followers on the Fediverse.
func (s *Service) FollowersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var response interface{}
	var err error
	if r.URL.Query().Get("page") != "" {
		response, err = s.getFollowersPage(r.URL.Query().Get("page"), r)
	} else {
		response, err = s.getInitialFollowersRequest(r)
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
	actorIRI := s.Models.MakeLocalIRIForAccount(accountName)
	publicKey := s.Crypto.GetPublicKey(actorIRI)

	if err := follower.WriteStreamResponse(response.(vocab.Type), w, s.Crypto, publicKey); err != nil {
		log.Errorln("unable to write stream response for followers handler", err)
	}
}

func (s *Service) getInitialFollowersRequest(r *http.Request) (vocab.ActivityStreamsOrderedCollection, error) {
	followerCount, _ := s.Follower.GetFollowerCount()
	collection := streams.NewActivityStreamsOrderedCollection()
	idProperty := streams.NewJSONLDIdProperty()
	id, err := s.createPageURL(r, nil)
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
	firstIRI, err := s.createPageURL(r, &page)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create first page property")
	}

	first.SetIRI(firstIRI)
	collection.SetActivityStreamsFirst(first)

	return collection, nil
}

func (s *Service) getFollowersPage(page string, r *http.Request) (vocab.ActivityStreamsOrderedCollectionPage, error) {
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse page number")
	}

	followerCount, err := s.Follower.GetFollowerCount()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get follower count")
	}

	followers, _, err := s.Follower.GetFederationFollowers(followersPageSize, (pageInt-1)*followersPageSize)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get federation followers")
	}

	collectionPage := streams.NewActivityStreamsOrderedCollectionPage()
	idProperty := streams.NewJSONLDIdProperty()
	id, err := s.createPageURL(r, &page)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create followers page ID")
	}
	idProperty.SetIRI(id)
	collectionPage.SetJSONLDId(idProperty)

	orderedItems := streams.NewActivityStreamsOrderedItemsProperty()

	for _, follower := range followers {
		u, _ := url.Parse(follower.ActorIRI)
		orderedItems.AppendIRI(u)
	}
	collectionPage.SetActivityStreamsOrderedItems(orderedItems)

	partOf := streams.NewActivityStreamsPartOfProperty()
	partOfIRI, err := s.createPageURL(r, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create partOf property for followers page")
	}

	partOf.SetIRI(partOfIRI)
	collectionPage.SetActivityStreamsPartOf(partOf)

	if pageInt*followersPageSize < int(followerCount) {
		next := streams.NewActivityStreamsNextProperty()
		nextPage := fmt.Sprintf("%d", pageInt+1)
		nextIRI, err := s.createPageURL(r, &nextPage)
		if err != nil {
			return nil, errors.Wrap(err, "unable to create next page property")
		}

		next.SetIRI(nextIRI)
		collectionPage.SetActivityStreamsNext(next)
	}

	return collectionPage, nil
}

func (s *Service) createPageURL(r *http.Request, page *string) (*url.URL, error) {
	domain := s.Persistence.Data.GetServerURL()
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
