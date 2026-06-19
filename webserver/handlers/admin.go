package handlers

import (
	"net/http"

	"github.com/owncast/owncast/webserver/handlers/admin"
	"github.com/owncast/owncast/webserver/handlers/generated"
	"github.com/owncast/owncast/webserver/router/middleware"
)

func (s *ServerInterfaceImpl) StatusAdmin(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.Status)(w, r)
}

func (s *ServerInterfaceImpl) StatusAdminOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.Status)(w, r)
}

func (s *ServerInterfaceImpl) DisconnectInboundConnection(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.DisconnectInboundConnection)(w, r)
}

func (s *ServerInterfaceImpl) DisconnectInboundConnectionOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.DisconnectInboundConnection)(w, r)
}

func (s *ServerInterfaceImpl) GetServerConfig(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetServerConfig)(w, r)
}

func (s *ServerInterfaceImpl) GetServerConfigOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetServerConfig)(w, r)
}

func (s *ServerInterfaceImpl) GetViewersOverTime(w http.ResponseWriter, r *http.Request, params generated.GetViewersOverTimeParams) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetViewersOverTime)(w, r)
}

func (s *ServerInterfaceImpl) GetViewersOverTimeOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetViewersOverTime)(w, r)
}

func (s *ServerInterfaceImpl) GetActiveViewers(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetActiveViewers)(w, r)
}

func (s *ServerInterfaceImpl) GetActiveViewersOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetActiveViewers)(w, r)
}

func (s *ServerInterfaceImpl) GetHardwareStats(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetHardwareStats)(w, r)
}

func (s *ServerInterfaceImpl) GetHardwareStatsOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetHardwareStats)(w, r)
}

func (s *ServerInterfaceImpl) GetConnectedChatClients(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetConnectedChatClients)(w, r)
}

func (s *ServerInterfaceImpl) GetConnectedChatClientsOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetConnectedChatClients)(w, r)
}

func (s *ServerInterfaceImpl) GetChatMessagesAdmin(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetChatMessages)(w, r)
}

func (s *ServerInterfaceImpl) GetChatMessagesAdminOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetChatMessages)(w, r)
}

func (s *ServerInterfaceImpl) UpdateMessageVisibilityAdmin(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.UpdateMessageVisibility)(w, r)
}

func (s *ServerInterfaceImpl) UpdateMessageVisibilityAdminOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.UpdateMessageVisibility)(w, r)
}

func (s *ServerInterfaceImpl) UpdateUserEnabledAdmin(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.UpdateUserEnabled)(w, r)
}

func (s *ServerInterfaceImpl) UpdateUserEnabledAdminOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.UpdateUserEnabled)(w, r)
}

func (s *ServerInterfaceImpl) GetDisabledUsers(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetDisabledUsers)(w, r)
}

func (s *ServerInterfaceImpl) GetDisabledUsersOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetDisabledUsers)(w, r)
}

func (s *ServerInterfaceImpl) BanIPAddress(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.BanIPAddress)(w, r)
}

func (s *ServerInterfaceImpl) BanIPAddressOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.BanIPAddress)(w, r)
}

func (s *ServerInterfaceImpl) UnbanIPAddress(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.UnBanIPAddress)(w, r)
}

func (s *ServerInterfaceImpl) UnbanIPAddressOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.UnBanIPAddress)(w, r)
}

func (s *ServerInterfaceImpl) GetIPAddressBans(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetIPAddressBans)(w, r)
}

func (s *ServerInterfaceImpl) GetIPAddressBansOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetIPAddressBans)(w, r)
}

func (s *ServerInterfaceImpl) UpdateUserModerator(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.UpdateUserModerator)(w, r)
}

func (s *ServerInterfaceImpl) UpdateUserModeratorOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.UpdateUserModerator)(w, r)
}

func (s *ServerInterfaceImpl) GetModerators(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetModerators)(w, r)
}

func (s *ServerInterfaceImpl) GetModeratorsOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetModerators)(w, r)
}

func (s *ServerInterfaceImpl) GetLogs(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(admin.GetLogs)(w, r)
}

func (s *ServerInterfaceImpl) GetLogsOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(admin.GetLogs)(w, r)
}

func (s *ServerInterfaceImpl) GetWarnings(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(admin.GetWarnings)(w, r)
}

func (s *ServerInterfaceImpl) GetWarningsOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(admin.GetWarnings)(w, r)
}

func (s *ServerInterfaceImpl) GetFollowersAdmin(w http.ResponseWriter, r *http.Request, params generated.GetFollowersAdminParams) {
	s.h.middleware.RequireAdminAuth(middleware.HandlePagination(s.h.GetFollowers))(w, r)
}

func (s *ServerInterfaceImpl) GetFollowersAdminOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(middleware.HandlePagination(s.h.GetFollowers))(w, r)
}

func (s *ServerInterfaceImpl) GetPendingFollowRequests(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetPendingFollowRequests)(w, r)
}

