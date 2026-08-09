package replays

import "time"

// HLSSegment represents a single recorded HLS segment.
type HLSSegment struct {
	ID                    string
	StreamID              string
	OutputConfigurationID string
	Path                  string

	// Duration is the segment's real media duration in seconds, as reported
	// by the transcoder in the variant playlist's EXTINF entry.
	Duration float64
	// MediaOffset is the segment's start position in media time: the running
	// sum of the durations of the segments before it in the same variant.
	MediaOffset float64
	// Bytes is the size of the segment file on disk.
	Bytes int64

	// Timestamp is the wall-clock time the segment was recorded. Used for
	// bookkeeping (healing crashed streams), not for playback timing.
	Timestamp time.Time
}
