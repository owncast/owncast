package stream

import (
	"math"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/models"
)

// setupStats restores stats from the previous session and starts the
// periodic save + viewer-prune goroutines.
func (s *Service) setupStats() error {
	saved := s.getSavedStats()
	s.stats = &saved

	statsSaveTicker := time.NewTicker(1 * time.Minute)
	go func() {
		for range statsSaveTicker.C {
			s.saveStats()
		}
	}()

	viewerCountPruneTicker := time.NewTicker(5 * time.Second)
	go func() {
		for range viewerCountPruneTicker.C {
			s.pruneViewerCount()
		}
	}()

	return nil
}

// IsStreamConnected reports whether the stream is actively delivering
// live HLS content. Returns false during the brief warm-up after an
// RTMP connection where HLS segments aren't ready yet.
func (s *Service) IsStreamConnected() bool {
	if s.stats == nil || !s.stats.StreamConnected {
		return false
	}

	// Kind of a hack. It takes a handful of seconds between an RTMP
	// connection and when HLS data is available. Account for that with
	// an artificial buffer of a few segments.
	timeSinceLastConnected := time.Since(s.stats.LastConnectTime.Time).Seconds()
	waitTime := math.Max(float64(s.configRepository.GetStreamLatencyLevel().SecondsPerSegment)*3.0, 7)
	if timeSinceLastConnected < waitTime {
		return false
	}

	return s.stats.StreamConnected
}

// RemoveChatClient removes a chat client from the active-clients record.
// Currently unused — kept for now during the migration; remove once
// confirmed dead.
func (s *Service) RemoveChatClient(clientID string) {
	log.Trace("Removing the client:", clientID)

	s.statsMu.Lock()
	delete(s.stats.ChatClients, clientID)
	s.statsMu.Unlock()
}

// SetViewerActive marks a viewer as currently watching and updates the
// session/overall peak counters. Silent no-op while no stream is live.
func (s *Service) SetViewerActive(viewer *models.Viewer) {
	if s.stats == nil || !s.stats.StreamConnected {
		return
	}

	s.statsMu.Lock()
	defer s.statsMu.Unlock()

	// Asynchronously, optionally, fetch GeoIP info.
	go func(viewer *models.Viewer) {
		viewer.Geo = s.geoIPClient.GetGeoFromIP(viewer.IPAddress)
	}(viewer)

	if _, exists := s.stats.Viewers[viewer.ClientID]; exists {
		s.stats.Viewers[viewer.ClientID].LastSeen = time.Now()
	} else {
		s.stats.Viewers[viewer.ClientID] = viewer
	}
	s.stats.SessionMaxViewerCount = int(math.Max(float64(len(s.stats.Viewers)), float64(s.stats.SessionMaxViewerCount)))
	s.stats.OverallMaxViewerCount = int(math.Max(float64(s.stats.SessionMaxViewerCount), float64(s.stats.OverallMaxViewerCount)))
}

// GetActiveViewers returns the currently-tracked viewers. The returned
// map is the live one — callers should not modify it.
func (s *Service) GetActiveViewers() map[string]*models.Viewer {
	if s.stats == nil {
		return nil
	}
	return s.stats.Viewers
}

func (s *Service) pruneViewerCount() {
	viewers := make(map[string]*models.Viewer)

	s.statsMu.Lock()
	defer s.statsMu.Unlock()

	for viewerID, viewer := range s.stats.Viewers {
		if time.Since(s.stats.Viewers[viewerID].LastSeen) < activeViewerPurgeTimeout {
			viewers[viewerID] = viewer
		}
	}

	s.stats.Viewers = viewers
}

func (s *Service) saveStats() {
	if err := s.configRepository.SetPeakOverallViewerCount(s.stats.OverallMaxViewerCount); err != nil {
		log.Errorln("error saving viewer count", err)
	}
	if err := s.configRepository.SetPeakSessionViewerCount(s.stats.SessionMaxViewerCount); err != nil {
		log.Errorln("error saving viewer count", err)
	}
	if s.stats.LastDisconnectTime != nil && s.stats.LastDisconnectTime.Valid {
		if err := s.configRepository.SetLastDisconnectTime(s.stats.LastDisconnectTime.Time); err != nil {
			log.Errorln("error saving disconnect time", err)
		}
	}
}

func (s *Service) getSavedStats() models.Stats {
	savedLastDisconnectTime, _ := s.configRepository.GetLastDisconnectTime()

	result := models.Stats{
		ChatClients:           make(map[string]models.Client),
		Viewers:               make(map[string]*models.Viewer),
		SessionMaxViewerCount: s.configRepository.GetPeakSessionViewerCount(),
		OverallMaxViewerCount: s.configRepository.GetPeakOverallViewerCount(),
		LastDisconnectTime:    savedLastDisconnectTime,
	}

	// If the stats were saved > 5min ago then ignore the
	// peak session count value, since the session is over.
	if result.LastDisconnectTime == nil || !result.LastDisconnectTime.Valid || time.Since(result.LastDisconnectTime.Time).Minutes() > 5 {
		result.SessionMaxViewerCount = 0
	}

	return result
}
