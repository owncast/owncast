package core

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/models"
	"github.com/gabek/owncast/termui"
	"github.com/gabek/owncast/utils"
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
	for clientID, lastConnectedtime := range _stats.Clients {
		timeSinceLastActive := time.Since(lastConnectedtime).Minutes()
		if timeSinceLastActive > 2 {
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
func SetClientActive(clientID string) {
	// if _, ok := s.clients[clientID]; !ok {
	// 	fmt.Println("Marking client active:", clientID, s.GetViewerCount()+1, "clients connected.")
	// }

	l.Lock()
	_stats.Clients[clientID] = time.Now()
	l.Unlock()
	_stats.SessionMaxViewerCount = int(math.Max(float64(len(_stats.Clients)), float64(_stats.SessionMaxViewerCount)))
	_stats.OverallMaxViewerCount = int(math.Max(float64(_stats.SessionMaxViewerCount), float64(_stats.OverallMaxViewerCount)))

	termui.SetClientList(_stats.Clients)
}

//RemoveClient removes a client from the active clients record
func RemoveClient(clientID string) {
	log.Trace("Removing the client:", clientID)

	delete(_stats.Clients, clientID)
	termui.SetClientList(_stats.Clients)
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
		Clients: make(map[string]time.Time),
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
