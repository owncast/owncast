package chat

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/user"
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
		"user_id" INTEGER,
		"body" TEXT,
		"eventType" TEXT,
		"hidden_at" DATE,
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

func addMessage(event events.UserMessageEvent) {
	saveEvent(event.Id, event.User.Id, event.Body, event.Type, event.Hidden, event.Timestamp)
}

func saveEvent(id string, userId string, body string, eventType string, hidden *time.Time, timestamp time.Time) {
	tx, err := _db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO messages(id, user_id, body, eventType, hidden_at, timestamp) values(?, ?, ?, ?, ?, ?)")

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(id, userId, body, eventType, hidden, timestamp); err != nil {
		log.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func getChat(query string) []events.UserMessageEvent {
	fmt.Println(query)

	history := make([]events.UserMessageEvent, 0)
	rows, err := _db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var userId string
		var body string
		var messageType models.EventType
		var hidden *time.Time
		var timestamp time.Time

		err = rows.Scan(&id, &userId, &body, &messageType, &hidden, &timestamp)
		if err != nil {
			log.Debugln(err)
			log.Error("There is a problem with the chat database.  Restore a backup of owncast.db or remove it and start over.")
			break
		}

		user := user.GetUserById(userId)
		message := events.UserMessageEvent{}
		message.Id = id
		message.User = user
		message.Body = body
		message.Type = messageType
		// message.Visible = visible == 1
		message.Timestamp = timestamp

		history = append(history, message)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return history
}

// func getChatModerationHistory() []events.Event {
// 	var query = "SELECT * FROM messages WHERE messageType == 'CHAT' AND datetime(timestamp) >=datetime('now', '-5 Hour')"
// 	return getChat(query)
// }

func GetChatHistory() []events.UserMessageEvent {
	// Get all messages sent within the past 5hrs, max 50
	var query = "SELECT * FROM (SELECT * FROM messages WHERE datetime(timestamp) >=datetime('now', '-5 Hour') AND hidden_at IS NULL ORDER BY timestamp DESC LIMIT 50) ORDER BY timestamp asc"
	return getChat(query)
}

// TODO: Fix this to set timestamp OR null for hidden_at
func saveMessageVisibility(messageIDs []string, visible bool) error {
	tx, err := _db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("UPDATE messages SET hidden_at=? WHERE id IN (?" + strings.Repeat(",?", len(messageIDs)-1) + ")")

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

func getMessageById(messageID string) (*events.UserMessageEvent, error) {
	var query = "SELECT * FROM messages WHERE id = ?"
	row := _db.QueryRow(query, messageID)

	var id string
	var author string
	var body string
	var messageType models.EventType
	var hiddenAt time.Time
	var timestamp time.Time

	err := row.Scan(&id, &author, &body, &messageType, &hiddenAt, &timestamp)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return &events.UserMessageEvent{
		events.Event{
			Id:        id,
			Timestamp: timestamp,
		},
		events.UserEvent{
			Hidden: &hiddenAt,
		},
		events.MessageEvent{
			Body: body,
		},
	}, nil
}
