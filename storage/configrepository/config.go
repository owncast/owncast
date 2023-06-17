package configrepository

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/config"
	"github.com/owncast/owncast/static"
	"github.com/owncast/owncast/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// GetExtraPageBodyContent will return the user-supplied body content.
func (cr *SqlConfigRepository) GetExtraPageBodyContent() string {
	content, err := cr.datastore.GetString(extraContentKey)
	if err != nil {
		log.Traceln(extraContentKey, err)
		return config.GetDefaults().PageBodyContent
	}

	return content
}

// SetExtraPageBodyContent will set the user-supplied body content.
func (cr *SqlConfigRepository) SetExtraPageBodyContent(content string) error {
	return cr.datastore.SetString(extraContentKey, content)
}

// GetStreamTitle will return the name of the current stream.
func (cr *SqlConfigRepository) GetStreamTitle() string {
	title, err := cr.datastore.GetString(streamTitleKey)
	if err != nil {
		return ""
	}

	return title
}

// SetStreamTitle will set the name of the current stream.
func (cr *SqlConfigRepository) SetStreamTitle(title string) error {
	return cr.datastore.SetString(streamTitleKey, title)
}

// GetAdminPassword will return the admin password.
func (cr *SqlConfigRepository) GetAdminPassword() string {
	key, _ := cr.datastore.GetString(adminPasswordKey)
	return key
}

// SetAdminPassword will set the admin password.
func (cr *SqlConfigRepository) SetAdminPassword(key string) error {
	return cr.datastore.SetString(adminPasswordKey, key)
}

