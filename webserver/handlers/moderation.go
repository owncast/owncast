package handlers

import (
	"net/http"

	"github.com/owncast/owncast/webserver/handlers/generated"
	"github.com/owncast/owncast/webserver/handlers/moderation"
	"github.com/owncast/owncast/webserver/router/middleware"
)

func (*ServerInterfaceImpl) GetUserDetails(w http.ResponseWriter, r *http.Request, userId string, params generated.GetUserDetailsParams) {
	middleware.RequireUserModerationScopeAccesstoken(moderation.GetUserDetails)(w, r)
}
