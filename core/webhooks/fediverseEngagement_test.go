package webhooks

import (
	"testing"
	"time"

	"github.com/owncast/owncast/activitypub/events"
	"github.com/owncast/owncast/models"
)

func TestSendFediverseEventFollow(t *testing.T) {
	checkPayload(t, models.FediverseEngagementFollow, func() {
		sendFediverseEventFollow(events.FediverseEngagementFollowEvent{
			Event: events.Event{
				Timestamp: time.Unix(72, 6).UTC(),
				ID: "id",
				Type: models.FediverseEngagementFollow,
			},
			Name: "be",
			Username: "be@witch.me",
		})
	}, `{
		"id": "id",
		"image": "",
		"name": "be",
		"timestamp": "1970-01-01T00:01:12.000000006Z",
		"type": "FEDIVERSE_ENGAGEMENT_FOLLOW",
		"username": "be@witch.me"
		}`)
}
