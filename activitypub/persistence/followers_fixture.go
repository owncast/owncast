//go:build fixture
// +build fixture

package persistence

import (
	"encoding/json"
	"fmt"
	"os"

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

	// Iterate over the followers array
	for _, follower := range followers {
		createFollow(follower.ActorIRI, follower.Inbox, "", follower.Name, follower.Username, follower.Image, nil, true)
	}
}
