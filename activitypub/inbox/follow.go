package inbox

import (
	"context"
	"fmt"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/activitypub/resolvers"
	"github.com/owncast/owncast/core/chat"

	log "github.com/sirupsen/logrus"
)

func handleFollowInboxRequest(c context.Context, activity vocab.ActivityStreamsFollow) error {
	follow, err := resolvers.MakeFollowRequest(activity, c)
	if err != nil {
		log.Errorln("unable to create follow inbox request", err)
		return err
	}

	if follow == nil {
		return fmt.Errorf("unable to handle request")
	}

	followRequest := *follow
	log.Println("follow request:", followRequest)

	if err := persistence.AddFollow(followRequest); err != nil {
		log.Errorln("unable to save follow request", err)
		return err
	}

	localAccountName := c.Value("account").(string)

	if err := requests.SendFollowAccept(followRequest, localAccountName); err != nil {
		return err
	}

	actorReference := activity.GetActivityStreamsActor()
	actor, _ := resolvers.GetResolvedPersonFromActor(actorReference)
	actorName := actor.GetActivityStreamsName().Begin().GetXMLSchemaString()
	actorIRI := actorReference.Begin().GetIRI().String()

	msg := fmt.Sprintf("[%s](%s) just **followed**!", actorName, actorIRI)
	chat.SendSystemMessage(msg, false)

	return nil
}

func handleUndoInboxRequest(c context.Context, activity vocab.ActivityStreamsUndo) error {
	// Determine if this is an undo of a follow, favorite, announce, etc.
	o := activity.GetActivityStreamsObject()
	for iter := o.Begin(); iter != o.End(); iter = iter.Next() {
		if iter.IsActivityStreamsFollow() {
			// This is an Unfollow request
			handleUnfollowRequest(c, activity)
		} else {
			log.Println("Undo", iter.GetType().GetTypeName(), "ignored")
			continue
		}
	}

	return nil
}

func handleUnfollowRequest(c context.Context, activity vocab.ActivityStreamsUndo) {
	request := resolvers.MakeUnFollowRequest(activity, c)
	if request == nil {
		log.Errorf("unable to handle unfollow request")
		return
	}

	unfollowRequest := *request
	log.Println("unfollow request:", unfollowRequest)

	persistence.RemoveFollow(unfollowRequest)
}

func handleLikeRequest(c context.Context, activity vocab.ActivityStreamsLike) error {
	log.Debugln("handleLikeRequest")

	object := activity.GetActivityStreamsObject()
	actorReference := activity.GetActivityStreamsActor()
	handleEngagementActivity(object, actorReference, "liked")
	return nil
}

func handleAnnounceRequest(c context.Context, activity vocab.ActivityStreamsAnnounce) error {
	log.Debugln("handleAnnounceRequest")

	object := activity.GetActivityStreamsObject()
	actorReference := activity.GetActivityStreamsActor()
	handleEngagementActivity(object, actorReference, "re-posted")
	return nil
}

func handleEngagementActivity(object vocab.ActivityStreamsObjectProperty, actorReference vocab.ActivityStreamsActorProperty, action string) {
	log.Debugln("handleEngagementActivity")

	for iter := object.Begin(); iter != object.End(); iter = iter.Next() {
		postIRI := iter.GetIRI().String()
		// Verify we actually sent this post.
		post, err := persistence.GetObjectByIRI(postIRI)
		if err != nil || post == "" {
			log.Errorln("Could not find post locally:", postIRI, err)
			// return
		}

		// Get actor of the Like
		actor, _ := resolvers.GetResolvedPersonFromActor(actorReference)

		// Send chat message
		actorName := actor.GetActivityStreamsName().Begin().GetXMLSchemaString()
		actorIRI := actorReference.Begin().GetIRI().String()

		msg := fmt.Sprintf("[%s](%s) just **%s** [this post](%s)", actorName, actorIRI, action, postIRI)
		chat.SendSystemMessage(msg, false)
	}
}