func (s *ServerInterfaceImpl) GetPendingFollowRequestsOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetPendingFollowRequests)(w, r)
}

func (s *ServerInterfaceImpl) GetBlockedAndRejectedFollowers(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetBlockedAndRejectedFollowers)(w, r)
}

func (s *ServerInterfaceImpl) GetBlockedAndRejectedFollowersOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetBlockedAndRejectedFollowers)(w, r)
}

func (s *ServerInterfaceImpl) GetDirectoryFollowers(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetDirectoryFollowers)(w, r)
}

func (s *ServerInterfaceImpl) GetDirectoryFollowersOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetDirectoryFollowers)(w, r)
}

func (s *ServerInterfaceImpl) ApproveFollower(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.ApproveFollower)(w, r)
}

func (s *ServerInterfaceImpl) ApproveFollowerOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.ApproveFollower)(w, r)
}

func (s *ServerInterfaceImpl) RemoveFollower(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.RemoveFollower)(w, r)
}

func (s *ServerInterfaceImpl) RemoveFollowerOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.RemoveFollower)(w, r)
}

func (s *ServerInterfaceImpl) UploadCustomEmoji(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(admin.UploadCustomEmoji)(w, r)
}

func (s *ServerInterfaceImpl) UploadCustomEmojiOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(admin.UploadCustomEmoji)(w, r)
}

func (s *ServerInterfaceImpl) DeleteCustomEmoji(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(admin.DeleteCustomEmoji)(w, r)
}

func (s *ServerInterfaceImpl) DeleteCustomEmojiOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(admin.DeleteCustomEmoji)(w, r)
}

func (s *ServerInterfaceImpl) GetWebhooks(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetWebhooks)(w, r)
}

func (s *ServerInterfaceImpl) GetWebhooksOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetWebhooks)(w, r)
}

func (s *ServerInterfaceImpl) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.DeleteWebhook)(w, r)
}

func (s *ServerInterfaceImpl) DeleteWebhookOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.DeleteWebhook)(w, r)
}

func (s *ServerInterfaceImpl) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.CreateWebhook)(w, r)
}

func (s *ServerInterfaceImpl) CreateWebhookOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.CreateWebhook)(w, r)
}

func (s *ServerInterfaceImpl) GetExternalAPIUsers(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetExternalAPIUsers)(w, r)
}

func (s *ServerInterfaceImpl) GetExternalAPIUsersOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetExternalAPIUsers)(w, r)
}

func (s *ServerInterfaceImpl) DeleteExternalAPIUser(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.DeleteExternalAPIUser)(w, r)
}

func (s *ServerInterfaceImpl) DeleteExternalAPIUserOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.DeleteExternalAPIUser)(w, r)
}

func (s *ServerInterfaceImpl) CreateExternalAPIUser(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.CreateExternalAPIUser)(w, r)
}

func (s *ServerInterfaceImpl) CreateExternalAPIUserOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.CreateExternalAPIUser)(w, r)
}

func (s *ServerInterfaceImpl) AutoUpdateOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(admin.AutoUpdateOptions)(w, r)
}

func (s *ServerInterfaceImpl) AutoUpdateOptionsOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(admin.AutoUpdateOptions)(w, r)
}

func (s *ServerInterfaceImpl) AutoUpdateStart(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(admin.AutoUpdateStart)(w, r)
}

func (s *ServerInterfaceImpl) AutoUpdateStartOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(admin.AutoUpdateStart)(w, r)
}

func (s *ServerInterfaceImpl) AutoUpdateForceQuit(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(admin.AutoUpdateForceQuit)(w, r)
}

func (s *ServerInterfaceImpl) AutoUpdateForceQuitOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(admin.AutoUpdateForceQuit)(w, r)
}

func (s *ServerInterfaceImpl) ResetYPRegistration(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.ResetYPRegistration)(w, r)
}

func (s *ServerInterfaceImpl) ResetYPRegistrationOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.ResetYPRegistration)(w, r)
}

func (s *ServerInterfaceImpl) GetVideoPlaybackMetrics(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetVideoPlaybackMetrics)(w, r)
}

func (s *ServerInterfaceImpl) GetVideoPlaybackMetricsOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.GetVideoPlaybackMetrics)(w, r)
}

func (s *ServerInterfaceImpl) SendFederatedMessage(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SendFederatedMessage)(w, r)
}

func (s *ServerInterfaceImpl) SendFederatedMessageOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.SendFederatedMessage)(w, r)
}

func (s *ServerInterfaceImpl) GetFederatedActions(w http.ResponseWriter, r *http.Request, params generated.GetFederatedActionsParams) {
	s.h.middleware.RequireAdminAuth(middleware.HandlePagination(s.h.admin.GetFederatedActions))(w, r)
}

func (s *ServerInterfaceImpl) GetFederatedActionsOptions(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(middleware.HandlePagination(s.h.admin.GetFederatedActions))(w, r)
}
