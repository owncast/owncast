package transcoder

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

const (
	pixelFormatYUV420P = "yuv420p"
	ffmpegEncoderX264  = "libx264"
)

// CodecProfile holds codec-dependent encoding parameters for an encoder.
type CodecProfile struct {
	PixelFormat  string
	Scaler       string
	ExtraFilters string
}

// Encoder represents a hardware/software video encoding backend (e.g., software CPU, NVENC, VA-API).
type Encoder interface {
	Name() string
	DisplayName() string
	GlobalFlags() []string
	ExtraArguments() []string
	ProfileForCodec(codec string) CodecProfile
	VariantFlags(v *HLSVariant) []string
	SupportedVideoCodecs() []string
	FFmpegEncoderForCodec(codec string) string
	GetPresetForLevel(l int) string
}

// SoftwareEncoder uses CPU-based encoding via libx264/libx265.
type SoftwareEncoder struct{}

func (e *SoftwareEncoder) Name() string {
	return "software"
}

func (e *SoftwareEncoder) DisplayName() string {
	return "Software (CPU)"
}

func (e *SoftwareEncoder) GlobalFlags() []string {
	return nil
}

func (e *SoftwareEncoder) ExtraArguments() []string {
	return []string{"-tune", "zerolatency"}
}

func (e *SoftwareEncoder) ProfileForCodec(codec string) CodecProfile {
	return CodecProfile{
		PixelFormat: pixelFormatYUV420P,
	}
}

func (e *SoftwareEncoder) VariantFlags(v *HLSVariant) []string {
	codec := v.videoCodec.Name()
	var flags []string

	if codec == videoCodecH264 {
		flags = append(flags,
			fmt.Sprintf("-x264-params:v:%d", v.index), "scenecut=0:open_gop=0",
		)
	}

	flags = append(flags,
		fmt.Sprintf("-bufsize:v:%d", v.index), fmt.Sprintf("%dk", v.getBufferSize()),
	)

	if codec == videoCodecH264 {
		flags = append(flags,
			fmt.Sprintf("-profile:v:%d", v.index), "high",
		)
	}

	return flags
}

func (e *SoftwareEncoder) SupportedVideoCodecs() []string {
	return []string{videoCodecH264}
}

func (e *SoftwareEncoder) FFmpegEncoderForCodec(codec string) string {
	switch codec {
	case videoCodecH264:
		return ffmpegEncoderX264
	default:
		log.Warnf("software encoder does not support codec %q, falling back to x264", codec)
		return ffmpegEncoderX264
	}
}

func (e *SoftwareEncoder) GetPresetForLevel(l int) string {
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
		log.Errorf("Invalid level for software preset %d, defaulting to %s", l, defaultPreset)
		return defaultPreset
	}
	return preset
}

var encoderRegistry = map[string]Encoder{
	"software":     &SoftwareEncoder{},
	"nvenc":        &NvencEncoder{},
	"vaapi":        &VaapiEncoder{},
	"qsv":          &QuicksyncEncoder{},
	"omx":          &OmxEncoder{},
	"v4l2m2m":      &V4l2m2mEncoder{},
	"videotoolbox": &VideoToolboxEncoder{},
}

func getEncoder(encoderType string) Encoder {
	if enc, ok := encoderRegistry[encoderType]; ok {
		return enc
	}
	return &SoftwareEncoder{}
}
