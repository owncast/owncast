package inbox

import (
	"context"
	"time"

	"code.superseriousbusiness.org/activity/streams/vocab"
	"github.com/pkg/errors"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	activityevents "github.com/owncast/owncast/services/activitypub/events"
	"github.com/owncast/owncast/services/chat/events"
)

type engagementRequest struct {
	objectError       string
	actorError        string
	maxAgeError       string
	resolveActorError string
	duplicateError    string
	saveError         string
	internalEventType models.EventType
	persistedEvent    events.EventType
	chatEvent         events.EventType
	chatAction        string
}

func (s *Service) acceptEngagementRequest(c context.Context, objectProperty vocab.ActivityStreamsObjectProperty, actorProperty vocab.ActivityStreamsActorProperty, request engagementRequest) error {
	objectIRI, err := apmodels.GetIRIStringFromObjectProperty(objectProperty)
	if err != nil {
		return errors.Wrap(err, request.objectError)
	}

	if actorProperty == nil || actorProperty.Empty() || actorProperty.Len() == 0 || actorProperty.At(0) == nil {
		return errors.New(request.actorError)
	}

	_, isLiveNotification, timestamp, err := s.persistence.GetObjectByIRI(objectIRI)
	if err != nil {
		return errors.Wrap(err, "Could not find post locally")
	}

	if time.Since(timestamp) > maxAgeForEngagement {
		return errors.New(request.maxAgeError)
	}

	actor, err := s.resolver.GetResolvedActorFromActorProperty(actorProperty)
	if err != nil {
		return errors.Wrap(err, request.resolveActorError)
	}
	actorIRI := actor.ActorIriString()
	if hasPreviouslyHandled, err := s.persistence.HasPreviouslyHandledInboundActivity(objectIRI, actorIRI, request.persistedEvent); hasPreviouslyHandled || err != nil {
		return errors.Wrap(err, request.duplicateError)
	}

	if err := s.persistence.SaveInboundFediverseActivity(objectIRI, actorIRI, request.persistedEvent, time.Now()); err != nil {
		return errors.Wrap(err, request.saveError)
	}

	s.publishFediverseEvent(c, request.internalEventType, &activityevents.FediverseEngagementEvent{
		Actor:  fediverseActorFromResolvedActor(actor),
		Target: &activityevents.FediverseTarget{URL: objectIRI},
	})

	return s.handleEngagementActivity(request.chatEvent, isLiveNotification, actor, request.chatAction)
}
