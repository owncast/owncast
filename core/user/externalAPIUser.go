package user

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/owncast/owncast/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

// ExternalAPIUser represents a single 3rd party integration that uses an access token.
// This struct mostly matches the User struct so they can be used interchangeably.
type ExternalAPIUser struct {
	ID           string     `json:"id"`
	AccessToken  string     `json:"accessToken"`
	DisplayName  string     `json:"displayName"`
	DisplayColor int        `json:"displayColor"`
	CreatedAt    time.Time  `json:"createdAt"`
	Scopes       []string   `json:"scopes"`
	Type         string     `json:"type,omitempty"` // Should be API
	LastUsedAt   *time.Time `json:"lastUsedAt,omitempty"`
	IsBot        bool       `json:"isBot"`
}

const (
	// ScopeCanSendChatMessages will allow sending chat messages as itself.
	ScopeCanSendChatMessages = "CAN_SEND_MESSAGES"
	// ScopeCanSendSystemMessages will allow sending chat messages as the system.
	ScopeCanSendSystemMessages = "CAN_SEND_SYSTEM_MESSAGES"
	// ScopeHasAdminAccess will allow performing administrative actions on the server.
	ScopeHasAdminAccess = "HAS_ADMIN_ACCESS"
)

// For a scope to be seen as "valid" it must live in this slice.
var validAccessTokenScopes = []string{
	ScopeCanSendChatMessages,
	ScopeCanSendSystemMessages,
	ScopeHasAdminAccess,
}

// InsertExternalAPIUser will add a new API user to the database.
func InsertExternalAPIUser(token string, name string, color int, scopes []string) error {
	log.Traceln("Adding new API user")

	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	scopesString := strings.Join(scopes, ",")
	id := shortid.MustGenerate()

	tx, err := _datastore.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO users(id, display_name, display_color, scopes, type, previous_names) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(id, name, color, scopesString, "API", name); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	if err := addAccessTokenForUser(token, id); err != nil {
		return errors.Wrap(err, "unable to save access token for new external api user")
	}

	return nil
}

// DeleteExternalAPIUser will delete a token from the database.
func DeleteExternalAPIUser(token string) error {
	log.Traceln("Deleting access token")

	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	tx, err := _datastore.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("UPDATE users SET disabled_at = CURRENT_TIMESTAMP WHERE id = (SELECT user_id FROM user_access_tokens WHERE token = ?)")
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

// GetExternalAPIUserForAccessTokenAndScope will determine if a specific token has access to perform a scoped action.
func GetExternalAPIUserForAccessTokenAndScope(token string, scope string) (*ExternalAPIUser, error) {
	// This will split the scopes from comma separated to individual rows
	// so we can efficiently find if a token supports a single scope.
	// This is SQLite specific, so if we ever support other database
	// backends we need to support other methods.
	query := `SELECT id, scopes, display_name, display_color, created_at, last_used FROM user_access_tokens, (
		WITH RECURSIVE split(id, scopes, display_name, display_color, created_at, last_used, disabled_at, scope, rest) AS (
		  SELECT id, scopes, display_name, display_color, created_at, last_used, disabled_at, '', scopes || ',' FROM users
		   UNION ALL
		  SELECT id, scopes, display_name, display_color, created_at, last_used, disabled_at,
				 substr(rest, 0, instr(rest, ',')),
				 substr(rest, instr(rest, ',')+1)
			FROM split
		   WHERE rest <> '')
		SELECT id, scopes, display_name, display_color, created_at, last_used, disabled_at, scope
		  FROM split
		 WHERE scope <> ''
		 ORDER BY scope
	  ) AS token WHERE user_access_tokens.token = ? AND token.scope = ?`

	row := _datastore.DB.QueryRow(query, token, scope)
	integration, err := makeExternalAPIUserFromRow(row)

	return integration, err
}

// GetIntegrationNameForAccessToken will return the integration name associated with a specific access token.
func GetIntegrationNameForAccessToken(token string) *string {
	name, err := _datastore.GetQueries().GetUserDisplayNameByToken(context.Background(), token)
	if err != nil {
		return nil
	}

	return &name
}

// GetExternalAPIUser will return all API users with access tokens.
func GetExternalAPIUser() ([]ExternalAPIUser, error) { //nolint
	// Get all messages sent within the past day
	query := "SELECT id, token, display_name, display_color, scopes, created_at, last_used FROM users, user_access_tokens WHERE user_access_tokens.user_id = id  AND type IS 'API' AND disabled_at IS NULL"

	rows, err := _datastore.DB.Query(query)
	if err != nil {
		return []ExternalAPIUser{}, err
	}
	defer rows.Close()

	integrations, err := makeExternalAPIUsersFromRows(rows)

	return integrations, err
}

// SetExternalAPIUserAccessTokenAsUsed will update the last used timestamp for a token.
func SetExternalAPIUserAccessTokenAsUsed(token string) error {
	tx, err := _datastore.DB.Begin()
	if err != nil {
		return err
	}
	// stmt, err := tx.Prepare("UPDATE users SET last_used = CURRENT_TIMESTAMP WHERE access_token = ?")
	stmt, err := tx.Prepare("UPDATE users SET last_used = CURRENT_TIMESTAMP WHERE id = (SELECT user_id FROM user_access_tokens WHERE token = ?)")
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

func makeExternalAPIUserFromRow(row *sql.Row) (*ExternalAPIUser, error) {
	var id string
	var displayName string
	var displayColor int
	var scopes string
	var createdAt time.Time
	var lastUsedAt *time.Time

	err := row.Scan(&id, &scopes, &displayName, &displayColor, &createdAt, &lastUsedAt)
	if err != nil {
		log.Debugln("unable to convert row to api user", err)
		return nil, err
	}

	integration := ExternalAPIUser{
		ID:           id,
		DisplayName:  displayName,
		DisplayColor: displayColor,
		CreatedAt:    createdAt,
		Scopes:       strings.Split(scopes, ","),
		LastUsedAt:   lastUsedAt,
	}

	return &integration, nil
}

func makeExternalAPIUsersFromRows(rows *sql.Rows) ([]ExternalAPIUser, error) {
	integrations := make([]ExternalAPIUser, 0)

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

		integration := ExternalAPIUser{
			ID:           id,
			AccessToken:  accessToken,
			DisplayName:  displayName,
			DisplayColor: displayColor,
			CreatedAt:    createdAt,
			Scopes:       strings.Split(scopes, ","),
			LastUsedAt:   lastUsedAt,
			IsBot:        true,
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
			return false
		}
	}
	return true
}
