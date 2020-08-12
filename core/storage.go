package core

import (
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/storageproviders"
)

var (
	usingExternalStorage = false
)

func setupStorage() error {
	if config.Config.S3.Enabled {
		_storage = &storageproviders.S3Storage{}
		usingExternalStorage = true
	}

	if usingExternalStorage {
		if err := _storage.Setup(); err != nil {
			return err
		}

		// go playlist.StartVideoContentMonitor(_storage)
	}

	return nil
}
