package handlers

import (
	"net/http"

	"github.com/owncast/owncast/core/rtmp"
	"github.com/owncast/owncast/webserver/handlers/admin"
	"github.com/owncast/owncast/webserver/handlers/generated"
	"github.com/owncast/owncast/webserver/router/middleware"
)

func (*ServerInterfaceImpl) StatusAdmin(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.Status)(w, r)
}

func (*ServerInterfaceImpl) StatusAdminOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.Status)(w, r)
}

func (*ServerInterfaceImpl) DisconnectInboundConnection(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.DisconnectInboundConnection)(w, r)
}

func (*ServerInterfaceImpl) DisconnectInboundConnectionOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.DisconnectInboundConnection)(w, r)
}

func (*ServerInterfaceImpl) GetServerConfig(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetServerConfig)(w, r)
}

func (*ServerInterfaceImpl) GetServerConfigOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetServerConfig)(w, r)
}

func (*ServerInterfaceImpl) GetViewersOverTime(w http.ResponseWriter, r *http.Request, params generated.GetViewersOverTimeParams) {
	middleware.RequireAdminAuth(admin.GetViewersOverTime)(w, r)
}

func (*ServerInterfaceImpl) GetViewersOverTimeOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetViewersOverTime)(w, r)
}

func (*ServerInterfaceImpl) GetActiveViewers(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetActiveViewers)(w, r)
}

func (*ServerInterfaceImpl) GetActiveViewersOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetActiveViewers)(w, r)
}

func (*ServerInterfaceImpl) GetHardwareStats(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetHardwareStats)(w, r)
}

func (*ServerInterfaceImpl) GetHardwareStatsOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetHardwareStats)(w, r)
}

func (*ServerInterfaceImpl) GetConnectedChatClients(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetConnectedChatClients)(w, r)
}

func (*ServerInterfaceImpl) GetConnectedChatClientsOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetConnectedChatClients)(w, r)
}

func (*ServerInterfaceImpl) GetChatMessagesAdmin(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetChatMessages)(w, r)
}

func (*ServerInterfaceImpl) GetChatMessagesAdminOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetChatMessages)(w, r)
}

func (*ServerInterfaceImpl) UpdateMessageVisibilityAdmin(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.UpdateMessageVisibility)(w, r)
}

func (*ServerInterfaceImpl) UpdateMessageVisibilityAdminOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.UpdateMessageVisibility)(w, r)
}

func (*ServerInterfaceImpl) UpdateUserEnabledAdmin(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.UpdateUserEnabled)(w, r)
}

func (*ServerInterfaceImpl) UpdateUserEnabledAdminOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.UpdateUserEnabled)(w, r)
}

func (*ServerInterfaceImpl) GetDisabledUsers(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetDisabledUsers)(w, r)
}

func (*ServerInterfaceImpl) GetDisabledUsersOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetDisabledUsers)(w, r)
}

func (*ServerInterfaceImpl) BanIPAddress(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.BanIPAddress)(w, r)
}

func (*ServerInterfaceImpl) BanIPAddressOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.BanIPAddress)(w, r)
}

func (*ServerInterfaceImpl) UnbanIPAddress(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.UnBanIPAddress)(w, r)
}

func (*ServerInterfaceImpl) UnbanIPAddressOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.UnBanIPAddress)(w, r)
}

func (*ServerInterfaceImpl) GetIPAddressBans(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetIPAddressBans)(w, r)
}

func (*ServerInterfaceImpl) GetIPAddressBansOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetIPAddressBans)(w, r)
}

func (*ServerInterfaceImpl) UpdateUserModerator(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.UpdateUserModerator)(w, r)
}

func (*ServerInterfaceImpl) UpdateUserModeratorOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.UpdateUserModerator)(w, r)
}

func (*ServerInterfaceImpl) GetModerators(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetModerators)(w, r)
}

func (*ServerInterfaceImpl) GetModeratorsOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetModerators)(w, r)
}

func (*ServerInterfaceImpl) GetLogs(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetLogs)(w, r)
}

func (*ServerInterfaceImpl) GetLogsOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetLogs)(w, r)
}

func (*ServerInterfaceImpl) GetWarnings(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetWarnings)(w, r)
}

func (*ServerInterfaceImpl) GetWarningsOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetWarnings)(w, r)
}

func (*ServerInterfaceImpl) GetFollowersAdmin(w http.ResponseWriter, r *http.Request, params generated.GetFollowersAdminParams) {
	middleware.RequireAdminAuth(middleware.HandlePagination(GetFollowers))(w, r)
}

func (*ServerInterfaceImpl) GetFollowersAdminOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(middleware.HandlePagination(GetFollowers))(w, r)
}

