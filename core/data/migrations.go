package data

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

func migrateToSchema5(db *sql.DB) {
	// Access tokens have been broken into its own table.

	// Authenticated bool added to the users table.
	stmt, err := db.Prepare("ALTER TABLE users ADD authenticated_at timestamp DEFAULT null ")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Warnln(err)
	}

	// Migrate the access tokens from the users table to the access tokens table.
	query := `SELECT id, access_token, created_at FROM users`
	rows, err := db.Query(query)
	if err != nil || rows.Err() != nil {
		log.Errorln("error migrating access tokens to schema v5", err, rows.Err())
		return
	}
	defer rows.Close()

	valueStrings := []string{}
	valueArgs := []interface{}{}

	var token string
	var userID string
	var timestamp time.Time
	for rows.Next() {
		if err := rows.Scan(&userID, &token, &timestamp); err != nil {
			log.Error("There is a problem reading the database.", err)
			return
		}

		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, userID, token, timestamp)
	}

	smt := `INSERT INTO user_access_tokens(token, user_id, timestamp) VALUES %s ON CONFLICT DO NOTHING`
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))
	// fmt.Println(smt)
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln("Error starting transaction", err)
	}
	_, err = tx.Exec(smt, valueArgs...)
	if err != nil {
		_ = tx.Rollback()
		log.Fatalln("Error inserting access tokens", err)
	}
	if err := tx.Commit(); err != nil {
		log.Fatalln("Error committing transaction", err)
	}

	// Remove old access token column from the users table.
	stmt, err = db.Prepare("ALTER TABLE users DROP COLUMN access_token;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Warnln(err)
	}
}

func migrateToSchema4(db *sql.DB) {
	// Access tokens have been broken into its own table.
	stmt, err := db.Prepare("ALTER TABLE ap_followers ADD COLUMN request_object BLOB")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Warnln(err)
	}
}

func migrateToSchema3(db *sql.DB) {
	// Since it's just a backlog of chat messages let's wipe the old messages
	// and recreate the table.

	// Drop the old messages table
	stmt, err := db.Prepare("DROP TABLE messages")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Warnln(err)
	}

	// Recreate it
	CreateMessagesTable(db)
}

func migrateToSchema2(db *sql.DB) {
	// Since it's just a backlog of chat messages let's wipe the old messages
	// and recreate the table.

	// Drop the old messages table
	stmt, err := db.Prepare("DROP TABLE messages")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Warnln(err)
	}

	// Recreate it
	CreateMessagesTable(db)
}

func migrateToSchema1(db *sql.DB) {
	// Since it's just a backlog of chat messages let's wipe the old messages
	// and recreate the table.

	// Drop the old messages table
	stmt, err := db.Prepare("DROP TABLE messages")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Warnln(err)
	}

	// Recreate it
	CreateMessagesTable(db)

	// Migrate access tokens to become chat users
	type oldAccessToken struct {
		accessToken string
		displayName string
		scopes      string
		createdAt   time.Time
		lastUsedAt  *time.Time
	}

	oldAccessTokens := make([]oldAccessToken, 0)

	query := `SELECT * FROM access_tokens`

	rows, err := db.Query(query)
	if err != nil || rows.Err() != nil {
		log.Errorln("error migrating access tokens to schema v1", err, rows.Err())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var token string
		var name string
		var scopes string
		var timestampString string
		var lastUsedString *string

		if err := rows.Scan(&token, &name, &scopes, &timestampString, &lastUsedString); err != nil {
			log.Error("There is a problem reading the database.", err)
			return
		}

		timestamp, err := time.Parse(time.RFC3339, timestampString)
		if err != nil {
			return
		}

		var lastUsed *time.Time
		if lastUsedString != nil {
			lastUsedTime, _ := time.Parse(time.RFC3339, *lastUsedString)
			lastUsed = &lastUsedTime
		}

		oldToken := oldAccessToken{
			accessToken: token,
			displayName: name,
			scopes:      scopes,
			createdAt:   timestamp,
			lastUsedAt:  lastUsed,
		}

		oldAccessTokens = append(oldAccessTokens, oldToken)
	}

	// Recreate them as users
	for _, token := range oldAccessTokens {
		color := utils.GenerateRandomDisplayColor()
		if err := insertAPIToken(db, token.accessToken, token.displayName, color, token.scopes); err != nil {
			log.Errorln("Error migrating access token", err)
		}
	}
}

func insertAPIToken(db *sql.DB, token string, name string, color int, scopes string) error {
	log.Debugln("Adding new access token:", name)

	id := shortid.MustGenerate()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO users(id, access_token, display_name, display_color, scopes, type) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(id, token, name, color, scopes, "API"); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
