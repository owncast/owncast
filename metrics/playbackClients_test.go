package metrics

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func TestPlaybackClientStateTracksLatestValues(t *testing.T) {
	s := newTestService()

	s.RegisterPlaybackReport(PlaybackReport{
		ClientID:       "cmcd-session",
		ViewerID:       "viewer-1",
		BandwidthKbps:  2500,
		LatencySeconds: 4.5,
		BitrateKbps:    1200,
		HasErrorCount:  true,
		ClientReported: true,
	})
	s.RegisterPlaybackReport(PlaybackReport{
		ClientID:       "cmcd-session",
		ViewerID:       "viewer-1",
		BandwidthKbps:  1800,
		BitrateKbps:    800,
		ErrorCount:     1,
		HasErrorCount:  true,
		ClientReported: true,
	})

	clients := s.GetActivePlaybackClients()
	if len(clients) != 1 {
		t.Fatalf("expected 1 playback client, got %d", len(clients))
	}

	client := clients[0]
	if client.ViewerID != "viewer-1" {
		t.Errorf("expected viewer identity to be kept, got %q", client.ViewerID)
	}
	// The newest measurement wins; there is no history here.
	if client.Measurement(client.BandwidthKbps, client.BandwidthAt) == nil || client.BandwidthKbps != 1800 {
		t.Errorf("expected latest bandwidth 1800, got %f", client.BandwidthKbps)
	}
	// Latency was only reported once, but recently enough to still stand.
	if client.Measurement(client.LatencySeconds, client.LatencyAt) == nil || client.LatencySeconds != 4.5 {
		t.Errorf("expected latency 4.5 to persist, got %f", client.LatencySeconds)
	}
	// Counters describe the whole playback session, not one report.
	if client.ErrorCount != 1 || !client.HasErrorCount {
		t.Errorf("expected 1 error, got %f", client.ErrorCount)
	}
	if client.QualityVariantChanges != 1 || !client.HasQualityChanges {
		t.Errorf("expected 1 quality variant change, got %f", client.QualityVariantChanges)
	}
	if !client.IsSelfReporting() {
		t.Error("expected a client-reported measurement to be self reporting")
	}
	if client.Measurement(client.DownloadSeconds, client.DownloadAt) != nil {
		t.Error("expected an unmeasured download duration to stay unknown")
	}
}

func TestPlaybackClientStateSurvivesAggregation(t *testing.T) {
	s := newTestService()

	s.RegisterPlaybackReport(PlaybackReport{
		ClientID:      "viewer-1",
		ViewerID:      "viewer-1",
		BandwidthKbps: 900,
		ErrorCount:    2,
		HasErrorCount: true,
	})

	// The aggregation pass clears the windowed maps it consumes. The live
	// per-client view must keep describing a client that is still playing.
	s.playbackErrorCount = prometheus.NewGauge(prometheus.GaugeOpts{Name: "test_playback_errors"})
	s.metrics.m.Lock()
	s.collectPlaybackErrorCount()
	s.collectLowestBandwidth()
	s.collectQualityVariantChanges()
	s.metrics.m.Unlock()

	clients := s.GetActivePlaybackClients()
	if len(clients) != 1 {
		t.Fatalf("expected the client to survive aggregation, got %d clients", len(clients))
	}
	if clients[0].BandwidthKbps != 900 {
		t.Errorf("expected bandwidth 900, got %f", clients[0].BandwidthKbps)
	}
	if clients[0].ErrorCount != 2 {
		t.Errorf("expected the error count to survive aggregation, got %f", clients[0].ErrorCount)
	}
}

func TestServerObservedClientIsNotSelfReporting(t *testing.T) {
	s := newTestService()

	s.RegisterPlaybackReport(PlaybackReport{
		ClientID:        "viewer-1",
		ViewerID:        "viewer-1",
		BandwidthKbps:   4000,
		DownloadSeconds: 0.3,
	})

	clients := s.GetActivePlaybackClients()
	if len(clients) != 1 {
		t.Fatalf("expected 1 playback client, got %d", len(clients))
	}
	if clients[0].IsSelfReporting() {
		t.Error("expected a server-observed client to not be self reporting")
	}
	// The server watching segments leave says nothing about the player's
	// errors or the quality it chose, so those must read as unknown
	// rather than as a healthy zero.
	if clients[0].HasErrorCount {
		t.Error("expected a server-observed client to have no error information")
	}
	if clients[0].HasQualityChanges {
		t.Error("expected a server-observed client to have no quality information")
	}
}

