package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/owncast/owncast/core/data"
)

// RemoteFollow handles a request to begin the remote follow redirect flow.
func RemoteFollow(w http.ResponseWriter, r *http.Request) {
	type followRequest struct {
		Account string `json:"account"`
	}

	type followResponse struct {
		RedirectURL string `json:"redirectUrl"`
	}

	var request followRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		WriteSimpleResponse(w, false, "unable to parse request")
		return
	}

	if request.Account == "" {
		WriteSimpleResponse(w, false, "Remote Fediverse account is required to follow.")
		return
	}

	localActorPath, _ := url.Parse(data.GetServerURL())
	localActorPath.Path = fmt.Sprintf("/federation/user/%s", data.GetDefaultFederationUsername())
	var template string
	links, err := getWebfingerLinks(request.Account)
	if err != nil {
		WriteSimpleResponse(w, false, err.Error())
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
		WriteSimpleResponse(w, false, "unable to determine remote follow information for "+request.Account)
		return
	}

	redirectURL := strings.Replace(template, "{uri}", localActorPath.String(), 1)
	response := followResponse{
		RedirectURL: redirectURL,
	}

	WriteResponse(w, response)
}

func getWebfingerLinks(account string) ([]map[string]interface{}, error) {
	type webfingerResponse struct {
		Links []map[string]interface{} `json:"links"`
	}

	account = strings.TrimLeft(account, "@") // remove any leading @
	accountComponents := strings.Split(account, "@")
	fediverseServer := accountComponents[1]

	// HTTPS is required.
	requestURL, err := url.Parse("https://" + fediverseServer)
	if err != nil {
		return nil, fmt.Errorf("unable to parse fediverse server host %s", fediverseServer)
	}

	requestURL.Path = "/.well-known/webfinger"
	query := requestURL.Query()
	query.Add("resource", fmt.Sprintf("acct:%s", account))
	requestURL.RawQuery = query.Encode()

	response, err := http.DefaultClient.Get(requestURL.String())
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var links webfingerResponse
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&links); err != nil {
		return nil, err
	}

	return links.Links, nil
}
