package inbox

import (
	"context"
	"time"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"

	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/internal/activitypub/persistence"
)

func (s *Service) handleAnnounceRequest(c context.Context, activity vocab.ActivityStreamsAnnounce, p *persistence.Service) error {
	object := activity.GetActivityStreamsObject()
	actorReference := activity.GetActivityStreamsActor()
	objectIRI := object.At(0).GetIRI().String()
	actorIRI := actorReference.At(0).GetIRI().String()

	if hasPreviouslyhandled, err := p.HasPreviouslyHandledInboundActivity(objectIRI, actorIRI, events.FediverseEngagementRepost); hasPreviouslyhandled || err != nil {
		return errors.Wrap(err, "inbound activity of share/re-post has already been handled")
	}

	// Shares need to match a post we had already sent.
	_, isLiveNotification, timestamp, err := p.GetObjectByIRI(objectIRI)
	if err != nil {
		return errors.Wrap(err, "Could not find post locally")
	}

	// Don't allow old activities to be liked
	if time.Since(timestamp) > maxAgeForEngagement {
		return errors.New("Activity is too old to be shared")
	}

	// Save as an accepted activity
	if err := p.SaveInboundFediverseActivity(objectIRI, actorIRI, events.FediverseEngagementRepost, time.Now()); err != nil {
		return errors.Wrap(err, "unable to save inbound share/re-post activity")
	}

	return s.handleEngagementActivity(events.FediverseEngagementRepost, isLiveNotification, actorReference, events.FediverseEngagementRepost)
}
