package persistence

import (
	"os"
	"testing"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

var followers = []models.Follower{}

func setup() {
	data.SetupPersistence(":memory:")
	_datastore = data.GetDatastore()
	createFederationFollowersTable()

	number := 100
	for i := 0; i < number; i++ {
		u := createFakeFollower()
		createFollow(u.ActorIRI, u.Inbox, "https://fake.fediverse.server/some/request", u.Name, u.Username, u.Image, nil, true)
		followers = append(followers, u)
	}
}

func TestQueryFollowers(t *testing.T) {
	f, total, err := GetFederationFollowers(10, 0)
	if err != nil {
		t.Errorf("Error querying followers: %s", err)
	}

	if len(f) != 10 {
		t.Errorf("Expected 10 followers, got %d", len(f))
	}

	if total != 100 {
		t.Errorf("Expected 100 followers, got %d", total)
	}
}

func TestQueryFollowersWithOffset(t *testing.T) {
	f, total, err := GetFederationFollowers(10, 10)
	if err != nil {
		t.Errorf("Error querying followers: %s", err)
	}

	if len(f) != 10 {
		t.Errorf("Expected 10 followers, got %d", len(f))
	}

	if total != 100 {
		t.Errorf("Expected 100 followers, got %d", total)
	}
}

func TestQueryFollowersWithOffsetAndLimit(t *testing.T) {
	f, total, err := GetFederationFollowers(10, 90)
	if err != nil {
		t.Errorf("Error querying followers: %s", err)
	}

	if len(f) != 10 {
		t.Errorf("Expected 10 followers, got %d", len(f))
	}

	if total != 100 {
		t.Errorf("Expected 100 followers, got %d", total)
	}
}

func TestQueryFollowersWithPagination(t *testing.T) {
	f, _, err := GetFederationFollowers(15, 10)
	if err != nil {
		t.Errorf("Error querying followers: %s", err)
	}

	comparisonFollowers := followers[10:25]
	if len(f) != len(comparisonFollowers) {
		t.Errorf("Expected %d followers, got %d", len(comparisonFollowers), len(f))
	}

	for i, follower := range f {
		if follower.ActorIRI != comparisonFollowers[i].ActorIRI {
			t.Errorf("Expected %s, got %s", comparisonFollowers[i].ActorIRI, follower.ActorIRI)
		}
	}
}

func createFakeFollower() models.Follower {
	user, _ := utils.GenerateRandomString(10)

	return models.Follower{
		ActorIRI:  "https://freedom.eagle/user/" + user,
		Inbox:     "https://fake.fediverse.server/user/" + user + "/inbox",
		Image:     "https://fake.fediverse.server/user/" + user + "/avatar.png",
		Name:      user,
		Username:  user,
		Timestamp: utils.NullTime{},
	}
}
