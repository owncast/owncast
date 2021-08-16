package admin

import (
	"net/http"

	"github.com/owncast/owncast/activitypub"
	"github.com/owncast/owncast/controllers"
)

func SendFederationMessage(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	message, ok := configValue.Value.(string)
	if !ok {
		controllers.WriteSimpleResponse(w, false, "unable to send message")
	}

	activitypub.SendPublicFederatedMessage(message)

	controllers.WriteSimpleResponse(w, true, "sent")
}
