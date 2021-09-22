//nolint:goconst

package transcoder

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Codec represents a supported codec on the system.
type Codec interface {
	Name() string
	DisplayName() string
	GlobalFlags() string
	PixelFormat() string
	ExtraArguments() string
	ExtraFilters() string
	VariantFlags(v *HLSVariant) string
	GetPresetForLevel(l int) string
}

var supportedCodecs = map[string]string{
	(&Libx264Codec{}).Name(): "libx264",
	(&OmxCodec{}).Name():     "omx",
	(&VaapiCodec{}).Name():   "vaapi",
	(&NvencCodec{}).Name():   "NVIDIA nvenc",
}

// Libx264Codec represents an instance of the Libx264 Codec.
type Libx264Codec struct {
}

// Name returns the codec name.
func (c *Libx264Codec) Name() string {
	return "libx264"
}

// DisplayName returns the human readable name of the codec.
func (c *Libx264Codec) DisplayName() string {
	return "x264"
}

// GlobalFlags are the global flags used with this codec in the transcoder.
func (c *Libx264Codec) GlobalFlags() string {
	return ""
}

// PixelFormat is the pixel format required for this codec.
func (c *Libx264Codec) PixelFormat() string {
	return "yuv420p" //nolint:goconst
}

// ExtraArguments are the extra arguments used with this codec in the transcoder.
func (c *Libx264Codec) ExtraArguments() string {
	return strings.Join([]string{
		"-tune", "zerolatency", // Option used for good for fast encoding and low-latency streaming (always includes iframes in each segment)
	}, " ")
}

// ExtraFilters are the extra filters required for this codec in the transcoder.
func (c *Libx264Codec) ExtraFilters() string {
	return ""
}

// VariantFlags returns a string representing a single variant processed by this codec.
func (c *Libx264Codec) VariantFlags(v *HLSVariant) string {
	bufferSize := int(float64(v.videoBitrate) * 1.2) // How often it checks the bitrate of encoded segments to see if it's too high/low.

	return strings.Join([]string{
		fmt.Sprintf("-x264-params:v:%d \"scenecut=0:open_gop=0\"", v.index), // How often the encoder checks the bitrate in order to meet average/max values
		fmt.Sprintf("-bufsize:v:%d %dk", v.index, bufferSize),
		fmt.Sprintf("-profile:v:%d %s", v.index, "high"), // Encoding profile
	}, " ")
}

// GetPresetForLevel returns the string preset for this codec given an integer level.
func (c *Libx264Codec) GetPresetForLevel(l int) string {
	presetMapping := []string{
		"ultrafast",
		"superfast",
		"veryfast",
		"faster",
		"fast",
	}

	if l >= len(presetMapping) {
		return "superfast" //nolint:goconst
	}

	return presetMapping[l]
}

// OmxCodec represents an instance of the Omx codec.
type OmxCodec struct {
}

// Name returns the codec name.
func (c *OmxCodec) Name() string {
	return "h264_omx"
}

// DisplayName returns the human readable name of the codec.
func (c *OmxCodec) DisplayName() string {
	return "OpenMAX (omx)"
}

// GlobalFlags are the global flags used with this codec in the transcoder.
func (c *OmxCodec) GlobalFlags() string {
	return ""
}

// PixelFormat is the pixel format required for this codec.
func (c *OmxCodec) PixelFormat() string {
	return "yuv420p"
}

// ExtraArguments are the extra arguments used with this codec in the transcoder.
func (c *OmxCodec) ExtraArguments() string {
	return strings.Join([]string{
		"-tune", "zerolatency", // Option used for good for fast encoding and low-latency streaming (always includes iframes in each segment)
	}, " ")
}

// ExtraFilters are the extra filters required for this codec in the transcoder.
func (c *OmxCodec) ExtraFilters() string {
	return ""
}

// VariantFlags returns a string representing a single variant processed by this codec.
func (c *OmxCodec) VariantFlags(v *HLSVariant) string {
	return ""
}

// GetPresetForLevel returns the string preset for this codec given an integer level.
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

// VaapiCodec represents an instance of the Vaapi codec.
type VaapiCodec struct {
}

// Name returns the codec name.
func (c *VaapiCodec) Name() string {
	return "h264_vaapi"
}

// DisplayName returns the human readable name of the codec.
func (c *VaapiCodec) DisplayName() string {
	return "VA-API"
}

