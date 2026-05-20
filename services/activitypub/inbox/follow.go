package inbox

import (
	"context"
	"fmt"
	"time"

	"github.com/go-fed/activity/streams/vocab"
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

	if err := s.followers.Add(followRequest, approved); err != nil {
		log.Errorln("unable to save follow request", err)
		return err
	}

	localAccountName := s.configRepository.GetDefaultFederationUsername()

	objectIRI, err := apmodels.GetIRIStringFromObjectProperty(activity.GetActivityStreamsObject())
	if err != nil {
		return errors.Wrap(err, "follow activity is missing object IRI")
	}

	actorIRI, err := apmodels.GetIRIStringFromActorProperty(activity.GetActivityStreamsActor())
	if err != nil {
		return errors.Wrap(err, "follow activity is missing actor IRI")
	}

	actorReference := activity.GetActivityStreamsActor()

	if approved {
		if err := requests.SendFollowAccept(s.workerpool, follow.Inbox, activity, localAccountName, s.builder, s.signer); err != nil {
			log.Errorln("unable to send follow accept", err)
			return err
		}
		go s.webhooks.SendFediverseEngagementFollowEvent(actorIRI)
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
		return s.handleEngagementActivity(events.FediverseEngagementFollow, false, actorReference, events.FediverseEngagementFollow)
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
