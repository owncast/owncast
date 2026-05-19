package followersrepository

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/url"
	"time"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/utils"
)

// FollowersRepository handles persistence of ActivityPub followers.
type FollowersRepository interface {
	// GetCount returns the number of followers.
	GetCount() (int64, error)
	// GetFollowers returns a paginated list of followers.
	GetFollowers(limit int, offset int) ([]models.Follower, int, error)
	// GetPendingFollowRequests returns pending follow requests.
	GetPendingFollowRequests() ([]models.Follower, error)
	// GetBlockedAndRejected returns blocked and rejected followers.
	GetBlockedAndRejected() ([]models.Follower, error)
	// GetUniqueDeliveryInboxes returns unique inbox URLs for delivery.
	GetUniqueDeliveryInboxes() ([]string, error)
	// GetByIRI returns a single follower by IRI.
	GetByIRI(iri string) (*apmodels.ActivityPubActor, error)
	// Add saves a new follow to the datastore.
	Add(follow apmodels.ActivityPubActor, approved bool) error
	// Remove removes a follow from the datastore.
	Remove(unfollow apmodels.ActivityPubActor) error
	// ApprovePreviousRequest approves a pending follow request.
	ApprovePreviousRequest(iri string) error
	// BlockOrReject blocks an existing follower or rejects a follow request.
	BlockOrReject(iri string) error
	// Update updates the details of a stored follower.
	Update(actorIRI string, inbox string, sharedInbox string, name string, username string, image string) error
	// GetFollowersToValidate returns followers needing validation, ordered by oldest validated first.
	GetFollowersToValidate(limit int) ([]models.Follower, error)
	// UpdateFollowerValidationSuccess marks a follower as successfully validated and clears failure timestamp.
	UpdateFollowerValidationSuccess(iri string) error
	// UpdateFollowerValidationFailure marks a validation failure, setting first failure time if not already set.
	UpdateFollowerValidationFailure(iri string) error
	// RemoveByIRI removes a follower directly by IRI string.
	RemoveByIRI(iri string) error
}

// SqlFollowersRepository is the SQL-based implementation of FollowersRepository.
type SqlFollowersRepository struct {
	datastore *datastore.Datastore
}

// New creates a new instance of the FollowersRepository.
func New(datastore *datastore.Datastore) FollowersRepository {
	r := SqlFollowersRepository{
		datastore: datastore,
	}
	return &r
}

// GetCount returns the number of followers.
func (r *SqlFollowersRepository) GetCount() (int64, error) {
	ctx := context.Background()
	return r.datastore.GetQueries().GetFollowerCount(ctx)
}

