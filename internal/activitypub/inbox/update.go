package inbox

import (
	"context"

	"github.com/go-fed/activity/streams/vocab"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/internal/activitypub/persistence"
	"github.com/owncast/owncast/internal/activitypub/resolvers"
)

func handleUpdateRequest(c context.Context, activity vocab.ActivityStreamsUpdate) error {
	// We only care about update events to followers.
	if !activity.GetActivityStreamsObject().At(0).IsActivityStreamsPerson() {
		return nil
	}

	actor, err := resolvers.GetResolvedActorFromActorProperty(activity.GetActivityStreamsActor())
	if err != nil {
		log.Errorln(err)
		return err
	}

	return persistence.UpdateFollower(actor.ActorIri.String(), actor.Inbox.String(), actor.Name, actor.FullUsername, actor.Image.String())
}
