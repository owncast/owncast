package tables

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// CreateFederatedServersTable creates the federated_servers table in the database.
func CreateFederatedServersTable(db *sql.DB) {
	log.Traceln("Creating federated_servers table...")

	createTableSQL := `CREATE TABLE IF NOT EXISTS federated_servers (
		"id" INTEGER NOT NULL PRIMARY KEY,
		"iri" TEXT NOT NULL UNIQUE,
		"name" TEXT,
		"logo_url" TEXT,
		"is_online" BOOLEAN DEFAULT FALSE,
		"stream_title" TEXT,
		"stream_description" TEXT,
		"stream_tags" TEXT,
		"thumbnail_url" TEXT,
		"last_seen_online" TIMESTAMP,
		"last_status_update" TIMESTAMP,
		"added_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"followed_at" TIMESTAMP,
		"pending" BOOLEAN DEFAULT TRUE,
		"username" TEXT,
		"display_name" TEXT,
		"summary" TEXT,
		"accepted_at" TIMESTAMP,
		"rejected_at" TIMESTAMP,
		"follow_status" TEXT DEFAULT 'pending'
	);`

	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatal("error creating federated_servers table", err)
	}

	// Create indexes
	createIndexes := []string{
		`CREATE INDEX IF NOT EXISTS federated_servers_iri ON federated_servers (iri);`,
		`CREATE INDEX IF NOT EXISTS federated_servers_is_online ON federated_servers (is_online);`,
		`CREATE INDEX IF NOT EXISTS federated_servers_last_seen ON federated_servers (last_seen_online);`,
	}

	for _, indexSQL := range createIndexes {
		if _, err := db.Exec(indexSQL); err != nil {
			log.Warnf("error creating index for federated_servers table: %v", err)
		}
	}
}
