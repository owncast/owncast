package rtmp

import (
	"io"
	"net"
	"strconv"
	"time"

	"github.com/nareix/joy5/format/flv"
	"github.com/nareix/joy5/format/flv/flvio"
	"github.com/nareix/joy5/format/rtmp"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/models"
)

// Start binds the RTMP listener and runs the accept loop. Blocks until
// the listener returns. setStreamAsConnected is invoked once per
// authenticated connection; setBroadcaster is invoked when the inbound
// metadata tag arrives.
func (s *Service) Start(setStreamAsConnected func(*io.PipeReader), setBroadcaster func(models.Broadcaster)) {
	s.setStreamAsConnected = setStreamAsConnected
	s.setBroadcaster = setBroadcaster

	bindAddr := s.configRepository.GetRTMPBindAddress()
	port := s.configRepository.GetRTMPPortNumber()
	listenAddress := net.JoinHostPort(bindAddr, strconv.Itoa(port))
	srv := rtmp.NewServer()
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatal(err)
	}

	srv.LogEvent = func(c *rtmp.Conn, nc net.Conn, e int) {
		es := rtmp.EventString[e]
		log.Traceln("RTMP", nc.LocalAddr(), nc.RemoteAddr(), es)
	}

	srv.HandleConn = s.handleConn

	log.Tracef("RTMP server is listening for incoming stream on port: %d", port)

	for {
		nc, err := lis.Accept()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		go srv.HandleNetConn(nc)
	}
}

// handleConn is fired when an inbound RTMP connection takes place.
func (s *Service) handleConn(c *rtmp.Conn, nc net.Conn) {
	c.LogTagEvent = func(isRead bool, t flvio.Tag) {
		if t.Type == flvio.TAG_AMF0 {
			log.Tracef("%+v\n", t.DebugFields())
			s.setCurrentBroadcasterInfo(t, nc.RemoteAddr().String())
		}
	}

	if s.hasInboundConnection {
		log.Errorln("stream already running; can not overtake an existing stream from", nc.RemoteAddr().String())
		_ = nc.Close()
		return
	}

	accessGranted := false
	validStreamingKeys := s.configRepository.GetStreamKeys()

	// If a stream key override was specified then use that instead.
	if s.cfg.TemporaryStreamKey != "" {
		validStreamingKeys = []models.StreamKey{{Key: &s.cfg.TemporaryStreamKey}}
	}

	for _, key := range validStreamingKeys {
		if key.Key != nil && secretMatch(*key.Key, c.URL.Path) {
			accessGranted = true
			break
		}
	}

	if !accessGranted {
		log.Errorln("invalid streaming key; rejecting incoming stream from", nc.RemoteAddr().String())
		_ = nc.Close()
		return
	}

	rtmpOut, rtmpIn := io.Pipe()
	s.pipe = rtmpIn
	log.Infoln("Inbound stream connected from", nc.RemoteAddr().String())
	s.setStreamAsConnected(rtmpOut)

	s.hasInboundConnection = true
	s.rtmpConnection = nc

	w := flv.NewMuxer(rtmpIn)

	for s.hasInboundConnection {
		// If we don't get a readable packet in 10 seconds give up and disconnect.
		if err := s.rtmpConnection.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
			log.Debugln(err)
		}

		pkt, err := c.ReadPacket()

		// Broadcaster disconnected.
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
	if !s.hasInboundConnection {
		return
	}

	log.Infoln("Inbound stream disconnected.")
	_ = conn.Close()
	_ = s.pipe.Close()
	s.hasInboundConnection = false
}

// Disconnect will force disconnect the current inbound RTMP connection.
func (s *Service) Disconnect() {
	if s.rtmpConnection == nil {
		return
	}

	log.Traceln("Inbound stream disconnect requested.")
	s.handleDisconnect(s.rtmpConnection)
}
