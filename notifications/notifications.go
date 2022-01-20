package notifications

import (
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/notifications/browser"
	"github.com/owncast/owncast/notifications/discord"
	"github.com/owncast/owncast/notifications/email"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	smtpServer = "in-v3.mailjet.com"
	smtpPort   = "587"
	username   = "username"
	password   = "password"
)

// Notifier is an instance of the live stream notifier.
type Notifier struct {
	datastore *data.Datastore
	browser   *browser.Browser
	discord   *discord.Discord
	email     *email.Email
}

// Setup will perform any pre-use setup for the notifier.
func Setup(datastore *data.Datastore) {
	createNotificationsTable(datastore.DB)

	pubKey, _ := data.GetBrowserPushPublicKey()
	privKey, _ := data.GetBrowserPushPrivateKey()

	// We need browser push keys so people can register for pushes.
	if pubKey == "" || privKey == "" {
		browserPrivateKey, browserPublicKey, err := browser.GenerateBrowserPushKeys()
		if err != nil {
			log.Errorln("unable to initialize browser push notification keys", err)
		}

		if err := data.SetBrowserPushPrivateKey(browserPrivateKey); err != nil {
			log.Errorln("unable to set browser push private key", err)
		}

		if err := data.SetBrowserPushPublicKey(browserPublicKey); err != nil {
			log.Errorln("unable to set browser push public key", err)
		}
	}

	// Enable browser push notifications by default.
	if !data.GetHasPerformedInitialNotificationsConfig() {
		_ = data.SetBrowserPushConfig(models.BrowserNotificationConfiguration{Enabled: true, GoLiveMessage: config.GetDefaults().FederationGoLiveMessage})
		_ = data.SetHasPerformedInitialNotificationsConfig(true)
	}
}

// New creates a new instance of the Notifier.
func New(datastore *data.Datastore) (*Notifier, error) {
	notifier := Notifier{
		datastore: datastore,
	}

	// Add browser push notifications
	if data.GetBrowserPushConfig().Enabled {
		publicKey, err := data.GetBrowserPushPublicKey()
		if err != nil || publicKey == "" {
			return nil, errors.Wrap(err, "notifier disabled, failed to get browser push public key")
		}

		privateKey, err := data.GetBrowserPushPrivateKey()
		if err != nil || privateKey == "" {
			return nil, errors.Wrap(err, "notifier disabled, failed to get browser push private key")
		}

		browserNotifier, err := browser.New(datastore, publicKey, privateKey)
		if err != nil {
			return nil, errors.Wrap(err, "error creating browser notifier")
		}
		notifier.browser = browserNotifier
	}

	// Add discord notifications
	discordConfig := data.GetDiscordConfig()
	if discordConfig.Enabled && discordConfig.Webhook != "" {
		var image string
		if serverURL := data.GetServerURL(); serverURL != "" {
			image = serverURL + "/images/owncast-logo.png"
		}
		discordNotifier, err := discord.New(
			data.GetServerName(),
			image,
			discordConfig.Webhook,
		)
		if err != nil {
			return nil, errors.Wrap(err, "error creating discord notifier")
		}
		notifier.discord = discordNotifier
	}

	// Add email notifications
	if emailConfig := data.GetMailjetConfiguration(); emailConfig.Enabled && emailConfig.FromAddress != "" && emailConfig.ListAddress != "" {
		e := email.New(emailConfig.FromAddress, emailConfig.SMTPServer, "587", emailConfig.Username, emailConfig.Password)
		notifier.email = e
	}

	return &notifier, nil
}

func (n *Notifier) notifyBrowserDestinations() {
	destinations, err := GetNotificationDestinationsForChannel(BrowserPushNotification)
	if err != nil {
		log.Errorln("error getting browser push notification destinations", err)
	}
	for _, destination := range destinations {
		if err := n.browser.Send(destination, data.GetServerName(), data.GetBrowserPushConfig().GoLiveMessage); err != nil {
			log.Errorln(err)
		}
	}
}

func (n *Notifier) notifyDiscord() {
	if err := n.discord.Send(data.GetDiscordConfig().GoLiveMessage); err != nil {
		log.Errorln("error sending discord message", err)
	}
}

func (n *Notifier) notifyEmail() {
	content, err := email.GenerateEmailContent()
	if err != nil {
		log.Errorln("unable to generate email notification content: ", err)
		return
	}

	emailConfig := data.GetMailjetConfiguration()
	if !emailConfig.Enabled {
		return
	}

	subject := emailConfig.GoLiveSubject
	if data.GetStreamTitle() != "" {
		subject += " - " + data.GetStreamTitle()
	}

	if err := n.email.Send([]string{emailConfig.ListAddress}, content, subject); err != nil {
		log.Errorln("unable to send email notification: ", err)
		return
	}
}

// Notify will fire the different notification channels.
func (n *Notifier) Notify() {
	if n.browser != nil {
		n.notifyBrowserDestinations()
	}

	if n.discord != nil {
		n.notifyDiscord()
	}

	if n.email != nil {
		n.notifyEmail()
	}
}
