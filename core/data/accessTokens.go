package data

import (
	"errors"
	"strings"
	"time"

	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

func createAccessTokensTable() {
	log.Traceln("Creating access_tokens table...")

	createTableSQL := `CREATE TABLE IF NOT EXISTS access_tokens (
		"token" string NOT NULL PRIMARY KEY,
		"name" string,
		"scopes" TEXT,
		"timestamp" DATETIME DEFAULT CURRENT_TIMESTAMP,
		"last_used" DATETIME
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

func InsertToken(token string, name string, scopes []string) error {
	log.Println("Adding new access token:", name)

	scopesString := strings.Join(scopes, ",")

	tx, err := _db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO access_tokens(token, name, scopes) values(?, ?, ?)")

	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(token, name, scopesString); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func DeleteToken(token string) error {
	log.Println("Deleting access token:", token)

	tx, err := _db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("DELETE FROM access_tokens WHERE token = ?")

	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(token)
	if err != nil {
		return err
	}

	if rowsDeleted, _ := result.RowsAffected(); rowsDeleted == 0 {
		tx.Rollback() //nolint
		return errors.New(token + " not found")
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func DoesTokenSupportScope(token string, scope string) (bool, error) {
	// This will split the scopes from comma separated to individual rows
	// so we can efficiently find if a token supports a single scope.
	// This is SQLite specific, so if we ever support other database
	// backends we need to support other methods.
	var query = `SELECT count(*) FROM (
		WITH RECURSIVE split(token, scope, rest) AS (
		  SELECT token, '', scopes || ',' FROM access_tokens
		   UNION ALL
		  SELECT token, 
				 substr(rest, 0, instr(rest, ',')),
				 substr(rest, instr(rest, ',')+1)
			FROM split
		   WHERE rest <> '')
		SELECT token, scope 
		  FROM split 
		 WHERE scope <> ''
		 ORDER BY token, scope
	  ) AS token WHERE token.token = ? AND token.scope = ?;`

	row := _db.QueryRow(query, token, scope)

	var count = 0
	err := row.Scan(&count)

	return count > 0, err
}

func GetAccessTokens() ([]models.AccessToken, error) {
	tokens := make([]models.AccessToken, 0)

	// Get all messages sent within the past day
	var query = "SELECT * FROM access_tokens"

	rows, err := _db.Query(query)
	if err != nil {
		return tokens, err
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
			return tokens, err
		}

		timestamp, err := time.Parse(time.RFC3339, timestampString)
		if err != nil {
			return tokens, err
		}

		var lastUsed *time.Time = nil
		if lastUsedString != nil {
			lastUsedTime, _ := time.Parse(time.RFC3339, *lastUsedString)
			lastUsed = &lastUsedTime
		}

		singleToken := models.AccessToken{
			Name:      name,
			Token:     token,
			Scopes:    strings.Split(scopes, ","),
			Timestamp: timestamp,
			LastUsed:  lastUsed,
		}

		tokens = append(tokens, singleToken)
	}

	if err := rows.Err(); err != nil {
		return tokens, err
	}

	return tokens, nil
}

func SetAccessTokenAsUsed(token string) error {
	tx, err := _db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("UPDATE access_tokens SET last_used = CURRENT_TIMESTAMP WHERE token = ?")

	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(token); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
