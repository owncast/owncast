package notificationsrepository

import (
	"context"
	"fmt"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/tables"

	"github.com/owncast/owncast/notifications/browser"
	"github.com/owncast/owncast/notifications/discord"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type NotificationsRepository interface {
	Notify()
	AddNotification(channel, destination string) error
	RemoveNotificationForChannel(channel, destination string) error
	GetNotificationDestinationsForChannel(channel string) ([]string, error)
}

// SqlNotificationsRepository is an instance of the live stream notifier.
type SqlNotificationsRepository struct {
	datastore        *data.Datastore
	browser          *browser.Browser
	discord          *discord.Discord
	configRepository configrepository.ConfigRepository
}

// NOTE: This is temporary during the transition period.
var temporaryGlobalInstance NotificationsRepository

// Get will return the user repository.
func Get() NotificationsRepository {
	if temporaryGlobalInstance == nil {
		i := New(data.GetDatastore())
		temporaryGlobalInstance = i
	}
	return temporaryGlobalInstance
}

// Setup will perform any pre-use setup for the notifier.
func Setup(datastore *data.Datastore) {
	tables.CreateNotificationsTable(datastore.DB)
	initializeBrowserPushIfNeeded()
}

func initializeBrowserPushIfNeeded() {
	configRepository := configrepository.Get()

	pubKey, _ := configRepository.GetBrowserPushPublicKey()
	privKey, _ := configRepository.GetBrowserPushPrivateKey()

	// We need browser push keys so people can register for pushes.
	if pubKey == "" || privKey == "" {
		browserPrivateKey, browserPublicKey, err := browser.GenerateBrowserPushKeys()
		if err != nil {
			log.Errorln("unable to initialize browser push notification keys", err)
		}

		if err := configRepository.SetBrowserPushPrivateKey(browserPrivateKey); err != nil {
			log.Errorln("unable to set browser push private key", err)
		}

		if err := configRepository.SetBrowserPushPublicKey(browserPublicKey); err != nil {
			log.Errorln("unable to set browser push public key", err)
		}
	}

	// Enable browser push notifications by default.
	if !configRepository.GetHasPerformedInitialNotificationsConfig() {
		_ = configRepository.SetBrowserPushConfig(models.BrowserNotificationConfiguration{Enabled: true, GoLiveMessage: config.GetDefaults().FederationGoLiveMessage})
		_ = configRepository.SetHasPerformedInitialNotificationsConfig(true)
	}
}

// New creates a new instance of the Notifier.
func New(datastore *data.Datastore) NotificationsRepository {
	notifier := SqlNotificationsRepository{
		datastore:        datastore,
		configRepository: configrepository.Get(),
	}

	if err := notifier.setupBrowserPush(); err != nil {
		log.Error(err)
	}
	if err := notifier.setupDiscord(); err != nil {
		log.Error(err)
	}

	return &notifier
}

func (n *SqlNotificationsRepository) setupBrowserPush() error {
	if n.configRepository.GetBrowserPushConfig().Enabled {
		publicKey, err := n.configRepository.GetBrowserPushPublicKey()
		if err != nil || publicKey == "" {
			return errors.Wrap(err, "browser notifier disabled, failed to get browser push public key")
		}

		privateKey, err := n.configRepository.GetBrowserPushPrivateKey()
		if err != nil || privateKey == "" {
			return errors.Wrap(err, "browser notifier disabled, failed to get browser push private key")
		}

		browserNotifier, err := browser.New(n.datastore, publicKey, privateKey)
		if err != nil {
			return errors.Wrap(err, "error creating browser notifier")
		}
		n.browser = browserNotifier
	}
	return nil
}

func (n *SqlNotificationsRepository) notifyBrowserPush() {
	destinations, err := n.GetNotificationDestinationsForChannel(BrowserPushNotification)
	if err != nil {
		log.Errorln("error getting browser push notification destinations", err)
	}
	for _, destination := range destinations {
		unsubscribed, err := n.browser.Send(destination, n.configRepository.GetServerName(), n.configRepository.GetBrowserPushConfig().GoLiveMessage)
		if unsubscribed {
			// If the error is "unsubscribed", then remove the destination from the database.
			if err := n.RemoveNotificationForChannel(BrowserPushNotification, destination); err != nil {
				log.Errorln(err)
			}
		} else if err != nil {
			log.Errorln(err)
		}
	}
}

func (n *SqlNotificationsRepository) setupDiscord() error {
	discordConfig := n.configRepository.GetDiscordConfig()
	if discordConfig.Enabled && discordConfig.Webhook != "" {
		var image string
		if serverURL := n.configRepository.GetServerURL(); serverURL != "" {
			image = serverURL + "/logo"
		}
		discordNotifier, err := discord.New(
			n.configRepository.GetServerName(),
			image,
			discordConfig.Webhook,
		)
		if err != nil {
			return errors.Wrap(err, "error creating discord notifier")
		}
		n.discord = discordNotifier
	}
	return nil
}

func (n *SqlNotificationsRepository) notifyDiscord() {
	goLiveMessage := n.configRepository.GetDiscordConfig().GoLiveMessage
	streamTitle := n.configRepository.GetStreamTitle()
	if streamTitle != "" {
		goLiveMessage += "\n" + streamTitle
	}
	message := fmt.Sprintf("%s\n\n%s", goLiveMessage, n.configRepository.GetServerURL())

	if err := n.discord.Send(message); err != nil {
		log.Errorln("error sending discord message", err)
	}
}

// Notify will fire the different notification channels.
func (n *SqlNotificationsRepository) Notify() {
	if n.browser != nil {
		n.notifyBrowserPush()
	}

	if n.discord != nil {
		n.notifyDiscord()
	}
}

// AddNotification saves a new user notification destination.
func (n *SqlNotificationsRepository) AddNotification(channel, destination string) error {
	return n.datastore.GetQueries().AddNotification(context.Background(), db.AddNotificationParams{
		Channel:     channel,
		Destination: destination,
	})
}

// RemoveNotificationForChannel removes a notification destination.
func (n *SqlNotificationsRepository) RemoveNotificationForChannel(channel, destination string) error {
	log.Debugln("Removing notification for channel", channel)

	return n.datastore.GetQueries().RemoveNotificationDestinationForChannel(context.Background(), db.RemoveNotificationDestinationForChannelParams{
		Channel:     channel,
		Destination: destination,
	})
}

// GetNotificationDestinationsForChannel will return a collection of
// destinations to notify for a given channel.
func (n *SqlNotificationsRepository) GetNotificationDestinationsForChannel(channel string) ([]string, error) {
	result, err := n.datastore.GetQueries().GetNotificationDestinationsForChannel(context.Background(), channel)
	if err != nil {
		return nil, errors.Wrap(err, "unable to query notification destinations for channel "+channel)
	}

	return result, nil
}
