package replays

import "time"

// HLSSegment represents a single HLS segment.
type HLSSegment struct {
	ID                    string
	StreamID              string
	Timestamp             time.Time
	RelativeTimestamp     float32
	OutputConfigurationID string
	Path                  string
}
