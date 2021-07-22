package core

import (
	"math"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/geoip"
	"github.com/owncast/owncast/models"
)

var l = &sync.RWMutex{}
var _activeViewerPurgeTimeout = time.Second * 10

func setupStats() error {
	s := getSavedStats()
	_stats = &s

	statsSaveTimer := time.NewTicker(1 * time.Minute)
	go func() {
		for range statsSaveTimer.C {
			saveStats()
		}
	}()

	viewerCountPruneTimer := time.NewTicker(5 * time.Second)
	go func() {
		for range viewerCountPruneTimer.C {
			pruneViewerCount()
		}
	}()

	return nil
}

// IsStreamConnected checks if the stream is connected or not.
func IsStreamConnected() bool {
	if !_stats.StreamConnected {
		return false
	}

	// Kind of a hack.  It takes a handful of seconds between a RTMP connection and when HLS data is available.
	// So account for that with an artificial buffer of four segments.
	timeSinceLastConnected := time.Since(_stats.LastConnectTime.Time).Seconds()
	waitTime := math.Max(float64(data.GetStreamLatencyLevel().SecondsPerSegment)*3.0, 7)
	if timeSinceLastConnected < waitTime {
		return false
	}

	return _stats.StreamConnected
}

// SetChatClientActive sets a client as active and connected.
func SetChatClientActive(client models.Client) {
	l.Lock()
	defer l.Unlock()

	// If this clientID already exists then update it.
	// Otherwise set a new one.
	if existingClient, ok := _stats.ChatClients[client.ClientID]; ok {
		existingClient.LastSeen = time.Now()
		existingClient.Username = client.Username
		existingClient.MessageCount = client.MessageCount
		existingClient.Geo = geoip.GetGeoFromIP(existingClient.IPAddress)
		_stats.ChatClients[client.ClientID] = existingClient
	} else {
		if client.Geo == nil {
			geoip.FetchGeoForIP(client.IPAddress)
		}
		_stats.ChatClients[client.ClientID] = client
	}
}

// RemoveChatClient removes a client from the active clients record.
func RemoveChatClient(clientID string) {
	log.Trace("Removing the client:", clientID)

	l.Lock()
	delete(_stats.ChatClients, clientID)
	l.Unlock()
}

// SetViewerIdActive sets a client as active and connected.
func SetViewerIdActive(id string) {
	l.Lock()
	defer l.Unlock()

	_stats.Viewers[id] = time.Now()

	// Don't update viewer counts if a live stream session is not active.
	if _stats.StreamConnected {
		_stats.SessionMaxViewerCount = int(math.Max(float64(len(_stats.Viewers)), float64(_stats.SessionMaxViewerCount)))
		_stats.OverallMaxViewerCount = int(math.Max(float64(_stats.SessionMaxViewerCount), float64(_stats.OverallMaxViewerCount)))
	}
}

func pruneViewerCount() {
	viewers := make(map[string]time.Time)

	l.Lock()
	defer l.Unlock()
	
	for viewerId := range _stats.Viewers {
		viewerLastSeenTime := _stats.Viewers[viewerId]
		if time.Since(viewerLastSeenTime) < _activeViewerPurgeTimeout {
			viewers[viewerId] = viewerLastSeenTime
		}
	}

	_stats.Viewers = viewers
}

func saveStats() {
	if err := data.SetPeakOverallViewerCount(_stats.OverallMaxViewerCount); err != nil {
		log.Errorln("error saving viewer count", err)
	}
	if err := data.SetPeakSessionViewerCount(_stats.SessionMaxViewerCount); err != nil {
		log.Errorln("error saving viewer count", err)
	}
	if _stats.LastDisconnectTime.Valid {
		if err := data.SetLastDisconnectTime(_stats.LastConnectTime.Time); err != nil {
			log.Errorln("error saving disconnect time", err)
		}
	}
}

func getSavedStats() models.Stats {
	savedLastDisconnectTime, _ := data.GetLastDisconnectTime()

	result := models.Stats{
		ChatClients:           make(map[string]models.Client),
		Viewers:               make(map[string]time.Time),
		SessionMaxViewerCount: data.GetPeakSessionViewerCount(),
		OverallMaxViewerCount: data.GetPeakOverallViewerCount(),
		LastDisconnectTime:    savedLastDisconnectTime,
	}

	// If the stats were saved > 5min ago then ignore the
	// peak session count value, since the session is over.
	if !result.LastDisconnectTime.Valid || time.Since(result.LastDisconnectTime.Time).Minutes() > 5 {
		result.SessionMaxViewerCount = 0
	}

	return result
}
