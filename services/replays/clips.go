package replays

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/utils"
)

// Clip represents a clip that has been created from a stream.
// A clip is a subset of a stream that has start and end seconds
// relative to the start of the stream.
type Clip struct {
	ID                string    `json:"id"`
	StreamID          string    `json:"stream_id"`
	ClippedBy         string    `json:"clipped_by,omitempty"`
	ClipTitle         string    `json:"title,omitempty"`
	StreamTitle       string    `json:"stream_title,omitempty"`
	RelativeStartTime float32   `json:"relativeStartTime"`
	RelativeEndTime   float32   `json:"relativeEndTime"`
	DurationSeconds   int       `json:"durationSeconds"`
	Manifest          string    `json:"manifest,omitempty"`
	Timestamp         time.Time `json:"timestamp"`
}

// GetAllClips will return all clips that have been recorded.
func (s *Service) GetAllClips() ([]*Clip, error) {
	clips, err := s.datastore.GetQueries().GetAllClips(context.Background())
	if err != nil {
		return nil, errors.WithMessage(err, "failure to get clips")
	}

	response := []*Clip{}
	for _, clip := range clips {
		c := Clip{
			ID:                clip.ID,
			ClipTitle:         clip.ClipTitle.String,
			StreamID:          clip.StreamID,
			StreamTitle:       clip.StreamTitle.String,
			RelativeStartTime: float32(clip.RelativeStartTime.Float64),
			RelativeEndTime:   float32(clip.RelativeEndTime.Float64),
			DurationSeconds:   numericToInt(clip.DurationSeconds),
			Timestamp:         clip.Timestamp.Time,
			Manifest:          fmt.Sprintf("/clip/%s", clip.ID),
		}
		response = append(response, &c)
	}
	return response, nil
}

// GetAllClipsForStream will return all clips that have been recorded for a stream.
func (s *Service) GetAllClipsForStream(streamID string) ([]*Clip, error) {
	clips, err := s.datastore.GetQueries().GetAllClipsForStream(context.Background(), streamID)
	if err != nil {
		return nil, errors.WithMessage(err, "failure to get clips")
	}

	response := []*Clip{}
	for _, clip := range clips {
		c := Clip{
			ID:                clip.ClipID,
			ClipTitle:         clip.ClipTitle.String,
			StreamID:          clip.StreamID,
			StreamTitle:       clip.StreamTitle.String,
			ClippedBy:         clip.ClippedBy.String,
			RelativeStartTime: float32(clip.RelativeStartTime.Float64),
			RelativeEndTime:   float32(clip.RelativeEndTime.Float64),
			Timestamp:         clip.Timestamp.Time,
			Manifest:          fmt.Sprintf("/clip/%s", clip.ClipID),
		}
		response = append(response, &c)
	}
	return response, nil
}

// AddClipForStream will save a new clip for a stream and return the new clip's
// ID and its duration in seconds.
func (s *Service) AddClipForStream(streamID, clipTitle, clippedBy string, relativeStartTimeSeconds, relativeEndTimeSeconds float32) (string, int, error) {
	playlistGenerator := s.NewPlaylistGenerator()

	// Verify this stream exists
	if _, err := playlistGenerator.GetStream(streamID); err != nil {
		return "", 0, errors.WithMessage(err, "stream not found")
	}

	// Verify this stream has at least one output configuration.
	configs, err := playlistGenerator.GetConfigurationsForStream(streamID)
	if err != nil {
		return "", 0, errors.WithMessage(err, "unable to get configurations for stream")
	}

	if len(configs) == 0 {
		return "", 0, errors.New("no configurations found for stream")
	}

	// We want the start and end seconds to be aligned to the segment so
	// round down and up the values to get a fully inclusive segment range.
	config := configs[0]
	segmentDuration := int(config.SegmentDuration)

	updatedRelativeStartTimeSeconds := utils.RoundDownToNearest(relativeStartTimeSeconds, segmentDuration)
	updatedRelativeEndTimeSeconds := utils.RoundUpToNearest(relativeEndTimeSeconds, segmentDuration)
	clipID := shortid.MustGenerate()
	duration := updatedRelativeEndTimeSeconds - updatedRelativeStartTimeSeconds

	err = s.datastore.GetQueries().InsertClip(context.Background(), db.InsertClipParams{
		ID:                clipID,
		StreamID:          streamID,
		ClipTitle:         sql.NullString{String: clipTitle, Valid: clipTitle != ""},
		RelativeStartTime: sql.NullFloat64{Float64: float64(updatedRelativeStartTimeSeconds), Valid: true},
		RelativeEndTime:   sql.NullFloat64{Float64: float64(updatedRelativeEndTimeSeconds), Valid: true},
		Timestamp:         sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return "", 0, errors.WithMessage(err, "failure to add clip")
	}

	return clipID, duration, nil
}

// GetFinalSegmentForStream will return the final known segment for a stream.
func (s *Service) GetFinalSegmentForStream(streamID string) (*HLSSegment, error) {
	segmentResponse, err := s.datastore.GetQueries().GetFinalSegmentForStream(context.Background(), streamID)
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

// numericToInt converts the dynamically-typed numeric result of a SQL
// expression (sqlc types computed columns as interface{}) into an int.
func numericToInt(v interface{}) int {
	switch n := v.(type) {
	case int64:
		return int(n)
	case int:
		return n
	case float64:
		return int(n)
	case float32:
		return int(n)
	default:
		return 0
	}
}
