package notifications

import (
	"context"
	"database/sql"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/db"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func createNotificationsTable(db *sql.DB) {
	log.Traceln("Creating federation followers table...")

	createTableSQL := `CREATE TABLE IF NOT EXISTS notifications (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"channel" TEXT NOT NULL,
		"destination" TEXT NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
		CREATE INDEX channel_index ON notifications (channel);`

	stmt, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Warnln("error executing sql creating followers table", createTableSQL, err)
	}
}

// AddNotification saves a new user notification destination.
func AddNotification(channel, destination string) error {
	return data.GetDatastore().GetQueries().AddNotification(context.Background(), db.AddNotificationParams{
		Channel:     channel,
		Destination: destination,
	})
}

// RemoveNotificationForChannel removes a notification destination..
func RemoveNotificationForChannel(channel, destination string) error {
	log.Println("Removing notification for channel", channel)
	return data.GetDatastore().GetQueries().RemoveNotificationDestinationForChannel(context.Background(), db.RemoveNotificationDestinationForChannelParams{
		Channel:     channel,
		Destination: destination,
	})
}

// GetNotificationDestinationsForChannel will return a collection of
// destinations to notify for a given channel.
func GetNotificationDestinationsForChannel(channel string) ([]string, error) {
	result, err := data.GetDatastore().GetQueries().GetNotificationDestinationsForChannel(context.Background(), channel)
	if err != nil {
		return nil, errors.Wrap(err, "unable to query notification destinations for channel "+channel)
	}

	return result, nil
}
