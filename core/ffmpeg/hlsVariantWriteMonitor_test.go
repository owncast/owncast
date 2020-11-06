package ffmpeg

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

// time.Sleep calls in this file are present to allow for synchronization
// to reproduce effects consistently.

func Test_hlsVariantWriteMonitor_playlistLagging(t *testing.T) {
	logger, hook := test.NewNullLogger()
	moment, _ := time.Parse(time.RFC3339, "2000-01-01T00:00:00Z")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	monitor := newHlsVariantWriteMonitor(ctx, logger, time.Second*5)

	monitor.VariantPlaylistWritten("hls/0", moment)
	monitor.SegmentWritten("hls/0", moment)
	monitor.SegmentWritten("hls/0", moment.Add(time.Second*4))
	monitor.SegmentWritten("hls/0", moment.Add(time.Second*8))

	time.Sleep(50 * time.Millisecond)

	if len(hook.Entries) != 1 {
		t.Fatalf("wrong number of log entries\nexpected: 1\n  actual: %d", len(hook.Entries))
	}

	if hook.LastEntry().Level != logrus.WarnLevel {
		t.Errorf("incorrect log error level output: %v", hook.LastEntry().Level)
	}

	expected := `HLS playlist at "hls/0" not updated with recent segments (8s drift)`
	if hook.LastEntry().Message != expected {
		t.Errorf("log message does not match\nexpected: %q\n  actual: %q", expected, hook.LastEntry().Message)
	}
}

func Test_hlsVariantWriteMonitor_segmentDelay(t *testing.T) {
	logger, hook := test.NewNullLogger()
	moment, _ := time.Parse(time.RFC3339, "2000-01-01T00:00:00Z")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	monitor := newHlsVariantWriteMonitor(ctx, logger, time.Second*5)

	monitor.VariantPlaylistWritten("hls/0", moment)
	monitor.SegmentWritten("hls/0", moment)
	monitor.SegmentWritten("hls/0", moment.Add(time.Second*20))

	time.Sleep(50 * time.Millisecond)

	if len(hook.Entries) != 0 {
		t.Fatalf("wrong number of log entries\nexpected: 0\n  actual: %d", len(hook.Entries))
	}
}

func Test_hlsVariantWriteMonitor_multiplePlaylists(t *testing.T) {
	logger, hook := test.NewNullLogger()
	moment, _ := time.Parse(time.RFC3339, "2000-01-01T00:00:00Z")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	monitor := newHlsVariantWriteMonitor(ctx, logger, time.Second*5)

	monitor.VariantPlaylistWritten("hls/0", moment)
	monitor.SegmentWritten("hls/0", moment)
	monitor.SegmentWritten("hls/0", moment.Add(time.Second*4))

	monitor.VariantPlaylistWritten("hls/1", moment)
	monitor.SegmentWritten("hls/1", moment)
	monitor.SegmentWritten("hls/1", moment.Add(time.Second*4))
	monitor.SegmentWritten("hls/1", moment.Add(time.Second*8))

	time.Sleep(50 * time.Millisecond)

	if len(hook.Entries) != 1 {
		t.Fatalf("wrong number of log entries\nexpected: 1\n  actual: %d", len(hook.Entries))
	}

	if hook.LastEntry().Level != logrus.WarnLevel {
		t.Errorf("incorrect log error level output: %v", hook.LastEntry().Level)
	}

	expected := `HLS playlist at "hls/1" not updated with recent segments (8s drift)`
	if hook.LastEntry().Message != expected {
		t.Errorf("log message does not match\nexpected: %q\n  actual: %q", expected, hook.LastEntry().Message)
	}
}

func Test_hlsVariantWriteMonitor_writeTimeouts(t *testing.T) {
	logger, hook := test.NewNullLogger()
	moment, _ := time.Parse(time.RFC3339, "2000-01-01T00:00:00Z")

	ctx, cancel := context.WithCancel(context.Background())

	monitor := newHlsVariantWriteMonitor(ctx, logger, time.Second*5)
	cancel() // stop the run loop
	time.Sleep(50 * time.Millisecond)

	monitor.VariantPlaylistWritten("hls/0", moment)
	time.Sleep(25 * time.Millisecond)
	monitor.SegmentWritten("hls/0", moment)
	time.Sleep(100 * time.Millisecond)

	if len(hook.Entries) != 2 {
		t.Log(hook.Entries)
		t.Fatalf("wrong number of log entries\nexpected: 2\n  actual: %d", len(hook.Entries))
	}

	expectedLogs := []struct {
		Level   logrus.Level
		Message string
	}{
		{
			logrus.ErrorLevel,
			`unable to monitor playlist write "hls/0" / 2000-01-01 00:00:00 +0000 UTC`,
		},
		{
			logrus.ErrorLevel,
			`unable to monitor segment write "hls/0" / 2000-01-01 00:00:00 +0000 UTC`,
		},
	}

	for i, expected := range expectedLogs {
		if hook.Entries[i].Level != expected.Level {
			t.Errorf("incorrect log error level for entry %d\nexpected: %v\n  actual: %v", i, expected.Level, hook.Entries[i].Level)
		}
		if hook.Entries[i].Message != expected.Message {
			t.Errorf("incorrect log message for entry %d\nexpected: %v\n  actual: %v", i, expected.Message, hook.Entries[i].Message)
		}
	}
}
