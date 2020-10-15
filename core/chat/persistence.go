package chat

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

var _db *sql.DB

func setupPersistence() {
	file := config.Config.ChatDatabaseFilePath
	// Create empty DB file if it doesn't exist.
	if !utils.DoesFileExists(file) {
		log.Traceln("Creating new chat history database at", file)

		_, err := os.Create(file)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	sqliteDatabase, _ := sql.Open("sqlite3", file)
	_db = sqliteDatabase
	createTable(sqliteDatabase)
}

func createTable(db *sql.DB) {
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
	stmt.Exec()
}

func addMessage(message models.ChatMessage) {
	tx, err := _db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO messages(id, author, body, messageType, visible, timestamp) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(message.ID, message.Author, message.Body, message.MessageType, 1, message.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()

	defer stmt.Close()
}

func getChatHistory() []models.ChatMessage {
	history := make([]models.ChatMessage, 0)

	// Get all messages sent within the past day
	rows, err := _db.Query("SELECT * FROM messages WHERE visible = 1 AND datetime(timestamp) >=datetime('now', '-1 Day')")
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
			log.Fatal(err)
		}

		message := models.ChatMessage{}
		message.ID = id
		message.Author = author
		message.Body = body
		message.MessageType = messageType
		message.Visible = visible == 1
		message.Timestamp = timestamp

		history = append(history, message)
	}

	return history
}
