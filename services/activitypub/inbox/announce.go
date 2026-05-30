package inbox

import (
	"context"
	"time"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"

	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/chat/events"
)

func (s *Service) handleAnnounceRequest(c context.Context, activity vocab.ActivityStreamsAnnounce) error {
	objectIRI, err := apmodels.GetIRIStringFromObjectProperty(activity.GetActivityStreamsObject())
	if err != nil {
		return errors.Wrap(err, "announce activity is missing object IRI")
	}

	actorIRI, err := apmodels.GetIRIStringFromActorProperty(activity.GetActivityStreamsActor())
	if err != nil {
		return errors.Wrap(err, "announce activity is missing actor IRI")
	}

	actorReference := activity.GetActivityStreamsActor()

	if hasPreviouslyhandled, err := s.persistence.HasPreviouslyHandledInboundActivity(objectIRI, actorIRI, events.FediverseEngagementRepost); hasPreviouslyhandled || err != nil {
		return errors.Wrap(err, "inbound activity of share/re-post has already been handled")
	}

	// Shares need to match a post we had already sent.
	_, isLiveNotification, timestamp, err := s.persistence.GetObjectByIRI(objectIRI)
	if err != nil {
		return errors.Wrap(err, "Could not find post locally")
	}

	// Don't allow old activities to be liked.
	if time.Since(timestamp) > maxAgeForEngagement {
		return errors.New("Activity is too old to be shared")
	}

	// Save as an accepted activity.
	if err := s.persistence.SaveInboundFediverseActivity(objectIRI, actorIRI, events.FediverseEngagementRepost, time.Now()); err != nil {
		return errors.Wrap(err, "unable to save inbound share/re-post activity")
	}

	return s.handleEngagementActivity(events.FediverseEngagementRepost, isLiveNotification, actorReference, events.FediverseEngagementRepost)
}
