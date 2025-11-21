package notificationsrepository

import (
	"fmt"
	"testing"

	"github.com/owncast/owncast/core/data"
)

var (
	testDatastore *data.Datastore
	testRepo      NotificationsRepository
)

func TestMain(m *testing.M) {
	// Create an in-memory database for testing
	if err := data.SetupPersistence(":memory:"); err != nil {
		panic(err)
	}

	// Get the shared datastore instance
	testDatastore = data.GetDatastore()

	// Setup the notifications repository
	Setup(testDatastore)
	testRepo = New(testDatastore)

	// Run tests
	m.Run()
}

func TestAddNotification(t *testing.T) {
	channel := "TEST_CHANNEL_ADD_UNIQUE"
	destination := "test-destination-add-unique-123"

	// Test adding a notification
	err := testRepo.AddNotification(channel, destination)
	if err != nil {
		t.Errorf("Failed to add notification: %v", err)
	}

	// Verify the notification was added
	destinations, err := testRepo.GetNotificationDestinationsForChannel(channel)
	if err != nil {
		t.Errorf("Failed to get notification destinations: %v", err)
	}

	if len(destinations) != 1 {
		t.Errorf("Expected 1 destination, got %d", len(destinations))
	}

	if destinations[0] != destination {
		t.Errorf("Expected destination %s, got %s", destination, destinations[0])
	}
}

