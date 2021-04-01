package transcoder

import (
	"bytes"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"net/http"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

// FileWriterReceiverServiceCallback are to be fired when transcoder responses are written to disk.
type FileWriterReceiverServiceCallback interface {
	SegmentWritten(localFilePath string)
	VariantPlaylistWritten(localFilePath string)
	MasterPlaylistWritten(localFilePath string)
}

// FileWriterReceiverService accepts transcoder responses via HTTP and fires the callbacks.
type FileWriterReceiverService struct {
	callbacks FileWriterReceiverServiceCallback
}

// SetupFileWriterReceiverService will start listening for transcoder responses.
func (s *FileWriterReceiverService) SetupFileWriterReceiverService(callbacks FileWriterReceiverServiceCallback) {
	s.callbacks = callbacks

	httpServer := http.NewServeMux()
	httpServer.HandleFunc("/", s.uploadHandler)

	localListenerAddress := "127.0.0.1:0"

	// go func() {
	listener, err := net.Listen("tcp", localListenerAddress)
	if err != nil {
		log.Fatalln("Unable to start internal video writing service", err)
	}

	listenerPort := strings.Split(listener.Addr().String(), ":")[1]
	config.InternalHLSListenerPort = listenerPort
	log.Traceln("Transcoder response service listening on: " + listenerPort)
	go func() {
		if err := http.Serve(listener, httpServer); err != nil {
			log.Fatalln("Unable to start internal video writing service", err)
		}
	}()
}

func (s *FileWriterReceiverService) uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	path := r.URL.Path
	writePath := filepath.Join(config.PrivateHLSStoragePath, path)

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r.Body)
	data := buf.Bytes()

	f, err := os.Create(writePath)
	if err != nil {
		returnError(err, w)
		return
	}

	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		returnError(err, w)
		return
	}

	s.fileWritten(writePath)
	w.WriteHeader(http.StatusOK)
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

func returnError(err error, w http.ResponseWriter) {
	log.Debugln(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError)+": "+err.Error(), http.StatusInternalServerError)
}
