package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/owncast/owncast/services/activitypub/webfinger"
	"github.com/owncast/owncast/webserver/handlers/generated"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// RemoteFollow handles a request to begin the remote follow redirect flow.
func (h *Handlers) RemoteFollow(w http.ResponseWriter, r *http.Request) {
	type followResponse struct {
		RedirectURL string `json:"redirectUrl"`
	}

	var request generated.RemoteFollowJSONRequestBody
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		webutils.WriteSimpleResponse(w, false, "unable to parse request")
		return
	}

	if request.Account == nil {
		webutils.WriteSimpleResponse(w, false, "account field is required")
		return
	}

	if *request.Account == "" {
		webutils.WriteSimpleResponse(w, false, "Remote Fediverse account is required to follow.")
		return
	}

	localActorPath := h.apBuilder.MakeLocalIRIForAccount(h.configRepository.GetDefaultFederationUsername())
	var template string
	links, err := webfinger.GetWebfingerLinks(*request.Account)
	if err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Acquire the remote follow redirect template.
	for _, link := range links {
		for k, v := range link {
			if k == "rel" && v == "http://ostatus.org/schema/1.0/subscribe" && link["template"] != nil {
				template = link["template"].(string)
			}
		}
	}

	if localActorPath == nil || localActorPath.String() == "" || template == "" {
		webutils.WriteSimpleResponse(w, false, "unable to determine remote follow information for "+*request.Account)
		return
	}

	redirectURL := strings.Replace(template, "{uri}", localActorPath.String(), 1)
	response := followResponse{
		RedirectURL: redirectURL,
	}

	webutils.WriteResponse(w, response)
}
