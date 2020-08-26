package storageproviders

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/owncast/owncast/core/playlist"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
)

var variants []models.Variant

//S3Storage is the s3 implementation of the ChunkStorageProvider
type S3Storage struct {
	sess *session.Session
	host string

	s3Endpoint        string
	s3ServingEndpoint string
	s3Region          string
	s3Bucket          string
	s3AccessKey       string
	s3Secret          string
	s3ACL             string
}

var _uploader *s3manager.Uploader

//Setup sets up the s3 storage for saving the video to s3
func (s *S3Storage) Setup() error {
	log.Trace("Setting up S3 for external storage of video...")

	if config.Config.S3.ServingEndpoint != "" {
		s.host = config.Config.S3.ServingEndpoint
	} else {
		s.host = fmt.Sprintf("%s/%s", config.Config.S3.Endpoint, config.Config.S3.Bucket)
	}

	variants = make([]models.Variant, len(config.Config.GetVideoStreamQualities()))
	for index := range variants {
		variants[index] = models.Variant{
			VariantIndex: index,
			Segments:     make(map[string]*models.Segment),
		}
	}

	s.s3Endpoint = config.Config.S3.Endpoint
	s.s3ServingEndpoint = config.Config.S3.ServingEndpoint
	s.s3Region = config.Config.S3.Region
	s.s3Bucket = config.Config.S3.Bucket
	s.s3AccessKey = config.Config.S3.AccessKey
	s.s3Secret = config.Config.S3.Secret
	s.s3ACL = config.Config.S3.ACL

	s.sess = s.connectAWS()

	_uploader = s3manager.NewUploader(s.sess)

	return nil
}

//Save saves the file to the s3 bucket
func (s *S3Storage) Save(filePath string, retryCount int) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	uploadInput := &s3manager.UploadInput{
		Bucket: aws.String(s.s3Bucket), // Bucket to be used
		Key:    aws.String(filePath),   // Name of the file to be saved
		Body:   file,                   // File
	}
	if s.s3ACL != "" {
		uploadInput.ACL = aws.String(s.s3ACL)
	}
	response, err := _uploader.Upload(uploadInput)

	if err != nil {
		log.Trace("error uploading:", err.Error())
		if retryCount < 4 {
			log.Trace("Retrying...")
			return s.Save(filePath, retryCount+1)
		}
	}

	return response.Location, nil
}

func (s *S3Storage) connectAWS() *session.Session {
	creds := credentials.NewStaticCredentials(s.s3AccessKey, s.s3Secret, "")
	_, err := creds.Get()
	if err != nil {
		log.Panicln(err)
	}

	sess, err := session.NewSession(
		&aws.Config{
			Region:           aws.String(s.s3Region),
			Credentials:      creds,
			Endpoint:         aws.String(s.s3Endpoint),
			S3ForcePathStyle: aws.Bool(true),
		},
	)

	if err != nil {
		log.Panicln(err)
	}
	return sess
}

// SegmentWritten is a callback to be notified that a single video segment has been written
func (s *S3Storage) SegmentWritten(localFilePath string) {
	go func() {
		_, err := s.Save(localFilePath, 0)

		if err != nil {
			log.Errorln("failed to save the file to the chunk storage.", err)
		}
	}()
}

func (s *S3Storage) VariantPlaylistWritten(localFilePath string) {
}

func (s *S3Storage) MasterPlaylistWritten(localFilePath string) {
	err := utils.Move(localFilePath, path.Join(config.Config.GetPublicHLSSavePath(), filepath.Base(localFilePath)))
	if err != nil {
		log.Errorln(err)
	}
}

// func (s *S3Storage) WriteVariantPlaylist(fullPath string) error {
func (s *S3Storage) GenerateRemotePlaylist(filePath string) error {
	playlistBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	playlistString := string(playlistBytes)
	scanner := bufio.NewScanner(strings.NewReader(playlistString))

	remoteHost := s.host
	if s.s3ServingEndpoint != "" {
		remoteHost = s.s3ServingEndpoint
	}

	newPlaylist := ""

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" && line[0:1] != "#" {
			line = fmt.Sprintf("%s/hls/%s", remoteHost, line)
		}

		newPlaylist = newPlaylist + line + "\n"
	}
	publicPath := filepath.Join(config.Config.GetPublicHLSSavePath(), filepath.Base(filePath))

	return playlist.WritePlaylist(newPlaylist, publicPath)
}
