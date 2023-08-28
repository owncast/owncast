package replays

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/utils"
	"github.com/teris-io/shortid"

	log "github.com/sirupsen/logrus"
)

type HLSRecorder struct {
	streamID  string
	startTime time.Time

	// The video variant configurations that were used for this stream.
	outputConfigurations []HLSOutputConfiguration

	datastore *data.Datastore
}

// NewRecording returns a new instance of the HLS recorder.
func NewRecording(streamID string) *HLSRecorder {
	// We don't support replaying offline clips.
	if streamID == "offline" {
		return nil
	}

	h := HLSRecorder{
		streamID:  streamID,
		startTime: time.Now(),
		datastore: data.GetDatastore(),
	}

	outputs := data.GetStreamOutputVariants()
	latency := data.GetStreamLatencyLevel()

	streamTitle := data.GetStreamTitle()
	validTitle := streamTitle != ""

	if err := h.datastore.GetQueries().InsertStream(context.Background(), db.InsertStreamParams{
		ID:          streamID,
		StartTime:   sql.NullTime{Time: h.startTime, Valid: true},
		StreamTitle: sql.NullString{String: streamTitle, Valid: validTitle},
	}); err != nil {
		log.Panicln(err)
	}

	// Create a reference of the output configurations that were used for this stream.
	for variantId, o := range outputs {
		configId := shortid.MustGenerate()

		if err := h.datastore.GetQueries().InsertOutputConfiguration(context.Background(), db.InsertOutputConfigurationParams{
			ID:               configId,
			Name:             o.Name,
			StreamID:         streamID,
			VariantID:        strconv.Itoa(variantId),
			SegmentDuration:  int32(latency.SecondsPerSegment),
			Bitrate:          int32(o.VideoBitrate),
			Framerate:        int32(o.Framerate),
			ResolutionWidth:  sql.NullInt32{Int32: int32(o.ScaledWidth), Valid: true},
			ResolutionHeight: sql.NullInt32{Int32: int32(o.ScaledHeight), Valid: true},
			Timestamp:        sql.NullTime{Time: time.Now(), Valid: true},
		}); err != nil {
			log.Panicln(err)
		}

		h.outputConfigurations = append(h.outputConfigurations, HLSOutputConfiguration{
			ID:              configId,
			Name:            o.Name,
			VideoBitrate:    o.VideoBitrate,
			ScaledWidth:     o.ScaledWidth,
			ScaledHeight:    o.ScaledHeight,
			Framerate:       o.Framerate,
			SegmentDuration: float64(latency.SegmentCount),
		})
	}
	return &h
}

// SegmentWritten is called when a segment is written to disk.
func (h *HLSRecorder) SegmentWritten(path string) {
	outputConfigurationIndexString := utils.GetIndexFromFilePath(path)
	outputConfigurationIndex, err := strconv.Atoi(outputConfigurationIndexString)
	if err != nil {
		log.Errorln("HLSRecorder segmentWritten error:", err)
		return
	}

	p := strings.ReplaceAll(path, "data/", "")
	relativeTimestamp := time.Since(h.startTime)

	if err := h.datastore.GetQueries().InsertSegment(context.Background(), db.InsertSegmentParams{
		ID:                    shortid.MustGenerate(),
		StreamID:              h.streamID,
		OutputConfigurationID: h.outputConfigurations[outputConfigurationIndex].ID,
		Path:                  p,
		RelativeTimestamp:     float32(relativeTimestamp.Seconds()),
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
