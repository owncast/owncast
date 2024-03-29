package handler

import (
	"net/http"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/handler/generated"
	"github.com/owncast/owncast/router/middleware"
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
