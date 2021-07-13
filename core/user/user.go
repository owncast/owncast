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

type User struct {
	Id            string     `json:"id"`
	AccessToken   string     `json:"-"`
	DisplayName   string     `json:"displayName"`
	DisplayColor  int        `json:"displayColor"`
	CreatedAt     time.Time  `json:"createdAt"`
	DisabledAt    *time.Time `json:"disabledAt,omitempty"`
	PreviousNames []string   `json:"previousNames"`
	NameChangedAt *time.Time `json:"nameChangedAt,omitempty"`
}

func (u *User) IsEnabled() bool {
	return u.DisabledAt == nil
}

func SetupUsers() {
	_datastore = data.GetDatastore()
}

func CreateAnonymousUser(username string) (*User, error) {
	id := shortid.MustGenerate()
	accessToken, err := utils.GenerateAccessToken()
	if err != nil {
		log.Errorln("Unable to create access token for new user")
		return nil, err
	}

	var displayName = username
	if displayName == "" {
		displayName = utils.GeneratePhrase()
	}

	displayColor := utils.GenerateRandomDisplayColor()

	user := &User{
		Id:           id,
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

func ChangeUsername(userId string, username string) {
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

	_, err = stmt.Exec(username, fmt.Sprintf(",%s", username), time.Now(), userId)
	if err != nil {
		log.Errorln(err)
	}

	if err := tx.Commit(); err != nil {
		log.Errorln("error changing display name of user", userId, err)
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

	log.Traceln("Creating new user", user.Id, user.DisplayName)

	stmt, err := tx.Prepare("INSERT INTO users(id, access_token, display_name, display_color, previous_names, created_at) values(?, ?, ?, ?, ?, ?)")

	if err != nil {
		log.Debugln(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id, user.AccessToken, user.DisplayName, user.DisplayColor, user.DisplayName, user.CreatedAt)
	if err != nil {
		log.Errorln("error creating new user", err)
	}

	return tx.Commit()
}

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

	query := "SELECT id, display_name, display_color, created_at, disabled_at, previous_names, namechanged_at FROM users WHERE access_token = ?"
	row := _datastore.DB.QueryRow(query, token)

	return getUserFromRow(row)
}

// GetUserById will return a user by a user ID.
func GetUserById(id string) *User {
	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	query := "SELECT id, display_name, display_color, created_at, disabled_at, previous_names, namechanged_at FROM users WHERE id = ?"
	row := _datastore.DB.QueryRow(query, id)
	if row == nil {
		log.Errorln(row)
		return nil
	}
	return getUserFromRow(row)
}

// GetDisabledUsers will return back all the currently disabled users that are not API users.
func GetDisabledUsers() []*User {
	query := "SELECT id, display_name, display_color, created_at, disabled_at, previous_names, namechanged_at FROM users WHERE disabled_at IS NOT NULL AND type IS NOT 'API'"

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

		if err := rows.Scan(&id, &displayName, &displayColor, &createdAt, &disabledAt, &previousUsernames, &userNameChangedAt); err != nil {
			log.Errorln("error creating collection of users from results", err)
			return nil
		}

		user := &User{
			Id:            id,
			DisplayName:   displayName,
			DisplayColor:  displayColor,
			CreatedAt:     createdAt,
			DisabledAt:    disabledAt,
			PreviousNames: strings.Split(previousUsernames, ","),
			NameChangedAt: userNameChangedAt,
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

	if err := row.Scan(&id, &displayName, &displayColor, &createdAt, &disabledAt, &previousUsernames, &userNameChangedAt); err != nil {
		return nil
	}

	return &User{
		Id:            id,
		DisplayName:   displayName,
		DisplayColor:  displayColor,
		CreatedAt:     createdAt,
		DisabledAt:    disabledAt,
		PreviousNames: strings.Split(previousUsernames, ","),
		NameChangedAt: userNameChangedAt,
	}
}
