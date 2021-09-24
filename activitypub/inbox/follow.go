package inbox

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/activitypub/resolvers"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/data"

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
	log.Println("follow request:", followRequest)

	if err := persistence.AddFollow(followRequest, approved); err != nil {
		log.Errorln("unable to save follow request", err)
		return err
	}

	localAccountName := c.Value("account").(string)

	if approved {
		if err := requests.SendFollowAccept(followRequest, localAccountName); err != nil {
			return err
		}
	}

	actorReference := activity.GetActivityStreamsActor()
	actor, _ := resolvers.GetResolvedPersonFromActor(actorReference)
	actorIRI := actorReference.Begin().GetIRI().String()

	msg := fmt.Sprintf("[%s](%s) just **followed**!", actor.Username, actorIRI)

	return chat.SendSystemMessage(msg, false)
}

func handleUndoInboxRequest(c context.Context, activity vocab.ActivityStreamsUndo) error {
	// Determine if this is an undo of a follow, favorite, announce, etc.
	o := activity.GetActivityStreamsObject()
	for iter := o.Begin(); iter != o.End(); iter = iter.Next() {
		if iter.IsActivityStreamsFollow() {
			// This is an Unfollow request
			if err := handleUnfollowRequest(c, activity); err != nil {
				return err
			}
		} else {
			log.Println("Undo", iter.GetType().GetTypeName(), "ignored")
			continue
		}
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
	log.Println("unfollow request:", unfollowRequest)

	return persistence.RemoveFollow(unfollowRequest)
}

func handleLikeRequest(c context.Context, activity vocab.ActivityStreamsLike) error {
	log.Debugln("handleLikeRequest")

	object := activity.GetActivityStreamsObject()
	actorReference := activity.GetActivityStreamsActor()
	return handleEngagementActivity(object, actorReference, "liked")
}

func handleAnnounceRequest(c context.Context, activity vocab.ActivityStreamsAnnounce) error {
	log.Debugln("handleAnnounceRequest")

	object := activity.GetActivityStreamsObject()
	actorReference := activity.GetActivityStreamsActor()
	return handleEngagementActivity(object, actorReference, "re-posted")
}

func handleEngagementActivity(object vocab.ActivityStreamsObjectProperty, actorReference vocab.ActivityStreamsActorProperty, action string) error {
	log.Debugln("handleEngagementActivity")

	IRIPath := object.At(0).GetIRI().Path

	// Verify we actually sent this post.
	post, err := persistence.GetObjectByIRI(IRIPath)
	if err != nil || post == "" {
		log.Errorln("Could not find post locally:", IRIPath, err)
		return fmt.Errorf("Could not find post locally: %s", IRIPath)
	}

	// Get actor of the Like
	actor, _ := resolvers.GetResolvedPersonFromActor(actorReference)

	// Send chat message
	actorIRI := actorReference.Begin().GetIRI().String()

	msg := fmt.Sprintf("[%s](%s) just **%s** one of %s's posts.", actor.Username, actorIRI, action, data.GetServerName())

	if err := chat.SendSystemMessage(msg, false); err != nil {
		return err
	}

	return nil
}
