package configrepository

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/static"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/handlers/generated"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type SqlConfigRepository struct {
	datastore *data.Datastore
}

// NOTE: This is temporary during the transition period.
var temporaryGlobalInstance ConfigRepository

// Get will return the user repository.
func Get() ConfigRepository {
	if temporaryGlobalInstance == nil {
		i := New(data.GetDatastore())
		temporaryGlobalInstance = i
	}
	return temporaryGlobalInstance
}

// New will create a new instance of the UserRepository.
func New(datastore *data.Datastore) ConfigRepository {
	r := SqlConfigRepository{
		datastore: datastore,
	}

	migrateDatastoreValues(datastore, &r)

	// Set the server initialization date if needed.
	if hasSetInitDate, _ := r.GetServerInitTime(); hasSetInitDate == nil || !hasSetInitDate.Valid {
		_ = r.SetServerInitTime(time.Now())
	}

	if !r.HasPopulatedDefaults() {
		r.PopulateDefaults()
	}

	return &r
}

// GetExtraPageBodyContent will return the user-supplied body content.
func (r *SqlConfigRepository) GetExtraPageBodyContent() string {
	content, err := r.datastore.GetString(extraContentKey)
	if err != nil {
		log.Traceln(extraContentKey, err)
		return config.GetDefaults().PageBodyContent
	}

	return content
}

// SetExtraPageBodyContent will set the user-supplied body content.
func (r *SqlConfigRepository) SetExtraPageBodyContent(content string) error {
	return r.datastore.SetString(extraContentKey, content)
}

// GetStreamTitle will return the name of the current stream.
func (r *SqlConfigRepository) GetStreamTitle() string {
	title, err := r.datastore.GetString(streamTitleKey)
	if err != nil {
		return ""
	}

	return title
}

// SetStreamTitle will set the name of the current stream.
func (r *SqlConfigRepository) SetStreamTitle(title string) error {
	return r.datastore.SetString(streamTitleKey, title)
}

// GetAdminPassword will return the admin password.
func (r *SqlConfigRepository) GetAdminPassword() string {
	key, _ := r.datastore.GetString(adminPasswordKey)
	return key
}

// SetAdminPassword will set the admin password.
func (r *SqlConfigRepository) SetAdminPassword(key string) error {
	hashed_pass, err := utils.HashPassword(key)
	if err != nil {
		return err
	}
	return r.datastore.SetString(adminPasswordKey, hashed_pass)
}

// GetLogoPath will return the path for the logo, relative to webroot.
func (r *SqlConfigRepository) GetLogoPath() string {
	logo, err := r.datastore.GetString(logoPathKey)
	if err != nil {
		log.Traceln(logoPathKey, err)
		return config.GetDefaults().Logo
	}

	if logo == "" {
		return config.GetDefaults().Logo
	}

	return logo
}

// SetLogoPath will set the path for the logo, relative to webroot.
func (r *SqlConfigRepository) SetLogoPath(logo string) error {
	return r.datastore.SetString(logoPathKey, logo)
}

// SetLogoUniquenessString will set the logo cache busting string.
func (r *SqlConfigRepository) SetLogoUniquenessString(uniqueness string) error {
	return r.datastore.SetString(logoUniquenessKey, uniqueness)
}

// GetLogoUniquenessString will return the logo cache busting string.
func (r *SqlConfigRepository) GetLogoUniquenessString() string {
	uniqueness, err := r.datastore.GetString(logoUniquenessKey)
	if err != nil {
		log.Traceln(logoUniquenessKey, err)
		return ""
	}

	return uniqueness
}

// GetServerSummary will return the server summary text.
func (r *SqlConfigRepository) GetServerSummary() string {
	summary, err := r.datastore.GetString(serverSummaryKey)
	if err != nil {
		log.Traceln(serverSummaryKey, err)
		return ""
	}

	return summary
}

// SetServerSummary will set the server summary text.
func (r *SqlConfigRepository) SetServerSummary(summary string) error {
	return r.datastore.SetString(serverSummaryKey, summary)
}

