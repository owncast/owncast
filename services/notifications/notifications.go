package notifications

import (
	"fmt"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/config"
	"github.com/owncast/owncast/services/notifications/browser"
	"github.com/owncast/owncast/services/notifications/discord"
	"github.com/owncast/owncast/storage/configrepository"
	"github.com/owncast/owncast/storage/data"
	"github.com/owncast/owncast/storage/notificationsrepository"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Notifier is an instance of the live stream notifier.
type Notifier struct {
	datastore *data.Store
	browser   *browser.Browser
	discord   *discord.Discord
}

var (
	configRepository        *configrepository.SqlConfigRepository
	notificationsRepository *notificationsrepository.SqlNotificationsRepository
)

// Setup will perform any pre-use setup for the notifier.
func Setup(datastore *data.Store) {
	configRepository = configrepository.Get()
	notificationsRepository = notificationsrepository.Get()

	initializeBrowserPushIfNeeded()
}

func initializeBrowserPushIfNeeded() {
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
func New(datastore *data.Store) (*Notifier, error) {
	notifier := Notifier{
		datastore: datastore,
	}

	if err := notifier.setupBrowserPush(); err != nil {
		log.Error(err)
	}
	if err := notifier.setupDiscord(); err != nil {
		log.Error(err)
	}

	return &notifier, nil
}

func (n *Notifier) setupBrowserPush() error {
	if configRepository.GetBrowserPushConfig().Enabled {
		publicKey, err := configRepository.GetBrowserPushPublicKey()
		if err != nil || publicKey == "" {
			return errors.Wrap(err, "browser notifier disabled, failed to get browser push public key")
		}

		privateKey, err := configRepository.GetBrowserPushPrivateKey()
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

func (n *Notifier) notifyBrowserPush() {
	destinations, err := notificationsRepository.GetNotificationDestinationsForChannel(BrowserPushNotification)
	if err != nil {
		log.Errorln("error getting browser push notification destinations", err)
	}
	for _, destination := range destinations {
		unsubscribed, err := n.browser.Send(destination, configRepository.GetServerName(), configRepository.GetBrowserPushConfig().GoLiveMessage)
		if unsubscribed {
			// If the error is "unsubscribed", then remove the destination from the database.
			if err := notificationsRepository.RemoveNotificationForChannel(BrowserPushNotification, destination); err != nil {
				log.Errorln(err)
			}
		} else if err != nil {
			log.Errorln(err)
		}
	}
}

func (n *Notifier) setupDiscord() error {
	discordConfig := configRepository.GetDiscordConfig()
	if discordConfig.Enabled && discordConfig.Webhook != "" {
		var image string
		if serverURL := configRepository.GetServerURL(); serverURL != "" {
			image = serverURL + "/logo"
		}
		discordNotifier, err := discord.New(
			configRepository.GetServerName(),
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

func (n *Notifier) notifyDiscord() {
	goLiveMessage := configRepository.GetDiscordConfig().GoLiveMessage
	streamTitle := configRepository.GetStreamTitle()
	if streamTitle != "" {
		goLiveMessage += "\n" + streamTitle
	}
	message := fmt.Sprintf("%s\n\n%s", goLiveMessage, configRepository.GetServerURL())

	if err := n.discord.Send(message); err != nil {
		log.Errorln("error sending discord message", err)
	}
}

// Notify will fire the different notification channels.
func (n *Notifier) Notify() {
	if n.browser != nil {
		n.notifyBrowserPush()
	}

	if n.discord != nil {
		n.notifyDiscord()
	}
}
