package stream

import (
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/storage"
)

// segmentProtectorSetter is implemented by storage providers that can be told
// which recorded segments to keep during cleanup.
type segmentProtectorSetter interface {
	SetSegmentProtector(protector models.SegmentProtector)
}

// setupStorage picks an HLS storage backend based on the S3 config and
// wires it into the HLS handler. Called once from Start().
func (s *Service) setupStorage() error {
	s3Config := s.configRepository.GetS3Config()

	if s3Config.Enabled {
		s.storage = storage.NewS3Storage(s.configRepository)
	} else {
		s.storage = storage.NewLocalStorage(s.configRepository)
	}

	// Cleanup keeps running while replays are enabled, so it needs to know
	// which segments a clip or replay still references.
	if setter, ok := s.storage.(segmentProtectorSetter); ok && s.replays != nil {
		setter.SetSegmentProtector(s.replays)
	}

	if err := s.storage.Setup(); err != nil {
		return err
	}

	s.handler.Storage = s.storage

	return nil
}

// protectedSegmentFilenames asks the replay ledger which segments must be kept
// during a cleanup or directory wipe. With no replay subsystem, or when the
// ledger cannot answer, nothing is protected.
func (s *Service) protectedSegmentFilenames() map[string]bool {
	if s.replays == nil {
		return nil
	}

	protected, err := s.replays.ProtectedSegmentFilenames()
	if err != nil {
		log.Warnln("unable to determine which video segments are in use:", err)
		return nil
	}

	return protected
}

// replayProtector returns the replay ledger as a segment protector, or nil
// when the replay subsystem is absent.
func (s *Service) replayProtector() models.SegmentProtector {
	if s.replays == nil {
		return nil
	}

	return s.replays
}
