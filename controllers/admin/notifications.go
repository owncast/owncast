package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/models"
)

// SetDiscordNotificationConfiguration will set the discord notification configuration.
func (c *Controller) SetDiscordNotificationConfiguration(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	type request struct {
		Value models.DiscordConfiguration `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var config request
	if err := decoder.Decode(&config); err != nil {
		c.Service.WriteSimpleResponse(w, false, "unable to update discord config with provided values")
		return
	}

	if err := c.Data.SetDiscordConfig(config.Value); err != nil {
		c.Service.WriteSimpleResponse(w, false, "unable to update discord config with provided values")
		return
	}

	c.Service.WriteSimpleResponse(w, true, "updated discord config with provided values")
}

// SetBrowserNotificationConfiguration will set the browser notification configuration.
func (c *Controller) SetBrowserNotificationConfiguration(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	type request struct {
		Value models.BrowserNotificationConfiguration `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var config request
	if err := decoder.Decode(&config); err != nil {
		c.Service.WriteSimpleResponse(w, false, "unable to update browser push config with provided values")
		return
	}

	if err := c.Data.SetBrowserPushConfig(config.Value); err != nil {
		c.Service.WriteSimpleResponse(w, false, "unable to update browser push config with provided values")
		return
	}

	c.Service.WriteSimpleResponse(w, true, "updated browser push config with provided values")
}
