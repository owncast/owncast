package inbox

import (
	"context"

	"code.superseriousbusiness.org/activity/streams/vocab"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/chat/events"
)

func (s *Service) handleAnnounceRequest(c context.Context, activity vocab.ActivityStreamsAnnounce) error {
	return s.acceptEngagementRequest(c, activity.GetActivityStreamsObject(), activity.GetActivityStreamsActor(), engagementRequest{
		objectError:       "announce activity is missing object IRI",
		actorError:        "announce activity is missing actor IRI",
		maxAgeError:       "Activity is too old to be shared",
		resolveActorError: "unable to resolve actor of share/re-post activity",
		saveError:         "unable to save inbound share/re-post activity",
		internalEventType: models.FediverseEngagementRepost,
		persistedEvent:    events.FediverseEngagementRepost,
		chatEvent:         events.FediverseEngagementRepost,
		chatAction:        events.FediverseEngagementRepost,
	})
}
