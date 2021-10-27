package storageproviders

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/playlist"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/owncast/owncast/config"

	"github.com/grafov/m3u8"
)

// S3Storage is the s3 implementation of a storage provider.
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
	s3ForcePathStyle  bool

	// If we try to upload a playlist but it is not yet on disk
	// then keep a reference to it here.
	queuedPlaylistUpdates map[string]string

	uploader *s3manager.Uploader
	s3Client *s3.S3
}

// NewS3Storage returns a new S3Storage instance.
func NewS3Storage() *S3Storage {
	return &S3Storage{
		queuedPlaylistUpdates: make(map[string]string),
	}
}

// Setup sets up the s3 storage for saving the video to s3.
func (s *S3Storage) Setup() error {
	log.Trace("Setting up S3 for external storage of video...")

	s3Config := data.GetS3Config()
	if s3Config.ServingEndpoint != "" {
		s.host = s3Config.ServingEndpoint
	} else {
		s.host = fmt.Sprintf("%s/%s", s3Config.Endpoint, s3Config.Bucket)
	}

	s.s3Endpoint = s3Config.Endpoint
	s.s3ServingEndpoint = s3Config.ServingEndpoint
	s.s3Region = s3Config.Region
	s.s3Bucket = s3Config.Bucket
	s.s3AccessKey = s3Config.AccessKey
	s.s3Secret = s3Config.Secret
	s.s3ACL = s3Config.ACL
	s.s3ForcePathStyle = s3Config.ForcePathStyle

	s.sess = s.connectAWS()

	s.uploader = s3manager.NewUploader(s.sess)
	s.s3Client = s3.New(s.sess)

	return nil
}

// SegmentWritten is called when a single segment of video is written.
func (s *S3Storage) SegmentWritten(localFilePath string) {
	index := utils.GetIndexFromFilePath(localFilePath)
	performanceMonitorKey := "s3upload-" + index
	utils.StartPerformanceMonitor(performanceMonitorKey)

	// Upload the segment
	if _, err := s.Save(localFilePath, 0); err != nil {
		log.Errorln(err)
		return
	}
	averagePerformance := utils.GetAveragePerformance(performanceMonitorKey)

	// Warn the user about long-running save operations
	if averagePerformance != 0 {
		if averagePerformance > float64(data.GetStreamLatencyLevel().SecondsPerSegment)*0.9 {
			log.Warnln("Possible slow uploads: average upload S3 save duration", averagePerformance, "s. troubleshoot this issue by visiting https://owncast.online/docs/troubleshooting/")
		}
	}

	// Upload the variant playlist for this segment
	// so the segments and the HLS playlist referencing
	// them are in sync.
	playlistPath := filepath.Join(filepath.Dir(localFilePath), "stream.m3u8")
	if _, err := s.Save(playlistPath, 0); err != nil {
		s.queuedPlaylistUpdates[playlistPath] = playlistPath
		if pErr, ok := err.(*os.PathError); ok {
			log.Debugln(pErr.Path, "does not yet exist locally when trying to upload to S3 storage.")
			return
		}
	}
}

// VariantPlaylistWritten is called when a variant hls playlist is written.
func (s *S3Storage) VariantPlaylistWritten(localFilePath string) {
	// We are uploading the variant playlist after uploading the segment
	// to make sure we're not referring to files in a playlist that don't
	// yet exist.  See SegmentWritten.
	if _, ok := s.queuedPlaylistUpdates[localFilePath]; ok {
		if _, err := s.Save(localFilePath, 0); err != nil {
			log.Errorln(err)
			s.queuedPlaylistUpdates[localFilePath] = localFilePath
		}
		delete(s.queuedPlaylistUpdates, localFilePath)
	}
}

// MasterPlaylistWritten is called when the master hls playlist is written.
func (s *S3Storage) MasterPlaylistWritten(localFilePath string) {
	// Rewrite the playlist to use absolute remote S3 URLs
	if err := s.rewriteRemotePlaylist(localFilePath); err != nil {
		log.Warnln(err)
	}
}

