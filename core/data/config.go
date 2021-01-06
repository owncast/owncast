package data

import (
	"strings"
	"time"

	"github.com/owncast/owncast/config"
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

// GetExtraPageBodyContent will return the user-supplied body content.
func GetExtraPageBodyContent() string {
	content, err := _datastore.GetString(EXTRA_CONTENT_KEY)
	if err != nil {
		log.Errorln(EXTRA_CONTENT_KEY, err)
		return ""
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
		log.Errorln(STREAM_TITLE_KEY, err)
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
		log.Errorln(SERVER_TITLE_KEY, err)
		return ""
	}

	return title
}

// SetServerTitle will set the title of the server.
func SetServerTitle(title string) error {
	return _datastore.SetString(SERVER_TITLE_KEY, title)
}

// GetStreamKey will return the inbound streaming password.
func GetStreamKey() string {
	title, err := _datastore.GetString(STREAM_KEY_KEY)
	if err != nil {
		log.Errorln(STREAM_KEY_KEY, err)
		return ""
	}

	return title
}

// SetStreamKey will set the inbound streaming password.
func SetStreamKey(key string) error {
	return _datastore.SetString(STREAM_KEY_KEY, key)
}

// GetLogoPath will return the path for the logo, relative to webroot.
func GetLogoPath() string {
	title, err := _datastore.GetString(LOGO_PATH_KEY)
	if err != nil {
		log.Errorln(LOGO_PATH_KEY, err)
		return ""
	}

	return title
}

// SetLogoPath will set the path for the logo, relative to webroot.
func SetLogoPath(key string) error {
	return _datastore.SetString(LOGO_PATH_KEY, key)
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
		return 8080
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
		return 8080
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
		return config.GetDefaults().YP.Enabled
	}

	return enabled
}

func SetDirectoryEnabled(enabled bool) error {
	return _datastore.SetBool(DIRECTORY_ENABLED_KEY, enabled)
}

func GetSocialHandles() []config.SocialHandle {
	var socialHandles []config.SocialHandle

	configEntry, err := _datastore.Get(SOCIAL_HANDLES_KEY)
	if err != nil {
		log.Errorln(SOCIAL_HANDLES_KEY, err)
		return socialHandles
	}

	if err := configEntry.getObject(socialHandles); err != nil {
		log.Errorln(err)
		return socialHandles
	}

	return socialHandles
}

func SetSocialHandles(socialHandles []config.SocialHandle) error {
	var configEntry = ConfigEntry{Key: SOCIAL_HANDLES_KEY, Value: socialHandles}
	return _datastore.Save(configEntry)
}

func GetPeakSessionViewerCount() int {
	count, err := _datastore.GetNumber(PEAK_VIEWERS_SESSION_KEY)
	if err != nil {
		log.Errorln(PEAK_VIEWERS_SESSION_KEY, err)
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
		log.Errorln(PEAK_VIEWERS_OVERALL_KEY, err)
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
