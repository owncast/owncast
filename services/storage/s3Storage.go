package storage

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/owncast/owncast/config"
)

// S3Storage is the s3 implementation of a storage provider.
type S3Storage struct {
	// If we try to upload a playlist but it is not yet on disk
	// then keep a reference to it here.
	queuedPlaylistUpdates map[string]string

	s3Client *s3.Client

	s3Secret string

	s3Bucket          string
	s3Region          string
	s3ServingEndpoint string
	s3AccessKey       string
	s3ACL             string
	s3PathPrefix      string

	s3Endpoint string
	host       string

	lock sync.Mutex

	configRepository models.EngineConfig

	performanceTracker *utils.PerformanceTracker

	s3ForcePathStyle bool
}

// NewS3Storage returns a new S3Storage instance.
func NewS3Storage(configRepository models.EngineConfig) *S3Storage {
	return &S3Storage{
		queuedPlaylistUpdates: make(map[string]string),
		lock:                  sync.Mutex{},
		configRepository:      configRepository,
		performanceTracker:    utils.NewPerformanceTracker(),
	}
}

// Setup sets up the s3 storage for saving the video to s3.
func (s *S3Storage) Setup() error {
	log.Trace("Setting up S3 for external storage of video...")
	s3Config := s.configRepository.GetS3Config()
	customVideoServingEndpoint := s.configRepository.GetVideoServingEndpoint()

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
	s.s3PathPrefix = s3Config.PathPrefix
	s.s3ForcePathStyle = s3Config.ForcePathStyle

	s.s3Client = s.createS3Client()

	return nil
}

// SegmentWritten is called when a single segment of video is written.
func (s *S3Storage) SegmentWritten(localFilePath string) {
	index := utils.GetIndexFromFilePath(localFilePath)
	performanceMonitorKey := "s3upload-" + index
	s.performanceTracker.StartPerformanceMonitor(performanceMonitorKey)

	// Upload the segment
	if _, err := s.Save(localFilePath, 0); err != nil {
		log.Errorln(err)
		return
	}
	averagePerformance := s.performanceTracker.GetAveragePerformance(performanceMonitorKey)

	// Warn the user about long-running save operations
	if averagePerformance != 0 {
		if averagePerformance > float64(s.configRepository.GetStreamLatencyLevel().SecondsPerSegment)*0.9 {
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
	s.lock.Lock()
	defer s.lock.Unlock()
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
	if err := rewritePlaylistLocations(localFilePath, s.host, s.remoteHLSPrefix()); err != nil {
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
	// Build the remote path using the shared HLS prefix.
	remotePath := s.remoteHLSPrefix() + normalizedPath

	maxAgeSeconds := utils.GetCacheDurationSecondsForPath(filePath)
	cacheControlHeader := fmt.Sprintf("max-age=%d", maxAgeSeconds)

	uploadInput := &s3.PutObjectInput{
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
		uploadInput.ACL = types.ObjectCannedACL(s.s3ACL)
	} else {
		// Default ACL
		uploadInput.ACL = types.ObjectCannedACLPublicRead
	}

	_, err = s.s3Client.PutObject(context.Background(), uploadInput)
	if err != nil {
		log.Traceln("error uploading segment", err.Error())
		if retryCount < 4 {
			log.Traceln("Retrying...")
			return s.Save(filePath, retryCount+1)
		}

		return "", fmt.Errorf("giving up uploading %s to object storage %s", filePath, s.s3Endpoint)
	}

	return url.JoinPath(s.host, remotePath)
}

// Cleanup will fire the different cleanup tasks required.
func (s *S3Storage) Cleanup() error {
	if err := s.RemoteCleanup(); err != nil {
		log.Errorln(err)
	}

	return localCleanup(4)
}

// RemoteCleanup will remove old files from the remote storage provider.
func (s *S3Storage) RemoteCleanup() error {
	// Determine how many files we should keep on S3 storage
	maxNumber := s.configRepository.GetStreamLatencyLevel().SegmentCount
	buffer := 20

	keys, err := s.getDeletableVideoSegmentsWithOffset(maxNumber + buffer)
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		s.deleteObjects(keys)
	}

	return nil
}

func (s *S3Storage) createS3Client() *s3.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConnsPerHost = 100

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}

	creds := credentials.NewStaticCredentialsProvider(s.s3AccessKey, s.s3Secret, "")

	client := s3.New(s3.Options{
		Region:       s.s3Region,
		Credentials:  aws.NewCredentialsCache(creds),
		HTTPClient:   httpClient,
		BaseEndpoint: aws.String(s.s3Endpoint),
		UsePathStyle: s.s3ForcePathStyle,
	})

	return client
}

