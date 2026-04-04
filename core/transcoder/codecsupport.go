package transcoder

import (
	"os/exec"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
)

// EncoderCodecSupport describes which codecs an encoder supports on this system.
type EncoderCodecSupport struct {
	EncoderType          string   `json:"encoderType"`
	EncoderDisplayName   string   `json:"encoderDisplayName"`
	SupportedVideoCodecs []string `json:"supportedVideoCodecs"`
}

// GetAvailableVideoEncoders probes ffmpeg to discover which encoder+codec combinations
// are available on this system. Returns only encoders that have at least one
// supported codec available.
func GetAvailableVideoEncoders(ffmpegPath string) []EncoderCodecSupport {
	ffmpegEncoders := getFFmpegEncoders(ffmpegPath)

	// Sort encoder names for deterministic output order.
	encoderNames := make([]string, 0, len(encoderRegistry))
	for name := range encoderRegistry {
		encoderNames = append(encoderNames, name)
	}
	sort.Strings(encoderNames)

	var available []EncoderCodecSupport
	for _, name := range encoderNames {
		enc := encoderRegistry[name]
		var supportedCodecs []string
		for _, codecName := range enc.SupportedVideoCodecs() {
			ffmpegEncoderName := enc.FFmpegEncoderForCodec(codecName)
			if _, ok := ffmpegEncoders[ffmpegEncoderName]; ok {
				supportedCodecs = append(supportedCodecs, codecName)
			}
		}
		if len(supportedCodecs) > 0 {
			available = append(available, EncoderCodecSupport{
				EncoderType:          enc.Name(),
				EncoderDisplayName:   enc.DisplayName(),
				SupportedVideoCodecs: supportedCodecs,
			})
		}
	}

	// Always include software encoder even if ffmpeg probe fails,
	// as libx264 is almost universally available.
	hasSoftware := false
	for _, enc := range available {
		if enc.EncoderType == "software" {
			hasSoftware = true
			break
		}
	}
	if !hasSoftware {
		available = append([]EncoderCodecSupport{{
			EncoderType:          "software",
			EncoderDisplayName:   "Software (CPU)",
			SupportedVideoCodecs: []string{"h264"},
		}}, available...)
	}

	return available
}

// GetAvailableAudioCodecs probes ffmpeg to discover which audio codecs are available.
func GetAvailableAudioCodecs(ffmpegPath string) []string {
	ffmpegEncoders := getFFmpegEncoders(ffmpegPath)
	var available []string
	for _, codec := range audioCodecRegistry {
		if _, ok := ffmpegEncoders[codec.FFmpegEncoderName()]; ok {
			available = append(available, codec.Name())
		}
	}
	// Always include AAC as it is built into ffmpeg.
	hasAAC := false
	for _, name := range available {
		if name == "aac" {
			hasAAC = true
			break
		}
	}
	if !hasAAC {
		available = append([]string{"aac"}, available...)
	}
	sort.Strings(available)
	return available
}

// getFFmpegEncoders runs `ffmpeg -encoders` and returns a set of encoder names.
func getFFmpegEncoders(ffmpegPath string) map[string]struct{} {
	result := make(map[string]struct{})

	cmd := exec.Command(ffmpegPath, "-encoders")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorln(err)
		return result
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '-' || line[0] == '=' {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			result[fields[1]] = struct{}{}
		}
	}

	return result
}
