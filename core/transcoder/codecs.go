package transcoder

import (
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

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

type VaapiCodec struct {
}

func (c *VaapiCodec) Name() string {
	return "h264_vaapi"
}

func (c *VaapiCodec) GlobalFlags() string {
	flags := []string{
		"-hwaccel vaapi",
		"-hwaccel_output_format vaapi",
		"-vaapi_device /dev/dri/renderD128",
	}

	return strings.Join(flags, " ")
}

func (c *VaapiCodec) PixelFormat() string {
	return "vaapi_vld"
}

func (c *VaapiCodec) ExtraArguments() string {
	return ""
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

type Codec interface {
	Name() string
	GlobalFlags() string
	PixelFormat() string
	ExtraArguments() string
}

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
			codecs = append(codecs, codec)
		}
	}

	return codecs
}
