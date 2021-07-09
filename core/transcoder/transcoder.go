package transcoder

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/logging"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
)

var _commandExec *exec.Cmd

// Transcoder is a single instance of a video transcoder.
type Transcoder struct {
	input                string
	stdin                *io.PipeReader
	segmentOutputPath    string
	playlistOutputPath   string
	variants             []HLSVariant
	appendToStream       bool
	ffmpegPath           string
	segmentIdentifier    string
	internalListenerPort string
	codec                Codec

	currentStreamOutputSettings []models.StreamOutputVariant
	currentLatencyLevel         models.LatencyLevel

	TranscoderCompleted func(error)
}

// HLSVariant is a combination of settings that results in a single HLS stream.
type HLSVariant struct {
	index int

	videoSize          VideoSize // Resizes the video via scaling
	framerate          int       // The output framerate
	videoBitrate       int       // The output bitrate
	isVideoPassthrough bool      // Override all settings and just copy the video stream

	audioBitrate       string // The audio bitrate
	isAudioPassthrough bool   // Override all settings and just copy the audio stream

	cpuUsageLevel int // The amount of hardware to use for encoding a stream
}

// VideoSize is the scaled size of the video output.
type VideoSize struct {
	Width  int
	Height int
}

// getString returns a WxH formatted getString for scaling video output.
func (v *VideoSize) getString() string {
	widthString := strconv.Itoa(v.Width)
	heightString := strconv.Itoa(v.Height)

	if widthString != "0" && heightString != "0" {
		return widthString + ":" + heightString
	} else if widthString != "0" {
		return widthString + ":-2"
	} else if heightString != "0" {
		return "-2:" + heightString
	}

	return ""
}

func (t *Transcoder) Stop() {
	log.Traceln("Transcoder STOP requested.")
	err := _commandExec.Process.Kill()
	if err != nil {
		log.Errorln(err)
	}
}

// Start will execute the transcoding process with the settings previously set.
func (t *Transcoder) Start() {
	_lastTranscoderLogMessage = ""

	command := t.getString()
	log.Infof("Video transcoder started using %s with %d stream variants.", t.codec.DisplayName(), len(t.variants))
	createVariantDirectories()

	if config.EnableDebugFeatures {
		log.Println(command)
	}

	_commandExec = exec.Command("sh", "-c", command)

	if t.stdin != nil {
		_commandExec.Stdin = t.stdin
	}

	stdout, err := _commandExec.StderrPipe()
	if err != nil {
		panic(err)
	}

	if err := _commandExec.Start(); err != nil {
		log.Errorln("Transcoder error.  See ", logging.GetTranscoderLogFilePath(), " for full output to debug.")
		log.Panicln(err, command)
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			handleTranscoderMessage(line)
		}
	}()

	err = _commandExec.Wait()
	if t.TranscoderCompleted != nil {
		t.TranscoderCompleted(err)
	}

	if err != nil {
		log.Errorln("transcoding error. look at ", logging.GetTranscoderLogFilePath(), " to help debug. your copy of ffmpeg may not support your selected codec of", t.codec.Name(), "https://owncast.online/docs/troubleshooting/#codecs")
	}
}

func (t *Transcoder) getString() string {
	var port = t.internalListenerPort
	localListenerAddress := "http://127.0.0.1:" + port

	hlsOptionFlags := []string{}

	if t.appendToStream {
		hlsOptionFlags = append(hlsOptionFlags, "append_list")
	}

	if t.segmentIdentifier == "" {
		t.segmentIdentifier = shortid.MustGenerate()
	}

	hlsOptionsString := ""
	if len(hlsOptionFlags) > 0 {
		hlsOptionsString = "-hls_flags " + strings.Join(hlsOptionFlags, "+")
	}
	ffmpegFlags := []string{
		fmt.Sprintf(`FFREPORT=file="%s":level=32`, logging.GetTranscoderLogFilePath()),
		t.ffmpegPath,
		"-hide_banner",
		"-loglevel warning",
		t.codec.GlobalFlags(),
		"-fflags +genpts", // Generate presentation time stamp if missing
		"-i ", t.input,

		t.getVariantsString(),

		// HLS Output
		"-f", "hls",

		"-hls_time", strconv.Itoa(t.currentLatencyLevel.SecondsPerSegment), // Length of each segment
		"-hls_list_size", strconv.Itoa(t.currentLatencyLevel.SegmentCount), // Max # in variant playlist
		hlsOptionsString,
		"-segment_format_options", "mpegts_flags=+initial_discontinuity:mpegts_copyts=1",

		// Video settings
		t.codec.ExtraArguments(),
		"-pix_fmt", t.codec.PixelFormat(),
		"-sc_threshold", "0", // Disable scene change detection for creating segments

		// Filenames
		"-master_pl_name", "stream.m3u8",
		"-strftime 1", // Support the use of strftime in filenames

		"-hls_segment_filename", localListenerAddress + "/%v/stream-" + t.segmentIdentifier + "%s.ts", // Send HLS segments back to us over HTTP
		"-max_muxing_queue_size", "400", // Workaround for Too many packets error: https://trac.ffmpeg.org/ticket/6375?cversion=0

		"-method PUT -http_persistent 0",         // HLS results sent back to us will be over PUTs
		localListenerAddress + "/%v/stream.m3u8", // Send HLS playlists back to us over HTTP
	}

	return strings.Join(ffmpegFlags, " ")
}

