package handler

import (
	"net/http"

	"github.com/owncast/owncast/controllers/auth/fediverse"
	"github.com/owncast/owncast/controllers/auth/indieauth"
	"github.com/owncast/owncast/handler/generated"
	"github.com/owncast/owncast/router/middleware"
)

func (*ServerInterfaceImpl) StartIndieAuthFlow(w http.ResponseWriter, r *http.Request, params generated.StartIndieAuthFlowParams) {
	middleware.RequireUserAccessToken(indieauth.StartAuthFlow)(w, r)
}

func (*ServerInterfaceImpl) HandleIndieAuthRedirect(w http.ResponseWriter, r *http.Request, params generated.HandleIndieAuthRedirectParams) {
	indieauth.HandleRedirect(w, r)
}

func (*ServerInterfaceImpl) HandleIndieAuthEndpointGet(w http.ResponseWriter, r *http.Request, params generated.HandleIndieAuthEndpointGetParams) {
	middleware.RequireAdminAuth(indieauth.HandleAuthEndpointGet)(w, r)
}

func (*ServerInterfaceImpl) HandleIndieAuthEndpointPost(w http.ResponseWriter, r *http.Request) {
	indieauth.HandleAuthEndpointPost(w, r)
}

func (*ServerInterfaceImpl) RegisterFediverseOTPRequest(w http.ResponseWriter, r *http.Request, params generated.RegisterFediverseOTPRequestParams) {
	middleware.RequireUserAccessToken(fediverse.RegisterFediverseOTPRequest)(w, r)
}

func (*ServerInterfaceImpl) VerifyFediverseOTPRequest(w http.ResponseWriter, r *http.Request) {
	fediverse.VerifyFediverseOTPRequest(w, r)
}
