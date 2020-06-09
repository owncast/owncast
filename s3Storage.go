package main

import (
	"bufio"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Storage struct {
	sess *session.Session
	host string

	s3Endpoint  string
	s3Region    string
	s3Bucket    string
	s3AccessKey string
	s3Secret    string
}

func (s *S3Storage) Setup(configuration Config) {
	log.Println("Setting up S3 for external storage of video...")

	s.s3Endpoint = configuration.S3.Endpoint
	s.s3Region = configuration.S3.Region
	s.s3Bucket = configuration.S3.Bucket
	s.s3AccessKey = configuration.S3.AccessKey
	s.s3Secret = configuration.S3.Secret

	s.sess = s.connectAWS()
}

func (s *S3Storage) Save(filePath string) string {
	// fmt.Println("Saving", filePath)

	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	uploader := s3manager.NewUploader(s.sess)

	response, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.s3Bucket), // Bucket to be used
		Key:    aws.String(filePath),   // Name of the file to be saved
		Body:   file,                   // File
	})

	if err != nil {
		panic(err)
	}

	// fmt.Println("Uploaded", filePath, "to", response.Location)

	return response.Location
}

func (s *S3Storage) GenerateRemotePlaylist(playlist string, variant Variant) string {
	var newPlaylist = ""

	scanner := bufio.NewScanner(strings.NewReader(playlist))
	for scanner.Scan() {
		line := scanner.Text()
		if line[0:1] != "#" {
			fullRemotePath := variant.getSegmentForFilename(line)
			if fullRemotePath != nil {
				line = fullRemotePath.RemoteID
			} else {
				line = ""
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
		panic(err)
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
		panic(err)
	}
	return sess
}
