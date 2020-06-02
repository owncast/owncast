package main

import (
	"io"
	"net"
	"strconv"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/yutopp/go-rtmp"
)

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
