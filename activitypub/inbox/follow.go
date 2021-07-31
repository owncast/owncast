package inbox

import (
	"context"
	"fmt"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/activitypub/resolvers"
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

	// We only care about undoing follows right now.
	request := resolvers.MakeUnFollowRequest(activity, c)
	if request == nil {
		return fmt.Errorf("unable to handle unfollow request")
	}

	unfollowRequest := *request
	fmt.Println("unfollow request:", unfollowRequest)

	return persistence.RemoveFollow(unfollowRequest)
}
