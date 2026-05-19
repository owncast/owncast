package handlers

import (
	"net/http"

	"github.com/owncast/owncast/webserver/handlers/generated"
	"github.com/owncast/owncast/webserver/router/middleware"
)

func (s *ServerInterfaceImpl) StartIndieAuthFlow(w http.ResponseWriter, r *http.Request, params generated.StartIndieAuthFlowParams) {
	middleware.RequireUserAccessToken(s.h.indieauth.StartAuthFlow)(w, r)
}

func (s *ServerInterfaceImpl) HandleIndieAuthRedirect(w http.ResponseWriter, r *http.Request, params generated.HandleIndieAuthRedirectParams) {
	s.h.indieauth.HandleRedirect(w, r)
}

func (s *ServerInterfaceImpl) HandleIndieAuthEndpointGet(w http.ResponseWriter, r *http.Request, params generated.HandleIndieAuthEndpointGetParams) {
	middleware.RequireAdminAuth(s.h.indieauth.HandleAuthEndpointGet)(w, r)
}

func (s *ServerInterfaceImpl) HandleIndieAuthEndpointPost(w http.ResponseWriter, r *http.Request) {
	s.h.indieauth.HandleAuthEndpointPost(w, r)
}

func (s *ServerInterfaceImpl) RegisterFediverseOTPRequest(w http.ResponseWriter, r *http.Request, params generated.RegisterFediverseOTPRequestParams) {
	middleware.RequireUserAccessToken(s.h.fediverse.RegisterFediverseOTPRequest)(w, r)
}

func (s *ServerInterfaceImpl) VerifyFediverseOTPRequest(w http.ResponseWriter, r *http.Request) {
	s.h.fediverse.VerifyFediverseOTPRequest(w, r)
}
