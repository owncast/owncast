package handlers

import (
	"net/http"
)

func (s *ServerInterfaceImpl) SetAdminPassword(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetAdminPassword)(w, r)
}

func (s *ServerInterfaceImpl) SetAdminPasswordOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetAdminPassword)(w, r)
}

func (s *ServerInterfaceImpl) SetStreamKeys(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetStreamKeys)(w, r)
}

func (s *ServerInterfaceImpl) SetStreamKeysOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetStreamKeys)(w, r)
}

func (s *ServerInterfaceImpl) SetExtraPageContent(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetExtraPageContent)(w, r)
}

func (s *ServerInterfaceImpl) SetExtraPageContentOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetExtraPageContent)(w, r)
}

func (s *ServerInterfaceImpl) SetStreamTitle(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetStreamTitle)(w, r)
}

func (s *ServerInterfaceImpl) SetStreamTitleOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetStreamTitle)(w, r)
}

func (s *ServerInterfaceImpl) SetServerName(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetServerName)(w, r)
}

func (s *ServerInterfaceImpl) SetServerNameOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetServerName)(w, r)
}

func (s *ServerInterfaceImpl) SetServerSummary(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetServerSummary)(w, r)
}

func (s *ServerInterfaceImpl) SetServerSummaryOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetServerSummary)(w, r)
}

func (s *ServerInterfaceImpl) SetCustomOfflineMessage(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetCustomOfflineMessage)(w, r)
}

func (s *ServerInterfaceImpl) SetCustomOfflineMessageOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetCustomOfflineMessage)(w, r)
}

func (s *ServerInterfaceImpl) SetServerWelcomeMessage(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetServerWelcomeMessage)(w, r)
}

func (s *ServerInterfaceImpl) SetServerWelcomeMessageOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetServerWelcomeMessage)(w, r)
}

func (s *ServerInterfaceImpl) SetChatDisabled(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetChatDisabled)(w, r)
}

func (s *ServerInterfaceImpl) SetChatDisabledOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetChatDisabled)(w, r)
}

func (s *ServerInterfaceImpl) SetChatJoinMessagesEnabled(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetChatJoinMessagesEnabled)(w, r)
}

func (s *ServerInterfaceImpl) SetChatJoinMessagesEnabledOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetChatJoinMessagesEnabled)(w, r)
}

func (s *ServerInterfaceImpl) SetEnableEstablishedChatUserMode(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetEnableEstablishedChatUserMode)(w, r)
}

func (s *ServerInterfaceImpl) SetEnableEstablishedChatUserModeOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetEnableEstablishedChatUserMode)(w, r)
}

func (s *ServerInterfaceImpl) SetForbiddenUsernameList(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetForbiddenUsernameList)(w, r)
}

func (s *ServerInterfaceImpl) SetForbiddenUsernameListOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetForbiddenUsernameList)(w, r)
}

func (s *ServerInterfaceImpl) SetSuggestedUsernameList(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetSuggestedUsernameList)(w, r)
}

func (s *ServerInterfaceImpl) SetSuggestedUsernameListOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetSuggestedUsernameList)(w, r)
}

func (s *ServerInterfaceImpl) SetChatSpamProtectionEnabled(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetChatSpamProtectionEnabled)(w, r)
}

func (s *ServerInterfaceImpl) SetChatSpamProtectionEnabledOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetChatSpamProtectionEnabled)(w, r)
}

func (s *ServerInterfaceImpl) SetChatSlurFilterEnabled(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetChatSlurFilterEnabled)(w, r)
}

func (s *ServerInterfaceImpl) SetChatSlurFilterEnabledOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetChatSlurFilterEnabled)(w, r)
}

func (s *ServerInterfaceImpl) SetChatRequireAuthentication(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetChatRequireAuthentication)(w, r)
}

func (s *ServerInterfaceImpl) SetChatRequireAuthenticationOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetChatRequireAuthentication)(w, r)
}

func (s *ServerInterfaceImpl) SetVideoCodec(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetVideoCodec)(w, r)
}

func (s *ServerInterfaceImpl) SetVideoCodecOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetVideoCodec)(w, r)
}

func (s *ServerInterfaceImpl) SetStreamLatencyLevel(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetStreamLatencyLevel)(w, r)
}

func (s *ServerInterfaceImpl) SetStreamLatencyLevelOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetStreamLatencyLevel)(w, r)
}

func (s *ServerInterfaceImpl) SetStreamOutputVariants(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetStreamOutputVariants)(w, r)
}

func (s *ServerInterfaceImpl) SetStreamOutputVariantsOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetStreamOutputVariants)(w, r)
}

func (s *ServerInterfaceImpl) SetCustomColorVariableValues(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetCustomColorVariableValues)(w, r)
}

func (s *ServerInterfaceImpl) SetCustomColorVariableValuesOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetCustomColorVariableValues)(w, r)
}

func (s *ServerInterfaceImpl) SetLogo(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetLogo)(w, r)
}

func (s *ServerInterfaceImpl) SetLogoOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetLogo)(w, r)
}

func (s *ServerInterfaceImpl) SetFavicon(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetFavicon)(w, r)
}

func (s *ServerInterfaceImpl) SetFaviconOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetFavicon)(w, r)
}

func (s *ServerInterfaceImpl) ResetFavicon(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.ResetFavicon)(w, r)
}

func (s *ServerInterfaceImpl) SetTags(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetTags)(w, r)
}

func (s *ServerInterfaceImpl) SetTagsOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetTags)(w, r)
}

func (s *ServerInterfaceImpl) SetFfmpegPath(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetFfmpegPath)(w, r)
}

