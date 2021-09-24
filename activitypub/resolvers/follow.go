package resolvers

import (
	"context"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	log "github.com/sirupsen/logrus"
)

func getPersonFromFollow(activity vocab.ActivityStreamsFollow) (apmodels.ActivityPubActor, error) {
	return GetResolvedPersonFromActor(activity.GetActivityStreamsActor())
}

// MakeFollowRequest will convert an inbound Follow request to our internal actor model.
func MakeFollowRequest(c context.Context, activity vocab.ActivityStreamsFollow) (*apmodels.ActivityPubActor, error) {
	person, err := getPersonFromFollow(activity)
	if err != nil {
		return nil, err
	}

	person.FollowIri = activity.GetJSONLDId().Get()

	return &person, nil
}

// MakeUnFollowRequest will convert an inbound Unfollow request to our internal actor model.
func MakeUnFollowRequest(c context.Context, activity vocab.ActivityStreamsUndo) *apmodels.ActivityPubActor {
	person, err := GetResolvedPersonFromActor(activity.GetActivityStreamsActor())
	if err != nil {
		log.Errorln("unable to resolve person from actor iri", person.ActorIri, err)
		return nil
	}

	return &person
}
