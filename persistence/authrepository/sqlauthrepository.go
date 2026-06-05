package authrepository

import (
	"github.com/owncast/owncast/services/datastore"
)

type SqlAuthRepository struct {
	datastore *datastore.Datastore
}

// New will create a new instance of the UserRepository.
func New(datastore *datastore.Datastore) *SqlAuthRepository {
	r := &SqlAuthRepository{
		datastore: datastore,
	}

	return r
}