// GlobalFlags are the global flags used with this codec in the transcoder.
func (c *VaapiCodec) GlobalFlags() string {
	flags := []string{
		"-vaapi_device", "/dev/dri/renderD128",
	}

	return strings.Join(flags, " ")
}

// PixelFormat is the pixel format required for this codec.
func (c *VaapiCodec) PixelFormat() string {
	return "vaapi_vld"
}

// ExtraFilters are the extra filters required for this codec in the transcoder.
func (c *VaapiCodec) ExtraFilters() string {
	return "format=nv12,hwupload"
}

// ExtraArguments are the extra arguments used with this codec in the transcoder.
func (c *VaapiCodec) ExtraArguments() string {
	return ""
}

// VariantFlags returns a string representing a single variant processed by this codec.
func (c *VaapiCodec) VariantFlags(v *HLSVariant) string {
	return ""
}

// GetPresetForLevel returns the string preset for this codec given an integer level.
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

// NvencCodec represents an instance of the Nvenc Codec.
type NvencCodec struct {
}

// Name returns the codec name.
func (c *NvencCodec) Name() string {
	return "h264_nvenc"
}

// DisplayName returns the human readable name of the codec.
func (c *NvencCodec) DisplayName() string {
	return "nvidia nvenc"
}

// GlobalFlags are the global flags used with this codec in the transcoder.
func (c *NvencCodec) GlobalFlags() string {
	flags := []string{
		"-hwaccel cuda",
	}

	return strings.Join(flags, " ")
}

// PixelFormat is the pixel format required for this codec.
func (c *NvencCodec) PixelFormat() string {
	return "yuv420p"
}

// ExtraArguments are the extra arguments used with this codec in the transcoder.
func (c *NvencCodec) ExtraArguments() string {
	return ""
}

// ExtraFilters are the extra filters required for this codec in the transcoder.
func (c *NvencCodec) ExtraFilters() string {
	return ""
}

// VariantFlags returns a string representing a single variant processed by this codec.
func (c *NvencCodec) VariantFlags(v *HLSVariant) string {
	tuning := "ll" // low latency
	return fmt.Sprintf("-tune:v:%d %s", v.index, tuning)
}

// GetPresetForLevel returns the string preset for this codec given an integer level.
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

// QuicksyncCodec represents an instance of the Intel Quicksync Codec.
type QuicksyncCodec struct {
}

// Name returns the codec name.
func (c *QuicksyncCodec) Name() string {
	return "h264_qsv"
}

// DisplayName returns the human readable name of the codec.
func (c *QuicksyncCodec) DisplayName() string {
	return "Intel QuickSync"
}

// GlobalFlags are the global flags used with this codec in the transcoder.
func (c *QuicksyncCodec) GlobalFlags() string {
	return ""
}

// PixelFormat is the pixel format required for this codec.
func (c *QuicksyncCodec) PixelFormat() string {
	return "nv12"
}

// ExtraArguments are the extra arguments used with this codec in the transcoder.
func (c *QuicksyncCodec) ExtraArguments() string {
	return ""
}

// ExtraFilters are the extra filters required for this codec in the transcoder.
func (c *QuicksyncCodec) ExtraFilters() string {
	return ""
}

// VariantFlags returns a string representing a single variant processed by this codec.
func (c *QuicksyncCodec) VariantFlags(v *HLSVariant) string {
	return ""
}

// GetPresetForLevel returns the string preset for this codec given an integer level.
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

// Video4Linux represents an instance of the V4L Codec.
type Video4Linux struct{}

// Name returns the codec name.
func (c *Video4Linux) Name() string {
	return "h264_v4l2m2m"
}

// DisplayName returns the human readable name of the codec.
func (c *Video4Linux) DisplayName() string {
	return "Video4Linux"
}

// GlobalFlags are the global flags used with this codec in the transcoder.
func (c *Video4Linux) GlobalFlags() string {
	return ""
}

// PixelFormat is the pixel format required for this codec.
func (c *Video4Linux) PixelFormat() string {
	return "nv21"
}

// ExtraArguments are the extra arguments used with this codec in the transcoder.
func (c *Video4Linux) ExtraArguments() string {
	return ""
}

// ExtraFilters are the extra filters required for this codec in the transcoder.
func (c *Video4Linux) ExtraFilters() string {
	return ""
}

// VariantFlags returns a string representing a single variant processed by this codec.
func (c *Video4Linux) VariantFlags(v *HLSVariant) string {
	return ""
}

// GetPresetForLevel returns the string preset for this codec given an integer level.
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
