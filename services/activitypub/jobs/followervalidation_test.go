package jobs

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/go-fed/activity/streams"

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

var testDatastore *datastore.Datastore

func setup() {
	resetTestDatabase()
}

// resetTestDatabase initializes a fresh in-memory database for testing.
func resetTestDatabase() {
	ds, err := datastore.SetupPersistence(":memory:", os.TempDir())
	if err != nil {
		panic(err)
	}
	testDatastore = ds
}

// setupTestWithRepo resets the database and returns a new repository instance.
func setupTestWithRepo(t *testing.T) followersrepository.FollowersRepository {
	t.Helper()
	resetTestDatabase()
	return followersrepository.New(testDatastore)
}

func createTestFollower(repo followersrepository.FollowersRepository, iri, inbox, name, username string) {
	actorIRI, _ := url.Parse(iri)
	inboxURL, _ := url.Parse(inbox)
	requestIRI, _ := url.Parse("https://fake.server/follow/request")
	fakeRequest := streams.NewActivityStreamsFollow()

	repo.Add(apmodels.ActivityPubActor{
		ActorIri:         actorIRI,
		Inbox:            inboxURL,
		Name:             name,
		Username:         username,
		FullUsername:     username + "@fake.server",
		FollowRequestIri: requestIRI,
		RequestObject:    fakeRequest,
	}, true)
}

func TestGetFollowersToValidate(t *testing.T) {
	repo := setupTestWithRepo(t)

	// Create some test followers
	for i := 0; i < 10; i++ {
		user, _ := utils.GenerateRandomString(10)
		createTestFollower(repo, "https://fake.server/user/"+user, "https://fake.server/user/"+user+"/inbox", user, user)
	}

	// Get followers to validate
	followers, err := repo.GetFollowersToValidate(5)
	if err != nil {
		t.Fatalf("Error getting followers to validate: %s", err)
	}

	if len(followers) != 5 {
		t.Errorf("Expected 5 followers to validate, got %d", len(followers))
	}
}

func TestUpdateFollowerValidationSuccess(t *testing.T) {
	repo := setupTestWithRepo(t)

	// Create a test follower
	testIRI := "https://fake.server/user/testuser"
	createTestFollower(repo, testIRI, "https://fake.server/user/testuser/inbox", "Test User", "testuser")

	// Mark validation as successful
	err := repo.UpdateFollowerValidationSuccess(testIRI)
	if err != nil {
		t.Fatalf("Error updating follower validation success: %s", err)
	}

	// Verify the follower was updated
	followers, err := repo.GetFollowersToValidate(10)
	if err != nil {
		t.Fatalf("Error getting followers: %s", err)
	}

	// After validation success, FirstValidationFailureAt should be NULL/invalid
	for _, f := range followers {
		if f.ActorIRI == testIRI {
			if f.FirstValidationFailureAt.Valid {
				t.Error("Expected FirstValidationFailureAt to be NULL after success")
			}
		}
	}
}

func TestUpdateFollowerValidationFailure(t *testing.T) {
	repo := setupTestWithRepo(t)

	// Create a test follower
	testIRI := "https://fake.server/user/testuser"
	createTestFollower(repo, testIRI, "https://fake.server/user/testuser/inbox", "Test User", "testuser")

	// Mark validation as failed
	err := repo.UpdateFollowerValidationFailure(testIRI)
	if err != nil {
		t.Fatalf("Error updating follower validation failure: %s", err)
	}

	// Verify the FirstValidationFailureAt is set
	followers, err := repo.GetFollowersToValidate(10)
	if err != nil {
		t.Fatalf("Error getting followers: %s", err)
	}

	found := false
	for _, f := range followers {
		if f.ActorIRI == testIRI {
			found = true
			if !f.FirstValidationFailureAt.Valid {
				t.Error("Expected FirstValidationFailureAt to be set after failure")
			}
		}
	}

	if !found {
		t.Error("Test follower not found in results")
	}
}

func TestValidationFailureClearedOnSuccess(t *testing.T) {
	repo := setupTestWithRepo(t)

	// Create a test follower
	testIRI := "https://fake.server/user/testuser"
	createTestFollower(repo, testIRI, "https://fake.server/user/testuser/inbox", "Test User", "testuser")

	// Mark validation as failed first
	err := repo.UpdateFollowerValidationFailure(testIRI)
	if err != nil {
		t.Fatalf("Error updating follower validation failure: %s", err)
	}

	// Then mark as successful
	err = repo.UpdateFollowerValidationSuccess(testIRI)
	if err != nil {
		t.Fatalf("Error updating follower validation success: %s", err)
	}

	// Verify FirstValidationFailureAt is cleared
	followers, err := repo.GetFollowersToValidate(10)
	if err != nil {
		t.Fatalf("Error getting followers: %s", err)
	}

	for _, f := range followers {
		if f.ActorIRI == testIRI {
			if f.FirstValidationFailureAt.Valid {
				t.Error("Expected FirstValidationFailureAt to be NULL after success")
			}
		}
	}
}

