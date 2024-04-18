package handler

import (
	"net/http"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/controllers/admin"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/router/middleware"
	"github.com/owncast/owncast/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (*ServerInterfaceImpl) SendSystemMessage(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(user.ScopeCanSendSystemMessages, admin.SendSystemMessage)(w, r)
}

func (*ServerInterfaceImpl) SendSystemMessageToConnectedClient(w http.ResponseWriter, r *http.Request, clientId int) {
	// TODO doing this hack to make the new system work with the old system
	// TODO this needs to be refactored asap
	r.Header[utils.RestURLPatternHeaderKey] = []string{`/api/integrations/chat/system/client/{clientId}`}
	middleware.RequireExternalAPIAccessToken(user.ScopeCanSendSystemMessages, admin.SendSystemMessageToConnectedClient)(w, r)
}

// Deprecated
func (*ServerInterfaceImpl) SendUserMessage(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(user.ScopeCanSendChatMessages, admin.SendUserMessage)(w, r)
}

func (*ServerInterfaceImpl) SendIntegrationChatMessage(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(user.ScopeCanSendChatMessages, admin.SendIntegrationChatMessage)(w, r)
}

func (*ServerInterfaceImpl) SendChatAction(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(user.ScopeCanSendSystemMessages, admin.SendChatAction)(w, r)
}

func (*ServerInterfaceImpl) ExternalUpdateMessageVisibility(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(user.ScopeHasAdminAccess, admin.ExternalUpdateMessageVisibility)(w, r)
}

func (*ServerInterfaceImpl) ExternalSetStreamTitle(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(user.ScopeHasAdminAccess, admin.ExternalSetStreamTitle)(w, r)
}

func (*ServerInterfaceImpl) ExternalGetChatMessages(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(user.ScopeHasAdminAccess, controllers.ExternalGetChatMessages)(w, r)
}

func (*ServerInterfaceImpl) ExternalGetConnectedChatClients(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(user.ScopeHasAdminAccess, admin.ExternalGetConnectedChatClients)(w, r)
}

func (*ServerInterfaceImpl) GetPrometheusAPI(w http.ResponseWriter, r *http.Request) {
	// TODO might need to bring this out of the codegen
	middleware.RequireAdminAuth(func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler()
	})(w, r)
}

func (*ServerInterfaceImpl) PostPrometheusAPI(w http.ResponseWriter, r *http.Request) {
	// TODO might need to bring this out of the codegen
	middleware.RequireAdminAuth(func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler()
	})(w, r)
}

func (*ServerInterfaceImpl) PutPrometheusAPI(w http.ResponseWriter, r *http.Request) {
	// TODO might need to bring this out of the codegen
	middleware.RequireAdminAuth(func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler()
	})(w, r)
}

func (*ServerInterfaceImpl) DeletePrometheusAPI(w http.ResponseWriter, r *http.Request) {
	// TODO might need to bring this out of the codegen
	middleware.RequireAdminAuth(func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler()
	})(w, r)
}
