// Package stream is the live-streaming lifecycle service: stream connect /
// disconnect handling, viewer stats, status reporting, storage provider
// setup, and the offline content state. It owns the state that used to
// live as package-level globals in the legacy `core` package; that package
// now exists only as a thin backward-compatibility shim that delegates
// here.
package stream

import (
	"context"
	"sync"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/geoip"
	"github.com/owncast/owncast/services/rtmp"
	"github.com/owncast/owncast/services/transcoder"
	"github.com/owncast/owncast/yp"
)

// activeViewerPurgeTimeout is the silence-window after which a viewer is
// dropped from the active set.
const activeViewerPurgeTimeout = 15 * time.Second

// Service is the stream subsystem. Construct one in main.go with New(),
// register it with the legacy core compat shim via core.Init(svc), then
// call Start(ctx) to bring up the storage, transcoder, RTMP listener and
// associated infrastructure.
type Service struct {
	// stats owns viewer-count, connect timing, etc. statsMu guards the
	// stats fields that are touched from goroutines (viewer-set updates,
	// chat-client removal).
	stats   *models.Stats
	statsMu sync.RWMutex

	// storage is the active HLS storage provider (local filesystem or S3).
	storage models.StorageProvider

	// transcoder holds the active ffmpeg child while a stream is online.
	// Replaced when a stream comes up, nil'd when it goes down.
	transcoder *transcoder.Transcoder

	// yp talks to the Owncast directory.
	yp *yp.YP

	// broadcaster carries the metadata of the currently-connected RTMP
	// source. Nil between streams.
	broadcaster *models.Broadcaster

	// currentBroadcast holds the in-flight broadcast settings (latency
	// level, output variants) captured at stream-connect time.
	currentBroadcast *models.CurrentBroadcast

	// offlineCleanupTimer fires once N minutes after a stream goes
	// offline to wipe HLS content. onlineCleanupTicker fires every N
	// minutes while a stream is up to prune old segments.
	offlineCleanupTimer *time.Timer
	onlineCleanupTicker *time.Ticker

	// onlineTimerCancelFunc cancels the delayed live-notification
	// dispatcher when a stream ends before the notification fires.
	onlineTimerCancelFunc context.CancelFunc

	// lastNotified is the last time we sent the live-going-online
	// notification; used to debounce re-notifications when a stream
	// flaps.
	lastNotified *time.Time

	// HLS transcoder plumbing. handler routes finished playlists/segments
	// to the storage provider; fileWriter feeds the transcoder's output
	// into handler.
	handler    transcoder.HLSHandler
	fileWriter transcoder.FileWriterReceiverService

	// thumbnailGen owns the periodic thumbnail+preview snapshotter that
	// runs while a stream is online. Lazily started on stream-connect,
	// stopped on stream-disconnect.
	thumbnailGen *transcoder.ThumbnailGenerator

	// geoIPClient performs the asynchronous viewer-geolocation lookup.
	geoIPClient *geoip.Client

	// rtmp is the inbound RTMP ingest service. Start hands it
	// stream-lifecycle callbacks; admin disconnect goes through it too.
	rtmp *rtmp.Service
}

// Deps is the explicit-dependency contract the service requires at
// construction time. As more collaborators (chat, transcoder, webhooks,
// etc.) migrate to constructor-injected services, they appear here.
type Deps struct {
	Rtmp *rtmp.Service
}

// New constructs an idle stream Service. Call Start(ctx) to bring up the
// storage backend, transcoder, RTMP listener, and associated lifecycle.
func New(deps Deps) *Service {
	return &Service{
		geoIPClient: geoip.NewClient(),
		fileWriter:  transcoder.FileWriterReceiverService{},
		rtmp:        deps.Rtmp,
	}
}
