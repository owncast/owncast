package inbox

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"code.superseriousbusiness.org/activity/streams/vocab"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/activitypub/requests"
	"github.com/owncast/owncast/services/chat/events"
)

func (s *Service) handleFollowInboxRequest(c context.Context, activity vocab.ActivityStreamsFollow) error {
	follow, err := s.resolver.MakeFollowRequest(c, activity)
	if err != nil {
		log.Errorln("unable to create follow inbox request", err)
		return err
	}

	if follow == nil {
		return fmt.Errorf("unable to handle request")
	}

	approved := !s.configRepository.GetFederationIsPrivate()

	followRequest := *follow

	// A directory follow (a server that sent the ns#directory marker so it can
	// list our stream) always requires explicit approval, regardless of whether
	// this server otherwise accepts follows automatically. Being listed by a
	// directory is a different relationship from gaining a fan, so the operator
	// opts in per directory. The Accept is sent later by the admin approval
	// flow, not here.
	if followRequest.IsDirectory {
		approved = false
	}

	actorIRI := followRequest.ActorIriString()
	if actorIRI == "" {
		return errors.New("follow activity is missing actor IRI")
	}

	localAccountName := s.configRepository.GetDefaultFederationUsername()
	handled, err := s.createOrHandleFollow(c, activity, follow, followRequest, approved, actorIRI, localAccountName)
	if err != nil {
		return err
	}
	if handled {
		return nil
	}
	return s.handleNewFollow(activity, followRequest, approved, actorIRI)
}

func (s *Service) createOrHandleFollow(c context.Context, activity vocab.ActivityStreamsFollow, follow *apmodels.ActivityPubActor, followRequest apmodels.ActivityPubActor, approved bool, actorIRI, localAccountName string) (bool, error) {
	// The duplicate check and insert must be one critical section. Some
	// ActivityPub servers deliver the same Follow concurrently.
	ds := s.persistence.Datastore()
	ds.DbLock.Lock()
	defer ds.DbLock.Unlock()

	if handled, err := s.respondToExistingFollow(activity, follow, actorIRI, localAccountName); handled {
		return true, err
	}
	if err := s.addFollow(c, activity, followRequest, approved, localAccountName); err != nil {
		return false, err
	}
	return false, nil
}

func (s *Service) handleNewFollow(activity vocab.ActivityStreamsFollow, followRequest apmodels.ActivityPubActor, approved bool, actorIRI string) error {
	objectIRI, err := apmodels.GetIRIStringFromObjectProperty(activity.GetActivityStreamsObject())
	if err != nil {
		return errors.Wrap(err, "follow activity is missing object IRI")
	}

	if approved && !followRequest.IsDirectory {
		go s.webhooks.SendFediverseEngagementFollowEvent(actorIRI)
	}

	// A directory follow is a listing relationship, not a fan follow. It is kept
	// and accepted above because we need it to deliver stream-status pings to
	// that directory, but it must not be surfaced as a new follower in chat or
	// the activity feed.
	if followRequest.IsDirectory {
		return nil
	}

	// If this request is approved and we have not previously sent an action to
	// chat due to a previous follow request, then do so.
	hasPreviouslyhandled := true // Default so we don't send anything if it fails.
	if approved {
		hasPreviouslyhandled, err = s.persistence.HasPreviouslyHandledInboundActivity(objectIRI, actorIRI, events.FediverseEngagementFollow)
		if err != nil {
			log.Errorln("error checking for previously handled follow activity", err)
		}
	}

	// Save this follow action to our activities table.
	if err := s.persistence.SaveInboundFediverseActivity(objectIRI, actorIRI, events.FediverseEngagementFollow, time.Now()); err != nil {
		return errors.Wrap(err, "unable to save inbound share/re-post activity")
	}

	// Send action to chat if it has not been previously handled.
	if !hasPreviouslyhandled {
		return s.handleEngagementActivity(events.FediverseEngagementFollow, false, followRequest, events.FediverseEngagementFollow)
	}

	return nil
}

func (s *Service) respondToExistingFollow(activity vocab.ActivityStreamsFollow, follow *apmodels.ActivityPubActor, actorIRI, localAccountName string) (bool, error) {
	existing, err := s.followers.GetByIRI(actorIRI)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return true, err
	}
	if existing.DisabledAt != nil {
		return true, requests.SendFollowReject(s.workerpool, follow.Inbox, activity, localAccountName, s.builder)
	}
	if existing.ApprovedAt != nil {
		return true, requests.SendFollowAccept(s.workerpool, follow.Inbox, activity, localAccountName, s.builder, s.configRepository, false)
	}
	return true, nil
}

func (s *Service) addFollow(c context.Context, activity vocab.ActivityStreamsFollow, follow apmodels.ActivityPubActor, approved bool, localAccountName string) error {
	ds := s.persistence.Datastore()
	if !approved {
		tx, err := ds.DB.BeginTx(c, nil)
		if err != nil {
			return errors.Wrap(err, "beginning follow transaction")
		}
		defer func() { _ = tx.Rollback() }()
		if err := s.followers.AddTx(c, tx, follow, false); err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return errors.Wrap(err, "committing follow transaction")
		}
		return nil
	}

	delivery, err := requests.MakeFollowAcceptDelivery(follow.Inbox, activity, localAccountName, s.builder, s.configRepository, false)
	if err != nil {
		return err
	}

	tx, err := ds.DB.BeginTx(c, nil)
	if err != nil {
		return errors.Wrap(err, "beginning follow transaction")
	}
	defer func() { _ = tx.Rollback() }()
	if err := s.followers.AddTx(c, tx, follow, true); err != nil {
		return err
	}
	if err := s.workerpool.EnqueueTx(c, tx, delivery); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "committing follow transaction")
	}
	return nil
}

func (s *Service) handleUnfollowRequest(c context.Context, activity vocab.ActivityStreamsUndo) error {
	request := s.resolver.MakeUnFollowRequest(c, activity)
	if request == nil {
		log.Errorf("unable to handle unfollow request")
		return errors.New("unable to handle unfollow request")
	}

	unfollowRequest := *request
	log.Traceln("unfollow request:", unfollowRequest)

	return s.followers.Remove(unfollowRequest)
}
