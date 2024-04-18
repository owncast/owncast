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

func (*ServerInterfaceImpl) UpdateMessageVisibilityAdmin(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.UpdateMessageVisibility)(w, r)
}

func (*ServerInterfaceImpl) UpdateUserEnabledAdmin(w http.ResponseWriter, r *http.Request) {
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

func (*ServerInterfaceImpl) ResetYPRegistration(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.ResetYPRegistration)(w, r)
}

func (*ServerInterfaceImpl) GetVideoPlaybackMetrics(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetVideoPlaybackMetrics)(w, r)
}

func (*ServerInterfaceImpl) SendFederatedMessage(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SendFederatedMessage)(w, r)
}