// GetServerWelcomeMessage will return the server welcome message text.
func (r *SqlConfigRepository) GetServerWelcomeMessage() string {
	welcomeMessage, err := r.datastore.GetString(serverWelcomeMessageKey)
	if err != nil {
		log.Traceln(serverWelcomeMessageKey, err)
		return config.GetDefaults().ServerWelcomeMessage
	}

	return welcomeMessage
}

// SetServerWelcomeMessage will set the server welcome message text.
func (r *SqlConfigRepository) SetServerWelcomeMessage(welcomeMessage string) error {
	return r.datastore.SetString(serverWelcomeMessageKey, welcomeMessage)
}

// GetServerName will return the server name text.
func (r *SqlConfigRepository) GetServerName() string {
	name, err := r.datastore.GetString(serverNameKey)
	if err != nil {
		log.Traceln(serverNameKey, err)
		return config.GetDefaults().Name
	}

	return name
}

// SetServerName will set the server name text.
func (r *SqlConfigRepository) SetServerName(name string) error {
	return r.datastore.SetString(serverNameKey, name)
}

// GetServerURL will return the server URL.
func (r *SqlConfigRepository) GetServerURL() string {
	url, err := r.datastore.GetString(serverURLKey)
	if err != nil {
		return ""
	}

	return url
}

// SetServerURL will set the server URL.
func (r *SqlConfigRepository) SetServerURL(url string) error {
	return r.datastore.SetString(serverURLKey, url)
}

// GetHTTPPortNumber will return the server HTTP port.
func (r *SqlConfigRepository) GetHTTPPortNumber() int {
	port, err := r.datastore.GetNumber(httpPortNumberKey)
	if err != nil {
		log.Traceln(httpPortNumberKey, err)
		return config.GetDefaults().WebServerPort
	}

	if port == 0 {
		return config.GetDefaults().WebServerPort
	}
	return int(port)
}

// SetWebsocketOverrideHost will set the host override for websockets.
func (r *SqlConfigRepository) SetWebsocketOverrideHost(host string) error {
	return r.datastore.SetString(websocketHostOverrideKey, host)
}

// GetWebsocketOverrideHost will return the host override for websockets.
func (r *SqlConfigRepository) GetWebsocketOverrideHost() string {
	host, _ := r.datastore.GetString(websocketHostOverrideKey)

	return host
}

// SetHTTPPortNumber will set the server HTTP port.
func (r *SqlConfigRepository) SetHTTPPortNumber(port float64) error {
	return r.datastore.SetNumber(httpPortNumberKey, port)
}

// GetHTTPListenAddress will return the HTTP listen address.
func (r *SqlConfigRepository) GetHTTPListenAddress() string {
	address, err := r.datastore.GetString(httpListenAddressKey)
	if err != nil {
		log.Traceln(httpListenAddressKey, err)
		return config.GetDefaults().WebServerIP
	}
	return address
}

// SetHTTPListenAddress will set the server HTTP listen address.
func (r *SqlConfigRepository) SetHTTPListenAddress(address string) error {
	return r.datastore.SetString(httpListenAddressKey, address)
}

// GetRTMPPortNumber will return the server RTMP port.
func (r *SqlConfigRepository) GetRTMPPortNumber() int {
	port, err := r.datastore.GetNumber(rtmpPortNumberKey)
	if err != nil {
		log.Traceln(rtmpPortNumberKey, err)
		return config.GetDefaults().RTMPServerPort
	}

	if port == 0 {
		return config.GetDefaults().RTMPServerPort
	}

	return int(port)
}

// SetRTMPPortNumber will set the server RTMP port.
func (r *SqlConfigRepository) SetRTMPPortNumber(port float64) error {
	return r.datastore.SetNumber(rtmpPortNumberKey, port)
}

// GetServerMetadataTags will return the metadata tags.
func (r *SqlConfigRepository) GetServerMetadataTags() []string {
	tagsString, err := r.datastore.GetString(serverMetadataTagsKey)
	if tagsString == "" {
		return []string{}
	}

	if err != nil {
		log.Traceln(serverMetadataTagsKey, err)
		return []string{}
	}

	return strings.Split(tagsString, ",")
}

// SetServerMetadataTags will return the metadata tags.
func (r *SqlConfigRepository) SetServerMetadataTags(tags []string) error {
	tagString := strings.Join(tags, ",")
	return r.datastore.SetString(serverMetadataTagsKey, tagString)
}

