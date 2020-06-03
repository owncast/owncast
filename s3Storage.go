package main

import (
	"bufio"
	"fmt"
	"net/url"
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

	s3Region    string
	s3Bucket    string
	s3AccessKey string
	s3Secret    string
}

func (s *S3Storage) Setup(configuration Config) {
	log.Println("Setting up S3 for external storage of video...")

	s.s3Region = configuration.S3.Region
	s.s3Bucket = configuration.S3.Bucket
	s.s3AccessKey = configuration.S3.AccessKey
	s.s3Secret = configuration.S3.Secret

	s.sess = s.connectAWS()
}

func (s *S3Storage) Save(filePath string) string {
	// fmt.Println("Saving", filePath)

	file, err := os.Open(filePath) // For read access.
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

	if s.host == "" {
		// Take note of the root host location so we can regenerate full
		// URLs to these files later when building the playlist in GenerateRemotePlaylist.
		url, err := url.Parse(response.Location)
		if err != nil {
			fmt.Println(err)
		}

		// The following is a bit of a hack to take the location URL string of that file
		// and get just the base URL without the file from it.
		pathComponents := strings.Split(url.Path, "/")
		pathComponents[len(pathComponents)-1] = ""
		pathString := strings.Join(pathComponents, "/")
		s.host = fmt.Sprintf("%s://%s%s", url.Scheme, url.Host, pathString)
	}

	// fmt.Println("Uploaded", filePath, "to", response.Location)

	return filePath
}

func (s *S3Storage) GenerateRemotePlaylist(playlist string, segments map[string]string) string {
	baseHost, err := url.Parse(s.host)
	baseHostComponents := []string{baseHost.Scheme + "://", baseHost.Host, baseHost.Path}

	verifyError(err)

	// baseHostString := fmt.Sprintf("%s://%s/%s", baseHost.Scheme, baseHost.Hostname, baseHost.Path)

	var newPlaylist = ""

	scanner := bufio.NewScanner(strings.NewReader(playlist))
	for scanner.Scan() {
		line := scanner.Text()
		if line[0:1] != "#" {
			urlComponents := baseHostComponents
			urlComponents = append(urlComponents, line)
			line = strings.Join(urlComponents, "") //path.Join(s.host, line)
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
			Region:      aws.String(s.s3Region),
			Credentials: creds,
		},
	)

	if err != nil {
		panic(err)
	}
	return sess
}
