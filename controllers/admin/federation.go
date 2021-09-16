package admin

import (
	"net/http"

	"github.com/owncast/owncast/activitypub"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core/data"
)

// SendFederationMessage will send a manual message to the fediverse.
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

	if err := activitypub.SendPublicFederatedMessage(message); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "sent")
}

// SetFederationEnabled will set if Federation features are enabled.
func SetFederationEnabled(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetFederationEnabled(configValue.Value.(bool)); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}
	controllers.WriteSimpleResponse(w, true, "federation features saved")
}

// SetFederationActivityPrivate will set if Federation features are private to followers.
func SetFederationActivityPrivate(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetFederationIsPrivate(configValue.Value.(bool)); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}
	controllers.WriteSimpleResponse(w, true, "federation private saved")
}

// SetFederationUsername will set the local actor username used for federation activities.
func SetFederationUsername(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetFederationUsername(configValue.Value.(string)); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "username saved")
}

// SetFederationGoLiveMessage will set the federated message sent when the streamer goes live.
func SetFederationGoLiveMessage(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetFederationGoLiveMessage(configValue.Value.(string)); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "message saved")
}

// ApproveFollower will approve a federated follow request.
func ApproveFollower(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	type approveFollowerRequest struct {
		FederationIRI string `json:"federationIRI"`
		Approved      bool   `json:"approved"`
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	approval := configValue.Value.(approveFollowerRequest)
	if err := persistence.ApprovePreviousFollowRequest(approval.FederationIRI); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}
	controllers.WriteSimpleResponse(w, true, "follow request approved")
}

// GetPendingFollowRequests will return a list of pending follow requests.
func GetPendingFollowRequests(w http.ResponseWriter, r *http.Request) {
	//requests := activitypub.GetPendingFollowRequests()

}
