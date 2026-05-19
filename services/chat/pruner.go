package chat

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Only keep recent messages so we don't keep more chat data than needed
// for privacy and efficiency reasons.
func runPruner() {
	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	log.Traceln("Removing chat messages older than", maxBacklogHours, "hours")

	deleteStatement := `DELETE FROM messages WHERE timestamp <= datetime('now', 'localtime', ?)`
	tx, err := _datastore.DB.Begin()
	if err != nil {
		log.Debugln(err)
		return
	}

	stmt, err := tx.Prepare(deleteStatement)
	if err != nil {
		log.Debugln(err)
		return
	}
	defer stmt.Close()

	if _, err = stmt.Exec(fmt.Sprintf("-%d hours", maxBacklogHours)); err != nil {
		log.Debugln(err)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Debugln(err)
		return
	}
}
