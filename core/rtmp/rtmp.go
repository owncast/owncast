package rtmp

import (
	"fmt"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
	yutmp "github.com/yutopp/go-rtmp"
)

var (
	//IsConnected whether there is a connection or not
	_isConnected = false
)

//Start starts the rtmp service, listening on port 1935
func Start() {
	port := 1935

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Panicf("Failed to resolve the tcp address for the rtmp service: %+v", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Panicf("Failed to acquire the tcp listener: %+v", err)
	}

	srv := yutmp.NewServer(&yutmp.ServerConfig{
		OnConnect: func(conn net.Conn) (io.ReadWriteCloser, *yutmp.ConnConfig) {
			l := log.StandardLogger()
			l.SetLevel(log.WarnLevel)

			return conn, &yutmp.ConnConfig{
				Handler: &Handler{},

				ControlState: yutmp.StreamControlStateConfig{
					DefaultBandwidthWindowSize: 6 * 1024 * 1024 / 8,
				},

				Logger: l,
			}
		},
	})

	log.Printf("RTMP server is listening for incoming stream on port: %d", port)
	if err := srv.Serve(listener); err != nil {
		log.Panicf("Failed to serve the rtmp service: %+v", err)
	}
}

//IsConnected gets whether there is an rtmp connection or not
//this is only a getter since it is controlled by the rtmp handler
func IsConnected() bool {
	return _isConnected
}
