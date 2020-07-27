package storageproviders

import (
	"bufio"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/models"
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
}

//Setup sets up the s3 storage for saving the video to s3
func (s *S3Storage) Setup() error {
	log.Trace("Setting up S3 for external storage of video...")

	s.s3Endpoint = config.Config.S3.Endpoint
	s.s3ServingEndpoint = config.Config.S3.ServingEndpoint
	s.s3Region = config.Config.S3.Region
	s.s3Bucket = config.Config.S3.Bucket
	s.s3AccessKey = config.Config.S3.AccessKey
	s.s3Secret = config.Config.S3.Secret

	s.sess = s.connectAWS()

	return nil
}

//Save saves the file to the s3 bucket
func (s *S3Storage) Save(filePath string, retryCount int) (string, error) {
	// fmt.Println("Saving", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	uploader := s3manager.NewUploader(s.sess)

	response, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.s3Bucket), // Bucket to be used
		Key:    aws.String(filePath),   // Name of the file to be saved
		Body:   file,                   // File
	})

	if err != nil {
		log.Trace("error uploading:", err.Error())
		if retryCount < 4 {
			log.Trace("Retrying...")
			return s.Save(filePath, retryCount+1)
		}
	}

	// fmt.Println("Uploaded", filePath, "to", response.Location)

	return response.Location, nil
}

//GenerateRemotePlaylist implements the 'GenerateRemotePlaylist' method
func (s *S3Storage) GenerateRemotePlaylist(playlist string, variant models.Variant) string {
	var newPlaylist = ""

	scanner := bufio.NewScanner(strings.NewReader(playlist))
	for scanner.Scan() {
		line := scanner.Text()
		if line[0:1] != "#" {
			fullRemotePath := variant.GetSegmentForFilename(line)
			if fullRemotePath == nil {
				line = ""
			} else if s.s3ServingEndpoint != "" {
				line = s.s3ServingEndpoint + "/" + fullRemotePath.RelativeUploadPath
			} else {
				line = fullRemotePath.RemoteID
			}
		}

		newPlaylist = newPlaylist + line + "\n"
	}

	return newPlaylist
}

func (s S3Storage) connectAWS() *session.Session {
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
