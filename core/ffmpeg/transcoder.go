package ffmpeg

import (
	"fmt"
	"os/exec"
	"path"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/utils"
)

// Transcoder is a single instance of a video transcoder
type Transcoder struct {
	Input              string
	SegmentOutputPath  string
	PlaylistOutputPath string
	Variants           []HLSVariant
	HLSPlaylistLength  int
}

// HLSVariant is a combination of settings that results in a single HLS stream
type HLSVariant struct {
	Index int

	VideoSize          VideoSize // Resizes the video via scaling
	Framerate          int       // The output framerate
	VideoBitrate       string    // The output bitrate
	IsVideoPassthrough bool      // Override all settings and just copy the video stream

	AudioBitrate       string // The audio bitrate
	IsAudioPassthrough bool   // Override all settings and just copy the audio stream

	EncoderPreset string // A collection of automatic settings for the encoder. https://trac.ffmpeg.org/wiki/Encode/H.264#crf
}

// VideoSize is the scaled size of the video output
type VideoSize struct {
	Width  int
	Height int
}

// String returns a WxH formatted string for scaling video output
func (v *VideoSize) String() string {
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

// Start will execute the transcoding process with the settings previously set.
func (t *Transcoder) Start() {
	command := t.getString()

	log.Printf("Video transcoder started with %d stream variants.", len(t.Variants))

	_, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		log.Panicln(err, command)
	}

	return
}

func (t *Transcoder) getString() string {
	hlsOptionFlags := []string{
		"delete_segments",
		"program_date_time",
		"temp_file",
	}
	ffmpegFlags := []string{
		"cat", t.Input, "|",
		config.Config.FFMpegPath,
		"-hide_banner",
		"-i pipe:",
		t.getVariantsString(),

		// HLS Output
		"-f", "hls",
		"-hls_time", strconv.Itoa(config.Config.VideoSettings.ChunkLengthInSeconds), // Length of each segment
		"-hls_list_size", strconv.Itoa(config.Config.Files.MaxNumberInPlaylist), // Max # in variant playlist
		"-hls_delete_threshold", "10", // Start deleting files after hls_list_size + 10
		"-hls_flags", strings.Join(hlsOptionFlags, "+"), // Specific options in HLS generation

		// Video settings
		"-tune", "zerolatency", // Option used for good for fast encoding and low-latency streaming (always includes iframes in each segment)
		// "-profile:v", "high", // Main – for standard definition (SD) to 640×480, High – for high definition (HD) to 1920×1080
		"-sc_threshold", "0", // Disable scene change detection for creating segments

		// Filenames
		"-master_pl_name", "stream.m3u8",
		"-strftime 1",                                                               // Support the use of strftime in filenames
		"-hls_segment_filename", path.Join(t.SegmentOutputPath, "/%v/stream-%s.ts"), // Each segment's filename
		"-max_muxing_queue_size", "400", // Workaround for Too many packets error: https://trac.ffmpeg.org/ticket/6375?cversion=0
		path.Join(t.SegmentOutputPath, "/%v/stream.m3u8"), // Each variant's playlist
	}

	return strings.Join(ffmpegFlags, " ")
}

func getVariantFromConfigQuality(quality config.StreamQuality, index int) HLSVariant {
	variant := HLSVariant{}
	variant.Index = index
	variant.IsAudioPassthrough = quality.IsAudioPassthrough
	variant.IsVideoPassthrough = quality.IsVideoPassthrough

	// If no audio bitrate is specified then we pass through original audio
	if quality.AudioBitrate == 0 {
		variant.IsAudioPassthrough = true
	}

	if quality.Bitrate == 0 {
		variant.IsVideoPassthrough = true
	}

	// If the video is being passed through then
	// don't continue to set options on the variant.
	if variant.IsVideoPassthrough {
		return variant
	}

	// Set a default, reasonable preset if one is not provided.
	// "superfast" and "ultrafast" are generally not recommended since they look bad.
	// https://trac.ffmpeg.org/wiki/Encode/H.264
	if quality.EncoderPreset != "" {
		variant.EncoderPreset = quality.EncoderPreset
	} else {
		variant.EncoderPreset = "veryfast"
	}

	variant.SetVideoBitrate(strconv.Itoa(quality.Bitrate) + "k")
	variant.SetAudioBitrate(strconv.Itoa(quality.AudioBitrate) + "k")
	variant.SetVideoScalingWidth(quality.ScaledWidth)
	variant.SetVideoScalingHeight(quality.ScaledHeight)
	variant.SetVideoFramerate(quality.Framerate)

	return variant
}

// NewTranscoder will return a new Transcoder, populated by the config
func NewTranscoder() Transcoder {
	transcoder := new(Transcoder)

	var outputPath string
	if config.Config.S3.Enabled || config.Config.IPFS.Enabled {
		// Segments are not available via the local HTTP server
		outputPath = config.Config.PrivateHLSPath
	} else {
		// Segments are available via the local HTTP server
		outputPath = config.Config.PublicHLSPath
	}

	transcoder.SegmentOutputPath = outputPath
	// Playlists are available via the local HTTP server
	transcoder.PlaylistOutputPath = config.Config.PublicHLSPath

	transcoder.Input = utils.GetTemporaryPipePath()

	for index, quality := range config.Config.VideoSettings.StreamQualities {
		variant := getVariantFromConfigQuality(quality, index)
		transcoder.AddVariant(variant)
	}

	return *transcoder
}