// GetDirectoryEnabled will return if this server should register to YP.
func (r *SqlConfigRepository) GetDirectoryEnabled() bool {
	enabled, err := r.datastore.GetBool(directoryEnabledKey)
	if err != nil {
		return config.GetDefaults().YPEnabled
	}

	return enabled
}

// SetDirectoryEnabled will set if this server should register to YP.
func (r *SqlConfigRepository) SetDirectoryEnabled(enabled bool) error {
	return r.datastore.SetBool(directoryEnabledKey, enabled)
}

// SetDirectoryRegistrationKey will set the YP protocol registration key.
func (r *SqlConfigRepository) SetDirectoryRegistrationKey(key string) error {
	return r.datastore.SetString(directoryRegistrationKeyKey, key)
}

// GetDirectoryRegistrationKey will return the YP protocol registration key.
func (r *SqlConfigRepository) GetDirectoryRegistrationKey() string {
	key, _ := r.datastore.GetString(directoryRegistrationKeyKey)
	return key
}

// GetSocialHandles will return the external social links.
func (r *SqlConfigRepository) GetSocialHandles() []models.SocialHandle {
	var socialHandles []models.SocialHandle

	configEntry, err := r.datastore.Get(socialHandlesKey)
	if err != nil {
		log.Traceln(socialHandlesKey, err)
		return socialHandles
	}

	if err := configEntry.GetObject(&socialHandles); err != nil {
		log.Traceln(err)
		return socialHandles
	}

	return socialHandles
}

// SetSocialHandles will set the external social links.
func (r *SqlConfigRepository) SetSocialHandles(socialHandles []models.SocialHandle) error {
	configEntry := models.ConfigEntry{Key: socialHandlesKey, Value: socialHandles}
	return r.datastore.Save(configEntry)
}

// GetPeakSessionViewerCount will return the max number of viewers for this stream.
func (r *SqlConfigRepository) GetPeakSessionViewerCount() int {
	count, err := r.datastore.GetNumber(peakViewersSessionKey)
	if err != nil {
		return 0
	}
	return int(count)
}

// SetPeakSessionViewerCount will set the max number of viewers for this stream.
func (r *SqlConfigRepository) SetPeakSessionViewerCount(count int) error {
	return r.datastore.SetNumber(peakViewersSessionKey, float64(count))
}

// GetPeakOverallViewerCount will return the overall max number of viewers.
func (r *SqlConfigRepository) GetPeakOverallViewerCount() int {
	count, err := r.datastore.GetNumber(peakViewersOverallKey)
	if err != nil {
		return 0
	}
	return int(count)
}

// SetPeakOverallViewerCount will set the overall max number of viewers.
func (r *SqlConfigRepository) SetPeakOverallViewerCount(count int) error {
	return r.datastore.SetNumber(peakViewersOverallKey, float64(count))
}

// GetLastDisconnectTime will return the time the last stream ended.
func (r *SqlConfigRepository) GetLastDisconnectTime() (*utils.NullTime, error) {
	var disconnectTime utils.NullTime

	configEntry, err := r.datastore.Get(lastDisconnectTimeKey)
	if err != nil {
		return nil, err
	}

	if err := configEntry.GetObject(&disconnectTime); err != nil {
		return nil, err
	}

	if !disconnectTime.Valid || disconnectTime.Time.IsZero() {
		return nil, err
	}

	return &disconnectTime, nil
}

// SetLastDisconnectTime will set the time the last stream ended.
func (r *SqlConfigRepository) SetLastDisconnectTime(disconnectTime time.Time) error {
	savedDisconnectTime := utils.NullTime{Time: disconnectTime, Valid: true}
	configEntry := models.ConfigEntry{Key: lastDisconnectTimeKey, Value: savedDisconnectTime}
	return r.datastore.Save(configEntry)
}

// SetNSFW will set if this stream has NSFW content.
func (r *SqlConfigRepository) SetNSFW(isNSFW bool) error {
	return r.datastore.SetBool(nsfwKey, isNSFW)
}

// GetNSFW will return if this stream has NSFW content.
func (r *SqlConfigRepository) GetNSFW() bool {
	nsfw, err := r.datastore.GetBool(nsfwKey)
	if err != nil {
		return false
	}
	return nsfw
}

