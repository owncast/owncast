package ffmpeg

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"net/http"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/utils"
	log "github.com/sirupsen/logrus"
)

// FileWriterReceiverServiceCallback are to be fired when transcoder responses are written to disk
type FileWriterReceiverServiceCallback interface {
	SegmentWritten(localFilePath string)
	VariantPlaylistWritten(localFilePath string)
	MasterPlaylistWritten(localFilePath string)
}

// FileWriterReceiverService accepts transcoder responses via HTTP and fires the callbacks
type FileWriterReceiverService struct {
	callbacks FileWriterReceiverServiceCallback
}

// SetupFileWriterReceiverService will start listening for transcoder responses
func (s *FileWriterReceiverService) SetupFileWriterReceiverService(callbacks FileWriterReceiverServiceCallback) {
	s.callbacks = callbacks

	httpServer := http.NewServeMux()
	httpServer.HandleFunc("/", s.uploadHandler)

	localListenerAddress := "127.0.0.1:" + strconv.Itoa(config.Config.GetPublicWebServerPort()+1)
	go http.ListenAndServe(localListenerAddress, httpServer)
	log.Debugln("Transcoder response listening on: " + localListenerAddress)
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
		log.Errorln(err)
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
