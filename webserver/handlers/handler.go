package handlers

import (
	"net/http"

	"github.com/owncast/owncast/webserver/handlers/admin"
	"github.com/owncast/owncast/webserver/handlers/generated"
	"github.com/owncast/owncast/webserver/router/middleware"

	"github.com/owncast/owncast/yp"
)

type ServerInterfaceImpl struct{}

// ensure ServerInterfaceImpl implements ServerInterface.
var _ generated.ServerInterface = &ServerInterfaceImpl{}

func New() *ServerInterfaceImpl {
	return &ServerInterfaceImpl{}
}

func (s *ServerInterfaceImpl) Handler() http.Handler {
	return generated.Handler(s)
}

func (*ServerInterfaceImpl) GetStatus(w http.ResponseWriter, r *http.Request) {
	GetStatus(w, r)
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

func (*ServerInterfaceImpl) UpdateMessageVisibility(w http.ResponseWriter, r *http.Request, params generated.UpdateMessageVisibilityParams) {
	middleware.RequireUserModerationScopeAccesstoken(admin.UpdateMessageVisibility)(w, r)
}

func (*ServerInterfaceImpl) UpdateUserEnabled(w http.ResponseWriter, r *http.Request, params generated.UpdateUserEnabledParams) {
	middleware.RequireUserModerationScopeAccesstoken(admin.UpdateUserEnabled)(w, r)
}

func (*ServerInterfaceImpl) GetWebConfig(w http.ResponseWriter, r *http.Request) {
	GetWebConfig(w, r)
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

func (*ServerInterfaceImpl) Ping(w http.ResponseWriter, r *http.Request) {
	Ping(w, r)
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

// Federated servers endpoints

func (*ServerInterfaceImpl) GetFederatedServers(w http.ResponseWriter, r *http.Request) {
	admin.GetFederatedServers(w, r)
}

func (*ServerInterfaceImpl) AddFederatedServer(w http.ResponseWriter, r *http.Request) {
	middleware.RequireAdminAuth(admin.AddFederatedServer)(w, r)
}

func (*ServerInterfaceImpl) AddFederatedServerOptions(w http.ResponseWriter, r *http.Request) {
	admin.AddFederatedServerOptions(w, r)
}

func (*ServerInterfaceImpl) RemoveFederatedServer(w http.ResponseWriter, r *http.Request, id int) {
	middleware.RequireAdminAuth(func(w http.ResponseWriter, r *http.Request) {
		admin.RemoveFederatedServer(w, r, id)
	})(w, r)
}

func (*ServerInterfaceImpl) RemoveFederatedServerOptions(w http.ResponseWriter, r *http.Request, id int) {
	admin.RemoveFederatedServerOptions(w, r, id)
}
