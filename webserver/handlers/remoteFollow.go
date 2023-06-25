package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/owncast/owncast/activitypub/webfinger"
	"github.com/owncast/owncast/webserver/responses"
)

// RemoteFollow handles a request to begin the remote follow redirect flow.
func (h *Handlers) RemoteFollow(w http.ResponseWriter, r *http.Request) {
	type followRequest struct {
		Account string `json:"account"`
	}

	type followResponse struct {
		RedirectURL string `json:"redirectUrl"`
	}

	var request followRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		responses.WriteSimpleResponse(w, false, "unable to parse request")
		return
	}

	if request.Account == "" {
		responses.WriteSimpleResponse(w, false, "Remote Fediverse account is required to follow.")
		return
	}

	localActorPath, _ := url.Parse(configRepository.GetServerURL())
	localActorPath.Path = fmt.Sprintf("/federation/user/%s", configRepository.GetDefaultFederationUsername())
	var template string
	links, err := webfinger.GetWebfingerLinks(request.Account)
	if err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
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

	if localActorPath.String() == "" || template == "" {
		responses.WriteSimpleResponse(w, false, "unable to determine remote follow information for "+request.Account)
		return
	}

	redirectURL := strings.Replace(template, "{uri}", localActorPath.String(), 1)
	response := followResponse{
		RedirectURL: redirectURL,
	}

	responses.WriteResponse(w, response)
}
