package ffmpeg

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"net/http"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/utils"
)

const localListenerAddress = "127.0.0.1:8089"

type FileWriterReceiverServiceCallback interface {
	SegmentWritten(localFilePath string)
	VariantPlaylistWritten(localFilePath string)
	MasterPlaylistWritten(localFilePath string)
}

type FileWriterReceiverService struct {
	callbacks FileWriterReceiverServiceCallback
}

func (s *FileWriterReceiverService) SetupFileWriterReceiverService(callbacks FileWriterReceiverServiceCallback) {
	s.callbacks = callbacks

	httpServer := http.NewServeMux()
	httpServer.HandleFunc("/", s.uploadHandler)

	go http.ListenAndServe(localListenerAddress, httpServer)
}

// By returning a handler, we have an elegant way of initializing path.
func (s *FileWriterReceiverService) uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	path := r.URL.Path
	writePath := filepath.Join(config.Config.GetPrivateHLSSavePath(), path)
	out, err := os.Create(writePath)

	defer out.Close()
	if err != nil {
		panic(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": "+err.Error(), http.StatusInternalServerError)
		return
	}

	io.Copy(out, r.Body)
	s.fileWritten(writePath)
}

func (s *FileWriterReceiverService) fileWritten(path string) {
	if utils.GetRelativePathFromAbsolutePath(path) == "hls/stream.m3u8" {
		s.callbacks.MasterPlaylistWritten(path)

	} else if strings.HasSuffix(path, ".ts") {
		s.callbacks.SegmentWritten(path)

	} else if strings.HasSuffix(path, ".m3u8") {
		s.callbacks.VariantPlaylistWritten(path)
	}
}
