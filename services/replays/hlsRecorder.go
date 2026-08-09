package replays

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/grafov/m3u8"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/utils"
)

// HLSRecorder persists the segments of a single live stream to the database so
// the stream can be replayed later.
type HLSRecorder struct {
	streamID  string
	startTime time.Time

	// The video variant configurations that were used for this stream.
	outputConfigurations []HLSOutputConfiguration

	// mu guards the per-variant media-time bookkeeping below.
	mu sync.Mutex
	// variantTimings tracks, per output configuration, the media-time offset
	// the next segment starts at and which segment files have already had
	// their timing recorded.
	variantTimings map[string]*variantTiming

	datastore *datastore.Datastore
}

// variantTiming is the running media-time state for one output configuration.
type variantTiming struct {
	nextOffset float64
	processed  map[string]bool
}

// newRecording returns a new instance of the HLS recorder for a stream, having
// already persisted the stream and its output configurations. Returns nil for
// the offline placeholder stream, which is never recorded.
func newRecording(streamID string, ds *datastore.Datastore, configRepository configrepository.ConfigRepository) *HLSRecorder {
	// We don't support replaying offline clips.
	if streamID == "offline" {
		return nil
	}

	log.Infoln("Recording replay of this stream:", streamID)

	h := HLSRecorder{
		streamID:       streamID,
		startTime:      time.Now(),
		datastore:      ds,
		variantTimings: map[string]*variantTiming{},
	}

	outputs := configRepository.GetStreamOutputVariants()
	latency := configRepository.GetStreamLatencyLevel()

	streamTitle := configRepository.GetStreamTitle()
	validTitle := streamTitle != ""

	if err := h.datastore.GetQueries().InsertStream(context.Background(), db.InsertStreamParams{
		ID:          streamID,
		StartTime:   sql.NullTime{Time: h.startTime, Valid: true},
		StreamTitle: sql.NullString{String: streamTitle, Valid: validTitle},
	}); err != nil {
		log.Errorln("unable to record stream:", err)
		return nil
	}

	// Create a reference of the output configurations that were used for this stream.
	for variantID, o := range outputs {
		configID := shortid.MustGenerate()

		if err := h.datastore.GetQueries().InsertOutputConfiguration(context.Background(), db.InsertOutputConfigurationParams{
			ID:               configID,
			Name:             o.Name,
			StreamID:         streamID,
			VariantID:        strconv.Itoa(variantID),
			SegmentDuration:  int64(latency.SecondsPerSegment),
			Bitrate:          int64(o.VideoBitrate),
			Framerate:        int64(o.GetFramerate()),
			ResolutionWidth:  sql.NullInt64{Int64: int64(o.ScaledWidth), Valid: true},
			ResolutionHeight: sql.NullInt64{Int64: int64(o.ScaledHeight), Valid: true},
			Timestamp:        sql.NullTime{Time: time.Now(), Valid: true},
		}); err != nil {
			log.Errorln("unable to record stream output configuration:", err)
			continue
		}

		h.outputConfigurations = append(h.outputConfigurations, HLSOutputConfiguration{
			ID:              configID,
			Name:            o.Name,
			VideoBitrate:    o.VideoBitrate,
			ScaledWidth:     o.ScaledWidth,
			ScaledHeight:    o.ScaledHeight,
			Framerate:       o.GetFramerate(),
			SegmentDuration: float64(latency.SecondsPerSegment),
		})
	}
	return &h
}

// SegmentWritten is called when a segment is written to storage. path is the
// public (relative) path the segment is served from and size is the segment
// file's size in bytes. The segment's media timing is not yet known; it is
// filled in by VariantPlaylistWritten when the transcoder reports the real
// EXTINF duration.
func (h *HLSRecorder) SegmentWritten(path string, size int64) {
	outputConfigurationIndexString := utils.GetIndexFromFilePath(path)
	outputConfigurationIndex, err := strconv.Atoi(outputConfigurationIndexString)
	if err != nil {
		log.Errorln("HLSRecorder segmentWritten error:", err)
		return
	}

	if outputConfigurationIndex >= len(h.outputConfigurations) {
		log.Errorln("HLSRecorder segmentWritten error: unknown output configuration index", outputConfigurationIndex)
		return
	}

	p := strings.ReplaceAll(path, "data/", "")

	if err := h.datastore.GetQueries().InsertSegment(context.Background(), db.InsertSegmentParams{
		ID:                    shortid.MustGenerate(),
		StreamID:              h.streamID,
		OutputConfigurationID: h.outputConfigurations[outputConfigurationIndex].ID,
		Path:                  p,
		Bytes:                 size,
		Timestamp:             sql.NullTime{Time: time.Now(), Valid: true},
	}); err != nil {
		log.Errorln(err)
	}
}

