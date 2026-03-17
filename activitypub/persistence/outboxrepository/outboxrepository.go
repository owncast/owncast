package outboxrepository

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
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

// OutboxRepository handles persistence of ActivityPub outbox items.
type OutboxRepository interface {
	// GetPostCount returns the number of posts in the outbox.
	GetPostCount() (int64, error)
	// GetOutbox returns an ordered collection of outbox items.
	GetOutbox(limit, offset int) (vocab.ActivityStreamsOrderedCollection, error)
	// Add stores a single payload to the outbox.
	Add(iri string, itemData []byte, typeString string, isLiveNotification bool) error
	// GetObjectByIRI returns a string representation of a single object by IRI.
	GetObjectByIRI(iri string) (string, bool, time.Time, error)
	// GetLocalPostCount returns the number of posts existing locally.
	GetLocalPostCount() (int64, error)
}

// SqlOutboxRepository is the SQL-based implementation of OutboxRepository.
type SqlOutboxRepository struct {
	datastore *data.Datastore
}

// NOTE: This is temporary during the transition period.
var temporaryGlobalInstance OutboxRepository

// Get returns the outbox repository singleton.
func Get() OutboxRepository {
	if temporaryGlobalInstance == nil {
		i := New(data.GetDatastore())
		temporaryGlobalInstance = i
	}
	return temporaryGlobalInstance
}

// New creates a new instance of the OutboxRepository.
func New(datastore *data.Datastore) OutboxRepository {
	r := SqlOutboxRepository{
		datastore: datastore,
	}
	return &r
}

// GetPostCount returns the number of posts in the outbox.
func (r *SqlOutboxRepository) GetPostCount() (int64, error) {
	ctx := context.Background()
	return r.datastore.GetQueries().GetLocalPostCount(ctx)
}

// GetOutbox returns an ordered collection of outbox items.
func (r *SqlOutboxRepository) GetOutbox(limit int, offset int) (vocab.ActivityStreamsOrderedCollection, error) {
	collection := streams.NewActivityStreamsOrderedCollection()
	orderedItems := streams.NewActivityStreamsOrderedItemsProperty()
	rows, err := r.datastore.GetQueries().GetOutboxWithOffset(
		context.Background(),
		db.GetOutboxWithOffsetParams{Limit: utils.SafeIntToInt32(limit), Offset: utils.SafeIntToInt32(offset)},
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

// Add stores a single payload to the outbox.
func (r *SqlOutboxRepository) Add(iri string, itemData []byte, typeString string, isLiveNotification bool) error {
	tx, err := r.datastore.DB.Begin()
	if err != nil {
		log.Debugln(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err = r.datastore.GetQueries().WithTx(tx).AddToOutbox(context.Background(), db.AddToOutboxParams{
		Iri:              iri,
		Value:            itemData,
		Type:             typeString,
		LiveNotification: sql.NullBool{Bool: isLiveNotification, Valid: true},
	}); err != nil {
		return fmt.Errorf("error creating new item in federation outbox %s", err)
	}

	return tx.Commit()
}

// GetObjectByIRI returns a string representation of a single object by IRI.
func (r *SqlOutboxRepository) GetObjectByIRI(iri string) (string, bool, time.Time, error) {
	row, err := r.datastore.GetQueries().GetObjectFromOutboxByIRI(context.Background(), iri)
	return string(row.Value), row.LiveNotification.Bool, row.CreatedAt.Time, err
}

// GetLocalPostCount returns the number of posts existing locally.
func (r *SqlOutboxRepository) GetLocalPostCount() (int64, error) {
	ctx := context.Background()
	return r.datastore.GetQueries().GetLocalPostCount(ctx)
}
