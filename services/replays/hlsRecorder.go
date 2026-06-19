package replays

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

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

	datastore *datastore.Datastore
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
		streamID:  streamID,
		startTime: time.Now(),
		datastore: ds,
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
// public (relative) path the segment is served from.
func (h *HLSRecorder) SegmentWritten(path string) {
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
	relativeTimestamp := time.Since(h.startTime)

	if err := h.datastore.GetQueries().InsertSegment(context.Background(), db.InsertSegmentParams{
		ID:                    shortid.MustGenerate(),
		StreamID:              h.streamID,
		OutputConfigurationID: h.outputConfigurations[outputConfigurationIndex].ID,
		Path:                  p,
		RelativeTimestamp:     relativeTimestamp.Seconds(),
		Timestamp:             sql.NullTime{Time: time.Now(), Valid: true},
	}); err != nil {
		log.Errorln(err)
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
