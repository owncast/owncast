package handler

import (
	"net/http"

	"github.com/owncast/owncast/controllers/admin"
	"github.com/owncast/owncast/handler/generated"
	"github.com/owncast/owncast/router/middleware"
)

func (*ServerInterfaceImpl) GetAdminStatus(w http.ResponseWriter, r *http.Request, params generated.GetAdminStatusParams) {
	middleware.RequireAdminAuth(admin.Status)(w, r)
}
