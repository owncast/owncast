package chat

import (
	"time"

	"github.com/owncast/owncast/core/data"
)

var _datastore *data.Datastore

const (
	maxBacklogHours = 2 // Keep backlog max hours worth of messages
)

func setupPersistence() {
	_datastore = data.GetDatastore()

	chatDataPruner := time.NewTicker(5 * time.Minute)
	go func() {
		runPruner()
		for range chatDataPruner.C {
			runPruner()
		}
	}()
}
