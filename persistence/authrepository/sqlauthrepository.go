package authrepository

import (
	"github.com/owncast/owncast/services/datastore"
)

type SqlAuthRepository struct {
	datastore *datastore.Datastore
}

// NOTE: This is temporary during the transition period.
var temporaryGlobalInstance AuthRepository

// Get will return the user repository.
func Get() AuthRepository {
	if temporaryGlobalInstance == nil {
		i := New(datastore.GetDatastore())
		temporaryGlobalInstance = i
	}
	return temporaryGlobalInstance
}

// New will create a new instance of the UserRepository.
func New(datastore *datastore.Datastore) *SqlAuthRepository {
	r := &SqlAuthRepository{
		datastore: datastore,
	}

	return r
}
