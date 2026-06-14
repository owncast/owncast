package stream

import (
	"github.com/owncast/owncast/services/storage"
)

// setupStorage picks an HLS storage backend based on the S3 config and
// wires it into the HLS handler. Called once from Start().
func (s *Service) setupStorage() error {
	s3Config := s.configRepository.GetS3Config()

	if s3Config.Enabled {
		s.storage = storage.NewS3Storage(s.configRepository)
	} else {
		s.storage = storage.NewLocalStorage(s.configRepository)
	}

	if err := s.storage.Setup(); err != nil {
		return err
	}

	s.handler.Storage = s.storage

	return nil
}
