package notifications

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/notifications/browser"
	"github.com/owncast/owncast/notifications/discord"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/notificationsrepository"
	"github.com/owncast/owncast/services/datastore"
)

// Service defines the interface for notification operations.
type Service interface {
	Notify()
	NotifyScheduledEvent(message string)
}

// notificationService handles notification dispatching and channel management.
type notificationService struct {
	repository       notificationsrepository.NotificationsRepository
	configRepository configrepository.ConfigRepository
	sendBrowser      func(subscription, title, body string) (bool, error)
	sendDiscord      func(message string) error
}

// New creates a new instance of the notification service.
func New(datastore *datastore.Datastore, configRepository configrepository.ConfigRepository) (Service, error) {
	service := &notificationService{
		repository:       notificationsrepository.New(datastore, configRepository),
		configRepository: configRepository,
	}

	if err := service.setupBrowserPush(datastore); err != nil {
		log.Error(err)
	}
	if err := service.setupDiscord(); err != nil {
		log.Error(err)
	}

	return service, nil
}

func (s *notificationService) setupBrowserPush(datastore *datastore.Datastore) error {
	if s.configRepository.GetBrowserPushConfig().Enabled {
		publicKey, err := s.configRepository.GetBrowserPushPublicKey()
		if err != nil || publicKey == "" {
			return errors.Wrap(err, "browser notifier disabled, failed to get browser push public key")
		}

		privateKey, err := s.configRepository.GetBrowserPushPrivateKey()
		if err != nil || privateKey == "" {
			return errors.Wrap(err, "browser notifier disabled, failed to get browser push private key")
		}

		browserNotifier, err := browser.New(datastore, publicKey, privateKey)
		if err != nil {
			return errors.Wrap(err, "error creating browser notifier")
		}
		s.sendBrowser = browserNotifier.Send
	}
	return nil
}

func (s *notificationService) setupDiscord() error {
	discordConfig := s.configRepository.GetDiscordConfig()
	if discordConfig.Enabled && discordConfig.Webhook != "" {
		var image string
		if serverURL := s.configRepository.GetServerURL(); serverURL != "" {
			image = serverURL + "/logo"
		}
		discordNotifier, err := discord.New(
			s.configRepository.GetServerName(),
			image,
			discordConfig.Webhook,
		)
		if err != nil {
			return errors.Wrap(err, "error creating discord notifier")
		}
		s.sendDiscord = discordNotifier.Send
	}
	return nil
}

func (s *notificationService) notifyBrowserPush(body string) {
	destinations, err := s.repository.GetNotificationDestinationsForChannel(notificationsrepository.BrowserPushNotification)
	if err != nil {
		log.Errorln("error getting browser push notification destinations", err)
	}
	for _, destination := range destinations {
		unsubscribed, err := s.sendBrowser(destination, s.configRepository.GetServerName(), body)
		if unsubscribed {
			// If the error is "unsubscribed", then remove the destination from the database.
			if err := s.repository.RemoveNotificationForChannel(notificationsrepository.BrowserPushNotification, destination); err != nil {
				log.Errorln(err)
			}
		} else if err != nil {
			log.Errorln(err)
		}
	}
}

func (s *notificationService) notifyDiscord(message string) {
	if err := s.sendDiscord(message); err != nil {
		log.Errorln("error sending discord message", err)
	}
}

func (s *notificationService) goLiveMessage() string {
	message := s.configRepository.GetDiscordConfig().GoLiveMessage
	streamTitle := s.configRepository.GetStreamTitle()
	if streamTitle != "" {
		message += "\n" + streamTitle
	}
	return fmt.Sprintf("%s\n\n%s", message, s.configRepository.GetServerURL())
}

// Notify sends the go-live message to each configured notification channel.
func (s *notificationService) Notify() {
	if s.sendBrowser != nil {
		s.notifyBrowserPush(s.configRepository.GetBrowserPushConfig().GoLiveMessage)
	}
	if s.sendDiscord != nil {
		s.notifyDiscord(s.goLiveMessage())
	}
}

// NotifyScheduledEvent sends a fully formatted scheduled-event reminder.
func (s *notificationService) NotifyScheduledEvent(message string) {
	if s.sendBrowser != nil {
		s.notifyBrowserPush(message)
	}
	if s.sendDiscord != nil {
		s.notifyDiscord(message)
	}
}

// Compile-time verification that notificationService implements Service.
var _ Service = (*notificationService)(nil)
