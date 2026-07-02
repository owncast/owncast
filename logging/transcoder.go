package logging

import (
	"io"
	"sync"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

// Caps the transcoder log's disk usage at roughly maxSize * rotationCount.
const (
	transcoderLogMaxSizeBytes  = 10 * 1024 * 1024
	transcoderLogRotationCount = 5
)

var (
	transcoderWriterOnce sync.Once
	transcoderWriter     io.Writer
)

// GetTranscoderLogWriter returns the shared size-rotated writer for ffmpeg
// transcoder output.
func GetTranscoderLogWriter(logDirectory string) io.Writer {
	transcoderWriterOnce.Do(func() {
		path := GetTranscoderLogFilePath(logDirectory)
		writer, err := rotatelogs.New(
			path+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(path),
			rotatelogs.WithRotationSize(transcoderLogMaxSizeBytes),
			rotatelogs.WithRotationCount(transcoderLogRotationCount),
		)
		if err != nil {
			log.Errorln("unable to create transcoder log", path, err)
			transcoderWriter = io.Discard
			return
		}
		transcoderWriter = writer
	})
	return transcoderWriter
}
