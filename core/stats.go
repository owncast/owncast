package core

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/geoip"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
)

const (
	statsFilePath = "stats.json"
)

var l = sync.Mutex{}

func setupStats() error {
	s, err := getSavedStats()
	if err != nil {
		return err
	}

	_stats = &s

	statsSaveTimer := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-statsSaveTimer.C:
				if err := saveStatsToFile(); err != nil {
					panic(err)
				}
			}
		}
	}()

	staleViewerPurgeTimer := time.NewTicker(3 * time.Second)
	go func() {
		for {
			select {
			case <-staleViewerPurgeTimer.C:
				purgeStaleViewers()
			}
		}
	}()

	return nil
}

func purgeStaleViewers() {
	for clientID, client := range _stats.Clients {
		if client.LastSeen.IsZero() {
			continue
		}

		timeSinceLastActive := time.Since(client.LastSeen).Minutes()
		if timeSinceLastActive > 1 {
			RemoveClient(clientID)
		}
	}
}

//IsStreamConnected checks if the stream is connected or not
func IsStreamConnected() bool {
	if !_stats.StreamConnected {
		return false
	}

	// Kind of a hack.  It takes a handful of seconds between a RTMP connection and when HLS data is available.
	// So account for that with an artificial buffer of four segments.
	timeSinceLastConnected := time.Since(_stats.LastConnectTime.Time).Seconds()
	if timeSinceLastConnected < float64(config.Config.GetVideoSegmentSecondsLength()*4.0) {
		return false
	}

	return _stats.StreamConnected
}

//SetClientActive sets a client as active and connected
func SetClientActive(client models.Client) {
	l.Lock()
	// If this clientID already exists then update it.
	// Otherwise set a new one.
	if existingClient, ok := _stats.Clients[client.ClientID]; ok {
		existingClient.LastSeen = time.Now()
		existingClient.Username = client.Username
		existingClient.MessageCount = client.MessageCount
		_stats.Clients[client.ClientID] = existingClient
	} else {
		if client.Geo == nil {
			geo := geoip.GetGeoFromIP(client.IPAddress)
			client.Geo = geo
		}
		_stats.Clients[client.ClientID] = client
	}
	l.Unlock()

	// Don't update viewer counts if a live stream session is not active.
	if _stats.StreamConnected {
		_stats.SessionMaxViewerCount = int(math.Max(float64(len(_stats.Clients)), float64(_stats.SessionMaxViewerCount)))
		_stats.OverallMaxViewerCount = int(math.Max(float64(_stats.SessionMaxViewerCount), float64(_stats.OverallMaxViewerCount)))
	}
}

//RemoveClient removes a client from the active clients record
func RemoveClient(clientID string) {
	log.Trace("Removing the client:", clientID)

	delete(_stats.Clients, clientID)
}

func GetClients() []models.Client {
	clients := make([]models.Client, 0)
	for _, client := range _stats.Clients {
		chatClient := chat.GetClient(client.ClientID)
		if chatClient != nil {
			clients = append(clients, chatClient.GetViewerClientFromChatClient())
		} else {
			clients = append(clients, client)
		}
	}
	return clients
}

func saveStatsToFile() error {
	jsonData, err := json.Marshal(_stats)
	if err != nil {
		return err
	}

	f, err := os.Create(statsFilePath)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := f.Write(jsonData); err != nil {
		return err
	}

	return nil
}

func getSavedStats() (models.Stats, error) {
	result := models.Stats{
		Clients: make(map[string]models.Client),
	}

	if !utils.DoesFileExists(statsFilePath) {
		return result, nil
	}

	jsonFile, err := ioutil.ReadFile(statsFilePath)
	if err != nil {
		return result, nil
	}

	if err := json.Unmarshal(jsonFile, &result); err != nil {
		return result, nil
	}

	return result, nil
}
