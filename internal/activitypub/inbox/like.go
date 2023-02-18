package inbox

import (
	"context"
	"fmt"
	"time"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"

	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/internal/activitypub/persistence"
)

func (s *Service) handleLikeRequest(c context.Context, activity vocab.ActivityStreamsLike) error {
	object := activity.GetActivityStreamsObject()
	actorReference := activity.GetActivityStreamsActor()
	objectIRI := object.At(0).GetIRI().String()
	actorIRI := actorReference.At(0).GetIRI().String()

	p, err := persistence.New(s.data)
	if err != nil {
		return fmt.Errorf("getting persistent data store: %v", err)
	}

	if hasPreviouslyhandled, err := p.HasPreviouslyHandledInboundActivity(objectIRI, actorIRI, events.FediverseEngagementLike); hasPreviouslyhandled || err != nil {
		return errors.Wrap(err, "inbound activity of like has already been handled")
	}

	// Likes need to match a post we had already sent.
	_, isLiveNotification, timestamp, err := p.GetObjectByIRI(objectIRI)
	if err != nil {
		return errors.Wrap(err, "Could not find post locally")
	}

	// Don't allow old activities to be liked
	if time.Since(timestamp) > maxAgeForEngagement {
		return errors.New("Activity is too old to be liked")
	}

	// Save as an accepted activity
	if err := p.SaveInboundFediverseActivity(objectIRI, actorIRI, events.FediverseEngagementLike, time.Now()); err != nil {
		return errors.Wrap(err, "unable to save inbound like activity")
	}

	return s.handleEngagementActivity(events.FediverseEngagementLike, isLiveNotification, actorReference, events.FediverseEngagementLike)
}
