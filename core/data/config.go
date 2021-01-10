package data

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

const EXTRA_CONTENT_KEY = "extra_page_content"
const STREAM_TITLE_KEY = "stream_title"
const SERVER_TITLE_KEY = "server_title"
const STREAM_KEY_KEY = "stream_key"
const LOGO_PATH_KEY = "logo_path"
const SERVER_SUMMARY_KEY = "server_summary"
const SERVER_NAME_KEY = "server_name"
const SERVER_URL_KEY = "server_url"
const HTTP_PORT_NUMBER_KEY = "http_port_number"
const RTMP_PORT_NUMBER_KEY = "rtmp_port_number"
const DISABLE_UPGRADE_CHECKS_KEY = "disable_upgrade_checks"
const SERVER_METADATA_TAGS_KEY = "server_metadata_tags"
const DIRECTORY_ENABLED_KEY = "directory_enabled"
const SOCIAL_HANDLES_KEY = "social_handles"
const PEAK_VIEWERS_SESSION_KEY = "peak_viewers_session"
const PEAK_VIEWERS_OVERALL_KEY = "peak_viewers_overall"
const LAST_DISCONNECT_TIME_KEY = "last_disconnect_time"
const FFMPEG_PATH_KEY = "ffmpeg_path"
const NSFW_KEY = "nsfw"
const S3_STORAGE_ENABLED_KEY = "s3_storage_enabled"
const S3_STORAGE_CONFIG_KEY = "s3_storage_config"
const VIDEO_SEGMENT_LENGTH_SECONDS_KEY = "video_segment_length_seconds"
const VIDEO_SEGMENTS_IN_PLAYLIST_KEY = "video_segments_in_seconds"
const VIDEO_STREAM_OUTPUT_VARIANTS_KEY = "video_stream_output_variants"

// GetExtraPageBodyContent will return the user-supplied body content.
func GetExtraPageBodyContent() string {
	content, err := _datastore.GetString(EXTRA_CONTENT_KEY)
	if err != nil {
		log.Errorln(EXTRA_CONTENT_KEY, err)
		return config.GetDefaults().PageBodyContent
	}

	return content
}

// SetExtraPageBodyContent will set the user-supplied body content.
func SetExtraPageBodyContent(content string) error {
	return _datastore.SetString(EXTRA_CONTENT_KEY, content)
}

// GetStreamTitle will return the name of the current stream.
func GetStreamTitle() string {
	title, err := _datastore.GetString(STREAM_TITLE_KEY)
	if err != nil {
		return ""
	}

	return title
}

// SetStreamTitle will set the name of the current stream.
func SetStreamTitle(title string) error {
	return _datastore.SetString(STREAM_TITLE_KEY, title)
}

// GetServerTitle will return the title of the server.
func GetServerTitle() string {
	title, err := _datastore.GetString(SERVER_TITLE_KEY)
	if err != nil {
		return config.GetDefaults().Title
	}

	return title
}

// SetServerTitle will set the title of the server.
func SetServerTitle(title string) error {
	return _datastore.SetString(SERVER_TITLE_KEY, title)
}

// GetStreamKey will return the inbound streaming password.
func GetStreamKey() string {
	key, err := _datastore.GetString(STREAM_KEY_KEY)
	if err != nil {
		log.Errorln(STREAM_KEY_KEY, err)
		return ""
	}

	return key
}

// SetStreamKey will set the inbound streaming password.
func SetStreamKey(key string) error {
	return _datastore.SetString(STREAM_KEY_KEY, key)
}

// GetLogoPath will return the path for the logo, relative to webroot.
func GetLogoPath() string {
	logo, err := _datastore.GetString(LOGO_PATH_KEY)
	if err != nil {
		log.Errorln(LOGO_PATH_KEY, err)
		return config.GetDefaults().InstanceDetails.Logo
	}

	if logo == "" {
		return config.GetDefaults().Logo
	}

	return logo
}

// SetLogoPath will set the path for the logo, relative to webroot.
func SetLogoPath(logo string) error {
	return _datastore.SetString(LOGO_PATH_KEY, logo)
}

func GetServerSummary() string {
	summary, err := _datastore.GetString(SERVER_SUMMARY_KEY)
	if err != nil {
		log.Errorln(SERVER_SUMMARY_KEY, err)
		return ""
	}

	return summary
}

func SetServerSummary(summary string) error {
	return _datastore.SetString(SERVER_SUMMARY_KEY, summary)
}

func GetServerName() string {
	name, err := _datastore.GetString(SERVER_NAME_KEY)
	if err != nil {
		log.Errorln(SERVER_NAME_KEY, err)
		return ""
	}

	return name
}

func SetServerName(name string) error {
	return _datastore.SetString(SERVER_NAME_KEY, name)
}

func GetServerURL() string {
	url, err := _datastore.GetString(SERVER_URL_KEY)
	if err != nil {
		log.Errorln(SERVER_URL_KEY, err)
		return ""
	}

	return url
}

