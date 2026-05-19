package chat

import (
	"time"

	"github.com/owncast/owncast/services/datastore"
)

var _datastore *datastore.Datastore

const (
	maxBacklogHours = 2 // Keep backlog max hours worth of messages
)

func setupPersistence() {
	_datastore = datastore.GetDatastore()

	chatDataPruner := time.NewTicker(5 * time.Minute)
	go func() {
		runPruner()
		for range chatDataPruner.C {
			runPruner()
		}
	}()
}
