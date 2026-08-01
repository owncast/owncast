package metrics

import (
	"sort"
	"time"
)

// How long a playback client stays in the live per-client view after its
// last measurement. Players report every ~10s and server-side observation
// happens on every segment request, so a client silent for this long has
// stopped playing.
const playbackClientStaleTimeout = 60 * time.Second

// How long a single measurement is presented as current. Values are
// carried forward between reports, but a player that stops reporting one
// of them — ours drops live latency once it exceeds 100 seconds — must
// not keep showing its last good number while its other values refresh.
const measurementFreshness = 30 * time.Second

// ClientPlaybackState is the latest playback health of a single playback
// client. Unlike the windowed aggregate values it is not reset by the
// periodic collection pass, so the admin can see what each client is
// experiencing right now.
//
// Every measured value carries the time it was taken, because "not
// measured" and "measured as zero" mean different things to a streamer: a
// player that never reported its latency is not a player watching at the
// live edge.
type ClientPlaybackState struct {
	// ClientID is the playback identity: the player-supplied CMCD session
	// ID when available, otherwise the request-derived viewer identity.
	ClientID string
	// ViewerID is the request-derived identity of the viewer this client
	// belongs to. Multiple playback clients can share one viewer identity
	// when they're behind the same address with the same user agent.
	ViewerID string

	BandwidthKbps   float64
	LatencySeconds  float64
	DownloadSeconds float64
	BitrateKbps     float64

	BandwidthAt time.Time
	LatencyAt   time.Time
	DownloadAt  time.Time
	BitrateAt   time.Time

	// ErrorCount and QualityVariantChanges accumulate for as long as the
	// client keeps playing, so they describe this playback session rather
	// than an arbitrary collection window. Their Has* flags are false for
	// a client the server can only observe: watching segments leave the
	// server says nothing about whether the player hit an error or
	// switched quality, and reporting zero would claim it did neither.
	ErrorCount            float64
	QualityVariantChanges float64
	HasErrorCount         bool
	HasQualityChanges     bool

	// PlayerState is the latest CMCD playback state token: p (playing), a
	// (paused), w (buffering), k (seeking), or e (ended).
	PlayerState string
	// ServerMeasurementStatus explains why server-derived measurements are
	// unavailable. It is empty when the server has a usable sample or the
	// player reports its own measurements.
	ServerMeasurementStatus string
	// LastClientReport is when the player last reported its own
	// measurements; zero for clients only observed server-side.
	LastClientReport time.Time
	// LastUpdate is when any measurement for this client last arrived.
	LastUpdate time.Time
}

// IsSelfReporting returns true when the player supplied its own
// measurements recently enough that the server is still standing down
// from observing it. Past that window the server has resumed measuring
// the client itself.
func (c ClientPlaybackState) IsSelfReporting() bool {
	return !c.LastClientReport.IsZero() &&
		time.Since(c.LastClientReport) < selfReportSuppressionWindow
}

// Measurement returns a value taken at the given time, or nil when it was
// never taken or has gone stale.
func (c ClientPlaybackState) Measurement(value float64, takenAt time.Time) *float64 {
	if takenAt.IsZero() || time.Since(takenAt) >= measurementFreshness {
		return nil
	}
	return &value
}

// ServerObservation tells the HLS handler how to record its own
// measurement of a segment it just served to a viewer.
type ServerObservation struct {
	// ClientID is the playback client to record against. It is the
	// player's CMCD session when that player recently reported one, so
	// requests it didn't decorate don't split into a second client.
	ClientID string
	// Suppressed is true for clients that measure their own downloads and
	// beacon them to us, which makes server-side observation redundant.
	Suppressed bool
	// SpeedKnown is true when the player's own throughput measurement is
	// current, so the server should contribute only the download
	// duration and leave the speed to the richer client value.
	SpeedKnown bool
}

// ServerObservationFor resolves how a viewer's served segment should be
// measured, in a single locked read since it runs per media request.
func (s *Service) ServerObservationFor(viewerID string) ServerObservation {
	s.metrics.m.Lock()
	defer s.metrics.m.Unlock()

	observation := ServerObservation{ClientID: viewerID}

	if t, ok := s.selfReportingClients[viewerID]; ok && time.Since(t) < selfReportSuppressionWindow {
		observation.Suppressed = true
	}

	if id, ok := s.viewerPlaybackClients[viewerID]; ok {
		if state, ok := s.clientPlayback[id]; ok && time.Since(state.LastUpdate) < playbackClientStaleTimeout {
			observation.ClientID = id
		}
	}

	if state, ok := s.clientPlayback[observation.ClientID]; ok {
		observation.SpeedKnown = state.IsSelfReporting() &&
			state.Measurement(state.BandwidthKbps, state.BandwidthAt) != nil
	}

	return observation
}

// clientPlaybackState returns the live state for a report's client,
// creating it if this is the client's first measurement, and stamps it as
// updated. Callers must hold s.metrics.m.
func (s *Service) clientPlaybackState(report PlaybackReport) *ClientPlaybackState {
	now := time.Now()

	state, ok := s.clientPlayback[report.ClientID]
	// A state that already went stale belongs to a finished playback
	// session. Pruning only runs with the aggregation pass, so start the
	// new session fresh here rather than carrying the old session's
	// counters and last known values into it.
	if !ok || now.Sub(state.LastUpdate) >= playbackClientStaleTimeout {
		state = &ClientPlaybackState{ClientID: report.ClientID}
		s.clientPlayback[report.ClientID] = state
	}

	// A viewer identity only arrives with reports that came from a media
	// request or beacon, so keep the last known one rather than clearing
	// it.
	if report.ViewerID != "" {
		state.ViewerID = report.ViewerID
		// Remember which playback client a viewer's player identifies as
		// so requests that player leaves undecorated are still measured
		// against it instead of appearing as a second client.
		if report.ClientReported && report.ClientID != report.ViewerID {
			s.viewerPlaybackClients[report.ViewerID] = report.ClientID
		}
	}

	state.LastUpdate = now
	if report.ClientReported {
		state.LastClientReport = now
	}

	return state
}

// prunePlaybackClients drops clients that have stopped playing. Callers
// must hold s.metrics.m.
func (s *Service) prunePlaybackClients() {
	for id, state := range s.clientPlayback {
		if time.Since(state.LastUpdate) >= playbackClientStaleTimeout {
			delete(s.clientPlayback, id)
		}
	}

	for viewerID, clientID := range s.viewerPlaybackClients {
		if _, ok := s.clientPlayback[clientID]; !ok {
			delete(s.viewerPlaybackClients, viewerID)
		}
	}
}

// GetActivePlaybackClients returns a copy of the playback state of every
// client measured within the freshness window, ordered by viewer and
// client identity so repeated polls return a stable list.
func (s *Service) GetActivePlaybackClients() []ClientPlaybackState {
	s.metrics.m.Lock()
	defer s.metrics.m.Unlock()

	clients := make([]ClientPlaybackState, 0, len(s.clientPlayback))
	for _, state := range s.clientPlayback {
		if time.Since(state.LastUpdate) >= playbackClientStaleTimeout {
			continue
		}
		clients = append(clients, *state)
	}

	sort.Slice(clients, func(i, j int) bool {
		if clients[i].ViewerID != clients[j].ViewerID {
			return clients[i].ViewerID < clients[j].ViewerID
		}
		return clients[i].ClientID < clients[j].ClientID
	})

	return clients
}
