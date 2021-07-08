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
		"last_used" DATETIME 
	);`

	stmt, err := _db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	if _, err = stmt.Exec(); err != nil {
		log.Warnln(err)
	}
}

// InsertWebhook will add a new webhook to the database.
func InsertWebhook(url string, events []models.EventType) (int, error) {
	log.Println("Adding new webhook:", url)

	eventsString := strings.Join(events, ",")

	tx, err := _db.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare("INSERT INTO webhooks(url, events) values(?, ?)")

	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	insertResult, err := stmt.Exec(url, eventsString)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	newID, err := insertResult.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(newID), err
}

// DeleteWebhook will delete a webhook from the database.
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
		_ = tx.Rollback()
		return errors.New(fmt.Sprint(id) + " not found")
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetWebhooksForEvent will return all of the webhooks that want to be notified about an event type.
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
	  ) AS webhook WHERE event IS "` + event + `"`

	rows, err := _db.Query(query)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var url string

		if err := rows.Scan(&url, &event); err != nil {
			log.Debugln(err)
			log.Error("There is a problem with the database.")
			break
		}

		singleWebhook := models.Webhook{
			URL: url,
		}

		webhooks = append(webhooks, singleWebhook)
	}
	return webhooks
}

// GetWebhooks will return all the webhooks.
func GetWebhooks() ([]models.Webhook, error) { //nolint
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
		var lastUsedString *string

		if err := rows.Scan(&id, &url, &events, &timestampString, &lastUsedString); err != nil {
			log.Error("There is a problem reading the database.", err)
			return webhooks, err
		}

		timestamp, err := time.Parse(time.RFC3339, timestampString)
		if err != nil {
			return webhooks, err
		}

		var lastUsed *time.Time = nil
		if lastUsedString != nil {
			lastUsedTime, _ := time.Parse(time.RFC3339, *lastUsedString)
			lastUsed = &lastUsedTime
		}

		singleWebhook := models.Webhook{
			ID:        id,
			URL:       url,
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

// SetWebhookAsUsed will update the last used time for a webhook.
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
