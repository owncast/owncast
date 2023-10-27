package replays

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

var (
	fakeStreamId        = shortid.MustGenerate()
	fakeSegmentCount    = 300
	fakeSegmentDuration = 0
	fakeStreamStartTime = time.Now()
	fakeConfigId        = ""
	fakeClipper         = shortid.MustGenerate()
	fakeClipStartTime   = 10
	fakeClipEndTime     = 15
)

func TestMain(m *testing.M) {
	if err := data.SetupPersistence(":memory:"); err != nil {
		panic(err)
	}

	Setup()
	populateFakeStream()
	m.Run()
}

func populateFakeStream() {
	queries := data.GetDatastore().GetQueries()

	recording := NewRecording(fakeStreamId)
	fakeConfigId = recording.outputConfigurations[0].ID
	fakeSegmentDuration = data.GetStreamLatencyLevel().SecondsPerSegment // Seconds

	for i := 0; i < fakeSegmentCount; i++ {
		fakeSegmentName := fmt.Sprintf("%s-%d.ts", fakeStreamId, i)
		if err := queries.InsertSegment(context.Background(), db.InsertSegmentParams{
			ID:                    shortid.MustGenerate(),
			StreamID:              fakeStreamId,
			OutputConfigurationID: fakeConfigId,
			Path:                  filepath.Join(fakeStreamId, fakeConfigId, "0", fakeSegmentName),
			RelativeTimestamp:     float32(i * fakeSegmentDuration),
			Timestamp:             sql.NullTime{Time: fakeStreamStartTime.Add(time.Duration(fakeSegmentDuration * i)), Valid: true},
		}); err != nil {
			log.Errorln(err)
		}
	}

	if err := queries.SetStreamEnded(context.Background(), db.SetStreamEndedParams{
		ID:      fakeStreamId,
		EndTime: sql.NullTime{Time: fakeStreamStartTime.Add(time.Duration(fakeSegmentDuration * fakeSegmentCount)), Valid: true},
	}); err != nil {
		log.Errorln(err)
	}
}

func TestStream(t *testing.T) {
	playlist := NewPlaylistGenerator()
	stream, err := playlist.GetStream(fakeStreamId)
	if err != nil {
		t.Error(err)
	}

	if stream.ID != fakeStreamId {
		t.Error("expected stream id", fakeStreamId, "got", stream.ID)
	}
}

func TestPlaylist(t *testing.T) {
	playlist := NewPlaylistGenerator()
	p, err := playlist.GenerateMediaPlaylistForStreamAndConfiguration(fakeStreamId, fakeConfigId)
	if p == nil {
		t.Error("expected playlist")
	}

	if err != nil {
		t.Error(err)
	}

	if len(p.Segments) != fakeSegmentCount {
		t.Error("expected", fakeSegmentCount, "segments, got", len(p.Segments))
	}
}

func TestClip(t *testing.T) {
	segmentDuration := data.GetStreamLatencyLevel().SecondsPerSegment
	playlist := NewPlaylistGenerator()
	clipId, _, err := AddClipForStream(fakeStreamId, "test clip", fakeClipper, float32(fakeClipStartTime), float32(fakeClipEndTime))
	if err != nil {
		t.Error(err)
	}

	clips, err := GetAllClips()
	if err != nil {
		t.Error(err)
	}

	if len(clips) != 1 {
		t.Error("expected 1 clip, got", len(clips))
	}

	clip := clips[0]

	if clip.ID != clipId {
		t.Error("expected clip id", clipId, "got", clip.ID)
	}

	if clip.Manifest != fmt.Sprintf("/clip/%s", clipId) {
		t.Error("expected manifest id", fmt.Sprintf("/clip/%s", clipId), "got", clip.Manifest)
	}

	expectedStartTime := float32(utils.RoundDownToNearest(float32(fakeClipStartTime), segmentDuration))
	if clip.RelativeStartTime != expectedStartTime {
		t.Error("expected clip start time", fakeClipStartTime, "got", clip.RelativeStartTime)
	}

	expectedEndTime := float32(utils.RoundUpToNearest(float32(fakeClipEndTime), segmentDuration))
	if clip.RelativeEndTime != expectedEndTime {
		t.Error("expected clip end time", fakeClipEndTime, "got", clip.RelativeEndTime)
	}

	expectedDuration := expectedEndTime - expectedStartTime
	if float32(clip.DurationSeconds) != expectedDuration {
		t.Error("expected clip duration", expectedDuration, "got", clip.DurationSeconds)
	}

	p, err := playlist.GenerateMediaPlaylistForClipAndConfiguration(clipId, fakeConfigId)
	if err != nil {
		t.Error(err)
	}
	if p == nil {
		t.Error("expected playlist")
	}

	expectedSegmentCount := 3
	if len(p.Segments) != expectedSegmentCount {
		t.Error("expected", expectedSegmentCount, "segments, got", len(p.Segments))
	}

	if p.TargetDuration != float64(fakeSegmentDuration) {
		t.Error("expected target duration of", fakeSegmentDuration, "got", p.TargetDuration)
	}
}
