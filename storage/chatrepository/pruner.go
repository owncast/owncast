package chatrepository

import (
	"fmt"
	"time"

	"github.com/owncast/owncast/storage/data"
	log "github.com/sirupsen/logrus"
)

const (
	maxBacklogHours  = 2  // Keep backlog max hours worth of messages
	maxBacklogNumber = 50 // Return max number of messages in history request
)

func (cr *ChatRepository) startPruner() {
	chatDataPruner := time.NewTicker(5 * time.Minute)
	go func() {
		cr.runPruner()
		for range chatDataPruner.C {
			cr.runPruner()
		}
	}()
}

// Only keep recent messages so we don't keep more chat data than needed
// for privacy and efficiency reasons.
func (cr *ChatRepository) runPruner() {
	_datastore := data.GetDatastore()

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
