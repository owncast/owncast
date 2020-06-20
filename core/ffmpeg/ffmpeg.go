package ffmpeg

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/utils"
)

//ShowStreamOfflineState generates and shows the stream's offline state
func ShowStreamOfflineState() error {
	fmt.Println("----- Stream offline!  Showing offline state!")

	var outputDir = config.Config.PublicHLSPath
	var variantPlaylistPath = config.Config.PublicHLSPath

	if config.Config.IPFS.Enabled || config.Config.S3.Enabled {
		outputDir = config.Config.PrivateHLSPath
		variantPlaylistPath = config.Config.PrivateHLSPath
	}

	outputDir = path.Join(outputDir, "%v")
	var variantPlaylistName = path.Join(variantPlaylistPath, "%v", "stream.m3u8")

	var videoMaps = make([]string, 0)
	var streamMaps = make([]string, 0)
	var videoMapsString = ""
	var streamMappingString = ""
	if config.Config.VideoSettings.EnablePassthrough || len(config.Config.VideoSettings.StreamQualities) == 0 {
		fmt.Println("Enabling passthrough video")
		videoMapsString = "-b:v 1200k -b:a 128k" // Since we're compositing multiple sources we can't infer bitrate, so pick something reasonable.
		streamMaps = append(streamMaps, fmt.Sprintf("v:%d", 0))
	} else {
		for index, quality := range config.Config.VideoSettings.StreamQualities {
			maxRate := math.Floor(float64(quality.Bitrate) * 0.8)
			videoMaps = append(videoMaps, fmt.Sprintf("-map v:0 -c:v:%d libx264 -b:v:%d %dk -maxrate %dk -bufsize %dk", index, index, int(quality.Bitrate), int(maxRate), int(maxRate)))
			streamMaps = append(streamMaps, fmt.Sprintf("v:%d", index))
			videoMapsString = strings.Join(videoMaps, " ")
		}
	}

	framerate := 25

	streamMappingString = "-var_stream_map \"" + strings.Join(streamMaps, " ") + "\""

	ffmpegFlags := []string{
		"-hide_banner",
		// "-stream_loop 100",
		// "-fflags", "+genpts",
		"-i", config.Config.VideoSettings.OfflineImage,
		"-i", "webroot/thumbnail.jpg",
		"-filter_complex", "\"[0:v]scale=2640:2360[bg];[bg][1:v]overlay=200:250:enable='between(t,0,3)'\"",
		videoMapsString, // All the different video variants
		"-f hls",
		// "-hls_list_size " + strconv.Itoa(config.Config.Files.MaxNumberInPlaylist),
		"-hls_time 4", // + strconv.Itoa(config.Config.VideoSettings.ChunkLengthInSeconds),
		"-hls_playlist_type", "event",
		"-master_pl_name", "stream.m3u8",
		"-strftime 1",
		"-use_localtime 1",
		"-hls_flags temp_file",
		"-tune", "zerolatency",
		"-g " + strconv.Itoa(framerate*2), " -keyint_min " + strconv.Itoa(framerate*2), // multiply your output frame rate * 2. For example, if your input is -framerate 30, then use -g 60
		"-framerate " + strconv.Itoa(framerate),
		"-preset " + config.Config.VideoSettings.EncoderPreset,
		"-sc_threshold 0",    // don't create key frames on scene change - only according to -g
		"-profile:v", "main", // Main – for standard definition (SD) to 640×480, High – for high definition (HD) to 1920×1080
		// "-movflags +faststart",
		"-pix_fmt yuv420p",

		streamMappingString,
		"-hls_segment_filename " + path.Join(outputDir, "offline-%s.ts"),
		// "-s", "720x480", // size
		variantPlaylistName,
	}

	ffmpegFlagsString := strings.Join(ffmpegFlags, " ")

	ffmpegCmd := config.Config.FFMpegPath + " " + ffmpegFlagsString

	// fmt.Println(ffmpegCmd)

	_, err := exec.Command("sh", "-c", ffmpegCmd).Output()

	return err
}

