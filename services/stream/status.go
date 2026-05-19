package stream

import (
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
)

// GetStatus returns a snapshot of the current stream state suitable for
// the public/admin status APIs.
func (s *Service) GetStatus() models.Status {
	if s.stats == nil {
		return models.Status{}
	}

	viewerCount := 0
	if s.IsStreamConnected() {
		viewerCount = len(s.stats.Viewers)
	}

	return models.Status{
		Online:                s.IsStreamConnected(),
		ViewerCount:           viewerCount,
		OverallMaxViewerCount: s.stats.OverallMaxViewerCount,
		SessionMaxViewerCount: s.stats.SessionMaxViewerCount,
		LastDisconnectTime:    s.stats.LastDisconnectTime,
		LastConnectTime:       s.stats.LastConnectTime,
		VersionNumber:         config.VersionNumber,
		StreamTitle:           s.configRepository.GetStreamTitle(),
	}
}

// GetCurrentBroadcast returns the in-flight broadcast settings, or nil
// if no stream is currently live.
func (s *Service) GetCurrentBroadcast() *models.CurrentBroadcast {
	return s.currentBroadcast
}

// setBroadcaster records the metadata of the inbound RTMP source. Called
// by the RTMP server's on-connect callback.
func (s *Service) setBroadcaster(broadcaster models.Broadcaster) {
	s.broadcaster = &broadcaster
}

// GetBroadcaster returns the active inbound broadcaster, or nil between
// streams.
func (s *Service) GetBroadcaster() *models.Broadcaster {
	return s.broadcaster
}
