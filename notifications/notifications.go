package notifications

import (
	"fmt"
	"strings"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/notifications/browser"
	"github.com/owncast/owncast/notifications/discord"
	"github.com/owncast/owncast/notifications/email"
	"github.com/owncast/owncast/notifications/twitter"
	"github.com/owncast/owncast/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Notifier is an instance of the live stream notifier.
type Notifier struct {
	datastore *data.Datastore
	browser   *browser.Browser
	discord   *discord.Discord
	email     *email.Email
	twitter   *twitter.Twitter
}

// Setup will perform any pre-use setup for the notifier.
func Setup(datastore *data.Datastore) {
	createNotificationsTable(datastore.DB)
	initializeBrowserPushIfNeeded()
}

func initializeBrowserPushIfNeeded() {
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

	if err := notifier.setupBrowserPush(); err != nil {
		log.Error(err)
	}
	if err := notifier.setupDiscord(); err != nil {
		log.Error(err)
	}
	if err := notifier.setupEmail(); err != nil {
		log.Error(err)
	}
	if err := notifier.setupTwitter(); err != nil {
		log.Errorln(err)
	}

	return &notifier, nil
}

func (n *Notifier) setupBrowserPush() error {
	if data.GetBrowserPushConfig().Enabled {
		publicKey, err := data.GetBrowserPushPublicKey()
		if err != nil || publicKey == "" {
			return errors.Wrap(err, "browser notifier disabled, failed to get browser push public key")
		}

		privateKey, err := data.GetBrowserPushPrivateKey()
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

func (n *Notifier) setupDiscord() error {
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
			return errors.Wrap(err, "error creating discord notifier")
		}
		n.discord = discordNotifier
	}
	return nil
}

func (n *Notifier) notifyDiscord() {
	goLiveMessage := data.GetDiscordConfig().GoLiveMessage
	streamTitle := data.GetStreamTitle()
	if streamTitle != "" {
		goLiveMessage += "\n" + streamTitle
	}
	message := fmt.Sprintf("%s\n\n%s", goLiveMessage, data.GetServerURL())

	if err := n.discord.Send(message); err != nil {
		log.Errorln("error sending discord message", err)
	}
}

func (n *Notifier) setupTwitter() error {
	if twitterConfig := data.GetTwitterConfiguration(); twitterConfig.Enabled {
		if t, err := twitter.New(twitterConfig.APIKey, twitterConfig.APISecret, twitterConfig.AccessToken, twitterConfig.AccessTokenSecret, twitterConfig.BearerToken); err == nil {
			n.twitter = t
		} else if err != nil {
			return errors.Wrap(err, "error creating twitter notifier")
		}
	}
	return nil
}

func (n *Notifier) notifyTwitter() {
	goLiveMessage := data.GetTwitterConfiguration().GoLiveMessage
	streamTitle := data.GetStreamTitle()
	if streamTitle != "" {
		goLiveMessage += "\n" + streamTitle
	}
	tagString := ""
	for _, tag := range utils.ShuffleStringSlice(data.GetServerMetadataTags()) {
		tagString = fmt.Sprintf("%s #%s", tagString, tag)
	}
	tagString = strings.TrimSpace(tagString)

	message := fmt.Sprintf("%s\n%s\n\n%s", goLiveMessage, data.GetServerURL(), tagString)

	if err := n.twitter.Notify(message); err != nil {
		log.Errorln("error sending twitter message", err)
	}
}

func (n *Notifier) setupEmail() error {
	if emailConfig := data.GetSMTPConfiguration(); emailConfig.Enabled && emailConfig.FromAddress != "" && emailConfig.ListAddress != "" {
		e, err := email.New()
		if err != nil {
			return errors.Wrap(err, "error creating email notifier")
		}
		n.email = e
	}

	return nil
}

func (n *Notifier) notifyEmail() {
	content, err := email.GenerateEmailContent()
	if err != nil {
		log.Errorln("unable to generate email notification content: ", err)
		return
	}

	emailConfig := data.GetSMTPConfiguration()
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
		n.notifyBrowserPush()
	}

	if n.discord != nil {
		n.notifyDiscord()
	}

	if n.email != nil {
		n.notifyEmail()
	}

	if n.twitter != nil {
		n.notifyTwitter()
	}
}
