package user

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type userHistoryEntry struct {
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

func createUserHistoryTable() {
	log.Traceln("Creating user history table...")

	createTableSQL := `CREATE TABLE IF NOT EXISTS user_history (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"display_name" string NOT NULL,
		"changed_at" DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if err := execSQL(createTableSQL); err != nil {
		log.Fatalln(err)
	}
}

func getUserHistory(id string) []*userHistoryEntry {
	entries := []*userHistoryEntry{}

	query := "SELECT display_name, changed_at FROM user_history WHERE id = ? ORDER BY changed_at DESC LIMIT 3"
	rows, err := _db.Query(query, id)

	if err != nil {
		log.Panicln(err)
	}

	for rows.Next() {
		var name string
		var changedAt string

		err = rows.Scan(&name, &changedAt)
		changedAtParsed, err := time.Parse("dunno", changedAt)
		if err != nil {
			log.Fatalln(err)
			continue
		}

		entry := userHistoryEntry{
			name,
			changedAtParsed,
		}
		entries = append(entries, &entry)
	}

	return entries
}
