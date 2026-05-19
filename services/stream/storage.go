package stream

import (
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/storage"
)

// setupStorage picks an HLS storage backend based on the S3 config and
// wires it into the HLS handler. Called once from Start().
func (s *Service) setupStorage() error {
	configRepository := configrepository.Get()
	s3Config := configRepository.GetS3Config()

	if s3Config.Enabled {
		s.storage = storage.NewS3Storage()
	} else {
		s.storage = storage.NewLocalStorage()
	}

	if err := s.storage.Setup(); err != nil {
		return err
	}

	s.handler.Storage = s.storage

	return nil
}