func getVariantFromConfigQuality(quality models.StreamOutputVariant, index int) HLSVariant {
	variant := HLSVariant{}
	variant.index = index
	variant.isAudioPassthrough = quality.IsAudioPassthrough
	variant.isVideoPassthrough = quality.IsVideoPassthrough

	// If no audio bitrate is specified then we pass through original audio
	if quality.AudioBitrate == 0 {
		variant.isAudioPassthrough = true
	}

	if quality.VideoBitrate == 0 {
		quality.VideoBitrate = 1200
	}

	// If the video is being passed through then
	// don't continue to set options on the variant.
	if variant.isVideoPassthrough {
		return variant
	}

	// Set a default, reasonable preset if one is not provided.
	// "superfast" and "ultrafast" are generally not recommended since they look bad.
	// https://trac.ffmpeg.org/wiki/Encode/H.264
	variant.cpuUsageLevel = quality.CPUUsageLevel

	variant.SetVideoBitrate(quality.VideoBitrate)
	variant.SetAudioBitrate(strconv.Itoa(quality.AudioBitrate) + "k")
	variant.SetVideoScalingWidth(quality.ScaledWidth)
	variant.SetVideoScalingHeight(quality.ScaledHeight)
	variant.SetVideoFramerate(quality.GetFramerate())

	return variant
}

// NewTranscoder will return a new Transcoder, populated by the config.
func NewTranscoder() *Transcoder {
	ffmpegPath := utils.ValidatedFfmpegPath(data.GetFfMpegPath())

	transcoder := new(Transcoder)
	transcoder.ffmpegPath = ffmpegPath
	transcoder.internalListenerPort = config.InternalHLSListenerPort

	transcoder.currentStreamOutputSettings = data.GetStreamOutputVariants()
	transcoder.currentLatencyLevel = data.GetStreamLatencyLevel()
	transcoder.codec = getCodec(data.GetVideoCodec())

	var outputPath string
	if data.GetS3Config().Enabled {
		// Segments are not available via the local HTTP server
		outputPath = config.PrivateHLSStoragePath
	} else {
		// Segments are available via the local HTTP server
		outputPath = config.PublicHLSStoragePath
	}

	transcoder.segmentOutputPath = outputPath
	// Playlists are available via the local HTTP server
	transcoder.playlistOutputPath = config.PublicHLSStoragePath

	transcoder.input = "pipe:0" // stdin

	for index, quality := range transcoder.currentStreamOutputSettings {
		variant := getVariantFromConfigQuality(quality, index)
		transcoder.AddVariant(variant)
	}

	return transcoder
}

// Uses `map` https://www.ffmpeg.org/ffmpeg-all.html#Stream-specifiers-1 https://www.ffmpeg.org/ffmpeg-all.html#Advanced-options
func (v *HLSVariant) getVariantString(t *Transcoder) string {
	variantEncoderCommands := []string{
		v.getVideoQualityString(t),
		v.getAudioQualityString(),
	}

	if (v.videoSize.Width != 0 || v.videoSize.Height != 0) && !v.isVideoPassthrough {
		// Order here matters, you must scale before changing hardware formats
		filters := []string{
			v.getScalingString(),
		}
		if t.codec.ExtraFilters() != "" {
			filters = append(filters, t.codec.ExtraFilters())
		}
		scalingAlgorithm := "bilinear"
		filterString := fmt.Sprintf("-sws_flags %s -filter:v:%d \"%s\"", scalingAlgorithm, v.index, strings.Join(filters, ","))
		variantEncoderCommands = append(variantEncoderCommands, filterString)
	} else if t.codec.ExtraFilters() != "" && !v.isVideoPassthrough {
		filterString := fmt.Sprintf("-filter:v:%d \"%s\"", v.index, t.codec.ExtraFilters())
		variantEncoderCommands = append(variantEncoderCommands, filterString)
	}

	preset := t.codec.GetPresetForLevel(v.cpuUsageLevel)
	if preset != "" {
		variantEncoderCommands = append(variantEncoderCommands, fmt.Sprintf("-preset %s", preset))
	}

	return strings.Join(variantEncoderCommands, " ")
}

// Get the command flags for the variants.
func (t *Transcoder) getVariantsString() string {
	var variantsCommandFlags = ""
	var variantsStreamMaps = " -var_stream_map \""

	for _, variant := range t.variants {
		variantsCommandFlags = variantsCommandFlags + " " + variant.getVariantString(t)
		singleVariantMap := ""
		singleVariantMap = fmt.Sprintf("v:%d,a:%d ", variant.index, variant.index)
		variantsStreamMaps = variantsStreamMaps + singleVariantMap
	}
	variantsCommandFlags = variantsCommandFlags + " " + variantsStreamMaps + "\""

	return variantsCommandFlags
}

