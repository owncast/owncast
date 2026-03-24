package transcoder

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/utils"
)

// FileWriterReceiverServiceCallback are to be fired when transcoder responses are written to disk.
type FileWriterReceiverServiceCallback interface {
	SegmentWritten(localFilePath string)
	VariantPlaylistWritten(localFilePath string)
	MasterPlaylistWritten(localFilePath string)
	InitSegmentWritten(localFilePath string)
}

// FileWriterReceiverService accepts transcoder responses via HTTP and fires the callbacks.
// It is intended to be the middleman between the transcoder and the storage provider and allows
// the transcoder process to be completely isolated and even run remotely in the future, as long
// as it can send HTTP requests to this service with the results.
type FileWriterReceiverService struct {
	// cfg receives the listener port discovered when the receiver binds
	// (the kernel-assigned ephemeral port). The transcoder later reads
	// the same field when it spins up its ffmpeg child.
	cfg       *config.Config
	callbacks FileWriterReceiverServiceCallback
}

// NewFileWriterReceiverService constructs an idle receiver. Call
// SetupFileWriterReceiverService to bind the listener.
func NewFileWriterReceiverService(cfg *config.Config) *FileWriterReceiverService {
	return &FileWriterReceiverService{cfg: cfg}
}

// SetupFileWriterReceiverService will start listening for transcoder responses.
func (s *FileWriterReceiverService) SetupFileWriterReceiverService(callbacks FileWriterReceiverServiceCallback) {
	s.callbacks = callbacks

	httpServer := http.NewServeMux()
	httpServer.HandleFunc("/", s.uploadHandler)

	host := s.cfg.InternalHLSListenerHost
	if host == "" {
		host = "127.0.0.1"
	}
	localListenerAddress := net.JoinHostPort(host, "0")

	listener, err := net.Listen("tcp", localListenerAddress)
	if err != nil {
		log.Fatalln("Unable to start internal video writing service", err)
	}

	_, listenerPort, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		log.Fatalln("Unable to determine internal video writing service port", err)
	}
	s.cfg.InternalHLSListenerPort = listenerPort
	log.Traceln("Transcoder response service listening on: " + listenerPort)
	go func() {
		//nolint: gosec
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

	relativePath, writePath, err := safeHLSUploadTarget(config.HLSStoragePath, r.URL.EscapedPath())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	root, err := os.OpenRoot(config.HLSStoragePath)
	if err != nil {
		returnError(err, w)
		return
	}
	defer root.Close()

	f, err := root.Create(relativePath)
	if err != nil {
		returnError(err, w)
		return
	}

	defer f.Close()

	if _, err := io.Copy(f, r.Body); err != nil {
		returnError(err, w)
		return
	}

	s.fileWritten(writePath)
	w.WriteHeader(http.StatusOK)
}

func (s *FileWriterReceiverService) fileWritten(path string) {
	if filepath.ToSlash(utils.GetRelativePathFromAbsolutePath(path)) == "hls/stream.m3u8" {
		s.callbacks.MasterPlaylistWritten(path)
	} else if strings.HasSuffix(path, ".m4s") {
		s.callbacks.SegmentWritten(path)
	} else if strings.HasSuffix(path, "init.mp4") {
		s.callbacks.InitSegmentWritten(path)
	} else if strings.HasSuffix(path, ".m3u8") {
		s.callbacks.VariantPlaylistWritten(path)
	}
}

func safeHLSUploadTarget(baseDir string, escapedRequestPath string) (string, string, error) {
	relativePath, err := safeHLSUploadRelativePath(escapedRequestPath)
	if err != nil {
		return "", "", err
	}

	return relativePath, filepath.Join(baseDir, filepath.FromSlash(relativePath)), nil
}

func safeHLSUploadRelativePath(escapedRequestPath string) (string, error) {
	requestPath, err := url.PathUnescape(escapedRequestPath)
	if err != nil {
		return "", errors.New("invalid HLS upload path")
	}

	relativePath, ok := strings.CutPrefix(requestPath, "/")
	if !ok || relativePath == "" || strings.HasPrefix(relativePath, "/") || strings.Contains(relativePath, `\`) {
		return "", errors.New("invalid HLS upload path")
	}

	for _, segment := range strings.Split(relativePath, "/") {
		if segment == "" || segment == "." || segment == ".." {
			return "", errors.New("invalid HLS upload path")
		}
	}

	cleanPath := path.Clean(relativePath)
	if !filepath.IsLocal(filepath.FromSlash(cleanPath)) {
		return "", errors.New("invalid HLS upload path")
	}

	return cleanPath, nil
}

func returnError(err error, w http.ResponseWriter) {
	log.Debugln(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError)+": "+err.Error(), http.StatusInternalServerError)
}
