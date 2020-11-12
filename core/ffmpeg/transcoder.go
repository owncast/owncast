package ffmpeg

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/utils"
)

var _commandExec *exec.Cmd

// Transcoder is a single instance of a video transcoder.
type Transcoder struct {
	input                string
	segmentOutputPath    string
	playlistOutputPath   string
	variants             []HLSVariant
	hlsPlaylistLength    int
	segmentLengthSeconds int
	appendToStream       bool
	ffmpegPath           string
	segmentIdentifier    string
	internalListenerPort int
	TranscoderCompleted  func(error)
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

	encoderPreset string // A collection of automatic settings for the encoder. https://trac.ffmpeg.org/wiki/Encode/H.264#crf
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
	command := t.getString()

	log.Tracef("Video transcoder started with %d stream variants.", len(t.variants))

	if config.Config.EnableDebugFeatures {
		log.Println(command)
	}

	_commandExec = exec.Command("sh", "-c", command)
	err := _commandExec.Start()
	if err != nil {
		log.Errorln("Transcoder error.  See transcoder.log for full output to debug.")
		log.Panicln(err, command)
	}

	err = _commandExec.Wait()
	if t.TranscoderCompleted != nil {
		t.TranscoderCompleted(err)
	}
}

func (t *Transcoder) getString() string {
	var port int
	if config.Config != nil {
		port = config.Config.GetPublicWebServerPort() + 1
	} else if t.internalListenerPort != 0 {
		port = t.internalListenerPort
	} else {
		log.Panicln("A internal port must be set for transcoder callback")
	}

	localListenerAddress := "http://127.0.0.1:" + strconv.Itoa(port)

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
		t.ffmpegPath,
		"-hide_banner",
		"-loglevel warning",
		"-i ", t.input,

		t.getVariantsString(),

		// HLS Output
		"-f", "hls",

		"-hls_time", strconv.Itoa(t.segmentLengthSeconds), // Length of each segment
		"-hls_list_size", strconv.Itoa(t.hlsPlaylistLength), // Max # in variant playlist
		"-hls_delete_threshold", "10", // Start deleting files after hls_list_size + 10
		hlsOptionsString,

		// Video settings
		"-tune", "zerolatency", // Option used for good for fast encoding and low-latency streaming (always includes iframes in each segment)
		"-pix_fmt", "yuv420p", // Force yuv420p color format
		"-profile:v", "high", // Main – for standard definition (SD) to 640×480, High – for high definition (HD) to 1920×1080
		"-sc_threshold", "0", // Disable scene change detection for creating segments

		// Filenames
		"-master_pl_name", "stream.m3u8",
		"-strftime 1", // Support the use of strftime in filenames

		"-hls_segment_filename", localListenerAddress + "/%v/stream-" + t.segmentIdentifier + "%s.ts", // Send HLS segments back to us over HTTP
		"-max_muxing_queue_size", "400", // Workaround for Too many packets error: https://trac.ffmpeg.org/ticket/6375?cversion=0

		"-method PUT -http_persistent 0",         // HLS results sent back to us will be over PUTs
		"-fflags +genpts",                        // Generate presentation time stamp if missing
		localListenerAddress + "/%v/stream.m3u8", // Send HLS playlists back to us over HTTP
		"2> transcoder.log",                      // Log to a file for debugging
	}

	return strings.Join(ffmpegFlags, " ")
}

func getVariantFromConfigQuality(quality config.StreamQuality, index int) HLSVariant {
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
	variant.encoderPreset = quality.GetEncoderPreset()

	variant.SetVideoBitrate(quality.VideoBitrate)
	variant.SetAudioBitrate(strconv.Itoa(quality.AudioBitrate) + "k")
	variant.SetVideoScalingWidth(quality.ScaledWidth)
	variant.SetVideoScalingHeight(quality.ScaledHeight)
	variant.SetVideoFramerate(quality.GetFramerate())

	return variant
}

