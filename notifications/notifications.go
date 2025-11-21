package notifications

import (
	"github.com/owncast/owncast/core/data"
	notificationsrepo "github.com/owncast/owncast/persistence/notificationsrepository"
	notificationsservice "github.com/owncast/owncast/services/notifications"
)

// Notifier is an instance of the live stream notifier.
// Deprecated: Use services/notifications.Service interface directly.
type Notifier struct {
	service notificationsservice.Service
}

// Setup will perform any pre-use setup for the notifier.
func Setup(datastore *data.Datastore) {
	notificationsrepo.Setup(datastore)
}

// New creates a new instance of the Notifier.
// Deprecated: Use services/notifications.New directly.
func New(datastore *data.Datastore) (*Notifier, error) {
	service, err := notificationsservice.New(datastore)
	if err != nil {
		return nil, err
	}

	return &Notifier{
		service: service,
	}, nil
}

// Notify will fire the different notification channels.
func (n *Notifier) Notify() {
	n.service.Notify()
}

// RemoveNotificationForChannel removes a notification destination.
func RemoveNotificationForChannel(channel, destination string) error {
	return notificationsrepo.Get().RemoveNotificationForChannel(channel, destination)
}

// GetNotificationDestinationsForChannel will return a collection of
// destinations to notify for a given channel.
func GetNotificationDestinationsForChannel(channel string) ([]string, error) {
	return notificationsrepo.Get().GetNotificationDestinationsForChannel(channel)
}

// AddNotification saves a new user notification destination.
func AddNotification(channel, destination string) error {
	return notificationsrepo.Get().AddNotification(channel, destination)
}
