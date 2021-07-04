package data

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

func migrateToSchema1(db *sql.DB) {
	// Since it's just a backlog of chat messages let's wipe the old messages
	// and recreate the table.

	// Drop the old messages table
	stmt, err := db.Prepare("DROP TABLE messages")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Warnln(err)
	}

	// Recreate it
	createTableSQL := `CREATE TABLE IF NOT EXISTS messages (
		"id" string NOT NULL PRIMARY KEY,
		"user_id" INTEGER,
		"body" TEXT,
		"eventType" TEXT,
		"hidden_at" DATE,
		"timestamp" DATE
	);`

	stmt, err = db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Warnln(err)
	}
}
