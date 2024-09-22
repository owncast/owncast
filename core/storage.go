package core

import (
	"github.com/owncast/owncast/core/storageproviders"
	"github.com/owncast/owncast/persistence/configrepository"
)

func setupStorage() error {
	configRepository := configrepository.Get()
	s3Config := configRepository.GetS3Config()

	if s3Config.Enabled {
		_storage = storageproviders.NewS3Storage()
	} else {
		_storage = storageproviders.NewLocalStorage()
	}

	if err := _storage.Setup(); err != nil {
		return err
	}

	handler.Storage = _storage

	return nil
}
