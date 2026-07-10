package metrics

import "testing"

func newTestService() *Service {
	s := New(Deps{})
	s.metrics = new(CollectedMetrics)
	return s
}

func TestRegisterPlaybackReport(t *testing.T) {
	s := newTestService()

	s.RegisterPlaybackReport(PlaybackReport{
		ClientID:       "a",
		BandwidthKbps:  500,
		LatencySeconds: 3.2,
		BitrateKbps:    1200,
		HasErrorCount:  true,
	})

	if s.windowedBandwidths["a"] != 500 {
		t.Errorf("expected bandwidth 500, got %f", s.windowedBandwidths["a"])
	}
	if s.windowedLatencies["a"] != 3.2 {
		t.Errorf("expected latency 3.2, got %f", s.windowedLatencies["a"])
	}
	// A zero error count must still create the client's entry: the stream
	// health overview uses the error map as its denominator, so healthy
	// clients that never register would make the overview disappear.
	if _, ok := s.windowedErrorCounts["a"]; !ok {
		t.Error("expected a zero error count entry for a healthy client")
	}
	if s.windowedQualityVariantChanges["a"] != 0 {
		t.Error("first bitrate observation must not count as a variant change")
	}

	// A bitrate change is a quality variant change; errors accumulate.
	s.RegisterPlaybackReport(PlaybackReport{
		ClientID:      "a",
		BitrateKbps:   2400,
		ErrorCount:    1,
		HasErrorCount: true,
	})
	if s.windowedQualityVariantChanges["a"] != 1 {
		t.Errorf("expected 1 variant change, got %f", s.windowedQualityVariantChanges["a"])
	}
	if s.windowedErrorCounts["a"] != 1 {
		t.Errorf("expected error count 1, got %f", s.windowedErrorCounts["a"])
	}

	// The same bitrate again is not a change.
	s.RegisterPlaybackReport(PlaybackReport{ClientID: "a", BitrateKbps: 2400})
	if s.windowedQualityVariantChanges["a"] != 1 {
		t.Errorf("expected variant changes to stay at 1, got %f", s.windowedQualityVariantChanges["a"])
	}

	// A report without HasErrorCount (server-side observation carries no
	// error information) must not add the client to the error denominator.
	s.RegisterPlaybackReport(PlaybackReport{ClientID: "b", BandwidthKbps: 900})
	if _, ok := s.windowedErrorCounts["b"]; ok {
		t.Error("server-observed clients must not enter the error denominator")
	}
}
