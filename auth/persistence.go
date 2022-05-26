package auth

import (
	"context"
	"strings"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/user"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/db"
)

var _datastore *data.Datastore

// Setup will initialize auth persistence.
func Setup(db *data.Datastore) {
	_datastore = db

	createTableSQL := `CREATE TABLE IF NOT EXISTS auth (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"user_id" TEXT NOT NULL,
    "token" TEXT NOT NULL,
    "type" TEXT NOT NULL,
		"timestamp" DATE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id)
	);CREATE INDEX auth_token ON auth (token);`

	stmt, err := db.DB.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		log.Fatalln(err)
	}
}

// AddAuth will add an external authentication token and type for a user.
func AddAuth(userID, authToken string, authType Type) error {
	return _datastore.GetQueries().AddAuthForUser(context.Background(), db.AddAuthForUserParams{
		UserID: userID,
		Token:  authToken,
		Type:   string(authType),
	})
}

// GetUserByAuth will return an existing user given auth details if a user
// has previously authenticated with that method.
func GetUserByAuth(authToken string, authType Type) *user.User {
	u, err := _datastore.GetQueries().GetUserByAuth(context.Background(), db.GetUserByAuthParams{
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

	return &user.User{
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
