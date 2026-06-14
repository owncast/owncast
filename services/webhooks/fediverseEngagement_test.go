package webhooks

import (
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/activitypub/events"
)

func TestSendFediverseEngagementEventFollow(t *testing.T) {
	checkPayload(t, models.FediverseEngagementFollow, func() {
		testSvc.sendFediverseEngagementEventFollow(events.FediverseEngagementFollowEvent{
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
		"serverURL": "http://localhost:8080",
		"timestamp": "1970-01-01T00:01:12.000000006Z",
		"username": "be@witch.me"
		}`)
}
