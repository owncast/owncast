package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

func startFfmpeg(configuration Config) {
	var outputDir = configuration.PublicHLSPath
	var hlsPlaylistName = path.Join(configuration.PublicHLSPath, "stream.m3u8")

	if configuration.IPFS.Enabled {
		outputDir = configuration.PrivateHLSPath
		hlsPlaylistName = path.Join(outputDir, "temp.m3u8")
	}

	log.Printf("Starting transcoder saving to /%s.", outputDir)

	ffmpegCmd := "cat streampipe.flv | " + configuration.FFMpegPath +
		" -hide_banner -i pipe: -vf scale=" + strconv.Itoa(configuration.VideoSettings.ResolutionWidth) + ":-2 -g 48 -keyint_min 48 -preset ultrafast -f hls -hls_list_size 30 -hls_time " +
		strconv.Itoa(configuration.VideoSettings.ChunkLengthInSeconds) + " -strftime 1 -use_localtime 1 -hls_segment_filename '" +
		outputDir + "/stream-%Y%m%d-%s.ts' -hls_flags delete_segments -segment_wrap 100 " + hlsPlaylistName

	_, err := exec.Command("bash", "-c", ffmpegCmd).Output()
	verifyError(err)
}

func verifyError(e error) {
	if e != nil {
		panic(e)
	}
}

func generateRemotePlaylist(playlist string, gateway string, segments map[string]string) string {
	for local, remote := range segments {
		playlist = strings.ReplaceAll(playlist, local, gateway+remote)
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
