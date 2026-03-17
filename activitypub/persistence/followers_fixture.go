//go:build fixture
// +build fixture

package persistence

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/go-fed/activity/streams"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/persistence/followersrepository"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

func addFollowersFixtureData() {
	log.Println("Adding followers fixture data...")
	file, err := os.Open("./test/fixture/followers_fixture.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var followers []models.Follower
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&followers)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	followersRepo := followersrepository.Get()

	for _, follower := range followers {
		actorIRI, _ := url.Parse(follower.ActorIRI)
		inboxURL, _ := url.Parse(follower.Inbox)
		imageURL, _ := url.Parse(follower.Image)
		requestIRI, _ := url.Parse("https://fixture.example.com/follow-request")
		fakeRequest := streams.NewActivityStreamsFollow()

		followersRepo.Add(apmodels.ActivityPubActor{
			ActorIri:         actorIRI,
			Inbox:            inboxURL,
			Name:             follower.Name,
			Username:         follower.Username,
			Image:            imageURL,
			FollowRequestIri: requestIRI,
			RequestObject:    fakeRequest,
		}, true)
	}
}
