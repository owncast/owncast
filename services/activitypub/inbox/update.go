package inbox

import (
	"context"

	"github.com/go-fed/activity/streams/vocab"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/services/activitypub/apmodels"
)

func (s *Service) handleUpdateRequest(c context.Context, activity vocab.ActivityStreamsUpdate) error {
	// We only care about update events to followers.
	if !apmodels.IsFirstObjectActivityStreamsPerson(activity.GetActivityStreamsObject()) {
		return nil
	}

	actor, err := s.resolver.GetResolvedActorFromActorProperty(activity.GetActivityStreamsActor())
	if err != nil {
		log.Errorln(err)
		return err
	}

	return s.followers.Update(actor.ActorIriString(), actor.InboxString(), actor.SharedInboxString(), actor.Name, actor.FullUsername, actor.ImageString())
}
