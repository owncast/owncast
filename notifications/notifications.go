package notifications

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/notifications/browser"
	"github.com/owncast/owncast/notifications/discord"
)

// Notifier is an instance of the live stream notifier.
type Notifier struct {
	data    *data.Service
	browser *browser.Browser
	discord *discord.Discord
}

// Setup will perform any pre-use setup for the notifier.
func (n *Notifier) Setup(datastore *data.Datastore) {
	n.createNotificationsTable(datastore.DB)
	n.initializeBrowserPushIfNeeded()
}

func (n *Notifier) initializeBrowserPushIfNeeded() {
	pubKey, _ := n.data.GetBrowserPushPublicKey()
	privKey, _ := n.data.GetBrowserPushPrivateKey()

	// We need browser push keys so people can register for pushes.
	if pubKey == "" || privKey == "" {
		browserPrivateKey, browserPublicKey, err := browser.GenerateBrowserPushKeys()
		if err != nil {
			log.Errorln("unable to initialize browser push notification keys", err)
		}

		if err := n.data.SetBrowserPushPrivateKey(browserPrivateKey); err != nil {
			log.Errorln("unable to set browser push private key", err)
		}

		if err := n.data.SetBrowserPushPublicKey(browserPublicKey); err != nil {
			log.Errorln("unable to set browser push public key", err)
		}
	}

	// Enable browser push notifications by default.
	if !n.data.GetHasPerformedInitialNotificationsConfig() {
		_ = n.data.SetBrowserPushConfig(models.BrowserNotificationConfiguration{Enabled: true, GoLiveMessage: config.GetDefaults().FederationGoLiveMessage})
		_ = n.data.SetHasPerformedInitialNotificationsConfig(true)
	}
}

// New creates a new instance of the Notifier.
func New(d *data.Service) (*Notifier, error) {
	notifier := Notifier{
		data: d,
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
	if n.data.GetBrowserPushConfig().Enabled {
		publicKey, err := n.data.GetBrowserPushPublicKey()
		if err != nil || publicKey == "" {
			return errors.Wrap(err, "browser notifier disabled, failed to get browser push public key")
		}

		privateKey, err := n.data.GetBrowserPushPrivateKey()
		if err != nil || privateKey == "" {
			return errors.Wrap(err, "browser notifier disabled, failed to get browser push private key")
		}

		browserNotifier, err := browser.New(n.data.Store, publicKey, privateKey)
		if err != nil {
			return errors.Wrap(err, "error creating browser notifier")
		}
		n.browser = browserNotifier
	}
	return nil
}

func (n *Notifier) notifyBrowserPush() {
	destinations, err := n.GetNotificationDestinationsForChannel(BrowserPushNotification)
	if err != nil {
		log.Errorln("error getting browser push notification destinations", err)
	}
	for _, destination := range destinations {
		unsubscribed, err := n.browser.Send(destination, n.data.GetServerName(), n.data.GetBrowserPushConfig().GoLiveMessage)
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

func (n *Notifier) setupDiscord() error {
	discordConfig := n.data.GetDiscordConfig()
	if discordConfig.Enabled && discordConfig.Webhook != "" {
		var image string
		if serverURL := n.data.GetServerURL(); serverURL != "" {
			image = serverURL + "/images/owncast-logo.png"
		}
		discordNotifier, err := discord.New(
			n.data.GetServerName(),
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
	goLiveMessage := n.data.GetDiscordConfig().GoLiveMessage
	streamTitle := n.data.GetStreamTitle()
	if streamTitle != "" {
		goLiveMessage += "\n" + streamTitle
	}
	message := fmt.Sprintf("%s\n\n%s", goLiveMessage, n.data.GetServerURL())

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
