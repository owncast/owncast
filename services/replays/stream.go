package replays

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
)

// Stream is a single recorded stream: a replay. Replays are the container
// clips are taken from; full-stream replay playback is admin-only for now.
type Stream struct {
	ID         string    `json:"id"`
	Title      string    `json:"title,omitempty"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime,omitempty"`
	InProgress bool      `json:"inProgress,omitempty"`
	Manifest   string    `json:"manifest,omitempty"`

	// DurationSeconds is the replay's recorded media duration.
	DurationSeconds float64 `json:"durationSeconds"`
	// TotalBytes is how much disk the replay's segments occupy.
	TotalBytes int64 `json:"totalBytes"`
	// ClipCount is how many clips reference this replay. Deleting the replay
	// deletes them too.
	ClipCount int `json:"clipCount"`
}

// GetStreams will return all replays, including the disk usage and clip counts
// the admin needs to manage them.
func (s *Service) GetStreams() ([]*Stream, error) {
	streams, err := s.datastore.GetQueries().GetStreamsWithDetails(context.Background())
	if err != nil {
		return nil, errors.WithMessage(err, "failure to get replays")
	}

	response := []*Stream{}
	for _, stream := range streams {
		s := Stream{
			ID:              stream.ID,
			Title:           stream.StreamTitle.String,
			StartTime:       stream.StartTime.Time,
			EndTime:         stream.EndTime.Time,
			InProgress:      !stream.EndTime.Valid,
			Manifest:        fmt.Sprintf("/replay/%s", stream.ID),
			DurationSeconds: stream.DurationSeconds,
			TotalBytes:      stream.TotalBytes,
			ClipCount:       int(stream.ClipCount),
		}
		response = append(response, &s)
	}
	return response, nil
}

// DeleteClip removes a single clip and its poster image. The recorded segments
// it referenced are left alone; they belong to the replay.
func (s *Service) DeleteClip(clipID string) error {
	rows, err := s.datastore.GetQueries().DeleteClip(context.Background(), clipID)
	if err != nil {
		return errors.WithMessage(err, "failure to delete clip")
	}

	if rows == 0 {
		return errors.New("clip not found")
	}

	RemoveClipThumbnail(clipID)

	return nil
}

// DeleteReplay removes a replay: its clips (and their posters), its segment
// files on disk, and all of its database rows. Segment files are deleted before
// the rows so a failure part way through leaves the ledger pointing at what
// still exists rather than orphaning files.
func (s *Service) DeleteReplay(streamID string) error {
	queries := s.datastore.GetQueries()
	ctx := context.Background()

	clipIDs, err := queries.GetClipIDsForStream(ctx, streamID)
	if err != nil {
		return errors.WithMessage(err, "unable to list clips for replay")
	}

	paths, err := queries.GetSegmentPathsForStream(ctx, streamID)
	if err != nil {
		return errors.WithMessage(err, "unable to list segments for replay")
	}

	for _, p := range paths {
		s.removeSegmentFile(p)
	}

	if err := queries.DeleteClipsForStream(ctx, streamID); err != nil {
		return errors.WithMessage(err, "unable to delete clips for replay")
	}

	if err := queries.DeleteSegmentsForStream(ctx, streamID); err != nil {
		return errors.WithMessage(err, "unable to delete segments for replay")
	}

	if err := queries.DeleteOutputConfigurationsForStream(ctx, streamID); err != nil {
		return errors.WithMessage(err, "unable to delete output configurations for replay")
	}

	rows, err := queries.DeleteStream(ctx, streamID)
	if err != nil {
		return errors.WithMessage(err, "unable to delete replay")
	}

	if rows == 0 {
		return errors.New("replay not found")
	}

	for _, clipID := range clipIDs {
		RemoveClipThumbnail(clipID)
	}

	return nil
}

// removeSegmentFile deletes a recorded segment from local storage. Remote
// (object storage) segments are stored as absolute URLs and are skipped: they
// are removed by the storage provider's own lifecycle, not from here.
func (s *Service) removeSegmentFile(storedPath string) {
	if storedPath == "" || strings.HasPrefix(storedPath, "http") {
		return
	}

	// Stored paths are relative to the data directory and must stay inside
	// it: a recorded path is never allowed to reach outside its own storage.
	// IsLocal is checked on the stored path before joining, since Join's
	// lexical Clean would strip any traversal before it could be caught.
	if !filepath.IsLocal(storedPath) {
		log.Warnln("refusing to delete segment outside of the data directory:", storedPath)
		return
	}

	fullPath := filepath.Join(config.DataDirectory, storedPath)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		log.Debugln("unable to remove segment file:", err)
	}
}
