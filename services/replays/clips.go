package replays

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/pkg/errors"
	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/utils"
)

// Clip represents a clip that has been created from a replay.
// A clip is a window of a replay with start and end seconds in media time,
// relative to the start of the stream.
type Clip struct {
	ID                string    `json:"id"`
	StreamID          string    `json:"streamId"`
	ClippedBy         string    `json:"clippedBy,omitempty"`
	ClipTitle         string    `json:"title,omitempty"`
	StreamTitle       string    `json:"streamTitle,omitempty"`
	RelativeStartTime float32   `json:"relativeStartTime"`
	RelativeEndTime   float32   `json:"relativeEndTime"`
	DurationSeconds   int       `json:"durationSeconds"`
	Manifest          string    `json:"manifest,omitempty"`
	Thumbnail         string    `json:"thumbnail,omitempty"`
	Timestamp         time.Time `json:"timestamp"`
}

// clipThumbnailURL returns the public path of a clip's poster image, or an
// empty string when it has none.
func clipThumbnailURL(clipID string) string {
	if !utils.DoesFileExists(ClipThumbnailPath(clipID)) {
		return ""
	}

	return "/clip-thumbnail/" + ClipThumbnailFilename(clipID)
}

// GetAllClips will return all clips that have been created.
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
			ClippedBy:         clip.ClippedBy.String,
			StreamTitle:       clip.StreamTitle.String,
			RelativeStartTime: float32(clip.RelativeStartTime.Float64),
			RelativeEndTime:   float32(clip.RelativeEndTime.Float64),
			DurationSeconds:   numericToInt(clip.DurationSeconds),
			Timestamp:         clip.Timestamp.Time,
			Manifest:          fmt.Sprintf("/clip/%s", clip.ID),
			Thumbnail:         clipThumbnailURL(clip.ID),
		}
		response = append(response, &c)
	}
	return response, nil
}

// GetAllClipsForStream will return all clips created from a single replay.
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
			Thumbnail:         clipThumbnailURL(clip.ClipID),
		}
		response = append(response, &c)
	}
	return response, nil
}

// AddClipForStream will save a new clip for a stream and return the new clip's
// ID and its duration in seconds. The requested window is stored verbatim;
// segment selection happens at playback time by media-time overlap, so the
// clip plays exactly the requested range without snapping to segment
// boundaries at creation.
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

	if relativeStartTimeSeconds < 0 || relativeEndTimeSeconds <= relativeStartTimeSeconds {
		return "", 0, errors.New("invalid clip time window")
	}

	clipID := shortid.MustGenerate()
	duration := int(math.Round(float64(relativeEndTimeSeconds - relativeStartTimeSeconds)))

	err = s.datastore.GetQueries().InsertClip(context.Background(), db.InsertClipParams{
		ID:                clipID,
		StreamID:          streamID,
		ClippedBy:         sql.NullString{String: clippedBy, Valid: clippedBy != ""},
		ClipTitle:         sql.NullString{String: clipTitle, Valid: clipTitle != ""},
		RelativeStartTime: sql.NullFloat64{Float64: float64(relativeStartTimeSeconds), Valid: true},
		RelativeEndTime:   sql.NullFloat64{Float64: float64(relativeEndTimeSeconds), Valid: true},
		Timestamp:         sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return "", 0, errors.WithMessage(err, "failure to add clip")
	}

	return clipID, duration, nil
}

// GetStreamMediaEnd returns the media-time end (in seconds) of everything
// recorded so far for a stream: the maximum end offset across every variant.
// Zero means no segment has known timing yet.
func (s *Service) GetStreamMediaEnd(streamID string) (float64, error) {
	end, err := s.datastore.GetQueries().GetStreamMediaEnd(context.Background(), streamID)
	if err != nil {
		return 0, errors.Wrap(err, "unable to get stream media end")
	}
	return end, nil
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
