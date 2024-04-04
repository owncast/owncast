package handler

import (
	"net/http"

	"github.com/owncast/owncast/controllers/admin"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/router/middleware"
)

func (*ServerInterfaceImpl) SendSystemMessage(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(user.ScopeCanSendSystemMessages, admin.SendSystemMessage)(w, r)
}
