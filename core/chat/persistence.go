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
	if _, err := stmt.Exec(); err != nil {
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

	if _, err := stmt.Exec(message.ID, message.Author, message.Body, message.MessageType, 1, message.Timestamp); err != nil {
		log.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func getChat(query string) []models.ChatEvent {
	history := make([]models.ChatEvent, 0)
	rows, err := _db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var author string
		var body string
		var messageType models.EventType
		var visible int
		var timestamp time.Time

		err = rows.Scan(&id, &author, &body, &messageType, &visible, &timestamp)
		if err != nil {
			log.Debugln(err)
			log.Error("There is a problem with the chat database.  Restore a backup of owncast.db or remove it and start over.")
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

func getChatModerationHistory() []models.ChatEvent {
	var query = "SELECT * FROM messages WHERE messageType == 'CHAT' AND datetime(timestamp) >=datetime('now', '-5 Hour')"
	return getChat(query)
}

func getChatHistory() []models.ChatEvent {
	// Get all messages sent within the past 5hrs, max 50
	var query = "SELECT * FROM (SELECT * FROM messages WHERE datetime(timestamp) >=datetime('now', '-5 Hour') AND visible = 1 ORDER BY timestamp DESC LIMIT 50) ORDER BY timestamp asc"
	return getChat(query)
}

func saveMessageVisibility(messageIDs []string, visible bool) error {
	tx, err := _db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("UPDATE messages SET visible=? WHERE id IN (?" + strings.Repeat(",?", len(messageIDs)-1) + ")")

	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()

	args := make([]interface{}, len(messageIDs)+1)
	args[0] = visible
	for i, id := range messageIDs {
		args[i+1] = id
	}

	if _, err := stmt.Exec(args...); err != nil {
		log.Fatal(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func getMessageById(messageID string) (models.ChatEvent, error) {
	var query = "SELECT * FROM messages WHERE id = ?"
	row := _db.QueryRow(query, messageID)

	var id string
	var author string
	var body string
	var messageType models.EventType
	var visible int
	var timestamp time.Time

	err := row.Scan(&id, &author, &body, &messageType, &visible, &timestamp)
	if err != nil {
		log.Errorln(err)
		return models.ChatEvent{}, err
	}

	return models.ChatEvent{
		ID:          id,
		Author:      author,
		Body:        body,
		MessageType: messageType,
		Visible:     visible == 1,
		Timestamp:   timestamp,
	}, nil
}
