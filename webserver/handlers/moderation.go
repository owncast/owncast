package handlers

import (
	"net/http"

	"github.com/owncast/owncast/webserver/handlers/generated"
)

func (s *ServerInterfaceImpl) GetUserDetails(w http.ResponseWriter, r *http.Request, userId string, params generated.GetUserDetailsParams) {
	s.h.middleware.RequireUserModerationScopeAccesstoken(s.h.moderation.GetUserDetails)(w, r)
}
