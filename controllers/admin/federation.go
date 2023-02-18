package admin

import (
	"net/http"

	"github.com/owncast/owncast/controllers"
)

// SendFederatedMessage will send a manual message to the fediverse.
func (c *Controller) SendFederatedMessage(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	message, ok := configValue.Value.(string)
	if !ok {
		c.WriteSimpleResponse(w, false, "unable to send message")
		return
	}

	if err := c.ActivityPub.SendPublicFederatedMessage(message); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "sent")
}

// SetFederationEnabled will set if Federation features are enabled.
func (c *Controller) SetFederationEnabled(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := c.Data.SetFederationEnabled(configValue.Value.(bool)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}
	c.WriteSimpleResponse(w, true, "federation features saved")
}

// SetFederationActivityPrivate will set if Federation features are private to followers.
func (c *Controller) SetFederationActivityPrivate(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := c.Data.SetFederationIsPrivate(configValue.Value.(bool)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update Fediverse followers about this change.
	if err := c.Outbox.UpdateFollowersWithAccountUpdates(c.ActivityPub.Persistence); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "federation private saved")
}

// SetFederationShowEngagement will set if Fedivese engagement shows in chat.
func (c *Controller) SetFederationShowEngagement(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := c.Data.SetFederationShowEngagement(configValue.Value.(bool)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}
	c.WriteSimpleResponse(w, true, "federation show engagement saved")
}

// SetFederationUsername will set the local actor username used for federation activities.
func (c *Controller) SetFederationUsername(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := c.Data.SetFederationUsername(configValue.Value.(string)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "username saved")
}

// SetFederationGoLiveMessage will set the federated message sent when the streamer goes live.
func (c *Controller) SetFederationGoLiveMessage(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := c.Data.SetFederationGoLiveMessage(configValue.Value.(string)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "message saved")
}

// SetFederationBlockDomains saves a list of domains to block on the Fediverse.
func (c *Controller) SetFederationBlockDomains(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValues, success := c.getValuesFromRequest(w, r)
	if !success {
		c.WriteSimpleResponse(w, false, "unable to handle provided domains")
		return
	}

	domainStrings := make([]string, 0)
	for _, domain := range configValues {
		domainStrings = append(domainStrings, domain.Value.(string))
	}

	if err := c.Data.SetBlockedFederatedDomains(domainStrings); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "saved")
}

// GetFederatedActions will return the saved list of accepted inbound
// federated activities.
func (c *Controller) GetFederatedActions(page int, pageSize int, w http.ResponseWriter, r *http.Request) {
	offset := pageSize * page

	activities, total, err := c.ActivityPub.Persistence.GetInboundActivities(pageSize, offset)
	if err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	response := controllers.PaginatedResponse{
		Total:   total,
		Results: activities,
	}

	c.WriteResponse(w, response)
}