func TestStalePlaybackClientsAreDroppedAndPruned(t *testing.T) {
	s := newTestService()

	s.RegisterPlaybackReport(PlaybackReport{ClientID: "gone", ViewerID: "gone", BandwidthKbps: 100})
	s.RegisterPlaybackReport(PlaybackReport{ClientID: "here", ViewerID: "here", BandwidthKbps: 200})

	s.metrics.m.Lock()
	s.clientPlayback["gone"].LastUpdate = time.Now().Add(-2 * playbackClientStaleTimeout)
	s.metrics.m.Unlock()

	clients := s.GetActivePlaybackClients()
	if len(clients) != 1 || clients[0].ClientID != "here" {
		t.Fatalf("expected only the active client, got %+v", clients)
	}

	s.metrics.m.Lock()
	s.prunePlaybackClients()
	remaining := len(s.clientPlayback)
	s.metrics.m.Unlock()

	if remaining != 1 {
		t.Errorf("expected the stale client to be pruned, %d remain", remaining)
	}
}

func TestGetActivePlaybackClientsIsOrdered(t *testing.T) {
	s := newTestService()

	// Two players behind one viewer identity plus a second viewer: the
	// list must not reshuffle between polls.
	s.RegisterPlaybackReport(PlaybackReport{ClientID: "sid-b", ViewerID: "viewer-1", BandwidthKbps: 1})
	s.RegisterPlaybackReport(PlaybackReport{ClientID: "sid-a", ViewerID: "viewer-1", BandwidthKbps: 2})
	s.RegisterPlaybackReport(PlaybackReport{ClientID: "sid-c", ViewerID: "viewer-0", BandwidthKbps: 3})

	want := []string{"sid-c", "sid-a", "sid-b"}
	for range 5 {
		clients := s.GetActivePlaybackClients()
		if len(clients) != len(want) {
			t.Fatalf("expected %d clients, got %d", len(want), len(clients))
		}
		for i, id := range want {
			if clients[i].ClientID != id {
				t.Fatalf("expected client %d to be %q, got %q", i, id, clients[i].ClientID)
			}
		}
	}
}

func TestStaleMeasurementsReadAsUnknown(t *testing.T) {
	s := newTestService()

	s.RegisterPlaybackReport(PlaybackReport{
		ClientID:       "player",
		ViewerID:       "viewer-1",
		LatencySeconds: 6,
		ClientReported: true,
	})

	// The player keeps reporting everything except its latency, which our
	// own player stops sending once latency is implausible. The stale
	// latency must not ride along as if it were current.
	s.metrics.m.Lock()
	s.clientPlayback["player"].LatencyAt = time.Now().Add(-2 * measurementFreshness)
	s.metrics.m.Unlock()

	s.RegisterPlaybackReport(PlaybackReport{
		ClientID:       "player",
		ViewerID:       "viewer-1",
		BandwidthKbps:  3000,
		ClientReported: true,
	})

	client := s.GetActivePlaybackClients()[0]
	if client.Measurement(client.LatencySeconds, client.LatencyAt) != nil {
		t.Error("expected a latency the player stopped reporting to read as unknown")
	}
	if client.Measurement(client.BandwidthKbps, client.BandwidthAt) == nil {
		t.Error("expected the freshly reported bandwidth to still be known")
	}
}

func TestResumedClientStartsANewSession(t *testing.T) {
	s := newTestService()

	s.RegisterPlaybackReport(PlaybackReport{
		ClientID:       "player",
		ViewerID:       "viewer-1",
		ErrorCount:     4,
		HasErrorCount:  true,
		ClientReported: true,
	})

	// Gone long enough to be stale, but not yet pruned by the two minute
	// aggregation pass.
	s.metrics.m.Lock()
	s.clientPlayback["player"].LastUpdate = time.Now().Add(-2 * playbackClientStaleTimeout)
	s.metrics.m.Unlock()

	s.RegisterPlaybackReport(PlaybackReport{
		ClientID:       "player",
		ViewerID:       "viewer-1",
		ErrorCount:     1,
		HasErrorCount:  true,
		ClientReported: true,
	})

	client := s.GetActivePlaybackClients()[0]
	if client.ErrorCount != 1 {
		t.Errorf("expected the previous session's errors to be dropped, got %f", client.ErrorCount)
	}
}

