package handler

import (
	"net/http"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/controllers/admin"
	"github.com/owncast/owncast/handler/generated"
	"github.com/owncast/owncast/router/middleware"
)

func (*ServerInterfaceImpl) GetAdminStatus(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.Status)(w, r)
}

func (*ServerInterfaceImpl) DisconnectInboundConnection(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.DisconnectInboundConnection)(w, r)
}

func (*ServerInterfaceImpl) GetServerConfig(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetServerConfig)(w, r)
}

func (*ServerInterfaceImpl) GetViewersOverTime(w http.ResponseWriter, r *http.Request, params generated.GetViewersOverTimeParams) {
	middleware.RequireAdminAuth(admin.GetViewersOverTime)(w, r)
}

func (*ServerInterfaceImpl) GetActiveViewers(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetActiveViewers)(w, r)
}

func (*ServerInterfaceImpl) GetHardwareStats(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetHardwareStats)(w, r)
}

func (*ServerInterfaceImpl) GetConnectedChatClients(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetConnectedChatClients)(w, r)
}

func (*ServerInterfaceImpl) GetAdminChatMessages(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetChatMessages)(w, r)
}

func (*ServerInterfaceImpl) GetLogs(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetLogs)(w, r)
}

func (*ServerInterfaceImpl) GetWarnings(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetWarnings)(w, r)
}

func (*ServerInterfaceImpl) UpdateMessageVisibility(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.UpdateMessageVisibility)(w, r)
}

func (*ServerInterfaceImpl) UpdateUserEnabled(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.UpdateUserEnabled)(w, r)
}

func (*ServerInterfaceImpl) BanIPAddress(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.BanIPAddress)(w, r)
}

func (*ServerInterfaceImpl) UnbanIPAddress(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.UnBanIPAddress)(w, r)
}

func (*ServerInterfaceImpl) GetIPAddressBans(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetIPAddressBans)(w, r)
}

func (*ServerInterfaceImpl) GetDisabledUsers(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetDisabledUsers)(w, r)
}

func (*ServerInterfaceImpl) UpdateUserModerator(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.UpdateUserModerator)(w, r)
}

func (*ServerInterfaceImpl) GetModerators(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetModerators)(w, r)
}

func (*ServerInterfaceImpl) GetFollowersAdmin(w http.ResponseWriter, r *http.Request, params generated.GetFollowersAdminParams) {
	// FIXME this calls the same function as `GetFollowers` but with admin auth
	middleware.RequireAdminAuth(middleware.HandlePagination(controllers.GetFollowers))(w, r)
}

func (*ServerInterfaceImpl) GetPendingFollowers(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetPendingFollowRequests)(w, r)
}

func (*ServerInterfaceImpl) GetBlockedAndRejectedFollowers(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetBlockedAndRejectedFollowers)(w, r)
}

func (*ServerInterfaceImpl) ApproveFollower(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.ApproveFollower)(w, r)
}

func (*ServerInterfaceImpl) UploadCustomEmoji(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.UploadCustomEmoji)(w, r)
}

func (*ServerInterfaceImpl) DeleteCustomEmoji(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.DeleteCustomEmoji)(w, r)
}

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

func (*ServerInterfaceImpl) GetWebhooks(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetWebhooks)(w, r)
}

func (*ServerInterfaceImpl) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.DeleteWebhook)(w, r)
}

func (*ServerInterfaceImpl) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.CreateWebhook)(w, r)
}

func (*ServerInterfaceImpl) GetExternalAPIUsers(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetExternalAPIUsers)(w, r)
}

func (*ServerInterfaceImpl) DeleteExternalAPIUser(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.DeleteExternalAPIUser)(w, r)
}

func (*ServerInterfaceImpl) CreateExternalAPIUser(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.CreateExternalAPIUser)(w, r)
}

func (*ServerInterfaceImpl) AutoUpdateOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.AutoUpdateOptions)(w, r)
}

func (*ServerInterfaceImpl) AutoUpdateStart(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.AutoUpdateStart)(w, r)
}

func (*ServerInterfaceImpl) AutoUpdateForceQuit(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.AutoUpdateForceQuit)(w, r)
}