// GetLogoPath will return the path for the logo, relative to webroot.
func (cr *SqlConfigRepository) GetLogoPath() string {
	logo, err := cr.datastore.GetString(logoPathKey)
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
func (cr *SqlConfigRepository) SetLogoPath(logo string) error {
	return cr.datastore.SetString(logoPathKey, logo)
}

// SetLogoUniquenessString will set the logo cache busting string.
func (cr *SqlConfigRepository) SetLogoUniquenessString(uniqueness string) error {
	return cr.datastore.SetString(logoUniquenessKey, uniqueness)
}

// GetLogoUniquenessString will return the logo cache busting string.
func (cr *SqlConfigRepository) GetLogoUniquenessString() string {
	uniqueness, err := cr.datastore.GetString(logoUniquenessKey)
	if err != nil {
		log.Traceln(logoUniquenessKey, err)
		return ""
	}

	return uniqueness
}

// GetServerSummary will return the server summary text.
func (cr *SqlConfigRepository) GetServerSummary() string {
	summary, err := cr.datastore.GetString(serverSummaryKey)
	if err != nil {
		log.Traceln(serverSummaryKey, err)
		return ""
	}

	return summary
}

// SetServerSummary will set the server summary text.
func (cr *SqlConfigRepository) SetServerSummary(summary string) error {
	return cr.datastore.SetString(serverSummaryKey, summary)
}

// GetServerWelcomeMessage will return the server welcome message text.
func (cr *SqlConfigRepository) GetServerWelcomeMessage() string {
	welcomeMessage, err := cr.datastore.GetString(serverWelcomeMessageKey)
	if err != nil {
		log.Traceln(serverWelcomeMessageKey, err)
		return config.GetDefaults().ServerWelcomeMessage
	}

	return welcomeMessage
}

// SetServerWelcomeMessage will set the server welcome message text.
func (cr *SqlConfigRepository) SetServerWelcomeMessage(welcomeMessage string) error {
	return cr.datastore.SetString(serverWelcomeMessageKey, welcomeMessage)
}

// GetServerName will return the server name text.
func (cr *SqlConfigRepository) GetServerName() string {
	name, err := cr.datastore.GetString(serverNameKey)
	if err != nil {
		log.Traceln(serverNameKey, err)
		return config.GetDefaults().Name
	}

	return name
}

// SetServerName will set the server name text.
func (cr *SqlConfigRepository) SetServerName(name string) error {
	return cr.datastore.SetString(serverNameKey, name)
}

// GetServerURL will return the server URL.
func (cr *SqlConfigRepository) GetServerURL() string {
	url, err := cr.datastore.GetString(serverURLKey)
	if err != nil {
		return ""
	}

	return url
}

// SetServerURL will set the server URL.
func (cr *SqlConfigRepository) SetServerURL(url string) error {
	return cr.datastore.SetString(serverURLKey, url)
}

// GetHTTPPortNumber will return the server HTTP port.
func (cr *SqlConfigRepository) GetHTTPPortNumber() int {
	port, err := cr.datastore.GetNumber(httpPortNumberKey)
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
func (cr *SqlConfigRepository) SetWebsocketOverrideHost(host string) error {
	return cr.datastore.SetString(websocketHostOverrideKey, host)
}

// GetWebsocketOverrideHost will return the host override for websockets.
func (cr *SqlConfigRepository) GetWebsocketOverrideHost() string {
	host, _ := cr.datastore.GetString(websocketHostOverrideKey)

	return host
}

// SetHTTPPortNumber will set the server HTTP port.
func (cr *SqlConfigRepository) SetHTTPPortNumber(port float64) error {
	return cr.datastore.SetNumber(httpPortNumberKey, port)
}

// GetHTTPListenAddress will return the HTTP listen address.
func (cr *SqlConfigRepository) GetHTTPListenAddress() string {
	address, err := cr.datastore.GetString(httpListenAddressKey)
	if err != nil {
		log.Traceln(httpListenAddressKey, err)
		return config.GetDefaults().WebServerIP
	}
	return address
}

// SetHTTPListenAddress will set the server HTTP listen address.
func (cr *SqlConfigRepository) SetHTTPListenAddress(address string) error {
	return cr.datastore.SetString(httpListenAddressKey, address)
}

// GetRTMPPortNumber will return the server RTMP port.
func (cr *SqlConfigRepository) GetRTMPPortNumber() int {
	port, err := cr.datastore.GetNumber(rtmpPortNumberKey)
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
func (cr *SqlConfigRepository) SetRTMPPortNumber(port float64) error {
	return cr.datastore.SetNumber(rtmpPortNumberKey, port)
}

// GetServerMetadataTags will return the metadata tags.
func (cr *SqlConfigRepository) GetServerMetadataTags() []string {
	tagsString, err := cr.datastore.GetString(serverMetadataTagsKey)
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
func (cr *SqlConfigRepository) SetServerMetadataTags(tags []string) error {
	tagString := strings.Join(tags, ",")
	return cr.datastore.SetString(serverMetadataTagsKey, tagString)
}

// GetDirectoryEnabled will return if this server should register to YP.
func (cr *SqlConfigRepository) GetDirectoryEnabled() bool {
	enabled, err := cr.datastore.GetBool(directoryEnabledKey)
	if err != nil {
		return config.GetDefaults().YPEnabled
	}

	return enabled
}

// SetDirectoryEnabled will set if this server should register to YP.
func (cr *SqlConfigRepository) SetDirectoryEnabled(enabled bool) error {
	return cr.datastore.SetBool(directoryEnabledKey, enabled)
}

// SetDirectoryRegistrationKey will set the YP protocol registration key.
func (cr *SqlConfigRepository) SetDirectoryRegistrationKey(key string) error {
	return cr.datastore.SetString(directoryRegistrationKeyKey, key)
}

// GetDirectoryRegistrationKey will return the YP protocol registration key.
func (cr *SqlConfigRepository) GetDirectoryRegistrationKey() string {
	key, _ := cr.datastore.GetString(directoryRegistrationKeyKey)
	return key
}

// GetSocialHandles will return the external social links.
func (cr *SqlConfigRepository) GetSocialHandles() []models.SocialHandle {
	var socialHandles []models.SocialHandle

	configEntry, err := cr.datastore.Get(socialHandlesKey)
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
func (cr *SqlConfigRepository) SetSocialHandles(socialHandles []models.SocialHandle) error {
	configEntry := models.ConfigEntry{Key: socialHandlesKey, Value: socialHandles}
	return cr.datastore.Save(configEntry)
}

// GetPeakSessionViewerCount will return the max number of viewers for this stream.
func (cr *SqlConfigRepository) GetPeakSessionViewerCount() int {
	count, err := cr.datastore.GetNumber(peakViewersSessionKey)
	if err != nil {
		return 0
	}
	return int(count)
}

// SetPeakSessionViewerCount will set the max number of viewers for this stream.
func (cr *SqlConfigRepository) SetPeakSessionViewerCount(count int) error {
	return cr.datastore.SetNumber(peakViewersSessionKey, float64(count))
}

// GetPeakOverallViewerCount will return the overall max number of viewers.
func (cr *SqlConfigRepository) GetPeakOverallViewerCount() int {
	count, err := cr.datastore.GetNumber(peakViewersOverallKey)
	if err != nil {
		return 0
	}
	return int(count)
}

// SetPeakOverallViewerCount will set the overall max number of viewers.
func (cr *SqlConfigRepository) SetPeakOverallViewerCount(count int) error {
	return cr.datastore.SetNumber(peakViewersOverallKey, float64(count))
}

// GetLastDisconnectTime will return the time the last stream ended.
func (cr *SqlConfigRepository) GetLastDisconnectTime() (*utils.NullTime, error) {
	var disconnectTime utils.NullTime

	configEntry, err := cr.datastore.Get(lastDisconnectTimeKey)
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
func (cr *SqlConfigRepository) SetLastDisconnectTime(disconnectTime time.Time) error {
	savedDisconnectTime := utils.NullTime{Time: disconnectTime, Valid: true}
	configEntry := models.ConfigEntry{Key: lastDisconnectTimeKey, Value: savedDisconnectTime}
	return cr.datastore.Save(configEntry)
}

// SetNSFW will set if this stream has NSFW content.
func (cr *SqlConfigRepository) SetNSFW(isNSFW bool) error {
	return cr.datastore.SetBool(nsfwKey, isNSFW)
}

// GetNSFW will return if this stream has NSFW content.
func (cr *SqlConfigRepository) GetNSFW() bool {
	nsfw, err := cr.datastore.GetBool(nsfwKey)
	if err != nil {
		return false
	}
	return nsfw
}

// SetFfmpegPath will set the custom ffmpeg path.
func (cr *SqlConfigRepository) SetFfmpegPath(path string) error {
	return cr.datastore.SetString(ffmpegPathKey, path)
}

// GetFfMpegPath will return the ffmpeg path.
func (cr *SqlConfigRepository) GetFfMpegPath() string {
	path, err := cr.datastore.GetString(ffmpegPathKey)
	if err != nil {
		return ""
	}
	return path
}

// GetS3Config will return the external storage configuration.
func (cr *SqlConfigRepository) GetS3Config() models.S3 {
	configEntry, err := cr.datastore.Get(s3StorageConfigKey)
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
func (cr *SqlConfigRepository) SetS3Config(config models.S3) error {
	configEntry := models.ConfigEntry{Key: s3StorageConfigKey, Value: config}
	return cr.datastore.Save(configEntry)
}

// GetStreamLatencyLevel will return the stream latency level.
func (cr *SqlConfigRepository) GetStreamLatencyLevel() models.LatencyLevel {
	level, err := cr.datastore.GetNumber(videoLatencyLevel)
	if err != nil {
		level = 2 // default
	} else if level > 4 {
		level = 4 // highest
	}

	return models.GetLatencyLevel(int(level))
}

// SetStreamLatencyLevel will set the stream latency level.
func (cr *SqlConfigRepository) SetStreamLatencyLevel(level float64) error {
	return cr.datastore.SetNumber(videoLatencyLevel, level)
}

// GetStreamOutputVariants will return all of the stream output variants.
func (cr *SqlConfigRepository) GetStreamOutputVariants() []models.StreamOutputVariant {
	configEntry, err := cr.datastore.Get(videoStreamOutputVariantsKey)
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
func (cr *SqlConfigRepository) SetStreamOutputVariants(variants []models.StreamOutputVariant) error {
	configEntry := models.ConfigEntry{Key: videoStreamOutputVariantsKey, Value: variants}
	return cr.datastore.Save(configEntry)
}

// SetChatDisabled will disable chat if set to true.
func (cr *SqlConfigRepository) SetChatDisabled(disabled bool) error {
	return cr.datastore.SetBool(chatDisabledKey, disabled)
}

// GetChatDisabled will return if chat is disabled.
func (cr *SqlConfigRepository) GetChatDisabled() bool {
	disabled, err := cr.datastore.GetBool(chatDisabledKey)
	if err == nil {
		return disabled
	}

	return false
}

// SetChatEstablishedUsersOnlyMode sets the state of established user only mode.
func (cr *SqlConfigRepository) SetChatEstablishedUsersOnlyMode(enabled bool) error {
	return cr.datastore.SetBool(chatEstablishedUsersOnlyModeKey, enabled)
}

// GetChatEstbalishedUsersOnlyMode returns the state of established user only mode.
func (cr *SqlConfigRepository) GetChatEstbalishedUsersOnlyMode() bool {
	enabled, err := cr.datastore.GetBool(chatEstablishedUsersOnlyModeKey)
	if err == nil {
		return enabled
	}

	return false
}

// GetExternalActions will return the registered external actions.
func (cr *SqlConfigRepository) GetExternalActions() []models.ExternalAction {
	configEntry, err := cr.datastore.Get(externalActionsKey)
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
func (cr *SqlConfigRepository) SetExternalActions(actions []models.ExternalAction) error {
	configEntry := models.ConfigEntry{Key: externalActionsKey, Value: actions}
	return cr.datastore.Save(configEntry)
}

// SetCustomStyles will save a string with CSS to insert into the page.
func (cr *SqlConfigRepository) SetCustomStyles(styles string) error {
	return cr.datastore.SetString(customStylesKey, styles)
}

// GetCustomStyles will return a string with CSS to insert into the page.
func (cr *SqlConfigRepository) GetCustomStyles() string {
	style, err := cr.datastore.GetString(customStylesKey)
	if err != nil {
		return ""
	}

	return style
}

// SetCustomJavascript will save a string with Javascript to insert into the page.
func (cr *SqlConfigRepository) SetCustomJavascript(styles string) error {
	return cr.datastore.SetString(customJavascriptKey, styles)
}

// GetCustomJavascript will return a string with Javascript to insert into the page.
func (cr *SqlConfigRepository) GetCustomJavascript() string {
	style, err := cr.datastore.GetString(customJavascriptKey)
	if err != nil {
		return ""
	}

	return style
}

// SetVideoCodec will set the codec used for video encoding.
func (cr *SqlConfigRepository) SetVideoCodec(codec string) error {
	return cr.datastore.SetString(videoCodecKey, codec)
}

// GetVideoCodec returns the codec to use for transcoding video.
func (cr *SqlConfigRepository) GetVideoCodec() string {
	codec, err := cr.datastore.GetString(videoCodecKey)
	if codec == "" || err != nil {
		return "libx264" // Default value
	}

	return codec
}

// VerifySettings will perform a sanity check for specific settings values.
func (cr *SqlConfigRepository) VerifySettings() error {
	c := config.GetConfig()
	if len(cr.GetStreamKeys()) == 0 && c.TemporaryStreamKey == "" {
		log.Errorln("No stream key set. Streaming is disabled. Please set one via the admin or command line arguments")
	}

	if cr.GetAdminPassword() == "" {
		return errors.New("no admin password set. Please set one via the admin or command line arguments")
	}

	logoPath := cr.GetLogoPath()
	if !utils.DoesFileExists(filepath.Join(config.DataDirectory, logoPath)) {
		log.Traceln(logoPath, "not found in the data directory. copying a default logo.")
		logo := static.GetLogo()
		if err := os.WriteFile(filepath.Join(config.DataDirectory, "logo.png"), logo, 0o600); err != nil {
			return errors.Wrap(err, "failed to write logo to disk")
		}
		if err := cr.SetLogoPath("logo.png"); err != nil {
			return errors.Wrap(err, "failed to save logo filename")
		}
	}

	return nil
}

// FindHighestVideoQualityIndex will return the highest quality from a slice of variants.
func (cr *SqlConfigRepository) FindHighestVideoQualityIndex(qualities []models.StreamOutputVariant) int {
	type IndexedQuality struct {
		index   int
		quality models.StreamOutputVariant
	}

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

// GetForbiddenUsernameList will return the blocked usernames as a comma separated string.
func (cr *SqlConfigRepository) GetForbiddenUsernameList() []string {
	usernames, err := cr.datastore.GetStringSlice(blockedUsernamesKey)
	if err != nil {
		return config.DefaultForbiddenUsernames
	}

	if len(usernames) == 0 {
		return config.DefaultForbiddenUsernames
	}

	return usernames
}

// SetForbiddenUsernameList set the username blocklist as a comma separated string.
func (cr *SqlConfigRepository) SetForbiddenUsernameList(usernames []string) error {
	return cr.datastore.SetStringSlice(blockedUsernamesKey, usernames)
}

// GetSuggestedUsernamesList will return the suggested usernames.
// If the number of suggested usernames is smaller than 10, the number pool is
// not used (see code in the CreateAnonymousUser function).
func (cr *SqlConfigRepository) GetSuggestedUsernamesList() []string {
	usernames, err := cr.datastore.GetStringSlice(suggestedUsernamesKey)

	if err != nil || len(usernames) == 0 {
		return []string{}
	}

	return usernames
}

// SetSuggestedUsernamesList sets the username suggestion list.
func (cr *SqlConfigRepository) SetSuggestedUsernamesList(usernames []string) error {
	return cr.datastore.SetStringSlice(suggestedUsernamesKey, usernames)
}

// GetServerInitTime will return when the server was first setup.
func (cr *SqlConfigRepository) GetServerInitTime() (*utils.NullTime, error) {
	var t utils.NullTime

	configEntry, err := cr.datastore.Get(serverInitDateKey)
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
func (cr *SqlConfigRepository) SetServerInitTime(t time.Time) error {
	nt := utils.NullTime{Time: t, Valid: true}
	configEntry := models.ConfigEntry{Key: serverInitDateKey, Value: nt}
	return cr.datastore.Save(configEntry)
}

// SetFederationEnabled will enable federation if set to true.
func (cr *SqlConfigRepository) SetFederationEnabled(enabled bool) error {
	return cr.datastore.SetBool(federationEnabledKey, enabled)
}

// GetFederationEnabled will return if federation is enabled.
func (cr *SqlConfigRepository) GetFederationEnabled() bool {
	enabled, err := cr.datastore.GetBool(federationEnabledKey)
	if err == nil {
		return enabled
	}

	return false
}

// SetFederationUsername will set the username used in federated activities.
func (cr *SqlConfigRepository) SetFederationUsername(username string) error {
	return cr.datastore.SetString(federationUsernameKey, username)
}

// GetFederationUsername will return the username used in federated activities.
func (cr *SqlConfigRepository) GetFederationUsername() string {
	username, err := cr.datastore.GetString(federationUsernameKey)
	if username == "" || err != nil {
		return config.GetDefaults().FederationUsername
	}

	return username
}

// SetFederationGoLiveMessage will set the message sent when going live.
func (cr *SqlConfigRepository) SetFederationGoLiveMessage(message string) error {
	return cr.datastore.SetString(federationGoLiveMessageKey, message)
}

// GetFederationGoLiveMessage will return the message sent when going live.
func (cr *SqlConfigRepository) GetFederationGoLiveMessage() string {
	// Empty message means it's disabled.
	message, err := cr.datastore.GetString(federationGoLiveMessageKey)
	if err != nil {
		log.Traceln("unable to fetch go live message.", err)
	}

	return message
}

// SetFederationIsPrivate will set if federation activity is private.
func (cr *SqlConfigRepository) SetFederationIsPrivate(isPrivate bool) error {
	return cr.datastore.SetBool(federationPrivateKey, isPrivate)
}

// GetFederationIsPrivate will return if federation is private.
func (cr *SqlConfigRepository) GetFederationIsPrivate() bool {
	isPrivate, err := cr.datastore.GetBool(federationPrivateKey)
	if err == nil {
		return isPrivate
	}

	return false
}

// SetFederationShowEngagement will set if fediverse engagement shows in chat.
func (cr *SqlConfigRepository) SetFederationShowEngagement(showEngagement bool) error {
	return cr.datastore.SetBool(federationShowEngagementKey, showEngagement)
}

// GetFederationShowEngagement will return if fediverse engagement shows in chat.
func (cr *SqlConfigRepository) GetFederationShowEngagement() bool {
	showEngagement, err := cr.datastore.GetBool(federationShowEngagementKey)
	if err == nil {
		return showEngagement
	}

	return true
}

// SetBlockedFederatedDomains will set the blocked federated domains.
func (cr *SqlConfigRepository) SetBlockedFederatedDomains(domains []string) error {
	return cr.datastore.SetString(federationBlockedDomainsKey, strings.Join(domains, ","))
}

// GetBlockedFederatedDomains will return a list of blocked federated domains.
func (cr *SqlConfigRepository) GetBlockedFederatedDomains() []string {
	domains, err := cr.datastore.GetString(federationBlockedDomainsKey)
	if err != nil {
		return []string{}
	}

	if domains == "" {
		return []string{}
	}

	return strings.Split(domains, ",")
}

// SetChatJoinMessagesEnabled will set if chat join messages are enabled.
func (cr *SqlConfigRepository) SetChatJoinMessagesEnabled(enabled bool) error {
	return cr.datastore.SetBool(chatJoinMessagesEnabledKey, enabled)
}

// GetChatJoinMessagesEnabled will return if chat join messages are enabled.
func (cr *SqlConfigRepository) GetChatJoinMessagesEnabled() bool {
	enabled, err := cr.datastore.GetBool(chatJoinMessagesEnabledKey)
	if err != nil {
		return true
	}

	return enabled
}

// SetNotificationsEnabled will save the enabled state of notifications.
func (cr *SqlConfigRepository) SetNotificationsEnabled(enabled bool) error {
	return cr.datastore.SetBool(notificationsEnabledKey, enabled)
}

// GetNotificationsEnabled will return the enabled state of notifications.
func (cr *SqlConfigRepository) GetNotificationsEnabled() bool {
	enabled, _ := cr.datastore.GetBool(notificationsEnabledKey)
	return enabled
}

// GetDiscordConfig will return the Discord configuration.
func (cr *SqlConfigRepository) GetDiscordConfig() models.DiscordConfiguration {
	configEntry, err := cr.datastore.Get(discordConfigurationKey)
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
func (cr *SqlConfigRepository) SetDiscordConfig(config models.DiscordConfiguration) error {
	configEntry := models.ConfigEntry{Key: discordConfigurationKey, Value: config}
	return cr.datastore.Save(configEntry)
}

// GetBrowserPushConfig will return the browser push configuration.
func (cr *SqlConfigRepository) GetBrowserPushConfig() models.BrowserNotificationConfiguration {
	configEntry, err := cr.datastore.Get(browserPushConfigurationKey)
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
func (cr *SqlConfigRepository) SetBrowserPushConfig(config models.BrowserNotificationConfiguration) error {
	configEntry := models.ConfigEntry{Key: browserPushConfigurationKey, Value: config}
	return cr.datastore.Save(configEntry)
}

// SetBrowserPushPublicKey will set the public key for browser pushes.
func (cr *SqlConfigRepository) SetBrowserPushPublicKey(key string) error {
	return cr.datastore.SetString(browserPushPublicKeyKey, key)
}

// GetBrowserPushPublicKey will return the public key for browser pushes.
func (cr *SqlConfigRepository) GetBrowserPushPublicKey() (string, error) {
	return cr.datastore.GetString(browserPushPublicKeyKey)
}

// SetBrowserPushPrivateKey will set the private key for browser pushes.
func (cr *SqlConfigRepository) SetBrowserPushPrivateKey(key string) error {
	return cr.datastore.SetString(browserPushPrivateKeyKey, key)
}

// GetBrowserPushPrivateKey will return the private key for browser pushes.
func (cr *SqlConfigRepository) GetBrowserPushPrivateKey() (string, error) {
	return cr.datastore.GetString(browserPushPrivateKeyKey)
}

// SetHasPerformedInitialNotificationsConfig sets when performed initial setup.
func (cr *SqlConfigRepository) SetHasPerformedInitialNotificationsConfig(hasConfigured bool) error {
	return cr.datastore.SetBool(hasConfiguredInitialNotificationsKey, true)
}

// GetHasPerformedInitialNotificationsConfig gets when performed initial setup.
func (cr *SqlConfigRepository) GetHasPerformedInitialNotificationsConfig() bool {
	configured, _ := cr.datastore.GetBool(hasConfiguredInitialNotificationsKey)
	return configured
}

// GetHideViewerCount will return if the viewer count shold be hidden.
func (cr *SqlConfigRepository) GetHideViewerCount() bool {
	hide, _ := cr.datastore.GetBool(hideViewerCountKey)
	return hide
}

// SetHideViewerCount will set if the viewer count should be hidden.
func (cr *SqlConfigRepository) SetHideViewerCount(hide bool) error {
	return cr.datastore.SetBool(hideViewerCountKey, hide)
}

// GetCustomOfflineMessage will return the custom offline message.
func (cr *SqlConfigRepository) GetCustomOfflineMessage() string {
	message, _ := cr.datastore.GetString(customOfflineMessageKey)
	return message
}

// SetCustomOfflineMessage will set the custom offline message.
func (cr *SqlConfigRepository) SetCustomOfflineMessage(message string) error {
	return cr.datastore.SetString(customOfflineMessageKey, message)
}

// SetCustomColorVariableValues sets CSS variable names and values.
func (cr *SqlConfigRepository) SetCustomColorVariableValues(variables map[string]string) error {
	return cr.datastore.SetStringMap(customColorVariableValuesKey, variables)
}

// GetCustomColorVariableValues gets CSS variable names and values.
func (cr *SqlConfigRepository) GetCustomColorVariableValues() map[string]string {
	values, _ := cr.datastore.GetStringMap(customColorVariableValuesKey)
	return values
}

// GetStreamKeys will return valid stream keys.
func (cr *SqlConfigRepository) GetStreamKeys() []models.StreamKey {
	configEntry, err := cr.datastore.Get(streamKeysKey)
	if err != nil {
		return []models.StreamKey{}
	}

	var streamKeys []models.StreamKey
	if err := configEntry.GetObject(&streamKeys); err != nil {
		return []models.StreamKey{}
	}

	return streamKeys
}

// SetStreamKeys will set valid stream keys.
func (cr *SqlConfigRepository) SetStreamKeys(actions []models.StreamKey) error {
	configEntry := models.ConfigEntry{Key: streamKeysKey, Value: actions}
	return cr.datastore.Save(configEntry)
}

// SetDisableSearchIndexing will set if the web server should be indexable.
func (cr *SqlConfigRepository) SetDisableSearchIndexing(disableSearchIndexing bool) error {
	return cr.datastore.SetBool(disableSearchIndexingKey, disableSearchIndexing)
}

// GetDisableSearchIndexing will return if the web server should be indexable.
func (cr *SqlConfigRepository) GetDisableSearchIndexing() bool {
	disableSearchIndexing, err := cr.datastore.GetBool(disableSearchIndexingKey)
	if err != nil {
		return false
	}
	return disableSearchIndexing
}

// GetVideoServingEndpoint returns the custom video endpont.
func (cr *SqlConfigRepository) GetVideoServingEndpoint() string {
	message, _ := cr.datastore.GetString(videoServingEndpointKey)
	return message
}

// SetVideoServingEndpoint sets the custom video endpoint.
func (cr *SqlConfigRepository) SetVideoServingEndpoint(message string) error {
	return cr.datastore.SetString(videoServingEndpointKey, message)
}
