package webhooks

import (
	"github.com/owncast/owncast/activitypub/events"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/models"
)

// SendFediverseEventFollow will send a user followed event to webhook destinations
func SendUserFollowedEvent(iri string) {
	follower, err := persistence.GetFollower(iri)
	if err != nil {
		return
	}
	userFollowedEvent := events.FediverseEngagementFollowEvent{}
	userFollowedEvent.SetDefaults()
	userFollowedEvent.Name = follower.Name
	userFollowedEvent.Username = follower.Username
	userFollowedEvent.Image = follower.Image.String()

	sendFediverseEventFollow(userFollowedEvent)
}

func sendFediverseEventFollow(event events.FediverseEngagementFollowEvent) {
	webhookEvent := WebhookEvent{
		Type:      models.FediverseEngagementFollow,
		EventData: event,
	}

	SendEventToWebhooks(webhookEvent)
}
