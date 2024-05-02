package handler

import (
	"net/http"

	"github.com/owncast/owncast/controllers/admin"
	"github.com/owncast/owncast/router/middleware"
)

func (*ServerInterfaceImpl) SetAdminPassword(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetAdminPassword)(w, r)
}

func (*ServerInterfaceImpl) SetAdminPasswordOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetAdminPassword)(w, r)
}

func (*ServerInterfaceImpl) SetStreamKeys(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetStreamKeys)(w, r)
}

func (*ServerInterfaceImpl) SetStreamKeysOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetStreamKeys)(w, r)
}

func (*ServerInterfaceImpl) SetExtraPageContent(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetExtraPageContent)(w, r)
}

func (*ServerInterfaceImpl) SetExtraPageContentOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetExtraPageContent)(w, r)
}

func (*ServerInterfaceImpl) SetStreamTitle(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetStreamTitle)(w, r)
}

func (*ServerInterfaceImpl) SetStreamTitleOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetStreamTitle)(w, r)
}

func (*ServerInterfaceImpl) SetServerName(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetServerName)(w, r)
}

func (*ServerInterfaceImpl) SetServerNameOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetServerName)(w, r)
}

func (*ServerInterfaceImpl) SetServerSummary(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetServerSummary)(w, r)
}

func (*ServerInterfaceImpl) SetServerSummaryOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetServerSummary)(w, r)
}

func (*ServerInterfaceImpl) SetCustomOfflineMessage(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetCustomOfflineMessage)(w, r)
}

func (*ServerInterfaceImpl) SetCustomOfflineMessageOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetCustomOfflineMessage)(w, r)
}

func (*ServerInterfaceImpl) SetServerWelcomeMessage(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetServerWelcomeMessage)(w, r)
}

func (*ServerInterfaceImpl) SetServerWelcomeMessageOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetServerWelcomeMessage)(w, r)
}

func (*ServerInterfaceImpl) SetChatDisabled(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetChatDisabled)(w, r)
}

func (*ServerInterfaceImpl) SetChatDisabledOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetChatDisabled)(w, r)
}

func (*ServerInterfaceImpl) SetChatJoinMessagesEnabled(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetChatJoinMessagesEnabled)(w, r)
}

func (*ServerInterfaceImpl) SetChatJoinMessagesEnabledOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetChatJoinMessagesEnabled)(w, r)
}

func (*ServerInterfaceImpl) SetEnableEstablishedChatUserMode(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetEnableEstablishedChatUserMode)(w, r)
}

func (*ServerInterfaceImpl) SetEnableEstablishedChatUserModeOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetEnableEstablishedChatUserMode)(w, r)
}

func (*ServerInterfaceImpl) SetForbiddenUsernameList(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetForbiddenUsernameList)(w, r)
}

func (*ServerInterfaceImpl) SetForbiddenUsernameListOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetForbiddenUsernameList)(w, r)
}

func (*ServerInterfaceImpl) SetSuggestedUsernameList(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetSuggestedUsernameList)(w, r)
}

func (*ServerInterfaceImpl) SetSuggestedUsernameListOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetSuggestedUsernameList)(w, r)
}

func (*ServerInterfaceImpl) SetChatSpamProtectionEnabled(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetChatSpamProtectionEnabled)(w, r)
}

func (*ServerInterfaceImpl) SetChatSpamProtectionEnabledOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetChatSpamProtectionEnabled)(w, r)
}

func (*ServerInterfaceImpl) SetChatSlurFilterEnabled(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetChatSlurFilterEnabled)(w, r)
}

func (*ServerInterfaceImpl) SetChatSlurFilterEnabledOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetChatSlurFilterEnabled)(w, r)
}

func (*ServerInterfaceImpl) SetVideoCodec(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetVideoCodec)(w, r)
}

func (*ServerInterfaceImpl) SetVideoCodecOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetVideoCodec)(w, r)
}

func (*ServerInterfaceImpl) SetStreamLatencyLevel(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetStreamLatencyLevel)(w, r)
}

func (*ServerInterfaceImpl) SetStreamLatencyLevelOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetStreamLatencyLevel)(w, r)
}

func (*ServerInterfaceImpl) SetStreamOutputVariants(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetStreamOutputVariants)(w, r)
}

func (*ServerInterfaceImpl) SetStreamOutputVariantsOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetStreamOutputVariants)(w, r)
}

func (*ServerInterfaceImpl) SetCustomColorVariableValues(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetCustomColorVariableValues)(w, r)
}

func (*ServerInterfaceImpl) SetCustomColorVariableValuesOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetCustomColorVariableValues)(w, r)
}

func (*ServerInterfaceImpl) SetLogo(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetLogo)(w, r)
}

func (*ServerInterfaceImpl) SetLogoOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetLogo)(w, r)
}

func (*ServerInterfaceImpl) SetTags(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetTags)(w, r)
}

func (*ServerInterfaceImpl) SetTagsOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetTags)(w, r)
}

func (*ServerInterfaceImpl) SetFfmpegPath(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetFfmpegPath)(w, r)
}

func (*ServerInterfaceImpl) SetFfmpegPathOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetFfmpegPath)(w, r)
}

