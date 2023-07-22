package inbox

import (
	"context"

	"github.com/go-fed/activity/streams/vocab"
	log "github.com/sirupsen/logrus"
)

func (api *APInbox) handleUpdateRequest(c context.Context, activity vocab.ActivityStreamsUpdate) error {
	// We only care about update events to followers.
	if !activity.GetActivityStreamsObject().At(0).IsActivityStreamsPerson() {
		return nil
	}

	actor, err := api.resolvers.GetResolvedActorFromActorProperty(activity.GetActivityStreamsActor())
	if err != nil {
		log.Errorln(err)
		return err
	}

	return api.federationRepository.UpdateFollower(actor.ActorIri.String(), actor.Inbox.String(), actor.Name, actor.FullUsername, actor.Image.String())
}
