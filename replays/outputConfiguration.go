package replays

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
