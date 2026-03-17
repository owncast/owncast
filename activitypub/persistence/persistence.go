package persistence

import (
	"github.com/owncast/owncast/core/data"

	log "github.com/sirupsen/logrus"
)

var _datastore *data.Datastore

// Setup will initialize the ActivityPub persistence layer with the provided datastore.
func Setup(datastore *data.Datastore) {
	_datastore = datastore
	createFederationFollowersTable()
	createFederationOutboxTable()
	createFederatedActivitiesTable()
	addFollowersFixtureData()
}

// GetDatastore returns the datastore instance for use by sub-repositories.
func GetDatastore() *data.Datastore {
	return _datastore
}

// createFederatedActivitiesTable will create the accepted
// activities table if needed.
func createFederatedActivitiesTable() {
	createTableSQL := `CREATE TABLE IF NOT EXISTS ap_accepted_activities (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"iri" TEXT NOT NULL,
    "actor" TEXT NOT NULL,
    "type" TEXT NOT NULL,
		"timestamp" TIMESTAMP NOT NULL
	);`

	_datastore.MustExec(createTableSQL)
	_datastore.MustExec(`CREATE INDEX IF NOT EXISTS idx_iri_actor_index ON ap_accepted_activities (iri,actor);`)
}

func createFederationOutboxTable() {
	log.Traceln("Creating federation outbox table...")
	createTableSQL := `CREATE TABLE IF NOT EXISTS ap_outbox (
		"iri" TEXT NOT NULL,
		"value" BLOB,
		"type" TEXT NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "live_notification" BOOLEAN DEFAULT FALSE,
		PRIMARY KEY (iri));`

	_datastore.MustExec(createTableSQL)
	_datastore.MustExec(`CREATE INDEX IF NOT EXISTS idx_iri ON ap_outbox (iri);`)
	_datastore.MustExec(`CREATE INDEX IF NOT EXISTS idx_type ON ap_outbox (type);`)
	_datastore.MustExec(`CREATE INDEX IF NOT EXISTS idx_live_notification ON ap_outbox (live_notification);`)
}
