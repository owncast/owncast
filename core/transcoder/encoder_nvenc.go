package transcoder

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// NvencEncoder uses NVIDIA GPU hardware encoding.
type NvencEncoder struct{}

func (e *NvencEncoder) Name() string {
	return "nvenc"
}

func (e *NvencEncoder) DisplayName() string {
	return "NVIDIA NVENC"
}

func (e *NvencEncoder) GlobalFlags() []string {
	return []string{
		"-hwaccel", "cuda",
	}
}

func (e *NvencEncoder) PixelFormat() string {
	return pixelFormatYUV420P
}

func (e *NvencEncoder) Scaler() string {
	return ""
}

func (e *NvencEncoder) ExtraFilters() string {
	return ""
}

func (e *NvencEncoder) ExtraArguments() []string {
	return nil
}

func (e *NvencEncoder) VariantFlags(v *HLSVariant) []string {
	return []string{
		fmt.Sprintf("-tune:v:%d", v.index), "ll",
	}
}

func (e *NvencEncoder) SupportedVideoCodecs() []string {
	return []string{videoCodecH264}
}

func (e *NvencEncoder) FFmpegEncoderForCodec(codec string) string {
	switch codec {
	case videoCodecH264:
		return "h264_nvenc"
	default:
		log.Warnf("nvenc encoder does not support codec %q, falling back to h264_nvenc", codec)
		return "h264_nvenc"
	}
}

func (e *NvencEncoder) GetPresetForLevel(l int) string {
	presetMapping := map[int]string{
		0: "p1",
		1: "p2",
		2: "p3",
		3: "p4",
		4: "p5",
	}
	preset, ok := presetMapping[l]
	if !ok {
		defaultPreset := presetMapping[2]
		log.Errorf("Invalid level for nvenc preset %d, defaulting to %s", l, defaultPreset)
		return defaultPreset
	}
	return preset
}
