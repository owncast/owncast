package transcoder

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

var supportedCodecs = map[string]string{
	"libx264":    "libx264",
	"h264_omx":   "omx",
	"h264_vaapi": "vaapi",
	"h264_nvenc": "NVIDEA nvenc",
	"h264_qsv":   "Intel Quicksync",
	// "h264_v4l2m2m": "Video4Linux",
}

type Libx264Codec struct {
}

func (c *Libx264Codec) Name() string {
	return "libx264"
}

func (c *Libx264Codec) GlobalFlags() string {
	return ""
}

func (c *Libx264Codec) PixelFormat() string {
	return "yuv420p"
}

func (c *Libx264Codec) ExtraArguments() string {
	return strings.Join([]string{
		"-tune", "zerolatency", // Option used for good for fast encoding and low-latency streaming (always includes iframes in each segment)
	}, " ")
}

func (c *Libx264Codec) ExtraFilters() string {
	return ""
}

func (c *Libx264Codec) VariantFlags(v *HLSVariant) string {
	bufferSize := int(float64(v.videoBitrate) * 1.2) // How often it checks the bitrate of encoded segments to see if it's too high/low.

	return strings.Join([]string{
		fmt.Sprintf("-bufsize:v:%d %dk", v.index, bufferSize), // How often the encoder checks the bitrate in order to meet average/max values
	}, " ")
}

func (c *Libx264Codec) GetPresetForLevel(l int) string {
	presetMapping := []string{
		"ultrafast",
		"superfast",
		"veryfast",
		"faster",
		"fast",
	}

	if l >= len(presetMapping) {
		return "superfast"
	}

	return presetMapping[l]
}

type OmxCodec struct {
}

func (c *OmxCodec) Name() string {
	return "h264_omx"
}

func (c *OmxCodec) GlobalFlags() string {
	return ""
}

func (c *OmxCodec) PixelFormat() string {
	return "yuv420p"
}

func (c *OmxCodec) ExtraArguments() string {
	return strings.Join([]string{
		"-tune", "zerolatency", // Option used for good for fast encoding and low-latency streaming (always includes iframes in each segment)
	}, " ")
}

func (c *OmxCodec) ExtraFilters() string {
	return ""
}

func (c *OmxCodec) VariantFlags(v *HLSVariant) string {
	return ""
}

func (c *OmxCodec) GetPresetForLevel(l int) string {
	presetMapping := []string{
		"ultrafast",
		"superfast",
		"veryfast",
		"faster",
		"fast",
	}

	if l >= len(presetMapping) {
		return "superfast"
	}

	return presetMapping[l]
}

type VaapiCodec struct {
}

func (c *VaapiCodec) Name() string {
	return "h264_vaapi"
}

func (c *VaapiCodec) GlobalFlags() string {
	flags := []string{
		// "-hwaccel", "vaapi",
		// "-hwaccel_output_format", "vaapi",
		"-vaapi_device", "/dev/dri/renderD128",
	}

	return strings.Join(flags, " ")
}

func (c *VaapiCodec) PixelFormat() string {
	return "vaapi_vld"
}

func (c *VaapiCodec) ExtraFilters() string {
	return "format=nv12,hwupload"
}

func (c *VaapiCodec) ExtraArguments() string {
	return ""
}

func (c *VaapiCodec) VariantFlags(v *HLSVariant) string {
	return ""
}

func (c *VaapiCodec) GetPresetForLevel(l int) string {
	presetMapping := []string{
		"ultrafast",
		"superfast",
		"veryfast",
		"faster",
		"fast",
	}

	if l >= len(presetMapping) {
		return "superfast"
	}

	return presetMapping[l]
}

type NvencCodec struct {
}

func (c *NvencCodec) Name() string {
	return "h264_nvenc"
}

func (c *NvencCodec) GlobalFlags() string {
	flags := []string{
		"-hwaccel cuda",
	}

	return strings.Join(flags, " ")
}

func (c *NvencCodec) PixelFormat() string {
	return "yuv420p"
}

func (c *NvencCodec) ExtraArguments() string {
	return ""
}

func (c *NvencCodec) ExtraFilters() string {
	return ""
}

func (c *NvencCodec) VariantFlags(v *HLSVariant) string {
	tuning := "ll" // low latency
	return fmt.Sprintf("-tune:v:%d %s", v.index, tuning)
}

func (c *NvencCodec) GetPresetForLevel(l int) string {
	presetMapping := []string{
		"p1",
		"p2",
		"p3",
		"p4",
		"p5",
	}

	if l >= len(presetMapping) {
		return "p3"
	}

	return presetMapping[l]
}

type QuicksyncCodec struct {
}

func (c *QuicksyncCodec) Name() string {
	return "h264_qsv"
}

func (c *QuicksyncCodec) GlobalFlags() string {
	return ""
}

func (c *QuicksyncCodec) PixelFormat() string {
	return "nv12"
}

func (c *QuicksyncCodec) ExtraArguments() string {
	return ""
}

func (c *QuicksyncCodec) ExtraFilters() string {
	return ""
}

func (c *QuicksyncCodec) VariantFlags(v *HLSVariant) string {
	return ""
}

func (c *QuicksyncCodec) GetPresetForLevel(l int) string {
	presetMapping := []string{
		"ultrafast",
		"superfast",
		"veryfast",
		"faster",
		"fast",
	}

	if l >= len(presetMapping) {
		return "superfast"
	}

	return presetMapping[l]
}

type Video4Linux struct{}

func (c *Video4Linux) Name() string {
	return "h264_v4l2m2m"
}

func (c *Video4Linux) GlobalFlags() string {
	return ""
}

func (c *Video4Linux) PixelFormat() string {
	return "nv21"
}

func (c *Video4Linux) ExtraArguments() string {
	return ""
}

func (c *Video4Linux) ExtraFilters() string {
	return ""
}

func (c *Video4Linux) VariantFlags(v *HLSVariant) string {
	return ""
}

func (c *Video4Linux) GetPresetForLevel(l int) string {
	presetMapping := []string{
		"ultrafast",
		"superfast",
		"veryfast",
		"faster",
		"fast",
	}

	if l >= len(presetMapping) {
		return "superfast"
	}

	return presetMapping[l]
}

// Codec represents a supported codec on the system.
type Codec interface {
	Name() string
	GlobalFlags() string
	PixelFormat() string
	ExtraArguments() string
	ExtraFilters() string
	VariantFlags(v *HLSVariant) string
	GetPresetForLevel(l int) string
}

// GetCodecs will return the supported codecs available on the system.
func GetCodecs(ffmpegPath string) []string {
	codecs := make([]string, 0)

	cmd := exec.Command(ffmpegPath, "-encoders")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorln(err)
		return codecs
	}

	response := string(out)
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		if strings.Contains(line, "H.264") {
			fields := strings.Fields(line)
			codec := fields[1]
			if _, supported := supportedCodecs[codec]; supported {
				codecs = append(codecs, codec)
			}
		}
	}

	return codecs
}

func getCodec(name string) Codec {
	switch name {
	case (&NvencCodec{}).Name():
		return &NvencCodec{}
	case (&VaapiCodec{}).Name():
		return &VaapiCodec{}
	case (&QuicksyncCodec{}).Name():
		return &QuicksyncCodec{}
	case (&OmxCodec{}).Name():
		return &OmxCodec{}
	case (&Video4Linux{}).Name():
		return &Video4Linux{}
	default:
		return &Libx264Codec{}
	}
}
