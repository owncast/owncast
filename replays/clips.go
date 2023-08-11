package replays

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/utils"
	"github.com/pkg/errors"
	"github.com/teris-io/shortid"
)

// Clip represents a clip that has been created from a stream.
// A clip is a subset of a stream that has has start and end seconds
// relative to the start of the stream.
type Clip struct {
	ID                string    `json:"id"`
	StreamId          string    `json:"stream_id"`
	ClippedBy         string    `json:"clipped_by,omitempty"`
	ClipTitle         string    `json:"title,omitempty"`
	StreamTitle       string    `json:"stream_title,omitempty"`
	RelativeStartTime float32   `json:"relativeStartTime"`
	RelativeEndTime   float32   `json:"relativeEndTime"`
	DurationSeconds   int       `json:"durationSeconds"`
	Manifest          string    `json:"manifest,omitempty"`
	Timestamp         time.Time `json:"timestamp"`
}

// GetClips will return all clips that have been recorded.
func GetAllClips() ([]*Clip, error) {
	clips, err := data.GetDatastore().GetQueries().GetAllClips(context.Background())
	if err != nil {
		return nil, errors.WithMessage(err, "failure to get clips")
	}

	response := []*Clip{}
	for _, clip := range clips {
		s := Clip{
			ID:                clip.ID,
			ClipTitle:         clip.ClipTitle.String,
			StreamId:          clip.StreamID,
			StreamTitle:       clip.StreamTitle.String,
			RelativeStartTime: float32(clip.RelativeStartTime.Float64),
			RelativeEndTime:   float32(clip.RelativeEndTime.Float64),
			DurationSeconds:   int(clip.DurationSeconds),
			Timestamp:         clip.Timestamp.Time,
			Manifest:          fmt.Sprintf("/clip/%s", clip.ID),
		}
		response = append(response, &s)
	}
	return response, nil
}

// GetAllClipsForStream will return all clips that have been recorded for a stream.
func GetAllClipsForStream(streamId string) ([]*Clip, error) {
	clips, err := data.GetDatastore().GetQueries().GetAllClipsForStream(context.Background(), streamId)
	if err != nil {
		return nil, errors.WithMessage(err, "failure to get clips")
	}

	response := []*Clip{}
	for _, clip := range clips {
		s := Clip{
			ID:                clip.ClipID,
			ClipTitle:         clip.ClipTitle.String,
			StreamTitle:       clip.StreamTitle.String,
			RelativeStartTime: float32(clip.RelativeStartTime.Float64),
			RelativeEndTime:   float32(clip.RelativeEndTime.Float64),
			Timestamp:         clip.Timestamp.Time,
			Manifest:          fmt.Sprintf("/clips/%s", clip.ClipID),
		}
		response = append(response, &s)
	}
	return response, nil
}

// AddClipForStream will save a new clip for a stream.
func AddClipForStream(streamId, clipTitle, clippedBy string, relativeStartTimeSeconds, relativeEndTimeSeconds float32) (string, int, error) {
	playlistGenerator := NewPlaylistGenerator()

	// Verify this stream exists
	if _, err := playlistGenerator.GetStream(streamId); err != nil {
		return "", 0, errors.WithMessage(err, "stream not found")
	}

	// Verify this stream has at least one output configuration.
	configs, err := playlistGenerator.GetConfigurationsForStream(streamId)
	if err != nil {
		return "", 0, errors.WithMessage(err, "unable to get configurations for stream")
	}

	if len(configs) == 0 {
		return "", 0, errors.New("no configurations found for stream")
	}

	// We want the start and end seconds to be aligned to the segment so
	// round up and down the values to get a fully inclusive segment range.
	config := configs[0]
	segmentDuration := int(config.SegmentDuration)

	updatedRelativeStartTimeSeconds := utils.RoundDownToNearest(relativeStartTimeSeconds, segmentDuration)
	updatedRelativeEndTimeSeconds := utils.RoundUpToNearest(relativeEndTimeSeconds, segmentDuration)
	clipId := shortid.MustGenerate()
	duration := updatedRelativeEndTimeSeconds - updatedRelativeStartTimeSeconds

	err = data.GetDatastore().GetQueries().InsertClip(context.Background(), db.InsertClipParams{
		ID:                clipId,
		StreamID:          streamId,
		ClipTitle:         sql.NullString{String: clipTitle, Valid: clipTitle != ""},
		RelativeStartTime: sql.NullFloat64{Float64: float64(updatedRelativeStartTimeSeconds), Valid: true},
		RelativeEndTime:   sql.NullFloat64{Float64: float64(updatedRelativeEndTimeSeconds), Valid: true},
		Timestamp:         sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return "", 0, errors.WithMessage(err, "failure to add clip")
	}

	return clipId, duration, nil
}

// GetFinalSegmentForStream will return the final known segment for a stream.
func GetFinalSegmentForStream(streamId string) (*HLSSegment, error) {
	segmentResponse, err := data.GetDatastore().GetQueries().GetFinalSegmentForStream(context.Background(), streamId)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get final segment for stream")
	}

	segment := HLSSegment{
		ID:                    segmentResponse.ID,
		StreamID:              segmentResponse.StreamID,
		OutputConfigurationID: segmentResponse.OutputConfigurationID,
		Path:                  segmentResponse.Path,
		RelativeTimestamp:     segmentResponse.RelativeTimestamp,
		Timestamp:             segmentResponse.Timestamp.Time,
	}

	return &segment, nil
}
