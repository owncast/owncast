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
	saveEvent(event.Id, event.User.Id, event.Body, event.Type, event.HiddenAt, event.Timestamp)
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
		var hiddenAt *time.Time
		var timestamp time.Time

		err = rows.Scan(&id, &userId, &body, &messageType, &hiddenAt, &timestamp)
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
		message.HiddenAt = hiddenAt
		message.Timestamp = timestamp

		history = append(history, message)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return history
}

func GetChatModerationHistory() []events.UserMessageEvent {
	// Get all messages sent within the past 5hrs
	var query = "SELECT * FROM messages WHERE datetime(timestamp) >=datetime('now', '-5 Hour') ORDER BY timestamp desc"
	return getChat(query)
}

func GetChatHistory() []events.UserMessageEvent {
	// Get all visible messages sent within the past 5hrs, max 50
	var query = "SELECT * FROM (SELECT * FROM messages WHERE datetime(timestamp) >=datetime('now', '-5 Hour') AND hidden_at IS NULL ORDER BY timestamp DESC LIMIT 50) ORDER BY timestamp asc"
	return getChat(query)
}

// SetMessageVisibilityForUserId will bulk change the visibility of messages for a user.
func SetMessageVisibilityForUserId(userID string, visible bool) error {
	tx, err := _db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	var disabledAt *time.Time
	if !visible {
		now := time.Now()
		disabledAt = &now
	} else {
		disabledAt = nil
	}

	// Set all messages by this user as hidden in the database
	stmt, err := tx.Prepare("UPDATE messages SET hidden_at=? WHERE user_id IS ?")
	defer stmt.Close()
	stmt.Exec(disabledAt, userID)

	if err != nil {
		log.Fatal(err)
		return err
	}

	// Get a list of IDs to send to the connected clients to hide
	// stmt, err = tx.Prepare("SELECT id, hidden_at FROM (SELECT * FROM messages WHERE datetime(timestamp) >=datetime('now', '-5 Hour') ORDER BY timestamp DESC LIMIT 50) ORDER BY timestamp asc WHERE user_id IS ?")
	// defer stmt.Close()
	// result, err := stmt.Exec(userID)

	// history := make([]events.UserMessageEvent, 0)
	// rows, err := _db.Query(query)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var id string
	// 	var hidden *time.Time
	// 	err = rows.Scan(&id, &hidden)
	// 	if err != nil {
	// 		log.Debugln(err)
	// 		log.Error("There is a problem with the chat database.  Restore a backup of owncast.db or remove it and start over.")
	// 		break
	// 	}

	// }

	return nil
}

// DisconnectUser will forcefully disconnect all clients belonging to a user by ID.
func DisconnectUser(userID string) {
	for _, client := range _server.clients {
		if client.User.Id == userID {
			client.close()
		}
	}
}

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

	var hiddenAt *time.Time
	if !visible {
		now := time.Now()
		hiddenAt = &now
	} else {
		hiddenAt = nil
	}

	args := make([]interface{}, len(messageIDs)+1)
	args[0] = hiddenAt
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
	var userId string
	var body string
	var eventType models.EventType
	var hiddenAt *time.Time
	var timestamp time.Time

	err := row.Scan(&id, &userId, &body, &eventType, &hiddenAt, &timestamp)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	user := user.GetUserById(userId)

	return &events.UserMessageEvent{
		events.Event{
			Type:      eventType,
			Id:        id,
			Timestamp: timestamp,
		},
		events.UserEvent{
			User:     user,
			HiddenAt: hiddenAt,
		},
		events.MessageEvent{
			OutboundEvent: nil,
			Body:          body,
		},
	}, nil
}
