package handler

import (
	"net/http"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/controllers/admin"
	"github.com/owncast/owncast/handler/generated"
	"github.com/owncast/owncast/router/middleware"
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
	controllers.GetStatus(w, r)
}

func (*ServerInterfaceImpl) GetCustomEmojiList(w http.ResponseWriter, r *http.Request) {
	controllers.GetCustomEmojiList(w, r)
}

func (*ServerInterfaceImpl) GetChatMessages(w http.ResponseWriter, r *http.Request, params generated.GetChatMessagesParams) {
	middleware.RequireUserAccessToken(controllers.GetChatMessages)(w, r)
}

func (*ServerInterfaceImpl) RegisterAnonymousChatUser(w http.ResponseWriter, r *http.Request, params generated.RegisterAnonymousChatUserParams) {
	controllers.RegisterAnonymousChatUser(w, r)
}

func (*ServerInterfaceImpl) RegisterAnonymousChatUserOptions(w http.ResponseWriter, r *http.Request) {
	controllers.RegisterAnonymousChatUser(w, r)
}

func (*ServerInterfaceImpl) UpdateMessageVisibility(w http.ResponseWriter, r *http.Request, params generated.UpdateMessageVisibilityParams) {
	middleware.RequireUserModerationScopeAccesstoken(admin.UpdateMessageVisibility)(w, r)
}

func (*ServerInterfaceImpl) UpdateUserEnabled(w http.ResponseWriter, r *http.Request, params generated.UpdateUserEnabledParams) {
	middleware.RequireUserModerationScopeAccesstoken(admin.UpdateUserEnabled)(w, r)
}

func (*ServerInterfaceImpl) GetWebConfig(w http.ResponseWriter, r *http.Request) {
	controllers.GetWebConfig(w, r)
}

func (*ServerInterfaceImpl) GetYPResponse(w http.ResponseWriter, r *http.Request) {
	yp.GetYPResponse(w, r)
}

func (*ServerInterfaceImpl) GetAllSocialPlatforms(w http.ResponseWriter, r *http.Request) {
	controllers.GetAllSocialPlatforms(w, r)
}

func (*ServerInterfaceImpl) GetVideoStreamOutputVariants(w http.ResponseWriter, r *http.Request) {
	controllers.GetVideoStreamOutputVariants(w, r)
}

func (*ServerInterfaceImpl) Ping(w http.ResponseWriter, r *http.Request) {
	controllers.Ping(w, r)
}

func (*ServerInterfaceImpl) RemoteFollow(w http.ResponseWriter, r *http.Request) {
	controllers.RemoteFollow(w, r)
}

func (*ServerInterfaceImpl) GetFollowers(w http.ResponseWriter, r *http.Request, params generated.GetFollowersParams) {
	middleware.HandlePagination(controllers.GetFollowers)(w, r)
}

func (*ServerInterfaceImpl) ReportPlaybackMetrics(w http.ResponseWriter, r *http.Request) {
	controllers.ReportPlaybackMetrics(w, r)
}

func (*ServerInterfaceImpl) RegisterForLiveNotifications(w http.ResponseWriter, r *http.Request, params generated.RegisterForLiveNotificationsParams) {
	middleware.RequireUserAccessToken(controllers.RegisterForLiveNotifications)(w, r)
}