func TestAddMultipleNotifications(t *testing.T) {
	channel := "MULTI_TEST_CHANNEL"
	destinations := []string{"dest-1", "dest-2", "dest-3"}

	// Add multiple notifications
	for _, dest := range destinations {
		err := testRepo.AddNotification(channel, dest)
		if err != nil {
			t.Errorf("Failed to add notification for destination %s: %v", dest, err)
		}
	}

	// Verify all notifications were added
	retrievedDests, err := testRepo.GetNotificationDestinationsForChannel(channel)
	if err != nil {
		t.Errorf("Failed to get notification destinations: %v", err)
	}

	if len(retrievedDests) != len(destinations) {
		t.Errorf("Expected %d destinations, got %d", len(destinations), len(retrievedDests))
	}

	// Check that all destinations are present
	for _, expectedDest := range destinations {
		found := false
		for _, retrievedDest := range retrievedDests {
			if retrievedDest == expectedDest {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected destination %s not found in retrieved destinations", expectedDest)
		}
	}
}

func TestRemoveNotificationForChannel(t *testing.T) {
	channel := "REMOVE_TEST_CHANNEL"
	destination1 := "remove-dest-1"
	destination2 := "remove-dest-2"

	// Add two notifications
	err := testRepo.AddNotification(channel, destination1)
	if err != nil {
		t.Errorf("Failed to add notification: %v", err)
	}

	err = testRepo.AddNotification(channel, destination2)
	if err != nil {
		t.Errorf("Failed to add notification: %v", err)
	}

	// Verify both were added
	destinations, err := testRepo.GetNotificationDestinationsForChannel(channel)
	if err != nil {
		t.Errorf("Failed to get notification destinations: %v", err)
	}

	if len(destinations) != 2 {
		t.Errorf("Expected 2 destinations before removal, got %d", len(destinations))
	}

	// Remove one notification
	err = testRepo.RemoveNotificationForChannel(channel, destination1)
	if err != nil {
		t.Errorf("Failed to remove notification: %v", err)
	}

	// Verify only one remains
	destinations, err = testRepo.GetNotificationDestinationsForChannel(channel)
	if err != nil {
		t.Errorf("Failed to get notification destinations after removal: %v", err)
	}

	if len(destinations) != 1 {
		t.Errorf("Expected 1 destination after removal, got %d", len(destinations))
	}

	if destinations[0] != destination2 {
		t.Errorf("Expected remaining destination %s, got %s", destination2, destinations[0])
	}
}

func TestGetNotificationDestinationsForNonExistentChannel(t *testing.T) {
	channel := "NON_EXISTENT_CHANNEL"

	destinations, err := testRepo.GetNotificationDestinationsForChannel(channel)
	if err != nil {
		t.Errorf("Failed to get notification destinations for non-existent channel: %v", err)
	}

	if len(destinations) != 0 {
		t.Errorf("Expected 0 destinations for non-existent channel, got %d", len(destinations))
	}
}

func TestRemoveNonExistentNotification(t *testing.T) {
	channel := "NON_EXISTENT_REMOVE_CHANNEL"
	destination := "non-existent-destination"

	// Try to remove a notification that doesn't exist
	err := testRepo.RemoveNotificationForChannel(channel, destination)
	if err != nil {
		t.Errorf("Removing non-existent notification should not return error, got: %v", err)
	}
}

func TestBrowserPushNotificationConstant(t *testing.T) {
	// Test that the constant is defined and has the expected value
	if BrowserPushNotification == "" {
		t.Error("BrowserPushNotification constant should not be empty")
	}

	expectedValue := "BROWSER_PUSH_NOTIFICATION"
	if BrowserPushNotification != expectedValue {
		t.Errorf("Expected BrowserPushNotification to be %s, got %s", expectedValue, BrowserPushNotification)
	}
}

func TestNotificationRepositoryInterface(t *testing.T) {
	// Test that our implementation satisfies the interface
	var _ NotificationsRepository = &SqlNotificationsRepository{}

	// Test that we can get the repository instance
	repo := Get()
	if repo == nil {
		t.Error("Get() should return a non-nil repository instance")
	}

	// Test that New creates a valid repository
	newRepo := New(testDatastore)
	if newRepo == nil {
		t.Error("New() should return a non-nil repository instance")
	}
}

func TestAddDuplicateNotification(t *testing.T) {
	channel := "DUPLICATE_TEST_CHANNEL"
	destination := "duplicate-dest"

	// Add the same notification twice
	err := testRepo.AddNotification(channel, destination)
	if err != nil {
		t.Errorf("Failed to add notification first time: %v", err)
	}

	err = testRepo.AddNotification(channel, destination)
	if err != nil {
		t.Errorf("Failed to add duplicate notification: %v", err)
	}

	// Check how many destinations we have (should handle duplicates gracefully)
	destinations, err := testRepo.GetNotificationDestinationsForChannel(channel)
	if err != nil {
		t.Errorf("Failed to get notification destinations: %v", err)
	}

	// The behavior for duplicates depends on database constraints
	// This test documents the current behavior
	if len(destinations) == 0 {
		t.Error("Should have at least one destination even with duplicates")
	}
}

func TestChannelIsolation(t *testing.T) {
	channel1 := "CHANNEL_1"
	channel2 := "CHANNEL_2"
	destination1 := "dest-for-channel-1"
	destination2 := "dest-for-channel-2"

	// Add notifications to different channels
	err := testRepo.AddNotification(channel1, destination1)
	if err != nil {
		t.Errorf("Failed to add notification to channel1: %v", err)
	}

	err = testRepo.AddNotification(channel2, destination2)
	if err != nil {
		t.Errorf("Failed to add notification to channel2: %v", err)
	}

	// Verify channel isolation
	destinations1, err := testRepo.GetNotificationDestinationsForChannel(channel1)
	if err != nil {
		t.Errorf("Failed to get destinations for channel1: %v", err)
	}

	destinations2, err := testRepo.GetNotificationDestinationsForChannel(channel2)
	if err != nil {
		t.Errorf("Failed to get destinations for channel2: %v", err)
	}

	// Each channel should only have its own destination
	if len(destinations1) != 1 || destinations1[0] != destination1 {
		t.Errorf("Channel1 should only have destination1, got: %v", destinations1)
	}

	if len(destinations2) != 1 || destinations2[0] != destination2 {
		t.Errorf("Channel2 should only have destination2, got: %v", destinations2)
	}

	// Remove from one channel shouldn't affect the other
	err = testRepo.RemoveNotificationForChannel(channel1, destination1)
	if err != nil {
		t.Errorf("Failed to remove notification from channel1: %v", err)
	}

	// Verify channel2 is unaffected
	destinations2After, err := testRepo.GetNotificationDestinationsForChannel(channel2)
	if err != nil {
		t.Errorf("Failed to get destinations for channel2 after removal: %v", err)
	}

	if len(destinations2After) != 1 || destinations2After[0] != destination2 {
		t.Errorf("Channel2 should still have destination2 after channel1 removal, got: %v", destinations2After)
	}
}

// Benchmark tests
func BenchmarkAddNotification(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		channel := fmt.Sprintf("BENCH_CHANNEL_%d", i)
		destination := fmt.Sprintf("bench_destination_%d", i)
		err := testRepo.AddNotification(channel, destination)
		if err != nil {
			b.Fatalf("Failed to add notification: %v", err)
		}
	}
}

func BenchmarkGetNotificationDestinationsForChannel(b *testing.B) {
	// Setup test data
	channel := "BENCH_GET_CHANNEL"
	for i := 0; i < 100; i++ {
		destination := fmt.Sprintf("bench_get_destination_%d", i)
		testRepo.AddNotification(channel, destination)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := testRepo.GetNotificationDestinationsForChannel(channel)
		if err != nil {
			b.Fatalf("Failed to get notification destinations: %v", err)
		}
	}
}

func BenchmarkRemoveNotificationForChannel(b *testing.B) {
	// Setup test data
	channel := "BENCH_REMOVE_CHANNEL"
	destinations := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		destination := fmt.Sprintf("bench_remove_destination_%d", i)
		destinations[i] = destination
		testRepo.AddNotification(channel, destination)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := testRepo.RemoveNotificationForChannel(channel, destinations[i])
		if err != nil {
			b.Fatalf("Failed to remove notification: %v", err)
		}
	}
}
