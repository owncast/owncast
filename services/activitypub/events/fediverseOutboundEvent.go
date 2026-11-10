package events

import (
	"context"
	"encoding/json"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/dispatcher"
)

// PublishOutbound reports an activity after its local state and delivery queue
// have been persisted. Authentication messages must never enter this path.
func PublishOutbound(ctx context.Context, target *dispatcher.Dispatcher, payload []byte) {
	if target == nil || len(payload) == 0 {
		return
	}
	target.Publish(ctx, dispatcher.Event{
		Type:    models.FediverseOutboundActivity,
		Payload: append(json.RawMessage(nil), payload...),
	})
}