func SetServerURL(url string) error {
	return _datastore.SetString(SERVER_URL_KEY, url)
}

func GetHTTPPortNumber() int {
	port, err := _datastore.GetNumber(HTTP_PORT_NUMBER_KEY)
	if err != nil {
		log.Errorln(HTTP_PORT_NUMBER_KEY, err)
		return config.GetDefaults().WebServerPort
	}

	if port == 0 {
		return config.GetDefaults().WebServerPort
	}
	return int(port)
}

func SetHTTPPortNumber(port int) error {
	return _datastore.SetNumber(HTTP_PORT_NUMBER_KEY, float32(port))
}

func GetRTMPPortNumber() int {
	port, err := _datastore.GetNumber(RTMP_PORT_NUMBER_KEY)
	if err != nil {
		log.Errorln(RTMP_PORT_NUMBER_KEY, err)
		return config.GetDefaults().RTMPServerPort
	}

	if port == 0 {
		return config.GetDefaults().RTMPServerPort
	}

	return int(port)
}

func SetRTMPPortNumber(port int) error {
	return _datastore.SetNumber(RTMP_PORT_NUMBER_KEY, float32(port))
}

func GetDisableUpgradeChecks() bool {
	disable, err := _datastore.GetBool(DISABLE_UPGRADE_CHECKS_KEY)
	if err != nil {
		return config.GetDefaults().DisableUpgradeChecks
	}

	return disable
}

func SetDisableUpgradeChecks(disable bool) error {
	return _datastore.SetBool(DISABLE_UPGRADE_CHECKS_KEY, disable)
}

func GetServerMetadataTags() []string {
	tagsString, err := _datastore.GetString(SERVER_METADATA_TAGS_KEY)
	if err != nil {
		log.Errorln(SERVER_METADATA_TAGS_KEY, err)
		return []string{}
	}

	return strings.Split(tagsString, ",")
}

func SetServerMetadataTags(tags []string) error {
	tagString := strings.Join(tags, ",")
	return _datastore.SetString(SERVER_METADATA_TAGS_KEY, tagString)
}

func GetDirectoryEnabled() bool {
	enabled, err := _datastore.GetBool(DIRECTORY_ENABLED_KEY)
	if err != nil {
		return config.GetDefaults().YPEnabled
	}

	return enabled
}

func SetDirectoryEnabled(enabled bool) error {
	return _datastore.SetBool(DIRECTORY_ENABLED_KEY, enabled)
}

func GetSocialHandles() []models.SocialHandle {
	var socialHandles []models.SocialHandle

	configEntry, err := _datastore.Get(SOCIAL_HANDLES_KEY)
	if err != nil {
		log.Errorln(SOCIAL_HANDLES_KEY, err)
		return socialHandles
	}

	if err := configEntry.getObject(&socialHandles); err != nil {
		log.Errorln(err)
		return socialHandles
	}

	return socialHandles
}

func SetSocialHandles(socialHandles []models.SocialHandle) error {
	var configEntry = ConfigEntry{Key: SOCIAL_HANDLES_KEY, Value: socialHandles}
	return _datastore.Save(configEntry)
}

func GetPeakSessionViewerCount() int {
	count, err := _datastore.GetNumber(PEAK_VIEWERS_SESSION_KEY)
	if err != nil {
		return 0
	}
	return int(count)
}

func SetPeakSessionViewerCount(count int) error {
	return _datastore.SetNumber(PEAK_VIEWERS_SESSION_KEY, float32(count))
}

func GetPeakOverallViewerCount() int {
	count, err := _datastore.GetNumber(PEAK_VIEWERS_OVERALL_KEY)
	if err != nil {
		return 0
	}
	return int(count)
}

func SetPeakOverallViewerCount(count int) error {
	return _datastore.SetNumber(PEAK_VIEWERS_OVERALL_KEY, float32(count))
}

func GetLastDisconnectTime() (time.Time, error) {
	var disconnectTime time.Time
	configEntry, err := _datastore.Get(LAST_DISCONNECT_TIME_KEY)
	if err != nil {
		return disconnectTime, err
	}

	if err := configEntry.getObject(disconnectTime); err != nil {
		return disconnectTime, err
	}

	return disconnectTime, nil
}

func SetLastDisconnectTime(disconnectTime time.Time) error {
	var configEntry = ConfigEntry{Key: LAST_DISCONNECT_TIME_KEY, Value: disconnectTime}
	return _datastore.Save(configEntry)
}

func SetNSFW(isNSFW bool) error {
	return _datastore.SetBool(NSFW_KEY, isNSFW)
}

func GetNSFW() bool {
	nsfw, err := _datastore.GetBool(NSFW_KEY)
	if err != nil {
		return false
	}
	return nsfw
}

func SetFfmpegPath(path string) error {
	return _datastore.SetString(FFMPEG_PATH_KEY, path)
}

