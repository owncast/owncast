package resolvers

import (
	"context"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	log "github.com/sirupsen/logrus"
)

func getPersonFromFollow(activity vocab.ActivityStreamsFollow, c context.Context) (vocab.ActivityStreamsPerson, error) {
	return GetResolvedPersonFromActor(activity.GetActivityStreamsActor())
}

func MakeFollowRequest(activity vocab.ActivityStreamsFollow, c context.Context) (*apmodels.ActivityPubActor, error) {
	person, err := getPersonFromFollow(activity, c)
	if err != nil {
		return nil, err
	}

	log.Println(activity.GetJSONLDId().Get().String())

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
		log.Errorln("unable to resolve person from actor iri", person.GetJSONLDId(), err)
		return nil
	}

	if person == nil {
		return nil
	}

	log.Println(activity.GetJSONLDId().Get().String())

	request := apmodels.ActivityPubActor{
		ActorIri:  person.GetJSONLDId().Get(),
		FollowIri: activity.GetJSONLDId().Get(),
		Inbox:     person.GetActivityStreamsInbox().GetIRI(),
	}

	return &request
}
