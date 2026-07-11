package inbox

import (
	"fmt"

	"github.com/microcosm-cc/bluemonday"

	"github.com/owncast/owncast/services/activitypub/apmodels"
	activityevents "github.com/owncast/owncast/services/activitypub/events"
	"github.com/owncast/owncast/services/chat/events"
)

// sanitizeActorName strips HTML tags from the ActivityPub actor display name.
// Falls back to the username if the display name is empty or entirely HTML.
func sanitizeActorName(displayName, username string) string {
	strict := bluemonday.StrictPolicy()
	name := strict.Sanitize(displayName)
	if name == "" {
		name = strict.Sanitize(username)
	}
	return name
}

func fediverseActorFromResolvedActor(actor apmodels.ActivityPubActor) activityevents.FediverseActor {
	return activityevents.FediverseActor{
		Name:   sanitizeActorName(actor.Name, actor.Username),
		Handle: actor.FullUsername,
		URL:    actor.ActorIriString(),
		Image:  actor.ImageString(),
	}
}

func (s *Service) handleEngagementActivity(eventType events.EventType, isLiveNotification bool, actor apmodels.ActivityPubActor, action string) error {
	// Do nothing if displaying engagement actions has been turned off.
	if !s.configRepository.GetFederationShowEngagement() {
		return nil
	}

	// Do nothing if chat is disabled
	if s.configRepository.GetChatDisabled() {
		return nil
	}

	// Send chat message
	actorName := sanitizeActorName(actor.Name, actor.Username)
	actorIRI := actor.ActorIriString()

	userPrefix := fmt.Sprintf("%s ", actorName)
	var suffix string
	if isLiveNotification && action == events.FediverseEngagementLike {
		suffix = "liked that this stream went live."
	} else if action == events.FediverseEngagementLike {
		suffix = fmt.Sprintf("liked a post from %s.", s.configRepository.GetServerName())
	} else if isLiveNotification && action == events.FediverseEngagementRepost {
		suffix = "shared this stream with their followers."
	} else if action == events.FediverseEngagementRepost {
		suffix = fmt.Sprintf("shared a post from %s.", s.configRepository.GetServerName())
	} else if action == events.FediverseEngagementFollow {
		suffix = "followed this stream."
	} else {
		return fmt.Errorf("could not handle event for sending to chat: %s", action)
	}
	body := fmt.Sprintf("%s %s", userPrefix, suffix)

	var image *string
	if imageStr := actor.ImageString(); imageStr != "" {
		image = &imageStr
	}

	if err := s.chat.SendFediverseAction(eventType, actor.FullUsername, image, body, actorIRI); err != nil {
		return err
	}

	return nil
}
