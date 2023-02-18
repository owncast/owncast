package follower

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/internal/activitypub/apmodels"
	"github.com/owncast/owncast/internal/activitypub/crypto"
	"github.com/owncast/owncast/internal/activitypub/workerpool"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
)

// New will initialize the ActivityPub persistence layer service
// with the provided Data.
func New(
	ds *data.Service,
	c *crypto.Service,
	m *apmodels.Service,
	w *workerpool.Service,
) (*Service, error) {
	if ds == nil {
		return nil, errors.New("Data store is nil")
	}

	s := &Service{}

	s.Data = ds
	s.crypto = c
	s.models = m
	s.worker = w

	s.createFederationFollowersTable()

	return s, nil
}

type Service struct {
	Data   *data.Service
	crypto *crypto.Service
	models *apmodels.Service
	worker *workerpool.Service
}

func (s *Service) createFederationFollowersTable() {
	const (
		sqlCreateTable = `CREATE TABLE IF NOT EXISTS ap_followers (
		"iri" TEXT NOT NULL,
		"inbox" TEXT NOT NULL,
		"name" TEXT,
		"username" TEXT NOT NULL,
		"image" TEXT,
    "request" TEXT NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"approved_at" TIMESTAMP,
    "disabled_at" TIMESTAMP,
    "request_object" BLOB,
		PRIMARY KEY (iri));`
	)

	log.Traceln("Creating federation followers table...")
	s.Data.Store.MustExec(sqlCreateTable)
	s.Data.Store.MustExec(`CREATE INDEX IF NOT EXISTS idx_iri ON ap_followers (iri);`)
	s.Data.Store.MustExec(`CREATE INDEX IF NOT EXISTS idx_approved_at ON ap_followers (approved_at);`)
}

// GetFollowerCount will return the number of followers we're keeping track of.
func (s *Service) GetFollowerCount() (int64, error) {
	ctx := context.Background()
	return s.Data.Store.GetQueries().GetFollowerCount(ctx)
}

// GetFederationFollowers will return a slice of the followers we keep track of locally.
func (s *Service) GetFederationFollowers(limit int, offset int) ([]models.Follower, int, error) {
	ctx := context.Background()
	total, err := s.Data.Store.GetQueries().GetFollowerCount(ctx)
	if err != nil {
		return nil, 0, errors.Wrap(err, "unable to fetch total number of followers")
	}

	followersResult, err := s.Data.Store.GetQueries().GetFederationFollowersWithOffset(ctx, db.GetFederationFollowersWithOffsetParams{
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
func (s *Service) GetPendingFollowRequests() ([]models.Follower, error) {
	pendingFollowersResult, err := s.Data.Store.GetQueries().GetFederationFollowerApprovalRequests(context.Background())
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
func (s *Service) GetBlockedAndRejectedFollowers() ([]models.Follower, error) {
	pendingFollowersResult, err := s.Data.Store.GetQueries().GetRejectedAndBlockedFollowers(context.Background())
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

// AddFollow will save a follow to the Data.
func (s *Service) AddFollow(follow apmodels.ActivityPubActor, approved bool) error {
	log.Traceln("Saving", follow.ActorIri, "as a follower.")
	var image string
	if follow.Image != nil {
		image = follow.Image.String()
	}

	followRequestObject, err := s.models.Serialize(follow.RequestObject)
	if err != nil {
		return errors.Wrap(err, "error serializing follow request object")
	}

	return s.createFollow(follow.ActorIri.String(), follow.Inbox.String(), follow.FollowRequestIri.String(), follow.Name, follow.Username, image, followRequestObject, approved)
}

// RemoveFollow will remove a follow from the Data.
func (s *Service) RemoveFollow(unfollow apmodels.ActivityPubActor) error {
	log.Traceln("Removing", unfollow.ActorIri, "as a follower.")
	return s.removeFollow(unfollow.ActorIri)
}

// GetFollower will return a single follower/request given an IRI.
func (s *Service) GetFollower(iri string) (*apmodels.ActivityPubActor, error) {
	result, err := s.Data.Store.GetQueries().GetFollowerByIRI(context.Background(), iri)
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
func (s *Service) ApprovePreviousFollowRequest(iri string) error {
	return s.Data.Store.GetQueries().ApproveFederationFollower(context.Background(), db.ApproveFederationFollowerParams{
		Iri: iri,
		ApprovedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
}

// BlockOrRejectFollower will block an existing follower or reject a follow request.
func (s *Service) BlockOrRejectFollower(iri string) error {
	return s.Data.Store.GetQueries().RejectFederationFollower(context.Background(), db.RejectFederationFollowerParams{
		Iri: iri,
		DisabledAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
}

func (s *Service) createFollow(actor, inbox, request, name, username, image string, requestObject []byte, approved bool) error {
	tx, err := s.Data.DB.Begin()
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

	if err = s.Data.Store.GetQueries().WithTx(tx).AddFollower(context.Background(), db.AddFollowerParams{
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
func (s *Service) UpdateFollower(actorIRI string, inbox string, name string, username string, image string) error {
	s.Data.Store.DbLock.Lock()
	defer s.Data.Store.DbLock.Unlock()

	tx, err := s.Data.DB.Begin()
	if err != nil {
		log.Debugln(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err = s.Data.Store.GetQueries().WithTx(tx).UpdateFollowerByIRI(context.Background(), db.UpdateFollowerByIRIParams{
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

func (s *Service) removeFollow(actor *url.URL) error {
	s.Data.Store.DbLock.Lock()
	defer s.Data.Store.DbLock.Unlock()

	tx, err := s.Data.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := s.Data.Store.GetQueries().WithTx(tx).RemoveFollowerByIRI(context.Background(), actor.String()); err != nil {
		return err
	}

	return tx.Commit()
}
