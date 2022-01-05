package inbox

import (
	"fmt"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/resolvers"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/data"
)

func handleEngagementActivity(eventType events.EventType, isLiveNotification bool, actorReference vocab.ActivityStreamsActorProperty, action string) error {
	// Do nothing if displaying engagement actions has been turned off.
	if !data.GetFederationShowEngagement() {
		return nil
	}

	// Do nothing if chat is disabled
	if data.GetChatDisabled() {
		return nil
	}

	// Get actor of the action
	actor, _ := resolvers.GetResolvedActorFromActorProperty(actorReference)

	// Send chat message
	actorName := actor.Name
	if actorName == "" {
		actorName = actor.Username
	}
	actorIRI := actorReference.Begin().GetIRI().String()

	userPrefix := fmt.Sprintf("%s ", actorName)
	var suffix string
	if isLiveNotification && action == events.FediverseEngagementLike {
		suffix = "liked that this stream went live."
	} else if action == events.FediverseEngagementLike {
		suffix = fmt.Sprintf("liked a post from %s.", data.GetServerName())
	} else if isLiveNotification && action == events.FediverseEngagementRepost {
		suffix = "shared this stream with their followers."
	} else if action == events.FediverseEngagementRepost {
		suffix = fmt.Sprintf("shared a post from %s.", data.GetServerName())
	} else if action == events.FediverseEngagementFollow {
		suffix = "followed this stream."
	} else {
		return fmt.Errorf("could not handle event for sending to chat: %s", action)
	}
	body := fmt.Sprintf("%s %s", userPrefix, suffix)

	var image *string
	if actor.Image != nil {
		s := actor.Image.String()
		image = &s
	}

	if err := chat.SendFediverseAction(eventType, actor.FullUsername, image, body, actorIRI); err != nil {
		return err
	}

	return nil
}
