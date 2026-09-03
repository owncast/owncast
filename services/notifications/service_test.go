package notifications

import (
	"errors"
	"testing"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/notificationsrepository"
)

type testNotificationRepository struct {
	notificationsrepository.NotificationsRepository
	destinations []string
	removed      []string
}

func (r *testNotificationRepository) GetNotificationDestinationsForChannel(string) ([]string, error) {
	return r.destinations, nil
}

func (r *testNotificationRepository) RemoveNotificationForChannel(_ string, destination string) error {
	r.removed = append(r.removed, destination)
	return nil
}

type testNotificationConfig struct {
	configrepository.ConfigRepository
}

func (testNotificationConfig) GetServerName() string { return "Owncast Test" }
func (testNotificationConfig) GetServerURL() string  { return "https://owncast.example" }
func (testNotificationConfig) GetStreamTitle() string {
	return "Current stream"
}

func (testNotificationConfig) GetBrowserPushConfig() models.BrowserNotificationConfiguration {
	return models.BrowserNotificationConfiguration{Enabled: true, GoLiveMessage: "Browser live"}
}

func (testNotificationConfig) GetDiscordConfig() models.DiscordConfiguration {
	return models.DiscordConfiguration{Enabled: true, GoLiveMessage: "Discord live"}
}

func TestNotifyScheduledEventSendsBrowserAndDiscord(t *testing.T) {
	repository := &testNotificationRepository{destinations: []string{"subscribed", "expired"}}
	var browserMessages []string
	var discordMessages []string
	service := &notificationService{
		repository:       repository,
		configRepository: testNotificationConfig{},
		sendBrowser: func(destination, title, body string) (bool, error) {
			if title != "Owncast Test" {
				t.Errorf("browser title = %q", title)
			}
			browserMessages = append(browserMessages, destination+":"+body)
			return destination == "expired", nil
		},
		sendDiscord: func(message string) error {
			discordMessages = append(discordMessages, message)
			return nil
		},
	}

	service.NotifyScheduledEvent("Event reminder")

	if len(browserMessages) != 2 || browserMessages[0] != "subscribed:Event reminder" || browserMessages[1] != "expired:Event reminder" {
		t.Errorf("browser messages = %#v", browserMessages)
	}
	if len(discordMessages) != 1 || discordMessages[0] != "Event reminder" {
		t.Errorf("discord messages = %#v", discordMessages)
	}
	if len(repository.removed) != 1 || repository.removed[0] != "expired" {
		t.Errorf("removed subscriptions = %#v", repository.removed)
	}
}

func TestNotifyScheduledEventContinuesAfterBrowserError(t *testing.T) {
	repository := &testNotificationRepository{destinations: []string{"broken", "working"}}
	browserCalls := 0
	discordCalled := false
	service := &notificationService{
		repository:       repository,
		configRepository: testNotificationConfig{},
		sendBrowser: func(destination, _, _ string) (bool, error) {
			browserCalls++
			if destination == "broken" {
				return false, errors.New("delivery failed")
			}
			return false, nil
		},
		sendDiscord: func(string) error {
			discordCalled = true
			return nil
		},
	}

	service.NotifyScheduledEvent("Event reminder")

	if browserCalls != 2 {
		t.Errorf("browser calls = %d, want 2", browserCalls)
	}
	if !discordCalled {
		t.Error("Discord was not notified after a browser delivery error")
	}
}

func TestNotifyKeepsGoLiveMessagesChannelSpecific(t *testing.T) {
	repository := &testNotificationRepository{destinations: []string{"browser"}}
	var browserMessage string
	var discordMessage string
	service := &notificationService{
		repository:       repository,
		configRepository: testNotificationConfig{},
		sendBrowser: func(_, _, body string) (bool, error) {
			browserMessage = body
			return false, nil
		},
		sendDiscord: func(message string) error {
			discordMessage = message
			return nil
		},
	}

	service.Notify()

	if browserMessage != "Browser live" {
		t.Errorf("browser go-live message = %q", browserMessage)
	}
	if discordMessage != "Discord live\nCurrent stream\n\nhttps://owncast.example" {
		t.Errorf("Discord go-live message = %q", discordMessage)
	}
}