//Start starts the ffmpeg process
func Start() error {
	var outputDir = config.Config.PublicHLSPath
	var variantPlaylistPath = config.Config.PublicHLSPath

	if config.Config.IPFS.Enabled || config.Config.S3.Enabled {
		outputDir = config.Config.PrivateHLSPath
		variantPlaylistPath = config.Config.PrivateHLSPath
	}

	outputDir = path.Join(outputDir, "%v")
	var variantPlaylistName = path.Join(variantPlaylistPath, "%v", "stream.m3u8")

	log.Printf("Starting transcoder saving to /%s.", variantPlaylistName)
	pipePath := utils.GetTemporaryPipePath()

	var videoMaps = make([]string, 0)
	var streamMaps = make([]string, 0)
	var audioMaps = make([]string, 0)
	var videoMapsString = ""
	var audioMapsString = ""
	var streamMappingString = ""
	var profileString = ""

	if config.Config.VideoSettings.EnablePassthrough || len(config.Config.VideoSettings.StreamQualities) == 0 {
		fmt.Println("Enabling passthrough video")
		streamMaps = append(streamMaps, fmt.Sprintf("v:%d,a:%d", 0, 0))
		videoMaps = append(videoMaps, "-map v:0 -c:v copy")
		videoMapsString = strings.Join(videoMaps, " ")
		audioMaps = append(audioMaps, "-map a:0")
		audioMapsString = strings.Join(audioMaps, " ") + " -c:a copy" // Pass through audio for all the variants, don't reencode

	} else {
		for index, quality := range config.Config.VideoSettings.StreamQualities {
			maxRate := math.Floor(float64(quality.Bitrate) * 0.8)
			videoMaps = append(videoMaps, fmt.Sprintf("-map v:0 -c:v:%d libx264 -b:v:%d %dk -maxrate %dk -bufsize %dk", index, index, int(quality.Bitrate), int(maxRate), int(maxRate)))
			streamMaps = append(streamMaps, fmt.Sprintf("v:%d,a:%d", index, index))
			videoMapsString = strings.Join(videoMaps, " ")
			audioMaps = append(audioMaps, "-map a:0")
			audioMapsString = strings.Join(audioMaps, " ") + " -c:a copy" // Pass through audio for all the variants, don't reencode
			profileString = "-profile:v high"                             // Main – for standard definition (SD) to 640×480, High – for high definition (HD) to 1920×1080
		}
	}

	framerate := 25

	streamMappingString = "-var_stream_map \"" + strings.Join(streamMaps, " ") + "\""
	ffmpegFlags := []string{
		"-hide_banner",
		// "-re",
		"-fflags", "+genpts",
		"-i pipe:",
		// "-vf scale=900:-2", // Re-enable in the future with a config to togging resizing?
		// "-sws_flags fast_bilinear",
		videoMapsString, // All the different video variants
		audioMapsString,
		"-master_pl_name stream.m3u8",
		"-framerate " + strconv.Itoa(framerate),
		"-g " + strconv.Itoa(framerate*2), " -keyint_min " + strconv.Itoa(framerate*2), // multiply your output frame rate * 2. For example, if your input is -framerate 30, then use -g 60
		// "-r 25",
		"-preset " + config.Config.VideoSettings.EncoderPreset,
		"-sc_threshold 0", // don't create key frames on scene change - only according to -g
		profileString,
		"-movflags +faststart",
		"-pix_fmt yuv420p",
		"-f hls",
		"-hls_list_size " + strconv.Itoa(config.Config.Files.MaxNumberInPlaylist),
		"-hls_delete_threshold 10", // Keep 10 unreferenced segments on disk before they're deleted.
		"-hls_time " + strconv.Itoa(config.Config.VideoSettings.ChunkLengthInSeconds),
		"-strftime 1",
		"-use_localtime 1",
		"-hls_playlist_type event",
		"-hls_segment_filename " + path.Join(outputDir, "stream-%Y%m%d-%s.ts"),
		"-hls_flags delete_segments+program_date_time+temp_file",
		"-tune zerolatency",
		// "-s", "720x480", // size

		streamMappingString,
		variantPlaylistName,
	}

	ffmpegFlagsString := strings.Join(ffmpegFlags, " ")

	ffmpegCmd := "cat " + pipePath + " | " + config.Config.FFMpegPath + " " + ffmpegFlagsString

	// fmt.Println(ffmpegCmd)

	_, err := exec.Command("sh", "-c", ffmpegCmd).Output()

	return err
}

//WritePlaylist writes the playlist to disk
func WritePlaylist(data string, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(data); err != nil {
		return err
	}

	return nil
}