// SetFfmpegPath will set the custom ffmpeg path.
func (r *SqlConfigRepository) SetFfmpegPath(path string) error {
	return r.datastore.SetString(ffmpegPathKey, path)
}

// GetFfMpegPath will return the ffmpeg path.
func (r *SqlConfigRepository) GetFfMpegPath() string {
	path, err := r.datastore.GetString(ffmpegPathKey)
	if err != nil {
		return ""
	}
	return path
}

// GetS3Config will return the external storage configuration.
func (r *SqlConfigRepository) GetS3Config() models.S3 {
	configEntry, err := r.datastore.Get(s3StorageConfigKey)
	if err != nil {
		return models.S3{Enabled: false}
	}

	var s3Config models.S3
	if err := configEntry.GetObject(&s3Config); err != nil {
		return models.S3{Enabled: false}
	}

	return s3Config
}

// SetS3Config will set the external storage configuration.
func (r *SqlConfigRepository) SetS3Config(config models.S3) error {
	configEntry := models.ConfigEntry{Key: s3StorageConfigKey, Value: config}
	return r.datastore.Save(configEntry)
}

// GetStreamLatencyLevel will return the stream latency level.
func (r *SqlConfigRepository) GetStreamLatencyLevel() models.LatencyLevel {
	level, err := r.datastore.GetNumber(videoLatencyLevel)
	if err != nil {
		level = 2 // default
	} else if level > 4 {
		level = 4 // highest
	}

	return models.GetLatencyLevel(int(level))
}

// SetStreamLatencyLevel will set the stream latency level.
func (r *SqlConfigRepository) SetStreamLatencyLevel(level float64) error {
	return r.datastore.SetNumber(videoLatencyLevel, level)
}

// GetStreamOutputVariants will return all of the stream output variants.
func (r *SqlConfigRepository) GetStreamOutputVariants() []models.StreamOutputVariant {
	configEntry, err := r.datastore.Get(videoStreamOutputVariantsKey)
	if err != nil {
		return config.GetDefaults().StreamVariants
	}

	var streamOutputVariants []models.StreamOutputVariant
	if err := configEntry.GetObject(&streamOutputVariants); err != nil {
		return config.GetDefaults().StreamVariants
	}

	if len(streamOutputVariants) == 0 {
		return config.GetDefaults().StreamVariants
	}

	return streamOutputVariants
}

// SetStreamOutputVariants will set the stream output variants.
func (r *SqlConfigRepository) SetStreamOutputVariants(variants []models.StreamOutputVariant) error {
	configEntry := models.ConfigEntry{Key: videoStreamOutputVariantsKey, Value: variants}
	return r.datastore.Save(configEntry)
}

// SetChatDisabled will disable chat if set to true.
func (r *SqlConfigRepository) SetChatDisabled(disabled bool) error {
	return r.datastore.SetBool(chatDisabledKey, disabled)
}

// GetChatDisabled will return if chat is disabled.
func (r *SqlConfigRepository) GetChatDisabled() bool {
	disabled, err := r.datastore.GetBool(chatDisabledKey)
	if err == nil {
		return disabled
	}

	return false
}

// SetChatEstablishedUsersOnlyMode sets the state of established user only mode.
func (r *SqlConfigRepository) SetChatEstablishedUsersOnlyMode(enabled bool) error {
	return r.datastore.SetBool(chatEstablishedUsersOnlyModeKey, enabled)
}

// GetChatEstbalishedUsersOnlyMode returns the state of established user only mode.
func (r *SqlConfigRepository) GetChatEstbalishedUsersOnlyMode() bool {
	enabled, err := r.datastore.GetBool(chatEstablishedUsersOnlyModeKey)
	if err == nil {
		return enabled
	}

	return false
}

// SetChatSpamProtectionEnabled will enable chat spam protection if set to true.
func (r *SqlConfigRepository) SetChatSpamProtectionEnabled(enabled bool) error {
	return r.datastore.SetBool(chatSpamProtectionEnabledKey, enabled)
}

// GetChatSpamProtectionEnabled will return if chat spam protection is enabled.
func (r *SqlConfigRepository) GetChatSpamProtectionEnabled() bool {
	enabled, err := r.datastore.GetBool(chatSpamProtectionEnabledKey)
	if err == nil {
		return enabled
	}

	return true
}

