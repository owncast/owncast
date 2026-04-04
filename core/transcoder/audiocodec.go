package transcoder

const audioCodecAAC = "aac"

// aacFFmpegEncoder happens to match audioCodecAAC, but represents the ffmpeg
// encoder name, which could be "libfdk_aac" or another AAC encoder instead.
const aacFFmpegEncoder = "aac"

// AudioCodec represents an audio compression format (e.g., AAC, Opus).
type AudioCodec interface {
	Name() string
	DisplayName() string
	FFmpegEncoderName() string
	BitstreamFilters() []string
}

// AACCodec implements AudioCodec for AAC audio.
type AACCodec struct{}

func (c *AACCodec) Name() string {
	return audioCodecAAC
}

func (c *AACCodec) DisplayName() string {
	return "AAC"
}

// FFmpegEncoderName returns the ffmpeg encoder to use. We use the built-in
// "aac" encoder rather than "libfdk_aac" since it is universally available.
func (c *AACCodec) FFmpegEncoderName() string {
	return aacFFmpegEncoder
}

func (c *AACCodec) BitstreamFilters() []string {
	return []string{
		// "aac_adtstoasc", // Required for fMP4 segments, not needed yet.
	}
}

var audioCodecRegistry = map[string]AudioCodec{
	audioCodecAAC: &AACCodec{},
}

func getAudioCodec(name string) AudioCodec {
	if codec, ok := audioCodecRegistry[name]; ok {
		return codec
	}
	return &AACCodec{}
}
