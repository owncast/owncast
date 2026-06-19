package replays

import "time"

// HLSSegment represents a single recorded HLS segment.
type HLSSegment struct {
	ID                    string
	StreamID              string
	Timestamp             time.Time
	RelativeTimestamp     float64
	OutputConfigurationID string
	Path                  string
}
