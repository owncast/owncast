package handler

import (
	"net/http"

	"github.com/owncast/owncast/controllers/admin"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/router/middleware"
	"github.com/owncast/owncast/utils"
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

func (*ServerInterfaceImpl) SendUserMessage(w http.ResponseWriter, r *http.Request) {
	middleware.RequireExternalAPIAccessToken(user.ScopeCanSendChatMessages, admin.SendUserMessage)(w, r)
}
