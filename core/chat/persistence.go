package chat

import (
	"database/sql"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

var _db *sql.DB

func setupPersistence() {
	_db = data.GetDatabase()
	createTable()
}

func createTable() {
	createTableSQL := `CREATE TABLE IF NOT EXISTS messages (
		"id" string NOT NULL PRIMARY KEY,
		"author" TEXT,
		"body" TEXT,
		"messageType" TEXT,
		"visible" INTEGER,
		"timestamp" DATE
	);`

	stmt, err := _db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Warnln(err)
	}
}

func addMessage(message models.ChatEvent) {
	tx, err := _db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO messages(id, author, body, messageType, visible, timestamp) values(?, ?, ?, ?, ?, ?)")

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(message.ID, message.Author, message.Body, message.MessageType, 1, message.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func getChatHistory(filtered bool) []models.ChatEvent {
	history := make([]models.ChatEvent, 0)

	// Get all messages sent within the past day
	var query = "SELECT * FROM messages WHERE messageType != 'SYSTEM' AND datetime(timestamp) >=datetime('now', '-1 Day')"
	if filtered {
		query = query + " AND visible = 1"
	}

	rows, err := _db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var author string
		var body string
		var messageType string
		var visible int
		var timestamp time.Time

		err = rows.Scan(&id, &author, &body, &messageType, &visible, &timestamp)
		if err != nil {
			log.Debugln(err)
			log.Error("There is a problem with the chat database.  Please delete chat.db and restart Owncast.")
			break
		}

		message := models.ChatEvent{}
		message.ID = id
		message.Author = author
		message.Body = body
		message.MessageType = messageType
		message.Visible = visible == 1
		message.Timestamp = timestamp

		history = append(history, message)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return history
}

func saveMessageVisibility(messageIDs []string, visible bool) error {
	tx, err := _db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("UPDATE messages SET visible=? WHERE id IN (?)")

	strIDList := strings.Join(messageIDs, ",")

	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(visible, strIDList); err != nil {
		log.Fatal(err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
