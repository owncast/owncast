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

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/internal/activitypub/follower"
	"github.com/owncast/owncast/models"
)

// New will initialize the ActivityPub persistence layer service
// with the provided Data.
func New(ds *data.Service) (*Service, error) {
	if ds == nil {
		return nil, errors.New("Data store is nil")
	}

	s := &Service{
		Data: ds,
	}

	s.createFederationOutboxTable()
	s.createFederatedActivitiesTable()

	return s, nil
}

type Service struct {
	Data     *data.Service
	follower *follower.Service
}

// createFederatedActivitiesTable will create the accepted
// activities table if needed.
func (s *Service) createFederatedActivitiesTable() {
	createTableSQL := `CREATE TABLE IF NOT EXISTS ap_accepted_activities (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"iri" TEXT NOT NULL,
    "actor" TEXT NOT NULL,
    "type" TEXT NOT NULL,
		"timestamp" TIMESTAMP NOT NULL
	);`

	s.Data.Store.MustExec(createTableSQL)
	s.Data.Store.MustExec(`CREATE INDEX IF NOT EXISTS idx_iri_actor_index ON ap_accepted_activities (iri,actor);`)
}

func (s *Service) createFederationOutboxTable() {
	log.Traceln("Creating federation outbox table...")
	createTableSQL := `CREATE TABLE IF NOT EXISTS ap_outbox (
		"iri" TEXT NOT NULL,
		"value" BLOB,
		"type" TEXT NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "live_notification" BOOLEAN DEFAULT FALSE,
		PRIMARY KEY (iri));`

	s.Data.Store.MustExec(createTableSQL)
	s.Data.Store.MustExec(`CREATE INDEX IF NOT EXISTS idx_iri ON ap_outbox (iri);`)
	s.Data.Store.MustExec(`CREATE INDEX IF NOT EXISTS idx_type ON ap_outbox (type);`)
	s.Data.Store.MustExec(`CREATE INDEX IF NOT EXISTS idx_live_notification ON ap_outbox (live_notification);`)
}

// GetOutboxPostCount will return the number of posts in the outbox.
func (s *Service) GetOutboxPostCount() (int64, error) {
	ctx := context.Background()
	return s.Data.Store.GetQueries().GetLocalPostCount(ctx)
}

// GetOutbox will return an instance of the outbox populated by stored items.
func (s *Service) GetOutbox(limit int, offset int) (vocab.ActivityStreamsOrderedCollection, error) {
	collection := streams.NewActivityStreamsOrderedCollection()
	orderedItems := streams.NewActivityStreamsOrderedItemsProperty()
	rows, err := s.Data.Store.GetQueries().GetOutboxWithOffset(
		context.Background(),
		db.GetOutboxWithOffsetParams{Limit: int32(limit), Offset: int32(offset)},
	)
	if err != nil {
		return collection, err
	}

	for _, value := range rows {
		createCallback := func(c context.Context, activity vocab.ActivityStreamsCreate) error {
			orderedItems.AppendActivityStreamsCreate(activity)
			return nil
		}

		if err := s.follower.Resolve(context.Background(), value, createCallback); err != nil {
			return collection, err
		}
	}

	return collection, nil
}

// AddToOutbox will store a single payload to the persistence layer.
func (s *Service) AddToOutbox(iri string, itemData []byte, typeString string, isLiveNotification bool) error {
	tx, err := s.Data.Store.DB.Begin()
	if err != nil {
		log.Debugln(err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err = s.Data.Store.GetQueries().WithTx(tx).AddToOutbox(context.Background(), db.AddToOutboxParams{
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
func (s *Service) GetObjectByIRI(iri string) (string, bool, time.Time, error) {
	row, err := s.Data.Store.GetQueries().GetObjectFromOutboxByIRI(context.Background(), iri)
	return string(row.Value), row.LiveNotification.Bool, row.CreatedAt.Time, err
}

// GetLocalPostCount will return the number of posts existing locally.
func (s *Service) GetLocalPostCount() (int64, error) {
	ctx := context.Background()
	return s.Data.Store.GetQueries().GetLocalPostCount(ctx)
}

// SaveInboundFediverseActivity will save an event to the ap_inbound_activities table.
func (s *Service) SaveInboundFediverseActivity(objectIRI string, actorIRI string, eventType string, timestamp time.Time) error {
	if err := s.Data.Store.GetQueries().AddToAcceptedActivities(context.Background(), db.AddToAcceptedActivitiesParams{
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
func (s *Service) GetInboundActivities(limit int, offset int) ([]models.FederatedActivity, int, error) {
	ctx := context.Background()
	rows, err := s.Data.Store.GetQueries().GetInboundActivitiesWithOffset(ctx, db.GetInboundActivitiesWithOffsetParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	activities := make([]models.FederatedActivity, 0)

	total, err := s.Data.Store.GetQueries().GetInboundActivityCount(context.Background())
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
func (s *Service) HasPreviouslyHandledInboundActivity(iri string, actorIRI string, eventType string) (bool, error) {
	exists, err := s.Data.Store.GetQueries().DoesInboundActivityExist(context.Background(), db.DoesInboundActivityExistParams{
		Iri:   iri,
		Actor: actorIRI,
		Type:  eventType,
	})
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}
