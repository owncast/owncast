package resolvers

import (
	"context"
	"fmt"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/services/activitypub/apmodels"
)

func (r *Resolver) getPersonFromFollow(activity vocab.ActivityStreamsFollow) (apmodels.ActivityPubActor, error) {
	return r.GetResolvedActorFromActorProperty(activity.GetActivityStreamsActor())
}

// MakeFollowRequest will convert an inbound Follow request to our internal actor model.
func (r *Resolver) MakeFollowRequest(c context.Context, activity vocab.ActivityStreamsFollow) (*apmodels.ActivityPubActor, error) {
	person, err := r.getPersonFromFollow(activity)
	if err != nil {
		return nil, errors.Wrap(err, "unable to resolve person from follow request")
	}

	hostname := person.ActorIriHostname()
	username := person.Username
	fullUsername := fmt.Sprintf("%s@%s", username, hostname)

	// Check if this is an Owncast server by looking for custom properties in the follow request
	unknownProps := activity.GetUnknownProperties()
	metadata := apmodels.ParseOwncastMetadata(unknownProps)

	followRequest := apmodels.ActivityPubActor{
		ActorIri:         person.ActorIri,
		FollowRequestIri: activity.GetJSONLDId().Get(),
		Inbox:            person.Inbox,
		SharedInbox:      person.SharedInbox,
		Name:             person.Name,
		Username:         fullUsername,
		Image:            person.Image,
		RequestObject:    activity,
		IsOwncastServer:  metadata.IsOwncastServer,
	}

	return &followRequest, nil
}

// MakeUnFollowRequest will convert an inbound Unfollow request to our internal actor model.
func (r *Resolver) MakeUnFollowRequest(c context.Context, activity vocab.ActivityStreamsUndo) *apmodels.ActivityPubActor {
	person, err := r.GetResolvedActorFromActorProperty(activity.GetActivityStreamsActor())
	if err != nil {
		log.Errorln("unable to resolve person from actor iri", person.ActorIri, err)
		return nil
	}

	unfollowRequest := apmodels.ActivityPubActor{
		ActorIri:         person.ActorIri,
		FollowRequestIri: activity.GetJSONLDId().Get(),
		Inbox:            person.Inbox,
		Name:             person.Name,
		Image:            person.Image,
	}

	return &unfollowRequest
}
