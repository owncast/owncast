package user

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/utils"
	"github.com/teris-io/shortid"

	log "github.com/sirupsen/logrus"
)

var _datastore *data.Datastore

const moderatorScopeKey = "MODERATOR"

// User represents a single chat user.
type User struct {
	ID            string     `json:"id"`
	AccessToken   string     `json:"-"`
	DisplayName   string     `json:"displayName"`
	DisplayColor  int        `json:"displayColor"`
	CreatedAt     time.Time  `json:"createdAt"`
	DisabledAt    *time.Time `json:"disabledAt,omitempty"`
	PreviousNames []string   `json:"previousNames"`
	NameChangedAt *time.Time `json:"nameChangedAt,omitempty"`
	Scopes        []string   `json:"scopes"`
}

// IsEnabled will return if this single user is enabled.
func (u *User) IsEnabled() bool {
	return u.DisabledAt == nil
}

// IsModerator will return if the user has moderation privileges.
func (u *User) IsModerator() bool {
	_, hasModerationScope := utils.FindInSlice(u.Scopes, moderatorScopeKey)
	return hasModerationScope
}

// SetupUsers will perform the initial initialization of the user package.
func SetupUsers() {
	_datastore = data.GetDatastore()
}

// CreateAnonymousUser will create a new anonymous user with the provided display name.
func CreateAnonymousUser(displayName string) (*User, error) {
	id := shortid.MustGenerate()
	accessToken, err := utils.GenerateAccessToken()
	if err != nil {
		log.Errorln("Unable to create access token for new user")
		return nil, err
	}

	if displayName == "" {
		suggestedUsernamesList := data.GetSuggestedUsernamesList()

		if len(suggestedUsernamesList) != 0 {
			index := utils.RandomIndex(len(suggestedUsernamesList))
			displayName = suggestedUsernamesList[index]
		} else {
			displayName = utils.GeneratePhrase()
		}
	}

	displayColor := utils.GenerateRandomDisplayColor()

	user := &User{
		ID:           id,
		AccessToken:  accessToken,
		DisplayName:  displayName,
		DisplayColor: displayColor,
		CreatedAt:    time.Now(),
	}

	if err := create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangeUsername will change the user associated to userID from one display name to another.
func ChangeUsername(userID string, username string) {
	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	tx, err := _datastore.DB.Begin()
	if err != nil {
		log.Debugln(err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Debugln(err)
		}
	}()

	stmt, err := tx.Prepare("UPDATE users SET display_name = ?, previous_names = previous_names || ?, namechanged_at = ? WHERE id = ?")
	if err != nil {
		log.Debugln(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, fmt.Sprintf(",%s", username), time.Now(), userID)
	if err != nil {
		log.Errorln(err)
	}

	if err := tx.Commit(); err != nil {
		log.Errorln("error changing display name of user", userID, err)
	}
}

func create(user *User) error {
	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	tx, err := _datastore.DB.Begin()
	if err != nil {
		log.Debugln(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare("INSERT INTO users(id, access_token, display_name, display_color, previous_names, created_at) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Debugln(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID, user.AccessToken, user.DisplayName, user.DisplayColor, user.DisplayName, user.CreatedAt)
	if err != nil {
		log.Errorln("error creating new user", err)
	}

	return tx.Commit()
}

// SetEnabled will set the enabled status of a single user by ID.
func SetEnabled(userID string, enabled bool) error {
	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	tx, err := _datastore.DB.Begin()
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
func GetUserByToken(token string) *User {
	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	query := "SELECT id, display_name, display_color, created_at, disabled_at, previous_names, namechanged_at, scopes FROM users WHERE access_token = ?"
	row := _datastore.DB.QueryRow(query, token)

	return getUserFromRow(row)
}

// SetModerator will add or remove moderator status for a single user by ID.
func SetModerator(userID string, isModerator bool) error {
	if isModerator {
		return addScopeToUser(userID, moderatorScopeKey)
	}

	return removeScopeFromUser(userID, moderatorScopeKey)
}

func addScopeToUser(userID string, scope string) error {
	u := GetUserByID(userID)
	scopesString := u.Scopes
	scopes := utils.StringSliceToMap(scopesString)
	scopes[scope] = true

	scopesSlice := utils.StringMapKeys(scopes)

	return setScopesOnUser(userID, scopesSlice)
}

func removeScopeFromUser(userID string, scope string) error {
	u := GetUserByID(userID)
	scopesString := u.Scopes
	scopes := utils.StringSliceToMap(scopesString)
	delete(scopes, scope)

	scopesSlice := utils.StringMapKeys(scopes)

	return setScopesOnUser(userID, scopesSlice)
}

func setScopesOnUser(userID string, scopes []string) error {
	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	tx, err := _datastore.DB.Begin()
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
func GetUserByID(id string) *User {
	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	query := "SELECT id, display_name, display_color, created_at, disabled_at, previous_names, namechanged_at, scopes FROM users WHERE id = ?"
	row := _datastore.DB.QueryRow(query, id)
	if row == nil {
		log.Errorln(row)
		return nil
	}
	return getUserFromRow(row)
}

// GetDisabledUsers will return back all the currently disabled users that are not API users.
func GetDisabledUsers() []*User {
	query := "SELECT id, display_name, scopes, display_color, created_at, disabled_at, previous_names, namechanged_at FROM users WHERE disabled_at IS NOT NULL AND type IS NOT 'API'"

	rows, err := _datastore.DB.Query(query)
	if err != nil {
		log.Errorln(err)
		return nil
	}
	defer rows.Close()

	users := getUsersFromRows(rows)

	sort.Slice(users, func(i, j int) bool {
		return users[i].DisabledAt.Before(*users[j].DisabledAt)
	})

	return users
}

// GetModeratorUsers will return a list of users with moderator access.
func GetModeratorUsers() []*User {
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

	rows, err := _datastore.DB.Query(query, moderatorScopeKey)
	if err != nil {
		log.Errorln(err)
		return nil
	}
	defer rows.Close()

	users := getUsersFromRows(rows)

	return users
}

func getUsersFromRows(rows *sql.Rows) []*User {
	users := make([]*User, 0)

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

		user := &User{
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

func getUserFromRow(row *sql.Row) *User {
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

	return &User{
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
