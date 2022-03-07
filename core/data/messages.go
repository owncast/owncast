package data

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// CreateMessagesTable will create the chat messages table if needed.
func CreateMessagesTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS messages (
		"id" string NOT NULL,
		"user_id" INTEGER,
		"body" TEXT,
		"eventType" TEXT,
		"hidden_at" DATETIME,
		"timestamp" DATETIME,
    "title" TEXT,
    "subtitle" TEXT,
    "image" TEXT,
    "link" TEXT,
		PRIMARY KEY (id)
	);CREATE INDEX index ON messages (id, user_id, hidden_at, timestamp);
	CREATE INDEX id ON messages (id);
	CREATE INDEX user_id ON messages (user_id);
	CREATE INDEX hidden_at ON messages (hidden_at);
	CREATE INDEX timestamp ON messages (timestamp);`

	stmt, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal("error creating chat messages table", err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		log.Fatal("error creating chat messages table", err)
	}
}

// GetMessagesCount will return the number of messages in the database.
func GetMessagesCount() int64 {
	query := `SELECT COUNT(*) FROM messages`
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
