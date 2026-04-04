package transcoder

import log "github.com/sirupsen/logrus"

// VaapiEncoder uses VA-API hardware encoding (Intel/AMD GPUs).
type VaapiEncoder struct{}

func (e *VaapiEncoder) Name() string {
	return "vaapi"
}

func (e *VaapiEncoder) DisplayName() string {
	return "VA-API"
}

func (e *VaapiEncoder) GlobalFlags() []string {
	return []string{
		"-hwaccel", "vaapi",
		"-hwaccel_output_format", "vaapi",
		"-vaapi_device", "/dev/dri/renderD128",
	}
}

func (e *VaapiEncoder) PixelFormat() string {
	return "vaapi"
}

func (e *VaapiEncoder) Scaler() string {
	return "scale_vaapi"
}

func (e *VaapiEncoder) ExtraFilters() string {
	return "hwupload=extra_hw_frames=64,format=vaapi"
}

func (e *VaapiEncoder) ExtraArguments() []string {
	return nil
}

func (e *VaapiEncoder) VariantFlags(v *HLSVariant) []string {
	return nil
}

func (e *VaapiEncoder) SupportedVideoCodecs() []string {
	return []string{videoCodecH264}
}

func (e *VaapiEncoder) FFmpegEncoderForCodec(codec string) string {
	switch codec {
	case videoCodecH264:
		return "h264_vaapi"
	default:
		log.Warnf("vaapi encoder does not support codec %q, falling back to h264_vaapi", codec)
		return "h264_vaapi"
	}
}

func (e *VaapiEncoder) GetPresetForLevel(l int) string {
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
		log.Errorf("Invalid level for vaapi preset %d, defaulting to %s", l, defaultPreset)
		return defaultPreset
	}
	return preset
}
