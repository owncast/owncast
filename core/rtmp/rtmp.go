package rtmp

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/core"
	"github.com/gabek/owncast/core/ffmpeg"
	"github.com/gabek/owncast/utils"

	"github.com/nareix/joy5/format/flv"
	"github.com/nareix/joy5/format/rtmp"

	log "github.com/sirupsen/logrus"
)

var (
	//IsConnected whether there is a connection or not
	_isConnected = false
	pipePath     = utils.GetTemporaryPipePath()
	filePipe     *os.File
)

//Start starts the rtmp service, listening on port 1935
func Start() {
	port := 1935

	server := rtmp.NewServer()

	server.LogEvent = func(conn *rtmp.Conn, nc net.Conn, e int) {
		log.Errorln("RTMP status:", rtmp.EventString[e])
	}

	server.OnNewConn = func(conn *rtmp.Conn) {
		log.Println("OnNewConn!", conn.FlashVer)
	}

	server.HandleConn = func(conn *rtmp.Conn, nc net.Conn) {
		if _isConnected {
			log.Errorln("stream already running; can not overtake an existing stream")
			nc.Close()
		}

		streamingKeyComponents := strings.Split(conn.URL.Path, "/")
		streamingKey := streamingKeyComponents[len(streamingKeyComponents)-1]

		if streamingKey != config.Config.VideoSettings.StreamingKey {
			log.Errorln("invalid streaming key; rejecting incoming stream")
			nc.Close()
			return
		}

		// Record streams as FLV
		syscall.Mkfifo(pipePath, 0666)
		file, err := os.OpenFile(pipePath, os.O_RDWR, os.ModeNamedPipe)
		if err != nil {
			log.Panicln(err)
		}

		filePipe = file
		fmt.Println(pipePath)

		muxer := flv.NewMuxer(filePipe)

		_isConnected = true
		core.SetStreamAsConnected()

		defer filePipe.Close()

		if err != nil {
			panic(err)
		}

		defer nc.Close()

		transcoder := ffmpeg.NewTranscoder()
		go transcoder.Start()

		for {
			pkt, err := conn.ReadPacket()
			if err != nil {
				log.Errorln(err)

				if err == io.EOF {
					_isConnected = false
					core.SetStreamAsDisconnected()
					break
				}

				return
			}

			// err = muxer.WriteFileHeader()
			// if err != nil {
			// 	log.Errorln(err)
			// }

			err = muxer.WritePacket(pkt)
			if err != nil {
				log.Errorln(err)
			}

		}
	}

	var lis net.Listener
	var err error
	if lis, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
		return
	}
	log.Printf("RTMP server is listening for incoming stream on port: %d", port)

	go func() {
		for {
			nc, err := lis.Accept()
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			go server.HandleNetConn(nc)
		}
	}()

}

//IsConnected gets whether there is an rtmp connection or not
//this is only a getter since it is controlled by the rtmp handler
func IsConnected() bool {
	return _isConnected
}
