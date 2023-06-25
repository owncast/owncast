package configrepository

import (
	"testing"

	"github.com/owncast/owncast/storage/data"
	"github.com/owncast/owncast/utils"
)

var (
	_datastore        *data.Store
	_configRepository *SqlConfigRepository
)

func TestMain(m *testing.M) {
	ds, err := data.NewStore(":memory")
	if err != nil {
		panic(err)
	}
	_datastore = ds

	_configRepository = New(_datastore)

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