func TestQualityVariantChangeSurvivesACollectionWindow(t *testing.T) {
	s := newTestService()

	s.RegisterPlaybackReport(PlaybackReport{ClientID: "player", ViewerID: "viewer-1", BitrateKbps: 1200, ClientReported: true})

	// The aggregation pass runs between the two bitrates. The switch still
	// has to be counted, for the client and for the window.
	s.metrics.m.Lock()
	s.collectQualityVariantChanges()
	s.metrics.m.Unlock()

	s.RegisterPlaybackReport(PlaybackReport{ClientID: "player", ViewerID: "viewer-1", BitrateKbps: 800, ClientReported: true})

	client := s.GetActivePlaybackClients()[0]
	if client.QualityVariantChanges != 1 {
		t.Errorf("expected the variant switch to be counted, got %f", client.QualityVariantChanges)
	}
	if s.windowedQualityVariantChanges["player"] != 1 {
		t.Errorf("expected the window to count the switch too, got %f", s.windowedQualityVariantChanges["player"])
	}
}

func TestServerObservationFollowsThePlayersSession(t *testing.T) {
	s := newTestService()

	// A CMCD player reporting its own throughput.
	s.RegisterPlaybackReport(PlaybackReport{
		ClientID:       "cmcd-session",
		ViewerID:       "viewer-1",
		BandwidthKbps:  2500,
		ClientReported: true,
	})

	// A request the same player didn't decorate must be measured against
	// the session it already has, not split into a second client, and its
	// speed must be left to the player's own measurement.
	observation := s.ServerObservationFor("viewer-1")
	if observation.ClientID != "cmcd-session" {
		t.Errorf("expected the player's session, got %q", observation.ClientID)
	}
	if !observation.SpeedKnown {
		t.Error("expected the player's own throughput to own the speed metric")
	}
	if observation.Suppressed {
		t.Error("expected observation to still provide the download duration")
	}

	// Once that player's throughput goes stale the server measures speed
	// itself again.
	s.metrics.m.Lock()
	s.clientPlayback["cmcd-session"].BandwidthAt = time.Now().Add(-2 * measurementFreshness)
	s.metrics.m.Unlock()

	if s.ServerObservationFor("viewer-1").SpeedKnown {
		t.Error("expected a stale client throughput to hand the speed metric back to the server")
	}
}

func TestServerObservationSuppressedForBeaconingClients(t *testing.T) {
	s := newTestService()

	s.RegisterSelfReportingClient("viewer-1")

	if !s.ServerObservationFor("viewer-1").Suppressed {
		t.Error("expected a client that beacons its own measurements to suppress observation")
	}
	if s.ServerObservationFor("viewer-2").Suppressed {
		t.Error("expected an unknown viewer to be observed")
	}
}

func TestPlaybackStateAndUnmeasurableServerSample(t *testing.T) {
	serverOnly := newTestService()
	serverOnly.RegisterUnmeasurableServerSample("viewer-1", "viewer-1", "unmeasurable")
	clients := serverOnly.GetActivePlaybackClients()
	if len(clients) != 1 {
		t.Fatalf("expected one server-only client, got %d", len(clients))
	}
	if clients[0].ServerMeasurementStatus != "unmeasurable" {
		t.Errorf("expected unmeasurable status, got %q", clients[0].ServerMeasurementStatus)
	}

	reporting := newTestService()
	reporting.RegisterPlaybackReport(PlaybackReport{
		ClientID:       "player",
		ViewerID:       "viewer-1",
		PlayerState:    "p",
		ClientReported: true,
	})
	clients = reporting.GetActivePlaybackClients()
	if len(clients) != 1 || clients[0].PlayerState != "p" {
		t.Errorf("expected one playing client, got %+v", clients)
	}

	// A relay sample must not overwrite the status of a player that is
	// currently reporting its own measurements.
	reporting.RegisterUnmeasurableServerSample("player", "viewer-1", "unmeasurable")
	clients = reporting.GetActivePlaybackClients()
	if clients[0].ServerMeasurementStatus != "" {
		t.Errorf("expected client reporting to stay clear, got %q", clients[0].ServerMeasurementStatus)
	}
}
