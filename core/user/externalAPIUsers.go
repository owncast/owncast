package user

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

// ExternalIntegration represents a single 3rd party integration that uses an access token.
type ExternalIntegration struct {
	Id           string     `json:"id"`
	AccessToken  string     `json:"accessToken"`
	DisplayName  string     `json:"displayName"`
	DisplayColor int        `json:"displayColor"`
	CreatedAt    time.Time  `json:"createdAt"`
	Scopes       []string   `json:"scopes"`
	Type         string     `json:"type,omitempty"` // STANDARD or API
	LastUsedAt   *time.Time `json:"lastUsedAt,omitempty"`
}

const (
	// ScopeCanSendUserMessages will allow sending chat messages as users.
	ScopeCanSendUserMessages = "CAN_SEND_MESSAGES"
	// ScopeCanSendSystemMessages will allow sending chat messages as the system.
	ScopeCanSendSystemMessages = "CAN_SEND_SYSTEM_MESSAGES"
	// ScopeHasAdminAccess will allow performing administrative actions on the server.
	ScopeHasAdminAccess = "HAS_ADMIN_ACCESS"
)

// For a scope to be seen as "valid" it must live in this slice.
var validAccessTokenScopes = []string{
	ScopeCanSendUserMessages,
	ScopeCanSendSystemMessages,
	ScopeHasAdminAccess,
}

// InsertToken will add a new token to the database.
func InsertAPIToken(token string, name string, color int, scopes []string) error {
	log.Println("Adding new access token:", name)

	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	scopesString := strings.Join(scopes, ",")
	id := shortid.MustGenerate()

	tx, err := _datastore.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO users(id, access_token, display_name, display_color, scopes, type, previous_names) values(?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(id, token, name, color, scopesString, "API", name); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		panic(err)
		return err
	}

	return nil
}

// DeleteAPIToken will delete a token from the database.
func DeleteAPIToken(token string) error {
	log.Println("Deleting access token:", token)

	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	tx, err := _datastore.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("UPDATE users SET disabled_at = ? WHERE access_token = ?")

	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(time.Now(), token)
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

// GetExternalIntegrationForAccessTokenAndScope will determine if a specific token has access to perform a scoped action.
func GetExternalIntegrationForAccessTokenAndScope(token string, scope string) (*ExternalIntegration, error) {
	// This will split the scopes from comma separated to individual rows
	// so we can efficiently find if a token supports a single scope.
	// This is SQLite specific, so if we ever support other database
	// backends we need to support other methods.
	var query = `SELECT id, access_token, scopes, display_name, display_color, created_at, last_used FROM (
		WITH RECURSIVE split(id, access_token, scopes, display_name, display_color, created_at, last_used, disabled_at, scope, rest) AS (
		  SELECT id, access_token, scopes, display_name, display_color, created_at, last_used, disabled_at, '', scopes || ',' FROM users
		   UNION ALL
		  SELECT id, access_token, scopes, display_name, display_color, created_at, last_used, disabled_at,
				 substr(rest, 0, instr(rest, ',')),
				 substr(rest, instr(rest, ',')+1)
			FROM split
		   WHERE rest <> '')
		SELECT id, access_token, scopes, display_name, display_color, created_at, last_used, disabled_at, scope 
		  FROM split 
		 WHERE scope <> ''
		 ORDER BY access_token, scope
	  ) AS token WHERE token.access_token = ? AND token.scope = ?`

	row := _datastore.DB.QueryRow(query, token, scope)
	integration, err := makeExternalIntegrationFromRow(row)

	return integration, err
}

// GetIntegrationNameForAccessToken will return the integration name associated with a specific access token.
func GetIntegrationNameForAccessToken(token string) *string {
	query := "SELECT display_name FROM users WHERE access_token IS ? AND disabled_at IS NULL"
	row := _datastore.DB.QueryRow(query, token)

	var name string
	err := row.Scan(&name)

	if err != nil {
		log.Warnln(err)
		return nil
	}

	return &name
}

// GetIntegrationAccessTokens will return all access tokens.
func GetIntegrationAccessTokens() ([]ExternalIntegration, error) { //nolint
	// Get all messages sent within the past day
	var query = "SELECT id, access_token, display_name, display_color, scopes, created_at, last_used FROM users WHERE type IS 'API' AND disabled_at IS NULL"

	rows, err := _datastore.DB.Query(query)
	if err != nil {
		return []ExternalIntegration{}, err
	}
	defer rows.Close()

	integrations, err := makeExternalIntegrationsFromRows(rows)

	return integrations, err
}

// SetIntegrationAccessTokenAsUsed will update the last used timestamp for a token.
func SetIntegrationAccessTokenAsUsed(token string) error {
	tx, err := _datastore.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("UPDATE users SET last_used = CURRENT_TIMESTAMP WHERE access_token = ?")

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

func makeExternalIntegrationFromRow(row *sql.Row) (*ExternalIntegration, error) {
	var id string
	var accessToken string
	var displayName string
	var displayColor int
	var scopes string
	var createdAt time.Time
	var lastUsedAt *time.Time

	err := row.Scan(&id, &accessToken, &scopes, &displayName, &displayColor, &createdAt, &lastUsedAt)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	integration := ExternalIntegration{
		Id:           id,
		AccessToken:  accessToken,
		DisplayName:  displayName,
		DisplayColor: displayColor,
		CreatedAt:    createdAt,
		Scopes:       strings.Split(scopes, ","),
		LastUsedAt:   lastUsedAt,
	}

	return &integration, nil
}

func makeExternalIntegrationsFromRows(rows *sql.Rows) ([]ExternalIntegration, error) {
	integrations := make([]ExternalIntegration, 0)

	for rows.Next() {
		var id string
		var accessToken string
		var displayName string
		var displayColor int
		var scopes string
		var createdAt time.Time
		var lastUsedAt *time.Time

		err := rows.Scan(&id, &accessToken, &displayName, &displayColor, &scopes, &createdAt, &lastUsedAt)
		if err != nil {
			log.Errorln(err)
			return nil, err
		}

		integration := ExternalIntegration{
			Id:           id,
			AccessToken:  accessToken,
			DisplayName:  displayName,
			DisplayColor: displayColor,
			CreatedAt:    createdAt,
			Scopes:       strings.Split(scopes, ","),
			LastUsedAt:   lastUsedAt,
		}
		integrations = append(integrations, integration)
	}

	return integrations, nil
}

// HasValidScopes will verify that all the scopes provided are valid.
func HasValidScopes(scopes []string) bool {
	for _, scope := range scopes {
		_, foundInSlice := utils.FindInSlice(validAccessTokenScopes, scope)
		if !foundInSlice {
			log.Errorln("Invalid scope", scope)
			return false
		}
	}
	return true
}
