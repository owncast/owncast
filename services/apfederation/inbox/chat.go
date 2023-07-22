package inbox

import (
	"fmt"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/models"
)

func (api *APInbox) handleEngagementActivity(eventType models.EventType, isLiveNotification bool, actorReference vocab.ActivityStreamsActorProperty, action string) error {
	// Do nothing if displaying engagement actions has been turned off.
	if !api.configRepository.GetFederationShowEngagement() {
		return nil
	}

	// Do nothing if chat is disabled
	if api.configRepository.GetChatDisabled() {
		return nil
	}

	// Get actor of the action
	actor, _ := api.resolvers.GetResolvedActorFromActorProperty(actorReference)

	// Send chat message
	actorName := actor.Name
	if actorName == "" {
		actorName = actor.Username
	}
	actorIRI := actorReference.Begin().GetIRI().String()

	userPrefix := fmt.Sprintf("%s ", actorName)
	var suffix string
	if isLiveNotification && action == models.FediverseEngagementLike {
		suffix = "liked that this stream went live."
	} else if action == models.FediverseEngagementLike {
		suffix = fmt.Sprintf("liked a post from %s.", api.configRepository.GetServerName())
	} else if isLiveNotification && action == models.FediverseEngagementRepost {
		suffix = "shared this stream with their followers."
	} else if action == models.FediverseEngagementRepost {
		suffix = fmt.Sprintf("shared a post from %s.", api.configRepository.GetServerName())
	} else if action == models.FediverseEngagementFollow {
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

	if err := api.chatService.SendFediverseAction(eventType, actor.FullUsername, image, body, actorIRI); err != nil {
		return err
	}

	return nil
}