// SetChatSlurFilterEnabled will enable the chat slur filter.
func (r *SqlConfigRepository) SetChatSlurFilterEnabled(enabled bool) error {
	return r.datastore.SetBool(chatSlurFilterEnabledKey, enabled)
}

// GetChatSlurFilterEnabled will return if the chat slur filter is enabled.
func (r *SqlConfigRepository) GetChatSlurFilterEnabled() bool {
	enabled, err := r.datastore.GetBool(chatSlurFilterEnabledKey)
	if err == nil {
		return enabled
	}

	return false
}

// GetExternalActions will return the registered external actions.
func (r *SqlConfigRepository) GetExternalActions() []models.ExternalAction {
	configEntry, err := r.datastore.Get(externalActionsKey)
	if err != nil {
		return []models.ExternalAction{}
	}

	var externalActions []models.ExternalAction
	if err := configEntry.GetObject(&externalActions); err != nil {
		return []models.ExternalAction{}
	}

	return externalActions
}

// SetExternalActions will save external actions.
func (r *SqlConfigRepository) SetExternalActions(actions []models.ExternalAction) error {
	configEntry := models.ConfigEntry{Key: externalActionsKey, Value: actions}
	return r.datastore.Save(configEntry)
}

// SetCustomStyles will save a string with CSS to insert into the page.
func (r *SqlConfigRepository) SetCustomStyles(styles string) error {
	return r.datastore.SetString(customStylesKey, styles)
}

// GetCustomStyles will return a string with CSS to insert into the page.
func (r *SqlConfigRepository) GetCustomStyles() string {
	style, err := r.datastore.GetString(customStylesKey)
	if err != nil {
		return ""
	}

	return style
}

// SetCustomJavascript will save a string with Javascript to insert into the page.
func (r *SqlConfigRepository) SetCustomJavascript(styles string) error {
	return r.datastore.SetString(customJavascriptKey, styles)
}

// GetCustomJavascript will return a string with Javascript to insert into the page.
func (r *SqlConfigRepository) GetCustomJavascript() string {
	style, err := r.datastore.GetString(customJavascriptKey)
	if err != nil {
		return ""
	}

	return style
}

// SetVideoCodec will set the codec used for video encoding.
func (r *SqlConfigRepository) SetVideoCodec(codec string) error {
	return r.datastore.SetString(videoCodecKey, codec)
}

// GetVideoCodec returns the codec to use for transcoding video.
func (r *SqlConfigRepository) GetVideoCodec() string {
	codec, err := r.datastore.GetString(videoCodecKey)
	if codec == "" || err != nil {
		return "libx264" // Default value
	}

	return codec
}

// VerifySettings will perform a sanity check for specific settings values.
func (r *SqlConfigRepository) VerifySettings() error {
	if len(r.GetStreamKeys()) == 0 && config.TemporaryStreamKey == "" {
		log.Errorln("No stream key set. Streaming is disabled. Please set one via the admin or command line arguments")
	}

	if r.GetAdminPassword() == "" {
		return errors.New("no admin password set. Please set one via the admin or command line arguments")
	}

	logoPath := r.GetLogoPath()
	if !utils.DoesFileExists(filepath.Join(config.DataDirectory, logoPath)) {
		log.Traceln(logoPath, "not found in the data directory. copying a default logo.")
		logo := static.GetLogo()
		if err := os.WriteFile(filepath.Join(config.DataDirectory, "logo.png"), logo, 0o600); err != nil {
			return errors.Wrap(err, "failed to write logo to disk")
		}
		if err := r.SetLogoPath("logo.png"); err != nil {
			return errors.Wrap(err, "failed to save logo filename")
		}
	}

	return nil
}

