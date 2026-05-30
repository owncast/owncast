//go:build fixture
// +build fixture

package persistence

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/go-fed/activity/streams"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
)

type fixtureFollower struct {
	ActorIRI string `json:"link"`
	Inbox    string `json:"inbox"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Image    string `json:"image"`
}

func (s *Service) addFollowersFixtureData() {
	log.Println("Adding followers fixture data...")
	file, err := os.Open("./test/fixture/followers_fixture.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var followers []fixtureFollower
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&followers); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	followersRepo := followersrepository.New(s.datastore)
	for _, f := range followers {
		actorIRI, _ := url.Parse(f.ActorIRI)
		inboxURL, _ := url.Parse(f.Inbox)
		if err := followersRepo.Add(apmodels.ActivityPubActor{
			ActorIri:      actorIRI,
			Inbox:         inboxURL,
			Name:          f.Name,
			Username:      f.Username,
			Image:         optionalURL(f.Image),
			RequestObject: streams.NewActivityStreamsFollow(),
		}, true); err != nil {
			log.Errorln("Error adding fixture follower:", err)
		}
	}
}

func optionalURL(s string) *url.URL {
	if s == "" {
		return nil
	}
	u, _ := url.Parse(s)
	return u
}
