package webhooks

import (
	"time"

	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
)

// SendStreamStatusEvent dispatches a stream-status event (started,
// stopped, title-updated, …) to all configured webhook destinations.
func (s *Service) SendStreamStatusEvent(eventType models.EventType) {
	s.sendStreamStatusEvent(eventType, shortid.MustGenerate(), time.Now())
}

func (s *Service) sendStreamStatusEvent(eventType models.EventType, id string, timestamp time.Time) {
	configRepository := configrepository.Get()

	s.SendEventToWebhooks(WebhookEvent{
		Type: eventType,
		EventData: map[string]interface{}{
			"id":          id,
			"name":        configRepository.GetServerName(),
			"summary":     configRepository.GetServerSummary(),
			"streamTitle": configRepository.GetStreamTitle(),
			"status":      s.getStatus(),
			"serverURL":   getServerURL(),
			"timestamp":   timestamp,
		},
	})
}