// FindHighestVideoQualityIndex will return the highest quality from a slice of variants.
func (r *SqlConfigRepository) FindHighestVideoQualityIndex(qualities []models.StreamOutputVariant) (int, bool) {
	type IndexedQuality struct {
		quality models.StreamOutputVariant
		index   int
	}

	if len(qualities) < 2 {
		return 0, qualities[0].IsVideoPassthrough
	}

	indexedQualities := make([]IndexedQuality, 0)
	for index, quality := range qualities {
		indexedQuality := IndexedQuality{quality, index}
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

	// nolint:gosec
	selectedQuality := indexedQualities[0]
	return selectedQuality.index, selectedQuality.quality.IsVideoPassthrough
}

// GetForbiddenUsernameList will return the blocked usernames as a comma separated string.
func (r *SqlConfigRepository) GetForbiddenUsernameList() []string {
	usernames, err := r.datastore.GetStringSlice(blockedUsernamesKey)
	if err != nil {
		return config.DefaultForbiddenUsernames
	}

	if len(usernames) == 0 {
		return config.DefaultForbiddenUsernames
	}

	return usernames
}

// SetForbiddenUsernameList set the username blocklist as a comma separated string.
func (r *SqlConfigRepository) SetForbiddenUsernameList(usernames []string) error {
	return r.datastore.SetStringSlice(blockedUsernamesKey, usernames)
}

// GetSuggestedUsernamesList will return the suggested usernames.
// If the number of suggested usernames is smaller than 10, the number pool is
// not used (see code in the CreateAnonymousUser function).
func (r *SqlConfigRepository) GetSuggestedUsernamesList() []string {
	usernames, err := r.datastore.GetStringSlice(suggestedUsernamesKey)

	if err != nil || len(usernames) == 0 {
		return []string{}
	}

	return usernames
}

// SetSuggestedUsernamesList sets the username suggestion list.
func (r *SqlConfigRepository) SetSuggestedUsernamesList(usernames []string) error {
	return r.datastore.SetStringSlice(suggestedUsernamesKey, usernames)
}

// GetServerInitTime will return when the server was first setup.
func (r *SqlConfigRepository) GetServerInitTime() (*utils.NullTime, error) {
	var t utils.NullTime

	configEntry, err := r.datastore.Get(serverInitDateKey)
	if err != nil {
		return nil, err
	}

	if err := configEntry.GetObject(&t); err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, err
	}

	return &t, nil
}

// SetServerInitTime will set when the server was first created.
func (r *SqlConfigRepository) SetServerInitTime(t time.Time) error {
	nt := utils.NullTime{Time: t, Valid: true}
	configEntry := models.ConfigEntry{Key: serverInitDateKey, Value: nt}
	return r.datastore.Save(configEntry)
}

// SetFederationEnabled will enable federation if set to true.
func (r *SqlConfigRepository) SetFederationEnabled(enabled bool) error {
	return r.datastore.SetBool(federationEnabledKey, enabled)
}

// GetFederationEnabled will return if federation is enabled.
func (r *SqlConfigRepository) GetFederationEnabled() bool {
	enabled, err := r.datastore.GetBool(federationEnabledKey)
	if err == nil {
		return enabled
	}

	return false
}

// SetFederationUsername will set the username used in federated activities.
func (r *SqlConfigRepository) SetFederationUsername(username string) error {
	return r.datastore.SetString(federationUsernameKey, username)
}

// GetFederationUsername will return the username used in federated activities.
func (r *SqlConfigRepository) GetFederationUsername() string {
	username, err := r.datastore.GetString(federationUsernameKey)
	if username == "" || err != nil {
		return config.GetDefaults().FederationUsername
	}

	return username
}

// SetFederationGoLiveMessage will set the message sent when going live.
func (r *SqlConfigRepository) SetFederationGoLiveMessage(message string) error {
	return r.datastore.SetString(federationGoLiveMessageKey, message)
}

// GetFederationGoLiveMessage will return the message sent when going live.
func (r *SqlConfigRepository) GetFederationGoLiveMessage() string {
	// Empty message means it's disabled.
	message, err := r.datastore.GetString(federationGoLiveMessageKey)
	if err != nil {
		log.Traceln("unable to fetch go live message.", err)
	}

	return message
}

// SetFederationIsPrivate will set if federation activity is private.
func (r *SqlConfigRepository) SetFederationIsPrivate(isPrivate bool) error {
	return r.datastore.SetBool(federationPrivateKey, isPrivate)
}

// GetFederationIsPrivate will return if federation is private.
func (r *SqlConfigRepository) GetFederationIsPrivate() bool {
	isPrivate, err := r.datastore.GetBool(federationPrivateKey)
	if err == nil {
		return isPrivate
	}

	return false
}

