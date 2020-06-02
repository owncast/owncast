package main

import (
	"net/http"
	"strconv"

	icore "github.com/ipfs/interface-go-ipfs-core"
	log "github.com/sirupsen/logrus"
)

var ipfs icore.CoreAPI
var configuration = getConfig()

func main() {
	checkConfig(configuration)
	// resetDirectories()

	var hlsDirectoryPath = configuration.PublicHLSPath

	log.Println("Starting up.  Please wait...")

	if configuration.IPFS.Enabled {
		hlsDirectoryPath = configuration.PrivateHLSPath
		enableIPFS()
		go monitorVideoContent(hlsDirectoryPath, configuration, &ipfs)
	}

	go startChatServer()

	startRTMPService()
}

func enableIPFS() {
	log.Println("Enabling IPFS support...")

	ipfsInstance, node, _ := createIPFSInstance()
	ipfs = *ipfsInstance

	createIPFSDirectory(ipfsInstance, "./hls")
	go startIPFSNode(ipfs, node)
}

func startChatServer() {
	// log.SetFlags(log.Lshortfile)

	// websocket server
	server := NewServer("/entry")
	go server.Listen()

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(configuration.WebServerPort), nil))
}
