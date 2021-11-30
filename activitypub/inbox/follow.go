package inbox

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/activitypub/resolvers"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/data"
	"github.com/teris-io/shortid"

	log "github.com/sirupsen/logrus"
)

func handleFollowInboxRequest(c context.Context, activity vocab.ActivityStreamsFollow) error {
	log.Printf("INBOX Follow: %v ", activity)

	follow, err := resolvers.MakeFollowRequest(c, activity)
	if err != nil {
		log.Errorln("unable to create follow inbox request", err)
		return err
	}

	if follow == nil {
		return fmt.Errorf("unable to handle request")
	}

	approved := !data.GetFollowApprovalRequired()

	followRequest := *follow
	log.Println("follow request:", followRequest)

	if err := persistence.AddFollow(followRequest, approved); err != nil {
		log.Errorln("unable to save follow request", err)
		return err
	}

	localAccountName := data.GetDefaultFederationUsername()

	if approved {
		if err := requests.SendFollowAccept(followRequest, localAccountName); err != nil {
			log.Errorln("unable to send follow accept", err)
			return err
		}
	}

	actorReference := activity.GetActivityStreamsActor()

	return handleEngagementActivity(activity.GetActivityStreamsObject(), actorReference, activity.GetJSONLDId().Get(), "follow")
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
	return handleEngagementActivity(object, actorReference, activity.GetJSONLDId().Get(), "liked")
}

func handleAnnounceRequest(c context.Context, activity vocab.ActivityStreamsAnnounce) error {
	log.Debugln("handleAnnounceRequest")

	object := activity.GetActivityStreamsObject()
	actorReference := activity.GetActivityStreamsActor()
	return handleEngagementActivity(object, actorReference, activity.GetJSONLDId().Get(), "re-posted")
}

func handleEngagementActivity(object vocab.ActivityStreamsObjectProperty, actorReference vocab.ActivityStreamsActorProperty, activityIRI *url.URL, action string) error {
	// Do nothing if displaying engagement actions has been turned off.
	if !data.GetFederationShowEngagement() {
		return nil
	}

	IRIPath := object.At(0).GetIRI().Path

	// for iter := object.Begin(); iter != object.End(); iter = iter.Next() {
	// Verify we actually sent this post.
	post, err := persistence.GetObjectByIRI(IRIPath)
	if err != nil || post == "" {
		log.Errorln("Could not find post locally:", IRIPath, err)
		// return err
	}

	// Get actor of the Like
	actor, _ := resolvers.GetResolvedActorFromActorProperty(actorReference)

	// Send chat message
	actorName := actor.Name
	actorIRI := actorReference.Begin().GetIRI().String()

	body := fmt.Sprintf("[%s](%s) just **%s** one of %s's posts.", actorName, actorIRI, action, data.GetServerName())
	var image string
	if actor.Image != nil {
		s := actor.Image.String()
		image = s
	}

	internalID := shortid.MustGenerate()
	if err := persistence.SaveFediverseActivity(internalID, activityIRI.String(), actorIRI, events.FediverseEngagementFollow, time.Now()); err != nil {
		return err
	}

	if err := chat.SendFediverseAction(events.FediverseEngagementFollow, actorName, actor.Username, image, body, actorIRI); err != nil {
		return err
	}

	return nil
}
