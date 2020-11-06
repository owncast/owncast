package ffmpeg

import (
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type hlsVariantWriteMonitor struct {
	playlistEvents chan hlsWriteEvent
	segmentEvents  chan hlsWriteEvent
	logger         *log.Logger
}

type hlsWriteEvent struct {
	Dir  string
	Time time.Time
}

func newHlsVariantWriteMonitor(ctx context.Context, logger *log.Logger, notificationThreshold time.Duration) *hlsVariantWriteMonitor {
	h := &hlsVariantWriteMonitor{make(chan hlsWriteEvent), make(chan hlsWriteEvent), logger}
	go h.run(ctx, notificationThreshold)
	return h
}

func (h *hlsVariantWriteMonitor) SegmentWritten(dir string, when time.Time) {
	select {
	case h.segmentEvents <- hlsWriteEvent{dir, when}:
	case <-time.After(50 * time.Millisecond):
		h.logger.Errorf("unable to monitor segment write %q / %v", dir, when)
	}
}

func (h *hlsVariantWriteMonitor) VariantPlaylistWritten(dir string, when time.Time) {
	select {
	case h.playlistEvents <- hlsWriteEvent{dir, when}:
	case <-time.After(50 * time.Millisecond):
		h.logger.Errorf("unable to monitor playlist write %q / %v", dir, when)
	}
}

type hlsVariantWriteHistory struct {
	lastSegmentWrite  time.Time
	lastPlaylistWrite time.Time
}

// Run watches the handler traffic for lagging playlist writes
func (h *hlsVariantWriteMonitor) run(ctx context.Context, notificationThreshold time.Duration) {
	var (
		event  hlsWriteEvent
		drift  time.Duration
		viable bool

		variant  *hlsVariantWriteHistory
		ok       bool
		variants map[string]*hlsVariantWriteHistory = make(map[string]*hlsVariantWriteHistory)
	)

	for {
		select {
		case <-ctx.Done():
			return

		case event = <-h.playlistEvents:
			variant, ok = variants[event.Dir]
			if !ok {
				variant = &hlsVariantWriteHistory{event.Time, event.Time}
			} else {
				variant.lastPlaylistWrite = event.Time
			}
			variants[event.Dir] = variant
			drift = 0
			viable = false
			continue

		case event = <-h.segmentEvents:
			variant, ok = variants[event.Dir]
			if !ok {
				variant = &hlsVariantWriteHistory{event.Time, event.Time}
				drift = 0
			} else {
				drift = event.Time.Sub(variant.lastPlaylistWrite)
				viable = event.Time.Sub(variant.lastSegmentWrite) < notificationThreshold
				variant.lastSegmentWrite = event.Time
			}
			variants[event.Dir] = variant

		}

		if drift >= notificationThreshold && viable {
			h.logger.Warnf("HLS playlist at %q not updated with recent segments (%v drift)", event.Dir, drift)
		}
	}
}
