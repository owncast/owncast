package webhooks

import (
	"testing"
	"time"

	"github.com/owncast/owncast/activitypub/events"
	"github.com/owncast/owncast/models"
)

func TestSendFediverseEngagementEventFollow(t *testing.T) {
	checkPayload(t, models.FediverseEngagementFollow, func() {
		sendFediverseEngagementEventFollow(events.FediverseEngagementFollowEvent{
			Event: events.Event{
				Timestamp: time.Unix(72, 6).UTC(),
				ID:        "id",
			},
			Name:     "be",
			Username: "be@witch.me",
		})
	}, `{
		"id": "id",
		"image": "",
		"name": "be",
		"timestamp": "1970-01-01T00:01:12.000000006Z",
		"username": "be@witch.me"
		}`)
}
