package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func startFfmpeg() {
	outputDir := "webroot"
	chunkLength := "4"

	log.Printf("Starting transcoder with segments saving to %s.", outputDir)

	// ffmpegCmd := "cat streampipe.flv | ffmpeg -hide_banner -i pipe: -vf scale=900:-2 -g 48 -keyint_min 48 -preset ultrafast -f hls -hls_list_size 30 -hls_time 10 -strftime 1 -use_localtime 1 -hls_segment_filename 'hls/stream-%Y%m%d-%s.ts' -hls_flags delete_segments -segment_wrap 100 hls/temp.m3u8"

	ffmpegCmd := "cat streampipe.flv | ffmpeg -hide_banner -i pipe: -vf scale=900:-2 -g 48 -keyint_min 48 -preset ultrafast -f hls -hls_list_size 30 -hls_time " + chunkLength + " -strftime 1 -use_localtime 1 -hls_segment_filename '" + outputDir + "/stream-%Y%m%d-%s.ts' -hls_flags delete_segments -segment_wrap 100 hls/temp.m3u8"
	fmt.Println(ffmpegCmd)

	_, err := exec.Command("bash", "-c", ffmpegCmd).Output()
	verifyError(err)
}

func verifyError(e error) {
	if e != nil {
		panic(e)
	}
}

func generateRemotePlaylist(playlist string, segments map[string]string) string {
	for local, remote := range segments {
		playlist = strings.ReplaceAll(playlist, local, "https://gateway.temporal.cloud"+remote)
	}
	return playlist
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
