package replays

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/utils"
)

var (
	testService         *Service
	testConfigRepo      configrepository.ConfigRepository
	fakeStreamID        = shortid.MustGenerate()
	fakeSegmentCount    = 300
	fakeSegmentDuration = 0
	fakeStreamStartTime = time.Now()
	fakeConfigID        = ""
	fakeClipper         = shortid.MustGenerate()
	fakeClipStartTime   = 10
	fakeClipEndTime     = 15
)

func TestMain(m *testing.M) {
	ds, err := datastore.SetupPersistence(":memory:", os.TempDir())
	if err != nil {
		panic(err)
	}

	testConfigRepo = configrepository.New(ds)
	testService = New(ds, testConfigRepo)
	testService.Setup()

	populateFakeStream()
	os.Exit(m.Run())
}

func populateFakeStream() {
	queries := testService.datastore.GetQueries()

	recording := testService.NewRecording(fakeStreamID)
	fakeConfigID = recording.outputConfigurations[0].ID
	fakeSegmentDuration = testConfigRepo.GetStreamLatencyLevel().SecondsPerSegment // Seconds

	for i := 0; i < fakeSegmentCount; i++ {
		fakeSegmentName := fmt.Sprintf("%s-%d.ts", fakeStreamID, i)
		if err := queries.InsertSegment(context.Background(), db.InsertSegmentParams{
			ID:                    shortid.MustGenerate(),
			StreamID:              fakeStreamID,
			OutputConfigurationID: fakeConfigID,
			Path:                  filepath.Join(fakeStreamID, fakeConfigID, "0", fakeSegmentName),
			RelativeTimestamp:     float64(i * fakeSegmentDuration),
			Timestamp:             sql.NullTime{Time: fakeStreamStartTime.Add(time.Duration(fakeSegmentDuration * i)), Valid: true},
		}); err != nil {
			log.Errorln(err)
		}
	}

	if err := queries.SetStreamEnded(context.Background(), db.SetStreamEndedParams{
		ID:      fakeStreamID,
		EndTime: sql.NullTime{Time: fakeStreamStartTime.Add(time.Duration(fakeSegmentDuration * fakeSegmentCount)), Valid: true},
	}); err != nil {
		log.Errorln(err)
	}
}

func TestStream(t *testing.T) {
	playlist := testService.NewPlaylistGenerator()
	stream, err := playlist.GetStream(fakeStreamID)
	if err != nil {
		t.Error(err)
	}

	if stream.ID != fakeStreamID {
		t.Error("expected stream id", fakeStreamID, "got", stream.ID)
	}
}

func TestPlaylist(t *testing.T) {
	playlist := testService.NewPlaylistGenerator()
	p, err := playlist.GenerateMediaPlaylistForStreamAndConfiguration(fakeStreamID, fakeConfigID)
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
	segmentDuration := testConfigRepo.GetStreamLatencyLevel().SecondsPerSegment
	playlist := testService.NewPlaylistGenerator()
	clipID, _, err := testService.AddClipForStream(fakeStreamID, "test clip", fakeClipper, float32(fakeClipStartTime), float32(fakeClipEndTime))
	if err != nil {
		t.Error(err)
	}

	clips, err := testService.GetAllClips()
	if err != nil {
		t.Error(err)
	}

	if len(clips) != 1 {
		t.Error("expected 1 clip, got", len(clips))
	}

	clip := clips[0]

	if clip.ID != clipID {
		t.Error("expected clip id", clipID, "got", clip.ID)
	}

	if clip.Manifest != fmt.Sprintf("/clip/%s", clipID) {
		t.Error("expected manifest id", fmt.Sprintf("/clip/%s", clipID), "got", clip.Manifest)
	}

	expectedStartTime := float32(utils.RoundDownToNearest(float32(fakeClipStartTime), segmentDuration))
	if clip.RelativeStartTime != expectedStartTime {
		t.Error("expected clip start time", expectedStartTime, "got", clip.RelativeStartTime)
	}

	expectedEndTime := float32(utils.RoundUpToNearest(float32(fakeClipEndTime), segmentDuration))
	if clip.RelativeEndTime != expectedEndTime {
		t.Error("expected clip end time", expectedEndTime, "got", clip.RelativeEndTime)
	}

	expectedDuration := expectedEndTime - expectedStartTime
	if float32(clip.DurationSeconds) != expectedDuration {
		t.Error("expected clip duration", expectedDuration, "got", clip.DurationSeconds)
	}

	p, err := playlist.GenerateMediaPlaylistForClipAndConfiguration(clipID, fakeConfigID)
	if err != nil {
		t.Error(err)
	}
	if p == nil {
		t.Error("expected playlist")
		return
	}

	expectedSegmentCount := 3
	if len(p.Segments) != expectedSegmentCount {
		t.Error("expected", expectedSegmentCount, "segments, got", len(p.Segments))
	}

	if p.TargetDuration != float64(fakeSegmentDuration) {
		t.Error("expected target duration of", fakeSegmentDuration, "got", p.TargetDuration)
	}
}
