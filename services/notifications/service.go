package notifications

import (
	"fmt"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/notifications/browser"
	"github.com/owncast/owncast/notifications/discord"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/notificationsrepository"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Service defines the interface for notification operations.
type Service interface {
	Notify()
}

// notificationService handles notification dispatching and channel management.
type notificationService struct {
	repository       notificationsrepository.NotificationsRepository
	configRepository configrepository.ConfigRepository
	browser          *browser.Browser
	discord          *discord.Discord
}

// New creates a new instance of the notification service.
func New(datastore *data.Datastore) (Service, error) {
	service := &notificationService{
		repository:       notificationsrepository.New(datastore),
		configRepository: configrepository.Get(),
	}

	if err := service.setupBrowserPush(datastore); err != nil {
		log.Error(err)
	}
	if err := service.setupDiscord(); err != nil {
		log.Error(err)
	}

	return service, nil
}

func (s *notificationService) setupBrowserPush(datastore *data.Datastore) error {
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
		s.browser = browserNotifier
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
		s.discord = discordNotifier
	}
	return nil
}

func (s *notificationService) notifyBrowserPush() {
	destinations, err := s.repository.GetNotificationDestinationsForChannel(notificationsrepository.BrowserPushNotification)
	if err != nil {
		log.Errorln("error getting browser push notification destinations", err)
	}
	for _, destination := range destinations {
		unsubscribed, err := s.browser.Send(destination, s.configRepository.GetServerName(), s.configRepository.GetBrowserPushConfig().GoLiveMessage)
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

func (s *notificationService) notifyDiscord() {
	goLiveMessage := s.configRepository.GetDiscordConfig().GoLiveMessage
	streamTitle := s.configRepository.GetStreamTitle()
	if streamTitle != "" {
		goLiveMessage += "\n" + streamTitle
	}
	message := fmt.Sprintf("%s\n\n%s", goLiveMessage, s.configRepository.GetServerURL())

	if err := s.discord.Send(message); err != nil {
		log.Errorln("error sending discord message", err)
	}
}

// Notify will fire the different notification channels.
func (s *notificationService) Notify() {
	if s.browser != nil {
		s.notifyBrowserPush()
	}

	if s.discord != nil {
		s.notifyDiscord()
	}
}

// Compile-time verification that notificationService implements Service.
var _ Service = (*notificationService)(nil)
