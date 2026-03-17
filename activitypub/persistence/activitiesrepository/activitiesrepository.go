package activitiesrepository

import (
	"context"
	"time"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
	"github.com/pkg/errors"
)

// ActivitiesRepository handles persistence of inbound federated activities.
type ActivitiesRepository interface {
	// Save stores an inbound federated activity.
	Save(objectIRI, actorIRI, eventType string, timestamp time.Time) error
	// GetInboundActivities returns a paginated list of inbound activities.
	GetInboundActivities(limit, offset int) ([]models.FederatedActivity, int, error)
	// HasPreviouslyHandled returns whether an inbound activity has already been handled.
	HasPreviouslyHandled(iri, actorIRI, eventType string) (bool, error)
}

// SqlActivitiesRepository is the SQL-based implementation of ActivitiesRepository.
type SqlActivitiesRepository struct {
	datastore *data.Datastore
}

// NOTE: This is temporary during the transition period.
var temporaryGlobalInstance ActivitiesRepository

// Get returns the activities repository singleton.
func Get() ActivitiesRepository {
	if temporaryGlobalInstance == nil {
		i := New(data.GetDatastore())
		temporaryGlobalInstance = i
	}
	return temporaryGlobalInstance
}

// New creates a new instance of the ActivitiesRepository.
func New(datastore *data.Datastore) ActivitiesRepository {
	r := SqlActivitiesRepository{
		datastore: datastore,
	}
	return &r
}

// Save stores an inbound federated activity.
func (r *SqlActivitiesRepository) Save(objectIRI string, actorIRI string, eventType string, timestamp time.Time) error {
	if err := r.datastore.GetQueries().AddToAcceptedActivities(context.Background(), db.AddToAcceptedActivitiesParams{
		Iri:       objectIRI,
		Actor:     actorIRI,
		Type:      eventType,
		Timestamp: timestamp,
	}); err != nil {
		return errors.Wrap(err, "error saving event "+objectIRI)
	}

	return nil
}

// GetInboundActivities returns a paginated list of inbound activities.
func (r *SqlActivitiesRepository) GetInboundActivities(limit int, offset int) ([]models.FederatedActivity, int, error) {
	ctx := context.Background()
	rows, err := r.datastore.GetQueries().GetInboundActivitiesWithOffset(ctx, db.GetInboundActivitiesWithOffsetParams{
		Limit:  utils.SafeIntToInt32(limit),
		Offset: utils.SafeIntToInt32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	activities := make([]models.FederatedActivity, 0)

	total, err := r.datastore.GetQueries().GetInboundActivityCount(context.Background())
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

// HasPreviouslyHandled returns whether an inbound activity has already been handled.
func (r *SqlActivitiesRepository) HasPreviouslyHandled(iri string, actorIRI string, eventType string) (bool, error) {
	exists, err := r.datastore.GetQueries().DoesInboundActivityExist(context.Background(), db.DoesInboundActivityExistParams{
		Iri:   iri,
		Actor: actorIRI,
		Type:  eventType,
	})
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}
