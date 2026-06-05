package webhooks

import (
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/activitypub/events"
)

// SendFediverseEngagementFollowEvent dispatches a user-followed event
// to webhook destinations, resolving the actor IRI through the
// followers repository the service was constructed with.
func (s *Service) SendFediverseEngagementFollowEvent(iri string) {
	follower, err := s.followers.GetByIRI(iri)
	if err != nil {
		return
	}
	userFollowedEvent := events.FediverseEngagementFollowEvent{}
	userFollowedEvent.SetDefaults()
	userFollowedEvent.Name = follower.Name
	userFollowedEvent.Username = follower.Username
	userFollowedEvent.Image = follower.Image.String()

	s.sendFediverseEngagementEventFollow(userFollowedEvent)
}

func (s *Service) sendFediverseEngagementEventFollow(event events.FediverseEngagementFollowEvent) {
	webhookEvent := WebhookEvent{
		Type:      models.FediverseEngagementFollow,
		EventData: event,
	}

	s.SendEventToWebhooks(webhookEvent)
}