func (*ServerInterfaceImpl) GetPendingFollowRequests(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetPendingFollowRequests)(w, r)
}

func (*ServerInterfaceImpl) GetPendingFollowRequestsOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetPendingFollowRequests)(w, r)
}

func (*ServerInterfaceImpl) GetBlockedAndRejectedFollowers(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetBlockedAndRejectedFollowers)(w, r)
}

func (*ServerInterfaceImpl) GetBlockedAndRejectedFollowersOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetBlockedAndRejectedFollowers)(w, r)
}

func (*ServerInterfaceImpl) ApproveFollower(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.ApproveFollower)(w, r)
}

func (*ServerInterfaceImpl) ApproveFollowerOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.ApproveFollower)(w, r)
}

func (*ServerInterfaceImpl) UploadCustomEmoji(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.UploadCustomEmoji)(w, r)
}

func (*ServerInterfaceImpl) UploadCustomEmojiOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.UploadCustomEmoji)(w, r)
}

func (*ServerInterfaceImpl) DeleteCustomEmoji(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.DeleteCustomEmoji)(w, r)
}

func (*ServerInterfaceImpl) DeleteCustomEmojiOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.DeleteCustomEmoji)(w, r)
}

func (*ServerInterfaceImpl) GetWebhooks(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetWebhooks)(w, r)
}

func (*ServerInterfaceImpl) GetWebhooksOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetWebhooks)(w, r)
}

func (*ServerInterfaceImpl) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.DeleteWebhook)(w, r)
}

func (*ServerInterfaceImpl) DeleteWebhookOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.DeleteWebhook)(w, r)
}

func (*ServerInterfaceImpl) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.CreateWebhook)(w, r)
}

func (*ServerInterfaceImpl) CreateWebhookOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.CreateWebhook)(w, r)
}

func (*ServerInterfaceImpl) GetExternalAPIUsers(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetExternalAPIUsers)(w, r)
}

func (*ServerInterfaceImpl) GetExternalAPIUsersOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetExternalAPIUsers)(w, r)
}

func (*ServerInterfaceImpl) DeleteExternalAPIUser(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.DeleteExternalAPIUser)(w, r)
}

func (*ServerInterfaceImpl) DeleteExternalAPIUserOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.DeleteExternalAPIUser)(w, r)
}

func (*ServerInterfaceImpl) CreateExternalAPIUser(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.CreateExternalAPIUser)(w, r)
}

func (*ServerInterfaceImpl) CreateExternalAPIUserOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.CreateExternalAPIUser)(w, r)
}

func (*ServerInterfaceImpl) AutoUpdateOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.AutoUpdateOptions)(w, r)
}

func (*ServerInterfaceImpl) AutoUpdateOptionsOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.AutoUpdateOptions)(w, r)
}

func (*ServerInterfaceImpl) AutoUpdateStart(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.AutoUpdateStart)(w, r)
}

func (*ServerInterfaceImpl) AutoUpdateStartOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.AutoUpdateStart)(w, r)
}

func (*ServerInterfaceImpl) AutoUpdateForceQuit(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.AutoUpdateForceQuit)(w, r)
}

func (*ServerInterfaceImpl) AutoUpdateForceQuitOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.AutoUpdateForceQuit)(w, r)
}

func (*ServerInterfaceImpl) ResetYPRegistration(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.ResetYPRegistration)(w, r)
}

func (*ServerInterfaceImpl) ResetYPRegistrationOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.ResetYPRegistration)(w, r)
}

func (*ServerInterfaceImpl) GetVideoPlaybackMetrics(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetVideoPlaybackMetrics)(w, r)
}

func (*ServerInterfaceImpl) GetVideoPlaybackMetricsOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.GetVideoPlaybackMetrics)(w, r)
}

func (*ServerInterfaceImpl) SendFederatedMessage(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SendFederatedMessage)(w, r)
}

func (*ServerInterfaceImpl) SendFederatedMessageOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.SendFederatedMessage)(w, r)
}

func (*ServerInterfaceImpl) GetFederatedActions(w http.ResponseWriter, r *http.Request, params generated.GetFederatedActionsParams) {
	middleware.RequireAdminAuth(middleware.HandlePagination(admin.GetFederatedActions))(w, r)
}

func (*ServerInterfaceImpl) GetFederatedActionsOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(middleware.HandlePagination(admin.GetFederatedActions))(w, r)
}

// DisconnectInboundConnection will force-disconnect an inbound stream.
func DisconnectInboundConnection(w http.ResponseWriter, r *http.Request) {
	rtmp.Disconnect()
	w.WriteHeader(http.StatusOK)
}
