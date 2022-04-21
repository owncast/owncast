package data

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

func createAccessTokenTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS user_access_tokens (
    "token" TEXT NOT NULL PRIMARY KEY,
    "user_id" TEXT NOT NULL,
    "timestamp" DATE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id)
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

func createUsersTable(db *sql.DB) {
	log.Traceln("Creating users table...")

	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" TEXT,
		"display_name" TEXT NOT NULL,
		"display_color" NUMBER NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"disabled_at" TIMESTAMP,
		"previous_names" TEXT DEFAULT '',
		"namechanged_at" TIMESTAMP,
    "authenticated_at" TIMESTAMP,
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

// GetUsersCount will return the number of users in the database.
func GetUsersCount() int64 {
	query := `SELECT COUNT(*) FROM users`
	rows, err := _db.Query(query)
	if err != nil || rows.Err() != nil {
		return 0
	}
	defer rows.Close()
	var count int64
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0
		}
	}
	return count
}
