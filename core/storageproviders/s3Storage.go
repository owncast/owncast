package storageproviders

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/owncast/owncast/core/playlist"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/owncast/owncast/config"

	"github.com/grafov/m3u8"
)

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

// GenerateRemotePlaylist will take a local playlist and rewrite it to have absolute URLs to remote locations.
func (s *S3Storage) GenerateRemotePlaylist(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	p := m3u8.NewMasterPlaylist()
	err = p.DecodeFrom(bufio.NewReader(f), false)

	for _, item := range p.Variants {
		item.URI = filepath.Join(s.host, "hls", item.URI)
		fmt.Println(item.URI)
	}

	publicPath := filepath.Join(config.Config.GetPublicHLSSavePath(), filepath.Base(filePath))

	newPlaylist := p.String()

	return playlist.WritePlaylist(newPlaylist, publicPath)
}
