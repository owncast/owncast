package userrepository

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/db"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
	"github.com/pkg/errors"
	"github.com/teris-io/shortid"

	log "github.com/sirupsen/logrus"
)

type UserRepository interface {
	ChangeUserColor(userID string, color int) error
	ChangeUsername(userID string, username string) error
	CreateAnonymousUser(displayName string) (*models.User, string, error)
	DeleteExternalAPIUser(token string) error
	GetDisabledUsers() []*models.User
	GetExternalAPIUser() ([]models.ExternalAPIUser, error)
	GetExternalAPIUserForAccessTokenAndScope(token string, scope string) (*models.ExternalAPIUser, error)
	GetModeratorUsers() []*models.User
	GetUserByID(id string) *models.User
	GetUserByToken(token string) *models.User
	InsertExternalAPIUser(token string, name string, color int, scopes []string) error
	IsDisplayNameAvailable(displayName string) (bool, error)
	SetAccessTokenToOwner(token, userID string) error
	SetEnabled(userID string, enabled bool) error
	SetModerator(userID string, isModerator bool) error
	SetUserAsAuthenticated(userID string) error
	HasValidScopes(scopes []string) bool
	GetUserByAuth(authToken string, authType models.AuthType) *models.User
	AddAuth(userID, authToken string, authType models.AuthType) error
	SetExternalAPIUserAccessTokenAsUsed(token string) error
	GetUsersCount() int
}

type SqlUserRepository struct {
	datastore *data.Datastore
}

// NOTE: This is temporary during the transition period.
var temporaryGlobalInstance UserRepository

// Get will return the user repository.
func Get() UserRepository {
	if temporaryGlobalInstance == nil {
		i := New(data.GetDatastore())
		temporaryGlobalInstance = i
	}
	return temporaryGlobalInstance
}

// New will create a new instance of the UserRepository.
func New(datastore *data.Datastore) UserRepository {
	r := SqlUserRepository{
		datastore: datastore,
	}

	return &r
}

// CreateAnonymousUser will create a new anonymous user with the provided display name.
func (r *SqlUserRepository) CreateAnonymousUser(displayName string) (*models.User, string, error) {
	if displayName == "" {
		return nil, "", errors.New("display name cannot be empty")
	}

	// Try to assign a name that was requested.
	// If name isn't available then generate a random one.
	if available, _ := r.IsDisplayNameAvailable(displayName); !available {
		rand, _ := utils.GenerateRandomString(3)
		displayName += rand
	}

	displayColor := utils.GenerateRandomDisplayColor(config.MaxUserColor)

	id := shortid.MustGenerate()
	user := &models.User{
		ID:           id,
		DisplayName:  displayName,
		DisplayColor: displayColor,
		CreatedAt:    time.Now(),
	}

	// Create new user.
	if err := r.create(user); err != nil {
		return nil, "", err
	}

	// Assign it an access token.
	accessToken, err := utils.GenerateAccessToken()
	if err != nil {
		log.Errorln("Unable to create access token for new user")
		return nil, "", err
	}
	if err := r.addAccessTokenForUser(accessToken, id); err != nil {
		return nil, "", errors.Wrap(err, "unable to save access token for new user")
	}

	return user, accessToken, nil
}

// IsDisplayNameAvailable will check if the proposed name is available for use.
func (r *SqlUserRepository) IsDisplayNameAvailable(displayName string) (bool, error) {
	if available, err := r.datastore.GetQueries().IsDisplayNameAvailable(context.Background(), displayName); err != nil {
		return false, errors.Wrap(err, "unable to check if display name is available")
	} else if available != 0 {
		return false, nil
	}

	return true, nil
}

// ChangeUsername will change the user associated to userID from one display name to another.
func (r *SqlUserRepository) ChangeUsername(userID string, username string) error {
	r.datastore.DbLock.Lock()
	defer r.datastore.DbLock.Unlock()

	if err := r.datastore.GetQueries().ChangeDisplayName(context.Background(), db.ChangeDisplayNameParams{
		DisplayName:   username,
		ID:            userID,
		PreviousNames: sql.NullString{String: fmt.Sprintf(",%s", username), Valid: true},
		NamechangedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}); err != nil {
		return errors.Wrap(err, "unable to change display name")
	}

	return nil
}