// SetFederationShowEngagement will set if fediverse engagement shows in chat.
func (r *SqlConfigRepository) SetFederationShowEngagement(showEngagement bool) error {
	return r.datastore.SetBool(federationShowEngagementKey, showEngagement)
}

// GetFederationShowEngagement will return if fediverse engagement shows in chat.
func (r *SqlConfigRepository) GetFederationShowEngagement() bool {
	showEngagement, err := r.datastore.GetBool(federationShowEngagementKey)
	if err == nil {
		return showEngagement
	}

	return true
}

// SetBlockedFederatedDomains will set the blocked federated domains.
func (r *SqlConfigRepository) SetBlockedFederatedDomains(domains []string) error {
	return r.datastore.SetString(federationBlockedDomainsKey, strings.Join(domains, ","))
}

// GetBlockedFederatedDomains will return a list of blocked federated domains.
func (r *SqlConfigRepository) GetBlockedFederatedDomains() []string {
	domains, err := r.datastore.GetString(federationBlockedDomainsKey)
	if err != nil {
		return []string{}
	}

	if domains == "" {
		return []string{}
	}

	return strings.Split(domains, ",")
}

// SetChatJoinMessagesEnabled will set if chat join messages are enabled.
func (r *SqlConfigRepository) SetChatJoinMessagesEnabled(enabled bool) error {
	return r.datastore.SetBool(chatJoinMessagesEnabledKey, enabled)
}

// GetChatJoinPartMessagesEnabled will return if chat join messages are enabled.
func (r *SqlConfigRepository) GetChatJoinPartMessagesEnabled() bool {
	enabled, err := r.datastore.GetBool(chatJoinMessagesEnabledKey)
	if err != nil {
		return true
	}

	return enabled
}

// SetNotificationsEnabled will save the enabled state of notifications.
func (r *SqlConfigRepository) SetNotificationsEnabled(enabled bool) error {
	return r.datastore.SetBool(notificationsEnabledKey, enabled)
}

// GetNotificationsEnabled will return the enabled state of notifications.
func (r *SqlConfigRepository) GetNotificationsEnabled() bool {
	enabled, _ := r.datastore.GetBool(notificationsEnabledKey)
	return enabled
}

// GetDiscordConfig will return the Discord configuration.
func (r *SqlConfigRepository) GetDiscordConfig() models.DiscordConfiguration {
	configEntry, err := r.datastore.Get(discordConfigurationKey)
	if err != nil {
		return models.DiscordConfiguration{Enabled: false}
	}

	var config models.DiscordConfiguration
	if err := configEntry.GetObject(&config); err != nil {
		return models.DiscordConfiguration{Enabled: false}
	}

	return config
}

// SetDiscordConfig will set the Discord configuration.
func (r *SqlConfigRepository) SetDiscordConfig(config models.DiscordConfiguration) error {
	configEntry := models.ConfigEntry{Key: discordConfigurationKey, Value: config}
	return r.datastore.Save(configEntry)
}

// GetBrowserPushConfig will return the browser push configuration.
func (r *SqlConfigRepository) GetBrowserPushConfig() models.BrowserNotificationConfiguration {
	configEntry, err := r.datastore.Get(browserPushConfigurationKey)
	if err != nil {
		return models.BrowserNotificationConfiguration{Enabled: false}
	}

	var config models.BrowserNotificationConfiguration
	if err := configEntry.GetObject(&config); err != nil {
		return models.BrowserNotificationConfiguration{Enabled: false}
	}

	return config
}

// SetBrowserPushConfig will set the browser push configuration.
func (r *SqlConfigRepository) SetBrowserPushConfig(config models.BrowserNotificationConfiguration) error {
	configEntry := models.ConfigEntry{Key: browserPushConfigurationKey, Value: config}
	return r.datastore.Save(configEntry)
}

// SetBrowserPushPublicKey will set the public key for browser pushes.
func (r *SqlConfigRepository) SetBrowserPushPublicKey(key string) error {
	return r.datastore.SetString(browserPushPublicKeyKey, key)
}

// GetBrowserPushPublicKey will return the public key for browser pushes.
func (r *SqlConfigRepository) GetBrowserPushPublicKey() (string, error) {
	return r.datastore.GetString(browserPushPublicKeyKey)
}

