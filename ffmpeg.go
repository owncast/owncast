package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
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
	for index, quality := range configuration.VideoSettings.StreamQualities {
		videoMaps = append(videoMaps, fmt.Sprintf("-map v:0 -c:v:%d libx264 -b:v:%d %s", index, index, quality.Bitrate))
		streamMaps = append(streamMaps, fmt.Sprintf("v:%d,a:%d", index, index))
		audioMaps = append(audioMaps, "-map a:0")
	}

	ffmpegFlags := []string{
		"-hide_banner",
		"-re",
		"-i pipe:",
		// "-vf scale=900:-2", // Re-enable in the future with a config to togging resizing?
		// "-sws_flags fast_bilinear",
		strings.Join(videoMaps, " "),                // All the different video variants
		strings.Join(audioMaps, " ") + " -c:a copy", // Audio for all the variants
		// strings.Join(audioMaps, " ") + " -c:a aac -b:a 192k -ac 2", // Audio for all the variants
		"-master_pl_name stream.m3u8",
		"-g 60", "-keyint_min 60", // create key frame (I-frame) every 48 frames (~2 seconds) - will later affect correct slicing of segments and alignment of renditions
		"-framerate 30",
		"-r 15",
		"-preset " + configuration.VideoSettings.EncoderPreset,
		"-sc_threshold 0", // don't create key frames on scene change - only according to -g
		"-profile:v main", // Main – for standard definition (SD) to 640×480, High – for high definition (HD) to 1920×1080
		"-f hls",
		"-hls_list_size 30",
		"-hls_time 10",
		"-strftime 1",
		"-use_localtime 1",
		"-hls_playlist_type event",
		"-hls_segment_filename " + path.Join(outputDir, "stream-%Y%m%d-%s.ts"),
		"-hls_flags delete_segments+program_date_time+temp_file",
		"-segment_wrap 100",
		"-tune zerolatency",

		// "-master_m3u8_publish_rate 5",
		"-var_stream_map \"" + strings.Join(streamMaps, " ") + "\"",
		variantPlaylistName,
	}

	ffmpegFlagsString := strings.Join(ffmpegFlags, " ")

	ffmpegCmd := "cat " + pipePath + " | " + configuration.FFMpegPath + " " + ffmpegFlagsString

	// fmt.Println(ffmpegCmd)

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
