package inbox

import (
	"context"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/persistence/followersrepository"
	"github.com/owncast/owncast/activitypub/resolvers"
	log "github.com/sirupsen/logrus"
)

func handleUpdateRequest(c context.Context, activity vocab.ActivityStreamsUpdate) error {
	// We only care about update events to followers.
	if !apmodels.IsFirstObjectActivityStreamsPerson(activity.GetActivityStreamsObject()) {
		return nil
	}

	actor, err := resolvers.GetResolvedActorFromActorProperty(activity.GetActivityStreamsActor())
	if err != nil {
		log.Errorln(err)
		return err
	}

	followersRepo := followersrepository.Get()
	return followersRepo.Update(actor.ActorIriString(), actor.InboxString(), actor.SharedInboxString(), actor.Name, actor.FullUsername, actor.ImageString())
}