// Uses `map` https://www.ffmpeg.org/ffmpeg-all.html#Stream-specifiers-1 https://www.ffmpeg.org/ffmpeg-all.html#Advanced-options
func (v *HLSVariant) getVariantString() string {
	variantEncoderCommands := []string{
		v.getVideoQualityString(),
		v.getAudioQualityString(),
	}

	if v.VideoSize.Width != 0 || v.VideoSize.Height != 0 {
		variantEncoderCommands = append(variantEncoderCommands, v.getScalingString())
	}

	if v.Framerate != 0 {
		variantEncoderCommands = append(variantEncoderCommands, fmt.Sprintf("-r %d", v.Framerate))
		// multiply your output frame rate * 2. For example, if your input is -framerate 30, then use -g 60
		variantEncoderCommands = append(variantEncoderCommands, "-g "+strconv.Itoa(v.Framerate*2))
		variantEncoderCommands = append(variantEncoderCommands, "-keyint_min "+strconv.Itoa(v.Framerate*2))
	}

	if v.EncoderPreset != "" {
		variantEncoderCommands = append(variantEncoderCommands, fmt.Sprintf("-preset %s", v.EncoderPreset))
	}

	return strings.Join(variantEncoderCommands, " ")
}

// Get the command flags for the variants
func (t *Transcoder) getVariantsString() string {
	var variantsCommandFlags = ""
	var variantsStreamMaps = " -var_stream_map \""

	for _, variant := range t.Variants {
		variantsCommandFlags = variantsCommandFlags + " " + variant.getVariantString()
		variantsStreamMaps = variantsStreamMaps + fmt.Sprintf("v:%d,a:%d ", variant.Index, variant.Index)
	}
	variantsCommandFlags = variantsCommandFlags + " " + variantsStreamMaps + "\""

	return variantsCommandFlags
}

// Video Scaling
// https://trac.ffmpeg.org/wiki/Scaling
// If we'd like to keep the aspect ratio, we need to specify only one component, either width or height.
// Some codecs require the size of width and height to be a multiple of n. You can achieve this by setting the width or height to -n.

// SetVideoScalingWidth will set the scaled video width of this variant
func (v *HLSVariant) SetVideoScalingWidth(width int) {
	v.VideoSize.Width = width
}

// SetVideoScalingHeight will set the scaled video height of this variant
func (v *HLSVariant) SetVideoScalingHeight(height int) {
	v.VideoSize.Height = height
}

func (v *HLSVariant) getScalingString() string {
	scalingAlgorithm := "bilinear"
	return fmt.Sprintf("-filter:v:%d \"scale=%s\" -sws_flags %s", v.Index, v.VideoSize.String(), scalingAlgorithm)
}

// Video Quality

// SetVideoBitrate will set the output bitrate of this variant's video
func (v *HLSVariant) SetVideoBitrate(bitrate string) {
	v.VideoBitrate = bitrate
}

func (v *HLSVariant) getVideoQualityString() string {
	if v.IsVideoPassthrough {
		return fmt.Sprintf("-map v:0 -c:v:%d copy", v.Index)
	}

	encoderCodec := "libx264"
	return fmt.Sprintf("-map v:0 -c:v:%d %s -b:v:%d %s", v.Index, encoderCodec, v.Index, v.VideoBitrate)
}

// SetVideoFramerate will set the output framerate of this variant's video
func (v *HLSVariant) SetVideoFramerate(framerate int) {
	v.Framerate = framerate
}

// SetEncoderPreset will set the video encoder preset of this variant
func (v *HLSVariant) SetEncoderPreset(preset string) {
	v.EncoderPreset = preset
}

// Audio Quality

// SetAudioBitrate will set the output framerate of this variant's audio
func (v *HLSVariant) SetAudioBitrate(bitrate string) {
	v.AudioBitrate = bitrate
}

func (v *HLSVariant) getAudioQualityString() string {
	if v.IsAudioPassthrough {
		return fmt.Sprintf("-map a:0 -c:a:%d copy", v.Index)
	}

	encoderCodec := "libfdk_aac"
	return fmt.Sprintf("-map a:0 -c:a:%d %s -profile:a aac_he -b:a:%d %s", v.Index, encoderCodec, v.Index, v.AudioBitrate)
}

// AddVariant adds a new HLS variant to include in the output
func (t *Transcoder) AddVariant(variant HLSVariant) {
	t.Variants = append(t.Variants, variant)
}

// SetInput sets the input stream on the filesystem
func (t *Transcoder) SetInput(input string) {
	t.Input = input
}

// SetOutputPath sets the root directory that should include playlists and video segments
func (t *Transcoder) SetOutputPath(output string) {
	t.SegmentOutputPath = output
}

// SetHLSPlaylistLength will set the max number of items in a HLS variant's playlist
func (t *Transcoder) SetHLSPlaylistLength(length int) {
	t.HLSPlaylistLength = length
}
