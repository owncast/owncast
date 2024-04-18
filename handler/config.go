package handler

import (
	"net/http"

	"github.com/owncast/owncast/controllers/admin"
	"github.com/owncast/owncast/router/middleware"
)

func (*ServerInterfaceImpl) SetAdminPassword(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetAdminPassword)(w, r)
}

func (*ServerInterfaceImpl) SetStreamKeys(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetStreamKeys)(w, r)
}

func (*ServerInterfaceImpl) SetExtraPageContent(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetExtraPageContent)(w, r)
}

func (*ServerInterfaceImpl) SetStreamTitle(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetStreamTitle)(w, r)
}

func (*ServerInterfaceImpl) SetServerName(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetServerName)(w, r)
}

func (*ServerInterfaceImpl) SetServerSummary(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetServerSummary)(w, r)
}

func (*ServerInterfaceImpl) SetCustomOfflineMessage(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetCustomOfflineMessage)(w, r)
}

func (*ServerInterfaceImpl) SetServerWelcomeMessage(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetServerWelcomeMessage)(w, r)
}

func (*ServerInterfaceImpl) SetChatDisabled(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetChatDisabled)(w, r)
}

func (*ServerInterfaceImpl) SetChatJoinMessagesEnabled(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetChatJoinMessagesEnabled)(w, r)
}

func (*ServerInterfaceImpl) SetEnableEstablishedChatUserMode(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetEnableEstablishedChatUserMode)(w, r)
}

func (*ServerInterfaceImpl) SetForbiddenUsernameList(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetForbiddenUsernameList)(w, r)
}

func (*ServerInterfaceImpl) SetSuggestedUsernameList(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetSuggestedUsernameList)(w, r)
}

func (*ServerInterfaceImpl) SetVideoCodec(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetVideoCodec)(w, r)
}

func (*ServerInterfaceImpl) SetStreamLatencyLevel(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetStreamLatencyLevel)(w, r)
}

func (*ServerInterfaceImpl) SetStreamOutputVariants(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetStreamOutputVariants)(w, r)
}

func (*ServerInterfaceImpl) SetCustomColorVariableValues(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetCustomColorVariableValues)(w, r)
}

func (*ServerInterfaceImpl) SetLogo(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetLogo)(w, r)
}

func (*ServerInterfaceImpl) SetTags(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetTags)(w, r)
}

func (*ServerInterfaceImpl) SetFfmpegPath(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetFfmpegPath)(w, r)
}

func (*ServerInterfaceImpl) SetWebServerPort(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetWebServerPort)(w, r)
}

func (*ServerInterfaceImpl) SetWebServerIP(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetWebServerIP)(w, r)
}

func (*ServerInterfaceImpl) SetRTMPServerPort(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetRTMPServerPort)(w, r)
}

func (*ServerInterfaceImpl) SetSocketHostOverride(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetSocketHostOverride)(w, r)
}

func (*ServerInterfaceImpl) SetVideoServingEndpoint(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetVideoServingEndpoint)(w, r)
}

func (*ServerInterfaceImpl) SetNSFW(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetNSFW)(w, r)
}

func (*ServerInterfaceImpl) SetDirectoryEnabled(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetDirectoryEnabled)(w, r)
}

func (*ServerInterfaceImpl) SetSocialHandles(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetSocialHandles)(w, r)
}

func (*ServerInterfaceImpl) SetS3Configuration(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetS3Configuration)(w, r)
}

func (*ServerInterfaceImpl) SetServerURL(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetServerURL)(w, r)
}

func (*ServerInterfaceImpl) SetExternalActions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetExternalActions)(w, r)
}

func (*ServerInterfaceImpl) SetCustomStyles(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetCustomStyles)(w, r)
}

func (*ServerInterfaceImpl) SetCustomJavascript(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetCustomJavascript)(w, r)
}

func (*ServerInterfaceImpl) SetHideViewerCount(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetHideViewerCount)(w, r)
}

func (*ServerInterfaceImpl) SetDisableSearchIndexing(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetDisableSearchIndexing)(w, r)
}

func (*ServerInterfaceImpl) SetFederationEnabled(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetFederationEnabled)(w, r)
}

func (*ServerInterfaceImpl) SetFederationActivityPrivate(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SetFederationActivityPrivate)(w, r)
}
