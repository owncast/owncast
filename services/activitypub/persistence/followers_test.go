package persistence

import (
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/go-fed/activity/streams"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/utils"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

// _datastore is the test-local datastore handle. Each setup sets it
// to a fresh in-memory database; the package no longer maintains a
// global one.
var _datastore *datastore.Datastore

var followers = []models.Follower{}

func setup() {
	ds, err := datastore.SetupPersistence(":memory:", os.TempDir())
	if err != nil {
		panic(err)
	}
	_datastore = ds

	followersRepo := followersrepository.New(_datastore)

	number := 100
	for i := 0; i < number; i++ {
		u := createFakeFollower()
		actorIRI, _ := url.Parse(u.ActorIRI)
		inboxURL, _ := url.Parse(u.Inbox)
		requestIRI, _ := url.Parse("https://fake.fediverse.server/some/request")
		fakeRequest := streams.NewActivityStreamsFollow()
		followersRepo.Add(apmodels.ActivityPubActor{
			ActorIri:         actorIRI,
			Inbox:            inboxURL,
			Name:             u.Name,
			Username:         u.Username,
			FollowRequestIri: requestIRI,
			RequestObject:    fakeRequest,
		}, true)
		followers = append(followers, u)
	}
}

func TestQueryFollowers(t *testing.T) {
	followersRepo := followersrepository.New(_datastore)
	f, total, err := followersRepo.GetFollowers(10, 0)
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
	followersRepo := followersrepository.New(_datastore)
	f, total, err := followersRepo.GetFollowers(10, 10)
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
	followersRepo := followersrepository.New(_datastore)
	f, total, err := followersRepo.GetFollowers(10, 90)
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
	followersRepo := followersrepository.New(_datastore)
	f, _, err := followersRepo.GetFollowers(15, 10)
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

func createTestFollower(followersRepo followersrepository.FollowersRepository, actor, inbox, sharedInbox, request, name, username string) {
	actorIRI, _ := url.Parse(actor)
	inboxURL, _ := url.Parse(inbox)
	var sharedInboxURL *url.URL
	if sharedInbox != "" {
		sharedInboxURL, _ = url.Parse(sharedInbox)
	}
	requestIRI, _ := url.Parse(request)
	fakeRequest := streams.NewActivityStreamsFollow()
	followersRepo.Add(apmodels.ActivityPubActor{
		ActorIri:         actorIRI,
		Inbox:            inboxURL,
		SharedInbox:      sharedInboxURL,
		Name:             name,
		Username:         username,
		FollowRequestIri: requestIRI,
		RequestObject:    fakeRequest,
	}, true)
}

func TestGetUniqueDeliveryInboxes(t *testing.T) {
	// Set up a fresh database for this test
	ds, err := datastore.SetupPersistence(":memory:", os.TempDir())
	if err != nil {
		t.Fatalf("setup error: %s", err)
	}
	_datastore = ds
	followersRepo := followersrepository.New(ds)

	// Create followers from server1 with a shared inbox (3 users, 1 shared inbox)
	server1SharedInbox := "https://server1.example.com/inbox"
	for i := 0; i < 3; i++ {
		user, _ := utils.GenerateRandomString(10)
		createTestFollower(
			followersRepo,
			"https://server1.example.com/user/"+user,
			"https://server1.example.com/user/"+user+"/inbox",
			server1SharedInbox,
			"https://server1.example.com/follow/"+user,
			user,
			user,
		)
	}

	// Create followers from server2 with a shared inbox (2 users, 1 shared inbox)
	server2SharedInbox := "https://server2.example.com/inbox"
	for i := 0; i < 2; i++ {
		user, _ := utils.GenerateRandomString(10)
		createTestFollower(
			followersRepo,
			"https://server2.example.com/user/"+user,
			"https://server2.example.com/user/"+user+"/inbox",
			server2SharedInbox,
			"https://server2.example.com/follow/"+user,
			user,
			user,
		)
	}

	// Create followers from server3 WITHOUT a shared inbox (2 users, 2 individual inboxes)
	for i := 0; i < 2; i++ {
		user, _ := utils.GenerateRandomString(10)
		createTestFollower(
			followersRepo,
			"https://server3.example.com/user/"+user,
			"https://server3.example.com/user/"+user+"/inbox",
			"",
			"https://server3.example.com/follow/"+user,
			user,
			user,
		)
	}

	// Total: 7 followers, but should result in 4 unique delivery inboxes:
	// - 1 shared inbox for server1
	// - 1 shared inbox for server2
	// - 2 individual inboxes for server3

	inboxes, err := followersRepo.GetUniqueDeliveryInboxes()
	if err != nil {
		t.Fatalf("Error getting unique delivery inboxes: %s", err)
	}

	if len(inboxes) != 4 {
		t.Errorf("Expected 4 unique delivery inboxes, got %d: %v", len(inboxes), inboxes)
	}

	// Verify the shared inboxes are included
	hasServer1Shared := false
	hasServer2Shared := false
	server3IndividualCount := 0

	for _, inbox := range inboxes {
		if inbox == server1SharedInbox {
			hasServer1Shared = true
		}
		if inbox == server2SharedInbox {
			hasServer2Shared = true
		}
		if len(inbox) > 0 && inbox != server1SharedInbox && inbox != server2SharedInbox {
			// Should be one of server3's individual inboxes
			if !strings.Contains(inbox, "server3.example.com") {
				t.Errorf("Unexpected inbox in results: %s", inbox)
			}
			server3IndividualCount++
		}
	}

	if !hasServer1Shared {
		t.Error("Expected server1 shared inbox to be in results")
	}
	if !hasServer2Shared {
		t.Error("Expected server2 shared inbox to be in results")
	}
	if server3IndividualCount != 2 {
		t.Errorf("Expected 2 individual inboxes from server3, got %d", server3IndividualCount)
	}
}

func TestSharedInboxPreferredOverIndividual(t *testing.T) {
	// Set up a fresh database for this test
	ds, err := datastore.SetupPersistence(":memory:", os.TempDir())
	if err != nil {
		t.Fatalf("setup error: %s", err)
	}
	_datastore = ds
	followersRepo := followersrepository.New(ds)

	// Create a single follower with both individual and shared inbox
	sharedInbox := "https://mastodon.social/inbox"
	individualInbox := "https://mastodon.social/users/testuser/inbox"

	createTestFollower(
		followersRepo,
		"https://mastodon.social/users/testuser",
		individualInbox,
		sharedInbox,
		"https://mastodon.social/follow/123",
		"Test User",
		"testuser",
	)

	inboxes, err := followersRepo.GetUniqueDeliveryInboxes()
	if err != nil {
		t.Fatalf("Error getting unique delivery inboxes: %s", err)
	}

	if len(inboxes) != 1 {
		t.Errorf("Expected 1 unique delivery inbox, got %d", len(inboxes))
	}

	// The shared inbox should be returned, not the individual inbox
	if inboxes[0] != sharedInbox {
		t.Errorf("Expected shared inbox %s, got %s", sharedInbox, inboxes[0])
	}
}

func TestIndividualInboxUsedWhenNoSharedInbox(t *testing.T) {
	// Set up a fresh database for this test
	ds, err := datastore.SetupPersistence(":memory:", os.TempDir())
	if err != nil {
		t.Fatalf("setup error: %s", err)
	}
	_datastore = ds
	followersRepo := followersrepository.New(ds)

	// Create a follower without a shared inbox
	individualInbox := "https://pleroma.example.com/users/testuser/inbox"

	createTestFollower(
		followersRepo,
		"https://pleroma.example.com/users/testuser",
		individualInbox,
		"",
		"https://pleroma.example.com/follow/123",
		"Test User",
		"testuser",
	)

	inboxes, err := followersRepo.GetUniqueDeliveryInboxes()
	if err != nil {
		t.Fatalf("Error getting unique delivery inboxes: %s", err)
	}

	if len(inboxes) != 1 {
		t.Errorf("Expected 1 unique delivery inbox, got %d", len(inboxes))
	}

	// The individual inbox should be returned when no shared inbox exists
	if inboxes[0] != individualInbox {
		t.Errorf("Expected individual inbox %s, got %s", individualInbox, inboxes[0])
	}
}
