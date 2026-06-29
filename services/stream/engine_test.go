package stream

import (
	"io"
	"testing"
	"time"

	"github.com/owncast/owncast/models"
)

// fakeRTMP is a stand-in for *rtmp.Service that records how the local engine
// drives it, without binding a real socket.
type fakeRTMP struct {
	started       chan struct{}
	onConnect     func(*io.PipeReader)
	onBroadcaster func(models.Broadcaster)
	disconnected  bool
}

func (f *fakeRTMP) Start(onConnect func(*io.PipeReader), onBroadcaster func(models.Broadcaster)) {
	f.onConnect = onConnect
	f.onBroadcaster = onBroadcaster
	close(f.started)
}

func (f *fakeRTMP) Disconnect() { f.disconnected = true }

var _ rtmpListener = (*fakeRTMP)(nil)

// TestLocalStreamEngineDrivesRTMP proves the local engine starts ingest by
// delegating to the RTMP listener with the service's own callbacks, and tears
// down by disconnecting it.
func TestLocalStreamEngineDrivesRTMP(t *testing.T) {
	f := &fakeRTMP{started: make(chan struct{})}

	connectCalled := false
	broadcasterCalled := false
	e := &localStreamEngine{
		rtmp:          f,
		onConnect:     func(*io.PipeReader) { connectCalled = true },
		onBroadcaster: func(models.Broadcaster) { broadcasterCalled = true },
	}

	if err := e.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}

	select {
	case <-f.started:
	case <-time.After(2 * time.Second):
		t.Fatal("engine never started the RTMP listener")
	}

	// The listener must have received the engine's own callbacks verbatim.
	f.onConnect(nil)
	f.onBroadcaster(models.Broadcaster{})
	if !connectCalled {
		t.Error("engine did not wire onConnect through to the listener")
	}
	if !broadcasterCalled {
		t.Error("engine did not wire onBroadcaster through to the listener")
	}

	if err := e.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}
	if !f.disconnected {
		t.Error("Stop did not disconnect the RTMP listener")
	}
}

// TestBroadcasterSetRecordsMetadata proves BroadcasterSet (a StreamEvents
// method) records the inbound metadata that the hot status read returns, with
// no other service dependencies wired.
func TestBroadcasterSetRecordsMetadata(t *testing.T) {
	s := &Service{}
	if s.GetBroadcaster() != nil {
		t.Fatal("expected no broadcaster before any stream")
	}

	s.BroadcasterSet(models.Broadcaster{
		RemoteAddr:    "1.2.3.4",
		Time:          time.Now(),
		StreamDetails: models.InboundStreamDetails{VideoCodec: "H.264"},
	})

	got := s.GetBroadcaster()
	if got == nil {
		t.Fatal("expected a broadcaster after BroadcasterSet")
	}
	if got.RemoteAddr != "1.2.3.4" {
		t.Errorf("RemoteAddr = %q, want 1.2.3.4", got.RemoteAddr)
	}
	if got.StreamDetails.VideoCodec != "H.264" {
		t.Errorf("VideoCodec = %q, want H.264", got.StreamDetails.VideoCodec)
	}
}

// TestGetCurrentBroadcast proves the in-flight broadcast settings recorded on
// connect are returned to the hot read path.
func TestGetCurrentBroadcast(t *testing.T) {
	s := &Service{}
	if s.GetCurrentBroadcast() != nil {
		t.Fatal("expected no current broadcast before any stream")
	}

	bc := &models.CurrentBroadcast{LatencyLevel: models.GetLatencyLevel(2)}
	s.currentBroadcast = bc
	if s.GetCurrentBroadcast() != bc {
		t.Error("GetCurrentBroadcast did not return the recorded broadcast")
	}
}
