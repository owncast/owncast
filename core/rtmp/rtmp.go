package rtmp

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/nareix/joy5/format/flv"
	"github.com/nareix/joy5/format/flv/flvio"
	log "github.com/sirupsen/logrus"

	"github.com/nareix/joy5/format/rtmp"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
)

var (
	//IsConnected whether there is a connection or not
	_isConnected = false
)

var _pipe *os.File
var _rtmpConnection net.Conn

//Start starts the rtmp service, listening on port 1935
func Start() {
	port := 1935
	s := rtmp.NewServer()
	var lis net.Listener
	var err error
	if lis, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
		return
	}

	s.LogEvent = func(c *rtmp.Conn, nc net.Conn, e int) {
		es := rtmp.EventString[e]
		log.Traceln(unsafe.Pointer(c), nc.LocalAddr(), nc.RemoteAddr(), es)
	}

	s.HandleConn = HandleConn

	if err != nil {
		log.Panicln(err)
	}
	log.Infof("RTMP server is listening for incoming stream on port: %d", port)

	for {
		nc, err := lis.Accept()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		go s.HandleNetConn(nc)
	}
}

func setCurrentBroadcasterInfo(t flvio.Tag, remoteAddr string) {
	data, err := getInboundDetailsFromMetadata(t.DebugFields())
	if err != nil {
		log.Errorln(err)
		return
	}

	broadcaster := models.Broadcaster{
		RemoteAddr: remoteAddr,
		Time:       time.Now(),
		StreamDetails: models.InboundStreamDetails{
			Width:          data.Width,
			Height:         data.Height,
			VideoBitrate:   int(data.VideoBitrate),
			VideoCodec:     getVideoCodec(data.VideoCodec),
			VideoFramerate: data.VideoFramerate,
			AudioBitrate:   int(data.AudioBitrate),
			AudioCodec:     getAudioCodec(data.AudioCodec),
			Encoder:        data.Encoder,
		},
	}

	core.SetBroadcaster(broadcaster)
}

func HandleConn(c *rtmp.Conn, nc net.Conn) {
	c.LogTagEvent = func(isRead bool, t flvio.Tag) {
		if t.Type == flvio.TAG_AMF0 {
			log.Tracef("%+v\n", t.DebugFields())
			setCurrentBroadcasterInfo(t, nc.RemoteAddr().String())
		}
	}

	if _isConnected {
		log.Errorln("stream already running; can not overtake an existing stream")
		nc.Close()
		return
	}

	streamingKeyComponents := strings.Split(c.URL.Path, "/")
	streamingKey := streamingKeyComponents[len(streamingKeyComponents)-1]
	if streamingKey != config.Config.VideoSettings.StreamingKey {
		log.Errorln("invalid streaming key; rejecting incoming stream")
		nc.Close()
		return
	}

	log.Infoln("Incoming RTMP connected.")

	pipePath := utils.GetTemporaryPipePath()
	syscall.Mkfifo(pipePath, 0666)

	_isConnected = true
	core.SetStreamAsConnected()
	_rtmpConnection = nc

	f, err := os.OpenFile(pipePath, os.O_RDWR, os.ModeNamedPipe)
	_pipe = f
	if err != nil {
		panic(err)
	}

	w := flv.NewMuxer(f)

	for {
		if !_isConnected {
			break
		}

		pkt, err := c.ReadPacket()
		if err == io.EOF {
			handleDisconnect(nc)
			break
		}

		if err := w.WritePacket(pkt); err != nil {
			panic(err)
		}
	}

}

func handleDisconnect(conn net.Conn) {
	log.Infoln("RTMP disconnected.")
	conn.Close()
	_pipe.Close()
	_isConnected = false
}

// Disconnect will force disconnect the current inbound RTMP connection.
func Disconnect() {
	if _rtmpConnection == nil {
		return
	}

	log.Infoln("Inbound stream disconnect requested.")
	handleDisconnect(_rtmpConnection)
}

//IsConnected gets whether there is an rtmp connection or not
//this is only a getter since it is controlled by the rtmp handler
func IsConnected() bool {
	return _isConnected
}