func (s *ServerInterfaceImpl) SetFfmpegPathOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetFfmpegPath)(w, r)
}

func (s *ServerInterfaceImpl) SetWebServerPort(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetWebServerPort)(w, r)
}

func (s *ServerInterfaceImpl) SetWebServerPortOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetWebServerPort)(w, r)
}

func (s *ServerInterfaceImpl) SetWebServerIP(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetWebServerIP)(w, r)
}

func (s *ServerInterfaceImpl) SetWebServerIPOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetWebServerIP)(w, r)
}

func (s *ServerInterfaceImpl) SetRTMPServerBindAddress(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetRTMPServerBindAddress)(w, r)
}

func (s *ServerInterfaceImpl) SetRTMPServerPort(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetRTMPServerPort)(w, r)
}

func (s *ServerInterfaceImpl) SetRTMPServerPortOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetRTMPServerPort)(w, r)
}

func (s *ServerInterfaceImpl) SetRTMPServerBindAddressOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetRTMPServerBindAddress)(w, r)
}

func (s *ServerInterfaceImpl) SetSocketHostOverride(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetSocketHostOverride)(w, r)
}

func (s *ServerInterfaceImpl) SetSocketHostOverrideOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetSocketHostOverride)(w, r)
}

func (s *ServerInterfaceImpl) SetVideoServingEndpoint(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetVideoServingEndpoint)(w, r)
}

func (s *ServerInterfaceImpl) SetVideoServingEndpointOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetVideoServingEndpoint)(w, r)
}

func (s *ServerInterfaceImpl) SetNSFW(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetNSFW)(w, r)
}

func (s *ServerInterfaceImpl) SetNSFWOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetNSFW)(w, r)
}

func (s *ServerInterfaceImpl) SetDirectoryEnabled(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetDirectoryEnabled)(w, r)
}

func (s *ServerInterfaceImpl) SetDirectoryEnabledOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetDirectoryEnabled)(w, r)
}

func (s *ServerInterfaceImpl) SetSocialHandles(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetSocialHandles)(w, r)
}

func (s *ServerInterfaceImpl) SetSocialHandlesOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetSocialHandles)(w, r)
}

func (s *ServerInterfaceImpl) SetS3Configuration(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetS3Configuration)(w, r)
}

func (s *ServerInterfaceImpl) SetS3ConfigurationOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetS3Configuration)(w, r)
}

func (s *ServerInterfaceImpl) SetServerURL(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetServerURL)(w, r)
}

func (s *ServerInterfaceImpl) SetServerURLOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetServerURL)(w, r)
}

func (s *ServerInterfaceImpl) SetExternalActions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetExternalActions)(w, r)
}

func (s *ServerInterfaceImpl) SetExternalActionsOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetExternalActions)(w, r)
}

func (s *ServerInterfaceImpl) SetCustomStyles(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetCustomStyles)(w, r)
}

func (s *ServerInterfaceImpl) SetCustomStylesOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetCustomStyles)(w, r)
}

func (s *ServerInterfaceImpl) SetCustomJavascript(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetCustomJavascript)(w, r)
}

func (s *ServerInterfaceImpl) SetCustomJavascriptOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetCustomJavascript)(w, r)
}

func (s *ServerInterfaceImpl) SetHideViewerCount(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetHideViewerCount)(w, r)
}

func (s *ServerInterfaceImpl) SetHideViewerCountOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetHideViewerCount)(w, r)
}

func (s *ServerInterfaceImpl) SetDisableSearchIndexing(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetDisableSearchIndexing)(w, r)
}

func (s *ServerInterfaceImpl) SetDisableSearchIndexingOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetDisableSearchIndexing)(w, r)
}

func (s *ServerInterfaceImpl) SetFederationEnabled(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetFederationEnabled)(w, r)
}

func (s *ServerInterfaceImpl) SetFederationEnabledOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetFederationEnabled)(w, r)
}

func (s *ServerInterfaceImpl) SetFederationActivityPrivate(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetFederationActivityPrivate)(w, r)
}

func (s *ServerInterfaceImpl) SetFederationActivityPrivateOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetFederationActivityPrivate)(w, r)
}

func (s *ServerInterfaceImpl) SetFederationShowEngagement(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetFederationShowEngagement)(w, r)
}

func (s *ServerInterfaceImpl) SetFederationShowEngagementOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetFederationShowEngagement)(w, r)
}

func (s *ServerInterfaceImpl) SetFederationUsername(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetFederationUsername)(w, r)
}

func (s *ServerInterfaceImpl) SetFederationUsernameOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetFederationUsername)(w, r)
}

func (s *ServerInterfaceImpl) SetFederationGoLiveMessage(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetFederationGoLiveMessage)(w, r)
}

func (s *ServerInterfaceImpl) SetFederationGoLiveMessageOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetFederationGoLiveMessage)(w, r)
}

func (s *ServerInterfaceImpl) SetFederationBlockDomains(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetFederationBlockDomains)(w, r)
}

func (s *ServerInterfaceImpl) SetFederationBlockDomainsOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetFederationBlockDomains)(w, r)
}

func (s *ServerInterfaceImpl) SetDiscordNotificationConfiguration(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetDiscordNotificationConfiguration)(w, r)
}

func (s *ServerInterfaceImpl) SetDiscordNotificationConfigurationOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetDiscordNotificationConfiguration)(w, r)
}

func (s *ServerInterfaceImpl) SetBrowserNotificationConfiguration(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetBrowserNotificationConfiguration)(w, r)
}

func (s *ServerInterfaceImpl) SetBrowserNotificationConfigurationOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SetBrowserNotificationConfiguration)(w, r)
}
