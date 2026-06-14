// Package rtmp is the inbound RTMP ingest service. It listens on the
// configured RTMP port, authenticates incoming streams against the
// configured stream keys, and pipes the FLV-muxed bytes to the stream
// service via callbacks supplied to Start.
//
// The stream service holds a *Service and calls Start(...) once at
// startup; admin handlers hold one to force-disconnect via Disconnect().
package rtmp

import (
	"io"
	"net"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
)

// Service owns the RTMP listener and the single in-flight inbound
// connection (RTMP allows one broadcaster at a time). Construct with
// New() and start with Start(...).
type Service struct {
	// hasInboundConnection is set true while an authenticated RTMP source
	// is connected. Used both as a single-broadcaster guard and as the
	// loop condition for the packet pump.
	hasInboundConnection bool

	// pipe is the write end of the FLV bytestream handed to the stream
	// service via setStreamAsConnected. Closed on disconnect.
	pipe *io.PipeWriter

	// rtmpConnection is the live net.Conn of the broadcaster. Held so
	// Disconnect() can force-close it and so the read deadline can be
	// reset between packets.
	rtmpConnection net.Conn

	// setStreamAsConnected is invoked once at authenticated-connect time
	// with the read end of the FLV pipe. Set by Start.
	setStreamAsConnected func(*io.PipeReader)

	// setBroadcaster is invoked when the inbound onMetaData arrives with
	// the parsed broadcaster details. Set by Start.
	setBroadcaster func(models.Broadcaster)

	// configRepository is consulted at listener bind for the RTMP port and at
	// connect time for the list of valid stream keys.
	configRepository configrepository.ConfigRepository

	// cfg supplies the optional temporary stream key override set at
	// startup via the --streamkey CLI flag.
	cfg *config.Config
}

// Deps is the explicit dependency contract for the RTMP service.
type Deps struct {
	ConfigRepository configrepository.ConfigRepository
	Config           *config.Config
}

// New constructs an idle RTMP service. Call Start(...) to bind the
// listener and begin accepting connections.
func New(deps Deps) *Service {
	return &Service{configRepository: deps.ConfigRepository, cfg: deps.Config}
}
