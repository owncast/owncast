package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/resolvers"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/models"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

var _datastore *data.Datastore

// Setup will initialize the ActivityPub persistence layer with the provided datastore.
// ActivityPub-related tables (ap_followers, ap_outbox, ap_accepted_activities)
// are created by the goose migrations package.
func Setup(datastore *data.Datastore) {
	_datastore = datastore
	addFollowersFixtureData()
}

// GetDatastore returns the datastore instance for use by sub-repositories.
func GetDatastore() *data.Datastore {
	return _datastore
}

// GetOutboxPostCount will return the number of posts in the outbox.
func GetOutboxPostCount() (int64, error) {
	ctx := context.Background()
	return _datastore.GetQueries().GetLocalPostCount(ctx)
}

// GetOutbox will return an instance of the outbox populated by stored items.
func GetOutbox(limit int, offset int) (vocab.ActivityStreamsOrderedCollection, error) {
	collection := streams.NewActivityStreamsOrderedCollection()
	orderedItems := streams.NewActivityStreamsOrderedItemsProperty()
	rows, err := _datastore.GetQueries().GetOutboxWithOffset(
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
		if err := resolvers.Resolve(context.Background(), value, createCallback); err != nil {
			return collection, err
		}
	}

	return collection, nil
}

// AddToOutbox will store a single payload to the persistence layer.
func AddToOutbox(iri string, itemData []byte, typeString string, isLiveNotification bool) error {
	tx, err := _datastore.DB.Begin()
	if err != nil {
		log.Debugln(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err = _datastore.GetQueries().WithTx(tx).AddToOutbox(context.Background(), db.AddToOutboxParams{
		Iri:              iri,
		Value:            itemData,
		Type:             typeString,
		LiveNotification: sql.NullBool{Bool: isLiveNotification, Valid: true},
	}); err != nil {
		return fmt.Errorf("error creating new item in federation outbox %s", err)
	}

	return tx.Commit()
}

// GetObjectByIRI will return a string representation of a single object by the IRI.
func GetObjectByIRI(iri string) (string, bool, time.Time, error) {
	row, err := _datastore.GetQueries().GetObjectFromOutboxByIRI(context.Background(), iri)
	return string(row.Value), row.LiveNotification.Bool, row.CreatedAt.Time, err
}

// GetLocalPostCount will return the number of posts existing locally.
func GetLocalPostCount() (int64, error) {
	ctx := context.Background()
	return _datastore.GetQueries().GetLocalPostCount(ctx)
}

// SaveInboundFediverseActivity will save an event to the ap_inbound_activities table.
func SaveInboundFediverseActivity(objectIRI string, actorIRI string, eventType string, timestamp time.Time) error {
	if err := _datastore.GetQueries().AddToAcceptedActivities(context.Background(), db.AddToAcceptedActivitiesParams{
		Iri:       objectIRI,
		Actor:     actorIRI,
		Type:      eventType,
		Timestamp: timestamp,
	}); err != nil {
		return errors.Wrap(err, "error saving event "+objectIRI)
	}

	return nil
}

// GetInboundActivities will return a collection of saved, federated activities
// limited and offset by the values provided to support pagination.
func GetInboundActivities(limit int, offset int) ([]models.FederatedActivity, int, error) {
	ctx := context.Background()
	rows, err := _datastore.GetQueries().GetInboundActivitiesWithOffset(ctx, db.GetInboundActivitiesWithOffsetParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	activities := make([]models.FederatedActivity, 0)

	total, err := _datastore.GetQueries().GetInboundActivityCount(context.Background())
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

// HasPreviouslyHandledInboundActivity will return if we have previously handled
// an inbound federated activity.
func HasPreviouslyHandledInboundActivity(iri string, actorIRI string, eventType string) (bool, error) {
	exists, err := _datastore.GetQueries().DoesInboundActivityExist(context.Background(), db.DoesInboundActivityExistParams{
		Iri:   iri,
		Actor: actorIRI,
		Type:  eventType,
	})
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}
