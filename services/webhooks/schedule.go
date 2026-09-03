package webhooks

import (
	"time"

	"github.com/owncast/owncast/models"
)

// SendScheduledEvent dispatches a scheduled-event lifecycle notification.
func (s *Service) SendScheduledEvent(event models.ScheduledEvent, action models.ScheduledEventWebhookAction) {
	s.sendScheduledEvent(event, action, time.Now())
}

func (s *Service) sendScheduledEvent(event models.ScheduledEvent, action models.ScheduledEventWebhookAction, timestamp time.Time) {
	s.SendEventToWebhooks(WebhookEvent{
		Type: models.ScheduledEvents,
		EventData: &WebhookScheduledEventData{
			BaseWebhookData: BaseWebhookData{
				Status:    s.getStatus(),
				ServerURL: s.serverURL(),
			},
			Event:     event,
			Action:    action,
			Timestamp: timestamp,
		},
	})
}
