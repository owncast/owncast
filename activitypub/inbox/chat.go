package inbox

import (
	"fmt"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/microcosm-cc/bluemonday"
	"github.com/owncast/owncast/activitypub/resolvers"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/persistence/configrepository"
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

func handleEngagementActivity(eventType events.EventType, isLiveNotification bool, actorReference vocab.ActivityStreamsActorProperty, action string) error {
	configRepository := configrepository.Get()

	// Do nothing if displaying engagement actions has been turned off.
	if !configRepository.GetFederationShowEngagement() {
		return nil
	}

	// Do nothing if chat is disabled
	if configRepository.GetChatDisabled() {
		return nil
	}

	// Get actor of the action
	actor, err := resolvers.GetResolvedActorFromActorProperty(actorReference)
	if err != nil {
		return fmt.Errorf("unable to resolve actor for engagement activity: %w", err)
	}

	// Send chat message
	actorName := sanitizeActorName(actor.Name, actor.Username)
	actorIRI := actor.ActorIriString()

	userPrefix := fmt.Sprintf("%s ", actorName)
	var suffix string
	if isLiveNotification && action == events.FediverseEngagementLike {
		suffix = "liked that this stream went live."
	} else if action == events.FediverseEngagementLike {
		suffix = fmt.Sprintf("liked a post from %s.", configRepository.GetServerName())
	} else if isLiveNotification && action == events.FediverseEngagementRepost {
		suffix = "shared this stream with their followers."
	} else if action == events.FediverseEngagementRepost {
		suffix = fmt.Sprintf("shared a post from %s.", configRepository.GetServerName())
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

	if err := chat.SendFediverseAction(eventType, actor.FullUsername, image, body, actorIRI); err != nil {
		return err
	}

	return nil
}
