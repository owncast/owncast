package transcoder

import log "github.com/sirupsen/logrus"

// QuicksyncEncoder uses Intel Quick Sync Video hardware encoding.
type QuicksyncEncoder struct{}

func (e *QuicksyncEncoder) Name() string {
	return "qsv"
}

func (e *QuicksyncEncoder) DisplayName() string {
	return "Intel QuickSync"
}

func (e *QuicksyncEncoder) GlobalFlags() []string {
	return []string{
		"-init_hw_device", "qsv=hw",
		"-filter_hw_device", "hw",
	}
}

func (e *QuicksyncEncoder) ExtraArguments() []string {
	return nil
}

func (e *QuicksyncEncoder) ProfileForCodec(codec string) CodecProfile {
	return CodecProfile{
		PixelFormat:  "qsv",
		Scaler:       "scale_qsv",
		ExtraFilters: "hwupload=extra_hw_frames=64,format=qsv",
	}
}

func (e *QuicksyncEncoder) VariantFlags(v *HLSVariant) []string {
	return nil
}

func (e *QuicksyncEncoder) SupportedVideoCodecs() []string {
	return []string{videoCodecH264}
}

func (e *QuicksyncEncoder) FFmpegEncoderForCodec(codec string) string {
	switch codec {
	case videoCodecH264:
		return "h264_qsv"
	default:
		log.Warnf("quicksync encoder does not support codec %q, falling back to h264_qsv", codec)
		return "h264_qsv"
	}
}

func (e *QuicksyncEncoder) GetPresetForLevel(l int) string {
	presetMapping := map[int]string{
		0: "veryfast",
		1: "fast",
		2: "medium",
		3: "slow",
		4: "veryslow",
	}
	preset, ok := presetMapping[l]
	if !ok {
		defaultPreset := presetMapping[2]
		log.Errorf("Invalid level for quicksync preset %d, defaulting to %s", l, defaultPreset)
		return defaultPreset
	}
	return preset
}
