package rtmp

import (
	"io"
	"net"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/format/ts"
	log "github.com/sirupsen/logrus"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/core"
	"github.com/gabek/owncast/core/ffmpeg"
	"github.com/gabek/owncast/utils"

	"github.com/nareix/joy4/format"
	"github.com/nareix/joy4/format/rtmp"
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
	if _isConnected {
		log.Errorln("stream already running; can not overtake an existing stream")
		conn.Close()
		return
	}

	streamingKeyComponents := strings.Split(conn.URL.Path, "/")
	streamingKey := streamingKeyComponents[len(streamingKeyComponents)-1]
	if streamingKey != config.Config.VideoSettings.StreamingKey {
		log.Errorln("invalid streaming key; rejecting incoming stream")
		conn.Close()
		return
	}

	log.Println("Incoming RTMP connected.")

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

	// Is this too fast?  Are there downsides to peeking
	// into the stream so frequently?
	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				error := connCheck(conn.NetConn())
				if error == io.EOF {
					handleDisconnect(conn)
				}
			}
		}
	}()
	muxer := ts.NewMuxer(f)
	avutil.CopyFile(muxer, conn)
}

// Proactively check if the RTMP connection is still active or not.
// Taken from https://stackoverflow.com/a/58664631.
func connCheck(conn net.Conn) error {
	var sysErr error = nil
	rc, err := conn.(syscall.Conn).SyscallConn()
	if err != nil {
		return err
	}
	err = rc.Read(func(fd uintptr) bool {
		var buf []byte = []byte{0}
		n, _, err := syscall.Recvfrom(int(fd), buf, syscall.MSG_PEEK|syscall.MSG_DONTWAIT)
		switch {
		case n == 0 && err == nil:
			sysErr = io.EOF
		case err == syscall.EAGAIN || err == syscall.EWOULDBLOCK:
			sysErr = nil
		default:
			sysErr = err
		}
		return true
	})
	if err != nil {
		return err
	}

	return sysErr
}

func handleDisconnect(conn *rtmp.Conn) {
	log.Println("RTMP disconnected.")
	conn.Close()
	_isConnected = false
	core.SetStreamAsDisconnected()
}

//IsConnected gets whether there is an rtmp connection or not
//this is only a getter since it is controlled by the rtmp handler
func IsConnected() bool {
	return _isConnected
}
