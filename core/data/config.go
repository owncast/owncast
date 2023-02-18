package data

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/static"
	"github.com/owncast/owncast/utils"
)

const (
	extraContentKey                      = "extra_page_content"
	streamTitleKey                       = "stream_title"
	adminPasswordKey                     = "admin_password_key"
	logoPathKey                          = "logo_path"
	logoUniquenessKey                    = "logo_uniqueness"
	serverSummaryKey                     = "server_summary"
	serverWelcomeMessageKey              = "server_welcome_message"
	serverNameKey                        = "server_name"
	serverURLKey                         = "server_url"
	httpPortNumberKey                    = "http_port_number"
	httpListenAddressKey                 = "http_listen_address"
	websocketHostOverrideKey             = "websocket_host_override"
	rtmpPortNumberKey                    = "rtmp_port_number"
	serverMetadataTagsKey                = "server_metadata_tags"
	directoryEnabledKey                  = "directory_enabled"
	directoryRegistrationKeyKey          = "directory_registration_key"
	socialHandlesKey                     = "social_handles"
	peakViewersSessionKey                = "peak_viewers_session"
	peakViewersOverallKey                = "peak_viewers_overall"
	lastDisconnectTimeKey                = "last_disconnect_time"
	ffmpegPathKey                        = "ffmpeg_path"
	nsfwKey                              = "nsfw"
	s3StorageConfigKey                   = "s3_storage_config"
	videoLatencyLevel                    = "video_latency_level"
	videoStreamOutputVariantsKey         = "video_stream_output_variants"
	chatDisabledKey                      = "chat_disabled"
	externalActionsKey                   = "external_actions"
	customStylesKey                      = "custom_styles"
	customJavascriptKey                  = "custom_javascript"
	videoCodecKey                        = "video_codec"
	blockedUsernamesKey                  = "blocked_usernames"
	publicKeyKey                         = "public_key"
	privateKeyKey                        = "private_key"
	serverInitDateKey                    = "server_init_date"
	federationEnabledKey                 = "federation_enabled"
	federationUsernameKey                = "federation_username"
	federationPrivateKey                 = "federation_private"
	federationGoLiveMessageKey           = "federation_go_live_message"
	federationShowEngagementKey          = "federation_show_engagement"
	federationBlockedDomainsKey          = "federation_blocked_domains"
	suggestedUsernamesKey                = "suggested_usernames"
	chatJoinMessagesEnabledKey           = "chat_join_messages_enabled"
	chatEstablishedUsersOnlyModeKey      = "chat_established_users_only_mode"
	notificationsEnabledKey              = "notifications_enabled"
	discordConfigurationKey              = "discord_configuration"
	browserPushConfigurationKey          = "browser_push_configuration"
	browserPushPublicKeyKey              = "browser_push_public_key"
	browserPushPrivateKeyKey             = "browser_push_private_key"
	hasConfiguredInitialNotificationsKey = "has_configured_initial_notifications"
	hideViewerCountKey                   = "hide_viewer_count"
	customOfflineMessageKey              = "custom_offline_message"
	customColorVariableValuesKey         = "custom_color_variable_values"
	streamKeysKey                        = "stream_keys"
)

// GetExtraPageBodyContent will return the user-supplied body content.
func (s *Service) GetExtraPageBodyContent() string {
	content, err := s.Store.GetString(extraContentKey)
	if err != nil {
		log.Traceln(extraContentKey, err)
		return config.GetDefaults().PageBodyContent
	}

	return content
}

// SetExtraPageBodyContent will set the user-supplied body content.
func (s *Service) SetExtraPageBodyContent(content string) error {
	return s.Store.SetString(extraContentKey, content)
}

// GetStreamTitle will return the name of the current stream.
func (s *Service) GetStreamTitle() string {
	title, err := s.Store.GetString(streamTitleKey)
	if err != nil {
		return ""
	}

	return title
}

// SetStreamTitle will set the name of the current stream.
func (s *Service) SetStreamTitle(title string) error {
	return s.Store.SetString(streamTitleKey, title)
}

// GetAdminPassword will return the admin password.
func (s *Service) GetAdminPassword() string {
	key, _ := s.Store.GetString(adminPasswordKey)
	return key
}

