/*
Viewer counting doesn't just count the number of websocket clients that are currently connected,
because people may be watching the stream outside of the web browser via any HLS video client.
Instead we keep track of requests and consider each unique IP as a "viewer".
As a signal, however, we do use the websocket disconnect from a client as a signal that a viewer
dropped and we call ViewerDisconnected().
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"
)

type Stats struct {
	streamConnected       bool      `json:"-"`
	SessionMaxViewerCount int       `json:"sessionMaxViewerCount"`
	OverallMaxViewerCount int       `json:"overallMaxViewerCount"`
	LastDisconnectTime    time.Time `json:"lastDisconnectTime"`

	clients map[string]time.Time
}

func (s *Stats) Setup() {
	s.clients = make(map[string]time.Time)

	statsSaveTimer := time.NewTicker(2 * time.Minute)
	go func() {
		for {
			select {
			case <-statsSaveTimer.C:
				s.save()
			}
		}
	}()

	staleViewerPurgeTimer := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-staleViewerPurgeTimer.C:
				s.purgeStaleViewers()
			}
		}
	}()
}

func (s *Stats) purgeStaleViewers() {
	for clientID, lastConnectedtime := range s.clients {
		timeSinceLastActive := time.Since(lastConnectedtime).Minutes()
		if timeSinceLastActive > 2 {
			s.ViewerDisconnected(clientID)
		}

	}
}

func (s *Stats) IsStreamConnected() bool {
	return s.streamConnected
}

func (s *Stats) GetViewerCount() int {
	return len(s.clients)
}

func (s *Stats) GetSessionMaxViewerCount() int {
	return s.SessionMaxViewerCount
}

func (s *Stats) GetOverallMaxViewerCount() int {
	return s.OverallMaxViewerCount
}

func (s *Stats) SetClientActive(clientID string) {
	fmt.Println("Marking client active:", clientID)

	s.clients[clientID] = time.Now()
	s.SessionMaxViewerCount = int(math.Max(float64(s.GetViewerCount()), float64(s.SessionMaxViewerCount)))
	s.OverallMaxViewerCount = int(math.Max(float64(s.SessionMaxViewerCount), float64(s.OverallMaxViewerCount)))

	fmt.Println("Now", s.GetViewerCount(), "clients connected.")
}

func (s *Stats) ViewerDisconnected(clientID string) {
	log.Println("Removed client", clientID)

	delete(s.clients, clientID)
}

func (s *Stats) StreamConnected() {
	s.streamConnected = true

	timeSinceDisconnect := time.Since(s.LastDisconnectTime).Minutes()
	if timeSinceDisconnect > 15 {
		s.SessionMaxViewerCount = 0
	}
}

func (s *Stats) StreamDisconnected() {
	s.streamConnected = false
	s.LastDisconnectTime = time.Now()
}

func (s *Stats) save() {
	jsonData, err := json.Marshal(&s)
	verifyError(err)

	f, err := os.Create("config/stats.json")
	defer f.Close()

	verifyError(err)

	_, err = f.Write(jsonData)
	verifyError(err)
}

func getSavedStats() *Stats {
	filePath := "config/stats.json"

	if !fileExists(filePath) {
		return &Stats{}
	}

	jsonFile, err := ioutil.ReadFile(filePath)

	var stats Stats
	err = json.Unmarshal(jsonFile, &stats)
	if err != nil {
		panic(err)
	}

	return &stats
}
