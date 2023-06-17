package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/webserver/requests"
	"github.com/owncast/owncast/webserver/responses"
)

// SetDiscordNotificationConfiguration will set the discord notification configuration.
func (h *Handlers) SetDiscordNotificationConfiguration(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	type request struct {
		Value models.DiscordConfiguration `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var config request
	if err := decoder.Decode(&config); err != nil {
		responses.WriteSimpleResponse(w, false, "unable to update discord config with provided values")
		return
	}

	if err := data.SetDiscordConfig(config.Value); err != nil {
		responses.WriteSimpleResponse(w, false, "unable to update discord config with provided values")
		return
	}

	responses.WriteSimpleResponse(w, true, "updated discord config with provided values")
}

// SetBrowserNotificationConfiguration will set the browser notification configuration.
func (h *Handlers) SetBrowserNotificationConfiguration(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	type request struct {
		Value models.BrowserNotificationConfiguration `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var config request
	if err := decoder.Decode(&config); err != nil {
		responses.WriteSimpleResponse(w, false, "unable to update browser push config with provided values")
		return
	}

	if err := data.SetBrowserPushConfig(config.Value); err != nil {
		responses.WriteSimpleResponse(w, false, "unable to update browser push config with provided values")
		return
	}

	responses.WriteSimpleResponse(w, true, "updated browser push config with provided values")
}
