package resolvers

import (
	"context"
	"fmt"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	log "github.com/sirupsen/logrus"
)

func getPersonFromFollow(activity vocab.ActivityStreamsFollow) (vocab.ActivityStreamsPerson, error) {
	return GetResolvedPersonFromActor(activity.GetActivityStreamsActor())
}

// MakeFollowRequest will convert an inbound Follow request to our internal actor model.
func MakeFollowRequest(c context.Context, activity vocab.ActivityStreamsFollow) (*apmodels.ActivityPubActor, error) {
	person, err := getPersonFromFollow(activity)
	if err != nil {
		return nil, err
	}

	hostname := person.GetJSONLDId().GetIRI().Hostname()
	username := person.GetActivityStreamsPreferredUsername().GetXMLSchemaString()
	fullUsername := fmt.Sprintf("%s@%s", username, hostname)

	followRequest := apmodels.ActivityPubActor{
		ActorIri:  person.GetJSONLDId().Get(),
		FollowIri: activity.GetJSONLDId().Get(),
		Inbox:     person.GetActivityStreamsInbox().GetIRI(),
		Name:      person.GetActivityStreamsName().At(0).GetXMLSchemaString(),
		Username:  fullUsername,
		Image:     person.GetActivityStreamsIcon().At(0).GetActivityStreamsImage().GetActivityStreamsUrl().Begin().GetIRI(),
	}

	return &followRequest, nil
}

// MakeUnFollowRequest will convert an inbound Unfollow request to our internal actor model.
func MakeUnFollowRequest(c context.Context, activity vocab.ActivityStreamsUndo) *apmodels.ActivityPubActor {
	person, err := GetResolvedPersonFromActor(activity.GetActivityStreamsActor())
	if err != nil {
		log.Errorln("unable to resolve person from actor iri", person.GetJSONLDId(), err)
		return nil
	}

	if person == nil {
		return nil
	}

	request := apmodels.ActivityPubActor{
		ActorIri:  person.GetJSONLDId().Get(),
		FollowIri: activity.GetJSONLDId().Get(),
		Inbox:     person.GetActivityStreamsInbox().GetIRI(),
		Name:      person.GetActivityStreamsName().Begin().GetXMLSchemaString(),
		Image:     person.GetActivityStreamsImage().Begin().GetIRI(),
	}

	return &request
}