func GetFfMpegPath() string {
	ffmpegPath, err := _datastore.GetString(FFMPEG_PATH_KEY)

	if err != nil {
		log.Warnln(ffmpegPath, "is an invalid path to ffmpeg will try to use a copy in your path, if possible")
	} else if ffmpegPath == "" {
		log.Debugln("there is no path to ffmpeg set.  will try to use a copy in your path, if possible")
	}

	// First look to see if ffmpeg is in the current working directory
	localCopy := "./ffmpeg"
	hasLocalCopyError := VerifyFFMpegPath(localCopy)
	if hasLocalCopyError == nil {
		// No error, so all is good.  Use the local copy.
		return localCopy
	}

	cmd := exec.Command("which", "ffmpeg")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln("Unable to determine path to ffmpeg. Please specify it in the config file.")
	}

	path := strings.TrimSpace(string(out))

	return path
}

func GetS3Config() models.S3 {
	config, err := _datastore.Get(S3_STORAGE_CONFIG_KEY)
	if err != nil {
		return models.S3{Enabled: false}
	}

	return config.Value.(models.S3)
}

func SetS3Config(config models.S3) error {
	var configEntry = ConfigEntry{Key: S3_STORAGE_CONFIG_KEY, Value: config}
	return _datastore.Save(configEntry)
}

func GetS3StorageEnabled() bool {
	enabled, err := _datastore.GetBool(S3_STORAGE_ENABLED_KEY)
	if err != nil {
		log.Errorln(err)
		return false
	}

	return enabled
}

func SetS3StorageEnabled(enabled bool) error {
	return _datastore.SetBool(S3_STORAGE_ENABLED_KEY, enabled)
}

func GetVideoSegmentLengthDuration() float32 {
	seconds, err := _datastore.GetNumber(VIDEO_SEGMENT_LENGTH_SECONDS_KEY)
	if err != nil || seconds == 0 {
		return 4
	}

	return seconds
}

func SetVideoSegmentLengthDuration(seconds float32) error {
	return _datastore.SetNumber(VIDEO_SEGMENT_LENGTH_SECONDS_KEY, seconds)
}

func GetVideoSegmentsInPlaylist() float32 {
	count, err := _datastore.GetNumber(VIDEO_SEGMENTS_IN_PLAYLIST_KEY)
	if err != nil || count == 0 {
		return 20
	}

	return count
}

func SetStreamOutputVariants(count float32) error {
	return _datastore.SetNumber(VIDEO_SEGMENTS_IN_PLAYLIST_KEY, count)
}

func GetStreamOutputVariants() []models.StreamOutputVariant {
	variants, err := _datastore.Get(VIDEO_STREAM_OUTPUT_VARIANTS_KEY)
	if err != nil || len(variants.Value.([]models.StreamOutputVariant)) == 0 {
		return config.GetDefaults().StreamVariants
	}

	return variants.Value.([]models.StreamOutputVariant)
}

func SetVideoOutputVariants(variants []models.StreamOutputVariant) error {
	var configEntry = ConfigEntry{Key: VIDEO_STREAM_OUTPUT_VARIANTS_KEY, Value: variants}
	return _datastore.Save(configEntry)
}

func VerifySettings() error {
	if GetStreamKey() == "" {
		return errors.New("No stream key set. Please set one in your config file.")
	}

	return nil
}

// VerifyFFMpegPath verifies that the path exists, is a file, and is executable.
func VerifyFFMpegPath(path string) error {
	stat, err := os.Stat(path)

	if os.IsNotExist(err) {
		return errors.New("ffmpeg path does not exist")
	}

	if err != nil {
		return fmt.Errorf("error while verifying the ffmpeg path: %s", err.Error())
	}

	if stat.IsDir() {
		return errors.New("ffmpeg path can not be a folder")
	}

	mode := stat.Mode()
	//source: https://stackoverflow.com/a/60128480
	if mode&0111 == 0 {
		return errors.New("ffmpeg path is not executable")
	}

	return nil
}

func FindHighestVideoQualityIndex() int {
	type IndexedQuality struct {
		index   int
		quality models.StreamOutputVariant
	}

	qualities := GetStreamOutputVariants()

	if len(qualities) < 2 {
		return 0
	}

	indexedQualities := make([]IndexedQuality, 0)
	for index, quality := range qualities {
		indexedQuality := IndexedQuality{index, quality}
		indexedQualities = append(indexedQualities, indexedQuality)
	}

	sort.Slice(indexedQualities, func(a, b int) bool {
		if indexedQualities[a].quality.IsVideoPassthrough && !indexedQualities[b].quality.IsVideoPassthrough {
			return true
		}

		if !indexedQualities[a].quality.IsVideoPassthrough && indexedQualities[b].quality.IsVideoPassthrough {
			return false
		}

		return indexedQualities[a].quality.VideoBitrate > indexedQualities[b].quality.VideoBitrate
	})

	return indexedQualities[0].index
}
