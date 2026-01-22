package persistence

import (
	log "github.com/sirupsen/logrus"
)

func createFederationFollowersTable() {
	log.Traceln("Creating federation followers table...")

	createTableSQL := `CREATE TABLE IF NOT EXISTS ap_followers (
		"iri" TEXT NOT NULL,
		"inbox" TEXT NOT NULL,
		"shared_inbox" TEXT,
		"name" TEXT,
		"username" TEXT NOT NULL,
		"image" TEXT,
    "request" TEXT NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"approved_at" TIMESTAMP,
    "disabled_at" TIMESTAMP,
    "request_object" BLOB,
		PRIMARY KEY (iri));`
	_datastore.MustExec(createTableSQL)
	_datastore.MustExec(`CREATE INDEX IF NOT EXISTS idx_iri ON ap_followers (iri);`)
	_datastore.MustExec(`CREATE INDEX IF NOT EXISTS idx_approved_at ON ap_followers (approved_at);`)
}