// VariantPlaylistWritten is called when the transcoder rewrites a variant's
// live playlist. The playlist carries the authoritative EXTINF duration for
// each segment, so this is where recorded segments get their real media
// timing: duration straight from EXTINF, and media offset as the running sum
// of every prior duration in the variant.
func (h *HLSRecorder) VariantPlaylistWritten(localFilePath string) {
	outputConfigurationIndexString := utils.GetIndexFromFilePath(localFilePath)
	outputConfigurationIndex, err := strconv.Atoi(outputConfigurationIndexString)
	if err != nil {
		log.Errorln("HLSRecorder variantPlaylistWritten error:", err)
		return
	}

	if outputConfigurationIndex >= len(h.outputConfigurations) {
		log.Errorln("HLSRecorder variantPlaylistWritten error: unknown output configuration index", outputConfigurationIndex)
		return
	}
	configID := h.outputConfigurations[outputConfigurationIndex].ID

	f, err := os.Open(localFilePath) // #nosec G304 -- path comes from the transcoder's own file writer
	if err != nil {
		log.Errorln("HLSRecorder variantPlaylistWritten error:", err)
		return
	}
	defer f.Close()

	playlist, listType, err := m3u8.DecodeFrom(f, true)
	if err != nil || listType != m3u8.MEDIA {
		log.Errorln("HLSRecorder variantPlaylistWritten: unable to parse variant playlist:", err)
		return
	}
	mediaPlaylist, ok := playlist.(*m3u8.MediaPlaylist)
	if !ok {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	timing := h.variantTimings[configID]
	if timing == nil {
		timing = &variantTiming{processed: map[string]bool{}}
		h.variantTimings[configID] = timing
	}

	// The live playlist is a sliding window: already-processed entries are
	// skipped, new entries at the tail get their timing assigned in order.
	window := map[string]bool{}
	for _, segment := range mediaPlaylist.Segments {
		if segment == nil {
			continue
		}
		filename := filepath.Base(segment.URI)
		window[filename] = true
		if timing.processed[filename] {
			continue
		}

		rows, err := h.datastore.GetQueries().SetSegmentTiming(context.Background(), db.SetSegmentTimingParams{
			Duration:              sql.NullFloat64{Float64: segment.Duration, Valid: true},
			MediaOffset:           sql.NullFloat64{Float64: timing.nextOffset, Valid: true},
			OutputConfigurationID: configID,
			Filename:              filename,
		})
		if err != nil {
			log.Errorln("HLSRecorder variantPlaylistWritten error:", err)
			return
		}
		if rows == 0 {
			// No stored segment row matches this playlist entry, so its
			// storage write failed. Stopping here keeps offsets contiguous;
			// the next playlist write retries.
			return
		}

		timing.processed[filename] = true
		timing.nextOffset += segment.Duration
	}

	// Entries that slid out of the live window never need to be looked at
	// again, so drop them from the dedupe set to keep it bounded.
	for filename := range timing.processed {
		if !window[filename] {
			delete(timing.processed, filename)
		}
	}
}

// StreamEnded is called when a stream is ended so the end time can be noted
// in the stream's metadata.
func (h *HLSRecorder) StreamEnded() {
	if err := h.datastore.GetQueries().SetStreamEnded(context.Background(), db.SetStreamEndedParams{
		ID:      h.streamID,
		EndTime: sql.NullTime{Time: time.Now(), Valid: true},
	}); err != nil {
		log.Errorln(err)
	}
}
