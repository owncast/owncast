package storageproviders

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
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
	// fmt.Println("Saving", filePath)

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

	// fmt.Println("Uploaded", filePath, "to", response.Location)

	return response.Location, nil
}

//GenerateRemotePlaylist implements the 'GenerateRemotePlaylist' method
func (s *S3Storage) GenerateRemotePlaylist(playlist string, variant models.Variant) string {
	var newPlaylist = ""

	scanner := bufio.NewScanner(strings.NewReader(playlist))
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" && line[0:1] != "#" {
			fullRemotePath := variant.GetSegmentForFilename(line)
			if fullRemotePath == nil {
				line = ""
			} else if s.s3ServingEndpoint != "" {
				line = fmt.Sprintf("%s/%s/%s", s.s3ServingEndpoint, config.PrivateHLSStoragePath, fullRemotePath.RelativeUploadPath)
			} else {
				line = fullRemotePath.RemoteID
			}
		}

		newPlaylist = newPlaylist + line + "\n"
	}

	return newPlaylist
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

func (s *S3Storage) SegmentWritten(localFilePath string) {
	newObjectPathChannel := make(chan string, 1)
	go func() {
		newObjectPath, err := s.Save(localFilePath, 0)

		if err != nil {
			log.Errorln("failed to save the file to the chunk storage.", err)
		}

		newObjectPathChannel <- newObjectPath
	}()

	newObjectPath := <-newObjectPathChannel

	segment := models.Segment{
		FullDiskPath:       localFilePath,
		RelativeUploadPath: utils.GetRelativePathFromAbsolutePath(localFilePath),
		RemoteID:           newObjectPath,
	}
	fmt.Println("Uploaded", segment.RelativeUploadPath, "as", newObjectPath)

	variants[segment.VariantIndex].Segments[filepath.Base(segment.RelativeUploadPath)] = &segment

	associatedVariantPlaylist := strings.ReplaceAll(localFilePath, path.Base(localFilePath), "stream.m3u8")
	s.writeVariantPlaylist(associatedVariantPlaylist)
}

func (s *S3Storage) VariantPlaylistWritten(localFilePath string) {
}

func (s *S3Storage) MasterPlaylistWritten(localFilePath string) {
	err := utils.Copy(localFilePath, path.Join(config.Config.GetPublicHLSSavePath(), "stream.m3u8"))
	if err != nil {
		log.Errorln(err)
	}
}

func (s *S3Storage) writeVariantPlaylist(fullPath string) error {
	relativePath := utils.GetRelativePathFromAbsolutePath(fullPath)
	variantIndex, err := strconv.Atoi(utils.GetIndexFromFilePath(relativePath))
	if err != nil {
		log.Errorln(err)
	}
	variant := variants[variantIndex]

	playlistBytes, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}

	playlistString := string(playlistBytes)
	playlistString = s.GenerateRemotePlaylist(playlistString, variant)

	fmt.Println("Wrote", relativePath)
	return playlist.WritePlaylist(playlistString, path.Join(config.Config.GetPublicHLSSavePath(), relativePath))
}