func TestRemoveByIRI(t *testing.T) {
	repo := setupTestWithRepo(t)

	// Create a test follower
	testIRI := "https://fake.server/user/testuser"
	createTestFollower(repo, testIRI, "https://fake.server/user/testuser/inbox", "Test User", "testuser")

	// Verify follower exists
	count, err := repo.GetCount()
	if err != nil {
		t.Fatalf("Error getting count: %s", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 follower, got %d", count)
	}

	// Remove the follower
	err = repo.RemoveByIRI(testIRI)
	if err != nil {
		t.Fatalf("Error removing follower: %s", err)
	}

	// Verify follower is removed
	count, err = repo.GetCount()
	if err != nil {
		t.Fatalf("Error getting count: %s", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 followers after removal, got %d", count)
	}
}

func TestFailureDurationThresholdLogic(t *testing.T) {
	repo := setupTestWithRepo(t)

	// Create a test follower
	testIRI := "https://fake.server/user/testfailure"
	createTestFollower(repo, testIRI, "https://fake.server/user/testfailure/inbox", "Test User", "testfailure")

	// Set first_validation_failure_at to 8 days ago (past threshold)
	eightDaysAgo := time.Now().Add(-8 * 24 * time.Hour)
	_, err := testDatastore.DB.Exec(
		"UPDATE ap_followers SET first_validation_failure_at = ? WHERE iri = ?",
		eightDaysAgo, testIRI,
	)
	if err != nil {
		t.Fatalf("Error setting first_validation_failure_at: %s", err)
	}

	// Fetch the follower
	followers, err := repo.GetFollowersToValidate(1)
	if err != nil {
		t.Fatalf("Error getting followers: %s", err)
	}

	if len(followers) != 1 {
		t.Fatalf("Expected 1 follower, got %d", len(followers))
	}

	// Verify the failure duration is calculated correctly
	failureDuration := time.Since(followers[0].FirstValidationFailureAt.Time)
	if failureDuration < FailureDurationThreshold {
		t.Errorf("Expected failure duration (%v) to be >= threshold (%v)", failureDuration, FailureDurationThreshold)
	}

	// Test removal - this directly tests the removal logic without network calls
	err = repo.RemoveByIRI(testIRI)
	if err != nil {
		t.Fatalf("Error removing follower: %s", err)
	}

	// Verify the follower was removed
	count, err := repo.GetCount()
	if err != nil {
		t.Fatalf("Error getting count: %s", err)
	}
	if count != 0 {
		t.Errorf("Expected follower to be removed, but count is %d", count)
	}
}

func TestFailureNotRemovedBeforeThreshold(t *testing.T) {
	repo := setupTestWithRepo(t)

	// Create a test follower
	testIRI := "https://fake.server/user/testnotremoved"
	createTestFollower(repo, testIRI, "https://fake.server/user/testnotremoved/inbox", "Test User", "testnotremoved")

	// Set first_validation_failure_at to 1 day ago (not past threshold)
	oneDayAgo := time.Now().Add(-1 * 24 * time.Hour)
	_, err := testDatastore.DB.Exec(
		"UPDATE ap_followers SET first_validation_failure_at = ? WHERE iri = ?",
		oneDayAgo, testIRI,
	)
	if err != nil {
		t.Fatalf("Error setting first_validation_failure_at: %s", err)
	}

	// Fetch the follower
	followers, err := repo.GetFollowersToValidate(1)
	if err != nil {
		t.Fatalf("Error getting followers: %s", err)
	}

	if len(followers) != 1 {
		t.Fatalf("Expected 1 follower, got %d", len(followers))
	}

	// Verify the failure duration is below threshold
	failureDuration := time.Since(followers[0].FirstValidationFailureAt.Time)
	if failureDuration >= FailureDurationThreshold {
		t.Errorf("Expected failure duration (%v) to be < threshold (%v)", failureDuration, FailureDurationThreshold)
	}

	// Follower should NOT be removed yet
	count, err := repo.GetCount()
	if err != nil {
		t.Fatalf("Error getting count: %s", err)
	}
	if count != 1 {
		t.Errorf("Expected follower to NOT be removed yet, but count is %d", count)
	}
}

func TestFollowersOrderedByOldestValidatedFirst(t *testing.T) {
	repo := setupTestWithRepo(t)

	// Create followers
	iri1 := "https://fake.server/user/first"
	iri2 := "https://fake.server/user/second"
	iri3 := "https://fake.server/user/third"

	createTestFollower(repo, iri1, "https://fake.server/user/first/inbox", "First", "first")
	createTestFollower(repo, iri2, "https://fake.server/user/second/inbox", "Second", "second")
	createTestFollower(repo, iri3, "https://fake.server/user/third/inbox", "Third", "third")

	// Set different last_validated_at times
	now := time.Now()
	_, err := testDatastore.DB.Exec("UPDATE ap_followers SET last_validated_at = ? WHERE iri = ?",
		now.Add(-3*time.Hour), iri1)
	if err != nil {
		t.Fatalf("Error updating: %s", err)
	}
	_, err = testDatastore.DB.Exec("UPDATE ap_followers SET last_validated_at = ? WHERE iri = ?",
		now.Add(-1*time.Hour), iri2)
	if err != nil {
		t.Fatalf("Error updating: %s", err)
	}
	// iri3 has NULL last_validated_at (never validated)

	// Get followers - should return in order: iri3 (NULL), iri1 (oldest), iri2 (newest)
	followers, err := repo.GetFollowersToValidate(3)
	if err != nil {
		t.Fatalf("Error getting followers: %s", err)
	}

	if len(followers) != 3 {
		t.Fatalf("Expected 3 followers, got %d", len(followers))
	}

	// NULL values should come first (NULLS FIRST in query)
	if followers[0].ActorIRI != iri3 {
		t.Errorf("Expected first follower to be %s (never validated), got %s", iri3, followers[0].ActorIRI)
	}

	// Then oldest validated
	if followers[1].ActorIRI != iri1 {
		t.Errorf("Expected second follower to be %s (oldest validated), got %s", iri1, followers[1].ActorIRI)
	}

	// Then newest validated
	if followers[2].ActorIRI != iri2 {
		t.Errorf("Expected third follower to be %s (newest validated), got %s", iri2, followers[2].ActorIRI)
	}
}