// SetBrowserPushPrivateKey will set the private key for browser pushes.
func (r *SqlConfigRepository) SetBrowserPushPrivateKey(key string) error {
	return r.datastore.SetString(browserPushPrivateKeyKey, key)
}

// GetBrowserPushPrivateKey will return the private key for browser pushes.
func (r *SqlConfigRepository) GetBrowserPushPrivateKey() (string, error) {
	return r.datastore.GetString(browserPushPrivateKeyKey)
}

// SetHasPerformedInitialNotificationsConfig sets when performed initial setup.
func (r *SqlConfigRepository) SetHasPerformedInitialNotificationsConfig(hasConfigured bool) error {
	return r.datastore.SetBool(hasConfiguredInitialNotificationsKey, true)
}

// GetHasPerformedInitialNotificationsConfig gets when performed initial setup.
func (r *SqlConfigRepository) GetHasPerformedInitialNotificationsConfig() bool {
	configured, _ := r.datastore.GetBool(hasConfiguredInitialNotificationsKey)
	return configured
}

// GetHideViewerCount will return if the viewer count shold be hidden.
func (r *SqlConfigRepository) GetHideViewerCount() bool {
	hide, _ := r.datastore.GetBool(hideViewerCountKey)
	return hide
}

// SetHideViewerCount will set if the viewer count should be hidden.
func (r *SqlConfigRepository) SetHideViewerCount(hide bool) error {
	return r.datastore.SetBool(hideViewerCountKey, hide)
}

// GetCustomOfflineMessage will return the custom offline message.
func (r *SqlConfigRepository) GetCustomOfflineMessage() string {
	message, _ := r.datastore.GetString(customOfflineMessageKey)
	return message
}

// SetCustomOfflineMessage will set the custom offline message.
func (r *SqlConfigRepository) SetCustomOfflineMessage(message string) error {
	return r.datastore.SetString(customOfflineMessageKey, message)
}

// SetCustomColorVariableValues sets CSS variable names and values.
func (r *SqlConfigRepository) SetCustomColorVariableValues(variables map[string]string) error {
	return r.datastore.SetStringMap(customColorVariableValuesKey, variables)
}

// GetCustomColorVariableValues gets CSS variable names and values.
func (r *SqlConfigRepository) GetCustomColorVariableValues() map[string]string {
	values, _ := r.datastore.GetStringMap(customColorVariableValuesKey)
	return values
}

// GetStreamKeys will return valid stream keys.
func (r *SqlConfigRepository) GetStreamKeys() []generated.StreamKey {
	configEntry, err := r.datastore.Get(streamKeysKey)
	if err != nil {
		return []generated.StreamKey{}
	}

	var streamKeys []generated.StreamKey
	if err := configEntry.GetObject(&streamKeys); err != nil {
		return []generated.StreamKey{}
	}

	return streamKeys
}

// SetStreamKeys will set valid stream keys.
func (r *SqlConfigRepository) SetStreamKeys(actions []generated.StreamKey) error {
	configEntry := models.ConfigEntry{Key: streamKeysKey, Value: actions}
	return r.datastore.Save(configEntry)
}

// SetDisableSearchIndexing will set if the web server should be indexable.
func (r *SqlConfigRepository) SetDisableSearchIndexing(disableSearchIndexing bool) error {
	return r.datastore.SetBool(disableSearchIndexingKey, disableSearchIndexing)
}

// GetDisableSearchIndexing will return if the web server should be indexable.
func (r *SqlConfigRepository) GetDisableSearchIndexing() bool {
	disableSearchIndexing, err := r.datastore.GetBool(disableSearchIndexingKey)
	if err != nil {
		return false
	}
	return disableSearchIndexing
}

// GetVideoServingEndpoint returns the custom video endpont.
func (r *SqlConfigRepository) GetVideoServingEndpoint() string {
	message, _ := r.datastore.GetString(videoServingEndpointKey)
	return message
}

// SetVideoServingEndpoint sets the custom video endpoint.
func (r *SqlConfigRepository) SetVideoServingEndpoint(message string) error {
	return r.datastore.SetString(videoServingEndpointKey, message)
}
