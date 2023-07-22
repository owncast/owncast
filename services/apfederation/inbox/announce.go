package inbox

import (
	"context"
	"time"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/models"
	"github.com/pkg/errors"
)

func (api *APInbox) handleAnnounceRequest(c context.Context, activity vocab.ActivityStreamsAnnounce) error {
	object := activity.GetActivityStreamsObject()
	actorReference := activity.GetActivityStreamsActor()
	objectIRI := object.At(0).GetIRI().String()
	actorIRI := actorReference.At(0).GetIRI().String()

	if hasPreviouslyhandled, err := api.federationRepository.HasPreviouslyHandledInboundActivity(objectIRI, actorIRI, models.FediverseEngagementRepost); hasPreviouslyhandled || err != nil {
		return errors.Wrap(err, "inbound activity of share/re-post has already been handled")
	}

	// Shares need to match a post we had already sent.
	_, isLiveNotification, timestamp, err := api.federationRepository.GetObjectByIRI(objectIRI)
	if err != nil {
		return errors.Wrap(err, "Could not find post locally")
	}

	// Don't allow old activities to be liked
	if time.Since(timestamp) > maxAgeForEngagement {
		return errors.New("Activity is too old to be shared")
	}

	// Save as an accepted activity
	if err := api.federationRepository.SaveInboundFediverseActivity(objectIRI, actorIRI, models.FediverseEngagementRepost, time.Now()); err != nil {
		return errors.Wrap(err, "unable to save inbound share/re-post activity")
	}

	return api.handleEngagementActivity(models.FediverseEngagementRepost, isLiveNotification, actorReference, models.FediverseEngagementRepost)
}
