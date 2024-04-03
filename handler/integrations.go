package handler

import (
	"net/http"

	"github.com/owncast/owncast/controllers/admin"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/handler/generated"
	"github.com/owncast/owncast/router/middleware"
)

func (*ServerInterfaceImpl) SendSystemMessage(w http.ResponseWriter, r *http.Request, params generated.SendSystemMessageParams) {
	middleware.RequireExternalAPIAccessToken(user.ScopeCanSendSystemMessages, admin.SendSystemMessage)(w, r)
}
