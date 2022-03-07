package admin

import (
	"net/http"

	"github.com/owncast/owncast/activitypub"
	"github.com/owncast/owncast/activitypub/outbox"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core/data"
)

// SendFederatedMessage will send a manual message to the fediverse.
func SendFederatedMessage(w http.ResponseWriter, r *http.Request) {
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

	// Update Fediverse followers about this change.
	if err := outbox.UpdateFollowersWithAccountUpdates(); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "federation private saved")
}

// SetFederationShowEngagement will set if Fedivese engagement shows in chat.
func SetFederationShowEngagement(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetFederationShowEngagement(configValue.Value.(bool)); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}
	controllers.WriteSimpleResponse(w, true, "federation show engagement saved")
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

// SetFederationBlockDomains saves a list of domains to block on the Fediverse.
func SetFederationBlockDomains(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValues, success := getValuesFromRequest(w, r)
	if !success {
		controllers.WriteSimpleResponse(w, false, "unable to handle provided domains")
		return
	}

	domainStrings := make([]string, 0)
	for _, domain := range configValues {
		domainStrings = append(domainStrings, domain.Value.(string))
	}

	if err := data.SetBlockedFederatedDomains(domainStrings); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "saved")
}

// GetFederatedActions will return the saved list of accepted inbound
// federated activities.
func GetFederatedActions(page int, pageSize int, w http.ResponseWriter, r *http.Request) {
	offset := pageSize * page

	activities, total, err := persistence.GetInboundActivities(pageSize, offset)
	if err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	response := controllers.PaginatedResponse{
		Total:   total,
		Results: activities,
	}

	controllers.WriteResponse(w, response)
}
