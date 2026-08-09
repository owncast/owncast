package replays

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/datastore"
)

var (
	testService         *Service
	testConfigRepo      configrepository.ConfigRepository
	fakeStreamID        = shortid.MustGenerate()
	fakeSegmentCount    = 300
	fakeStreamStartTime = time.Now()
	fakeConfigID        = ""
	fakeClipper         = shortid.MustGenerate()
	fakeClipStartTime   = 10
	fakeClipEndTime     = 15
)

// fakeSegmentDurationAt returns the simulated real (EXTINF) duration of
// segment i: alternating 2.5s and 1.5s, so media offsets exercise non-integer
// arithmetic.
func fakeSegmentDurationAt(i int) float64 {
	if i%2 == 0 {
		return 2.5
	}
	return 1.5
}

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

// populateFakeStream drives the real recording path: segments are reported
// via SegmentWritten (no timing yet), then a variant playlist write supplies
// the authoritative EXTINF durations, exactly like the live transcoder does.
func populateFakeStream() {
	queries := testService.datastore.GetQueries()

	recording := testService.NewRecording(fakeStreamID)
	fakeConfigID = recording.outputConfigurations[0].ID

	playlistDir := filepath.Join(os.TempDir(), "replay-test", "hls", "0")
	if err := os.MkdirAll(playlistDir, 0o750); err != nil {
		panic(err)
	}

	playlist := "#EXTM3U\n#EXT-X-VERSION:6\n#EXT-X-TARGETDURATION:3\n#EXT-X-MEDIA-SEQUENCE:0\n"
	for i := range fakeSegmentCount {
		fakeSegmentName := fmt.Sprintf("stream-%s-%d.ts", fakeStreamID, i)
		recording.SegmentWritten(filepath.Join("hls", "0", fakeSegmentName), 1000)
		playlist += fmt.Sprintf("#EXTINF:%.3f,\n%s\n", fakeSegmentDurationAt(i), fakeSegmentName)
	}
	playlist += "#EXT-X-ENDLIST\n"

	playlistPath := filepath.Join(playlistDir, "stream.m3u8")
	if err := os.WriteFile(playlistPath, []byte(playlist), 0o600); err != nil {
		panic(err)
	}
	recording.VariantPlaylistWritten(playlistPath)

	// One extra segment whose playlist entry never arrived: it must stay
	// invisible to playback until its timing is known.
	recording.SegmentWritten(filepath.Join("hls", "0", "stream-pending.ts"), 1000)

	if err := queries.SetStreamEnded(context.Background(), db.SetStreamEndedParams{
		ID:      fakeStreamID,
		EndTime: sql.NullTime{Time: fakeStreamStartTime.Add(20 * time.Minute), Valid: true},
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
		t.Fatal("expected playlist")
	}

	if err != nil {
		t.Error(err)
	}

	// The segment with no known timing is excluded.
	if int(p.Count()) != fakeSegmentCount {
		t.Error("expected", fakeSegmentCount, "segments, got", p.Count())
	}

	// EXTINF must reflect the real recorded durations, not the configured
	// segment length.
	if p.Segments[0].Duration != fakeSegmentDurationAt(0) {
		t.Error("expected first segment duration", fakeSegmentDurationAt(0), "got", p.Segments[0].Duration)
	}
	if p.Segments[1].Duration != fakeSegmentDurationAt(1) {
		t.Error("expected second segment duration", fakeSegmentDurationAt(1), "got", p.Segments[1].Duration)
	}

	if p.TargetDuration != math.Ceil(fakeSegmentDurationAt(0)) {
		t.Error("expected target duration", math.Ceil(fakeSegmentDurationAt(0)), "got", p.TargetDuration)
	}
}

func TestStreamMediaEnd(t *testing.T) {
	var expected float64
	for i := range fakeSegmentCount {
		expected += fakeSegmentDurationAt(i)
	}

	end, err := testService.GetStreamMediaEnd(fakeStreamID)
	if err != nil {
		t.Fatal(err)
	}

	if math.Abs(end-expected) > 0.001 {
		t.Error("expected media end", expected, "got", end)
	}
}

func TestClip(t *testing.T) {
	playlist := testService.NewPlaylistGenerator()
	clipID, duration, err := testService.AddClipForStream(fakeStreamID, "test clip", fakeClipper, float32(fakeClipStartTime), float32(fakeClipEndTime))
	if err != nil {
		t.Error(err)
	}

	clips, err := testService.GetAllClips()
	if err != nil {
		t.Error(err)
	}

	if len(clips) != 1 {
		t.Fatal("expected 1 clip, got", len(clips))
	}

	clip := clips[0]

	if clip.ID != clipID {
		t.Error("expected clip id", clipID, "got", clip.ID)
	}

	if clip.Manifest != fmt.Sprintf("/clip/%s", clipID) {
		t.Error("expected manifest id", fmt.Sprintf("/clip/%s", clipID), "got", clip.Manifest)
	}

	// The requested window is stored verbatim -- no snapping at creation.
	if clip.RelativeStartTime != float32(fakeClipStartTime) {
		t.Error("expected clip start time", fakeClipStartTime, "got", clip.RelativeStartTime)
	}

	if clip.RelativeEndTime != float32(fakeClipEndTime) {
		t.Error("expected clip end time", fakeClipEndTime, "got", clip.RelativeEndTime)
	}

	if duration != fakeClipEndTime-fakeClipStartTime {
		t.Error("expected clip duration", fakeClipEndTime-fakeClipStartTime, "got", duration)
	}

	p, err := playlist.GenerateMediaPlaylistForClipAndConfiguration(clipID, fakeConfigID)
	if err != nil {
		t.Error(err)
	}
	if p == nil {
		t.Fatal("expected playlist")
	}

	// Segments overlapping [10, 15) in media time. With durations
	// 2.5/1.5/2.5/... the offsets run 0, 2.5, 4, 6.5, 8, 10.5, 12, 14.5, 16:
	// the segments starting at 8, 10.5, 12 and 14.5 all contribute media to
	// the window.
	expectedSegmentCount := 4
	if int(p.Count()) != expectedSegmentCount {
		t.Error("expected", expectedSegmentCount, "segments, got", p.Count())
	}

	// The first segment of the clip is the one whose media range contains
	// second 10: the segment starting at offset 8.
	expectedFirstSegment := fmt.Sprintf("stream-%s-4.ts", fakeStreamID)
	if filepath.Base(p.Segments[0].URI) != expectedFirstSegment {
		t.Error("expected first clip segment", expectedFirstSegment, "got", filepath.Base(p.Segments[0].URI))
	}

	if p.TargetDuration != 3 {
		t.Error("expected target duration of 3, got", p.TargetDuration)
	}
}

func TestInvalidClipWindows(t *testing.T) {
	if _, _, err := testService.AddClipForStream(fakeStreamID, "bad clip", "", -1, 10); err == nil {
		t.Error("expected error for negative start time")
	}

	if _, _, err := testService.AddClipForStream(fakeStreamID, "bad clip", "", 20, 10); err == nil {
		t.Error("expected error for end before start")
	}

	if _, _, err := testService.AddClipForStream("nonexistent-stream", "bad clip", "", 0, 10); err == nil {
		t.Error("expected error for unknown stream")
	}
}
