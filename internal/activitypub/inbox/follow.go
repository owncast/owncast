package inbox

import (
	"context"
	"fmt"
	"time"

	"github.com/go-fed/activity/streams/vocab"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/core/chat/events"
)

func (s *Service) handleFollowInboxRequest(c context.Context, activity vocab.ActivityStreamsFollow) error {
	follow, err := s.follower.MakeFollowRequest(c, activity)
	if err != nil {
		log.Errorln("unable to create follow inbox request", err)
		return err
	}

	if follow == nil {
		return fmt.Errorf("unable to handle request")
	}

	approved := !s.data.GetFederationIsPrivate()

	followRequest := *follow

	if err := s.follower.AddFollow(followRequest, approved); err != nil {
		log.Errorln("unable to save follow request", err)
		return err
	}

	localAccountName := s.data.GetDefaultFederationUsername()

	if approved {
		if err := s.follower.SendFollowAccept(follow.Inbox, activity, localAccountName); err != nil {
			log.Errorln("unable to send follow accept", err)
			return err
		}
	}

	// Save as an accepted activity
	actorReference := activity.GetActivityStreamsActor()
	object := activity.GetActivityStreamsObject()
	objectIRI := object.At(0).GetIRI().String()
	actorIRI := actorReference.At(0).GetIRI().String()

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
		return s.handleEngagementActivity(events.FediverseEngagementFollow, false, actorReference, events.FediverseEngagementFollow)
	}

	return nil
}

func (s *Service) handleUnfollowRequest(c context.Context, activity vocab.ActivityStreamsUndo) error {
	request := s.follower.MakeUnFollowRequest(c, activity)
	if request == nil {
		log.Errorf("unable to handle unfollow request")
		return errors.New("unable to handle unfollow request")
	}

	unfollowRequest := *request
	log.Traceln("unfollow request:", unfollowRequest)

	return s.follower.RemoveFollow(unfollowRequest)
}