// SetAdminPassword will set the admin password.
func (s *Service) SetAdminPassword(key string) error {
	return s.Store.SetString(adminPasswordKey, key)
}

// GetLogoPath will return the path for the logo, relative to webroot.
func (s *Service) GetLogoPath() string {
	logo, err := s.Store.GetString(logoPathKey)
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
func (s *Service) SetLogoPath(logo string) error {
	return s.Store.SetString(logoPathKey, logo)
}

// SetLogoUniquenessString will set the logo cache busting string.
func (s *Service) SetLogoUniquenessString(uniqueness string) error {
	return s.Store.SetString(logoUniquenessKey, uniqueness)
}

// GetLogoUniquenessString will return the logo cache busting string.
func (s *Service) GetLogoUniquenessString() string {
	uniqueness, err := s.Store.GetString(logoUniquenessKey)
	if err != nil {
		log.Traceln(logoUniquenessKey, err)
		return ""
	}

	return uniqueness
}

// GetServerSummary will return the server summary text.
func (s *Service) GetServerSummary() string {
	summary, err := s.Store.GetString(serverSummaryKey)
	if err != nil {
		log.Traceln(serverSummaryKey, err)
		return ""
	}

	return summary
}

// SetServerSummary will set the server summary text.
func (s *Service) SetServerSummary(summary string) error {
	return s.Store.SetString(serverSummaryKey, summary)
}

// GetServerWelcomeMessage will return the server welcome message text.
func (s *Service) GetServerWelcomeMessage() string {
	welcomeMessage, err := s.Store.GetString(serverWelcomeMessageKey)
	if err != nil {
		log.Traceln(serverWelcomeMessageKey, err)
		return config.GetDefaults().ServerWelcomeMessage
	}

	return welcomeMessage
}

// SetServerWelcomeMessage will set the server welcome message text.
func (s *Service) SetServerWelcomeMessage(welcomeMessage string) error {
	return s.Store.SetString(serverWelcomeMessageKey, welcomeMessage)
}

// GetServerName will return the server name text.
func (s *Service) GetServerName() string {
	name, err := s.Store.GetString(serverNameKey)
	if err != nil {
		log.Traceln(serverNameKey, err)
		return config.GetDefaults().Name
	}

	return name
}

// SetServerName will set the server name text.
func (s *Service) SetServerName(name string) error {
	return s.Store.SetString(serverNameKey, name)
}

// GetServerURL will return the server URL.
func (s *Service) GetServerURL() string {
	url, err := s.Store.GetString(serverURLKey)
	if err != nil {
		return ""
	}

	return url
}

// SetServerURL will set the server URL.
func (s *Service) SetServerURL(url string) error {
	return s.Store.SetString(serverURLKey, url)
}

// GetHTTPPortNumber will return the server HTTP port.
func (s *Service) GetHTTPPortNumber() int {
	port, err := s.Store.GetNumber(httpPortNumberKey)
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
func (s *Service) SetWebsocketOverrideHost(host string) error {
	return s.Store.SetString(websocketHostOverrideKey, host)
}

// GetWebsocketOverrideHost will return the host override for websockets.
func (s *Service) GetWebsocketOverrideHost() string {
	host, _ := s.Store.GetString(websocketHostOverrideKey)

	return host
}

// SetHTTPPortNumber will set the server HTTP port.
func (s *Service) SetHTTPPortNumber(port float64) error {
	return s.Store.SetNumber(httpPortNumberKey, port)
}

// GetHTTPListenAddress will return the HTTP listen address.
func (s *Service) GetHTTPListenAddress() string {
	address, err := s.Store.GetString(httpListenAddressKey)
	if err != nil {
		log.Traceln(httpListenAddressKey, err)
		return config.GetDefaults().WebServerIP
	}
	return address
}

// SetHTTPListenAddress will set the server HTTP listen address.
func (s *Service) SetHTTPListenAddress(address string) error {
	return s.Store.SetString(httpListenAddressKey, address)
}

// GetRTMPPortNumber will return the server RTMP port.
func (s *Service) GetRTMPPortNumber() int {
	port, err := s.Store.GetNumber(rtmpPortNumberKey)
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
func (s *Service) SetRTMPPortNumber(port float64) error {
	return s.Store.SetNumber(rtmpPortNumberKey, port)
}

// GetServerMetadataTags will return the metadata tags.
func (s *Service) GetServerMetadataTags() []string {
	tagsString, err := s.Store.GetString(serverMetadataTagsKey)
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
func (s *Service) SetServerMetadataTags(tags []string) error {
	tagString := strings.Join(tags, ",")
	return s.Store.SetString(serverMetadataTagsKey, tagString)
}

// GetDirectoryEnabled will return if this server should register to YP.
func (s *Service) GetDirectoryEnabled() bool {
	enabled, err := s.Store.GetBool(directoryEnabledKey)
	if err != nil {
		return config.GetDefaults().YPEnabled
	}

	return enabled
}

// SetDirectoryEnabled will set if this server should register to YP.
func (s *Service) SetDirectoryEnabled(enabled bool) error {
	return s.Store.SetBool(directoryEnabledKey, enabled)
}

// SetDirectoryRegistrationKey will set the YP protocol registration key.
func (s *Service) SetDirectoryRegistrationKey(key string) error {
	return s.Store.SetString(directoryRegistrationKeyKey, key)
}

// GetDirectoryRegistrationKey will return the YP protocol registration key.
func (s *Service) GetDirectoryRegistrationKey() string {
	key, _ := s.Store.GetString(directoryRegistrationKeyKey)
	return key
}

// GetSocialHandles will return the external social links.
func (s *Service) GetSocialHandles() []models.SocialHandle {
	var socialHandles []models.SocialHandle

	configEntry, err := s.Store.Get(socialHandlesKey)
	if err != nil {
		log.Traceln(socialHandlesKey, err)
		return socialHandles
	}

	if err := configEntry.getObject(&socialHandles); err != nil {
		log.Traceln(err)
		return socialHandles
	}

	return socialHandles
}

// SetSocialHandles will set the external social links.
func (s *Service) SetSocialHandles(socialHandles []models.SocialHandle) error {
	configEntry := ConfigEntry{Key: socialHandlesKey, Value: socialHandles}
	return s.Store.Save(configEntry)
}

// GetPeakSessionViewerCount will return the max number of viewers for this stream.
func (s *Service) GetPeakSessionViewerCount() int {
	count, err := s.Store.GetNumber(peakViewersSessionKey)
	if err != nil {
		return 0
	}
	return int(count)
}

// SetPeakSessionViewerCount will set the max number of viewers for this stream.
func (s *Service) SetPeakSessionViewerCount(count int) error {
	return s.Store.SetNumber(peakViewersSessionKey, float64(count))
}

// GetPeakOverallViewerCount will return the overall max number of viewers.
func (s *Service) GetPeakOverallViewerCount() int {
	count, err := s.Store.GetNumber(peakViewersOverallKey)
	if err != nil {
		return 0
	}
	return int(count)
}

// SetPeakOverallViewerCount will set the overall max number of viewers.
func (s *Service) SetPeakOverallViewerCount(count int) error {
	return s.Store.SetNumber(peakViewersOverallKey, float64(count))
}

// GetLastDisconnectTime will return the time the last stream ended.
func (s *Service) GetLastDisconnectTime() (*utils.NullTime, error) {
	var disconnectTime utils.NullTime

	configEntry, err := s.Store.Get(lastDisconnectTimeKey)
	if err != nil {
		return nil, err
	}

	if err := configEntry.getObject(&disconnectTime); err != nil {
		return nil, err
	}

	if !disconnectTime.Valid || disconnectTime.Time.IsZero() {
		return nil, err
	}

	return &disconnectTime, nil
}

// SetLastDisconnectTime will set the time the last stream ended.
func (s *Service) SetLastDisconnectTime(disconnectTime time.Time) error {
	savedDisconnectTime := utils.NullTime{Time: disconnectTime, Valid: true}
	configEntry := ConfigEntry{Key: lastDisconnectTimeKey, Value: savedDisconnectTime}
	return s.Store.Save(configEntry)
}

// SetNSFW will set if this stream has NSFW content.
func (s *Service) SetNSFW(isNSFW bool) error {
	return s.Store.SetBool(nsfwKey, isNSFW)
}

// GetNSFW will return if this stream has NSFW content.
func (s *Service) GetNSFW() bool {
	nsfw, err := s.Store.GetBool(nsfwKey)
	if err != nil {
		return false
	}
	return nsfw
}

// SetFfmpegPath will set the custom ffmpeg path.
func (s *Service) SetFfmpegPath(path string) error {
	return s.Store.SetString(ffmpegPathKey, path)
}

// GetFfMpegPath will return the ffmpeg path.
func (s *Service) GetFfMpegPath() string {
	path, err := s.Store.GetString(ffmpegPathKey)
	if err != nil {
		return ""
	}
	return path
}

// GetS3Config will return the external storage configuration.
func (s *Service) GetS3Config() models.S3 {
	configEntry, err := s.Store.Get(s3StorageConfigKey)
	if err != nil {
		return models.S3{Enabled: false}
	}

	var s3Config models.S3
	if err := configEntry.getObject(&s3Config); err != nil {
		return models.S3{Enabled: false}
	}

	return s3Config
}

// SetS3Config will set the external storage configuration.
func (s *Service) SetS3Config(config models.S3) error {
	configEntry := ConfigEntry{Key: s3StorageConfigKey, Value: config}
	return s.Store.Save(configEntry)
}

// GetStreamLatencyLevel will return the stream latency level.
func (s *Service) GetStreamLatencyLevel() models.LatencyLevel {
	level, err := s.Store.GetNumber(videoLatencyLevel)
	if err != nil {
		level = 2 // default
	} else if level > 4 {
		level = 4 // highest
	}

	return models.GetLatencyLevel(int(level))
}

// SetStreamLatencyLevel will set the stream latency level.
func (s *Service) SetStreamLatencyLevel(level float64) error {
	return s.Store.SetNumber(videoLatencyLevel, level)
}

// GetStreamOutputVariants will return all of the stream output variants.
func (s *Service) GetStreamOutputVariants() []models.StreamOutputVariant {
	configEntry, err := s.Store.Get(videoStreamOutputVariantsKey)
	if err != nil {
		return config.GetDefaults().StreamVariants
	}

	var streamOutputVariants []models.StreamOutputVariant
	if err := configEntry.getObject(&streamOutputVariants); err != nil {
		return config.GetDefaults().StreamVariants
	}

	if len(streamOutputVariants) == 0 {
		return config.GetDefaults().StreamVariants
	}

	return streamOutputVariants
}

// SetStreamOutputVariants will set the stream output variants.
func (s *Service) SetStreamOutputVariants(variants []models.StreamOutputVariant) error {
	configEntry := ConfigEntry{Key: videoStreamOutputVariantsKey, Value: variants}
	return s.Store.Save(configEntry)
}

// SetChatDisabled will disable chat if set to true.
func (s *Service) SetChatDisabled(disabled bool) error {
	return s.Store.SetBool(chatDisabledKey, disabled)
}

// GetChatDisabled will return if chat is disabled.
func (s *Service) GetChatDisabled() bool {
	disabled, err := s.Store.GetBool(chatDisabledKey)
	if err == nil {
		return disabled
	}

	return false
}

// SetChatEstablishedUsersOnlyMode sets the state of established user only mode.
func (s *Service) SetChatEstablishedUsersOnlyMode(enabled bool) error {
	return s.Store.SetBool(chatEstablishedUsersOnlyModeKey, enabled)
}

// GetChatEstbalishedUsersOnlyMode returns the state of established user only mode.
func (s *Service) GetChatEstbalishedUsersOnlyMode() bool {
	enabled, err := s.Store.GetBool(chatEstablishedUsersOnlyModeKey)
	if err == nil {
		return enabled
	}

	return false
}

// GetExternalActions will return the registered external actions.
func (s *Service) GetExternalActions() []models.ExternalAction {
	configEntry, err := s.Store.Get(externalActionsKey)
	if err != nil {
		return []models.ExternalAction{}
	}

	var externalActions []models.ExternalAction
	if err := configEntry.getObject(&externalActions); err != nil {
		return []models.ExternalAction{}
	}

	return externalActions
}

// SetExternalActions will save external actions.
func (s *Service) SetExternalActions(actions []models.ExternalAction) error {
	configEntry := ConfigEntry{Key: externalActionsKey, Value: actions}
	return s.Store.Save(configEntry)
}

// SetCustomStyles will save a string with CSS to insert into the page.
func (s *Service) SetCustomStyles(styles string) error {
	return s.Store.SetString(customStylesKey, styles)
}

// GetCustomStyles will return a string with CSS to insert into the page.
func (s *Service) GetCustomStyles() string {
	style, err := s.Store.GetString(customStylesKey)
	if err != nil {
		return ""
	}

	return style
}

// SetCustomJavascript will save a string with Javascript to insert into the page.
func (s *Service) SetCustomJavascript(styles string) error {
	return s.Store.SetString(customJavascriptKey, styles)
}

// GetCustomJavascript will return a string with Javascript to insert into the page.
func (s *Service) GetCustomJavascript() string {
	style, err := s.Store.GetString(customJavascriptKey)
	if err != nil {
		return ""
	}

	return style
}

// SetVideoCodec will set the codec used for video encoding.
func (s *Service) SetVideoCodec(codec string) error {
	return s.Store.SetString(videoCodecKey, codec)
}

// GetVideoCodec returns the codec to use for transcoding video.
func (s *Service) GetVideoCodec() string {
	codec, err := s.Store.GetString(videoCodecKey)
	if codec == "" || err != nil {
		return "libx264" // Default value
	}

	return codec
}

// VerifySettings will perform a sanity check for specific settings values.
func (s *Service) VerifySettings() error {
	if len(s.GetStreamKeys()) == 0 && config.TemporaryStreamKey == "" {
		log.Errorln("No stream key set. Streaming is disabled. Please set one via the admin or command line arguments")
	}

	if s.GetAdminPassword() == "" {
		return errors.New("no admin password set. Please set one via the admin or command line arguments")
	}

	logoPath := s.GetLogoPath()
	if !utils.DoesFileExists(filepath.Join(config.DataDirectory, logoPath)) {
		log.Traceln(logoPath, "not found in the data directory. copying a default logo.")
		logo := static.GetLogo()
		if err := os.WriteFile(filepath.Join(config.DataDirectory, "logo.png"), logo, 0o600); err != nil {
			return errors.Wrap(err, "failed to write logo to disk")
		}
		if err := s.SetLogoPath("logo.png"); err != nil {
			return errors.Wrap(err, "failed to save logo filename")
		}
	}

	return nil
}

// FindHighestVideoQualityIndex will return the highest quality from a slice of variants.
func (s *Service) FindHighestVideoQualityIndex(qualities []models.StreamOutputVariant) int {
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
func (s *Service) GetForbiddenUsernameList() []string {
	usernames, err := s.Store.GetStringSlice(blockedUsernamesKey)
	if err != nil {
		return config.DefaultForbiddenUsernames
	}

	if len(usernames) == 0 {
		return config.DefaultForbiddenUsernames
	}

	return usernames
}

// SetForbiddenUsernameList set the username blocklist as a comma separated string.
func (s *Service) SetForbiddenUsernameList(usernames []string) error {
	return s.Store.SetStringSlice(blockedUsernamesKey, usernames)
}

// GetSuggestedUsernamesList will return the suggested usernames.
// If the number of suggested usernames is smaller than 10, the number pool is
// not used (see code in the CreateAnonymousUser function).
func (s *Service) GetSuggestedUsernamesList() []string {
	usernames, err := s.Store.GetStringSlice(suggestedUsernamesKey)

	if err != nil || len(usernames) == 0 {
		return []string{}
	}

	return usernames
}

// SetSuggestedUsernamesList sets the username suggestion list.
func (s *Service) SetSuggestedUsernamesList(usernames []string) error {
	return s.Store.SetStringSlice(suggestedUsernamesKey, usernames)
}

// GetServerInitTime will return when the server was first setup.
func (s *Service) GetServerInitTime() (*utils.NullTime, error) {
	var t utils.NullTime

	configEntry, err := s.Store.Get(serverInitDateKey)
	if err != nil {
		return nil, err
	}

	if err := configEntry.getObject(&t); err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, err
	}

	return &t, nil
}

// SetServerInitTime will set when the server was first created.
func (s *Service) SetServerInitTime(t time.Time) error {
	nt := utils.NullTime{Time: t, Valid: true}
	configEntry := ConfigEntry{Key: serverInitDateKey, Value: nt}
	return s.Store.Save(configEntry)
}

// SetFederationEnabled will enable federation if set to true.
func (s *Service) SetFederationEnabled(enabled bool) error {
	return s.Store.SetBool(federationEnabledKey, enabled)
}

// GetFederationEnabled will return if federation is enabled.
func (s *Service) GetFederationEnabled() bool {
	enabled, err := s.Store.GetBool(federationEnabledKey)
	if err == nil {
		return enabled
	}

	return false
}

// SetFederationUsername will set the username used in federated activities.
func (s *Service) SetFederationUsername(username string) error {
	return s.Store.SetString(federationUsernameKey, username)
}

// GetFederationUsername will return the username used in federated activities.
func (s *Service) GetFederationUsername() string {
	username, err := s.Store.GetString(federationUsernameKey)
	if username == "" || err != nil {
		return config.GetDefaults().FederationUsername
	}

	return username
}

// SetFederationGoLiveMessage will set the message sent when going live.
func (s *Service) SetFederationGoLiveMessage(message string) error {
	return s.Store.SetString(federationGoLiveMessageKey, message)
}

// GetFederationGoLiveMessage will return the message sent when going live.
func (s *Service) GetFederationGoLiveMessage() string {
	// Empty message means it's disabled.
	message, err := s.Store.GetString(federationGoLiveMessageKey)
	if err != nil {
		log.Traceln("unable to fetch go live message.", err)
	}

	return message
}

// SetFederationIsPrivate will set if federation activity is private.
func (s *Service) SetFederationIsPrivate(isPrivate bool) error {
	return s.Store.SetBool(federationPrivateKey, isPrivate)
}

// GetFederationIsPrivate will return if federation is private.
func (s *Service) GetFederationIsPrivate() bool {
	isPrivate, err := s.Store.GetBool(federationPrivateKey)
	if err == nil {
		return isPrivate
	}

	return false
}

// SetFederationShowEngagement will set if fediverse engagement shows in chat.
func (s *Service) SetFederationShowEngagement(showEngagement bool) error {
	return s.Store.SetBool(federationShowEngagementKey, showEngagement)
}

// GetFederationShowEngagement will return if fediverse engagement shows in chat.
func (s *Service) GetFederationShowEngagement() bool {
	showEngagement, err := s.Store.GetBool(federationShowEngagementKey)
	if err == nil {
		return showEngagement
	}

	return true
}

// SetBlockedFederatedDomains will set the blocked federated domains.
func (s *Service) SetBlockedFederatedDomains(domains []string) error {
	return s.Store.SetString(federationBlockedDomainsKey, strings.Join(domains, ","))
}

// GetBlockedFederatedDomains will return a list of blocked federated domains.
func (s *Service) GetBlockedFederatedDomains() []string {
	domains, err := s.Store.GetString(federationBlockedDomainsKey)
	if err != nil {
		return []string{}
	}

	if domains == "" {
		return []string{}
	}

	return strings.Split(domains, ",")
}

// SetChatJoinMessagesEnabled will set if chat join messages are enabled.
func (s *Service) SetChatJoinMessagesEnabled(enabled bool) error {
	return s.Store.SetBool(chatJoinMessagesEnabledKey, enabled)
}

// GetChatJoinMessagesEnabled will return if chat join messages are enabled.
func (s *Service) GetChatJoinMessagesEnabled() bool {
	enabled, err := s.Store.GetBool(chatJoinMessagesEnabledKey)
	if err != nil {
		return true
	}

	return enabled
}

// SetNotificationsEnabled will save the enabled state of notifications.
func (s *Service) SetNotificationsEnabled(enabled bool) error {
	return s.Store.SetBool(notificationsEnabledKey, enabled)
}

// GetNotificationsEnabled will return the enabled state of notifications.
func (s *Service) GetNotificationsEnabled() bool {
	enabled, _ := s.Store.GetBool(notificationsEnabledKey)
	return enabled
}

// GetDiscordConfig will return the Discord configuration.
func (s *Service) GetDiscordConfig() models.DiscordConfiguration {
	configEntry, err := s.Store.Get(discordConfigurationKey)
	if err != nil {
		return models.DiscordConfiguration{Enabled: false}
	}

	var config models.DiscordConfiguration
	if err := configEntry.getObject(&config); err != nil {
		return models.DiscordConfiguration{Enabled: false}
	}

	return config
}

// SetDiscordConfig will set the Discord configuration.
func (s *Service) SetDiscordConfig(config models.DiscordConfiguration) error {
	configEntry := ConfigEntry{Key: discordConfigurationKey, Value: config}
	return s.Store.Save(configEntry)
}

// GetBrowserPushConfig will return the browser push configuration.
func (s *Service) GetBrowserPushConfig() models.BrowserNotificationConfiguration {
	configEntry, err := s.Store.Get(browserPushConfigurationKey)
	if err != nil {
		return models.BrowserNotificationConfiguration{Enabled: false}
	}

	var config models.BrowserNotificationConfiguration
	if err := configEntry.getObject(&config); err != nil {
		return models.BrowserNotificationConfiguration{Enabled: false}
	}

	return config
}

// SetBrowserPushConfig will set the browser push configuration.
func (s *Service) SetBrowserPushConfig(config models.BrowserNotificationConfiguration) error {
	configEntry := ConfigEntry{Key: browserPushConfigurationKey, Value: config}
	return s.Store.Save(configEntry)
}

// SetBrowserPushPublicKey will set the public key for browser pushes.
func (s *Service) SetBrowserPushPublicKey(key string) error {
	return s.Store.SetString(browserPushPublicKeyKey, key)
}

// GetBrowserPushPublicKey will return the public key for browser pushes.
func (s *Service) GetBrowserPushPublicKey() (string, error) {
	return s.Store.GetString(browserPushPublicKeyKey)
}

// SetBrowserPushPrivateKey will set the private key for browser pushes.
func (s *Service) SetBrowserPushPrivateKey(key string) error {
	return s.Store.SetString(browserPushPrivateKeyKey, key)
}

// GetBrowserPushPrivateKey will return the private key for browser pushes.
func (s *Service) GetBrowserPushPrivateKey() (string, error) {
	return s.Store.GetString(browserPushPrivateKeyKey)
}

// SetHasPerformedInitialNotificationsConfig sets when performed initial setup.
func (s *Service) SetHasPerformedInitialNotificationsConfig(hasConfigured bool) error {
	return s.Store.SetBool(hasConfiguredInitialNotificationsKey, true)
}

// GetHasPerformedInitialNotificationsConfig gets when performed initial setup.
func (s *Service) GetHasPerformedInitialNotificationsConfig() bool {
	configured, _ := s.Store.GetBool(hasConfiguredInitialNotificationsKey)
	return configured
}

// GetHideViewerCount will return if the viewer count shold be hidden.
func (s *Service) GetHideViewerCount() bool {
	hide, _ := s.Store.GetBool(hideViewerCountKey)
	return hide
}

// SetHideViewerCount will set if the viewer count should be hidden.
func (s *Service) SetHideViewerCount(hide bool) error {
	return s.Store.SetBool(hideViewerCountKey, hide)
}

// GetCustomOfflineMessage will return the custom offline message.
func (s *Service) GetCustomOfflineMessage() string {
	message, _ := s.Store.GetString(customOfflineMessageKey)
	return message
}

// SetCustomOfflineMessage will set the custom offline message.
func (s *Service) SetCustomOfflineMessage(message string) error {
	return s.Store.SetString(customOfflineMessageKey, message)
}

// SetCustomColorVariableValues sets CSS variable names and values.
func (s *Service) SetCustomColorVariableValues(variables map[string]string) error {
	return s.Store.SetStringMap(customColorVariableValuesKey, variables)
}

// GetCustomColorVariableValues gets CSS variable names and values.
func (s *Service) GetCustomColorVariableValues() map[string]string {
	values, _ := s.Store.GetStringMap(customColorVariableValuesKey)
	return values
}

// GetStreamKeys will return valid stream keys.
func (s *Service) GetStreamKeys() []models.StreamKey {
	configEntry, err := s.Store.Get(streamKeysKey)
	if err != nil {
		return []models.StreamKey{}
	}

	var streamKeys []models.StreamKey
	if err := configEntry.getObject(&streamKeys); err != nil {
		return []models.StreamKey{}
	}

	return streamKeys
}

// SetStreamKeys will set valid stream keys.
func (s *Service) SetStreamKeys(actions []models.StreamKey) error {
	configEntry := ConfigEntry{Key: streamKeysKey, Value: actions}
	return s.Store.Save(configEntry)
}
