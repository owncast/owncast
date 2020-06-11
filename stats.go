package main

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
	"time"
)

type Stats struct {
	streamConnected       bool      `json:"-"`
	ViewerCount           int       `json:"viewerCount"`
	SessionMaxViewerCount int       `json:"sessionMaxViewerCount"`
	OverallMaxViewerCount int       `json:"overallMaxViewerCount"`
	LastDisconnectTime    time.Time `json:"lastDisconnectTime"`
}

func (s *Stats) Setup() {
	ticker := time.NewTicker(2 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				s.save()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *Stats) IsStreamConnected() bool {
	return s.streamConnected
}

func (s *Stats) SetViewerCount(count int) {
	s.ViewerCount = count
	s.SessionMaxViewerCount = int(math.Max(float64(s.ViewerCount), float64(s.SessionMaxViewerCount)))
	s.OverallMaxViewerCount = int(math.Max(float64(s.SessionMaxViewerCount), float64(s.OverallMaxViewerCount)))
}

func (s *Stats) GetViewerCount() int {
	return s.ViewerCount
}

func (s *Stats) GetSessionMaxViewerCount() int {
	return s.SessionMaxViewerCount
}

func (s *Stats) GetOverallMaxViewerCount() int {
	return s.OverallMaxViewerCount
}

func (s *Stats) ViewerConnected() {
}

func (s *Stats) ViewerDisconnected() {
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
