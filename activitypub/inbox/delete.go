package inbox

import (
	"context"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/resolvers"
	log "github.com/sirupsen/logrus"
)

func handleDeleteRequest(c context.Context, activity vocab.ActivityStreamsDelete) error {
	// Only need to handle deletes for followers
	if !activity.GetActivityStreamsObject().At(0).IsActivityStreamsPerson() {
		return nil
	}

	actor, err := resolvers.GetResolvedActorFromActorProperty(activity.GetActivityStreamsActor())
	if err != nil {
		log.Error(err)
		return err
	}

	return persistence.RemoveFollow(actor)
}
