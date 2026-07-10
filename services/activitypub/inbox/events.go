package inbox

import (
	"context"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/dispatcher"
)

func (s *Service) publishFediverseEvent(ctx context.Context, eventType models.EventType, payload any) {
	if s.events == nil {
		return
	}

	s.events.Publish(ctx, dispatcher.Event{Type: eventType, Payload: payload})
}
