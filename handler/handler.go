package handler

import (
	"net/http"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/handler/generated"
	"github.com/owncast/owncast/router/middleware"
	"github.com/owncast/owncast/yp"
)

type ServerInterfaceImpl struct {}

func New() *ServerInterfaceImpl {
	return &ServerInterfaceImpl{}
}

func (s *ServerInterfaceImpl) Handler() http.Handler {
	return generated.Handler(s)
}

func (*ServerInterfaceImpl) GetStatus(w http.ResponseWriter, r *http.Request) {
	controllers.GetStatus(w, r)
}

func (*ServerInterfaceImpl) GetEmoji(w http.ResponseWriter, r *http.Request) {
	controllers.GetCustomEmojiList(w, r)
}

func (*ServerInterfaceImpl) GetChatList(w http.ResponseWriter, r *http.Request) {
	middleware.RequireUserAccessToken(controllers.GetChatMessages)(w, r)
}

func (*ServerInterfaceImpl) OptionsChatRegister(w http.ResponseWriter, r *http.Request) {
	// TODO this uses the same method as POST, should break that method up
	controllers.RegisterAnonymousChatUser(w, r)
}

func (*ServerInterfaceImpl) RegisterAnonymousChatUser(w http.ResponseWriter, r *http.Request, params generated.RegisterAnonymousChatUserParams) {
	controllers.RegisterAnonymousChatUser(w, r)
}

func (*ServerInterfaceImpl) GetConfig(w http.ResponseWriter, r *http.Request) {
	controllers.GetWebConfig(w, r)
}

func (*ServerInterfaceImpl) GetYP(w http.ResponseWriter, r *http.Request) {
	yp.GetYPResponse(w, r)
}

func (*ServerInterfaceImpl) GetSocialPlatforms(w http.ResponseWriter, r *http.Request) {
	controllers.GetAllSocialPlatforms(w, r)
}

func (*ServerInterfaceImpl) GetVideoVariants(w http.ResponseWriter, r *http.Request) {
	controllers.GetVideoStreamOutputVariants(w, r)
}

func (*ServerInterfaceImpl) Ping(w http.ResponseWriter, r *http.Request) {
	controllers.Ping(w, r)
}

func (*ServerInterfaceImpl) RemoteFollow(w http.ResponseWriter, r *http.Request) {
	controllers.RemoteFollow(w, r)
}

func (*ServerInterfaceImpl) GetFollowers(w http.ResponseWriter, r *http.Request, params generated.GetFollowersParams) {
	// TODO instead of using the mw, we can use the params object
	middleware.HandlePagination(controllers.GetFollowers)(w, r)
}
