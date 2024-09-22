package authrepository

import (
	"github.com/owncast/owncast/core/data"
)

type SqlAuthRepository struct {
	datastore *data.Datastore
}

// NOTE: This is temporary during the transition period.
var temporaryGlobalInstance AuthRepository

// Get will return the user repository.
func Get() AuthRepository {
	if temporaryGlobalInstance == nil {
		i := New(data.GetDatastore())
		temporaryGlobalInstance = i
	}
	return temporaryGlobalInstance
}

// New will create a new instance of the UserRepository.
func New(datastore *data.Datastore) *SqlAuthRepository {
	r := &SqlAuthRepository{
		datastore: datastore,
	}

	return r
}
