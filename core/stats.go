package core

import (
	"math"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/geoip"
	"github.com/owncast/owncast/models"
)

const (
	activeViewerPurgeTimeout = time.Second * 15
)

func (s *Service) setupStats() error {
	savedLastDisconnectTime, _ := s.Data.GetLastDisconnectTime()

	s.geoIpClient = geoip.NewClient()

	s.stats = &models.Stats{
		ChatClients:           make(map[string]models.Client),
		Viewers:               make(map[string]*models.Viewer),
		SessionMaxViewerCount: s.Data.GetPeakSessionViewerCount(),
		OverallMaxViewerCount: s.Data.GetPeakOverallViewerCount(),
		LastDisconnectTime:    savedLastDisconnectTime,
	}

	// If the stats were saved > 5min ago then ignore the
	// peak session count value, since the session is over.
	if s.stats.LastDisconnectTime == nil || !s.stats.LastDisconnectTime.Valid || time.Since(s.stats.LastDisconnectTime.Time).Minutes() > 5 {
		s.stats.SessionMaxViewerCount = 0
	}

	statsSaveTimer := time.NewTicker(1 * time.Minute)
	go func() {
		for range statsSaveTimer.C {
			s.saveStats()
		}
	}()

	viewerCountPruneTimer := time.NewTicker(5 * time.Second)
	go func() {
		for range viewerCountPruneTimer.C {
			s.pruneViewerCount()
		}
	}()

	return nil
}

// IsStreamConnected checks if the stream is connected or not.
func (s *Service) IsStreamConnected() bool {
	if !s.stats.StreamConnected {
		return false
	}

	// Kind of a hack.  It takes a handful of seconds between a RTMP connection and when HLS Data is available.
	// So account for that with an artificial buffer of four segments.
	timeSinceLastConnected := time.Since(s.stats.LastConnectTime.Time).Seconds()
	waitTime := math.Max(float64(s.Data.GetStreamLatencyLevel().SecondsPerSegment)*3.0, 7)
	if timeSinceLastConnected < waitTime {
		return false
	}

	return s.stats.StreamConnected
}

// RemoveChatClient removes a client from the active clients record.
func (s *Service) RemoveChatClient(clientID string) {
	log.Trace("Removing the client:", clientID)

	s.mux.Lock()
	delete(s.stats.ChatClients, clientID)
	s.mux.Unlock()
}

// SetViewerActive sets a client as active and connected.
func (s *Service) SetViewerActive(viewer *models.Viewer) {
	// Don't update viewer counts if a live stream session is not active.
	if !s.stats.StreamConnected {
		return
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	// Asynchronously, optionally, fetch GeoIP s.Data.
	go func(viewer *models.Viewer) {
		viewer.Geo = s.geoIpClient.GetGeoFromIP(viewer.IPAddress)
	}(viewer)

	if _, exists := s.stats.Viewers[viewer.ClientID]; exists {
		s.stats.Viewers[viewer.ClientID].LastSeen = time.Now()
	} else {
		s.stats.Viewers[viewer.ClientID] = viewer
	}
	s.stats.SessionMaxViewerCount = int(math.Max(float64(len(s.stats.Viewers)), float64(s.stats.SessionMaxViewerCount)))
	s.stats.OverallMaxViewerCount = int(math.Max(float64(s.stats.SessionMaxViewerCount), float64(s.stats.OverallMaxViewerCount)))
}

// GetActiveViewers will return the active viewers.
func (s *Service) GetActiveViewers() map[string]*models.Viewer {
	return s.stats.Viewers
}

func (s *Service) pruneViewerCount() {
	viewers := make(map[string]*models.Viewer)

	s.mux.Lock()
	defer s.mux.Unlock()

	for viewerID, viewer := range s.stats.Viewers {
		viewerLastSeenTime := s.stats.Viewers[viewerID].LastSeen
		if time.Since(viewerLastSeenTime) < activeViewerPurgeTimeout {
			viewers[viewerID] = viewer
		}
	}

	s.stats.Viewers = viewers
}

func (s *Service) saveStats() {
	if err := s.Data.SetPeakOverallViewerCount(s.stats.OverallMaxViewerCount); err != nil {
		log.Errorln("error saving viewer count", err)
	}
	if err := s.Data.SetPeakSessionViewerCount(s.stats.SessionMaxViewerCount); err != nil {
		log.Errorln("error saving viewer count", err)
	}
	if s.stats.LastDisconnectTime != nil && s.stats.LastDisconnectTime.Valid {
		if err := s.Data.SetLastDisconnectTime(s.stats.LastDisconnectTime.Time); err != nil {
			log.Errorln("error saving disconnect time", err)
		}
	}
}