func (*ServerInterfaceImpl) SetWebServerPort(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetWebServerPort)(w, r)
}

func (*ServerInterfaceImpl) SetWebServerPortOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetWebServerPort)(w, r)
}

func (*ServerInterfaceImpl) SetWebServerIP(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetWebServerIP)(w, r)
}

func (*ServerInterfaceImpl) SetWebServerIPOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetWebServerIP)(w, r)
}

func (*ServerInterfaceImpl) SetRTMPServerPort(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetRTMPServerPort)(w, r)
}

func (*ServerInterfaceImpl) SetRTMPServerPortOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetRTMPServerPort)(w, r)
}

func (*ServerInterfaceImpl) SetSocketHostOverride(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetSocketHostOverride)(w, r)
}

func (*ServerInterfaceImpl) SetSocketHostOverrideOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetSocketHostOverride)(w, r)
}

func (*ServerInterfaceImpl) SetVideoServingEndpoint(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetVideoServingEndpoint)(w, r)
}

func (*ServerInterfaceImpl) SetVideoServingEndpointOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetVideoServingEndpoint)(w, r)
}

func (*ServerInterfaceImpl) SetNSFW(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetNSFW)(w, r)
}

func (*ServerInterfaceImpl) SetNSFWOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetNSFW)(w, r)
}

func (*ServerInterfaceImpl) SetDirectoryEnabled(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetDirectoryEnabled)(w, r)
}

func (*ServerInterfaceImpl) SetDirectoryEnabledOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetDirectoryEnabled)(w, r)
}

func (*ServerInterfaceImpl) SetSocialHandles(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetSocialHandles)(w, r)
}

func (*ServerInterfaceImpl) SetSocialHandlesOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetSocialHandles)(w, r)
}

func (*ServerInterfaceImpl) SetS3Configuration(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetS3Configuration)(w, r)
}

func (*ServerInterfaceImpl) SetS3ConfigurationOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetS3Configuration)(w, r)
}

func (*ServerInterfaceImpl) SetServerURL(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetServerURL)(w, r)
}

func (*ServerInterfaceImpl) SetServerURLOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetServerURL)(w, r)
}

func (*ServerInterfaceImpl) SetExternalActions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetExternalActions)(w, r)
}

func (*ServerInterfaceImpl) SetExternalActionsOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetExternalActions)(w, r)
}

func (*ServerInterfaceImpl) SetCustomStyles(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetCustomStyles)(w, r)
}

func (*ServerInterfaceImpl) SetCustomStylesOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetCustomStyles)(w, r)
}

func (*ServerInterfaceImpl) SetCustomJavascript(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetCustomJavascript)(w, r)
}

func (*ServerInterfaceImpl) SetCustomJavascriptOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetCustomJavascript)(w, r)
}

func (*ServerInterfaceImpl) SetHideViewerCount(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetHideViewerCount)(w, r)
}

func (*ServerInterfaceImpl) SetHideViewerCountOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetHideViewerCount)(w, r)
}

func (*ServerInterfaceImpl) SetDisableSearchIndexing(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetDisableSearchIndexing)(w, r)
}

func (*ServerInterfaceImpl) SetDisableSearchIndexingOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetDisableSearchIndexing)(w, r)
}

func (*ServerInterfaceImpl) SetFederationEnabled(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetFederationEnabled)(w, r)
}

func (*ServerInterfaceImpl) SetFederationEnabledOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetFederationEnabled)(w, r)
}

func (*ServerInterfaceImpl) SetFederationActivityPrivate(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetFederationActivityPrivate)(w, r)
}

func (*ServerInterfaceImpl) SetFederationActivityPrivateOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetFederationActivityPrivate)(w, r)
}

func (*ServerInterfaceImpl) SetFederationShowEngagement(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetFederationShowEngagement)(w, r)
}

func (*ServerInterfaceImpl) SetFederationShowEngagementOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetFederationShowEngagement)(w, r)
}

func (*ServerInterfaceImpl) SetFederationUsername(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetFederationUsername)(w, r)
}

func (*ServerInterfaceImpl) SetFederationUsernameOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetFederationUsername)(w, r)
}

func (*ServerInterfaceImpl) SetFederationGoLiveMessage(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetFederationGoLiveMessage)(w, r)
}

func (*ServerInterfaceImpl) SetFederationGoLiveMessageOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetFederationGoLiveMessage)(w, r)
}

func (*ServerInterfaceImpl) SetFederationBlockDomains(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetFederationBlockDomains)(w, r)
}

func (*ServerInterfaceImpl) SetFederationBlockDomainsOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetFederationBlockDomains)(w, r)
}

func (*ServerInterfaceImpl) SetDiscordNotificationConfiguration(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetDiscordNotificationConfiguration)(w, r)
}

func (*ServerInterfaceImpl) SetDiscordNotificationConfigurationOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetDiscordNotificationConfiguration)(w, r)
}

func (*ServerInterfaceImpl) SetBrowserNotificationConfiguration(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetBrowserNotificationConfiguration)(w, r)
}

func (*ServerInterfaceImpl) SetBrowserNotificationConfigurationOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetBrowserNotificationConfiguration)(w, r)
}
