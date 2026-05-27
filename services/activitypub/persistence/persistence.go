// Package persistence is the storage layer for the ActivityPub
// subsystem: federated outbox posts, accepted inbound activities, and
// the followers fixture loader. Construct one in main.go (or via
// activitypub.New) with the datastore; everything else goes through the
// returned *Service.
package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/models"
	apresolvers "github.com/owncast/owncast/services/activitypub/resolvers"
	"github.com/owncast/owncast/services/datastore"
)

// Service owns the ActivityPub persistence operations.
type Service struct {
	datastore *datastore.Datastore
	resolver  *apresolvers.Resolver
}

// New constructs a persistence Service bound to the given datastore.
// Side effect: loads followers fixture data (no-op in the default build
// tag; populated when the `fixture` build tag is set). The resolver is
// required for outbox deserialization in GetOutbox; pass nil only if
// you do not call that method.
func New(ds *datastore.Datastore, resolver *apresolvers.Resolver) *Service {
	s := &Service{datastore: ds, resolver: resolver}
	s.addFollowersFixtureData()
	return s
}

// Datastore returns the underlying datastore. Provided so sub-services
// (followers repository, etc.) can be constructed from the same handle.
func (s *Service) Datastore() *datastore.Datastore {
	return s.datastore
}

// GetOutboxPostCount returns the number of posts in the outbox.
func (s *Service) GetOutboxPostCount() (int64, error) {
	ctx := context.Background()
	return s.datastore.GetQueries().GetLocalPostCount(ctx)
}

// GetOutbox returns an instance of the outbox populated by stored items.
func (s *Service) GetOutbox(limit int, offset int) (vocab.ActivityStreamsOrderedCollection, error) {
	collection := streams.NewActivityStreamsOrderedCollection()
	orderedItems := streams.NewActivityStreamsOrderedItemsProperty()
	rows, err := s.datastore.GetQueries().GetOutboxWithOffset(
		context.Background(),
		db.GetOutboxWithOffsetParams{Limit: int64(limit), Offset: int64(offset)},
	)
	if err != nil {
		return collection, err
	}

	for _, value := range rows {
		createCallback := func(c context.Context, activity vocab.ActivityStreamsCreate) error {
			orderedItems.AppendActivityStreamsCreate(activity)
			return nil
		}
		if err := s.resolver.Resolve(context.Background(), value, createCallback); err != nil {
			return collection, err
		}
	}

	return collection, nil
}

// AddToOutbox stores a single payload to the persistence layer.
func (s *Service) AddToOutbox(iri string, itemData []byte, typeString string, isLiveNotification bool) error {
	tx, err := s.datastore.DB.Begin()
	if err != nil {
		log.Debugln(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err = s.datastore.GetQueries().WithTx(tx).AddToOutbox(context.Background(), db.AddToOutboxParams{
		Iri:              iri,
		Value:            itemData,
		Type:             typeString,
		LiveNotification: sql.NullBool{Bool: isLiveNotification, Valid: true},
	}); err != nil {
		return fmt.Errorf("error creating new item in federation outbox %s", err)
	}

	return tx.Commit()
}

// GetObjectByIRI returns a string representation of a single object by the IRI.
func (s *Service) GetObjectByIRI(iri string) (string, bool, time.Time, error) {
	row, err := s.datastore.GetQueries().GetObjectFromOutboxByIRI(context.Background(), iri)
	return string(row.Value), row.LiveNotification.Bool, row.CreatedAt.Time, err
}

// GetLocalPostCount returns the number of posts existing locally.
func (s *Service) GetLocalPostCount() (int64, error) {
	ctx := context.Background()
	return s.datastore.GetQueries().GetLocalPostCount(ctx)
}

// SaveInboundFediverseActivity saves an event to the ap_inbound_activities table.
func (s *Service) SaveInboundFediverseActivity(objectIRI string, actorIRI string, eventType string, timestamp time.Time) error {
	if err := s.datastore.GetQueries().AddToAcceptedActivities(context.Background(), db.AddToAcceptedActivitiesParams{
		Iri:       objectIRI,
		Actor:     actorIRI,
		Type:      eventType,
		Timestamp: timestamp,
	}); err != nil {
		return errors.Wrap(err, "error saving event "+objectIRI)
	}

	return nil
}

// GetInboundActivities returns a paginated collection of saved
// federated activities along with the total count.
func (s *Service) GetInboundActivities(limit int, offset int) ([]models.FederatedActivity, int, error) {
	ctx := context.Background()
	rows, err := s.datastore.GetQueries().GetInboundActivitiesWithOffset(ctx, db.GetInboundActivitiesWithOffsetParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	activities := make([]models.FederatedActivity, 0)

	total, err := s.datastore.GetQueries().GetInboundActivityCount(context.Background())
	if err != nil {
		return nil, 0, errors.Wrap(err, "unable to fetch total activity count")
	}

	for _, row := range rows {
		singleActivity := models.FederatedActivity{
			IRI:       row.Iri,
			ActorIRI:  row.Actor,
			Type:      row.Type,
			Timestamp: row.Timestamp,
		}
		activities = append(activities, singleActivity)
	}

	return activities, int(total), nil
}

// HasPreviouslyHandledInboundActivity reports whether we have previously
// handled an inbound federated activity matching iri/actorIRI/eventType.
func (s *Service) HasPreviouslyHandledInboundActivity(iri string, actorIRI string, eventType string) (bool, error) {
	exists, err := s.datastore.GetQueries().DoesInboundActivityExist(context.Background(), db.DoesInboundActivityExistParams{
		Iri:   iri,
		Actor: actorIRI,
		Type:  eventType,
	})
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}
