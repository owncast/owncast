package main

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// Build-time injected values
var GitCommit string = "unknown"
var BuildVersion string = "0.0.0"
var BuildType string = "localdev"

var storage ChunkStorage
var configuration = getConfig()
var server *Server
var stats *Stats

var usingExternalStorage = false

func main() {
	log.StandardLogger().Printf("Owncast v%s/%s (%s)", BuildVersion, BuildType, GitCommit)

	checkConfig(configuration)
	stats = getSavedStats()
	stats.Setup()

	if configuration.IPFS.Enabled {
		storage = &IPFSStorage{}
		usingExternalStorage = true
	} else if configuration.S3.Enabled {
		storage = &S3Storage{}
		usingExternalStorage = true
	}

	if usingExternalStorage {
		storage.Setup(configuration)
		// hlsDirectoryPath = configuration.PrivateHLSPath
		go monitorVideoContent(configuration.PrivateHLSPath, configuration, storage)
	}

	resetDirectories(configuration)
	go startRTMPService()

	startChatServer()
}

func startChatServer() {
	// log.SetFlags(log.Lshortfile)

	// websocket server
	server = NewServer("/entry")
	go server.Listen()

	// static files
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		http.ServeFile(w, r, path.Join("webroot", r.URL.Path))

		if path.Ext(r.URL.Path) == ".m3u8" {
			clientID := getClientIDFromRequest(r)
			stats.SetClientActive(clientID)
		}
	})

	http.HandleFunc("/status", getStatus)

	log.Printf("Starting public web server on port %d", configuration.WebServerPort)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(configuration.WebServerPort), nil))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	status := Status{
		Online:                stats.IsStreamConnected(),
		ViewerCount:           stats.GetViewerCount(),
		OverallMaxViewerCount: stats.GetOverallMaxViewerCount(),
		SessionMaxViewerCount: stats.GetSessionMaxViewerCount(),
	}
	json.NewEncoder(w).Encode(status)
}

func streamConnected() {
	stats.StreamConnected()

	chunkPath := configuration.PublicHLSPath
	if usingExternalStorage {
		chunkPath = configuration.PrivateHLSPath
	}
	startThumbnailGenerator(chunkPath)
}

func streamDisconnected() {
	stats.StreamDisconnected()
	if configuration.EnableOfflineImage {
		showStreamOfflineState(configuration)
	}
}

func viewerAdded(clientID string) {
	stats.SetClientActive(clientID)
}

func viewerRemoved(clientID string) {
	stats.ViewerDisconnected(clientID)
}
