package activitypub

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/internal/activitypub/follower"
	"github.com/owncast/owncast/internal/activitypub/persistence"
)

const (
	outboxPageSize = 50
)

// OutboxHandler will handle requests for the local ActivityPub outbox.
func (s *Service) OutboxHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var response interface{}
	var err error
	if r.URL.Query().Get("page") != "" {
		response, err = s.getOutboxPage(r.URL.Query().Get("page"), r)
	} else {
		response, err = s.getInitialOutboxHandler(r)
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
		log.Errorln("unable to write stream response for outbox handler", err)
	}
}

// ActorObjectHandler will handle the request for a single ActivityPub object.
func (s *Service) ActorObjectHandler(w http.ResponseWriter, r *http.Request) {
	object, _, _, err := s.Persistence.GetObjectByIRI(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
		// controllers.WriteSimpleResponse(w, false, err.Error())
	}

	if _, err := w.Write([]byte(object)); err != nil {
		log.Errorln(err)
	}
}

func (s *Service) getInitialOutboxHandler(r *http.Request) (vocab.ActivityStreamsOrderedCollection, error) {
	collection := streams.NewActivityStreamsOrderedCollection()

	idProperty := streams.NewJSONLDIdProperty()
	id, err := s.createPageURL(r, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create followers page property")
	}
	idProperty.SetIRI(id)
	collection.SetJSONLDId(idProperty)

	p, err := persistence.New(s.Persistence.Data)
	if err != nil {
		return nil, fmt.Errorf("getting persistent data store: %v", err)
	}

	totalPosts, err := p.GetOutboxPostCount()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get outbox post count")
	}
	totalItemsProperty := streams.NewActivityStreamsTotalItemsProperty()
	totalItemsProperty.Set(int(totalPosts))
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

func (s *Service) getOutboxPage(page string, r *http.Request) (vocab.ActivityStreamsOrderedCollectionPage, error) {
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse page number")
	}

	p, err := persistence.New(s.Persistence.Data)
	if err != nil {
		return nil, fmt.Errorf("getting persistent data store: %v", err)
	}

	postCount, err := p.GetOutboxPostCount()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get outbox post count")
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

	outboxItems, err := p.GetOutbox(outboxPageSize, (pageInt-1)*outboxPageSize)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get federation followers")
	}
	orderedItems.AppendActivityStreamsOrderedCollection(outboxItems)
	collectionPage.SetActivityStreamsOrderedItems(orderedItems)

	partOf := streams.NewActivityStreamsPartOfProperty()
	partOfIRI, err := s.createPageURL(r, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create partOf property for outbox page")
	}

	partOf.SetIRI(partOfIRI)
	collectionPage.SetActivityStreamsPartOf(partOf)

	if pageInt*followersPageSize < int(postCount) {
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
