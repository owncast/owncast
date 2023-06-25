package handlers

import (
	"net/http"

	"github.com/owncast/owncast/activitypub"
	"github.com/owncast/owncast/activitypub/outbox"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/webserver/requests"
	"github.com/owncast/owncast/webserver/responses"
)

// SendFederatedMessage will send a manual message to the fediverse.
func (h *Handlers) SendFederatedMessage(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	message, ok := configValue.Value.(string)
	if !ok {
		responses.WriteSimpleResponse(w, false, "unable to send message")
		return
	}

	if err := activitypub.SendPublicFederatedMessage(message); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "sent")
}

// SetFederationEnabled will set if Federation features are enabled.
func (h *Handlers) SetFederationEnabled(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	if err := configRepository.SetFederationEnabled(configValue.Value.(bool)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}
	responses.WriteSimpleResponse(w, true, "federation features saved")
}

// SetFederationActivityPrivate will set if Federation features are private to followers.
func (h *Handlers) SetFederationActivityPrivate(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	if err := configRepository.SetFederationIsPrivate(configValue.Value.(bool)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update Fediverse followers about this change.
	if err := outbox.UpdateFollowersWithAccountUpdates(); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "federation private saved")
}

// SetFederationShowEngagement will set if Fedivese engagement shows in chat.
func (h *Handlers) SetFederationShowEngagement(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	if err := configRepository.SetFederationShowEngagement(configValue.Value.(bool)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}
	responses.WriteSimpleResponse(w, true, "federation show engagement saved")
}

// SetFederationUsername will set the local actor username used for federation activities.
func (h *Handlers) SetFederationUsername(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	if err := configRepository.SetFederationUsername(configValue.Value.(string)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "username saved")
}

// SetFederationGoLiveMessage will set the federated message sent when the streamer goes live.
func (h *Handlers) SetFederationGoLiveMessage(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	if err := configRepository.SetFederationGoLiveMessage(configValue.Value.(string)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "message saved")
}

// SetFederationBlockDomains saves a list of domains to block on the Fediverse.
func (h *Handlers) SetFederationBlockDomains(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValues, success := requests.GetValuesFromRequest(w, r)
	if !success {
		responses.WriteSimpleResponse(w, false, "unable to handle provided domains")
		return
	}

	domainStrings := make([]string, 0)
	for _, domain := range configValues {
		domainStrings = append(domainStrings, domain.Value.(string))
	}

	if err := configRepository.SetBlockedFederatedDomains(domainStrings); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "saved")
}

// GetFederatedActions will return the saved list of accepted inbound
// federated activities.
func (h *Handlers) GetFederatedActions(page int, pageSize int, w http.ResponseWriter, r *http.Request) {
	offset := pageSize * page

	activities, total, err := persistence.GetInboundActivities(pageSize, offset)
	if err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	response := responses.PaginatedResponse{
		Total:   total,
		Results: activities,
	}

	responses.WriteResponse(w, response)
}
