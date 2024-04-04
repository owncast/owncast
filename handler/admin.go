package handler

import (
	"net/http"

	"github.com/owncast/owncast/controllers/admin"
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
