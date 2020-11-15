package core

import (
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/storageproviders"
)

func setupStorage() error {
	handler.Storage = _storage

	if config.Config.S3.Enabled {
		_storage = &storageproviders.S3Storage{}
	} else {
		_storage = &storageproviders.LocalStorage{}
	}

	if err := _storage.Setup(); err != nil {
		return err
	}

	return nil
}
