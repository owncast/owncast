package federationrepository

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/apfederation/apmodels"
	"github.com/owncast/owncast/services/apfederation/resolvers"
	"github.com/owncast/owncast/storage/data"
	"github.com/owncast/owncast/storage/sqlstorage"
	"github.com/owncast/owncast/utils"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

type FederationRepository struct {
	datastore *data.Store
}

func New(datastore *data.Store) *FederationRepository {
	r := &FederationRepository{
		datastore: datastore,
	}

	return r
}

// NOTE: This is temporary during the transition period.
var temporaryGlobalInstance *FederationRepository

// GetUserRepository will return the user repository.
func Get() *FederationRepository {
	if temporaryGlobalInstance == nil {
		i := New(data.GetDatastore())
		temporaryGlobalInstance = i
	}
	return temporaryGlobalInstance
}

// GetFollowerCount will return the number of followers we're keeping track of.
func (f *FederationRepository) GetFollowerCount() (int64, error) {
	ctx := context.Background()
	return f.datastore.GetQueries().GetFollowerCount(ctx)
}

// GetFederationFollowers will return a slice of the followers we keep track of locally.
func (f *FederationRepository) GetFederationFollowers(limit int, offset int) ([]models.Follower, int, error) {
	ctx := context.Background()
	total, err := f.datastore.GetQueries().GetFollowerCount(ctx)
	if err != nil {
		return nil, 0, errors.Wrap(err, "unable to fetch total number of followers")
	}

	followersResult, err := f.datastore.GetQueries().GetFederationFollowersWithOffset(ctx, sqlstorage.GetFederationFollowersWithOffsetParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	followers := make([]models.Follower, 0)

	for _, row := range followersResult {
		singleFollower := models.Follower{
			Name:      row.Name.String,
			Username:  row.Username,
			Image:     row.Image.String,
			ActorIRI:  row.Iri,
			Inbox:     row.Inbox,
			Timestamp: utils.NullTime(row.CreatedAt),
		}

		followers = append(followers, singleFollower)
	}

	return followers, int(total), nil
}

// GetPendingFollowRequests will return pending follow requests.
func (f *FederationRepository) GetPendingFollowRequests() ([]models.Follower, error) {
	pendingFollowersResult, err := f.datastore.GetQueries().GetFederationFollowerApprovalRequests(context.Background())
	if err != nil {
		return nil, err
	}

	followers := make([]models.Follower, 0)

	for _, row := range pendingFollowersResult {
		singleFollower := models.Follower{
			Name:      row.Name.String,
			Username:  row.Username,
			Image:     row.Image.String,
			ActorIRI:  row.Iri,
			Inbox:     row.Inbox,
			Timestamp: utils.NullTime{Time: row.CreatedAt.Time, Valid: true},
		}
		followers = append(followers, singleFollower)
	}

	return followers, nil
}

// GetBlockedAndRejectedFollowers will return blocked and rejected followers.
func (f *FederationRepository) GetBlockedAndRejectedFollowers() ([]models.Follower, error) {
	pendingFollowersResult, err := f.datastore.GetQueries().GetRejectedAndBlockedFollowers(context.Background())
	if err != nil {
		return nil, err
	}

	followers := make([]models.Follower, 0)

	for _, row := range pendingFollowersResult {
		singleFollower := models.Follower{
			Name:       row.Name.String,
			Username:   row.Username,
			Image:      row.Image.String,
			ActorIRI:   row.Iri,
			DisabledAt: utils.NullTime{Time: row.DisabledAt.Time, Valid: true},
			Timestamp:  utils.NullTime{Time: row.CreatedAt.Time, Valid: true},
		}
		followers = append(followers, singleFollower)
	}

	return followers, nil
}

// AddFollow will save a follow to the datastore.
func (f *FederationRepository) AddFollow(follow apmodels.ActivityPubActor, approved bool) error {
	log.Traceln("Saving", follow.ActorIri, "as a follower.")
	var image string
	if follow.Image != nil {
		image = follow.Image.String()
	}

	followRequestObject, err := apmodels.Serialize(follow.RequestObject)
	if err != nil {
		return errors.Wrap(err, "error serializing follow request object")
	}

	return f.createFollow(follow.ActorIri.String(), follow.Inbox.String(), follow.FollowRequestIri.String(), follow.Name, follow.Username, image, followRequestObject, approved)
}

// RemoveFollow will remove a follow from the datastore.
func (f *FederationRepository) RemoveFollow(unfollow apmodels.ActivityPubActor) error {
	log.Traceln("Removing", unfollow.ActorIri, "as a follower.")
	return f.removeFollow(unfollow.ActorIri)
}

// GetFollower will return a single follower/request given an IRI.
func (f *FederationRepository) GetFollower(iri string) (*apmodels.ActivityPubActor, error) {
	result, err := f.datastore.GetQueries().GetFollowerByIRI(context.Background(), iri)
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
func (f *FederationRepository) ApprovePreviousFollowRequest(iri string) error {
	return f.datastore.GetQueries().ApproveFederationFollower(context.Background(), sqlstorage.ApproveFederationFollowerParams{
		Iri: iri,
		ApprovedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
}

// BlockOrRejectFollower will block an existing follower or reject a follow request.
func (f *FederationRepository) BlockOrRejectFollower(iri string) error {
	return f.datastore.GetQueries().RejectFederationFollower(context.Background(), sqlstorage.RejectFederationFollowerParams{
		Iri: iri,
		DisabledAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
}

func (f *FederationRepository) createFollow(actor, inbox, request, name, username, image string, requestObject []byte, approved bool) error {
	tx, err := f.datastore.DB.Begin()
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

	if err = f.datastore.GetQueries().WithTx(tx).AddFollower(context.Background(), sqlstorage.AddFollowerParams{
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
func (f *FederationRepository) UpdateFollower(actorIRI string, inbox string, name string, username string, image string) error {
	f.datastore.DbLock.Lock()
	defer f.datastore.DbLock.Unlock()

	tx, err := f.datastore.DB.Begin()
	if err != nil {
		log.Debugln(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err = f.datastore.GetQueries().WithTx(tx).UpdateFollowerByIRI(context.Background(), sqlstorage.UpdateFollowerByIRIParams{
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

func (f *FederationRepository) removeFollow(actor *url.URL) error {
	f.datastore.DbLock.Lock()
	defer f.datastore.DbLock.Unlock()

	tx, err := f.datastore.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := f.datastore.GetQueries().WithTx(tx).RemoveFollowerByIRI(context.Background(), actor.String()); err != nil {
		return err
	}

	return tx.Commit()
}

// GetOutboxPostCount will return the number of posts in the outbox.
func (f *FederationRepository) GetOutboxPostCount() (int64, error) {
	ctx := context.Background()
	return f.datastore.GetQueries().GetLocalPostCount(ctx)
}

// GetOutbox will return an instance of the outbox populated by stored items.
func (f *FederationRepository) GetOutbox(limit int, offset int) (vocab.ActivityStreamsOrderedCollection, error) {
	collection := streams.NewActivityStreamsOrderedCollection()
	orderedItems := streams.NewActivityStreamsOrderedItemsProperty()
	r := resolvers.Get()

	rows, err := f.datastore.GetQueries().GetOutboxWithOffset(
		context.Background(),
		sqlstorage.GetOutboxWithOffsetParams{Limit: int32(limit), Offset: int32(offset)},
	)
	if err != nil {
		return collection, err
	}

	for _, value := range rows {
		createCallback := func(c context.Context, activity vocab.ActivityStreamsCreate) error {
			orderedItems.AppendActivityStreamsCreate(activity)
			return nil
		}
		if err := r.Resolve(context.Background(), value, createCallback); err != nil {
			return collection, err
		}
	}

	return collection, nil
}

// AddToOutbox will store a single payload to the persistence layer.
func (f *FederationRepository) AddToOutbox(iri string, itemData []byte, typeString string, isLiveNotification bool) error {
	tx, err := f.datastore.DB.Begin()
	if err != nil {
		log.Debugln(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err = f.datastore.GetQueries().WithTx(tx).AddToOutbox(context.Background(), sqlstorage.AddToOutboxParams{
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
func (f *FederationRepository) GetObjectByIRI(iri string) (string, bool, time.Time, error) {
	row, err := f.datastore.GetQueries().GetObjectFromOutboxByIRI(context.Background(), iri)
	return string(row.Value), row.LiveNotification.Bool, row.CreatedAt.Time, err
}

// GetLocalPostCount will return the number of posts existing locally.
func (f *FederationRepository) GetLocalPostCount() (int64, error) {
	ctx := context.Background()
	return f.datastore.GetQueries().GetLocalPostCount(ctx)
}

// SaveInboundFediverseActivity will save an event to the ap_inbound_activities table.
func (f *FederationRepository) SaveInboundFediverseActivity(objectIRI string, actorIRI string, eventType string, timestamp time.Time) error {
	if err := f.datastore.GetQueries().AddToAcceptedActivities(context.Background(), sqlstorage.AddToAcceptedActivitiesParams{
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
func (f *FederationRepository) GetInboundActivities(limit int, offset int) ([]models.FederatedActivity, int, error) {
	ctx := context.Background()
	rows, err := f.datastore.GetQueries().GetInboundActivitiesWithOffset(ctx, sqlstorage.GetInboundActivitiesWithOffsetParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	activities := make([]models.FederatedActivity, 0)

	total, err := f.datastore.GetQueries().GetInboundActivityCount(context.Background())
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
func (f *FederationRepository) HasPreviouslyHandledInboundActivity(iri string, actorIRI string, eventType string) (bool, error) {
	exists, err := f.datastore.GetQueries().DoesInboundActivityExist(context.Background(), sqlstorage.DoesInboundActivityExistParams{
		Iri:   iri,
		Actor: actorIRI,
		Type:  eventType,
	})
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}
