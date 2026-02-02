package webhooks

import (
	"github.com/owncast/owncast/activitypub/events"
	"github.com/owncast/owncast/activitypub/persistence/followersrepository"
	"github.com/owncast/owncast/models"
)

// SendFediverseEventFollow will send a user followed event to webhook
// destinations.
func SendFediverseEngagementFollowEvent(iri string) {
	followersRepo := followersrepository.Get()
	follower, err := followersRepo.GetByIRI(iri)
	if err != nil {
		return
	}
	userFollowedEvent := events.FediverseEngagementFollowEvent{}
	userFollowedEvent.SetDefaults()
	userFollowedEvent.Name = follower.Name
	userFollowedEvent.Username = follower.Username
	userFollowedEvent.Image = follower.Image.String()

	sendFediverseEngagementEventFollow(userFollowedEvent)
}

func sendFediverseEngagementEventFollow(event events.FediverseEngagementFollowEvent) {
	webhookEvent := WebhookEvent{
		Type:      models.FediverseEngagementFollow,
		EventData: event,
	}

	SendEventToWebhooks(webhookEvent)
}
