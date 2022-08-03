package data

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// MustExec will execute a SQL statement on a provided database instance.
func MustExec(s string, db *sql.DB) {
	stmt, err := db.Prepare(s)
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Warnln(err)
	}
}