// ChangeUserColor will change the user associated to userID from one display name to another.
func (r *SqlUserRepository) ChangeUserColor(userID string, color int) error {
	r.datastore.DbLock.Lock()
	defer r.datastore.DbLock.Unlock()

	if err := r.datastore.GetQueries().ChangeDisplayColor(context.Background(), db.ChangeDisplayColorParams{
		DisplayColor: color,
		ID:           userID,
	}); err != nil {
		return errors.Wrap(err, "unable to change display color")
	}

	return nil
}

func (r *SqlUserRepository) addAccessTokenForUser(accessToken, userID string) error {
	return r.datastore.GetQueries().AddAccessTokenForUser(context.Background(), db.AddAccessTokenForUserParams{
		Token:  accessToken,
		UserID: userID,
	})
}

func (r *SqlUserRepository) create(user *models.User) error {
	r.datastore.DbLock.Lock()
	defer r.datastore.DbLock.Unlock()

	tx, err := r.datastore.DB.Begin()
	if err != nil {
		log.Debugln(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare("INSERT INTO users(id, display_name, display_color, previous_names, created_at) values(?, ?, ?, ?, ?)")
	if err != nil {
		log.Debugln(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID, user.DisplayName, user.DisplayColor, user.DisplayName, user.CreatedAt)
	if err != nil {
		log.Errorln("error creating new user", err)
		return err
	}

	return tx.Commit()
}

// SetEnabled will set the enabled status of a single user by ID.
func (r *SqlUserRepository) SetEnabled(userID string, enabled bool) error {
	r.datastore.DbLock.Lock()
	defer r.datastore.DbLock.Unlock()

	tx, err := r.datastore.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback() //nolint

	var stmt *sql.Stmt
	if !enabled {
		stmt, err = tx.Prepare("UPDATE users SET disabled_at=DATETIME('now', 'localtime') WHERE id IS ?")
	} else {
		stmt, err = tx.Prepare("UPDATE users SET disabled_at=null WHERE id IS ?")
	}

	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err := stmt.Exec(userID); err != nil {
		return err
	}

	return tx.Commit()
}

// GetUserByToken will return a user by an access token.
func (r *SqlUserRepository) GetUserByToken(token string) *models.User {
	u, err := r.datastore.GetQueries().GetUserByAccessToken(context.Background(), token)
	if err != nil {
		return nil
	}

	var scopes []string
	if u.Scopes.Valid {
		scopes = strings.Split(u.Scopes.String, ",")
	}

	var disabledAt *time.Time
	if u.DisabledAt.Valid {
		disabledAt = &u.DisabledAt.Time
	}

	var authenticatedAt *time.Time
	if u.AuthenticatedAt.Valid {
		authenticatedAt = &u.AuthenticatedAt.Time
	}

	return &models.User{
		ID:              u.ID,
		DisplayName:     u.DisplayName,
		DisplayColor:    int(u.DisplayColor),
		CreatedAt:       u.CreatedAt.Time,
		DisabledAt:      disabledAt,
		PreviousNames:   strings.Split(u.PreviousNames.String, ","),
		NameChangedAt:   &u.NamechangedAt.Time,
		AuthenticatedAt: authenticatedAt,
		Authenticated:   authenticatedAt != nil,
		Scopes:          scopes,
	}
}

// SetAccessTokenToOwner will reassign an access token to be owned by a
// different user. Used for logging in with external auth.
func (r *SqlUserRepository) SetAccessTokenToOwner(token, userID string) error {
	return r.datastore.GetQueries().SetAccessTokenToOwner(context.Background(), db.SetAccessTokenToOwnerParams{
		UserID: userID,
		Token:  token,
	})
}

// SetUserAsAuthenticated will mark that a user has been authenticated
// in some way.
func (r *SqlUserRepository) SetUserAsAuthenticated(userID string) error {
	return errors.Wrap(r.datastore.GetQueries().SetUserAsAuthenticated(context.Background(), userID), "unable to set user as authenticated")
}

// AddAuth will add an external authentication token and type for a user.
func (r *SqlUserRepository) AddAuth(userID, authToken string, authType models.AuthType) error {
	return r.datastore.GetQueries().AddAuthForUser(context.Background(), db.AddAuthForUserParams{
		UserID: userID,
		Token:  authToken,
		Type:   string(authType),
	})
}

// GetUserByAuth will return an existing user given auth details if a user
// has previously authenticated with that method.
func (r *SqlUserRepository) GetUserByAuth(authToken string, authType models.AuthType) *models.User {
	u, err := r.datastore.GetQueries().GetUserByAuth(context.Background(), db.GetUserByAuthParams{
		Token: authToken,
		Type:  string(authType),
	})
	if err != nil {
		return nil
	}

	var scopes []string
	if u.Scopes.Valid {
		scopes = strings.Split(u.Scopes.String, ",")
	}

	return &models.User{
		ID:              u.ID,
		DisplayName:     u.DisplayName,
		DisplayColor:    int(u.DisplayColor),
		CreatedAt:       u.CreatedAt.Time,
		DisabledAt:      &u.DisabledAt.Time,
		PreviousNames:   strings.Split(u.PreviousNames.String, ","),
		NameChangedAt:   &u.NamechangedAt.Time,
		AuthenticatedAt: &u.AuthenticatedAt.Time,
		Scopes:          scopes,
	}
}

// SetModerator will add or remove moderator status for a single user by ID.
func (r *SqlUserRepository) SetModerator(userID string, isModerator bool) error {
	if isModerator {
		return r.addScopeToUser(userID, models.ModeratorScopeKey)
	}

	return r.removeScopeFromUser(userID, models.ModeratorScopeKey)
}

func (r *SqlUserRepository) addScopeToUser(userID string, scope string) error {
	u := r.GetUserByID(userID)
	if u == nil {
		return errors.New("user not found when modifying scope")
	}

	scopesString := u.Scopes
	scopes := utils.StringSliceToMap(scopesString)
	scopes[scope] = true

	scopesSlice := utils.StringMapKeys(scopes)

	return r.setScopesOnUser(userID, scopesSlice)
}

func (r *SqlUserRepository) removeScopeFromUser(userID string, scope string) error {
	u := r.GetUserByID(userID)
	scopesString := u.Scopes
	scopes := utils.StringSliceToMap(scopesString)
	delete(scopes, scope)

	scopesSlice := utils.StringMapKeys(scopes)

	return r.setScopesOnUser(userID, scopesSlice)
}

func (r *SqlUserRepository) setScopesOnUser(userID string, scopes []string) error {
	r.datastore.DbLock.Lock()
	defer r.datastore.DbLock.Unlock()

	tx, err := r.datastore.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback() //nolint

	scopesSliceString := strings.TrimSpace(strings.Join(scopes, ","))
	stmt, err := tx.Prepare("UPDATE users SET scopes=? WHERE id IS ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	var val *string
	if scopesSliceString == "" {
		val = nil
	} else {
		val = &scopesSliceString
	}

	if _, err := stmt.Exec(val, userID); err != nil {
		return err
	}

	return tx.Commit()
}

// GetUserByID will return a user by a user ID.
func (r *SqlUserRepository) GetUserByID(id string) *models.User {
	r.datastore.DbLock.Lock()
	defer r.datastore.DbLock.Unlock()

	query := "SELECT id, display_name, display_color, created_at, disabled_at, previous_names, namechanged_at, scopes FROM users WHERE id = ?"
	row := r.datastore.DB.QueryRow(query, id)
	if row == nil {
		log.Errorln(row)
		return nil
	}
	return r.getUserFromRow(row)
}

// GetDisabledUsers will return back all the currently disabled users that are not API users.
func (r *SqlUserRepository) GetDisabledUsers() []*models.User {
	query := "SELECT id, display_name, scopes, display_color, created_at, disabled_at, previous_names, namechanged_at FROM users WHERE disabled_at IS NOT NULL AND type IS NOT 'API'"

	rows, err := r.datastore.DB.Query(query)
	if err != nil {
		log.Errorln(err)
		return nil
	}
	defer rows.Close()

	users := r.getUsersFromRows(rows)

	sort.Slice(users, func(i, j int) bool {
		return users[i].DisabledAt.Before(*users[j].DisabledAt)
	})

	return users
}

// GetModeratorUsers will return a list of users with moderator access.
func (r *SqlUserRepository) GetModeratorUsers() []*models.User {
	query := `SELECT id, display_name, scopes, display_color, created_at, disabled_at, previous_names, namechanged_at FROM (
		WITH RECURSIVE split(id, display_name, scopes, display_color, created_at, disabled_at, previous_names, namechanged_at, scope, rest) AS (
		  SELECT id, display_name, scopes, display_color, created_at, disabled_at, previous_names, namechanged_at, '', scopes || ',' FROM users
		   UNION ALL
		  SELECT id, display_name, scopes, display_color, created_at, disabled_at, previous_names, namechanged_at,
				 substr(rest, 0, instr(rest, ',')),
				 substr(rest, instr(rest, ',')+1)
			FROM split
		   WHERE rest <> '')
		SELECT id, display_name, scopes, display_color, created_at, disabled_at, previous_names, namechanged_at, scope
		  FROM split
		 WHERE scope <> ''
		 ORDER BY created_at
	  ) AS token WHERE token.scope = ?`

	rows, err := r.datastore.DB.Query(query, models.ModeratorScopeKey)
	if err != nil {
		log.Errorln(err)
		return nil
	}
	defer rows.Close()

	users := r.getUsersFromRows(rows)

	return users
}

func (r *SqlUserRepository) getUsersFromRows(rows *sql.Rows) []*models.User {
	users := make([]*models.User, 0)

	for rows.Next() {
		var id string
		var displayName string
		var displayColor int
		var createdAt time.Time
		var disabledAt *time.Time
		var previousUsernames string
		var userNameChangedAt *time.Time
		var scopesString *string

		if err := rows.Scan(&id, &displayName, &scopesString, &displayColor, &createdAt, &disabledAt, &previousUsernames, &userNameChangedAt); err != nil {
			log.Errorln("error creating collection of users from results", err)
			return nil
		}

		var scopes []string
		if scopesString != nil {
			scopes = strings.Split(*scopesString, ",")
		}

		user := &models.User{
			ID:            id,
			DisplayName:   displayName,
			DisplayColor:  displayColor,
			CreatedAt:     createdAt,
			DisabledAt:    disabledAt,
			PreviousNames: strings.Split(previousUsernames, ","),
			NameChangedAt: userNameChangedAt,
			Scopes:        scopes,
		}
		users = append(users, user)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].CreatedAt.Before(users[j].CreatedAt)
	})

	return users
}

func (r *SqlUserRepository) getUserFromRow(row *sql.Row) *models.User {
	var id string
	var displayName string
	var displayColor int
	var createdAt time.Time
	var disabledAt *time.Time
	var previousUsernames string
	var userNameChangedAt *time.Time
	var scopesString *string

	if err := row.Scan(&id, &displayName, &displayColor, &createdAt, &disabledAt, &previousUsernames, &userNameChangedAt, &scopesString); err != nil {
		return nil
	}

	var scopes []string
	if scopesString != nil {
		scopes = strings.Split(*scopesString, ",")
	}

	return &models.User{
		ID:            id,
		DisplayName:   displayName,
		DisplayColor:  displayColor,
		CreatedAt:     createdAt,
		DisabledAt:    disabledAt,
		PreviousNames: strings.Split(previousUsernames, ","),
		NameChangedAt: userNameChangedAt,
		Scopes:        scopes,
	}
}

// InsertExternalAPIUser will add a new API user to the database.
func (r *SqlUserRepository) InsertExternalAPIUser(token string, name string, color int, scopes []string) error {
	log.Traceln("Adding new API user")

	r.datastore.DbLock.Lock()
	defer r.datastore.DbLock.Unlock()

	scopesString := strings.Join(scopes, ",")
	id := shortid.MustGenerate()

	tx, err := r.datastore.DB.Begin()
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

	if err := r.addAccessTokenForUser(token, id); err != nil {
		return errors.Wrap(err, "unable to save access token for new external api user")
	}

	return nil
}

// DeleteExternalAPIUser will delete a token from the database.
func (r *SqlUserRepository) DeleteExternalAPIUser(token string) error {
	log.Traceln("Deleting access token")

	r.datastore.DbLock.Lock()
	defer r.datastore.DbLock.Unlock()

	tx, err := r.datastore.DB.Begin()
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
func (r *SqlUserRepository) GetExternalAPIUserForAccessTokenAndScope(token string, scope string) (*models.ExternalAPIUser, error) {
	// This will split the scopes from comma separated to individual rows
	// so we can efficiently find if a token supports a single scope.
	// This is SQLite specific, so if we ever support other database
	// backends we need to support other methods.
	query := `SELECT
  id,
	scopes,
  display_name,
  display_color,
  created_at,
  last_used
FROM
  user_access_tokens
  INNER JOIN (
    WITH RECURSIVE split(
      id,
      scopes,
      display_name,
      display_color,
      created_at,
      last_used,
      disabled_at,
      scope,
      rest
    ) AS (
      SELECT
        id,
        scopes,
        display_name,
        display_color,
        created_at,
        last_used,
        disabled_at,
        '',
        scopes || ','
      FROM
        users AS u
      UNION ALL
      SELECT
        id,
        scopes,
        display_name,
        display_color,
        created_at,
        last_used,
        disabled_at,
        substr(rest, 0, instr(rest, ',')),
        substr(rest, instr(rest, ',') + 1)
      FROM
        split
      WHERE
        rest <> ''
    )
    SELECT
      id,
      display_name,
      display_color,
      created_at,
      last_used,
      disabled_at,
      scopes,
      scope
    FROM
      split
    WHERE
      scope <> ''
  ) ON user_access_tokens.user_id = id
WHERE
  disabled_at IS NULL
  AND token = ?
  AND scope = ?;`

	row := r.datastore.DB.QueryRow(query, token, scope)
	integration, err := r.makeExternalAPIUserFromRow(row)

	return integration, err
}

// GetIntegrationNameForAccessToken will return the integration name associated with a specific access token.
func (r *SqlUserRepository) GetIntegrationNameForAccessToken(token string) *string {
	name, err := r.datastore.GetQueries().GetUserDisplayNameByToken(context.Background(), token)
	if err != nil {
		return nil
	}

	return &name
}

// GetExternalAPIUser will return all API users with access tokens.
func (r *SqlUserRepository) GetExternalAPIUser() ([]models.ExternalAPIUser, error) { //nolint
	query := "SELECT id, token, display_name, display_color, scopes, created_at, last_used FROM users, user_access_tokens WHERE user_access_tokens.user_id = id  AND type IS 'API' AND disabled_at IS NULL"

	rows, err := r.datastore.DB.Query(query)
	if err != nil {
		return []models.ExternalAPIUser{}, err
	}
	defer rows.Close()

	integrations, err := r.makeExternalAPIUsersFromRows(rows)

	return integrations, err
}

// SetExternalAPIUserAccessTokenAsUsed will update the last used timestamp for a token.
func (r *SqlUserRepository) SetExternalAPIUserAccessTokenAsUsed(token string) error {
	tx, err := r.datastore.DB.Begin()
	if err != nil {
		return err
	}
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

func (r *SqlUserRepository) makeExternalAPIUserFromRow(row *sql.Row) (*models.ExternalAPIUser, error) {
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

	integration := models.ExternalAPIUser{
		ID:           id,
		DisplayName:  displayName,
		DisplayColor: displayColor,
		CreatedAt:    createdAt,
		Scopes:       strings.Split(scopes, ","),
		LastUsedAt:   lastUsedAt,
	}

	return &integration, nil
}

func (r *SqlUserRepository) makeExternalAPIUsersFromRows(rows *sql.Rows) ([]models.ExternalAPIUser, error) {
	integrations := make([]models.ExternalAPIUser, 0)

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

		integration := models.ExternalAPIUser{
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
func (r *SqlUserRepository) HasValidScopes(scopes []string) bool {
	// For a scope to be seen as "valid" it must live in this slice.
	validAccessTokenScopes := []string{
		models.ScopeCanSendChatMessages,
		models.ScopeCanSendSystemMessages,
		models.ScopeHasAdminAccess,
	}

	for _, scope := range scopes {
		_, foundInSlice := utils.FindInSlice(validAccessTokenScopes, scope)
		if !foundInSlice {
			return false
		}
	}
	return true
}

// GetUsersCount will return the number of users in the database.
func (r *SqlUserRepository) GetUsersCount() int {
	query := `SELECT COUNT(*) FROM users`
	rows, err := r.datastore.DB.Query(query)
	if err != nil || rows.Err() != nil {
		return 0
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0
		}
	}
	return count
}
