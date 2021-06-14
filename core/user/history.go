package user

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type usernameHistoryEntry struct {
	DisplayName string    `json:"displayName"`
	ChangedAt   time.Time `json:"changedAt"`
}

func addNameHistory(userId string, name string) error {
	tx, err := _db.Begin()
	stmt, err := tx.Prepare("INSERT INTO user_history(user_id, display_name) values(?, ?)")

	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, name)
	if err != nil {
		panic(err)
	}

	invalidateUserCache(userId)

	return tx.Commit()
}

func createUsernameHistoryTable() {
	log.Traceln("Creating user history table...")

	createTableSQL := `CREATE TABLE IF NOT EXISTS user_history (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"user_id" STRING NOT NULL,
		"display_name" STRING NOT NULL,
		"changed_at" DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if err := execSQL(createTableSQL); err != nil {
		log.Fatalln(err)
	}
}

func getUsernameHistory(id string) []*usernameHistoryEntry {
	entries := []*usernameHistoryEntry{}

	query := "SELECT display_name, changed_at FROM user_history WHERE user_id IS ? ORDER BY changed_at DESC LIMIT 3"
	rows, err := _db.Query(query, id)

	if err != nil {
		log.Panicln(err)
	}

	for rows.Next() {
		var name string
		var changedAt time.Time

		err = rows.Scan(&name, &changedAt)
		// changedAtParsed, err := time.Parse("dunno", changedAt)
		if err != nil {
			log.Fatalln(err)
			continue
		}

		entry := usernameHistoryEntry{
			name,
			changedAt,
		}
		entries = append(entries, &entry)
	}

	return entries
}
