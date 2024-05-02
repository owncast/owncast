package handler

import (
	"net/http"

	"github.com/owncast/owncast/controllers/moderation"
	"github.com/owncast/owncast/handler/generated"
	"github.com/owncast/owncast/router/middleware"
)

func (*ServerInterfaceImpl) GetUserDetails(w http.ResponseWriter, r *http.Request, userId string, params generated.GetUserDetailsParams) {
	middleware.RequireUserModerationScopeAccesstoken(moderation.GetUserDetails)(w, r)
}
