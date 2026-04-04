package transcoder

import log "github.com/sirupsen/logrus"

// OmxEncoder uses OpenMAX hardware encoding (Raspberry Pi / ARM).
type OmxEncoder struct{}

func (e *OmxEncoder) Name() string {
	return "omx"
}

func (e *OmxEncoder) DisplayName() string {
	return "OpenMAX (omx)"
}

func (e *OmxEncoder) GlobalFlags() []string {
	return nil
}

func (e *OmxEncoder) PixelFormat() string {
	return pixelFormatYUV420P
}

func (e *OmxEncoder) Scaler() string {
	return ""
}

func (e *OmxEncoder) ExtraFilters() string {
	return ""
}

func (e *OmxEncoder) ExtraArguments() []string {
	return []string{
		"-tune", "zerolatency",
	}
}

func (e *OmxEncoder) VariantFlags(v *HLSVariant) []string {
	return nil
}

func (e *OmxEncoder) SupportedVideoCodecs() []string {
	return []string{videoCodecH264}
}

func (e *OmxEncoder) FFmpegEncoderForCodec(codec string) string {
	switch codec {
	case videoCodecH264:
		return "h264_omx"
	default:
		log.Warnf("omx encoder does not support codec %q, falling back to h264_omx", codec)
		return "h264_omx"
	}
}

func (e *OmxEncoder) GetPresetForLevel(l int) string {
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
		log.Errorf("Invalid level for omx preset %d, defaulting to %s", l, defaultPreset)
		return defaultPreset
	}
	return preset
}
