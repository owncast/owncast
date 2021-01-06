package core

import (
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/storageproviders"
)

func setupStorage() error {
	handler.Storage = _storage
	s3Config := data.GetS3Config()

	if s3Config.Enabled {
		_storage = &storageproviders.S3Storage{}
	} else {
		_storage = &storageproviders.LocalStorage{}
	}

	if err := _storage.Setup(); err != nil {
		return err
	}

	return nil
}
