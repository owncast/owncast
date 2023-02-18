package rtmp

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/nareix/joy5/format/flv"
	"github.com/nareix/joy5/format/flv/flvio"
	log "github.com/sirupsen/logrus"

	"github.com/nareix/joy5/format/rtmp"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
)

func New(d *data.Service) (*Service, error) {
	s := &Service{}

	s.data = d

	return s, nil
}

type Service struct {
	data                      *data.Service
	_hasInboundRTMPConnection bool
	_pipe                     *io.PipeWriter
	_rtmpConnection           net.Conn
	_setStreamAsConnected     func(*io.PipeReader)
	_setBroadcaster           func(models.Broadcaster)
}

// Start starts the rtmp service, listening on specified RTMP port.
func (s *Service) Start(setStreamAsConnected func(*io.PipeReader), setBroadcaster func(models.Broadcaster)) {
	s._setStreamAsConnected = setStreamAsConnected
	s._setBroadcaster = setBroadcaster

	port := s.data.GetRTMPPortNumber()

	rtmpServer := rtmp.NewServer()
	var lis net.Listener
	var err error
	if lis, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal(err)
	}

	rtmpServer.LogEvent = func(c *rtmp.Conn, nc net.Conn, e int) {
		es := rtmp.EventString[e]
		log.Traceln("RTMP", nc.LocalAddr(), nc.RemoteAddr(), es)
	}

	rtmpServer.HandleConn = s.HandleConn

	if err != nil {
		log.Panicln(err)
	}
	log.Tracef("RTMP server is listening for incoming stream on port: %d", port)

	for {
		nc, err := lis.Accept()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		go rtmpServer.HandleNetConn(nc)
	}
}

// HandleConn is fired when an inbound RTMP connection takes place.
func (s *Service) HandleConn(c *rtmp.Conn, nc net.Conn) {
	c.LogTagEvent = func(isRead bool, t flvio.Tag) {
		if t.Type == flvio.TAG_AMF0 {
			log.Tracef("%+v\n", t.DebugFields())
			s.setCurrentBroadcasterInfo(t, nc.RemoteAddr().String())
		}
	}

	if s._hasInboundRTMPConnection {
		log.Errorln("stream already running; can not overtake an existing stream")
		_ = nc.Close()
		return
	}

	accessGranted := false
	validStreamingKeys := s.data.GetStreamKeys()

	for _, key := range validStreamingKeys {
		if secretMatch(key.Key, c.URL.Path) {
			accessGranted = true
			break
		}
	}

	// Test against the temporary key if it was set at runtime.
	if config.TemporaryStreamKey != "" && secretMatch(config.TemporaryStreamKey, c.URL.Path) {
		accessGranted = true
	}

	if !accessGranted {
		log.Errorln("invalid streaming key; rejecting incoming stream")
		_ = nc.Close()
		return
	}

	rtmpOut, rtmpIn := io.Pipe()
	s._pipe = rtmpIn
	log.Infoln("Inbound stream connected.")
	s._setStreamAsConnected(rtmpOut)

	s._hasInboundRTMPConnection = true
	s._rtmpConnection = nc

	w := flv.NewMuxer(rtmpIn)

	for {
		if !s._hasInboundRTMPConnection {
			break
		}

		// If we don't get a readable packet in 10 seconds give up and disconnect
		if err := s._rtmpConnection.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
			log.Debugln(err)
		}

		pkt, err := c.ReadPacket()

		// Broadcaster disconnected
		if err == io.EOF {
			s.handleDisconnect(nc)
			return
		}

		// Read timeout.  Disconnect.
		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			log.Debugln("Timeout reading the inbound stream from the broadcaster.  Assuming that they disconnected and ending the stream.")
			s.handleDisconnect(nc)
			return
		}

		if err := w.WritePacket(pkt); err != nil {
			log.Errorln("unable to write rtmp packet", err)
			s.handleDisconnect(nc)
			return
		}
	}
}

func (s *Service) handleDisconnect(conn net.Conn) {
	if !s._hasInboundRTMPConnection {
		return
	}

	log.Infoln("Inbound stream disconnected.")
	_ = conn.Close()
	_ = s._pipe.Close()
	s._hasInboundRTMPConnection = false
}

// Disconnect will force disconnect the current inbound RTMP connection.
func (s *Service) Disconnect() {
	if s._rtmpConnection == nil {
		return
	}

	log.Traceln("Inbound stream disconnect requested.")
	s.handleDisconnect(s._rtmpConnection)
}
