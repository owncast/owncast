package inbox

import (
	"context"
	"time"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/core/chat/events"
	"github.com/pkg/errors"
)

func handleAnnounceRequest(c context.Context, activity vocab.ActivityStreamsAnnounce) error {
	object := activity.GetActivityStreamsObject()
	actorReference := activity.GetActivityStreamsActor()
	objectIRI := object.At(0).GetIRI().String()
	actorIRI := actorReference.At(0).GetIRI().String()

	if hasPreviouslyhandled, err := persistence.HasPreviouslyHandledInboundActivity(objectIRI, actorIRI, events.FediverseEngagementRepost); hasPreviouslyhandled || err != nil {
		return errors.Wrap(err, "inbound activity of share/re-post has already been handled")
	}

	// Save as an accepted activity
	if err := persistence.SaveInboundFediverseActivity(objectIRI, actorIRI, events.FediverseEngagementRepost, time.Now()); err != nil {
		return errors.Wrap(err, "unable to save inbound share/re-post activity")
	}

	return handleEngagementActivity(events.FediverseEngagementRepost, object, actorReference, activity.GetJSONLDId().Get(), events.FediverseEngagementRepost)
}
