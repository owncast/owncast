package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func pipeTest() {
	ffmpegCmd := "cat streampipe.flv | ffmpeg -hide_banner -i pipe: -preset ultrafast -f hls -hls_list_size 10 -hls_time 10 -strftime 1 -hls_segment_filename 'hls/stream-%Y%m%d-%s.ts' -hls_flags delete_segments -segment_wrap 100 hls/temp.m3u8"

	out, err := exec.Command("bash", "-c", ffmpegCmd).Output()
	verifyError(err)
	fmt.Println(string(out))
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
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.WriteString(data)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
}
