package inbox

import (
	"context"
	"time"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"

	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/chat/events"
)

func (s *Service) handleLikeRequest(c context.Context, activity vocab.ActivityStreamsLike) error {
	objectIRI, err := apmodels.GetIRIStringFromObjectProperty(activity.GetActivityStreamsObject())
	if err != nil {
		return errors.Wrap(err, "like activity is missing object IRI")
	}

	actorIRI, err := apmodels.GetIRIStringFromActorProperty(activity.GetActivityStreamsActor())
	if err != nil {
		return errors.Wrap(err, "like activity is missing actor IRI")
	}

	actorReference := activity.GetActivityStreamsActor()

	if hasPreviouslyhandled, err := s.persistence.HasPreviouslyHandledInboundActivity(objectIRI, actorIRI, events.FediverseEngagementLike); hasPreviouslyhandled || err != nil {
		return errors.Wrap(err, "inbound activity of like has already been handled")
	}

	// Likes need to match a post we had already sent.
	_, isLiveNotification, timestamp, err := s.persistence.GetObjectByIRI(objectIRI)
	if err != nil {
		return errors.Wrap(err, "Could not find post locally")
	}

	// Don't allow old activities to be liked.
	if time.Since(timestamp) > maxAgeForEngagement {
		return errors.New("Activity is too old to be liked")
	}

	// Save as an accepted activity.
	if err := s.persistence.SaveInboundFediverseActivity(objectIRI, actorIRI, events.FediverseEngagementLike, time.Now()); err != nil {
		return errors.Wrap(err, "unable to save inbound like activity")
	}

	return s.handleEngagementActivity(events.FediverseEngagementLike, isLiveNotification, actorReference, events.FediverseEngagementLike)
}
