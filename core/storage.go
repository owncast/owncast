package core

import (
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/storageproviders"
)

func setupStorage() error {
	s3Config := data.GetS3Config()

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
