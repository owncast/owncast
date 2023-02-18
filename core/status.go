package core

import (
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
)

// GetStatus gets the status of the system.
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
		StreamTitle:           s.Data.GetStreamTitle(),
	}
}

// GetCurrentBroadcast will return the currently active broadcast.
func (s *Service) GetCurrentBroadcast() *models.CurrentBroadcast {
	return s.streamState._currentBroadcast
}

// SetBroadcaster will store the current inbound broadcasting details.
func (s *Service) SetBroadcaster(b models.Broadcaster) {
	s.broadcaster = &b
}

// GetBroadcaster will return the details of the currently active broadcaster.
func (s *Service) GetBroadcaster() *models.Broadcaster {
	return s.broadcaster
}