// Save saves the file to the s3 bucket.
func (s *S3Storage) Save(filePath string, retryCount int) (string, error) {
	file, err := os.Open(filePath) // nolint
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Convert the local path to the variant/file path by stripping the local storage location.
	normalizedPath := strings.TrimPrefix(filePath, config.HLSStoragePath)
	// Build the remote path by adding the "hls" path prefix.
	remotePath := strings.Join([]string{"hls", normalizedPath}, "")

	maxAgeSeconds := utils.GetCacheDurationSecondsForPath(filePath)
	cacheControlHeader := fmt.Sprintf("max-age=%d", maxAgeSeconds)
	uploadInput := &s3manager.UploadInput{
		Bucket:       aws.String(s.s3Bucket), // Bucket to be used
		Key:          aws.String(remotePath), // Name of the file to be saved
		Body:         file,                   // File
		CacheControl: &cacheControlHeader,
	}

	if s.s3ACL != "" {
		uploadInput.ACL = aws.String(s.s3ACL)
	} else {
		// Default ACL
		uploadInput.ACL = aws.String("public-read")
	}

	response, err := s.uploader.Upload(uploadInput)

	if err != nil {
		log.Traceln("error uploading:", filePath, err.Error())
		if retryCount < 4 {
			log.Traceln("Retrying...")
			return s.Save(filePath, retryCount+1)
		}

		log.Warnln("Giving up on", filePath, err)
		return "", fmt.Errorf("Giving up on %s", filePath)
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
			S3ForcePathStyle: aws.Bool(s.s3ForcePathStyle),
		},
	)

	if err != nil {
		log.Panicln(err)
	}
	return sess
}

// CleanupOldContent will delete old files from the S3 Bucket that are no longer being referenced
// in the stream.
func (s *S3Storage) CleanupOldContent(baseDirectory string) {
	keys := s.retrieveKeysToDeleteButLeftNewest(50)
	s.deleteObjectsWithGiven(keys)
}

func (s *S3Storage) retrieveKeysToDeleteButLeftNewest(offset int) []*s3.ObjectIdentifier {

	allObjects := s.retrieveAllObjects()
	tsObjects := s.filterTsObjectsFrom(allObjects)
	s.sortObjectsByModificationDateDescending(tsObjects)
	objectsToDelete := s.ommitNewest(tsObjects, offset)

	return s.retrieveKeysFrom(objectsToDelete)
}

func (s *S3Storage) deleteObjectsWithGiven(keys []*s3.ObjectIdentifier) {
	deleteObjectsRequest := &s3.DeleteObjectsInput{
		Bucket: aws.String(s.s3Bucket),
		Delete: &s3.Delete{
			Objects: keys,
			Quiet:   aws.Bool(false),
		},
	}

	deleteObjectsResponse, err := s.s3Client.DeleteObjects(deleteObjectsRequest)

	if err != nil {
		log.Debugf("Unable to delete objects from bucket %q, %v\n", s.s3Bucket, err)
	}

	log.Traceln(deleteObjectsResponse)
}

func (s *S3Storage) retrieveAllObjects() []*s3.Object {

	allObjectsListRequest := &s3.ListObjectsInput{
		Bucket: aws.String(s.s3Bucket),
	}

	allObjectsListResponse, err := s.s3Client.ListObjects(allObjectsListRequest)
	if err != nil {
		log.Debugf("Unable to list items in bucket %q, %v", s.s3Bucket, err)
		return nil
	}

	return allObjectsListResponse.Contents
}

func (s *S3Storage) filterTsObjectsFrom(objects []*s3.Object) []*s3.Object {

	var result []*s3.Object

	for _, item := range objects {
		if strings.HasSuffix(*item.Key, ".ts") {
			result = append(result, item)
		}
	}

	return result
}

func (s *S3Storage) sortObjectsByModificationDateDescending(objects []*s3.Object) {
	sort.Slice(objects, func(i, j int) bool {
		return objects[i].LastModified.After(*objects[j].LastModified)
	})
}

func (s *S3Storage) ommitNewest(objects []*s3.Object, offset int) []*s3.Object {
	return objects[offset:]
}

func (s *S3Storage) retrieveKeysFrom(objects []*s3.Object) []*s3.ObjectIdentifier {
	var result []*s3.ObjectIdentifier

	for _, item := range objects {
		result = append(result, &s3.ObjectIdentifier{
			Key: item.Key,
		})
	}

	return result
}

// rewriteRemotePlaylist will take a local playlist and rewrite it to have absolute URLs to remote locations.
func (s *S3Storage) rewriteRemotePlaylist(filePath string) error {
	f, err := os.Open(filePath) // nolint
	if err != nil {
		log.Fatalln(err)
	}

	p := m3u8.NewMasterPlaylist()
	if err := p.DecodeFrom(bufio.NewReader(f), false); err != nil {
		log.Warnln(err)
	}

	for _, item := range p.Variants {
		item.URI = s.host + filepath.Join("/hls", item.URI)
	}

	publicPath := filepath.Join(config.HLSStoragePath, filepath.Base(filePath))

	newPlaylist := p.String()

	return playlist.WritePlaylist(newPlaylist, publicPath)
}
