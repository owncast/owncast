package transcoder

import log "github.com/sirupsen/logrus"

// VideoToolboxEncoder uses Apple VideoToolbox hardware encoding (macOS).
type VideoToolboxEncoder struct{}

func (e *VideoToolboxEncoder) Name() string {
	return "videotoolbox"
}

func (e *VideoToolboxEncoder) DisplayName() string {
	return "VideoToolbox"
}

func (e *VideoToolboxEncoder) GlobalFlags() []string {
	return nil
}

func (e *VideoToolboxEncoder) ExtraArguments() []string {
	return nil
}

func (e *VideoToolboxEncoder) ProfileForCodec(codec string) CodecProfile {
	return CodecProfile{
		PixelFormat: "nv12",
	}
}

func (e *VideoToolboxEncoder) VariantFlags(v *HLSVariant) []string {
	if v.cpuUsageLevel >= 3 {
		return nil
	}
	return []string{"-realtime", "true"}
}

func (e *VideoToolboxEncoder) SupportedVideoCodecs() []string {
	return []string{videoCodecH264}
}

func (e *VideoToolboxEncoder) FFmpegEncoderForCodec(codec string) string {
	switch codec {
	case videoCodecH264:
		return "h264_videotoolbox"
	default:
		log.Warnf("videotoolbox encoder does not support codec %q, falling back to h264_videotoolbox", codec)
		return "h264_videotoolbox"
	}
}

func (e *VideoToolboxEncoder) GetPresetForLevel(l int) string {
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
		log.Errorf("Invalid level for videotoolbox preset %d, defaulting to %s", l, defaultPreset)
		return defaultPreset
	}
	return preset
}
