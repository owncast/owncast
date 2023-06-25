package core

import (
	"github.com/owncast/owncast/storage/configrepository"
	"github.com/owncast/owncast/video/storageproviders"
)

func setupStorage() error {
	config := configrepository.Get()
	s3Config := config.GetS3Config()

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
