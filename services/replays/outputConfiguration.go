package replays

import "github.com/pkg/errors"

// HLSOutputConfiguration describes a single video output (variant) of a stream
// as it was configured at recording time.
type HLSOutputConfiguration struct {
	ID              string
	StreamID        string
	VariantID       string
	Name            string
	VideoBitrate    int
	ScaledWidth     int
	ScaledHeight    int
	Framerate       int
	SegmentDuration float64
}

// Validate ensures the configuration has the minimum information needed to
// build a playlist variant.
func (config *HLSOutputConfiguration) Validate() error {
	if config.VideoBitrate == 0 {
		return errors.New("video bitrate is unavailable")
	}

	if config.Framerate == 0 {
		return errors.New("video framerate is unavailable")
	}

	return nil
}
