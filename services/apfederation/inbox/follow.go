package inbox

import (
	"context"
	"fmt"
	"time"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/apfederation/outbox"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

func (api *APInbox) handleFollowInboxRequest(c context.Context, activity vocab.ActivityStreamsFollow) error {
	follow, err := api.resolvers.MakeFollowRequest(c, activity)
	if err != nil {
		log.Errorln("unable to create follow inbox request", err)
		return err
	}

	if follow == nil {
		return fmt.Errorf("unable to handle request")
	}

	approved := !api.configRepository.GetFederationIsPrivate()

	followRequest := *follow

	if err := api.federationRepository.AddFollow(followRequest, approved); err != nil {
		log.Errorln("unable to save follow request", err)
		return err
	}

	localAccountName := api.configRepository.GetDefaultFederationUsername()

	ob := outbox.Get()

	if approved {
		if err := ob.SendFollowAccept(follow.Inbox, activity, localAccountName); err != nil {
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
		hasPreviouslyhandled, err = api.federationRepository.HasPreviouslyHandledInboundActivity(objectIRI, actorIRI, models.FediverseEngagementFollow)
		if err != nil {
			log.Errorln("error checking for previously handled follow activity", err)
		}
	}

	// Save this follow action to our activities table.
	if err := api.federationRepository.SaveInboundFediverseActivity(objectIRI, actorIRI, models.FediverseEngagementFollow, time.Now()); err != nil {
		return errors.Wrap(err, "unable to save inbound share/re-post activity")
	}

	// Send action to chat if it has not been previously handled.
	if !hasPreviouslyhandled {
		return api.handleEngagementActivity(models.FediverseEngagementFollow, false, actorReference, models.FediverseEngagementFollow)
	}

	return nil
}

func (api *APInbox) handleUnfollowRequest(c context.Context, activity vocab.ActivityStreamsUndo) error {
	request := api.resolvers.MakeUnFollowRequest(c, activity)
	if request == nil {
		log.Errorf("unable to handle unfollow request")
		return errors.New("unable to handle unfollow request")
	}

	unfollowRequest := *request
	log.Traceln("unfollow request:", unfollowRequest)

	return api.federationRepository.RemoveFollow(unfollowRequest)
}