// GetFollowers returns a paginated list of followers.
func (r *SqlFollowersRepository) GetFollowers(limit int, offset int) ([]models.Follower, int, error) {
	ctx := context.Background()
	total, err := r.datastore.GetQueries().GetFollowerCount(ctx)
	if err != nil {
		return nil, 0, errors.Wrap(err, "unable to fetch total number of followers")
	}

	followersResult, err := r.datastore.GetQueries().GetFederationFollowersWithOffset(ctx, db.GetFederationFollowersWithOffsetParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	followers := make([]models.Follower, 0)

	for _, row := range followersResult {
		singleFollower := models.Follower{
			Name:        row.Name.String,
			Username:    row.Username,
			Image:       row.Image.String,
			ActorIRI:    row.Iri,
			Inbox:       row.Inbox,
			SharedInbox: row.SharedInbox.String,
			Timestamp:   utils.NullTime(row.CreatedAt),
		}

		followers = append(followers, singleFollower)
	}

	return followers, int(total), nil
}

// GetPendingFollowRequests returns pending follow requests.
func (r *SqlFollowersRepository) GetPendingFollowRequests() ([]models.Follower, error) {
	pendingFollowersResult, err := r.datastore.GetQueries().GetFederationFollowerApprovalRequests(context.Background())
	if err != nil {
		return nil, err
	}

	followers := make([]models.Follower, 0)

	for _, row := range pendingFollowersResult {
		singleFollower := models.Follower{
			Name:        row.Name.String,
			Username:    row.Username,
			Image:       row.Image.String,
			ActorIRI:    row.Iri,
			Inbox:       row.Inbox,
			SharedInbox: row.SharedInbox.String,
			Timestamp:   utils.NullTime{Time: row.CreatedAt.Time, Valid: true},
		}
		followers = append(followers, singleFollower)
	}

	return followers, nil
}

// GetBlockedAndRejected returns blocked and rejected followers.
func (r *SqlFollowersRepository) GetBlockedAndRejected() ([]models.Follower, error) {
	pendingFollowersResult, err := r.datastore.GetQueries().GetRejectedAndBlockedFollowers(context.Background())
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

// GetUniqueDeliveryInboxes returns unique inbox URLs for delivery.
func (r *SqlFollowersRepository) GetUniqueDeliveryInboxes() ([]string, error) {
	ctx := context.Background()
	return r.datastore.GetQueries().GetUniqueDeliveryInboxes(ctx)
}

// GetByIRI returns a single follower by IRI.
func (r *SqlFollowersRepository) GetByIRI(iri string) (*apmodels.ActivityPubActor, error) {
	result, err := r.datastore.GetQueries().GetFollowerByIRI(context.Background(), iri)
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

	var sharedInbox *url.URL
	if result.SharedInbox.Valid && result.SharedInbox.String != "" {
		sharedInbox, err = url.Parse(result.SharedInbox.String)
		if err != nil {
			log.Warnln("error parsing shared inbox, ignoring:", err)
		}
	}

	requestObjectBytes := result.RequestObject
	var followRequestObject vocab.ActivityStreamsFollow

	resolver, err := streams.NewJSONResolver(func(c context.Context, followObject vocab.ActivityStreamsFollow) error {
		followRequestObject = followObject
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating JSON resolver")
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(requestObjectBytes, &jsonMap)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling follow request object")
	}

	err = resolver.Resolve(context.Background(), jsonMap)
	if err != nil {
		return nil, errors.Wrap(err, "error resolving follow request object")
	}

	image, _ := url.Parse(result.Image.String)

	var disabledAt *time.Time
	if result.DisabledAt.Valid {
		disabledAt = &result.DisabledAt.Time
	}

	follower := apmodels.ActivityPubActor{
		ActorIri:         iriURL,
		Inbox:            inbox,
		SharedInbox:      sharedInbox,
		Name:             result.Name.String,
		Username:         result.Username,
		Image:            image,
		FollowRequestIri: followIRI,
		DisabledAt:       disabledAt,
		RequestObject:    followRequestObject,
	}

	return &follower, nil
}

// Add saves a new follow to the datastore.
func (r *SqlFollowersRepository) Add(follow apmodels.ActivityPubActor, approved bool) error {
	if err := follow.Validate(); err != nil {
		return errors.Wrap(err, "cannot add invalid follow")
	}

	log.Traceln("Saving", follow.ActorIriString(), "as a follower.")

	followRequestObject, err := apmodels.Serialize(follow.RequestObject)
	if err != nil {
		return errors.Wrap(err, "error serializing follow request object")
	}

	return r.createFollow(follow.ActorIriString(), follow.InboxString(), follow.SharedInboxString(), follow.FollowRequestIriString(), follow.Name, follow.Username, follow.ImageString(), followRequestObject, approved)
}

// Remove removes a follow from the datastore.
func (r *SqlFollowersRepository) Remove(unfollow apmodels.ActivityPubActor) error {
	if err := unfollow.Validate(); err != nil {
		return errors.Wrap(err, "cannot remove invalid follow")
	}
	log.Traceln("Removing", unfollow.ActorIriString(), "as a follower.")
	return r.removeFollow(unfollow.ActorIri)
}

// ApprovePreviousRequest approves a pending follow request.
func (r *SqlFollowersRepository) ApprovePreviousRequest(iri string) error {
	return r.datastore.GetQueries().ApproveFederationFollower(context.Background(), db.ApproveFederationFollowerParams{
		Iri: iri,
		ApprovedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
}

// BlockOrReject blocks an existing follower or rejects a follow request.
func (r *SqlFollowersRepository) BlockOrReject(iri string) error {
	return r.datastore.GetQueries().RejectFederationFollower(context.Background(), db.RejectFederationFollowerParams{
		Iri: iri,
		DisabledAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
}

// Update updates the details of a stored follower.
func (r *SqlFollowersRepository) Update(actorIRI string, inbox string, sharedInbox string, name string, username string, image string) error {
	r.datastore.DbLock.Lock()
	defer r.datastore.DbLock.Unlock()

	tx, err := r.datastore.DB.Begin()
	if err != nil {
		return errors.Wrap(err, "error beginning transaction")
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err = r.datastore.GetQueries().WithTx(tx).UpdateFollowerByIRI(context.Background(), db.UpdateFollowerByIRIParams{
		Inbox:       inbox,
		SharedInbox: sql.NullString{String: sharedInbox, Valid: sharedInbox != ""},
		Name:        sql.NullString{String: name, Valid: true},
		Username:    username,
		Image:       sql.NullString{String: image, Valid: true},
		Iri:         actorIRI,
	}); err != nil {
		return errors.Wrap(err, "error updating follower "+actorIRI)
	}

	return tx.Commit()
}

func (r *SqlFollowersRepository) createFollow(actor, inbox, sharedInbox, request, name, username, image string, requestObject []byte, approved bool) error {
	r.datastore.DbLock.Lock()
	defer r.datastore.DbLock.Unlock()

	tx, err := r.datastore.DB.Begin()
	if err != nil {
		return errors.Wrap(err, "error beginning transaction")
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

	if err = r.datastore.GetQueries().WithTx(tx).AddFollower(context.Background(), db.AddFollowerParams{
		Iri:           actor,
		Inbox:         inbox,
		SharedInbox:   sql.NullString{String: sharedInbox, Valid: sharedInbox != ""},
		Name:          sql.NullString{String: name, Valid: true},
		Username:      username,
		Image:         sql.NullString{String: image, Valid: true},
		ApprovedAt:    approvedAt,
		Request:       request,
		RequestObject: requestObject,
	}); err != nil {
		return errors.Wrap(err, "error creating new federation follow")
	}

	return tx.Commit()
}

func (r *SqlFollowersRepository) removeFollow(actor *url.URL) error {
	r.datastore.DbLock.Lock()
	defer r.datastore.DbLock.Unlock()

	tx, err := r.datastore.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := r.datastore.GetQueries().WithTx(tx).RemoveFollowerByIRI(context.Background(), actor.String()); err != nil {
		return err
	}

	return tx.Commit()
}

// GetFollowersToValidate returns followers needing validation, ordered by oldest validated first.
func (r *SqlFollowersRepository) GetFollowersToValidate(limit int) ([]models.Follower, error) {
	ctx := context.Background()
	followersResult, err := r.datastore.GetQueries().GetFollowersToValidate(ctx, int64(limit))
	if err != nil {
		return nil, err
	}

	followers := make([]models.Follower, 0)

	for _, row := range followersResult {
		singleFollower := models.Follower{
			Name:                     row.Name.String,
			Username:                 row.Username,
			Image:                    row.Image.String,
			ActorIRI:                 row.Iri,
			Inbox:                    row.Inbox,
			SharedInbox:              row.SharedInbox.String,
			FirstValidationFailureAt: utils.NullTime(row.FirstValidationFailureAt),
		}

		followers = append(followers, singleFollower)
	}

	return followers, nil
}

// UpdateFollowerValidationSuccess marks a follower as successfully validated and clears failure timestamp.
func (r *SqlFollowersRepository) UpdateFollowerValidationSuccess(iri string) error {
	r.datastore.DbLock.Lock()
	defer r.datastore.DbLock.Unlock()

	ctx := context.Background()
	tx, err := r.datastore.DB.Begin()
	if err != nil {
		return errors.Wrap(err, "error beginning transaction")
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := r.datastore.GetQueries().WithTx(tx).UpdateFollowerValidationSuccess(ctx, db.UpdateFollowerValidationSuccessParams{
		LastValidatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Iri:             iri,
	}); err != nil {
		return errors.Wrap(err, "error updating follower validation success")
	}

	return tx.Commit()
}

// UpdateFollowerValidationFailure marks a validation failure, setting first failure time if not already set.
func (r *SqlFollowersRepository) UpdateFollowerValidationFailure(iri string) error {
	r.datastore.DbLock.Lock()
	defer r.datastore.DbLock.Unlock()

	ctx := context.Background()
	tx, err := r.datastore.DB.Begin()
	if err != nil {
		return errors.Wrap(err, "error beginning transaction")
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := r.datastore.GetQueries().WithTx(tx).UpdateFollowerValidationFailure(ctx, db.UpdateFollowerValidationFailureParams{
		LastValidatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Iri:             iri,
	}); err != nil {
		return errors.Wrap(err, "error updating follower validation failure")
	}

	return tx.Commit()
}

// RemoveByIRI removes a follower directly by IRI string.
func (r *SqlFollowersRepository) RemoveByIRI(iri string) error {
	r.datastore.DbLock.Lock()
	defer r.datastore.DbLock.Unlock()

	tx, err := r.datastore.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := r.datastore.GetQueries().WithTx(tx).RemoveFollowerByIRI(context.Background(), iri); err != nil {
		return err
	}

	return tx.Commit()
}
