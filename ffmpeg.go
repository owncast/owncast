package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

func startFfmpeg(configuration Config) {
	var outputDir = configuration.PublicHLSPath
	var variantPlaylistPath = configuration.PublicHLSPath

	if configuration.IPFS.Enabled || configuration.S3.Enabled {
		outputDir = configuration.PrivateHLSPath
		variantPlaylistPath = configuration.PrivateHLSPath
	}

	outputDir = path.Join(outputDir, "%v")

	// var masterPlaylistName = path.Join(configuration.PublicHLSPath, "%v", "stream.m3u8")
	var variantPlaylistName = path.Join(variantPlaylistPath, "%v", "stream.m3u8")
	// var variantRootPath = configuration.PublicHLSPath

	// variantRootPath = path.Join(variantRootPath, "%v")
	// variantPlaylistName := path.Join("%v", "stream.m3u8")

	log.Printf("Starting transcoder saving to /%s.", variantPlaylistName)
	pipePath := getTempPipePath()

	var videoMaps = make([]string, 0)
	var streamMaps = make([]string, 0)
	var audioMaps = make([]string, 0)
	var videoMapsString = ""
	var audioMapsString = ""
	var streamMappingString = ""
	if configuration.VideoSettings.EnablePassthrough || len(configuration.VideoSettings.StreamQualities) == 0 {
		fmt.Println("Enabling passthrough video")
		streamMaps = append(streamMaps, fmt.Sprintf("v:%d,a:%d", 0, 0))
	} else {
		for index, quality := range configuration.VideoSettings.StreamQualities {
			maxRate := math.Floor(float64(quality.Bitrate) * 0.8)
			videoMaps = append(videoMaps, fmt.Sprintf("-map v:0 -c:v:%d libx264 -b:v:%d %dk -maxrate %dk -bufsize %dk", index, index, int(quality.Bitrate), int(maxRate), int(maxRate)))
			streamMaps = append(streamMaps, fmt.Sprintf("v:%d,a:%d", index, index))
			videoMapsString = strings.Join(videoMaps, " ")
			audioMaps = append(audioMaps, "-map a:0")
			audioMapsString = strings.Join(audioMaps, " ") + " -c:a copy" // Pass through audio for all the variants, don't reencode
		}
	}

	framerate := 25

	streamMappingString = "-var_stream_map \"" + strings.Join(streamMaps, " ") + "\""
	ffmpegFlags := []string{
		"-hide_banner",
		"-re",
		"-i pipe:",
		// "-vf scale=900:-2", // Re-enable in the future with a config to togging resizing?
		// "-sws_flags fast_bilinear",
		videoMapsString, // All the different video variants
		audioMapsString,
		"-master_pl_name stream.m3u8",
		"-framerate " + strconv.Itoa(framerate),
		"-g " + strconv.Itoa(framerate*2), " -keyint_min " + strconv.Itoa(framerate*2), // multiply your output frame rate * 2. For example, if your input is -framerate 30, then use -g 60
		// "-r 25",
		"-preset " + configuration.VideoSettings.EncoderPreset,
		"-sc_threshold 0", // don't create key frames on scene change - only according to -g
		"-profile:v high", // Main – for standard definition (SD) to 640×480, High – for high definition (HD) to 1920×1080
		"-movflags +faststart",
		"-pix_fmt yuv420p",
		"-f hls",
		"-hls_list_size " + strconv.Itoa(configuration.Files.MaxNumberInPlaylist),
		"-hls_time " + strconv.Itoa(configuration.VideoSettings.ChunkLengthInSeconds),
		"-strftime 1",
		"-use_localtime 1",
		"-hls_playlist_type event",
		"-hls_segment_filename " + path.Join(outputDir, "stream-%Y%m%d-%s.ts"),
		"-hls_flags delete_segments+program_date_time+temp_file",
		"-tune zerolatency",

		streamMappingString,
		variantPlaylistName,
	}

	ffmpegFlagsString := strings.Join(ffmpegFlags, " ")

	ffmpegCmd := "cat " + pipePath + " | " + configuration.FFMpegPath + " " + ffmpegFlagsString

	fmt.Println(ffmpegCmd)

	_, err := exec.Command("sh", "-c", ffmpegCmd).Output()
	fmt.Println(err)
	verifyError(err)
}

func writePlaylist(data string, filePath string) {
	f, err := os.Create(filePath)
	defer f.Close()

	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.WriteString(data)
	if err != nil {
		fmt.Println(err)
		return
	}
}
