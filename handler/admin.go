package handler

import (
	"net/http"

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
