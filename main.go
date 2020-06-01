package main

import (
	"io"
	"net"
	"net/http"
	"strconv"

	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/yutopp/go-rtmp"
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

func startRTMPService() {
	port := 1935
	log.Printf("RTMP server is listening for incoming stream on port %d.\n", port)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Panicf("Failed: %+v", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Panicf("Failed: %+v", err)
	}

	srv := rtmp.NewServer(&rtmp.ServerConfig{
		OnConnect: func(conn net.Conn) (io.ReadWriteCloser, *rtmp.ConnConfig) {
			l := log.StandardLogger()
			l.SetLevel(logrus.WarnLevel)

			h := &Handler{}

			return conn, &rtmp.ConnConfig{
				Handler: h,

				ControlState: rtmp.StreamControlStateConfig{
					DefaultBandwidthWindowSize: 6 * 1024 * 1024 / 8,
				},

				Logger: l,
			}
		},
	})
	if err := srv.Serve(listener); err != nil {
		log.Panicf("Failed: %+v", err)
	}
}
