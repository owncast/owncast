package inbox

import (
	"context"

	"code.superseriousbusiness.org/activity/streams/vocab"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/chat/events"
)

func (s *Service) handleLikeRequest(c context.Context, activity vocab.ActivityStreamsLike) error {
	return s.acceptEngagementRequest(c, activity.GetActivityStreamsObject(), activity.GetActivityStreamsActor(), engagementRequest{
		objectError:       "like activity is missing object IRI",
		actorError:        "like activity is missing actor IRI",
		maxAgeError:       "Activity is too old to be liked",
		resolveActorError: "unable to resolve actor of like activity",
		duplicateError:    "inbound activity of like has already been handled",
		saveError:         "unable to save inbound like activity",
		internalEventType: models.FediverseEngagementLike,
		persistedEvent:    events.FediverseEngagementLike,
		chatEvent:         events.FediverseEngagementLike,
		chatAction:        events.FediverseEngagementLike,
	})
}
