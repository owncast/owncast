package notificationsrepository

import (
	"context"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/datastore"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/notifications/browser"
)

type NotificationsRepository interface {
	AddNotification(channel, destination string) error
	RemoveNotificationForChannel(channel, destination string) error
	GetNotificationDestinationsForChannel(channel string) ([]string, error)
	Setup()
}

// SqlNotificationsRepository handles database operations for notifications.
type SqlNotificationsRepository struct {
	datastore        *datastore.Datastore
	configRepository configrepository.ConfigRepository
}

// Setup will perform any pre-use setup for the notifier.
// The notifications table itself is created by the goose migrations package.
func (n *SqlNotificationsRepository) Setup() {
	n.initializeBrowserPushIfNeeded()
}

func (n *SqlNotificationsRepository) initializeBrowserPushIfNeeded() {
	pubKey, _ := n.configRepository.GetBrowserPushPublicKey()
	privKey, _ := n.configRepository.GetBrowserPushPrivateKey()

	// We need browser push keys so people can register for pushes.
	if pubKey == "" || privKey == "" {
		browserPrivateKey, browserPublicKey, err := browser.GenerateBrowserPushKeys()
		if err != nil {
			log.Errorln("unable to initialize browser push notification keys", err)
		}

		if err := n.configRepository.SetBrowserPushPrivateKey(browserPrivateKey); err != nil {
			log.Errorln("unable to set browser push private key", err)
		}

		if err := n.configRepository.SetBrowserPushPublicKey(browserPublicKey); err != nil {
			log.Errorln("unable to set browser push public key", err)
		}
	}

	// Enable browser push notifications by default.
	if !n.configRepository.GetHasPerformedInitialNotificationsConfig() {
		_ = n.configRepository.SetBrowserPushConfig(models.BrowserNotificationConfiguration{Enabled: true, GoLiveMessage: config.GetDefaults().FederationGoLiveMessage})
		_ = n.configRepository.SetHasPerformedInitialNotificationsConfig(true)
	}
}

// New creates a new instance of the NotificationsRepository.
func New(datastore *datastore.Datastore, configRepository configrepository.ConfigRepository) NotificationsRepository {
	return &SqlNotificationsRepository{
		datastore:        datastore,
		configRepository: configRepository,
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
