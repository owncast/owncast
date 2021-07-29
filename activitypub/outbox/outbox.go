package outbox

import (
	"encoding/json"

	"github.com/go-fed/activity/streams"
	"github.com/owncast/owncast/activitypub/models"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/core/data"
)

func SendLive() {

}

// SendPublic will send a public message to all followers.
func SendPublic(textContent string) {
	localActor := models.MakeLocalIRIForAccount(data.GetDefaultFederationUsername())
	message := models.CreateMessage(textContent, localActor)

	var jsonmap map[string]interface{}
	jsonmap, _ = streams.Serialize(message)
	b, _ := json.Marshal(jsonmap)

	followers, err := persistence.GetFederationFollowers()
	if err != nil {
		panic(err)
	}

	for _, follower := range followers {
		if _, err := requests.PostSignedRequest(b, follower.Inbox, localActor); err != nil {
			panic(err)
		}
	}
}
