package notificationsrepository

import (
	"testing"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/datastore"
)

var (
	integrationTestDatastore *datastore.Datastore
	integrationRepo          NotificationsRepository
)

func init() {
	// Use the shared test datastore from the main test file
	// This ensures we're using the same database setup
	integrationTestDatastore = testDatastore
	if integrationTestDatastore != nil {
		integrationRepo = New(integrationTestDatastore, configrepository.New(integrationTestDatastore))
	}
}

func TestIntegrationSetup(t *testing.T) {
	// This test ensures integration setup works properly
	// The actual setup is done in the init function below
	if integrationTestDatastore == nil {
		integrationTestDatastore = testDatastore
		integrationRepo = New(integrationTestDatastore, configrepository.New(integrationTestDatastore))
		integrationRepo.Setup()
	}

	if integrationRepo == nil {
		t.Error("Integration repository should be initialized")
	}
}

func TestBrowserPushSetupIntegration(t *testing.T) {
	// Test that browser push keys are generated during setup
	configRepository := configrepository.New(testDatastore)
	pubKey, err := configRepository.GetBrowserPushPublicKey()
	if err != nil {
		t.Errorf("Should be able to get browser push public key: %v", err)
	}

	privKey, err := configRepository.GetBrowserPushPrivateKey()
	if err != nil {
		t.Errorf("Should be able to get browser push private key: %v", err)
	}

	// Keys should be generated during setup
	if pubKey == "" || privKey == "" {
		t.Error("Browser push keys should be generated during setup")
	}
}

func TestBrowserPushConfigurationIntegration(t *testing.T) {
	configRepository := configrepository.New(testDatastore)

	// Test that browser push is enabled by default
	browserConfig := configRepository.GetBrowserPushConfig()
	if !browserConfig.Enabled {
		t.Error("Browser push should be enabled by default")
	}

	// Test that initial notification configuration flag is set
	hasConfigured := configRepository.GetHasPerformedInitialNotificationsConfig()
	if !hasConfigured {
		t.Error("Should have performed initial notifications configuration")
	}
}

func TestDiscordConfigurationIntegration(t *testing.T) {
	configRepository := configrepository.New(testDatastore)

	// Test Discord configuration setup
	discordConfig := models.DiscordConfiguration{
		Enabled:       true,
		Webhook:       "https://discord.com/api/webhooks/test",
		GoLiveMessage: "Test stream is live!",
	}

	err := configRepository.SetDiscordConfig(discordConfig)
	if err != nil {
		t.Errorf("Failed to set Discord configuration: %v", err)
	}

	// Test with disabled Discord
	disabledConfig := models.DiscordConfiguration{
		Enabled: false,
		Webhook: "",
	}
	err = configRepository.SetDiscordConfig(disabledConfig)
	if err != nil {
		t.Errorf("Failed to disable Discord configuration: %v", err)
	}
}

func TestNotificationWorkflowIntegration(t *testing.T) {
	channel := BrowserPushNotification
	destination := "integration-test-endpoint"

	// Add a notification
	err := integrationRepo.AddNotification(channel, destination)
	if err != nil {
		t.Errorf("Failed to add notification: %v", err)
	}

	// Verify it was added
	destinations, err := integrationRepo.GetNotificationDestinationsForChannel(channel)
	if err != nil {
		t.Errorf("Failed to get notification destinations: %v", err)
	}

	found := false
	for _, dest := range destinations {
		if dest == destination {
			found = true
			break
		}
	}

	if !found {
		t.Error("Notification destination should be found after adding")
	}

	// Clean up
	err = integrationRepo.RemoveNotificationForChannel(channel, destination)
	if err != nil {
		t.Errorf("Failed to remove notification: %v", err)
	}

	// Verify removal
	destinationsAfter, err := integrationRepo.GetNotificationDestinationsForChannel(channel)
	if err != nil {
		t.Errorf("Failed to get notification destinations after removal: %v", err)
	}

	for _, dest := range destinationsAfter {
		if dest == destination {
			t.Error("Notification destination should be removed")
		}
	}
}

func TestDatabaseTransactionIntegrity(t *testing.T) {
	channel := "TRANSACTION_TEST_CHANNEL"
	destinations := []string{"dest1", "dest2", "dest3"}

	// Add multiple notifications
	for _, dest := range destinations {
		err := integrationRepo.AddNotification(channel, dest)
		if err != nil {
			t.Errorf("Failed to add notification %s: %v", dest, err)
		}
	}

	// Verify all were added
	retrievedDests, err := integrationRepo.GetNotificationDestinationsForChannel(channel)
	if err != nil {
		t.Errorf("Failed to retrieve destinations: %v", err)
	}

	if len(retrievedDests) != len(destinations) {
		t.Errorf("Expected %d destinations, got %d", len(destinations), len(retrievedDests))
	}

	// Remove them one by one and verify state consistency
	for i, dest := range destinations {
		err := integrationRepo.RemoveNotificationForChannel(channel, dest)
		if err != nil {
			t.Errorf("Failed to remove notification %s: %v", dest, err)
		}

		remaining, err := integrationRepo.GetNotificationDestinationsForChannel(channel)
		if err != nil {
			t.Errorf("Failed to get remaining destinations: %v", err)
		}

		expectedRemaining := len(destinations) - (i + 1)
		if len(remaining) != expectedRemaining {
			t.Errorf("Expected %d remaining destinations, got %d", expectedRemaining, len(remaining))
		}
	}
}

func TestErrorHandling(t *testing.T) {
	// Test with empty channel
	err := integrationRepo.AddNotification("", "test-destination")
	if err != nil {
		// Empty channel might be allowed - this documents the behavior
		t.Logf("Empty channel behavior: %v", err)
	}

	// Test with empty destination
	err = integrationRepo.AddNotification("TEST_CHANNEL", "")
	if err != nil {
		// Empty destination might be allowed - this documents the behavior
		t.Logf("Empty destination behavior: %v", err)
	}

	// Test removing from empty channel
	err = integrationRepo.RemoveNotificationForChannel("", "test-destination")
	if err != nil {
		t.Logf("Remove from empty channel behavior: %v", err)
	}

	// Test getting destinations for empty channel
	destinations, err := integrationRepo.GetNotificationDestinationsForChannel("")
	if err != nil {
		t.Errorf("Getting destinations for empty channel should not error: %v", err)
	}

	if len(destinations) > 0 {
		t.Logf("Empty channel returned %d destinations", len(destinations))
	}
}
