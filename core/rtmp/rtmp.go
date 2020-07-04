package rtmp

import (
	"os"
	"strings"
	"syscall"

	"github.com/Seize/joy4/av/avutil"
	"github.com/Seize/joy4/format/ts"
	log "github.com/sirupsen/logrus"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/core"
	"github.com/gabek/owncast/core/ffmpeg"
	"github.com/gabek/owncast/utils"

	"github.com/Seize/joy4/format"
	"github.com/Seize/joy4/format/rtmp"
)

var (
	//IsConnected whether there is a connection or not
	_isConnected = false
)

func init() {
	format.RegisterAll()
}

//Start starts the rtmp service, listening on port 1935
func Start() {

	port := 1935
	server := &rtmp.Server{}

	server.HandlePublish = handlePublish

	error := server.ListenAndServe()
	if error != nil {
		log.Panicln(error)
	}
	log.Printf("RTMP server is listening for incoming stream on port: %d", port)
}

func handlePublish(conn *rtmp.Conn) {
	// Commented out temporarily because I have no way to set _isConnected to false after RTMP is closed.
	// if _isConnected {
	// 	log.Errorln("stream already running; can not overtake an existing stream")
	// 	conn.Close()
	// 	return
	// }

	streamingKeyComponents := strings.Split(conn.URL.Path, "/")
	streamingKey := streamingKeyComponents[len(streamingKeyComponents)-1]
	if streamingKey != config.Config.VideoSettings.StreamingKey {
		log.Errorln("invalid streaming key; rejecting incoming stream")
		conn.Close()
		return
	}

	pipePath := utils.GetTemporaryPipePath()
	syscall.Mkfifo(pipePath, 0666)
	transcoder := ffmpeg.NewTranscoder()
	go transcoder.Start()

	_isConnected = true
	core.SetStreamAsConnected()

	f, err := os.OpenFile(pipePath, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		panic(err)
	}

	muxer := ts.NewMuxer(f)
	avutil.CopyFile(muxer, conn)
}

//IsConnected gets whether there is an rtmp connection or not
//this is only a getter since it is controlled by the rtmp handler
func IsConnected() bool {
	return _isConnected
}
