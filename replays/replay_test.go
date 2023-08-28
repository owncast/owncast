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
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

var (
	fakeStreamId        = shortid.MustGenerate()
	fakeSegmentCount    = 300
	fakeSegmentDuration = 4 // Seconds
	fakeStreamStartTime = time.Now()
	fakeConfigId        = ""
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