func (s *S3Storage) getDeletableVideoSegmentsWithOffset(offset int) ([]s3object, error) {
	objectsToDelete, err := s.retrieveAllVideoSegments()
	if err != nil {
		return nil, err
	}

	if len(objectsToDelete) <= offset {
		return nil, nil
	}

	return objectsToDelete[offset:], nil
}

// s3MaxDeleteKeys is the maximum number of keys per DeleteObjects request.
const s3MaxDeleteKeys = 1000

func (s *S3Storage) deleteObjects(objects []s3object) {
	keys := make([]types.ObjectIdentifier, len(objects))
	for i, object := range objects {
		keys[i] = types.ObjectIdentifier{Key: aws.String(object.key)}
	}

	log.Debugln("Deleting", len(keys), "objects from S3 bucket:", s.s3Bucket)

	for i := 0; i < len(keys); i += s3MaxDeleteKeys {
		end := i + s3MaxDeleteKeys
		if end > len(keys) {
			end = len(keys)
		}

		resp, err := s.s3Client.DeleteObjects(context.Background(), &s3.DeleteObjectsInput{
			Bucket: aws.String(s.s3Bucket),
			Delete: &types.Delete{
				Objects: keys[i:end],
				Quiet:   aws.Bool(true),
			},
		})
		if err != nil {
			log.Errorf("Unable to delete objects from bucket %q, %v\n", s.s3Bucket, err)
		} else if len(resp.Errors) > 0 {
			log.Errorf("Failed to delete %d objects from bucket %q, first error: %v\n", len(resp.Errors), s.s3Bucket, resp.Errors[0])
		}
	}
}

func (s *S3Storage) retrieveAllVideoSegments() ([]s3object, error) {
	paginator := s3.NewListObjectsV2Paginator(s.s3Client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.s3Bucket),
		Prefix: aws.String(s.remoteHLSListingPrefix()),
	})

	var allObjects []s3object
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, errors.Wrap(err, "Unable to fetch list of items in bucket for cleanup")
		}
		for _, item := range page.Contents {
			if strings.HasSuffix(*item.Key, ".ts") {
				allObjects = append(allObjects, s3object{
					key:          *item.Key,
					lastModified: *item.LastModified,
				})
			}
		}
	}

	// Sort the results by timestamp, newest first.
	sort.Slice(allObjects, func(i, j int) bool {
		return allObjects[i].lastModified.After(allObjects[j].lastModified)
	})

	return allObjects, nil
}

type s3object struct {
	lastModified time.Time
	key          string
}

// remoteHLSPrefix returns the S3 key prefix under which HLS segments are
// stored, without a trailing slash. normalizedPath from Save() starts with
// "/" so concatenation produces "hls/0/segment.ts" or "myprefix/hls/0/segment.ts".
func (s *S3Storage) remoteHLSPrefix() string {
	prefix := "hls"
	if s.s3PathPrefix != "" {
		prefix = strings.Trim(s.s3PathPrefix, "/") + "/" + prefix
	}
	return prefix
}

// remoteHLSListingPrefix returns remoteHLSPrefix with a trailing slash,
// ensuring S3 listing is directory-scoped and won't match sibling keys
// like "hls-archive/" or "hls_old/".
func (s *S3Storage) remoteHLSListingPrefix() string {
	return s.remoteHLSPrefix() + "/"
}
