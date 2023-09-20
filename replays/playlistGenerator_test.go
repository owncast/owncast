package replays

import (
	"testing"
	"time"

	"github.com/grafov/m3u8"
)

var (
	generator = NewPlaylistGenerator()
	config    = []HLSOutputConfiguration{
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

var segments = []HLSSegment{
	{
		ID:                    "testSegmentId",
		StreamID:              "testStreamId",
		Timestamp:             time.Now(),
		OutputConfigurationID: "testOutputConfigId",
		Path:                  "hls/testStreamId/testOutputConfigId/testSegmentId.ts",
	},
}

func TestMasterPlaylist(t *testing.T) {
	playlist := generator.createNewMasterPlaylist()

	mediaPlaylists, err := generator.createMediaPlaylistForConfigurationAndSegments(&config[0], time.Now(), false, segments)
	playlist.Append("test", mediaPlaylists, m3u8.VariantParams{
		Bandwidth: uint32(config[0].VideoBitrate),
		FrameRate: float64(config[0].Framerate),
	})
	mediaPlaylists.Close()

	if err != nil {
		t.Error(err)
	}

	if playlist.Version() != 6 {
		t.Error("expected version 6, got", playlist.Version())
	}

	if !playlist.IndependentSegments() {
		t.Error("expected independent segments")
	}

	if playlist.Variants[0].Bandwidth != uint32(config[0].VideoBitrate) {
		t.Error("expected bandwidth", config[0].VideoBitrate, "got", playlist.Variants[0].Bandwidth)
	}

	if playlist.Variants[0].FrameRate != float64(config[0].Framerate) {
		t.Error("expected framerate", config[0].Framerate, "got", playlist.Variants[0].FrameRate)
	}
}

func TestCompletedMediaPlaylist(t *testing.T) {
	startTime := segments[0].Timestamp
	conf := config[0]

	// Create a completed media playlist.
	playlist, err := generator.createMediaPlaylistForConfigurationAndSegments(&conf, startTime, false, segments)
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
	if int(playlist.Count()) != len(segments) {
		t.Error("expected", len(segments), "segments, got", playlist.Count())
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
	if playlist.Segments[0].URI != "/"+segments[0].Path {
		t.Error("expected segment URI", segments[0].Path, "got", playlist.Segments[0].URI)
	}
}

func TestInProgressMediaPlaylist(t *testing.T) {
	startTime := segments[0].Timestamp
	conf := config[0]

	// Create a completed media playlist.
	playlist, err := generator.createMediaPlaylistForConfigurationAndSegments(&conf, startTime, true, segments)
	if err != nil {
		t.Error(err)
	}

	// Verify it's marked as cachable.
	if playlist.Custom != nil && playlist.Custom["#EXT-X-ALLOW-CACHE"].String() == "#EXT-X-ALLOW-CACHE" {
		t.Error("expected non-achable playlist when stream is still in progress")
	}

	// Verify it has the correct number of segments in the media playlist.
	if int(playlist.Count()) != len(segments) {
		t.Error("expected", len(segments), "segments, got", playlist.Count())
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
