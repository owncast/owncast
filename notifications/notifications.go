package notifications

import (
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/persistence/configrepository"
	notificationsrepo "github.com/owncast/owncast/persistence/notificationsrepository"
)

// Notifier is an instance of the live stream notifier.
type Notifier struct {
	repository       notificationsrepo.NotificationsRepository
	configRepository configrepository.ConfigRepository
}

// Setup will perform any pre-use setup for the notifier.
func Setup(datastore *data.Datastore) {
	notificationsrepo.Setup(datastore)
}

// New creates a new instance of the Notifier.
func New(datastore *data.Datastore) (*Notifier, error) {
	notifier := Notifier{
		repository:       notificationsrepo.New(datastore),
		configRepository: configrepository.Get(),
	}

	return &notifier, nil
}

// Notify will fire the different notification channels.
func (n *Notifier) Notify() {
	n.repository.Notify()
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
