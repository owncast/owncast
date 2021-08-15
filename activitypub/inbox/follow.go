package inbox

import (
	"context"
	"fmt"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/activitypub/resolvers"

	log "github.com/sirupsen/logrus"
)

func handleFollowInboxRequest(c context.Context, activity vocab.ActivityStreamsFollow) error {
	fmt.Println("followCallback fired!")

	follow, err := resolvers.MakeFollowRequest(activity, c)
	if err != nil {
		panic(err)
	}

	if follow == nil {
		return fmt.Errorf("unable to handle request")
	}

	followRequest := *follow
	fmt.Println("follow request:", followRequest)

	if err := persistence.AddFollow(followRequest); err != nil {
		panic(err)
	}

	localAccountName := c.Value("account").(string)

	if err := requests.SendFollowAccept(followRequest, localAccountName); err != nil {
		return err
	}

	return nil
}

func handleUndoInboxRequest(c context.Context, activity vocab.ActivityStreamsUndo) error {
	fmt.Println("handleUndoInboxRequest fired!")

	// Determine if this is an undo of a follow, favorite, announce, etc.
	o := activity.GetActivityStreamsObject()
	for iter := o.Begin(); iter != o.End(); iter = iter.Next() {
		if iter.IsActivityStreamsFollow() {
			// This is an Unfollow request
			handleUnfollowRequest(c, activity)
		} else {
			log.Println("Undo", iter.GetType().GetTypeName(), "ignored")
		}
	}

	return nil
}

func handleUnfollowRequest(c context.Context, activity vocab.ActivityStreamsUndo) {
	request := resolvers.MakeUnFollowRequest(activity, c)
	if request == nil {
		log.Errorf("unable to handle unfollow request")
	}

	unfollowRequest := *request
	log.Println("unfollow request:", unfollowRequest)

	persistence.RemoveFollow(unfollowRequest)
}
