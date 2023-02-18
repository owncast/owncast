package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/owncast/owncast/internal/activitypub/webfinger"
)

// RemoteFollow handles a request to begin the remote follow redirect flow.
func (s *Service) RemoteFollow(w http.ResponseWriter, r *http.Request) {
	type followRequest struct {
		Account string `json:"account"`
	}

	type followResponse struct {
		RedirectURL string `json:"redirectUrl"`
	}

	var request followRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		s.WriteSimpleResponse(w, false, "unable to parse request")
		return
	}

	if request.Account == "" {
		s.WriteSimpleResponse(w, false, "Remote Fediverse account is required to follow.")
		return
	}

	localActorPath, _ := url.Parse(s.Data.GetServerURL())
	localActorPath.Path = fmt.Sprintf("/federation/user/%s", s.Data.GetDefaultFederationUsername())
	var template string
	links, err := webfinger.GetWebfingerLinks(request.Account)
	if err != nil {
		s.WriteSimpleResponse(w, false, err.Error())
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
		s.WriteSimpleResponse(w, false, "unable to determine remote follow information for "+request.Account)
		return
	}

	redirectURL := strings.Replace(template, "{uri}", localActorPath.String(), 1)
	response := followResponse{
		RedirectURL: redirectURL,
	}

	s.WriteResponse(w, response)
}
