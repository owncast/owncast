package resolvers

import (
	"context"
	"fmt"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
)

func getPersonFromFollow(activity vocab.ActivityStreamsFollow, c context.Context) (vocab.ActivityStreamsPerson, error) {
	fmt.Println("getPersonFromFollow...")

	return GetResolvedPersonFromActor(activity.GetActivityStreamsActor())
}

func MakeFollowRequest(activity vocab.ActivityStreamsFollow, c context.Context) (*apmodels.ActivityPubActor, error) {
	person, err := getPersonFromFollow(activity, c)
	if err != nil {
		return nil, err
	}

	fmt.Println(activity.GetJSONLDId().Get().String())

	followRequest := apmodels.ActivityPubActor{
		ActorIri:  person.GetJSONLDId().Get(),
		FollowIri: activity.GetJSONLDId().Get(),
		Inbox:     person.GetActivityStreamsInbox().GetIRI(),
	}

	return &followRequest, nil
}

func MakeUnFollowRequest(activity vocab.ActivityStreamsUndo, c context.Context) *apmodels.ActivityPubActor {
	person, err := GetResolvedPersonFromActor(activity.GetActivityStreamsActor())
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if person == nil {
		return nil
	}

	fmt.Println(activity.GetJSONLDId().Get().String())

	request := apmodels.ActivityPubActor{
		ActorIri:  person.GetJSONLDId().Get(),
		FollowIri: activity.GetJSONLDId().Get(),
		Inbox:     person.GetActivityStreamsInbox().GetIRI(),
	}

	return &request
}
