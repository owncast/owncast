package transcoder

import log "github.com/sirupsen/logrus"

// V4l2m2mEncoder uses Video4Linux memory-to-memory hardware encoding.
type V4l2m2mEncoder struct{}

func (e *V4l2m2mEncoder) Name() string {
	return "v4l2m2m"
}

func (e *V4l2m2mEncoder) DisplayName() string {
	return "Video4Linux"
}

func (e *V4l2m2mEncoder) GlobalFlags() []string {
	return nil
}

func (e *V4l2m2mEncoder) PixelFormat() string {
	return "nv21"
}

func (e *V4l2m2mEncoder) Scaler() string {
	return ""
}

func (e *V4l2m2mEncoder) ExtraFilters() string {
	return ""
}

func (e *V4l2m2mEncoder) ExtraArguments() []string {
	return nil
}

func (e *V4l2m2mEncoder) VariantFlags(v *HLSVariant) []string {
	return nil
}

func (e *V4l2m2mEncoder) SupportedVideoCodecs() []string {
	return []string{videoCodecH264}
}

func (e *V4l2m2mEncoder) FFmpegEncoderForCodec(codec string) string {
	switch codec {
	case videoCodecH264:
		return "h264_v4l2m2m"
	default:
		log.Warnf("v4l2m2m encoder does not support codec %q, falling back to h264_v4l2m2m", codec)
		return "h264_v4l2m2m"
	}
}

func (e *V4l2m2mEncoder) GetPresetForLevel(l int) string {
	presetMapping := map[int]string{
		0: "ultrafast",
		1: "superfast",
		2: "veryfast",
		3: "faster",
		4: "fast",
	}
	preset, ok := presetMapping[l]
	if !ok {
		defaultPreset := presetMapping[1]
		log.Errorf("Invalid level for v4l preset %d, defaulting to %s", l, defaultPreset)
		return defaultPreset
	}
	return preset
}
