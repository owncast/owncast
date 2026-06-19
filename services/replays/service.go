// Package replays records live streams to disk and database so they can be
// played back later in whole (a replay) or in part (a clip).
//
// The recording side (HLSRecorder) is driven by the stream lifecycle: as the
// transcoder writes HLS segments, each one is noted in the database along with
// the output configuration it belongs to and its offset from the start of the
// stream. The playback side (PlaylistGenerator) reconstructs HLS master and
// media playlists on demand from that recorded metadata.
//
// The whole subsystem is gated behind the config.EnableReplayFeatures flag and
// is disabled by default.
package replays

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/datastore"
)

// Service is the replay subsystem. Construct one with New in the composition
// root and hand it to the stream service (for recording) and the web handlers
// (for playback and clip creation).
type Service struct {
	datastore        *datastore.Datastore
	configRepository configrepository.ConfigRepository
}

// New constructs a replay Service backed by the provided datastore and config
// repository.
func New(ds *datastore.Datastore, configRepository configrepository.ConfigRepository) *Service {
	return &Service{
		datastore:        ds,
		configRepository: configRepository,
	}
}

// Setup performs one-time startup work for the replay subsystem.
func (s *Service) Setup() {
	s.fixUnfinishedStreams()
}

// fixUnfinishedStreams will find streams with no end time (left over from a
// previous run that was interrupted) and give them an end time based on the
// last segment recorded for that stream.
func (s *Service) fixUnfinishedStreams() {
	if err := s.datastore.GetQueries().FixUnfinishedStreams(context.Background()); err != nil {
		log.Warnln(err)
	}
}

// NewRecording returns a recorder that will persist the segments of the stream
// identified by streamID. Returns nil for the offline placeholder stream,
// which is never recorded.
func (s *Service) NewRecording(streamID string) *HLSRecorder {
	return newRecording(streamID, s.datastore, s.configRepository)
}

// NewPlaylistGenerator returns a generator capable of reconstructing replay and
// clip playlists from recorded metadata.
func (s *Service) NewPlaylistGenerator() *PlaylistGenerator {
	return &PlaylistGenerator{datastore: s.datastore}
}
