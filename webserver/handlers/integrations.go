package handlers

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/webserver/router/middleware"
)

func (s *ServerInterfaceImpl) SendSystemMessage(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(models.ScopeCanSendSystemMessages, s.h.admin.SendSystemMessage)(w, r)
}

func (s *ServerInterfaceImpl) SendSystemMessageOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(models.ScopeCanSendSystemMessages, s.h.admin.SendSystemMessage)(w, r)
}

func (s *ServerInterfaceImpl) SendSystemMessageToConnectedClient(w http.ResponseWriter, r *http.Request, clientId int) {
	middleware.RequireExternalAPIAccessToken(models.ScopeCanSendSystemMessages, s.h.admin.SendSystemMessageToConnectedClient)(w, r)
}

func (s *ServerInterfaceImpl) SendSystemMessageToConnectedClientOptions(w http.ResponseWriter, r *http.Request, clientId int) {
	middleware.RequireExternalAPIAccessToken(models.ScopeCanSendSystemMessages, s.h.admin.SendSystemMessageToConnectedClient)(w, r)
}

// Deprecated.
func (s *ServerInterfaceImpl) SendUserMessage(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(models.ScopeCanSendChatMessages, s.h.admin.SendUserMessage)(w, r)
}

// Deprecated.
func (s *ServerInterfaceImpl) SendUserMessageOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(models.ScopeCanSendChatMessages, s.h.admin.SendUserMessage)(w, r)
}

func (s *ServerInterfaceImpl) SendIntegrationChatMessage(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(models.ScopeCanSendChatMessages, s.h.admin.SendIntegrationChatMessage)(w, r)
}

func (s *ServerInterfaceImpl) SendIntegrationChatMessageOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(models.ScopeCanSendChatMessages, s.h.admin.SendIntegrationChatMessage)(w, r)
}

func (s *ServerInterfaceImpl) SendChatAction(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(models.ScopeCanSendSystemMessages, s.h.admin.SendChatAction)(w, r)
}

func (s *ServerInterfaceImpl) SendChatActionOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(models.ScopeCanSendSystemMessages, s.h.admin.SendChatAction)(w, r)
}

func (s *ServerInterfaceImpl) ExternalUpdateMessageVisibility(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(models.ScopeHasAdminAccess, s.h.admin.ExternalUpdateMessageVisibility)(w, r)
}

func (s *ServerInterfaceImpl) ExternalUpdateMessageVisibilityOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(models.ScopeHasAdminAccess, s.h.admin.ExternalUpdateMessageVisibility)(w, r)
}

func (s *ServerInterfaceImpl) ExternalSetStreamTitle(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(models.ScopeHasAdminAccess, s.h.admin.ExternalSetStreamTitle)(w, r)
}

func (s *ServerInterfaceImpl) ExternalSetStreamTitleOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(models.ScopeHasAdminAccess, s.h.admin.ExternalSetStreamTitle)(w, r)
}

func (s *ServerInterfaceImpl) ExternalGetUserDetails(w http.ResponseWriter, r *http.Request, userId string) {
	middleware.RequireExternalAPIAccessToken(models.ScopeHasAdminAccess, s.h.moderation.ExternalGetUserDetails)(w, r)
}

func (*ServerInterfaceImpl) ExternalGetChatMessages(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(models.ScopeHasAdminAccess, ExternalGetChatMessages)(w, r)
}

func (s *ServerInterfaceImpl) ExternalGetChatMessagesOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(models.ScopeHasAdminAccess, ExternalGetChatMessages)(w, r)
}

func (s *ServerInterfaceImpl) ExternalGetConnectedChatClients(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(models.ScopeHasAdminAccess, s.h.admin.ExternalGetConnectedChatClients)(w, r)
}

func (s *ServerInterfaceImpl) ExternalGetConnectedChatClientsOptions(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(models.ScopeHasAdminAccess, s.h.admin.ExternalGetConnectedChatClients)(w, r)
}

func (s *ServerInterfaceImpl) ExternalGetStatus(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(models.ScopeHasAdminAccess, s.h.admin.ExternalGetStatus)(w, r)
}

func (*ServerInterfaceImpl) GetPrometheusAPI(w http.ResponseWriter, r *http.Request) {
	// might need to bring this out of the codegen
	middleware.RequireAdminAuth(func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})(w, r)
}

func (*ServerInterfaceImpl) PostPrometheusAPI(w http.ResponseWriter, r *http.Request) {
	// might need to bring this out of the codegen
	middleware.RequireAdminAuth(func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})(w, r)
}

func (*ServerInterfaceImpl) PutPrometheusAPI(w http.ResponseWriter, r *http.Request) {
	// might need to bring this out of the codegen
	middleware.RequireAdminAuth(func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})(w, r)
}

func (*ServerInterfaceImpl) DeletePrometheusAPI(w http.ResponseWriter, r *http.Request) {
	// might need to bring this out of the codegen
	middleware.RequireAdminAuth(func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})(w, r)
}

func (*ServerInterfaceImpl) OptionsPrometheusAPI(w http.ResponseWriter, r *http.Request) {
	// might need to bring this out of the codegen
	middleware.RequireAdminAuth(func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})(w, r)
}
