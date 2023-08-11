package replays

import "github.com/pkg/errors"

type HLSOutputConfiguration struct {
	ID              string
	StreamId        string
	VariantId       string
	Name            string
	VideoBitrate    int
	ScaledWidth     int
	ScaledHeight    int
	Framerate       int
	SegmentDuration float64
}

func (config *HLSOutputConfiguration) Validate() error {
	if config.VideoBitrate == 0 {
		return errors.New("video bitrate is unavailable")
	}

	if config.Framerate == 0 {
		return errors.New("video framerate is unavailable")
	}

	return nil
}