// NewTranscoder will return a new Transcoder, populated by the config.
func NewTranscoder() *Transcoder {
	transcoder := new(Transcoder)
	transcoder.ffmpegPath = config.Config.GetFFMpegPath()
	transcoder.hlsPlaylistLength = config.Config.GetMaxNumberOfReferencedSegmentsInPlaylist()

	var outputPath string
	if config.Config.S3.Enabled {
		// Segments are not available via the local HTTP server
		outputPath = config.PrivateHLSStoragePath
	} else {
		// Segments are available via the local HTTP server
		outputPath = config.PublicHLSStoragePath
	}

	transcoder.segmentOutputPath = outputPath
	// Playlists are available via the local HTTP server
	transcoder.playlistOutputPath = config.PublicHLSStoragePath

	transcoder.input = utils.GetTemporaryPipePath()
	transcoder.segmentLengthSeconds = config.Config.GetVideoSegmentSecondsLength()

	qualities := config.Config.GetVideoStreamQualities()
	for index, quality := range qualities {
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

	if v.videoSize.Width != 0 || v.videoSize.Height != 0 {
		variantEncoderCommands = append(variantEncoderCommands, v.getScalingString())
	}

	variantEncoderCommands = append(variantEncoderCommands, fmt.Sprintf("-r %d", v.framerate))

	if v.encoderPreset != "" {
		variantEncoderCommands = append(variantEncoderCommands, fmt.Sprintf("-preset %s", v.encoderPreset))
	}

	return strings.Join(variantEncoderCommands, " ")
}

// Get the command flags for the variants.
func (t *Transcoder) getVariantsString() string {
	var variantsCommandFlags = ""
	var variantsStreamMaps = " -var_stream_map \""

	for _, variant := range t.variants {
		variantsCommandFlags = variantsCommandFlags + " " + variant.getVariantString(t)
		variantsStreamMaps = variantsStreamMaps + fmt.Sprintf("v:%d,a:%d ", variant.index, variant.index)
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
	scalingAlgorithm := "bilinear"
	return fmt.Sprintf("-filter:v:%d \"scale=%s\" -sws_flags %s", v.index, v.videoSize.getString(), scalingAlgorithm)
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

	encoderCodec := "libx264"

	// -1 to work around segments being generated slightly larger than expected.
	// https://trac.ffmpeg.org/ticket/6915?replyto=58#comment:57
	gop := (t.segmentLengthSeconds * v.framerate) - 1

	// For limiting the output bitrate
	// https://trac.ffmpeg.org/wiki/Limiting%20the%20output%20bitrate
	// https://developer.apple.com/documentation/http_live_streaming/about_apple_s_http_live_streaming_tools
	// Adjust the max & buffer size until the output bitrate doesn't exceed the ~+10% that Apple's media validator
	// complains about.
	maxBitrate := int(float64(v.videoBitrate) * 1.06) // Max is a ~+10% over specified bitrate.
	bufferSize := int(float64(v.videoBitrate) * 1.2)  // How often it checks the bitrate of encoded segments to see if it's too high/low.

	cmd := []string{
		"-map v:0",
		fmt.Sprintf("-c:v:%d %s", v.index, encoderCodec),                                                      // Video codec used for this variant
		fmt.Sprintf("-b:v:%d %dk", v.index, v.videoBitrate),                                                   // The average bitrate for this variant
		fmt.Sprintf("-maxrate:v:%d %dk", v.index, maxBitrate),                                                 // The max bitrate allowed for this variant
		fmt.Sprintf("-bufsize:v:%d %dk", v.index, bufferSize),                                                 // How often the encoder checks the bitrate in order to meet average/max values
		fmt.Sprintf("-g:v:%d %d", v.index, gop),                                                               // How often i-frames are encoded into the segments
		fmt.Sprintf("-x264-params:v:%d \"scenecut=0:open_gop=0:min-keyint=%d:keyint=%d\"", v.index, gop, gop), // How often i-frames are encoded into the segments
	}

	return strings.Join(cmd, " ")
}

// SetVideoFramerate will set the output framerate of this variant's video.
func (v *HLSVariant) SetVideoFramerate(framerate int) {
	v.framerate = framerate
}

// SetEncoderPreset will set the video encoder preset of this variant.
func (v *HLSVariant) SetEncoderPreset(preset string) {
	v.encoderPreset = preset
}

// Audio Quality

// SetAudioBitrate will set the output framerate of this variant's audio.
func (v *HLSVariant) SetAudioBitrate(bitrate string) {
	v.audioBitrate = bitrate
}

func (v *HLSVariant) getAudioQualityString() string {
	if v.isAudioPassthrough {
		return fmt.Sprintf("-map a:0 -c:a:%d copy", v.index)
	}

	// libfdk_aac is not a part of every ffmpeg install, so use "aac" instead
	encoderCodec := "aac"
	return fmt.Sprintf("-map a:0 -c:a:%d %s -b:a:%d %s", v.index, encoderCodec, v.index, v.audioBitrate)
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

// SetOutputPath sets the root directory that should include playlists and video segments.
func (t *Transcoder) SetOutputPath(output string) {
	t.segmentOutputPath = output
}

// SetHLSPlaylistLength will set the max number of items in a HLS variant's playlist.
func (t *Transcoder) SetHLSPlaylistLength(length int) {
	t.hlsPlaylistLength = length
}

// SetSegmentLength Specifies the number of seconds each segment should be.
func (t *Transcoder) SetSegmentLength(seconds int) {
	t.segmentLengthSeconds = seconds
}

// SetAppendToStream enables appending to the HLS stream instead of overwriting.
func (t *Transcoder) SetAppendToStream(append bool) {
	t.appendToStream = append
}

// SetIdentifer enables appending a unique identifier to segment file name.
func (t *Transcoder) SetIdentifier(output string) {
	t.segmentIdentifier = output
}

func (t *Transcoder) SetInternalHTTPPort(port int) {
	t.internalListenerPort = port
}
