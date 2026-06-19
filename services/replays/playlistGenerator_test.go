package replays

import (
	"testing"
	"time"

	"github.com/grafov/m3u8"
)

// These tests exercise the pure playlist-construction helpers, which do not
// touch the datastore, so a zero-value generator is sufficient.
var (
	pureGenerator = &PlaylistGenerator{}
	testConfigs   = []HLSOutputConfiguration{
		{
			ID:           "1",
			VideoBitrate: 1000,
			Framerate:    30,
		},
		{
			ID:           "2",
			VideoBitrate: 2000,
			Framerate:    30,
		},
	}
)

var testSegments = []HLSSegment{
	{
		ID:                    "testSegmentId",
		StreamID:              "testStreamId",
		Timestamp:             time.Now(),
		OutputConfigurationID: "testOutputConfigId",
		Path:                  "hls/testStreamId/testOutputConfigId/testSegmentId.ts",
	},
}

func TestMasterPlaylist(t *testing.T) {
	playlist := pureGenerator.createNewMasterPlaylist()

	mediaPlaylists, err := pureGenerator.createMediaPlaylistForConfigurationAndSegments(&testConfigs[0], time.Now(), false, testSegments)
	if err != nil {
		t.Error(err)
	}
	playlist.Append("test", mediaPlaylists, m3u8.VariantParams{
		Bandwidth: uint32(testConfigs[0].VideoBitrate),
		FrameRate: float64(testConfigs[0].Framerate),
	})
	mediaPlaylists.Close()

	if playlist.Version() != 6 {
		t.Error("expected version 6, got", playlist.Version())
	}

	if !playlist.IndependentSegments() {
		t.Error("expected independent segments")
	}

	if playlist.Variants[0].Bandwidth != uint32(testConfigs[0].VideoBitrate) {
		t.Error("expected bandwidth", testConfigs[0].VideoBitrate, "got", playlist.Variants[0].Bandwidth)
	}

	if playlist.Variants[0].FrameRate != float64(testConfigs[0].Framerate) {
		t.Error("expected framerate", testConfigs[0].Framerate, "got", playlist.Variants[0].FrameRate)
	}
}

func TestCompletedMediaPlaylist(t *testing.T) {
	startTime := testSegments[0].Timestamp
	conf := testConfigs[0]

	// Create a completed media playlist.
	playlist, err := pureGenerator.createMediaPlaylistForConfigurationAndSegments(&conf, startTime, false, testSegments)
	if err != nil {
		t.Error(err)
	}

	if playlist.TargetDuration != conf.SegmentDuration {
		t.Error("expected target duration", conf.SegmentDuration, "got", playlist.TargetDuration)
	}

	// Verify it's marked as cachable.
	if playlist.Custom["#EXT-X-ALLOW-CACHE"].String() != "#EXT-X-ALLOW-CACHE" {
		t.Error("expected cachable playlist, tag not set")
	}

	// Verify it has the correct number of segments in the media playlist.
	if int(playlist.Count()) != len(testSegments) {
		t.Error("expected", len(testSegments), "segments, got", playlist.Count())
	}

	// Test the playlist version.
	if playlist.Version() != 6 {
		t.Error("expected version 6, got", playlist.Version())
	}

	// Verify the playlist type
	if playlist.MediaType != m3u8.VOD {
		t.Error("expected VOD playlist type, got type", playlist.MediaType)
	}

	// Verify the first segment URI.
	if playlist.Segments[0].URI != "/"+testSegments[0].Path {
		t.Error("expected segment URI", testSegments[0].Path, "got", playlist.Segments[0].URI)
	}
}

func TestInProgressMediaPlaylist(t *testing.T) {
	startTime := testSegments[0].Timestamp
	conf := testConfigs[0]

	// Create an in-progress media playlist.
	playlist, err := pureGenerator.createMediaPlaylistForConfigurationAndSegments(&conf, startTime, true, testSegments)
	if err != nil {
		t.Error(err)
	}

	// Verify it's NOT marked as cachable while in progress.
	if playlist.Custom != nil && playlist.Custom["#EXT-X-ALLOW-CACHE"] != nil && playlist.Custom["#EXT-X-ALLOW-CACHE"].String() == "#EXT-X-ALLOW-CACHE" {
		t.Error("expected non-cachable playlist when stream is still in progress")
	}

	// Verify it has the correct number of segments in the media playlist.
	if int(playlist.Count()) != len(testSegments) {
		t.Error("expected", len(testSegments), "segments, got", playlist.Count())
	}

	// Test the playlist version.
	if playlist.Version() != 6 {
		t.Error("expected version 6, got", playlist.Version())
	}

	// Verify the playlist type
	if playlist.MediaType != m3u8.EVENT {
		t.Error("expected EVENT playlist type, got type", playlist.MediaType)
	}
}
