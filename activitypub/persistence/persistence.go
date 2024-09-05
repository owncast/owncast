package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/resolvers"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/models"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

var _datastore *data.Datastore

// Setup will initialize the ActivityPub persistence layer with the provided datastore.
func Setup(datastore *data.Datastore) {
	_datastore = datastore
	createFederationFollowersTable()
	createFederationOutboxTable()
	createFederatedActivitiesTable()
	addFollowersFixtureData()
}

// AddFollow will save a follow to the datastore.
func AddFollow(follow apmodels.ActivityPubActor, approved bool) error {
	log.Traceln("Saving", follow.ActorIri, "as a follower.")
	var image string
	if follow.Image != nil {
		image = follow.Image.String()
	}

	followRequestObject, err := apmodels.Serialize(follow.RequestObject)
	if err != nil {
		return errors.Wrap(err, "error serializing follow request object")
	}

	return createFollow(follow.ActorIri.String(), follow.Inbox.String(), follow.FollowRequestIri.String(), follow.Name, follow.Username, image, followRequestObject, approved)
}

// RemoveFollow will remove a follow from the datastore.
func RemoveFollow(unfollow apmodels.ActivityPubActor) error {
	log.Traceln("Removing", unfollow.ActorIri, "as a follower.")
	return removeFollow(unfollow.ActorIri)
}

// GetFollower will return a single follower/request given an IRI.
func GetFollower(iri string) (*apmodels.ActivityPubActor, error) {
	result, err := _datastore.GetQueries().GetFollowerByIRI(context.Background(), iri)
	if err != nil {
		return nil, err
	}

	followIRI, err := url.Parse(result.Request)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing follow request IRI")
	}

	iriURL, err := url.Parse(result.Iri)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing actor IRI")
	}

	inbox, err := url.Parse(result.Inbox)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing acting inbox")
	}

	image, _ := url.Parse(result.Image.String)

	var disabledAt *time.Time
	if result.DisabledAt.Valid {
		disabledAt = &result.DisabledAt.Time
	}

	follower := apmodels.ActivityPubActor{
		ActorIri:         iriURL,
		Inbox:            inbox,
		Name:             result.Name.String,
		Username:         result.Username,
		Image:            image,
		FollowRequestIri: followIRI,
		DisabledAt:       disabledAt,
	}

	return &follower, nil
}

// ApprovePreviousFollowRequest will approve a follow request.
func ApprovePreviousFollowRequest(iri string) error {
	return _datastore.GetQueries().ApproveFederationFollower(context.Background(), db.ApproveFederationFollowerParams{
		Iri: iri,
		ApprovedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
}

// BlockOrRejectFollower will block an existing follower or reject a follow request.
func BlockOrRejectFollower(iri string) error {
	return _datastore.GetQueries().RejectFederationFollower(context.Background(), db.RejectFederationFollowerParams{
		Iri: iri,
		DisabledAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
}

func createFollow(actor, inbox, request, name, username, image string, requestObject []byte, approved bool) error {
	tx, err := _datastore.DB.Begin()
	if err != nil {
		log.Debugln(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var approvedAt sql.NullTime
	if approved {
		approvedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	if err = _datastore.GetQueries().WithTx(tx).AddFollower(context.Background(), db.AddFollowerParams{
		Iri:           actor,
		Inbox:         inbox,
		Name:          sql.NullString{String: name, Valid: true},
		Username:      username,
		Image:         sql.NullString{String: image, Valid: true},
		ApprovedAt:    approvedAt,
		Request:       request,
		RequestObject: requestObject,
	}); err != nil {
		log.Errorln("error creating new federation follow: ", err)
	}

	return tx.Commit()
}

// UpdateFollower will update the details of a stored follower given an IRI.
func UpdateFollower(actorIRI string, inbox string, name string, username string, image string) error {
	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	tx, err := _datastore.DB.Begin()
	if err != nil {
		log.Debugln(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err = _datastore.GetQueries().WithTx(tx).UpdateFollowerByIRI(context.Background(), db.UpdateFollowerByIRIParams{
		Inbox:    inbox,
		Name:     sql.NullString{String: name, Valid: true},
		Username: username,
		Image:    sql.NullString{String: image, Valid: true},
		Iri:      actorIRI,
	}); err != nil {
		return fmt.Errorf("error updating follower %s %s", actorIRI, err)
	}

	return tx.Commit()
}

func removeFollow(actor *url.URL) error {
	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	tx, err := _datastore.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := _datastore.GetQueries().WithTx(tx).RemoveFollowerByIRI(context.Background(), actor.String()); err != nil {
		return err
	}

	return tx.Commit()
}

// createFederatedActivitiesTable will create the accepted
// activities table if needed.
func createFederatedActivitiesTable() {
	createTableSQL := `CREATE TABLE IF NOT EXISTS ap_accepted_activities (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"iri" TEXT NOT NULL,
    "actor" TEXT NOT NULL,
    "type" TEXT NOT NULL,
		"timestamp" TIMESTAMP NOT NULL
	);`

	_datastore.MustExec(createTableSQL)
	_datastore.MustExec(`CREATE INDEX IF NOT EXISTS idx_iri_actor_index ON ap_accepted_activities (iri,actor);`)
}

func createFederationOutboxTable() {
	log.Traceln("Creating federation outbox table...")
	createTableSQL := `CREATE TABLE IF NOT EXISTS ap_outbox (
		"iri" TEXT NOT NULL,
		"value" BLOB,
		"type" TEXT NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "live_notification" BOOLEAN DEFAULT FALSE,
		PRIMARY KEY (iri));`

	_datastore.MustExec(createTableSQL)
	_datastore.MustExec(`CREATE INDEX IF NOT EXISTS idx_iri ON ap_outbox (iri);`)
	_datastore.MustExec(`CREATE INDEX IF NOT EXISTS idx_type ON ap_outbox (type);`)
	_datastore.MustExec(`CREATE INDEX IF NOT EXISTS idx_live_notification ON ap_outbox (live_notification);`)
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
		db.GetOutboxWithOffsetParams{Limit: limit, Offset: offset},
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
		Limit:  limit,
		Offset: offset,
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
