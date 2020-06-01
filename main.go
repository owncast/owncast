package main

import (
	"io"
	"net"
	"net/http"

	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/yutopp/go-rtmp"
)

var ipfs icore.CoreAPI

func main() {
	resetDirectories()

	ipfsInstance, node, _ := createIPFSInstance()
	ipfs = *ipfsInstance

	createIPFSDirectory(ipfsInstance, "./hls")
	// touch("hls/stream.m3u8")

	go startIPFSNode(ipfs, node)
	go monitorVideoContent("./hls/", ipfsInstance)
	go startChatServer()

	startRTMPService()
}

func startChatServer() {
	// log.SetFlags(log.Lshortfile)

	// websocket server
	server := NewServer("/entry")
	go server.Listen()

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func startRTMPService() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1935")
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
