package auth

import (
	"github.com/owncast/owncast/core/data"
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
	);`
	_datastore.MustExec(createTableSQL)
	_datastore.MustExec(`CREATE INDEX IF NOT EXISTS idx_auth_token ON auth (token);`)
}
