package transcoder

const videoCodecH264 = "h264"

// VideoCodec represents a video compression format (e.g., H.264, H.265).
type VideoCodec interface {
	Name() string
	DisplayName() string
}

// H264Codec implements VideoCodec for the H.264 format.
type H264Codec struct{}

func (c *H264Codec) Name() string {
	return videoCodecH264
}

func (c *H264Codec) DisplayName() string {
	return "H.264"
}

var videoCodecRegistry = map[string]VideoCodec{
	videoCodecH264: &H264Codec{},
}

func getVideoCodec(name string) VideoCodec {
	if codec, ok := videoCodecRegistry[name]; ok {
		return codec
	}
	return &H264Codec{}
}