// Video Scaling
// https://trac.ffmpeg.org/wiki/Scaling
// If we'd like to keep the aspect ratio, we need to specify only one component, either width or height.
// Some codecs require the size of width and height to be a multiple of n. You can achieve this by setting the width or height to -n.

// SetVideoScalingWidth will set the scaled video width of this variant.
func (v *HLSVariant) SetVideoScalingWidth(width int) {
	v.videoSize.Width = width
}

// SetVideoScalingHeight will set the scaled video height of this variant.
func (v *HLSVariant) SetVideoScalingHeight(height int) {
	v.videoSize.Height = height
}

func (v *HLSVariant) getScalingString() string {
	return fmt.Sprintf("scale=%s", v.videoSize.getString())
}

// Video Quality

// SetVideoBitrate will set the output bitrate of this variant's video.
func (v *HLSVariant) SetVideoBitrate(bitrate int) {
	v.videoBitrate = bitrate
}

func (v *HLSVariant) getVideoQualityString(t *Transcoder) string {
	if v.isVideoPassthrough {
		return fmt.Sprintf("-map v:0 -c:v:%d copy", v.index)
	}

	gop := v.framerate * t.currentLatencyLevel.SecondsPerSegment // force an i-frame every segment

	// For limiting the output bitrate
	// https://trac.ffmpeg.org/wiki/Limiting%20the%20output%20bitrate
	// https://developer.apple.com/documentation/http_live_streaming/about_apple_s_http_live_streaming_tools
	// Adjust the max & buffer size until the output bitrate doesn't exceed the ~+10% that Apple's media validator
	// complains about.
	maxBitrate := int(float64(v.videoBitrate) * 1.06) // Max is a ~+10% over specified bitrate.

	cmd := []string{
		"-map v:0",
		fmt.Sprintf("-c:v:%d %s", v.index, t.codec.Name()),    // Video codec used for this variant
		fmt.Sprintf("-b:v:%d %dk", v.index, v.videoBitrate),   // The average bitrate for this variant
		fmt.Sprintf("-maxrate:v:%d %dk", v.index, maxBitrate), // The max bitrate allowed for this variant
		fmt.Sprintf("-g:v:%d %d", v.index, gop),               // Suggested interval where i-frames are encoded into the segments
		fmt.Sprintf("-keyint_min:v:%d %d", v.index, gop),      // minimum i-keyframe interval
		fmt.Sprintf("-r:v:%d %d", v.index, v.framerate),
		t.codec.VariantFlags(v),
	}

	return strings.Join(cmd, " ")
}

// SetVideoFramerate will set the output framerate of this variant's video.
func (v *HLSVariant) SetVideoFramerate(framerate int) {
	v.framerate = framerate
}

// SetCPUUsageLevel will set the hardware usage of this variant.
func (v *HLSVariant) SetCPUUsageLevel(level int) {
	v.cpuUsageLevel = level
}

// Audio Quality

// SetAudioBitrate will set the output framerate of this variant's audio.
func (v *HLSVariant) SetAudioBitrate(bitrate string) {
	v.audioBitrate = bitrate
}

func (v *HLSVariant) getAudioQualityString() string {
	if v.isAudioPassthrough {
		return fmt.Sprintf("-map a:0? -c:a:%d copy", v.index)
	}

	// libfdk_aac is not a part of every ffmpeg install, so use "aac" instead
	encoderCodec := "aac"
	return fmt.Sprintf("-map a:0? -c:a:%d %s -b:a:%d %s", v.index, encoderCodec, v.index, v.audioBitrate)
}

// AddVariant adds a new HLS variant to include in the output.
func (t *Transcoder) AddVariant(variant HLSVariant) {
	variant.index = len(t.variants)
	t.variants = append(t.variants, variant)
}

// SetInput sets the input stream on the filesystem.
func (t *Transcoder) SetInput(input string) {
	t.input = input
}

// SetStdin sets the Stdin of the ffmpeg command.
func (t *Transcoder) SetStdin(rtmp *io.PipeReader) {
	t.stdin = rtmp
}

// SetOutputPath sets the root directory that should include playlists and video segments.
func (t *Transcoder) SetOutputPath(output string) {
	t.segmentOutputPath = output
}

// SetAppendToStream enables appending to the HLS stream instead of overwriting.
func (t *Transcoder) SetAppendToStream(append bool) {
	t.appendToStream = append
}

// SetIdentifer enables appending a unique identifier to segment file name.
func (t *Transcoder) SetIdentifier(output string) {
	t.segmentIdentifier = output
}

func (t *Transcoder) SetInternalHTTPPort(port string) {
	t.internalListenerPort = port
}

func (t *Transcoder) SetCodec(codecName string) {
	t.codec = getCodec(codecName)
}
