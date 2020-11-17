package ffmpeg

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strconv"
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

	localListenerAddress := "127.0.0.1:" + strconv.Itoa(config.Config.GetPublicWebServerPort()+1)

	go func() {
		if err := http.ListenAndServe(localListenerAddress, httpServer); err != nil {
			log.Fatal(err)
		}
	}()

	log.Traceln("Transcoder response listening on: " + localListenerAddress)
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

var _inWarningState = false

func (s *FileWriterReceiverService) fileWritten(path string) {
	index := utils.GetIndexFromFilePath(path)

	if utils.GetRelativePathFromAbsolutePath(path) == "hls/stream.m3u8" {
		s.callbacks.MasterPlaylistWritten(path)
	} else if strings.HasSuffix(path, ".ts") {
		performanceMonitorKey := "segmentWritten-" + index
		averagePerformance := utils.GetAveragePerformance(performanceMonitorKey)

		utils.StartPerformanceMonitor(performanceMonitorKey)
		s.callbacks.SegmentWritten(path)

		if averagePerformance != 0 && averagePerformance > float64(config.Config.GetVideoSegmentSecondsLength())*1.1 {
			if !_inWarningState {
				log.Warnln("slow encoding for variant", index, "if this continues you may see buffering or errors. troubleshoot this issue by visiting https://owncast.online/docs/troubleshooting/")
				_inWarningState = true
			}
		} else {
			_inWarningState = false
		}
	} else if strings.HasSuffix(path, ".m3u8") {
		s.callbacks.VariantPlaylistWritten(path)
	}
}

func returnError(err error, w http.ResponseWriter) {
	log.Errorln(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError)+": "+err.Error(), http.StatusInternalServerError)
}
