package userrepository

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/services/datastore"

	"github.com/pkg/errors"
	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"

	log "github.com/sirupsen/logrus"
)

type UserRepository interface {
	ChangeUserColor(userID string, color int) error
	ChangeUsername(userID string, username string) error
	CreateAnonymousUser(displayName string) (*models.User, string, error)
	DeleteExternalAPIUser(token string) error
	GetDisabledUsers() []*models.User
	GetUsers() []*models.User
	GetUsersPaginated(offset int, limit int, search string, status string) ([]*models.User, int, error)
	DeleteUser(userID string) error
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
	UserRegisteredByPlugin(pluginName, userID string) bool
	AddAccessTokenForUser(accessToken, userID string) error
	SetUserScopes(userID string, scopes []string) error
	SetExternalAPIUserAccessTokenAsUsed(token string) error
	GetUsersCount() int
}

type SqlUserRepository struct {
	datastore *datastore.Datastore
}

// New will create a new instance of the UserRepository.
func New(datastore *datastore.Datastore) UserRepository {
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
		DisplayColor: int64(color),
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

// AddAccessTokenForUser associates a new access token with an existing user.
// Used by the viewer-auth gate to mint a session token the gate cookie carries.
func (r *SqlUserRepository) AddAccessTokenForUser(accessToken, userID string) error {
	return r.addAccessTokenForUser(accessToken, userID)
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

// userFromColumns builds a *models.User from the column set shared by the
// sqlc user lookups (GetUserByAccessToken, GetUserByID, GetUsers). IsBot is
// derived in SQL (users.type = 'API') so callers don't second-guess it.
func userFromColumns(
	id, displayName string,
	displayColor int64,
	createdAt, disabledAt sql.NullTime,
	previousNames sql.NullString,
	namechangedAt, authenticatedAt sql.NullTime,
	scopes sql.NullString,
	isBot bool,
) *models.User {
	var scopeSlice []string
	if scopes.Valid {
		scopeSlice = strings.Split(scopes.String, ",")
	}

	var disabled *time.Time
	if disabledAt.Valid {
		disabled = &disabledAt.Time
	}

	var authAt *time.Time
	if authenticatedAt.Valid {
		authAt = &authenticatedAt.Time
	}

	return &models.User{
		ID:              id,
		DisplayName:     displayName,
		DisplayColor:    int(displayColor),
		CreatedAt:       createdAt.Time,
		DisabledAt:      disabled,
		PreviousNames:   strings.Split(previousNames.String, ","),
		NameChangedAt:   &namechangedAt.Time,
		AuthenticatedAt: authAt,
		Authenticated:   authAt != nil,
		Scopes:          scopeSlice,
		IsBot:           isBot,
	}
}

// GetUserByToken will return a user by an access token.
func (r *SqlUserRepository) GetUserByToken(token string) *models.User {
	u, err := r.datastore.GetQueries().GetUserByAccessToken(context.Background(), token)
	if err != nil {
		return nil
	}
	return userFromColumns(u.ID, u.DisplayName, u.DisplayColor, u.CreatedAt, u.DisabledAt, u.PreviousNames, u.NamechangedAt, u.AuthenticatedAt, u.Scopes, u.IsBot)
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

// UserRegisteredByPlugin reports whether userID belongs to a user the named
// plugin registered via owncast.users.register — i.e. the user has a
// plugin.auth identity namespaced to that plugin's slug (RegisterUser stores
// it as "<slug>:<authId>"). The viewer-auth gate uses this to confine
// owncast.auth.grantSession to a plugin's own users, so a gate plugin can't
// mint a session impersonating an arbitrary existing user (e.g. a moderator).
//
// Matching on the "<slug>:" prefix is safe: plugin slugs are validated to
// lowercase letters/digits/hyphens, so they can't carry SQL LIKE wildcards.
// Fails closed (returns false) on any query error.
func (r *SqlUserRepository) UserRegisteredByPlugin(pluginName, userID string) bool {
	count, err := r.datastore.GetQueries().CountUserAuthByTypeAndTokenPrefix(context.Background(), db.CountUserAuthByTypeAndTokenPrefixParams{
		UserID: userID,
		Type:   string(models.PluginAuth),
		Token:  pluginName + ":%",
	})
	if err != nil {
		log.Errorln("checking plugin user ownership:", err)
		return false
	}
	return count > 0
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

// SetUserScopes replaces a user's full scope set. Used by viewer-auth plugins
// (owncast.users.register) to map an external provider's roles onto Owncast
// scopes (e.g. granting MODERATOR). Callers should validate with
// HasValidScopes first.
func (r *SqlUserRepository) SetUserScopes(userID string, scopes []string) error {
	return r.setScopesOnUser(userID, scopes)
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
	u, err := r.datastore.GetQueries().GetUserByID(context.Background(), id)
	if err != nil {
		return nil
	}
	return userFromColumns(u.ID, u.DisplayName, u.DisplayColor, u.CreatedAt, u.DisabledAt, u.PreviousNames, u.NamechangedAt, u.AuthenticatedAt, u.Scopes, u.IsBot)
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

// GetUsers will return all users, most-recently-created first.
func (r *SqlUserRepository) GetUsers() []*models.User {
	rows, err := r.datastore.GetQueries().GetUsers(context.Background())
	if err != nil {
		log.Errorln(err)
		return nil
	}

	users := make([]*models.User, 0, len(rows))
	for _, u := range rows {
		users = append(users, userFromColumns(u.ID, u.DisplayName, u.DisplayColor, u.CreatedAt, u.DisabledAt, u.PreviousNames, u.NamechangedAt, u.AuthenticatedAt, u.Scopes, u.IsBot))
	}
	return users
}

// GetUsersPaginated returns a page of users of every type (chat viewers,
// authenticated/plugin users, and API integrations), most-recently-created
// first, filtered to display names containing search and an optional status
// ("" / "all" = every user; else "active", "banned", "moderators", "bots"). It also
// returns the total number of users matching the filter so the admin
// user-management page can paginate. The search arg is always Valid (even when
// empty) so the LIKE binds to ” and matches every user rather than NULL.
func (r *SqlUserRepository) GetUsersPaginated(offset int, limit int, search string, status string) ([]*models.User, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	ctx := context.Background()
	searchArg := sql.NullString{String: search, Valid: true}

	total, err := r.datastore.GetQueries().CountUsers(ctx, db.CountUsersParams{
		Search: searchArg,
		Status: status,
	})
	if err != nil {
		return nil, 0, errors.Wrap(err, "unable to count users")
	}

	rows, err := r.datastore.GetQueries().GetUsersPaginated(ctx, db.GetUsersPaginatedParams{
		Search:     searchArg,
		Status:     status,
		PageLimit:  int64(limit),
		PageOffset: int64(offset),
	})
	if err != nil {
		return nil, 0, errors.Wrap(err, "unable to query users")
	}

	users := make([]*models.User, 0, len(rows))
	for _, u := range rows {
		users = append(users, userFromColumns(u.ID, u.DisplayName, u.DisplayColor, u.CreatedAt, u.DisabledAt, u.PreviousNames, u.NamechangedAt, u.AuthenticatedAt, u.Scopes, u.IsBot))
	}

	if err := r.attachAuthProviders(ctx, users); err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}

// attachAuthProviders fills each user's AuthProviders with friendly labels for
// the external auth methods they signed in with (IndieAuth, Fediverse, or a
// viewer-auth plugin's slug). One query covers the whole page; anonymous users
// have no auth rows and are left with an empty slice.
func (r *SqlUserRepository) attachAuthProviders(ctx context.Context, users []*models.User) error {
	if len(users) == 0 {
		return nil
	}

	ids := make([]string, 0, len(users))
	for _, u := range users {
		ids = append(ids, u.ID)
	}

	authRows, err := r.datastore.GetQueries().GetAuthForUsers(ctx, ids)
	if err != nil {
		return errors.Wrap(err, "unable to load user auth providers")
	}

	// Distinct provider labels per user id.
	byUser := map[string][]string{}
	seen := map[string]map[string]bool{}
	for _, a := range authRows {
		label := authProviderLabel(a.Type, a.Token)
		if seen[a.UserID] == nil {
			seen[a.UserID] = map[string]bool{}
		}
		if seen[a.UserID][label] {
			continue
		}
		seen[a.UserID][label] = true
		byUser[a.UserID] = append(byUser[a.UserID], label)
	}

	for _, u := range users {
		u.AuthProviders = byUser[u.ID]
	}

	return nil
}

// authProviderLabel turns an auth row into a human-friendly provider name. For
// plugin auth the token is namespaced as "<plugin-slug>:<externalId>", so the
// plugin slug is surfaced rather than the opaque "plugin.auth" type.
func authProviderLabel(authType, token string) string {
	switch models.AuthType(authType) {
	case models.IndieAuth:
		return "IndieAuth"
	case models.Fediverse:
		return "Fediverse"
	case models.PluginAuth:
		if i := strings.IndexByte(token, ':'); i > 0 {
			return token[:i]
		}
		return "Plugin"
	default:
		return authType
	}
}

// DeleteUser permanently removes a user along with everything tied to their
// identity: access tokens, external/plugin auth identities, and chat messages.
// Unlike SetEnabled(false) (a reversible ban), this cannot be undone. All
// removals run on a single transaction (via the sqlc Queries.WithTx) so a user
// is never left half-deleted.
func (r *SqlUserRepository) DeleteUser(userID string) error {
	r.datastore.DbLock.Lock()
	defer r.datastore.DbLock.Unlock()

	tx, err := r.datastore.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint

	q := r.datastore.GetQueries().WithTx(tx)
	ctx := context.Background()

	// Remove rows that reference the user before the user row itself.
	if err := q.DeleteUserAccessTokens(ctx, userID); err != nil {
		return errors.Wrap(err, "unable to delete user access tokens")
	}
	if err := q.DeleteUserAuth(ctx, userID); err != nil {
		return errors.Wrap(err, "unable to delete user auth")
	}
	if err := q.DeleteUserMessages(ctx, sql.NullString{String: userID, Valid: true}); err != nil {
		return errors.Wrap(err, "unable to delete user messages")
	}

	rowsDeleted, err := q.DeleteUserByID(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "unable to delete user")
	}
	if rowsDeleted == 0 {
		return errors.New("user " + userID + " not found")
	}

	return tx.Commit()
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
//
// Only true third-party API integrations (users.type = 'API', created via the
// admin Access Tokens UI) may authenticate to the scoped external API. Without
// the type filter, any user that happened to carry an admin scope on a regular
// access token — e.g. a user created and granted a session by a viewer-auth
// plugin — would also pass, which is not the intent of this Bearer-token API.
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
      WHERE
        u.type = 'API'
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
