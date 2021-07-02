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
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO messages(id, user_id, body, eventType, hidden_at, timestamp) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	if _, err = stmt.Exec(id, userId, body, eventType, hidden, timestamp); err != nil {
		log.Fatal(err)
	}
	if err = tx.Commit(); err != nil {
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

// SetMessageVisibilityForUserId will bulk change the visibility of messages for a user
// and then send out visibility changed events to chat clients.
func SetMessageVisibilityForUserId(userID string, visible bool) error {
	// Get a list of IDs from this user within the 5hr window to send to the connected clients to hide
	ids := make([]string, 0)
	query := fmt.Sprintf("SELECT * FROM messages WHERE user_id IS '%s'", userID)
	messages := getChat(query)

	if len(messages) == 0 {
		return nil
	}

	for _, message := range messages {
		ids = append(ids, message.Id)
	}

	// Tell the clients to hide/show these messages.
	return SetMessagesVisibility(ids, visible)
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
