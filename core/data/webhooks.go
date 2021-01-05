package data

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

func createWebhooksTable() {
	log.Traceln("Creating webhooks table...")

	createTableSQL := `CREATE TABLE IF NOT EXISTS webhooks (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"url" string NOT NULL,
		"events" TEXT NOT NULL,
		"timestamp" DATETIME DEFAULT CURRENT_TIMESTAMP,
		"last_used" DATETIME DEFAULT CURRENT_TIMESTAMP
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

func InsertWebhook(url string, events []models.EventType) error {
	log.Println("Adding new webhook:", url)

	eventsString := strings.Join(events, ",")

	tx, err := _db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO webhooks(url, events) values(?, ?)")

	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(url, eventsString); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func DeleteWebhook(id int) error {
	log.Println("Deleting webhook:", id)

	tx, err := _db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("DELETE FROM webhooks WHERE id = ?")

	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	if rowsDeleted, _ := result.RowsAffected(); rowsDeleted == 0 {
		tx.Rollback() //nolint
		return errors.New(fmt.Sprint(id) + " not found")
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func GetWebhooksForEvent(event models.EventType) []models.Webhook {
	webhooks := make([]models.Webhook, 0)

	var query = `SELECT * FROM (
		WITH RECURSIVE split(url, event, rest) AS (
		  SELECT url, '', events || ',' FROM webhooks
		   UNION ALL
		  SELECT url, 
				 substr(rest, 0, instr(rest, ',')),
				 substr(rest, instr(rest, ',')+1)
			FROM split
		   WHERE rest <> '')
		SELECT url, event 
		  FROM split 
		 WHERE event <> ''
	  ) AS webhook WHERE event IS "?"`

	rows, err := _db.Query(query, event)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var url string
		var events string
		var timestamp time.Time
		var lastUsed time.Time

		err = rows.Scan(&id, &url, &events, &timestamp, &lastUsed)
		if err != nil {
			log.Debugln(err)
			log.Error("There is a problem with the database.")
			break
		}

		singleWebhook := models.Webhook{
			ID:        id,
			Url:       url,
			Events:    strings.Split(events, ","),
			Timestamp: timestamp,
			LastUsed:  lastUsed,
		}

		webhooks = append(webhooks, singleWebhook)
	}
	return webhooks
}

func GetWebhooks() ([]models.Webhook, error) {
	webhooks := make([]models.Webhook, 0)

	var query = "SELECT * FROM webhooks"

	rows, err := _db.Query(query)
	if err != nil {
		return webhooks, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var url string
		var events string
		var timestampString string
		var lastUsedString string

		if err := rows.Scan(&id, &url, &events, &timestampString, &lastUsedString); err != nil {
			log.Error("There is a problem reading the database.", err)
			return webhooks, err
		}

		timestamp, err := time.Parse(time.RFC3339, timestampString)
		if err != nil {
			return webhooks, err
		}

		lastUsed, _ := time.Parse(time.RFC3339, lastUsedString)

		singleWebhook := models.Webhook{
			ID:        id,
			Url:       url,
			Events:    strings.Split(events, ","),
			Timestamp: timestamp,
			LastUsed:  lastUsed,
		}

		webhooks = append(webhooks, singleWebhook)
	}

	if err := rows.Err(); err != nil {
		return webhooks, err
	}

	return webhooks, nil
}

func SetWebhookAsUsed(id string) error {
	tx, err := _db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("UPDATE webhooks SET last_used = CURRENT_TIMESTAMP WHERE id = ?")

	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(id); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
