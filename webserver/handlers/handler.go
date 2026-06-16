package handlers

import (
	"net/http"

	"github.com/owncast/owncast/webserver/handlers/generated"
	"github.com/owncast/owncast/webserver/router/middleware"
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

func (s *ServerInterfaceImpl) GetChatMessages(w http.ResponseWriter, r *http.Request, params generated.GetChatMessagesParams) {
	s.h.middleware.RequireUserAccessToken(s.h.GetChatMessages)(w, r)
}

func (s *ServerInterfaceImpl) RegisterAnonymousChatUser(w http.ResponseWriter, r *http.Request, params generated.RegisterAnonymousChatUserParams) {
	s.h.RegisterAnonymousChatUser(w, r)
}

func (s *ServerInterfaceImpl) RegisterAnonymousChatUserOptions(w http.ResponseWriter, r *http.Request) {
	s.h.RegisterAnonymousChatUser(w, r)
}

func (s *ServerInterfaceImpl) UpdateMessageVisibility(w http.ResponseWriter, r *http.Request, params generated.UpdateMessageVisibilityParams) {
	s.h.middleware.RequireUserModerationScopeAccesstoken(s.h.admin.UpdateMessageVisibility)(w, r)
}

func (s *ServerInterfaceImpl) UpdateUserEnabled(w http.ResponseWriter, r *http.Request, params generated.UpdateUserEnabledParams) {
	s.h.middleware.RequireUserModerationScopeAccesstoken(s.h.admin.UpdateUserEnabled)(w, r)
}

func (s *ServerInterfaceImpl) GetWebConfig(w http.ResponseWriter, r *http.Request) {
	s.h.GetWebConfig(w, r)
}

func (s *ServerInterfaceImpl) GetYPResponse(w http.ResponseWriter, r *http.Request) {
	s.h.yp.GetYPResponse(w, r)
}

func (*ServerInterfaceImpl) GetAllSocialPlatforms(w http.ResponseWriter, r *http.Request) {
	GetAllSocialPlatforms(w, r)
}

func (s *ServerInterfaceImpl) GetVideoStreamOutputVariants(w http.ResponseWriter, r *http.Request) {
	s.h.GetVideoStreamOutputVariants(w, r)
}

func (s *ServerInterfaceImpl) Ping(w http.ResponseWriter, r *http.Request) {
	s.h.Ping(w, r)
}

func (s *ServerInterfaceImpl) RemoteFollow(w http.ResponseWriter, r *http.Request) {
	s.h.RemoteFollow(w, r)
}

func (s *ServerInterfaceImpl) GetFollowers(w http.ResponseWriter, r *http.Request, params generated.GetFollowersParams) {
	middleware.HandlePagination(s.h.GetFollowers)(w, r)
}

func (s *ServerInterfaceImpl) ReportPlaybackMetrics(w http.ResponseWriter, r *http.Request) {
	s.h.ReportPlaybackMetrics(w, r)
}

func (s *ServerInterfaceImpl) RegisterForLiveNotifications(w http.ResponseWriter, r *http.Request, params generated.RegisterForLiveNotificationsParams) {
	s.h.middleware.RequireUserAccessToken(s.h.RegisterForLiveNotifications)(w, r)
}

// Federated servers endpoints

func (s *ServerInterfaceImpl) GetFederatedServers(w http.ResponseWriter, r *http.Request) {
	s.h.admin.GetFederatedServers(w, r)
}

func (s *ServerInterfaceImpl) AddFederatedServer(w http.ResponseWriter, r *http.Request) {
	s.h.middleware.RequireAdminAuth(s.h.admin.AddFederatedServer)(w, r)
}

func (s *ServerInterfaceImpl) AddFederatedServerOptions(w http.ResponseWriter, r *http.Request) {
	s.h.admin.AddFederatedServerOptions(w, r)
}

func (s *ServerInterfaceImpl) RemoveFederatedServer(w http.ResponseWriter, r *http.Request, id int) {
	s.h.middleware.RequireAdminAuth(func(w http.ResponseWriter, r *http.Request) {
		s.h.admin.RemoveFederatedServer(w, r, id)
	})(w, r)
}

func (s *ServerInterfaceImpl) RemoveFederatedServerOptions(w http.ResponseWriter, r *http.Request, id int) {
	s.h.admin.RemoveFederatedServerOptions(w, r, id)
}
