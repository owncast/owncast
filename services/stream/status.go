package stream

import (
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
)

// GetStatus returns a snapshot of the current stream state suitable for
// the public/admin status APIs.
func (s *Service) GetStatus() models.Status {
	waitTime := float64(s.configRepository.GetStreamLatencyLevel().SecondsPerSegment) * 3.0
	if waitTime < 7 {
		waitTime = 7
	}

	s.statsMu.RLock()
	if s.stats == nil {
		s.statsMu.RUnlock()
		return models.Status{}
	}

	online := s.isStreamConnectedLocked(waitTime)
	viewerCount := 0
	if online {
		viewerCount = len(s.stats.Viewers)
	}
	status := models.Status{
		Online:                online,
		ViewerCount:           viewerCount,
		OverallMaxViewerCount: s.stats.OverallMaxViewerCount,
		SessionMaxViewerCount: s.stats.SessionMaxViewerCount,
		LastDisconnectTime:    s.stats.LastDisconnectTime,
		LastConnectTime:       s.stats.LastConnectTime,
		VersionNumber:         config.VersionNumber,
	}
	s.statsMu.RUnlock()

	status.StreamTitle = s.configRepository.GetStreamTitle()
	return status
}

// GetCurrentBroadcast returns the in-flight broadcast settings, or nil
// if no stream is currently live.
func (s *Service) GetCurrentBroadcast() *models.CurrentBroadcast {
	return s.currentBroadcast
}

// BroadcasterSet records the metadata of the inbound source. Implements
// StreamEvents; the local engine calls it from the RTMP metadata callback, a
// remote engine from its signaling channel.
func (s *Service) BroadcasterSet(broadcaster models.Broadcaster) {
	s.broadcaster = &broadcaster
}

// GetBroadcaster returns the active inbound broadcaster, or nil between
// streams.
func (s *Service) GetBroadcaster() *models.Broadcaster {
	return s.broadcaster
}
