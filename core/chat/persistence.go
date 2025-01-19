package chat

import (
	"time"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/persistence/authrepository"
	"github.com/owncast/owncast/persistence/tables"
)

var _datastore *data.Datastore

const (
	maxBacklogHours = 2 // Keep backlog max hours worth of messages
)

func setupPersistence() {
	_datastore = data.GetDatastore()
	tables.CreateMessagesTable(_datastore.DB)

	authRepository := authrepository.Get()
	authRepository.CreateBanIPTable(_datastore.DB)

	chatDataPruner := time.NewTicker(5 * time.Minute)
	go func() {
		runPruner()
		for range chatDataPruner.C {
			runPruner()
		}
	}()
}
