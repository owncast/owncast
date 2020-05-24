package main

import (
	"fmt"
	"os/exec"

	"github.com/r3b-fish/goffmpeg/transcoder"
)

func main() {

	// Create new instance of transcoder
	trans := new(transcoder.Transcoder)

	// Initialize an empty transcoder
	err := trans.InitializeEmptyTranscoder()
	if err != nil {
		panic(err)
	}

	trans.MediaFile().SetHlsListSize(10)
	trans.MediaFile().SetHlsSegmentDuration(10)
	trans.MediaFile().SetPreset("ultrafast")
	trans.MediaFile().SetOutputPath("./test.m3u8")
	trans.MediaFile().SetOutputFormat("hls")
	trans.MediaFile().SetHlsSegmentFilename("stream%05d.ts")
	// Create a command such that its output should be passed as stdin to ffmpeg
	cmd := exec.Command("cat", "../rfBd56ti2SMtYvSgD5xAV0YU99zampta7Z7S575KLkIZ9PYk.flv")

	// Create an input pipe to write to, which will return *io.PipeWriter
	w, err := trans.CreateInputPipe()

	cmd.Stdout = w

	// Create an output pipe to read from, which will return *io.PipeReader.
	// Must also specify the output container format
	// r, err := trans.CreateOutputPipe("hls")

	// wg := &sync.WaitGroup{}
	// wg.Add(1)
	// go func() {
	// 	defer r.Close()
	// 	defer wg.Done()

	// 	// Read data from output pipe
	// 	data, err := ioutil.ReadAll(r)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	// err = ioutil.WriteFile("test.tmp", data, 0644)
	// 	// fmt.Println(len(data))
	// }()

	go func() {
		fmt.Println("1")

		defer w.Close()
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
		fmt.Println("2")
		// Handle error...
	}()

	// Returns a channel to get the transcoding progress

	progress := trans.Output()
	// Example of printing transcoding progress
	for msg := range progress {
		fmt.Println(msg)
	}

	done := trans.Run(true)

	// This channel is used to wait for the transcoding process to end
	err = <-done
	if err != nil {
		panic(err)
	}

	// wg.Wait()
}
