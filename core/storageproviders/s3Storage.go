package storageproviders

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/owncast/owncast/config"
)

// S3Storage is the s3 implementation of a storage provider.
type S3Storage struct {
	sess     *session.Session
	s3Client *s3.S3
	host     string

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
	customVideoServingEndpoint := data.GetVideoServingEndpoint()

	if customVideoServingEndpoint != "" {
		s.host = customVideoServingEndpoint
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
	s.s3Client = s3.New(s.sess)

	s.uploader = s3manager.NewUploader(s.sess)

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
	if err := rewriteRemotePlaylist(localFilePath, s.host); err != nil {
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

	if path.Ext(filePath) == ".m3u8" {
		noCacheHeader := "no-cache, no-store, must-revalidate"
		contentType := "application/x-mpegURL"

		uploadInput.CacheControl = &noCacheHeader
		uploadInput.ContentType = &contentType
	}

	if s.s3ACL != "" {
		uploadInput.ACL = aws.String(s.s3ACL)
	} else {
		// Default ACL
		uploadInput.ACL = aws.String("public-read")
	}

	response, err := s.uploader.Upload(uploadInput)
	if err != nil {
		log.Traceln("error uploading segment", err.Error())
		if retryCount < 4 {
			log.Traceln("Retrying...")
			return s.Save(filePath, retryCount+1)
		}

		return "", fmt.Errorf("Giving up uploading %s to object storage %s", filePath, s.s3Endpoint)
	}

	return response.Location, nil
}

func (s *S3Storage) Cleanup() error {
	// Determine how many files we should keep on S3 storage
	maxNumber := data.GetStreamLatencyLevel().SegmentCount
	buffer := 20

	keys, err := s.getDeletableVideoSegmentsWithOffset(maxNumber + buffer)
	if err != nil {
		return err
	}

	s.deleteObjects(keys)

	return nil
}

func (s *S3Storage) connectAWS() *session.Session {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConnsPerHost = 100

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}

	creds := credentials.NewStaticCredentials(s.s3AccessKey, s.s3Secret, "")
	_, err := creds.Get()
	if err != nil {
		log.Panicln(err)
	}

	sess, err := session.NewSession(
		&aws.Config{
			HTTPClient:       httpClient,
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

func (s *S3Storage) getDeletableVideoSegmentsWithOffset(offset int) ([]s3object, error) {
	objectsToDelete, err := s.retrieveAllVideoSegments()
	if err != nil {
		return nil, err
	}

	objectsToDelete = objectsToDelete[offset : len(objectsToDelete)-1]

	return objectsToDelete, nil
}

func (s *S3Storage) deleteObjects(objects []s3object) {
	keys := make([]*s3.ObjectIdentifier, len(objects))
	for i, object := range objects {
		keys[i] = &s3.ObjectIdentifier{Key: aws.String(object.key)}
	}

	log.Debugln("Deleting", len(keys), "objects from S3 bucket:", s.s3Bucket)

	deleteObjectsRequest := &s3.DeleteObjectsInput{
		Bucket: aws.String(s.s3Bucket),
		Delete: &s3.Delete{
			Objects: keys,
			Quiet:   aws.Bool(true),
		},
	}

	_, err := s.s3Client.DeleteObjects(deleteObjectsRequest)
	if err != nil {
		log.Errorf("Unable to delete objects from bucket %q, %v\n", s.s3Bucket, err)
	}
}

func (s *S3Storage) retrieveAllVideoSegments() ([]s3object, error) {
	allObjectsListRequest := &s3.ListObjectsInput{
		Bucket: aws.String(s.s3Bucket),
	}

	// Fetch all objects in the bucket
	allObjectsListResponse, err := s.s3Client.ListObjects(allObjectsListRequest)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to fetch list of items in bucket for cleanup")
	}

	// Filter out non-video segments
	allObjects := []s3object{}
	for _, item := range allObjectsListResponse.Contents {
		if !strings.HasSuffix(*item.Key, ".ts") {
			continue
		}

		allObjects = append(allObjects, s3object{
			key:          *item.Key,
			lastModified: *item.LastModified,
		})
	}

	// Sort the results by timestamp
	sort.Slice(allObjects, func(i, j int) bool {
		return allObjects[i].lastModified.After(allObjects[j].lastModified)
	})

	return allObjects, nil
}

type s3object struct {
	key          string
	lastModified time.Time
}
