package follower

import (
	"context"
	"fmt"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/internal/activitypub/apmodels"
)

func (s *Service) getPersonFromFollow(activity vocab.ActivityStreamsFollow) (apmodels.ActivityPubActor, error) {
	return s.GetResolvedActorFromActorProperty(activity.GetActivityStreamsActor())
}

// MakeFollowRequest will convert an inbound Follow request to our internal actor model.
func (s *Service) MakeFollowRequest(c context.Context, activity vocab.ActivityStreamsFollow) (*apmodels.ActivityPubActor, error) {
	person, err := s.getPersonFromFollow(activity)
	if err != nil {
		return nil, errors.New("unable to resolve person from follow request: " + err.Error())
	}

	hostname := person.ActorIri.Hostname()
	username := person.Username
	fullUsername := fmt.Sprintf("%s@%s", username, hostname)

	followRequest := apmodels.ActivityPubActor{
		ActorIri:         person.ActorIri,
		FollowRequestIri: activity.GetJSONLDId().Get(),
		Inbox:            person.Inbox,
		Name:             person.Name,
		Username:         fullUsername,
		Image:            person.Image,
		RequestObject:    activity,
	}

	return &followRequest, nil
}

// MakeUnFollowRequest will convert an inbound Unfollow request to our internal actor model.
func (s *Service) MakeUnFollowRequest(c context.Context, activity vocab.ActivityStreamsUndo) *apmodels.ActivityPubActor {
	person, err := s.GetResolvedActorFromActorProperty(activity.GetActivityStreamsActor())
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
