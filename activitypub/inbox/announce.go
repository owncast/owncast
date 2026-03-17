package inbox

import (
	"context"
	"time"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/persistence/activitiesrepository"
	"github.com/owncast/owncast/activitypub/persistence/outboxrepository"
	"github.com/owncast/owncast/core/chat/events"
	"github.com/pkg/errors"
)

func handleAnnounceRequest(c context.Context, activity vocab.ActivityStreamsAnnounce) error {
	objectIRI, err := apmodels.GetIRIStringFromObjectProperty(activity.GetActivityStreamsObject())
	if err != nil {
		return errors.Wrap(err, "announce activity is missing object IRI")
	}

	actorIRI, err := apmodels.GetIRIStringFromActorProperty(activity.GetActivityStreamsActor())
	if err != nil {
		return errors.Wrap(err, "announce activity is missing actor IRI")
	}

	actorReference := activity.GetActivityStreamsActor()

	activitiesRepo := activitiesrepository.Get()
	if hasPreviouslyhandled, err := activitiesRepo.HasPreviouslyHandled(objectIRI, actorIRI, events.FediverseEngagementRepost); hasPreviouslyhandled || err != nil {
		return errors.Wrap(err, "inbound activity of share/re-post has already been handled")
	}

	// Shares need to match a post we had already sent.
	outboxRepo := outboxrepository.Get()
	_, isLiveNotification, timestamp, err := outboxRepo.GetObjectByIRI(objectIRI)
	if err != nil {
		return errors.Wrap(err, "Could not find post locally")
	}

	// Don't allow old activities to be liked
	if time.Since(timestamp) > maxAgeForEngagement {
		return errors.New("Activity is too old to be shared")
	}

	// Save as an accepted activity
	if err := activitiesRepo.Save(objectIRI, actorIRI, events.FediverseEngagementRepost, time.Now()); err != nil {
		return errors.Wrap(err, "unable to save inbound share/re-post activity")
	}

	return handleEngagementActivity(events.FediverseEngagementRepost, isLiveNotification, actorReference, events.FediverseEngagementRepost)
}
