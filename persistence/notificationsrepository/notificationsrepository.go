package notificationsrepository

import (
	"context"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/tables"

	"github.com/owncast/owncast/notifications/browser"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type NotificationsRepository interface {
	AddNotification(channel, destination string) error
	RemoveNotificationForChannel(channel, destination string) error
	GetNotificationDestinationsForChannel(channel string) ([]string, error)
}

// SqlNotificationsRepository handles database operations for notifications.
type SqlNotificationsRepository struct {
	datastore *data.Datastore
}

// NOTE: This is temporary during the transition period.
var temporaryGlobalInstance NotificationsRepository

// Get will return the notifications repository.
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

// New creates a new instance of the NotificationsRepository.
func New(datastore *data.Datastore) NotificationsRepository {
	return &SqlNotificationsRepository{
		datastore: datastore,
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
