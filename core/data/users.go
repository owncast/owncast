package data

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

func createUsersTable(db *sql.DB) {
	log.Traceln("Creating users table...")

	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" TEXT,
		"access_token" string NOT NULL,
		"display_name" TEXT NOT NULL,
		"display_color" NUMBER NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"disabled_at" TIMESTAMP,
		"previous_names" TEXT DEFAULT '',
		"namechanged_at" TIMESTAMP,
		"scopes" TEXT,
		"type" TEXT DEFAULT 'STANDARD',
		"last_used" DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	);CREATE INDEX index ON users (id, access_token, disabled_at);
	CREATE INDEX id ON users (id);
  CREATE INDEX id_disabled ON users (id, disabled_at);
	CREATE INDEX access_token ON users (access_token);
	CREATE INDEX disabled_at ON USERS (disabled_at);`

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
