package handler

import (
	"net/http"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/handler/generated"
)

type ServerInterfaceImpl struct {}

func New() *ServerInterfaceImpl {
	return &ServerInterfaceImpl{}
}

func (s *ServerInterfaceImpl) Handler() http.Handler {
	return generated.Handler(s)
}

func (s *ServerInterfaceImpl) Status(w http.ResponseWriter, r *http.Request) {
	controllers.GetStatus(w, r)
}
