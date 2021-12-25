package inbox

import (
	"context"
	"time"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/core/chat/events"
	"github.com/pkg/errors"
)

func handleLikeRequest(c context.Context, activity vocab.ActivityStreamsLike) error {
	object := activity.GetActivityStreamsObject()
	actorReference := activity.GetActivityStreamsActor()
	objectIRI := object.At(0).GetIRI().String()
	actorIRI := actorReference.At(0).GetIRI().String()

	if hasPreviouslyhandled, err := persistence.HasPreviouslyHandledInboundActivity(objectIRI, actorIRI, events.FediverseEngagementLike); hasPreviouslyhandled || err != nil {
		return errors.Wrap(err, "inbound activity of like has already been handled")
	}

	// Save as an accepted activity
	if err := persistence.SaveInboundFediverseActivity(objectIRI, actorIRI, events.FediverseEngagementLike, time.Now()); err != nil {
		return errors.Wrap(err, "unable to save inbound like activity")
	}

	return handleEngagementActivity(events.FediverseEngagementLike, object, actorReference, activity.GetJSONLDId().Get(), events.FediverseEngagementLike)
}
