package notificationsrepository

import (
	"context"

	"github.com/owncast/owncast/storage/data"
	"github.com/owncast/owncast/storage/sqlstorage"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type NotificationsRepository interface{}

type SqlNotificationsRepository struct {
	datastore *data.Store
}

func New(datastore *data.Store) *SqlNotificationsRepository {
	return &SqlNotificationsRepository{datastore}
}

var temporaryGlobalInstance *SqlNotificationsRepository

func Get() *SqlNotificationsRepository {
	if temporaryGlobalInstance == nil {
		temporaryGlobalInstance = &SqlNotificationsRepository{}
	}
	return temporaryGlobalInstance
}

// AddNotification saves a new user notification destination.
func (r *SqlNotificationsRepository) AddNotification(channel, destination string) error {
	return data.GetDatastore().GetQueries().AddNotification(context.Background(), sqlstorage.AddNotificationParams{
		Channel:     channel,
		Destination: destination,
	})
}

// RemoveNotificationForChannel removes a notification destination.
func (r *SqlNotificationsRepository) RemoveNotificationForChannel(channel, destination string) error {
	log.Debugln("Removing notification for channel", channel)
	return data.GetDatastore().GetQueries().RemoveNotificationDestinationForChannel(context.Background(), sqlstorage.RemoveNotificationDestinationForChannelParams{
		Channel:     channel,
		Destination: destination,
	})
}

// GetNotificationDestinationsForChannel will return a collection of
// destinations to notify for a given channel.
func (r *SqlNotificationsRepository) GetNotificationDestinationsForChannel(channel string) ([]string, error) {
	result, err := data.GetDatastore().GetQueries().GetNotificationDestinationsForChannel(context.Background(), channel)
	if err != nil {
		return nil, errors.Wrap(err, "unable to query notification destinations for channel "+channel)
	}

	return result, nil
}
