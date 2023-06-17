package configrepository

import (
	"os"
	"testing"

	"github.com/owncast/owncast/storage/datastore"
	"github.com/owncast/owncast/utils"
)

var (
	_datastore        *datastore.Datastore
	_configRepository *SqlConfigRepository
)

func TestMain(m *testing.M) {
	dbFile, err := os.CreateTemp(os.TempDir(), ":memory:")
	if err != nil {
		panic(err)
	}

	ds, err := datastore.NewDatastore(dbFile.Name())
	if err != nil {
		panic(err)
	}
	_datastore = ds

	_configRepository = NewConfigRepository(_datastore)

	m.Run()
}

func TestSetName(t *testing.T) {
	value, _ := utils.GenerateRandomString(10)
	err := _configRepository.SetServerName(value)
	if err != nil {
		t.Error(err)
	}

	result := _configRepository.GetServerName()
	if result != value {
		t.Error("expected", value, "but test returned", result)
	}
}
