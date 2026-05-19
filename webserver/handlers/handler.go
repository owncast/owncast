package handlers

import (
	"net/http"

	"github.com/owncast/owncast/webserver/handlers/generated"
	"github.com/owncast/owncast/webserver/router/middleware"

	"github.com/owncast/owncast/yp"
)

// ServerInterfaceImpl is the OpenAPI-generated ServerInterface
// implementation. It holds the dependency-bearing handler sets and
// delegates each generated method to either a free function in this
// package or a method on one of those sets.
//
// As more handlers migrate to needing injected services, their
// delegations here switch from `pkg.X(w, r)` to `s.h.X(w, r)` or
// `s.h.admin.X(w, r)`.
type ServerInterfaceImpl struct {
	h *Handlers
}

// ensure ServerInterfaceImpl implements ServerInterface.
var _ generated.ServerInterface = &ServerInterfaceImpl{}

func New(h *Handlers) *ServerInterfaceImpl {
	return &ServerInterfaceImpl{h: h}
}

func (s *ServerInterfaceImpl) Handler() http.Handler {
	return generated.Handler(s)
}

func (s *ServerInterfaceImpl) GetStatus(w http.ResponseWriter, r *http.Request) {
	s.h.GetStatus(w, r)
}

func (*ServerInterfaceImpl) GetCustomEmojiList(w http.ResponseWriter, r *http.Request) {
	GetCustomEmojiList(w, r)
}

func (*ServerInterfaceImpl) GetChatMessages(w http.ResponseWriter, r *http.Request, params generated.GetChatMessagesParams) {
	middleware.RequireUserAccessToken(GetChatMessages)(w, r)
}

func (*ServerInterfaceImpl) RegisterAnonymousChatUser(w http.ResponseWriter, r *http.Request, params generated.RegisterAnonymousChatUserParams) {
	RegisterAnonymousChatUser(w, r)
}

func (*ServerInterfaceImpl) RegisterAnonymousChatUserOptions(w http.ResponseWriter, r *http.Request) {
	RegisterAnonymousChatUser(w, r)
}

func (s *ServerInterfaceImpl) UpdateMessageVisibility(w http.ResponseWriter, r *http.Request, params generated.UpdateMessageVisibilityParams) {
	middleware.RequireUserModerationScopeAccesstoken(s.h.admin.UpdateMessageVisibility)(w, r)
}

func (s *ServerInterfaceImpl) UpdateUserEnabled(w http.ResponseWriter, r *http.Request, params generated.UpdateUserEnabledParams) {
	middleware.RequireUserModerationScopeAccesstoken(s.h.admin.UpdateUserEnabled)(w, r)
}

func (s *ServerInterfaceImpl) GetWebConfig(w http.ResponseWriter, r *http.Request) {
	s.h.GetWebConfig(w, r)
}

func (*ServerInterfaceImpl) GetYPResponse(w http.ResponseWriter, r *http.Request) {
	yp.GetYPResponse(w, r)
}

func (*ServerInterfaceImpl) GetAllSocialPlatforms(w http.ResponseWriter, r *http.Request) {
	GetAllSocialPlatforms(w, r)
}

func (*ServerInterfaceImpl) GetVideoStreamOutputVariants(w http.ResponseWriter, r *http.Request) {
	GetVideoStreamOutputVariants(w, r)
}

func (s *ServerInterfaceImpl) Ping(w http.ResponseWriter, r *http.Request) {
	s.h.Ping(w, r)
}

func (*ServerInterfaceImpl) RemoteFollow(w http.ResponseWriter, r *http.Request) {
	RemoteFollow(w, r)
}

func (*ServerInterfaceImpl) GetFollowers(w http.ResponseWriter, r *http.Request, params generated.GetFollowersParams) {
	middleware.HandlePagination(GetFollowers)(w, r)
}

func (*ServerInterfaceImpl) ReportPlaybackMetrics(w http.ResponseWriter, r *http.Request) {
	ReportPlaybackMetrics(w, r)
}

func (*ServerInterfaceImpl) RegisterForLiveNotifications(w http.ResponseWriter, r *http.Request, params generated.RegisterForLiveNotificationsParams) {
	middleware.RequireUserAccessToken(RegisterForLiveNotifications)(w, r)
}
