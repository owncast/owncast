package data

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

func createUsersTable(db *sql.DB) {
	log.Traceln("Creating users table...")

	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" TEXT PRIMARY KEY,
		"access_token" string NOT NULL,
		"display_name" TEXT NOT NULL,
		"display_color" NUMBER NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"disabled_at" TIMESTAMP,
		"previous_names" TEXT DEFAULT '',
		"namechanged_at" TIMESTAMP,
		"scopes" TEXT,
		"type" TEXT DEFAULT 'STANDARD',
		"last_used" DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	stmt, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Warnln(err)
	}
}
