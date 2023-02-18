package webhooks

import (
	"time"

	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/models"
)

// SendStreamStatusEvent will send all webhook destinations the current stream status.
func (s *Service) SendStreamStatusEvent(eventType models.EventType) {
	s.sendStreamStatusEvent(eventType, shortid.MustGenerate(), time.Now())
}

func (s *Service) sendStreamStatusEvent(eventType models.EventType, id string, timestamp time.Time) {
	s.SendEventToWebhooks(WebhookEvent{
		Type: eventType,
		EventData: map[string]interface{}{
			"id":          id,
			"name":        s.data.GetServerName(),
			"summary":     s.data.GetServerSummary(),
			"streamTitle": s.data.GetStreamTitle(),
			"timestamp":   timestamp,
		},
	})
}
