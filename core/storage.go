package core

import (
	"fmt"

	"github.com/owncast/owncast/core/storageproviders"
)

func (s *Service) setupStorage() error {
	s.Storage = storageproviders.NewLocalStorage(s.Transcoder)

	if cfgS3 := s.ActivityPub.Persistence.Data.GetS3Config(); cfgS3.Enabled {
		s.Storage = storageproviders.NewS3Storage(s.ActivityPub.Persistence.Data)
	}

	if err := s.Storage.Setup(); err != nil {
		return fmt.Errorf("setting up Storage: %v", err)
	}

	s.handler.Storage = s.Storage

	return nil
}
