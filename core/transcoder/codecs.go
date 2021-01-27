package transcoder

import (
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// ./ffmpeg -hwaccel qsv -c:v h264_qsv -i ~/Videos/BigBuckBunny.mp4 -vf hwdownload,format=nv12 -pix_fmt yuv420p -c:v h264_qsv test.mp4

// ffmpeg -y -loglevel debug -init_hw_device qsv=hw -filter_hw_device hw -hwaccel qsv -hwaccel_output_format qsv \
//-i simpsons.mp4 -vf 'format=qsv,hwupload=extra_hw_frames=64'  \
//-c:v h264_qsv \
//-bf 3 -b:v 15M -maxrate:v 15M -bufsize:v 2M -r:v 30 -c:a copy -f mp4 trolled.mp4

type QuickSync struct {
	BaseCodec
}

func (c *QuickSync) Name() string {
	return "Intel QuickSync"
}

type Nvenc struct {
	BaseCodec
}

func (c *Nvenc) Name() string {
	return "NVENC"
}

type Codec interface {
	Name() string
	ToString() string
}

type BaseCodec struct {
	Bitrate int
	Quality string
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
