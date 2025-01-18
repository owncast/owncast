package tables

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

func CreateWebhooksTable(db *sql.DB) {
	log.Traceln("Creating webhooks table...")

	createTableSQL := `CREATE TABLE IF NOT EXISTS webhooks (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"url" string NOT NULL,
		"events" TEXT NOT NULL,
		"timestamp" DATETIME DEFAULT CURRENT_TIMESTAMP,
		"last_used" DATETIME
	);`

	stmt, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	if _, err = stmt.Exec(); err != nil {
		log.Warnln(err)
	}
}
