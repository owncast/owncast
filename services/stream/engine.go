package stream

import (
	"io"

	"github.com/owncast/owncast/models"
)

// StreamEngine drives broadcaster ingest and transcoding. The local
// implementation runs RTMP ingest plus ffmpeg in-process (today's behavior);
// a future remote implementation manages a separate engine process. Core
// holds a StreamEngine and reacts to what it reports through StreamEvents, so
// core never needs to know where ingest and transcoding actually happen.
type StreamEngine interface {
	// Start begins accepting an inbound broadcast.
	Start() error
	// Stop force-disconnects any inbound broadcast and tears the engine down.
	Stop() error
}

// StreamEvents are the always-on, core-side reactions an engine drives as a
// broadcast comes and goes. The local engine triggers these in-process; a
// remote engine will call the same methods from its signaling channel, so
// core's handling of a stream going live/offline is identical regardless of
// where ingest happened.
type StreamEvents interface {
	// StreamConnected reports a broadcast going live, carrying the settings
	// it is being transcoded with.
	StreamConnected(broadcast *models.CurrentBroadcast)
	// BroadcasterSet reports the inbound source's metadata (codec, etc.).
	BroadcasterSet(broadcaster models.Broadcaster)
	// StreamDisconnected reports the broadcast ending. err is the transcode
	// exit cause if any, nil for a clean end.
	StreamDisconnected(err error)
}

// Compile-time guarantees for the local wiring.
var (
	_ StreamEngine = (*localStreamEngine)(nil)
	_ StreamEvents = (*Service)(nil)
)

// rtmpListener is the slice of the RTMP service the local engine drives.
// Keeping it an interface (which *rtmp.Service satisfies structurally) lets
// the engine be exercised without binding a real socket.
type rtmpListener interface {
	Start(onConnect func(*io.PipeReader), onBroadcaster func(models.Broadcaster))
	Disconnect()
}

// localStreamEngine is the in-process engine: it binds the RTMP listener and
// hands each inbound connection to the stream service's connect handler, which
// spins up the transcoder. This is today's single-process path expressed
// behind the StreamEngine seam; a remote engine replaces it in a later release.
type localStreamEngine struct {
	rtmp          rtmpListener
	onConnect     func(*io.PipeReader)
	onBroadcaster func(models.Broadcaster)
}

// Start binds the RTMP listener and begins accepting connections. The accept
// loop runs in its own goroutine, matching the prior inline behavior.
func (e *localStreamEngine) Start() error {
	go e.rtmp.Start(e.onConnect, e.onBroadcaster)
	return nil
}

// Stop force-disconnects the current inbound broadcaster, if any.
func (e *localStreamEngine) Stop() error {
	e.rtmp.Disconnect()
	return nil
}
