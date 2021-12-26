package inbox

import (
	"context"
	"fmt"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/activitypub/resolvers"
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/data"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

func handleFollowInboxRequest(c context.Context, activity vocab.ActivityStreamsFollow) error {
	follow, err := resolvers.MakeFollowRequest(c, activity)
	if err != nil {
		log.Errorln("unable to create follow inbox request", err)
		return err
	}

	if follow == nil {
		return fmt.Errorf("unable to handle request")
	}

	approved := !data.GetFederationIsPrivate()

	followRequest := *follow

	if err := persistence.AddFollow(followRequest, approved); err != nil {
		log.Errorln("unable to save follow request", err)
		return err
	}

	localAccountName := data.GetDefaultFederationUsername()

	if approved {
		if err := requests.SendFollowAccept(follow.Inbox, follow.FollowRequestIri, localAccountName); err != nil {
			log.Errorln("unable to send follow accept", err)
			return err
		}
	}

	actorReference := activity.GetActivityStreamsActor()

	if approved {
		return handleEngagementActivity(events.FediverseEngagementFollow, activity.GetActivityStreamsObject(), actorReference, activity.GetJSONLDId().Get(), events.FediverseEngagementFollow)
	}

	return nil
}

func handleUnfollowRequest(c context.Context, activity vocab.ActivityStreamsUndo) error {
	request := resolvers.MakeUnFollowRequest(c, activity)
	if request == nil {
		log.Errorf("unable to handle unfollow request")
		return errors.New("unable to handle unfollow request")
	}

	unfollowRequest := *request
	log.Traceln("unfollow request:", unfollowRequest)

	return persistence.RemoveFollow(unfollowRequest)
}
